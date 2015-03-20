package instructions

import (
	"github.com/zxh0/jvm.go/jvmgo/jvm/rtda"
	rtc "github.com/zxh0/jvm.go/jvmgo/jvm/rtda/class"
)

// Fetch field from object
type getfield struct {
	Index16Instruction
	field *rtc.Field
}

func (self *getfield) Execute(frame *rtda.Frame) {
	if self.field == nil {
		cp := frame.Method().ConstantPool()
		kFieldRef := cp.GetConstant(self.index).(*rtc.ConstantFieldref)
		self.field = kFieldRef.InstanceField()
	}

	stack := frame.OperandStack()
	ref := stack.PopRef()
	if ref == nil {
		frame.Thread().ThrowNPE()
		return
	}

	val := self.field.GetValue(ref)
	stack.Push(val)
}

// Get static field from class
type getstatic struct {
	Index16Instruction
	field *rtc.Field
}

func (self *getstatic) Execute(frame *rtda.Frame) {
	if self.field == nil {
		cp := frame.Method().Class().ConstantPool()
		kFieldRef := cp.GetConstant(self.index).(*rtc.ConstantFieldref)
		self.field = kFieldRef.StaticField()
	}

	class := self.field.Class()
	if class.InitializationNotStarted() {
		frame.RevertNextPC() // undo getstatic
		frame.Thread().InitClass(class)
		return
	}

	val := self.field.GetStaticValue()
	frame.OperandStack().Push(val)
}
