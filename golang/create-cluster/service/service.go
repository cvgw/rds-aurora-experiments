package service

import (
	"context"
	"errors"
	"time"

	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/cvgw/rds-aurora-experiments/golang/create-cluster/factory"
	"github.com/cvgw/rds-aurora-experiments/golang/create-cluster/request"
	log "github.com/sirupsen/logrus"
)

func HandleRequest(svc *rds.RDS, req request.ClusterRequest) error {
	dbSubnetGroup, err := factory.UpdateOrCreateDBSubnetGroup(
		svc,
		req.GroupName,
		req.GroupDescription,
		req.Subnets,
	)
	if err != nil {
		return err
	}
	log.Info(dbSubnetGroup)

	clusterFactoryInput := factory.NewDBClusterFactoryInput{
		ClusterId:        req.ClusterId,
		Engine:           req.Engine,
		EngineVersion:    req.EngineVersion,
		MasterUsername:   req.MasterUsername,
		MasterUserPass:   req.MasterUserPass,
		SecurityGroupIds: req.SgIds,
		SubnetGroupName:  dbSubnetGroup.DBSubnetGroupName,
	}

	cluster, err := updateOrCreateCluster(svc, clusterFactoryInput, req.ReadyTimeout)
	if err != nil {
		return err
	}

	instanceFactory := factory.DBInstanceFactory{}
	instanceFactory.SetSvc(svc).
		SetInstanceIdentifier(req.InstanceIdentifier).
		SetClusterIdentifier(*cluster.DBClusterIdentifier).
		SetEngine(*cluster.Engine).
		SetInstanceClass(req.InstanceClass)

	_, err = updateOrCreateInstance(instanceFactory, req.ReadyTimeout, svc)
	if err != nil {
		return err
	}
	return nil
}

func updateOrCreateCluster(svc *rds.RDS, input factory.NewDBClusterFactoryInput, rTimeout int) (*rds.DBCluster, error) {
	clusterFactory := factory.NewDBClusterFactory(input)
	cluster, err := clusterFactory.UpdateOrCreateDBCluster(svc)
	if err != nil {
		return nil, err
	}
	log.Info(cluster)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(rTimeout)*time.Minute)
	defer cancel()
	ready := factory.WaitForClusterReady(ctx, svc, cluster)

	if !ready {
		return nil, errors.New("cluster not ready within timeout")
	}

	return cluster, nil
}

func updateOrCreateInstance(
	f factory.DBInstanceFactory, rTimeout int, svc *rds.RDS,
) (*rds.DBInstance, error) {
	instance, err := f.UpdateOrCreateDBClusterInstance()
	if err != nil {
		return nil, err
	}
	log.Info(instance)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(rTimeout)*time.Minute)
	defer cancel()
	ready := factory.WaitForInstanceReady(ctx, svc, instance)
	if !ready {
		return nil, errors.New("instance not ready within timeout")
	}

	return instance, nil
}
