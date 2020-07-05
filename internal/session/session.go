package session

import (
	"net/http"

	gs "github.com/gorilla/sessions"
)

type SessionImpl struct {
	changed bool
	rolling bool
	writer  *http.ResponseWriter
	request *http.Request
	session *gs.Session
}

func (s *SessionImpl) Values() http.Header {
	result := http.Header{}
	for rawKey, rawValue := range s.session.Values {
		if key, ok := rawKey.(string); ok {
			if values, ok := rawValue.([]string); ok {
				result[key] = values
			}
		}
	}
	return result
}

func (s *SessionImpl) Set(key string, values []string) {
	s.changed = true
	s.session.Values[key] = values
}

func (s *SessionImpl) Del(key string) {
	s.changed = true
	delete(s.session.Values, key)
}

func (s *SessionImpl) Save() error {
	return s.session.Save(s.request, *s.writer)
}

func (s *SessionImpl) Rolling() bool {
	return s.rolling
}

func (s *SessionImpl) Changed() bool {
	return s.changed
}
