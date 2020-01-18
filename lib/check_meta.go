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
	Namespace      string `short:"n" long:"namespace" required:"true"  description:"Uses the metadata for the specified namespace"`
	MetaKey        string `short:"k" long:"key"       required:"true"  description:"The value matching the specified key is used for comparison"`
	Expected       string `short:"e" long:"expected"  required:"true"  description:"Compares with the specified expected value"`
	IsRegex        bool   `          long:"regex"     required:"false" description:"Compare with regular expression if specified (Enable only for string type value)"`
	GreaterThan    bool   `          long:"gt"        required:"false" description:"Compare as 'actual > expected' (Enable only for number type value)"`
	LessThan       bool   `          long:"lt"        required:"false" description:"Compare as 'actual < expected' (Enable only for number type value)"`
	GreaterOrEqual bool   `          long:"ge"        required:"false" description:"Compare as 'actual >= expected' (Enable only for number type value)"`
	LessOrEqual    bool   `          long:"le"        required:"false" description:"Compare as 'actual <= expected' (Enable only for number type value)"`
	apiKey         string
	hostID         string
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
	msg := ""
	status := checkers.OK

	switch actual.(type) {
	case string:
		status, msg = checkStringValue(actual.(string))
	case float64:
		if !isValidNumberComparisonOption() {
			return checkers.NewChecker(checkers.UNKNOWN, "gt/lt/ge/le options are only one can be specified")
		}
		status, msg = checkNumberValue(actual.(float64))
	case bool:
		status, msg = checkBooleanTypeValue(actual.(bool))
	default:
		status = checkers.UNKNOWN
		msg = fmt.Sprintf("unsupported type value: type=%T, value=%v", actual, actual)
	}

	return checkers.NewChecker(status, msg)
}

func checkStringValue(actual string) (checkers.Status, string) {
	var result bool
	var typeRegex string = ""

	reason := "matched"
	status := checkers.OK

	if opts.IsRegex {
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
		reason = "does not matched"
		status = checkers.CRITICAL
	}

	return status, fmt.Sprintf("%sstring %s: key=%s, actual=%s, expected=%s", typeRegex, reason, opts.MetaKey, actual, opts.Expected)
}

func checkNumberValue(actual float64) (checkers.Status, string) {
	var result bool
	var op string

	reason := "matched"
	status := checkers.OK

	expected, err := strconv.ParseFloat(opts.Expected, 64)
	if err != nil {
		return checkers.UNKNOWN, err.Error()
	}

	if opts.GreaterThan {
		result = actual > expected
		op = ">"
	} else if opts.LessThan {
		result = actual < expected
		op = "<"
	} else if opts.GreaterOrEqual {
		result = actual >= expected
		op = ">="
	} else if opts.LessOrEqual {
		result = actual <= expected
		op = "<="
	} else {
		result = actual == expected
		op = "=="
	}

	if !result {
		reason = "does not matched"
		status = checkers.CRITICAL
	}

	return status, fmt.Sprintf("number value %s: key=%s, actual(%f) %s expected(%f)", reason, opts.MetaKey, actual, op, expected)
}

func checkBooleanTypeValue(actual bool) (checkers.Status, string) {
	reason := "matched"
	status := checkers.OK

	expected, err := strconv.ParseBool(opts.Expected)
	if err != nil {
		return checkers.UNKNOWN, err.Error()
	}

	if actual != expected {
		reason = "does not matched"
		status = checkers.CRITICAL
	}

	return status, fmt.Sprintf("boolean value %s: key=%s, actual=%t expected=%t", reason, opts.MetaKey, actual, expected)
}

func isValidNumberComparisonOption() bool {
	optCnt := 0
	if opts.GreaterThan {
		optCnt++
	}
	if opts.LessThan {
		optCnt++
	}
	if opts.GreaterOrEqual {
		optCnt++
	}
	if opts.LessOrEqual {
		optCnt++
	}
	return optCnt <= 1
}
