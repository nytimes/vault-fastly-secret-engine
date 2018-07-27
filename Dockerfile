FROM vault:0.10.1

RUN mkdir /tmp/vault-plugins 
COPY vault-fastly-secret-engine /tmp/vault-plugins
RUN echo 'plugin_directory = "/tmp/vault-plugins"' >> /vault/config/plugin.hcl
