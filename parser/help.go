package parser

import (
	"context"
	"fmt"
)

func Iter(v interface{}, iterFunc func(v interface{})) {
	if v == nil {
		return
	}

	switch item := v.(type) {
	case []interface{}:
		for _, i := range item {
			Iter(i, iterFunc)
		}

	default:
		iterFunc(v)
	}
}

func toIfaceSlice(v interface{}) []interface{} {
	if v == nil {
		return nil
	}
	return v.([]interface{})
}

func toFirst(v interface{}) interface{} {
	if v == nil {
		return nil
	}

	if vList, exist := v.([]interface{}); exist {
		return vList[0]
	}

	return v

}

func parseString(v interface{}) string {
	if v == nil {
		return ""
	}
	switch ifs := v.(type) {
	case []interface{}:
		b := make([]byte, len(ifs))
		for i, v := range ifs {
			b[i] = v.([]uint8)[0]
		}
		return string(b)
	case []byte:
		return string(ifs)
	default:
		return fmt.Sprintf("%s", ifs)
	}
}

func ResolveModel(ctx context.Context, ident string) func() *Type {
	return nil
}

func ResolveView(ctx context.Context, ident string) func() *View {
	return nil
}

var nowCount = 0

func GetCurrentCount() int {
	nowCount += 1
	return nowCount
}
