FROM scratch
COPY wioctl /
ENTRYPOINT ["/wioctl"]

