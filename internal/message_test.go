package internal

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMessage_ToJson(t *testing.T) {
	message := Message{Message: "testmessage", UserName: "testuser", System: false, Transferred: false}
	assert.Equal(t, []byte(`{"message":"testmessage","user_name":"testuser","system":false,"transferred":true}`), message.ToTransferredJson())
}
