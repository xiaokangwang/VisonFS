package main

import (
	"encoding/gob"
	"fmt"
	"os"
	"time"

	"github.com/ld9999999999/go-interfacetools"
	"github.com/nahanni/go-ucl"
	"github.com/xiaokangwang/VisonFS/conf"
	"github.com/xiaokangwang/VisonFS/filesystem"
	"github.com/xiaokangwang/VisonFS/instanceadm"
	"github.com/xiaokangwang/VisonFS/transfer"
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
	interfacetools.CopyOut(result, &conffile)
	fmt.Println(conffile)
	var fi instanceadm.Instance
	fi.Prepare(conffile.Gitpath, conffile.Pubdir, conffile.Prvdir, conffile.Prvpass, conffile.Tmpdir, conffile.UploadPrefix)
	insi := fi.Launch()
	switch os.Args[2] {
	case "mkdir":
		insi.Mkdir(os.Args[3])
	case "rm":
		insi.Rm(os.Args[3])
	case "push":
		task := transfer.NewTask(os.Args[3], os.Args[4], true, insi)
		ProgressTask(task, fi)
	case "pull":
		task := transfer.NewTask(os.Args[3], os.Args[4], false, insi)
		ProgressTask(task, fi)
	case "resume":
		f, err := os.Open("resume")
		if err != nil {
			panic(err)
		}
		o := gob.NewDecoder(f)
		var Task transfer.Transfer
		o.Decode(&Task)
		Task.PushFileInstance(insi)
		ProgressTask(&Task, fi)
	case "purge":
		fi.Purge()
	case "mount":
		filesystem.Mount(insi, os.Args[3])
		time.Sleep(1 * time.Second)
	case "ls":
		dir, err := insi.Ls(os.Args[3])
		if err != nil {
			panic(err)
		}
		for _, v := range dir {
			if v == nil {
				continue
			}
			fmt.Printf("\n\n%v Dir: %v\n\n", v.Name(), v.IsDir())
		}
	}
}
func ProgressTask(task *transfer.Transfer, fi instanceadm.Instance) {
	for task.HasNext() {
		task.ProcessBlock()
		f, err := os.Create("resume")
		if err != nil {
			panic(err)
		}
		e := gob.NewEncoder(f)
		e.Encode(task)
		f.Close()
		fi.Purge()
	}
	os.Remove("resume")
}
