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

type fieldFlags int

const (
	fElement fieldFlags = 1 << iota
	fAuto
	fOmitEmpty

	fMode = fElement | fAuto | fOmitEmpty
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

			// For embedded structs, embed its fields.
			//if f.Anonymous {
			//	t := f.Type
			//	if t.Kind() == reflect.Ptr {
			//		t = t.Elem()
			//	}
			//	if t.Kind() == reflect.Struct {
			//		inner, err := getTypeInfo(t)
			//		if err != nil {
			//			return nil, err
			//		}
			//		if tinfo.daltable == nil {
			//			tinfo.daltable = inner.daltable
			//		}
			//		for _, finfo := range inner.fields {
			//			finfo.idx = append([]int{i}, finfo.idx...)
			//			if err := addfieldInfo(typ, tinfo, &finfo); err != nil {
			//				return nil, err
			//			}
			//		}
			//		continue
			//	}
			//}

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
	if len(tokens) == 1 {
		finfo.flags = fElement
	} else {
		tag = tokens[0]
		for _, flag := range tokens[1:] {
			switch flag {
			case "auto":
				finfo.flags |= fAuto
			case "omitempty":
				finfo.flags |= fOmitEmpty
			}
		}

		// Validate the flags used.
		//valid := true
		//switch mode := finfo.flags & fMode; mode {
		//case 0:
		//	finfo.flags |= fElement
		//case fAttr, fCharData, fInnerXml, fComment, fAny:
		//	if f.Name == "XMLName" || tag != "" && mode != fAttr {
		//		valid = false
		//	}
		//default:
		//	// This will also catch multiple modes in a single field.
		//	valid = false
		//}
		//if finfo.flags&fMode == fAny {
		//	finfo.flags |= fElement
		//}
		//if finfo.flags&fOmitEmpty != 0 && finfo.flags&(fElement|fAttr) == 0 {
		//	valid = false
		//}
		//if !valid {
		//	return nil, fmt.Errorf("xml: invalid tag in field %s of type %s: %q",
		//		f.Name, typ, f.Tag.Get("xml"))
		//}
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
