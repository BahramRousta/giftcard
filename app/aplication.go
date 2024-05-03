package app

import (
	"context"
	"fmt"
	"giftCard/config"
	"giftCard/internal/adaptor/giftcard"
	"giftCard/internal/adaptor/postgres"
	delivery "giftCard/internal/modules/customer/delivery/http"
	"giftCard/internal/modules/customer/repository"
	"giftCard/internal/modules/customer/usecase"
	orderModule "giftCard/internal/modules/order"
	shopModule "giftCard/internal/modules/shop"
	"giftCard/internal/server"
	"go.uber.org/fx"
	"log"
)

// Start Application func
func Start() {
	fmt.Println("\n\n--------------------------------")
	// if you go code crashed we get error and line
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// init configs

	//fxNew :=
	fx.New(
		fx.Provide(config.C),
		fx.Provide(usecase.NewCustomerUseCase),
		fx.Provide(delivery.NewCustomerInfoHandler),
		fx.Provide(repository.NewWalletRepository),
		fx.Provide(postgres.DB),
		//customerModule.Module,
		orderModule.Module,
		shopModule.Module,
		fx.Provide(giftcard.NewGiftCard),
		fx.Provide(server.NewServer),
		fx.Invoke(serve),
	).Run()

	//if err := fxNew.Run(); err != nil {
	//	log.Println(err)
	//	return
	//}
	//if val := <-fxNew.Done(); val == os.Interrupt {
	//	return
	//}
	//
	//stopCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//defer cancel()
	//
	//if err := fxNew.Stop(stopCtx); err != nil {
	//	log.Println(err)
	//	return
	//}

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
