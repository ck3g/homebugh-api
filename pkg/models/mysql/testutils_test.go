package mysql

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"gopkg.in/yaml.v2"
)

type dbConfig struct {
	Database struct {
		Test struct {
			Database string `yaml:"database"`
			Username string `yaml:"username"`
			Password string `yaml:"password"`
		} `yaml:"test"`
	} `yaml:"database"`
}

func dsn() string {
	f, err := os.Open("../../../config.yml")
	if err != nil {
		panic("cannot load the config file")
	}
	defer f.Close()

	var cfg dbConfig
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		panic("cannot parse the config file")
	}

	c := cfg.Database.Test
	return fmt.Sprintf("%s:%s@/%s?parseTime=true&multiStatements=true", c.Username, c.Password, c.Database)
}

func newTestDB(t *testing.T) (*sql.DB, func()) {
	db, err := sql.Open("mysql", dsn())
	if err != nil {
		t.Fatal(err)
	}

	script, err := ioutil.ReadFile("./testdata/setup.sql")
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec(string(script))
	if err != nil {
		t.Fatal(err)
	}

	return db, func() {
		script, err := ioutil.ReadFile("./testdata/teardown.sql")
		if err != nil {
			t.Fatal(err)
		}

		_, err = db.Exec(string(script))
		if err != nil {
			t.Fatal(err)
		}

		db.Close()
	}
}
