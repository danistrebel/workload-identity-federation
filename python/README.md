# Python Example

This example illustrates how to use the Google Authentication Library for Python and provides a workaround for using workload identity federation with ECS on AWS Fargate.
The ECS workaround is automatically enabled if the application runs in a Fargate context if the `AWS_EXECUTION_ENV` is set to `AWS_ECS_FARGATE`.

## Run it locally

```sh
git clone https://github.com/danistrebel/workload-identity-federation.git
cd workload-identity-federation/python
pip install -r requirements.txt
python app.py
```