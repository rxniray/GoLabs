package serializer 

import (
	"encoding/json"
	"gopkg.in/yaml.v3"
)

type Server struct {
	Host       string   `json:"host" yaml:"host"`
	Port       int      `json:"port" yaml:"port"`
	Debug      bool     `json:"debug" yaml:"debug"`
	AllowedIPs []string `json:"allowed_ips" yaml:"allowed_ips"`
}

func ToYAML(v any) (string, error) {
	bytes, err := yaml.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func ToJSON(v any) (string, error) {
	bytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}