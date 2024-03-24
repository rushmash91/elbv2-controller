// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

// Code generated by ack-generate. DO NOT EDIT.

package rule

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	ackv1alpha1 "github.com/aws-controllers-k8s/runtime/apis/core/v1alpha1"
	ackerr "github.com/aws-controllers-k8s/runtime/pkg/errors"
	acktypes "github.com/aws-controllers-k8s/runtime/pkg/types"

	svcapitypes "github.com/aws-controllers-k8s/elbv2-controller/apis/v1alpha1"
)

// ClearResolvedReferences removes any reference values that were made
// concrete in the spec. It returns a copy of the input AWSResource which
// contains the original *Ref values, but none of their respective concrete
// values.
func (rm *resourceManager) ClearResolvedReferences(res acktypes.AWSResource) acktypes.AWSResource {
	ko := rm.concreteResource(res).ko.DeepCopy()

	for f0idx, f0iter := range ko.Spec.Actions {
		if f0iter.TargetGroupRef != nil {
			ko.Spec.Actions[f0idx].TargetGroupARN = nil
		}
	}

	if ko.Spec.ListenerRef != nil {
		ko.Spec.ListenerARN = nil
	}

	return &resource{ko}
}

// ResolveReferences finds if there are any Reference field(s) present
// inside AWSResource passed in the parameter and attempts to resolve those
// reference field(s) into their respective target field(s). It returns a
// copy of the input AWSResource with resolved reference(s), a boolean which
// is set to true if the resource contains any references (regardless of if
// they are resolved successfully) and an error if the passed AWSResource's
// reference field(s) could not be resolved.
func (rm *resourceManager) ResolveReferences(
	ctx context.Context,
	apiReader client.Reader,
	res acktypes.AWSResource,
) (acktypes.AWSResource, bool, error) {
	namespace := res.MetaObject().GetNamespace()
	ko := rm.concreteResource(res).ko

	resourceHasReferences := false
	err := validateReferenceFields(ko)
	if fieldHasReferences, err := rm.resolveReferenceForActions_TargetGroupARN(ctx, apiReader, namespace, ko); err != nil {
		return &resource{ko}, (resourceHasReferences || fieldHasReferences), err
	} else {
		resourceHasReferences = resourceHasReferences || fieldHasReferences
	}

	if fieldHasReferences, err := rm.resolveReferenceForListenerARN(ctx, apiReader, namespace, ko); err != nil {
		return &resource{ko}, (resourceHasReferences || fieldHasReferences), err
	} else {
		resourceHasReferences = resourceHasReferences || fieldHasReferences
	}

	return &resource{ko}, resourceHasReferences, err
}

// validateReferenceFields validates the reference field and corresponding
// identifier field.
func validateReferenceFields(ko *svcapitypes.Rule) error {

	for _, f0iter := range ko.Spec.Actions {
		if f0iter.TargetGroupRef != nil && f0iter.TargetGroupARN != nil {
			return ackerr.ResourceReferenceAndIDNotSupportedFor("Actions.TargetGroupARN", "Actions.TargetGroupRef")
		}
	}

	if ko.Spec.ListenerRef != nil && ko.Spec.ListenerARN != nil {
		return ackerr.ResourceReferenceAndIDNotSupportedFor("ListenerARN", "ListenerRef")
	}
	if ko.Spec.ListenerRef == nil && ko.Spec.ListenerARN == nil {
		return ackerr.ResourceReferenceOrIDRequiredFor("ListenerARN", "ListenerRef")
	}
	return nil
}

// resolveReferenceForActions_TargetGroupARN reads the resource referenced
// from Actions.TargetGroupRef field and sets the Actions.TargetGroupARN
// from referenced resource. Returns a boolean indicating whether a reference
// contains references, or an error
func (rm *resourceManager) resolveReferenceForActions_TargetGroupARN(
	ctx context.Context,
	apiReader client.Reader,
	namespace string,
	ko *svcapitypes.Rule,
) (hasReferences bool, err error) {
	for f0idx, f0iter := range ko.Spec.Actions {
		if f0iter.TargetGroupRef != nil && f0iter.TargetGroupRef.From != nil {
			hasReferences = true
			arr := f0iter.TargetGroupRef.From
			if arr.Name == nil || *arr.Name == "" {
				return hasReferences, fmt.Errorf("provided resource reference is nil or empty: Actions.TargetGroupRef")
			}
			obj := &svcapitypes.TargetGroup{}
			if err := getReferencedResourceState_TargetGroup(ctx, apiReader, obj, *arr.Name, namespace); err != nil {
				return hasReferences, err
			}
			ko.Spec.Actions[f0idx].TargetGroupARN = (*string)(obj.Status.ACKResourceMetadata.ARN)
		}
	}

	return hasReferences, nil
}

