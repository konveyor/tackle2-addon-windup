FROM registry.access.redhat.com/ubi8/go-toolset:1.16.12 as builder
ENV GOPATH=$APP_ROOT
COPY --chown=1001:0 . .
RUN make cmd

FROM registry.access.redhat.com/ubi8/ubi-minimal
USER root
RUN microdnf -y update && microdnf -y clean all
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
java-11-openjdk-headless \
openssh-clients \
unzip \
wget \
git \
subversion \
maven \
&& microdnf -y clean all
ARG WINDUP=https://repo1.maven.org/maven2/org/jboss/windup/tackle-cli/6.1.2.Final/tackle-cli-6.1.2.Final-offline.zip
RUN wget -qO /opt/windup.zip $WINDUP \
 && unzip /opt/windup.zip -d /opt \
 && rm /opt/windup.zip \
 && ln -s /opt/tackle-cli-*/bin/windup-cli /opt/windup

ENV HOME=/working \
    JAVA_HOME="/usr/lib/jvm/jre-11" \
    JAVA_VENDOR="openjdk" \
    JAVA_VERSION="11"
WORKDIR /working
COPY --from=builder /opt/app-root/src/bin/addon /usr/local/bin/addon
ENTRYPOINT ["/usr/local/bin/addon"]
