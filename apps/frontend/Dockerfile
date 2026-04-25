FROM node:24-alpine AS builder

ARG PROJECT=shell
WORKDIR /app

COPY package.json package-lock.json ./
RUN npm ci

COPY . .
RUN npx nx build ${PROJECT} --configuration=production

# --- Shell Runtime ---
FROM nginx:alpine AS mfe-shell-runtime

ARG MF_COMMON_URL=https://arda.io.vn/common/mf-manifest.json
ARG GIT_SHA=unknown

COPY --from=builder /app/dist/apps/shell/browser/ /usr/share/nginx/html/

RUN sed -i "s|http://localhost:3001/mf-manifest.json|${MF_COMMON_URL}|g" \
    /usr/share/nginx/html/module-federation.manifest.json

COPY apps/shell/nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]

# --- Common Runtime ---
FROM nginx:alpine AS mfe-common-runtime

COPY --from=builder /app/dist/apps/common/browser/ /usr/share/nginx/html/
COPY apps/common/nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
