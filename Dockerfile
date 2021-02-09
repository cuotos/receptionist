FROM alpine
EXPOSE 8080
COPY receptionist /receptionist
CMD ["/receptionist"]