package refl

import (
	"fmt"
	"iter"
	"reflect"
)

// IterStructFieldsImplementingInterface returns an iterator that yields fields of the target struct that implement E
// todo: create tests
func IterStructFieldsImplementingInterface[E any](val any) iter.Seq[E] {
	iface := ReflectRoot(reflect.TypeOf((*E)(nil)))

	if iface.Kind() != reflect.Interface {
		panic(fmt.Sprintf("val must be an interface; got: %s", iface.Kind()))
	}

	rv := ReflectRoot(reflect.ValueOf(val))
	if rv.Type().Kind() != reflect.Struct {
		panic(fmt.Sprintf("val must be a struct, got: %s", rv.Type().Kind()))
	}

	// filter on interface implementing E
	flt := func(v reflect.Value) bool {
		return v.Type().Implements(iface)
	}

	return func(yield func(E) bool) {
		for v := range IterReflectStructFields(rv, flt) {
			if !yield(v.Interface().(E)) {
				return
			}
		}
	}
}

// todo: create test
func IterReflectStructFields(v reflect.Value, filter func(reflect.Value) bool) iter.Seq[reflect.Value] {
	if filter == nil {
		filter = func(reflect.Value) bool { return true }
	}

	rv := ReflectRoot(reflect.ValueOf(v.Interface()))
	return func(yield func(reflect.Value) bool) {
		for i := 0; i < rv.Type().NumField(); i++ {
			fv := ReflectRoot(rv.Field(i))

			if fv.Kind() == reflect.Invalid {
				// skip invalid values
				continue
			}

			fvIface := ReflectRoot(reflect.ValueOf(fv.Interface()))
			ft := fv.Type()

			switch true {
			case fvIface.Type().Kind() == reflect.Struct:
				IterReflectStructFields(rv.Field(i), filter)(yield)

			case filter(rv.Field(i)):
				if !yield(rv.Field(i)) {
					return
				}

			case ft.Kind() == reflect.Slice:
				IterReflectSliceValues(rv.Field(i), filter)(yield)
			}
		}

		if filter(v) {
			if !yield(v) {
				return
			}
		}
	}
}

func IterReflectSliceValues(v reflect.Value, filter func(reflect.Value) bool) iter.Seq[reflect.Value] {
	if filter == nil {
		return v.Seq()
	}

	return func(yield func(reflect.Value) bool) {
		for elem := range v.Seq() {
			if filter(elem) {
				if !yield(elem) {
					return
				}
			}
		}
	}
}

// reflectKindAndElem is an interface that allows us to unwrap pointers with ReflectRoot
// as both reflect.Value and reflect.Type use this pattern
type reflectKindAndElem[E any] interface {
	Kind() reflect.Kind
	Elem() E
}

// ReflectRoot returns the base reflect value or type by recusively unwrapping pointers
func ReflectRoot[E reflectKindAndElem[E]](t E) E {
	if t.Kind() == reflect.Ptr {
		return ReflectRoot[E](t.Elem())
	}
	return t
}
