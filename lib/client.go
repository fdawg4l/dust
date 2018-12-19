package lib

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

const (
	defaultHost    = "dust.local"
	defaultTimeout = 5 * time.Minute
)

type Sensor struct {
	host   string
	ctx    context.Context
	cancel context.CancelFunc
	t      *time.Ticker

	Datum chan *Datum
	Err   error
}

func NewSensor(ctx context.Context, host string, every time.Duration) *Sensor {
	h := defaultHost
	if host != "" {
		h = host
	}

	ctx, cancel := context.WithCancel(ctx)
	s := &Sensor{
		host:   h,
		ctx:    ctx,
		cancel: cancel,
		Datum:  make(chan *Datum),
		t:      time.NewTicker(every),
	}

	go s.worker()

	return s
}

func (s *Sensor) worker() {
	for {
		select {
		case <-s.t.C:
			ctx, cancel := context.WithTimeout(s.ctx, defaultTimeout)
			d, err := s.do(ctx)
			cancel()
			if err != nil {
				s.Err = err
				return
			}

			s.Datum <- d
		case <-s.ctx.Done():
			return
		}
	}
}

func (s *Sensor) do(ctx context.Context) (*Datum, error) {
	req, err := http.NewRequest("GET", s.host, nil)
	if err != nil {
		return nil, err
	}
	c := http.DefaultClient
	resp, err := c.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	d := new(Datum)
	dec := json.NewDecoder(resp.Body)
	if err = dec.Decode(d); err != nil {
		return nil, err
	}

	return d, nil
}

func (s *Sensor) Stop() {
	s.cancel()
	s.t.Stop()
	close(s.Datum)
}
