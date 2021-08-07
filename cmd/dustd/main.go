package main

import (
	"context"
	"log"
	"time"

	"github.com/fdawg4l/dust/lib"
	"github.com/fdawg4l/dust/lib/db"

	"github.com/spf13/viper"
)

var Config = struct {
	SensorHost   string
	DBHost       string
	Org          string
	Bucket       string
	Token        string
	SampleEveryM int64
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

	Config.SensorHost = viper.GetString("SensorHost")
	Config.DBHost = viper.GetString("DBHost")
	Config.SampleEveryM = viper.GetInt64("SampleEveryM")
	Config.Org = viper.GetString("Org")
	Config.Token = viper.GetString("Token")
	Config.Bucket = viper.GetString("Bucket")

	log.Printf("Sensor: %s", Config.SensorHost)
	log.Printf("DB: %s", Config.DBHost)
	log.Printf("freq: %d", Config.SampleEveryM)
	log.Printf("Org: %s", Config.Org)
	log.Printf("Token: %s", Config.Token)
	log.Printf("Bucket: %s", Config.Bucket)

	freq := time.Duration(Config.SampleEveryM) * time.Minute

	sensor := lib.NewSensor(context.Background(), Config.SensorHost, freq)

	dbclient, err := db.NewClient(context.Background(), Config.DBHost, Config.Bucket, Config.Org, Config.Token)
	if err != nil {
		log.Fatalf(err.Error())
	}

	for data := range sensor.Datum {
		ctx, cancel := context.WithTimeout(context.Background(), freq)
		if err = dbclient.Write(ctx, data.Host, data.Map()); err != nil {
			log.Printf("db error: " + err.Error())
		}
		cancel()
	}
}
