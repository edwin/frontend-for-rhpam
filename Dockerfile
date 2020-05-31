FROM registry.access.redhat.com/ubi8/go-toolset
ENV GOPATH=$APP_ROOT
ENV GOBIN=$APP_ROOT/bin
ADD . $GOPATH/src/frontendForPam/
RUN go install frontendForPam

FROM registry.access.redhat.com/ubi8/ubi-minimal
COPY --from=0 /opt/app-root/bin/frontendForPam /usr/local/bin/frontendForPam
COPY --from=0 /opt/app-root/src/frontendForPam/static /usr/local/bin/static
WORKDIR /usr/local/bin
CMD frontendForPam
EXPOSE 8080