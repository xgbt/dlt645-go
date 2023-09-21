package main

import (
	"log"
	"os"

	dlt "github.com/xgbt/dlt645-go"
)

const (
	rtuDevice = "/dev/ttyS9"
	Address   = 304257140001
)

func main() {
	handler := dlt.NewClient2007Handler(rtuDevice)
	handler.BaudRate = 4800
	handler.DataBits = 8
	handler.Parity = "N"
	handler.StopBits = 1
	handler.RS485.Enabled = false
	handler.SlaveAddr = Address // [6]uint8{0x99, 0x99, 0x99, 0x99, 0x99, 0x99}
	handler.Logger = log.New(os.Stdout, "dlt645: ", log.LstdFlags)
	err := handler.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer handler.Close()

	client := dlt.NewClient(handler)
	results, err := client.ReadData(0x00000000, 0, 0, 0, 0, 0, 0)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%#v\n", results)
}
