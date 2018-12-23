package main

import (
	"context"
	"log"
	"time"

	"github.com/fdawg4l/dust/lib"
	"github.com/fdawg4l/dust/lib/db"

	"github.com/spf13/viper"
)

var cfg = struct {
	SensorHost   string
	DBHost       string
	SampleEveryM uint8
	Debug        bool
}{}

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
	err := viper.Unmarshal(&cfg)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}
	log.Printf("Config is %#v", cfg)

	sensor := lib.NewSensor(context.Background(), cfg.SensorHost, time.Duration(cfg.SampleEveryM)*time.Minute)
	dbclient, err := db.NewClient(context.Background(), cfg.DBHost, "air_quality", "", "")
	if err != nil {
		log.Fatalf(err.Error())
	}

	for data := range sensor.Datum {
		if err = dbclient.Write(data.Host, data.Map()); err != nil {
			log.Printf("db error: " + err.Error())
		}
	}
}
