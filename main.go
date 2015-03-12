package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"

	"github.com/alecthomas/kingpin"
)

type Size int64

func (size Size) String() string {
	nSize := float64(size)
	for _, unit := range []string{"", "Ki", "Mi", "Gi", "Ti", "Pi", "Ei", "Zi"} {
		if nSize < 1024.0 {
			return fmt.Sprintf("\"%.2f%sB\"", nSize, unit)
		}
		nSize /= 1024.0
	}
	return fmt.Sprintf("\"%.1fYiB\"", nSize)
}

func (size Size) MarshalJSON() ([]byte, error) {
	return []byte(size.String()), nil
}

type DirInfo struct {
	FullPath string     `json:"fullPath"`
	Path     string     `json:"path"`
	Size     Size       `json:"size"`
	SubDirs  []*DirInfo `json:"subDirs,omitempty"`
}

func visit(dir *DirInfo, minSize Size) {
	var err error
	var file *os.File
	if file, err = os.Open(dir.FullPath); err != nil {
		return
	}
	var files []os.FileInfo
	if files, err = file.Readdir(-1); err != nil {
		file.Close()
		return
	}
	file.Close()
	var fi os.FileInfo
	var subDir *DirInfo
	for _, fi = range files {
		if fi.IsDir() {
			subDir = &DirInfo{
				Path:     fi.Name(),
				FullPath: filepath.Join(dir.FullPath, fi.Name()),
			}
			visit(subDir, minSize)
			dir.Size += subDir.Size
			if subDir.Size >= minSize {
				if dir.SubDirs == nil {
					dir.SubDirs = make([]*DirInfo, 0, 2)
				}
				dir.SubDirs = append(dir.SubDirs, subDir)
			}
		} else if fi.Mode().IsRegular() {
			dir.Size += Size(fi.Size())
		}
	}
	if dir.SubDirs != nil &&
		len(dir.SubDirs) == 1 {
		subDir = dir.SubDirs[0]
		dir.FullPath = subDir.FullPath
		dir.Path = filepath.Join(dir.Path, subDir.Path)
		dir.Size = subDir.Size
		dir.SubDirs = subDir.SubDirs
	}
}

const (
	B   int64 = 1
	KB        = 1000
	KiB       = 1024
	MB        = 1000 * KB
	MiB       = 1024 * KiB
	GB        = 1000 * MB
	GiB       = 1024 * MiB
	TB        = 1000 * GB
	TiB       = 1024 * GiB
)

var (
	fPath = kingpin.Flag("path", "The root directory.").Short('p').
		Required().String()
	fSize = kingpin.Flag("size", "The minimum size per directory. Eg. 500MiB").
		OverrideDefaultFromEnvar("500MiB").Short('s').String()
	fOut    = kingpin.Flag("out", "The output file. Writes to Stdout when not specified.").Short('o').String()
	fIndent = kingpin.Flag("indent", "Indent the outputted JSON.").
		OverrideDefaultFromEnvar("true").Short('i').Bool()
)

const defaultSize = Size(500 * MiB)

var modifiers = map[string]int64{
	"kb": KB, "kib": KiB,
	"mb": MB, "mib": MiB,
	"gb": GB, "gib": GiB,
	"tb": TB, "tib": TiB,
}

func parseSize() Size {
	if fSize == nil {
		return defaultSize
	}
	var b bytes.Buffer
	input := strings.TrimSpace(strings.ToLower(*fSize))
	dec := false
	var i int
	var r rune
	for i, r = range input {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
			continue
		}
		if r == '.' && !dec {
			b.WriteRune('.')
			dec = true
			continue
		}
		if b.Len() > 0 {
			// skip spaces
			if unicode.IsSpace(r) {
				continue
			}
			break
		}
	}
	if b.Len() == 0 {
		return defaultSize
	}
	ffSize, _ := strconv.ParseFloat(b.String(), 64)
	if mod, ok := modifiers[input[i:]]; ok {
		return Size(ffSize * float64(mod))
	}
	return Size(ffSize)
}

func main() {
	kingpin.Version("0.0.1")
	kingpin.Parse()

	var err error
	path := strings.TrimSpace(*fPath)
	if _, err = os.Stat(path); err != nil {
		os.Stderr.WriteString(err.Error())
		return
	}

	var out string
	if fOut != nil {
		out = strings.TrimSpace(*fOut)
	}
	indent := (fIndent != nil && (*fIndent))
	isHtml := false
	switch strings.ToLower(filepath.Ext(out)) {
	case ".json":
	case ".html":
		isHtml = true
	default:
		if out != "" {
			os.Stderr.WriteString("-out file name must have an extension of .json or .html")
			return
		}
	}

	dir := &DirInfo{FullPath: path, Path: path, SubDirs: make([]*DirInfo, 0)}
	visit(dir, parseSize())

	var w io.Writer
	if out == "" {
		w = os.Stdout
	} else {
		var f *os.File
		if f, err = os.Create(out); err != nil {
			os.Stderr.WriteString(err.Error())
			return
		}
		w = f
		defer f.Close()
	}

	var p []byte
	if indent && !isHtml {
		p, err = json.MarshalIndent(dir, "", "  ")
	} else {
		p, err = json.Marshal(dir)
	}
	if err != nil {
		os.Stderr.WriteString(err.Error())
		return
	}

	if out != "" && isHtml {
		tmpl.Execute(w, string(p))
	} else {
		w.Write(p)
	}
}
