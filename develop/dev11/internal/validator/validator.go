package validator

import (
	"errors"
	"net/url"
	"reflect"
	"strings"
)

type Validator struct {
	e any
}

func NewValidator(e any) (*Validator, error) {
	ptr := reflect.TypeOf(e)

	if ptr.Kind() != reflect.Pointer {
		return nil, errors.New("Should be pointer")
	}

	return &Validator{e}, nil
}

func (v *Validator) IsFormCorrect(values url.Values) (bool, error) {
	rt := reflect.TypeOf(v.e).Elem()

	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)

		v := strings.Split(f.Tag.Get(`event`), ",")[0] // use split to ignore tag "options" like omitempty, etc.
		if v == "" {
			continue
		}
		if _, found := values[v]; !found {
			return false, nil
		}
		if len(values[v]) != 1 {
			return false, nil
		}
	}

	return true, nil
}
