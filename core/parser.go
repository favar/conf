package core

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
)

type Parser interface {
	Get(key string, node interface{}) error
	GetNode(key string) (interface{}, error)
	Clear()
}

type parser struct {
	sync.Mutex
	source    map[string]interface{}
	converter Converter
	Nodes     map[string]reflect.Value
}

func (p *parser) Clear() {
	p.Lock()
	defer p.Unlock()
	p.Nodes = make(map[string]reflect.Value)
}

func (p *parser) GetNode(key string) (interface{}, error) {
	if key == "" {
		return nil, errors.New("key is emtpy string")
	}

	keys := strings.Split(key, ".")

	var node interface{} = nil
	node = p.source
	for _, k := range keys {
		if mp, ok := node.(map[string]interface{}); ok {
			if n, ok := mp[k]; ok {
				node = n
			} else {
				return nil, NotFoundError{fmt.Sprintf("%s not found", key)}
			}
		} else {
			return nil, fmt.Errorf("node:`%v` is not map[string]interface{}", node)
		}
	}
	return node, nil
}

func (p *parser) Get(key string, addr interface{}) (err error) {
	p.Lock()
	defer p.Unlock()

	value := reflect.ValueOf(addr)
	valueKind := value.Kind()
	if valueKind != reflect.Ptr {
		return errors.New("valueAddress is not Pointer")
	}

	if cacheElem, ok := p.Nodes[key]; ok {
		value.Elem().Set(cacheElem)
		return
	}

	var node interface{}
	node, err = p.GetNode(key)
	if err != nil {
		return err
	}
	elem := value.Elem()
	err = p.mapper(node, &elem)
	p.Nodes[key] = elem
	return err
}

func (p *parser) mapper(from interface{}, to *reflect.Value) (err error) {
	valueNode := reflect.ValueOf(from)
	switch from.(type) {
	case float64:
		err = p.converter.ConvertTo(from, *to)
	case bool:
		to.Set(valueNode)
	case string:
		to.Set(valueNode)
	case map[string]interface{}:
		kv, _ := from.(map[string]interface{})
		for k, v := range kv {
			cv := to.FieldByName(k)
			if !cv.IsValid() {
				continue
			}
			err = p.mapper(v, &cv)
			if err != nil {
				break
			}
		}
	case []interface{}:
		arr := from.([]interface{})
		newTo := reflect.MakeSlice(to.Type(), len(arr), len(arr))
		for i, ele := range arr {
			newTo.Index(i).Set(reflect.ValueOf(ele))
		}
		to.Set(newTo)
	}
	return
}

func DefaultParser(s map[string]interface{}, converter Converter) *parser {
	return &parser{
		source:    s,
		converter: converter,
		Nodes:     make(map[string]reflect.Value),
	}
}

func JsonParser(s map[string]interface{}) *parser {
	return DefaultParser(s, JsonConverter())
}
