package cmd

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/bpicode/fritzctl/meta"
	"github.com/mitchellh/cli"
	"github.com/stretchr/testify/assert"
)

// TestCommands is a unit test that runs most commands.
func TestCommands(t *testing.T) {
	testCases := []struct {
		cmd  cli.Command
		args []string
		srv  *httptest.Server
	}{
		{cmd: &pingCommand{}, srv: serverAnswering("testdata/loginresponse_test.xml")},
		{cmd: &listSwitchesCommand{}, srv: serverAnswering("testdata/loginresponse_test.xml", "testdata/loginresponse_test.xml", "testdata/devicelist_test.xml")},
		{cmd: &listThermostatsCommand{}, srv: serverAnswering("testdata/loginresponse_test.xml", "testdata/loginresponse_test.xml", "testdata/devicelist_test.xml")},
		{cmd: &switchOnCommand{}, args: []string{"My device"}, srv: serverAnswering("testdata/loginresponse_test.xml", "testdata/loginresponse_test.xml", "testdata/devicelist_test.xml", "testdata/answer_switch_on_test")},
		{cmd: &switchOffCommand{}, args: []string{"My device"}, srv: serverAnswering("testdata/loginresponse_test.xml", "testdata/loginresponse_test.xml", "testdata/devicelist_test.xml", "testdata/answer_switch_on_test")},
		{cmd: &temperatureCommand{}, args: []string{"19.5", "My device"}, srv: serverAnswering("testdata/loginresponse_test.xml", "testdata/loginresponse_test.xml", "testdata/devicelist_test.xml", "testdata/answer_switch_on_test")},
		{cmd: &toggleCommand{}, args: []string{"My device"}, srv: serverAnswering("testdata/loginresponse_test.xml", "testdata/loginresponse_test.xml", "testdata/devicelist_test.xml", "testdata/answer_switch_on_test")},
		{cmd: &sessionIDCommand{}, args: []string{}, srv: serverAnswering("testdata/loginresponse_test.xml", "testdata/loginresponse_test.xml")},
		{cmd: &listLandevicesCommand{}, args: []string{}, srv: serverAnswering("testdata/loginresponse_test.xml", "testdata/loginresponse_test.xml", "testdata/landevices_test.json")},
		{cmd: &listLogsCommand{}, args: []string{}, srv: serverAnswering("testdata/loginresponse_test.xml", "testdata/loginresponse_test.xml", "testdata/logs_test.json")},
	}
	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("Test run command %d", i), func(t *testing.T) {
			l, err := net.Listen("tcp", ":61666")
			assert.NoError(t, err)
			testCase.srv.Listener = l
			testCase.srv.Start()
			defer testCase.srv.Close()
			exitCode := testCase.cmd.Run(testCase.args)
			assert.Equal(t, 0, exitCode)
		})
	}
}

// TestConfigure tests the interactive configuration.
func TestConfigure(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "test_fritzctl")
	defer os.Remove(tempDir)
	assert.NoError(t, err)
	meta.DefaultConfigDir = tempDir
	cmd := configureCommand{}
	i := cmd.Run([]string{})
	assert.Equal(t, 0, i)
}

func serverAnswering(answers ...string) *httptest.Server {
	meta.ConfigDir = "testdata"
	meta.ConfigFilename = "config_localhost_test.json"
	it := 0
	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ch, _ := os.Open(answers[it%len(answers)])
		defer ch.Close()
		it++
		io.Copy(w, ch)
	}))
	return server
}

// TestCommandsHaveHelp ensures that every command provides
// a help text.
func TestCommandsHaveHelp(t *testing.T) {
	c := cli.NewCLI(meta.ApplicationName, meta.Version)
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"configure":       Configure,
		"listswitches":    ListSwitches,
		"listthermostats": ListThermostats,
		"listlandevices":  ListLandevices,
		"listlogs":        ListLogs,
		"ping":            Ping,
		"sessionid":       SessionID,
		"switchon":        SwitchOnDevice,
		"switchoff":       SwitchOffDevice,
		"toggle":          ToggleDevice,
		"temperature":     Temperature,
	}
	for i, command := range c.Commands {
		t.Run(fmt.Sprintf("Test help of command %s", i), func(t *testing.T) {
			com, err := command()
			assert.NoError(t, err)
			help := com.Help()
			fmt.Printf("Help on command %s: '%s'\n", i, help)
			assert.NotEmpty(t, help)
		})
	}
}

// TestCommandsHaveSynopsis ensures that every command provides
// short a synopsis text.
func TestCommandsHaveSynopsis(t *testing.T) {
	c := cli.NewCLI(meta.ApplicationName, meta.Version)
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"configure":       Configure,
		"listswitches":    ListSwitches,
		"listthermostats": ListThermostats,
		"listlandevices":  ListLandevices,
		"listlogs":        ListLogs,
		"ping":            Ping,
		"sessionid":       SessionID,
		"switchon":        SwitchOnDevice,
		"switchoff":       SwitchOffDevice,
		"toggle":          ToggleDevice,
		"temperature":     Temperature,
	}
	for i, command := range c.Commands {
		t.Run(fmt.Sprintf("Test synopsis of command %s", i), func(t *testing.T) {
			com, err := command()
			assert.NoError(t, err)
			syn := com.Synopsis()
			fmt.Printf("Synopsis on command '%s': '%s'\n", i, syn)
			assert.NotEmpty(t, syn)
		})
	}
}