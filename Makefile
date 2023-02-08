all:	node

node:
	go build app\services\node\node.go
	copy node.exe distribute\node.exe
	del node.exe
	copy app\services\node\startnode.cmd distribute\startnode.cmd

clean:	cleannode

cleannode:
	del distribute\node.exe


