// Copyright 2011 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package dal

import (
	"reflect"
	"strings"
	"sync"
)

// typeInfo holds details for the schema representation of a type.
type typeInfo struct {
	daltable *fieldInfo
	fields   []fieldInfo
}

// fieldInfo holds details for the column representation of a single field.
type fieldInfo struct {
	idx   []int
	name  string
	flags fieldFlags
}

type Table struct {
	Name, Encoding, Locale string
}

type fieldFlags int

const (
	fElement fieldFlags = 1 << iota
	fAuto
	fOmitEmpty

	QuerySelect = fElement
	QueryInsert = fElement | fAuto
	fMode       = fElement | fAuto | fOmitEmpty
)

var tinfoMap = make(map[reflect.Type]*typeInfo)
var tinfoLock sync.RWMutex
var tableType = reflect.TypeOf(Table{})

func getTypeInfo(typ reflect.Type) (*typeInfo, error) {
	tinfoLock.RLock()
	tinfo, ok := tinfoMap[typ]
	tinfoLock.RUnlock()
	if ok {
		return tinfo, nil
	}
	tinfo = &typeInfo{}

	if typ.Kind() == reflect.Struct && typ != tableType {
		n := typ.NumField()
		for i := 0; i < n; i++ {
			f := typ.Field(i)
			if f.PkgPath != "" || f.Tag.Get("dal") == "-" {
				continue // Private field
			}

			// FIXME: Support embedded structs

			finfo, err := structfieldInfo(typ, &f)
			if err != nil {
				return nil, err
			}
			if f.Name == "DALTable" {
				tinfo.daltable = finfo
				continue
			}
			tinfo.fields = append(tinfo.fields, *finfo)
		}
	}

	// If no DALTable field exists on the struct, set the DALTable name to be name of the struct
	if tinfo.daltable == nil {
		tinfo.daltable = &fieldInfo{
			idx:  []int{},
			name: typ.Name(),
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

	tag := f.Tag.Get("dal")

	// Parse flags.
	tokens := strings.Split(tag, ",")

	finfo.flags = fElement

	if len(tokens) > 1 {
		tag = tokens[0]
		for _, flag := range tokens[1:] {
			switch flag {
			case "auto":
				finfo.flags |= fAuto
			case "omitempty":
				finfo.flags |= fOmitEmpty
			}
		}
		// FIXME: Validate the flags used.
	}

	if f.Name == "DALTable" {
		// The DALTable field records the table name. Don't
		// process it as usual because its name should default to
		// empty rather than to the field name.
		finfo.name = tag
		return finfo, nil
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
