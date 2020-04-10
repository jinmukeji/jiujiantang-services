FROM 949191617935.dkr.ecr.cn-north-1.amazonaws.com.cn/jm-app/jm-app-base:2.0.1
LABEL   com.jinmuhealth.service.name="svc-sem-gw" \
        com.jinmuhealth.service.type="srv"
ADD ./svc-sem-gw_linux_amd64 /svc-sem-gw_linux_amd64
ENTRYPOINT [ "/svc-sem-gw_linux_amd64"]
