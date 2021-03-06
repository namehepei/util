package thinkJson

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func TestGetJsonObject(t *testing.T) {
	//fmt.Println(MustGetJsonObject(nil))

	m := map[string]interface{}{
		"test":  1,
		"hello": "hello",
		"str":   "ok",
		"m": map[string]interface{}{
			"num": 1,
		},
	}

	for i := 0; i < 3; i++ {
		fmt.Println(string(MustMarshal(m)))
	}
	fmt.Println(MapJsonBySortKey(m, nil))
	data, _ := json.Marshal(m)

	fmt.Println(MustGetJsonObject(data)["m"])
	//fmt.Println(GetJsonObject(data).TransMapStringString())
	//fmt.Println(MustGetJsonObject(data).MustGetJsonObject("m").MustGetString("num"))
}

func happenErr() (s map[string]string, err error) {
	return s, ErrNotGetValue{reflect.TypeOf(s), "hello"}

}
