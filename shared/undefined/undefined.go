package undefined

import (
	"database/sql"
	"database/sql/driver"
	"encoding"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

// Verify interface implementations
var (
	_ json.Marshaler           = (*Undefined[any])(nil)
	_ json.Unmarshaler         = (*Undefined[any])(nil)
	_ encoding.TextUnmarshaler = (*Undefined[any])(nil)
	_ driver.Valuer            = (*Undefined[any])(nil)
	_ sql.Scanner              = (*Undefined[any])(nil)
)

// Undefined is a generic wrapper type that can represent a value that may be explicitly
// undefined or unset, which is particularly useful for:
//   - JSON marshaling/unmarshaling where fields can be omitted
//   - Database operations where NULL values need to be distinguished from zero values
//
// Supported types for T include:
//   - Basic types: int64, float64, bool, string
//   - time.Time for timestamp handling
//   - []byte for binary data
//   - nil for explicit NULL values in databases
//
// The zero value of Undefined[T] is considered undefined (valid = false).
type Undefined[U any] struct {
	value U
	valid bool
}

func New[U any](value U) Undefined[U] {
	return Undefined[U]{
		value: value,
		valid: true,
	}
}

func (u Undefined[U]) Get() U {
	return u.value
}

func (u *Undefined[U]) Set(value U) {
	u.value = value
	u.valid = true
}

func (u *Undefined[U]) Unset() {
	var v U
	u.value = v
	u.valid = false
}

// Implement json.Unmarshaler
func (u *Undefined[U]) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &u.value); err != nil {
		return err
	}

	u.valid = true
	return nil
}

// Implement json.Marshaler
func (u Undefined[U]) MarshalJSON() ([]byte, error) {
	data, err := json.Marshal(u.value)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Implement encoding.TextUnmarshaler
func (u *Undefined[U]) UnmarshalText(text []byte) error {
	u.valid = len(text) > 0
	if textUnmarshaler, ok := any(&u.value).(encoding.TextUnmarshaler); ok {
		if err := textUnmarshaler.UnmarshalText(text); err != nil {
			return err
		}
		u.valid = true
		return nil
	}

	return errors.New("Undefined: cannot unmarshal text: underlying value doesn't implement encoding.TextUnmarshaler")
}

// Implement driver.Valuer
func (u Undefined[U]) Value() (driver.Value, error) {
	if !u.valid {
		return nil, nil
	}

	if valuer, ok := any(u.value).(driver.Valuer); ok {
		v, err := valuer.Value()
		return v, err
	}
	return u.value, nil
}

// Implement sql.Scanner
func (u *Undefined[U]) Scan(src any) error {
	u.valid = true

	if src == nil {
		var zero U
		u.value = zero
		return nil
	}

	if val, ok := src.(U); ok {
		u.value = val
		return nil
	}

	val := reflect.ValueOf(src)
	typ := reflect.TypeOf((*U)(nil)).Elem()

	if typ.Kind() == reflect.Interface && val.Type().Implements(typ) {
		u.value = val.Interface().(U)
		return nil
	}
	if val.Type().ConvertibleTo(typ) {
		u.value = val.Convert(typ).Interface().(U)
		return nil
	}

	u.valid = false
	return fmt.Errorf("Undefined: Scan() incompatible types (src: %T, dst: %T)", src, u.value)
}

func (u Undefined[U]) IsUndefined() bool {
	return !u.valid
}

func (u Undefined[U]) Ptr() *U {
	if u.valid {
		return &u.value
	}

	return nil
}

func (u Undefined[U]) Equal(other Undefined[U]) bool {
	if u.valid != other.valid {
		return false
	}
	if !u.valid {
		return true
	}

	return reflect.DeepEqual(u.value, other.value)
}
