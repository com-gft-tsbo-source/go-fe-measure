FROM golang:alpine AS buildimage

ARG COMPONENT
ARG MODULE
ARG NAME
ARG VERSION
ARG RELEASE_COMMIT

RUN apk add --update --no-cache gcc libc-dev git

COPY go/common/ /go/src/common/
COPY go/go.mod /go/src/
COPY "${COMPONENT}/${MODULE}/${NAME}/go/" "/go/src/${NAME}/"

WORKDIR /go/src

RUN mkdir -p ../bin
RUN go get ./...
RUN go build -ldflags "-linkmode external -extldflags -static -X com.gft.tsbo-training.src.go/common/ms-framework/microservice._build_stamp=${VERSION} -X com.gft.tsbo-training.src.go/common/ms-framework/microservice._build_commit=${RELEASE_COMMIT}" -a -o "../bin/${NAME}" "${NAME}/cmd/main.go"

# ----------------------------------------------------------------------------
# ----------------------------------------------------------------------------
# ----------------------------------------------------------------------------

FROM alpine

RUN apk add --update --no-cache curl bash

ARG DOMAIN
ARG CUSTOMER
ARG PROJECT
ARG COMPONENT
ARG MODULE
ARG NAME
ARG VERSION
ARG RELEASE_TYPE
ARG RELEASE_COMMIT
ARG RELEASE_STAMP

COPY --from=buildimage "/go/bin/${NAME}" "/${NAME}"
COPY "${COMPONENT}/${MODULE}/${NAME}/docker/entrypoint.sh" /
COPY "${COMPONENT}/${MODULE}/${NAME}/docker/index.html" /workdir/
COPY "${COMPONENT}/${MODULE}/${NAME}/docker/static/" /workdir/static/
RUN chmod -v u+x "/${NAME}" /entrypoint.sh
RUN sed -i "s/<NAME>/${NAME}/g" /entrypoint.sh
WORKDIR /workdir
ENTRYPOINT [ "/entrypoint.sh" ]

ENV APP_DOMAIN="${DOMAIN}" \
    APP_CUSTOMER="${CUSTOMER}" \
    APP_PROJECT="${PROJECT}" \
    APP_COMPONENT="${COMPONENT}" \
    APP_MODULE="${MODULE}" \
    APP_NAME="${NAME}" \
    APP_VERSION="${VERSION}" \
    APP_RELEASE_TYPE="${RELEASE_TYPE}" \
    APP_RELEASE_COMMIT="${RELEASE_COMMIT}" \
    APP_RELEASE_STAMP="${RELEASE_STAMP}"

LABEL app.kubernetes.io/part-of=${PROJECT} \
      app.kubernetes.io/component=${COMPONENT}.${MODULE} \
      app.kubernetes.io/name=${NAME} \
      app.kubernetes.io/version=${VERSION} \
      project=${PROJECT} \
      ${PROJECT}/domain=${DOMAIN} \
      ${PROJECT}/customer=${CUSTOMER} \
      ${PROJECT}/project=${PROJECT} \
      ${PROJECT}/component=${COMPONENT} \
      ${PROJECT}/module=${MODULE} \
      ${PROJECT}/name=${NAME} \
      ${PROJECT}/version=${VERSION} \
      ${PROJECT}/release/type=${RELEASE_TYPE} \
      ${PROJECT}/release/commit=${RELEASE_COMMIT} \
      ${PROJECT}/release/stamp=${RELEASE_STAMP}

EXPOSE 8080
