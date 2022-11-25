# EsbuildAngular

## TODO
- OK -> [SCSS](https://siongui.github.io/2016/01/28/go-compile-sass-scss/)
- OK -> Hoisting injects from the constructor.
- OK -> Compile component-level .css and .scss correctly.
- OK -> OS-specific run.
- OK -> SCSS: using smaller Dart compiler.
- OK -> Live server writen in Go.
- IP -> Set --base-href option.

## Tips
### Communitace between two go programs
https://blog.kowalczyk.info/article/wOYk/advanced-command-execution-in-go-with-osexec.html
https://stackoverflow.com/questions/60640083/exec-wait-with-a-modified-stdin-waits-indefinitely
https://michelenasti.com/2020/09/16/how-to-read-and-write-from-stdin-and-stdout-in-go.html

### Run files based on OS
```json
"scripts": {
    "test": "run-script-os",
    "test:darwin:linux": "export NODE_ENV=test && mocha",
    "test:win32": "SET NODE_ENV=test&& mocha"
}
```

