# check-meta

## Install

```sh
$ sudo mkr plugin install tukaelu/check-meta
```

## Usage

```
Usage:
  check-meta [OPTIONS]

Application Options:
  -n, --namespace= Uses the metadata for the specified namespace
  -k, --key=       The value matching the specified key is used for comparison
  -e, --expected=  Compares with the specified expected value

Help Options:
  -h, --help       Show this help message
```

Supported expected value types are

- string
- float64
- bool

## Configuration

```
[plugin.checks.meta-namespace-key]
command = ["/path/to/check-meta", "-n", "namespace", "-k", "key", "-e", "expected"]
```

## Sample
```
# GET /api/v0/hosts/<hostId>/metadata/namespace
# {
#   "key1": "value1",
#   "key2": 1000,
#   "key3": true,
# }

## OK
[plugin.checks.meta_string]
command = ["/path/to/check-meta", "-n", "namespace", "-k", "key1", "-e", "value1"]

## CRITICAL
[plugin.checks.meta_string]
command = ["/path/to/check-meta", "-n", "namespace", "-k", "key2", "-e", 1001]
```