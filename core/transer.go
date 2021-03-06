package core

import (
	"errors"
	"fmt"
	"reflect"
)

type Converter interface {
	ConvertTo(data interface{}, value reflect.Value) error
}

type jsonBasicConverter struct {
}

func (j jsonBasicConverter) ConvertTo(d interface{}, value reflect.Value) (err error) {
	data, ok := d.(float64)
	if !ok {
		return errors.New("parameter d is not kind of float64")
	}

	var v reflect.Value
	switch value.Kind() {
	case reflect.Uint:
		v = reflect.ValueOf(uint(data))
	case reflect.Uint8:
		v = reflect.ValueOf(uint8(data))
	case reflect.Uint16:
		v = reflect.ValueOf(uint16(data))
	case reflect.Uint32:
		v = reflect.ValueOf(uint32(data))
	case reflect.Uint64:
		v = reflect.ValueOf(uint64(data))
	case reflect.Int:
		v = reflect.ValueOf(int(data))
	case reflect.Int8:
		v = reflect.ValueOf(int8(data))
	case reflect.Int16:
		v = reflect.ValueOf(int16(data))
	case reflect.Int32:
		v = reflect.ValueOf(int32(data))
	case reflect.Int64:
		v = reflect.ValueOf(int64(data))
	case reflect.Float32:
		v = reflect.ValueOf(float32(data))
	case reflect.Float64:
		v = reflect.ValueOf(data)
	default:
		err = fmt.Errorf("converter not support type %s", value.Kind().String())
	}
	value.Set(v)
	return err
}

func JsonConverter() Converter {
	return jsonBasicConverter{}
}
