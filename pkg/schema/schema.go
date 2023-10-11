package schema

type ConfigMap struct {
	Name     string                 `json:"name"`
	Metadata map[string]interface{} `json:"metadata"`
}
