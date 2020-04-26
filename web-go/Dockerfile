FROM 949191617935.dkr.ecr.cn-north-1.amazonaws.com.cn/jm-app/jm-app-base:2.0.1
LABEL   com.jinmuhealth.service.name="web-go" \
        com.jinmuhealth.service.type="web"
ADD ./web-go_linux_amd64 /web-go_linux_amd64
ADD ./web-go/resource.yml ./resource.yml
ENTRYPOINT [ "/web-go_linux_amd64", "--server_address=:9100","--x_config_file=/resource.yml"]
