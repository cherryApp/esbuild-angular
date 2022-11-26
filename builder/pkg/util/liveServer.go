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

	"github.com/gorilla/websocket"
)

var WsConn *websocket.Conn
var upgrader = websocket.Upgrader{} // use default options
var buildOptions api.BuildOptions
var indexFile *os.File

const WsScript = `<script type="text/javascript">
            (function() {
				const wsOrigin = document.head.querySelector("base").href.replace(/^http/, 'ws');
                ws = new WebSocket(wsOrigin + '__ws');
				ws.onmessage = function(evt) {
					if (evt.data === 'command:refresh') {
						location.reload();
					}
				}
				ws.onerror = function(evt) {
					console.log("Websocket Error, Live-Server: " + evt.data);
				}
			})();
        </script>
`

func serveSPS(fs http.FileSystem) http.Handler {
	fileServer := http.FileServer(fs)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Websocket
		if r.URL.Path == "/__ws" {
			serveWs(w, r, func(conn *websocket.Conn) {
				WsConn = conn
			})
			return
		}

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

			case <-folderWatcher.Moved():
				callback("Items have been moved")

			}
		}

	}()
}

func serveWs(w http.ResponseWriter, r *http.Request, callback func(*websocket.Conn)) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("upgrade:", err)
		return
	}

	callback(c)

	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			// fmt.Println("read:", err)
			break
		}

		err = c.WriteMessage(mt, message)
		if err != nil {
			fmt.Println("write:", err)
			break
		}
	}
}

func RefreshLiveServerPage() {
	if WsConn == nil {
		return
	}

	if err := WsConn.WriteMessage(1, []byte("command:refresh")); err != nil {
		fmt.Println("Error in Websocket:", err)
	}
}

func LiveServer(_buildOptions api.BuildOptions) {
	buildOptions = _buildOptions
	var addr = "127.0.0.1:" + fmt.Sprintf("%v", GetRuntimeOption("port"))
	fmt.Println("LiveServer runs on: http://" + addr)
	http.ListenAndServe(addr, serveSPS(http.Dir(buildOptions.Outdir)))
}
