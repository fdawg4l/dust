package db

import (
	"context"
	"time"

	influxdb "github.com/influxdata/influxdb-client-go/v2"
)

type Client struct {
	Host   string `mapstructure: "INFLUXHOST"`
	Bucket string `mapstructure: "INFLUXBUCKET"`
	Org    string `mapstructure: "INFLUXORG"`
	Token  string `mapstructure: "INFLUXTOKEN"`
}

func NewClient(ctx context.Context, host, bucket, org, token string) (*Client, error) {
	return &Client{
		Host:   host,
		Bucket: bucket,
		Org:    org,
		Token:  token,
	}, nil
}

func (c *Client) Write(ctx context.Context, host string, fields map[string]interface{}) error {
	clnt := influxdb.NewClient(c.Host, c.Token)
	defer clnt.Close()

	writeAPI := clnt.WriteAPIBlocking(c.Org, c.Bucket)

	tags := map[string]string{
		"host": host,
	}

	pt := influxdb.NewPoint(host, tags, fields, time.Now())

	if err := writeAPI.WritePoint(ctx, pt); err != nil {
		return err
	}

	return nil
}
