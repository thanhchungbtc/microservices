FROM node:alpine

ENV CI=true
WORKDIR /app
COPY . .
RUN npm install
CMD ["npm", "run", "serve"]
