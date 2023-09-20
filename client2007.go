package dlt645

import (
	"fmt"
	"io"
	"time"
)

const (
	rtuMinSize          = 10
	rtuMaxSize          = 200 + rtuMinSize + 2
	rtuErrExceptionSize = 13

	ReadDataDomainMaxSize  = 200
	WriteDataDomainMaxSize = 50
)

const (
	CommunicationRate19200 = 0x40 // binary 0100 0000
	CommunicationRate9600  = 0x20 // binary 0010 0000
	CommunicationRate4800  = 0x10 // binary 0001 0000
	CommunicationRate2400  = 0x08 // binary 0000 1000
	CommunicationRate1200  = 0x04 // binary 0000 0100
	CommunicationRate600   = 0x02 // binary 0000 0010
)

type Client2007Handler struct {
	rtuPackager
	rtuSerialTransporter
}

func NewClient2007Handler(address string) *Client2007Handler {
	handler := &Client2007Handler{}
	handler.Address = address
	handler.Timeout = serialTimeout
	handler.IdleTimeout = serialIdleTimeout
	return handler
}

func Client2007(address string) Client {
	handler := NewClient2007Handler(address)
	return NewClient(handler)
}

type rtuPackager struct {
	SlaveAddr uint64
}

// Encode encodes a DTL645-2007 frame:
//
// StartSymbol   : 1 byte
// Address       : 6 byte
// StartSymbol2  : 1 byte
// ControlCode   : 1 byte
// DataLen       : 1 byte
// Data          : n byte
// CheckSUm      : 1 byte
// EndSymbol     : 1 byte
func (dtl *rtuPackager) Encode(frame *FramePayLoad) (raw []byte, err error) {
	dataDomainLen := len(frame.Data)

	rawLen := 12 + dataDomainLen

	if dataDomainLen > ReadDataDomainMaxSize {
		err = fmt.Errorf("dtl: length of data domain '%v' must not be bigger than '%v'", dataDomainLen, ReadDataDomainMaxSize)
		return
	}

	raw = make([]byte, rawLen)

	raw[0] = FrameHead
	slaveAddrBCD := BCDFromUint(dtl.SlaveAddr, 6)
	raw[1] = slaveAddrBCD[5]
	raw[2] = slaveAddrBCD[4]
	raw[3] = slaveAddrBCD[3]
	raw[4] = slaveAddrBCD[2]
	raw[5] = slaveAddrBCD[1]
	raw[6] = slaveAddrBCD[0]
	raw[7] = FrameHead
	// controlCode
	// 8 bit   : 0 master send   1 slave send
	// 7 biy   : 0 slave ok   1 slave err
	// 6 bit   : 0 have not follow-up data    1 have follow-up data
	// 1-5 bit : function code
	controlCode := frame.FunctionCode
	controlCode = controlCode & 0x1F // 0001 1111
	raw[8] = controlCode

	raw[9] = byte(dataDomainLen)
	dataDomain := make([]byte, len(frame.Data))
	for k, v := range frame.Data {
		dataDomain[k] = v + 0x33
	}
	copy(raw[10:], dataDomain)

	// append check sum
	checkSum := generateCheckSum(raw[:rawLen-2])
	raw[rawLen-2] = checkSum
	raw[rawLen-1] = FrameTail

	return
}

// Decode a DTL645-2007 frame:
//
// StartSymbol   : 1 byte
// Address       : 6 byte
// StartSymbol2  : 1 byte
// ControlCode   : 1 byte
// DataLen       : 1 byte
// Data          : n byte
// CheckSUm      : 1 byte
// EndSymbol     : 1 byte
func (dlt *rtuPackager) Decode(raw []byte) (payload *FramePayLoad, err error) {
	length := len(raw)
	// Calculate checksum
	checkSum := generateCheckSum(raw[:length-2])

	if checkSum != raw[length-2] {
		err = fmt.Errorf("dlt645: response check sum '%v' does not match expected '%v'", raw[length-2], checkSum)
		return
	}
	// Function code & data
	payload = &FramePayLoad{}
	payload.HasFollowUpData = (raw[8]&0x20)>>5 != 0 // 0010 0000
	payload.FunctionCode = raw[8] & 0x1F            // 0001 1111
	dataDomain := make([]byte, len(raw)-12)
	for k, v := range raw[10 : length-2] {
		dataDomain[k] = v - 0x33
	}
	payload.Data = dataDomain
	// check err word
	IsSlaveErr := (raw[8]&0x40)>>6 != 0 // 0100 0000
	if IsSlaveErr && len(dataDomain) > 0 {
		err = responseError(payload.FunctionCode, dataDomain[0])
		return
	}

	return
}

