package dba

import (
	"reflect"
	"strings"
	"sync"
)

// typeInfo holds details for the schema representation of a type.
type typeInfo struct {
	fields []fieldInfo
}

// fieldInfo holds details for the column representation of a single field.
type fieldInfo struct {
	idx   []int
	name  string
	flags fieldFlags
}

type fieldFlags int

const (
	fElement fieldFlags = 1 << iota
	fOmitEmpty

	fMode = fElement | fOmitEmpty
)

var tinfoMap = make(map[reflect.Type]*typeInfo)
var tinfoLock sync.RWMutex

func getTypeInfo(typ reflect.Type) (*typeInfo, error) {
	tinfoLock.RLock()
	tinfo, ok := tinfoMap[typ]
	tinfoLock.RUnlock()
	if ok {
		return tinfo, nil
	}
	tinfo = &typeInfo{}

	if typ.Kind() == reflect.Struct {
		n := typ.NumField()
		for i := 0; i < n; i++ {
			f := typ.Field(i)
			if f.PkgPath != "" || f.Tag.Get("dba") == "-" {
				continue // Private field
			}

			// FIXME: Support embedded structs

			finfo, err := structfieldInfo(typ, &f)
			if err != nil {
				return nil, err
			}
			tinfo.fields = append(tinfo.fields, *finfo)
		}
	}

	tinfoLock.Lock()
	tinfoMap[typ] = tinfo
	tinfoLock.Unlock()
	return tinfo, nil
}

// structfieldInfo builds and returns a fieldInfo for f.
func structfieldInfo(typ reflect.Type, f *reflect.StructField) (*fieldInfo, error) {
	finfo := &fieldInfo{idx: f.Index}

	tag := f.Tag.Get("dba")

	// Parse flags.
	tokens := strings.Split(tag, ",")

	finfo.flags = fElement

	if len(tokens) > 1 {
		tag = tokens[0]
		for _, flag := range tokens[1:] {
			switch flag {
			case "omitempty":
				finfo.flags |= fOmitEmpty
			}
		}
		// FIXME: Validate the flags used.
	}

	if tag == "" {
		// If the name part of the tag is completely empty, get
		// default from or field name otherwise.
		finfo.name = f.Name
		return finfo, nil
	}

	finfo.name = tag
	return finfo, nil
}
