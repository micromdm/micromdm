# depsim container
FROM ubuntu:14.04
ADD depsim /usr/bin
RUN mkdir /etc/depsim

EXPOSE 8080
VOLUME ["/etc/depsim"]
ENTRYPOINT ["depsim"]
CMD ["start"]
