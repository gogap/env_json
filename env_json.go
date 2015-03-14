package env_json

import (
	"encoding/json"
	"github.com/gogap/env_strings"
)

const (
	ENV_JSON_KEY = "ENV_JSON_CONFIG"
	ENV_JSON_EXT = ".env"
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
	envStrings := env_strings.NewEnvStrings(p.envName, p.envExt)

	strData := ""
	if strData, err = envStrings.Execute(string(data)); err != nil {
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
	envJson := NewEnvJson(ENV_JSON_KEY, ENV_JSON_EXT)
	return envJson.Unmarshal(data, v)
}
