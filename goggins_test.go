package goggins

import (
    "testing"
)

func TestYaml(t *testing.T) {
    _, err := Init("config.yaml")
    if err != nil {
	t.Errorf("error parsing config from yaml file: %s", err)
    }
}
