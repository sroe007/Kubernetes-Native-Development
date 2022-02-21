package model

import (
	kubdevv1alpha1 "apress.com/m/v2/api/v1alpha1"
	"fmt"
	networkingv1 "k8s.io/api/networking/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func BackendIngress(cr *kubdevv1alpha1.LocalNewsApp) *networkingv1.Ingress {
	return &networkingv1.Ingress{
		ObjectMeta: v12.ObjectMeta{
			Name:      "news-backend",
			Namespace: cr.Namespace,
			Labels: map[string]string{
				"app": "news-backend",
			},
		},
		Spec: getBackendIngressSpec(cr),
	}
}

func getBackendIngressSpec(cr *kubdevv1alpha1.LocalNewsApp) networkingv1.IngressSpec {
	if cr.Spec.LocalNews.MinikubeIp == "" {
		cr.Spec.LocalNews.MinikubeIp = "fill-in-minikube-ip"
	}
	if cr.Spec.LocalNews.Domain == "" {
		cr.Spec.LocalNews.Domain = "nip.io"
	}
	if cr.Spec.NewsBackend.ServicePort < 1 {
		cr.Spec.NewsBackend.ServicePort = 8080
	}
	pathTypePrefix := networkingv1.PathTypePrefix
	spec := networkingv1.IngressSpec{
		Rules: []networkingv1.IngressRule{
			{
				Host: fmt.Sprintf("news-backend.%s.%s", cr.Spec.LocalNews.MinikubeIp, cr.Spec.LocalNews.Domain),
				IngressRuleValue: networkingv1.IngressRuleValue{
					HTTP: &networkingv1.HTTPIngressRuleValue{
						Paths: []networkingv1.HTTPIngressPath{
							{
								Path:     "/",
								PathType: &pathTypePrefix,
								Backend: networkingv1.IngressBackend{
									Service: &networkingv1.IngressServiceBackend{
										Name: "news-backend",
										Port: networkingv1.ServiceBackendPort{
											Number: cr.Spec.NewsBackend.ServicePort,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	return spec
}

func ReconcileBackendIngress(cr *kubdevv1alpha1.LocalNewsApp, currentState *networkingv1.Ingress) *networkingv1.Ingress {
	reconciled := currentState.DeepCopy()
	reconciled.Spec = getBackendIngressSpec(cr)
	return reconciled
}
