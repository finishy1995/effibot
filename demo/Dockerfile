FROM node:alpine as builder

WORKDIR /code
ADD demo/package.json /code
#RUN npm config set registry http://registry.npm.taobao.org/
RUN npm install
ADD demo /code
RUN npm run build

FROM nginx:latest

COPY --from=builder /code/dist /usr/share/nginx/html
COPY demo/nginx/default.conf /etc/nginx/conf.d/default.conf

EXPOSE 4000
