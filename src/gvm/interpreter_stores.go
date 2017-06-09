package gvm

/*54 (0X36)*/
func ISTORE(opcode uint8, f *StackFrame, t *Thread, c *JavaClass, m *JavaMethod) {
	index := m.code[f.pc+1]
	f.storeVar(uint(index), f.pop().(t_int))
}

/*55 (0X37)*/
func LSTORE(opcode uint8, f *StackFrame, t *Thread, c *JavaClass, m *JavaMethod) {
	index := m.code[f.pc+1]
	f.storeVar(uint(index), f.pop().(t_long))
}

/*56 (0X38)*/
func FSTORE(opcode uint8, f *StackFrame, t *Thread, c *JavaClass, m *JavaMethod) {
	index := m.code[f.pc+1]
	f.storeVar(uint(index), f.pop().(t_float))
}

/*57 (0X39)*/
func DSTORE(opcode uint8, f *StackFrame, t *Thread, c *JavaClass, m *JavaMethod) {
	index := m.code[f.pc+1]
	f.storeVar(uint(index), f.pop().(t_double))
}

/*58 (0X3A)*/
func ASTORE(opcode uint8, f *StackFrame, t *Thread, c *JavaClass, m *JavaMethod) {
	index := m.code[f.pc+1]
	f.storeVar(uint(index), f.pop().(Reference))
}

/*59 (0X3B)*/
func ISTORE_0(opcode uint8, f *StackFrame, t *Thread, c *JavaClass, m *JavaMethod) {
	f.storeVar(0, f.pop().(t_int))
}

/*60 (0X3C)*/
func ISTORE_1(opcode uint8, f *StackFrame, t *Thread, c *JavaClass, m *JavaMethod) {
	f.storeVar(1, f.pop().(t_int))
}

/*61 (0X3D)*/
func ISTORE_2(opcode uint8, f *StackFrame, t *Thread, c *JavaClass, m *JavaMethod) {
	f.storeVar(2, f.pop().(t_int))
}

/*62 (0X3E)*/
func ISTORE_3(opcode uint8, f *StackFrame, t *Thread, c *JavaClass, m *JavaMethod) {
	f.storeVar(3, f.pop().(t_int))
}

/*63 (0X3F)*/
func LSTORE_0(opcode uint8, f *StackFrame, t *Thread, c *JavaClass, m *JavaMethod) {
	f.storeVar(0, f.pop().(t_long))
}

/*64 (0X40)*/
func LSTORE_1(opcode uint8, f *StackFrame, t *Thread, c *JavaClass, m *JavaMethod) {
	f.storeVar(1, f.pop().(t_long))
}

/*65 (0X41)*/
func LSTORE_2(opcode uint8, f *StackFrame, t *Thread, c *JavaClass, m *JavaMethod) {
	f.storeVar(2, f.pop().(t_long))
}

/*66 (0X42)*/
func LSTORE_3(opcode uint8, f *StackFrame, t *Thread, c *JavaClass, m *JavaMethod) {
	f.storeVar(3, f.pop().(t_long))
}

/*67 (0X43)*/
func FSTORE_0(opcode uint8, f *StackFrame, t *Thread, c *JavaClass, m *JavaMethod) {
	f.storeVar(0, f.pop().(t_float))
}

/*68 (0X44)*/
func FSTORE_1(opcode uint8, f *StackFrame, t *Thread, c *JavaClass, m *JavaMethod) {
	f.storeVar(1, f.pop().(t_float))
}

/*69 (0X45)*/
func FSTORE_2(opcode uint8, f *StackFrame, t *Thread, c *JavaClass, m *JavaMethod) {
	f.storeVar(2, f.pop().(t_float))
}

/*70 (0X46)*/
func FSTORE_3(opcode uint8, f *StackFrame, t *Thread, c *JavaClass, m *JavaMethod) {
	f.storeVar(3, f.pop().(t_float))
}

/*71 (0X47)*/
func DSTORE_0(opcode uint8, f *StackFrame, t *Thread, c *JavaClass, m *JavaMethod) {
	f.storeVar(0, f.pop().(t_double))
}

