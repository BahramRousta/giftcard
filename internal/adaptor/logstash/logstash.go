package logstash

import (
	"context"
	"giftcard/config"
	"log"
	"net"
	"time"

	"go.uber.org/fx"
)

type LogStash struct {
	conn net.Conn
}

func NewLogStash(lc fx.Lifecycle) *LogStash {
	var err error
	logstash := LogStash{}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logstash.conn, err = net.DialTimeout("udp", config.C().Logstash.Endpoint, time.Duration(config.C().Logstash.Timeout)*time.Second)
			if err != nil {
				return err
			}
			log.Printf("logstash connected successfully \n")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return logstash.conn.Close()
		},
	})
	return &logstash
}

func (l *LogStash) Write(p []byte) (int, error) {
	return l.conn.Write(p)
}
