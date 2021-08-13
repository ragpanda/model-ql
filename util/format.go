package util

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/spf13/cast"
)

func Display(data interface{}) string {
	dataValue := reflect.Indirect(reflect.ValueOf(data))
	if !dataValue.IsValid() {
		return fmt.Sprintf("%+v", data)
	}

	if u8t, ok := data.([]uint8); ok {
		return fmt.Sprintf("(%T)%s", data, string(u8t))
	}

	kind := dataValue.Type().Kind()
	if kind == reflect.Struct || kind == reflect.Slice || kind == reflect.Map {
		if dataValue.Type() == reflect.TypeOf(sync.Map{}) {
			syncMap := data.(sync.Map)
			showStr := []string{}
			syncMap.Range(func(key, value interface{}) bool {
				showStr = append(showStr, fmt.Sprintf("\"%v\":%s", key, Display(value)))
				return true
			})
			return fmt.Sprintf("{%s}", strings.Join(showStr, ","))
		}

		b, _ := json.Marshal(data)
		result := string(b)
		if result == "" {
			result = fmt.Sprintf("%+v", data)
		}

		return fmt.Sprintf("(%T)%s", data, result)

	} else {
		result, err := cast.ToStringE(dataValue.Interface())
		if err != nil {
			return fmt.Sprintf("(%T)%s", data, dataValue.Interface())
		}
		return fmt.Sprintf("(%T)%s", data, result)
	}

}
