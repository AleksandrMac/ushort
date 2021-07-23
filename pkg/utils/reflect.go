package utils

import (
	"fmt"
	"reflect"
)

func UpdateStruct(in interface{}, values map[string]interface{}) (err error) {
	if in == nil {
		return fmt.Errorf("\"in\" is nil")
	}

	if values == nil {
		return fmt.Errorf("\"values\" is empty")
	}

	val := reflect.ValueOf(in)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return fmt.Errorf("\"in\" is not struct")
	}

	for key, v := range values {
		tmp := val.FieldByName(key)
		// если в структуке отсутствует поле с именем key то мы его пропускаем
		if tmp.Kind() == reflect.Invalid {
			continue
		}
		if tmp.Type() == nil {
			err = fmt.Errorf("field \"%v\" is not correct\n%w", key, err)
			continue
		}
		newVal := reflect.ValueOf(v)
		if tmp.Type().Kind() != newVal.Type().Kind() {
			err = fmt.Errorf("type field \"%v\" is not correct, expected %v, current: %v\n%w", key, tmp.Type().Kind(), newVal.Type().Kind(), err)
			continue
		}
		tmp.Set(newVal)
	}
	return nil
}

func FieldsFromStruct(in interface{}) ([]string, error) {
	if in == nil {
		return nil, fmt.Errorf("\"in\" is nil")
	}

	val := reflect.TypeOf(in)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("\"in\" is not struct")
	}

	out := make([]string, val.NumField())
	for i := 0; i < len(out); i++ {
		out[i] = val.Field(i).Name
	}
	return out, nil
}

func Value(in interface{}, tag, field string) (interface{}, error) {
	if in == nil {
		return nil, fmt.Errorf("\"in\" is nil")
	}

	typ := reflect.TypeOf(in)
	val := reflect.ValueOf(in)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = val.Elem()
	}

	if typ.Kind() != reflect.Struct {
		return nil, fmt.Errorf("\"in\" is not struct")
	}

	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		if f.Tag.Get(tag) == field {
			return val.Field(i).Interface(), nil
		}
	}
	return nil, fmt.Errorf("в структуре %T не найдено поле %s", in, field)
}
