# Server configuration files

This directory contains the files used by the server running in the module.

## Directory Structure :

```
.
├── content.lua
├── nginx.conf
└── pubkey.gpg
```

- **content.lua**  : The *content_by_lua* file optionally used by nginx in this case to import Lua modules under *example-module/src/lib/section* and provide an interface to them.

- **nginx.conf** : The nginx configuration file used by the nginx server. This file helps in ensuring that the module honours the Section module contract : https://www.section.io/docs/platform/reference/proxy-contract/
- **pubkey.gpg** : The gpg key file used to install Openresty for the example module.
