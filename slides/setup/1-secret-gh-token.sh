
echo -n "my super secret github token" | gcloud secrets create "${SERVICE_NAME}-token" \
    --replication-policy="automatic" \  # automatic or user-managed
    --data-file=- 

# Do not save your command history

# doc: https://cloud.google.com/secret-manager/docs/create-secret#secretmanager-quickstart-gcloud
