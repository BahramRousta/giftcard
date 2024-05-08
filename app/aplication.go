package app

import (
	"context"
	"fmt"
	"giftcard/config"
	"giftcard/internal/adaptor/giftcard"
	"giftcard/internal/adaptor/logstash"
	"giftcard/internal/adaptor/postgres"
	"giftcard/internal/adaptor/redis"
	"giftcard/internal/adaptor/trace"
	customerModule "giftcard/internal/modules/customer"
	orderModule "giftcard/internal/modules/order"
	shopModule "giftcard/internal/modules/shop"
	"giftcard/internal/server"
	"giftcard/pkg/logger"
	"go.uber.org/fx"
	"log"
	"os"
	"time"
)

// Start Application func
func Start() {
	fmt.Println("\n\n--------------------------------")
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	fxNew := fx.New(
		fx.Provide(config.C),
		fx.Provide(postgres.DB),
		fx.Provide(redis.NewRedis),
		customerModule.Module,
		orderModule.Module,
		shopModule.Module,
		fx.Provide(giftcard.NewGiftCard),
		//fx.Provide(config.NewLogger),
		fx.Provide(logstash.NewLogStash),
		fx.Invoke(trace.InitGlobalTracer),
		fx.Invoke(logger.InitGlobalLogger),
		fx.Provide(server.NewServer),
		fx.Invoke(serve),
	)

	if err := fxNew.Start(context.Background()); err != nil {
		log.Println(err)
		return
	}
	if val := <-fxNew.Done(); val == os.Interrupt {
		return
	}

	stopCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := fxNew.Stop(stopCtx); err != nil {
		log.Println(err)
		return
	}

}

func serve(lc fx.Lifecycle, server server.IServer) {
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			return server.Run()
		},
		OnStop: func(ctx context.Context) error {
			return server.Shutdown()
		},
	})
}
