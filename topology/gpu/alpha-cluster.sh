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

gcloud container clusters create ${CLUSTER_NAME} \
    --location=${LOCATION} \
    --cluster-version=1.36.0-gke.1379000 \
    --enable-kubernetes-alpha --no-enable-autoupgrade --no-enable-autorepair \
    --alpha-cluster-feature-gates=DRAResourcePoolStatus=true,GenericWorkload=true,WorkloadWithJob=true,DRAWorkloadResourceClaims=true,GangScheduling=true \
    --logging SYSTEM,WORKLOAD,API_SERVER,SCHEDULER,CONTROLLER_MANAGER
