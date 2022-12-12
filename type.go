package egorm

import "reflect"

func NotContainZeroValues(value any) bool {
	return !ContainZeroValues(value)
}

// ContainZeroValues
// Checks whether the argument contains zero value or not.
// If the argument is a struct, it is determined if the inner field // contains zero value or not.
// If there is a pointer in the inner field, the zero value is determined based on whether the pointer is nil or not.
// And whether it is a zero value or not is not considered. (This is similar to the behavior of json.)
// For other than struct, it conforms to the isZero() function of reflect.
func ContainZeroValues(value any) bool {
	if value == nil {
		return true
	}
	rv := reflect.ValueOf(value)
	// structの場合, ゼロ値を含むかどうかを確認する
	switch rv.Kind() {
	case reflect.Struct:
		return containZeroValuesInStruct(rv)
	default:
		return rv.IsZero()
	}
}

// containZeroValuesInStruct
// Determine if zero value exists in the internal field of struct
// If even one zero value exists, return true.
// If the inner field is a pointer, the zero value is determined by whether the pointer is nil or not.
// (ex) In case of field: str *string, str: nil -> isZero:true, str: "" -> isZero:false,
func containZeroValuesInStruct(rv reflect.Value) bool {

	for i := 0; i < rv.NumField(); i++ {
		fieldValue := rv.Field(i)
		if fieldValue.IsZero() {
			return true
		}
	}
	return false
}
