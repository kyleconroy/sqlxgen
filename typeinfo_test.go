package dba

import (
	"reflect"
	"testing"
)

type library struct {
	Name     string `dba:"branch_name"`
	Location string
	Key      int   `dba:"key,auto"`
	Books    []int `dba:"-"`
}

func TestColumnNames(t *testing.T) {
	typ := reflect.ValueOf(library{}).Type()
	tinfo, err := getTypeInfo(typ)

	if err != nil || tinfo == nil {
		t.Fatal(err)
	}

	if len(tinfo.fields) != 3 {
		t.Fatal("tinfo should have 3 fields, not", len(tinfo.fields))
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
