package model

import (
	"fmt"
	"reflect"
	"time"

	"github.com/evergreen-ci/evergreen"
	"github.com/evergreen-ci/evergreen/model"
	"github.com/evergreen-ci/evergreen/model/patch"
	"github.com/evergreen-ci/evergreen/util"
	"github.com/evergreen-ci/utility"
	"github.com/mongodb/grip"
	"github.com/mongodb/grip/recovery"
	"github.com/pkg/errors"
)

// publicProjectFields are the fields needed by the UI
// on base_angular and the menu
type UIProjectFields struct {
	Id          string `json:"id"`
	Identifier  string `json:"identifier"`
	DisplayName string `json:"display_name"`
	Repo        string `json:"repo_name"`
	Owner       string `json:"owner_name"`
}

type APITriggerDefinition struct {
	Project           *string `json:"project"`
	Level             *string `json:"level"` //build or task
	DefinitionID      *string `json:"definition_id"`
	BuildVariantRegex *string `json:"variant_regex"`
	TaskRegex         *string `json:"task_regex"`
	Status            *string `json:"status"`
	DateCutoff        *int    `json:"date_cutoff"`
	ConfigFile        *string `json:"config_file"`
	Alias             *string `json:"alias"`
}

func (t *APITriggerDefinition) ToService() model.TriggerDefinition {
	return model.TriggerDefinition{
		Project:           utility.FromStringPtr(t.Project),
		Level:             utility.FromStringPtr(t.Level),
		DefinitionID:      utility.FromStringPtr(t.DefinitionID),
		BuildVariantRegex: utility.FromStringPtr(t.BuildVariantRegex),
		TaskRegex:         utility.FromStringPtr(t.TaskRegex),
		Status:            utility.FromStringPtr(t.Status),
		ConfigFile:        utility.FromStringPtr(t.ConfigFile),
		Alias:             utility.FromStringPtr(t.Alias),
		DateCutoff:        t.DateCutoff,
	}
}

func (t *APITriggerDefinition) BuildFromService(triggerDef model.TriggerDefinition) {
	t.Project = utility.ToStringPtr(triggerDef.Project)
	t.Level = utility.ToStringPtr(triggerDef.Level)
	t.DefinitionID = utility.ToStringPtr(triggerDef.DefinitionID)
	t.BuildVariantRegex = utility.ToStringPtr(triggerDef.BuildVariantRegex)
	t.TaskRegex = utility.ToStringPtr(triggerDef.TaskRegex)
	t.Status = utility.ToStringPtr(triggerDef.Status)
	t.ConfigFile = utility.ToStringPtr(triggerDef.ConfigFile)
	t.Alias = utility.ToStringPtr(triggerDef.Alias)
	t.DateCutoff = triggerDef.DateCutoff
}

type APIPatchTriggerDefinition struct {
	Alias                  *string            `json:"alias"`
	ChildProjectId         *string            `json:"child_project_id"`
	ChildProjectIdentifier *string            `json:"child_project_identifier"`
	TaskSpecifiers         []APITaskSpecifier `json:"task_specifiers"`
	Status                 *string            `json:"status,omitempty"`
	ParentAsModule         *string            `json:"parent_as_module,omitempty"`
	VariantsTasks          []VariantTask      `json:"variants_tasks,omitempty"`
}

func (t *APIPatchTriggerDefinition) BuildFromService(def patch.PatchTriggerDefinition) error {
	t.ChildProjectId = utility.ToStringPtr(def.ChildProject) // we store the real ID in the child project field
	identifier, err := model.GetIdentifierForProject(def.ChildProject)
	if err != nil {
		return errors.Wrapf(err, "getting identifier for child project '%s'", def.ChildProject)
	}
	t.ChildProjectIdentifier = utility.ToStringPtr(identifier)
	// not sure in which direction this should go
	t.Alias = utility.ToStringPtr(def.Alias)
	t.Status = utility.ToStringPtr(def.Status)
	t.ParentAsModule = utility.ToStringPtr(def.ParentAsModule)
	var specifiers []APITaskSpecifier
	for _, ts := range def.TaskSpecifiers {
		specifier := APITaskSpecifier{}
		specifier.BuildFromService(ts)
		specifiers = append(specifiers, specifier)
	}
	t.TaskSpecifiers = specifiers
	return nil
}

