# vppsim container
FROM ubuntu:14.04
ADD vppsim /usr/bin
RUN mkdir /etc/vppsim

EXPOSE 8080
VOLUME ["/etc/vppsim"]
ENTRYPOINT ["vppsim"]
CMD ["start"]
