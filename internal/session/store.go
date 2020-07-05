package session

import (
	"net/http"

	s "github.com/gorilla/sessions"

	"github.com/the-bailiff/bailiff/internal/config"
	"github.com/the-bailiff/bailiff/internal/proxy"
)

type StoreImpl struct {
	name    string
	rolling bool
	store   s.Store
}

func ProvideStore(cfg config.Config, baseStore s.Store) proxy.Store {
	return &StoreImpl{
		name:    cfg.Cookie.Name,
		rolling: cfg.Rolling,
		store:   baseStore,
	}
}

func (s *StoreImpl) Get(w *http.ResponseWriter, r *http.Request) (proxy.Session, error) {
	session, err := s.store.Get(r, s.name)
	if err != nil {
		return nil, err
	}
	return &SessionImpl{
		changed: false,
		rolling: s.rolling,
		writer:  w,
		request: r,
		session: session,
	}, nil
}
