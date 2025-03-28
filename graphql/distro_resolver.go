package graphql

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/evergreen-ci/birch"
	"github.com/evergreen-ci/evergreen"
	"github.com/evergreen-ci/evergreen/model/distro"
	"github.com/evergreen-ci/evergreen/rest/model"
	"github.com/evergreen-ci/utility"
)

// Communication is the resolver for the communication field.
func (r *bootstrapSettingsResolver) Communication(ctx context.Context, obj *model.APIBootstrapSettings) (CommunicationMethod, error) {
	if obj == nil {
		return "", InternalServerError.Send(ctx, "distro undefined when attempting to resolve communication method")
	}

	switch utility.FromStringPtr(obj.Communication) {
	case distro.CommunicationMethodLegacySSH:
		return CommunicationMethodLegacySSH, nil
	case distro.CommunicationMethodSSH:
		return CommunicationMethodSSH, nil
	case distro.CommunicationMethodRPC:
		return CommunicationMethodRPC, nil
	default:
		return "", InputValidationError.Send(ctx, fmt.Sprintf("communication method '%s' is invalid", utility.FromStringPtr(obj.Communication)))
	}
}

// Method is the resolver for the method field.
func (r *bootstrapSettingsResolver) Method(ctx context.Context, obj *model.APIBootstrapSettings) (BootstrapMethod, error) {
	if obj == nil {
		return "", InternalServerError.Send(ctx, "distro undefined when attempting to resolve bootstrap method")
	}

	switch utility.FromStringPtr(obj.Method) {
	case distro.BootstrapMethodLegacySSH:
		return BootstrapMethodLegacySSH, nil
	case distro.BootstrapMethodSSH:
		return BootstrapMethodSSH, nil
	case distro.BootstrapMethodUserData:
		return BootstrapMethodUserData, nil
	default:
		return "", InputValidationError.Send(ctx, fmt.Sprintf("bootstrap method '%s' is invalid", utility.FromStringPtr(obj.Method)))
	}
}

// Version is the resolver for the version field.
func (r *dispatcherSettingsResolver) Version(ctx context.Context, obj *model.APIDispatcherSettings) (DispatcherVersion, error) {
	if obj == nil {
		return "", InternalServerError.Send(ctx, "distro undefined when attempting to resolve dispatcher version")
	}

	switch utility.FromStringPtr(obj.Version) {
	case evergreen.DispatcherVersionRevisedWithDependencies:
		return DispatcherVersionRevisedWithDependencies, nil
	default:
		return "", InputValidationError.Send(ctx, fmt.Sprintf("dispatcher version '%s' is invalid", utility.FromStringPtr(obj.Version)))
	}
}

// Arch is the resolver for the arch field.
func (r *distroResolver) Arch(ctx context.Context, obj *model.APIDistro) (Arch, error) {
	if obj == nil {
		return "", InternalServerError.Send(ctx, "distro undefined when attempting to resolve arch")
	}

	switch utility.FromStringPtr(obj.Arch) {
	case evergreen.ArchDarwinAmd64:
		return ArchOsx64Bit, nil
	case evergreen.ArchDarwinArm64:
		return ArchOsxArm64Bit, nil
	case evergreen.ArchLinuxAmd64:
		return ArchLinux64Bit, nil
	case evergreen.ArchLinuxArm64:
		return ArchLinuxArm64Bit, nil
	case evergreen.ArchLinuxPpc64le:
		return ArchLinuxPpc64Bit, nil
	case evergreen.ArchLinuxS390x:
		return ArchLinuxZseries, nil
	case evergreen.ArchWindowsAmd64:
		return ArchWindows64Bit, nil
	default:
		return "", InputValidationError.Send(ctx, fmt.Sprintf("arch '%s' is invalid", utility.FromStringPtr(obj.Arch)))
	}
}

// Provider is the resolver for the provider field.
func (r *distroResolver) Provider(ctx context.Context, obj *model.APIDistro) (Provider, error) {
	if obj == nil {
		return "", InternalServerError.Send(ctx, "distro undefined when attempting to resolve provider")
	}

	switch utility.FromStringPtr(obj.Provider) {
	case evergreen.ProviderNameDocker:
		return ProviderDocker, nil
	case evergreen.ProviderNameEc2Fleet:
		return ProviderEc2Fleet, nil
	case evergreen.ProviderNameEc2OnDemand:
		return ProviderEc2OnDemand, nil
	case evergreen.ProviderNameStatic:
		return ProviderStatic, nil
	default:
		return "", InputValidationError.Send(ctx, fmt.Sprintf("provider '%s' is invalid", utility.FromStringPtr(obj.Provider)))
	}
}

