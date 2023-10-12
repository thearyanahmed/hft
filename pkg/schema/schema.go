package schema

type ConfigMap struct {
	Name     string                 `json:"name"`
	Metadata map[string]interface{} `json:"metadata"`
}

type FilterOptions struct {
	Limit      int32
	Conditions map[string]string

	// If set to true, it will return all values if Conditions is empty.
	// Therefore, it'll not match anything and just select all the values until the given Limit
	SelectAllIfConditionsAreEmpty bool
}
