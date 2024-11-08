package priest

import (
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

func doAction(
	dirs []string,
	packageName, functionName, outputFile string,
	showStat, isWatch bool,
) error {
	if len(dirs) == 0 {
		wd, _ := os.Getwd()
		dirs = append(dirs, wd)
	}

	for i := range dirs {
		dirs[i], _ = filepath.Abs(dirs[i])
	}

	if !filepath.IsAbs(outputFile) {
		outputFile, _ = filepath.Abs(outputFile)
	}

	loader := autoload{
		scanDir:      dirs,
		packageName:  packageName,
		functionName: functionName,
		outputFile:   outputFile,
	}
	err := loader.fillModuleInfo()
	if err != nil {
		log.Fatalf("loader.fillModuleInfo() err:%v", err)
		return err
	}
	err = loader.firstGenerate()
	if err != nil {
		log.Fatalf("loader.firstGenerate() err:%v", err)
		return err
	}

	if isWatch {
		log.Println("watch mode...")
		doWatch(loader.reGenerate, dirs, outputFile)
	}
	return nil
}
