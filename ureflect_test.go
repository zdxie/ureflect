package reflect

import (
	"ureflect/testdata"
	"reflect"
	"testing"
)

func TestStructOrFuncName(t *testing.T) {
	name := StructOrFuncName(testdata.Function1)
	if name != "testdata.Function1" {
		t.Errorf("Name error: expect testdata.Function1 but get: %s", name)
		return
	}
	name = StructOrFuncName(testdata.Function2)
	if name != "testdata.Function2" {
		t.Errorf("Name error: expect testdata.Function2 but get: %s", name)
		return
	}
	name = StructOrFuncName(testdata.Stc1{})
	if name != "testdata.Stc1" {
		t.Errorf("Name error: expect testdata.Stc1 but get: %s", name)
		return
	}
	name = StructOrFuncName(&testdata.Stc1{})
	if name != "testdata.Stc1" {
		t.Errorf("Name error: expect testdata.Stc1 but get: %s", name)
		return
	}
	name = StructOrFuncName(testdata.Stc2{}.S)
	if name != "testdata.Stc3" {
		t.Errorf("Name error: expect testdata.Stc3 but get: %s", name)
		return
	}
}

var (
	goldArray1 = []int{}
	goldArray2 = [2]string{}
)

func TestArrayType(t *testing.T) {
	typ := ArrayType(goldArray1)
	if typ != reflect.Int {
		t.Errorf("ArrayType error: expect int but get %s", typ)
		return
	}
	typ = ArrayType(goldArray2)
	if typ != reflect.String {
		t.Errorf("ArrayType error: expect string but get %s", typ)
		return
	}
}

var (
	goldMap1 = map[string]int{}
	goldMap2 = map[string]string{}
)

func TestMapType(t *testing.T) {
	ktyp, vtyp := MapType(goldMap1)
	if vtyp != reflect.Int {
		t.Errorf("MapType error: value type: expect int but get %s", vtyp)
		return
	}
	if ktyp != reflect.String {
		t.Errorf("MapType error: key type: expect string but get %s", ktyp)
		return
	}
	ktyp, vtyp = MapType(goldMap2)
	if vtyp != reflect.String {
		t.Errorf("MapType error: value type: expect string but get %s", vtyp)
		return
	}
	if ktyp != reflect.String {
		t.Errorf("MapType error: key type: expect string but get %s", ktyp)
		return
	}
}

type ForHaveTag struct {
	t1 string  `test:"tt1"`
	t2 int64   `test:"tt2"`
	t3 int     `test:"tt3"`
	t4 float64 `test:"tt4"`
	t5 string
}

func TestHaveTag(t *testing.T) {
	for _, tagInfo := range HaveTag(&ForHaveTag{}, "test") {
		switch tagInfo.FieldName {
		case "t1":
			if tagInfo.TagContent != "tt1" {
				t.Errorf("Expect tt1 but get %s", tagInfo.TagContent)
				return
			}
		case "t2":
			if tagInfo.TagContent != "tt2" {
				t.Errorf("Expect tt2 but get %s", tagInfo.TagContent)
				return
			}
		case "t3":
			if tagInfo.TagContent != "tt3" {
				t.Errorf("Expect tt3 but get %s", tagInfo.TagContent)
				return
			}
		case "t4":
			if tagInfo.TagContent != "tt4" {
				t.Errorf("Expect tt4 but get %s", tagInfo.TagContent)
				return
			}
		default:
			t.Errorf("Don't expect %s", tagInfo.FieldName)
		}
	}
}

type stc2ForClone struct {
	c int
}

type stc1ForClone struct {
	a int
	b string
	c stc2ForClone
	d *stc2ForClone
}

func TestClone(t *testing.T) {
	goldenStc := &stc1ForClone{
		a: 1,
		b: "test",
		c: stc2ForClone{
			c: 2,
		},
		d: &stc2ForClone{
			c: 3,
		},
	}
	clone := Clone(goldenStc)
	if !reflect.DeepEqual(clone, goldenStc) {
		t.Errorf("Expect equal struct: %v but get: %v", goldenStc, clone)
		return
	}
}

type stcForValueByKey struct {
	F1 int
	F2 string `reflect:"ff2"`
	F3 []int
	F4 int `reflect:"-"`
}

func TestValueByKey(t *testing.T) {
	var (
		goldenMap1 = map[string]string{
			"t1": "test1",
			"t2": "test2",
			"t3": "test3",
		}
		goldenStc1 = stcForValueByKey{
			F1: 1,
			F2: "test f2",
			F3: []int{1, 2, 3},
		}
		invalid = 1
	)
	_, ok, err := ValueByKey(goldenStc1, "F1")
	if err != nil {
		t.Errorf("ValueByKey error: %v", err)
		return
	}
	if !ok {
		t.Errorf("Expect F1 but get none")
		return
	}
	_, ok, err = ValueByKey(goldenStc1, "ff2")
	if err != nil {
		t.Errorf("ValueByKey error: %v", err)
		return
	}
	if !ok {
		t.Errorf("Expect ff2 but get none")
		return
	}
	_, ok, err = ValueByKey(goldenStc1, "F3")
	if err != nil {
		t.Errorf("ValueByKey error: %v", err)
		return
	}
	if !ok {
		t.Errorf("Expect F3 but get none")
		return
	}
	_, ok, err = ValueByKey(goldenStc1, "F4")
	if err != nil {
		t.Errorf("ValueByKey error: %v", err)
		return
	}
	if ok {
		t.Errorf("Ignore F4 but get the value")
		return
	}
	_, ok, err = ValueByKey(goldenMap1, "t1")
	if err != nil {
		t.Errorf("ValueByKey error: %v", err)
		return
	}
	if !ok {
		t.Errorf("Expect t1 but get none")
		return
	}
	_, ok, err = ValueByKey(invalid, "ff")
	if err == nil {
		t.Errorf("Expect error type but get right")
		return
	}
}
