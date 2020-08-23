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
	toSave      http.Header
	toSaveClean http.Header
	toDel       http.Header
	toDelClean  http.Header
}

func (i *Interceptor) load() {
	for key, values := range *i.header {
		if strings.HasPrefix(key, i.setPrefix) {
			cleanKey := strings.Replace(key, i.setPrefix, "", 1)
			for _, value := range values {
				i.toSave.Add(key, value)
				i.toSaveClean.Add(cleanKey, value)
			}
		} else if strings.HasPrefix(key, i.delPrefix) {
			cleanKey := strings.Replace(key, i.delPrefix, "", 1)
			for _, value := range values {
				i.toDel.Add(key, value)
				i.toDelClean.Add(cleanKey, value)
			}
		}
	}
}

func (i *Interceptor) HeadersToSave() http.Header {
	return i.toSaveClean
}

func (i *Interceptor) HeadersToDel() http.Header {
	return i.toDelClean
}

func (i *Interceptor) ClearDistinctiveHeaders() {
	for k := range i.toSave {
		i.header.Del(k)
	}
	for k := range i.toDel {
		i.header.Del(k)
	}
}
