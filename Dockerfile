FROM alpine:3.17
EXPOSE 8080
COPY receptionist /receptionist

ENV TLSCERTFILE ""
ENV TLSKEYFILE ""

CMD ["/receptionist"]