package factory

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/rds"
	log "github.com/sirupsen/logrus"
)

const (
	waitSleepTime = 10
	requiredReady = 4
)

func WaitForClusterReady(ctx context.Context, svc *rds.RDS, cluster *rds.DBCluster) bool {
	var readyCount int

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
				log.Infof("cluster ready test %d/%d", readyCount+1, requiredReady)
				readyCount++
			} else {
				readyCount = 0
				log.Infof("cluster not ready: status %s", *dbCluster.Status)
			}

			if readyCount == requiredReady {
				log.Info("cluster ready and stable")
				return true
			}

			time.Sleep(waitSleepTime * time.Second)
		}
	}
}

func WaitForInstanceReady(ctx context.Context, svc *rds.RDS, instance *rds.DBInstance) bool {
	var readyCount int

	identifier := instance.DBInstanceIdentifier

	for {
		select {
		case <-ctx.Done():
			log.Warn("context expired")
			return false
		default:
			instance, err := findDBClusterInstance(svc, identifier)
			if err != nil {
				return false
			}

			if *instance.DBInstanceStatus == "available" {
				log.Infof("instance ready test %d/%d", readyCount+1, requiredReady)
				readyCount++
			} else {
				readyCount = 0
				log.Infof("instance not ready: status %s", *instance.DBInstanceStatus)
			}

			if readyCount == requiredReady {
				log.Info("instance ready and stable")
				return true
			}

			time.Sleep(waitSleepTime * time.Second)
		}
	}
}
