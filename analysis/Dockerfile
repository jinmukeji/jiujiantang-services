FROM 949191617935.dkr.ecr.cn-north-1.amazonaws.com.cn/jm-app/jm-app-base:2.0.1
LABEL   com.jinmuhealth.service.name="svc-analysis" \
        com.jinmuhealth.service.type="srv"
ADD ./svc-analysis_linux_amd64 /svc-analysis_linux_amd64
ADD ./ae_data_v2/lookups /lookups
ADD ./ae_data_v2/lua_src /lua_src
ADD ./ae_data_v2/question /question
ADD ./ae_data_v2/presets.yaml /presets.yaml
ENTRYPOINT [ "/svc-analysis_linux_amd64", "--x_lua_src_path=/lua_src","--x_templates_dir=/lookups","--x_question_dir=/question","--x_presets_file_path=/presets.yaml"]
