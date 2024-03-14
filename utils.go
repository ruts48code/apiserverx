package main

import (
	"database/sql"
	"errors"
	"log"

	util "github.com/ruts48code/utils4ruts"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"gopkg.in/yaml.v3"
)

type (
	Conf struct {
		DBUsername string       `yaml:"dbusername"`
		DBPassword string       `yaml:"dbpassword"`
		DBName     string       `yaml:"dbname"`
		DBParam    string       `yaml:"dbparam"`
		DBType     string       `yaml:"dbtype"`
		DBS        []string     `yaml:"dbs"`
		Elogin     EloginStruct `yaml:"elogin"`
	}

	EloginStruct struct {
		Expire int `yaml:"expire"`
		Clean  int `yaml:"clean"`
	}
)

func processConfig() {
	confdata := util.ReadFile("/etc/apiserver.yml")
	yaml.Unmarshal(confdata, &conf)
}

func getDBS() (*sql.DB, error) {
	dbN := util.RandomArrayString(conf.DBS)
	dbConnect := false
	var db *sql.DB
	var err error
	qstring := ""
	for i := range dbN {
		switch conf.DBType {
		case "postgres":
			if conf.DBPassword == "" {
				qstring = "postgres://" + conf.DBUsername + "@" + dbN[i] + "/" + conf.DBName
			} else {
				qstring = "postgres://" + conf.DBUsername + ":" + conf.DBPassword + "@" + dbN[i] + "/" + conf.DBName
			}
			if conf.DBParam != "" {
				qstring = qstring + "?" + conf.DBParam
			}
		case "mysql":
			if conf.DBPassword == "" {
				qstring = conf.DBUsername + "@tcp(" + dbN[i] + ":3306)/" + conf.DBName
			} else {
				qstring = conf.DBUsername + ":" + conf.DBPassword + "@tcp(" + dbN[i] + ":3306)/" + conf.DBName
			}
			if conf.DBParam != "" {
				qstring = qstring + "?" + conf.DBParam
			}
		}

		log.Printf("Try db %s\n", dbN[i])
		db, err = sql.Open(conf.DBType, qstring)
		if err != nil {
			log.Printf("Error: fail to open db %s - %v\n", dbN[i], err)
			continue
		}

		err = db.Ping()
		if err != nil {
			log.Printf("Error: fail to ping db %s - %v\n", dbN[i], err)
			db.Close()
			continue
		}
		log.Printf("Connect to db %s\n", dbN[i])
		dbConnect = true
		break
	}
	if !dbConnect {
		log.Printf("Error: Cannot connect to db\n")
		return nil, errors.New("cannot connect to db")
	}
	return db, nil
}
