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

func (r *InMemoryRepository) Store(configMap schema.ConfigMap) error {
	r.configs = append(r.configs, configMap)

	return nil
}

func (r *InMemoryRepository) Find(options *schema.FilterOptions) ([]schema.ConfigMap, error) {
	result := make([]schema.ConfigMap, 0)

	count := int32(0)

	for _, config := range r.configs {
		fmt.Printf("count %d options.limit %d\n", count, options.Limit)
		if !options.SelectAll && count == options.Limit {
			break
		}

		result = append(result, config)
		count++
	}

	return result, nil
}

func (r *InMemoryRepository) Update(entity schema.ConfigMap) error {
	index, found := r.findIndexByName(entity.Name)

	if !found {
		return fmt.Errorf("no resoruce found with index name %s", entity.Name)
	}

	r.configs[index] = entity

	return nil
}

func (r *InMemoryRepository) Delete(entity schema.ConfigMap) error {
	index, found := r.findIndexByName(entity.Name)

	if !found {
		return fmt.Errorf("no resoruce found with index name %s", entity.Name)
	}

	r.configs = append(r.configs[:index], r.configs[index+1:]...)

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
