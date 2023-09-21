package utils

func pow100(power byte) uint64 {
	res := uint64(1)
	for i := byte(0); i < power; i++ {
		res *= 100
	}
	return res
}

func BCDFromUint(value uint64, size int) []byte {
	buf := make([]byte, size)
	if value > 0 {
		remainder := value
		for pos := size - 1; pos >= 0 && remainder > 0; pos-- {
			tail := byte(remainder % 100)
			hi, lo := tail/10, tail%10
			buf[pos] = byte(hi<<4 + lo)
			remainder = remainder / 100
		}
	}
	return buf
}

func BCDFromUint8(value uint8) byte {
	return BCDFromUint(uint64(value), 1)[0]
}

func BCDFromUint16(value uint16) []byte {
	return BCDFromUint(uint64(value), 2)
}

func BCDFromUint32(value uint32) []byte {
	return BCDFromUint(uint64(value), 4)
}

func BCDFromUint64(value uint64) []byte {
	return BCDFromUint(value, 8)
}

func BCDToUint8(value byte) uint8 {
	return uint8(toUint([]byte{value}, 1))
}

func BCDToUint16(value []byte) uint16 {
	return uint16(toUint(value, 2))
}

func BCDToUint32(value []byte) uint32 {
	return uint32(toUint(value, 4))
}

func BCDToUint64(value []byte) uint64 {
	return toUint(value, 8)
}

func toUint(bytes []byte, size int) uint64 {
	len := len(bytes)
	if len > size {
		bytes = bytes[len-size:]
	}
	res := uint64(0)
	for i, b := range bytes {
		hi, lo := b>>4, b&0x0f
		if hi > 9 || lo > 9 {
			return 0
		}
		res += uint64(hi*10+lo) * pow100(byte(len-i)-1)
	}
	return res
}
