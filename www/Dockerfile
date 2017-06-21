FROM scratch
MAINTAINER Pavel Snagovsky <pavel@snagovsky.com>
LABEL description="Castle test deployment frontend."

ADD www /

EXPOSE 80
# NOTE: App required argument is the backend end point, must be provided at
#       container runtime.
ENTRYPOINT ["/www"]
# Default argument
CMD ["--help"]