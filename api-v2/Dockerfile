FROM 949191617935.dkr.ecr.cn-north-1.amazonaws.com.cn/jm-app/jm-app-base:2.0.1
LABEL   com.jinmuhealth.service.name="web-rest-api-v2" \
        com.jinmuhealth.service.type="api"
ADD ./web-rest-api-v2_linux_amd64 /web-rest-api-v2_linux_amd64
ENTRYPOINT [ "/web-rest-api-v2_linux_amd64"]
