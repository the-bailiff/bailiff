package proxy

import (
	"net/http"
	"net/http/httputil"
)

type Session interface {
	Values() http.Header
	Set(key string, values []string)
	Del(key string)
	Save() error
	Rolling() bool
	Changed() bool
}

type RequestTransformer interface {
	Transform(*http.Request)
}

type RequestTransformerFactory interface {
	New(Session) RequestTransformer
}

type ResponseTransformer interface {
	Transform(*http.Response) error
}

type ResponseTransformerFactory interface {
	New(Session) ResponseTransformer
}

type Store interface {
	Get(*http.ResponseWriter, *http.Request) (Session, error)
}

type Handler struct {
	Store
	RequestTransformerFactory
	ResponseTransformerFactory
	HeadersInterceptorFactory
}

func ProvideHandler(
	store Store,
	reqTrFac RequestTransformerFactory,
	resTrFac ResponseTransformerFactory,
	interceptorFac HeadersInterceptorFactory,
) http.Handler {
	return &Handler{
		store,
		reqTrFac,
		resTrFac,
		interceptorFac,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	session, err := h.Store.Get(&w, r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	reqTransformer := h.RequestTransformerFactory.New(session)
	resTransformer := h.ResponseTransformerFactory.New(session)

	rp := &httputil.ReverseProxy{
		Director:       reqTransformer.Transform,
		ModifyResponse: resTransformer.Transform,
	}

	rp.ServeHTTP(w, r)
}
