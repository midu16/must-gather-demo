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

```bash
# oc get pods -A -o wide -w | grep -i must
default must-gather-rpxln    0/2     Pending             0              0s      <none>         <none>                             <none>           <none>
default must-gather-rpxln    0/2     Pending             0              0s      <none>         hub-ctlplane-1.5g-deployment.lab   <none>           <none>
default must-gather-rpxln    0/2     ContainerCreating   0              0s      172.16.30.21   hub-ctlplane-1.5g-deployment.lab   <none>           <none>
default must-gather-rpxln    2/2     Running             0              2s      172.16.30.21   hub-ctlplane-1.5g-deployment.lab   <none>           <none>
```

