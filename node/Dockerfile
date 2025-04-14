FROM python:3.13.3-alpine3.21

ARG PIP_NAME=pip.txt
ARG APP_ROOT=/app/
ARG APP_NAME=go-job-node
ENV RUN=${APP_ROOT}${APP_NAME}
USER root

WORKDIR ${APP_ROOT}

COPY ./${APP_NAME}  ${APP_ROOT}
COPY ./${PIP_NAME} ${APP_ROOT}

RUN mkdir ~/.pip \
    && mkdir -p ${APP_ROOT}/config \
    && printf '[global]\nindex-url=https://mirrors.aliyun.com/pypi/simple/' > ~/.pip/pip.conf \
    && pip install -no-cache-dir -r ${APP_ROOT}${PIP_NAME}

CMD ["/bin/sh", "-c", "$RUN"]


