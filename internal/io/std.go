package io

import (
	"fmt"
	"io"
	"os"
)

type Standard interface {
	Args() []string
	Exit(code int)
	Stdout() io.Writer
	Stderr() io.Writer
	Stdin() io.Reader
	Print(a ...any) (n int, err error)
	Printf(format string, a ...any) (n int, err error)
	Println(a ...any) (n int, err error)
}

type RealStandard struct{}

func (self *RealStandard) Args() []string {
	return os.Args
}

func (self *RealStandard) Exit(code int) {
	os.Exit(code)
}

func (self *RealStandard) Stdout() io.Writer {
	return os.Stdout
}

func (self *RealStandard) Stderr() io.Writer {
	return os.Stderr
}

func (self *RealStandard) Stdin() io.Reader {
	return os.Stdin
}

func (self *RealStandard) Print(a ...any) (n int, err error) {
	return fmt.Print(a...)
}

func (self *RealStandard) Printf(format string, a ...any) (n int, err error) {
	return fmt.Printf(format, a...)
}

func (self *RealStandard) Println(a ...any) (n int, err error) {
	return fmt.Println(a...)
}
