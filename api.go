package dlt645

type Client interface {
	// read data
	ReadData(dataMarker uint32, blockQuantity uint8, year, month, day, hour, minute uint8) (results []byte, err error)
	// write data
	WriteData(dataMarker uint32, passwordPermission uint8, password uint32, operatorCode uint32, data []byte) (results []byte, err error)
	// read communication address
	ReadCommunicationAddress() (results []byte, err error)
	// write communication address
	WriteCommunicationAddress(commAddr uint64) (results []byte, err error)
	// broadcast timing
	BroadcastTiming(year, month, day, hour, minute, second uint8) (err error)
	// freeze command
	FreezeCommand(month, day, hour, minute uint8) (results []byte, err error)
	// change communication speed
	ChangeCommunicationRate(Word uint8) (results []byte, err error)
	// change password
	ChangePassword(dataMarker uint32, oldPasswordPermission uint8, oldPassword uint32, newPasswordPermission uint8, newPassword uint32) (results []byte, err error)
	// Clear the maximum demand
	ClearMaximumDemand(passwordPermission uint8, password uint32, operatorCode uint32) (results []byte, err error)
	// Clear the ammeter
	ClearAmmeter(passwordPermission uint8, password uint32, operatorCode uint32) (results []byte, err error)
	// Clear the event
	ClearEvent(dataMarker uint32, passwordPermission uint8, password uint32, operatorCode uint32) (results []byte, err error)
}
