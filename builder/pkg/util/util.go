package util

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"reflect"
	"runtime"

	"github.com/evanw/esbuild/pkg/api"

	gojsonq "github.com/thedevsaddam/gojsonq/v2"
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

func GetNestedProp(currentInterface interface{}, keyArray []string) {

	for i := 0; i < len(keyArray); i++ {
		keys := GetInterfaceKeys(currentInterface)
		fmt.Println(keys)
	}

	// projects := GetInterfaceKeys(angularJson["projects"])
	// project := angularJson["projects"].(map[string]interface{})[projects[0]]
	// architect := project.(map[string]interface{})["architect"]
	// build := architect.(map[string]interface{})["build"]
	// options := build.(map[string]interface{})["options"]
	// main := options.(map[string]interface{})["main"].(string)
}

func ReadJsonFile(filePath string) *gojsonq.JSONQ {
	jsonFile, err := os.Open(filePath)
	Check(err)

	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	jsonObject := gojsonq.New().FromString(string(byteValue))

	fmt.Println(jsonObject.FromInterface("projects").First())

	return jsonObject
}

func GetEsbuildOptions(workingDir string) (api.BuildOptions, AngularOptions) {

	outPath := path.Join(workingDir, "dist", "project")
	// srcPath := path.Join(workingDir, "src")

	bundle := flag.Bool("bundle", true, "bundle the result")
	splitting := flag.Bool("splitting", true, "splitting the result")
	write := flag.Bool("write", true, "write the result")
	minify := flag.Bool("minify", false, "minify the result")

	serve := flag.Bool("serve", false, "start the devserver")
	port := flag.Int("port", 4200, "devserver port")

	flag.Parse()

	// var angularJson = ReadJsonFile(path.Join(workingDir, "angular.json"))
	// projects := GetInterfaceKeys(angularJson["projects"])
	// project := angularJson["projects"].(map[string]interface{})[projects[0]]
	// architect := project.(map[string]interface{})["architect"]
	// build := architect.(map[string]interface{})["build"]
	// options := build.(map[string]interface{})["options"]
	// main := options.(map[string]interface{})["main"].(string)

	// GetNestedProp(angularJson, []string{"projects", projects[0]})

	angularJson := ReadJsonFile(path.Join(workingDir, "angular.json"))

	project := angularJson.First()
	fmt.Println(project)
	// main := project.Find("architect.build.options.main")

	var buildOptions = api.BuildOptions{
		// EntryPoints: []string{path.Join(workingDir, main.(string))},
		EntryPoints: []string{path.Join(workingDir, "src", "main.ts")},
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
