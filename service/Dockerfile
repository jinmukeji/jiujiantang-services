FROM 949191617935.dkr.ecr.cn-north-1.amazonaws.com.cn/jm-app/jm-app-base:2.0.1
LABEL   com.jinmuhealth.service.name="svc-biz-core" \
        com.jinmuhealth.service.type="srv"
ADD ./svc-biz-core_linux_amd64 /svc-biz-core_linux_amd64 
ADD ./ae_data /ae_data
ENTRYPOINT [ "/svc-biz-core_linux_amd64","--x_ae_config_dir=/ae_data"]
