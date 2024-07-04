package config

import (
	"github.com/mitchellh/mapstructure"
)

// AppConf application level properties of app
type AppConf struct {
	MinAge string `mapstructure:"kyc.kyc.auto.age.min"`
}

func (a *AppConf) Set(m map[string]string) error {
	if err := mapstructure.WeakDecode(m, a); err != nil {
		return err
	}
	return nil
}
func (a *AppConf) GetMinAge() string {
	return a.MinAge
}
func New() *AppConf {
	var a AppConf
	return &a
}
