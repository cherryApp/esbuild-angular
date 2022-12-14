package util

import (
	"flag"
	"fmt"
	"os"
	"path"
	"reflect"
	"runtime"

	"github.com/evanw/esbuild/pkg/api"

	gojsonq "github.com/thedevsaddam/gojsonq/v2"
)

var AngularOptions *gojsonq.JSONQ

var RuntimeOptions = make(map[string]interface{})

var ProjectName string

func CopyFile(sourceFile string, destinationFile string) {
	input, err := os.ReadFile(sourceFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = os.WriteFile(destinationFile, input, 0644)
	if err != nil {
		fmt.Println("Error creating", destinationFile)
		fmt.Println(err)
		return
	}
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

func GetInterfaceKeys(currentInterface interface{}) []string {
	keyList := map[int]string{}
	iter := reflect.ValueOf(currentInterface).MapRange()
	i := 0
	for iter.Next() {
		key := iter.Key().Interface()
		// value := iter.Value().Interface()
		keyList[i] = key.(string)
		i += 1
	}

	keyArray := make([]string, 0, len(keyList))
	for _, value := range keyList {
		keyArray = append(keyArray, value)
	}

	return keyArray
}

func ArrayContains(items []string, element string) bool {
	for _, x := range items {
		if x == element {
			return true
		}
	}
	return false
}

func GetProjectOption(key string) interface{} {
	return AngularOptions.Copy().Find(ProjectName + "." + key)
}

func GetRuntimeOption(key string) interface{} {
	return RuntimeOptions[key]
}

func CheckBuildPath(outputDir string) {
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		fmt.Println("Creating output directory:", outputDir)
		err := os.MkdirAll(outputDir, os.ModePerm)
		if err != nil {
			fmt.Println("ERROR creating output directory:", outputDir)
			os.Exit(1)
		}
	}
}

func GetEsbuildOptions(workingDir string) api.BuildOptions {

	// Parse angular.json
	AngularOptions = gojsonq.New().File(path.Join(workingDir, "angular.json"))
	projectNames := GetInterfaceKeys(AngularOptions.Copy().Find("projects"))

	// Set flags.
	bundle := flag.Bool("bundle", true, "bundle the result")
	splitting := flag.Bool("splitting", true, "splitting the result")
	write := flag.Bool("write", true, "write the result")
	minify := flag.Bool("minify", true, "minify the result")

	project := flag.String("project", projectNames[0], "project name")

	// Runtime options
	serve := flag.Bool("serve", false, "start the devserver")
	port := flag.Int("port", 4200, "devserver port")
	baseHref := flag.String("base-href", "/", "Base url for the application being built.")

	flag.Parse()

	// Set paths.
	ProjectName = "projects." + *project
	main := GetProjectOption("architect.build.options.main")
	outputPath := GetProjectOption("architect.build.options.outputPath")

	RuntimeOptions["port"] = *port
	RuntimeOptions["serve"] = *serve
	RuntimeOptions["base-href"] = *baseHref

	CheckBuildPath(outputPath.(string))

	var buildOptions = api.BuildOptions{
		EntryPoints: []string{path.Join(workingDir, main.(string))},
		Format:      api.FormatESModule,
		Bundle:      *bundle,
		Outdir:      outputPath.(string),
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
		MangleProps:  "_$",
		MinifySyntax: *minify,
	}

	return buildOptions

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
