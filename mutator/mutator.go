package mutator

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

func serve(w http.ResponseWriter, r *http.Request) {
	log.Infof("Got request from host %v, path %v, header %v, uri %v, rawquery %v", r.Host, r.URL.Path, r.Header, r.URL.String(), r.URL.RawQuery)
}
