/*
 * Copyright 2019 IBM Corporation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package esindex

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"

	ibmcloudv1alpha1 "github.com/ibm/cloud-operators/pkg/apis/ibmcloud/v1alpha1"
	esindexv1alpha1 "github.com/ibm/esindex-operator/pkg/apis/ibmcloud/v1alpha1"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// RestResult is a struct for REST call result
type RestResult struct {
	StatusCode int
	Body       string
	ErrorType  string
}

// EsConnection is a struct for elastic search connection
type EsConnection struct {
	HTTPS EsHTTPS `json:"https"`
}

// EsHTTPS is a struct for elastic search https connection
type EsHTTPS struct {
	Composed []string `json:"composed"`
}

// ErrorTypeEsURINotFound - elastic search uri is not found in binding secret
const ErrorTypeEsURINotFound string = "EsUriNotFound"

// restCallFunc : common rest call fun
func restCallFunc(rsString string, postBody []byte, method string, header string, token string, expectReturn bool) (RestResult, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	restClient := http.Client{
		Timeout:   time.Second * 300,
		Transport: tr,
	}
	u, _ := url.ParseRequestURI(rsString)
	urlStr := u.String()
	var req *http.Request
	if postBody != nil {

		req, _ = http.NewRequest(method, urlStr, bytes.NewBuffer(postBody))
	} else {
		req, _ = http.NewRequest(method, urlStr, nil)
	}

	if token != "" {
		if header == "" {
			req.Header.Set("Authorization", token)
		} else {
			req.Header.Set(header, token)
		}
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := restClient.Do(req)
	if err != nil {
		return RestResult{}, err
	}
	defer res.Body.Close()

	if expectReturn {
		body, err := ioutil.ReadAll(res.Body)
		result := RestResult{StatusCode: res.StatusCode, Body: string(body[:])}
		return result, err
	}
	return RestResult{}, nil
}

// createIndex : create an index on elastic search
func (r *ReconcileEsIndex) createIndex(obj *esindexv1alpha1.EsIndex) (RestResult, error) {
	uri, err := r.getESUri(obj)
	if err != nil {
		return RestResult{ErrorType: ErrorTypeEsURINotFound}, err
	}
	uri0 := uri + "/" + obj.Spec.IndexName
	var indexObj IndexCreate
	indexObj.Settings.NumberOfShards = obj.Spec.NumberOfShards
	indexObj.Settings.NumberOfReplicas = obj.Spec.NumberOfReplicas
	putBody, _ := json.Marshal(&indexObj)

	if obj.Spec.BindOnly {
		// bind only, check if the index exists
		resp, err := restCallFunc(uri0, putBody, "GET", "", "", true)
		return resp, err
	}
	// create index on elastic search
	resp, err := restCallFunc(uri0, putBody, "PUT", "", "", true)
	return resp, err

}

// getIndex : get index from elastic search
func (r *ReconcileEsIndex) getIndex(obj *esindexv1alpha1.EsIndex) (RestResult, error) {
	uri, err := r.getESUri(obj)
	if err != nil {
		return RestResult{ErrorType: ErrorTypeEsURINotFound, Body: err.Error()}, err
	}
	uri0 := uri + "/" + obj.Spec.IndexName
	var body []byte
	resp, err := restCallFunc(uri0, body, "GET", "", "", true)
	return resp, err
}

// deleteIndex : delete an index on elastic search
func (r *ReconcileEsIndex) deleteIndex(obj *esindexv1alpha1.EsIndex) (RestResult, error) {

	if obj.Spec.BindOnly || obj.Status.State == ResourceStateFailed {
		//do nothing on remote
		logt.Info("bindOnly and do nothing for deletion", "indexName", obj.Spec.IndexName)
		return RestResult{StatusCode: 200}, nil
	}

	uri, err := r.getESUri(obj)
	if err != nil {
		return RestResult{ErrorType: ErrorTypeEsURINotFound, Body: err.Error()}, err
	}
	uri0 := uri + "/" + obj.Spec.IndexName
	var body []byte
	resp, err := restCallFunc(uri0, body, "DELETE", "", "", true)
	return resp, err
}

// getEsURIComposed gets ES URI from user specified source
func (r *ReconcileEsIndex) getEsURIComposed(namespace string, composed esindexv1alpha1.CredSource) (string, error) {
	if composed.ConfigMapKeyRef == nil { //not configmap
		if composed.SecretKeyRef == nil { //not secret
			return "", fmt.Errorf("err: neither secretKeyRef nor configMapKeyRef is specified. must specify one of them")
		}
		secret, err := r.getSecret(namespace, composed.SecretKeyRef.Name)
		if err != nil {
			return "", err
		}
		uri, err := r.getValueFromSecret(secret, composed.SecretKeyRef.Key)
		return uri, err
	}
	// get from configmap
	configmap, err := r.getConfigMap(namespace, composed.ConfigMapKeyRef.Name)
	if err != nil {
		return "", err
	}
	uri, err := r.getValueFromConfigMap(configmap, composed.ConfigMapKeyRef.Key)
	return uri, err
}

// getESUri : returns elastic search URI
func (r *ReconcileEsIndex) getESUri(obj *esindexv1alpha1.EsIndex) (string, error) {
	if obj.Spec.BindingFrom.Name == "" { //cretential is not from binding
		if reflect.DeepEqual(obj.Spec.EsURIComposed, esindexv1alpha1.CredSource{}) { //empty object
			return "", fmt.Errorf("err: neither bindingFrom nor esURIComposed is specified. must specify one of them")
		}
		uri, err := r.getEsURIComposed(obj.ObjectMeta.Namespace, obj.Spec.EsURIComposed)
		return uri, err
	}

	// cretential is specified in binding
	binding, err := r.getBinding(obj.ObjectMeta.Namespace, obj.Spec.BindingFrom.Name)
	if err != nil || len(binding.OwnerReferences) < 1 {
		logt.Info("failed to get binding", "bindingName", obj.Spec.BindingFrom.Name)
		return "", err
	}
	var secretname string
	if binding.Spec.SecretName == "" {
		secretname = obj.Spec.BindingFrom.Name
	} else {
		secretname = binding.Spec.SecretName
	}
	logt.Info("secretname from binding", "name", secretname)
	secrt, err := r.getSecret(obj.ObjectMeta.Namespace, secretname)
	if err != nil {
		logt.Info("getSecret error", "secretName", secretname)
		return "", err
	}
	// get uri from secret.data and decode
	datajson, _ := json.Marshal(secrt.Data)
	var mydat map[string]interface{}
	if err := json.Unmarshal(datajson, &mydat); err != nil {
		logt.Error(err, "json.Unmarshal of elastic search secret data failed", "secretName", obj.Spec.BindingFrom.Name)
		return "", err
	}
	if mydat["connection"] == nil {
		logt.Info("elastic search credentials not found in secret, nil connection", "secretName", obj.Spec.BindingFrom.Name)
		return "", fmt.Errorf("err: elastic search credentials not found in secret")
	}

	connection, err := base64.StdEncoding.DecodeString(mydat["connection"].(string))
	if err != nil {
		logt.Error(err, "base64 decode failed", "connectionBase64encoded", mydat["connection"].(string))
		return "", err
	}
	var conn EsConnection
	if err := json.Unmarshal(connection, &conn); err != nil {
		logt.Error(err, "json.Unmarshal of decoded connection failed")
		return "", err
	}
	if conn.HTTPS.Composed == nil || conn.HTTPS.Composed[0] == "" {
		return "", fmt.Errorf("err: elastic search composed uri not found in secret")
	}
	return conn.HTTPS.Composed[0], nil
}

// setCRDOwnerReference : set owner reference for index deletion upon deletion of it's owner
// currently owner reference info can be obtained from service's secret.
func (r *ReconcileEsIndex) setCRDOwnerReference(obj *esindexv1alpha1.EsIndex) error {
	binding, err := r.getBinding(obj.ObjectMeta.Namespace, obj.Spec.BindingFrom.Name)
	if err != nil || len(binding.OwnerReferences) < 1 {
		logt.Info("failed to get binding's OwnerReference", "bindingName", obj.Spec.BindingFrom.Name)
		return err
	}

	boolPtr := func(b bool) *bool { return &b }
	ownerReference := &metav1.OwnerReference{
		APIVersion: binding.OwnerReferences[0].APIVersion,
		Kind:       binding.OwnerReferences[0].Kind,
		Name:       binding.OwnerReferences[0].Name,
		UID:        binding.OwnerReferences[0].UID,
		Controller: boolPtr(true),
	}
	metaobj, _ := meta.Accessor(obj)
	existingRefs := metaobj.GetOwnerReferences()
	existingRefs = append(existingRefs, *ownerReference)
	metaobj.SetOwnerReferences(existingRefs)
	return nil
}

func (r *ReconcileEsIndex) getSecret(namespace string, secretname string) (*v1.Secret, error) {
	var secret v1.Secret
	name := strings.TrimSpace(secretname) //remove whitespaces
	if err := r.Client.Get(context.Background(), client.ObjectKey{Namespace: namespace, Name: name}, &secret); err != nil {
		return nil, err
	}
	return &secret, nil
}

func (r *ReconcileEsIndex) getConfigMap(namespace string, configmaptname string) (*v1.ConfigMap, error) {
	var conf v1.ConfigMap
	name := strings.TrimSpace(configmaptname) //remove whitespaces
	if err := r.Client.Get(context.Background(), client.ObjectKey{Namespace: namespace, Name: name}, &conf); err != nil {
		logt.Error(err, " getConfigMap failed")
		return nil, err
	}
	return &conf, nil
}

// getBinding: get Binding object
func (r *ReconcileEsIndex) getBinding(namespace string, bindingName string) (*ibmcloudv1alpha1.Binding, error) {
	binding := ibmcloudv1alpha1.Binding{}

	if err := r.Client.Get(context.Background(), client.ObjectKey{Namespace: namespace, Name: bindingName}, &binding); err != nil {
		logt.Error(err, "binding object not found", "bindingName", bindingName, "namespace", namespace)
		return nil, err
	}
	return &binding, nil
}

func (r *ReconcileEsIndex) getValueFromSecret(secret *v1.Secret, key string) (string, error) {

	datajson, _ := json.Marshal(secret.Data)
	var mydat map[string]interface{}
	if err := json.Unmarshal(datajson, &mydat); err != nil {
		logt.Error(err, "json.Unmarshal of elastic search secret data failed", "secretName", secret.ObjectMeta.Name)
		return "", err
	}
	if mydat[key] == nil {
		logt.Info("elastic search URI not found in secret", "secretName", secret.ObjectMeta.Name)
		return "", fmt.Errorf("err: elastic search credentials not found in secret %v", secret.ObjectMeta.Name)
	}

	uri, err := base64.StdEncoding.DecodeString(mydat[key].(string))
	if err != nil {
		logt.Error(err, "base64 decode failed", "connectionBase64encoded", mydat[key].(string))
		return "", err
	}
	return string(uri), nil
}

func (r *ReconcileEsIndex) getValueFromConfigMap(configmap *v1.ConfigMap, key string) (string, error) {

	datajson, _ := json.Marshal(configmap.Data)
	var mydat map[string]interface{}
	if err := json.Unmarshal(datajson, &mydat); err != nil {
		logt.Error(err, "json.Unmarshal of configmap data failed", "configmapName", configmap.ObjectMeta.Name)
		return "", err
	}
	if mydat[key] == nil {
		logt.Info("key not found in configmap", "configmapName", configmap.ObjectMeta.Name, "key", key)
		return "", fmt.Errorf("err: elastic search credentials not found in configmap %v", configmap.ObjectMeta.Name)
	}
	return mydat[key].(string), nil
}
