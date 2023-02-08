cd app\\services\\node
delete node.exe
go build node.go
cd ..\\..\\..
start  "NODE:%1" app\services\node\node.exe --source %2 --node %1
pause
 