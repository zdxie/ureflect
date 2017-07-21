// Copyright 2017 The Package Author zdxie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Implements some utility functions base on reflect
package reflect

import (
	"fmt"
	"reflect"
	"regexp"
	"runtime"
	"strings"
)

var namesKind = map[string]reflect.Kind{
	"invalid":        reflect.Invalid,
	"bool":           reflect.Bool,
	"int":            reflect.Int,
	"int8":           reflect.Int8,
	"int16":          reflect.Int16,
	"int32":          reflect.Int32,
	"int64":          reflect.Int64,
	"uint":           reflect.Uint,
	"uint8":          reflect.Uint8,
	"uint16":         reflect.Uint16,
	"uint32":         reflect.Uint32,
	"uint64":         reflect.Uint64,
	"uintptr":        reflect.Uintptr,
	"float32":        reflect.Float32,
	"float64":        reflect.Float64,
	"complex64":      reflect.Complex64,
	"complex128":     reflect.Complex128,
	"array":          reflect.Array,
	"chan":           reflect.Chan,
	"func":           reflect.Func,
	"interface":      reflect.Interface,
	"map":            reflect.Map,
	"ptr":            reflect.Ptr,
	"slice":          reflect.Slice,
	"string":         reflect.String,
	"struct":         reflect.Struct,
	"unsafe.Pointer": reflect.UnsafePointer,
}

var (
	arrayRegexp = regexp.MustCompile(`\s*\[.*\]\s*`)
	mapRegexp   = regexp.MustCompile(`\s*map\[|\]\s*`)
)

const (
	mapKeyIndex   = 1
	mapValueIndex = 2
)

// Gets element of the pointer obj points to.
func elem(obj reflect.Value) reflect.Value {
	for obj.Kind() == reflect.Ptr {
		obj = obj.Elem()
	}
	return obj
}

// Returns the type of interface obj or that the pointer obj points to.
func TypeOf(obj interface{}) reflect.Type {
	t := reflect.TypeOf(obj)
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

// Returns the value of interface obj or that the pointer obj points to.
func ValueOf(obj interface{}) reflect.Value {
	v, ok := obj.(reflect.Value)
	if ok {
		return v
	} else {
		v = reflect.ValueOf(obj)
	}
	return elem(v)
}

// Returns the field value of interface obj or that the pointer obj points to.
func Field(obj interface{}, field string) reflect.Value {
	f := ValueOf(obj).FieldByName(field)
	return elem(f)
}

// Gets name of struct or function.
func StructOrFuncName(obj interface{}) string {
	v := ValueOf(obj)
	var t reflect.Type
	if v.Kind() != reflect.Invalid {
		t = v.Type()
	} else {
		t = TypeOf(obj)
	}
	var name string
	switch t.Kind() {
	case reflect.Func:
		name = runtime.FuncForPC(v.Pointer()).Name()
	case reflect.Struct:
		name = strings.Join([]string{t.PkgPath(), t.Name()}, ".")
	default:
		return ""
	}
	return name[strings.LastIndex(name, "/")+1:]
}

func kind(name string) reflect.Kind {
	k, ok := namesKind[name]
	if !ok {
		k = reflect.Invalid
	}
	return k
}

// Gets elements type of array or slice.
func ArrayType(array interface{}) reflect.Kind {
	name := arrayRegexp.ReplaceAllString(ValueOf(array).Type().String(), "")
	return kind(name)
}

// Gets types of key and value of map.
func MapType(array interface{}) (reflect.Kind, reflect.Kind) {
	strs := mapRegexp.Split(ValueOf(array).Type().String(), -1)
	return kind(strs[mapKeyIndex]), kind(strs[mapValueIndex])
}

type FieldWithTag struct {
	Field      reflect.Value
	FieldName  string
	TagContent string
}

// Gets the content of tagName and field related information.
func HaveTag(obj interface{}, tagName string) []*FieldWithTag {
	// obj is a struct or reflect.value
	v := ValueOf(obj)
	t := v.Type()
	tagInfoSlice := []*FieldWithTag{}
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		tagContent := field.Tag.Get(tagName)
		if tagContent != "" {
			tagInfoSlice = append(tagInfoSlice, &FieldWithTag{
				Field:      v.Field(i),
				FieldName:  field.Name,
				TagContent: tagContent,
			})
		}
	}
	return tagInfoSlice
}

// Clones a struct.
func Clone(obj interface{}) interface{} {
	clone := reflect.New(TypeOf(obj)).Elem()
	clone.Set(ValueOf(obj))
	return clone.Addr().Interface()
}

// Converts struct to map and supports reflect tag.
func StructToMap(stc interface{}) map[string]interface{} {
	v := ValueOf(stc)
	t := v.Type()
	result := make(map[string]interface{})
	for i := 0; i < v.NumField(); i++ {
		result[t.Field(i).Name] = elem(v.Field(i)).Interface()
	}
	fields := HaveTag(v, "reflect")
	for _, field := range fields {
		if field.TagContent != "-" {
			result[field.TagContent] = field.Field.Interface()
		} else {
			delete(result, field.FieldName)
		}
	}
	return result
}

// Gets value of map or struct by key.
func ValueByKey(obj interface{}, key string) (interface{}, bool, error) {
	v := ValueOf(obj)
	t := v.Type()
	var value reflect.Value
	switch t.Kind() {
	case reflect.Struct:
		ret, ok := StructToMap(obj)[key]
		if !ok {
			return nil, false, nil
		}
		return ret, true, nil
	case reflect.Map:
		value = v.MapIndex(ValueOf(key))
	default:
		return nil, false, fmt.Errorf("Invalid object type.")
	}
	if !value.IsValid() {
		return nil, false, nil
	}
	return value.Interface(), true, nil
}
