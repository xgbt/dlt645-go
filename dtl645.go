/*
Package dlt645 provides a client for dlt645
*/
package dlt645

import "fmt"

const (
	// len limit 5-Bit
	FunctionCodeReadData                  = 0x11 // binary 0001 0001
	FunctionCodeReadFollowUpData          = 0x12 // binary 0001 0010
	FunctionCodeWriteData                 = 0x14 // binary 0001 0100
	FunctionCodeReadCommunicationAddress  = 0x13 // binary 0001 0011
	FunctionCodeWriteCommunicationAddress = 0x15 // binary 0001 0101
	FunctionCodeBroadcastTiming           = 0x8  // binary 0000 1000
	FunctionCodeFreezeCommand             = 0x16 // binary 0001 0110
	FunctionCodeChangeCommunicationRate   = 0x17 // binary 0001 0111
	FunctionCodeChangePassword            = 0x18 // binary 0001 1000
	FunctionCodeClearMaximumDemand        = 0x19 // binary 0001 1001
	FunctionCodeClearAmmeter              = 0x1A // binary 0001 1010
	FunctionCodeClearEvent                = 0x1B // binary 0001 1011
)

const (
	ExceptionCodeRatesExceedsLimit              = 0x40 // binary 0100 0000 费率数超过限制
	ExceptionCodeDayPeriodsExceedsThreshold     = 0x20 // binary 0010 0000 日时段数超
	ExceptionCodeTimeZonesYearExceedsThreshold  = 0x10 // binary 0001 0000 年时区数超
	ExceptionCodeCommunicationRateCannotChanged = 0x08 // binary 0000 1000 通讯速率不能更改
	ExceptionCodeIllegalPassword                = 0x04 // binary 0000 0100 密码错误或者没有权限
	ExceptionCodeRequestWithoutData             = 0x02 // binary 0000 0010 请求无数据
	ExceptionCodeOtherError                     = 0x01 // binary 0000 0001 其他错误
)

const (
	FrameHead = 0x68
	FrameTail = 0x16
)

const (
	BroadcastAddressDomain = 0x999999999999
)

// DLTError implements error interface
type DltError struct {
	FunctionCode  byte
	ExceptionCode byte
}

// Error converts known dlt645 exception code to error message
func (e *DltError) Error() string {
	var name string
	switch e.ExceptionCode {
	case ExceptionCodeRatesExceedsLimit:
		name = "The number of rates exceeds the limit"
	case ExceptionCodeDayPeriodsExceedsThreshold:
		name = "The number of day periods exceeds the threshold"
	case ExceptionCodeTimeZonesYearExceedsThreshold:
		name = "The number of time zones in the year exceeds the threshold"
	case ExceptionCodeCommunicationRateCannotChanged:
		name = "The communication rate cannot be changed"
	case ExceptionCodeIllegalPassword:
		name = "Incorrect password or no permission"
	case ExceptionCodeRequestWithoutData:
		name = "Request without data"
	case ExceptionCodeOtherError:
		name = "Other error"
	default:
		name = "Unknown"
	}
	return fmt.Sprintf("dlt645: exception '%v' (%s), function '%v'", e.ExceptionCode, name, e.FunctionCode)
}

// controlCode
// 8 bit   : 0 master send   1 slave send
// 7 biy   : 0 slave ok   1 slave err
// 6 bit   : 0 have not follow-up data    1 have follow-up data
// 1-5 bit : function code
type FramePayLoad struct {
	HasFollowUpData bool
	FunctionCode    byte
	Data            []byte
}

// Packager specifies the communication layer.
type Packager interface {
	Encode(frame *FramePayLoad) (adu []byte, err error)
	Decode(adu []byte) (frame *FramePayLoad, err error)
	Verify(aduRequest []byte, aduResponse []byte) (err error)
}

// Transporter specifies the transport layer.
type Transporter interface {
	Send(request []byte) (response []byte, err error)
	SendNotResponse(request []byte) (err error)
}
