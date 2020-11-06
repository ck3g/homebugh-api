package main

import "fmt"

type config struct {
	TLS struct {
		CertPemFile string `yaml:"cert_pem_file"`
		KeyPemFile  string `yaml:"key_pem_file"`
	} `yaml:"tls"`
	Server struct {
		Addr string `yaml:"addr"`
	} `yaml:"server"`
	Database struct {
		Production  dbConfig `yaml:"production"`
		Development dbConfig `yaml:"development"`
		Test        dbConfig `yaml:"test"`
	} `yaml:"database"`
}

type dbConfig struct {
	Driver   string `yaml:"driver"`
	Database string `yaml:"database"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

func (cfg *config) dsn(env string) string {
	c := map[string]dbConfig{
		"production":  cfg.Database.Production,
		"development": cfg.Database.Development,
		"test":        cfg.Database.Test,
	}[env]

	// "user:password@/database?parseTime=true"
	return fmt.Sprintf("%s:%s@/%s?parseTime=true", c.Username, c.Password, c.Database)
}