func (t *APIPatchTriggerDefinition) ToService() patch.PatchTriggerDefinition {
	trigger := patch.PatchTriggerDefinition{}

	trigger.ChildProject = utility.FromStringPtr(t.ChildProjectIdentifier) // we'll fix this to be the ID in case it's changed
	trigger.Status = utility.FromStringPtr(t.Status)
	trigger.Alias = utility.FromStringPtr(t.Alias)
	trigger.ParentAsModule = utility.FromStringPtr(t.ParentAsModule)
	var specifiers []patch.TaskSpecifier
	for _, ts := range t.TaskSpecifiers {
		specifiers = append(specifiers, ts.ToService())
	}
	trigger.TaskSpecifiers = specifiers
	return trigger
}

type APITaskSpecifier struct {
	PatchAlias   *string `json:"patch_alias,omitempty"`
	TaskRegex    *string `json:"task_regex,omitempty"`
	VariantRegex *string `json:"variant_regex,omitempty"`
}

func (ts *APITaskSpecifier) BuildFromService(def patch.TaskSpecifier) {
	ts.PatchAlias = utility.ToStringPtr(def.PatchAlias)
	ts.TaskRegex = utility.ToStringPtr(def.TaskRegex)
	ts.VariantRegex = utility.ToStringPtr(def.VariantRegex)
}

func (t *APITaskSpecifier) ToService() patch.TaskSpecifier {
	specifier := patch.TaskSpecifier{}

	specifier.PatchAlias = utility.FromStringPtr(t.PatchAlias)
	specifier.TaskRegex = utility.FromStringPtr(t.TaskRegex)
	specifier.VariantRegex = utility.FromStringPtr(t.VariantRegex)
	return specifier
}

type APIPeriodicBuildDefinition struct {
	ID            *string    `json:"id"`
	ConfigFile    *string    `json:"config_file"`
	IntervalHours *int       `json:"interval_hours"`
	Alias         *string    `json:"alias,omitempty"`
	Message       *string    `json:"message,omitempty"`
	NextRunTime   *time.Time `json:"next_run_time,omitempty"`
}

type APICommitQueueParams struct {
	Enabled       *bool   `json:"enabled"`
	RequireSigned *bool   `json:"require_signed"`
	MergeMethod   *string `json:"merge_method"`
	Message       *string `json:"message"`
}

func (bd *APIPeriodicBuildDefinition) ToService() model.PeriodicBuildDefinition {
	buildDef := model.PeriodicBuildDefinition{}
	buildDef.ID = utility.FromStringPtr(bd.ID)
	buildDef.ConfigFile = utility.FromStringPtr(bd.ConfigFile)
	buildDef.IntervalHours = utility.FromIntPtr(bd.IntervalHours)
	buildDef.Alias = utility.FromStringPtr(bd.Alias)
	buildDef.Message = utility.FromStringPtr(bd.Message)
	buildDef.NextRunTime = utility.FromTimePtr(bd.NextRunTime)
	return buildDef
}

func (bd *APIPeriodicBuildDefinition) BuildFromService(params model.PeriodicBuildDefinition) {
	bd.ID = utility.ToStringPtr(params.ID)
	bd.ConfigFile = utility.ToStringPtr(params.ConfigFile)
	bd.IntervalHours = utility.ToIntPtr(params.IntervalHours)
	bd.Alias = utility.ToStringPtr(params.Alias)
	bd.Message = utility.ToStringPtr(params.Message)
	bd.NextRunTime = utility.ToTimePtr(params.NextRunTime)
}

func (cqParams *APICommitQueueParams) BuildFromService(params model.CommitQueueParams) {
	cqParams.Enabled = utility.BoolPtrCopy(params.Enabled)
	cqParams.RequireSigned = utility.BoolPtrCopy(params.RequireSigned)
	cqParams.MergeMethod = utility.ToStringPtr(params.MergeMethod)
	cqParams.Message = utility.ToStringPtr(params.Message)
}

