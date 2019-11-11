package mutator

import (
	"net/http"

	log "github.com/sirupsen/logrus"
	v1beta1 "k8s.io/api/admission/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"polycube.network/polycube-sidecar-injector/types"
)

var (
	deserializer runtime.Decoder
	defaulter    runtime.ObjectDefaulter
	patchOps     []types.PatchOperation
)

func init() {
	runtimeScheme := runtime.NewScheme()
	codecs := serializer.NewCodecFactory(runtimeScheme)
	deserializer = codecs.UniversalDeserializer()
	patchOps = []types.PatchOperation{}
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

	log.Infof("Pod to enrich has name %s, on namespace %s", pod.ObjectMeta.GenerateName, adRev.Request.Namespace)

	if !requiresMutation(pod, adRev) {
		adRev.Response = &v1beta1.AdmissionResponse{
			Allowed: true,
		}

		writeRespone(w, adRev)
	}
}

// requiresMutation checks if the mutation is actually needed
func requiresMutation(pod *corev1.Pod, adRev *v1beta1.AdmissionReview) bool {
	om := pod.ObjectMeta

	// Check the namespace.
	// NOTE: I still haven't understood why, but the namespace in
	// pod.ObjectMeta is empty. So we must take it from the admission
	// review itself
	namespace := om.Namespace
	if len(namespace) == 0 {
		namespace = adRev.Request.Namespace
	}

	// Namespace is ok?
	if namespace == "kube-system" || namespace == "kube-public" {
		log.Infof("Pod %s is namespace %s, so it will be skipped", om.GenerateName, adRev.Request.Namespace)
		return false
	}

	// Has the appropriate label?
	val, exists := om.Annotations["polycube.network/sidecar"]
	if !exists {
		log.Infof("Pod %s does not have the annotation polycube.network/sidecar, so it will be skipped", om.GenerateName)
		return false
	}
	if val != "enabled" {
		log.Infof("Pod %s has annotation polycube.network/sidecar but the value is not recognized (%s), so it will be skipped", om.GenerateName, val)
		return false
	}

	// Already done?
	val, exists = om.Annotations["polycube.network/sidecar-status"]
	if exists && val == "injected" {
		log.Infof("Pod %s already has annotation polycube.network/sidecar-status=injected, so it will be skipped", om.GenerateName)
		return false
	}

	return true
}
