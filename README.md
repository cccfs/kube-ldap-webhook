# Kubernetes LDAP Webhook Authentication

# Configuration

## kubernetes
Configure your kubernetes apiserver to use the kube-ldap-webhook webhook for authentication using the following configuration file.

* kube-apiserver.yaml
```yaml
- --enable-bootstrap-token-auth=true
- --authentication-token-webhook-config-file=/etc/kubernetes/webhook-token-auth-config.yaml
```

* webhook-token-auth-config.yaml
```yaml
# clusters refers to the remote service.
clusters:
- name: webhook-token-auth-cluster
  cluster:
    server: https://kube-ldap-webhook.example/token
    insecure-skip-tls-verify: False

# users refers to the API server's webhook configuration.
users:
- name: webhook-token-auth-user

# kubeconfig files require a context. Provide one for the API server.
current-context: webhook-token-auth
contexts:
- context:
    cluster: webhook-token-auth-cluster
    user: webhook-token-auth-user
  name: webhook-token-auth
```

## kubectl
```
TOKEN=$(curl https://kube-ldap-webhook.example.com/auth -u zhangsan)

kubectl config set-cluster <cluster name> --server=https://<apiserver url>:6443

kubectl config set-context <cluster name> --cluster=<cluster name> --user=<your user name>

kubectl config set-credentials <your user name> --token="$TOKEN"

kubectl config use-context test

//Change $HOME/.kube/config Add certificate-authority-data
```

#### Example
```yaml
apiVersion: v1
clusters:
- cluster:
    server: https://<apiserver url>:6443
    certificate-authority-data: "LS0tLS1CRUdJTiBDRVJUSnVaWFJsJBS3A5CkdnSUhvaUVFN1Vrdk1kS0tLQo="
  name: <cluster name>
contexts:
- context:
    cluster: <cluster name>
    user: <your user name>
  name: <cluster name>
current-context: <cluster name>
kind: Config
preferences: {}
users:
- name: <your user name>
  user:
    token: eyJhbGkxNDg1L

```

## client-go credential plugin