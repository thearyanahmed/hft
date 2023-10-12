package repository

import (
	"fmt"

	"github.com/hellofreshdevtests/HFtest-platform-engineering-thearyanahmed/pkg/schema"
)

type InMemoryRepository struct {
	configs []schema.ConfigMap
}

func NewInMemoryRepository() *InMemoryRepository {
	configMaps := make([]schema.ConfigMap, 0)

	return &InMemoryRepository{
		configs: configMaps,
	}
}

func (r *InMemoryRepository) Store(configMap schema.ConfigMap) (schema.ConfigMap, error) {
	r.configs = append(r.configs, configMap)

	return configMap, nil
}

func (r *InMemoryRepository) Find(options *schema.FilterOptions) ([]schema.ConfigMap, error) {
	result := make([]schema.ConfigMap, 0)

	if options.Limit == 0 {
		return result, nil
	}

	count := int32(0)

	queryLength := len(options.Conditions)

	for _, config := range r.configs {
		if count == options.Limit {
			break
		}

		if options.SelectAllIfConditionsAreEmpty && queryLength == 0 {
			result = append(result, config)
			count++

			continue
		}

		dataMap := make(map[string]interface{})

		dataMap["name"] = config.Name
		dataMap["metadata"] = config.Metadata

		dataMap = flatten(dataMap, "")

		matchesAllConditions := true

		for key, value := range options.Conditions {
			if val, ok := dataMap[key]; !ok || val != value {
				matchesAllConditions = false
				break
			}
		}

		if matchesAllConditions {
			result = append(result, config)
			count++
		}
	}

	return result, nil
}

func flatten(data interface{}, prefix string) map[string]interface{} {
	flattened := make(map[string]interface{})

	switch v := data.(type) {
	case map[string]interface{}:
		for key, val := range v {
			newKey := key
			if prefix != "" {
				newKey = prefix + "." + key
			}
			submap := flatten(val, newKey)
			for k, v := range submap {
				flattened[k] = v
			}
		}
	default:
		flattened[prefix] = v
	}

	return flattened
}

func (r *InMemoryRepository) Update(name string, entity schema.ConfigMap) (schema.ConfigMap, error) {
	index, found := r.findIndexByName(name)

	if !found {
		return schema.ConfigMap{}, fmt.Errorf("no resource found with name %s", name)
	}

	r.configs[index] = entity

	return entity, nil
}

func (r *InMemoryRepository) Delete(name string) error {
	index, found := r.findIndexByName(name)

	if !found {
		return fmt.Errorf("no resource found with name %s", name)
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
