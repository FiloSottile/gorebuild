package main

import (
	"flag"
	"fmt"
	"go/build"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	dry := flag.Bool("n", false, "don't build, just print the package names")
	flag.Parse()

	bins := flag.Args()
	if len(bins) == 0 {
		fi, err := ioutil.ReadDir(build.Default.GOPATH + "/bin")
		if err != nil {
			log.Fatal(err)
		}
		for _, f := range fi {
			if f.IsDir() {
				continue
			}
			bins = append(bins, build.Default.GOPATH+"/bin/"+f.Name())
		}
	}

	var binDir string
	if !*dry {
		dir, err := ioutil.TempDir("", "gorebuild")
		if err != nil {
			log.Fatal(err)
		}
		defer os.RemoveAll(dir)
		binDir = dir
	}

	for _, file := range bins {
		path, err := getMainPath(file)
		if err != nil {
			log.Fatal(err)
		}
		importPath := stripPath(path)
		if *dry {
			fmt.Println(stripPath(path))
		} else {
			cmd := exec.Command("go", "install", "-v", importPath)
			cmd.Env = append(os.Environ(), "GOBIN="+binDir)
			cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
			cmd.Run()
			fi, err := ioutil.ReadDir(binDir)
			if err != nil {
				log.Fatal(err)
			}
			for _, f := range fi {
				err := os.Rename(filepath.Join(binDir, f.Name()), build.Default.GOPATH+"/bin/"+f.Name())
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}
}

func stripPath(path string) string {
	dir := filepath.Dir(path)
	return strings.TrimPrefix(dir, build.Default.GOPATH+"/src/")
}