func (cqParams *APICommitQueueParams) ToService() model.CommitQueueParams {
	serviceParams := model.CommitQueueParams{}
	serviceParams.Enabled = utility.BoolPtrCopy(cqParams.Enabled)
	serviceParams.RequireSigned = utility.BoolPtrCopy(cqParams.RequireSigned)
	serviceParams.MergeMethod = utility.FromStringPtr(cqParams.MergeMethod)
	serviceParams.Message = utility.FromStringPtr(cqParams.Message)

	return serviceParams
}

type APIBuildBaronSettings struct {
	TicketCreateProject     *string   `bson:"ticket_create_project" json:"ticket_create_project"`
	TicketSearchProjects    []*string `bson:"ticket_search_projects" json:"ticket_search_projects"`
	BFSuggestionServer      *string   `bson:"bf_suggestion_server" json:"bf_suggestion_server"`
	BFSuggestionUsername    *string   `bson:"bf_suggestion_username" json:"bf_suggestion_username"`
	BFSuggestionPassword    *string   `bson:"bf_suggestion_password" json:"bf_suggestion_password"`
	BFSuggestionTimeoutSecs *int      `bson:"bf_suggestion_timeout_secs" json:"bf_suggestion_timeout_secs"`
	BFSuggestionFeaturesURL *string   `bson:"bf_suggestion_features_url" json:"bf_suggestion_features_url"`
}

func (bb *APIBuildBaronSettings) BuildFromService(def evergreen.BuildBaronSettings) {
	bb.TicketCreateProject = utility.ToStringPtr(def.TicketCreateProject)
	bb.TicketSearchProjects = utility.ToStringPtrSlice(def.TicketSearchProjects)
	bb.BFSuggestionServer = utility.ToStringPtr(def.BFSuggestionServer)
	bb.BFSuggestionUsername = utility.ToStringPtr(def.BFSuggestionUsername)
	bb.BFSuggestionPassword = utility.ToStringPtr(def.BFSuggestionPassword)
	bb.BFSuggestionTimeoutSecs = utility.ToIntPtr(def.BFSuggestionTimeoutSecs)
	bb.BFSuggestionFeaturesURL = utility.ToStringPtr(def.BFSuggestionFeaturesURL)
}

func (bb *APIBuildBaronSettings) ToService() evergreen.BuildBaronSettings {
	buildBaron := evergreen.BuildBaronSettings{}
	buildBaron.TicketCreateProject = utility.FromStringPtr(bb.TicketCreateProject)
	buildBaron.TicketSearchProjects = utility.FromStringPtrSlice(bb.TicketSearchProjects)
	buildBaron.BFSuggestionServer = utility.FromStringPtr(bb.BFSuggestionServer)
	buildBaron.BFSuggestionUsername = utility.FromStringPtr(bb.BFSuggestionUsername)
	buildBaron.BFSuggestionPassword = utility.FromStringPtr(bb.BFSuggestionPassword)
	buildBaron.BFSuggestionTimeoutSecs = utility.FromIntPtr(bb.BFSuggestionTimeoutSecs)
	buildBaron.BFSuggestionFeaturesURL = utility.FromStringPtr(bb.BFSuggestionFeaturesURL)
	return buildBaron
}

type APITaskAnnotationSettings struct {
	JiraCustomFields  []APIJiraField `bson:"jira_custom_fields" json:"jira_custom_fields"`
	FileTicketWebhook APIWebHook     `bson:"web_hook" json:"web_hook"`
}

type APIWebHook struct {
	Endpoint *string `bson:"endpoint" json:"endpoint"`
	Secret   *string `bson:"secret" json:"secret"`
}

type APIJiraField struct {
	Field       *string `bson:"field" json:"field"`
	DisplayText *string `bson:"display_text" json:"display_text"`
}

