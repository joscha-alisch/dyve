sed -i '' "s/\$BACKEND/$DYVE_BACKEND/" /etc/nginx/conf.d/default.conf
nginx -g "daemon off;"