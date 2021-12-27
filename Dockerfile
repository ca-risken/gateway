FROM golang:1.16 as builder
WORKDIR /go/src/github.com/ca-risken/gateway/
ADD *.go go.* ./
ENV CGO_ENABLED=0 GOPRIVATE="github.com/ca-risken/*"
ARG GITHUB_USER
ARG GITHUB_TOKEN
RUN echo "machine github.com" > ~/.netrc
RUN echo "login $GITHUB_USER" >> ~/.netrc
RUN echo "password $GITHUB_TOKEN" >> ~/.netrc
RUN go build && cp gateway /go/bin/

FROM public.ecr.aws/risken/base/risken-base:v0.0.1
RUN mkdir -p /usr/local/gateway/doc
COPY --from=builder /go/bin/gateway /usr/local/gateway/bin/
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
WORKDIR /usr/local/gateway
ENTRYPOINT ["/usr/local/bin/env-injector"]
CMD ["bin/gateway"]
