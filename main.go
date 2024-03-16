package main

import (
	"log"
	"time"

	otp "github.com/ruts48code/otp4ruts"

	utils "github.com/ruts48code/utils4ruts"
)

type (
	Conf struct {
		DBS    []string     `yaml:"dbs"`
		OTP    OTPStruct    `yaml:"otp"`
		Elogin EloginStruct `yaml:"elogin"`
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
)

var (
	conf Conf
)

func CleanTokenElogin() {
	for {
		log.Printf("Start Clean Token Elogin\n")
		otphex := otp.TimeOTPxHex([]byte(conf.OTP.Key), conf.OTP.Size)
		log.Printf("OTP = %s\n", otphex)
		utils.HTTPGet("https://api.rmutsv.ac.th/elogin/clean/" + otphex)
		log.Printf("Clean Token Elogin Finnish\n")
		time.Sleep(time.Duration(conf.Elogin.Clean) * time.Second)
	}
}

func main() {
	utils.ProcessConfig("/etc/apiserver.yml", &conf)
	go CleanTokenElogin()
	select {}
}
