package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/patogallaiov/boltz-poc/connectors"
	"github.com/patogallaiov/boltz-poc/http"
	"github.com/patogallaiov/boltz-poc/storage"
	log "github.com/sirupsen/logrus"
	"github.com/tkanos/gonfig"
)

var (
	cfg config
	srv http.Server
)

func loadConfig() {
	err := gonfig.GetConf("config.json", &cfg)

	if err != nil {
		log.Fatalf("error loading config file: %v", err)
	}
}

func initLogger() {
	if cfg.Debug {
		log.SetLevel(log.DebugLevel)
	}
	if cfg.LogFile == "" {
		return
	}
	f, err := os.OpenFile(cfg.LogFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)

	if err != nil {
		log.Error(fmt.Sprintf("cannot open file %v: ", cfg.LogFile), err)
	} else {
		log.SetOutput(f)
	}
}

func startServer(rsk *connectors.RSK, db *storage.DB) {

	srv = http.New(rsk, db)
	log.Debug("registering local provider (this might take a while)")
	port := cfg.Server.Port

	if port == 0 {
		port = 8080
	}
	go func() {
		err := srv.Start(port)

		if err != nil {
			log.Error("server error: ", err.Error())
		}
	}()
}

func main() {
	loadConfig()
	initLogger()
	rand.Seed(time.Now().UnixNano())

	log.Info("starting liquidity provider server")
	log.Debugf("loaded config %+v", cfg)

	db, err := storage.Connect(cfg.DB.Path)
	if err != nil {
		log.Fatal("error connecting to DB: ", err)
	}

	rsk, err := connectors.NewRSK(cfg.ErpKeys)
	if err != nil {
		log.Fatal("RSK error: ", err)
	}

	err = rsk.Connect(cfg.RSK.Endpoint, cfg.RSK.ChainId)
	if err != nil {
		log.Fatal("error connecting to RSK: ", err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	startServer(rsk, db)

	<-done

	srv.Shutdown()
	rsk.Close()

	err = db.Close()
	if err != nil {
		log.Fatal("error closing DB connection: ", err)
	}
}
