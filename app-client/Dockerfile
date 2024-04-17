FROM node:21.7-alpine3.18
LABEL authors="jorka"

WORKDIR /usr/src/app-client

COPY ./package.json ./
COPY ./vite.config.js ./

RUN node -v

RUN apk add --update npm
RUN npm -v

RUN npm install
RUN npm fund

COPY ./ ./

EXPOSE 8000

CMD ["npm", "run", "dev"]