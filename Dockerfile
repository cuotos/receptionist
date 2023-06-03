FROM alpine:3.18
EXPOSE 8080
COPY receptionist /receptionist

ENV TLSCERTFILE ""
ENV TLSKEYFILE ""

CMD ["/receptionist"]