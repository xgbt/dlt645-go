package dlt645

import (
	"encoding/binary"
	"fmt"
)

type ClientHandler interface {
	Packager
	Transporter
}

type client struct {
	packager    Packager
	transporter Transporter
}

func NewClient(handler ClientHandler) Client {
	return &client{packager: handler, transporter: handler}
}

// ReadData
func (dtl *client) ReadData(dataMarker uint32, blockQuantity uint8, year, month, day, hour, minute uint8) (results []byte, err error) {
	uintArray := []interface{}{dataMarker}
	if blockQuantity > 0 {
		uintArray = append(uintArray, blockQuantity)
	}
	if blockQuantity > 0 && year > 0 {
		uintArray = append(uintArray, blockQuantity, year, month, day, hour, minute)
	}

	request := FramePayLoad{
		FunctionCode: byte(FunctionCodeReadData),
		Data:         uintArrayToDataDomain(uintArray...),
	}
	response, err := dtl.send(&request)
	if err != nil {
		return
	}
	results = append(results, response.Data...)

	for seq := uint8(1); response.HasFollowUpData && seq <= 255; seq++ {
		request := FramePayLoad{
			FunctionCode: byte(FunctionCodeReadFollowUpData),
			Data:         uintArrayToDataDomain(dataMarker, seq),
		}
		response, err = dtl.send(&request)
		if err != nil {
			return
		}
		results = append(results, response.Data...)
	}

	results = results[4:]
	return
}

// WriteData
func (dtl *client) WriteData(dataMarker uint32, passwordPermission uint8, password uint32, operatorCode uint32, data []byte) (results []byte, err error) {
	if passwordPermission > 9 {
		err = fmt.Errorf("dlt645: password permission '%v must be between '%v' and '%v',", passwordPermission, "0", "9")
		return
	}
	if (password >> 24) > 0 {
		err = fmt.Errorf("dlt645: password '%v' must be less than '%v'", password, "3byte")
		return
	}

	request := FramePayLoad{
		FunctionCode: byte(FunctionCodeWriteData),
		Data:         uintArrayToDataDomain(dataMarker, uint32(passwordPermission)<<24|password, operatorCode, data),
	}
	response, err := dtl.send(&request)
	if err != nil {
		return
	}

	results = response.Data
	return
}

// ReadCommunicationAddress
func (dtl *client) ReadCommunicationAddress() (results []byte, err error) {
	request := FramePayLoad{
		FunctionCode: byte(FunctionCodeReadCommunicationAddress),
	}
	response, err := dtl.send(&request)
	if err != nil {
		return
	}
	results = response.Data

	return
}

// WriteCommunicationAddress
func (dtl *client) WriteCommunicationAddress(commAddr uint64) (results []byte, err error) {
	if (commAddr >> 48) > 0 {
		err = fmt.Errorf("dlt645: communication address '%v' must be between '%v' and '%v',", commAddr, "0byte", "6byte")
		return
	}

	request := FramePayLoad{
		FunctionCode: byte(FunctionCodeWriteCommunicationAddress),
		Data:         uintArrayToDataDomain(commAddr),
	}
	response, err := dtl.send(&request)
	if err != nil {
		return
	}
	results = response.Data
	return
}

// BroadcastTiming
func (dtl *client) BroadcastTiming(year, month, day, hour, minute, second uint8) (err error) {
	request := FramePayLoad{
		FunctionCode: byte(FunctionCodeBroadcastTiming),
		Data:         uintArrayToDataDomain(year, month, day, hour, minute, second),
	}

	_, err = dtl.send(&request, true)
	if err != nil {
		return
	}

	return
}

// FreezeCommand
func (dtl *client) FreezeCommand(month, day, hour, minute uint8) (results []byte, err error) {
	request := FramePayLoad{
		FunctionCode: byte(FunctionCodeFreezeCommand),
		Data:         uintArrayToDataDomain(month, day, hour, minute),
	}

	response, err := dtl.send(&request)
	if err != nil {
		return
	}
	results = response.Data

	return
}

// change communication speed  Word 特征字
func (dtl *client) ChangeCommunicationRate(word uint8) (results []byte, err error) {
	request := FramePayLoad{
		FunctionCode: byte(FunctionCodeChangeCommunicationRate),
		Data:         uintArrayToDataDomain(word),
	}
	response, err := dtl.send(&request)
	if err != nil {
		return
	}
	results = response.Data

	return
}

// change password
func (dtl *client) ChangePassword(dataMarker uint32, oldPasswordPermission uint8, oldPassword uint32, newPasswordPermission uint8, newPassword uint32) (results []byte, err error) {
	if oldPasswordPermission > 9 || newPasswordPermission > 9 {
		err = fmt.Errorf("dlt645: old or new password permission '%v'&'%v' must be between '%v' and '%v',", oldPasswordPermission, newPasswordPermission, "0", "9")
		return
	}
	if (oldPassword>>24) > 0 || (newPassword>>24) > 0 {
		err = fmt.Errorf("dlt645: old or new password '%v'&'%v' must be less than '%v'", oldPassword, newPassword, "3byte")
		return
	}

	request := FramePayLoad{
		FunctionCode: byte(FunctionCodeChangePassword),
		Data:         uintArrayToDataDomain(dataMarker, uint32(oldPasswordPermission)<<24|oldPassword, uint32(newPasswordPermission)<<24|newPassword),
	}
	response, err := dtl.send(&request)
	if err != nil {
		return
	}
	results = response.Data

	return
}

