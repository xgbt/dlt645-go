dlt645-go
=========
使用go实现的多功能电能表通信协议

支持的指令
-------------------

* Read Data
* Write Data
* Read Communication Address
* Write Communication Address
* Broadcast Timing
* Freeze Command
* Change Communication Rate
* Change Password
* Clear Maximum Demand
* Clear Ammeter
* Clear Event


版本支持
-----------------
- [x] DL/T 645 2007 
- [ ] DL/T 645 1997

用法
-----
Basic usage:
```go
// default configuration is 19200, 8, 1, even
handler := dlt.NewClient2007Handler(rtuDevice)
err := handler.Connect()
defer handler.Close()
client := dlt.NewClient(handler)
results, err := client.ReadData(00000000, 0, 0, 0, 0, 0, 0)
```

Advanced usage:
```go
handler := dlt.NewClient2007Handler(rtuDevice)
handler.BaudRate = 4800
handler.DataBits = 8
handler.Parity = "N"
handler.StopBits = 1
handler.RS485.Enabled = true
handler.SlaveAddr = 304257140001
err := handler.Connect()
defer handler.Close()

client := dlt.NewClient(handler)
results, err := client.ReadData(00000000, 0, 0, 0, 0, 0, 0)
```

References
----------
* [DLT645-2007](https://www.toky.com.cn/up_pic/2020_12_15_12243_142130.pdf)
