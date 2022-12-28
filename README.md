# slogger [![Go Report Card](https://goreportcard.com/badge/github.com/BitlyTwiser/slogger)](https://goreportcard.com/report/github.com/BitlyTwiser/slogger)
```
  .--.--.     ,--,                                                      
 /  /    '. ,--.'|                                                      
|  :  /`. / |  | :     ,---.                                    __  ,-. 
;  |  |--`  :  : '    '   ,'\   ,----._,.  ,----._,.          ,' ,'/ /| 
|  :  ;_    |  ' |   /   /   | /   /  ' / /   /  ' /   ,---.  '  | |' | 
 \  \    `. '  | |  .   ; ,. :|   :     ||   :     |  /     \ |  |   ,' 
  `----.   \|  | :  '   | |: :|   | .\  .|   | .\  . /    /  |'  :  /   
  __ \  \  |'  : |__'   | .; :.   ; ';  |.   ; ';  |.    ' / ||  | '    
 /  /`--'  /|  | '.'|   :    |'   .   . |'   .   . |'   ;   /|;  : |    
'--'.     / ;  :    ;\   \  /  `---`-'| | `---`-'| |'   |  / ||  , ;    
  `--'---'  |  ,   /  `----'   .'__/\_: | .'__/\_: ||   :    | ---'     
             ---`-'            |   :    : |   :    : \   \  /           
                                \   \  /   \   \  /   `----'            
                                 `--`-'     `--`-'                      
```
The logger library using the experimental Golang slog package

## Usage:
- The package was designed to be a simple interface over the slog package extracing out some if the discoverability features and crafting a simple, easy logger.
- creating a new logger with ```NewLogger``` function. The input is an io.Writer, so stdout or a file work (any any io.Writer implementation)
- All variadic values passed to the ```LogEvent```/```LogError``` functions will be treated as key value pairs, any additional values, that are not mapped to their preceeding values, are added as "misc" fields to the JSON return.

### Example Usage:
- create a logger with ```slogger.NewLogger(os.Stdout)```
- or log to a file:
```
f, _ := os.OpenFile(testFile, os.O_RDWR|os.O_TRUNC, os.ModeAppend)
log := slogger.NewLogger(f)
stdoutLogger.LogEvent("info", "Send data to file") 
```
- slogger will accept any number of arbitrary elements and create a json structure from the arguments to the LogEvent functions.
- i.e.
```
stdoutLogger.LogEvent("info", "Something", "one", "two", "Another", false, "four", true)
```
output:
```
stdoutLogger.LogEvent("info", "Something four", "one", "two", "Another", false, "four", true)
```

- You can also pass a map to  the LogEvent functions as unput:
```
stdoutLogger.LogEvent("info", "Something four", map[string]any{"one": false})
```
output: 
```
{"time":"2022-12-27T19:19:22.711393515-08:00","level":"INFO","msg":"Something four","one":false}
```
- You can even pass a combination of map strings:
```
stdoutLogger.LogEvent("warn", "aaaa", "key", "value", "AnotherKey", false, "four", 123123, map[string]any{"testOne": 42069, "testTwo": false})
```
output:
```
stdoutLogger.LogEvent("warn", "aaaa", "key", "value", "AnotherKey", false, "four", 123123, map[string]any{"testOne": 42069, "testTwo": false})
```

- Example of Misc fields:
```
{"time":"2022-12-27T18:56:27.543409183-08:00","level":"INFO","msg":"Something four","one":"two","Another":false,"four":true,"miscFields":"bob"}
```
## Output:
- Example test output:
```
{"time":"2022-12-27T18:56:27.543409183-08:00","level":"INFO","msg":"Something four","one":"two","Another":false,"four":true,"miscFields":"bob"}
{"time":"2022-12-27T18:56:27.543564927-08:00","level":"INFO","msg":"Something four","one":"two","Another":false,"four":true}
{"time":"2022-12-27T18:56:27.543572397-08:00","level":"INFO","msg":"Something four","one":false,"Another":false,"four":true}
{"time":"2022-12-27T18:56:27.543577644-08:00","level":"INFO","msg":"A thing","key":"value","AnotherKey":false,"four":123123,"more":"","MORETEST":123123,"TestAgain":false,"miscFields":123123}
{"time":"2022-12-27T18:56:27.543583676-08:00","level":"WARN","msg":"aaaa","AnotherKey":false,"four":123123,"testOne":42069,"testTwo":false,"key":"value"}
{"time":"2022-12-27T18:56:27.543587977-08:00","level":"ERROR","msg":"error","err":"something died"}
{"time":"2022-12-27T18:56:27.54359305-08:00","level":"ERROR","msg":"error","1":2,"3":"masdasd","err":"something died"}
```
