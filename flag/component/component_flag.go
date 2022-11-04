package component

import "github.com/entgigi/upgrade-operator.git/api/v1alpha1"

type ComponentFlag struct {
	ComponentName string
	Flag          string
	// returns the reference to the related ImageOverride field in EntandoAppV2
	ImageOverrideGetter func(entandoApp *v1alpha1.EntandoAppV2) *string
}

var ComponentFlags = []ComponentFlag{
	mkComponentFlag("DeApp", "image-de-app", func(entandoApp *v1alpha1.EntandoAppV2) *string {
		return &entandoApp.Spec.DeApp.ImageOverride
	}),
	mkComponentFlag("AppBuilder", "image-app-builder", func(entandoApp *v1alpha1.EntandoAppV2) *string {
		return &entandoApp.Spec.AppBuilder.ImageOverride
	}),
	mkComponentFlag("ComponentManager", "image-component-manager", func(entandoApp *v1alpha1.EntandoAppV2) *string {
		return &entandoApp.Spec.ComponentManager.ImageOverride
	}),
	mkComponentFlag("Keycloak", "image-keycloak", func(entandoApp *v1alpha1.EntandoAppV2) *string {
		return &entandoApp.Spec.Keycloak.ImageOverride
	}),
	mkComponentFlag("K8sService", "image-k8s-service", func(entandoApp *v1alpha1.EntandoAppV2) *string {
		return &entandoApp.Spec.K8sService.ImageOverride
	}),
	mkComponentFlag("K8sPluginController", "image-k8s-plugin-controller", func(entandoApp *v1alpha1.EntandoAppV2) *string {
		return &entandoApp.Spec.K8sPluginController.ImageOverride
	}),
	mkComponentFlag("K8sAppPluginLinkController", "image-k8s-app-plugin-link-controller", func(entandoApp *v1alpha1.EntandoAppV2) *string {
		return &entandoApp.Spec.K8sAppPluginLinkController.ImageOverride
	}),
}

func mkComponentFlag(name, flag string, getter func(entandoApp *v1alpha1.EntandoAppV2) *string) ComponentFlag {
	return ComponentFlag{
		ComponentName:       name,
		Flag:                flag,
		ImageOverrideGetter: getter,
	}
}
