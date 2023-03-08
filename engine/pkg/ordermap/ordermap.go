package ordermap

// OrderMap 有序的map
type OrderMap struct {
	keys         []string
	originalKeys []string
	m            map[string]interface{}
}

// NewOrderMap 创建一个有序的map
func NewOrderMap() *OrderMap {
	return &OrderMap{
		keys:         make([]string, 0),
		originalKeys: make([]string, 0),
		m:            make(map[string]interface{}),
	}
}

// Set 设置key value
func (om *OrderMap) Set(key string, value interface{}) {
	if _, ok := om.m[key]; ok {
		for i, k := range om.keys {
			if k == key {
				om.keys[i] = key
				break
			}
		}
		om.m[key] = value
		return
	}
	om.originalKeys = append(om.originalKeys, key)
	om.keys = append(om.keys, key)
	om.m[key] = value
}

// Get 获取key对应的value
func (om *OrderMap) Get(key string) (interface{}, bool) {
	if _, ok := om.m[key]; !ok {
		return nil, false
	}
	return om.m[key], true
}

// MustGet 获取key对应的value，如果不存在则panic
func (om *OrderMap) MustGet(key string) interface{} {
	val, ok := om.m[key]
	if !ok {
		panic("key not found")
	}
	return val
}

// Delete 删除key
func (om *OrderMap) Delete(key string) {
	if _, ok := om.m[key]; !ok {
		return
	}
	delete(om.m, key)
	for i, k := range om.keys {
		if k == key {
			om.keys = append(om.keys[:i], om.keys[i+1:]...)
			break
		}
	}
}

// Keys 获取所有的key
func (om *OrderMap) Keys() []string {
	return om.originalKeys
}

// Values 获取所有的value
func (om *OrderMap) Values() []interface{} {
	values := make([]interface{}, len(om.keys))
	for i, k := range om.keys {
		values[i] = om.m[k]
	}
	return values
}

// Len 返回map中键值对的数量
func (om *OrderMap) Len() int {
	return len(om.keys)
}

// Clear 清空map
func (om *OrderMap) Clear() {
	om.keys = make([]string, 0)
	om.originalKeys = make([]string, 0)
	om.m = make(map[string]interface{})
}

// Each 遍历
func (om *OrderMap) Each(fn func(key string, value interface{})) {
	for _, k := range om.keys {
		fn(k, om.m[k])
	}
}
