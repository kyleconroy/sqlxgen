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
	DALTable Table  `dal:"libraries"`
	Name     string `dal:"branch_name"`
	Location string
	Key      int    `dal:"key,auto"`
	Books    []book `dal:"-"`
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

func TestColumnNames(t *testing.T) {
	typ := reflect.ValueOf(library{}).Type()
	tinfo, err := getTypeInfo(typ)

	if err != nil || tinfo == nil {
		t.Fatal(err)
	}

	if len(tinfo.fields) != 3 {
		t.Fatal("tinfo should have 2 fields, not", len(tinfo.fields))
	}

	if tinfo.fields[0].name != "branch_name" {
		t.Errorf("tinfo[0] should be branch_name, not %+v", tinfo.fields[1])
	}

	if tinfo.fields[1].name != "Location" {
		t.Errorf("tinfo[1] should be Location, not %+v", tinfo.fields[1])
	}

	if tinfo.fields[2].name != "key" {
		t.Errorf("tinfo[2] should be Location, not %+v", tinfo.fields[2])
	}

}

//func TestValues(t *testing.T) {
//	_, columns, _ := Expression(&library{}, QuerySelect)
//
//	if !reflect.DeepEqual(columns, []string{"libraries.branch_name libraries.Location libraries.key"}) {
//		t.Errorf("%s is wrong", columns)
//	}
//
//}
