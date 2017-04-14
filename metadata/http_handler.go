package metadata

import (
	"encoding/json"
	"net"
	"net/http"
	"time"
)

var transport = &http.Transport{

	Proxy: http.ProxyFromEnvironment,
	DialContext: (&net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 60 * time.Second,
		DualStack: true,
	}).DialContext,
	MaxIdleConns:          1000,
	IdleConnTimeout:       60 * time.Second,
	TLSHandshakeTimeout:   10 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
	MaxIdleConnsPerHost:   1000,
}

type HttpHandler struct {
	mp PublishService
}

func NewHttpHandler(mp PublishService) *HttpHandler {
	return &HttpHandler{mp: mp}
}

func (h *HttpHandler) Publish(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var ids []Content
	err := decoder.Decode(&ids)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	errorCh := make(chan error)
	doneCh := make(chan bool)
	go h.mp.SendMetadataJob(ids, errorCh, doneCh)
	for {
		select {
		case err := <-errorCh:
			log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
		case <-doneCh:
			log.Infof("Finished importing %d contents", len(ids))
			return
		}
	}
}