// ProviderSettingsList is the resolver for the providerSettingsList field.
func (r *distroResolver) ProviderSettingsList(ctx context.Context, obj *model.APIDistro) ([]map[string]any, error) {
	settings := []map[string]any{}
	for _, entry := range obj.ProviderSettingsList {
		settings = append(settings, entry.ExportMap())
	}

	return settings, nil
}

// Version is the resolver for the version field.
func (r *finderSettingsResolver) Version(ctx context.Context, obj *model.APIFinderSettings) (FinderVersion, error) {
	if obj == nil {
		return "", InternalServerError.Send(ctx, "distro undefined when attempting to resolve finder version")
	}

	switch utility.FromStringPtr(obj.Version) {
	case evergreen.FinderVersionLegacy:
		return FinderVersionLegacy, nil
	case evergreen.FinderVersionParallel:
		return FinderVersionParallel, nil
	case evergreen.FinderVersionPipeline:
		return FinderVersionPipeline, nil
	case evergreen.FinderVersionAlternate:
		return FinderVersionAlternate, nil
	default:
		return "", InputValidationError.Send(ctx, fmt.Sprintf("finder version '%s' is invalid", utility.FromStringPtr(obj.Version)))
	}
}

// FeedbackRule is the resolver for the feedbackRule field.
func (r *hostAllocatorSettingsResolver) FeedbackRule(ctx context.Context, obj *model.APIHostAllocatorSettings) (FeedbackRule, error) {
	if obj == nil {
		return "", InternalServerError.Send(ctx, "distro undefined when attempting to resolve feedback rule")
	}

	switch utility.FromStringPtr(obj.FeedbackRule) {
	case evergreen.HostAllocatorWaitsOverThreshFeedback:
		return FeedbackRuleWaitsOverThresh, nil
	case evergreen.HostAllocatorNoFeedback:
		return FeedbackRuleNoFeedback, nil
	case evergreen.HostAllocatorUseDefaultFeedback:
		return FeedbackRuleDefault, nil
	default:
		return "", InputValidationError.Send(ctx, fmt.Sprintf("feedback rule '%s' is invalid", utility.FromStringPtr(obj.FeedbackRule)))
	}
}

// HostsOverallocatedRule is the resolver for the hostsOverallocatedRule field.
func (r *hostAllocatorSettingsResolver) HostsOverallocatedRule(ctx context.Context, obj *model.APIHostAllocatorSettings) (OverallocatedRule, error) {
	if obj == nil {
		return "", InternalServerError.Send(ctx, "distro undefined when attempting to resolve overallocated rule")
	}

	switch utility.FromStringPtr(obj.HostsOverallocatedRule) {
	case evergreen.HostsOverallocatedTerminate:
		return OverallocatedRuleTerminate, nil
	case evergreen.HostsOverallocatedIgnore:
		return OverallocatedRuleIgnore, nil
	case evergreen.HostsOverallocatedUseDefault:
		return OverallocatedRuleDefault, nil
	default:
		return "", InputValidationError.Send(ctx, fmt.Sprintf("overallocated rule '%s' is invalid", utility.FromStringPtr(obj.HostsOverallocatedRule)))
	}
}

// RoundingRule is the resolver for the roundingRule field.
func (r *hostAllocatorSettingsResolver) RoundingRule(ctx context.Context, obj *model.APIHostAllocatorSettings) (RoundingRule, error) {
	if obj == nil {
		return "", InternalServerError.Send(ctx, "distro undefined when attempting to resolve rounding rule")
	}

	switch utility.FromStringPtr(obj.RoundingRule) {
	case evergreen.HostAllocatorRoundDown:
		return RoundingRuleDown, nil
	case evergreen.HostAllocatorRoundUp:
		return RoundingRuleUp, nil
	case evergreen.HostAllocatorRoundDefault:
		return RoundingRuleDefault, nil
	default:
		return "", InputValidationError.Send(ctx, fmt.Sprintf("rounding rule '%s' is invalid", utility.FromStringPtr(obj.RoundingRule)))
	}
}

// Version is the resolver for the version field.
func (r *hostAllocatorSettingsResolver) Version(ctx context.Context, obj *model.APIHostAllocatorSettings) (HostAllocatorVersion, error) {
	if obj == nil {
		return "", InternalServerError.Send(ctx, "distro undefined when attempting to resolve host allocator version")
	}

	switch utility.FromStringPtr(obj.Version) {
	case evergreen.HostAllocatorUtilization:
		return HostAllocatorVersionUtilization, nil
	default:
		return "", InputValidationError.Send(ctx, fmt.Sprintf("host allocator version '%s' is invalid", utility.FromStringPtr(obj.Version)))
	}
}

