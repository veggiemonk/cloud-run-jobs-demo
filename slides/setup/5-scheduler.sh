gcloud iam service-accounts create "${SERVICE_NAME}-scheduler" \
    --description="cloud scheduler for ${SERVICE_NAME}" \
    --display-name="${SERVICE_NAME}-scheduler"

gcloud projects add-iam-policy-binding "${PROJECT_ID}" \
      --member "serviceAccount:${SERVICE_NAME}-scheduler@${PROJECT_ID}.iam.gserviceaccount.com" \
      --role "roles/workflow.invoker"

gcloud projects add-iam-policy-binding "${PROJECT_ID}" \
      --member "serviceAccount:${SERVICE_NAME}-scheduler@${PROJECT_ID}.iam.gserviceaccount.com" \
      --role "roles/logging.logWriter"


# At 06:00 on every 7th day-of-month.
# https://crontab.guru/#0_6_*/6_*_*
gcloud scheduler jobs create http "${SERVICE_NAME}-scheduler" \
   --location europe-west3 \ # There is no cloud scheduler in eu-north1 :(
   --schedule="0 6 */6 * *" \
   --uri="https://workflowexecutions.googleapis.com/v1/projects/${PROJECT}/locations/${REGION_WIDE}/workflows/${SERVICE_NAME}-scheduler/executions" \
   --http-method POST \
   --oauth-service-account-email "${SERVICE_NAME}-scheduler@${PROJECT_ID}.iam.gserviceaccount.com"

