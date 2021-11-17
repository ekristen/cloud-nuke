package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudtrail"
	"github.com/gruntwork-io/cloud-nuke/logging"
	"github.com/gruntwork-io/go-commons/errors"
)

func getAllCloudTrailTrails(session *session.Session, region string) ([]*string, error) {
	svc := cloudtrail.New(session)

	output, err := svc.DescribeTrails(&cloudtrail.DescribeTrailsInput{})
	if err != nil {
		return nil, errors.WithStackTrace(err)
	}

	var arns []*string
	for _, trail := range output.TrailList {
		if aws.BoolValue(trail.IsOrganizationTrail) == true {
			continue
		}

		if trail.HomeRegion != aws.String(region) {
			continue
		}

		arns = append(arns, trail.TrailARN)
	}

	return arns, nil
}

func nukeAllCloudTrailTrails(session *session.Session, identifiers []*string) error {
	svc := cloudtrail.New(session)

	if len(identifiers) == 0 {
		logging.Logger.Infof("No CloudTrail Trails to nuke in region %s", *session.Config.Region)
		return nil
	}

	logging.Logger.Infof("Deleting all CloudTrail Trails in region %s", *session.Config.Region)

	var deletedTrails = 0

	for _, arn := range identifiers {
		_, err := svc.DeleteTrail(&cloudtrail.DeleteTrailInput{
			Name: arn,
		})
		if err != nil {
			logging.Logger.Errorf("[Failed] %s", err)
		} else {
			logging.Logger.Infof("[OK] CloudTrail Trail %s terminated in %s", arn, *session.Config.Region)
			deletedTrails++
		}

	}

	logging.Logger.Infof("[OK] %d CloudTrail Trail(s) terminated in %s", deletedTrails, *session.Config.Region)

	return nil
}
