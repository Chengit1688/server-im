package model

const (
	RTCRequestExpireTime = 60
	RTCRetainExpireTime  = 15
)

type RTCType int

const (
	RTCTypeNone  RTCType = 0
	RTCTypeAudio RTCType = 1
	RTCTypeVideo RTCType = 2
)

type RTCStatusType int

const (
	RTCStatusTypeNone        RTCStatusType = 0
	RTCStatusTypeRequest     RTCStatusType = 1
	RTCStatusTypeCancel      RTCStatusType = 2
	RTCStatusTypeAgree       RTCStatusType = 3
	RTCStatusTypeDisagree    RTCStatusType = 4
	RTCStatusTypeFinish      RTCStatusType = 5
	RTCStatusTypeNotResponse RTCStatusType = 6
	RTCStatusTypeBusy        RTCStatusType = 7
	RTCStatusTypeAbort       RTCStatusType = 8
)

type RTCOperationType int

const (
	RTCOperationTypeNone        RTCOperationType = 0
	RTCOperationTypeCancel      RTCOperationType = 1
	RTCOperationTypeAgree       RTCOperationType = 2
	RTCOperationTypeDisagree    RTCOperationType = 3
	RTCOperationTypeFinish      RTCOperationType = 4
	RTCOperationTypeSwitch      RTCOperationType = 5
	RTCOperationTypeNotResponse RTCOperationType = 6
	RTCOperationTypeAbort       RTCOperationType = 7
)
