FROM registry.access.redhat.com/ubi8/ubi-minimal:8.2-301.1592810506


ENV OPERATOR=/usr/local/bin/ibm-licensing-hub-operator \
    USER_UID=1001 \
    USER_NAME=ibm-licensing-hub-operator

# install operator binary
COPY build/_output/bin/ibm-licensing-hub-operator ${OPERATOR}

COPY build/bin /usr/local/bin
RUN  /usr/local/bin/user_setup

ENTRYPOINT ["/usr/local/bin/entrypoint"]

USER ${USER_UID}