func (ta *APITaskAnnotationSettings) ToService() evergreen.AnnotationsSettings {
	res := evergreen.AnnotationsSettings{}
	webhook := evergreen.WebHook{}
	webhook.Secret = utility.FromStringPtr(ta.FileTicketWebhook.Secret)
	webhook.Endpoint = utility.FromStringPtr(ta.FileTicketWebhook.Endpoint)
	res.FileTicketWebhook = webhook
	for _, apiJiraField := range ta.JiraCustomFields {
		jiraField := evergreen.JiraField{}
		jiraField.Field = utility.FromStringPtr(apiJiraField.Field)
		jiraField.DisplayText = utility.FromStringPtr(apiJiraField.DisplayText)
		res.JiraCustomFields = append(res.JiraCustomFields, jiraField)
	}
	return res
}

func (ta *APITaskAnnotationSettings) BuildFromService(config evergreen.AnnotationsSettings) {
	apiWebhook := APIWebHook{}
	apiWebhook.Secret = utility.ToStringPtr(config.FileTicketWebhook.Secret)
	apiWebhook.Endpoint = utility.ToStringPtr(config.FileTicketWebhook.Endpoint)
	ta.FileTicketWebhook = apiWebhook
	for _, jiraField := range config.JiraCustomFields {
		apiJiraField := APIJiraField{}
		apiJiraField.Field = utility.ToStringPtr(jiraField.Field)
		apiJiraField.DisplayText = utility.ToStringPtr(jiraField.DisplayText)
		ta.JiraCustomFields = append(ta.JiraCustomFields, apiJiraField)
	}
}

type APITaskSyncOptions struct {
	ConfigEnabled *bool `json:"config_enabled"`
	PatchEnabled  *bool `json:"patch_enabled"`
}

func (opts *APITaskSyncOptions) BuildFromService(in model.TaskSyncOptions) {
	opts.ConfigEnabled = utility.BoolPtrCopy(in.ConfigEnabled)
	opts.PatchEnabled = utility.BoolPtrCopy(in.PatchEnabled)
}

func (opts *APITaskSyncOptions) ToService() model.TaskSyncOptions {
	return model.TaskSyncOptions{
		ConfigEnabled: utility.BoolPtrCopy(opts.ConfigEnabled),
		PatchEnabled:  utility.BoolPtrCopy(opts.PatchEnabled),
	}
}

type APIWorkstationConfig struct {
	SetupCommands []APIWorkstationSetupCommand `bson:"setup_commands" json:"setup_commands"`
	GitClone      *bool                        `bson:"git_clone" json:"git_clone"`
}

type APIContainerCredential struct {
	Username *string `bson:"username" json:"username"`
	Password *string `bson:"password" json:"password"`
}

func (cr *APIContainerCredential) BuildFromService(h model.ContainerCredential) {
	cr.Username = utility.ToStringPtr(h.Username)
	cr.Password = utility.ToStringPtr(h.Password)
}

func (cr *APIContainerCredential) ToService() model.ContainerCredential {
	return model.ContainerCredential{
		Username: utility.FromStringPtr(cr.Username),
		Password: utility.FromStringPtr(cr.Password),
	}
}

type APIContainerResources struct {
	MemoryMB *int `bson:"memory_mb" json:"memory_mb"`
	CPU      *int `bson:"cpu" json:"cpu"`
}

func (cr *APIContainerResources) BuildFromService(h model.ContainerResources) {
	cr.MemoryMB = utility.ToIntPtr(h.MemoryMB)
	cr.CPU = utility.ToIntPtr(h.CPU)
}

func (cr *APIContainerResources) ToService() model.ContainerResources {
	return model.ContainerResources{
		MemoryMB: utility.FromIntPtr(cr.MemoryMB),
		CPU:      utility.FromIntPtr(cr.CPU),
	}
}

type APIWorkstationSetupCommand struct {
	Command   *string `bson:"command" json:"command"`
	Directory *string `bson:"directory" json:"directory"`
}

func (c *APIWorkstationConfig) ToService() model.WorkstationConfig {
	res := model.WorkstationConfig{}
	res.GitClone = utility.BoolPtrCopy(c.GitClone)
	if c.SetupCommands != nil {
		res.SetupCommands = []model.WorkstationSetupCommand{}
	}
	for _, apiCmd := range c.SetupCommands {
		cmd := model.WorkstationSetupCommand{}
		cmd.Command = utility.FromStringPtr(apiCmd.Command)
		cmd.Directory = utility.FromStringPtr(apiCmd.Directory)
		res.SetupCommands = append(res.SetupCommands, cmd)
	}
	return res
}

