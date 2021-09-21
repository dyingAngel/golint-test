package tcpclient

import (
	"fmt"
	"net"
	"time"

	"github.com/dyingAngel/golint-test/utils/net/tcpclient/mock"

	"github.com/spf13/viper"
)

// TCPDoer ...
type TCPDoer interface {
	TCPDo() (string, error)
}

// TCPParam ...
type TCPParam struct {
	// tcp address
	Address string `json:"address"`
	// tcp content type
	ContentType string `json:"content_type"`
	// tcp Body Message
	Body string `json:"body"`
	// tcp timeout, in miliseconds
	Timeout int `json:"timeout"`
	// tcp max read package size
	MaxReadPackageSize int `json:"max_read_package_size"`
}

// TCPDo ...
func (tcpParam *TCPParam) TCPDo() (string, error) {
	if viper.Get("env") == "testing" {
		return mock.RespTCPISO, nil
	}
	timeout := time.Duration(tcpParam.Timeout) * time.Second

	conn, err := net.DialTimeout("tcp", tcpParam.Address, timeout)
	if err != nil {
		return "", fmt.Errorf("Dial: " + err.Error())
	}

	// Close connection when this function ends
	defer func() {
		conn.Close()
	}()

	_, errWrite := conn.Write([]byte(tcpParam.Body))
	if errWrite != nil {
		return "", fmt.Errorf("Write: " + errWrite.Error())
	}

	errSetRead := conn.SetReadDeadline(time.Now().Add(timeout))
	if errSetRead != nil {
		return "", fmt.Errorf("SetReadDeadline: " + errSetRead.Error())
	}

	buff := make([]byte, tcpParam.MaxReadPackageSize)
	dataLength, errRead := conn.Read(buff)
	message := string(buff[:dataLength])

	if errRead != nil {
		return "", fmt.Errorf("Read: " + errRead.Error())
	}

	return message, nil
}
