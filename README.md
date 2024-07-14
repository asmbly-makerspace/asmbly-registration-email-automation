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
