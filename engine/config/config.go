package config

import (
	"reflect"
	"strings"

	"github.com/yqchilde/pkgs/log"

	"github.com/yqchilde/wxbot/engine/robot"
)

type Config map[string]any

type Plugin interface {
	// OnRegister 注册后发生
	OnRegister()
	// OnEvent 产生event发生
	OnEvent(msg *robot.Message)
}

func (c Config) Unmarshal(s any) {
	defer func() {
		if err := recover(); err != nil {
			log.Error("unmarshal error:", err)
		}
	}()
	if s == nil {
		return
	}
	var el reflect.Value
	if v, ok := s.(reflect.Value); ok {
		el = v
	} else {
		el = reflect.ValueOf(s)
	}
	if el.Kind() == reflect.Pointer {
		el = el.Elem()
	}
	t := el.Type()
	if t.Kind() == reflect.Map {
		for k, v := range c {
			el.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(v).Convert(t.Elem()))
		}
		return
	}
	//字段映射
	nameMap := make(map[string]string)
	for i, j := 0, t.NumField(); i < j; i++ {
		name := t.Field(i).Name
		tag := t.Field(i).Tag.Get("yaml")
		nameMap[tag] = name
	}
	for k, v := range c {
		name, ok := nameMap[k]
		if !ok {
			if k != "pluginmagic" {
				log.Warnf("%v plugin not found config: %v", t.Name(), k)
			}
			continue
		}
		// 需要被写入的字段
		fv := el.FieldByName(name)
		ft := fv.Type()
		// 先处理值是数组的情况
		if value := reflect.ValueOf(v); value.Kind() == reflect.Slice {
			l := value.Len()
			s := reflect.MakeSlice(ft, l, value.Cap())
			for i := 0; i < l; i++ {
				fv := value.Index(i)
				if ft == reflect.TypeOf(c) {
					fv.FieldByName("Unmarshal").Call([]reflect.Value{fv})
				} else {
					item := s.Index(i)
					if fv.Kind() == reflect.Interface {
						item.Set(reflect.ValueOf(fv.Interface()).Convert(item.Type()))
					} else {
						item.Set(fv)
					}
				}
			}
			fv.Set(s)
		} else if child, ok := v.(Config); ok { //然后处理值是递归情况（map)
			if fv.Kind() == reflect.Map {
				if fv.IsNil() {
					fv.Set(reflect.MakeMap(ft))
				}
			}
			child.Unmarshal(fv)
		} else {
			switch fv.Kind() {
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				fv.SetUint(uint64(value.Int()))
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				fv.SetInt(value.Int())
			case reflect.Float32, reflect.Float64:
				fv.SetFloat(value.Float())
			case reflect.Slice: //值是单值，但类型是数组，默认解析为一个元素的数组
				s := reflect.MakeSlice(ft, 1, 1)
				s.Index(0).Set(value)
				fv.Set(s)
			default:
				fv.Set(value)
			}
		}
	}
}

func (c *Config) Set(key string, value any) {
	if *c == nil {
		*c = Config{key: value}
	} else {
		(*c)[key] = value
	}
}

func (c Config) Get(key string) any {
	v, _ := c[key]
	return v
}

func (c Config) Has(key string) (ok bool) {
	_, ok = c[key]
	return
}

func (c Config) HasChild(key string) (ok bool) {
	_, ok = c[key].(Config)
	return ok
}

func (c Config) GetChild(key string) Config {
	if v, ok := c[strings.ToLower(key)]; ok && v != nil {
		return v.(Config)
	}
	return nil
}

func (c Config) Assign(source Config) {
	for k, v := range source {
		switch m := c[k].(type) {
		case Config:
			switch vv := v.(type) {
			case Config:
				m.Assign(vv)
			case map[string]any:
				m.Assign(vv)
			}
		default:
			c[k] = v
		}
	}
}

func Struct2Config(s any) (config Config) {
	config = make(Config)
	var t reflect.Type
	var v reflect.Value
	if vv, ok := s.(reflect.Value); ok {
		v = vv
		t = vv.Type()
	} else {
		t = reflect.TypeOf(s)
		v = reflect.ValueOf(s)
		if t.Kind() == reflect.Pointer {
			v = v.Elem()
			t = t.Elem()
		}
	}
	for i, j := 0, t.NumField(); i < j; i++ {
		ft := t.Field(i)
		if !ft.IsExported() {
			continue
		}
		name := ft.Name
		switch ft.Type.Kind() {
		case reflect.Struct:
			config[name] = Struct2Config(v.Field(i))
		case reflect.Slice:
			fallthrough
		default:
			reflect.ValueOf(config).SetMapIndex(reflect.ValueOf(name), v.Field(i))
		}
	}
	return
}
