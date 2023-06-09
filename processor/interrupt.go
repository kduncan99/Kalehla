package processor

import "kalehla/types"

const (
	HardwareDefaultInterruptClass uint = iota
	HardwareCheckInterruptClass
	DiagnosticInterruptClass
	ReferenceViolationInterruptClass = iota + 5
	AddressingExceptionInterruptClass
	TerminalAddressingExceptionInterruptClass
	RCSGenericStackUnderOverflowInterruptClass
	SignalInterruptClass
	TestAndSetInterruptClass
	InvalidInstructionInterruptClass
	PageExceptionInterruptClass
	ArithmeticExceptionInterruptClass
	DataExceptionInterruptClass
	OperationTrapInterruptClass
	BreakpointInterruptClass
	QuantumTimerInterruptClass
	PageZeroedInterruptClass = iota + 7
	SoftwareBreakInterruptClass
	JumpHistoryFullInterruptClass
	DayClockInterruptClass = iota + 8
	PerformanceMonitoringInterruptClass
	IPLInterruptClass
	UPIInitialInterruptClass
	UPINormalInterruptClass
)

const (
	ReferenceViolationGRS           = 00
	ReferenceViolationStorageLimits = 01
	ReferenceViolationReadAccess    = 02
	ReferenceViolationWriteAccess   = 03
)

const (
	AddressingExceptionFatal                            = 00
	AddressingExceptionGateGBitSet                      = 01
	AddressingExceptionEnterAccessDenied                = 02
	AddressingExceptionInvalidSourceLBDI                = 03
	AddressingExceptionGateBankBoundaryViolation        = 04
	AddressingExceptionInvalidISValue                   = 05
	AddressingExceptionGOTOInhibit                      = 06
	AddressingExceptionGeneralQueuingViolation          = 07
	AddressingExceptionMaxCountEnq                      = 010
	AddressingExceptionIndirectGBitSet                  = 011
	AddressingExceptionInactiveQueuebDListEmpty         = 013
	AddressingExceptionUpdateInProgress                 = 014
	AddressingExceptionQueueBankRepositoryFull          = 015
	AddressingExceptionBDTypeInvalid                    = 016
	AddressingExceptionAccessDeniedPosternOrDataExpanse = 017
	//	There are others...
)

const (
	RCSGenericStackOverflow  = 00
	RCSGenericStackUnderflow = 01
)

const (
	InvalidInstructionBadFunctionCode  = 00
	InvalidInstructionX0Linkage        = 00
	InvalidInstructionLBUUsesB0OrB1    = 00
	InvalidInstructionLBUDUsesB0       = 00
	InvalidInstructionBadPP            = 01
	InvalidInstructionEXRInvalidTarget = 03
)

type Interrupt interface {
	GetClass() uint
	GetShortStatusField() uint
	GetStatusWord0() types.Word36
	GetStatusWord1() types.Word36
	IsDeferrable() bool
}

// Class 8 Reference Violation -----------------------------------------------------------------------------------------

// ssf values:
//	bits 0-1: Entry Type
//				0: GRS violation with insufficient PP (see 2.3.7)
//					JGD j-field concatenated with a-field is a GRS location
//					SRS, LRS a-field indicates a GRS address
//					All other GRS locations developed as an instructionType operand caused by any instructions
//						other than JGD, SRS, or LRS
//				1: Storage Limits violation
//				2: Read Access violation
//				3: Write Access violation
//  bits 2-4: reserved
//	bits 5: true if this occurred during an instructionType fetch

type ReferenceViolationInterrupt struct {
	shortStatusField uint
}

func (i *ReferenceViolationInterrupt) GetClass() uint {
	return ReferenceViolationInterruptClass
}

func (i *ReferenceViolationInterrupt) GetShortStatusField() uint {
	return i.shortStatusField
}

func (i *ReferenceViolationInterrupt) GetStatusWord0() types.Word36 {
	return 0
}

func (i *ReferenceViolationInterrupt) GetStatusWord1() types.Word36 {
	return 0
}

func (i *ReferenceViolationInterrupt) IsDeferrable() bool {
	return false
}

func NewReferenceViolationInterrupt(entryType uint, fetchOperation bool) *ReferenceViolationInterrupt {
	ssf := (entryType & 03) << 4
	if fetchOperation {
		ssf |= 01
	}
	return &ReferenceViolationInterrupt{
		shortStatusField: ssf,
	}
}

// Class 8 Addressing Exception ----------------------------------------------------------------------------------------

