#!/usr/bin/env bash

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
