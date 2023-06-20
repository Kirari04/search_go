# search_go

![Screenshot](image.png)

## Build

```bash
go build main.go
```

## Default Config

| env      | default | accepts |
| -------- | ------- | ------- |
| ROOTDIR  | C:\     | string  |
| MAXDEBTH | 10      | int     |
| SILENT   | 1       | bool    |
| REGEX    | 1       | bool    |

## Run With Custom Config

In Powershell

```shell
$env:MAXDEBTH=10
$env:SILENT=0
$env:REGEX=1
$env:ROOTDIR='C:\'
.\main.exe
```

In Bash

```bash
MAXDEBTH=10 SILENT=0 REGEX=1 ROOTDIR='C:\' ./main
```

## Example Search Term

Listing all .mp4 files using regex

```
^*.mp4$
```