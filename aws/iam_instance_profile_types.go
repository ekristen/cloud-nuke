package aws

import (
	awsgo "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gruntwork-io/go-commons/errors"
)

// IAMInstanceProfiles - represents all IAMInstanceProfiles on the AWS Account
type IAMInstanceProfiles struct {
	ProfileNames []string
}

// ResourceName - the simple name of the aws resource
func (r IAMInstanceProfiles) ResourceName() string {
	return "iam-instance-profile"
}

// ResourceIdentifiers - The IAM UserNames
func (r IAMInstanceProfiles) ResourceIdentifiers() []string {
	return r.ProfileNames
}

// Tentative batch size to ensure AWS doesn't throttle
func (r IAMInstanceProfiles) MaxBatchSize() int {
	return 200
}

// Nuke - nuke 'em all!!!
func (r IAMInstanceProfiles) Nuke(session *session.Session, identifiers []string) error {
	if err := nukeAllIamInstanceProfiles(session, awsgo.StringSlice(identifiers)); err != nil {
		return errors.WithStackTrace(err)
	}

	return nil
}
