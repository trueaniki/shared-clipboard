package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"path"
	"syscall"

	"github.com/trueaniki/admiral"
	parsehotkeys "github.com/trueaniki/go-parse-hotkeys"
	"github.com/trueaniki/gopeers"
	sharedclipboard "github.com/trueaniki/shared-clipboard"
	"go.uber.org/zap"
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
	Init  Init  `type:"command" name:"init" description:"Initialize the config file at ~/.shared-clipboard.conf"`
}

type Init struct{}

type Start struct {
	Network string `type:"flag" name:"network" alias:"n" description:"Network to scan in CIDR format" required:"true"`
	Conf    string `type:"flag" name:"conf" alias:"c" description:"Path to config file"`
	Logfile string `type:"flag" name:"logfile" alias:"l" description:"Path to log file"`
}

type Stop struct{}

const daemonPort = 17893

var defaultShareHk = "Ctrl+Shift+S"
var defaultAdoptHk = "Ctrl+Shift+A"

var hotkeys *sharedclipboard.Hotkeys

func main() {
	shareHk, _ := parsehotkeys.Parse(defaultShareHk, "+")
	adoptHk, _ := parsehotkeys.Parse(defaultAdoptHk, "+")
	hotkeys = &sharedclipboard.Hotkeys{
		HKShare: shareHk,
		HKAdopt: adoptHk,
	}

	conf := &Conf{}
	a := admiral.New(appName, appDesc)
	a.Configure(conf)

	// Handle version flag
	a.Flag("version").Handle(func(_ any) {
		fmt.Println(version)
		os.Exit(0)
	})

	// Handle init command
	a.Command("init").Handle(func(_ interface{}) {
		confpath := path.Join(getHomeDir(), ".shared-clipboard.conf")
		// Check if file exists
		if _, err := os.Stat(confpath); !os.IsNotExist(err) {
			fmt.Println("Config file already exists at", confpath)
			os.Exit(0)
		}
		f, err := os.Create(confpath)
		if err != nil {
			printAndExit(err)
		}
		defer f.Close()
		_, err = f.WriteString(fmt.Sprintf("# Share=%s\n# Adopt=%s\n", defaultShareHk, defaultAdoptHk))
		if err != nil {
			printAndExit(err)
		}
		fmt.Println("Config file created at", confpath)
		os.Exit(0)
	})

	// Handle start command
	a.Command("start").Handle(func(opts interface{}) {
		args := opts.(*Start)

		cmdArgs := []string{"-n", args.Network}
		if args.Conf != "" {
			cmdArgs = append(cmdArgs, "-c", args.Conf)
		}
		cmd := exec.Command(os.Args[0], cmdArgs...)
		cmd.Env = append(os.Environ(), "DAEMON=true")
		cmd.SysProcAttr = &syscall.SysProcAttr{
			Setsid: true,
		}
		logfile, err := os.OpenFile(args.Logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
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

	// Handle stop command
	a.Command("stop").Handle(func(args interface{}) {
		conn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", daemonPort))
		if err != nil {
			printAndExit(err)
		}
		conn.Close()
		os.Exit(0)
	})

	// Parse the command line arguments
	_, err := a.Parse(os.Args)
	if err != nil {
		printAndExit(err)
	}

	// If DAEMON env is set, start the daemon
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

	// The actual application starts here
	fmt.Println("Starting shared clipboard")

	// Load config file
	confpath := path.Join(getHomeDir(), ".shared-clipboard.conf")
	// Check if file exists
	if _, err := os.Stat(confpath); os.IsNotExist(err) {
		confpath = ""
	}
	if conf.Conf != "" {
		confpath = conf.Conf
	}

	if confpath != "" {
		f, err := os.Open(conf.Conf)
		if err != nil {
			printAndExit(err)
		}
		defer f.Close()
		content, err := io.ReadAll(f)
		if err != nil {
			printAndExit(err)
		}
		hks, err := sharedclipboard.ParseHotkeys(string(content))
		if err != nil {
			printAndExit(err)
		}
		if hks != nil {
			hotkeys = hks
		}
	}

	// Start the application
	start(conf.Network)
}

func start(network string) {
	err := clipboard.Init()
	if err != nil {
		printAndExit(err)
	}
	locals := gopeers.PingSweep(network)
	peer := gopeers.NewPeer()

	logger, err := zap.NewProduction()
	if err != nil {
		printAndExit(err)
	}
	peer.SetLogger(logger)
	peer.Discover(locals)
	peer.Listen()

	sharedclipboard.Listen(peer, hotkeys)
}

func printAndExit(err error) {
	fmt.Println(err)
	os.Exit(1)
}

func getHomeDir() string {
	home := os.Getenv("HOME")
	if home == "" {
		home = os.Getenv("USERPROFILE")
	}
	return home
}
