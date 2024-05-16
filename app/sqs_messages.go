package app

/*** Send Message Response */

type SendMessageResult struct {
	MD5OfMessageAttributes string `xml:"MD5OfMessageAttributes"`
	MD5OfMessageBody       string `xml:"MD5OfMessageBody"`
	MessageId              string `xml:"MessageId"`
	SequenceNumber         string `xml:"SequenceNumber"`
}

type SendMessageResponse struct {
	Xmlns    string            `xml:"xmlns,attr"`
	Result   SendMessageResult `xml:"SendMessageResult"`
	Metadata ResponseMetadata  `xml:"ResponseMetadata"`
}

/*** Delete Message Response */
type DeleteMessageResponse struct {
	Xmlns    string           `xml:"xmlns,attr,omitempty"`
	Metadata ResponseMetadata `xml:"ResponseMetadata,omitempty"`
}

type DeleteQueueResponse struct {
	Xmlns    string           `xml:"xmlns,attr,omitempty"`
	Metadata ResponseMetadata `xml:"ResponseMetadata,omitempty"`
}

type DeleteMessageBatchResultEntry struct {
	Id string `xml:"Id"`
}

type SendMessageBatchResultEntry struct {
	Id                     string `xml:"Id"`
	MessageId              string `xml:"MessageId"`
	MD5OfMessageBody       string `xml:"MD5OfMessageBody,omitempty"`
	MD5OfMessageAttributes string `xml:"MD5OfMessageAttributes,omitempty"`
	SequenceNumber         string `xml:"SequenceNumber"`
}

type BatchResultErrorEntry struct {
	Code        string `xml:"Code"`
	Id          string `xml:"Id"`
	Message     string `xml:"Message,omitempty"`
	SenderFault bool   `xml:"SenderFault"`
}

type DeleteMessageBatchResult struct {
	Entry []DeleteMessageBatchResultEntry `xml:"DeleteMessageBatchResultEntry"`
	Error []BatchResultErrorEntry         `xml:"BatchResultErrorEntry,omitempty"`
}

/*** Delete Message Batch Response */
type DeleteMessageBatchResponse struct {
	Xmlns    string                   `xml:"xmlns,attr,omitempty"`
	Result   DeleteMessageBatchResult `xml:"DeleteMessageBatchResult"`
	Metadata ResponseMetadata         `xml:"ResponseMetadata,omitempty"`
}

type SendMessageBatchResult struct {
	Entry []SendMessageBatchResultEntry `xml:"SendMessageBatchResultEntry"`
	Error []BatchResultErrorEntry       `xml:"BatchResultErrorEntry,omitempty"`
}

/*** Delete Message Batch Response */
type SendMessageBatchResponse struct {
	Xmlns    string                 `xml:"xmlns,attr,omitempty"`
	Result   SendMessageBatchResult `xml:"SendMessageBatchResult"`
	Metadata ResponseMetadata       `xml:"ResponseMetadata,omitempty"`
}

/*** Purge Queue Response */
type PurgeQueueResponse struct {
	Xmlns    string           `xml:"xmlns,attr,omitempty"`
	Metadata ResponseMetadata `xml:"ResponseMetadata,omitempty"`
}

/*** Get Queue Url Response */
type GetQueueUrlResult struct {
	QueueUrl string `xml:"QueueUrl,omitempty"`
}

type GetQueueUrlResponse struct {
	Xmlns    string            `xml:"xmlns,attr,omitempty"`
	Result   GetQueueUrlResult `xml:"GetQueueUrlResult"`
	Metadata ResponseMetadata  `xml:"ResponseMetadata,omitempty"`
}

type SetQueueAttributesResponse struct {
	Xmlns    string           `xml:"xmlns,attr,omitempty"`
	Metadata ResponseMetadata `xml:"ResponseMetadata,omitempty"`
}
