package schema

type ConfigMap struct {
	Name     string                 `json:"name"`
	Metadata map[string]interface{} `json:"metadata"`
}

type FilterOptions struct {
	Limit     int32
	Query     map[string]string // metadata.allergens.eggs=true
	SelectAll bool
}
