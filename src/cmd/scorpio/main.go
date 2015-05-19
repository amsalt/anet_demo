package main

import (
	getopt "github.com/kesselborn/go-getopt"

	"config"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"server"
	"syscall"
)

var (
	_VERSION_   = "unknown"
	_BUILDDATE_ = "unknown"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	optionDefinition := getopt.Options{
		"description",
		getopt.Definitions{
			{"config|c", "config file", getopt.IsConfigFile | getopt.ExampleIsDefault, "conf/server.conf"},
			{"version|v", "show version", getopt.Optional | getopt.Flag, nil},
		},
	}

	options, _, _, e := optionDefinition.ParseCommandLine()
	help, wantsHelp := options["help"]
	if e != nil || wantsHelp {
		exit_code := 0
		switch {
		case wantsHelp && help.String == "usage":
			fmt.Print(optionDefinition.Usage())
		case wantsHelp && help.String == "help":
			fmt.Print(optionDefinition.Help())
		default:
			fmt.Println("**** Error: ", e.Error(), "\n", optionDefinition.Help())
			exit_code = e.ErrorCode
		}
		os.Exit(exit_code)
	}
	version, showVersion := options["version"]
	if showVersion && version.Bool {
		fmt.Printf("server version %s\n%s\n", _VERSION_, _BUILDDATE_)
		os.Exit(0)
	}

	cfg, err := config.Load(options["config"].String)
	if err != nil {
		log.Println("load config falied:", err)
		os.Exit(0)
	}
	if app, err := server.NewApp(cfg); err != nil {
		log.Println(err.Error())
	} else {
		sc := make(chan os.Signal, 1)
		signal.Notify(sc,
			os.Kill,
			os.Interrupt,
			syscall.SIGHUP,
			syscall.SIGINT,
			syscall.SIGTERM,
			syscall.SIGQUIT)
		go app.Run()
		<-sc
		app.Close()
	}
}
