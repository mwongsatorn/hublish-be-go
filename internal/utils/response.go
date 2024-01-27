package utils

import (
	"maps"
	"reflect"
)

func ResponseOmitFilter(st interface{}, omitFields []string) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	v := reflect.ValueOf(st)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		tag := field.Tag.Get("json")

		if contains(omitFields, tag) {
			continue
		}

		if field.Anonymous && field.Type.Kind() == reflect.Struct {
			embeddedResult, err := ResponseOmitFilter(v.Field(i).Interface(), omitFields)
			if err != nil {
				return nil, err
			}
			maps.Copy(result, embeddedResult)
		} else {
			result[tag] = v.Field(i).Interface()
		}
	}

	return result, nil
}

func contains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}
