package main

import (
	"flag"
	"fmt"
	"os"
	"os/user"
	"path"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"syscall"
	"time"

	"edls/fileinfo"

	"github.com/fatih/color"
	"golang.org/x/exp/constraints"
)

func main() {
	// Filters
	flagPattern := flag.String("p", "", "filter by pattern")
	flagAll := flag.Bool("a", false, "all files including hide files")
	flagNumberRecords := flag.Int("n", 0, "number of records")

	// order
	hasOrderByTime := flag.Bool("t", false, "sort by time, newest first")
	hasOrderBySize := flag.Bool("s", false, "sort by file size, largest first")
	hasOrderReverse := flag.Bool("r", false, "reverse order while sorting")

	// TODO explicar flag --help

	flag.Parse()

	path := flag.Arg(0)
	if path == "" {
		path = "."
	}

	dirs, err := os.ReadDir(path)
	if err != nil {
		panic(err)
	}

	fs := []file{}
	for _, dir := range dirs {
		info, err := dir.Info()
		if err != nil {
			panic(err)
		}

		if *flagPattern != "" {
			isMatched, err := regexp.MatchString(*flagPattern, dir.Name())
			if err != nil {
				panic(err)
			}

			if !isMatched {
				continue
			}
		}

		f := file{
			name:             dir.Name(),
			isDir:            dir.IsDir(),
			size:             info.Size(),
			mode:             info.Mode().String(),
			modificationTime: info.ModTime(),
		}

		setFileType(&f)
		setIsHidden(&f, path)

		if f.isHidden && !*flagAll {
			continue
		}

		stat, ok := info.Sys().(*syscall.Stat_t)
		if ok {
			if u, err := user.LookupId(fmt.Sprintf("%d", stat.Uid)); err == nil {
				f.userName = u.Username
			}

			if g, err := user.LookupGroupId(fmt.Sprintf("%d", stat.Gid)); err == nil {
				f.groupName = g.Name
			}
		}

		fs = append(fs, f)
	}

	if !*hasOrderBySize || !*hasOrderByTime {
		orderByName(fs, *hasOrderReverse)
	}

	if *hasOrderBySize && !*hasOrderByTime {
		orderBySize(fs, *hasOrderReverse)
	}

	if *hasOrderByTime {
		orderByTime(fs, *hasOrderReverse)
	}

	if *flagNumberRecords == 0 || *flagNumberRecords > len(fs) {
		*flagNumberRecords = len(fs)
	}

	printList(fs, *flagNumberRecords, path)
}

func mySort[T constraints.Ordered](i, j T, isReverse bool) bool {
	if isReverse {
		return i > j
	}

	return i < j
}

func orderByName(files []file, isReverse bool) {
	sort.SliceStable(files, func(i, j int) bool {
		return mySort(strings.ToLower(files[i].name), strings.ToLower(files[j].name), isReverse)
	})
}

func orderBySize(files []file, isReverse bool) {
	sort.SliceStable(files, func(i, j int) bool {
		return mySort(files[i].size, files[j].size, isReverse)
	})
}

func orderByTime(files []file, isReverse bool) {
	sort.SliceStable(files, func(i, j int) bool {
		return mySort(files[i].modificationTime.Unix(), files[j].modificationTime.Unix(), isReverse)
	})
}

// func orderBySize(files []file, isReverse bool) {
// 	sort.SliceStable(files, func(i, j int) bool {
// 		if isReverse {
// 			return files[i].size > files[j].size
// 		}

// 		return files[i].size < files[j].size
// 	})
// }

// func orderByTime(files []file, isReverse bool) {
// 	sort.SliceStable(files, func(i, j int) bool {
// 		if isReverse {
// 			return files[i].modificationTime.After(files[j].modificationTime)
// 		}

// 		return files[i].modificationTime.Before(files[j].modificationTime)
// 	})
// }

func printList(files []file, nRecords int, route string) {
	for _, file := range files[:nRecords] {
		style := mapStyleByFileType[file.fileType]

		// TODO logica para formato MB, KB etc

		fmt.Printf("%s %-10s %-10s %10d %s %s %s%s %s\n",
			file.mode, file.userName, file.groupName, file.size,
			file.modificationTime.Format(time.DateTime), style.icon,
			setColor(file.name, style.color), style.symbol, markHidden(file.isHidden))
	}
}

func setColor(nameFile string, styleColor color.Attribute) string {
	switch styleColor {
	case color.FgBlue:
		return blue(nameFile)
	case color.FgGreen:
		return green(nameFile)
	case color.FgRed:
		return red(nameFile)
	case color.FgMagenta:
		return magenta(nameFile)
	case color.FgCyan:
		return cyan(nameFile)
	}

	return nameFile
}

func setIsHidden(f *file, basePath string) {
	filePath := f.name

	if runtime.GOOS == windows {
		filePath = path.Join(basePath, filePath)
	}

	f.isHidden = fileinfo.IsHidden(filePath)
}

func setFileType(f *file) {
	switch {
	case isLink(*f):
		f.fileType = fileLink
	case f.isDir:
		f.fileType = fileDirectory
	case isExec(*f):
		f.fileType = fileExecutable
	case isCompress(*f):
		f.fileType = fileCompress
	case isImage(*f):
		f.fileType = fileImage
	default:
		f.fileType = fileRegular
	}
}

func isLink(f file) bool {
	return strings.HasPrefix(strings.ToLower(f.mode), "l")
}

func isExec(f file) bool {
	if runtime.GOOS == windows {
		return strings.HasSuffix(f.name, exe)
	}

	return strings.Contains(f.mode, "x")
}

func isCompress(f file) bool {
	return strings.HasSuffix(f.name, zip) || strings.HasSuffix(f.name, gz) ||
		strings.HasSuffix(f.name, tar) || strings.HasSuffix(f.name, rar) ||
		strings.HasSuffix(f.name, deb)
}

func isImage(f file) bool {
	return strings.HasSuffix(f.name, png) || strings.HasSuffix(f.name, jpg) ||
		strings.HasSuffix(f.name, gif)
}

func markHidden(isHidden bool) string {
	if !isHidden {
		return ""
	}

	return yellow("ø")
}
