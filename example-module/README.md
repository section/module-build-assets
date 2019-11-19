# Example Module


## Directory Structure :

```
├── Dockerfile
├── .section-proxy.yaml
└── src
    ├── module source files
```

- **Dockerfile** : Used by the script *create-module-assets.sh* to build the module image. The image built using this file is run at Section's edge locations.
- **section-proxy.yaml** : The YAML file defining the POD used to run the module. This is further integrated with the Section base template before being published to Section's edge locations. This file is a YAML-formatted file with a custom schema for describing how a module will operate on the Section platform.. Further explanation of the metadata can is provided here : https://github.com/section-io/module-build-assets/blob/master/example-module/section-proxy.md

##### The resulting module built should strictly follow the proxy contract as defined by Section : https://www.section.io/docs/platform/reference/proxy-contract/
