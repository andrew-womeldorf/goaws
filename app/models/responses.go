package models

import "github.com/Admiral-Piett/goaws/app"

// NOTE: Every response in here MUST implement the `AbstractResponseBody` interface in order to be used
//  in `encodeResponse`

/*** Error Responses ***/
type ErrorResult struct {
	Type    string `xml:"Type,omitempty"`
	Code    string `xml:"Code,omitempty"`
	Message string `xml:"Message,omitempty"`
}

type ErrorResponse struct {
	Result    ErrorResult `xml:"Error"`
	RequestId string      `xml:"RequestId"`
}

func (r ErrorResponse) GetResult() interface{} {
	return r.Result
}

func (r ErrorResponse) GetRequestId() string {
	return r.RequestId
}

/*** Receive Message Response */
type ReceiveMessageResult struct {
	Message []*ResultMessage `json:"Message" xml:"Message,omitempty"`
}

type ReceiveMessageResponse struct {
	Xmlns    string               `xml:"xmlns,attr"`
	Result   ReceiveMessageResult `xml:"ReceiveMessageResult"`
	Metadata app.ResponseMetadata `xml:"ResponseMetadata"`
}

type ResultMessage struct {
	MessageId              string                    `xml:"MessageId,omitempty"`
	ReceiptHandle          string                    `xml:"ReceiptHandle,omitempty"`
	MD5OfBody              string                    `xml:"MD5OfBody,omitempty"`
	Body                   []byte                    `xml:"Body,omitempty"`
	MD5OfMessageAttributes string                    `xml:"MD5OfMessageAttributes,omitempty"`
	MessageAttributes      []*ResultMessageAttribute `xml:"MessageAttribute,omitempty"`
	Attributes             []*ResultAttribute        `xml:"Attribute,omitempty"`
}

type ResultMessageAttributeValue struct {
	DataType    string `xml:"DataType,omitempty"`
	StringValue string `xml:"StringValue,omitempty"`
	BinaryValue string `xml:"BinaryValue,omitempty"`
}

type ResultMessageAttribute struct {
	Name  string                       `xml:"Name,omitempty"`
	Value *ResultMessageAttributeValue `xml:"Value,omitempty"`
}

type ResultAttribute struct {
	Name  string `xml:"Name,omitempty"`
	Value string `xml:"Value,omitempty"`
}

type ChangeMessageVisibilityResult struct {
	Xmlns    string               `xml:"xmlns,attr"`
	Metadata app.ResponseMetadata `xml:"ResponseMetadata"`
}

func (r ReceiveMessageResponse) GetResult() interface{} {
	return r.Result
}

func (r ReceiveMessageResponse) GetRequestId() string {
	return r.Metadata.RequestId
}

/*** Create Queue Response */
type CreateQueueResult struct {
	QueueUrl string `json:"QueueUrl" xml:"QueueUrl"`
}

type CreateQueueResponse struct {
	Xmlns    string               `xml:"xmlns,attr"`
	Result   CreateQueueResult    `xml:"CreateQueueResult"`
	Metadata app.ResponseMetadata `xml:"ResponseMetadata"`
}

func (r CreateQueueResponse) GetResult() interface{} {
	return r.Result
}

func (r CreateQueueResponse) GetRequestId() string {
	return r.Metadata.RequestId
}

/*** List Queues Response */
type ListQueuesResult struct {
	// NOTE: the old XML sdks depend on QueueUrl, and the new JSON ones need QueueUrls
	QueueUrls []string `json:"QueueUrls" xml:"QueueUrl"`
}

type ListQueuesResponse struct {
	Xmlns    string               `xml:"xmlns,attr"`
	Result   ListQueuesResult     `xml:"ListQueuesResult"`
	Metadata app.ResponseMetadata `xml:"ResponseMetadata"`
}

func (r ListQueuesResponse) GetResult() interface{} {
	return r.Result
}

func (r ListQueuesResponse) GetRequestId() string {
	return r.Metadata.RequestId
}

/*** Get Queue Attributes ***/
type Attribute struct {
	Name  string `xml:"Name,omitempty"`
	Value string `xml:"Value,omitempty"`
}

type GetQueueAttributesResult struct {
	/* VisibilityTimeout, DelaySeconds, ReceiveMessageWaitTimeSeconds, ApproximateNumberOfMessages
	   ApproximateNumberOfMessagesNotVisible, CreatedTimestamp, LastModifiedTimestamp, QueueArn */
	Attrs []Attribute `xml:"Attribute,omitempty"`
}

type GetQueueAttributesResponse struct {
	Xmlns    string                   `xml:"xmlns,attr,omitempty"`
	Result   GetQueueAttributesResult `xml:"GetQueueAttributesResult"`
	Metadata app.ResponseMetadata     `xml:"ResponseMetadata,omitempty"`
}

func (r GetQueueAttributesResponse) GetResult() interface{} {
	result := map[string]string{}
	for _, attr := range r.Result.Attrs {
		result[attr.Name] = attr.Value
	}
	return map[string]map[string]string{"Attributes": result}
}

func (r GetQueueAttributesResponse) GetRequestId() string {
	return r.Metadata.RequestId
}
