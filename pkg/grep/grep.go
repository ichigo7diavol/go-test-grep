package grep

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

var (
	ErrFileInfoIsNil = errors.New("internal error. file info is nil")
)

type Config struct {
	Expression string
	Target     string

	IsRecursive bool
	IsVerbose   bool
}

type extendedFileInfo struct {
	FileInfo os.FileInfo
	Path     string
}

type grepWorkerJobData struct {
	FileInfo extendedFileInfo
}

type grepWorkerResultData struct {
	ResultInfo resultInfo
	Error      error
}

type resultInfo struct {
}

func Execute(config Config) error {
	info, err := os.Stat(config.Target)
	if err != nil {
		return err
	}
	extInfo := extendedFileInfo{
		FileInfo: info,
		Path:     config.Target,
	}
	if info.IsDir() {
		wg := sync.WaitGroup{}

		results := make(chan grepWorkerResultData)
		jobs := make(chan grepWorkerJobData)

		for i := 0; i < runtime.NumCPU(); i++ {
			wg.Add(1)
			go processFileWorker(config, jobs, results, &wg)
		}
		go func() {
			defer close(results)
			wg.Wait()
		}()
		go func() {
			defer close(jobs)
			processDirectoryWithFunction(
				config.Target,
				config.IsRecursive,
				func(path string, entry os.DirEntry) error {
					if entry.IsDir() {
						return nil
					}
					jobs <- grepWorkerJobData{
						FileInfo: extendedFileInfo{
							FileInfo: info,
							Path:     path,
						},
					}
					return nil
				},
			)
		}()

		for result := range results {
			fmt.Println(result)
		}
	} else {
		result, err := processFile(extInfo, config)
		if err != nil {
			return err
		}
		fmt.Println(result)
	}
	return nil
}

func processFile(info extendedFileInfo, config Config) (resultInfo, error) {
	return resultInfo{}, nil
}

func processDirectoryWithFunction(path string, isRecursive bool, fn func(path string, entry os.DirEntry) error) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", path)
	}
	if isRecursive {
		return filepath.WalkDir(path, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			return fn(path, d)
		})
	} else {
		entries, err := os.ReadDir(path)
		if err != nil {
			return err
		}
		for _, e := range entries {
			fullPath := path + string(os.PathSeparator) + e.Name()
			err := fn(fullPath, e)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func processFileWorker(
	config Config,
	jobs <-chan grepWorkerJobData,
	results chan<- grepWorkerResultData,
	wg *sync.WaitGroup,
) {
	defer wg.Done()
	for j := range jobs {
		result, err := processFile(j.FileInfo, config)
		results <- grepWorkerResultData{
			ResultInfo: result,
			Error:      err,
		}
	}
}
