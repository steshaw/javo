package jago

import (
	"unsafe"
	"strings"
	"math"
)

type ClassLoader struct {
	classPath     *ClassPath
	classCache    map[string]*Class
	parent        *ClassLoader
}

func NewClassLoader(str string, parent *ClassLoader) *ClassLoader {
	return &ClassLoader{
		NewClassPath(str),
		make(map[string]*Class),
		parent}
}

func (this *ClassLoader) CreateClass(className string) *Class {
	return this.internalCreateClass(className, true)
}

func (this *ClassLoader) internalCreateClass(className string, requireInitialize bool) *Class {
	if class, found := this.classCache[className]; found {
		return class
	}

	CLASSLOAD_LOG.Debug(__times(__indention, "   ") + "↳ %s \n", className)
	__indention++

	var class *Class
	if string(className[0]) == JVM_SIGNATURE_ARRAY {
		class = this.defineArrayClass(className)
	} else {
		clazz := this.loadClass(className)
		class = clazz
	}

	// attach a java.lang.Class object
	// set classloader
	class.classLoader = this
	class.classObject = NewJavaLangClass(class)

	// eager linkage
	this.link(class)

	// initialize only all the super class and interface prepared
	// now can initialize together
	if requireInitialize {
		this.initialize(class)
	}

	__indention--
	return class
}

var __indention = 0
func __times(t int, str string) string {
	ret := ""
	for i := 0; i < t; i++ {
		ret += str
	}
	return ret
}

func (this *ClassLoader) defineArrayClass(className string) *Class {

	arrayClass := &Class{
			name: className,
			superClassName: "java/lang/Object",
			interfaceNames: []string{"java/io/Serializable", "java/lang/Cloneable"}}

	this.classCache[className] = arrayClass

	arrayClass.accessFlags = 0
	arrayClass.superClass = this.internalCreateClass("java/lang/Object", false)
	arrayClass.interfaces = []*Class{
		this.internalCreateClass("java/io/Serializable", false),
		this.internalCreateClass("java/lang/Cloneable", false)}

	componentDescriptor := string(className[1])
	switch componentDescriptor {
	case JVM_SIGNATURE_BYTE:
		{
			arrayClass.componentType = BYTE_TYPE
			arrayClass.elementType = BYTE_TYPE
		}
	case JVM_SIGNATURE_CHAR:
		{
			arrayClass.componentType = CHAR_TYPE
			arrayClass.elementType = CHAR_TYPE
		}
	case JVM_SIGNATURE_SHORT:
		{
			arrayClass.componentType = SHORT_TYPE
			arrayClass.elementType = SHORT_TYPE
		}
	case JVM_SIGNATURE_INT:
		{
			arrayClass.componentType = INT_TYPE
			arrayClass.elementType = INT_TYPE
		}
	case JVM_SIGNATURE_LONG:
		{
			arrayClass.componentType = LONG_TYPE
			arrayClass.elementType = LONG_TYPE
		}
	case JVM_SIGNATURE_FLOAT:
		{
			arrayClass.componentType = FLOAT_TYPE
			arrayClass.elementType = FLOAT_TYPE
		}
	case JVM_SIGNATURE_DOUBLE:
		{
			arrayClass.componentType = DOUBLE_TYPE
			arrayClass.elementType = DOUBLE_TYPE
		}
	case JVM_SIGNATURE_BOOLEAN:
		{
			arrayClass.componentType = BOOLEAN_TYPE
			arrayClass.elementType = BOOLEAN_TYPE
		}
	case JVM_SIGNATURE_CLASS:
		{
			arrayClass.componentType = BOOTSTRAP_CLASSLOADER.CreateClass(className[2:len(className)-1])
			arrayClass.elementType = arrayClass.componentType
		}
	case JVM_SIGNATURE_ARRAY:
		{
			arrayClass.componentType = BOOTSTRAP_CLASSLOADER.CreateClass(className[1:])
			arrayClass.elementType = arrayClass.componentType.(*Class).elementType
		}
	}
	dimensions := 1
	class := arrayClass
	for ;; {
		componentType, ok := class.componentType.(*Class)
		if ok {
			class = componentType
			dimensions++
		} else {
			break
		}
	}
	arrayClass.dimensions = dimensions

	arrayClass.defined = true
	return arrayClass
}

