package main

import (
	"SimpleFTP/common"
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

const maxBinarySize = 4096
const helpStr = "Help:\t[command] [args]\ncd [path]\n"

func main() {
	// 获取用户身份信息
	var user string
	flag.StringVar(&user, "user", "", "user")
	flag.Parse()
	if len(user) == 0 {
		log.Println(common.AuthenticationErr)
		return
	}

	fmt.Print("Password:")
	var pwd string
	input := bufio.NewScanner(os.Stdin)
	if input.Scan() {
		pwd = input.Text()
	}

	// 连接到ftp服务器
	ftpClient, err := NewDail(common.Address)
	if err != nil {
		fmt.Println(err)
		return
	}
	if ftpClient == nil {
		return
	}
	defer ftpClient.Close()

	// 用户认证
	err = ftpClient.Auth(user, pwd)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 监听命令行输入
	for input.Scan() && !ftpClient.Exit {
		text := input.Text()
		args := strings.Split(strings.TrimSpace(text), " ")
		if len(args) == 0 {
			printHelp()
			continue
		}

		command := args[0]
		if len(args) > 1 {
			args = args[1:]
		} else {
			args = nil
		}

		err = ftpClient.handleCommand(command, args)
		if err != nil {
			log.Println(err)
		}
	}

	return
}

func printHelp() {
	log.Printf(helpStr)
}
