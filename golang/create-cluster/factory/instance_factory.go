package factory

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/rds"
	log "github.com/sirupsen/logrus"
)

func NewDBInstanceFactory(cluster *rds.DBCluster, id, class string) *dbInstanceFactory {
	f := &dbInstanceFactory{}

	f.engine = cluster.Engine
	f.instanceIdentifier = aws.String(id)
	f.instanceClass = aws.String(class)

	f.clusterIdentifier = cluster.DBClusterIdentifier

	return f
}

type dbInstanceFactory struct {
	instanceIdentifier *string
	clusterIdentifier  *string
	engine             *string
	instanceClass      *string
}

func (f *dbInstanceFactory) FindOrCreateDBClusterInstance(svc *rds.RDS) (*rds.DBInstance, error) {

	instance, err := findDBClusterInstance(svc, f.instanceIdentifier)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case rds.ErrCodeDBInstanceNotFoundFault:
				log.Info(rds.ErrCodeDBInstanceNotFoundFault, aerr.Error())
				return f.createDBClusterInstance(svc)
			default:
				log.Warn(aerr.Error())
				return nil, aerr
			}
		} else {
			log.Warn(err)
			return nil, err
		}
	}

	return instance, nil
}

func (f *dbInstanceFactory) createDBClusterInstance(svc *rds.RDS) (*rds.DBInstance, error) {

	instanceInput := &rds.CreateDBInstanceInput{
		DBInstanceIdentifier: f.instanceIdentifier,
		DBClusterIdentifier:  f.clusterIdentifier,
		Engine:               f.engine,
		DBInstanceClass:      f.instanceClass,
	}

	instanceOutput, err := svc.CreateDBInstance(instanceInput)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return instanceOutput.DBInstance, nil
}
