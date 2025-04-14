FROM python:3.13.3-alpine3.21

ARG APP_ROOT=/app/
ARG APP_NAME=go-job-master
ENV RUN=${APP_ROOT}${APP_NAME}
USER root

WORKDIR ${APP_ROOT}

COPY ./${APP_NAME}  ${APP_ROOT}

CMD ["/bin/sh", "-c", "$RUN"]

