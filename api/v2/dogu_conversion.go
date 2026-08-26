package v2

import (
	"errors"
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/cloudogu/k8s-dogu-lib/v2/api/v3beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

// Compile-Time Interface Assertion
var _ conversion.Convertible = &Dogu{}

// These annotations are used to transport the v3beta1-only DoguApiVersion, Values... fields through a v2
// view of a Dogu without losing them, since v3beta1.Dogu is the conversion hub and v2 has no
// fields to hold them natively. They are a pure transport mechanism: ConvertFrom sets them,
// ConvertTo consumes and removes them again, so they never appear in a persisted v3beta1 object.
const (
	doguApiVersionAnnotationKey   = "k8s.cloudogu.com/v3beta1-doguApiVersion"
	valuesAnnotationKey           = "k8s.cloudogu.com/v3beta1-values"
	statusAppVersionAnnotationKey = "k8s.cloudogu.com/v3beta1-status-appVersion"
	doguNamespaceDelimiter        = "/"
)

// ConvertTo converts this v2 Dogu to the v3beta1 hub representation. This runs whenever a v2
// object is written (e.g., by an old client), since v3beta1 is the storage version.
func (d *Dogu) ConvertTo(dstRaw conversion.Hub) error {
	dst, ok := dstRaw.(*v3beta1.Dogu)
	if !ok {
		return fmt.Errorf("expected conversion target *v3beta1.Dogu, got %T", dstRaw)
	}

	dst.ObjectMeta = *d.ObjectMeta.DeepCopy()

	if strings.Count(d.Spec.Name, "/") != 1 {
		return errors.New(fmt.Sprintf("v2 dogu name %q should only contain one namespace delimiter %q", d.Spec.Name, doguNamespaceDelimiter))
	}

	doguNamespaceSplit := strings.Split(d.Spec.Name, doguNamespaceDelimiter)

	dst.Spec.DoguNamespace = doguNamespaceSplit[0]
	dst.Spec.Name = doguNamespaceSplit[1]
	dst.Spec.Version = d.Spec.Version
	dst.Spec.Resources = convertResourcesToV3beta1(d.Spec.Resources)
	dst.Spec.Security = convertSecurityToV3beta1(d.Spec.Security)
	dst.Spec.SupportMode = d.Spec.SupportMode
	dst.Spec.ExportMode = d.Spec.ExportMode
	dst.Spec.Stopped = d.Spec.Stopped
	dst.Spec.PauseReconciliation = d.Spec.PauseReconciliation
	dst.Spec.UpgradeConfig = v3beta1.UpgradeConfig{
		AllowNamespaceSwitch: d.Spec.UpgradeConfig.AllowNamespaceSwitch,
		ForceUpgrade:         d.Spec.UpgradeConfig.ForceUpgrade,
	}
	dst.Spec.AdditionalIngressAnnotations = maps.Clone(map[string]string(d.Spec.AdditionalIngressAnnotations))
	dst.Spec.AdditionalMounts = convertMountsToV3beta1(d.Spec.AdditionalMounts)

	// DoguApiVersion/Values only exist on the v3beta1 hub. Restore them from the transport annotations
	// if this v2 object still carries them (e.g., it originated from a v3beta1 object that was
	// viewed/round-tripped as v2); otherwise this is a genuinely v2-native object and defaults to
	// "v2". The CRD's +kubebuilder:default marker does not apply during conversion, so the
	// default must be set explicitly here.
	dst.Spec.DoguApiVersion = v3beta1.DoguApiVersionV2
	if doguApiVersion, foundApiVersion := d.Annotations[doguApiVersionAnnotationKey]; foundApiVersion {
		dst.Spec.DoguApiVersion = v3beta1.DoguApiVersion(doguApiVersion)
		delete(dst.Annotations, doguApiVersionAnnotationKey)
	}
	dst.Spec.Values = runtime.RawExtension{}
	if values, foundValues := d.Annotations[valuesAnnotationKey]; foundValues {
		dst.Spec.Values = runtime.RawExtension{Raw: []byte(values)}
		delete(dst.Annotations, valuesAnnotationKey)
	}

	dst.Status = convertStatusToV3beta1(d.Status, d.Annotations, dst.Annotations)
	return nil
}

// ConvertFrom converts the v3beta1 hub representation into this v2 Dogu. This runs whenever a
// Dogu stored as v3beta1 is served to a v2 client.
func (d *Dogu) ConvertFrom(srcRaw conversion.Hub) error {
	src, ok := srcRaw.(*v3beta1.Dogu)
	if !ok {
		return fmt.Errorf("expected conversion source *v3beta1.Dogu, got %T", srcRaw)
	}

	d.ObjectMeta = *src.ObjectMeta.DeepCopy()

	d.Spec.Name = fmt.Sprintf("%s%s%s", src.Spec.DoguNamespace, doguNamespaceDelimiter, src.Spec.Name)
	d.Spec.Version = src.Spec.Version
	d.Spec.Resources = convertResourcesFromV3beta1(src.Spec.Resources)
	d.Spec.Security = convertSecurityFromV3beta1(src.Spec.Security)
	d.Spec.SupportMode = src.Spec.SupportMode
	d.Spec.ExportMode = src.Spec.ExportMode
	d.Spec.Stopped = src.Spec.Stopped
	d.Spec.PauseReconciliation = src.Spec.PauseReconciliation
	d.Spec.UpgradeConfig = UpgradeConfig{
		AllowNamespaceSwitch: src.Spec.UpgradeConfig.AllowNamespaceSwitch,
		ForceUpgrade:         src.Spec.UpgradeConfig.ForceUpgrade,
	}
	d.Spec.AdditionalIngressAnnotations = maps.Clone(map[string]string(src.Spec.AdditionalIngressAnnotations))
	d.Spec.AdditionalMounts = convertMountsFromV3beta1(src.Spec.AdditionalMounts)

	// DoguApiVersion/Values have no field on v2. Stash them as transport annotations so a subsequent
	// v2->v3beta1 conversion (ConvertTo) can restore them, unless they're already at their
	// zero/default value, in which case a genuinely v2 dogu round-trips without annotation
	// clutter.
	if src.Spec.DoguApiVersion != "" && src.Spec.DoguApiVersion != v3beta1.DoguApiVersionV2 {
		if d.Annotations == nil {
			d.Annotations = map[string]string{}
		}
		d.Annotations[doguApiVersionAnnotationKey] = string(src.Spec.DoguApiVersion)
	}
	// TODO Check this values as annotation value
	if len(src.Spec.Values.Raw) > 0 {
		if d.Annotations == nil {
			d.Annotations = map[string]string{}
		}
		d.Annotations[valuesAnnotationKey] = string(src.Spec.Values.Raw)
	}

	if src.Status.AppVersion != "" {
		if d.Annotations == nil {
			d.Annotations = map[string]string{}
		}
		d.Annotations[statusAppVersionAnnotationKey] = src.Status.AppVersion
	}

	d.Status = convertStatusFromV3beta1(src.Status, d.Annotations)

	return nil
}

func convertResourcesToV3beta1(r DoguResources) v3beta1.DoguResources {
	return v3beta1.DoguResources{
		DataVolumeSize:    r.DataVolumeSize,
		MinDataVolumeSize: r.MinDataVolumeSize,
		StorageClassName:  clonePtr(r.StorageClassName),
	}
}

func convertResourcesFromV3beta1(r v3beta1.DoguResources) DoguResources {
	return DoguResources{
		DataVolumeSize:    r.DataVolumeSize,
		MinDataVolumeSize: r.MinDataVolumeSize,
		StorageClassName:  clonePtr(r.StorageClassName),
	}
}

func convertSecurityToV3beta1(s Security) v3beta1.Security {
	result := v3beta1.Security{
		Capabilities: v3beta1.Capabilities{
			Add:  slices.Clone(s.Capabilities.Add),
			Drop: slices.Clone(s.Capabilities.Drop),
		},
		RunAsNonRoot:           clonePtr(s.RunAsNonRoot),
		ReadOnlyRootFileSystem: clonePtr(s.ReadOnlyRootFileSystem),
	}
	if s.SELinuxOptions != nil {
		result.SELinuxOptions = &v3beta1.SELinuxOptions{
			User:  s.SELinuxOptions.User,
			Role:  s.SELinuxOptions.Role,
			Type:  s.SELinuxOptions.Type,
			Level: s.SELinuxOptions.Level,
		}
	}
	if s.SeccompProfile != nil {
		result.SeccompProfile = &v3beta1.SeccompProfile{
			Type:             v3beta1.SeccompProfileType(s.SeccompProfile.Type),
			LocalhostProfile: clonePtr(s.SeccompProfile.LocalhostProfile),
		}
	}
	if s.AppArmorProfile != nil {
		result.AppArmorProfile = &v3beta1.AppArmorProfile{
			Type:             v3beta1.AppArmorProfileType(s.AppArmorProfile.Type),
			LocalhostProfile: clonePtr(s.AppArmorProfile.LocalhostProfile),
		}
	}
	return result
}

func convertSecurityFromV3beta1(s v3beta1.Security) Security {
	result := Security{
		Capabilities: Capabilities{
			Add:  slices.Clone(s.Capabilities.Add),
			Drop: slices.Clone(s.Capabilities.Drop),
		},
		RunAsNonRoot:           clonePtr(s.RunAsNonRoot),
		ReadOnlyRootFileSystem: clonePtr(s.ReadOnlyRootFileSystem),
	}
	if s.SELinuxOptions != nil {
		result.SELinuxOptions = &SELinuxOptions{
			User:  s.SELinuxOptions.User,
			Role:  s.SELinuxOptions.Role,
			Type:  s.SELinuxOptions.Type,
			Level: s.SELinuxOptions.Level,
		}
	}
	if s.SeccompProfile != nil {
		result.SeccompProfile = &SeccompProfile{
			Type:             SeccompProfileType(s.SeccompProfile.Type),
			LocalhostProfile: clonePtr(s.SeccompProfile.LocalhostProfile),
		}
	}
	if s.AppArmorProfile != nil {
		result.AppArmorProfile = &AppArmorProfile{
			Type:             AppArmorProfileType(s.AppArmorProfile.Type),
			LocalhostProfile: clonePtr(s.AppArmorProfile.LocalhostProfile),
		}
	}
	return result
}

func convertMountsToV3beta1(mounts []DataMount) []v3beta1.DataMount {
	if mounts == nil {
		return nil
	}
	result := make([]v3beta1.DataMount, len(mounts))
	for i, m := range mounts {
		result[i] = v3beta1.DataMount{
			SourceType: v3beta1.DataSourceType(m.SourceType),
			Name:       m.Name,
			Volume:     m.Volume,
			Subfolder:  m.Subfolder,
		}
	}
	return result
}

func convertMountsFromV3beta1(mounts []v3beta1.DataMount) []DataMount {
	if mounts == nil {
		return nil
	}
	result := make([]DataMount, len(mounts))
	for i, m := range mounts {
		result[i] = DataMount{
			SourceType: DataSourceType(m.SourceType),
			Name:       m.Name,
			Volume:     m.Volume,
			Subfolder:  m.Subfolder,
		}
	}
	return result
}

func convertStatusToV3beta1(s DoguStatus, srcAnnotations map[string]string, dstAnnotations map[string]string) v3beta1.DoguStatus {
	status := v3beta1.DoguStatus{
		Status:           s.Status,
		RequeueTime:      s.RequeueTime,
		RequeuePhase:     s.RequeuePhase,
		StartedAt:        s.StartedAt,
		Health:           v3beta1.HealthStatus(s.Health),
		InstalledVersion: s.InstalledVersion,
		Stopped:          s.Stopped,
		ExportMode:       s.ExportMode,
		DataVolumeSize:   clonePtr(s.DataVolumeSize),
		Conditions:       slices.Clone(s.Conditions),
	}

	if appVersion, foundAppVersion := srcAnnotations[statusAppVersionAnnotationKey]; foundAppVersion {
		status.AppVersion = appVersion
		delete(dstAnnotations, statusAppVersionAnnotationKey)
	}

	return status
}

func convertStatusFromV3beta1(s v3beta1.DoguStatus, annotations map[string]string) DoguStatus {
	status := DoguStatus{
		Status:           s.Status,
		RequeueTime:      s.RequeueTime,
		RequeuePhase:     s.RequeuePhase,
		StartedAt:        s.StartedAt,
		Health:           HealthStatus(s.Health),
		InstalledVersion: s.InstalledVersion,
		Stopped:          s.Stopped,
		ExportMode:       s.ExportMode,
		DataVolumeSize:   clonePtr(s.DataVolumeSize),
		Conditions:       slices.Clone(s.Conditions),
	}

	if s.AppVersion != "" {
		if annotations == nil {
			annotations = map[string]string{}
		}
		annotations[statusAppVersionAnnotationKey] = s.AppVersion
	}

	return status
}

// clonePtr returns a pointer to a copy of the value v points to, or nil if v is nil. It avoids
// aliasing a pointer field between the source and destination of a conversion.
func clonePtr[T any](v *T) *T {
	if v == nil {
		return nil
	}
	return new(*v)
}
