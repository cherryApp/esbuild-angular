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
var buildOptions api.BuildOptions

func rebuild(start time.Time) api.BuildResult {
	result := api.Build(buildOptions)
	elapsed := time.Since(start)

	fmt.Printf("Project built in %s", elapsed)
	fmt.Println()

	if len(result.Errors) > 0 {
		fmt.Printf("%+v\n", result.Errors)
		os.Exit(1)
	}

  return result
}

func main() {
	start := time.Now()

	wd, _ := os.Getwd()
	workingDir = wd

	buildOptions = util.GetEsbuildOptions(workingDir)

	indexFilePath := path.Join(
		workingDir,
		util.GetProjectOption("architect.build.options.index").(string),
	)

	srcPath = path.Dir(buildOptions.EntryPoints[0])

	buildOptions.Plugins = []api.Plugin{
		plugin.GetAssetManager(
			workingDir,
			util.GetProjectOption("architect.build.options.assets").([]interface{}),
			buildOptions.Outdir,
		),
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

	// Build and serve
	if util.GetRuntimeOption("serve").(bool) {
    // buildOptions.Write = false;
    rebuild(start)
    util.FileWatcher(buildOptions, func(message string) {
			fmt.Println(message)
			start := time.Now()
			rebuild(start)
      util.RefreshLiveServerPage()
		})
		util.LiveServer(buildOptions)
    } else {
      rebuild(start)
    }

}
