package repository

import (
	"testing"

	"github.com/hellofreshdevtests/HFtest-platform-engineering-thearyanahmed/pkg/schema"
	"github.com/stretchr/testify/assert"
)

func TestInMemoryRepository_Store(t *testing.T) {
	repo := NewInMemoryRepository()

	configMap := schema.ConfigMap{
		Name:     "TestConfig",
		Metadata: map[string]interface{}{"key": "value"},
	}

	storedConfig, err := repo.Store(configMap)

	assert.NoError(t, err)
	assert.Equal(t, configMap, storedConfig)
	assert.Len(t, repo.configs, 1)
}

func TestInMemoryRepository_Find(t *testing.T) {
	repo := NewInMemoryRepository()

	configMap := schema.ConfigMap{
		Name:     "TestConfig",
		Metadata: map[string]interface{}{"key": "value"},
	}

	_, _ = repo.Store(configMap)

	options := &schema.FilterOptions{
		Limit: 1,
	}

	result, err := repo.Find(options)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestInMemoryRepository_Update(t *testing.T) {
	repo := NewInMemoryRepository()

	configMap := schema.ConfigMap{
		Name:     "TestConfig",
		Metadata: map[string]interface{}{"key": "value"},
	}

	_, _ = repo.Store(configMap)

	newConfig := schema.ConfigMap{
		Name:     "UpdatedConfig",
		Metadata: map[string]interface{}{"newKey": "newValue"},
	}

	updatedConfig, err := repo.Update("TestConfig", newConfig)

	assert.NoError(t, err)
	assert.Equal(t, newConfig, updatedConfig)

	updatedConfig, err = repo.Update("NonExistentConfig", newConfig)

	assert.Error(t, err)
	assert.EqualError(t, err, "no resource found with name NonExistentConfig")
}

func TestInMemoryRepository_Delete(t *testing.T) {
	repo := NewInMemoryRepository()

	configMap := schema.ConfigMap{
		Name:     "TestConfig",
		Metadata: map[string]interface{}{"key": "value"},
	}

	_, _ = repo.Store(configMap)

	err := repo.Delete("TestConfig")

	assert.NoError(t, err)
	assert.Len(t, repo.configs, 0)

	err = repo.Delete("NonExistentConfig")

	assert.Error(t, err)
	assert.EqualError(t, err, "no resource found with name NonExistentConfig")
}

func TestInMemoryRepository_Exists(t *testing.T) {
	repo := NewInMemoryRepository()

	configMap := schema.ConfigMap{
		Name:     "TestConfig",
		Metadata: map[string]interface{}{"key": "value"},
	}

	_, _ = repo.Store(configMap)

	assert.True(t, repo.Exists("TestConfig"))
	assert.False(t, repo.Exists("NonExistentConfig"))
}

func TestInMemoryRepository_findIndexByName(t *testing.T) {
	repo := NewInMemoryRepository()

	// Add some sample data to the repository
	configMap1 := schema.ConfigMap{Name: "Config1"}
	configMap2 := schema.ConfigMap{Name: "Config2"}
	configMap3 := schema.ConfigMap{Name: "Config3"}

	repo.configs = []schema.ConfigMap{configMap1, configMap2, configMap3}

	// Test when the name is found
	index, found := repo.findIndexByName("Config2")
	assert.True(t, found)
	assert.Equal(t, 1, index)

	// Test when the name is not found
	index, found = repo.findIndexByName("NonExistentConfig")
	assert.False(t, found)
	assert.Equal(t, 0, index)
}

func TestFlatten(t *testing.T) {
	// Test case 1: When the input is a map with nested maps
	input := map[string]interface{}{
		"key1": "value1",
		"key2": map[string]interface{}{
			"nestedKey1": "nestedValue1",
			"nestedKey2": "nestedValue2",
		},
	}

	expectedOutput := map[string]interface{}{
		"key1":            "value1",
		"key2.nestedKey1": "nestedValue1",
		"key2.nestedKey2": "nestedValue2",
	}

	result := flatten(input, "")

	assert.Equal(t, expectedOutput, result)

	// Test case 2: When the input is a simple value
	simpleInput := "simpleValue"
	expectedOutput = map[string]interface{}{
		"": "simpleValue",
	}

	result = flatten(simpleInput, "")

	assert.Equal(t, expectedOutput, result)
}
