package aws

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	awsgo "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/gruntwork-io/cloud-nuke/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListIamInstanceProfiles(t *testing.T) {
	t.Parallel()

	region, err := getRandomRegion()
	require.NoError(t, err)

	session, err := session.NewSession(&awsgo.Config{
		Region: awsgo.String(region)},
	)
	require.NoError(t, err)

	intanceProfileNames, err := getAllIamInstanceProfiles(session)
	require.NoError(t, err)

	assert.NotEmpty(t, intanceProfileNames)
}

func createTestInstanceProfile(t *testing.T, session *session.Session, name string) error {
	svc := iam.New(session)

	input := &iam.CreateInstanceProfileInput{
		InstanceProfileName: aws.String(name),
	}

	_, err := svc.CreateInstanceProfile(input)
	require.NoError(t, err)

	return nil
}

func TestCreateIamInstanceProfile(t *testing.T) {
	t.Parallel()

	region, err := getRandomRegion()
	require.NoError(t, err)

	session, err := session.NewSession(&awsgo.Config{
		Region: awsgo.String(region)},
	)
	require.NoError(t, err)

	name := "cloud-nuke-test-" + util.UniqueID()
	profileNames, err := getAllIamInstanceProfiles(session)
	require.NoError(t, err)
	assert.NotContains(t, awsgo.StringValueSlice(profileNames), name)

	err = createTestInstanceProfile(t, session, name)
	defer nukeAllIamInstanceProfiles(session, []*string{&name})
	require.NoError(t, err)

	profileNames, err = getAllIamInstanceProfiles(session)
	require.NoError(t, err)
	assert.Contains(t, awsgo.StringValueSlice(profileNames), name)
}

func TestNukeIamInstanceProfile(t *testing.T) {
	t.Parallel()

	region, err := getRandomRegion()
	require.NoError(t, err)

	session, err := session.NewSession(&awsgo.Config{
		Region: awsgo.String(region)},
	)
	require.NoError(t, err)

	name := "cloud-nuke-test-" + util.UniqueID()
	err = createTestInstanceProfile(t, session, name)
	require.NoError(t, err)

	err = nukeAllIamInstanceProfiles(session, []*string{&name})
	require.NoError(t, err)
}
