package parser2

import (
	"back/utils"
	"errors"
	"fmt"
	"slices"
)

type valType int
type valPimitiveType int

const (
	valTypeObj valType = iota
	valTypeArray

	valTypeNumber
	valTypeBool
	valTypeString
	valTypeNull
)

var primitiveTypes = []valType{valTypeNumber, valTypeBool, valTypeString, valTypeNull}

// implemented by val, and keyval
type addable interface {
	add(child any) error
}

type keyVal struct {
	KeyName string `json:"key_name"`
	Val     val    `json:"val"`
}
func (k *keyVal) add(child any) error {
	v, ok := child.(val)

	if !ok {
		return errors.New("key val node excepts only vals")
	}

	k.Val = v

	return nil
}
var _ addable = &keyVal{}

type val struct {
	// We capitalize the field names so the JSON package can see them
	ValueType valType `json:"value_type"`

	// omitempty prevents null/empty fields from cluttering your JSON
	ValObj       *valObj `json:"val_obj,omitempty"`
	ValArr       *valArr `json:"val_arr,omitempty"`
	ValPrimitive *string `json:"val_primitive,omitempty"`
}
func (v *val) add(child any) error {
	if v.ValueType == valTypeObj {
		utils.Assert(v.ValObj != nil, "object must be non nil in obj type value")

		k, ok := child.(keyVal)

		if !ok {
			return errors.New("parent node object excepts only keyvals")
		}

		v.ValObj.KeyVals = append(v.ValObj.KeyVals, k)

		return nil
	}

	if v.ValueType == valTypeArray {
		utils.Assert(v.ValArr != nil, "arr must be non nil in arr type value")

		k, ok := child.(val)

		if !ok {
			return errors.New("parent node array excepts only vals")
		}

		v.ValArr.Vals = append(v.ValArr.Vals, k)

		return nil
	}

	return errors.New("append operation is not suported on this type " + string((v.ValueType)))
}

var _ addable = &val{}

type valArr struct {
	Vals []val `json:"vals"`
}

type valObj struct {
	KeyVals []keyVal `json:"key_vals"`
}

type valsParams interface {
	string | int64 | float64 | bool
}

func createPrimitiveVal[V valsParams](valType valType, v V) (val, error) {
	utils.Assert(slices.Contains(primitiveTypes, valType),
		"create primitive accepts only primitive types")

	str := fmt.Sprintf("%v", v)

	return val{
		ValueType:    valType,
		ValPrimitive: &str,
	}, nil

}

func createArrVal() val {
	return val{
		ValueType: valTypeArray,
		ValArr: &valArr{
			Vals: make([]val, 0),
		},
	}
}

func createObjVal() val {
	return val{
		ValueType: valTypeObj,
		ValObj: &valObj{
			KeyVals: make([]keyVal, 0),
		},
	}
}

func createKeyVal(keyName string) keyVal {
	return keyVal{
		KeyName: keyName,
	}
}
