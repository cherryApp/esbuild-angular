package main

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/evanw/esbuild/pkg/api"

	"cherryapp/angular/pkg/plugin"
	"cherryapp/angular/pkg/util"
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

	buildOptions := api.BuildOptions{
		EntryPoints: []string{path.Join(srcPath, "main.ts")},
		Format:      api.FormatESModule,
		Bundle:      true,
		Outdir:      outPath,
		Platform:    api.PlatformBrowser,
		Splitting:   true,
		Target:      api.Target(8),
		Write:       true,
		TreeShaking: api.TreeShakingTrue,
		Loader: map[string]api.Loader{
			".html": api.LoaderText,
			".css":  api.LoaderText,
		},
		Sourcemap:    api.SourceMapExternal,
		MinifySyntax: true,
		Plugins: []api.Plugin{
			plugin.GetIndexFileProcessor(outPath, srcPath),
			plugin.GetMainManager(),
			plugin.AngularComponentDecoratorPlugin,
		},
		AbsWorkingDir: workingDir,
	}

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
