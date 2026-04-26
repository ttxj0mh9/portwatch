package alert

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sns"
)

// mockSNSPublisher records the last publish input and optionally returns an error.
type mockSNSPublisher struct {
	lastInput *sns.PublishInput
	errToReturn error
}

func (m *mockSNSPublisher) Publish(_ context.Context, params *sns.PublishInput, _ ...func(*sns.Options)) (*sns.PublishOutput, error) {
	m.lastInput = params
	return &sns.PublishOutput{}, m.errToReturn
}

func TestNewSNSHandler_MissingARN(t *testing.T) {
	_, err := NewSNSHandler("")
	if err == nil {
		t.Fatal("expected error for empty topic ARN")
	}
}

func TestSNSHandler_Send_Success(t *testing.T) {
	mock := &mockSNSPublisher{}
	h := &SNSHandler{client: mock, topicARN: "arn:aws:sns:us-east-1:123456789012:portwatch"}

	e := Event{
		Kind:      KindOpened,
		Port:      8080,
		Timestamp: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
	}

	if err := h.Send(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mock.lastInput == nil {
		t.Fatal("expected Publish to be called")
	}
	if *mock.lastInput.TopicArn != h.topicARN {
		t.Errorf("topic ARN mismatch: got %s", *mock.lastInput.TopicArn)
	}
}

func TestSNSHandler_Send_PublishError(t *testing.T) {
	mock := &mockSNSPublisher{errToReturn: errors.New("network error")}
	h := &SNSHandler{client: mock, topicARN: "arn:aws:sns:us-east-1:123456789012:portwatch"}

	e := Event{Kind: KindClosed, Port: 443, Timestamp: time.Now()}
	if err := h.Send(e); err == nil {
		t.Fatal("expected error when publish fails")
	}
}

func TestSNSHandler_Send_SubjectContainsPort(t *testing.T) {
	mock := &mockSNSPublisher{}
	h := &SNSHandler{client: mock, topicARN: "arn:aws:sns:us-east-1:123456789012:portwatch"}

	e := Event{Kind: KindOpened, Port: 9090, Timestamp: time.Now()}
	_ = h.Send(e)

	if mock.lastInput == nil {
		t.Fatal("expected Publish to be called")
	}
	subject := *mock.lastInput.Subject
	if subject == "" {
		t.Error("expected non-empty subject")
	}
}
