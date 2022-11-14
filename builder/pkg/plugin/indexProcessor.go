package plugin

import (
	"os"
	"path"
	"regexp"

	"github.com/evanw/esbuild/pkg/api"

	"cherryApp/esbuild-angular/pkg/util"
)

func GetIndexFileProcessor(srcPath string, outPath string) api.Plugin {
	return api.Plugin{
		Name: "indexProcessor",
		Setup: func(build api.PluginBuild) {

			indexFileContent := ""

			build.OnStart(func() (api.OnStartResult, error) {
				indexContent, err := os.ReadFile(path.Join(srcPath, "index.html"))
				util.Check(err)
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
				util.Check(err)
			})

		},
	}
}
