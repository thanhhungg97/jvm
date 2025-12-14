package classfile

import (
	"fmt"
	"os"
)

// ClassFile represents a parsed Java class file
type ClassFile struct {
	Magic        uint32
	MinorVersion uint16
	MajorVersion uint16
	ConstantPool ConstantPool
	AccessFlags  uint16
	ThisClass    uint16
	SuperClass   uint16
	Interfaces   []uint16
	Fields       []*FieldInfo
	Methods      []*MethodInfo
	Attributes   []*AttributeInfo
}

// FieldInfo represents a field in the class
type FieldInfo struct {
	AccessFlags     uint16
	NameIndex       uint16
	DescriptorIndex uint16
	Attributes      []*AttributeInfo
}

// MethodInfo represents a method in the class
type MethodInfo struct {
	AccessFlags     uint16
	NameIndex       uint16
	DescriptorIndex uint16
	Attributes      []*AttributeInfo
}

// AttributeInfo represents an attribute
type AttributeInfo struct {
	NameIndex uint16
	Info      []byte
}

// CodeAttribute represents the Code attribute of a method
type CodeAttribute struct {
	MaxStack       uint16
	MaxLocals      uint16
	Code           []byte
	ExceptionTable []*ExceptionTableEntry
	Attributes     []*AttributeInfo
}

// ExceptionTableEntry represents an exception handler
type ExceptionTableEntry struct {
	StartPC   uint16
	EndPC     uint16
	HandlerPC uint16
	CatchType uint16
}

// Parse reads a class file from bytes
func Parse(data []byte) (*ClassFile, error) {
	reader := NewClassReader(data)
	cf := &ClassFile{}

	cf.Magic = reader.ReadU4()
	if cf.Magic != 0xCAFEBABE {
		return nil, fmt.Errorf("invalid class file: bad magic number %X", cf.Magic)
	}

	cf.MinorVersion = reader.ReadU2()
	cf.MajorVersion = reader.ReadU2()
	cf.ConstantPool = ReadConstantPool(reader)
	cf.AccessFlags = reader.ReadU2()
	cf.ThisClass = reader.ReadU2()
	cf.SuperClass = reader.ReadU2()
	cf.Interfaces = reader.ReadU2s()
	cf.Fields = readFields(reader, cf.ConstantPool)
	cf.Methods = readMethods(reader, cf.ConstantPool)
	cf.Attributes = readAttributes(reader, cf.ConstantPool)

	return cf, nil
}

// ParseFile reads a class file from disk
func ParseFile(path string) (*ClassFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return Parse(data)
}

func readFields(reader *ClassReader, cp ConstantPool) []*FieldInfo {
	count := reader.ReadU2()
	fields := make([]*FieldInfo, count)
	for i := range fields {
		fields[i] = &FieldInfo{
			AccessFlags:     reader.ReadU2(),
			NameIndex:       reader.ReadU2(),
			DescriptorIndex: reader.ReadU2(),
			Attributes:      readAttributes(reader, cp),
		}
	}
	return fields
}

func readMethods(reader *ClassReader, cp ConstantPool) []*MethodInfo {
	count := reader.ReadU2()
	methods := make([]*MethodInfo, count)
	for i := range methods {
		methods[i] = &MethodInfo{
			AccessFlags:     reader.ReadU2(),
			NameIndex:       reader.ReadU2(),
			DescriptorIndex: reader.ReadU2(),
			Attributes:      readAttributes(reader, cp),
		}
	}
	return methods
}

func readAttributes(reader *ClassReader, cp ConstantPool) []*AttributeInfo {
	count := reader.ReadU2()
	attrs := make([]*AttributeInfo, count)
	for i := range attrs {
		nameIndex := reader.ReadU2()
		length := reader.ReadU4()
		info := reader.ReadBytes(int(length))
		attrs[i] = &AttributeInfo{
			NameIndex: nameIndex,
			Info:      info,
		}
	}
	return attrs
}

// ClassName returns the name of this class
func (cf *ClassFile) ClassName() string {
	return cf.ConstantPool.GetClassName(cf.ThisClass)
}

// SuperClassName returns the name of the superclass
func (cf *ClassFile) SuperClassName() string {
	if cf.SuperClass == 0 {
		return ""
	}
	return cf.ConstantPool.GetClassName(cf.SuperClass)
}

// GetMethod finds a method by name and descriptor
func (cf *ClassFile) GetMethod(name, descriptor string) *MethodInfo {
	for _, method := range cf.Methods {
		methodName := cf.ConstantPool.GetUtf8(method.NameIndex)
		methodDesc := cf.ConstantPool.GetUtf8(method.DescriptorIndex)
		if methodName == name && (descriptor == "" || methodDesc == descriptor) {
			return method
		}
	}
	return nil
}

// Name returns the method name
func (m *MethodInfo) Name(cp ConstantPool) string {
	return cp.GetUtf8(m.NameIndex)
}

// Descriptor returns the method descriptor
func (m *MethodInfo) Descriptor(cp ConstantPool) string {
	return cp.GetUtf8(m.DescriptorIndex)
}

// GetCodeAttribute extracts the Code attribute from a method
func (m *MethodInfo) GetCodeAttribute(cp ConstantPool) *CodeAttribute {
	for _, attr := range m.Attributes {
		name := cp.GetUtf8(attr.NameIndex)
		if name == "Code" {
			return parseCodeAttribute(attr.Info)
		}
	}
	return nil
}

func parseCodeAttribute(data []byte) *CodeAttribute {
	reader := NewClassReader(data)
	code := &CodeAttribute{}
	code.MaxStack = reader.ReadU2()
	code.MaxLocals = reader.ReadU2()

	codeLength := reader.ReadU4()
	code.Code = reader.ReadBytes(int(codeLength))

	exceptionTableLength := reader.ReadU2()
	code.ExceptionTable = make([]*ExceptionTableEntry, exceptionTableLength)
	for i := range code.ExceptionTable {
		code.ExceptionTable[i] = &ExceptionTableEntry{
			StartPC:   reader.ReadU2(),
			EndPC:     reader.ReadU2(),
			HandlerPC: reader.ReadU2(),
			CatchType: reader.ReadU2(),
		}
	}

	// Read attributes (LineNumberTable, etc.)
	attrCount := reader.ReadU2()
	code.Attributes = make([]*AttributeInfo, attrCount)
	for i := range code.Attributes {
		nameIndex := reader.ReadU2()
		length := reader.ReadU4()
		info := reader.ReadBytes(int(length))
		code.Attributes[i] = &AttributeInfo{
			NameIndex: nameIndex,
			Info:      info,
		}
	}

	return code
}
