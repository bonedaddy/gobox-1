package main

import (
	"flag"
	"log"
	"net/http"
	_ "net/http/pprof" // !DEBUG
	"os"
	"os/signal"
	"syscall"

	"github.com/guerillagrow/gobox/controllers"
	"github.com/guerillagrow/gobox/models"
	_ "github.com/guerillagrow/gobox/routers"

	"github.com/astaxie/beego"
)

func main() {

	sigs := make(chan os.Signal, 1)
	//done := make(chan bool, 1)

	models.ARG_ConfigFile = flag.String("c", "./conf/raspberrypi.json", "Config file")
	models.ARG_DBFile = flag.String("d", "./conf/main.db", "Database file")
	models.ARG_Debug = flag.Bool("debug", false, "Debug mode")
	flag.Parse()

	if *models.ARG_Debug == true {
		go func() {
			// !DEBUG
			log.Println(http.ListenAndServe(":6060", nil))
		}()
	}

	models.Init() // Initialize models database etc.

	go func() {
		sig := <-sigs
		if sig == os.Interrupt || sig == os.Kill || sig == syscall.SIGTERM {
			models.GoBox.Stop()
			os.Exit(1)
		}
		//done <- true
	}()

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	models.GoBox.Start()

	beego.BConfig.WebConfig.Session.SessionOn = true
	beego.BConfig.WebConfig.Session.SessionAutoSetCookie = true
	beego.BConfig.MaxMemory = 1 << 26
	//beego.BConfig.CopyRequestBody = true // required for json request bodies!

	beego.ErrorController(&controllers.ErrorController{})
	beego.Run()

}