// ssf values:
//	000 Fatal addressing exception
//	001 G-bit set in gate bank descriptor
//	002 Enter access denied by gate bank descriptor or by gate, or queuing instruction access denied
//	003 invalid source L,BDI or BDT limit error for L,BDI supplied by user instruction
//  004 gate bank boundary violation or gate input offset not within gate bd limits
//	005 invalid IS value
//	006 GOTO inhibit set in gate
//	007 General queuing instruction violation
//	010 MaxCount exceeded on ENQ/ENQF
//	011 G-bit set in indirect bank descriptor
//	013 Inactive QBD list empty on DEQ/DEQW
//	014 Update in progress set in queue structure

type AddressingExceptionInterrupt struct {
	shortStatusField     uint
	interruptStatusWord1 types.Word36
}

func (i *AddressingExceptionInterrupt) GetClass() uint {
	return AddressingExceptionInterruptClass
}

func (i *AddressingExceptionInterrupt) GetShortStatusField() uint {
	return i.shortStatusField
}

func (i *AddressingExceptionInterrupt) GetStatusWord0() types.Word36 {
	return 0
}

func (i *AddressingExceptionInterrupt) GetStatusWord1() types.Word36 {
	return i.interruptStatusWord1
}

func (i *AddressingExceptionInterrupt) IsDeferrable() bool {
	return false
}

func NewAddressingExceptionInterrupt(shortStatusField uint, sourceBankLevel uint, sourceBankDescriptorIndex uint) *AddressingExceptionInterrupt {
	isw1 := types.Word36(sourceBankLevel&07) << 33
	isw1 |= types.Word36(sourceBankDescriptorIndex&077777) << 18
	return &AddressingExceptionInterrupt{
		shortStatusField:     shortStatusField,
		interruptStatusWord1: isw1,
	}
}

// Class 11 RCS/Generic Stack Under/Overflow ---------------------------------------------------------------------------

// ssf values:
//
//	0 Generic stack or RCS overflow
//	1 Generic stack or RCS underrflow
//
// ISW0:
//	Bits 0-5 (S1): BREG (base register causing trouble) - when the RCS causes the interrupt, BREG will be 25
//  Bits 12-35:    Relative address (n/a for BREG 25)
//                  When BREG != 25 and ssf == 0, this field contains Xm - Xi - d of the X register specified
//                      by the BUY instruction
//                  When BREG != 25 and ssf == 1, this field contains Xm of the X register specified
//                      by the SELL instruction

type RCSGenericStackUnderOverflowInterrupt struct {
	shortStatusField     uint
	interruptStatusWord0 types.Word36
}

func (i *RCSGenericStackUnderOverflowInterrupt) GetClass() uint {
	return RCSGenericStackUnderOverflowInterruptClass
}

func (i *RCSGenericStackUnderOverflowInterrupt) GetShortStatusField() uint {
	return i.shortStatusField
}

func (i *RCSGenericStackUnderOverflowInterrupt) GetStatusWord0() types.Word36 {
	return i.interruptStatusWord0
}

func (i *RCSGenericStackUnderOverflowInterrupt) GetStatusWord1() types.Word36 {
	return 0
}

func (i *RCSGenericStackUnderOverflowInterrupt) IsDeferrable() bool {
	return false
}

func NewRCSGenericStackUnderOverflowInterrupt(shortStatusField uint, baseRegister uint, relativeAddress uint) *RCSGenericStackUnderOverflowInterrupt {
	isw0 := (types.Word36(baseRegister) << 30) | types.Word36(relativeAddress)
	return &RCSGenericStackUnderOverflowInterrupt{
		shortStatusField:     shortStatusField,
		interruptStatusWord0: isw0,
	}
}

// Class 14 Invalid Instruction ----------------------------------------------------------------------------------------

// ssf values:
//
//	0 function code not defined, direct execution or as a target of EXR
//		or LBJ/LIJ/LDJ uses X0
//		or LBU uses B0 or B1
//		or LBUD uses B0
//	1 insufficient processor privilege
//	3 EXR target invalid (other than as above for value 0)
//	4 compatibility trap (we don't do this)

type InvalidInstructionInterrupt struct {
	shortStatusField uint
}

func (i *InvalidInstructionInterrupt) GetClass() uint {
	return InvalidInstructionInterruptClass
}

func (i *InvalidInstructionInterrupt) GetShortStatusField() uint {
	return i.shortStatusField
}

func (i *InvalidInstructionInterrupt) GetStatusWord0() types.Word36 {
	return 0
}

func (i *InvalidInstructionInterrupt) GetStatusWord1() types.Word36 {
	return 0
}

func (i *InvalidInstructionInterrupt) IsDeferrable() bool {
	return false
}

func NewInvalidInstructionInterrupt(shortStatusField uint) *InvalidInstructionInterrupt {
	return &InvalidInstructionInterrupt{
		shortStatusField: shortStatusField,
	}
}
