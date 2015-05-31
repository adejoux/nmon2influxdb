package main

import (
	"github.com/codegangsta/cli"
	"strings"
)

//
//helper functions
//
func check(e error) {
	if e != nil {
		panic(e)
	}
}

func ReplaceComma(s string) string {
	return "<tr><td>" + strings.Replace(s, ",", "</td><td>", 1) + "</td></tr>"
}

type Params struct {
	Filepath string
	NoDisks  bool
	CpuAll   bool
	Host     string
	User     string
	Port     string
	Db       string
	Password string
	Template string
	Debug    bool
}

func ParseParameters(c *cli.Context) (params *Params) {
	return &Params{Filepath: c.Args()[0],
		NoDisks:  c.Bool("nodisks"),
		CpuAll:   c.Bool("cpus"),
		Debug:    c.GlobalBool("debug"),
		Host:     c.GlobalString("host"),
		User:     c.GlobalString("user"),
		Port:     c.GlobalString("port"),
		Db:       c.GlobalString("db"),
		Password: c.GlobalString("pass"),
		Template: c.String("template"),
	}
}
