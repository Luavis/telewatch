package main

import (
	"errors"
	"fmt"
	"github.com/luavis/telewatch"
	daemon "gopkg.in/sevlyar/go-daemon.v0"
	"os"
	"strconv"
)

const ProgramName = "telewatch"

func help(message string) {

	fmt.Println(message)
	fmt.Println("Usage:")
	fmt.Printf("  %s register\n", ProgramName)
	fmt.Printf("  %s alert\n", ProgramName)
	fmt.Printf("  %s [options] [commands]\n", ProgramName)
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  -n, --interval <secs>: seconds to wait between updates")
	fmt.Println("  -d: demonize watch process")
}

type watchOptions struct {
	interval int
	daemonize bool
	command []string
}

type watchOptionParseMode uint8
const (
	watchOptionParseModeNormal = watchOptionParseMode(iota)
	watchOptionParseModeInterval
)

func parseWatchOption(args []string) (watchOptions, error) {
	cursor := 0
	mode := watchOptionParseModeNormal
	options := watchOptions{
		interval: 2,
		daemonize: false,
		command: []string {},
	}

	for cursor < len(args) {
		arg := args[cursor]
		switch mode {
		case watchOptionParseModeNormal:
			if arg[0] == '-' && len(arg) >= 2 {
				argSubstring := arg[1:]
				switch argSubstring {
				case "n", "-interval":
					mode = watchOptionParseModeInterval
				case "d":
					options.daemonize = true
				default:
					return options, fmt.Errorf("invalid option: %s", arg)
				}
			} else {
				options.command = args[cursor:]
				return options, nil
			}
		case watchOptionParseModeInterval:
			secs, err := strconv.Atoi(arg)
			if err != nil {
				return options, errors.New("interval option must be number")
			}

			options.interval = secs
			mode = watchOptionParseModeNormal
		default:
			panic("Invalid parse mode")
		}

		cursor += 1
	}

	return options, nil
}

func main() {
	argLength := len(os.Args)
	if argLength < 2 {
		help("")
		return
	}

	first := os.Args[1]
	config, err := telewatch.LoadConfigurationFromHomeDirectory()
	if err != nil {
		telewatch.PrintErrorAndExit("Fail to load config file", err)
	}

	if first == "register" {
		if argLength > 2 {
			help(fmt.Sprintf("Invalid arguments %s", os.Args[2:]))
			return
		}
		err := telewatch.RegisterChatId(config)
		if err != nil {
			telewatch.PrintErrorAndExit("Fail to register chat id", err)
		} else {
			fmt.Println("Chat id is successfully saved")
		}
	} else if first == "alert" {
		if argLength > 2 {
			help(fmt.Sprintf("Invalid arguments %s", os.Args[2:]))
			return
		}
		err := telewatch.Alert(config, "Alert!")
		if err != nil {
			telewatch.PrintErrorAndExit("Fail to alert", err)
		}
	} else {
		if argLength < 3 {
			help("")
			return
		}

		options, err := parseWatchOption(os.Args[1:])
		if err != nil {
			help(err.Error())
			return
		}

		if len(options.command) < 1 {
			help("")
			return
		}

		if options.daemonize {
			context := new(daemon.Context)
			child, _ := context.Reborn()

			// daemonize process
			if child == nil {
				defer context.Release()
				err = telewatch.Watch(config, options.interval, options.command)
			}
		} else {
			err = telewatch.Watch(config, options.interval, options.command)
		}

		if err != nil {
			telewatch.PrintErrorAndExit("Fail to watch", err)
		}
	}
}
