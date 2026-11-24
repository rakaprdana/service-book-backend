package utils

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func BindFromJSON(dest any, filename, path string) error {
	v := viper.New()
	v.SetConfigType("json")
	v.AddConfigPath(path)
	v.SetConfigFile(filename)

	if err := v.ReadInConfig(); err != nil {
		return err
	}

	if err := v.Unmarshal(&dest); err != nil {
		logrus.Errorf("Failed to unmarshal: %v", err)
		return err
	}

	return nil
}
