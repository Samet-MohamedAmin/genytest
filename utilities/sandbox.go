package utilities

import (
	"genytest/config"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
)

func getSrcFiles(path string) []string {
	srcFiles, _ := filepath.Glob(path + "/*")
	excludeRegex := regexp.MustCompile(config.ExcludePattern)
	filteredFiles := []string{}
	for _, file := range srcFiles {
		if !excludeRegex.MatchString(file) {
			filteredFiles = append(filteredFiles, file)
		}
	}

	return filteredFiles
}

func CreateSandboxedEnv(h string) (dir string) {
	dir = path.Join(config.BaseTmpDir, h)

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		panic(err)
	}

	srcFiles := getSrcFiles(config.BasePath)
	args := append([]string{"cp", "-r"}, append(srcFiles, dir)...)
	c := exec.Command(args[0], args[1:]...)
	c.Dir = config.BasePath
	RunCommand(c, false)

	return
}
