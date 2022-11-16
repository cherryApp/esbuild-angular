package plugin

import (
	"os"
	"path"
  "fmt"

	cp "github.com/otiai10/copy"

	cpf "github.com/nmrshll/go-cp"

	"github.com/evanw/esbuild/pkg/api"

  "cherryApp/esbuild-angular/pkg/util"
)

func GetAssetManager(workingDir string, assets []interface{}, outPath string) api.Plugin {
	return api.Plugin{
		Name: "assetManager",
		Setup: func(build api.PluginBuild) {

      // Copy assets.
			build.OnStart(func() (api.OnStartResult, error) {
        for _, v := range assets {
          assetPath := path.Join(workingDir, v.(string))
          file, err := os.Open(assetPath)
          if err != nil {
            continue
          }

          fileInfo, _ := file.Stat()
          if fileInfo.IsDir() {
            cp.Copy(assetPath, path.Join(outPath, path.Base(assetPath) ) )
          } else {
            dst := path.Join(outPath, path.Base(assetPath) )
            cpf.CopyFile(assetPath, dst)
          }
        }
				return api.OnStartResult{}, nil
			})

      // Parse and copy css files.
      build.OnStart(func() (api.OnStartResult, error) {
        cssContent := ""

        err := os.WriteFile(
          path.Join(outPath, "main.css"),
          []byte(cssContent),
          0660,
        )
        util.Check(err)

				return api.OnStartResult{}, nil
			})

      // Merge and copy .js files.
      build.OnStart(func() (api.OnStartResult, error) {
        vendorJSContent := ""

        // Get scripts options from the angular.json
        scripts := util.GetProjectOption("architect.build.options.scripts")
        for _, v := range scripts.([]interface{}) {
          scriptPath := path.Join(workingDir, v.(string));
          scriptContent, err := os.ReadFile(scriptPath)
          if err != nil {
            fmt.Println("ERROR! Wrong filepath in: ", scriptPath)
          } else {
            vendorJSContent += string(scriptContent) + "\n\n"
          }
        }

        err := os.WriteFile(
          path.Join(outPath, "vendor.js"),
          []byte(vendorJSContent),
          0660,
        )
        util.Check(err)

				return api.OnStartResult{}, nil
			})

		},
	}
}
