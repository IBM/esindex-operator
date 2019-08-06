
[![Build Status](https://travis-ci.com/IBM/esindex-operator.svg?branch=master)](https://travis-ci.com/IBM/esindex-operator)

# IBM Cloud Operator for ElasticSearch Indices
The IBM Cloud Operator for Elastic Search Indices, as part of IBM Cloud operators, provides a Kubernetes CRD-Based API to manage the lifecycle of Elastic Search indices. It allows to provision elasticsearch indices from your Kubernetes cluster, using the ESIndex CRD. 

The Elastic Search access credentials can be specified in requests via reference to a Binding, Secret, or ConfigMap resource. The Binding resource is managed by IBM Cloud Binding Operator in conjuction with IBM Cloud Service Operator. Details can be found at https://github.com/IBM/cloud-operators. 

## Supported Features

* **Creation and Deletion** - Creates, deletes and monitors indices on Elastic Search service.

* **Credentials by Reference** - Elasticsearch access credentials can be provided using IBMCloud Binding, Secret, or ConfigMap.

## Install

To install the operator, run the following script:

```
curl -sL https://raw.githubusercontent.com/IBM/esindex-operator/master/hack/install-operator.sh | bash 
```
This will install the latest version of the operator. It will run in `ibmcloud-operators` namespace. To see its status, run this command:
```
kubectl get pod -n ibmcloud-operators
```

## Uninstall

```
curl -sL https://raw.githubusercontent.com/IBM/esindex-operator/master/hack/uninstall-operator.sh | bash 
```

## Use 

Sample yaml files are provided under [config/samples](https://github.com/IBM/esindex-operator/tree/master/config/samples). To create an index, you must already have an elasticsearch service instance and a source of elasticsearch access credential. The operator supports three options to specify the source of elasticsearch access credentials:

 - Binding.ibmcloud.ibm.com
 - Secret
 - ConfigMap
 
You may choose to use anyone of these options.  `esindex.yaml`, `esindex_secret.yaml` and `esindex_configmap.yaml` contain examples for each of them, respectively. The following commands assume the use of Binding.ibmcloud.ibm.com as the source.

1. Create an elasticserch service instance on IBM Cloud:

```
kubectl apply -f config/samples/elasticsearch.yaml
```

2. Create a binding instance:

```
kubectl apply -f config/samples/elasticsearch_binding.yaml
```

3. Create an index on the elasticsearch:

```
kubectl apply -f config/samples/esindex.yaml
```

## Troubleshooting

To find the current git revision for the operator, type:

```
kubectl exec -n ibmcloud-operators $(kubectl get pod -l " app=esindex-operator" -n ibmcloud-operators -o jsonpath='{.items[0].metadata.name}') -- cat git-rev
```

## Learn more about how to contribute

- [contributions](./CONTRIBUTING.md)
