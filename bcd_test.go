package dlt645

import (
	"testing"
)

func TestSean(t *testing.T) {
	//201512120120
	result := BCDFromUint(100, 2)

	// a := 0x49 - 0x33
	// t.Fatalf("0x%x", a)
	t.Fatalf("0x%x", result)
}
