package main

import (
	"embed"
	"fmt"
	"io/fs"
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

// func serveHome(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodGet {
// 		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 		return
// 	}

// 	indexPath := path.Join(buildOptions.Outdir, "index.html")

// 	http.ServeFile(w, r, indexPath)
// }

func BuildHTTPFS() http.FileSystem {
	build, err := fs.Sub(BuildFs, buildOptions.Outdir)
	if err != nil {
		panic(err)
	}
	return http.FS(build)
}

func handleSPA(w http.ResponseWriter, r *http.Request) {
	http.FileServer(BuildHTTPFS()).ServeHTTP(w, r)
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

	// Build or serve
	if util.GetRuntimeOption("serve").(bool) {
		var addr = "127.0.0.1:" + fmt.Sprintf("%v", util.GetRuntimeOption("port"))
		http.HandleFunc("/", handleSPA)
		http.ListenAndServe(addr, nil)
		// fs := http.FileServer(http.Dir(buildOptions.Outdir))
		// http.Handle("/", fs)
		// err := http.ListenAndServe(addr, nil)
		// if err != nil {
		// 	panic(err)
		// }

		// http.HandleFunc("/", serveHome)
		// http.HandleFunc("/ws", serveWs)
		// http.ListenAndServe(addr, nil)

		// buildOptions.Watch = &api.WatchMode{}
		// server, err := api.Serve(api.ServeOptions{
		// 	Servedir: buildOptions.Outdir,
		// 	Port:     uint16(util.GetRuntimeOption("port").(int)),
		// },
		// 	api.BuildOptions{},
		// )

		// if err != nil {
		// 	panic(err)
		// }

		// result := api.Build(buildOptions)
		// elapsed := time.Since(start)

		// fmt.Printf("Project built in %s", elapsed)
		// fmt.Println()

		// if len(result.Errors) > 0 {
		// 	fmt.Printf("%+v\n", result.Errors)
		// 	os.Exit(1)
		// }

		// fmt.Printf("Server running in: http://localhost:%d", server.Port)
		// fmt.Println()
		// server.Wait()
	} else {
		result := api.Build(buildOptions)
		elapsed := time.Since(start)

		fmt.Printf("Project built in %s", elapsed)
		fmt.Println()

		if len(result.Errors) > 0 {
			fmt.Printf("%+v\n", result.Errors)
			os.Exit(1)
		}
	}

}
