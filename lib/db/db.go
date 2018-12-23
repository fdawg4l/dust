package db

import (
	"context"
	"time"

	client "github.com/influxdata/influxdb/client/v2"
)

type Client struct {
	dbhost, databaseName, username, password string
}

func NewClient(ctx context.Context, dbhost, dbname, username, password string) (*Client, error) {
	return &Client{
		dbhost:       dbhost,
		databaseName: dbname,
		username:     username,
		password:     password,
	}, nil
}

func (c *Client) Write(host string, fields map[string]interface{}) error {
	clnt, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     c.dbhost,
		Username: c.username,
		Password: c.password,
	})

	if err != nil {
		return err
	}

	defer clnt.Close()

	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  c.databaseName,
		Precision: "s",
	})

	tags := map[string]string{
		"host": host,
	}

	pt, err := client.NewPoint("air_quality", tags, fields, time.Now())
	if err != nil {
		return err
	}

	bp.AddPoint(pt)

	if err = clnt.Write(bp); err != nil {
		return err
	}

	return nil
}