/*72 (0X48)*/
func DSTORE_1(opcode uint8, f *StackFrame, t *Thread, c *JavaClass, m *JavaMethod) {
	f.storeVar(1, f.pop().(t_double))
}

/*73 (0X49)*/
func DSTORE_2(opcode uint8, f *StackFrame, t *Thread, c *JavaClass, m *JavaMethod) {
	f.storeVar(2, f.pop().(t_double))
}

/*74 (0X4A)*/
func DSTORE_3(opcode uint8, f *StackFrame, t *Thread, c *JavaClass, m *JavaMethod) {
	f.storeVar(3, f.pop().(t_double))
}

/*75 (0X4B)*/
func ASTORE_0(opcode uint8, f *StackFrame, t *Thread, c *JavaClass, m *JavaMethod) {
	f.storeVar(0, f.pop().(Reference))
}

/*76 (0X4C)*/
func ASTORE_1(opcode uint8, f *StackFrame, t *Thread, c *JavaClass, m *JavaMethod) {
	f.storeVar(1, f.pop().(Reference))
}

/*77 (0X4D)*/
func ASTORE_2(opcode uint8, f *StackFrame, t *Thread, c *JavaClass, m *JavaMethod) {
	f.storeVar(2, f.pop().(Reference))
}

/*78 (0X4E)*/
func ASTORE_3(opcode uint8, f *StackFrame, t *Thread, c *JavaClass, m *JavaMethod) {
	f.storeVar(3, f.pop().(Reference))
}

/*79 (0X4F)*/
func IASTORE(opcode uint8, f *StackFrame, t *Thread, c *JavaClass, m *JavaMethod) {
	value := f.pop().(t_int)
	index := f.pop().(t_int)
	arrayref := f.pop().(*t_array)
	arrayref.elements[index] = value
	//TODO check component type and boundary
}

/*80 (0X50)*/
func LASTORE(opcode uint8, f *StackFrame, t *Thread, c *JavaClass, m *JavaMethod) {
	value := f.pop().(t_long)
	index := f.pop().(t_int)
	arrayref := f.pop().(*t_array)
	arrayref.elements[index] = value
	//TODO check component type and boundary
}

/*81 (0X51)*/
func FASTORE(opcode uint8, f *StackFrame, t *Thread, c *JavaClass, m *JavaMethod) {
	value := f.pop().(t_float)
	index := f.pop().(t_int)
	arrayref := f.pop().(*t_array)
	arrayref.elements[index] = value
	//TODO check component type and boundary
}

/*82 (0X52)*/
func DASTORE(opcode uint8, f *StackFrame, t *Thread, c *JavaClass, m *JavaMethod) {
	value := f.pop().(t_double)
	index := f.pop().(t_int)
	arrayref := f.pop().(*t_array)
	arrayref.elements[index] = value
	//TODO check component type and boundary
}

/*83 (0X53)*/
func AASTORE(opcode uint8, f *StackFrame, t *Thread, c *JavaClass, m *JavaMethod) {
	value := f.pop().(*t_array)
	index := f.pop().(t_int)
	arrayref := f.pop().(*t_array)
	arrayref.elements[index] = value
	//TODO check component type, boundary and subtypes
}

/*84 (0X54)*/
func BASTORE(opcode uint8, f *StackFrame, t *Thread, c *JavaClass, m *JavaMethod) {
	value := f.pop().(t_int)
	index := f.pop().(t_int)
	arrayref := f.pop().(*t_array)
	arrayref.elements[index] = t_byte(value)
	//TODO check component type and boundary
}

/*85 (0X55)*/
func CASTORE(opcode uint8, f *StackFrame, t *Thread, c *JavaClass, m *JavaMethod) {
	value := f.pop().(t_int)
	index := f.pop().(t_int)
	arrayref := f.pop().(*t_array)
	arrayref.elements[index] = t_char(value)
	//TODO check component type and boundary
}

/*86 (0X56)*/
func SASTORE(opcode uint8, f *StackFrame, t *Thread, c *JavaClass, m *JavaMethod) {
	value := f.pop().(t_int)
	index := f.pop().(t_int)
	arrayref := f.pop().(*t_array)
	arrayref.elements[index] = t_short(value)
	//TODO check component type and boundary
}
