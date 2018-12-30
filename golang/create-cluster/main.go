package main

import (
	"context"
	"os"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/cvgw/rds-aurora-experiments/golang/create-cluster/factory"
)

const (
	groupDescriptionVar = "SUBNET_GROUP_DESCRIPTION"
	groupNameVar        = "SUBNET_GROUP_NAME"
	engineVar           = "ENGINE"
	engineVersionVar    = "ENGINE_VERSION"
	masterUsernameVar   = "MASTER_USERNAME"
	masterUserPassVar   = "MASTER_USER_PASSWORD"
	clusterIdVar        = "CLUSTER_ID"
	awsRegionVar        = "AWS_REGION"
	awsProfileVar       = "AWS_PROFILE"
	readyTimeoutVar     = "READY_TIMEOUT_MINUTES"
	instanceIdVar       = "INSTANCE_ID"
	instanceClassVar    = "INSTANCE_CLASS"
	sgIdsVar            = "SECURITY_GROUP_IDS"
	subnetsVar          = "SUBNETS"
)

func main() {
	region := os.Getenv(awsRegionVar)
	profile := os.Getenv(awsProfileVar)
	instanceIdentifier := os.Getenv(instanceIdVar)
	instanceClass := os.Getenv(instanceClassVar)
	clusterId := os.Getenv(clusterIdVar)
	engine := os.Getenv(engineVar)
	engineVersion := os.Getenv(engineVersionVar)
	masterUsername := os.Getenv(masterUsernameVar)
	masterUserPass := os.Getenv(masterUserPassVar)
	groupDescription := os.Getenv(groupDescriptionVar)
	groupName := os.Getenv(groupNameVar)
	readyTimeout := os.Getenv(readyTimeoutVar)
	sgIds := make([]string, 0)
	subnets := make([]string, 0)

	for _, sgI := range strings.Split(os.Getenv(sgIdsVar), ",") {
		sgIds = append(sgIds, sgI)
	}
	for _, subnet := range strings.Split(os.Getenv(subnetsVar), ",") {
		subnets = append(subnets, subnet)
	}

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config:  aws.Config{Region: aws.String(region)},
		Profile: profile,
	}))
	svc := rds.New(sess)

	dbSubnetGroup, err := factory.FindOrCreateDBSubnetGroup(svc, groupName, groupDescription, subnets)
	if err != nil {
		log.Fatal(err)
	}
	log.Info(dbSubnetGroup)

	clusterFactoryInput := factory.NewDBClusterFactoryInput{
		ClusterId:        clusterId,
		Engine:           engine,
		EngineVersion:    engineVersion,
		MasterUsername:   masterUsername,
		MasterUserPass:   masterUserPass,
		SecurityGroupIds: sgIds,
		SubnetGroupName:  dbSubnetGroup.DBSubnetGroupName,
	}

	clusterFactory := factory.NewDBClusterFactory(clusterFactoryInput)
	cluster, err := clusterFactory.FindOrCreateDBCluster(svc)
	if err != nil {
		log.Fatal(err)
	}
	log.Info(cluster)

	rTimeout, err := strconv.Atoi(readyTimeout)
	if err != nil {
		log.Warn(err)
		rTimeout = 1
	}

	func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(rTimeout)*time.Minute)
		defer cancel()
		ready := factory.WaitForClusterReady(ctx, svc, cluster)

		if !ready {
			log.Fatal("cluster not ready within timeout")
		}
	}()

	instanceFactory := factory.NewDBInstanceFactory(cluster, instanceIdentifier, instanceClass)
	instance, err := instanceFactory.FindOrCreateDBClusterInstance(svc)
	if err != nil {
		log.Fatal(err)
	}
	log.Info(instance)

	func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(rTimeout)*time.Minute)
		defer cancel()
		ready := factory.WaitForInstanceReady(ctx, svc, instance)
		if !ready {
			log.Fatal("instance not ready within timeout")
		}
	}()

	log.Info("success")
}
