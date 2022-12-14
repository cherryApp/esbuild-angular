package plugin

import (
	"os"
	"path"
	"regexp"

	"github.com/evanw/esbuild/pkg/api"

	"cherryApp/esbuild-angular/pkg/util"
)

var regexpBaseHref = regexp.MustCompile(`\<base *href\=\"\/\"\>`)

func GetIndexFileProcessor(indexFilePath string, outPath string) api.Plugin {
	return api.Plugin{
		Name: "indexProcessor",
		Setup: func(build api.PluginBuild) {

			indexFileContent := ""

			build.OnStart(func() (api.OnStartResult, error) {
				indexContent, err := os.ReadFile(indexFilePath)
				util.Check(err)
				indexFileContent = string(indexContent)
				return api.OnStartResult{}, nil
			})

			build.OnEnd(func(result *api.BuildResult) {
				reg := regexp.MustCompile(`(?im)\<\/body\>`)
				baseHref := util.GetRuntimeOption("base-href").(string)
				indexFileContent = reg.ReplaceAllString(
					indexFileContent,
					`<script data-version="0.2" src="vendor.js"></script>
					<script data-version="0.2" type="module" src="main.js"></script>
					</body>`,
				)

				if util.GetRuntimeOption("serve").(bool) {
					indexFileContent = reg.ReplaceAllString(
						indexFileContent,
						util.WsScript+`
						</body>`,
					)
				}

				// <base href="/"> baseHref
				if baseHref != "/" {
					indexFileContent = regexpBaseHref.ReplaceAllString(
						indexFileContent,
						`<base href="`+baseHref+`">`,
					)
				}

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
