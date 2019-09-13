FROM scratch
ADD ipboss /
COPY cfg.cfg /cfg.cfg
LABEL IPboss 1.0
ENTRYPOINT ["/ipboss"]
