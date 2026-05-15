# Copyright The Kubernetes Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

CLUSTER_NAME=${CLUSTER_NAME:-dra-alpha}
LOCATION=${LOCATION:-us-west4}
GPU=${GPU:-B200}
name=${1:-gb200-00}
zone=${2:-b}

AUTOREPAIR=""
if gcloud container clusters describe $CLUSTER_NAME --location $LOCATION | grep -q alpha ; then
	AUTOREPAIR=" --no-enable-autorepair --no-enable-autoupgrade "
fi

gcloud container node-pools create fake-${name} \
    --cluster=${CLUSTER_NAME} \
    --location=${LOCATION} \
    --node-locations=${LOCATION}-${zone} \
    --machine-type="e2-highcpu-4" \
    --num-nodes=18 \
    --node-labels=gke-no-default-nvidia-gpu-device-plugin=true,nvidia.com/gpu.present=true,cloud.google.com/gke-nvidia-gpu-dra-driver=true,example.com/nvldomain=${name},example.com/fake-gpu=${GPU} $AUTOREPAIR

# Create ConfigMap for mock NVML
CONFIG_MAP_NAME="mock-nvml-config-${name}"
MOCK_CONFIG_FILE="/tmp/mock-nvml-config-${name}.yaml"

cat <<'EOF' > $MOCK_CONFIG_FILE
# Mock NVML Configuration: B200
# Full configuration for nvidia-smi -x -q compatibility
# Use: export MOCK_NVML_CONFIG=/path/to/this/file.yaml

version: "1.0"

system:
  driver_version: "535.288.01"
  nvml_version: "12.535.288.01"
  cuda_version: "12.6"
  cuda_version_major: 12
  cuda_version_minor: 6

device_defaults:
  name: "NVIDIA B200"
  brand: "nvidia"
  serial: "1562849203750"
  board_part_number: "699-2G540-0200-000"
  vbios_version: "96.00.A0.00.01"
  architecture: "blackwell"
  compute_capability:
    major: 10
    minor: 0
  num_gpu_cores: 18432
  inforom:
    image_version: "B200.0200.00.01"
    oem_object: "2.1"
    ecc_object: "7.16"
    pwr_object: "1.0"
  memory:
    total_bytes: 206158430208
    reserved_bytes: 1073741824
    free_bytes: 205084688384
    used_bytes: 0
    memory_bus_width: 8192
  bar1_memory:
    total_bytes: 274877906944
    free_bytes: 274877906944
    used_bytes: 0
  pci:
    device_id: 0x234010DE
    subsystem_id: 0x181810DE
  pcie:
    max_link_gen: 6
    current_link_gen: 6
    max_link_width: 16
    current_link_width: 16
    replay_counter: 0
    tx_throughput_kbps: 0
    rx_throughput_kbps: 0
  power:
    management_supported: true
    management_mode: "enabled"
    default_limit_mw: 1000000
    enforced_limit_mw: 1000000
    min_limit_mw: 400000
    max_limit_mw: 1200000
    current_draw_mw: 130000
    power_state: "P0"
    total_energy_consumption_mj: 900000
  thermal:
    temperature_gpu_c: 35
    temperature_memory_c: 33
    shutdown_threshold_c: 95
    slowdown_threshold_c: 90
    max_operating_c: 85
    target_temperature_c: 85
  fan:
    count: 0
    speed_percent: "N/A"
    target_speed_percent: "N/A"
  clocks:
    graphics_current: 345
    graphics_max: 2100
    graphics_app: 2100
    graphics_app_default: 2100
    sm_current: 345
    sm_max: 2100
    memory_current: 2500
    memory_max: 2500
    memory_app: 2500
    memory_app_default: 2500
    video_current: 1200
    video_max: 2100
  clocks_throttle_reasons:
    gpu_idle: true
    applications_clocks_setting: false
    sw_power_cap: false
    hw_slowdown: false
    hw_thermal_slowdown: false
    hw_power_brake_slowdown: false
    sync_boost: false
    sw_thermal_slowdown: false
    display_clocks_setting: false
  supported_clocks:
    memory_clocks:
      - freq_mhz: 2500
        graphics_clocks: [345, 690, 1035, 1380, 1725, 1890, 2100]
  performance_state: "P0"
  utilization:
    gpu: 0
    memory: 0
    encoder: 0
    decoder: 0
    jpeg: 0
    ofa: 0
  encoder_stats:
    session_count: 0
    average_fps: 0
    average_latency_us: 0
  fbc_stats:
    session_count: 0
    average_fps: 0
    average_latency_us: 0
  ecc:
    mode_current: "enabled"
    mode_pending: "enabled"
    default_mode: "enabled"
    errors:
      volatile:
        single_bit:
          device_memory: 0
          l1_cache: 0
          l2_cache: 0
          register_file: 0
          texture_memory: 0
          total: 0
        double_bit:
          device_memory: 0
          l1_cache: 0
          l2_cache: 0
          register_file: 0
          texture_memory: 0
          total: 0
      aggregate:
        single_bit:
          device_memory: 0
          l1_cache: 0
          l2_cache: 0
          register_file: 0
          texture_memory: 0
          total: 0
        double_bit:
          device_memory: 0
          l1_cache: 0
          l2_cache: 0
          register_file: 0
          texture_memory: 0
          total: 0
  retired_pages:
    single_bit_retirement:
      count: 0
      addresses: []
    double_bit_retirement:
      count: 0
      addresses: []
    pending_blacklist: false
    pending_retirement: false
  remapped_rows:
    correctable: 0
    uncorrectable: 0
    pending: false
    failure_occurred: false
  display:
    mode: "disabled"
    active: "disabled"
  persistence_mode: "enabled"
  compute_mode: "default"
  mig:
    mode_current: "disabled"
    mode_pending: "disabled"
    max_gpu_instances: 7
  gpu_operation_mode:
    current: "all_on"
    pending: "all_on"
  driver_model:
    current: "N/A"
    pending: "N/A"
  accounting:
    mode: "disabled"
    buffer_size: 4000
  virtualization:
    mode: "none"
    host_vgpu_mode: "N/A"
  gsp_firmware:
    mode: "enabled"
    version: "535.288.01"
  features:
    transformer_engine: true
    fp4_support: true
    fp8_support: true
    confidential_compute: true
    decompression_engine: true
    fifth_gen_tensor_cores: true
  processes: []

