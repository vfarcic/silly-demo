FROM node:22-alpine

LABEL org.opencontainers.image.source=https://github.com/vfarcic/silly-demo
WORKDIR /app
COPY package.json .
RUN npm install
COPY public/ /app/public
COPY src/ /app/src
ARG VERSION
ENV VERSION=$VERSION

CMD ["npm", "start", "--host", "0.0.0.0"]
