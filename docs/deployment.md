# 🌐 Deployment

The backend is deployed to a [AWS EC2 instance](https://aws.amazon.com/ec2/) with a reverse-proxy using [nginx](https://www.nginx.com).

The EC2 instance is allocated an elastic IP that is routed to by [Route 53](https://aws.amazon.com/route53/).

Signed SSL certificate for the subdomain is provided by [Let's Encrypt](https://letsencrypt.org/).

_To get started with deploying your app following this configuration, you can try giving the links below a read._

## Useful links:

- [REST Api with EC2](https://docs.aws.amazon.com/AWSEC2/latest/APIReference/making-api-requests.html#using-libraries)
- [Reverse proxy architecture ](https://aws.amazon.com/blogs/architecture/serving-content-using-fully-managed-reverse-proxy-architecture/)
