package main

import (
	"errors"
	"flag"
	"github.com/1f349/gomvn"
	"github.com/1f349/gomvn/routes"
	exitReload "github.com/MrMelon54/exit-reload"
	"gopkg.in/yaml.v3"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type startupConfig struct {
	Name       string   `yaml:"name"`
	Listen     string   `yaml:"listen"`
	Repository []string `yaml:"repository"`
}

func main() {
	configPath := flag.String("conf", "", "/path/to/config.yml : path to the config file")
	flag.Parse()

	log.Println("[GoMVN] Starting...")

	if *configPath == "" {
		log.Fatal("[GoMVN] Error: config flag is missing")
		return
	}

	openConf, err := os.Open(*configPath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Fatal("[Violet] Error: missing config file")
		} else {
			log.Fatal("[Violet] Error: open config file: ", err)
		}
		return
	}

	var config startupConfig
	err = yaml.NewDecoder(openConf).Decode(&config)
	if err != nil {
		log.Fatal("[GoMVN] Error: invalid config file: ", err)
	}

	// working directory is the parent of the config file
	wd := filepath.Dir(*configPath)
	db, err := gomvn.InitDB(filepath.Join(wd, "gomvn.sqlite3.db"))
	if err != nil {
		log.Fatal("[GoMVN] Error: invalid database: ", err)
	}
	repoBasePath := filepath.Join(wd, "repositories")
	err = os.MkdirAll(repoBasePath, os.ModePerm)
	if err != nil {
		log.Fatal("[GoMVN] Error: failed to create repositories directory: ", err)
	}

	srv := &http.Server{
		Addr:              config.Listen,
		Handler:           routes.Router(db, config.Name, repoBasePath, config.Repository),
		ReadTimeout:       time.Minute,
		ReadHeaderTimeout: time.Minute,
		WriteTimeout:      time.Minute,
		IdleTimeout:       time.Minute,
		MaxHeaderBytes:    5000,
	}
	go func() {
		err = srv.ListenAndServe()
		if err != nil && !errors.Is(http.ErrServerClosed, err) {
			log.Println("Serve HTTP Error:", err)
		}
	}()

	exitReload.ExitReload("GoMVN", func() {}, func() {
		err := srv.Close()
		if err != nil {
			log.Println(err)
		}
	})
}
