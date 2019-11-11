package types

// PatchOperation is the JSON patch operation that must be taken
// For more information, visit: http://jsonpatch.com/
type PatchOperation struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value,omitempty"`
}
