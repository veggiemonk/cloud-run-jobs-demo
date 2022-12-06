#!/usr/bin/env bash

#
# This script is more of a template than a script that you can run.
#

set -x


PROJECT_ID=${PROJECT_ID:-"CHANGE_ME"}
SERVICE_NAME="github-star"
REGION=${REGION:-"eu-north1"}
REGION_WIDE=${REGION_WIDE:-"europe-north1"}


#       .o.       ooooooooo.   ooooo
#      .888.      `888   `Y88. `888'
#     .8"888.      888   .d88'  888
#    .8' `888.     888ooo88P'   888
#   .88ooo8888.    888          888
#  .8'     `888.   888          888
# o88o     o8888o o888o        o888o
#
gcloud services enable \
    artifactregistry.googleapis.com \
    cloudbuild.googleapis.com \
    cloudscheduler.googleapis.com \
    containerregistry.googleapis.com \
    run.googleapis.com \
    secretmanager.googleapis.com \
    workflows.googleapis.com

# oooooooooo.  ooooo     ooo ooooo ooooo        oooooooooo.
# `888'   `Y8b `888'     `8' `888' `888'        `888'   `Y8b
#  888     888  888       8   888   888          888      888
#  888oooo888'  888       8   888   888          888      888
#  888    `88b  888       8   888   888          888      888
#  888    .88P  `88.    .8'   888   888       o  888     d88'
# o888bood8P'     `YbodP'    o888o o888ooooood8 o888bood8P'

gcloud artifacts repositories create "${SERVICE_NAME}" \
  --repository-format=docker \
  --location="$REGION" \
  --description="Docker repository for ${SERVICE_NAME}"


gcloud iam service-accounts create "${SERVICE_NAME}-build" \
    --description="cloud workflow for ${SERVICE_NAME}" \
    --display-name="${SERVICE_NAME}-build"

gcloud projects add-iam-policy-binding "${PROJECT_ID}" \
    --member="serviceAccount:${SERVICE_NAME}-build@${PROJECT_ID}.iam.gserviceaccount.com" \
    --role="roles/artifactregistry.repoAdmin"

gcloud projects add-iam-policy-binding "${PROJECT_ID}" \
    --member="serviceAccount:${SERVICE_NAME}-build@${PROJECT_ID}.iam.gserviceaccount.com" \
    --role="roles/run.admin"

gcloud projects add-iam-policy-binding "${PROJECT_ID}" \
    --member="serviceAccount:${SERVICE_NAME}-build@${PROJECT_ID}.iam.gserviceaccount.com" \
    --role="roles/iam.serviceAccountUser"

gcloud projects add-iam-policy-binding "${PROJECT_ID}" \
    --member="serviceAccount:${SERVICE_NAME}-build@${PROJECT_ID}.iam.gserviceaccount.com" \
    --role="roles/cloudbuild.builds.builder"

gcloud projects add-iam-policy-binding "${PROJECT_ID}" \
    --member="serviceAccount:${SERVICE_NAME}-build@${PROJECT_ID}.iam.gserviceaccount.com" \
    --role="roles/logging.logWriter"

gcloud projects add-iam-policy-binding "${PROJECT_ID}" \
    --member="serviceAccount:${SERVICE_NAME}-build@${PROJECT_ID}.iam.gserviceaccount.com" \
    --role="roles/workflows.editor"

gcloud beta builds triggers import --source=build-trigger.yaml


# ooooooooo.   ooooo     ooo ooooo      ooo
# `888   `Y88. `888'     `8' `888b.     `8'
#  888   .d88'  888       8   8 `88b.    8
#  888ooo88P'   888       8   8   `88b.  8
#  888`88b.     888       8   8     `88b.8
#  888  `88b.   `88.    .8'   8       `888
# o888o  o888o    `YbodP'    o8o        `8

# create service accounts for cloud run jobs

gcloud iam service-accounts create "${SERVICE_NAME}-jobs" \
    --description="cloud run jobs for ${SERVICE_NAME}" \
    --display-name="${SERVICE_NAME}-jobs"

# set permissions for service account

gcloud projects add-iam-policy-binding "${PROJECT_ID}" \
    --member="serviceAccount:${SERVICE_NAME}-jobs@${PROJECT_ID}.iam.gserviceaccount.com" \
    --role="roles/secretmanager.secretAccessor"

gcloud projects add-iam-policy-binding "${PROJECT_ID}" \
    --member="serviceAccount:${SERVICE_NAME}-jobs@${PROJECT_ID}.iam.gserviceaccount.com" \
    --role="roles/logging.logWriter"

