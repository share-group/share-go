package util

import (
	"fmt"
	"log"
	"net"
	"os"
	"runtime/debug"
	"strings"
)

type systemUtil struct{}

var SystemUtil = newSystemUtil()

func newSystemUtil() *systemUtil {
	return &systemUtil{}
}

/**
 * 模拟三元表达式，简化代码
 * @param condition 条件
 * @param trueVal 条件为真时的值
 * @param falseVal 条件为假时的值
 */
func (s *systemUtil) If(condition bool, trueVal any, falseVal any) any {
	if condition {
		return trueVal
	}
	return falseVal
}

/**
 * 封装一下协程，保证协程在出现异常的情况下程序不崩溃
 * @param businessFunc 处理业务的方法 (不可为空)
 * @param catchFunc 异常处理的方法 (可为空)
 * @param finallyFunc 无论异常与否都会最终执行的方法 (可为空)
 */
func (s *systemUtil) Goroutine(funs ...func()) {
	if len(funs) <= 0 || len(funs) > 3 {
		panic("参数错误")
		return
	}

	// 未传方法进来的，给一个空方法
	for i := len(funs); i < 3; i++ {
		funs = append(funs, func() {})
	}

	businessFunc := funs[0]
	catchFunc := funs[1]
	finallyFunc := funs[2]
	go func(f func()) {
		defer func() {
			if e := recover(); e != nil {
				message := fmt.Sprintf("协程捕获异常: %v \n 协程堆栈信息：\n %v", e, string(debug.Stack()))
				log.Println(message)
				catchFunc()
				finallyFunc()
			}
		}()
		businessFunc()
	}(businessFunc)
	finallyFunc()
}

// 获取本机IP
func (s *systemUtil) GetLocalIp() string {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		log.Println("net.Interfaces failed, err:", err.Error())
		return ""
	}
	for i := 0; i < len(netInterfaces); i++ {
		if (netInterfaces[i].Flags & net.FlagUp) != 0 {
			addrs, _ := netInterfaces[i].Addrs()
			for _, address := range addrs {
				if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						return ipnet.IP.String()
					}
				}
			}
		}
	}

	return ""
}

// 获取主机名
func (s *systemUtil) GetHostName() string {
	hostname, _ := os.Hostname()
	return strings.TrimSpace(hostname)
}
