# King Family Photos

A Serverless app for distributing digital photo frame content. Built in Golang, CI/CD in GitHub Actions, hosted in AWS and deployed with [serverless](https://www.serverless.com/).

<img src="docs/diagrams/king_family_photos.png" alt="high level infrastructure diagram"
style="max-width:600px;">

1. Photos on home media server are synced to `backup bucket`
1. New photos in bucket trigger `resizePhoto` lambda to copy smaller resolution version of image to `display bucket`
1. Digital photo frame downloads new images every night at `00:00` and restarts
1. Photos removed from `backup bucket` trigger `removePhoto` lambda to remove photo from `display bucket`

## Requirements

- golang
- serverless framework
- AWS account and SDK

## Deployment

The application requires an S3 bucket to consume events from. This bucket exists outside the application in order to avoid accidental deletion. Photos uploaded to this bucket will be ingested by the `resizePhoto` lambda. The bucket should be named `APP_NAME-live-ingest`. App name can be edited in `serverless.yaml`

Application can be deployed in `dev` and `live` environments with the `makefile`

```bash
make deploy-dev

make deploy-live
```

Environments can be torn down with:

```bash
make teardown-dev

make teardown-live
```

## CI/CD

Unit tests, integration tests and deployment can be handled by `GitHub Actions`. To do this, you will need to generate `AWS_ACESS_KEY` and `AWS_SECRET_ACCESS_KEY` for the GitHub service. They should be stored in github as secrets named `AWS_KEY` and `AWS_SECRET` respectively.

## Photo Frame

The photo frame is a raspberry pi and display with raspberry pi os installed, configured to run `scripts/slideshow.sh` on startup. The photo frame should reboot regularly to ensure that photos are up to date.

### Bill of Materials

Matrials listed are those used in original project and provided as a guideline. Any linux computer with display and internet access should be able to run the requred scripts.

- [raspberry pi 3 A+](https://thepihut.com/products/raspberry-pi-3-model-a-plus?variant=13584708763710&currency=GBP&utm_medium=product_sync&utm_source=google&utm_content=sag_organic&utm_campaign=sag_organic&gad_source=1&gclid=CjwKCAjwzIK1BhAuEiwAHQmU3gK1jxzDaA6VRVTLlTprheUi7p_0XxGUF0SXrxza_1SNpiqLLVX6sRoCH4QQAvD_BwE)
- [pibow touchscreen frame](https://shop.pimoroni.com/products/raspberry-pi-7-touchscreen-display-frame?variant=6337432321)
- [raspberry pi display](https://shop.pimoroni.com/products/raspberry-pi-7-touchscreen-display-with-frame?variant=2677960835082)
- micro SD
- micro-usb power supply
- keyboard and mouse for setup

## Display Bucket Permissions

The photo frame will need the `aws-cli` installed, configured with a user that has the
`king-family-photos-live-PhotoBucketSyncPolicy` attached to it.
