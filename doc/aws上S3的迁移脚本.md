## s3上的迁移脚本

```sh
aws s3 mv s3://jm-measurement-data/dev/ s3://jm-measurement-data/dev/spec-v1/ --include "*.txt" --exclude "spec-v1"  --recursive
aws s3 mv s3://jm-measurement-data/testing/ s3://jm-measurement-data/testing/spec-v1/ --include "*.txt" --exclude "spec-v1"  --recursive
aws s3 mv s3://jm-measurement-data/staging/ s3://jm-measurement-data/staging/spec-v1/ --include "*.txt" --exclude "spec-v1"  --recursive
aws s3 mv s3://jm-measurement-data/sandbox/ s3://jm-measurement-data/sandbox/spec-v1/ --include "*.txt" --exclude "spec-v1"  --recursive
# 发布时执行
aws s3 mv s3://jm-measurement-data/production/ s3://jm-measurement-data/production/spec-v1/ --include "*.txt" --exclude "spec-v1"  --recursive
```

