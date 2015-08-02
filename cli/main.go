package main

import (
	"github.com/BurntSushi/toml"
	"github.com/fatih/color"
	"github.com/natefinch/pie"
	"io/ioutil"
	"net/rpc/jsonrpc"
)

type PluginCommand struct {
	Command string
	Args    []string
}

type Config struct {
	Differ   PluginCommand
	Locator  PluginCommand
	Executor PluginCommand
}

func findChangedFiles(config Config) ([]string, error) {
	var err error
	codec := jsonrpc.NewClientCodec

	plugin, err := pie.StartProviderCodec(codec, os.StdErr, config.Differ.Command, config.Differ.Args...)

	if err != nil {
		return nil, err
	}

	// TODO interact with plugin
}

func determineTestsToExecute(config Config, changedFiles []string) ([]string, error) {
}

func executeTests(config Config, tests []string) error {
}

func main() {
	var conf Config
	var err error

	if configFileContents, err := ioutil.ReadFile("local"); err != nil {
		color.Red("Loading configuration failed!")
		return
	}

	if _, err := toml.Decode(configFileContents, &conf); err != nil {
		color.Red("Loading configuration failed!")
		return
	}

	color.Blue("Scanning for changed files...")
	if changedFiles, err := findChangedFiles(conf); err != nil {
		color.Red("Diffing failed!")
		return
	}

	color.Blue("Determining what tests to run...")
	if tests, err := determineTestsToExecute(conf, changedFiles); err != nil {
		color.Red("Test locating failed!")
		return
	}

	color.Blue("Executing tests.")
	if err = executeTests(conf, tests); err != nil {
		color.Red("Running tests failed!")
		return
	}

	color.Green("All done!")
}
