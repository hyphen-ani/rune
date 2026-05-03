package cli

import (
	"github.com/fatih/color"
)

var (
	success = color.New(color.FgGreen).SprintFunc()
	error   = color.New(color.FgRed).SprintFunc()
	info    = color.New(color.FgCyan).SprintFunc()
	warn    = color.New(color.FgYellow).SprintFunc()
)

func Success(msg string) {
	println(success("✔ " + msg))
}

func Error(msg string) {
	println(error("✖ " + msg))
}

func Info(msg string) {
	println(info("➜ " + msg))
}

func Warn(msg string) {
	println(warn("! " + msg))
}
