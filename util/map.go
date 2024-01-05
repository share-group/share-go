package util

type mapUtil struct{}

var MapUtil = newMapUtil()

func newMapUtil() *mapUtil {
	return &mapUtil{}
}

// GetValueFromMap 从映射中获取值，如果键不存在则返回默认值
func (s *mapUtil) GetValueFromMap(m map[string]any, key string, defaultValue any) any {
	if m == nil {
		return defaultValue
	}
	if val, ok := m[key]; ok {
		if val == nil {
			return defaultValue
		}
		return val
	}
	return defaultValue
}

// ContainsKey 判断映射中是否包含指定键
func (s *mapUtil) ContainsKey(m map[string]any, key string) bool {
	if _, ok := m[key]; ok {
		return true
	}
	return false
}
