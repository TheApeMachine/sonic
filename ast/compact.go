package ast

/*
compactGenericSlice returns a typed slice when every element shares the same
Go kind, otherwise it returns a compacted []interface{}.
*/
func compactGenericSlice(items []interface{}) interface{} {
	if len(items) == 0 {
		return items
	}

	compacted := make([]interface{}, len(items))

	for index, item := range items {
		compacted[index] = compactGenericValue(item)
	}

	if out, ok := compactHomogeneousSlice(compacted); ok {
		return out
	}

	return compacted
}

/*
compactGenericValue recursively compacts nested generic containers.
*/
func compactGenericValue(value interface{}) interface{} {
	switch typed := value.(type) {
	case []interface{}:
		return compactGenericSlice(typed)
	case map[string]interface{}:
		return compactGenericObject(typed)
	default:
		return value
	}
}

/*
compactGenericObject keeps map[string]interface{} but compacts nested values.
*/
func compactGenericObject(values map[string]interface{}) map[string]interface{} {
	if len(values) == 0 {
		return values
	}

	out := make(map[string]interface{}, len(values))

	for key, value := range values {
		out[key] = compactGenericValue(value)
	}

	return out
}

func compactHomogeneousSlice(items []interface{}) (interface{}, bool) {
	if len(items) == 0 {
		return []interface{}{}, true
	}

	switch items[0].(type) {
	case float64:
		out := make([]float64, len(items))

		for index, item := range items {
			sample, ok := item.(float64)

			if !ok {
				return nil, false
			}

			out[index] = sample
		}

		return out, true
	case string:
		out := make([]string, len(items))

		for index, item := range items {
			sample, ok := item.(string)

			if !ok {
				return nil, false
			}

			out[index] = sample
		}

		return out, true
	case bool:
		out := make([]bool, len(items))

		for index, item := range items {
			sample, ok := item.(bool)

			if !ok {
				return nil, false
			}

			out[index] = sample
		}

		return out, true
	case int:
		out := make([]int, len(items))

		for index, item := range items {
			sample, ok := item.(int)

			if !ok {
				return nil, false
			}

			out[index] = sample
		}

		return out, true
	case int64:
		out := make([]int64, len(items))

		for index, item := range items {
			sample, ok := item.(int64)

			if !ok {
				return nil, false
			}

			out[index] = sample
		}

		return out, true
	default:
		return nil, false
	}
}
