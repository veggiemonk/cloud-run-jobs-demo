gcloud iam service-accounts create "${SERVICE_NAME}-workflow" \
    --description="cloud workflow for ${SERVICE_NAME}" \
    --display-name="${SERVICE_NAME}-workflow"

gcloud projects add-iam-policy-binding "${SERVICE_NAME}-workflow" \
    --member="serviceAccount:${SERVICE_NAME}-workflow@$PROJECT_ID.iam.gserviceaccount.com" \
    --role="roles/run.invoker"

gcloud workflows deploy "${SERVICE_NAME}-workflow" \
      --location="${REGION}" \
      --source=workflow.yaml \
      --service-account"=${SERVICE_NAME}-workflow@${PROJECT_ID}.iam.gserviceaccount.com"
