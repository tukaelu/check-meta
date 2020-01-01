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

type chkOpts struct {
	Namespace string `short:"n" long:"namespace" required:"true" description:""`
	MetaKey   string `short:"k" long:"key" required:"true" description:""`
	Expected  string `short:"e" long:"expected" required:"true" description:""`
	apiKey    string
	hostID    string
}

func Do() {
	ckr := run(os.Args[1:])
	ckr.Name = "Meta"
	ckr.Exit()
}

func run(args []string) *checkers.Checker {
	opts := &chkOpts{}
	_, err := flags.ParseArgs(opts, args)
	if err != nil {
		return checkers.Unknown(err.Error())
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

	value, err := getHostMetaData(opts.apiKey, opts.hostID, opts.Namespace, opts.MetaKey)
	if err != nil {
		return checkers.Critical(err.Error())
	}

	return checkMetaValue(opts.Expected, value, opts.MetaKey)
}

func loadConfig() (*config.Config, error) {
	conf, err := config.LoadConfig(config.DefaultConfig.Conffile)
	if err != nil {
		return nil, fmt.Errorf("failed to load the config file: %s", err)
	}
	return conf, nil
}

func getHostMetaData(apiKey string, hostID string, namespace string, key string) (interface{}, error) {
	cli := mkr.NewClient(apiKey)
	meta, err := cli.GetHostMetaData(hostID, namespace)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%#v\n", meta)

	value, ok := meta.HostMetaData.(map[string]interface{})[key]
	if !ok {
		return nil, fmt.Errorf("meta key does not exists: %s", key)
	}
	return value, nil
}

func checkMetaValue(expected string, actual interface{}, metaKey string) *checkers.Checker {

	status := checkers.OK
	msg := fmt.Sprintf("matched")

	switch actual.(type) {
	case string:
		if actual != expected {
			status = checkers.CRITICAL
			msg = fmt.Sprintf("matched string: meta-key=%s, expected=%s, actual=%s", metaKey, expected, actual)
		}
	case float64:
		if converted, err := strconv.ParseFloat(expected, 64); err != nil {
			status = checkers.UNKNOWN
			msg = err.Error()
		} else if converted != actual {
			status = checkers.CRITICAL
			msg = fmt.Sprintf("matched float64 value: meta-key=%s, expected=%f, actual=%f", metaKey, converted, actual)
		}
	case bool:
		if converted, err := strconv.ParseBool(expected); err != nil {
			status = checkers.UNKNOWN
			msg = err.Error()
		} else if converted != actual {
			status = checkers.CRITICAL
			msg = fmt.Sprintf("matched boolean value: meta-key=%s, expected=%t, actual=%t", metaKey, converted, actual)
		}
	default:
		status = checkers.UNKNOWN
		msg = fmt.Sprintf("unsupported type value: type=%T, value=%v", actual, actual)
	}

	return checkers.NewChecker(status, msg)
}
