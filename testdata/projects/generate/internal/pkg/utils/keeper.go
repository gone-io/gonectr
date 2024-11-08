package utils

var m = make(map[string][]any)

func Put(key string, value any) {
	m[key] = append(m[key], value)
}

func Get(key string) (list []any) {
	return m[key]
}
