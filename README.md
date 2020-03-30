lunchmoneybridge
============

A utility to bridge Monzo transactions to [Lunch Money](https://lunchmoney.app) ([my referral](https://my.lunchmoney.app/refer/y490cbkx))

lunchmoneybridge comes with an application that is designed to run on [Google Cloud Run](https://cloud.google.com/run), and tools to back-import all your current Monzo transactions into Lunch Money (**tools are due to be published soon**)

Setting up the bridge
=====================

You require:

* a Google (cloud) account
* a Monzo UK account (this is untested on Monzo USA, please open an issue if you try!)
* a clone of this repo locally to build and push to Google Container Registry


1. Head to [https://developers.monzo.com/](https://developers.monzo.com/) and create a new oauth client. Make it confidential and put the redirect URL as `http://127.0.0.1:45679/auth/return`. Note down your Client ID, Owner ID and Client Secret.
2. With your locally cloned repo, you want to build the docker container into Google Container Registry. Firstly, ensure you have [gsutil setup and installed](https://cloud.google.com/storage/docs/gsutil) and run `gcloud builds submit --tag gcr.io/$PROJECT_ID/lunchmoneybridge` in your working directory. This may take a few minutes, but once done you will have the container in your GCR.
3. Head to [Google Cloud Storage](https://console.cloud.google.com/storage) and make a bucket, standard storage class. Note down the name of this bucket (and your region!).
4. Head to [Google Cloud Run](https://console.cloud.google.com/run) and create a service.

* Deployment platform: Cloud run (fully managed)
* Region: whatever you chose for Google Cloud Storage
* Service name: lunchmoneybridge
* Authentication: Allow unauthenticated invocations
* Container image: press select, find the image you built in step 2
* _Show advanced settings_
* Maximum requests per container: 10
* Memory allocated: 128MiB
* Maximum number of instances: 5
Environment variables:

* `BUCKET_NAME` - your Google Storage bucket name
* `MONZO_CLIENT_ID`
* `MONZO_CLIENT_SECRET`
* `LUNCHMONEY_TOKEN` - your Lunch Money API token
* `LUNCHMONEY_ASSET_ID` - your Lunch Money Asset ID which you want to allocate all the transactions to (make a manual one called Monzo)
* `ACCOUNT_ID` - your Monzo Owner ID which we noted down

* Create

5. Once created, click into the service and you'll have a URL at the top. Go to that URL, /auth. This should redirect you to the Monzo authentication. Once you click on the Monzo link in your email, it'll redirect you to a 127.0.0.1:xxxxx page that'll not work as we're hosting this locally. This is OK, simply replace 127.0.0.1:xxxxx with your cloud run URL and it'll process the authentication data.
6. Check your Monzo app and authorize the client
7. Go back to [Monzo Developers](https://developers.monzo.com/) and create a new webhook, pointing to $CLOUD_RUN_URL/webhook



