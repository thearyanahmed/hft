package serializer

import (
	"encoding/json"

	"github.com/hellofreshdevtests/HFtest-platform-engineering-thearyanahmed/pkg/schema"
	"github.com/thedevsaddam/govalidator"
)

type StoreConfigRequest struct {
	Name     string `json:"name" schema:"name"`
	Metadata string `json:"metadata" schema:"metadata"`
}

func (r *StoreConfigRequest) Rules() govalidator.MapData {
	return govalidator.MapData{
		"name":     []string{"required", "alpha_dash"},
		"metadata": []string{"required", "json"},
	}
}

func (r *StoreConfigRequest) ToConfigMapSchema() (schema.ConfigMap, error) {
	var data map[string]interface{}

	err := json.Unmarshal([]byte(r.Metadata), &data)
	if err != nil {
		return schema.ConfigMap{}, err
	}

	convertNestedMaps(data)

	valid := schema.ConfigMap{
		Name:     r.Name,
		Metadata: data,
	}

	return valid, nil
}

// @todo note add layer count
func convertNestedMaps(data interface{}) {
	switch v := data.(type) {
	case map[string]interface{}:
		for _, val := range v {
			if subMap, ok := val.(map[string]interface{}); ok {
				convertNestedMaps(subMap)
			}
		}
	}
}
