package main

import (
	"embed"
	"fmt"
	"net/http"
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
var BuildFs embed.FS
var indexFile *os.File

// func serveHome(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodGet {
// 		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 		return
// 	}

// 	indexPath := path.Join(buildOptions.Outdir, "index.html")

// 	http.ServeFile(w, r, indexPath)
// }

// func BuildHTTPFS() http.FileSystem {
// 	build, err := fs.Sub(BuildFs, buildOptions.Outdir)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return http.FS(build)
// }

// func handleSPA(w http.ResponseWriter, r *http.Request) {
// 	http.FileServer(BuildHTTPFS()).ServeHTTP(w, r)
// }

func serveSPS(fs http.FileSystem) http.Handler {
	fileServer := http.FileServer(fs)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := fs.Open(path.Clean(r.URL.Path)) // Do not allow path traversals.
		if err != nil {
			indexPath := path.Join(buildOptions.Outdir, "index.html")
			if indexFile == nil {
				indexFile, _ = os.Open(indexPath)
			}
			http.ServeContent(w, r, indexPath, time.Time{}, indexFile)
			return
		}
		fileServer.ServeHTTP(w, r)
	})
}

func rebuild(start time.Time) {
	result := api.Build(buildOptions)
	elapsed := time.Since(start)

	fmt.Printf("Project built in %s", elapsed)
	fmt.Println()

	if len(result.Errors) > 0 {
		fmt.Printf("%+v\n", result.Errors)
		os.Exit(1)
	}
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
	rebuild(start)
	if util.GetRuntimeOption("serve").(bool) {
		var addr = "127.0.0.1:" + fmt.Sprintf("%v", util.GetRuntimeOption("port"))
		http.ListenAndServe(addr, serveSPS(http.Dir(buildOptions.Outdir)))
	}

}
