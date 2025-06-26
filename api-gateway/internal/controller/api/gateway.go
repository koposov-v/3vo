package api

import (
	orderv1 "api-gateway/pkg/order/v1"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Gateway struct {
	authServiceURL string
	client         *http.Client
	orderClient    orderv1.OrderServiceClient
	logger         *logrus.Logger
}

func NewGateway(
	orderClient orderv1.OrderServiceClient,
	authServiceURL string,
	logger *logrus.Logger,
) *Gateway {
	return &Gateway{
		client:         &http.Client{}, //TODO::Надо бы в adapter - но я думаю это и так понятно
		orderClient:    orderClient,
		authServiceURL: authServiceURL,
		logger:         logger,
	}
}
