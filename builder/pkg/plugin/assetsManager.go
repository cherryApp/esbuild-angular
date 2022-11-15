package plugin

import (
	"os"
	"path"

	cp "github.com/otiai10/copy"

	cpf "github.com/nmrshll/go-cp"

	"github.com/evanw/esbuild/pkg/api"
)

func GetAssetManager(workingDir string, assets []interface{}, outPath string) api.Plugin {
	return api.Plugin{
		Name: "assetManager",
		Setup: func(build api.PluginBuild) {

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

		},
	}
}
