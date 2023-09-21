package utils

func GenerateCheckSum(msg []byte) uint8 {
	cs := uint8(0)

	for _, v := range msg {
		cs += v
	}
	return cs
}
