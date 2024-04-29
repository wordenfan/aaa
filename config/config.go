package viper_conf

import (
	"github.com/spf13/viper"
	"log"
	"os"
	"path/filepath"
	"tk-boot-worden/zap_log"
)

// 全局配置变量
var GlobalConf *Config

type Config struct {
	viper    *viper.Viper
	LogConf  *LogConf
	GrpcConf *GrpcConf
	EtcdConf *EtcdConf
}

type LogConf struct {
	Debug_path string
	Info_path  string
	Warn_path  string
}
type GrpcConf struct {
	Addr    string
	Name    string
	Version string
	Weight  int64
}
type EtcdConf struct {
	AddrS []string
}

func (conf *Config) InitConfig() (err error) {
	conf.viper = viper.New()
	workDir, _ := os.Getwd()
	conf.viper.SetConfigName("viper")
	conf.viper.SetConfigType("yaml")
	conf.viper.AddConfigPath(workDir + "/config")
	err = conf.viper.ReadInConfig()
	if err != nil {
		log.Fatalln(err)
		return err
	}
	conf.setLogConf()
	conf.setGrpcConf()
	conf.setEtcdConf()
	GlobalConf = conf
	return nil
}
func (c *Config) setEtcdConf() {
	etcdCf := &EtcdConf{}
	var addrS []string
	err := c.viper.UnmarshalKey("etcd.addrs", &addrS)
	if err != nil {
		log.Fatalln(err)
	}
	etcdCf.AddrS = addrS
	c.EtcdConf = etcdCf
}
func (c *Config) setLogConf() {
	logCf := &LogConf{}
	workDir, _ := os.Getwd()
	logCf.Debug_path = filepath.Join(workDir, c.viper.GetString("zap.debugFileName"))
	logCf.Info_path = filepath.Join(workDir, c.viper.GetString("zap.infoFileName"))
	logCf.Warn_path = filepath.Join(workDir, c.viper.GetString("zap.warnFileName"))

	c.LogConf = logCf
}

func (c *Config) setGrpcConf() {
	grpcCf := &GrpcConf{}
	grpcCf.Name = "worden_grpc"
	grpcCf.Addr = c.viper.GetString("grpc.addr")
	grpcCf.Version = c.viper.GetString("grpc.version")
	grpcCf.Weight = c.viper.GetInt64("grpc.weight")
	c.GrpcConf = grpcCf
}

func GetZapLogConf() *zap_log.ZapLogConfig {
	zap_log_cfg := &zap_log.ZapLogConfig{
		DebugFileName: GlobalConf.LogConf.Debug_path,
		InfoFileName:  GlobalConf.LogConf.Info_path,
		WarnFileName:  GlobalConf.LogConf.Warn_path,
		MaxSize:       500,
		MaxAge:        28,
		MaxBackups:    3,
	}
	return zap_log_cfg
}
