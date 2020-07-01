package copycat

import (
	"reflect"
	"testing"
)

func TestDeepCopy(t *testing.T) {

	type subStruct struct {
		X string
		Y int
		Z uint
	}

	type srcStruct struct {
		A string
		B int8
		C int16
		D int32
		E int64
		F uint8
		G uint16
		H uint32
		I uint64
		J float32
		K float64
		L bool
		M *string
		N **string
		O []string
		P []int16
		Q *subStruct
		R *subStruct
		S map[string]string
		T map[string]string
	}
	type dstStruct struct {
		A string
		B int8
		C int16
		D int32
		E int64
		F uint8
		G uint16
		H uint32
		I uint64
		J float32
		K float64
		L bool
		M string
		N string
		O []string
		P [2]int32
		Q subStruct
		R *subStruct
		S map[string]string
		T *map[string]string
	}

	str := "this is a string"
	str2 := "this is another string"
	strPtr := &str2

	src := srcStruct{
		A: "abc",
		B: -1,
		C: -2,
		D: -3,
		E: -4,
		F: 1,
		G: 2,
		H: 3,
		I: 4,
		J: 0.1,
		K: 0.2,
		L: true,
		M: &str,
		N: &strPtr,
		O: []string{"a", "b", "c"},
		P: []int16{1, 2, 3},
		Q: &subStruct{X: "x", Y: 100, Z: 200},
		R: &subStruct{X: "x2", Y: 300, Z: 400},
		S: map[string]string{"v1": "k1", "v2": "k2"},
		T: map[string]string{"v3": "k3", "v4": "k4"},
	}

	expected := dstStruct{
		A: "abc",
		B: -1,
		C: -2,
		D: -3,
		E: -4,
		F: 1,
		G: 2,
		H: 3,
		I: 4,
		J: 0.1,
		K: 0.2,
		L: true,
		M: str,
		N: str2,
		O: []string{"a", "b", "c"},
		P: [2]int32{1, 2},
		Q: subStruct{X: "x", Y: 100, Z: 200},
		R: &subStruct{X: "x2", Y: 300, Z: 400},
		S: map[string]string{"v1": "k1", "v2": "k2"},
		T: &map[string]string{"v3": "k3", "v4": "k4"},
	}

	dst := &dstStruct{}
	DeepCopy(dst, src)

	if !reflect.DeepEqual(*dst, expected) {
		t.Errorf("dst != expected\ndst:\n%+v\nexpected:\n%+v", *dst, expected)
	}
}
