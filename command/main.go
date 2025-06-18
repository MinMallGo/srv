package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
)

var mkName string

func main() {
	projectType := flag.String("type", "", "web or srv")
	name := flag.String("name", "", "micro name")
	flag.Parse()
	if len(*name) == 0 {
		panic("please input name")
	}

	mkName = *name

	switch *projectType {
	case "web":
		mkdir(pjMap("web"))
	case "srv":
		mkdir(pjMap("srv"))
	default:
		panic("unsupported project type")
	}
}

func mkdir(dirs []string) bool {
	curPath, _ := os.Getwd()
	curPath = path.Join(curPath, mkName)
	err := os.Mkdir(curPath, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return false
	}

	for _, dir := range dirs {
		yes := strings.HasSuffix(dir, ".go")
		if !yes {
			err := os.Mkdir(path.Join(curPath, dir), 0755)
			if err != nil {
				if errors.Is(err, os.ErrExist) {
					continue
				}
				log.Println("mkdir error:", err)
				return false
			}
			continue
		}
		err := os.WriteFile(path.Join(curPath, dir), []byte(mainTemplate()), 0755)
		fmt.Println(path.Join(curPath, dir))
		if err != nil {
			//if errors.Is(err, os.ErrExist) {
			//	continue
			//}
			log.Println("mk file with error:", err)
			return false
		}
	}
	return true
}

func mainTemplate() string {
	return `package main
func main(){
	// generate by script
}
`
}

// 1. 用于创建文件
func pjMap(pjType string) []string {
	web := []string{
		"api",
		"config",
		"forms",
		"global",
		"initialize",
		"middleware",
		"proto",
		"router",
		"utils",
		"validater",
		"main.go",
	}

	srv := []string{
		"global",
		"handler",
		"model",
		"proto",
		"test",
		"main.go",
	}

	switch pjType {
	case "web":
		return web
	case "srv":
		return srv
	default:
		return []string{}
	}
}
