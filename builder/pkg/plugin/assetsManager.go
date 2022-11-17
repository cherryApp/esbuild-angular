package plugin

import (
	"os"
	"path"
  "path/filepath"
  "fmt"
  "regexp"
  "strings"

	cp "github.com/otiai10/copy"

	cpf "github.com/nmrshll/go-cp"

	"github.com/evanw/esbuild/pkg/api"

  "github.com/wellington/go-libsass"

  "cherryApp/esbuild-angular/pkg/util"
)

// https://github.com/skeeto/w64devkit/releases/download/v1.17.0/w64devkit-1.17.0.zip
// go install github.com/wellington/go-libsass

var regexpDataUrl = regexp.MustCompile(`data\:`);
var regexpPathSeparator = regexp.MustCompile(`\\|\/`);
var regexpSourcemap = regexp.MustCompile(`\/\*.*sourceMappingURL\=.*\*\/`)
var regexpScssFile = regexp.MustCompile(`\.scss$`)

// Compile sass files.
func SassCompiler(outPath string, sassFile string) string {

  return ""
}

func UrlUnpacker(workingDir string, outPath string, cssPath string, cssContent string) string {
  var re = regexp.MustCompile(`(?m)url\(['"]?([^\)'"\?]*)[\"\?\)]?`)
  var matches = re.FindAllStringSubmatch(cssContent, -1)
  if len(matches) == 0 {
    return cssContent
  }

  cssContent = regexpSourcemap.ReplaceAllString(cssContent, "")

  parentDir := filepath.Dir(cssPath)
  pathSeparator := string(os.PathSeparator)
  for _, match := range matches {
    if regexpDataUrl.MatchString(match[0]) {
      continue
    }

    // [url('../fonts/fontawesome-webfont.eot?, ../fonts/fontawesome-webfont.eot]
    urlPath := regexpPathSeparator.ReplaceAllString(match[1], pathSeparator)
    sourcePath := path.Join(parentDir, urlPath)
    fileName := filepath.Base(urlPath)
    targetPath := path.Join(outPath, fileName)

    cpf.CopyFile(sourcePath, targetPath)
    cssContent = strings.Replace(cssContent, match[1], fileName, 2)
  }

  return cssContent
}

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

        styles := util.GetProjectOption("architect.build.options.styles")
        for _, v := range styles.([]interface{}) {
          stylePath := path.Join(workingDir, v.(string));
          styleContent, err := os.ReadFile(stylePath)
          if err != nil {
            fmt.Println("ERROR! In angular.json styles wrong filepath:", stylePath)
          } else {
            content := ""
            if regexpScssFile.MatchString(v.(string)) {
              content = SassCompiler(outPath, stylePath)
            } else {
              content = UrlUnpacker(workingDir, outPath, stylePath, string(styleContent))
            }

            cssContent += content + "\n\n"
          }
        }

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
            fmt.Println("ERROR! In angular.json scripts wrong filepath:", scriptPath)
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
