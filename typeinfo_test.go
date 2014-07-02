package dal

import (
	"reflect"
	"testing"
)

type book struct {
}

type empty struct {
	DALTable Table
}

type library struct {
	DALTable Table `dal:"libraries"`
}

func TestTableName(t *testing.T) {

	tests := map[string]interface{}{
		"":          empty{},
		"book":      book{},
		"libraries": library{},
	}

	for k, v := range tests {
		typ := reflect.ValueOf(v).Type()
		tinfo, err := getTypeInfo(typ)

		if err != nil || tinfo == nil {
			t.Error(err)
			continue
		}

		if tinfo.daltable == nil || tinfo.daltable.name != k {
			t.Errorf("%s != %s", tinfo.daltable, k)
		}
	}
}
