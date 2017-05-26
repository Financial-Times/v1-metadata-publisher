package metadata

import (
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
	"fmt"
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
	mp   PublishService
	size int
}

func NewHttpHandler(mp PublishService, size int) *HttpHandler {
	return &HttpHandler{mp: mp, size: size}
}

func (h *HttpHandler) Publish(w http.ResponseWriter, r *http.Request) {
	// decoder := json.NewDecoder(r.Body)
	// var ids []string
	var contents []Content
	// err := decoder.Decode(&ids)
	// if err != nil {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return
	// }
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	bodyString := string(body)
	ids := strings.Split(bodyString, ",")
	for _, id := range ids {
		fmt.Println(id)
		contents = append(contents, Content{UUID: id})
	}
	defer r.Body.Close()

	errorCh := make(chan error)
	doneCh := make(chan bool)
	var b []Content
	var prog int
	for _, c := range contents {
		prog++
		b = append(b, c)
		if prog%h.size == 0 {
			go h.mp.SendMetadataJob(b, errorCh, doneCh)
			wait(errorCh, doneCh)
			b = []Content{}
		}
	}

	if len(b) > 0 {
		go h.mp.SendMetadataJob(b, errorCh, doneCh)
		wait(errorCh, doneCh)
	}

	// for {
	// 	select {
	// 	case err := <-errorCh:
	// 		log.Error(err)
	// 		// w.WriteHeader(http.StatusInternalServerError)
	// 	case <-doneCh:
	// 		log.Infof("Finished importing %d contents", len(ids))
	// 		return
	// 	}
	// }
}
