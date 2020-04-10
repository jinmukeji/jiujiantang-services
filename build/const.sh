#!/usr/bin/env bash

# AWS ECR 端点
AWS_ECR_REPO=949191617935.dkr.ecr.cn-north-1.amazonaws.com.cn

# Docker Image 构建使用的命名空间
DOCKER_IMAGE_NAMESPACE=jm-app

# Docker Compose 运行容器时使用的项目名称
DOCKER_COMPOSE_PROJECT_NAME="jm"

DEFAULT_IMAGE_FILTER="label=maintainer=tech@jinmuhealth.com"
