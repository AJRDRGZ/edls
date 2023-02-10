package main

import (
	"time"

	"github.com/fatih/color"
)

// os
const (
	windows = "windows"
)

// file type
const (
	fileRegular int = iota
	fileDirectory
	fileExecutable
	fileCompress
	fileImage
	fileLink
)

// file extensions
const (
	exe = ".exe"
	zip = ".zip"
	gz  = ".gz"
	tar = ".tar"
	rar = ".rar"
	deb = ".deb"
	png = ".png"
	jpg = ".jpg"
	gif = ".gif"
)

type file struct {
	name             string
	fileType         int
	isDir            bool
	isHidden         bool
	userName         string
	groupName        string
	size             int64
	modificationTime time.Time
	mode             string
}

type styleFileType struct {
	symbol string
	color  color.Attribute
	icon   string
}

var mapStyleByFileType = map[int]styleFileType{
	fileRegular:    {icon: "📄"},
	fileDirectory:  {icon: "📂", color: color.FgBlue, symbol: "/"},
	fileExecutable: {icon: "🚀", color: color.FgGreen, symbol: "*"},
	fileCompress:   {icon: "📦", color: color.FgRed},
	fileImage:      {icon: "🖼️", color: color.FgMagenta},
	fileLink:       {icon: "🔗", color: color.FgCyan},
}

// funciones de color
var (
	blue    = color.New(color.FgBlue).Add(color.Bold).SprintFunc()
	green   = color.New(color.FgGreen).Add(color.Bold).SprintFunc()
	red     = color.New(color.FgRed).Add(color.Bold).SprintFunc()
	magenta = color.New(color.FgMagenta).Add(color.Bold).SprintFunc()
	cyan    = color.New(color.FgCyan).Add(color.Bold).SprintFunc()
	yellow  = color.New(color.FgYellow).SprintFunc()
)
