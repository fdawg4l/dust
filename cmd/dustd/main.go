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
	viper.SetEnvPrefix("DUST")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Error reading config file, %s", err)
	}

	log.Printf("Sensor: %s", viper.GetString("SensorHost"))
	log.Printf("DB: %s", viper.GetString("DBHost"))
	log.Printf("freq: %d", viper.GetInt64("SampleEveryM"))

	sensor := lib.NewSensor(context.Background(),
		viper.GetString("SensorHost"),
		time.Duration(viper.GetInt64("SampleEveryM"))*time.Minute)

	dbclient, err := db.NewClient(context.Background(),
		viper.GetString("DBHost"),
		"air_quality", "", "")
	if err != nil {
		log.Fatalf(err.Error())
	}

	for data := range sensor.Datum {
		if err = dbclient.Write(data.Host, data.Map()); err != nil {
			log.Printf("db error: " + err.Error())
		}
	}
}
