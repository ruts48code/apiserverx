package main

import (
	"log"
	"time"

	otp "github.com/ruts48code/otp4ruts"

	utils "github.com/ruts48code/utils4ruts"
)

type (
	Conf struct {
		DBS     []string      `yaml:"dbs"`
		OTP     OTPStruct     `yaml:"otp"`
		Elogin  EloginStruct  `yaml:"elogin"`
		Student StudentStruct `yaml:"student"`
		ArsDB   ARSStruct     `yaml:"ars"`
	}
	OTPStruct struct {
		Key      string `yaml:"key"`
		Size     int    `yaml:"size"`
		Interval int    `yaml:"interval"`
	}
	EloginStruct struct {
		Expire int `yaml:"expire"`
		Clean  int `yaml:"clean"`
	}
	StudentStruct struct {
		Cache StudentCacheStruct `yaml:"cache"`
	}

	StudentCacheStruct struct {
		Update int `yaml:"update"`
		Clean  int `yaml:"clean"`
	}
	ARSStruct struct {
		DB     string `yaml:"db"`
		Update int    `yaml:"update"`
		Clean  int    `yaml:"clean"`
	}
)

var (
	conf Conf
)

func CronLoop(name string, web string, interval int) {
	for {
		log.Printf("Start %s\n", name)
		webx := web + otp.TimeOTPxHex([]byte(conf.OTP.Key), conf.OTP.Size)
		log.Printf("web = %s\n", webx)
		utils.HTTPGet(webx)
		log.Printf("%s Finnish\n", name)
		time.Sleep(time.Duration(interval) * time.Second)
	}
}

func main() {
	utils.ProcessConfig("/etc/apiserver.yml", &conf)
	go CronLoop("Clean Token", "https://api.rmutsv.ac.th/elogin/clean/", conf.Elogin.Clean)
	go CronLoop("Update Student Process", "https://api.rmutsv.ac.th/student/report/processalldata/", conf.Student.Cache.Update)
	go CronLoop("Clean Student Process", "https://api.rmutsv.ac.th/student/report/cleanalldata/", conf.Student.Cache.Clean)
	go CronLoop("Update ARS", "https://api.rmutsv.ac.th/ars/process/", conf.ArsDB.Update)
	go CronLoop("Clean ARS", "https://api.rmutsv.ac.th/ars/clean/", conf.ArsDB.Clean)
	select {}
}
