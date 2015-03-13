package env_json

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

const (
	ENV_JSON_ENV_NAME = "ENV_JSON_CONFIG"
	ENV_JSON_ENV_EXT  = ".env"
)

type EnvJson struct {
	envName string
	envExt  string
}

func NewEnvJson(envName string, envExt string) *EnvJson {
	if envName == "" {
		panic("env_json: env name could not be nil")
	}

	return &EnvJson{
		envName: envName,
		envExt:  envExt,
	}
}

func (p *EnvJson) Marshal(v interface{}) (data []byte, err error) {
	return json.Marshal(v)
}

func (p *EnvJson) MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	return json.MarshalIndent(v, prefix, indent)
}

func (p *EnvJson) Unmarshal(data []byte, v interface{}) (err error) {
	strConfigFiles := os.Getenv(p.envName)

	configFiles := strings.Split(strConfigFiles, ";")

	files := []string{}

	if strConfigFiles == "" || len(files) == 0 {
		return json.Unmarshal(data, v)
	}

	for _, confFile := range configFiles {
		var fi os.FileInfo
		if fi, err = os.Stat(confFile); err != nil {
			return
		}

		if fi.IsDir() {
			var dir *os.File
			if dir, err = os.Open(confFile); err != nil {
				return
			}

			var names []string
			if names, err = dir.Readdirnames(-1); err != nil {
				return
			}

			for _, name := range names {
				if ext := filepath.Ext(name); ext == p.envExt {
					filePath := strings.TrimRight(confFile, "/")
					files = append(files, filePath+"/"+name)
				}
			}
		} else {
			if ext := filepath.Ext(confFile); ext == p.envExt {
				files = append(files, confFile)
			}
		}
	}

	envs := map[string]map[string]interface{}{}

	for _, file := range files {
		var data []byte
		if data, err = ioutil.ReadFile(file); err != nil {

			return
		}

		env := map[string]interface{}{}
		if err = json.Unmarshal(data, &env); err != nil {
			return
		}

		envs[file] = env
	}

	allEnvs := map[string]interface{}{}

	for file, env := range envs {
		for envKey, envVal := range env {
			if _, exist := allEnvs[envKey]; exist {
				err = fmt.Errorf("env key of %s already exist, env file: %s", envKey, file)
				return
			} else {
				allEnvs[envKey] = envVal
			}
		}
	}

	var tpl *template.Template

	if tpl, err = template.New("env_json").Parse(string(data)); err != nil {
		return
	}
	var buf bytes.Buffer
	if err = tpl.Execute(&buf, allEnvs); err != nil {
		return
	}

	strData := buf.String()

	if strings.Contains(strData, "<no value>") {
		err = fmt.Errorf("some env value did not exist")
		return
	}

	err = json.Unmarshal([]byte(strData), v)

	return
}

func Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	return json.MarshalIndent(v, prefix, indent)
}

func Unmarshal(data []byte, v interface{}) error {
	envJson := NewEnvJson(ENV_JSON_ENV_NAME, ENV_JSON_ENV_EXT)
	return envJson.Unmarshal(data, v)
}
