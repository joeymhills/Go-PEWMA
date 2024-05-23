package goggins

import (
    "testing"
)

func TestYamlParse(t *testing.T) {
    s, err := Init("config.yaml")
    if err != nil {
	t.Errorf("error parsing config from yaml file: %s", err)
    }
    err = s.StartServer()
    if err != nil {
	t.Error(err)
    }
}
