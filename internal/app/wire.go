//+build wireinject

package app

import (
	"net/http"

	"github.com/google/wire"

	"github.com/the-bailiff/bailiff/internal/config"
	"github.com/the-bailiff/bailiff/internal/headers"
	"github.com/the-bailiff/bailiff/internal/proxy"
	"github.com/the-bailiff/bailiff/internal/session"
	"github.com/the-bailiff/bailiff/internal/store"
)

type App struct {
	http.Handler
	config.Config
}

func InitApp() (*App, error) {
	wire.Build(
		config.ProvideConfig,
		headers.ProvideInterceptorFactory,
		proxy.ProvideRequestTransformerFactory,
		proxy.ProvideResponseTransformerFactory,
		proxy.ProvideHandler,
		session.ProvideStore,
		store.ProvideBaseStore,
		wire.Struct(new(App), "Handler", "Config"),
	)

	return new(App), nil
}
