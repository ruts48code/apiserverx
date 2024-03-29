package main

import (
	"log"
	"time"

	otp "github.com/ruts48code/otp4ruts"

	utils "github.com/ruts48code/utils4ruts"
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
	if utils.FileExist("/etc/apiserver.hcl") {
		utils.ProcessConfigHCL("/etc/apiserver.hcl", &conf)
		log.Printf("Load /etc/apiserver.hcl sucessfully\n")
	} else if utils.FileExist("/etc/apiserver.yml") {
		utils.ProcessConfig("/etc/apiserver.yml", &conf)
		log.Printf("Load /etc/apiserver.yml sucessfully\n")
	} else {
		log.Printf("Error: cannot load configurationfile\n")
		return
	}

	go CronLoop("Clean Token", "https://api.rmutsv.ac.th/elogin/clean/", conf.Elogin.Clean)
	go CronLoop("Update Student Process", "https://api.rmutsv.ac.th/student/report/processalldata/", conf.Student.Cache.Update)
	go CronLoop("Clean Student Process", "https://api.rmutsv.ac.th/student/report/cleanalldata/", conf.Student.Cache.Clean)
	go CronLoop("Update ARS", "https://api.rmutsv.ac.th/ars/process/", conf.ArsDB.Update)
	go CronLoop("Clean ARS", "https://api.rmutsv.ac.th/ars/clean/", conf.ArsDB.Clean)
	select {}
}
