import 'dart:io';
import 'package:sass/sass.dart' as sass;

void main(List<String> arguments) {
  // Stopwatch stopwatch = new Stopwatch()..start();
  var result = sass.compileToResult(
    arguments[1],
    // logger: logger,
    // importers: importers,
    loadPaths: [arguments[0]],
    // packageConfig: packageConfig,
    // functions: functions,
    // style: style,
    // quietDeps: quietDeps,
    verbose: true,
    // sourceMap: null,
    // charset: 'utf8',
  );
  // new File(arguments[1]).writeAsStringSync(result);
  print(result.css);
  // print('doSomething() executed in ${stopwatch.elapsed}');
}
