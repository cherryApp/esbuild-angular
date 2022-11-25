package util

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/evanw/esbuild/pkg/api"

	"cherryApp/esbuild-angular/pkg/fswatch"
)

var buildOptions api.BuildOptions
var indexFile *os.File

func serveSPS(fs http.FileSystem) http.Handler {
	fileServer := http.FileServer(fs)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := fs.Open(path.Clean(r.URL.Path)) // Do not allow path traversals.
		if err != nil {
			indexPath := path.Join(buildOptions.Outdir, "index.html")
			if indexFile == nil {
				indexFile, _ = os.Open(indexPath)
			}
			http.ServeContent(w, r, indexPath, time.Time{}, indexFile)
			return
		}
		fileServer.ServeHTTP(w, r)
	})
}

func FileWatcher(_buildOptions api.BuildOptions, callback func(string)) {
	go func() {

		recurse := true // include all sub directories

		skipDotFilesAndFolders := func(path string) bool {
			return strings.HasPrefix(filepath.Base(path), ".")
		}

		checkIntervalInSeconds := 2

		folderWatcher := fswatch.NewFolderWatcher(
			path.Join(_buildOptions.SourceRoot, "src"),
			recurse,
			skipDotFilesAndFolders,
			checkIntervalInSeconds,
		)

		folderWatcher.Start()

		for folderWatcher.IsRunning() {

			select {

			case <-folderWatcher.Modified():
				callback("New or modified items detected")
				// fmt.Println("New or modified items detected")

			case <-folderWatcher.Moved():
				callback("Items have been moved")
				// fmt.Println("Items have been moved")

				// case changes := <-folderWatcher.ChangeDetails():

				// 	fmt.Printf("%s\n", changes.String())
				// 	fmt.Printf("New: %#v\n", changes.New())
				// 	fmt.Printf("Modified: %#v\n", changes.Modified())
				// 	fmt.Printf("Moved: %#v\n", changes.Moved())

			}
		}

	}()
}

func LiveServer(_buildOptions api.BuildOptions) {
	buildOptions = _buildOptions
	var addr = "127.0.0.1:" + fmt.Sprintf("%v", GetRuntimeOption("port"))
	http.ListenAndServe(addr, serveSPS(http.Dir(buildOptions.Outdir)))
}
