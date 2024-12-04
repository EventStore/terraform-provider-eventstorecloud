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
