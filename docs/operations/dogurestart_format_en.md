# Dogurestart format

The Dogurestart-CR can be used to explicitly restart Dogus. The Dogu operator then scales
the corresponding Dogu deployment down and then up again.

All fields of a Dogurestart-CR are described below and illustrated with examples.

## Complete example

```yaml
apiVersion: k8s.cloudogu.com/v2
kind: DoguRestart
metadata:
  generateName: usermgt-restart-
spec:
 doguName: usermgt
```

Please note: `generateName` can be used to generate a unique name. However, this does not work
with `kubectl apply` but only `kubectl create`.

## doguName

* Required
* Data type: string
* Content: The `doguName` field specifies the name of the dogu to be restarted.
* Example: `“doguName”: “usermgt”`