# Accessing GCP Services from AWS ECS

The current workload identity federation implementation for AWS in the [Google Auth Python Library](https://google-auth.readthedocs.io/en/master/) is based on the EC2 metadata server.

## Python

To make workload identity federation work in ECS, we need to apply a few tweaks as shown in [gcp_aws_credentials.py](./python/gcp_aws_credentials.py) to customize how the AWS token and client ID are obtained.

This workaround is following the guidance from the auth library maintainers as mentioned in [this comment](https://github.com/googleapis/google-auth-library-python/pull/1556#issuecomment-2259334622).

## Authenticating the ECS service in GCP IAM

Most likely you'll want to authenticate the AWS service in GCP via the task role:

```sh
principalSet://iam.googleapis.com/projects/{GCP_PROJECT_NUMBER}/locations/global/workloadIdentityPools/{WORKLOAD_IDENTITY_POOL_ID}/attribute.aws_role/arn:aws:sts::{AWS_PROJECT_ID}:assumed-role/{ECS_TASK_ROLE}
```

You could also use the individual task run ID

```sh
principal://iam.googleapis.com/projects/{GCP_PROJECT_NUMBER}/locations/global/workloadIdentityPools/{WORKLOAD_IDENTITY_POOL_ID}/subject/arn:aws:sts::{AWS_PROJECT_ID}:assumed-role/{ECS_TASK_ROLE}/{ECS_TASK_ID}
```

## Disclaimer

This code is provided as a proof of concept and is not optimized for production use. 

In a production environment, consider implementing:

* **Caching of access credentials:** The current implementation fetches new credentials for every request. Caching credentials would significantly improve performance.
* **Error handling and retries:** Robust error handling and retry mechanisms are crucial for production-level reliability.
* **Security best practices:** Review and implement appropriate security measures for handling credentials and accessing GCP services.

This code should not be used in production without careful consideration and implementation of the above points.

