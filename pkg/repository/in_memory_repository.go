package repository

import (
	"encoding/json"
	"fmt"

	"github.com/hellofreshdevtests/HFtest-platform-engineering-thearyanahmed/pkg/schema"
)

type InMemoryRepository struct {
	configs []schema.ConfigMap
}

type FilterOptions struct {
	Limit     int32
	Query     map[string]string // metadata.allergens.eggs=true
	SelectAll bool
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

func (r *InMemoryRepository) Store(configMap schema.ConfigMap) error {
	r.configs = append(r.configs, configMap)

	return nil
}

func (r *InMemoryRepository) Find() ([]schema.ConfigMap, error) {
	return r.configs, nil
}

func (r *InMemoryRepository) Update() error {
	name := "some-name"
	entity := schema.ConfigMap{}

	index, found := r.findIndexByName(name)

	if !found {
		return fmt.Errorf("no resoruce found with index name %s", name)
	}

	r.configs[index] = entity

	return nil
}

// returns index, foundOrNot
func (r *InMemoryRepository) findIndexByName(name string) (int, bool) {
	for i, config := range r.configs {
		if config.Name == name {
			return i, true
		}
	}

	return 0, false
}

func (r *InMemoryRepository) Delete() error {
	name := "hello"

	index, found := r.findIndexByName(name)

	if !found {
		return fmt.Errorf("no resoruce found with index name %s", name)
	}

	r.configs = append(r.configs[:index], r.configs[index+1:]...)

	return nil
}