gcloud beta run jobs create "${SERVICE_NAME}" \
  --image "${ARTIFACT}/${SERVICE_NAME}:${SHORT_SHA}" \
  --tasks 4 \
  --task-timeout 3600s \
  --service-account "${SERVICE_NAME}-jobs@h${PROJECT_ID}.iam.gserviceaccount.com" \
  --memory 2Gi \
  --cpu 2 \
  --max-retries 2 \
  --parallelism 2 \
  --region "${REGION}" \
  --set-secrets GITHUB_TOKEN="${SERVICE_NAME}-token:latest" \
  --labels env=prod,app="${SERVICE_NAME}"

# oooooo   oooooo     oooo   .oooooo.   ooooooooo.   oooo    oooo oooooooooooo ooooo          .oooooo.   oooooo   oooooo     oooo
#  `888.    `888.     .8'   d8P'  `Y8b  `888   `Y88. `888   .8P'  `888'     `8 `888'         d8P'  `Y8b   `888.    `888.     .8'
#   `888.   .8888.   .8'   888      888  888   .d88'  888  d8'     888          888         888      888   `888.   .8888.   .8'
#    `888  .8'`888. .8'    888      888  888ooo88P'   88888[       888oooo8     888         888      888    `888  .8'`888. .8'
#     `888.8'  `888.8'     888      888  888`88b.     888`88b.     888    "     888         888      888     `888.8'  `888.8'
#      `888'    `888'      `88b    d88'  888  `88b.   888  `88b.   888          888       o `88b    d88'      `888'    `888'
#       `8'      `8'        `Y8bood8P'  o888o  o888o o888o  o888o o888o        o888ooooood8  `Y8bood8P'        `8'      `8'

gcloud iam service-accounts create "${SERVICE_NAME}-workflow" \
    --description="cloud workflow for ${SERVICE_NAME}" \
    --display-name="${SERVICE_NAME}-workflow"

# set permissions for service account

gcloud projects add-iam-policy-binding "${SERVICE_NAME}-workflow" \
    --member="serviceAccount:${SERVICE_NAME}-workflow@$PROJECT_ID.iam.gserviceaccount.com" \
    --role="roles/run.invoker"

gcloud workflows deploy "${SERVICE_NAME}-workflow" \
      --location="${REGION}" \
      --source=workflow.yaml \
      --service-account"=${SERVICE_NAME}-workflow@${PROJECT_ID}.iam.gserviceaccount.com"


#  .oooooo..o   .oooooo.   ooooo   ooooo oooooooooooo oooooooooo.   ooooo     ooo ooooo        oooooooooooo ooooooooo.
# d8P'    `Y8  d8P'  `Y8b  `888'   `888' `888'     `8 `888'   `Y8b  `888'     `8' `888'        `888'     `8 `888   `Y88.
# Y88bo.      888           888     888   888          888      888  888       8   888          888          888   .d88'
#  `"Y8888o.  888           888ooooo888   888oooo8     888      888  888       8   888          888oooo8     888ooo88P'
#      `"Y88b 888           888     888   888    "     888      888  888       8   888          888    "     888`88b.
# oo     .d8P `88b    ooo   888     888   888       o  888     d88'  `88.    .8'   888       o  888       o  888  `88b.
# 8""88888P'   `Y8bood8P'  o888o   o888o o888ooooood8 o888bood8P'      `YbodP'    o888ooooood8 o888ooooood8 o888o  o888o

gcloud iam service-accounts create "${SERVICE_NAME}-scheduler" \
    --description="cloud scheduler for ${SERVICE_NAME}" \
    --display-name="${SERVICE_NAME}-scheduler"

gcloud projects add-iam-policy-binding "${PROJECT_ID}" \
      --member "serviceAccount:${SERVICE_NAME}-scheduler@${PROJECT_ID}.iam.gserviceaccount.com" \
      --role "roles/workflow.invoker"

gcloud projects add-iam-policy-binding "${PROJECT_ID}" \
      --member "serviceAccount:${SERVICE_NAME}-scheduler@${PROJECT_ID}.iam.gserviceaccount.com" \
      --role "roles/logging.logWriter"


# There is no cloud scheduler in eu-north1
# At 06:00 on every 7th day-of-month.
# https://crontab.guru/#0_6_*/6_*_*
gcloud scheduler jobs create http ${SERVICE_NAME}-scheduler \
   --location europe-west3 \
   --schedule="0 6 */6 * *" \
   --uri="https://workflowexecutions.googleapis.com/v1/projects/${PROJECT}/locations/${REGION_WIDE}/workflows/${SERVICE_NAME}-scheduler/executions" \
   --http-method POST \
   --oauth-service-account-email "${SERVICE_NAME}-scheduler@${PROJECT_ID}.iam.gserviceaccount.com"


## TIPS
#
# List the service account roles
#gcloud projects get-iam-policy $PROJECT  \
#--flatten="bindings[].members" \
#--format='table(bindings.role)' \
#--filter="bindings.members:$ROLE"
