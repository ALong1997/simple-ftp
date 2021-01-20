package main

import (
	"bufio"
	"log"
	"os"
	"path"
	"runtime"
	"strings"
	"sync"
)

type userInfo struct {
	name string
	pwd  string
	home string
}

var defaultDir = map[string]string{
	"windows": "D:/ALong1108/SimpleFTP/server",
	"unix":    "home/xxx",
}

var lock sync.Once // 初始化users一次
var users []userInfo

func init() {
	lock.Do(initUsers)
}

func initUsers() {
	cwd, ok := defaultDir[runtime.GOOS]
	if !ok {
		log.Fatal("Unsupported system.")
	}

	f, err := os.Open(cwd + "/users.txt")
	if err != nil {
		log.Fatal("failed to load users' information.", err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()
		userinfo := strings.Split(line, ";;")
		if len(userinfo) < 3 {
			continue
		}
		home := path.Join(cwd, userinfo[2])
		f, err := os.Open(home)
		if err != nil && os.IsNotExist(err) {
			err = os.Mkdir(home, os.ModePerm)
			if err != nil {
				log.Fatal("failed to make directory", home)
			}
		} else {
			f.Close()
		}
		users = append(users, userInfo{userinfo[0], userinfo[1], home})
	}
}

// 验证用户名和密码，返回验证结果true/false和验证通过后的用户家目录
func Validate(name string, pwd string) (pass bool, home string) {
	if len(users) <= 0 {
		return
	}

	for _, info := range users {
		if info.name == name && info.pwd == pwd {
			return true, info.home
		}
	}
	return
}
