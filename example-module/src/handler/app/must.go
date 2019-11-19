package main

import (
	"io"
	"io/ioutil"
	"log"
	"os"
)

// creates the directory at the given path with given permissions
func mustMkdirAll(path string, perm os.FileMode) {
	err := os.MkdirAll(path, perm)
	if err != nil {
		log.Panicf("os.MkdirAll('%s', %o) error: %#v\n", path, perm, err)
	}
}

// creates a symbolic link, newname to oldname
// for e.g. error.log -> stderr, access.log -> stdout
func mustSymlink(oldname, newname string) {
	err := os.Symlink(oldname, newname)
	if err != nil {
		log.Panicf("os.Symlink('%s', '%s') error: %#v\n",
			oldname,
			newname,
			err)
	}
}

// writes to the given file and sets the given permissions
func mustWriteFile(filename string, data []byte, perm os.FileMode) {
	err := ioutil.WriteFile(filename, data, perm)
	if err != nil {
		log.Panicf("ioutil.WriteFile('%s', len(data):%d, %o) error: %#v\n",
			filename,
			len(data),
			perm,
			err)

	}
}

// close the file
func mustClose(c io.Closer) {
	err := c.Close()
	if err != nil {
		log.Panic(err)
	}
}
