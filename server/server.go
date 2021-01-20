package main

import (
	"SimpleFTP/common"
	"encoding/binary"
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", common.Address)
	if err != nil {
		log.Fatal(err)
	}

	for {
		con, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleFunc(con)
	}
}

func handleFunc(con net.Conn) {
	ftpServer := NewServer(con, "", "")
	if ftpServer == nil {
		return
	}
	defer ftpServer.Close()

	// 身份验证
	// 读取用户名
	var length uint32
	err := ftpServer.Read(&length)
	if err != nil {
		err = ftpServer.Write(uint32(0))
		if err != nil {
			log.Println(err)
		}
		return
	}
	user := make([]byte, length-uint32(binary.Size(length)))
	err = ftpServer.Read(user)
	if err != nil {
		err = ftpServer.Write(uint32(0))
		if err != nil {
			log.Println(err)
		}
		return
	}

	// 读取密码
	err = ftpServer.Read(&length)
	if err != nil {
		err = ftpServer.Write(uint32(0))
		if err != nil {
			log.Println(err)
		}
		return
	}
	pwd := make([]byte, length-uint32(binary.Size(length)))
	err = ftpServer.Read(pwd)
	if err != nil {
		err = ftpServer.Write(uint32(0))
		if err != nil {
			log.Println(err)
		}
		return
	}

	// 验证用户名密码获取家目录
	validated, cwd := Validate(common.Sbyte2str(user), common.Sbyte2str(pwd))
	if !validated {
		err = ftpServer.Write(uint32(0))
		if err != nil {
			log.Println(err)
		}
		return
	}

	home := common.Str2sbyte(cwd)
	err = ftpServer.Write(uint32(binary.Size(home)))
	if err != nil {
		log.Println(err)
		return
	}
	err = ftpServer.Write(home)
	if err != nil {
		log.Println(err)
		return
	}

	ftpServer.Cwd = cwd
	ftpServer.Home = cwd

	// 循环监听命令请求
	for !ftpServer.Exit {
		var length uint32
		err = ftpServer.Read(&length)
		if err != nil {
			log.Println(err)
			return
		}
		var cmdId uint8
		err = ftpServer.Read(&cmdId)
		if err != nil {
			log.Println(err)
			return
		}
		args := make([]byte, length-uint32(binary.Size(cmdId))-uint32(binary.Size(length)))
		err = ftpServer.Read(args)
		if err != nil {
			log.Println(err)
			return
		}

		switch cmdId {
		case common.CdId:
			err = ftpServer.HandleCd(args)
		case common.LsId:
			err = ftpServer.HandleLs(args)
		case common.ExitId:
			err = ftpServer.HandleExit(args)
		case common.MkdirId:
			err = ftpServer.HandleMkdir(args)
		case common.PutId:
			err = ftpServer.HandlePut(args)
		case common.GetId:
			err = ftpServer.HandleGet(args)
		default:
			err = ftpServer.WriteContent([]byte(common.InvalidCommandErr.Error()))
		}

		if err != nil {
			log.Println(err)
		}
	}
}
