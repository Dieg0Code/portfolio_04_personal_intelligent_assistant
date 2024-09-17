# Serverless RAG Diary

![infra](infra.png)

## Configuration

### Environment Variables

Need to set the following environment variables:

```env
SUPABASE_URL=yoursupabaseurl
SUPABASE_KEY=yoursupabasekey
OPENAI_API_KEY=youropenaikey
```

*In my case, I'm using Github secrets to store these values.*

Also need to create an s3 bucket and a DynamoDB table for terraform state.

S3 bucket:
```bash
aws s3api create-bucket --bucket terraform-state-rag-diary --region sa-east-1 --create-bucket-configuration LocationConstraint=sa-east-1
```

Enable versioning for the bucket (optional):
```bash
aws s3api put-bucket-versioning --bucket terraform-state-rag-diary --versioning-configuration Status=Enabled
```

DynamoDB table:
```bash
aws dynamodb create-table \
    --table-name terraform_locks_diary \
    --attribute-definitions AttributeName=LockID,AttributeType=S \
    --key-schema AttributeName=LockID,KeyType=HASH \
    --billing-mode PAY_PER_REQUEST \
    --region sa-east-1
```