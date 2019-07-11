#!/bin/bash
#
# Copyright 2019 IBM Corp. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

# check prerequisites
IBMCLOUDOP=$(kubectl get crd | grep -c ^bindings.ibmcloud.ibm.com)
if [[ "$IBMCLOUDOP" -lt 1 ]]
then 
  echo "missing prerequisites: ibmcloud operators for service and binding."
  echo "please follow this link to install it: https://github.com/IBM/cloud-operators"
  exit
else echo "found the prerequisite bindings.ibmcloud.ibm.com"
fi