// Clear the maximum demand
func (dtl *client) ClearMaximumDemand(passwordPermission uint8, password uint32, operatorCode uint32) (results []byte, err error) {
	if passwordPermission > 9 {
		err = fmt.Errorf("dlt645: password permission '%v' must be between '%v' and '%v',", passwordPermission, "0", "9")
		return
	}
	if (password >> 24) > 0 {
		err = fmt.Errorf("dlt645: password '%v' must be less than '%v'", password, "3byte")
		return
	}

	request := FramePayLoad{
		FunctionCode: byte(FunctionCodeClearMaximumDemand),
		Data:         uintArrayToDataDomain(uint32(passwordPermission)<<24|password, operatorCode),
	}
	response, err := dtl.send(&request)
	if err != nil {
		return
	}
	results = response.Data

	return
}

// Clear the ammeter
func (dtl *client) ClearAmmeter(passwordPermission uint8, password uint32, operatorCode uint32) (results []byte, err error) {
	if passwordPermission > 9 {
		err = fmt.Errorf("dlt645: password permission '%v' must be between '%v' and '%v',", passwordPermission, "0", "9")
		return
	}
	if (password >> 24) > 0 {
		err = fmt.Errorf("dlt645: password '%v' must be less than '%v'", password, "3byte")
		return
	}

	request := FramePayLoad{
		FunctionCode: byte(FunctionCodeClearAmmeter),
		Data:         uintArrayToDataDomain(uint32(passwordPermission)<<24|password, operatorCode),
	}
	response, err := dtl.send(&request)
	if err != nil {
		return
	}
	results = response.Data

	return
}

// Clear the event
func (dtl *client) ClearEvent(dataMarker uint32, passwordPermission uint8, password uint32, operatorCode uint32) (results []byte, err error) {
	if passwordPermission > 9 {
		err = fmt.Errorf("dlt645: password permission '%v' must be between '%v' and '%v',", passwordPermission, "0", "9")
		return
	}
	if (password >> 24) > 0 {
		err = fmt.Errorf("dlt645: password '%v' must be less than '%v'", password, "3byte")
		return
	}

	request := FramePayLoad{
		FunctionCode: byte(FunctionCodeClearAmmeter),
		Data:         uintArrayToDataDomain(uint32(passwordPermission)<<24|password, operatorCode, dataMarker),
	}
	response, err := dtl.send(&request)
	if err != nil {
		return
	}
	results = response.Data

	return
}

// (dtl *client) send
//
// conditions `true` no response required
func (dtl *client) send(request *FramePayLoad, conditions ...interface{}) (response *FramePayLoad, err error) {
	rawRequest, err := dtl.packager.Encode(request)
	if err != nil {
		return
	}

	if len(conditions) > 0 {
		if nrr, _ := conditions[0].(bool); nrr {
			err = dtl.transporter.SendNotResponse(rawRequest)
			return
		}
	}

	dltResponse, err := dtl.transporter.Send(rawRequest)
	if err != nil {
		return
	}
	if err = dtl.packager.Verify(rawRequest, dltResponse); err != nil {
		return
	}
	response, err = dtl.packager.Decode(dltResponse)
	if err != nil {
		return
	}

	return
}

// func (dtl *client) sendNotResponse(request *FramePayLoad) (err error) {
// 	rawRequest, err := dtl.packager.Encode(request)
// 	if err != nil {
// 		return
// 	}

// 	err = dtl.transporter.SendNotResponse(rawRequest)
// 	if err != nil {
// 		return
// 	}

// 	return
// }

// func uintArrayToDataDomain(n ...interface{}) []byte {
// 	return uintArrayToBytes(n...)
// }

func uintArrayToDataDomain(n ...interface{}) []byte {
	var b []byte
	for _, v := range n {
		b = append(b, uintToBytes(v)...)
	}
	return b
}

func uintToBytes(n interface{}) (data []byte) {
	switch n := n.(type) {
	case uint8:
		data = append(data, n)
	case uint16:
		data = make([]byte, 2)
		binary.BigEndian.PutUint16(data, n)
	case uint32:
		data = make([]byte, 4)
		binary.BigEndian.PutUint32(data, n)
	case uint64:
		data = make([]byte, 8)
		binary.BigEndian.PutUint64(data, n)
	}
	return
}

func dataBlock(value ...uint16) []byte {
	data := make([]byte, 2*len(value))
	for i, v := range value {
		binary.BigEndian.PutUint16(data[i*2:], v)
	}
	return data
}

func responseError(functionCode, errWord byte) error {
	dltError := &DltError{FunctionCode: functionCode}
	dltError.ExceptionCode = errWord

	return dltError
}
