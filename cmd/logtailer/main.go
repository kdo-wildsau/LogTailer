package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"

	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/kdo-wildsau/logTailer/pkg/logtailer"
	"github.com/logdna/logdna-go/logger"
	"github.com/nxadm/tail"
)

func main() {

	var (
		configFile = flag.String("c", "config.json", "path to logfile. Could be file, folder, or has a wildcard *")
	)
	flag.Parse()

	var config logtailer.Config
	filesNotification := make(chan string)
	printChannel := make(chan string)

	file, err := os.ReadFile(*configFile)
	if err != nil {
		panic("Configfile could be found")
	}

	if err = json.Unmarshal(file, &config); err != nil {
		panic(err.Error())
	}

	if config.MemzoIngestionKey == "" {
		panic("No MemzoIngestionKey could be found in the config file")
	}

	if config.MemzoInstanceName == "" {
		config.MemzoInstanceName = "NO NAME GIVEN"
		fmt.Println("Warning: MemzoInstanceName is not set")
	}

	hostname, _ := os.Hostname()
	options := logger.Options{
		App:           config.MemzoInstanceName,
		FlushInterval: 10 * time.Second,
		SendTimeout:   0,
		Hostname:      hostname,
		IndexMeta:     false,
		Level:         "fatal",
	}

	memzoClient, err := logger.NewLogger(options, config.MemzoIngestionKey)
	if err != nil {
		panic(err)
	}
	defer memzoClient.Close()

	go fileWatcher(config, filesNotification)

	for {
		select {
		case file := <-filesNotification:
			fmt.Println("New File ", file)

			t, err := tail.TailFile(file, tail.Config{
				Follow:   true,
				ReOpen:   false,
				Poll:     true,
				Location: &tail.SeekInfo{Offset: 0, Whence: io.SeekEnd}, // <- line changed
			})
			if err != nil {
				fmt.Println(err)
			}

			go func(t *tail.Tail) {

				for l := range t.Lines {
					if l.Text != "" {
						printChannel <- l.Text
					}

				}

				if err := t.Err(); err != nil {
					fmt.Println(err)
				}
			}(t)

		case line := <-printChannel:
			memzoClient.Debug(line)
		}
	}
}

func fileWatcher(config logtailer.Config, newFile chan string) {
	t := time.Duration(config.LogPathUpdateTime) * time.Second

	var fileList []string

	for {
		files, _ := listFiles(config.LogPath, true)
		for _, f := range files {
			if !contains(fileList, f) {
				fileList = append(fileList, f)
				newFile <- f
			}
		}
		time.Sleep(t)
	}
}

func evaluatePathPattern(pattern string) (path string, prefix string, suffix string, wildcard bool) {
	path = filepath.Dir(pattern)
	file := filepath.Base(pattern)
	wildcard = false
	usingWildcard := strings.Contains(file, "*")

	if usingWildcard {
		s := strings.SplitN(file, "*", 2)
		prefix = s[0]

		if len(s) > 1 {
			suffix = s[1]
			wildcard = strings.Contains(suffix, "*")
			suffix = strings.ReplaceAll(suffix, "*", "")
		}
	}

	return path, prefix, suffix, wildcard
}

func listFiles(pattern string, ignoreDotFiles bool) ([]string, error) {
	var files []string
	startpath, prefix, suffix, wildcard := evaluatePathPattern(pattern)

	err := filepath.Walk(startpath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		filename := filepath.Base(path)
		dirName := filepath.Dir(path)
		appendFile := true
		_, postPath, _ := strings.Cut(dirName, startpath)

		// no Dirs will be add
		appendFile = appendFile && !info.IsDir()

		// check file prefix, but add only if path is exact equal
		if prefix != "" {
			appendFile = appendFile && postPath == ""
			appendFile = appendFile && strings.HasPrefix(filename, prefix)
		}

		if suffix != "" {
			if !wildcard {
				appendFile = appendFile && postPath == ""
			}
			appendFile = appendFile && strings.HasSuffix(path, suffix)
		}

		if ignoreDotFiles {
			appendFile = appendFile && !strings.HasPrefix(filename, ".")
			s := strings.Split(dirName, "/")
			for _, dirPart := range s {
				appendFile = appendFile && !strings.HasPrefix(dirPart, ".")
			}

		}

		if appendFile {
			files = append(files, path)
		}

		return nil
	})

	return files, err
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
