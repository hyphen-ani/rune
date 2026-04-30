package cli

import (
	"github.com/fatih/color"
)

var (
	success = color.New(color.FgGreen).SprintFunc()
	errc    = color.New(color.FgRed).SprintFunc()
	info    = color.New(color.FgCyan).SprintFunc()
)

func Success(msg string) {
	println(success("✔ " + msg))
}

func Error(msg string) {
	println(errc("✖ " + msg))
}

func Info(msg string) {
	println(info("➜ " + msg))
}
