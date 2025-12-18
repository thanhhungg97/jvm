package interpreter

// OpcodeCategory represents the category of a bytecode instruction
type OpcodeCategory uint8

const (
	CategoryUnknown OpcodeCategory = iota
	CategoryConst                  // Constants: iconst, ldc, etc.
	CategoryLoad                   // Loads: iload, aload, etc.
	CategoryStore                  // Stores: istore, astore, etc.
	CategoryMath                   // Arithmetic, bitwise, stack ops
	CategoryControl                // Branches, returns
	CategoryArray                  // Array operations
	CategoryObject                 // Object operations, fields
	CategoryInvoke                 // Method invocations
)

// opcodeCategories maps each opcode to its category
var opcodeCategories [256]OpcodeCategory

func init() {
	// Constants
	for _, op := range []uint8{NOP, ACONST_NULL, ICONST_M1, ICONST_0, ICONST_1, ICONST_2,
		ICONST_3, ICONST_4, ICONST_5, LCONST_0, LCONST_1, BIPUSH, SIPUSH, LDC, LDC_W, LDC2_W} {
		opcodeCategories[op] = CategoryConst
	}

	// Loads
	for _, op := range []uint8{ILOAD, LLOAD, ALOAD, ILOAD_0, ILOAD_1, ILOAD_2, ILOAD_3,
		LLOAD_0, LLOAD_1, LLOAD_2, LLOAD_3, ALOAD_0, ALOAD_1, ALOAD_2, ALOAD_3} {
		opcodeCategories[op] = CategoryLoad
	}

	// Stores
	for _, op := range []uint8{ISTORE, LSTORE, ASTORE, ISTORE_0, ISTORE_1, ISTORE_2, ISTORE_3,
		LSTORE_0, LSTORE_1, LSTORE_2, LSTORE_3, ASTORE_0, ASTORE_1, ASTORE_2, ASTORE_3} {
		opcodeCategories[op] = CategoryStore
	}

	// Math/Stack
	for _, op := range []uint8{POP, POP2, DUP, DUP_X1, DUP_X2, DUP2, SWAP,
		IADD, LADD, ISUB, LSUB, IMUL, LMUL, IDIV, LDIV, IREM, LREM, INEG, LNEG,
		ISHL, LSHL, ISHR, LSHR, IUSHR, LUSHR, IAND, LAND, IOR, LOR, IXOR, LXOR,
		IINC, I2L, I2F, I2D, L2I, LCMP} {
		opcodeCategories[op] = CategoryMath
	}

	// Control
	for _, op := range []uint8{IFEQ, IFNE, IFLT, IFGE, IFGT, IFLE,
		IF_ICMPEQ, IF_ICMPNE, IF_ICMPLT, IF_ICMPGE, IF_ICMPGT, IF_ICMPLE,
		IF_ACMPEQ, IF_ACMPNE, GOTO, JSR, RET, TABLESWITCH, LOOKUPSWITCH,
		IRETURN, LRETURN, FRETURN, DRETURN, ARETURN, RETURN, IFNULL, IFNONNULL, GOTO_W} {
		opcodeCategories[op] = CategoryControl
	}

	// Array
	for _, op := range []uint8{IALOAD, LALOAD, FALOAD, DALOAD, AALOAD, BALOAD, CALOAD, SALOAD,
		IASTORE, LASTORE, FASTORE, DASTORE, AASTORE, BASTORE, CASTORE, SASTORE,
		NEWARRAY, ANEWARRAY, ARRAYLENGTH} {
		opcodeCategories[op] = CategoryArray
	}

	// Object
	for _, op := range []uint8{GETSTATIC, PUTSTATIC, GETFIELD, PUTFIELD, NEW,
		ATHROW, CHECKCAST, INSTANCEOF, MONITORENTER, MONITOREXIT} {
		opcodeCategories[op] = CategoryObject
	}

	// Invoke
	for _, op := range []uint8{INVOKEVIRTUAL, INVOKESPECIAL, INVOKESTATIC, INVOKEINTERFACE, INVOKEDYNAMIC} {
		opcodeCategories[op] = CategoryInvoke
	}
}

// Category returns the category of an opcode
func Category(opcode uint8) OpcodeCategory {
	return opcodeCategories[opcode]
}

