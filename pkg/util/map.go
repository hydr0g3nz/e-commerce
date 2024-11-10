package util

func MapDeleteNilOrZero(m map[string]interface{}) {
	for k, v := range m {
		if v == nil || IsZeroValue(v) {
			delete(m, k)
		}
	}
}

// func MapDeleteEmpty(m map[string]interface{}) map[string]interface{} {
// 	for k, v := range m {
// 		if v == "" {
// 			delete(m, k)
// 		}
// 	}
// 	return m
// }
