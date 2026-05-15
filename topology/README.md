# Multi-Node Topology Experimentation

This directory contains experiments and proof-of-concept DRA resource examples
and drivers for managing multi-node topology (as opposed to intra-node
topology).

## Relevant Background Links

### GitHub Issues & Pull Requests
* [Health-Aware Topology · Issue #6006 · kubernetes/enhancements](https://github.com/kubernetes/enhancements/issues/6006)
* [DRA: ResourceClaim Support for Workloads · Issue #5729 · kubernetes/enhancements](https://github.com/kubernetes/enhancements/issues/5729#issuecomment-4337498700)
* [KEP-5732: Topology-aware workload scheduling · Pull Request #5733 · kubernetes/enhancements](https://github.com/kubernetes/enhancements/pull/5733)
* [PodGroup / Workload API Integration with ComputeDomain · Issue #934 · kubernetes-sigs/nvidia-dra-driver-gpu](https://github.com/kubernetes-sigs/nvidia-dra-driver-gpu/issues/934)
* [dra-evolution: partitioning of devices · Issue #20 · kubernetes-sigs/wg-device-management](https://github.com/kubernetes-sigs/wg-device-management/issues/20#issuecomment-2168189769)

### GitHub Repositories
* [kubernetes-sigs/dra-driver-topology](https://github.com/kubernetes-sigs/dra-driver-topology)
* [kubernetes-sigs/dra-driver-google-tpu](https://github.com/kubernetes-sigs/dra-driver-google-tpu)
* [kubernetes-sigs/dra-driver-nvidia-gpu](https://github.com/kubernetes-sigs/dra-driver-nvidia-gpu)
* [thameem-abbas/dra-llm-d-usability](https://github.com/thameem-abbas/dra-llm-d-usability)

### Conference Presentations (PDFs)
* [More Nodes, More Problems: Solving Multi-Host GPU/TPU Scheduling With Dynamic Resource Allocation (KubeCon EU 2025)](https://static.sched.com/hosted_files/kccnceu2025/3d/KubeConEU25_%20More%20Nodes%2C%20More%20Problems.pdf?_gl=1*1qzyra4*_gcl_au*MTk5MDY2OTU3MC4xNzQ2NzM1MTc4*FPAU*MTk5MDY2OTU3MC4xNzQ2NzM1MTc4)
* [Gaining more control over node scheduling with the Topology/Block Plugin (SLUG 24)](https://slurm.schedmd.com/SLUG24/NVIDIA-Craig_Tierney.pdf)

### Google Docs
* [\[PUBLIC\] Extensible Allocation Algorithms in DRA](https://docs.google.com/document/u/0/d/1GKhH-dDlMziun5p-fdTLXpNj-K9UK1Hi2nuk0SxztVg/edit?resourcekey=0-ZztAemssXwSJyK5bsn0tgQ)
* [\[Public\] Modeling Topology and Multi-Node Logical Devices](https://docs.google.com/document/u/0/d/1Fg9ughIRMtt1HmDqiGWV-w9OKdrcKf_PsH4TjuP8Y40/edit)
* [\[PUBLIC\] TPU Slice Topology Allocation Model](https://docs.google.com/document/u/0/d/1Fxrb_Buv6prKze1DgcNAAyriBi8zioYm3DXvVTswZT4/edit)
* [\[PUBLIC\] Using DRA for the TAS Infrastructure Model](https://docs.google.com/document/u/0/d/1IitNIFQMqwdHUJUkwZEK3Amd2hHE0f4ZGzZpuikQR6A/edit)
* [device-based topology-aware scheduling](https://docs.google.com/document/d/11rC_qDtArIOx_ZQfM-G4H5qIYFm0rX_GUz1yPGPrv3k/edit?usp=sharing)
* [\[External\] Kubernetes support for GH200 / GB200](https://docs.google.com/document/d/1PrdDofsPFVJuZvcv-vtlI9n2eAh-YVf_fRQLIVmDwVY/edit?usp=sharing)
* [\[Public\] Multi-host DRA for TPUs](https://docs.google.com/document/d/10S2RefDffPLCCB0sssQnbLb_g_Q_vYmwJAsHPYCfq7U/edit?usp=sharing)
