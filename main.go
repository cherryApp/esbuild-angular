package main

import (
	"fmt"
	"os"
	"path"
	"time"

	"cherryapp/angular/pkg/plugin"
	"cherryapp/angular/pkg/util"

	"github.com/evanw/esbuild/pkg/api"
)

// Global variables
var workingDir string
var srcPath string
var outPath string

func main() {
	start := time.Now()

	wd, _ := os.Getwd()
	workingDir = wd
	srcPath = path.Join(workingDir, "src")
	outPath = path.Join(workingDir, "dist", "project")

	buildOptions, _ := util.GetEsbuildOptions(srcPath, outPath)

	buildOptions.Plugins = []api.Plugin{
		plugin.GetIndexFileProcessor(srcPath, outPath),
		plugin.GetMainManager(),
		plugin.AngularComponentDecoratorPlugin,
	}

	buildOptions.AbsWorkingDir = workingDir

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
