package mutator

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
	v1beta1 "k8s.io/api/admission/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// decodeBody checks if the request body is in the correct format and return
// the admission review
func decodeBody(w http.ResponseWriter, request *http.Request) (*v1beta1.AdmissionReview, error) {
	var decodedBody []byte
	adRev := v1beta1.AdmissionReview{}

	if request.Body == nil {
		return nil, errors.New("Body is nil")
	}

	defer request.Body.Close()

	// Read it
	_decodedBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		io.Copy(ioutil.Discard, request.Body)
		return nil, fmt.Errorf("An error occurred while reading body of the request: %v", err)
	}

	decodedBody = _decodedBody

	// No body?
	if len(decodedBody) == 0 {
		http.Error(w, "Body is empty", http.StatusBadRequest)
		return nil, errors.New("Body is empty")
	}

	// Is the content-type correct?
	contentType := request.Header.Get("Content-Type")
	if contentType != "application/json" {
		http.Error(w, "Invalid Content-Type. Accepted: application/json", http.StatusUnsupportedMediaType)
		return nil, errors.New("Invalid Content-Type. Accepted: application/json")
	}

	// Deserialize it
	if _, _, err := deserializer.Decode(decodedBody, nil, &adRev); err != nil {
		return nil, fmt.Errorf("An error occurred while de-coding object: %v", err)
	}

	if adRev.Request == nil {
		return nil, errors.New("The request in the admission review is nil")
	}

	return &adRev, nil
}

// extractPod gets the pod object from the admission review
func extractPod(adRev *v1beta1.AdmissionReview) (*corev1.Pod, error) {
	req := adRev.Request
	var pod corev1.Pod

	// Decode it
	if err := json.Unmarshal(req.Object.Raw, &pod); err != nil {
		log.Errorf("An error occurred while decoding pod from the admission review : %v", err)
		return nil, fmt.Errorf("An error occurred while decoding pod from the admission review : %v", err)
	}

	return &pod, nil
}

// writeErrorResponse is just a convenient function to writeResponse when
// having to write a response
func writeErrorResponse(w http.ResponseWriter, err error) {
	log.Errorln(err.Error())

	adRev := &v1beta1.AdmissionReview{
		Response: &v1beta1.AdmissionResponse{
			Result: &metav1.Status{
				Message: err.Error(),
			},
		},
	}

	writeResponse(w, adRev)
}

// writeResponse writes the http response
func writeResponse(w http.ResponseWriter, adRev *v1beta1.AdmissionReview) {
	// First, encode it to json
	resp, err := json.Marshal(adRev)
	if err != nil {
		log.Errorf("An error occurred while encoding response: %v", err)
		http.Error(w, fmt.Sprintf("An error occurred while encoding response: %v", err), http.StatusInternalServerError)
	}

	// Write the response
	if _, err := w.Write(resp); err != nil {
		log.Errorf("An error occurred while writing response: %v", err)
		http.Error(w, fmt.Sprintf("An error occurred while writing response: %v", err), http.StatusInternalServerError)
	}
}

// buildResponse builds the complete admission review (i.e. with the response)
// to be sent as a response
func buildResponse() (*v1beta1.AdmissionReview, error) {
	pt := v1beta1.PatchTypeJSONPatch
	admissionReviewResp := v1beta1.AdmissionReview{}
	admissionReviewResp.Response = &v1beta1.AdmissionResponse{
		//UID:       reqUID,
		Allowed:   true,
		Patch:     marshalledPatchOps,
		PatchType: &pt,
	}

	// NOTE: the documentation says that Patch should be base64-ed, but I found
	// it to be working like this anyway. So, for now, I am leaving it
	// like this.
	return &admissionReviewResp, nil
}

func MarshalPatchOperations() {
	// Marshal the patch operations
	patchBytes, err := json.Marshal(patchOps)
	if err != nil {
		log.Fatalf("An error occurred while trygin to marshal the sidecars: %v", err)
	}

	marshalledPatchOps = patchBytes
}
