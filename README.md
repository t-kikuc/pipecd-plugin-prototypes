# pipecd-plugin-prototypes
[DO NOT USE IN PROD] Prototypes of PipeCD plugin

## Temporary Usage

### Preparation

1. Run pipedv0
2. Register any one app with pipedv0 on UI
   1. Because pipedv1 cannot register an app yet...
   2. _Note: Do not specify plugin config in app.pipecd.yaml yet._
3. Create config of pipedv1

Example configuration (v0 and v1 compatible):
```yaml
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  # v0 config for registering an app
  platformProviders:
    - name: terraform-dev
      type: TERRAFORM # Any type will be OK because this will be ignored in v1.

  # v1 config for plugins
  plugins:
    - name: ECSPRESSO
      port: 7003 # Any unused port
      url: file:///<HOME>/.piped/plugins/ecspresso # Replace <HOME> to your home dir
      deployTargets: # This is requried if this plugin requires `config`.
        - name: dt1
          config:
            version: 2.4.5
```

To use remote binary, you can use `url: <URL>`.
I want to release prototype plugin binaries in this repo soon...

### Run pipedv1 with plugins

_Note: Set `PipelineStage.Visible: true` for stages to show in the v0-UI._

1. Build plugins
```sh
make build/plugin
```

2. Run pipedv1
3. Override the app config

Example of new app config:
```yaml
apiVersion: pipecd.dev/v1beta1
kind: TerraformApp # for registering an app by pipedv0. This will be ignored in v1.
spec:
  name: v1-try-ecspresso
  input: # Depends on plugin
    config: ecspresso.yml
  pipeline:
    stages:
      - name: ECSPRESSO_DEPLOY # Define your stages
```

4. Push the app config  -> A new Deployment will be triggered!!

