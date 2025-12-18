package runtime

import "fmt"

// ArrayType represents the type of a primitive array
type ArrayType int

const (
	ArrayTypeBoolean ArrayType = 4
	ArrayTypeChar    ArrayType = 5
	ArrayTypeFloat   ArrayType = 6
	ArrayTypeDouble  ArrayType = 7
	ArrayTypeByte    ArrayType = 8
	ArrayTypeShort   ArrayType = 9
	ArrayTypeInt     ArrayType = 10
	ArrayTypeLong    ArrayType = 11
)

// Array represents a JVM array (primitive or reference)
type Array struct {
	Type       ArrayType // For primitive arrays
	ClassName  string    // For reference arrays (ANEWARRAY)
	Length     int32     // Array length
	Ints       []int32   // For int/boolean/byte/char/short arrays
	Longs      []int64   // For long arrays
	Floats     []float32 // For float arrays
	Doubles    []float64 // For double arrays
	References []any     // For reference arrays
}

// NewPrimitiveArray creates a new primitive array
func NewPrimitiveArray(atype ArrayType, length int32) *Array {
	arr := &Array{
		Type:   atype,
		Length: length,
	}

	switch atype {
	case ArrayTypeBoolean, ArrayTypeByte, ArrayTypeChar, ArrayTypeShort, ArrayTypeInt:
		arr.Ints = make([]int32, length)
	case ArrayTypeLong:
		arr.Longs = make([]int64, length)
	case ArrayTypeFloat:
		arr.Floats = make([]float32, length)
	case ArrayTypeDouble:
		arr.Doubles = make([]float64, length)
	}

	return arr
}

// NewReferenceArray creates a new reference array
func NewReferenceArray(className string, length int32) *Array {
	return &Array{
		ClassName:  className,
		Length:     length,
		References: make([]any, length),
	}
}

// GetInt returns an int value from the array
func (a *Array) GetInt(index int32) int32 {
	return a.Ints[index]
}

// SetInt sets an int value in the array
func (a *Array) SetInt(index int32, val int32) {
	a.Ints[index] = val
}

// GetLong returns a long value from the array
func (a *Array) GetLong(index int32) int64 {
	return a.Longs[index]
}

// SetLong sets a long value in the array
func (a *Array) SetLong(index int32, val int64) {
	a.Longs[index] = val
}

// GetFloat returns a float value from the array
func (a *Array) GetFloat(index int32) float32 {
	return a.Floats[index]
}

// SetFloat sets a float value in the array
func (a *Array) SetFloat(index int32, val float32) {
	a.Floats[index] = val
}

// GetDouble returns a double value from the array
func (a *Array) GetDouble(index int32) float64 {
	return a.Doubles[index]
}

// SetDouble sets a double value in the array
func (a *Array) SetDouble(index int32, val float64) {
	a.Doubles[index] = val
}

// GetRef returns a reference from the array
func (a *Array) GetRef(index int32) any {
	return a.References[index]
}

// SetRef sets a reference in the array
func (a *Array) SetRef(index int32, val any) {
	a.References[index] = val
}

// IsRefArray returns true if this is a reference array
func (a *Array) IsRefArray() bool {
	return a.References != nil
}

// String returns a string representation of the array
func (a *Array) String() string {
	if a.IsRefArray() {
		return fmt.Sprintf("[L%s;@%p", a.ClassName, a)
	}
	return fmt.Sprintf("[%c@%p", arrayTypeChar(a.Type), a)
}

func arrayTypeChar(t ArrayType) rune {
	switch t {
	case ArrayTypeBoolean:
		return 'Z'
	case ArrayTypeByte:
		return 'B'
	case ArrayTypeChar:
		return 'C'
	case ArrayTypeShort:
		return 'S'
	case ArrayTypeInt:
		return 'I'
	case ArrayTypeLong:
		return 'J'
	case ArrayTypeFloat:
		return 'F'
	case ArrayTypeDouble:
		return 'D'
	default:
		return '?'
	}
}

