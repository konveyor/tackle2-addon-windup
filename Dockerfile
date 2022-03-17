FROM registry.access.redhat.com/ubi8/go-toolset:1.16.7 as builder
ENV GOPATH=$APP_ROOT
COPY --chown=1001:0 . .
RUN make cmd

FROM registry.access.redhat.com/ubi8/ubi-minimal
USER root
RUN echo -e "[centos8]" \
 "\nname = centos8" \
 "\nbaseurl = http://mirror.centos.org/centos/8-stream/AppStream/x86_64/os/" \
 "\nenabled = 1" \
 "\ngpgcheck = 0" > /etc/yum.repos.d/centos.repo
RUN echo -e "[WandiscoSVN]" \
 "\nname=Wandisco SVN Repo" \
 "\nbaseurl=http://opensource.wandisco.com/centos/6/svn-1.9/RPMS/$basearch/" \
 "\nenabled=1" \
 "\ngpgcheck=0" > /etc/yum.repos.d/wandisco.repo
RUN microdnf -y install \
  java-1.8.0-openjdk-headless \
  unzip \
  wget \
  git \
  subversion \
  maven \
 && microdnf -y clean all
ARG MTA=https://repo1.maven.org/maven2/org/jboss/windup/mta-cli/5.2.1.Final/mta-cli-5.2.1.Final-offline.zip
RUN wget -qO /opt/windup.zip $MTA \
 && unzip /opt/windup.zip -d /opt \
 && rm /opt/windup.zip \
 && ln -s /opt/mta-cli-5.2.1.Final/bin/mta-cli /opt/windup
ENV HOME=/working \
    JAVA_HOME="/usr/lib/jvm/jre" \
    JAVA_VENDOR="openjdk" \
    JAVA_VERSION="1.8.0"
WORKDIR /working
COPY --from=builder /opt/app-root/src/bin/addon /usr/local/bin/addon
ENTRYPOINT ["/usr/local/bin/addon"]
