// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package app

import (
	"github.com/the-bailiff/bailiff/internal/config"
	"github.com/the-bailiff/bailiff/internal/headers"
	"github.com/the-bailiff/bailiff/internal/proxy"
	"github.com/the-bailiff/bailiff/internal/session"
	"github.com/the-bailiff/bailiff/internal/store"
	"net/http"
)

// Injectors from wire.go:

func InitApp() (*App, error) {
	configConfig, err := config.ProvideConfig()
	if err != nil {
		return nil, err
	}
	sessionsStore, err := store.ProvideBaseStore(configConfig)
	if err != nil {
		return nil, err
	}
	proxyStore := session.ProvideStore(configConfig, sessionsStore)
	requestTransformerFactory, err := proxy.ProvideRequestTransformerFactory(configConfig)
	if err != nil {
		return nil, err
	}
	headersInterceptorFactory := headers.ProvideInterceptorFactory(configConfig)
	responseTransformerFactory, err := proxy.ProvideResponseTransformerFactory(headersInterceptorFactory)
	if err != nil {
		return nil, err
	}
	handler := proxy.ProvideHandler(proxyStore, requestTransformerFactory, responseTransformerFactory, headersInterceptorFactory)
	app := &App{
		Handler: handler,
		Config:  configConfig,
	}
	return app, nil
}

// wire.go:

type App struct {
	http.Handler
	config.Config
}