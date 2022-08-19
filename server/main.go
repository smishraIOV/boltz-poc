package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/btcsuite/btcd/chaincfg"
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

func startServer(boltz *connectors.Boltz, rsk *connectors.RSK, db *storage.DB) {

	srv = http.New(boltz, rsk, db)
	log.Debug("registering server (this might take a while)")
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

	log.Info("starting boltz-poc server")
	log.Debugf("loaded config %+v", cfg)

	// INIT DB
	db, err := storage.Connect(cfg.DB.Path)
	if err != nil {
		log.Fatal("error connecting to DB: ", err)
	}

	// INIT RSK
	rsk, err := connectors.NewRSK(cfg.ErpKeys)
	if err != nil {
		log.Fatal("RSK error: ", err)
	}

	err = rsk.Connect(cfg.RSK.Endpoint, cfg.RSK.ChainId)
	if err != nil {
		log.Fatal("error connecting to RSK: ", err)
	}

	// INIT Boltz
	boltz, err := connectors.NewBoltz(cfg.Boltz.Endpoint, &chaincfg.SimNetParams)
	if err != nil {
		log.Fatal("Boltz error: ", err)
	}

	info, errBoltz := boltz.GetReverseSwapInfo()
	if errBoltz != nil {
		log.Fatal("error GetReverseSwapInfo to Boltz: ", errBoltz)
	}
	log.Debugf("Verified connection to Boltz GetReverseSwapInfo -> %+v", info)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	startServer(boltz, rsk, db)

	<-done

	srv.Shutdown()
	rsk.Close()

	err = db.Close()
	if err != nil {
		log.Fatal("error closing DB connection: ", err)
	}
}
