package config

import (
	"errors"
	"os"
)

type Conf struct {
	Authentik struct {
		Scheme string
		Host   string
		Token  string
	}
	Outline struct {
		URL   string
		Token string
	}
	App struct {
		GroupPrefix       string
		GroupNameSelector string
	}
}

func Init() (*Conf, error) {
	conf := &Conf{}
	// Authentik
	getString("AUTHENTIK_SCHEME", "http", &conf.Authentik.Scheme)
	if ok := getString("AUTHENTIK_HOST", "", &conf.Authentik.Host); !ok {
		return nil, errors.New("authentik.host is required")
	}
	if ok := getString("AUTHENTIK_TOKEN", "", &conf.Authentik.Token); !ok {
		return nil, errors.New("authentik.token is required")
	}

	// Outline
	if ok := getString("OUTLINE_URL", "", &conf.Outline.URL); !ok {
		return nil, errors.New("outline.url is required")
	}
	if ok := getString("OUTLINE_TOKEN", "", &conf.Outline.Token); !ok {
		return nil, errors.New("outline.token is required")
	}

	// App
	getString("GROUP_PREFIX", "outline_", &conf.App.GroupPrefix)
	getString("GROUP_NAME_SELECTOR", "name", &conf.App.GroupNameSelector)

	return conf, nil
}

func getString(key, defaultValue string, val *string) bool {
	value := os.Getenv(key)
	if value == "" {
		*val = defaultValue
		return false
	}
	*val = value
	return true
}
