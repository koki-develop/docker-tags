package awsutil

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

type SessionConfig struct {
	Profile string
	Region  string
}

func NewSession(cfg *SessionConfig) (*session.Session, error) {
	awscfg := &aws.Config{
		Region: aws.String(cfg.Region),
	}
	if cfg.Profile != "" {
		awscfg.Credentials = credentials.NewSharedCredentials("", cfg.Profile)
	}

	sess, err := session.NewSession(awscfg)
	if err != nil {
		return nil, err
	}

	return sess, nil
}
