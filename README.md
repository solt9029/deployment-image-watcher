# deployment-image-watcher

- deployment-image-watcher is a custom controller that notifies you of deployment's template container image update on your slack channel.
- this custom controller has been created with operator-sdk (v0.15.2).

## using deployment

```sh
vi secret.yaml # add your slack api token.
vi operator.yaml # add your target slack channel.
kubectl apply -f deploy/
```

## using operator-sdk (local)

```sh
operator-sdk run --local
```
