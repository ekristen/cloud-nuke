package aws

import (
	awsgo "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gruntwork-io/go-commons/errors"
)

// CloudTrails - represents all ec2 instances
type CloudTrails struct {
	Arns []string
}

// ResourceName - the simple name of the aws resource
func (instance CloudTrails) ResourceName() string {
	return "cloudtrail-trail"
}

// ResourceIdentifiers - The instance ids of the ec2 instances
func (instance CloudTrails) ResourceIdentifiers() []string {
	return instance.Arns
}

func (instance CloudTrails) MaxBatchSize() int {
	// Tentative batch size to ensure AWS doesn't throttle
	return 200
}

// Nuke - nuke 'em all!!!
func (instance CloudTrails) Nuke(session *session.Session, identifiers []string) error {
	if err := nukeAllCloudTrailTrails(session, awsgo.StringSlice(identifiers)); err != nil {
		return errors.WithStackTrace(err)
	}

	return nil
}
