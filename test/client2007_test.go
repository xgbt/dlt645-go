package test

import (
	"log"
	"os"
	"testing"

	dlt "github.com/xgbt/dlt645-go"
)

const (
	rtuDevice = "/dev/ttyUSB0"
)

func TestClient2007AdvancedUsage(t *testing.T) {
	handler := dlt.NewClient2007Handler(rtuDevice)
	handler.BaudRate = 4800
	handler.DataBits = 8
	handler.Parity = "N"
	handler.StopBits = 1
	handler.RS485.Enabled = true
	handler.SlaveAddr = 304257140001 // [6]uint8{0x99, 0x99, 0x99, 0x99, 0x99, 0x99}
	handler.Logger = log.New(os.Stdout, "dlt645: ", log.LstdFlags)
	err := handler.Connect()
	if err != nil {
		t.Fatal(err)
	}
	defer handler.Close()

	client := dlt.NewClient(handler)
	results, err := client.ReadData(00000000, 0, 0, 0, 0, 0, 0)
	if err != nil || results == nil {
		t.Fatal(err, results)
	}
}
