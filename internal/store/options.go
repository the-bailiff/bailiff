package store

import (
	"errors"
	"net/http"

	"github.com/gorilla/sessions"

	"github.com/the-bailiff/bailiff/internal/config"
)

const sameSiteDefault = "Default"
const sameSiteLax = "Lax"
const sameSiteStrict = "Strict"
const sameSiteNone = "None"

type options struct {
	Name string
	sessions.Options
}

func provideOptions(cfg config.Config) (options, error) {
	cookie := cfg.Cookie

	opts := options{
		Name: cookie.Name,
		Options: sessions.Options{
			Path:     cookie.Path,
			Domain:   cookie.Domain,
			MaxAge:   cookie.MaxAge,
			Secure:   cookie.Secure,
			HttpOnly: cookie.HttpOnly,
			SameSite: http.SameSiteDefaultMode,
		},
	}

	switch cookie.SameSite {
	case sameSiteDefault:
		opts.SameSite = http.SameSiteDefaultMode
	case sameSiteLax:
		opts.SameSite = http.SameSiteLaxMode
	case sameSiteStrict:
		opts.SameSite = http.SameSiteStrictMode
	case sameSiteNone:
		opts.SameSite = http.SameSiteNoneMode
	case "":
		// do nothing, already set default
	default:
		return options{}, errors.New("Invalid value for Cookie.SameSite: " + cookie.SameSite)
	}

	return opts, nil
}
