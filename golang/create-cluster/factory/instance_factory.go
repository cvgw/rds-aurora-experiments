package factory

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/rds"
	log "github.com/sirupsen/logrus"
)

type DBInstanceFactory struct {
	svc                *rds.RDS
	instanceIdentifier *string
	clusterIdentifier  *string
	allocatedStorage   *int64
	engine             *string
	instanceClass      *string
}

func (f *DBInstanceFactory) SetSvc(v *rds.RDS) *DBInstanceFactory {
	f.svc = v
	return f
}

func (f *DBInstanceFactory) SetInstanceIdentifier(v string) *DBInstanceFactory {
	f.instanceIdentifier = aws.String(v)
	return f
}

func (f *DBInstanceFactory) SetClusterIdentifier(v string) *DBInstanceFactory {
	f.clusterIdentifier = aws.String(v)
	return f
}

func (f *DBInstanceFactory) SetAllocatedStorage(v int64) *DBInstanceFactory {
	f.allocatedStorage = aws.Int64(v)
	return f
}

func (f *DBInstanceFactory) SetEngine(v string) *DBInstanceFactory {
	f.engine = aws.String(v)
	return f
}

func (f *DBInstanceFactory) SetInstanceClass(v string) *DBInstanceFactory {
	f.instanceClass = aws.String(v)
	return f
}

func (f *DBInstanceFactory) UpdateOrCreateDBClusterInstance() (*rds.DBInstance, error) {

	instance, err := findDBClusterInstance(f.svc, f.instanceIdentifier)
	if err != nil {
		if err == notFoundErr {
			log.Info("cluster instance does not exist")
			return f.createDBClusterInstance()
		}
	}

	instance, err = f.updateDBInstance(instance)
	if err != nil {
		return nil, err
	}

	return instance, nil
}

func (f *DBInstanceFactory) createDBClusterInstance() (*rds.DBInstance, error) {

	instanceInput := &rds.CreateDBInstanceInput{
		DBInstanceIdentifier: f.instanceIdentifier,
		DBClusterIdentifier:  f.clusterIdentifier,
		Engine:               f.engine,
		DBInstanceClass:      f.instanceClass,
	}

	instanceOutput, err := f.svc.CreateDBInstance(instanceInput)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return instanceOutput.DBInstance, nil
}

func (f *DBInstanceFactory) updateDBInstance(instance *rds.DBInstance) (
	*rds.DBInstance, error,
) {
	input := &rds.ModifyDBInstanceInput{
		//AllocatedStorage:           aws.Int64(10),
		ApplyImmediately: aws.Bool(true),
		//BackupRetentionPeriod:      aws.Int64(1),
		DBInstanceClass:      f.instanceClass,
		DBInstanceIdentifier: instance.DBInstanceIdentifier,
		//MasterUserPassword:         aws.String("mynewpassword"),
		//PreferredBackupWindow:      aws.String("04:00-04:30"),
		//PreferredMaintenanceWindow: aws.String("Tue:05:00-Tue:05:30"),
	}

	result, err := f.svc.ModifyDBInstance(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case rds.ErrCodeInvalidDBInstanceStateFault:
				log.Warn(rds.ErrCodeInvalidDBInstanceStateFault, aerr.Error())
				return nil, aerr
			case rds.ErrCodeInvalidDBSecurityGroupStateFault:
				log.Warn(rds.ErrCodeInvalidDBSecurityGroupStateFault, aerr.Error())
				return nil, aerr
			case rds.ErrCodeDBInstanceAlreadyExistsFault:
				log.Warn(rds.ErrCodeDBInstanceAlreadyExistsFault, aerr.Error())
				return nil, aerr
			case rds.ErrCodeDBInstanceNotFoundFault:
				log.Warn(rds.ErrCodeDBInstanceNotFoundFault, aerr.Error())
				return nil, aerr
			case rds.ErrCodeDBSecurityGroupNotFoundFault:
				log.Warn(rds.ErrCodeDBSecurityGroupNotFoundFault, aerr.Error())
				return nil, aerr
			case rds.ErrCodeDBParameterGroupNotFoundFault:
				log.Warn(rds.ErrCodeDBParameterGroupNotFoundFault, aerr.Error())
				return nil, aerr
			case rds.ErrCodeInsufficientDBInstanceCapacityFault:
				log.Warn(rds.ErrCodeInsufficientDBInstanceCapacityFault, aerr.Error())
				return nil, aerr
			case rds.ErrCodeStorageQuotaExceededFault:
				log.Warn(rds.ErrCodeStorageQuotaExceededFault, aerr.Error())
				return nil, aerr
			case rds.ErrCodeInvalidVPCNetworkStateFault:
				log.Warn(rds.ErrCodeInvalidVPCNetworkStateFault, aerr.Error())
				return nil, aerr
			case rds.ErrCodeProvisionedIopsNotAvailableInAZFault:
				log.Warn(rds.ErrCodeProvisionedIopsNotAvailableInAZFault, aerr.Error())
				return nil, aerr
			case rds.ErrCodeOptionGroupNotFoundFault:
				log.Warn(rds.ErrCodeOptionGroupNotFoundFault, aerr.Error())
				return nil, aerr
			case rds.ErrCodeDBUpgradeDependencyFailureFault:
				log.Warn(rds.ErrCodeDBUpgradeDependencyFailureFault, aerr.Error())
				return nil, aerr
			case rds.ErrCodeStorageTypeNotSupportedFault:
				log.Warn(rds.ErrCodeStorageTypeNotSupportedFault, aerr.Error())
				return nil, aerr
			case rds.ErrCodeAuthorizationNotFoundFault:
				log.Warn(rds.ErrCodeAuthorizationNotFoundFault, aerr.Error())
				return nil, aerr
			case rds.ErrCodeCertificateNotFoundFault:
				log.Warn(rds.ErrCodeCertificateNotFoundFault, aerr.Error())
				return nil, aerr
			case rds.ErrCodeDomainNotFoundFault:
				log.Warn(rds.ErrCodeDomainNotFoundFault, aerr.Error())
				return nil, aerr
			case rds.ErrCodeBackupPolicyNotFoundFault:
				log.Warn(rds.ErrCodeBackupPolicyNotFoundFault, aerr.Error())
				return nil, aerr
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

	return result.DBInstance, nil
}
