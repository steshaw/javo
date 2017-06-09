package gvm

/*
cp_info {
    u1 tag;
    u1 info[];
}
 */
type RuntimeConstantPoolInfo interface {
	resolve() RuntimeConstantPoolInfo
}

/*
CONSTANT_Class_info {
    u1 tag;
    u2 name_index;
}
 */
type RuntimeConstantClassInfo struct {
	hostClass       *JavaClass
	nameIndex       u2
	resolved        bool
	class           *JavaClass
}

func (this *RuntimeConstantClassInfo) resolve() RuntimeConstantPoolInfo {
	if !this.resolved {
		name := this.hostClass.constantPool[this.nameIndex].resolve().(*RuntimeConstantUtf8Info).value

		if this == this.hostClass.constantPool[this.hostClass.thisClass] {
			this.class = this.hostClass // current class
		} else {
			this.class = bootstrapClassLoader.load(name)
		}

		this.resolved = true
	}
	return this
}

/*
CONSTANT_Fieldref_info {
    u1 tag;
    u2 class_index;
    u2 name_and_type_index;
}
 */
type RuntimeConstantFieldrefInfo struct {
	hostClass           *JavaClass
	classIndex          u2
	nameAndTypeIndex    u2
	resolved            bool
	field               *JavaField
}

func (this *RuntimeConstantFieldrefInfo) resolve() RuntimeConstantPoolInfo  {
	if !this.resolved {
		rcpClass := this.hostClass.constantPool[this.classIndex].resolve().(*RuntimeConstantClassInfo)

		nameAndType := this.hostClass.constantPool[this.nameAndTypeIndex].resolve().(*RuntimeConstantNameAndTypeInfo)
		this.field = rcpClass.class.findField(nameAndType.toString())
		this.resolved = true
	}
	return this
}

/*
CONSTANT_Methodref_info {
    u1 tag;
    u2 class_index;
    u2 name_and_type_index;
}
 */
type RuntimeConstantMethodrefInfo struct {
	hostClass       *JavaClass
	classIndex       u2
	nameAndTypeIndex u2
	resolved         bool
	method           *JavaMethod
}

func (this *RuntimeConstantMethodrefInfo) resolve() RuntimeConstantPoolInfo  {
	if !this.resolved {
		rcpClass := this.hostClass.constantPool[this.classIndex].resolve().(*RuntimeConstantClassInfo)

		nameAndType := this.hostClass.constantPool[this.nameAndTypeIndex].resolve().(*RuntimeConstantNameAndTypeInfo)
		this.method = rcpClass.class.findMethod(nameAndType.toString())
		this.resolved = true
	}
	return this
}


/*
CONSTANT_InterfaceMethodref_info {
    u1 tag;
    u2 class_index;
    u2 name_and_type_index;
}
 */
type RuntimeConstantInterfaceMethodrefInfo struct {
	hostClass       *JavaClass
	classIndex       u2
	nameAndTypeIndex u2
	resolved         bool
	method           *JavaMethod
}

func (this *RuntimeConstantInterfaceMethodrefInfo) resolve() RuntimeConstantPoolInfo  {
	if !this.resolved {
		//TODO

		this.resolved = true
	}
	return this
}

/*
CONSTANT_String_info {
    u1 tag;
    u2 string_index;
}
 */
type RuntimeConstantStringInfo struct {
	hostClass       *JavaClass
	stringIndex     u2
	resolved        bool
	//value           []t_char
	value           *java_lang_string
}

func (this *RuntimeConstantStringInfo) resolve() RuntimeConstantPoolInfo {
	if !this.resolved {
		utf8string := this.hostClass.constantPool[this.stringIndex].resolve().(*RuntimeConstantUtf8Info).value
		javastring, found := stringTable[utf8string]
		if !found {
			javastring = newJavaLangString(utf8string)
			stringTable[utf8string] = javastring
		}
		this.value = javastring
		this.resolved = true
	}
	return this
}

/*
CONSTANT_Integer_info {
    u1 tag;
    u4 bytes;
}
 */
type RuntimeConstantIntegerInfo struct {
	hostClass *JavaClass
	bytes     u4
	resolved  bool
	value     t_int
}


func (this *RuntimeConstantIntegerInfo) resolve() RuntimeConstantPoolInfo  {
	if !this.resolved {
		this.value = t_int(this.bytes)
		this.resolved = true
	}
	return this
}

/*
CONSTANT_Float_info {
    u1 tag;
    u4 bytes;
}
 */
