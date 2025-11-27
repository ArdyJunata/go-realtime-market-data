package handler

type Service interface{}

func NewHandler(service Service) *handler {
	return &handler{
		service: service,
	}
}

type handler struct {
	service Service
}
