# Tenant



提供不同的租户隔离方式

- 强隔离：独立apiserver、controller-manager、etcd

- 弱隔离：namespace隔离，共享apiserver、controller-manager、etcd



```go
type Tenant struct {
  metav1.TypeMeta   `json:",inline"`
  metav1.ObjectMeta `json:"metadata,omitempty"`

  Spec   TenantSpec   `json:"spec,omitempty"`
  Status TenantStatus `json:"status,omitempty"`
}

type TenantSpec struct {
  IsolationLevel string `json:"isolationLevel,omitempty"`
}

type TenantStatus struct {
  Phase      string             `json:"phase,omitempty"`
  Conditions []metav1.Condition `json:"conditions,omitempty"`
}
```