// getReferencedResourceState_TargetGroup looks up whether a referenced resource
// exists and is in a ACK.ResourceSynced=True state. If the referenced resource does exist and is
// in a Synced state, returns nil, otherwise returns `ackerr.ResourceReferenceTerminalFor` or
// `ResourceReferenceNotSyncedFor` depending on if the resource is in a Terminal state.
func getReferencedResourceState_TargetGroup(
	ctx context.Context,
	apiReader client.Reader,
	obj *svcapitypes.TargetGroup,
	name string, // the Kubernetes name of the referenced resource
	namespace string, // the Kubernetes namespace of the referenced resource
) error {
	namespacedName := types.NamespacedName{
		Namespace: namespace,
		Name:      name,
	}
	err := apiReader.Get(ctx, namespacedName, obj)
	if err != nil {
		return err
	}
	var refResourceSynced, refResourceTerminal bool
	for _, cond := range obj.Status.Conditions {
		if cond.Type == ackv1alpha1.ConditionTypeResourceSynced &&
			cond.Status == corev1.ConditionTrue {
			refResourceSynced = true
		}
		if cond.Type == ackv1alpha1.ConditionTypeTerminal &&
			cond.Status == corev1.ConditionTrue {
			return ackerr.ResourceReferenceTerminalFor(
				"TargetGroup",
				namespace, name)
		}
	}
	if refResourceTerminal {
		return ackerr.ResourceReferenceTerminalFor(
			"TargetGroup",
			namespace, name)
	}
	if !refResourceSynced {
		return ackerr.ResourceReferenceNotSyncedFor(
			"TargetGroup",
			namespace, name)
	}
	if obj.Status.ACKResourceMetadata == nil || obj.Status.ACKResourceMetadata.ARN == nil {
		return ackerr.ResourceReferenceMissingTargetFieldFor(
			"TargetGroup",
			namespace, name,
			"Status.ACKResourceMetadata.ARN")
	}
	return nil
}

// resolveReferenceForListenerARN reads the resource referenced
// from ListenerRef field and sets the ListenerARN
// from referenced resource. Returns a boolean indicating whether a reference
// contains references, or an error
func (rm *resourceManager) resolveReferenceForListenerARN(
	ctx context.Context,
	apiReader client.Reader,
	namespace string,
	ko *svcapitypes.Rule,
) (hasReferences bool, err error) {
	if ko.Spec.ListenerRef != nil && ko.Spec.ListenerRef.From != nil {
		hasReferences = true
		arr := ko.Spec.ListenerRef.From
		if arr.Name == nil || *arr.Name == "" {
			return hasReferences, fmt.Errorf("provided resource reference is nil or empty: ListenerRef")
		}
		obj := &svcapitypes.Listener{}
		if err := getReferencedResourceState_Listener(ctx, apiReader, obj, *arr.Name, namespace); err != nil {
			return hasReferences, err
		}
		ko.Spec.ListenerARN = (*string)(obj.Status.ACKResourceMetadata.ARN)
	}

	return hasReferences, nil
}

// getReferencedResourceState_Listener looks up whether a referenced resource
// exists and is in a ACK.ResourceSynced=True state. If the referenced resource does exist and is
// in a Synced state, returns nil, otherwise returns `ackerr.ResourceReferenceTerminalFor` or
// `ResourceReferenceNotSyncedFor` depending on if the resource is in a Terminal state.
func getReferencedResourceState_Listener(
	ctx context.Context,
	apiReader client.Reader,
	obj *svcapitypes.Listener,
	name string, // the Kubernetes name of the referenced resource
	namespace string, // the Kubernetes namespace of the referenced resource
) error {
	namespacedName := types.NamespacedName{
		Namespace: namespace,
		Name:      name,
	}
	err := apiReader.Get(ctx, namespacedName, obj)
	if err != nil {
		return err
	}
	var refResourceSynced, refResourceTerminal bool
	for _, cond := range obj.Status.Conditions {
		if cond.Type == ackv1alpha1.ConditionTypeResourceSynced &&
			cond.Status == corev1.ConditionTrue {
			refResourceSynced = true
		}
		if cond.Type == ackv1alpha1.ConditionTypeTerminal &&
			cond.Status == corev1.ConditionTrue {
			return ackerr.ResourceReferenceTerminalFor(
				"Listener",
				namespace, name)
		}
	}
	if refResourceTerminal {
		return ackerr.ResourceReferenceTerminalFor(
			"Listener",
			namespace, name)
	}
	if !refResourceSynced {
		return ackerr.ResourceReferenceNotSyncedFor(
			"Listener",
			namespace, name)
	}
	if obj.Status.ACKResourceMetadata == nil || obj.Status.ACKResourceMetadata.ARN == nil {
		return ackerr.ResourceReferenceMissingTargetFieldFor(
			"Listener",
			namespace, name,
			"Status.ACKResourceMetadata.ARN")
	}
	return nil
}
