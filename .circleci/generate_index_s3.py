import boto3
import re
import os
import tempfile

S3_BUCKET_NAME = 'datadog-reliability-env'
client = boto3.client('s3', aws_access_key_id=os.getenv('AWS_ACCESS_KEY_ID'),aws_secret_access_key=os.getenv('AWS_SECRET_ACCESS_KEY'))
transfer = boto3.s3.transfer.S3Transfer(client)

# write information used by the build
with tempfile.NamedTemporaryFile(mode='w') as fp:
  for line in [os.getenv('CIRCLE_BRANCH'), os.getenv('CIRCLE_SHA1'), name, os.getenv('CIRCLE_USERNAME')]:
    fp.write(f'{line}\n')
  fp.seek(0)
  transfer.upload_file(fp.name, S3_BUCKET_NAME, 'go/index.txt')
