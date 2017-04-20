package db

import (
	"reflect"
)

// Merge .
func Merge(src, dst interface{}, force bool) error {

	const EPS = 1e-12

	typSrc := reflect.TypeOf(src)
	typDst := reflect.TypeOf(dst)
	if !reflect.DeepEqual(typSrc, typDst) {
		panic("not same type")
	}

	valSrc := reflect.ValueOf(src)
	valDst := reflect.ValueOf(dst)
	valSrc = reflect.Indirect(valSrc)
	valDst = reflect.Indirect(valDst)

	for i := 0; i < valSrc.NumField(); i++ {
		fSrc := valSrc.Field(i)
		fDst := valDst.Field(i)

		should := force

		if !should {
			switch fSrc.Kind() {
			case reflect.Int:
				fallthrough
			case reflect.Int16:
				fallthrough
			case reflect.Int32:
				fallthrough
			case reflect.Int64:
				n1 := fSrc.Int()
				n2 := fDst.Int()
				if n1 > 0 && n2 <= 0 {
					should = true
				}
				break
			case reflect.Uint:
				fallthrough
			case reflect.Uint16:
				fallthrough
			case reflect.Uint32:
				fallthrough
			case reflect.Uint64:
				n1 := fSrc.Uint()
				n2 := fDst.Uint()
				if n1 > 0 && n2 <= 0 {
					should = true
				}
				break

			case reflect.Float32:
				fallthrough
			case reflect.Float64:
				n1 := fSrc.Float()
				n2 := fDst.Float()
				if n1 > EPS && n2 < EPS {
					should = true
				}
				break
			case reflect.String:
				s1 := fSrc.String()
				s2 := fDst.String()
				if s1 != "" && s2 != "" {
					should = true
				}
				break
			}
		}

		if should {
			fDst.Set(fSrc)
		}
	}
	return nil
}

/*
type point struct {
	X int32
	Y int32
	Z int32
}

func main() {
	a := point{
		X: 7,
		Y: 8,
		Z: 9,
	}

	b := point{
		X: 1,
		Z: 2,
	}

	fmt.Printf("a: %+v\n", a)
	fmt.Printf("b: %+v\n", b)
	Merge(&a, &b, true)
	fmt.Printf("c: %+v\n", b)
}
*/
