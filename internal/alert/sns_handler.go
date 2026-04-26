package alert

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

// snsPublisher abstracts the SNS Publish call for testing.
type snsPublisher interface {
	Publish(ctx context.Context, params *sns.PublishInput, optFns ...func(*sns.Options)) (*sns.PublishOutput, error)
}

// SNSHandler sends alerts to an AWS SNS topic.
type SNSHandler struct {
	client   snsPublisher
	topicARN string
}

// NewSNSHandler creates an SNSHandler using the default AWS credential chain.
// topicARN must be a valid SNS topic ARN.
func NewSNSHandler(topicARN string) (*SNSHandler, error) {
	if topicARN == "" {
		return nil, fmt.Errorf("sns: topic ARN must not be empty")
	}
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, fmt.Errorf("sns: failed to load AWS config: %w", err)
	}
	return &SNSHandler{
		client:   sns.NewFromConfig(cfg),
		topicARN: topicARN,
	}, nil
}

// Send publishes an alert event to the configured SNS topic.
func (h *SNSHandler) Send(e Event) error {
	subject := fmt.Sprintf("[portwatch] %s port %d", e.Kind, e.Port)
	body := FormatAlert(e)
	_, err := h.client.Publish(context.Background(), &sns.PublishInput{
		TopicArn: aws.String(h.topicARN),
		Subject:  aws.String(subject),
		Message:  aws.String(body),
	})
	if err != nil {
		return fmt.Errorf("sns: publish failed: %w", err)
	}
	return nil
}
