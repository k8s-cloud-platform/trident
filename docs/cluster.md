# Cluster



提供不同的集群对接方式

- virtual-kubelet：部署到vk的节点
- OCM：部署到OCM集群



```go
type Cluster struct {
  metav1.TypeMeta   `json:",inline"`
  metav1.ObjectMeta `json:"metadata,omitempty"`

  Spec   ClusterSpec   `json:"spec,omitempty"`
  Status ClusterStatus `json:"status,omitempty"`
}

type ClusterSpec struct {
  Type string `json:"type,omitempty"`
}

type ClusterStatus struct {
}
```

