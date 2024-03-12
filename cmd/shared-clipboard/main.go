package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"syscall"

	"github.com/trueaniki/admiral"
	"github.com/trueaniki/gopeers"
	sharedclipboard "github.com/trueaniki/shared-clipboard"
	"golang.design/x/clipboard"
)

const appName = "shared-clipboard"
const appDesc = "Shared Clipboard"
const version = "dev"

type Conf struct {
	Version bool   `type:"flag" name:"version" alias:"v" description:"Print version and exit"`
	Network string `type:"flag" name:"network" alias:"n" description:"Network to scan in CIDR format" required:"true"`
	Conf    string `type:"flag" name:"conf" alias:"c" description:"Path to config file"`

	Start Start `type:"command" name:"start" description:"Start the shared clipboard daemon"`
	Stop  Stop  `type:"command" name:"stop" description:"Stop the shared clipboard daemon"`
}

type Start struct {
	Network string `type:"flag" name:"network" alias:"n" description:"Network to scan in CIDR format" required:"true"`
	Conf    string `type:"flag" name:"conf" alias:"c" description:"Path to config file"`
	Logfile string `type:"flag" name:"logfile" alias:"l" description:"Path to log file"`
}

type Stop struct {
}

const daemonPort = 17893

func main() {
	conf := &Conf{}
	a := admiral.New(appName, appDesc)
	a.Configure(conf)
	a.Flag("version").Handle(func(_ any) {
		fmt.Println(version)
		os.Exit(0)
	})

	a.Command("start").Handle(func(args interface{}) {
		args = args.(*Start)

		cmd := exec.Command(os.Args[0], "-n", args.(*Start).Network)
		cmd.Env = append(os.Environ(), "DAEMON=true")
		cmd.SysProcAttr = &syscall.SysProcAttr{
			Setsid: true,
		}
		logfile, err := os.OpenFile(args.(*Start).Logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			printAndExit(err)
		}
		defer logfile.Close()

		cmd.Stdout = logfile
		cmd.Stderr = logfile

		fmt.Println("Starting daemon")
		err = cmd.Start()
		if err != nil {
			printAndExit(err)
		}
		os.Exit(0)
	})

	a.Command("stop").Handle(func(args interface{}) {
		conn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", daemonPort))
		if err != nil {
			printAndExit(err)
		}
		conn.Close()
		os.Exit(0)
	})

	_, err := a.Parse(os.Args)
	if err != nil {
		printAndExit(err)
	}

	if os.Getenv("DAEMON") == "true" {
		go func() {
			// Listen for stop signal
			l, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", daemonPort))
			if err != nil {
				printAndExit(err)
			}
			defer l.Close()
			for {
				conn, err := l.Accept()
				if err != nil {
					printAndExit(err)
				}
				conn.Close()
				fmt.Println("Stopping due to stop signal")
				os.Exit(0)
			}
		}()
	}

	fmt.Println("Starting shared clipboard daemon")
	start(conf.Network)
}

func printAndExit(err error) {
	fmt.Println(err)
	os.Exit(1)
}

func start(network string) {
	err := clipboard.Init()
	if err != nil {
		printAndExit(err)
	}
	os.Stdout.WriteString("Clipboard initialized\n")
	locals := gopeers.PingSweep(network)
	peer := gopeers.NewPeer(locals)
	peer.Start()

	sharedclipboard.Listen(peer)
}
