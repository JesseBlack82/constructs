package webapp

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscertificatemanager"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscloudfront"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscloudfrontorigins"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3deployment"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type ConstructProps struct {
	awscdk.StackProps
	CertificateArn string
	Aliases        []*string
	SourcePath     string
}

func NewWebAppConstruct(scope constructs.Construct, id *string, props *ConstructProps) constructs.Construct {
	construct := constructs.NewConstruct(scope, id)
	bucket := awss3.NewBucket(scope, jsii.String("EcommercePoc"), &awss3.BucketProps{
		AccessControl: awss3.BucketAccessControl_PRIVATE,
	})
	originIdentity := awscloudfront.NewOriginAccessIdentity(scope, jsii.String("OriginAccessIdentity"), nil)
	bucket.GrantRead(originIdentity, nil)
	distribution := awscloudfront.NewDistribution(scope, jsii.String("EcommercePocDistribution"), &awscloudfront.DistributionProps{
		DefaultRootObject: jsii.String("index.html"),
		Certificate:       awscertificatemanager.Certificate_FromCertificateArn(scope, jsii.String("EcommercePocWebappCertificate"), &props.CertificateArn),
		DomainNames:       &props.Aliases,
		DefaultBehavior: &awscloudfront.BehaviorOptions{
			Origin: awscloudfrontorigins.NewS3Origin(bucket, &awscloudfrontorigins.S3OriginProps{
				OriginAccessIdentity: originIdentity,
			}),
		},
		ErrorResponses: &[]*awscloudfront.ErrorResponse{
			{
				HttpStatus:         jsii.Number(404),
				Ttl:                awscdk.Duration_Seconds(jsii.Number(10)),
				ResponseHttpStatus: jsii.Number(200),
				ResponsePagePath:   jsii.String("/index.html"),
			},
			{
				HttpStatus:         jsii.Number(403),
				Ttl:                awscdk.Duration_Seconds(jsii.Number(10)),
				ResponseHttpStatus: jsii.Number(200),
				ResponsePagePath:   jsii.String("/index.html"),
			},
		},
	})

	awss3deployment.NewBucketDeployment(scope, jsii.String("EcommercePocDeployment"), &awss3deployment.BucketDeploymentProps{
		DestinationBucket: bucket,
		Sources: &[]awss3deployment.ISource{
			awss3deployment.Source_Asset(&props.SourcePath, nil),
		},
		Distribution: distribution,
		DistributionPaths: &[]*string{
			jsii.String("/*"),
		},
	})
	return construct
}
