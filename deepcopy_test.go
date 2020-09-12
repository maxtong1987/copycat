package copycat

import (
	"fmt"
	"reflect"
	"testing"
	"time"
	"unsafe"
)

func TestDeepCopy(t *testing.T) {

	testSimple(t)

	testSlice(t)

	testArray(t)

	testPointer(t)

	testMap(t)

	structInStruct(t)

	testSliceToArray(t)

	testArrayToSlice(t)

	testPointerToConcrete(t)

	testDoublePointer(t)

	testPreserveHierarchy(t)

	testCombo(t)

	testSelfReferencing(t)

	testPrivate(t)
}

func BenchmarkDeepCopy_simple(b *testing.B) {
	src := newSimple()
	newDst := func() interface{} { return &simple{} }
	benchmarkDeepCopy(b, newDst, src)
}

func BenchmarkDeepCopy_advanceSlice(b *testing.B) {
	src := newAdvanceSlice()
	newDst := func() interface{} { return &advanceSlice{} }
	benchmarkDeepCopy(b, newDst, src)
}

func BenchmarkDeepCopy_advanceArray(b *testing.B) {
	src := newAdvanceArray()
	newDst := func() interface{} { return &advanceArray{} }
	benchmarkDeepCopy(b, newDst, src)
}

func benchmarkDeepCopy(
	b *testing.B,
	newDst func() interface{},
	src interface{},
) {
	for n := 0; n < b.N; n++ {
		dst := newDst()
		DeepCopy(dst, src)
	}
}

type simple struct {
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
	U string
}

func newSimple() *simple {
	return &simple{
		A: A,
		B: B,
		C: C,
		D: D,
		E: E,
		F: F,
		G: G,
		H: H,
		I: I,
		J: J,
		K: K,
		L: L,
		U: U,
	}
}

type advanceSlice struct {
	A []string
	B []int8
	C []int16
	D []int32
	E []int64
	F []uint8
	G []uint16
	H []uint32
	I []uint64
	J []float32
	K []float64
	L []bool
}

func newAdvanceSlice() *advanceSlice {
	return &advanceSlice{
		A: []string{A, A, A, A, A, A},
		B: []int8{B, B, B, B, B, B},
		C: []int16{C, C, C, C, C, C},
		D: []int32{D, D, D, D, D, D},
		E: []int64{E, E, E, E, E, E},
		F: []uint8{F, F, F, F, F, F},
		G: []uint16{G, G, G, G, G, G},
		H: []uint32{H, H, H, H, H, H},
		I: []uint64{I, I, I, I, I, I},
		J: []float32{J, J, J, J, J, J},
		K: []float64{K, K, K, K, K, K},
		L: []bool{L, L, L, L, L, L},
	}
}

type advanceArray struct {
	A [2]string
	B [2]int8
	C [2]int16
	D [2]int32
	E [2]int64
	F [2]uint8
	G [2]uint16
	H [2]uint32
	I [2]uint64
	J [2]float32
	K [2]float64
	L [2]bool
}

func newAdvanceArray() *advanceArray {
	return &advanceArray{
		A: [2]string{A},
		B: [2]int8{B},
		C: [2]int16{C},
		D: [2]int32{D},
		E: [2]int64{E},
		F: [2]uint8{F},
		G: [2]uint16{G},
		H: [2]uint32{H},
		I: [2]uint64{I},
		J: [2]float32{J},
		K: [2]float64{K},
		L: [2]bool{L},
	}
}

type advancePointer struct {
	A *string
	B *int8
	C *int16
	D *int32
	E *int64
	F *uint8
	G *uint16
	H *uint32
	I *uint64
	J *float32
	K *float64
	L *bool
}

type advanceDoublePointer struct {
	A **string
	B **int8
	C **int16
	D **int32
	E **int64
	F **uint8
	G **uint16
	H **uint32
	I **uint64
	J **float32
	K **float64
	L **bool
}

type advanceMap struct {
	A map[string]string
	B map[string]int8
	C map[string]int16
	D map[string]int32
	E map[string]int64
	F map[string]uint8
	G map[string]uint16
	H map[string]uint32
	I map[string]uint64
	J map[string]float32
	K map[string]float64
	L map[string]bool
}

type fooBar interface {
	Foo()
	Bar()
}

type fooBarImpl struct{}

