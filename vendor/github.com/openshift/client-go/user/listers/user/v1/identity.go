// Code generated by lister-gen. DO NOT EDIT.

package v1

import (
	v1 "github.com/openshift/api/user/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// IdentityLister helps list Identities.
// All objects returned here must be treated as read-only.
type IdentityLister interface {
	// List lists all Identities in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1.Identity, err error)
	// Get retrieves the Identity from the index for a given name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1.Identity, error)
	IdentityListerExpansion
}

// identityLister implements the IdentityLister interface.
type identityLister struct {
	indexer cache.Indexer
}

// NewIdentityLister returns a new IdentityLister.
func NewIdentityLister(indexer cache.Indexer) IdentityLister {
	return &identityLister{indexer: indexer}
}

// List lists all Identities in the indexer.
func (s *identityLister) List(selector labels.Selector) (ret []*v1.Identity, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.Identity))
	})
	return ret, err
}

// Get retrieves the Identity from the index for a given name.
func (s *identityLister) Get(name string) (*v1.Identity, error) {
	obj, exists, err := s.indexer.GetByKey(name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1.Resource("identity"), name)
	}
	return obj.(*v1.Identity), nil
}