func (c *APIWorkstationConfig) BuildFromService(config model.WorkstationConfig) {
	c.GitClone = utility.BoolPtrCopy(config.GitClone)
	if config.SetupCommands != nil {
		c.SetupCommands = []APIWorkstationSetupCommand{}
	}
	for _, cmd := range config.SetupCommands {
		apiCmd := APIWorkstationSetupCommand{}
		apiCmd.Command = utility.ToStringPtr(cmd.Command)
		apiCmd.Directory = utility.ToStringPtr(cmd.Directory)
		c.SetupCommands = append(c.SetupCommands, apiCmd)
	}
}

type APIParameterInfo struct {
	Key         *string `json:"key"`
	Value       *string `json:"value"`
	Description *string `json:"description"`
}

func (c *APIParameterInfo) ToService() model.ParameterInfo {
	res := model.ParameterInfo{}
	res.Key = utility.FromStringPtr(c.Key)
	res.Value = utility.FromStringPtr(c.Value)
	res.Description = utility.FromStringPtr(c.Description)
	return res
}

func (c *APIParameterInfo) BuildFromService(info model.ParameterInfo) {
	c.Key = utility.ToStringPtr(info.Key)
	c.Value = utility.ToStringPtr(info.Value)
	c.Description = utility.ToStringPtr(info.Description)
}

type APIProjectRef struct {
	Id                          *string                   `json:"id"`
	Owner                       *string                   `json:"owner_name"`
	Repo                        *string                   `json:"repo_name"`
	Branch                      *string                   `json:"branch_name"`
	Enabled                     *bool                     `json:"enabled"`
	Private                     *bool                     `json:"private"`
	BatchTime                   int                       `json:"batch_time"`
	RemotePath                  *string                   `json:"remote_path"`
	SpawnHostScriptPath         *string                   `json:"spawn_host_script_path"`
	Identifier                  *string                   `json:"identifier"`
	DisplayName                 *string                   `json:"display_name"`
	DeactivatePrevious          *bool                     `json:"deactivate_previous"`
	TracksPushEvents            *bool                     `json:"tracks_push_events"`
	PRTestingEnabled            *bool                     `json:"pr_testing_enabled"`
	ManualPRTestingEnabled      *bool                     `json:"manual_pr_testing_enabled"`
	GitTagVersionsEnabled       *bool                     `json:"git_tag_versions_enabled"`
	GithubChecksEnabled         *bool                     `json:"github_checks_enabled"`
	CedarTestResultsEnabled     *bool                     `json:"cedar_test_results_enabled"`
	UseRepoSettings             *bool                     `json:"use_repo_settings"`
	RepoRefId                   *string                   `json:"repo_ref_id"`
	DefaultLogger               *string                   `json:"default_logger"`
	CommitQueue                 APICommitQueueParams      `json:"commit_queue"`
	TaskSync                    APITaskSyncOptions        `json:"task_sync"`
	TaskAnnotationSettings      APITaskAnnotationSettings `json:"task_annotation_settings"`
	BuildBaronSettings          APIBuildBaronSettings     `json:"build_baron_settings"`
	PerfEnabled                 *bool                     `json:"perf_enabled"`
	Hidden                      *bool                     `json:"hidden"`
	PatchingDisabled            *bool                     `json:"patching_disabled"`
	RepotrackerDisabled         *bool                     `json:"repotracker_disabled"`
	DispatchingDisabled         *bool                     `json:"dispatching_disabled"`
	VersionControlEnabled       *bool                     `json:"version_control_enabled"`
	DisabledStatsCache          *bool                     `json:"disabled_stats_cache"`
	FilesIgnoredFromCache       []*string                 `json:"files_ignored_from_cache"`
	Admins                      []*string                 `json:"admins"`
	DeleteAdmins                []*string                 `json:"delete_admins,omitempty"`
	GitTagAuthorizedUsers       []*string                 `json:"git_tag_authorized_users" bson:"git_tag_authorized_users"`
	DeleteGitTagAuthorizedUsers []*string                 `json:"delete_git_tag_authorized_users,omitempty" bson:"delete_git_tag_authorized_users,omitempty"`
	GitTagAuthorizedTeams       []*string                 `json:"git_tag_authorized_teams" bson:"git_tag_authorized_teams"`
	DeleteGitTagAuthorizedTeams []*string                 `json:"delete_git_tag_authorized_teams,omitempty" bson:"delete_git_tag_authorized_teams,omitempty"`
	NotifyOnBuildFailure        *bool                     `json:"notify_on_failure"`
	Restricted                  *bool                     `json:"restricted"`
	Revision                    *string                   `json:"revision"`

	Triggers             []APITriggerDefinition       `json:"triggers"`
	GithubTriggerAliases []*string                    `json:"github_trigger_aliases"`
	PatchTriggerAliases  []APIPatchTriggerDefinition  `json:"patch_trigger_aliases"`
	Aliases              []APIProjectAlias            `json:"aliases"`
	Variables            APIProjectVars               `json:"variables"`
	WorkstationConfig    APIWorkstationConfig         `json:"workstation_config"`
	Subscriptions        []APISubscription            `json:"subscriptions"`
	DeleteSubscriptions  []*string                    `json:"delete_subscriptions,omitempty"`
	PeriodicBuilds       []APIPeriodicBuildDefinition `json:"periodic_builds,omitempty"`
}