func (fooBarImpl) Foo() {}
func (fooBarImpl) Bar() {}

type specialType struct {
	M complex64
	N complex128
	O uintptr
	P chan struct{}
	Q func()
	R fooBar
	S unsafe.Pointer
}

var (
	A = "this is a string"
	B = int8(-1)
	C = int16(-2)
	D = int32(-3)
	E = int64(-4)
	F = uint8(5)
	G = uint16(6)
	H = uint32(7)
	I = uint64(8)
	J = float32(1.11)
	K = float64(2.22)
	L = true
	O = uintptr(0)
	P = make(chan struct{})
	Q = func() { fmt.Sprintln("Hello!") }
	R = fooBarImpl{}
	S = unsafe.Pointer(nil)
	U = "this is another string"
)

func testSimple(t *testing.T) {

	src := *newSimple()

	dst := &simple{}

	expect := src

	fmt.Println("case: simple")
	DeepCopy(dst, src)

	if !reflect.DeepEqual(*dst, expect) {
		t.Errorf("dst != expect\ndst:\n%+v\nexpected:\n%+v", *dst, expect)
	}
}

func testSlice(t *testing.T) {

	src := newAdvanceSlice()

	dst := &advanceSlice{}

	expect := *src

	fmt.Println("case: slice")
	DeepCopy(dst, src)

	if !reflect.DeepEqual(*dst, expect) {
		t.Errorf("dst != expect\ndst:\n%+v\nexpected:\n%+v", *dst, expect)
	}
}

func testArray(t *testing.T) {

	src := advanceArray{
		A: [2]string{A},
		B: [2]int8{B},
		C: [2]int16{C},
		D: [2]int32{D},
		E: [2]int64{E},
		F: [2]uint8{F},
		G: [2]uint16{G},
		H: [2]uint32{H},
		I: [2]uint64{I},
		J: [2]float32{J},
		K: [2]float64{K},
		L: [2]bool{L},
	}

	dst := &advanceArray{}

	expect := src

	fmt.Println("case: array")
	DeepCopy(dst, src)

	if !reflect.DeepEqual(*dst, expect) {
		t.Errorf("dst != expect\ndst:\n%+v\nexpected:\n%+v", *dst, expect)
	}
}

func testPointer(t *testing.T) {

	src := &advancePointer{
		A: &A,
		B: &B,
		C: &C,
		D: &D,
		E: &E,
		F: &F,
		G: &G,
		H: &H,
		I: &I,
		J: &J,
		K: &K,
		L: &L,
	}

	dst := &advancePointer{}

	expect := *src

	fmt.Println("case: pointer")
	DeepCopy(dst, src)

	if !reflect.DeepEqual(*dst, expect) {
		t.Errorf("dst != expect\ndst:\n%+v\nexpected:\n%+v", *dst, expect)
	}
}

func testMap(t *testing.T) {

	src := advanceMap{
		A: map[string]string{"k": A},
		B: map[string]int8{"B": B},
		C: map[string]int16{"C": C},
		D: map[string]int32{"D": D},
		E: map[string]int64{"E": E},
		F: map[string]uint8{"F": F},
		G: map[string]uint16{"G": G},
		H: map[string]uint32{"H": H},
		I: map[string]uint64{"I": I},
		J: map[string]float32{"J": J},
		K: map[string]float64{"K": K},
		L: map[string]bool{"L": L},
	}

	dst := &advanceMap{}

	expect := src

	fmt.Println("case: map")
	DeepCopy(dst, src)

	if !reflect.DeepEqual(*dst, expect) {
		t.Errorf("dst != expect\ndst:\n%+v\nexpected:\n%+v", *dst, expect)
	}
}

func structInStruct(t *testing.T) {

	type advance struct {
		A simple
		B simple
	}

	src := advance{
		A: simple{
			A: A,
			B: B,
			C: C,
			D: D,
			E: E,
			F: F,
			G: G,
			H: H,
			I: I,
			J: J,
			K: K,
			L: L,
		},
		B: simple{
			A: A,
			B: B,
			C: C,
			D: D,
			E: E,
			F: F,
			G: G,
			H: H,
			I: I,
			J: J,
			K: K,
			L: L,
		},
	}

	dst := &advance{}

	expect := src

	fmt.Println("case: strcut in struct")
	DeepCopy(dst, src)

	if !reflect.DeepEqual(*dst, expect) {
		t.Errorf("dst != expect\ndst:\n%+v\nexpected:\n%+v", *dst, expect)
	}
}

