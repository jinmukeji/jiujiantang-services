FROM 949191617935.dkr.ecr.cn-north-1.amazonaws.com.cn/jm-app/jm-app-base:2.0.1
LABEL   com.jinmuhealth.service.name="svc-device" \
        com.jinmuhealth.service.type="srv"
ADD ./svc-device_linux_amd64 /svc-device_linux_amd64
ENTRYPOINT [ "/svc-device_linux_amd64"]

