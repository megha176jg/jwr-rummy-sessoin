package config

import (
	"github.com/mitchellh/mapstructure"
)

// AppConf application level properties of app
type AppConf struct {
	SessionTokenLength int `mapstructure:"kyc.kyc.auto.sessiontoken.length"`
}

func (a *AppConf) Set(m map[string]string) error {
	if err := mapstructure.WeakDecode(m, a); err != nil {
		return err
	}
	return nil
}
func (a *AppConf) GetSessionTokenLength() int {
	return a.SessionTokenLength
}
func New() *AppConf {
	var a AppConf
	return &a
}
