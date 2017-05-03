package podmon

import (
	"context"
	"sync"

	log "github.com/Sirupsen/logrus"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

// SNSEndpoint extends the AWS SNS type for our interface
type SNSEndpoint struct {
	*sns.SNS
}

// SNSMessage ...
type SNSMessage struct {
	Subject  string
	Message  string
	TopicARN string
}

// NewSNSEndpoint returns an SNS publisher
func NewSNSEndpoint() *SNSEndpoint {
	sess := session.Must(
		session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))
	return &SNSEndpoint{sns.New(sess)}
}

// EventLoop for firing messages
func (s *SNSEndpoint) EventLoop(ctx context.Context, wg *sync.WaitGroup, c chan SNSMessage) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			break
		case msg := <-c:
			log.Debugf("Sending the following to SNS: %#v", msg)
			resp, err := s.Send(msg)
			if err != nil {
				log.Errorf("SNS error: %s", err)
			} else {
				log.Debugf("Got a %d from sending the following to SNS: %#v", resp, msg)
			}
		}
	}
}

// Send sends a message to an SNS endpoint
func (s SNSEndpoint) Send(msg SNSMessage) (string, error) {
	params := &sns.PublishInput{
		Subject:  aws.String(msg.Subject),
		Message:  aws.String(msg.Message),
		TopicArn: aws.String(msg.TopicARN),
	}
	resp, err := s.Publish(params)
	return resp.String(), err
}