devices:
  - index: 0
    uuid: "GPU-b2000000-0000-0000-0000-000000000000"
    serial: "1562849203700"
    pci:
      bus_id: "00000000:1A:00.0"
    minor_number: 0
  - index: 1
    uuid: "GPU-b2000000-0000-0000-0000-000000000001"
    serial: "1562849203701"
    pci:
      bus_id: "00000000:1B:00.0"
    minor_number: 1
  - index: 2
    uuid: "GPU-b2000000-0000-0000-0000-000000000002"
    serial: "1562849203702"
    pci:
      bus_id: "00000000:1C:00.0"
    minor_number: 2
  - index: 3
    uuid: "GPU-b2000000-0000-0000-0000-000000000003"
    serial: "1562849203703"
    pci:
      bus_id: "00000000:1D:00.0"
    minor_number: 3

nvlink:
  version: 5
  links_per_gpu: 18
  bandwidth_per_link_gbps: 100
  links:
    - link: 0
      state: "active"
      remote_device_type: "gpu"
      remote_pci_bus_id: "00000000:1B:00.0"
EOF

if [ "$GPU" = "B300" ]; then
    # Mock B300 by replacing the name in the config
    sed -i 's/NVIDIA B200/NVIDIA B300/g' $MOCK_CONFIG_FILE
fi

NS="dra-driver-nvidia-gpu-mock-${name}"

kubectl create ns $NS

# Create ConfigMap in the cluster
kubectl create configmap $CONFIG_MAP_NAME --from-file=config.yaml=$MOCK_CONFIG_FILE -n $NS --dry-run=client -o yaml | kubectl apply -f -

# Deploy the driver to the created node pool
SCRIPT_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &> /dev/null && pwd)

# Create a temporary values file for Helm
HELM_VALUES_FILE="/tmp/helm-values-${name}.yaml"
cat <<EOF > $HELM_VALUES_FILE
image:
  repository: us-central1-docker.pkg.dev/jbelamaric-dev/jbelamaric-dev/dra-driver-nvidia-gpu-mock
  tag: latest
  pullPolicy: Always

mockNvml:
  configMaps:
    - name: $CONFIG_MAP_NAME
      mountPath: /var/lib/nvml-mock/config.yaml
      subPath: config.yaml

kubeletPlugin:
  priorityClassName: ""
  nodeSelector:
    example.com/nvldomain: "${name}"
  containers:
    gpus:
      env:
        - name: MOCK_NVML_CONFIG
          value: /var/lib/nvml-mock/config.yaml
    computeDomains:
      env:
        - name: MOCK_NVML_CONFIG
          value: /var/lib/nvml-mock/config.yaml
EOF

helm upgrade --install dra-driver-nvidia-gpu-${name} $SCRIPT_DIR/dra-mock/dra-driver-nvidia-gpu \
    --namespace $NS \
    --set 'gpuResourcesEnabledOverride=true' \
    -f $HELM_VALUES_FILE

