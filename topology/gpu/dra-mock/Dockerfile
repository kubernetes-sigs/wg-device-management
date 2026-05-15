# Stage 1: Build mocknvml shared library
FROM golang:1.26 as mocknvml-builder

RUN apt-get update && apt-get install -y git build-essential

# Clone the k8s-test-infra repo
RUN git clone https://github.com/NVIDIA/k8s-test-infra /k8s-test-infra
WORKDIR /k8s-test-infra

# Build the mocknvml shared library
RUN make -C pkg/gpu/mocknvml

# Stage 2: Build dra-driver-nvidia-gpu
FROM golang:1.26 as driver-builder

RUN apt-get update && apt-get install -y git build-essential make

# Clone the DRA driver repo
RUN git clone https://github.com/kubernetes-sigs/dra-driver-nvidia-gpu /dra-driver-nvidia-gpu
WORKDIR /dra-driver-nvidia-gpu
COPY driver-patches/cdi.go cmd/gpu-kubelet-plugin/cdi.go

# Build the driver commands. We try to use the Makefile target 'cmds'
# and fallback to a generic go build if needed.
RUN mkdir -p /build/bin && \
    (make PREFIX=/build/bin cmds || go build -o /build/bin/ ./cmd/...)

# Stage 3: Final image
FROM ubuntu:22.04

RUN apt-get update && apt-get install -y ca-certificates nvidia-utils-535 && rm -rf /var/lib/apt/lists/*
# Copy the mock library from Stage 1
COPY --from=mocknvml-builder /k8s-test-infra/pkg/gpu/mocknvml/libnvidia-ml.so /usr/lib/x86_64-linux-gnu/libnvidia-ml.so.mock
RUN ln -sf /usr/lib/x86_64-linux-gnu/libnvidia-ml.so.mock /usr/lib/x86_64-linux-gnu/libnvidia-ml.so
RUN ln -sf /usr/lib/x86_64-linux-gnu/libnvidia-ml.so.mock /usr/lib/x86_64-linux-gnu/libnvidia-ml.so.1

# Copy all binaries from Stage 2 to /usr/bin/
COPY --from=driver-builder /build/bin/* /usr/bin/
# Copy the prestart script
COPY --from=driver-builder /dra-driver-nvidia-gpu/hack/kubelet-plugin-prestart.sh /usr/bin/
RUN chmod +x /usr/bin/kubelet-plugin-prestart.sh
# Patch the script to exit early for mock environment
RUN sed -i '2i exit 0' /usr/bin/kubelet-plugin-prestart.sh

# Create dummy nvidia-cdi-hook
RUN echo '#!/bin/sh' > /usr/bin/nvidia-cdi-hook && \
    echo 'exit 0' >> /usr/bin/nvidia-cdi-hook && \
    chmod +x /usr/bin/nvidia-cdi-hook

# Copy other scripts and templates
COPY --from=driver-builder /dra-driver-nvidia-gpu/scripts/bind_to_driver.sh /usr/bin/
COPY --from=driver-builder /dra-driver-nvidia-gpu/scripts/unbind_from_driver.sh /usr/bin/
COPY --from=driver-builder /dra-driver-nvidia-gpu/templates /templates

# Set LD_LIBRARY_PATH to force using the mock library
ENV LD_LIBRARY_PATH=/usr/local/lib:$LD_LIBRARY_PATH

# Set a default entrypoint (the chart usually overrides this with the specific binary)
ENTRYPOINT ["gpu-kubelet-plugin"]
