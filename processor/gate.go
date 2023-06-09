package processor

import "kalehla/types"

type Gate struct {
	generalAccessPermissions *AccessPermissions
	specialAccessPermissions *AccessPermissions
	libFlag                  bool
	gotoInhibit              bool
	designatorInhibit        bool
	accessKeyInhibit         bool
	latentParameter0Inhibit  bool
	latentParameter1Inhibit  bool
	accessLock               *AccessLock
	targetLevel              uint
	targetBDI                uint
	targetOffset             uint
	basicModeBaseRegister    uint                // basic mode register is actually this field + 12
	designatorRegisterValue  *DesignatorRegister //	only bits 12-17 are significant
	newAccessKey             *AccessKey
	latentParameterValue0    uint64
	latentParameterValue1    uint64
}

func NewGateFromStorage(buffer []types.Word36) *Gate {
	g := Gate{}
	g.generalAccessPermissions = NewAccessPermissions(buffer[0]&0_400000_000000 != 0, false, false)
	g.specialAccessPermissions = NewAccessPermissions(buffer[0]&0_040000_000000 != 0, false, false)
	g.libFlag = buffer[0]&0_040_000000 != 0
	g.gotoInhibit = buffer[0]&0_020_000000 != 0
	g.designatorInhibit = buffer[0]&0_010_000000 != 0
	g.accessKeyInhibit = buffer[0]&0_004_000000 != 0
	g.latentParameter1Inhibit = buffer[0]&0_002_000000 != 0
	g.latentParameter1Inhibit = buffer[0]&0_001_000000 != 0
	g.accessLock = NewAccessLockFromComposite(uint(buffer[0] & 0777777))
	g.targetLevel = uint(buffer[1] >> 33)
	g.targetBDI = uint(buffer[1]>>18) & 077777
	g.targetOffset = uint(buffer[1]) & 0777777
	g.basicModeBaseRegister = uint(buffer[2]>>24) & 03
	g.designatorRegisterValue = NewDesignatorRegisterFromComposite(uint64(buffer[2]) & 0_000077_000000)
	g.newAccessKey = NewAccessKeyFromComposite(uint(buffer[2] & 0777777))
	g.latentParameterValue0 = uint64(buffer[3])
	g.latentParameterValue1 = uint64(buffer[4])
	return &g
}
