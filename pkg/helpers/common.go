package helpers

import (
	"fmt"
	"strconv"
)

const subComponent = "TypeCast"

func SafeMapCast(item interface{}) map[string]interface{} {
	if item != nil {
		if casted, ok := item.(map[string]interface{}); ok {
			return casted
		}
	}
	return nil
}
func SafeSelectFromMap(item map[string]interface{}, key string) interface{} {
	if item != nil {
		if val, ok := item[key]; ok {
			return val
		}
	}
	return nil
}

func SafeCastSlice(item interface{}) []interface{} {
	if item != nil {
		if casted, ok := item.([]interface{}); ok {
			return casted
		}
	}
	return nil
}

func SafeCastString(item interface{}) string {
	if casted, ok := item.(string); ok {
		return casted
	}
	return ""
}

func SafeConvertInt(item interface{}) int {
	if str := SafeCastString(item); str != "" {
		num, err := strconv.Atoi(str)
		if err != nil {
			LogError(fmt.Sprintf("Failed to convert %v to int", item), subComponent)
			return -1
		}
		return num
	}
	LogError(fmt.Sprintf("Failed to convert %v to int", item), subComponent)
	return -1
}

func SafeConvertFloat(item interface{}) float64 {
	if str := SafeCastString(item); str != "" {
		num, err := strconv.ParseFloat(str, 32)
		if err != nil {
			LogError(fmt.Sprintf("Failed to convert %v to int", item), subComponent)
			return -1
		}
		return num
	}
	LogError(fmt.Sprintf("Failed to convert %v to int", item), subComponent)
	return -1
}

func SafeCastBool(item interface{}) bool {
	if casted, ok := item.(bool); ok {
		return casted
	}
	return false
}