// Version is the resolver for the version field.
func (r *plannerSettingsResolver) Version(ctx context.Context, obj *model.APIPlannerSettings) (PlannerVersion, error) {
	if obj == nil {
		return "", InternalServerError.Send(ctx, "distro undefined when attempting to resolve planner version")
	}

	switch utility.FromStringPtr(obj.Version) {
	case evergreen.PlannerVersionTunable:
		return PlannerVersionTunable, nil
	default:
		return "", InputValidationError.Send(ctx, fmt.Sprintf("planner version '%s' is invalid", utility.FromStringPtr(obj.Version)))
	}
}

// Communication is the resolver for the communication field.
func (r *bootstrapSettingsInputResolver) Communication(ctx context.Context, obj *model.APIBootstrapSettings, data CommunicationMethod) error {
	switch data {
	case CommunicationMethodLegacySSH:
		obj.Communication = utility.ToStringPtr(distro.CommunicationMethodLegacySSH)
	case CommunicationMethodSSH:
		obj.Communication = utility.ToStringPtr(distro.CommunicationMethodSSH)
	case CommunicationMethodRPC:
		obj.Communication = utility.ToStringPtr(distro.CommunicationMethodRPC)
	default:
		return InputValidationError.Send(ctx, fmt.Sprintf("communication method '%s' is invalid", data))
	}
	return nil
}

// Method is the resolver for the method field.
func (r *bootstrapSettingsInputResolver) Method(ctx context.Context, obj *model.APIBootstrapSettings, data BootstrapMethod) error {
	switch data {
	case BootstrapMethodLegacySSH:
		obj.Method = utility.ToStringPtr(distro.BootstrapMethodLegacySSH)
	case BootstrapMethodSSH:
		obj.Method = utility.ToStringPtr(distro.BootstrapMethodSSH)
	case BootstrapMethodUserData:
		obj.Method = utility.ToStringPtr(distro.BootstrapMethodUserData)
	default:
		return InputValidationError.Send(ctx, fmt.Sprintf("bootstrap method '%s' is invalid", data))
	}
	return nil
}

// Version is the resolver for the version field.
func (r *dispatcherSettingsInputResolver) Version(ctx context.Context, obj *model.APIDispatcherSettings, data DispatcherVersion) error {
	switch data {
	case DispatcherVersionRevisedWithDependencies:
		obj.Version = utility.ToStringPtr(evergreen.DispatcherVersionRevisedWithDependencies)
	default:
		return InputValidationError.Send(ctx, fmt.Sprintf("dispatcher version '%s' is invalid", data))
	}
	return nil
}

// Arch is the resolver for the arch field.
func (r *distroInputResolver) Arch(ctx context.Context, obj *model.APIDistro, data Arch) error {
	switch data {
	case ArchOsx64Bit:
		obj.Arch = utility.ToStringPtr(evergreen.ArchDarwinAmd64)
	case ArchOsxArm64Bit:
		obj.Arch = utility.ToStringPtr(evergreen.ArchDarwinArm64)
	case ArchLinux64Bit:
		obj.Arch = utility.ToStringPtr(evergreen.ArchLinuxAmd64)
	case ArchLinuxArm64Bit:
		obj.Arch = utility.ToStringPtr(evergreen.ArchLinuxArm64)
	case ArchLinuxPpc64Bit:
		obj.Arch = utility.ToStringPtr(evergreen.ArchLinuxPpc64le)
	case ArchLinuxZseries:
		obj.Arch = utility.ToStringPtr(evergreen.ArchLinuxS390x)
	case ArchWindows64Bit:
		obj.Arch = utility.ToStringPtr(evergreen.ArchWindowsAmd64)
	default:
		return InputValidationError.Send(ctx, fmt.Sprintf("arch '%s' is invalid", data))
	}
	return nil
}

// Provider is the resolver for the provider field.
func (r *distroInputResolver) Provider(ctx context.Context, obj *model.APIDistro, data Provider) error {
	switch data {
	case ProviderDocker:
		obj.Provider = utility.ToStringPtr(evergreen.ProviderNameDocker)
	case ProviderEc2Fleet:
		obj.Provider = utility.ToStringPtr(evergreen.ProviderNameEc2Fleet)
	case ProviderEc2OnDemand:
		obj.Provider = utility.ToStringPtr(evergreen.ProviderNameEc2OnDemand)
	case ProviderStatic:
		obj.Provider = utility.ToStringPtr(evergreen.ProviderNameStatic)
	default:
		return InputValidationError.Send(ctx, fmt.Sprintf("provider '%s' is invalid", data))
	}
	return nil
}

