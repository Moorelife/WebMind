package system

import (
	"net"
	"reflect"
	"testing"
)

var TestAddress = net.TCPAddr{
	IP:   []byte{1, 2, 3, 4},
	Port: 14285,
	Zone: "",
}

func TestNewNode(t *testing.T) {
	type args struct {
		address net.TCPAddr
	}
	tests := []struct {
		name string
		args args
		want *Node
	}{
		{name: "NewNodeGoodFlow", args: args{address: TestAddress}, want: &Node{Address: TestAddress}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewNode(tt.args.address); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewNode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNode_ToJSON(t *testing.T) {
	expected := `{
  "Address": {
    "IP": "1.2.3.4",
    "Port": 14285,
    "Zone": ""
  }
}`
	type fields struct {
		Address net.TCPAddr
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{name: "ToJSONGoodFlow", fields: fields{Address: TestAddress}, want: string(expected)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &Node{
				Address: tt.fields.Address,
			}
			got := n.ToJSON()
			if got != tt.want {
				t.Errorf("ToJSON() \n%v\n\n%v\n", got, tt.want)
			}
		})
	}
}
