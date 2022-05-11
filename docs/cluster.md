# Cluster



提供不同的集群对接方式

- host：部署到宿主集群
- virtual-kubelet：部署到vk的节点
- OCM：部署到OCM worker集群



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

