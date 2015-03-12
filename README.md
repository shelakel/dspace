# dspace
dspace is a (Visual) Disk Space Analyzer written in Go

dspace was written to solve a simple problem:

 > dspace allows you to visually analyze used disk space by identifying directories exceeding a certain size.

## Usage

```bash
usage: dspace --path=PATH [<flags>]

Flags:
  --help           Show help.
  -p, --path=PATH  The root directory.
  -s, --size=SIZE  The minimum size per directory. Eg. 500MiB
  -o, --out=OUT    The output file. Writes to Stdout when not specified.
  -i, --indent      Indent the outputted JSON.
  --version        Show application version.
```

 - Only .html and .json file extensions are supported for -o/--out files.
 - Using a .html extension will generate a HTML file to visually traverse the file system starting at the root -p/--path

This tool is an evening project to solve a problem I'm having at work of running out of space on my 240GB primary disk SSD.

## License (MIT)

See LICENSE