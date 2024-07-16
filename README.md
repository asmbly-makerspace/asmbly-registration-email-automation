# Asmbly Class Registration Email Automation

This app is deployed as an AWS Lambda using the AWS CDK.

Requests to the lambda will search for an email template in the mail service
(currently Mailjet) corresponding to the event registration class name. If a
matching template is found, the lambda will send an email with that template to
the registrant.

The lambda is triggered by an event registration webhook in Neon.

## Useful commands

- `cdk deploy` deploy this stack to your default AWS account/region
- `cdk diff` compare deployed stack with current state
- `cdk synth` emits the synthesized CloudFormation template
- `go test` run unit tests

## Integration tests

### Installing AWS RIE

To run the lambda function locally, you will need to install the AWS runtime
interface emulator locally. Run the following command:

```
mkdir -p ~/.aws-lambda-rie && \
    curl -Lo ~/.aws-lambda-rie/aws-lambda-rie https://github.com/aws/aws-lambda-runtime-interface-emulator/releases/latest/download/aws-lambda-rie-arm64 && \
    chmod +x ~/.aws-lambda-rie/aws-lambda-rie
```

### Build the lambda container image

Build the lambda container image to run locally:

```
docker build --platform=linux/arm64 --load --no-cache -t emails-lambda-image:latest ./lambda-image
```

### Run the lambda container locally

Next, start a local container (from the project root directory containing the
.env file) running the lambda function:

```
docker run --env-file ./.env --platform linux/arm64 -d -v ~/.aws-lambda-rie:/aws-lambda -p 9000:8080 \    --entrypoint /aws-lambda/aws-lambda-rie \    emails-lambda-image:latest \        /app
```

### Test the function

The lambda function can now be tested locally at
http://localhost:9000/2015-03-31/functions/function/invocations

Hit the local endpoint with test events, e.g. :

```
curl "http://localhost:9000/2015-03-31/functions/function/invocations" -d '{"headers":"hello world!", "body": "testData"}'
```

In the future, these integration test events will run as Go tests
