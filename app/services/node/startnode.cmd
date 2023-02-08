delete node.exe
go build node.go
start  "NODE:%1" node.exe --source %2 --node %1
pause
 