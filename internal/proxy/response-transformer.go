package proxy

import (
	"net/http"
)

// Deps

type HeadersInterceptor interface {
	HeadersToSave() http.Header
	HeadersToDel() http.Header
	ClearDistinctiveHeaders()
}

type HeadersInterceptorFactory interface {
	NewInterceptor(*http.Header) HeadersInterceptor
}

// Provider

type ResponseTransformerFactoryImpl struct {
	HeadersInterceptorFactory
}

func ProvideResponseTransformerFactory(
	headersInterceptorFactory HeadersInterceptorFactory,
) (ResponseTransformerFactory, error) {
	return &ResponseTransformerFactoryImpl{
		headersInterceptorFactory,
	}, nil
}

func (f *ResponseTransformerFactoryImpl) New(
	session Session,
) ResponseTransformer {
	return &ResponseTransformerImpl{
		f.HeadersInterceptorFactory,
		session,
	}
}

// Transformer

type ResponseTransformerImpl struct {
	HeadersInterceptorFactory
	Session Session
}

func (t *ResponseTransformerImpl) Transform(res *http.Response) error {
	t.handleHeaders(res)
	return t.checkSession()
}

func (t *ResponseTransformerImpl) handleHeaders(res *http.Response) {
	headersInterceptor := t.NewInterceptor(&res.Header)
	for key, values := range headersInterceptor.HeadersToSave() {
		t.Session.Set(key, values)
	}
	for key := range headersInterceptor.HeadersToDel() {
		t.Session.Del(key)
	}
	headersInterceptor.ClearDistinctiveHeaders()
}

func (t *ResponseTransformerImpl) checkSession() error {
	if t.Session.Changed() || t.Session.Rolling() {
		if err := t.Session.Save(); err != nil {
			return err
		}
	}

	return nil
}
