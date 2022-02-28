FROM scratch
ADD mlog_linux_amd64 mlog_linux_amd64
EXPOSE 8003
ENTRYPOINT ["./mlog_linux_amd64"]