// verify verifies response length and slave id.
//
// StartSymbol   : 1 byte
// Address       : 6 byte
// StartSymbol2  : 1 byte
// ControlCode   : 1 byte
// DataLen       : 1 byte
// Data          : n byte
// CheckSUm      : 1 byte
// EndSymbol     : 1 byte
func (dlt *rtuPackager) Verify(request []byte, response []byte) (err error) {
	length := len(response)
	// Minimum size (including address, function and CRC)
	if length < rtuMinSize {
		err = fmt.Errorf("dlt: response length '%v' does not meet minimum '%v'", length, rtuMinSize)
		return
	}
	// Slave address must match
	for k, v := range response[1:7] {
		if v != request[k] {
			err = fmt.Errorf("dlt: response slave id '%v' does not match request '%v'", response[1:7], request[1:7])
			return
		}
	}

	return
}

// modify slave address domain for broadcast command.
// func (dtl *rtuPackager) ModifySlaveAddress(commAddr uint64) (err error) {
// 	if (commAddr >> 48) > 0 {
// 		err = fmt.Errorf("dlt645: communication address '%v' must be between '%v' and '%v',", commAddr, "0byte", "6byte")
// 		return
// 	}
// 	dtl.SlaveAddr = commAddr

// 	return
// }

type rtuSerialTransporter struct {
	serialPort
}

func (dlt *rtuSerialTransporter) Send(request []byte) (response []byte, err error) {
	// make sure port is connected
	if err = dlt.serialPort.connect(); err != nil {
		return
	}

	dlt.serialPort.lastActivity = time.Now()
	dlt.serialPort.startCloseTimer()

	raw := append([]byte{0xfe, 0xfe, 0xfe, 0xfe}, request...)

	// Send the request
	dlt.serialPort.logf("dlt: sending % x\n", raw)
	if _, err = dlt.port.Write(raw); err != nil {
		return
	}
	// controlCode
	// 8 bit   : 0 master send   1 slave send
	// 7 biy   : 0 slave ok   1 slave err
	// 6 bit   : 0 have not follow-up data    1 have follow-up data
	// 1-5 bit : function code
	// functionCode := request[8] & 0xE0 // 0001 1111

	time.Sleep(dlt.calculateDelay(len(raw) + rtuMaxSize))

	var n, n1 int
	var data [rtuMaxSize]byte
	// read frame head
	n, err = io.ReadAtLeast(dlt.port, data[:], rtuMinSize)
	if err != nil {
		return
	}
	dataDomainLen := int(data[9])
	bytesToRead := rtuMinSize + dataDomainLen + 2
	// read remaining data
	if n < bytesToRead {
		if bytesToRead <= rtuMaxSize {
			if bytesToRead > n {
				n1, err = io.ReadFull(dlt.port, data[n:bytesToRead])
				n += n1
			}
		}
	}
	if err != nil {
		return
	}

	response, err = dlt.ProcessPacket(data[:bytesToRead])
	if err != nil {
		return
	}

	dlt.serialPort.logf("dlt: received % x\n", response)
	return
}

func (dlt *rtuSerialTransporter) SendNotResponse(request []byte) (err error) {
	// make sure port is connected
	if err = dlt.serialPort.connect(); err != nil {
		return
	}

	dlt.serialPort.lastActivity = time.Now()
	dlt.serialPort.startCloseTimer()

	raw := append([]byte{0xfe, 0xfe, 0xfe, 0xfe}, request...)

	// Send the request
	dlt.serialPort.logf("dlt: sending % x\n", raw)
	if _, err = dlt.port.Write(raw); err != nil {
		return
	}
	// controlCode
	// 8 bit   : 0 master send   1 slave send
	// 7 biy   : 0 slave ok   1 slave err
	// 6 bit   : 0 have not follow-up data    1 have follow-up data
	// 1-5 bit : function code
	// functionCode := request[8] & 0xE0 // 0001 1111

	if err != nil {
		return
	}

	return
}

func (dlt *rtuSerialTransporter) ProcessPacket(data []byte) (result []byte, err error) {
	var frameStart int = -1
	var frameEnd int = -1

	for i, b := range data {
		if b == FrameHead && frameStart == -1 {
			frameStart = i
			continue
		}

		if b == FrameTail && frameStart != -1 {
			frameEnd = i
			break
		}
	}
	if frameStart == -1 || frameEnd == -1 {
		err = fmt.Errorf("dlt645: is not valid frame")
		return
	}
	result = data[frameStart:frameEnd]
	return
}

func (dlt *rtuSerialTransporter) calculateDelay(chars int) time.Duration {
	var characterDelay, frameDelay int // us

	if dlt.BaudRate <= 0 || dlt.BaudRate > 19200 {
		characterDelay = 750
		frameDelay = 1750
	} else {
		characterDelay = 15000000 / dlt.BaudRate
		frameDelay = 35000000 / dlt.BaudRate
	}
	return time.Duration(characterDelay*chars+frameDelay) * time.Microsecond
}