// ToService returns a service layer ProjectRef using the data from APIProjectRef
func (p *APIProjectRef) ToService() model.ProjectRef {
	projectRef := model.ProjectRef{
		Owner:                   utility.FromStringPtr(p.Owner),
		Repo:                    utility.FromStringPtr(p.Repo),
		Branch:                  utility.FromStringPtr(p.Branch),
		Enabled:                 utility.BoolPtrCopy(p.Enabled),
		Private:                 utility.BoolPtrCopy(p.Private),
		Restricted:              utility.BoolPtrCopy(p.Restricted),
		BatchTime:               p.BatchTime,
		RemotePath:              utility.FromStringPtr(p.RemotePath),
		Id:                      utility.FromStringPtr(p.Id),
		Identifier:              utility.FromStringPtr(p.Identifier),
		DisplayName:             utility.FromStringPtr(p.DisplayName),
		DeactivatePrevious:      utility.BoolPtrCopy(p.DeactivatePrevious),
		TracksPushEvents:        utility.BoolPtrCopy(p.TracksPushEvents),
		DefaultLogger:           utility.FromStringPtr(p.DefaultLogger),
		PRTestingEnabled:        utility.BoolPtrCopy(p.PRTestingEnabled),
		ManualPRTestingEnabled:  utility.BoolPtrCopy(p.ManualPRTestingEnabled),
		GitTagVersionsEnabled:   utility.BoolPtrCopy(p.GitTagVersionsEnabled),
		GithubChecksEnabled:     utility.BoolPtrCopy(p.GithubChecksEnabled),
		CedarTestResultsEnabled: utility.BoolPtrCopy(p.CedarTestResultsEnabled),
		RepoRefId:               utility.FromStringPtr(p.RepoRefId),
		CommitQueue:             p.CommitQueue.ToService(),
		TaskSync:                p.TaskSync.ToService(),
		WorkstationConfig:       p.WorkstationConfig.ToService(),
		BuildBaronSettings:      p.BuildBaronSettings.ToService(),
		TaskAnnotationSettings:  p.TaskAnnotationSettings.ToService(),
		PerfEnabled:             utility.BoolPtrCopy(p.PerfEnabled),
		Hidden:                  utility.BoolPtrCopy(p.Hidden),
		PatchingDisabled:        utility.BoolPtrCopy(p.PatchingDisabled),
		RepotrackerDisabled:     utility.BoolPtrCopy(p.RepotrackerDisabled),
		DispatchingDisabled:     utility.BoolPtrCopy(p.DispatchingDisabled),
		VersionControlEnabled:   utility.BoolPtrCopy(p.VersionControlEnabled),
		DisabledStatsCache:      utility.BoolPtrCopy(p.DisabledStatsCache),
		FilesIgnoredFromCache:   utility.FromStringPtrSlice(p.FilesIgnoredFromCache),
		NotifyOnBuildFailure:    utility.BoolPtrCopy(p.NotifyOnBuildFailure),
		SpawnHostScriptPath:     utility.FromStringPtr(p.SpawnHostScriptPath),
		Admins:                  utility.FromStringPtrSlice(p.Admins),
		GitTagAuthorizedUsers:   utility.FromStringPtrSlice(p.GitTagAuthorizedUsers),
		GitTagAuthorizedTeams:   utility.FromStringPtrSlice(p.GitTagAuthorizedTeams),
		GithubTriggerAliases:    utility.FromStringPtrSlice(p.GithubTriggerAliases),
	}

	// Copy triggers
	if p.Triggers != nil {
		triggers := []model.TriggerDefinition{}
		for _, t := range p.Triggers {
			triggers = append(triggers, t.ToService())
		}
		projectRef.Triggers = triggers
	}

	// Copy periodic builds
	if p.PeriodicBuilds != nil {
		builds := []model.PeriodicBuildDefinition{}
		for _, b := range p.PeriodicBuilds {
			builds = append(builds, b.ToService())
		}
		projectRef.PeriodicBuilds = builds
	}

	if p.PatchTriggerAliases != nil {
		patchTriggers := []patch.PatchTriggerDefinition{}
		for _, a := range p.PatchTriggerAliases {
			patchTriggers = append(patchTriggers, a.ToService())
		}
		projectRef.PatchTriggerAliases = patchTriggers
	}
	return projectRef
}

