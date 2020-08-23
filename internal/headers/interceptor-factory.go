package headers

import (
	"net/http"

	"github.com/the-bailiff/bailiff/internal/config"
	"github.com/the-bailiff/bailiff/internal/proxy"
)

type InterceptorFactory struct {
	ValuePrefix string
	SetPrefix   string
	DelPrefix   string
}

func ProvideInterceptorFactory(cfg config.Config) proxy.HeadersInterceptorFactory {
	return &InterceptorFactory{
		ValuePrefix: cfg.Headers.ValuePrefix,
		SetPrefix:   cfg.Headers.SetPrefix,
		DelPrefix:   cfg.Headers.DelPrefix,
	}
}

func (f *InterceptorFactory) NewInterceptor(h *http.Header) proxy.HeadersInterceptor {
	i := &Interceptor{
		header:      h,
		valuePrefix: f.ValuePrefix,
		setPrefix:   f.SetPrefix,
		delPrefix:   f.DelPrefix,
		toSave:      make(http.Header),
		toSaveClean: make(http.Header),
		toDel:       make(http.Header),
		toDelClean:  make(http.Header),
	}

	i.load()

	return i
}
