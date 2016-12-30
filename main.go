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
)

var (
	goPathBin = filepath.Join(build.Default.GOPATH, "bin")
	goPathSrc = filepath.Join(build.Default.GOPATH, "src")
)

func main() {
	dry := flag.Bool("n", false, "don't build, just print the package names")
	flag.Parse()

	bins := flag.Args()
	if len(bins) == 0 {
		fi, err := ioutil.ReadDir(goPathBin)
		if err != nil {
			log.Fatal(err)
		}
		for _, f := range fi {
			if f.IsDir() {
				continue
			}
			bins = append(bins, filepath.Join(goPathBin, f.Name()))
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
			log.Printf("Skipping %s: %s", file, err)
			continue
		}
		importPath, err := filepath.Rel(goPathSrc, filepath.Dir(path))
		if err != nil {
			log.Fatal(err)
		}
		if *dry {
			fmt.Println(importPath)
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
				err := os.Rename(
					filepath.Join(binDir, f.Name()),
					filepath.Join(goPathBin, f.Name()))
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}
}
