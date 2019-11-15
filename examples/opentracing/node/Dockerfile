FROM node

WORKDIR /usr/src/node

COPY server.js package.json ./
RUN npm install
EXPOSE 8080

ENTRYPOINT ["npm", "start"]