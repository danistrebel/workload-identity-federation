from google.auth import aws
from google.auth import exceptions
import boto3
import os
from google.auth import environment_vars

class CustomAwsSecurityCredentialsSupplier(aws.AwsSecurityCredentialsSupplier):

    def get_aws_security_credentials(self, context, request):
        aws_credentials = boto3.Session().get_credentials().get_frozen_credentials()

        audience = context.audience
        try:
            return aws.AwsSecurityCredentials(aws_credentials.access_key, aws_credentials.secret_key, aws_credentials.token)
        except Exception as e:
            raise exceptions.RefreshError(e, retryable=True)

    def get_aws_region(self, context, request):
        return "us-east-1"

credentials = aws.Credentials(
    f"//iam.googleapis.com/projects/{os.getenv('GCP_PROJECT_NUMBER')}/locations/global/workloadIdentityPools/{os.getenv('WORKLOAD_IDENTITY_POOL_ID')}/providers/{os.getenv('WORKLOAD_IDENTITY_PROVIDER_ID')}",
    "urn:ietf:params:aws:token-type:aws4_request",
    aws_security_credentials_supplier=CustomAwsSecurityCredentialsSupplier(),
    scopes=['https://www.googleapis.com/auth/cloud-platform']
)