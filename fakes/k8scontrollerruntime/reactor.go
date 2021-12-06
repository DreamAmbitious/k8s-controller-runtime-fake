
package k8scontrolruntime

import (
	"context"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
	"strings"
)

// FakeReactingCtrlRuntimeClient to fake the response and
// runtime objects
type FakeReactingCtrlRuntimeClient struct {
	client.Client
	FakeRunTimeMock map[string]FakeRunTimeMock
}

// FakeRunTimeMock replacement for reactors
// TestNameObject follows the method_Kind convention
type FakeRunTimeMock struct {
	TestNameObject rune
	MockFunc       func() (client.Object, error)
}

// Create reconciler's create
func (p *FakeReactingCtrlRuntimeClient) Create(
	ctx context.Context,
	obj client.Object,
	opts ...client.CreateOption) error {
	gvk, err := apiutil.GVKForObject(obj, p.Scheme())
	typeMeta := gvk.GroupKind().Kind
	for k, v := range p.FakeRunTimeMock {
		result := strings.Split(k, "_")
		if result[0] == "create" {
			if typeMeta == result[1] {
				obj, err = v.MockFunc()
				if err != nil {
					return err
				}
			}
		}
	}
	return p.Client.Create(ctx, obj, opts...)
}

// Update reconciler's update
func (p *FakeReactingCtrlRuntimeClient) Update(
	ctx context.Context,
	obj client.Object,
	opts ...client.UpdateOption) error {
	typeMeta := obj.GetObjectKind().GroupVersionKind()

	for k, v := range p.FakeRunTimeMock {
		result := strings.Split(k, "_")
		if result[0] == "update" {
			if typeMeta.Kind == result[1] {
				obj, err := v.MockFunc()
				if err != nil {
					return err
				}
				return p.Client.Update(ctx, obj, opts...)
			}
		}
	}
	return p.Client.Update(ctx, obj, opts...)
}

// Delete reconciler's delete
func (p *FakeReactingCtrlRuntimeClient) Delete(
	ctx context.Context,
	obj client.Object, opts ...client.DeleteOption) error {
	typeMeta := obj.GetObjectKind().GroupVersionKind()
	var err error
	for k, v := range p.FakeRunTimeMock {
		result := strings.Split(k, "_")
		if result[0] == "delete" {
			if typeMeta.Kind == result[1] {
				obj, err = v.MockFunc()
				if err != nil {
					return err
				}
			}
		}
	}
	return p.Client.Delete(ctx, obj, opts...)
}

// Get reconciler's Get
func (p *FakeReactingCtrlRuntimeClient) Get(
	ctx context.Context,
	key client.ObjectKey,
	obj client.Object) error {
	gvk, _ := apiutil.GVKForObject(obj, p.Scheme())
	typeMeta := gvk.GroupKind().Kind
	for k, v := range p.FakeRunTimeMock {
		result := strings.Split(k, "_")
		if result[0] == "get" {
			if typeMeta == result[1] {
				obj, err := v.MockFunc()
				if err != nil {
					return err
				}
				if obj != nil {
					return p.Client.Get(ctx, key, obj)
				}
			}
		}
	}
	return p.Client.Get(ctx, key, obj)
}