package plugin

import (
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/evanw/esbuild/pkg/api"

	"cherryApp/esbuild-angular/pkg/util"
)

var regexpSeparator = regexp.MustCompile(`\\|\/`)
var regexpPathCleaner = regexp.MustCompile(`[\'\"]`)
var regexpStyleUrls = regexp.MustCompile(`(?m)^ *styleUrls *: *\[(['"].*['"])\]`)
var regexpExtractStyleUrls = regexp.MustCompile(`styleUrls *: *\[(['"].*['"])\]`)
var regexpConstructor = regexp.MustCompile(`(?m)constructor *\(([^\)]*)`)
var regexpConstructorParameterCleaner = regexp.MustCompile(`[\n\r]`)

func InjectStyle(parentPath string, contents string) string {
	parentPath = regexpSeparator.ReplaceAllString(parentPath, "/")

	cssContent := ""
	scssPaths := regexpExtractStyleUrls.FindAllSubmatch([]byte(contents), -1)
	for _, match := range scssPaths {
		cleanedName := regexpPathCleaner.ReplaceAllString(string(match[1]), "")
		scssPath := path.Join(path.Dir(parentPath), cleanedName)
		scssPath = path.Clean(scssPath)
		cssContent += SassCompiler("", scssPath)
	}

	return regexpStyleUrls.ReplaceAllString(
		contents,
		"styles: [`"+cssContent+"`],",
	)
}

func AddInjects(contents string) string {
	matches := regexpConstructor.FindAllSubmatch([]byte(contents), -1)
	if len(matches) < 1 {
		return contents
	}

	for _, match := range matches {
		flat := regexpConstructorParameterCleaner.ReplaceAll(match[1], []byte(""))
		params := strings.Split(string(flat), ",")
		injectedParams := []string{}
		for _, param := range params {
			if len(param) < 3 {
				continue
			}
			phrase := strings.Split(param, ":")
			injectedParams = append(
				injectedParams,
				("@Inject(" + phrase[1] + ") " + param),
			)
		}

		contents = strings.ReplaceAll(
			contents,
			string(match[1]),
			strings.Join(injectedParams, ","),
		)
	}

	return contents
}

var AngularComponentDecoratorPlugin = api.Plugin{
	Name: "componentDecorator",
	Setup: func(build api.PluginBuild) {
		build.OnLoad(api.OnLoadOptions{Filter: `src.*\.(component|pipe|service|directive|guard|module)\.ts$`},
			func(args api.OnLoadArgs) (api.OnLoadResult, error) {
				source, err := os.ReadFile(args.Path)
				util.Check(err)

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

				if regexpStyleUrls.Match([]byte(contents)) {
					contents = InjectStyle(args.Path, contents)
				}

				contents = AddInjects(contents)

				return api.OnLoadResult{
					Contents: &contents,
					Loader:   api.LoaderTS,
				}, nil
			})
	},
}
