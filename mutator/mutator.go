package mutator

import (
	"net/http"

	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

var (
	deserializer runtime.Decoder
	defaulter    runtime.ObjectDefaulter
)

func init() {
	runtimeScheme := runtime.NewScheme()
	codecs := serializer.NewCodecFactory(runtimeScheme)
	deserializer = codecs.UniversalDeserializer()
}

// serve is the function that is called when a new request is made
func serve(w http.ResponseWriter, r *http.Request) {
	//-----------------------------------------
	// Decode the body
	//-----------------------------------------

	adRev, err := decodeBody(w, r)
	if err != nil {
		writeErrorResponse(w, err)
		return
	}

	//-----------------------------------------
	// Get the pod
	//-----------------------------------------

	pod, err := extractPod(adRev)
	if err != nil {
		writeErrorResponse(w, err)
		return
	}

	log.Infof("Pod to enrich has name %s, on namespace %s", pod.Name, pod.Namespace)
}
