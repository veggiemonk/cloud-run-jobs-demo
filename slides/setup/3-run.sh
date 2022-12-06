gcloud iam service-accounts create "${SERVICE_NAME}-jobs" \
    --description="cloud run jobs for ${SERVICE_NAME}" \
    --display-name="${SERVICE_NAME}-jobs"

# + IAM Roles ....

gcloud beta run jobs create "${SERVICE_NAME}" \
  --image "${ARTIFACT}/${SERVICE_NAME}:${SHORT_SHA}" \
  --tasks 4 \
  --task-timeout 3600s \
  --service-account "${SERVICE_NAME}-jobs@${PROJECT_ID}.iam.gserviceaccount.com" \
  --memory 2Gi \
  --cpu 2 \
  --max-retries 2 \
  --parallelism 2 \
  --region "${REGION}" \
  --set-secrets GITHUB_TOKEN="${SERVICE_NAME}-token:latest" \
  --labels env=prod,app="${SERVICE_NAME}"
