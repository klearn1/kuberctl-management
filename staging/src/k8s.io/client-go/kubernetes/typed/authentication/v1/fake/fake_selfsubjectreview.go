/*
Copyright The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	v1 "k8s.io/api/authentication/v1"
	gentype "k8s.io/client-go/gentype"
	clientset "k8s.io/client-go/kubernetes/typed/authentication/v1"
)

// fakeSelfSubjectReviews implements SelfSubjectReviewInterface
type fakeSelfSubjectReviews struct {
	*gentype.FakeClient[*v1.SelfSubjectReview]
	Fake *FakeAuthenticationV1
}

func newFakeSelfSubjectReviews(fake *FakeAuthenticationV1) clientset.SelfSubjectReviewInterface {
	return &fakeSelfSubjectReviews{
		gentype.NewFakeClient[*v1.SelfSubjectReview](
			fake.Fake,
			"",
			v1.SchemeGroupVersion.WithResource("selfsubjectreviews"),
			v1.SchemeGroupVersion.WithKind("SelfSubjectReview"),
			func() *v1.SelfSubjectReview { return &v1.SelfSubjectReview{} },
		),
		fake,
	}
}