func testSliceToArray(t *testing.T) {

	src := advanceSlice{
		A: []string{A},
		B: []int8{B},
		C: []int16{C},
		D: []int32{D},
		E: []int64{E},
		F: []uint8{F},
		G: []uint16{G},
		H: []uint32{H},
		I: []uint64{I},
		J: []float32{J},
		K: []float64{K},
		L: []bool{L},
	}

	expect := advanceArray{
		A: [2]string{A},
		B: [2]int8{B},
		C: [2]int16{C},
		D: [2]int32{D},
		E: [2]int64{E},
		F: [2]uint8{F},
		G: [2]uint16{G},
		H: [2]uint32{H},
		I: [2]uint64{I},
		J: [2]float32{J},
		K: [2]float64{K},
		L: [2]bool{L},
	}

	dst := &advanceArray{}

	fmt.Println("case: slice to array")
	DeepCopy(dst, src)

	if !reflect.DeepEqual(*dst, expect) {
		t.Errorf("dst != expect\ndst:\n%+v\nexpected:\n%+v", *dst, expect)
	}
}

func testArrayToSlice(t *testing.T) {

	src := newAdvanceArray()

	expect := advanceSlice{
		A: []string{A, ""},
		B: []int8{B, 0},
		C: []int16{C, 0},
		D: []int32{D, 0},
		E: []int64{E, 0},
		F: []uint8{F, 0},
		G: []uint16{G, 0},
		H: []uint32{H, 0},
		I: []uint64{I, 0},
		J: []float32{J, 0},
		K: []float64{K, 0},
		L: []bool{L, false},
	}

	dst := &advanceArray{}

	fmt.Println("case: array to slice")
	DeepCopy(dst, src)

	dstStr := fmt.Sprintf("%+v", *dst)
	expectedStr := fmt.Sprintf("%+v", expect)
	// can't use reflect.DeepEqual here
	if dstStr != expectedStr {
		t.Errorf("dst != expect\ndst:\n%+v\nexpected:\n%+v", *dst, expect)
	}
}

func testPointerToConcrete(t *testing.T) {
	src := &advancePointer{
		A: &A,
		B: &B,
		C: &C,
		D: &D,
		E: &E,
		F: &F,
		G: &G,
		H: &H,
		I: &I,
		J: &J,
		K: &K,
		L: &L,
	}

	dst := &simple{}

	expect := simple{
		A: A,
		B: B,
		C: C,
		D: D,
		E: E,
		F: F,
		G: G,
		H: H,
		I: I,
		J: J,
		K: K,
		L: L,
	}

	fmt.Println("case: pointer to contrect")
	DeepCopy(dst, src)

	if !reflect.DeepEqual(*dst, expect) {
		t.Errorf("dst != expect\ndst:\n%+v\nexpected:\n%+v", *dst, expect)
	}
}

func testDoublePointer(t *testing.T) {

	a := &A
	b := &B
	c := &C
	d := &D
	e := &E
	f := &F
	g := &G
	h := &H
	i := &I
	j := &J
	k := &K
	l := &L

	src := &advanceDoublePointer{
		A: &a,
		B: &b,
		C: &c,
		D: &d,
		E: &e,
		F: &f,
		G: &g,
		H: &h,
		I: &i,
		J: &j,
		K: &k,
		L: &l,
	}

	dst := &simple{}

	expect := simple{
		A: A,
		B: B,
		C: C,
		D: D,
		E: E,
		F: F,
		G: G,
		H: H,
		I: I,
		J: J,
		K: K,
		L: L,
	}

	fmt.Println("case: double pointer to concrete")
	DeepCopy(dst, src)

	if !reflect.DeepEqual(*dst, expect) {
		t.Errorf("dst != expect\ndst:\n%+v\nexpected:\n%+v", *dst, expect)
	}
}

