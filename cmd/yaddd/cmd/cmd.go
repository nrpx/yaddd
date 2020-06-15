package cmd

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"path"
	"strings"

	"yaddd/internal/config"
	"yaddd/internal/core"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// Путь к файлу конфигурации по умолчанию.
const defaultConfigPath = "/etc/yaddd/config.yml"

// Неизвестное расширение файла.
var errUnknownFileExt = errors.New("unknown file extension")

// Основная команда.
type cmd struct{ conf config.Config }

// Загрузка файла конфигурации.
func (c *cmd) loadConfig(configFile string) (err error) {
	var confData []byte

	fileExt := strings.ToLower(path.Ext(configFile))

	var unmarshal func([]byte, interface{}) error

	switch fileExt {
	case ".yml", ".yaml":
		unmarshal = yaml.Unmarshal
	case ".json":
		unmarshal = json.Unmarshal
	default:
		return fmt.Errorf("read config: %w: %s",
			errUnknownFileExt, fileExt)
	}

	confData, err = ioutil.ReadFile(configFile)

	if err != nil {
		return fmt.Errorf("read config: %w", err)
	}

	if err = unmarshal(confData, &c.conf); err != nil {
		err = fmt.Errorf("parse config: %w", err)
	}

	return err
}

// Инициализация сервиса с учетом переданных флагов.
func newCmd(args []string) (c *cmd, err error) {
	configFile, currentIp, debugMode := parseFlags(args)
	if debugMode {
		logrus.SetLevel(logrus.DebugLevel)
	}

	c = &cmd{}

	if err = c.loadConfig(configFile); err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	c.conf.CurrentIP = currentIp

	return c, nil
}

// Парсинг флагов командной строки.
func parseFlags(args []string) (configFile string, currentIP string, debugMode bool) {
	fs := flag.NewFlagSet(args[0], flag.ExitOnError)

	fs.StringVar(&configFile, "conf", defaultConfigPath, "use specified configuration file")
	fs.StringVar(&currentIP, "ip", "", "find A-record with specified IP address")
	fs.BoolVar(&debugMode, "debug", false, "enable debug mode")

	err := fs.Parse(args[2:])

	if err != nil {
		panic(err)
	}

	return
}

// Запуск сервиса.
func Run(args []string) (err error) {
	var c *cmd

	if c, err = newCmd(args); err != nil {
		return err
	}

	return core.StartService(&c.conf)
}
