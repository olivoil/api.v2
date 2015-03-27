package api

type GetterSetter interface {
	Getter
	Setter
}

// Value getter
type Getter interface {
	Get(string) string
	GetAll(string) []string
}

// Value setter
type Setter interface {
	Set(key string, value string)
	Append(key string, value ...string)
}

// Converts a map or url.Values into a rest.Getter interface
type Values map[string][]string

func (v Values) Get(key string) string {
	if v == nil {
		return ""
	}
	vs, ok := v[key]
	if !ok || len(vs) == 0 {
		return ""
	}
	return vs[0]
}

func (v Values) GetAll(key string) []string {
	if v == nil {
		return []string{}
	}
	vs, ok := v[key]
	if !ok || len(vs) == 0 {
		return []string{}
	}
	return vs
}

func (v Values) Set(key string, value string) {
	v[key] = []string{value}
}

func (v Values) Append(key string, value ...string) {
	vs, ok := v[key]
	if !ok || len(vs) == 0 {
		v[key] = value
		return
	}
	v[key] = append(vs, value...)
}