func (this *ClassLoader) loadClass(className string) *Class  {
	class := this.findClass(className)

	// TODO delegation

	// set classloader
	class.classLoader = this
	return class
}

func (this *ClassLoader) findClass(className string) *Class  {
	bytecode, err := this.classPath.ReadClass(className)
	if err != nil {
		VM_throw("java.lang.ClassNotFoundException", className)
	}

	//If L creates C directly, we say that L defines C
	class := this.defineClass(bytecode)
	return class
}

func (this *ClassLoader) defineClass(bytecode []byte) *Class  {
	classfile := NewClassFile(bytecode)

	class := &Class{}

	class.accessFlags = uint16(classfile.accessFlags)
	class.name = classfile.cpUtf8(classfile.constantPool[classfile.thisClass].(*ConstantClassInfo).nameIndex)
	// add to classcache
	this.classCache[class.name] = class

	if classfile.superClass == 0 {
		class.superClassName = ""
	} else {
		class.superClassName = classfile.cpUtf8(classfile.constantPool[classfile.superClass].(*ConstantClassInfo).nameIndex)
	}

	//class.superClass = ??
	class.interfaceNames = make([]string, len(classfile.interfaces))
	for i, interfaceIndex := range classfile.interfaces {
		class.interfaceNames[i] = classfile.cpUtf8(classfile.constantPool[interfaceIndex].(*ConstantClassInfo).nameIndex)
	}
	//class.interfaces = ??

	constantPool := make([]Constant, classfile.constantPoolCount+1)
	class.constantPool = constantPool

	// start loading
	for i, _ := range classfile.constantPool {
		if class.name == "_000_util/TreePrinter" && i == 25 {
			print("ss")
		}
		constInfo := classfile.constantPool[i]
		var constant Constant
		switch constInfo.(type) {
		case *ConstantClassInfo:
			constantClassInfo := constInfo.(*ConstantClassInfo)
			constant = &ClassRef{
				class,
				classfile.cpUtf8(constantClassInfo.nameIndex),
				nil}
		case *ConstantFieldrefInfo:
			constantFieldrefInfo := constInfo.(*ConstantFieldrefInfo)
			nameAndTypeInfo := classfile.constantPool[constantFieldrefInfo.nameAndTypeIndex].(*ConstantNameAndTypeInfo)
			name := classfile.cpUtf8(nameAndTypeInfo.nameIndex)
			descriptor := classfile.cpUtf8(nameAndTypeInfo.descriptorIndex)
			className := classfile.cpUtf8(classfile.constantPool[constantFieldrefInfo.classIndex].(*ConstantClassInfo).nameIndex)
			constant = &FieldRef{
				MemberRef{
					class,
					className,
					nil,
					name,
				descriptor},
				nil}
		case *ConstantMethodrefInfo:
			constantMethodrefInfo := constInfo.(*ConstantMethodrefInfo)
			nameAndTypeInfo := classfile.constantPool[constantMethodrefInfo.nameAndTypeIndex].(*ConstantNameAndTypeInfo)
			name := classfile.cpUtf8(nameAndTypeInfo.nameIndex)
			descriptor := classfile.cpUtf8(nameAndTypeInfo.descriptorIndex)
			className := classfile.cpUtf8(classfile.constantPool[constantMethodrefInfo.classIndex].(*ConstantClassInfo).nameIndex)
			constant = &MethodRef{
				MemberRef{
					class,
					className,
					nil,
					name,
					descriptor},
				nil}
		case *ConstantInterfaceMethodrefInfo:
			constantInterfaceMethodrefInfo := constInfo.(*ConstantInterfaceMethodrefInfo)
			nameAndTypeInfo := classfile.constantPool[constantInterfaceMethodrefInfo.nameAndTypeIndex].(*ConstantNameAndTypeInfo)
			name := classfile.cpUtf8(nameAndTypeInfo.nameIndex)
			descriptor := classfile.cpUtf8(nameAndTypeInfo.descriptorIndex)
			className := classfile.cpUtf8(classfile.constantPool[constantInterfaceMethodrefInfo.classIndex].(*ConstantClassInfo).nameIndex)
			constant = &InterfaceMethodRef{
				MemberRef{
					class,
					className,
					nil,
					name,
					descriptor},
				nil}
		case *ConstantUtf8Info:
			constantUtf8Info := constInfo.(*ConstantUtf8Info)
			constant = &UTF8Constant{class,u2s(constantUtf8Info.bytes)}
		case *ConstantStringInfo:
			constantStringInfo := constInfo.(*ConstantStringInfo)
			constant = &StringConstant{class, classfile.cpUtf8(constantStringInfo.stringIndex), NULL}
		case *ConstantIntegerInfo:
			constantIntegerInfo := constInfo.(*ConstantIntegerInfo)
			constant = &IntegerConstant{class, Int(constantIntegerInfo.bytes)}
		case *ConstantFloatInfo:
			constantFloatInfo := constInfo.(*ConstantFloatInfo)
			constant = &FloatConstant{class, Float(constantFloatInfo.bytes)}
		case *ConstantLongInfo:
			constantLongInfo := constInfo.(*ConstantLongInfo)
			l := int64(uint64(constantLongInfo.highBytes) << 32 + uint64(constantLongInfo.lowBytes))
			constant = &LongConstant{class, Long(l)}
		case *ConstantDoubleInfo:
			constantDoubleInfo := constInfo.(*ConstantDoubleInfo)
			bits := uint64(constantDoubleInfo.highBytes) << 32 + uint64(constantDoubleInfo.lowBytes)
			d := math.Float64frombits(bits)
			//var d float64
			//if l == 0x7ff0000000000000 {
			//	d = math.Inf(1)
			//} else if l == 0xfff0000000000000 {
			//	d = math.Inf(-1)
			//} else if (l >= 0x7ff0000000000001 && l <= 0x7fffffffffffffff) || (l >= 0xfff0000000000001 && l <= 0xffffffffffffffff) {
			//	d = math.NaN()
			//} else {
			//	s := -1
			//	if ((l >> 63) == 0) {
			//		s = 1
			//	}
			//	e := (l >> 52) & 0x7ff
			//	m := (l & 0xfffffffffffff) | 0x10000000000000
			//	if e == 0 {
			//		m = (l & 0xfffffffffffff) << 1
			//	}
			//	d = s * m * math.2e-1075
			//}
			constant = &DoubleConstant{class, Double(d)}
		case *ConstantNameAndTypeInfo:
			constantNameAndTypeInfo := constInfo.(*ConstantNameAndTypeInfo)
			constant = &NameAndTypeConstant{class,classfile.cpUtf8(constantNameAndTypeInfo.nameIndex),classfile.cpUtf8(constantNameAndTypeInfo.descriptorIndex)}
		case *ConstantMethodTypeInfo:
			constantMethodTypeInfo := constInfo.(*ConstantMethodTypeInfo)
			constant = &MethodTypeConstant{class,classfile.cpUtf8(constantMethodTypeInfo.descriptorIndex)}
		case *ConstantInvokeDynamicInfo:
			constantInvokeDynamicInfo := constInfo.(*ConstantInvokeDynamicInfo)
			nameAndTypeInfo := classfile.constantPool[constantInvokeDynamicInfo.nameAndTypeIndex].(*ConstantNameAndTypeInfo)
			name := classfile.cpUtf8(nameAndTypeInfo.nameIndex)
			descriptor := classfile.cpUtf8(nameAndTypeInfo.descriptorIndex)
			constant = &InvokeDynamicConstant{class,
				"",//TODO classfile.cpUtf8(constantInvokeDynamicInfo.bootstrapMethodAttrIndex),
				name,
				descriptor}
		}

		class.constantPool[i] = constant
	}

	for k := u2(0); k < classfile.attributeCount; k++ {
		attributeInfo := classfile.attributes[k]
		switch attributeInfo.(type) {
		case *SourceFileAttribue:
			sourceFileAttribue := attributeInfo.(*SourceFileAttribue)
			class.sourceFile = classfile.cpUtf8(sourceFileAttribue.sourceFileIndex)
		}
	}


	class.fields = make([]*Field, len(classfile.fields))
	//class.staticVars = ??

	for i, fieldInfo := range classfile.fields {
		field := &Field{}
		field.accessFlags = uint16(fieldInfo.accessFlags)
		field.name = classfile.cpUtf8(fieldInfo.nameIndex)
		field.descriptor = classfile.cpUtf8(fieldInfo.descriptorIndex)
		field.class = class
		//field.index = ??
		class.fields[i] = field
	}

	class.methods = make([]*Method, len(classfile.methods))

	for i, methodInfo := range classfile.methods {
		method := &Method{}
		method.accessFlags = uint16(methodInfo.accessFlags)
		method.name = classfile.cpUtf8(methodInfo.nameIndex)
		method.descriptor = classfile.cpUtf8(methodInfo.descriptorIndex)
		method.class = class

		for j := u2(0); j < methodInfo.attributeCount; j++ {
			attributeInfo := methodInfo.attributes[j]
			switch attributeInfo.(type) {
			case *CodeAttribute:
				codeAttribute := attributeInfo.(*CodeAttribute)
				method.maxStack = uint(codeAttribute.maxStack)
				method.maxLocals = uint(codeAttribute.maxLocals)
				method.code = *(*[]uint8)(unsafe.Pointer(&codeAttribute.code))
				method.exceptions = make([]*ExceptionHandler, codeAttribute.exceptionTableLength)
				for ei, ete := range codeAttribute.exceptionTable {
					exceptionHandler := &ExceptionHandler{int(ete.startPc), int(ete.endPc), int(ete.handlerPc), int(ete.catchType)}
					method.exceptions[ei] = exceptionHandler
				}
				for k := u2(0); k < codeAttribute.attributesCount; k++ {
					codeAttributeAttribute := codeAttribute.attributes[k]
					switch codeAttributeAttribute.(type) {
					case *LocalVariableTableAttribute:
						localVariableTableAttribute := codeAttributeAttribute.(*LocalVariableTableAttribute)
						method.localVars = make([]*LocalVariable, localVariableTableAttribute.localVariableTableLength)
						for m := u2(0); m < localVariableTableAttribute.localVariableTableLength; m++ {
							method.localVars[m] = &LocalVariable{
								method,
								uint16(localVariableTableAttribute.localVariableTable[m].startPc),
								uint16(localVariableTableAttribute.localVariableTable[m].length),
								uint16(localVariableTableAttribute.localVariableTable[m].index),
								classfile.cpUtf8(localVariableTableAttribute.localVariableTable[m].nameIndex),
								classfile.cpUtf8(localVariableTableAttribute.localVariableTable[m].descriptorIndex)}
						}
					case *LineNumberTableAttribute:
						lineNumberTableAttribute := codeAttributeAttribute.(*LineNumberTableAttribute)
						method.lineNumbers = make([]*LineNumber, lineNumberTableAttribute.lineNumberTableLength)
						for i, lineNumberTableEntry := range lineNumberTableAttribute.lineNumberTable {
							method.lineNumbers[i] = &LineNumber{int(lineNumberTableEntry.startPc), int(lineNumberTableEntry.lineNumber)}
						}
					}
				}

			}
		}

		// start parse parameters and return
		arr := strings.Split(method.descriptor, ")")
		parametersStr := arr[0][1:]
		returnDescritor := arr[1]

		var parameterDescriptors []string

		for i := 0; i < len(parametersStr); {
			ch := rune(parametersStr[i])
			switch ch {
			case 'B', 'C', 'D', 'F', 'I', 'J', 'S', 'Z':
				parameterDescriptors = append(parameterDescriptors, string(ch))
				i++
			case 'L':
			Ref: for j := i+1; j < len(parametersStr); j++ {
				switch rune(parametersStr[j]) {
				case ';':
					parameterDescriptors = append(parameterDescriptors, string(parametersStr[i:j+1]))
					i = j+1
					break Ref
				}
			}
			case '[':
			Arr: for j := i+1; j < len(parametersStr); j++ {
				switch rune(parametersStr[j]) {
				case '[':
					continue
				case 'B', 'C', 'D', 'F', 'I', 'J', 'S', 'Z':
					parameterDescriptors = append(parameterDescriptors, string(parametersStr[i:j+1]))
					i = j+1
					break Arr
				case 'L':
					for k := j+1; j < len(parametersStr); k++ {
						switch rune(parametersStr[k]) {
						case ';':
							parameterDescriptors = append(parameterDescriptors, string(parametersStr[i:k+1]))
							i = k+1
							break Arr
						}
					}
				}
			}
			}
		}
		method.parameterDescriptors = parameterDescriptors
		method.returnDescriptor = returnDescritor

		class.methods[i] = method
	}


	// resolve super class
	if class.superClassName != "" {
		class.superClass = this.internalCreateClass(class.superClassName, false)
	}

	// calculate static variables and instance variable count
	// must be immediately after resolving super class
	maxInstanceVars := 0  // include all the ancestry
	maxStaticVars := 0

	if class.superClass != nil {
		maxInstanceVars = class.superClass.maxInstanceVars
	}
	for _, field := range class.fields {
		if field.IsStatic() {
			field.index = maxStaticVars
			maxStaticVars++
		} else {
			field.index = maxInstanceVars
			maxInstanceVars++
		}
	}
	class.maxInstanceVars = maxInstanceVars
	class.maxStaticVars = maxStaticVars

	// resolve interfaces
	class.interfaces = make([]*Class, len(class.interfaceNames))
	for i, interfaceName := range class.interfaceNames {
		class.interfaces[i] = this.internalCreateClass(interfaceName, false)
	}

	return class
}

