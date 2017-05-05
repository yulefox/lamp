package mock_hello

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/proto"
	pb "github.com/yulefox/lamp/kits/cli/proto"
	"golang.org/x/net/context"
)

type rpcMsg struct {
	msg proto.Message
}

func (r *rpcMsg) Matches(msg interface{}) bool {
	m, ok := msg.(proto.Message)
	if !ok {
		return false
	}
	return proto.Equal(m, r.msg)
}

func (r *rpcMsg) String() string {
	return fmt.Sprintf("is %s", r.msg)
}

func TestSayHelloClient(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := NewMockGreeterClient(ctrl)
	req := &pb.Hello_Req{Name: "unit_test"}
	mockClient.EXPECT().SayHello(
		gomock.Any(),
		&rpcMsg{msg: req},
	).Return(&pb.Hello_Res{Message: "mocked interface"}, nil)
	testSayHelloClient(t, mockClient)
}

func testSayHelloClient(t *testing.T, client *MockGreeterClient) {
	r, err := client.SayHello(context.Background(), &pb.Hello_Req{Name: "unit_test"})
	if err != nil || r.Message != "mocked interface" {
		t.Errorf("mocking failed")
	}
	t.Log("Res: ", r.Message)
}

func TestSayHelloServer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockServer := NewMockGreeterServer(ctrl)
	req := &pb.Hello_Req{Name: "unit_test"}
	mockServer.EXPECT().SayHello(
		gomock.Any(),
		&rpcMsg{msg: req},
	).Return(&pb.Hello_Res{Message: "mocked interface"}, nil)
	testSayHelloServer(t, mockServer)
}

func testSayHelloServer(t *testing.T, server *MockGreeterServer) {
	r, err := server.SayHello(context.Background(), &pb.Hello_Req{Name: "unit_test"})
	if err != nil || r.Message != "mocked interface" {
		t.Errorf("mocking failed")
	}
	t.Log("Res: ", r.Message)
}
