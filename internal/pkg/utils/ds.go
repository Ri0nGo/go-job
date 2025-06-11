package utils

// RemoveDuplicate 切片去重
func RemoveDuplicate[T comparable](slice []T) []T {
	seen := make(map[T]struct{})
	result := make([]T, 0)
	for _, val := range slice {
		if _, ok := seen[val]; !ok {
			seen[val] = struct{}{}
			result = append(result, val)
		}
	}
	return result
}

// PtrToVal 接收一个指针，转为对应的值类型
// 若接收的是空指针，则转成对应类型的零值
func PtrToVal[T any](val *T) T {
	var result T
	if val == nil {
		return result
	}
	return *val
}
