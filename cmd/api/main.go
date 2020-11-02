package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/ck3g/homebugh-api/pkg/models"
	"github.com/ck3g/homebugh-api/pkg/models/mysql"
	"gopkg.in/yaml.v2"
)

type application struct {
	users models.UserStorage
}

type config struct {
	TLS struct {
		CertPemFile string `yaml:"cert_pem_file"`
		KeyPemFile  string `yaml:"key_pem_file"`
	} `yaml:"tls"`
	Server struct {
		Addr string `yaml:"addr"`
	} `yaml:"server"`
}

func main() {
	configFile := flag.String("config", "./config.yml", "The app configuration file path")
	flag.Parse()

	f, err := os.Open(*configFile)
	if err != nil {
		flag.PrintDefaults()
		panic("cannot load the config file")
	}
	defer f.Close()

	var cfg config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		panic("cannot parse the config file")
	}

	dsn := "root@/homebugh_test?parseTime=true"
	db, err := openDB(dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	app := &application{
		users: &mysql.UserModel{DB: db},
	}

	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	srv := &http.Server{
		Addr:      cfg.Server.Addr,
		Handler:   app.routes(),
		TLSConfig: tlsConfig,
	}

	fmt.Printf("Listening on %s, CTRL+C to stop\n", srv.Addr)
	err = srv.ListenAndServeTLS(cfg.TLS.CertPemFile, cfg.TLS.KeyPemFile)
	if err != nil {
		fmt.Println(err)
		panic("Cannot start the API server")
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
