FROM nginx:alpine

# Install gettext for environment variable substitution
RUN apk add --no-cache gettext

# Create the templates directory and add the configuration template using a heredoc
RUN mkdir -p /etc/nginx/templates && \
    cat <<'EOF' > /etc/nginx/templates/default.conf.template
server {
    listen 80;
    server_name ${DOMAIN};

    location / {
        proxy_pass http://${HOST};
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
}
EOF

# Set default values for environment variables
ENV DOMAIN="example.larek.tech" \
    HOST="127.0.0.1:8080"

# When container starts, envsubst will replace variables in the template
# and nginx will start in foreground
CMD ["sh", "-c", "envsubst '$DOMAIN $HOST' < /etc/nginx/templates/default.conf.template > /etc/nginx/conf.d/default.conf && nginx -g 'daemon off;'"]