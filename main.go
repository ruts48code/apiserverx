package main

import (
	"log"
	"time"

	dbs "github.com/ruts48code/dbs4ruts"
	utils "github.com/ruts48code/utils4ruts"
)

type (
	Conf struct {
		DBS    []string     `yaml:"dbs"`
		Elogin EloginStruct `yaml:"elogin"`
	}

	EloginStruct struct {
		Expire int `yaml:"expire"`
		Clean  int `yaml:"clean"`
	}
)

var (
	conf Conf
)

func CleanTokenElogin() {
	for {
		db, err := dbs.OpenDBS(conf.DBS)
		if err != nil {
			log.Printf("Error: %v\n", err)
			time.Sleep(time.Duration(conf.Elogin.Clean) * time.Second)
			continue
		}

		ts := utils.GetTimeStamp(time.Now().Add(time.Duration(conf.Elogin.Expire) * time.Second * -1))
		_, err = db.Exec("DELETE FROM token WHERE timestamp < ?;", ts)
		if err != nil {
			log.Printf("Error: %v\n", err)
		} else {
			log.Printf("Clean Successful\n")
		}
		db.Close()
		time.Sleep(time.Duration(conf.Elogin.Clean) * time.Second)
	}
}

func main() {
	utils.ProcessConfig("/etc/apiserver.yml", &conf)
	go CleanTokenElogin()
	select {}
}
