package classfile

// Constant pool tags
const (
	CONSTANT_Utf8               = 1
	CONSTANT_Integer            = 3
	CONSTANT_Float              = 4
	CONSTANT_Long               = 5
	CONSTANT_Double             = 6
	CONSTANT_Class              = 7
	CONSTANT_String             = 8
	CONSTANT_Fieldref           = 9
	CONSTANT_Methodref          = 10
	CONSTANT_InterfaceMethodref = 11
	CONSTANT_NameAndType        = 12
	CONSTANT_MethodHandle       = 15
	CONSTANT_MethodType         = 16
	CONSTANT_InvokeDynamic      = 18
)

// ConstantInfo is the base interface for all constant pool entries
type ConstantInfo interface {
	Tag() uint8
}

// ConstantPool holds all constant pool entries
type ConstantPool []ConstantInfo

// ReadConstantPool parses the constant pool from class data
func ReadConstantPool(reader *ClassReader) ConstantPool {
	count := int(reader.ReadU2())
	cp := make(ConstantPool, count)

	// The constant pool index starts at 1
	for i := 1; i < count; i++ {
		tag := reader.ReadU1()
		cp[i] = readConstantInfo(reader, tag)

		// Long and Double take two slots
		if tag == CONSTANT_Long || tag == CONSTANT_Double {
			i++
		}
	}
	return cp
}

func readConstantInfo(reader *ClassReader, tag uint8) ConstantInfo {
	switch tag {
	case CONSTANT_Utf8:
		return &ConstantUtf8Info{tag: tag, Value: readUtf8(reader)}
	case CONSTANT_Integer:
		return &ConstantIntegerInfo{tag: tag, Value: int32(reader.ReadU4())}
	case CONSTANT_Float:
		return &ConstantFloatInfo{tag: tag, Value: reader.ReadU4()}
	case CONSTANT_Long:
		high := uint64(reader.ReadU4())
		low := uint64(reader.ReadU4())
		return &ConstantLongInfo{tag: tag, Value: int64(high<<32 | low)}
	case CONSTANT_Double:
		high := uint64(reader.ReadU4())
		low := uint64(reader.ReadU4())
		return &ConstantDoubleInfo{tag: tag, Value: high<<32 | low}
	case CONSTANT_Class:
		return &ConstantClassInfo{tag: tag, NameIndex: reader.ReadU2()}
	case CONSTANT_String:
		return &ConstantStringInfo{tag: tag, StringIndex: reader.ReadU2()}
	case CONSTANT_Fieldref:
		return &ConstantFieldrefInfo{
			tag:              tag,
			ClassIndex:       reader.ReadU2(),
			NameAndTypeIndex: reader.ReadU2(),
		}
	case CONSTANT_Methodref:
		return &ConstantMethodrefInfo{
			tag:              tag,
			ClassIndex:       reader.ReadU2(),
			NameAndTypeIndex: reader.ReadU2(),
		}
	case CONSTANT_InterfaceMethodref:
		return &ConstantInterfaceMethodrefInfo{
			tag:              tag,
			ClassIndex:       reader.ReadU2(),
			NameAndTypeIndex: reader.ReadU2(),
		}
	case CONSTANT_NameAndType:
		return &ConstantNameAndTypeInfo{
			tag:             tag,
			NameIndex:       reader.ReadU2(),
			DescriptorIndex: reader.ReadU2(),
		}
	case CONSTANT_MethodHandle:
		return &ConstantMethodHandleInfo{
			tag:            tag,
			ReferenceKind:  reader.ReadU1(),
			ReferenceIndex: reader.ReadU2(),
		}
	case CONSTANT_MethodType:
		return &ConstantMethodTypeInfo{tag: tag, DescriptorIndex: reader.ReadU2()}
	case CONSTANT_InvokeDynamic:
		return &ConstantInvokeDynamicInfo{
			tag:                  tag,
			BootstrapMethodIndex: reader.ReadU2(),
			NameAndTypeIndex:     reader.ReadU2(),
		}
	default:
		panic("Unknown constant pool tag")
	}
}