func (p *APIProjectRef) BuildFromService(projectRef model.ProjectRef) error {
	p.Owner = utility.ToStringPtr(projectRef.Owner)
	p.Repo = utility.ToStringPtr(projectRef.Repo)
	p.Branch = utility.ToStringPtr(projectRef.Branch)
	p.Enabled = utility.BoolPtrCopy(projectRef.Enabled)
	p.Private = utility.BoolPtrCopy(projectRef.Private)
	p.Restricted = utility.BoolPtrCopy(projectRef.Restricted)
	p.BatchTime = projectRef.BatchTime
	p.RemotePath = utility.ToStringPtr(projectRef.RemotePath)
	p.Id = utility.ToStringPtr(projectRef.Id)
	p.Identifier = utility.ToStringPtr(projectRef.Identifier)
	p.DisplayName = utility.ToStringPtr(projectRef.DisplayName)
	p.DeactivatePrevious = projectRef.DeactivatePrevious
	p.TracksPushEvents = utility.BoolPtrCopy(projectRef.TracksPushEvents)
	p.DefaultLogger = utility.ToStringPtr(projectRef.DefaultLogger)
	p.PRTestingEnabled = utility.BoolPtrCopy(projectRef.PRTestingEnabled)
	p.ManualPRTestingEnabled = utility.BoolPtrCopy(projectRef.ManualPRTestingEnabled)
	p.GitTagVersionsEnabled = utility.BoolPtrCopy(projectRef.GitTagVersionsEnabled)
	p.GithubChecksEnabled = utility.BoolPtrCopy(projectRef.GithubChecksEnabled)
	p.CedarTestResultsEnabled = utility.BoolPtrCopy(projectRef.CedarTestResultsEnabled)
	p.UseRepoSettings = utility.ToBoolPtr(projectRef.UseRepoSettings())
	p.RepoRefId = utility.ToStringPtr(projectRef.RepoRefId)
	p.PerfEnabled = utility.BoolPtrCopy(projectRef.PerfEnabled)
	p.Hidden = utility.BoolPtrCopy(projectRef.Hidden)
	p.PatchingDisabled = utility.BoolPtrCopy(projectRef.PatchingDisabled)
	p.RepotrackerDisabled = utility.BoolPtrCopy(projectRef.RepotrackerDisabled)
	p.DispatchingDisabled = utility.BoolPtrCopy(projectRef.DispatchingDisabled)
	p.VersionControlEnabled = utility.BoolPtrCopy(projectRef.VersionControlEnabled)
	p.DisabledStatsCache = utility.BoolPtrCopy(projectRef.DisabledStatsCache)
	p.FilesIgnoredFromCache = utility.ToStringPtrSlice(projectRef.FilesIgnoredFromCache)
	p.NotifyOnBuildFailure = utility.BoolPtrCopy(projectRef.NotifyOnBuildFailure)
	p.SpawnHostScriptPath = utility.ToStringPtr(projectRef.SpawnHostScriptPath)
	p.Admins = utility.ToStringPtrSlice(projectRef.Admins)
	p.GitTagAuthorizedUsers = utility.ToStringPtrSlice(projectRef.GitTagAuthorizedUsers)
	p.GitTagAuthorizedTeams = utility.ToStringPtrSlice(projectRef.GitTagAuthorizedTeams)
	p.GithubTriggerAliases = utility.ToStringPtrSlice(projectRef.GithubTriggerAliases)

	cq := APICommitQueueParams{}
	cq.BuildFromService(projectRef.CommitQueue)
	p.CommitQueue = cq

	var taskSync APITaskSyncOptions
	taskSync.BuildFromService(projectRef.TaskSync)
	p.TaskSync = taskSync

	workstationConfig := APIWorkstationConfig{}
	workstationConfig.BuildFromService(projectRef.WorkstationConfig)
	p.WorkstationConfig = workstationConfig

	buildbaronConfig := APIBuildBaronSettings{}
	buildbaronConfig.BuildFromService(projectRef.BuildBaronSettings)
	p.BuildBaronSettings = buildbaronConfig

	taskAnnotationConfig := APITaskAnnotationSettings{}
	taskAnnotationConfig.BuildFromService(projectRef.TaskAnnotationSettings)
	p.TaskAnnotationSettings = taskAnnotationConfig

	// Copy triggers
	if projectRef.Triggers != nil {
		triggers := []APITriggerDefinition{}
		for _, t := range projectRef.Triggers {
			apiTrigger := APITriggerDefinition{}
			apiTrigger.BuildFromService(t)
			triggers = append(triggers, apiTrigger)
		}
		p.Triggers = triggers
	}

	// copy periodic builds
	if projectRef.PeriodicBuilds != nil {
		periodicBuilds := []APIPeriodicBuildDefinition{}
		for _, pb := range projectRef.PeriodicBuilds {
			periodicBuild := APIPeriodicBuildDefinition{}
			periodicBuild.BuildFromService(pb)
			periodicBuilds = append(periodicBuilds, periodicBuild)
		}
		p.PeriodicBuilds = periodicBuilds
	}

	if projectRef.PatchTriggerAliases != nil {
		patchTriggers := []APIPatchTriggerDefinition{}
		for idx, a := range projectRef.PatchTriggerAliases {
			trigger := APIPatchTriggerDefinition{}
			if err := trigger.BuildFromService(a); err != nil {
				return errors.Wrapf(err, "converting patch trigger alias at index %d to service model", idx)
			}
			patchTriggers = append(patchTriggers, trigger)
		}
		p.PatchTriggerAliases = patchTriggers
	}
	return nil
}

// DefaultUnsetBooleans is used to set booleans to their default value.
func (pRef *APIProjectRef) DefaultUnsetBooleans() {
	reflected := reflect.ValueOf(pRef).Elem()
	recursivelyDefaultBooleans(reflected)
}

func recursivelyDefaultBooleans(structToSet reflect.Value) {
	var err error
	var i int
	defer func() {
		grip.Error(recovery.HandlePanicWithError(recover(), err, fmt.Sprintf("panicked while recursively defaulting booleans for field number %d", i)))
	}()
	falseType := reflect.TypeOf(false)
	// Iterate through each field of the struct.
	for i = 0; i < structToSet.NumField(); i++ {
		field := structToSet.Field(i)

		// If it's a boolean pointer, set the default recursively.
		if field.Type() == reflect.PtrTo(falseType) && util.IsFieldUndefined(field) {
			field.Set(reflect.New(falseType))

		} else if field.Kind() == reflect.Struct {
			recursivelyDefaultBooleans(field)
		}
	}
}
