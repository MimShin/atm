package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

type Config struct {
	Values     map[string]interface{}
	configPath string
}

func NewConfig() *Config {
	c := Config{}
	return &c
}

func (c *Config) ReadConfigFile(path string) (err error) {
	jsonFile, err := os.Open(path)
	if err != nil {
		log.Errorf("unable to open config file: path=%s, error=%v", path, err)
		return
	}
	defer jsonFile.Close()

	confBytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Errorf("unable to read config file: path=%s, error=%v", path, err)
		return
	}

	log.Debugf("%s", confBytes)
	if err = json.Unmarshal(confBytes, &c.Values); err != nil {
		log.Errorf("unable to unmarshal config file: path=%s, error=%v", path, err)
		return
	}

	c.configPath = path
	return
}

func (c *Config) StrValue(key string, defaultValue string) string {
	if v, exists := c.Values[key]; exists {
		switch t := v.(type) {
		case string:
			return t
		default:
			log.Warnf("%s is not string: type=%T, value=%v", key, t, t)
		}
	}
	return defaultValue
}

func LogLevel(logLevel string) log.Level {
	switch strings.ToLower(logLevel) {
	case "trace":
		return log.TraceLevel
	case "debug":
		return log.DebugLevel
	case "info":
		return log.InfoLevel
	case "warn":
		return log.WarnLevel
	case "error":
		return log.ErrorLevel
	case "fatal":
		return log.FatalLevel
	case "panic":
		return log.PanicLevel
	}
	return log.InfoLevel
}
