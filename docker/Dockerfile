FROM opencurvedocker/curve-manager:base-debian11

RUN mkdir -p /curve-manager/conf /curve-manager/logs /curve-manager/db

# prepare frontend
COPY website /curve-manager/website
RUN cd /curve-manager/website && npm i && npm run build-only

# prepare backend
COPY pigeon /curve-manager/pigeon
COPY nginx/nginx.conf /etc/nginx/nginx.conf
COPY pigeon.yaml /curve-manager/conf/pigeon.yaml

# run script
COPY entrypoint.sh /entrypoint.sh
CMD ["/entrypoint.sh"]
