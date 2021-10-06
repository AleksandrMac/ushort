package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type utilsStructTest struct {
	FieldString string `json:"field_string"`
	FieldInt    int
	Slice       []int
	Object      struct {
		NestedField int
	}
}

func TestUpdateStruct(t *testing.T) {
	str := utilsStructTest{
		FieldString: "stroka",
		FieldInt:    107,
		Slice:       []int{112, 107, 207},
		Object:      struct{ NestedField int }{NestedField: 302},
	}

	v := map[string]interface{}{
		"FieldString": "NewString",
		"FieldInt":    555,
		"Slice":       []int{1, 2, 3},
		"Object":      struct{ NestedField int }{NestedField: 777},
	}

	expectStr := utilsStructTest{
		FieldString: "NewString",
		FieldInt:    555,
		Slice:       []int{1, 2, 3},
		Object:      struct{ NestedField int }{NestedField: 777},
	}

	err := UpdateStruct(&str, v)
	assert.NoError(t, err, "TestUpdateStruct")
	assert.Equal(t, expectStr, str)
}

func TestFieldsFromStruct(t *testing.T) {
	str := utilsStructTest{}
	expect := []string{"FieldString", "FieldInt", "Slice", "Object"}

	slice, _ := FieldsFromStruct(str)
	assert.Equal(t, expect, slice)
}
