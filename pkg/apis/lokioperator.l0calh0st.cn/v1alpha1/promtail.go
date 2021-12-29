package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:defaulter-gen=true

// Promtail defines Promtail deployment
type Promtail struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PromtailSpec   `json:"spec"`
	Status PromtailStatus `json:"status"`
}

// PromtailSpec describes the specification of Promtail applications using kubernetes as a cluster manager
type PromtailSpec struct {
	Image     string `json:"image"`
	ConfigMap string `json:"configMap"`
	Config PromtailConfig `json:"config"`
}

// PromtailStatus describes the current status of Promtail applications
type PromtailStatus struct {
	// todo, write your code
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PromtailList carries a list of Promtail objects
type PromtailList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Promtail `json:"items"`
}
