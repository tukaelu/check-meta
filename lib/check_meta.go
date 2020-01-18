package checkmeta

import (
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/jessevdk/go-flags"
	"github.com/mackerelio/checkers"
	"github.com/mackerelio/mackerel-agent/config"
	mkr "github.com/mackerelio/mackerel-client-go"
)

var opts struct {
	Namespace string `short:"n" long:"namespace" required:"true"  description:"Uses the metadata for the specified namespace"`
	MetaKey   string `short:"k" long:"key"       required:"true"  description:"The value matching the specified key is used for comparison"`
	Expected  string `short:"e" long:"expected"  required:"true"  description:"Compares with the specified expected value"`
	IsRegex   bool   `          long:"regex"     required:"false" description:"Compare with regular expression if specified (Enable only for string type value)"`
	apiKey    string
	hostID    string
}

func Do() {
	ckr := run(os.Args[1:])
	ckr.Name = "Meta"
	ckr.Exit()
}

func run(args []string) *checkers.Checker {
	_, err := flags.ParseArgs(&opts, args)
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
		status, msg = checkStringValue(opts.MetaKey, opts.Expected, actual.(string), opts.IsRegex)
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

func checkStringValue(key string, expected string, actual string, isRegex bool) (checkers.Status, string) {
	var result bool
	var reason string
	var typeRegex string = ""

	if isRegex {
		regExpected, err := regexp.Compile(opts.Expected)
		if err != nil {
			return checkers.UNKNOWN, err.Error()
		}
		result = regExpected.MatchString(actual)
		typeRegex = "regex-"
	} else {
		result = opts.Expected == actual
	}

	if !result {
		reason = "unmatched %sstring value: key=%s, expected=%s, actual=%s"
		return checkers.CRITICAL, fmt.Sprintf(reason, typeRegex, key, expected, actual)
	}
	reason = fmt.Sprintf("%sstring matched: key=%s, expected=%s, actual=%s", typeRegex, key, expected, actual)

	return checkers.OK, reason
}