// ProviderSettingsList is the resolver for the providerSettingsList field.
func (r *distroInputResolver) ProviderSettingsList(ctx context.Context, obj *model.APIDistro, data []map[string]any) error {
	settings := []*birch.Document{}
	for _, entry := range data {
		newEntry, err := json.Marshal(entry)
		if err != nil {
			return InternalServerError.Send(ctx, fmt.Sprintf("marshalling provider settings entry: %s", err.Error()))
		}
		doc := &birch.Document{}
		if err = json.Unmarshal(newEntry, doc); err != nil {
			return InternalServerError.Send(ctx, fmt.Sprintf("converting map to birch: %s", err.Error()))
		}
		settings = append(settings, doc)
	}
	obj.ProviderSettingsList = settings
	return nil
}

// Version is the resolver for the version field.
func (r *finderSettingsInputResolver) Version(ctx context.Context, obj *model.APIFinderSettings, data FinderVersion) error {
	switch data {
	case FinderVersionLegacy:
		obj.Version = utility.ToStringPtr(evergreen.FinderVersionLegacy)
	case FinderVersionParallel:
		obj.Version = utility.ToStringPtr(evergreen.FinderVersionParallel)
	case FinderVersionPipeline:
		obj.Version = utility.ToStringPtr(evergreen.FinderVersionPipeline)
	case FinderVersionAlternate:
		obj.Version = utility.ToStringPtr(evergreen.FinderVersionAlternate)
	default:
		return InputValidationError.Send(ctx, fmt.Sprintf("finder version '%s' is invalid", data))
	}
	return nil
}

// AcceptableHostIdleTime is the resolver for the acceptableHostIdleTime field.
func (r *hostAllocatorSettingsInputResolver) AcceptableHostIdleTime(ctx context.Context, obj *model.APIHostAllocatorSettings, data int) error {
	obj.AcceptableHostIdleTime = model.NewAPIDuration(time.Duration(data) * time.Millisecond)
	return nil
}

// FeedbackRule is the resolver for the feedbackRule field.
func (r *hostAllocatorSettingsInputResolver) FeedbackRule(ctx context.Context, obj *model.APIHostAllocatorSettings, data FeedbackRule) error {
	switch data {
	case FeedbackRuleWaitsOverThresh:
		obj.FeedbackRule = utility.ToStringPtr(evergreen.HostAllocatorWaitsOverThreshFeedback)
	case FeedbackRuleNoFeedback:
		obj.FeedbackRule = utility.ToStringPtr(evergreen.HostAllocatorNoFeedback)
	case FeedbackRuleDefault:
		obj.FeedbackRule = utility.ToStringPtr(evergreen.HostAllocatorUseDefaultFeedback)
	default:
		return InputValidationError.Send(ctx, fmt.Sprintf("feedback rule '%s' is invalid", data))
	}
	return nil
}

// HostsOverallocatedRule is the resolver for the hostsOverallocatedRule field.
func (r *hostAllocatorSettingsInputResolver) HostsOverallocatedRule(ctx context.Context, obj *model.APIHostAllocatorSettings, data OverallocatedRule) error {
	switch data {
	case OverallocatedRuleTerminate:
		obj.HostsOverallocatedRule = utility.ToStringPtr(evergreen.HostsOverallocatedTerminate)
	case OverallocatedRuleIgnore:
		obj.HostsOverallocatedRule = utility.ToStringPtr(evergreen.HostsOverallocatedIgnore)
	case OverallocatedRuleDefault:
		obj.HostsOverallocatedRule = utility.ToStringPtr(evergreen.HostsOverallocatedUseDefault)
	default:
		return InputValidationError.Send(ctx, fmt.Sprintf("overallocated rule '%s' is invalid", data))
	}
	return nil
}

// RoundingRule is the resolver for the roundingRule field.
func (r *hostAllocatorSettingsInputResolver) RoundingRule(ctx context.Context, obj *model.APIHostAllocatorSettings, data RoundingRule) error {
	switch data {
	case RoundingRuleDown:
		obj.RoundingRule = utility.ToStringPtr(evergreen.HostAllocatorRoundDown)
	case RoundingRuleUp:
		obj.RoundingRule = utility.ToStringPtr(evergreen.HostAllocatorRoundUp)
	case RoundingRuleDefault:
		obj.RoundingRule = utility.ToStringPtr(evergreen.HostAllocatorRoundDefault)
	default:
		return InputValidationError.Send(ctx, fmt.Sprintf("rounding rule '%s' is invalid", data))
	}
	return nil
}

