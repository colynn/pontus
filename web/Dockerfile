FROM 192.168.99.108:5000/library/node:12.2.0-alpine as builder
LABEL MAINTAINER='colynn.liu@gmail.com'

WORKDIR /root/app
COPY ./package.json package-lock.json ./
ENV NPM_CONFIG_LOGLEVEL warn
# ENV NPM_CONFIG_REGISTRY https://registry.npm.taobao.org

# RUN npm i --no-package-lock -d

# config registry to taobao, and node-sass to taobao(otherwise it will get from github)
RUN npm config set registry http://192.168.104.144:4873 && npm config set sass_binary_site https://npm.taobao.org/mirrors/node-sass/ 

# -d show info log
RUN npm i -d --registry=http://192.168.104.144:4873 --from-lock-file
# RUN npm i -d

COPY . .
RUN npm run build:prod

FROM 192.168.99.108:5000/library/nginx:1.16.0-alpine
LABEL MAINTAINER='colynn.liu@gmail.com'

ADD ./deploy/nginx/default.conf /etc/nginx/conf.d/default.conf
COPY --from=builder /root/app/dist /usr/share/nginx/html

# docs 
COPY docs/_book /usr/share/nginx/html/docs