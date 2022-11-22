FROM node:14.20.0-alpine as builder

WORKDIR /root/app
COPY ./package*.json ./
ENV NPM_CONFIG_LOGLEVEL warn

# RUN npm i --no-package-lock -d
# config registry to taobao, and node-sass to taobao(otherwise it will get from github)
RUN npm install -g cnpm  --registry-https://registry.npm.taobao.org
# RUN npm config set sass_binary_site https://npm.taobao.org/mirrors/node-sass/ 
RUN cnpm install 

COPY . .
RUN npm run build:prod

FROM nginx:1.21.6-alpine

LABEL MAINTAINER='colynn.liu@gmail.com'

ADD ./deploy/nginx/default.conf /etc/nginx/conf.d/default.conf
COPY --from=builder /root/app/dist /usr/share/nginx/html