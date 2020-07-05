package headers

import (
	"net/http"
	"strings"
)

type Interceptor struct {
	header      *http.Header
	valuePrefix string
	setPrefix   string
	delPrefix   string
	toSave      *http.Header
	toSaveClean *http.Header
	toDel       *http.Header
	toDelClean  *http.Header
}

func (m *Interceptor) load() {
	for key, values := range *m.header {
		if strings.HasPrefix(key, m.setPrefix) {
			cleanKey := strings.Replace(key, m.setPrefix, "", 1)
			for _, value := range values {
				m.toSave.Add(key, value)
				m.toSaveClean.Add(cleanKey, value)
			}
		} else if strings.HasPrefix(key, m.delPrefix) {
			cleanKey := strings.Replace(key, m.delPrefix, "", 1)
			for _, value := range values {
				m.toDel.Add(key, value)
				m.toDelClean.Add(cleanKey, value)
			}
		}
	}
}

func (m *Interceptor) HeadersToSave() http.Header {
	return *m.toSaveClean
}

func (m *Interceptor) HeadersToDel() http.Header {
	return *m.toDelClean
}

func (m *Interceptor) ClearDistinctiveHeaders() {
	for k := range *m.toSave {
		m.header.Del(k)
	}
	for k := range *m.toDel {
		m.header.Del(k)
	}
}
