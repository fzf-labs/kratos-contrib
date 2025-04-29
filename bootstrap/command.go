package bootstrap

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
)

type Command struct {
	Conf       string // 引导配置文件路径，默认为：../../configs
	Env        string // 开发环境：dev、debug……
	ConfigHost string // 远程配置服务端地址
	ConfigType string // 远程配置服务端类型
	Daemon     bool   // 是否转为守护进程
}

func NewCommand() *Command {
	f := new(Command)
	flag.StringVar(&f.Conf, "conf", "../../configs", "config path, eg: -conf ../../configs")
	flag.StringVar(&f.Env, "env", "dev", "runtime environment, eg: -env dev")
	flag.StringVar(&f.ConfigHost, "chost", "0.0.0.0:8500", "config server host, eg: -chost 0.0.0.0:8500")
	flag.StringVar(&f.ConfigType, "ctype", "file", "config server host, eg: -ctype consul")
	flag.BoolVar(&f.Daemon, "d", false, "run app as a daemon with -d=true.")
	if f.Daemon {
		BeDaemon("-d")
	}
	flag.Parse()
	return f
}

// BeDaemon 将当前进程转为守护进程
func BeDaemon(arg string) {
	subProcess(stripSlice(os.Args, arg))
	fmt.Printf("[*] Daemon running in PID: %d PPID: %d\n", os.Getpid(), os.Getppid())
	os.Exit(0)
}
func stripSlice(slice []string, element string) []string {
	for i := 0; i < len(slice); {
		if slice[i] == element && i != len(slice)-1 {
			slice = append(slice[:i], slice[i+1:]...)
		} else if slice[i] == element && i == len(slice)-1 {
			slice = slice[:i]
		} else {
			i++
		}
	}
	return slice
}

//nolint:gosec
func subProcess(args []string) *exec.Cmd {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "[-] Error: %s\n", err)
	}
	return cmd
}