type RuntimeConstantFloatInfo struct {
	hostClass *JavaClass
	bytes     u4
	resolved  bool
	value     t_float
}

func (this *RuntimeConstantFloatInfo) resolve() RuntimeConstantPoolInfo  {
	if !this.resolved {
		this.value = t_float(this.bytes)
		this.resolved = true
	}
	return this
}

/*
CONSTANT_Long_info {
    u1 tag;
    u4 high_bytes;
    u4 low_bytes;
}
 */
type RuntimeConstantLongInfo struct {
	hostClass *JavaClass
	highBytes u4
	lowBytes  u4
	resolved  bool
	value     t_long
}

func (this *RuntimeConstantLongInfo) resolve() RuntimeConstantPoolInfo  {
	if !this.resolved {
		this.value = t_long(this.highBytes << 8 | this.lowBytes)
		this.resolved = true
	}
	return this
}

/*
CONSTANT_Double_info {
    u1 tag;
    u4 high_bytes;
    u4 low_bytes;
}
 */
type RuntimeConstantDoubleInfo struct {
	hostClass *JavaClass
	highBytes u4
	lowBytes  u4
	resolved  bool
	value     t_double
}


func (this *RuntimeConstantDoubleInfo) resolve() RuntimeConstantPoolInfo  {
	if !this.resolved {
		this.value = t_double(this.highBytes << 8 | this.lowBytes)
		this.resolved = true
	}
	return this
}

/*
CONSTANT_NameAndType_info {
    u1 tag;
    u2 name_index;
    u2 descriptor_index;
}
 */
type RuntimeConstantNameAndTypeInfo struct {
	hostClass       *JavaClass
	nameIndex       u2
	descriptorIndex u2
	resolved        bool
	name            string
	descriptor      string
}

func (this *RuntimeConstantNameAndTypeInfo) resolve() RuntimeConstantPoolInfo  {
	if !this.resolved {
		this.name = this.hostClass.constantPool[this.nameIndex].resolve().(*RuntimeConstantUtf8Info).value
		this.descriptor = this.hostClass.constantPool[this.descriptorIndex].resolve().(*RuntimeConstantUtf8Info).value
		this.resolved = true
	}
	return this
}

func (this *RuntimeConstantNameAndTypeInfo) toString() string  {
	return this.name + this.descriptor
}

/*
CONSTANT_Utf8_info {
    u1 tag;
    u2 length;
    u1 bytes[length];
}
 */
type RuntimeConstantUtf8Info struct {
	hostClass       *JavaClass
	length          u2
	bytes           []u1
	resolved        bool
	value           string
}

func (this *RuntimeConstantUtf8Info) resolve() RuntimeConstantPoolInfo {
	if !this.resolved {
		this.value = u2s(this.bytes)
		this.resolved = true
	}
	return this
}

/*
CONSTANT_MethodHandle_info {
    u1 tag;
    u1 reference_kind;
    u2 reference_index;
}
 */
type RuntimeConstantMethodHandleInfo struct {
	hostClass       *JavaClass
	referenceKind   u1
	referenceIndex  u2
	resolved        bool
	//TODO
}

func (this *RuntimeConstantMethodHandleInfo) resolve() RuntimeConstantPoolInfo {
	if !this.resolved {
		//TODO
	}
	return this
}

/*
CONSTANT_MethodType_info {
    u1 tag;
    u2 descriptor_index;
}
 */
type RuntimeConstantMethodTypeInfo struct {
	hostClass       *JavaClass
	descriptorIndex u2
	resolved        bool
	descriptor      string
}

func (this *RuntimeConstantMethodTypeInfo) resolve() RuntimeConstantPoolInfo {
	if !this.resolved {
		this.descriptor = this.hostClass.constantPool[this.descriptorIndex].resolve().(*RuntimeConstantUtf8Info).value
		this.resolved = true
	}
	return this
}

/*
CONSTANT_InvokeDynamic_info {
    u1 tag;
    u2 bootstrap_method_attr_index;
    u2 name_and_type_index;
}
 */
type RuntimeConstantInvokeDynamicInfo struct {
	hostClass                *JavaClass
	bootstrapMethodAttrIndex u2
	nameAndTypeIndex         u2
	resolved                 bool
	//TODO
}


func (this *RuntimeConstantInvokeDynamicInfo) resolve() RuntimeConstantPoolInfo {
	if !this.resolved {
		//TODO
	}
	return this
}