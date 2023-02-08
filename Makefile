all:	node

node:
	go build app\services\node\node.go

clean:	cleannode

cleannode:
	del node.exe


