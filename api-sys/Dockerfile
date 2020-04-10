FROM 949191617935.dkr.ecr.cn-north-1.amazonaws.com.cn/jm-app/jm-app-base:2.0.1
LABEL   com.jinmuhealth.service.name="web-rest-sys" \
        com.jinmuhealth.service.type="api"
ADD ./web-rest-sys_linux_amd64 /web-rest-sys_linux_amd64
ADD ./sys_data/config.yml /config.yml
ENTRYPOINT [ "/web-rest-sys_linux_amd64","--x_config_file=/config.yml"]
