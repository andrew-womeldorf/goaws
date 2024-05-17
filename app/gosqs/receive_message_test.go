package gosqs

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"testing"
	"time"

	"github.com/Admiral-Piett/goaws/app/models"
	"github.com/stretchr/testify/assert"
)

func TestReceiveMessageWaitTimeEnforcedV1(t *testing.T) {
	// create a queue
	req, err := http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	form := url.Values{}
	form.Add("Action", "CreateQueue")
	form.Add("QueueName", "waiting-queue")
	form.Add("Attribute.1.Name", "ReceiveMessageWaitTimeSeconds")
	form.Add("Attribute.1.Value", "2")
	form.Add("Version", "2012-11-05")
	req.PostForm = form

	rr := httptest.NewRecorder()
	status, _ := CreateQueueV1(req)

	assert.Equal(t, http.StatusOK, status)

	// receive message ensure delay
	req, err = http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	form = url.Values{}
	form.Add("Action", "ReceiveMessage")
	form.Add("QueueUrl", "http://localhost:4100/queue/waiting-queue")
	form.Add("Version", "2012-11-05")
	req.PostForm = form

	rr = httptest.NewRecorder()

	start := time.Now()
	status, _ = ReceiveMessageV1(req)
	elapsed := time.Since(start)

	assert.Equal(t, http.StatusOK, status)
	if elapsed < 2*time.Second {
		t.Fatal("handler didn't wait ReceiveMessageWaitTimeSeconds")
	}

	// send a message
	req, err = http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	form = url.Values{}
	form.Add("Action", "SendMessage")
	form.Add("QueueUrl", "http://localhost:4100/queue/waiting-queue")
	form.Add("MessageBody", "1")
	form.Add("Version", "2012-11-05")
	req.PostForm = form

	rr = httptest.NewRecorder()
	http.HandlerFunc(SendMessage).ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got \n%v want %v",
			status, http.StatusOK)
	}

	// receive message
	req, err = http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	form = url.Values{}
	form.Add("Action", "ReceiveMessage")
	form.Add("QueueUrl", "http://localhost:4100/queue/waiting-queue")
	form.Add("Version", "2012-11-05")
	req.PostForm = form

	rr = httptest.NewRecorder()

	start = time.Now()
	status, _ = ReceiveMessageV1(req)
	elapsed = time.Since(start)

	assert.Equal(t, http.StatusOK, status)
	if elapsed > 1*time.Second {
		t.Fatal("handler waited when message was available, expected not to wait")
	}
}

func TestReceiveMessage_CanceledByClientV1(t *testing.T) {
	// create a queue
	req, err := http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	form := url.Values{}
	form.Add("Action", "CreateQueue")
	form.Add("QueueName", "cancel-queue")
	form.Add("Attribute.1.Name", "ReceiveMessageWaitTimeSeconds")
	form.Add("Attribute.1.Value", "20")
	form.Add("Version", "2012-11-05")
	req.PostForm = form

	rr := httptest.NewRecorder()
	status, _ := CreateQueueV1(req)

	assert.Equal(t, http.StatusOK, status)

	var wg sync.WaitGroup
	ctx, cancelReceive := context.WithCancel(context.Background())

	wg.Add(1)
	go func() {
		defer wg.Done()
		// receive message (that will be canceled)
		req, err := http.NewRequest("POST", "/", nil)
		req = req.WithContext(ctx)
		if err != nil {
			t.Fatal(err)
		}

		form := url.Values{}
		form.Add("Action", "ReceiveMessage")
		form.Add("QueueUrl", "http://localhost:4100/queue/cancel-queue")
		form.Add("Version", "2012-11-05")
		req.PostForm = form

		status, resp := ReceiveMessageV1(req)
		assert.Equal(t, http.StatusOK, status)

		if len(resp.GetResult().(models.ReceiveMessageResult).Messages) != 0 {
			t.Fatal("expecting this ReceiveMessage() to not pickup this message as it should canceled before the Send()")
		}
	}()
	time.Sleep(100 * time.Millisecond) // let enought time for the Receive go to wait mode
	cancelReceive()                    // cancel the first ReceiveMessage(), make sure it will not pickup the sent message below
	time.Sleep(5 * time.Millisecond)

	// send a message
	req, err = http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	form = url.Values{}
	form.Add("Action", "SendMessage")
	form.Add("QueueUrl", "http://localhost:4100/queue/cancel-queue")
	form.Add("MessageBody", "12345")
	form.Add("Version", "2012-11-05")
	req.PostForm = form

	rr = httptest.NewRecorder()
	http.HandlerFunc(SendMessage).ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got \n%v want %v",
			status, http.StatusOK)
	}

	// receive message
	req, err = http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	form = url.Values{}
	form.Add("Action", "ReceiveMessage")
	form.Add("QueueUrl", "http://localhost:4100/queue/cancel-queue")
	form.Add("Version", "2012-11-05")
	req.PostForm = form

	rr = httptest.NewRecorder()

	start := time.Now()
	status, resp := ReceiveMessageV1(req)
	assert.Equal(t, http.StatusOK, status)
	elapsed := time.Since(start)

	result, ok := resp.GetResult().(models.ReceiveMessageResult)
	if !ok {
		t.Fatal("handler should return a message")
	}

	if len(result.Messages) == 0 || string(result.Messages[0].Body) == "12345\n" {
		t.Fatal("handler should return a message")
	}
	if elapsed > 1*time.Second {
		t.Fatal("handler waited when message was available, expected not to wait")
	}

	if timedout := waitTimeout(&wg, 2*time.Second); timedout {
		t.Errorf("expected ReceiveMessage() in goroutine to exit quickly due to cancelReceive() called")
	}
}

