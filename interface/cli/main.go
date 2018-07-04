package main

import (
	"os"

	"github.com/ld9999999999/go-interfacetools"
	"github.com/nahanni/go-ucl"
	"github.com/xiaokangwang/VisonFS/conf"
	"github.com/xiaokangwang/VisonFS/instanceadm"
)

func main() {
	var conffile conf.Configure
	cfg, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}

	p := ucl.NewParser(cfg)
	result, err := p.Ucl()
	if err != nil {
		panic(err)
	}
	interfacetools.CopyOut(result, conffile)
	var fi instanceadm.Instance
	fi.Prepare(conffile.Gitpath, conffile.Pubdir, conffile.Prvdir, conffile.Prvpass)
	insi := fi.Launch()
	switch os.Args[2] {
	case "mkdir":
		insi.Mkdir(os.Args[3], os.Args[4])
	case "rm":
		insi.Rm(os.Args[3], os.Args[4])
	case "push":
	case "pull":
	case "purge":
	}
}
