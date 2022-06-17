#!/usr/bin/env bash
cat <<EOF > /workspace/dist/env.js
window.env = {
    clientId: '${CLIENT_ID}',
    authUri: '${AUTH_URI}',
    tokenUri: '${TOKEN_URI}',
    callbackUrl: '${CALLBACK_URL}',
    oauthResourceServer: '${OAUTH_RESOURCE_SERVER}',
};
EOF
nginx -p /workspace -c /workspace/nginx.conf