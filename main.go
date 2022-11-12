package main

import (
	"fmt"
	"os"
	"path"
	"regexp"
	"time"

	"github.com/evanw/esbuild/pkg/api"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// func joinPath(items ...string) string {
// 	sep := filepath.Separator
// 	return strings.Join(items, string(sep))
// }

func statPath(filePath string) bool {
	if _, err := os.Stat(filePath); err == nil {
		return true
	}
	return false
}

// Global variables
var workingDir string
var indexFileContent string
var srcPath string
var outPath string

var indexFileProcessor = api.Plugin{
	Name: "indexProcessor",
	Setup: func(build api.PluginBuild) {

		build.OnStart(func() (api.OnStartResult, error) {
			indexContent, err := os.ReadFile(path.Join(srcPath, "index.html"))
			check(err)
			indexFileContent = string(indexContent)
			return api.OnStartResult{}, nil
		})

		build.OnEnd(func(result *api.BuildResult) {
			reg := regexp.MustCompile(`(?im)\<\/body\>`)
			indexFileContent = reg.ReplaceAllString(
				indexFileContent,
				`<script data-version="0.2" src="vendor.js"></script>
				<script data-version="0.2" type="module" src="main.js"></script>
				</body>`,
			)

			reg = regexp.MustCompile(`(?im)\<\/head\>`)
			indexFileContent = reg.ReplaceAllString(
				indexFileContent,
				`<link rel="stylesheet" href="main.css">          
				</head>`,
			)

			err := os.WriteFile(
				path.Join(outPath, "index.html"),
				[]byte(indexFileContent),
				0644,
			)
			check(err)
		})

	},
}

var zoneJSIncluder = api.Plugin{
	Name: "zoneJs",
	Setup: func(build api.PluginBuild) {
		build.OnLoad(api.OnLoadOptions{Filter: `main\.ts$`},
			func(args api.OnLoadArgs) (api.OnLoadResult, error) {
				mainTs, err := os.ReadFile(args.Path)
				check(err)
				contents := "import 'zone.js';\n" + string(mainTs)

				return api.OnLoadResult{
					Contents: &contents,
					Loader:   api.LoaderTS,
				}, nil
			})
	},
}

var angularComponentDecoratorPlugin = api.Plugin{
	Name: "componentDecorator",
	Setup: func(build api.PluginBuild) {
		build.OnLoad(api.OnLoadOptions{Filter: `src.*\.(component|pipe|service|directive|guard|module)\.ts$`},
			func(args api.OnLoadArgs) (api.OnLoadResult, error) {
				source, err := os.ReadFile(args.Path)
				check(err)

				contents := string(source)
				// componentName := ""
				// componentID := ""

				reg := regexp.MustCompile(`module\.ts$`)
				if reg.Match([]byte(args.Path)) {
					contents = "import '@angular/compiler';\n" + contents
				}

				templateReg := regexp.MustCompile(`(?m)^ *templateUrl *\: *['"]*([^'"]*)['"]`)
				if templateReg.Match([]byte(contents)) {
					templateUrl := templateReg.FindStringSubmatch(contents)
					contents = templateReg.ReplaceAllString(
						contents,
						"template: templateSource || ''",
					)
					contents = "import templateSource from '" +
						templateUrl[1] + "';\n" + contents
				}

				styleReg := regexp.MustCompile(`(?m)^ *styleUrls *: *\[(['"].*['"])\]`)
				if styleReg.Match([]byte(contents)) {
					contents = styleReg.ReplaceAllString(
						contents,
						"styleUrls: [],",
					)
				}

				return api.OnLoadResult{
					Contents: &contents,
					Loader:   api.LoaderTS,
				}, nil
			})
	},
}

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
			indexFileProcessor,
			zoneJSIncluder,
			angularComponentDecoratorPlugin,
		},
		AbsWorkingDir: workingDir,
	}

	tsConfigPath := path.Join(workingDir, "tsconfig.json")
	if statPath(tsConfigPath) {
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
