package resources

import (
	"context"

	apierrs "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type patchCondition[T client.Object] func(ctx context.Context, proposed, current T) bool

type Request[T client.Object] struct {
	shouldCreate bool
	c            client.Client
}

func NewRequest[T client.Object](c client.Client) *Request[T] {
	return &Request[T]{c: c}
}

// CreateIfNotFound creates the given object if it doesn't already exist
func CreateIfNotFound[T client.Object](ctx context.Context, c client.Client, obj T) error {
	return NewRequest[T](c).
		CreateIfNotFound().
		Execute(ctx, obj)
}

func (r *Request[T]) Execute(
	ctx context.Context,
	proposed T,
) error {
	// Get the current status of the object
	current := proposed.DeepCopyObject().(T)
	err := r.c.Get(ctx, types.NamespacedName{Namespace: proposed.GetNamespace(), Name: proposed.GetName()}, current)
	switch {
	case apierrs.IsNotFound(err) && r.shouldCreate:
		return r.c.Create(ctx, proposed)
	case err != nil:
		return err
	}

	return nil
}

func (r *Request[T]) CreateIfNotFound() *Request[T] {
	r.shouldCreate = true
	return r
}
