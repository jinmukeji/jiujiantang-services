FROM 949191617935.dkr.ecr.cn-north-1.amazonaws.com.cn/jm-app/jm-app-base:2.0.1
LABEL   com.jinmuhealth.service.name="web-rest-websocket" \
        com.jinmuhealth.service.type="web"
ADD ./web-rest-websocket_linux_amd64 /web-rest-websocket_linux_amd64
ENTRYPOINT [ "/web-rest-websocket_linux_amd64"]
