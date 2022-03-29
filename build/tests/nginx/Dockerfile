FROM nginx:1.21

ENV NGINX_TEST_DIR ./build/tests/nginx

COPY ${NGINX_TEST_DIR}/nginx.conf /etc/nginx/conf.d/default.conf
COPY ${NGINX_TEST_DIR}/img/* /usr/share/nginx/html/
