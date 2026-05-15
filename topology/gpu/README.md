# GPU Multi-Node Logical Device Experiments

This directory contains experiments that use multi-node logical devices to represent NVL72 instances.

## Cluster Setup

To set up a test cluster, you need to create a cluster with Kubernetes alpha features enabled and specific feature gates. You can see an example of this in `alpha-cluster.sh`.

The required feature gates are:
- `DRAResourcePoolStatus=true`
- `GenericWorkload=true`
- `WorkloadWithJob=true`
- `DRAWorkloadResourceClaims=true`
- `GangScheduling=true`
- `DRADeviceTaintRules=true`

You must also enable the resource.k8s.io/v1beta2/devicetaintrules API.

### Deployment Steps

To set up the environment, follow these steps:

1.  **Build the mock NVIDIA DRA driver:**
    This driver uses `mocknvlm` to simulate GPUs on non-GPU nodes. Build it and push it to your container registry.
2.  **Deploy base resources:**
    Use the helm chart (found in `dra-mock/`) to deploy the base resources for the DRA driver.
3.  **Create node pools:**
    Use the `fake-nvl72.sh` script to create node pools that represent NVL72s. This script will also deploy individual instances of the mock driver.
4.  **Deploy `dra-driver-noop`:**
    Build and deploy the `dra-driver-noop` from [gke-labs/dra-drivers](https://github.com/gke-labs/dra-drivers), configuring it to register itself as `nvl.example.com`.

### Patched Scheduler

Note that Kubernetes versions `1.36.0` or `1.36.1` will need a patched `kube-scheduler` that includes the fix from [PR #139017](https://github.com/kubernetes/kubernetes/pull/139017).

Once that is built, it can be deployed using the `dra-scheduler.yaml` in this repo (located at `dra-mock/dra-scheduler.yaml`). The different job manifests in this directory are currently set up to use that scheduler.

## Running Experiments

### Whole-Node NVL72 Model

Currently, there is one experiment that models NVL72s as DRA devices with consumable capacity containing the number of GPUs. This is in the "wholenode" manifests.

This experiment is so named because it does not allow consumption of individual GPUs, only all GPUs in a node.

#### How to try it out:

1.  Apply the setup manifest:
    ```bash
    kubectl apply -f wholenode-setup.yaml
    ```
2.  Apply the individual jobs to see how PodGroup-based gang scheduling works with this, allocating the NVL resources and sharing them amongst the Pods:
    ```bash
    kubectl apply -f wholenode-job0.yaml
    # Apply other wholenode-job*.yaml files as needed
    ```
3.  **Test GPU Failure:**
    You can start a job, and then use a `DeviceTaintRule` to see how a GPU failure will result in an automatic reschedule of the Pod on a node in the same NVL72 (assuming one is available). To do this, modify and apply `taint.yaml`:
    ```bash
    kubectl apply -f taint.yaml
    ```
