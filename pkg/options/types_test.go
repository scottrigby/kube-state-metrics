/*
Copyright 2018 The Kubernetes Authors All rights reserved.

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

package options

import (
	"reflect"
	"testing"
)

func TestResourceSetSet(t *testing.T) {
	tests := []struct {
		Desc        string
		Value       string
		Wanted      ResourceSet
		WantedError bool
	}{
		{
			Desc:        "empty resources",
			Value:       "",
			Wanted:      ResourceSet{},
			WantedError: false,
		},
		{
			Desc:  "normal resources",
			Value: "configmaps,cronjobs,daemonsets,deployments",
			Wanted: ResourceSet(map[string]struct{}{
				"configmaps":  {},
				"cronjobs":    {},
				"daemonsets":  {},
				"deployments": {},
			}),
			WantedError: false,
		},
	}

	for _, test := range tests {
		cs := &ResourceSet{}
		gotError := cs.Set(test.Value)
		if !(((gotError == nil && !test.WantedError) || (gotError != nil && test.WantedError)) && reflect.DeepEqual(*cs, test.Wanted)) {
			t.Errorf("Test error for Desc: %s. Want: %+v. Got: %+v. Wanted Error: %v, Got Error: %v", test.Desc, test.Wanted, *cs, test.WantedError, gotError)
		}
	}
}

func TestNamespaceListSet(t *testing.T) {
	tests := []struct {
		Desc   string
		Value  string
		Wanted NamespaceList
	}{
		{
			Desc:   "empty namespacelist",
			Value:  "",
			Wanted: NamespaceList{},
		},
		{
			Desc:  "normal namespacelist",
			Value: "default, kube-system",
			Wanted: NamespaceList([]string{
				"default",
				"kube-system",
			}),
		},
	}

	for _, test := range tests {
		ns := &NamespaceList{}
		gotError := ns.Set(test.Value)
		if gotError != nil || !reflect.DeepEqual(*ns, test.Wanted) {
			t.Errorf("Test error for Desc: %s. Want: %+v. Got: %+v. Got Error: %v", test.Desc, test.Wanted, *ns, gotError)
		}
	}
}

func TestMetricSetSet(t *testing.T) {
	tests := []struct {
		Desc   string
		Value  string
		Wanted MetricSet
	}{
		{
			Desc:   "empty metrics",
			Value:  "",
			Wanted: MetricSet{},
		},
		{
			Desc:  "normal metrics",
			Value: "kube_cronjob_info, kube_cronjob_labels, kube_daemonset_labels",
			Wanted: MetricSet(map[string]struct{}{
				"kube_cronjob_info":     {},
				"kube_cronjob_labels":   {},
				"kube_daemonset_labels": {},
			}),
		},
	}

	for _, test := range tests {
		ms := &MetricSet{}
		gotError := ms.Set(test.Value)
		if gotError != nil || !reflect.DeepEqual(*ms, test.Wanted) {
			t.Errorf("Test error for Desc: %s. Want: %+v. Got: %+v. Got Error: %v", test.Desc, test.Wanted, *ms, gotError)
		}
	}
}

func TestLabelsAllowListSet(t *testing.T) {
	tests := []struct {
		Desc   string
		Value  string
		Wanted LabelsAllowList
		err    bool
	}{
		{
			Desc:   "empty metrics",
			Value:  "",
			Wanted: LabelsAllowList{},
		},
		{
			Desc:  "Deny all labels for metric",
			Value: "kube_cronjob_info=[]",
			Wanted: LabelsAllowList{
				"kube_cronjob_info": {},
			},
		},
		{
			Desc:   "[invalid] space delimited",
			Value:  "kube_cronjob_info=[somelabel,label2] kube_cronjob_labels=[label3,label4]",
			Wanted: LabelsAllowList(map[string][]string{}),
			err:    true,
		},
		{
			Desc:   "[invalid] normal missing bracket",
			Value:  "kube_cronjob_info=[somelabel,label2],kube_cronjob_labels=label3,label4]",
			Wanted: LabelsAllowList(map[string][]string{}),
			err:    true,
		},

		{
			Desc:   "[invalid] no comma between metrics",
			Value:  "kube_cronjob_info=[somelabel,label2]kube_cronjob_labels=[label3,label4]",
			Wanted: LabelsAllowList(map[string][]string{}),
			err:    true,
		},
		{
			Desc:   "[invalid] no '=' between name and label list",
			Value:  "kube_cronjob_info[somelabel,label2]kube_cronjob_labels=[label3,label4]",
			Wanted: LabelsAllowList(map[string][]string{}),
			err:    true,
		},
		{
			Desc:  "normal metrics",
			Value: "kube_cronjob_info=[somelabel,label2]",
			Wanted: LabelsAllowList(map[string][]string{
				"kube_cronjob_info": {
					"somelabel",
					"label2",
				}}),
		},
		{
			Desc:  "normal metrics",
			Value: "kube_cronjob_info=[somelabel,label2],kube_cronjob_count=[somelabel,label2],kube_cronjob_desc=[somelabel,label2]",
			Wanted: LabelsAllowList(map[string][]string{
				"kube_cronjob_info": {
					"somelabel",
					"label2"},
				"kube_cronjob_count": {
					"somelabel",
					"label2"},
				"kube_cronjob_desc": {
					"somelabel",
					"label2"}}),
		},
		{
			Desc:  "normal metrics, with empty allow labels",
			Value: "kube_cronjob_info=[somelabel,label2],kube_cronjob_count=[]",
			Wanted: LabelsAllowList(map[string][]string{
				"kube_cronjob_info": {
					"somelabel",
					"label2",
				},
				"kube_cronjob_count": {}}),
		},
	}

	for _, test := range tests {
		lal := &LabelsAllowList{}
		gotError := lal.Set(test.Value)
		if gotError != nil && !test.err || !reflect.DeepEqual(*lal, test.Wanted) {
			t.Errorf("Test error for Desc: %s. Want: %+v. Got: %+v. Got Error: %v", test.Desc, test.Wanted, *lal, gotError)
		}
	}
}
