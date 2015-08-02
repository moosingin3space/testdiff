package testdiff

import (
	"github.com/BurntSushi/toml"
	"github.com/fatih/color"
	"github.com/natefinch/pie"
	"io/ioutil"
	"net/rpc/jsonrpc"
	"os"
	"os/exec"
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
	var changedFiles []string
	var err error
	var workingDir string
	codec := jsonrpc.NewClientCodec

	plugin, err := pie.StartProviderCodec(codec, os.Stderr, config.Differ.Command, config.Differ.Args...)

	if err != nil {
		return nil, err
	}

	if workingDir, err = os.Getwd(); err != nil {
		return nil, err
	}
	if err = plugin.Call("Differ.Diff", workingDir, changedFiles); err != nil {
		return nil, err
	}
	return changedFiles, nil
}

func determineTestsToExecute(config Config, changedFiles []string) ([]string, error) {
	var testsToRun []string
	var err error
	codec := jsonrpc.NewClientCodec

	plugin, err := pie.StartProviderCodec(codec, os.Stderr, config.Locator.Command, config.Locator.Args...)

	if err != nil {
		return nil, err
	}

	if err = plugin.Call("Locator.Locate", changedFiles, testsToRun); err != nil {
		return nil, err
	}
	return testsToRun, nil
}

func executeTests(config Config, tests []string) error {
	var commands []string
	var err error
	codec := jsonrpc.NewClientCodec

	plugin, err := pie.StartProviderCodec(codec, os.Stderr, config.Executor.Command, config.Executor.Args...)

	if err != nil {
		return err
	}

	if err = plugin.Call("Executor.GetCommands", tests, commands); err != nil {
		return err
	}

	for _, command := range commands {
		// FIXME this only works on Unixes
		cmd := exec.Command("sh", "-c", command)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	var configFileContents []byte
	var conf Config
	var changedFiles []string
	var tests []string
	var err error

	if configFileContents, err = ioutil.ReadFile("testdiff.toml"); err != nil {
		color.Red("Loading configuration failed!")
		return
	}

	if _, err := toml.Decode(string(configFileContents), &conf); err != nil {
		color.Red("Loading configuration failed!")
		return
	}

	color.Blue("Scanning for changed files...")
	if changedFiles, err = findChangedFiles(conf); err != nil {
		color.Red("Diffing failed!")
		return
	}

	color.Blue("Determining what tests to run...")
	if tests, err = determineTestsToExecute(conf, changedFiles); err != nil {
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
