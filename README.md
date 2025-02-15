# must-gather-demo

## Table of Content



## Building the code

```bash
go build -o must-gather ./cmd/must-gather
```

## Running the code

- Running the code collection :
```bash
# ./must-gather --image=quay.io/openshift-release-dev/ocp-v4.0-art-dev@sha256:aaac3feab704eb100776366ccbed8eaf9c7c0b9dea0bce597495fce1225d592f --host-network=true
Successfully created pod "must-gather-rpxln" in namespace "default"
```
Since the namespace was not provided the `default` one its being used.

- On the cluster side:

One of my cluster node its in a `NotReady` state:
```bash
# oc get nodes 
NAME                               STATUS     ROLES                         AGE    VERSION
hub-ctlplane-0.5g-deployment.lab   NotReady   control-plane,master,worker   3d4h   v1.30.7
hub-ctlplane-1.5g-deployment.lab   Ready      control-plane,master,worker   3d4h   v1.30.7
hub-ctlplane-2.5g-deployment.lab   Ready      control-plane,master,worker   3d4h   v1.30.7
```
Despite this situation, the `must-gather` pod its not allocated in a random manner, but to the first `Ready` node:

```bash
# oc get pods -A -o wide -w | grep -i must
default must-gather-rpxln    0/2     Pending             0              0s      <none>         <none>                             <none>           <none>
default must-gather-rpxln    0/2     Pending             0              0s      <none>         hub-ctlplane-1.5g-deployment.lab   <none>           <none>
default must-gather-rpxln    0/2     ContainerCreating   0              0s      172.16.30.21   hub-ctlplane-1.5g-deployment.lab   <none>           <none>
default must-gather-rpxln    2/2     Running             0              2s      172.16.30.21   hub-ctlplane-1.5g-deployment.lab   <none>           <none>
```