func testPreserveHierarchy(t *testing.T) {
	type node struct {
		Name     string
		Children []*node
	}

	nodeA := &node{Name: "nodeA"}
	nodeB := &node{Name: "nodeB"}
	nodeC := &node{Name: "nodeC"}
	nodeD := &node{Name: "nodeD"}

	nodeA.Children = []*node{nodeB, nodeC}
	nodeB.Children = []*node{nodeC, nodeD}
	nodeC.Children = []*node{nodeD}

	src := nodeA
	dst := &node{}

	fmt.Println("case: PreserveHierarchy")

	DeepCopy(dst, src)

	childrenA := dst.Children
	if len(childrenA) != 2 {
		t.Errorf("len(childrenA) should be %v", 2)
		return
	}

	dstB, dstC := childrenA[0], childrenA[1]
	if dstB == nil || dstB.Name != "nodeB" {
		t.Errorf("invalid dstB: %v", dstB)
		return
	}

	if dstC == nil || dstC.Name != "nodeC" {
		t.Errorf("invalid dstC: %v", dstC)
		return
	}

	BChildren := dstB.Children
	if len(BChildren) != 2 {
		t.Errorf("len(BChildren) should be %v", 2)
		return
	}

	CChildren := dstC.Children
	if len(CChildren) != 1 {
		t.Errorf("len(CChildren) should be %v", 1)
		return
	}

	dstD := BChildren[1]
	if dstD == nil || dstD.Name != "nodeD" {
		t.Errorf("invalid dstD: %v", dstD)
		return
	}

	if !reflect.DeepEqual(BChildren[0], dstC) {
		t.Errorf("BChildren[0] != dstC: %+v  %+v", *BChildren[0], *dstC)
		return
	}

	if !reflect.DeepEqual(CChildren[0], dstD) {
		t.Errorf("CChildren[0] != dstD: %v  %v", *CChildren[0], *dstD)
		return
	}
}

func testCombo(t *testing.T) {

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
		T *map[string]string
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

	str := A
	str2 := A + " 2"
	strPtr := &str2

	src := srcStruct{
		A: A,
		B: B,
		C: C,
		D: D,
		E: E,
		F: F,
		G: G,
		H: H,
		I: I,
		J: J,
		K: K,
		L: L,
		M: &str,
		N: &strPtr,
		O: []string{"a", "b", "c"},
		P: []int16{1, 2, 3},
		Q: &subStruct{X: "x", Y: 100, Z: 200},
		R: &subStruct{X: "x2", Y: 300, Z: 400},
		S: map[string]string{"v1": "k1", "v2": "k2"},
		T: &map[string]string{"v3": "k3", "v4": "k4"},
	}

	expect := dstStruct{
		A: A,
		B: B,
		C: C,
		D: D,
		E: E,
		F: F,
		G: G,
		H: H,
		I: I,
		J: J,
		K: K,
		L: L,
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

	fmt.Println("case: combo")
	DeepCopy(dst, src)

	if !reflect.DeepEqual(*dst, expect) {
		t.Errorf("dst != expect\ndst:\n%+v\nexpected:\n%+v", *dst, expect)
	}
}

func testSelfReferencing(t *testing.T) {
	type selfReference struct {
		Me   simple
		Next *selfReference
	}
	src := &selfReference{
		Me: *newSimple(),
	}
	src.Next = src
	expect := selfReference{
		Me: *newSimple(),
	}
	expect.Next = &expect

	dst := &selfReference{}
	fmt.Println("case: self-referencing struct")

	timeout := make(chan struct{})
	go func() {
		time.Sleep(time.Second * 2)
		timeout <- struct{}{}
	}()
	done := make(chan struct{})
	go func() {
		DeepCopy(dst, src)
		done <- struct{}{}
	}()

	select {
	case <-done:
		if !reflect.DeepEqual(*dst, expect) {
			t.Errorf("dst != expect\ndst:\n%+v\nexpected:\n%+v", *dst, expect)
		}
	case <-timeout:
		t.Errorf("timeout")
	}
}

func testPrivate(t *testing.T) {
	type privateStruct struct {
		A string
		a string
		b *string
	}
	b := "b"
	src := &privateStruct{
		A: "A",
		a: "a",
		b: &b,
	}
	expect := privateStruct{
		A: "A",
	}
	dst := &privateStruct{}
	fmt.Println("case: private variables")
	DeepCopy(dst, src)

	if !reflect.DeepEqual(*dst, expect) {
		t.Errorf("dst != expect\ndst:\n%+v\nexpected:\n%+v", *dst, expect)
	}
}
