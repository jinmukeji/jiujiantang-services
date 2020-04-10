FROM 949191617935.dkr.ecr.cn-north-1.amazonaws.com.cn/jm-app/jm-app-base:2.0.1
LABEL   com.jinmuhealth.service.name="svc-sms-gw" \
        com.jinmuhealth.service.type="srv"
ADD ./svc-sms-gw_linux_amd64 /svc-sms-gw_linux_amd64
ENTRYPOINT [ "/svc-sms-gw_linux_amd64"]
