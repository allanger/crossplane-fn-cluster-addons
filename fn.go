package main

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/crossplane/crossplane-fn-cluster-addons/input/v1beta1"
	"github.com/crossplane/crossplane-runtime/pkg/errors"
	"github.com/crossplane/crossplane-runtime/pkg/logging"
	fnv1beta1 "github.com/crossplane/function-sdk-go/proto/v1beta1"
	"github.com/crossplane/function-sdk-go/request"
	"github.com/crossplane/function-sdk-go/response"
)

// Function returns whatever response you ask it to.
type Function struct {
	fnv1beta1.UnimplementedFunctionRunnerServiceServer

	log logging.Logger
}

// RunFunction runs the Function.
func (f *Function) RunFunction(_ context.Context, req *fnv1beta1.RunFunctionRequest) (*fnv1beta1.RunFunctionResponse, error) {
	f.log.Info("Running function", "tag", req.GetMeta().GetTag())

	rsp := response.To(req, response.DefaultTTL)

	in := &v1beta1.Input{}
	if err := request.GetInput(req, in); err != nil {
		response.Fatal(rsp, errors.Wrapf(err, "cannot get Function input from %T", req))
		return rsp, nil
	}
	sess := session.Must(session.NewSession())
	svc := eks.New(sess)
	params := &eks.ListAddonsInput{
		ClusterName: (*string)(&in.Spec.ClusterRef),
		MaxResults: aws.Int64(1),
	}
	res, err := svc.ListAddons(params)
	if err != nil {
		return nil, err
	}

	for _, addon := range res.Addons {
		f.log.Info("Found an addon", "addon", addon)
	}
	// TODO: Add your Function logic here!
	response.Normalf(rsp, "I was run with input %q!", in.Spec.ClusterRef)
	f.log.Info("I was run!", "input", in.Spec.ClusterRef)

	return rsp, nil
}
