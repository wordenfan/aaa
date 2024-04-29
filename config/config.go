package config

import (
	"github.com/spf13/viper"
	"log"
	"os"
)

var AppConf = InitConfig()

type Config struct {
	viper   *viper.Viper
	LogConf *LogConf
}

type LogConf struct {
	Debug_path string
	Info_path  string
	Warn_path  string
}

func InitConfig() *Config {
	v := viper.New()
	conf := &Config{viper: v}
	workDir, _ := os.Getwd()
	conf.viper.SetConfigName("viper")
	conf.viper.SetConfigType("yaml")
	conf.viper.AddConfigPath(workDir + "/config")
	err := conf.viper.ReadInConfig()
	if err != nil {
		log.Fatalln(err)
		return nil
	}
	conf.readLogConf()
	return conf
}
func (c *Config) readLogConf() {
	logCf := &LogConf{}
	workDir, _ := os.Getwd()
	logCf.Debug_path = workDir + c.viper.GetString("zap.debugFileName")
	logCf.Info_path = workDir + c.viper.GetString("zap.infoFileName")
	logCf.Warn_path = workDir + c.viper.GetString("zap.warnFileName")

	c.LogConf = logCf
}
