package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/user"

	"github.com/creack/pty"
	"github.com/gorilla/websocket"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// WebSocket 升级器
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许跨域连接（可以根据需求调整）
	},
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// 升级 HTTP 连接为 WebSocket 连接
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()

	// 获取用户的 home 目录
	usr, err := user.Current()
	if err != nil {
		log.Fatalf("Failed to get current user: %v", err)
	}

	// 切换到用户的 home 目录
	err = os.Chdir(usr.HomeDir)
	if err != nil {
		log.Fatalf("Failed to change directory to home: %v", err)
	}
	log.Printf("Changed working directory to: %s", usr.HomeDir)

	shell := GetShell()
	// 创建一个 Bash Shell
	cmd := exec.Command(shell)

	// 设置环境变量，确保 $TERM 被正确传递
	cmd.Env = append(os.Environ(),
		"TERM=xterm-256color",
		"COLORTERM=truecolor",
		"FORCE_COLOR=3",
	)

	// 单独处理 PATH，确保它不会被覆盖
	newPath := "/opt/homebrew/bin:/opt/homebrew/sbin:/usr/local/bin:/usr/local/go/bin:" + getCurrentPath()
	cmd.Env = append(cmd.Env, "PATH="+newPath)

	// 创建伪终端
	ptmx, err := pty.Start(cmd)
	if err != nil {
		log.Println("PTY start error:", err)
		return
	}
	defer ptmx.Close()

	// 设置初始尺寸（可选）
	// if err := pty.Setsize(ptmx, &pty.Winsize{
	// 	Rows: 40,  // 初始行数
	// 	Cols: 160, // 初始列数
	// }); err != nil {
	// 	log.Println("Failed to set initial PTY size:", err)
	// }

	// 将伪终端的输出发送到 WebSocket 客户端
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := ptmx.Read(buf)
			if err != nil {
				if err == io.EOF {
					log.Println("PTY closed")
				} else {
					log.Println("PTY read error:", err)
				}
				return
			}
			if n > 0 {
				err := conn.WriteMessage(websocket.TextMessage, buf[:n])
				if err != nil {
					log.Println("WebSocket write error:", err)
					return
				}
			}
		}
	}()

	// 接收 WebSocket 客户端的输入并写入伪终端
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("WebSocket read error:", err)
			break
		}

		// 解析消息内容
		var msg struct {
			Type string `json:"type"`
			Cols int    `json:"cols"`
			Rows int    `json:"rows"`
		}
		if err := json.Unmarshal(message, &msg); err == nil && msg.Type == "resize" {
			// 更新伪终端尺寸
			if err := pty.Setsize(ptmx, &pty.Winsize{
				Rows: uint16(msg.Rows),
				Cols: uint16(msg.Cols),
			}); err != nil {
				log.Println("Failed to resize PTY:", err)
			}
			continue
		}

		_, err = ptmx.Write(message)
		if err != nil {
			log.Println("PTY write error:", err)
			break
		}
	}
}

func IsWSLInstalled() bool {
	// 检查 wsl.exe 是否存在于 PATH 中
	_, err := exec.LookPath("wsl.exe")
	return err == nil
}

// 根据操作系统来判断使用zsh还是bash还是powershell
func GetShell() string {
	// 根据操作系统判断
	if os.Getenv("SHELL") == "/bin/zsh" {
		return "zsh"
	} else {
		return "bash"
	}
}

// 获取环境变量
func getCurrentPath() string {
	// 尝试从不同来源获取 PATH
	if path := os.Getenv("PATH"); path != "" {
		return path
	}
	// 默认 PATH
	return "/usr/bin:/bin:/usr/sbin:/sbin"
}
