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
)

var (
	conf Conf
)

func CleanTokenElogin() {
	for {
		log.Printf("Start Clean Token Elogin\n")
		web := "https://api.rmutsv.ac.th/elogin/clean/" + otp.TimeOTPxHex([]byte(conf.OTP.Key), conf.OTP.Size)
		log.Printf("web = %s\n", web)
		utils.HTTPGet(web)
		log.Printf("Clean Token Elogin Finnish\n")
		time.Sleep(time.Duration(conf.Elogin.Clean) * time.Second)
	}
}

func UpdateStudentProcess() {
	for {
		log.Printf("Start Update Student Process\n")
		web := "https://api.rmutsv.ac.th/student/report/processalldata/" + otp.TimeOTPxHex([]byte(conf.OTP.Key), conf.OTP.Size)
		log.Printf("web = %s\n", web)
		utils.HTTPGet(web)
		log.Printf("Update Student Process Finnish\n")
		time.Sleep(time.Duration(conf.Student.Cache.Update) * time.Second)
	}
}

func CleanStudentProcess() {
	for {
		log.Printf("Start Clean Student Process\n")
		web := "https://api.rmutsv.ac.th/student/report/cleanalldata/" + otp.TimeOTPxHex([]byte(conf.OTP.Key), conf.OTP.Size)
		log.Printf("web = %s\n", web)
		utils.HTTPGet(web)
		log.Printf("Clean Student Process Finnish\n")
		time.Sleep(time.Duration(conf.Student.Cache.Clean) * time.Second)
	}
}

func main() {
	utils.ProcessConfig("/etc/apiserver.yml", &conf)
	go CleanTokenElogin()
	go UpdateStudentProcess()
	go CleanStudentProcess()
	select {}
}
