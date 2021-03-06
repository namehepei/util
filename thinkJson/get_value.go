package thinkJson

import (
	"errors"
	"reflect"
	"time"
	"util/timeUtil"
)

func (jsonObject JsonObject) GetInterface(key string) interface{} {
	return jsonObject[key]
}

func (jsonObject JsonObject) GetJsonObject(key string) (jsonObjectUnder JsonObject, err error) {
	jsonObjectUnder, ok := jsonObject[key].(map[string]interface{})
	if ok {
		return jsonObjectUnder, nil
	} else {
		return nil, ErrNotGetValue{jsonObjectUnder, key}
	}
}

func (jsonObject JsonObject) GetBool(key string) (bool, error) {
	flag, ok := jsonObject[key].(bool)
	if ok {
		return flag, nil
	} else {
		return false, errors.New("json:not get bool from " + key)
	}
}

func (jsonObject JsonObject) GetString(key string) (string, error) {
	str, ok := jsonObject[key].(string)
	if ok {
		return str, nil
	} else {
		return "", errors.New("json:not get string from " + key)
	}
}

func (jsonObject JsonObject) GetInt(key string) (int, error) {
	num, ok := jsonObject[key].(float64)
	if ok {
		return int(num), nil
	} else {
		return 0, errors.New("json:not get int from " + key)
	}
}

func (jsonObject JsonObject) GetFloat64(key string) (float64, error) {
	num, ok := jsonObject[key].(float64)
	if ok {
		return num, nil
	} else {
		return 0.0, errors.New("json:not get float from " + key)
	}
}

func (jsonObject JsonObject) GetTime(key string) (time.Time, error) {
	datetime, ok := jsonObject[key].(string)
	if ok {
		return timeUtil.GetTimeFromString(datetime), nil
	} else {
		return time.Time{}, errors.New("json:not get time.Time from " + key)
	}
}

func (jsonObject JsonObject) GetList(key string) ([]JsonObject, error) {
	list, ok := jsonObject[key].([]interface{})
	if !ok {
		return make([]JsonObject, 0), errors.New("json:not get []interface{} from " + key)
	}
	var jsonObjectSlice = make([]JsonObject, 0, len(list))
	for i := 0; i < len(list); i++ {
		jsonObjectSlice = append(jsonObjectSlice, list[i].(map[string]interface{}))
	}
	return jsonObjectSlice, nil
}

func (jsonObject JsonObject) GetStruct(ptr interface{}) {
	v := reflect.ValueOf(ptr).Elem() // the struct variable
	// 获取结构体字段
	for i := 0; i < v.NumField(); i++ {
		fieldInfo := v.Type().Field(i) // a reflect.StructField
		name := fieldInfo.Name
		// json,tag
		tag := fieldInfo.Tag
		jsonName := tag.Get("json")
		if jsonName == "" {
			jsonName = name
		}
		//
		field := v.FieldByName(name)
		switch field.Kind() {
		case reflect.String:
			str, _ := jsonObject.GetString(jsonName)
			field.SetString(str)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			num, _ := jsonObject.GetInt(jsonName)
			field.SetInt(int64(num))
		case reflect.Float64, reflect.Float32:
			num, _ := jsonObject.GetFloat64(jsonName)
			field.SetFloat(num)
		case reflect.Bool:
			flag, _ := jsonObject.GetBool(jsonName)
			field.SetBool(flag)
		default:
			field.SetPointer(nil)
		}
	}
}
