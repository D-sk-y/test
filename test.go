package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	ch       = make(chan string) // 无缓冲
	mu       sync.Mutex
	receiver bool // 接收者是否在等待
)

func main() {
	fmt.Println("命令:")
	fmt.Println("  start  - 启动接收者（开始监听）")
	fmt.Println("  stop   - 停止接收者")
	fmt.Println("  send   - 阻塞发送")
	fmt.Println("  try    - 非阻塞发送（select+default）")
	fmt.Println("  wait   - 超时发送（ctx超时150ms）")
	fmt.Println("  quit   - 退出")
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		cmd := strings.TrimSpace(scanner.Text())

		switch cmd {
		case "start":
			mu.Lock()
			if receiver {
				fmt.Println("  接收者已在运行")
			} else {
				receiver = true
				go runReceiver()
				fmt.Println("  接收者已启动，等待消息中...")
			}
			mu.Unlock()

		case "stop":
			mu.Lock()
			receiver = false
			mu.Unlock()
			fmt.Println("  接收者已停止")

		case "send":
			fmt.Println("  阻塞发送 'hello'...")
			ch <- "hello"
			fmt.Println("  → 发送成功")

		case "try":
			select {
			case ch <- "hello":
				fmt.Println("  → 发送成功")
			default:
				fmt.Println("  → 丢弃：接收方未就绪")
			}

		case "wait":
			ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
			defer cancel()
			fmt.Println("  超时发送（150ms）...")
			select {
			case ch <- "hello":
				fmt.Println("  → 发送成功")
			case <-ctx.Done():
				fmt.Println("  → 超时放弃")
			}

		case "quit":
			return

		default:
			fmt.Println("  未知命令")
		}
	}
}

func runReceiver() {
	for {
		mu.Lock()
		if !receiver {
			mu.Unlock()
			return
		}
		mu.Unlock()

		select {
		case msg := <-ch:
			fmt.Printf("  ★ 接收者收到: %s\n", msg)
		case <-time.After(500 * time.Millisecond):
			// 每 500ms 检查一下是否被 stop
		}
	}
}
