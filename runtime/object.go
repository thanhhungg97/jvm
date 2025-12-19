package runtime

import (
	"fmt"
	"simplejvm/classfile"
)

// Object represents a JVM object instance
type Object struct {
	Class      *classfile.ClassFile // The class this object is an instance of
	Fields     map[string]any       // Instance fields (fieldName -> value)
	FieldSlots map[string]int64     // Primitive field values
}

// NewObject creates a new object instance
func NewObject(class *classfile.ClassFile) *Object {
	obj := &Object{
		Class:      class,
		Fields:     make(map[string]any),
		FieldSlots: make(map[string]int64),
	}

	// Initialize fields with default values
	for _, field := range class.Fields {
		fieldName := class.ConstantPool.GetUtf8(field.NameIndex)
		descriptor := class.ConstantPool.GetUtf8(field.DescriptorIndex)

		// Skip static fields (handled separately)
		if field.AccessFlags&0x0008 != 0 { // ACC_STATIC
			continue
		}

		// Initialize based on type
		switch descriptor[0] {
		case 'B', 'C', 'I', 'S', 'Z': // byte, char, int, short, boolean
			obj.FieldSlots[fieldName] = 0
		case 'J': // long
			obj.FieldSlots[fieldName] = 0
		case 'F': // float
			obj.FieldSlots[fieldName] = 0
		case 'D': // double
			obj.FieldSlots[fieldName] = 0
		case 'L', '[': // object or array
			obj.Fields[fieldName] = nil
		}
	}

	return obj
}

// GetFieldInt gets an int field value
func (o *Object) GetFieldInt(name string) int32 {
	return int32(o.FieldSlots[name])
}

// SetFieldInt sets an int field value
func (o *Object) SetFieldInt(name string, val int32) {
	o.FieldSlots[name] = int64(val)
}

// GetFieldLong gets a long field value
func (o *Object) GetFieldLong(name string) int64 {
	return o.FieldSlots[name]
}

// SetFieldLong sets a long field value
func (o *Object) SetFieldLong(name string, val int64) {
	o.FieldSlots[name] = val
}

// GetFieldRef gets a reference field value
func (o *Object) GetFieldRef(name string) any {
	return o.Fields[name]
}

// SetFieldRef sets a reference field value
func (o *Object) SetFieldRef(name string, val any) {
	o.Fields[name] = val
}

// ClassName returns the class name of this object
func (o *Object) ClassName() string {
	if o.Class != nil {
		return o.Class.ClassName()
	}
	return "<unknown>"
}

// String returns a string representation of the object
func (o *Object) String() string {
	return fmt.Sprintf("%s@%p", o.ClassName(), o)
}

// IsInstanceOf checks if this object is an instance of the given class
func (o *Object) IsInstanceOf(className string) bool {
	if o.Class == nil {
		return false
	}
	// Simple check - just compare class names
	// TODO: Handle inheritance hierarchy
	return o.Class.ClassName() == className
}
