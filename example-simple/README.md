# Example Module
Contains an example module that can be deployed on the Section Edge Compute Platform. For more complete instructions, there is a tutorial posted here: https://www.section.io/blog/tutorial-deploy-containerized-workload-to-the-edge/


## Directory Structure :

```
create-module-assets.sh
example-simple
├── proxy
    └── nginx.conf
├── proxy_config
├── .section-proxy.yml
├── Dockerfile
├── prepare.sh
└── validate.sh

```

- **create-module-assets.sh** : A simple bash script that takes a single optional argument. It will call either example-full and example-simple (for this tutorial, we’ll be using example-simple). The script will package your workload source code and then generate the configuration files required to deploy onto the platform. Those configuration files are injected into the pod running the git-daemon. We’ll see how that works in just a bit.
- **example-simple/Dockerfile** :This is a standard Dockerfile. You can use it to configure the container running your workload.
- **example-simple/.section-proxy.yaml** : This file contains the configuration information for your workload. The parts you need to pay particular attention too are the image, names, and container definitions. When Section operates your module, we will tune and recommend additional settings. (Further explanation of the metadata can is provided here : https://github.com/section-io/module-build-assets/blob/master/example-module/section-proxy.md)
- **example-simple/prepare.sh** and **example-simple/validate.sh** : These are scripts that are automatically called before the script is deployed during a scaling event (prepare.sh) and before it’s available to serve content (validate.sh). These scripts should contain whatever is required for you to say your workload is prepped and ready to run. If you have no logic to run at these lifecycle stages, they must still be present and return an exit_val of 0. Any other exit_val indicates an invalid condition and will stop the deployment of the module.
- **example-simple/proxy_config/** : A folder that contains (by convention) any files that you’d like to see deployed for your workload. These files are unique to an environment. So if you have an API key or other configuration that may be different for a module configuration, you can deploy it into the module via this folder.
-  **example-simple/proxy/nginx.conf** : Things in this folder are available on all modules. nginx.conf is an example of one of those files required when setting up nginx. This file is common to all deployments of your module.

##### The resulting module built should strictly follow the proxy contract as defined by Section : https://www.section.io/docs/platform/reference/proxy-contract/
