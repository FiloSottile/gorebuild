package main

import (
	"flag"
	"fmt"
	"go/build"
	"io"
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
			log.Printf("Skipping %s: %s", file, err)
			continue
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
				err := moveFile(filepath.Join(binDir, f.Name()), build.Default.GOPATH+"/bin/"+f.Name())
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

// moveFile safely moves files across different file systems.
func moveFile(src, dest string) error {
	if err := os.Rename(src, dest); err != nil {
		if _, ok := err.(*os.LinkError); ok {
			// Looks like we are trying to move files across system volumes. Try to
			// copy/move/delete strategy.

			// Open source file for reading.
			srcFile, err := os.Open(src)
			if err != nil {
				return err
			}
			defer srcFile.Close()
			// Get source file permissions.
			info, err := srcFile.Stat()
			if err != nil {
				return err
			}
			// Create a temporary file with unique name in the destination's directory.
			tmpFile, err := ioutil.TempFile(filepath.Dir(dest), filepath.Base(dest))
			if err != nil {
				return err
			}
			tmpPath := tmpFile.Name()
			// Always remove temporary file.
			defer os.Remove(tmpPath)
			// Copy bytes.
			if _, err = io.Copy(tmpFile, srcFile); err != nil {
				defer tmpFile.Close() // prevents file descriptor leak.
				return err
			}
			// Close tmpFile to flush all writes.
			if err = tmpFile.Close(); err != nil {
				return err
			}
			// Atomically rename. This is safe now since both files are in the same directory.
			if err = os.Rename(tmpPath, dest); err != nil {
				return err
			}
			// Change permission bits on destination to match the source file.
			if err = os.Chmod(dest, info.Mode().Perm()); err != nil {
				return err
			}
			// Finally, if all is good, delete source file.
			return os.Remove(src)
		}
	}
	return nil
}
