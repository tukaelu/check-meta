package checkmeta

import (
	"fmt"
	"os"
	"strconv"

	"github.com/jessevdk/go-flags"
	"github.com/mackerelio/checkers"
	"github.com/mackerelio/mackerel-agent/config"
	mkr "github.com/mackerelio/mackerel-client-go"
)

var opts struct {
	Namespace string `short:"n" long:"namespace" required:"true" description:"Uses the metadata for the specified namespace"`
	MetaKey   string `short:"k" long:"key"       required:"true" description:"The value matching the specified key is used for comparison"`
	Expected  string `short:"e" long:"expected"  required:"true" description:"Compares with the specified expected value"`
	apiKey    string
	hostID    string
}

func Do() {
	ckr := run(os.Args[1:])
	ckr.Name = "Meta"
	ckr.Exit()
}

func run(args []string) *checkers.Checker {
	_, err := flags.ParseArgs(opts, args)
	if err != nil {
		os.Exit(1)
	}

	conf, err := loadConfig()
	if err != nil {
		return checkers.Unknown(err.Error())
	}
	opts.apiKey = conf.Apikey

	hostID, err := conf.LoadHostID()
	if err != nil {
		return checkers.Unknown(err.Error())
	}
	opts.hostID = hostID

	value, err := getHostMetaData()
	if err != nil {
		return checkers.Critical(err.Error())
	}

	return checkMetaValue(value)
}

func loadConfig() (*config.Config, error) {
	conf, err := config.LoadConfig(config.DefaultConfig.Conffile)
	if err != nil {
		return nil, fmt.Errorf("failed to load the config file: %s", err)
	}
	return conf, nil
}

func getHostMetaData() (interface{}, error) {
	cli := mkr.NewClient(opts.apiKey)
	meta, err := cli.GetHostMetaData(opts.hostID, opts.Namespace)
	if err != nil {
		return nil, err
	}

	value, ok := meta.HostMetaData.(map[string]interface{})[opts.MetaKey]
	if !ok {
		return nil, fmt.Errorf("meta key does not exists: %s", opts.MetaKey)
	}
	return value, nil
}

func checkMetaValue(actual interface{}) *checkers.Checker {

	status := checkers.OK
	msg := fmt.Sprintf("matched value")

	switch actual.(type) {
	case string:
		if actual != opts.Expected {
			status = checkers.CRITICAL
			msg = fmt.Sprintf("unmatched string value: key=%s, expected=%s, actual=%s", opts.MetaKey, opts.Expected, actual)
		}
	case float64:
		if converted, err := strconv.ParseFloat(opts.Expected, 64); err != nil {
			status = checkers.UNKNOWN
			msg = err.Error()
		} else if converted != actual {
			status = checkers.CRITICAL
			msg = fmt.Sprintf("unmatched float64 value: key=%s, expected=%f, actual=%f", opts.MetaKey, converted, actual)
		}
	case bool:
		if converted, err := strconv.ParseBool(opts.Expected); err != nil {
			status = checkers.UNKNOWN
			msg = err.Error()
		} else if converted != actual {
			status = checkers.CRITICAL
			msg = fmt.Sprintf("unmatched boolean value: key=%s, expected=%t, actual=%t", opts.MetaKey, converted, actual)
		}
	default:
		status = checkers.UNKNOWN
		msg = fmt.Sprintf("unsupported type value: type=%T, value=%v", actual, actual)
	}

	return checkers.NewChecker(status, msg)
}