func (this *ClassLoader) link(class *Class)  {
	if class.linked {
		return
	}
	this.verify(class)
	this.prepare(class)

	class.linked = true
	// we resolve each symbolic class in a class or interface individually when it is used ("lazy" or "late" resolution)
	// So SymbolRef all implements a method PrimitiveType resolve()
	//for _, constant := range this.constantPool {
	//	switch constant.(type) {
	//	case SymbolRef:
	//		constant.(SymbolRef).resolve()
	//	}
	//}
	//this.resolve(class)
}

// invoke <clinit> to execute initialization code
func (this *ClassLoader) initialize(class *Class) {
	if class.initialized {
		return
	}

	// initialize its supper class and interfaces first
	if class.superClass != nil {
		this.initialize(class.superClass)
	}
	for _, iface := range class.interfaces {
		this.initialize(iface)
	}

	thread := THREAD_MANAGER.currentThread

	clinit := class.GetMethod("<clinit>", "()V")
	if clinit != nil {
		//inStack := false
		//for _, frame := range thread.vmStack {
		//	if frame.method == clinit {
		//		inStack = true
		//		break
		//	}
		//}
		//if !inStack {

		// always initialize super class
		if class.superClass != nil {
			this.initialize(class.superClass)
		}
		CLASSLOAD_LOG.Debug(__times(__indention, "   ") + "%s \n", clinit.Qualifier())
			VM_invokeMethod(thread, clinit)
		//}
	}

	class.initialized = true
}

func (this *ClassLoader) verify(class *Class) {
	//TODO
}

// initialize static variables to default values: no need to execute code
func (this *ClassLoader) prepare(class *Class)  {
	// Initialize static variables
	class.staticVars = make([]Value, class.maxStaticVars)
	for _, field := range class.fields {
		if field.IsStatic() {
			class.staticVars[field.index] = field.defaultValue()
		}
	}
}






