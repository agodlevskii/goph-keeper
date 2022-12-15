package configs

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"reflect"

	"github.com/caarlos0/env"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

var ErrDifferentTypes = errors.New("unknown configs type")
var ErrMissing = errors.New("configs path is missing")
var ErrUnknownType = errors.New("unknown configs type")

func UpdateConfigFromEnv(cfg any) error {
	return env.Parse(cfg)
}

func UpdateConfigFromFile(cfg any, fCfg any, fPath string) error {
	var err error
	if fPath == "" {
		return ErrMissing
	}
	if fCfg, err = getConfigFromFile(filepath.Clean(fPath), fCfg); err != nil {
		return err
	}
	return setConfigFromFile(cfg, fCfg)
}

func getConfigsDirPath() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(wd, "..", "..", "..", "configs"), nil
}

func getConfigFromFile(fName string, fCfg any) (any, error) {
	fName = filepath.Clean(fName)
	cfgPath, err := getConfigsDirPath()
	if err != nil {
		return nil, err
	}

	cfgFile, err := os.Open(filepath.Join(cfgPath, fName))
	if err != nil {
		return nil, err
	}
	defer closeFile(cfgFile)

	cfgBytes, err := io.ReadAll(cfgFile)
	if err != nil {
		return nil, err
	}

	switch filepath.Ext(fName) {
	case ".json":
		err = json.Unmarshal(cfgBytes, fCfg)
	case ".yaml":
	case ".yml":
		err = yaml.Unmarshal(cfgBytes, fCfg)
	default:
		err = ErrUnknownType
	}

	return fCfg, err
}

func setConfigFromFile(cfg any, fCfg any) error {
	rCfg := reflect.ValueOf(cfg)
	if rCfg.Kind() == reflect.Pointer {
		rCfg = reflect.Indirect(rCfg)
	}

	rFileCfg := reflect.ValueOf(fCfg)
	if rFileCfg.Kind() == reflect.Pointer {
		rFileCfg = reflect.Indirect(rFileCfg)
	}

	if rCfg.Type() != rFileCfg.Type() {
		return ErrDifferentTypes
	}

	for i := 0; i < rCfg.NumField(); i++ {
		rField := rCfg.Type().Field(i).Name
		rValue := rCfg.FieldByName(rField)

		if rValue.IsZero() && rValue.CanSet() {
			if fileValue := rFileCfg.FieldByName(rField); !fileValue.IsZero() {
				rValue.Set(fileValue)
			}
		}
	}

	return nil
}

func closeFile(f *os.File) {
	if err := f.Close(); err != nil {
		log.Error(err)
	}
}
