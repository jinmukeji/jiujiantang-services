FROM 949191617935.dkr.ecr.cn-north-1.amazonaws.com.cn/jm-app/jm-app-base:2.0.1
LABEL   com.jinmuhealth.service.name="svc-subscription" \
        com.jinmuhealth.service.type="srv"
ADD ./svc-subscription_linux_amd64 /svc-subscription_linux_amd64
ENTRYPOINT [ "/svc-subscription_linux_amd64"]