func readUtf8(reader *ClassReader) string {
	length := reader.ReadU2()
	bytes := reader.ReadBytes(int(length))
	return string(bytes)
}

// Constant pool entry types

type ConstantUtf8Info struct {
	tag   uint8
	Value string
}

func (c *ConstantUtf8Info) Tag() uint8 { return c.tag }

type ConstantIntegerInfo struct {
	tag   uint8
	Value int32
}

func (c *ConstantIntegerInfo) Tag() uint8 { return c.tag }

type ConstantFloatInfo struct {
	tag   uint8
	Value uint32
}

func (c *ConstantFloatInfo) Tag() uint8 { return c.tag }

type ConstantLongInfo struct {
	tag   uint8
	Value int64
}

func (c *ConstantLongInfo) Tag() uint8 { return c.tag }

type ConstantDoubleInfo struct {
	tag   uint8
	Value uint64
}

func (c *ConstantDoubleInfo) Tag() uint8 { return c.tag }

type ConstantClassInfo struct {
	tag       uint8
	NameIndex uint16
}

func (c *ConstantClassInfo) Tag() uint8 { return c.tag }

type ConstantStringInfo struct {
	tag         uint8
	StringIndex uint16
}

func (c *ConstantStringInfo) Tag() uint8 { return c.tag }

type ConstantFieldrefInfo struct {
	tag              uint8
	ClassIndex       uint16
	NameAndTypeIndex uint16
}

func (c *ConstantFieldrefInfo) Tag() uint8 { return c.tag }

type ConstantMethodrefInfo struct {
	tag              uint8
	ClassIndex       uint16
	NameAndTypeIndex uint16
}

func (c *ConstantMethodrefInfo) Tag() uint8 { return c.tag }

type ConstantInterfaceMethodrefInfo struct {
	tag              uint8
	ClassIndex       uint16
	NameAndTypeIndex uint16
}

func (c *ConstantInterfaceMethodrefInfo) Tag() uint8 { return c.tag }

type ConstantNameAndTypeInfo struct {
	tag             uint8
	NameIndex       uint16
	DescriptorIndex uint16
}

func (c *ConstantNameAndTypeInfo) Tag() uint8 { return c.tag }

type ConstantMethodHandleInfo struct {
	tag            uint8
	ReferenceKind  uint8
	ReferenceIndex uint16
}

func (c *ConstantMethodHandleInfo) Tag() uint8 { return c.tag }

type ConstantMethodTypeInfo struct {
	tag             uint8
	DescriptorIndex uint16
}

func (c *ConstantMethodTypeInfo) Tag() uint8 { return c.tag }

type ConstantInvokeDynamicInfo struct {
	tag                  uint8
	BootstrapMethodIndex uint16
	NameAndTypeIndex     uint16
}

func (c *ConstantInvokeDynamicInfo) Tag() uint8 { return c.tag }

// Helper methods for ConstantPool

// GetUtf8 retrieves a UTF8 string from the constant pool
func (cp ConstantPool) GetUtf8(index uint16) string {
	if utf8, ok := cp[index].(*ConstantUtf8Info); ok {
		return utf8.Value
	}
	return ""
}

// GetClassName retrieves a class name from the constant pool
func (cp ConstantPool) GetClassName(index uint16) string {
	if classInfo, ok := cp[index].(*ConstantClassInfo); ok {
		return cp.GetUtf8(classInfo.NameIndex)
	}
	return ""
}

// GetNameAndType retrieves name and type descriptor
func (cp ConstantPool) GetNameAndType(index uint16) (string, string) {
	if nat, ok := cp[index].(*ConstantNameAndTypeInfo); ok {
		return cp.GetUtf8(nat.NameIndex), cp.GetUtf8(nat.DescriptorIndex)
	}
	return "", ""
}
