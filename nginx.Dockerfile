FROM nginx:alpine

COPY ./deployment/nginx/nginx.conf /etc/nginx/nginx.conf

RUN rm -rf /usr/share/nginx/html/*
COPY frontend/ /usr/share/nginx/html/

EXPOSE 8081
CMD ["nginx", "-g", "daemon off;"]