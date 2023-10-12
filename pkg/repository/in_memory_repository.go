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

	// @todo remove
	configMaps = fakeData()

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
	}
	`
	var configMaps []schema.ConfigMap

	var configMap schema.ConfigMap
	err := json.Unmarshal([]byte(jsonData), &configMap)
	if err != nil {
		fmt.Println("Error:", err)
		return configMaps
	}

	for i := 0; i < 5; i++ {
		configMap.Name = fmt.Sprintf("%s-%d", configMap.Name, i)
		configMaps = append(configMaps, configMap)
	}

	return configMaps
}

func (r *InMemoryRepository) Store(configMap schema.ConfigMap) (schema.ConfigMap, error) {
	fmt.Println("new config", configMap)
	fmt.Println("len", len(r.configs))

	r.configs = append(r.configs, configMap)
	fmt.Println("after len", len(r.configs))

	return configMap, nil
}

func (r *InMemoryRepository) Find(options *schema.FilterOptions) ([]schema.ConfigMap, error) {
	result := make([]schema.ConfigMap, 0)

	count := int32(0)

	queryLength := len(options.Conditions)

	for _, config := range r.configs {
		if options.SelectAllIfConditionsAreEmpty && queryLength == 0 {
			if count == options.Limit {
				break
			}

			result = append(result, config)
			count++

		} else {

			if count == options.Limit {
				break
			}

			for k, v := range options.Conditions {
				if k == "name" && config.Name == v {
					result = append(result, config)
					count++

					continue
				}
			}
		}
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

func (r *InMemoryRepository) Delete(name string) error {
	index, found := r.findIndexByName(name)

	if !found {
		return fmt.Errorf("no resoruce found with index name %s", name)
	}

	r.configs = append(r.configs[:index], r.configs[index+1:]...)

	return nil
}

func (r *InMemoryRepository) Exists(name string) bool {
	_, found := r.findIndexByName(name)

	return found
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