func TestReceiveMessage_WithConcurrentDeleteQueueV1(t *testing.T) {
	// create a queue
	req, err := http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	form := url.Values{}
	form.Add("Action", "CreateQueue")
	form.Add("QueueName", "waiting-queue")
	form.Add("Attribute.1.Name", "ReceiveMessageWaitTimeSeconds")
	form.Add("Attribute.1.Value", "1")
	form.Add("Version", "2012-11-05")
	req.PostForm = form

	rr := httptest.NewRecorder()
	status, _ := CreateQueueV1(req)

	assert.Equal(t, http.StatusOK, status)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got \n%v want %v",
			status, http.StatusOK)
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		// receive message
		req, err := http.NewRequest("POST", "/", nil)
		if err != nil {
			t.Fatal(err)
		}

		form := url.Values{}
		form.Add("Action", "ReceiveMessage")
		form.Add("QueueUrl", "http://localhost:4100/queue/waiting-queue")
		form.Add("Version", "2012-11-05")
		req.PostForm = form

		status, resp := ReceiveMessageV1(req)
		assert.Equal(t, http.StatusBadRequest, status)

		// Check the response body is what we expect.
		expected := "QueueNotFound"
		result := resp.GetResult().(models.ErrorResult)
		if result.Type != "Not Found" {
			t.Errorf("handler returned unexpected body: got %v want %v",
				result.Message, expected)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(10 * time.Millisecond) // 10ms to let the ReceiveMessage() block
		// delete queue message
		req, err := http.NewRequest("POST", "/", nil)
		if err != nil {
			t.Fatal(err)
		}
		form := url.Values{}
		form.Add("Action", "DeleteQueue")
		form.Add("QueueUrl", "http://localhost:4100/queue/waiting-queue")
		form.Add("Version", "2012-11-05")
		req.PostForm = form

		rr := httptest.NewRecorder()
		http.HandlerFunc(DeleteQueue).ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got \n%v want %v",
				status, http.StatusOK)
		}
	}()

	if timedout := waitTimeout(&wg, 2*time.Second); timedout {
		t.Errorf("concurrent handlers timeout, expecting both to return within timeout")
	}
}

func TestReceiveMessageDelaySecondsV1(t *testing.T) {
	// create a queue
	req, err := http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	form := url.Values{}
	form.Add("Action", "CreateQueue")
	form.Add("QueueName", "delay-seconds-queue")
	form.Add("Attribute.1.Name", "DelaySeconds")
	form.Add("Attribute.1.Value", "2")
	form.Add("Version", "2012-11-05")
	req.PostForm = form

	rr := httptest.NewRecorder()
	status, _ := CreateQueueV1(req)

	assert.Equal(t, http.StatusOK, status)

	// send a message
	req, err = http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	form = url.Values{}
	form.Add("Action", "SendMessage")
	form.Add("QueueUrl", "http://localhost:4100/queue/delay-seconds-queue")
	form.Add("MessageBody", "1")
	form.Add("Version", "2012-11-05")
	req.PostForm = form
	rr = httptest.NewRecorder()
	http.HandlerFunc(SendMessage).ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got \n%v want %v",
			status, http.StatusOK)
	}

	// receive message before delay is up
	req, err = http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	form = url.Values{}
	form.Add("Action", "ReceiveMessage")
	form.Add("QueueUrl", "http://localhost:4100/queue/delay-seconds-queue")
	form.Add("Version", "2012-11-05")
	req.PostForm = form
	rr = httptest.NewRecorder()
	status, _ = ReceiveMessageV1(req)
	assert.Equal(t, http.StatusOK, status)

	// receive message with wait should return after delay
	req, err = http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	form = url.Values{}
	form.Add("Action", "ReceiveMessage")
	form.Add("QueueUrl", "http://localhost:4100/queue/delay-seconds-queue")
	form.Add("WaitTimeSeconds", "10")
	form.Add("Version", "2012-11-05")
	req.PostForm = form
	rr = httptest.NewRecorder()
	start := time.Now()
	status, _ = ReceiveMessageV1(req)
	elapsed := time.Since(start)
	assert.Equal(t, http.StatusOK, status)
	if elapsed < 1*time.Second {
		t.Errorf("handler didn't wait at all")
	}
	if elapsed > 4*time.Second {
		t.Errorf("handler didn't need to wait all WaitTimeSeconds=10, only DelaySeconds=2")
	}
}
