package main

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/evanw/esbuild/pkg/api"

	"cherryApp/esbuild-angular/pkg/plugin"
	"cherryApp/esbuild-angular/pkg/util"
)

// Global variables
var workingDir string
var srcPath string

func main() {
	start := time.Now()

	wd, _ := os.Getwd()
	workingDir = wd

	buildOptions := util.GetEsbuildOptions(workingDir)

  indexFilePath := path.Join(
    workingDir,
    util.GetProjectOption("architect.build.options.index").(string),
  )

  srcPath = path.Dir( buildOptions.EntryPoints[0] )

	buildOptions.Plugins = []api.Plugin{
		plugin.GetIndexFileProcessor(indexFilePath, buildOptions.Outdir),
		plugin.GetMainManager(),
		plugin.AngularComponentDecoratorPlugin,
	}

	buildOptions.AbsWorkingDir = workingDir
  buildOptions.EntryPoints = []string{path.Join(srcPath, "main.ts")}

	tsConfigPath := path.Join(workingDir, "tsconfig.json")
	if util.StatPath(tsConfigPath) {
		buildOptions.Tsconfig = tsConfigPath
	}

	result := api.Build(buildOptions)

	elapsed := time.Since(start)

	fmt.Println()
	fmt.Printf("Project built in %s", elapsed)

	if len(result.Errors) > 0 {
		fmt.Printf("%+v\n", result.Errors)
		os.Exit(1)
	}
}
