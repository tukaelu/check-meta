package checkmeta

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"

	"github.com/jessevdk/go-flags"
	"github.com/mackerelio/checkers"
	"github.com/mackerelio/mackerel-agent/config"
	mkr "github.com/mackerelio/mackerel-client-go"
)

var opts struct {
	Namespace        string `short:"n" long:"namespace"         required:"true"  value-name:"NAMESPACE"      description:"Uses the metadata for the specified namespace"`
	MetaKey          string `short:"k" long:"key"               required:"true"  value-name:"KEY"            description:"The value matching the specified key is used for comparison"`
	Expected         string `short:"e" long:"expected"          required:"false" value-name:"EXPECTED-VALUE" description:"Compares with the specified expected value"`
	IsRegex          bool   `          long:"regex"             required:"false"                             description:"Compare with regular expression if specified (Enable only for string type value)"`
	GreaterThan      bool   `          long:"gt"                required:"false"                             description:"Compare as 'actual > expected' (Enable only for number type value)"`
	LessThan         bool   `          long:"lt"                required:"false"                             description:"Compare as 'actual < expected' (Enable only for number type value)"`
	GreaterOrEqual   bool   `          long:"ge"                required:"false"                             description:"Compare as 'actual >= expected' (Enable only for number type value)"`
	LessOrEqual      bool   `          long:"le"                required:"false"                             description:"Compare as 'actual <= expected' (Enable only for number type value)"`
	CompareNamespace string `short:"N" long:"compare-namespace" required:"false" value-name:"NAMESPACE"      description:"Uses the metadata for the specified namespace to compare"`
	CompareMetaKey   string `short:"K" long:"compare-key"       required:"false" value-name:"KEY"            description:"Uses the metadata value that matches the specified key as the expected value"`
	apiKey           string
	hostID           string
	compareMetaValue interface{}
}

func Do() {
	ckr := run(os.Args[1:])
	ckr.Name = "Meta"
	ckr.Exit()
}

func run(args []string) *checkers.Checker {
	origArgs := make([]string, len(args))
	copy(origArgs, args)

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

	value, err := getHostMetaData(opts.hostID, opts.Namespace, opts.MetaKey)
	if err != nil {
		return checkers.Critical(err.Error())
	}

	if shouldUseMetadataForCompare() {
		if opts.CompareNamespace == "" {
			opts.CompareNamespace = opts.Namespace
		}
		if opts.CompareMetaKey == "" {
			opts.CompareMetaKey = opts.MetaKey
		}

		cacheFile := getCacheFile(origArgs)

		compareMetaValue, err := getHostMetaData(opts.hostID, opts.CompareNamespace, opts.CompareMetaKey)
		if err != nil {
			cache, err := loadCache(cacheFile)
			if err != nil {
				return checkers.Unknown(err.Error())
			}
			if cache.Expected == nil {
				return checkers.Unknown("there is no data in the cache.")
			}
			compareMetaValue = cache.Expected
		} else {
			if err := saveCache(cacheFile, &cache{Options: origArgs, Expected: compareMetaValue}); err != nil {
				log.Printf("failed to saveCache: %s", err)
			}
		}
		opts.compareMetaValue = compareMetaValue
	}

	if opts.Expected == "" && opts.compareMetaValue == nil {
		return checkers.Unknown("Expected value not specified. Specify the --exptected or --compare-namespace/--compare-key options.")
	}

	return checkMetadata(value)
}

func loadConfig() (*config.Config, error) {
	conf, err := config.LoadConfig(config.DefaultConfig.Conffile)
	if err != nil {
		return nil, fmt.Errorf("failed to load the config file: %s", err)
	}
	return conf, nil
}

func getHostMetaData(hostID, namespace, metaKey string) (interface{}, error) {
	cli := mkr.NewClient(opts.apiKey)
	meta, err := cli.GetHostMetaData(hostID, namespace)
	if err != nil {
		return nil, err
	}

	value, ok := meta.HostMetaData.(map[string]interface{})[metaKey]
	if !ok {
		return nil, fmt.Errorf("meta key does not exists: %s", metaKey)
	}
	return value, nil
}

func checkMetadata(actual interface{}) *checkers.Checker {
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
	expected := ""

	if shouldUseMetadataForCompare() {
		compareMetaValue, err := opts.compareMetaValue.(string)
		if !err {
			return checkers.UNKNOWN, fmt.Sprintf("unmatched compare type: actual=string, expected=%T", opts.compareMetaValue)
		}
		expected = compareMetaValue
	} else {
		expected = opts.Expected
	}

	if opts.IsRegex {
		regExpected, err := regexp.Compile(expected)
		if err != nil {
			return checkers.UNKNOWN, err.Error()
		}
		result = regExpected.MatchString(actual)
		typeRegex = "regex-"
	} else {
		result = expected == actual
	}

	if !result {
		reason = "does not matched"
		status = checkers.CRITICAL
	}

	return status, fmt.Sprintf("%sstring %s: key=%s, actual=%s, expected=%s", typeRegex, reason, opts.MetaKey, actual, expected)
}

func checkNumberValue(actual float64) (checkers.Status, string) {
	var result bool
	var op string

	reason := "matched"
	status := checkers.OK
	expected := float64(0)

	if shouldUseMetadataForCompare() {
		compareMetaValue, err := opts.compareMetaValue.(float64)
		if !err {
			return checkers.UNKNOWN, fmt.Sprintf("unmatched compare type: actual=float64, expected=%T", opts.compareMetaValue)
		}
		expected = compareMetaValue
	} else {
		var err error
		expected, err = strconv.ParseFloat(opts.Expected, 64)
		if err != nil {
			return checkers.UNKNOWN, err.Error()
		}
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
	expected := true

	if shouldUseMetadataForCompare() {
		compareMetaValue, err := opts.compareMetaValue.(bool)
		if !err {
			return checkers.UNKNOWN, fmt.Sprintf("unmatched compare type: actual=bool, expected=%T", opts.compareMetaValue)
		}
		expected = compareMetaValue
	} else {
		var err error
		expected, err = strconv.ParseBool(opts.Expected)
		if err != nil {
			return checkers.UNKNOWN, err.Error()
		}
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

func shouldUseMetadataForCompare() bool {
	return !(opts.CompareNamespace == "" && opts.CompareMetaKey == "")
}
