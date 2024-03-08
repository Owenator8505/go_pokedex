package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/Owenator8505/pokedexcli/api"
)

type CliCommand struct {
	name     string
	desc     string
	callback func(*CliConfigs) error
}

type CliConfigs struct {
	opts            map[string]CliCommand
	signals         chan os.Signal
	locationPayload api.LocationPayload
}

func CommandHelp(params *CliConfigs) error {
	fmt.Printf("%s\n\n", "Welcome to the Pokedex!")
	fmt.Println("Select an option to continue:")

	for _, opt := range params.opts {
		fmt.Printf("%s: %s \n", opt.name, opt.desc)
	}

	return nil
}

func CommandExit(params *CliConfigs) error {
	params.signals <- syscall.SIGINT
	return nil
}

func CommandMap(params *CliConfigs) error {
	api.GetLocationsHandler(&params.locationPayload)
	log.Printf("Location Payload: %v", params.locationPayload)
	return nil
}

func CommandGetter(key string, opts *CliConfigs) *CliCommand {
	value, ok := opts.opts[key]
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
		case <-opts.signals:
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
				_ = command.callback(&opts)
				fmt.Print("\nPokedex > ")
			} else {
				fmt.Println("Invalid command")
			}
		}
	}
}

func InitOpts(configs *CliConfigs) {
	configs.opts = map[string]CliCommand{
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
		"map": {
			name:     "map",
			desc:     "List location areas",
			callback: CommandMap,
		},
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	configs.signals = signals

	locationPayload := api.LocationPayload{}
	locationPayload.Id = 0
	locationPayload.Name = ""
	locationPayload.Limit = 20
	locationPayload.Offset = 0
	configs.locationPayload = locationPayload
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	commandOpts := CliConfigs{}
	InitOpts(&commandOpts)
	StartREPL(reader, commandOpts)
}
