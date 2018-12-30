package factory

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/rds"
	log "github.com/sirupsen/logrus"
)

type NewDBClusterFactoryInput struct {
	ClusterId        string
	Engine           string
	EngineVersion    string
	MasterUsername   string
	MasterUserPass   string
	SecurityGroupIds []string
	SubnetGroupName  *string
}

func NewDBClusterFactory(input NewDBClusterFactoryInput) *dbClusterFactory {
	f := &dbClusterFactory{}

	f.clusterIdentifier = aws.String(input.ClusterId)
	f.engine = aws.String(input.Engine)
	f.engineVersion = aws.String(input.EngineVersion)
	f.masterUsername = aws.String(input.MasterUsername)
	f.masterUserPass = aws.String(input.MasterUserPass)

	f.subnetGroupName = input.SubnetGroupName

	sIds := make([]*string, 0)
	for _, i := range input.SecurityGroupIds {
		sIds = append(sIds, aws.String(i))
	}
	f.securityGroupIds = sIds

	return f
}

func (f *dbClusterFactory) UpdateOrCreateDBCluster(svc *rds.RDS) (*rds.DBCluster, error) {
	dbCluster, err := findDBCluster(svc, f.clusterIdentifier)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case rds.ErrCodeDBClusterNotFoundFault:
				log.Info(rds.ErrCodeDBClusterNotFoundFault, aerr.Error())
				return f.createDBCluster(svc)
			default:
				log.Warn(aerr.Error())
				return nil, aerr
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			log.Warn(err.Error())
			return nil, err
		}
	}

	dbCluster, err = f.updateDBCluster(svc, dbCluster)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case rds.ErrCodeDBClusterNotFoundFault:
				log.Warn(rds.ErrCodeDBClusterNotFoundFault, aerr.Error())
				return nil, aerr
			case rds.ErrCodeInvalidDBClusterStateFault:
				log.Warn(rds.ErrCodeInvalidDBClusterStateFault, aerr.Error())
				return nil, aerr
			case rds.ErrCodeStorageQuotaExceededFault:
				log.Warn(rds.ErrCodeStorageQuotaExceededFault, aerr.Error())
				return nil, aerr
			case rds.ErrCodeDBSubnetGroupNotFoundFault:
				log.Warn(rds.ErrCodeDBSubnetGroupNotFoundFault, aerr.Error())
				return nil, aerr
			case rds.ErrCodeInvalidVPCNetworkStateFault:
				log.Warn(rds.ErrCodeInvalidVPCNetworkStateFault, aerr.Error())
				return nil, aerr
			case rds.ErrCodeInvalidDBSubnetGroupStateFault:
				log.Warn(rds.ErrCodeInvalidDBSubnetGroupStateFault, aerr.Error())
				return nil, aerr
			case rds.ErrCodeInvalidSubnet:
				log.Warn(rds.ErrCodeInvalidSubnet, aerr.Error())
				return nil, aerr
			case rds.ErrCodeDBClusterParameterGroupNotFoundFault:
				log.Warn(rds.ErrCodeDBClusterParameterGroupNotFoundFault, aerr.Error())
				return nil, aerr
			case rds.ErrCodeInvalidDBSecurityGroupStateFault:
				log.Warn(rds.ErrCodeInvalidDBSecurityGroupStateFault, aerr.Error())
				return nil, aerr
			case rds.ErrCodeInvalidDBInstanceStateFault:
				log.Warn(rds.ErrCodeInvalidDBInstanceStateFault, aerr.Error())
				return nil, aerr
			case rds.ErrCodeDBClusterAlreadyExistsFault:
				log.Warn(rds.ErrCodeDBClusterAlreadyExistsFault, aerr.Error())
				return nil, aerr
			default:
				log.Warn(aerr.Error())
				return nil, aerr
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			log.Warn(err)
			return nil, err
		}
	}

	return dbCluster, nil
}

type dbClusterFactory struct {
	clusterIdentifier *string
	subnetGroupName   *string
	securityGroupIds  []*string
	engine            *string
	engineVersion     *string
	masterUsername    *string
	masterUserPass    *string
}

func (f *dbClusterFactory) createDBCluster(svc *rds.RDS) (*rds.DBCluster, error) {
	clusterInput := &rds.CreateDBClusterInput{
		DBClusterIdentifier: f.clusterIdentifier,
		Engine:              f.engine,
		EngineVersion:       f.engineVersion,
		MasterUsername:      f.masterUsername,
		MasterUserPassword:  f.masterUserPass,
		DBSubnetGroupName:   f.subnetGroupName,
		VpcSecurityGroupIds: f.securityGroupIds,
	}

	clusterOutput, err := svc.CreateDBCluster(clusterInput)
	if err != nil {
		log.Warn(err)
		return nil, err
	}

	return clusterOutput.DBCluster, nil
}

func (f *dbClusterFactory) updateDBCluster(svc *rds.RDS, dbCluster *rds.DBCluster) (*rds.DBCluster, error) {
	input := &rds.ModifyDBClusterInput{
		ApplyImmediately:    aws.Bool(true),
		DBClusterIdentifier: dbCluster.DBClusterIdentifier,
		MasterUserPassword:  f.masterUserPass,
		VpcSecurityGroupIds: f.securityGroupIds,
	}

	if *dbCluster.EngineVersion != *f.engineVersion {
		input.EngineVersion = f.engineVersion
	}

	result, err := svc.ModifyDBCluster(input)
	if err != nil {
		return nil, err
	}

	return result.DBCluster, nil
}
