# Module Source Files

## Directory Structure : 

```
├── environment_variables.lua.gotemplate
├── module_schema.json
├── handler
│   └── app
│       ├── Example GO Handler
├── lib
│   └── section
│       ├── Example LUA modules.

├── proxy
│   ├── content.lua
│   ├── nginx.conf
│   └── pubkey.gpg
└── tests
    ├── invalid-module.json
    ├── resolv.conf
    └── valid-module.json
```

- **environment_variables.lua.gotemplate**  : GO template file used for generating the module configuration. This will be generated using the custom configuration provided by the user in their Section git directory. Example of a configuration : https://www.section.io/docs/modules/kraken/reference/kraken-advanced-config/

- **module_schema.json** : JSON schema of the module configuration used for testing by the GO Handler.
