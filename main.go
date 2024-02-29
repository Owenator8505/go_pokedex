package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

type CliCommand struct {
	name     string
	desc     string
	callback func(CliConfigs) error
}

type CliConfigs struct {
	Opts    map[string]CliCommand
	Signals chan os.Signal
}

func CommandHelp(params CliConfigs) error {
	fmt.Printf("%s\n\n", "Welcome to the Pokedex!")
	fmt.Println("Select an option to continue:")

	for _, opt := range params.Opts {
		fmt.Printf("%s: %s \n", opt.name, opt.desc)
	}

	return nil
}

func CommandExit(params CliConfigs) error {
	params.Signals <- syscall.SIGINT
	return nil
}

func CommandGetter(key string, opts *CliConfigs) *CliCommand {
	value, ok := opts.Opts[key]
	if ok {
		return &value
	}

	return nil
}

func SanitizeInput(input string) string {
	return strings.ToLower(strings.TrimSpace(input))
}

func StartREPL(reader *bufio.Reader, opts CliConfigs) {
	fmt.Print("Pokedex > ")
loop:
	for {
		select {
		case <-opts.Signals:
			fmt.Println("C-Ya later!")
			break loop
		default:
			key, err := reader.ReadString('\n')
			key = SanitizeInput(key)
			if err != nil {
				panic(err)
			}

			command := CommandGetter(key, &opts)

			if command != nil {
				_ = command.callback(opts)
				fmt.Print("\nPokedex > ")
			}
		}
	}
}

func InitOpts(commandOpts *CliConfigs) {
	commandOpts.Opts = map[string]CliCommand{
		"help": {
			name:     "help",
			desc:     "Display a help message",
			callback: CommandHelp,
		},
		"exit": {
			name:     "exit",
			desc:     "Exit the Pokedex",
			callback: CommandExit,
		},
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	commandOpts.Signals = signals
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	commandOpts := CliConfigs{}
	InitOpts(&commandOpts)
	StartREPL(reader, commandOpts)
}
