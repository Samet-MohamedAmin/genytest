package utilities

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"os/exec"
)

type Set map[string]any

func Hash(input string) string {
	h := sha1.New()
	h.Write([]byte(input))
	return hex.EncodeToString(h.Sum(nil))
}

func CleanFile(fo *os.File) {
	if err := fo.Truncate(0); err != nil {
		fmt.Println(err)
	}
	fo.Seek(0, 0)
}

func DeleteFile(filename string) error {
	err := os.Remove(filename)
	if err != nil && !os.IsNotExist(err) {
		// An error occurred other than "file does not exist"
		return err
	}
	return nil
}

func CleanDir(dir string) {
	err := os.RemoveAll(dir)
	if err != nil {
		fmt.Println("Error deleting directory:", err)
		return
	}

	fmt.Println("Directory deleted successfully")
}

func RunCommand(c *exec.Cmd, skipError bool) string {
	// if len(customPath) != 0 {
	// 	cmd.Dir = customPath
	// }

	var stdBuffer bytes.Buffer
	mw := io.MultiWriter(&stdBuffer) //, os.Stdout)

	c.Stdout = mw
	c.Stderr = mw

	if err := c.Run(); !skipError && err != nil {
		fmt.Println(stdBuffer.String())
		fmt.Println(err.Error())
		panic(err)
	}

	return stdBuffer.String()
}

func (s Set) Add(item string) {
	s[item] = nil
}

func (s Set) GetItems() (items []string) {
	for key := range s {
		items = append(items, key)
	}
	return
}
