package api

import (
	"errors"
	"reflect"
	"strings"
	"time"

	"github.com/spf13/cast"
)

type RequestParser interface {
	ParseRequest(*Req) (*Model, error)
}

type ResponseMarshaller interface {
	Body(*Model) interface{}
	Headers(*Model) map[string]string
	Status(*Model) int
}

// Model wraps model data to perform action onto
// and params used to populate the response
//
// Example:
//   Model{
//     Data: []*User{user1, user2},
//     Query: {
//       "currentUserID": 12,
//       "total": 32,
//       "limit": 5,
//       "offset": 30,
//     },
//     Response: {
//       "links": map[string]string{
//         "prev": "/users?limit=5&offset=25",
//         "first": "/users?limit=5&offset=0",
//       },
//     }
//   }
//
//   Model{
//     Data: &Article{ID: 1},
//     Query: {
//       "currentUserID": 12,
//     },
//   }
//
type Model struct {
	Data     interface{}
	Query    Meta
	Response Meta
}

type Meta map[string]interface{}

// Expose request context
func (m Meta) Set(key string, value interface{}) {
	m[strings.ToLower(key)] = value
}

func (m Meta) Del(key string) {
	delete(m, strings.ToLower(key))
}

func (m Meta) Has(key string) bool {
	_, ok := m[strings.ToLower(key)]
	return ok
}

func (m Meta) Clear() {
	for key, _ := range m {
		m.Del(key)
	}
}

func (m Meta) Get(key string) interface{} {
	key = strings.ToLower(key)
	val, ok := m[key]

	if !ok {
		return nil
	}

	switch val.(type) {
	case bool:
		return cast.ToBool(val)
	case string:
		return cast.ToString(val)
	case int64, int32, int16, int8, int:
		return cast.ToInt(val)
	case float64, float32:
		return cast.ToFloat64(val)
	case time.Time:
		return cast.ToTime(val)
	case time.Duration:
		return cast.ToDuration(val)
	case []string:
		return val
	}
	return val
}

func (m Meta) GetString(key string) string {
	return cast.ToString(m.Get(key))
}

func (m Meta) GetBool(key string) bool {
	return cast.ToBool(m.Get(key))
}

func (m Meta) GetInt(key string) int {
	return cast.ToInt(m.Get(key))
}

func (m Meta) GetFloat64(key string) float64 {
	return cast.ToFloat64(m.Get(key))
}

func (m Meta) GetTime(key string) time.Time {
	return cast.ToTime(m.Get(key))
}

func (m Meta) GetDuration(key string) time.Duration {
	return cast.ToDuration(m.Get(key))
}

func (m Meta) GetStringSlice(key string) []string {
	return cast.ToStringSlice(m.Get(key))
}

func (m Meta) GetStringMap(key string) map[string]interface{} {
	return cast.ToStringMap(m.Get(key))
}

func (m Meta) GetStringMapString(key string) map[string]string {
	return cast.ToStringMapString(m.Get(key))
}

func (m Meta) GetStringMapBool(key string) map[string]bool {
	return cast.ToStringMapBool(m.Get(key))
}

// Populate marshals the value of key into ptr
// Convenience method for a generic-type Message map
func (m Meta) Populate(key string, ptr interface{}) (err error) {
	val := m.Get(key)
	if val == nil {
		return
	}

	// don't panic, return the error
	defer func() {
		if r := recover(); r != nil {
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("Unknown panic")
			}
		}
	}()

	// lookup ptr's type
	typ := reflect.TypeOf(ptr)
	if typ.Kind() != reflect.Ptr {
		err = errors.New("Populate(key, ptr): ptr must be a pointer")
		return
	}
	typ = typ.Elem()

	// point the pointer at the new value
	ptr = reflect.ValueOf(val).Convert(typ)
	return
}
