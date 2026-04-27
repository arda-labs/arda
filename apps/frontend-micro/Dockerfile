FROM node:24-alpine AS builder

ARG PROJECT=shell
WORKDIR /app

COPY package.json package-lock.json ./
RUN npm ci

COPY . .
RUN npx ng build ${PROJECT}

# --- Runtime environments --- #

FROM nginx:alpine AS mfe-shell-runtime
COPY --from=builder /app/dist/shell/browser/ /usr/share/nginx/html/
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]

FROM nginx:alpine AS mfe-iam-runtime
COPY --from=builder /app/dist/iam/browser/ /usr/share/nginx/html/
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]

FROM nginx:alpine AS mfe-common-runtime
COPY --from=builder /app/dist/common/browser/ /usr/share/nginx/html/
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
