package util

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/evanw/esbuild/pkg/api"
)

type AngularOptions struct {
	serve bool
	port  int
}

func Check(e error) {
	if e != nil {
		panic(e)
	}
}

func StatPath(filePath string) bool {
	if _, err := os.Stat(filePath); err == nil {
		return true
	}
	return false
}

func GetEsbuildOptions(srcPath string, outPath string) (api.BuildOptions, AngularOptions) {

	bundle := flag.Bool("bundle", true, "bundle the result")
	splitting := flag.Bool("splitting", true, "splitting the result")
	write := flag.Bool("write", true, "write the result")
	minify := flag.Bool("minify", false, "minify the result")

	serve := flag.Bool("serve", false, "start the devserver")
	port := flag.Int("port", 4200, "devserver port")

	// var svar string
	// flag.StringVar(&svar, "svar", "bar", "a string var")

	flag.Parse()

	var buildOptions = api.BuildOptions{
		EntryPoints: []string{"D:/Projects/esbuild-angular/src/main.ts"},
		Format:      api.FormatESModule,
		Bundle:      *bundle,
		Outdir:      outPath,
		Platform:    api.PlatformBrowser,
		Splitting:   *splitting,
		Target:      api.Target(8),
		Write:       *write,
		TreeShaking: api.TreeShakingTrue,
		Loader: map[string]api.Loader{
			".html": api.LoaderText,
			".css":  api.LoaderText,
		},
		Sourcemap:    api.SourceMapExternal,
		MinifySyntax: *minify,
	}

	var angularOptions = AngularOptions{
		serve: *serve,
		port:  *port,
	}

	return buildOptions, angularOptions

}

func GetAngularOptions(srcPath string, outPath string) AngularOptions {

	serve := flag.Bool("serve", false, "start the devserver")
	port := flag.Int("port", 4200, "devserver port")

	// var svar string
	// flag.StringVar(&svar, "svar", "bar", "a string var")

	flag.Parse()

	return AngularOptions{
		serve: *serve,
		port:  *port,
	}
}

// PrintMemUsage outputs the current, total and OS memory being used. As well as the number
// of garage collection cycles completed.
func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Println()
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
