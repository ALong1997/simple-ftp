package common

import (
	"encoding/binary"
	"net"
)

var littleEndian = binary.LittleEndian

type FtpConn struct {
	Con  net.Conn
	Cwd  string
	Home string
	Exit bool
}

func (ftpCon *FtpConn) Read(data interface{}) error {
	return binary.Read(ftpCon.Con, littleEndian, data)
}

func (ftpCon *FtpConn) Write(data interface{}) error {
	if dataStr, ok := data.(string); ok {
		data = Str2sbyte(dataStr)
	}
	return binary.Write(ftpCon.Con, littleEndian, data)
}

func (ftpCon *FtpConn) WriteCommand(cmdId uint8, args string) error {
	var length uint32
	length = uint32(binary.Size(length)+binary.Size(cmdId)) + uint32(len(args))

	err := ftpCon.Write(length)
	if err != nil {
		return err
	}
	err = ftpCon.Write(cmdId)
	if err != nil {
		return err
	}
	err = ftpCon.Write(args)
	if err != nil {
		return err
	}

	return nil
}

func (ftpCon *FtpConn) WriteContent(content []byte) error {
	var length = uint32(len(content))

	// length == 0
	if length == 0 {
		return ftpCon.Write(length)
	}

	// write length
	length = length + uint32(binary.Size(length))
	err := ftpCon.Write(length)
	if err != nil {
		return err
	}

	// write content
	err = ftpCon.Write(content)
	if err != nil {
		return err
	}

	return nil
}
