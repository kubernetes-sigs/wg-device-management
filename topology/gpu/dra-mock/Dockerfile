# Stage 1: Build mocknvml and fake-imex
FROM golang:1.26 as mocknvml-builder

RUN apt-get update && apt-get install -y git build-essential

RUN git clone https://github.com/NVIDIA/k8s-test-infra /k8s-test-infra
WORKDIR /k8s-test-infra

RUN make -C pkg/gpu/mocknvml

RUN mkdir -p /build/bin
RUN go build -o /build/bin/nvidia-imex ./cmd/fake-imex/daemon
RUN go build -o /build/bin/nvidia-imex-ctl ./cmd/fake-imex/ctl

# Stage 2: Build dra-driver-nvidia-gpu
FROM golang:1.26 as driver-builder

RUN apt-get update && apt-get install -y git build-essential make

# Clone the DRA driver repo
RUN git clone https://github.com/kubernetes-sigs/dra-driver-nvidia-gpu /dra-driver-nvidia-gpu
WORKDIR /dra-driver-nvidia-gpu
COPY driver-patches/cdi.go cmd/gpu-kubelet-plugin/cdi.go
COPY driver-patches/cdi-computedomain.go cmd/compute-domain-kubelet-plugin/cdi.go

# Build the driver commands.
RUN mkdir -p /build/bin && \
    (make PREFIX=/build/bin cmds || go build -o /build/bin/ ./cmd/...)

# Stage 3: Final image
FROM ubuntu:22.04

RUN apt-get update && apt-get install -y ca-certificates nvidia-utils-570 && rm -rf /var/lib/apt/lists/*
# Copy the mock library from Stage 1
COPY --from=mocknvml-builder /k8s-test-infra/pkg/gpu/mocknvml/libnvidia-ml.so /usr/lib/x86_64-linux-gnu/libnvidia-ml-mock.so
RUN ln -sf /usr/lib/x86_64-linux-gnu/libnvidia-ml-mock.so /usr/lib/x86_64-linux-gnu/libnvidia-ml.so
RUN ln -sf /usr/lib/x86_64-linux-gnu/libnvidia-ml-mock.so /usr/lib/x86_64-linux-gnu/libnvidia-ml.so.1
RUN cp /usr/lib/x86_64-linux-gnu/libnvidia-ml-mock.so /usr/lib/x86_64-linux-gnu/libnvidia-ml.so.570.158.01

# Find installed libcuda.so.570.x and copy to requested version
RUN LIBCUDA_FILE=$(find /usr/lib/x86_64-linux-gnu/ -name "libcuda.so.570.*" | head -n 1) && \
    if [ -n "$LIBCUDA_FILE" ]; then \
        cp $LIBCUDA_FILE /usr/lib/x86_64-linux-gnu/libcuda.so.570.158.01; \
    else \
        echo "Warning: libcuda.so.570.* not found"; \
    fi

# Copy all binaries from Stage 2 to /usr/bin/
COPY --from=driver-builder /build/bin/* /usr/bin/

# Copy fake imex binaries from Stage 1
COPY --from=mocknvml-builder /build/bin/nvidia-imex /usr/bin/
COPY --from=mocknvml-builder /build/bin/nvidia-imex-ctl /usr/bin/

# Copy the prestart script
COPY --from=driver-builder /dra-driver-nvidia-gpu/hack/kubelet-plugin-prestart.sh /usr/bin/
RUN chmod +x /usr/bin/kubelet-plugin-prestart.sh
RUN sed -i '2i exit 0' /usr/bin/kubelet-plugin-prestart.sh

# Create dummy nvidia-cdi-hook
RUN echo '#!/bin/sh' > /usr/bin/nvidia-cdi-hook && \
    echo 'exit 0' >> /usr/bin/nvidia-cdi-hook && \
    chmod +x /usr/bin/nvidia-cdi-hook

# Copy other scripts and templates
COPY --from=driver-builder /dra-driver-nvidia-gpu/scripts/bind_to_driver.sh /usr/bin/
COPY --from=driver-builder /dra-driver-nvidia-gpu/scripts/unbind_from_driver.sh /usr/bin/
COPY --from=driver-builder /dra-driver-nvidia-gpu/templates /templates

ENV LD_LIBRARY_PATH=/usr/local/lib:$LD_LIBRARY_PATH
ENTRYPOINT ["gpu-kubelet-plugin"]
