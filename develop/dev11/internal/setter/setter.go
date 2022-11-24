package setter

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type Setter struct {
}

func NewSetter() Setter {
	return Setter{}
}

func (s Setter) getTegsName(e any) (map[string]string, error) {
	ptr := reflect.TypeOf(e)

	if ptr.Kind() != reflect.Pointer {
		return nil, errors.New("Should be pointer")
	}

	rt := ptr.Elem()

	tegs := map[string]string{}

	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)

		v := strings.Split(f.Tag.Get(`event`), ",")[0] // use split to ignore tag "options" like omitempty, etc.

		if v == "" {
			continue
		}

		tegs[f.Name] = v
	}

	return tegs, nil
}

func (s *Setter) SetFields(e any, params url.Values) error {
	if reflect.TypeOf(e).Kind() != reflect.Pointer {
		return errors.New("Should be pointer")
	}

	rv := reflect.ValueOf(e)

	tegs, err := s.getTegsName(e)
	if err != nil {
		return err
	}

	elem := rv.Elem()
	for key, value := range tegs {
		f := elem.FieldByName(key)

		if f.Kind() == reflect.Int {
			num, err := strconv.Atoi(params[value][0])
			if err != nil {
				return err
			}
			if f.OverflowInt(int64(num)) {
				return errors.New("Bad value")
			}
			f.SetInt(int64(num))
		}
		if f.Kind() == reflect.String {
			f.SetString(params[value][0])
		}
		if f.Kind() == reflect.ValueOf(time.Time{}).Kind() {
			parse, err := time.Parse("2006.02.01", params[value][0])
			if err != nil {
				return err
			}
			fmt.Println(parse)
			f.Set(reflect.ValueOf(parse))
		}
	}

	return nil
}
