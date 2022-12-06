gcloud artifacts repositories create "${SERVICE_NAME}" \
  --repository-format=docker \
  --location="$REGION" \
  --description="Docker repository for ${SERVICE_NAME}"

# Artifact repo: ${_REGION}-docker.pkg.dev/${PROJECT_ID}/${_SERVICE_NAME}
# image:         ${_REGION}-docker.pkg.dev/${PROJECT_ID}/${_SERVICE_NAME}/${_SERVICE_NAME}:${_TAG}

gcloud iam service-accounts create "${SERVICE_NAME}-build" \
    --description="cloud workflow for ${SERVICE_NAME}" \
    --display-name="${SERVICE_NAME}-build"

# + IAM Roles ....

gcloud beta builds triggers import --source=build-trigger.yaml
