package main

import (
	"fmt"
	"testing"

	cv "github.com/smartystreets/goconvey/convey"
)

func TestSliceToList(t *testing.T) {

	cv.Convey("Given a parsable golang source file with struct containing a slice", t, func() {
		cv.Convey("then the slice should be converted to a List() in the capnp output", func() {

			ex0 := `
type s1 struct {
  MyInts []int
}`
			cv.So(ExtractString2String(ex0), ShouldStartWithModuloWhiteSpace, `struct S1Capn { myInts  @0:   List(Int64); } `)

		})
	})
}

func TestSliceOfStructToList(t *testing.T) {

	cv.Convey("Given a parsable golang source file with struct containing a slice of struct bbb", t, func() {
		cv.Convey("then the slice should be converted to a List(Bbb) in the capnp output", func() {

			ex0 := `
type bbb struct {}
type s1 struct {
  MyBees []bbb
}`
			out0 := ExtractString2String(ex0)

			fmt.Printf("out0: '%s'\n", out0)

			cv.So(out0, ShouldStartWithModuloWhiteSpace, `struct BbbCapn { } struct S1Capn { myBees  @0:   List(BbbCapn); } `)

		})
	})
}

func TestSliceOfPointerToList(t *testing.T) {

	cv.Convey("Given a parsable golang source file with struct containing a slice of pointers to struct big", t, func() {
		cv.Convey("then the slice should be converted to a List(Big) in the capnp output", func() {

			ex0 := `
type big struct {}
type s1 struct {
  MyBigs []*big
}`
			cv.So(ExtractString2String(ex0), ShouldStartWithModuloWhiteSpace, `struct BigCapn { } struct S1Capn { myBigs  @0:   List(BigCapn); } `)

		})
	})
}

func TestSliceOfByteBecomesData(t *testing.T) {

	cv.Convey("Given golang src with []byte", t, func() {
		cv.Convey("then the slice should be converted to Data, not List(Byte), in the capnp output", func() {

			ex0 := `
type s1 struct {
  MyData []byte
}`
			cv.So(ExtractString2String(ex0), ShouldStartWithModuloWhiteSpace, `struct S1Capn { myData  @0:   Data; } `)

		})
	})
}

func TestStructWithSliceOfOtherStructs(t *testing.T) {

	cv.Convey("Given a go struct containing MyBigs []Big, where Big is another struct", t, func() {
		cv.Convey("then then the CapnToGo() translation code should call the CapnToGo translation function over each slice member during translation", func() {

			in0 := `
type Big struct {}
type s1 struct {
  MyBigs []Big
}`

			expect0 := `
struct BigCapn { }
struct S1Capn { myBigs  @0:   List(BigCapn); } 

func BigCapnToGo(src BigCapn, dest *Big) *Big { 
    if dest == nil { 
      dest = &Big{} 
    }
    return dest
}
func BigGoToCapn(seg *capn.Segment, src *Big) BigCapn { 
    dest := NewBigCapn(seg)
    return dest
}   

func S1CapnToGo(src S1Capn, dest *s1) *s1 {
	if dest == nil {
		dest = &s1{}
	}
    var n int

    // MyBigs
	n = src.MyBigs().Len()
	dest.MyBigs = make([]Big, n)
	for i := 0; i < n; i++ {
        BigCapnToGo(src.MyBigs().At(i), &dest.Big[i])
    }

`
			cv.So(ExtractString2String(in0), ShouldStartWithModuloWhiteSpace, expect0)

		})
	})
}
