package server

type IServer interface {
	SetUpServer(container DeliveryContainer)
	Shutdown() error
	Run() error
}
