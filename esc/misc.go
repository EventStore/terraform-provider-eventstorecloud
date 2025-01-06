package esc

// In terraform when we read back values they will always be of type []interface{}
// even if we passed a []string originally. This takes []interface{} and builds
// a []string by casting each element individually.
func interfaceToStringList(value interface{}) []string {
	list := value.([]interface{})
	result := []string{}
	for _, element := range list {
		result = append(result, element.(string))
	}
	return result
}

// In terraform when we read back values they will always be of type []interface{}
// even if we passed a []map[string]interface{} originally. This takes []interface{}
// and builds a []map[string]interface{} by casting each element individually.
func interfaceToMapList(value interface{}) []map[string]interface{} {
	list := value.([]interface{})
	result := []map[string]interface{}{}
	for _, element := range list {
		aMap, ok := element.(map[string]interface{})
		if !ok {
			aMap = map[string]interface{}{}
		}
		result = append(result, aMap)
	}
	return result
}
