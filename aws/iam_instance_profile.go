package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/gruntwork-io/cloud-nuke/logging"
	"github.com/gruntwork-io/go-commons/errors"
	"github.com/hashicorp/go-multierror"
)

// List all IAM users in the AWS account and returns a slice of the UserNames
func getAllIamInstanceProfiles(session *session.Session) ([]*string, error) {
	svc := iam.New(session)

	var profileNames []*string

	// TODO: Probably use ListRoles together with ListRolesPages in case there are lots of roles
	output, err := svc.ListInstanceProfiles(&iam.ListInstanceProfilesInput{})
	if err != nil {
		return nil, errors.WithStackTrace(err)
	}

	for _, profile := range output.InstanceProfiles {
		profileNames = append(profileNames, profile.InstanceProfileName)
	}

	return profileNames, nil
}

func detachRoles(svc *iam.IAM, profileName *string) error {
	resp, err := svc.GetInstanceProfile(&iam.GetInstanceProfileInput{
		InstanceProfileName: profileName,
	})
	if err != nil {
		return errors.WithStackTrace(err)
	}

	for _, role := range resp.InstanceProfile.Roles {
		_, err := svc.RemoveRoleFromInstanceProfile(&iam.RemoveRoleFromInstanceProfileInput{
			InstanceProfileName: profileName,
			RoleName:            role.RoleName,
		})
		if err != nil {
			return errors.WithStackTrace(err)
		}

		logging.Logger.Infof("[OK] Removed Role %s from Instance Profile %s", aws.StringValue(role.RoleName), aws.StringValue(profileName))
	}

	return nil
}

func deleteIamInstanceProfile(svc *iam.IAM, profileName *string) error {
	_, err := svc.DeleteInstanceProfile(&iam.DeleteInstanceProfileInput{
		InstanceProfileName: profileName,
	})
	if err != nil {
		return errors.WithStackTrace(err)
	}

	return nil
}

// Nuke a single Instance Profile
func nukeInstanceProfile(svc *iam.IAM, profileName *string) error {
	// Functions used to really nuke an IAM Instance Profile as a profile can have many attached
	// items we need delete/detach them before actually deleting it.
	// NOTE: The actual role deletion should always be the last one. This way we
	// can guarantee that it will fail if we forgot to delete/detach an item.
	functions := []func(svc *iam.IAM, profileName *string) error{
		detachRoles,
		deleteIamInstanceProfile,
	}

	for _, fn := range functions {
		if err := fn(svc, profileName); err != nil {
			return err
		}
	}

	return nil
}

// Delete all IAM Instance Profiles
func nukeAllIamInstanceProfiles(session *session.Session, profileNames []*string) error {
	if len(profileNames) == 0 {
		logging.Logger.Info("No IAM Instance Profiles to nuke")
		return nil
	}

	logging.Logger.Info("Deleting all IAM Instance Profiles")

	deletedResources := 0
	svc := iam.New(session)
	multiErr := new(multierror.Error)

	for _, profileName := range profileNames {
		err := nukeInstanceProfile(svc, profileName)
		if err != nil {
			logging.Logger.Errorf("[Failed] %s", err)
			multierror.Append(multiErr, err)
		} else {
			deletedResources++
			logging.Logger.Infof("Deleted IAM Instace Profile: %s", *profileName)
		}
	}

	logging.Logger.Infof("[OK] %d IAM Instance Profile(s) terminated", deletedResources)
	return multiErr.ErrorOrNil()
}
