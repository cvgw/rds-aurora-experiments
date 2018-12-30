package request

import (
	"os"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
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

type ClusterRequest struct {
	Region             string
	Profile            string
	InstanceIdentifier string
	InstanceClass      string
	ClusterId          string
	Engine             string
	EngineVersion      string
	MasterUsername     string
	MasterUserPass     string
	GroupDescription   string
	GroupName          string
	ReadyTimeout       int
	SgIds              []string
	Subnets            []string
}

func NewRequest() ClusterRequest {
	req := ClusterRequest{}
	req.Region = os.Getenv(awsRegionVar)
	req.Profile = os.Getenv(awsProfileVar)
	req.InstanceIdentifier = os.Getenv(instanceIdVar)
	req.InstanceClass = os.Getenv(instanceClassVar)
	req.ClusterId = os.Getenv(clusterIdVar)
	req.Engine = os.Getenv(engineVar)
	req.EngineVersion = os.Getenv(engineVersionVar)
	req.MasterUsername = os.Getenv(masterUsernameVar)
	req.MasterUserPass = os.Getenv(masterUserPassVar)
	req.GroupDescription = os.Getenv(groupDescriptionVar)
	req.GroupName = os.Getenv(groupNameVar)

	readyTimeout := os.Getenv(readyTimeoutVar)

	sgIds := make([]string, 0)
	subnets := make([]string, 0)

	for _, sgI := range strings.Split(os.Getenv(sgIdsVar), ",") {
		sgIds = append(sgIds, sgI)
	}
	for _, subnet := range strings.Split(os.Getenv(subnetsVar), ",") {
		subnets = append(subnets, subnet)
	}

	req.SgIds = sgIds
	req.Subnets = subnets

	rTimeout, err := strconv.Atoi(readyTimeout)
	if err != nil {
		log.Warn(err)
		rTimeout = 1
	}
	req.ReadyTimeout = rTimeout

	return req
}
