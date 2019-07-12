
[![Build Status](https://travis-ci.com/IBM/esindex-operator.svg?branch=master)](https://travis-ci.com/IBM/esindex-operator)

# IBM Cloud Operator for ElasticSearch Indices
The IBM Cloud Operator for Elastic Search Indices provides a Kubernetes CRD-Based API to manage the lifecycle of indices of the Elastic Search service (databases-for-elasticsearch) in IBM public cloud. This operator allows to provision elasticsearch indices from your Kubernetes cluster, using the ESIndex CRD.

## Supported Features

* **Creation and Deletion** - Creates or deletes indices on the subscribed elasticsearch service using the credentials managed by IBM Cloud Operators for Service and Binding.

* **Bind Only Mode** - Allows access to existing elasticsearch indices

## Prerequisites 

* Install IBM Cloud Operators for Services and Bindings: https://github.com/IBM/cloud-operators

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

Sample yaml files are provided under config/samples. To create an index, you must have an elasticsearch service instance and a binding instance exsit. Run the following commands to create them.

Create an elasticserch service instance:

```
kubectl apply -f config/samples/elasticsearch.yaml
```

Create a binding instance:

```
kubectl apply -f config/samples/elasticsearch_binding.yaml
```

Create an index on the elasticsearch 

```
kubectl apply -f config/samples/esindex.yaml
```

## Troubleshooting

To find the current git revision for the operator, type:

```
kubectl exec -n ibmcloud-operators $(kubectl get pod -l " app=ibmcloud-esindex-operator" -n ibmcloud-operators -o jsonpath='{.items[0].metadata.name}') -- cat git-rev
```

## Learn more about how to contribute

- [contributions](./CONTRIBUTING.md)
