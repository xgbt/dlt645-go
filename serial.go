package dlt645

import (
	"io"
	"log"
	"sync"
	"time"

	"github.com/goburrow/serial"
)

const (
	serialTimeout     = 5 * time.Second
	serialIdleTimeout = 60 * time.Second
)

// serialPort has configuration and I/O controller.
type serialPort struct {
	serial.Config

	Logger      *log.Logger
	IdleTimeout time.Duration

	mu           sync.Mutex
	port         io.ReadWriteCloser
	lastActivity time.Time
	closeTimer   *time.Timer
}

func (dlt *serialPort) Connect() (err error) {
	dlt.mu.Lock()
	defer dlt.mu.Unlock()

	return dlt.connect()
}

func (dlt *serialPort) connect() error {
	if dlt.port == nil {
		port, err := serial.Open(&dlt.Config)
		if err != nil {
			return err
		}
		dlt.port = port
	}
	return nil
}

func (dlt *serialPort) Close() (err error) {
	dlt.mu.Lock()
	defer dlt.mu.Unlock()

	return dlt.close()
}

func (dlt *serialPort) close() (err error) {
	if dlt.port != nil {
		err = dlt.port.Close()
		dlt.port = nil
	}
	return
}

func (dlt *serialPort) logf(format string, v ...interface{}) {
	if dlt.Logger != nil {
		dlt.Logger.Printf(format, v...)
	}
}

func (dlt *serialPort) startCloseTimer() {
	if dlt.IdleTimeout <= 0 {
		return
	}
	if dlt.closeTimer == nil {
		dlt.closeTimer = time.AfterFunc(dlt.IdleTimeout, dlt.closeIdle)
	} else {
		dlt.closeTimer.Reset(dlt.IdleTimeout)
	}
}

func (dlt *serialPort) closeIdle() {
	dlt.mu.Lock()
	defer dlt.mu.Unlock()

	if dlt.IdleTimeout <= 0 {
		return
	}
	idle := time.Since(dlt.lastActivity) // time.Now().Sub(dlt.lastActivity)
	if idle >= dlt.IdleTimeout {
		dlt.logf("dlt645: closing connection due to idle timeout: %v", idle)
		dlt.close()
	}
}
