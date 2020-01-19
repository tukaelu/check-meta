# check-meta

## Description

Check that the host metadata is the expected value.

## Installation

```sh
$ sudo mkr plugin install tukaelu/check-meta
```

## Usage

```
Usage:
  check-meta [OPTIONS]

Application Options:
  -n, --namespace=NAMESPACE            Uses the metadata for the specified namespace
  -k, --key=KEY                        The value matching the specified key is used for comparison
  -e, --expected=EXPECTED-VALUE        Compares with the specified expected value
      --regex                          Compare with regular expression if specified (Enable only for string type value)
      --gt                             Compare as 'actual > expected' (Enable only for number type value)
      --lt                             Compare as 'actual < expected' (Enable only for number type value)
      --ge                             Compare as 'actual >= expected' (Enable only for number type value)
      --le                             Compare as 'actual <= expected' (Enable only for number type value)
  -N, --compare-namespace=NAMESPACE    Uses the metadata for the specified namespace to compare
  -K, --compare-key=KEY                Uses the metadata value that matches the specified key as the expected value

Help Options:
  -h, --help       Show this help message
```

Supported expected value types are...

- string
- float64 (JSON number)
- bool

## Configuration

```
[plugin.checks.meta-namespace-key]
command = ["/path/to/check-meta", "--namespace", "namespace", "--key", "key", "--expected", "expected"]

[plugin.checks.meta-namespace-key-short]
command = ["/path/to/check-meta", "-n", "namespace", "-k", "key", "-e", "expected"]
```

## Examples
```
# GET /api/v0/hosts/<hostId>/metadata/namespace
# {
#   "key1": "value1",
#   "key2": 1000,
#   "key3": true,
#   "key4": "value1",
# }

## OK (match)
[plugin.checks.meta_string]
command = ["/path/to/check-meta", "-n", "namespace", "-k", "key1", "-e", "value1"]

## OK (regex match)
[plugin.checks.meta_string_regex]
command = ["/path/to/check-meta", "-n", "namespace", "-k", "key1", "-e", "value[0-1]{1}", "--regex"]

## OK (number lower than expected)
[plugin.checks.meta_greater_than]
command = ["/path/to/check-meta", "-n", "namespace", "-k", "key2", "-e", "2000", "--lt"]

## OK (boolean match)
[plugin.checks.meta_greater_than]
command = ["/path/to/check-meta", "-n", "namespace", "-k", "key3", "-e", "true"]

## OK (compare with metadata)
[plugin.checks.meta_compare_metadata]
command = ["/path/to/check-meta", "-n", "namespace", "-k", "key1", "-K", "key4"]

## CRITICAL (does not match)
[plugin.checks.meta_string]
command = ["/path/to/check-meta", "-n", "namespace", "-k", "key2", "-e", 1001]
```

## For more information
- Execute `check-meta -h` and you can get command line options.

## Misc.
- [Metadata - Mackerel API Documents (v0)](https://mackerel.io/api-docs/entry/metadata)
