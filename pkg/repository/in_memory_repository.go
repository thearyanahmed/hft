package repository

import (
	"encoding/json"
	"fmt"

	"github.com/hellofreshdevtests/HFtest-platform-engineering-thearyanahmed/pkg/schema"
)

type InMemoryRepository struct {
	configs []schema.ConfigMap
}

func NewInMemoryRepository() *InMemoryRepository {
	configMaps := make([]schema.ConfigMap, 0)
	configMaps = fakeData()

	fmt.Println(configMaps)

	return &InMemoryRepository{
		configs: configMaps,
	}
}

func fakeData() []schema.ConfigMap {
	jsonData := `
	{
		"name": "datacenter-1",
		"metadata": {
			"monitoring": {
				"enabled": "true"
			},
			"limits": {
				"cpu": {
					"enabled": "false",
					"value": "300m"
				}
			}
		}
	}`
	var configMaps []schema.ConfigMap

	var configMap schema.ConfigMap
	err := json.Unmarshal([]byte(jsonData), &configMap)
	if err != nil {
		fmt.Println("Error:", err)
		return configMaps
	}

	configMaps = append(configMaps, configMap)

	return configMaps
}

func (r *InMemoryRepository) Store() error {
	return nil
}

func (r *InMemoryRepository) Find() ([]schema.ConfigMap, error) {
	return r.configs, nil
}

func (r *InMemoryRepository) Update() error {
	return nil
}
func (r *InMemoryRepository) Delete() error {
	return nil
}
