cd app\\services\\node
delete node.exe
go build node.go
start  "NODE:%1" node.exe --source %2 --port %1
pause
 