package main

import (
	"atm/atm"
	"atm/config"
	"atm/server"
	"flag"

	log "github.com/sirupsen/logrus"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	defaultDBPath     = "atm.db"
	defaultLogLevel   = "debug"
	defaultListenAddr = "localhost:15000"
	defaultAdminPin   = "0000"
)

func main() {
	// Command line arguments
	configFile := ""
	flag.StringVar(&configFile, "config", "", "path for the config file")
	flag.Parse()

	// Read config file
	if configFile == "" {
		log.Fatal("config file missing")
	}
	conf := config.NewConfig()
	if err := conf.ReadConfigFile(configFile); err != nil {
		log.Fatalf("error reading config file %q", configFile)
	}

	// set log level
	log.SetLevel(config.LogLevel(conf.StrValue("logLevel", defaultLogLevel)))

	// setup database
	dbPath := conf.StrValue("dbPath", defaultDBPath)
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatal("Database connection failed", err)
	}

	// create new ATM with DB
	adminPin := conf.StrValue("adminPin", defaultAdminPin)
	atm := atm.NewATM(db, adminPin)
	if atm == nil {
		log.Fatal("Error initializing ATM")
	}

	// start http server for ATM
	listenAddr := conf.StrValue("listenAddr", defaultListenAddr)
	atmServer := server.NewAtmServer(listenAddr, atm)
	atmServer.InitRouter()
}
