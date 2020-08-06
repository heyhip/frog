package frog

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strconv"
	"strings"
	"unicode"
)

// 结构体转map
func StructToMap(obj interface{}) map[string]interface{} {
	objType := reflect.TypeOf(obj)
	objValue := reflect.ValueOf(obj)

	m := make(map[string]interface{})

	n := objType.NumField()
	for i := 0; i < n; i++ {
		if objType.Field(i).Type.Kind() == reflect.Struct {
			tmp := StructToMap(objValue.Field(i).Interface())
			for k, v := range tmp {
				m[Camel2Case(k)] = v
			}
		} else {
			m[Camel2Case(objType.Field(i).Name)] = objValue.Field(i).Interface()
		}
	}

	return m
}

// 将多维map的key处理为小写带下划线
func MapCamel2Case(data interface{}) (map[string]interface{}, bool) {

	info := make(map[string]interface{})
	if m, ok := data.(map[string]interface{}); ok {
		for k, v := range m {
			if reflect.TypeOf(v).Kind() == reflect.Map {
				m, ok = MapCamel2Case(v.(map[string]interface{}))
				if ok {
					info[Camel2Case(k)] = m
				}
			} else if reflect.TypeOf(v).Kind() == reflect.Slice {
				info[Camel2Case(k)] = sliceCamelToMap(v.([]interface{}))
			} else {
				info[Camel2Case(k)] = v
			}
		}
	}
	return info, true
}

// 处理slice
func sliceCamelToMap(data []interface{}) []interface{} {
	s := make([]interface{}, 0)
	for _, v := range data {
		if reflect.TypeOf(v).Kind() == reflect.Map {
			tmp, ok := MapCamel2Case(v)
			if ok {
				s = append(s, tmp)
			}
		} else if reflect.TypeOf(v).Kind() == reflect.Slice {
			s = append(s, sliceCamelToMap(v.([]interface{})))
		} else {
			s = append(s, v)
		}
	}
	return s
}

// 将多维小写带下划线map的key处理为驼峰
func MapCase2Camel(data interface{}) (map[string]interface{}, bool) {

	info := make(map[string]interface{})
	if m, ok := data.(map[string]interface{}); ok {
		for k, v := range m {
			if reflect.TypeOf(v).Kind() == reflect.Map {
				m, ok = MapCase2Camel(v.(map[string]interface{}))
				if ok {
					info[Case2Camel(k)] = m
				}
			} else if reflect.TypeOf(v).Kind() == reflect.Slice {
				info[Case2Camel(k)] = sliceCaseToMap(v.([]interface{}))
			} else {
				info[Case2Camel(k)] = v
			}
		}
	}
	return info, true
}

// 处理slice
func sliceCaseToMap(data []interface{}) []interface{} {
	s := make([]interface{}, 0)
	for _, v := range data {
		if reflect.TypeOf(v).Kind() == reflect.Map {
			tmp, ok := MapCase2Camel(v)
			if ok {
				s = append(s, tmp)
			}
		} else if reflect.TypeOf(v).Kind() == reflect.Slice {
			s = append(s, sliceCaseToMap(v.([]interface{})))
		} else {
			s = append(s, v)
		}
	}
	return s
}

// 驼峰转小写带下划线
func Camel2Case(str string) string {
	buffer := bytes.Buffer{}
	for i, r := range str {
		if unicode.IsUpper(r) {
			if i != 0 {
				buffer.WriteString("_")
			}
			buffer.WriteRune(unicode.ToLower(r))
		} else {
			buffer.WriteRune(r)
		}
	}
	return buffer.String()
}

// 小写带下划线转驼峰
func Case2Camel(str string) string {
	str = strings.Replace(str, "_", " ", -1)
	str = strings.Title(str)
	return strings.Replace(str, " ", "", -1)
}

/***************************/
func MapToStruct(m map[string]interface{}, s interface{}) ([]byte, error, bool) {
	mp := getMapByStructInterface(m, s)
	j, e := json.Marshal(mp)
	if e != nil {
		return []byte(nil), e, false
	}
	return j, nil, true
}

func getMapByStructInterface(m map[string]interface{}, s interface{}) map[string]interface{} {

	maps := make(map[string]interface{})

	st := reflect.TypeOf(s)
	sv := reflect.ValueOf(s)

	stNum := st.NumField()
	for i := 0; i < stNum; i++ {

		// 结构体参数名称
		name := st.Field(i).Name
		// 结构体参数类型
		t := st.Field(i).Type

		value, ok := m[name]
		if !ok {
			if t.Kind() != reflect.Struct {
				continue
			}
		}

		if t.Kind() == reflect.Struct {
			tmp := getMapByStructInterface(m, sv.Field(i).Interface())
			for k, v := range tmp {
				maps[k] = v
			}
		} else {
			if reflect.TypeOf(value).Kind() == reflect.String {
				maps[name] = getTypeByReflectType(t, value.(string))
			} else {
				maps[name] = value
			}
		}

	}

	return maps
}

func MapStringToStruct(m map[string]string, s interface{}) ([]byte, error, bool) {
	mp := getMapByStructString(m, s)
	j, e := json.Marshal(mp)
	if e != nil {
		return []byte(nil), e, false
	}
	return j, nil, true
}

func getMapByStructString(m map[string]string, s interface{}) map[string]interface{} {

	maps := make(map[string]interface{})

	st := reflect.TypeOf(s)
	sv := reflect.ValueOf(s)

	stNum := st.NumField()
	for i := 0; i < stNum; i++ {

		// 结构体参数名称
		name := st.Field(i).Name
		// 结构体参数类型
		t := st.Field(i).Type

		value, ok := m[name]
		if !ok {
			if t.Kind() != reflect.Struct {
				continue
			}
		}

		if t.Kind() == reflect.Struct {
			tmp := getMapByStructString(m, sv.Field(i).Interface())
			for k, v := range tmp {
				maps[k] = v
			}
		} else {
			maps[name] = getTypeByReflectType(t, value)
		}

	}

	return maps
}

/***************************/

func getTypeByReflectType(rt reflect.Type, s string) (i interface{}) {
	switch rt.Kind() {
	case reflect.String:
		i = s
	case reflect.Int64:
		i, _ = strconv.ParseInt(s, 10, 64)
	case reflect.Int:
		i, _ = strconv.Atoi(s)
	case reflect.Int8:
		it, _ := strconv.ParseInt(s, 10, 64)
		i = int8(it)
	case reflect.Int16:
		it, _ := strconv.ParseInt(s, 10, 64)
		i = int16(it)
	case reflect.Int32:
		it, _ := strconv.ParseInt(s, 10, 64)
		i = int32(it)
	case reflect.Bool:
		it, _ := strconv.ParseBool(s)
		i = it
	case reflect.Uint:
		it, _ := strconv.ParseUint(s, 10, 64)
		i = uint(it)
	case reflect.Uint8:
		it, _ := strconv.ParseUint(s, 10, 64)
		i = uint8(it)
	case reflect.Uint16:
		it, _ := strconv.ParseUint(s, 10, 64)
		i = uint16(it)
	case reflect.Uint32:
		it, _ := strconv.ParseUint(s, 10, 64)
		i = uint32(it)
	case reflect.Uint64:
		it, _ := strconv.ParseUint(s, 10, 64)
		i = it
	case reflect.Float32:
		it, _ := strconv.ParseFloat(s, 64)
		i = float32(it)
	case reflect.Float64:
		it, _ := strconv.ParseFloat(s, 64)
		i = it
	}
	return
}