// JVM Bytecode opcodes
const (
	// Constants
	NOP         = 0x00
	ACONST_NULL = 0x01
	ICONST_M1   = 0x02
	ICONST_0    = 0x03
	ICONST_1    = 0x04
	ICONST_2    = 0x05
	ICONST_3    = 0x06
	ICONST_4    = 0x07
	ICONST_5    = 0x08
	LCONST_0    = 0x09
	LCONST_1    = 0x0A
	BIPUSH      = 0x10
	SIPUSH      = 0x11
	LDC         = 0x12
	LDC_W       = 0x13
	LDC2_W      = 0x14

	// Loads
	ILOAD   = 0x15
	LLOAD   = 0x16
	ALOAD   = 0x19
	ILOAD_0 = 0x1A
	ILOAD_1 = 0x1B
	ILOAD_2 = 0x1C
	ILOAD_3 = 0x1D
	LLOAD_0 = 0x1E
	LLOAD_1 = 0x1F
	LLOAD_2 = 0x20
	LLOAD_3 = 0x21
	ALOAD_0 = 0x2A
	ALOAD_1 = 0x2B
	ALOAD_2 = 0x2C
	ALOAD_3 = 0x2D

	// Stores
	ISTORE   = 0x36
	LSTORE   = 0x37
	ASTORE   = 0x3A
	ISTORE_0 = 0x3B
	ISTORE_1 = 0x3C
	ISTORE_2 = 0x3D
	ISTORE_3 = 0x3E
	LSTORE_0 = 0x3F
	LSTORE_1 = 0x40
	LSTORE_2 = 0x41
	LSTORE_3 = 0x42
	ASTORE_0 = 0x4B
	ASTORE_1 = 0x4C
	ASTORE_2 = 0x4D
	ASTORE_3 = 0x4E

	// Stack
	POP    = 0x57
	POP2   = 0x58
	DUP    = 0x59
	DUP_X1 = 0x5A
	DUP_X2 = 0x5B
	DUP2   = 0x5C
	SWAP   = 0x5F

	// Arithmetic
	IADD = 0x60
	LADD = 0x61
	ISUB = 0x64
	LSUB = 0x65
	IMUL = 0x68
	LMUL = 0x69
	IDIV = 0x6C
	LDIV = 0x6D
	IREM = 0x70
	LREM = 0x71
	INEG = 0x74
	LNEG = 0x75

	// Shifts
	ISHL  = 0x78
	LSHL  = 0x79
	ISHR  = 0x7A
	LSHR  = 0x7B
	IUSHR = 0x7C
	LUSHR = 0x7D

	// Bitwise
	IAND = 0x7E
	LAND = 0x7F
	IOR  = 0x80
	LOR  = 0x81
	IXOR = 0x82
	LXOR = 0x83

	// Increment
	IINC = 0x84

	// Conversions
	I2L = 0x85
	I2F = 0x86
	I2D = 0x87
	L2I = 0x88

	// Comparisons
	LCMP      = 0x94
	IFEQ      = 0x99
	IFNE      = 0x9A
	IFLT      = 0x9B
	IFGE      = 0x9C
	IFGT      = 0x9D
	IFLE      = 0x9E
	IF_ICMPEQ = 0x9F
	IF_ICMPNE = 0xA0
	IF_ICMPLT = 0xA1
	IF_ICMPGE = 0xA2
	IF_ICMPGT = 0xA3
	IF_ICMPLE = 0xA4
	IF_ACMPEQ = 0xA5
	IF_ACMPNE = 0xA6

	// Control
	GOTO         = 0xA7
	JSR          = 0xA8
	RET          = 0xA9
	TABLESWITCH  = 0xAA
	LOOKUPSWITCH = 0xAB
	IRETURN      = 0xAC
	LRETURN      = 0xAD
	FRETURN      = 0xAE
	DRETURN      = 0xAF
	ARETURN      = 0xB0
	RETURN       = 0xB1

	// Array loads
	IALOAD = 0x2E
	LALOAD = 0x2F
	FALOAD = 0x30
	DALOAD = 0x31
	AALOAD = 0x32
	BALOAD = 0x33
	CALOAD = 0x34
	SALOAD = 0x35

	// Array stores
	IASTORE = 0x4F
	LASTORE = 0x50
	FASTORE = 0x51
	DASTORE = 0x52
	AASTORE = 0x53
	BASTORE = 0x54
	CASTORE = 0x55
	SASTORE = 0x56

	// References
	GETSTATIC       = 0xB2
	PUTSTATIC       = 0xB3
	GETFIELD        = 0xB4
	PUTFIELD        = 0xB5
	INVOKEVIRTUAL   = 0xB6
	INVOKESPECIAL   = 0xB7
	INVOKESTATIC    = 0xB8
	INVOKEINTERFACE = 0xB9
	INVOKEDYNAMIC   = 0xBA
	NEW             = 0xBB
	NEWARRAY        = 0xBC
	ANEWARRAY       = 0xBD
	ARRAYLENGTH     = 0xBE
	ATHROW          = 0xBF
	CHECKCAST       = 0xC0
	INSTANCEOF      = 0xC1

	// Synchronization
	MONITORENTER = 0xC2
	MONITOREXIT  = 0xC3

	// Extended
	IFNULL    = 0xC6
	IFNONNULL = 0xC7
	GOTO_W    = 0xC8
)
