FROM golang AS builder
WORKDIR /go/src/github.com/CyberAgent/mimosa-gateway/
ADD *.go go.* ./
ENV CGO_ENABLED=0 GO111MODULE=on GOPRIVATE="github.com/CyberAgent/*"
ARG GITHUB_USER
ARG GITHUB_TOKEN
RUN echo "machine github.com" > ~/.netrc
RUN echo "login $GITHUB_USER" >> ~/.netrc
RUN echo "password $GITHUB_TOKEN" >> ~/.netrc
RUN go install -v

FROM alpine
RUN apk add --no-cache ca-certificates tzdata \
  && mkdir -p /usr/local/mimosa-gateway/doc \
  && cd /tmp/ && wget -q https://github.com/gassara-kys/env-injector/releases/download/v0.0.6/env-injector_0.0.6_Linux_x86_64.tar.gz \
  && tar xvfz env-injector_0.0.6_Linux_x86_64.tar.gz && mv env-injector /usr/local/bin/ && rm -rf /tmp/*
COPY --from=builder /go/bin/mimosa-gateway /usr/local/mimosa-gateway/bin/
ENV PORT= \
  DEBUG= \
  USER_IDENTITY_HEADER=x-amzn-oidc-identity \
  OIDC_DATA_HEADER=x-amzn-oidc-data \
  IDP_PROVIDER_NAME= \
  FINDING_SVC_ADDR= \
  IAM_SVC_ADDR= \
  PROJECT_SVC_ADDR= \
  ALERT_SVC_ADDR= \
  REPORT_SVC_ADDR= \
  AWS_SVC_ADDR= \
  AWS_ACTIVITY_SVC_ADDR= \
  OSINT_SVC_ADDR= \
  DIAGNOSIS_SVC_ADDR= \
  CODE_SVC_ADDR= \
  GOOGLE_SVC_ADDR= \
  TZ=Asia/Tokyo
WORKDIR /usr/local/mimosa-gateway
ENTRYPOINT ["/usr/local/bin/env-injector"]
CMD ["bin/mimosa-gateway"]
