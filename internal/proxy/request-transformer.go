package proxy

import (
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/the-bailiff/bailiff/internal/config"
)

type RequestTransformerFactoryImpl struct {
	prefix string
	url    *url.URL
}

func ProvideRequestTransformerFactory(cfg config.Config) (RequestTransformerFactory, error) {
	dest, err := url.Parse(cfg.Proxy)
	if err != nil {
		return nil, errors.New("invalid remote address")
	}
	return &RequestTransformerFactoryImpl{cfg.Headers.ValuePrefix, dest}, nil
}

func (f *RequestTransformerFactoryImpl) New(s Session) RequestTransformer {
	return &RequestTransformerImpl{f.prefix, f.url, s.Values()}
}

type RequestTransformerImpl struct {
	prefix string
	url    *url.URL
	values http.Header
}

func (t *RequestTransformerImpl) Transform(r *http.Request) {
	t.reverse(r)
	t.clear(r)
	t.fill(r)
}

func (t *RequestTransformerImpl) reverse(r *http.Request) {
	r.URL.Host = t.url.Host
	r.URL.Scheme = t.url.Scheme
}

func (t *RequestTransformerImpl) clear(r *http.Request) {
	for key := range r.Header {
		if strings.HasPrefix(key, t.prefix) {
			r.Header.Del(key)
		}
	}
}

func (t *RequestTransformerImpl) fill(r *http.Request) {
	for key, valuesOfKey := range t.values {
		for _, value := range valuesOfKey {
			r.Header.Add(t.prefix+key, value)
		}
	}
}
