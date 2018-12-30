package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/cvgw/rds-aurora-experiments/golang/create-cluster/request"
	"github.com/cvgw/rds-aurora-experiments/golang/create-cluster/service"
)

func main() {
	req := request.NewRequest()

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config:  aws.Config{Region: aws.String(req.Region)},
		Profile: req.Profile,
	}))
	svc := rds.New(sess)

	err := service.HandleRequest(svc, req)
	if err != nil {
		log.Fatal(err)
	}

	log.Info("success")
}
