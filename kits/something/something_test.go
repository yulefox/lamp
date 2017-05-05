package something

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/yulefox/lamp/kits/something/mock_something"
)

func TestMyThing(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockObj := mock_something.NewMockMyInterface(mockCtrl)
	c := mockObj.EXPECT().SomeMethod(4, "blah")
	c.AnyTimes()
}
