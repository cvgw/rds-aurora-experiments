package factory

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/rds"
	log "github.com/sirupsen/logrus"
)

func WaitForClusterReady(ctx context.Context, svc *rds.RDS, cluster *rds.DBCluster) bool {
	clusterIdentifier := cluster.DBClusterIdentifier
	for {
		select {
		case <-ctx.Done():
			log.Warn("context expired")
			return false
		default:
			dbCluster, err := findDBCluster(svc, clusterIdentifier)
			if err != nil {
				if aerr, ok := err.(awserr.Error); ok {
					switch aerr.Code() {
					case rds.ErrCodeDBClusterNotFoundFault:
						log.Info(rds.ErrCodeDBClusterNotFoundFault, aerr.Error())
						return false
					default:
						log.Warn(aerr.Error())
						return false
					}
				} else {
					// Print the error, cast err to awserr.Error to get the Code and
					// Message from an error.
					log.Warn(err.Error())
					return false
				}
			}

			if *dbCluster.Status == "available" {
				log.Info("cluster ready")
				return true
			}

			log.Infof("cluster not ready: status %s", *dbCluster.Status)
			time.Sleep(30 * time.Second)
		}
	}
}

func WaitForInstanceReady(ctx context.Context, svc *rds.RDS, instance *rds.DBInstance) bool {
	identifier := instance.DBInstanceIdentifier
	for {
		select {
		case <-ctx.Done():
			log.Warn("context expired")
			return false
		default:
			instance, err := findDBClusterInstance(svc, identifier)
			if err != nil {
				if aerr, ok := err.(awserr.Error); ok {
					switch aerr.Code() {
					case rds.ErrCodeDBInstanceNotFoundFault:
						log.Info(rds.ErrCodeDBInstanceNotFoundFault, aerr.Error())
						return false
					default:
						log.Warn(aerr.Error())
						return false
					}
				} else {
					// Print the error, cast err to awserr.Error to get the Code and
					// Message from an error.
					log.Warn(err.Error())
					return false
				}
			}

			if *instance.DBInstanceStatus == "available" {
				log.Info("instance ready")
				return true
			}

			log.Infof("instance not ready: status %s", *instance.DBInstanceStatus)
			time.Sleep(30 * time.Second)
		}
	}
}
