package gosqs

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/Admiral-Piett/goaws/app"
	"github.com/stretchr/testify/assert"
)

func TestChangeMessageVisibility_POST_SUCCESS(t *testing.T) {
	req, err := http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	app.SyncQueues.Queues["testing"] = &app.Queue{Name: "testing"}
	app.SyncQueues.Queues["testing"].Messages = []app.Message{{
		MessageBody:   []byte("test1"),
		ReceiptHandle: "123",
	}}

	form := url.Values{}
	form.Add("Action", "ChangeMessageVisibility")
	form.Add("QueueUrl", "http://localhost:4100/queue/testing")
	form.Add("VisibilityTimeout", "0")
	form.Add("ReceiptHandle", "123")
	form.Add("Version", "2012-11-05")
	req.PostForm = form

	status, _ := ChangeMessageVisibilityV1(req)
	assert.Equal(t, status, http.StatusOK)
}
