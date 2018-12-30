package factory

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/rds"
	log "github.com/sirupsen/logrus"
)

func UpdateOrCreateDBSubnetGroup(svc *rds.RDS, groupName, groupDescription string, subnets []string) (*rds.DBSubnetGroup, error) {
	var subnetGroup *rds.DBSubnetGroup

	subnetGroupName := aws.String(groupName)

	descGroupsInput := &rds.DescribeDBSubnetGroupsInput{
		DBSubnetGroupName: subnetGroupName,
	}

	descGroupsOutput, err := svc.DescribeDBSubnetGroups(descGroupsInput)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case rds.ErrCodeDBSubnetGroupNotFoundFault:
				log.Info(rds.ErrCodeDBSubnetGroupNotFoundFault, aerr.Error())

				return createSubnetGroup(svc, subnetGroupName, groupDescription, subnets)
			default:
				log.Warn(aerr)
				return nil, aerr
			}
		} else {
			return nil, err
		}
	}

	subnetGroup = descGroupsOutput.DBSubnetGroups[0]

	return subnetGroup, nil
}

func createSubnetGroup(svc *rds.RDS, subnetGroupName *string, groupDescription string, subnetIds []string) (*rds.DBSubnetGroup, error) {
	sIds := make([]*string, 0)
	for _, i := range subnetIds {
		sIds = append(sIds, aws.String(i))
	}

	groupInput := &rds.CreateDBSubnetGroupInput{
		DBSubnetGroupName:        subnetGroupName,
		DBSubnetGroupDescription: aws.String(groupDescription),
		SubnetIds:                sIds,
	}

	groupOutput, err := svc.CreateDBSubnetGroup(groupInput)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case rds.ErrCodeDBSubnetGroupAlreadyExistsFault:
				log.Warn(rds.ErrCodeDBSubnetGroupAlreadyExistsFault, aerr.Error())
				return nil, aerr
			case rds.ErrCodeDBSubnetGroupQuotaExceededFault:
				log.Warn(rds.ErrCodeDBSubnetGroupQuotaExceededFault, aerr.Error())
				return nil, aerr
			case rds.ErrCodeDBSubnetQuotaExceededFault:
				log.Warn(rds.ErrCodeDBSubnetQuotaExceededFault, aerr.Error())
				return nil, aerr
			case rds.ErrCodeDBSubnetGroupDoesNotCoverEnoughAZs:
				log.Warn(rds.ErrCodeDBSubnetGroupDoesNotCoverEnoughAZs, aerr.Error())
				return nil, aerr
			case rds.ErrCodeInvalidSubnet:
				log.Warn(rds.ErrCodeInvalidSubnet, aerr.Error())
				return nil, aerr
			default:
				log.Warn(aerr)
				return nil, aerr
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			log.Warn(err)
			return nil, aerr
		}
	}

	return groupOutput.DBSubnetGroup, nil
}

func findDBCluster(svc *rds.RDS, clusterIdentifier *string) (*rds.DBCluster, error) {
	descClustersInput := &rds.DescribeDBClustersInput{
		DBClusterIdentifier: clusterIdentifier,
	}

	descClusterOuput, err := svc.DescribeDBClusters(descClustersInput)
	if err != nil {
		return nil, err
	}

	return descClusterOuput.DBClusters[0], nil
}

func findDBClusterInstance(svc *rds.RDS, instanceIdentifier *string) (*rds.DBInstance, error) {
	descInstancesInput := &rds.DescribeDBInstancesInput{
		DBInstanceIdentifier: instanceIdentifier,
	}

	descInstancesOuput, err := svc.DescribeDBInstances(descInstancesInput)
	if err != nil {
		return nil, err
	}

	return descInstancesOuput.DBInstances[0], nil
}
