# Remove GC logs

They always say "turn on the JVM GC logs, they cost nothing really", but they never say how hard is to filter out
all those lines.

This application removes those lines by default, but you can use regex to remove other types of lines.

Algorithm is:

```go
if filteredOut(line) {
    s = PassThrough
} else if isLogLine(line) {
    s = Reading
}
if s == Reading {
    fmt.Println(line)
}
```

## Installation

```
go get -u github.com/milanaleksic/gclog-cleaner
```

## Usage

```
milan@MilanMBP (master) ~/SourceCode/log-gc-cleaner → gclog-cleaner --help
Usage of gclog-cleaner:
  -begin-pattern string
        pattern that should match beginning of all log lines (default "\\d{4}-\\d{2}-\\d{2} \\d{2}:\\d{2}:\\d{2}")
  -exclusions value
        exclusion patterns (default: only one - '\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}.*\[GC')
  -input-file string
        which file to process (default - stdin)
```

## Example

Example with progress provided via `pv` utility

```bash
pv -i 0.1 biglogfile.txt | gclog-cleaner \
  --exclusions '\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}.*\[GC' \
  --exclusions "Creating an interceptor chain" \
  --exclusions 'request:84' \
  --exclusions 'DEBUG Sdk' \ 
  > clean.txt
```
