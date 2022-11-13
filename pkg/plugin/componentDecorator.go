package plugin

import (
	"os"
	"regexp"

	"github.com/evanw/esbuild/pkg/api"

	"cherryapp/angular/pkg/util"
)

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
