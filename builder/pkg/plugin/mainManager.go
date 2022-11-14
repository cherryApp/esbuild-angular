package plugin

import (
	"os"

	"github.com/evanw/esbuild/pkg/api"

	"cherryApp/esbuild-angular/pkg/util"
)

func GetMainManager() api.Plugin {
	return api.Plugin{
		Name: "zoneJs",
		Setup: func(build api.PluginBuild) {
			build.OnLoad(api.OnLoadOptions{Filter: `main\.ts$`},
				func(args api.OnLoadArgs) (api.OnLoadResult, error) {
					mainTs, err := os.ReadFile(args.Path)
					util.Check(err)
					contents := "import 'zone.js';\n" + string(mainTs)

					return api.OnLoadResult{
						Contents: &contents,
						Loader:   api.LoaderTS,
					}, nil
				})
		},
	}
}
