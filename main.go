package main

import (
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"

	utils "github.com/ruts48code/utils4ruts"
)

var (
	conf Conf
)

func CleanTokenElogin() {
	for {
		db, err := getDBS()
		if err != nil {
			log.Printf("Error: %v\n", err)
			time.Sleep(time.Duration(conf.Elogin.Clean) * time.Second)
			continue
		}

		ts := utils.GetTimeStamp(time.Now().Add(time.Duration(conf.Elogin.Expire) * time.Second * -1))
		qstring := ""
		switch conf.DBType {
		case "postgres":
			qstring = "DELETE FROM token WHERE timestamp < $1;"
		case "mysql":
			qstring = "DELETE FROM token WHERE timestamp < ?;"
		}
		log.Printf("query is %s -- time is %s\n", qstring, ts)
		_, err = db.Exec(qstring, ts)
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
