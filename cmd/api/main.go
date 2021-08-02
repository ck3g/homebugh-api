package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/ck3g/homebugh-api/pkg/models"
	"github.com/ck3g/homebugh-api/pkg/models/mysql"
	"gopkg.in/yaml.v2"
)

type application struct {
	environment string
	models      models.Models
}

func main() {
	configFile := flag.String("config", "./config.yml", "The app configuration file path")
	environment := flag.String("env", "development", "Current app environment. [production|development|test]")
	flag.Parse()

	env := setEnv(*environment)

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

	db, err := openDB(cfg.dsn(env))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	app := &application{
		environment: env,
		models: models.Models{
			Users: &mysql.UserModel{DB: db},
		},
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

	fmt.Println()
	fmt.Printf("The application starting in %s\n", env)
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

func setEnv(env string) string {
	for _, e := range []string{"test", "development", "production"} {
		if e == strings.ToLower(env) {
			return e
		}
	}

	return "development"
}