// Version is the resolver for the version field.
func (r *hostAllocatorSettingsInputResolver) Version(ctx context.Context, obj *model.APIHostAllocatorSettings, data HostAllocatorVersion) error {
	switch data {
	case HostAllocatorVersionUtilization:
		obj.Version = utility.ToStringPtr(evergreen.HostAllocatorUtilization)
	default:
		return InputValidationError.Send(ctx, fmt.Sprintf("host allocator version '%s' is invalid", data))
	}
	return nil
}

// TargetTime is the resolver for the targetTime field.
func (r *plannerSettingsInputResolver) TargetTime(ctx context.Context, obj *model.APIPlannerSettings, data int) error {
	obj.TargetTime = model.NewAPIDuration(time.Duration(data) * time.Millisecond)
	return nil
}

// Version is the resolver for the version field.
func (r *plannerSettingsInputResolver) Version(ctx context.Context, obj *model.APIPlannerSettings, data PlannerVersion) error {
	switch data {
	case PlannerVersionTunable:
		obj.Version = utility.ToStringPtr(evergreen.PlannerVersionTunable)
	default:
		return InputValidationError.Send(ctx, fmt.Sprintf("planner version '%s' is invalid", data))
	}
	return nil
}

// BootstrapSettings returns BootstrapSettingsResolver implementation.
func (r *Resolver) BootstrapSettings() BootstrapSettingsResolver {
	return &bootstrapSettingsResolver{r}
}

// DispatcherSettings returns DispatcherSettingsResolver implementation.
func (r *Resolver) DispatcherSettings() DispatcherSettingsResolver {
	return &dispatcherSettingsResolver{r}
}

// Distro returns DistroResolver implementation.
func (r *Resolver) Distro() DistroResolver { return &distroResolver{r} }

// FinderSettings returns FinderSettingsResolver implementation.
func (r *Resolver) FinderSettings() FinderSettingsResolver { return &finderSettingsResolver{r} }

// HostAllocatorSettings returns HostAllocatorSettingsResolver implementation.
func (r *Resolver) HostAllocatorSettings() HostAllocatorSettingsResolver {
	return &hostAllocatorSettingsResolver{r}
}

// PlannerSettings returns PlannerSettingsResolver implementation.
func (r *Resolver) PlannerSettings() PlannerSettingsResolver { return &plannerSettingsResolver{r} }

// BootstrapSettingsInput returns BootstrapSettingsInputResolver implementation.
func (r *Resolver) BootstrapSettingsInput() BootstrapSettingsInputResolver {
	return &bootstrapSettingsInputResolver{r}
}

// DispatcherSettingsInput returns DispatcherSettingsInputResolver implementation.
func (r *Resolver) DispatcherSettingsInput() DispatcherSettingsInputResolver {
	return &dispatcherSettingsInputResolver{r}
}

// DistroInput returns DistroInputResolver implementation.
func (r *Resolver) DistroInput() DistroInputResolver { return &distroInputResolver{r} }

// FinderSettingsInput returns FinderSettingsInputResolver implementation.
func (r *Resolver) FinderSettingsInput() FinderSettingsInputResolver {
	return &finderSettingsInputResolver{r}
}

// HostAllocatorSettingsInput returns HostAllocatorSettingsInputResolver implementation.
func (r *Resolver) HostAllocatorSettingsInput() HostAllocatorSettingsInputResolver {
	return &hostAllocatorSettingsInputResolver{r}
}

// PlannerSettingsInput returns PlannerSettingsInputResolver implementation.
func (r *Resolver) PlannerSettingsInput() PlannerSettingsInputResolver {
	return &plannerSettingsInputResolver{r}
}

type bootstrapSettingsResolver struct{ *Resolver }
type dispatcherSettingsResolver struct{ *Resolver }
type distroResolver struct{ *Resolver }
type finderSettingsResolver struct{ *Resolver }
type hostAllocatorSettingsResolver struct{ *Resolver }
type plannerSettingsResolver struct{ *Resolver }
type bootstrapSettingsInputResolver struct{ *Resolver }
type dispatcherSettingsInputResolver struct{ *Resolver }
type distroInputResolver struct{ *Resolver }
type finderSettingsInputResolver struct{ *Resolver }
type hostAllocatorSettingsInputResolver struct{ *Resolver }
type plannerSettingsInputResolver struct{ *Resolver }
