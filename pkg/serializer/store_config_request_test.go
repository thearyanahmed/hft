package serializer

import (
	"testing"

	"github.com/hellofreshdevtests/HFtest-platform-engineering-thearyanahmed/pkg/schema"
	"github.com/stretchr/testify/assert"
)

func TestStoreConfigRequest_Rules(t *testing.T) {
	request := StoreConfigRequest{}
	rules := request.Rules()

	assert.Contains(t, rules, "name")
	assert.Contains(t, rules, "metadata")
}

func TestStoreConfigRequest_ToConfigMapSchema(t *testing.T) {
	jsonData := `{"key1": "value1", "key2": {"nestedKey": "nestedValue"}}`
	request := StoreConfigRequest{
		Name:     "TestConfig",
		Metadata: jsonData,
	}

	expectedConfigMap := schema.ConfigMap{
		Name: "TestConfig",
		Metadata: map[string]interface{}{
			"key1": "value1",
			"key2": map[string]interface{}{
				"nestedKey": "nestedValue",
			},
		},
	}

	configMap, err := request.ToConfigMapSchema()

	assert.NoError(t, err)
	assert.Equal(t, expectedConfigMap, configMap)
}

func TestConvertNestedMaps(t *testing.T) {
	// Test case 1: Nested map with one level of nesting
	input := map[string]interface{}{
		"key1": "value1",
		"key2": map[string]interface{}{
			"nestedKey": "nestedValue",
		},
	}

	convertNestedMaps(input)

	expectedOutput := map[string]interface{}{
		"key1": "value1",
		"key2": map[string]interface{}{
			"nestedKey": "nestedValue",
		},
	}

	assert.Equal(t, expectedOutput, input)
}
