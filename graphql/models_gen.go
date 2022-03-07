// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package graphql

import (
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/evergreen-ci/evergreen/apimodels"
	"github.com/evergreen-ci/evergreen/model/host"
	"github.com/evergreen-ci/evergreen/rest/model"
	"github.com/evergreen-ci/evergreen/thirdparty"
)

type AbortInfo struct {
	User                    string `json:"user"`
	TaskID                  string `json:"taskID"`
	TaskDisplayName         string `json:"taskDisplayName"`
	BuildVariantDisplayName string `json:"buildVariantDisplayName"`
	NewVersion              string `json:"newVersion"`
	PrClosed                bool   `json:"prClosed"`
}

type BaseTaskMetadata struct {
	BaseTaskDuration *model.APIDuration `json:"baseTaskDuration"`
	BaseTaskLink     string             `json:"baseTaskLink"`
}

type BaseTaskResult struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

type BuildBaron struct {
	SearchReturnInfo        *thirdparty.SearchReturnInfo `json:"searchReturnInfo"`
	BuildBaronConfigured    bool                         `json:"buildBaronConfigured"`
	BbTicketCreationDefined bool                         `json:"bbTicketCreationDefined"`
}

type BuildVariantOptions struct {
	Variants []string `json:"variants"`
	Tasks    []string `json:"tasks"`
	Statuses []string `json:"statuses"`
}

type Dependency struct {
	Name           string         `json:"name"`
	MetStatus      MetStatus      `json:"metStatus"`
	RequiredStatus RequiredStatus `json:"requiredStatus"`
	BuildVariant   string         `json:"buildVariant"`
	TaskID         string         `json:"taskId"`
	UILink         string         `json:"uiLink"`
}

type DisplayTask struct {
	Name      string   `json:"Name"`
	ExecTasks []string `json:"ExecTasks"`
}

type EditSpawnHostInput struct {
	HostID              string          `json:"hostId"`
	DisplayName         *string         `json:"displayName"`
	Expiration          *time.Time      `json:"expiration"`
	NoExpiration        *bool           `json:"noExpiration"`
	InstanceType        *string         `json:"instanceType"`
	AddedInstanceTags   []*host.Tag     `json:"addedInstanceTags"`
	DeletedInstanceTags []*host.Tag     `json:"deletedInstanceTags"`
	Volume              *string         `json:"volume"`
	ServicePassword     *string         `json:"servicePassword"`
	PublicKey           *PublicKeyInput `json:"publicKey"`
	SavePublicKey       *bool           `json:"savePublicKey"`
}

type GroupedBuildVariant struct {
	Variant     string           `json:"variant"`
	DisplayName string           `json:"displayName"`
	Tasks       []*model.APITask `json:"tasks"`
}

type GroupedFiles struct {
	TaskName *string          `json:"taskName"`
	Files    []*model.APIFile `json:"files"`
}

type GroupedProjects struct {
	GroupDisplayName string                 `json:"groupDisplayName"`
	Name             string                 `json:"name"`
	Repo             *model.APIProjectRef   `json:"repo"`
	Projects         []*model.APIProjectRef `json:"projects"`
}

type HostEvents struct {
	EventLogEntries []*model.HostAPIEventLogEntry `json:"eventLogEntries"`
	Count           int                           `json:"count"`
}

type HostsResponse struct {
	FilteredHostsCount *int             `json:"filteredHostsCount"`
	TotalHostsCount    int              `json:"totalHostsCount"`
	Hosts              []*model.APIHost `json:"hosts"`
}

type MainlineCommitVersion struct {
	Version          *model.APIVersion   `json:"version"`
	RolledUpVersions []*model.APIVersion `json:"rolledUpVersions"`
}

type MainlineCommits struct {
	NextPageOrderNumber *int                     `json:"nextPageOrderNumber"`
	PrevPageOrderNumber *int                     `json:"prevPageOrderNumber"`
	Versions            []*MainlineCommitVersion `json:"versions"`
}

type MainlineCommitsOptions struct {
	ProjectID       string   `json:"projectID"`
	Limit           *int     `json:"limit"`
	SkipOrderNumber *int     `json:"skipOrderNumber"`
	ShouldCollapse  *bool    `json:"shouldCollapse"`
	Requesters      []string `json:"requesters"`
}

type Manifest struct {
	ID              string                 `json:"id"`
	Revision        string                 `json:"revision"`
	Project         string                 `json:"project"`
	Branch          string                 `json:"branch"`
	IsBase          bool                   `json:"isBase"`
	ModuleOverrides map[string]string      `json:"moduleOverrides"`
	Modules         map[string]interface{} `json:"modules"`
}

type MoveProjectInput struct {
	ProjectID string `json:"projectId"`
	NewOwner  string `json:"newOwner"`
	NewRepo   string `json:"newRepo"`
}

type PatchConfigure struct {
	Description         string                `json:"description"`
	VariantsTasks       []*VariantTasks       `json:"variantsTasks"`
	Parameters          []*model.APIParameter `json:"parameters"`
	PatchTriggerAliases []string              `json:"patchTriggerAliases"`
}

type PatchDuration struct {
	Makespan  *string    `json:"makespan"`
	TimeTaken *string    `json:"timeTaken"`
	Time      *PatchTime `json:"time"`
}

type PatchMetadata struct {
	Author  string `json:"author"`
	PatchID string `json:"patchID"`
}

type PatchProject struct {
	Variants []*ProjectBuildVariant `json:"variants"`
}

type PatchTasks struct {
	Tasks []*model.APITask `json:"tasks"`
	Count int              `json:"count"`
}

type PatchTime struct {
	Started     *string `json:"started"`
	Finished    *string `json:"finished"`
	SubmittedAt string  `json:"submittedAt"`
}

type Patches struct {
	Patches            []*model.APIPatch `json:"patches"`
	FilteredPatchCount int               `json:"filteredPatchCount"`
}

type PatchesInput struct {
	Limit              int      `json:"limit"`
	Page               int      `json:"page"`
	PatchName          string   `json:"patchName"`
	Statuses           []string `json:"statuses"`
	IncludeCommitQueue *bool    `json:"includeCommitQueue"`
	OnlyCommitQueue    *bool    `json:"onlyCommitQueue"`
}

type ProjectBuildVariant struct {
	Name        string   `json:"name"`
	DisplayName string   `json:"displayName"`
	Tasks       []string `json:"tasks"`
}

type ProjectEvents struct {
	EventLogEntries []*model.APIProjectEvent `json:"eventLogEntries"`
	Count           int                      `json:"count"`
}

type ProjectPermissions struct {
	CreateProject bool `json:"CreateProject"`
}

type PublicKeyInput struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

type RepoEvents struct {
	EventLogEntries []*model.APIProjectEvent `json:"eventLogEntries"`
	Count           int                      `json:"count"`
}

type SortOrder struct {
	Key       TaskSortCategory `json:"Key"`
	Direction SortDirection    `json:"Direction"`
}

type SpawnHostInput struct {
	DistroID                string          `json:"distroId"`
	Region                  string          `json:"region"`
	SavePublicKey           bool            `json:"savePublicKey"`
	PublicKey               *PublicKeyInput `json:"publicKey"`
	UserDataScript          *string         `json:"userDataScript"`
	Expiration              *time.Time      `json:"expiration"`
	NoExpiration            bool            `json:"noExpiration"`
	SetUpScript             *string         `json:"setUpScript"`
	IsVirtualWorkStation    bool            `json:"isVirtualWorkStation"`
	HomeVolumeSize          *int            `json:"homeVolumeSize"`
	VolumeID                *string         `json:"volumeId"`
	TaskID                  *string         `json:"taskId"`
	UseProjectSetupScript   *bool           `json:"useProjectSetupScript"`
	UseTaskConfig           *bool           `json:"useTaskConfig"`
	SpawnHostsStartedByTask *bool           `json:"spawnHostsStartedByTask"`
	TaskSync                *bool           `json:"taskSync"`
}

type SpawnVolumeInput struct {
	AvailabilityZone string     `json:"availabilityZone"`
	Size             int        `json:"size"`
	Type             string     `json:"type"`
	Expiration       *time.Time `json:"expiration"`
	NoExpiration     *bool      `json:"noExpiration"`
	Host             *string    `json:"host"`
}

type Subscriber struct {
	GithubPRSubscriber    *model.APIGithubPRSubscriber    `json:"githubPRSubscriber"`
	GithubCheckSubscriber *model.APIGithubCheckSubscriber `json:"githubCheckSubscriber"`
	WebhookSubscriber     *model.APIWebhookSubscriber     `json:"webhookSubscriber"`
	JiraIssueSubscriber   *model.APIJIRAIssueSubscriber   `json:"jiraIssueSubscriber"`
	JiraCommentSubscriber *string                         `json:"jiraCommentSubscriber"`
	EmailSubscriber       *string                         `json:"emailSubscriber"`
	SlackSubscriber       *string                         `json:"slackSubscriber"`
}

type TaskFiles struct {
	FileCount    int             `json:"fileCount"`
	GroupedFiles []*GroupedFiles `json:"groupedFiles"`
}

type TaskLogs struct {
	TaskID        string                        `json:"taskId"`
	Execution     int                           `json:"execution"`
	DefaultLogger string                        `json:"defaultLogger"`
	EventLogs     []*model.TaskAPIEventLogEntry `json:"eventLogs"`
	TaskLogs      []*apimodels.LogMessage       `json:"taskLogs"`
	SystemLogs    []*apimodels.LogMessage       `json:"systemLogs"`
	AgentLogs     []*apimodels.LogMessage       `json:"agentLogs"`
	AllLogs       []*apimodels.LogMessage       `json:"allLogs"`
}

type TaskQueueDistro struct {
	ID         string `json:"id"`
	QueueCount int    `json:"queueCount"`
	TaskCount  int    `json:"taskCount"`
	HostCount  int    `json:"hostCount"`
}

type TaskResult struct {
	ID                      string           `json:"id"`
	Execution               int              `json:"execution"`
	Aborted                 bool             `json:"aborted"`
	DisplayName             string           `json:"displayName"`
	Version                 string           `json:"version"`
	Status                  string           `json:"status"`
	BaseStatus              *string          `json:"baseStatus"`
	BaseTask                *BaseTaskResult  `json:"baseTask"`
	BuildVariant            string           `json:"buildVariant"`
	BuildVariantDisplayName string           `json:"buildVariantDisplayName"`
	Blocked                 bool             `json:"blocked"`
	ExecutionTasksFull      []*model.APITask `json:"executionTasksFull"`
}

type TaskTestResult struct {
	TotalTestCount    int              `json:"totalTestCount"`
	FilteredTestCount int              `json:"filteredTestCount"`
	TestResults       []*model.APITest `json:"testResults"`
}

type TaskTestResultSample struct {
	TaskID                  string   `json:"taskId"`
	Execution               int      `json:"execution"`
	TotalTestCount          int      `json:"totalTestCount"`
	MatchingFailedTestNames []string `json:"matchingFailedTestNames"`
}

type TestFilter struct {
	TestName   string `json:"testName"`
	TestStatus string `json:"testStatus"`
}

type UpdateVolumeInput struct {
	Expiration   *time.Time `json:"expiration"`
	NoExpiration *bool      `json:"noExpiration"`
	Name         *string    `json:"name"`
	VolumeID     string     `json:"volumeId"`
}

type UserConfig struct {
	User          string `json:"user"`
	APIKey        string `json:"api_key"`
	APIServerHost string `json:"api_server_host"`
	UIServerHost  string `json:"ui_server_host"`
}

type UserPatches struct {
	Patches            []*model.APIPatch `json:"patches"`
	FilteredPatchCount int               `json:"filteredPatchCount"`
}

type VariantTasks struct {
	Variant      string         `json:"variant"`
	Tasks        []string       `json:"tasks"`
	DisplayTasks []*DisplayTask `json:"displayTasks"`
}

type VersionTiming struct {
	Makespan  *model.APIDuration `json:"makespan"`
	TimeTaken *model.APIDuration `json:"timeTaken"`
}

type VolumeHost struct {
	VolumeID string `json:"volumeId"`
	HostID   string `json:"hostId"`
}

type HostSortBy string

const (
	HostSortByID          HostSortBy = "ID"
	HostSortByDistro      HostSortBy = "DISTRO"
	HostSortByCurrentTask HostSortBy = "CURRENT_TASK"
	HostSortByStatus      HostSortBy = "STATUS"
	HostSortByElapsed     HostSortBy = "ELAPSED"
	HostSortByUptime      HostSortBy = "UPTIME"
	HostSortByIDLeTime    HostSortBy = "IDLE_TIME"
	HostSortByOwner       HostSortBy = "OWNER"
)

var AllHostSortBy = []HostSortBy{
	HostSortByID,
	HostSortByDistro,
	HostSortByCurrentTask,
	HostSortByStatus,
	HostSortByElapsed,
	HostSortByUptime,
	HostSortByIDLeTime,
	HostSortByOwner,
}

func (e HostSortBy) IsValid() bool {
	switch e {
	case HostSortByID, HostSortByDistro, HostSortByCurrentTask, HostSortByStatus, HostSortByElapsed, HostSortByUptime, HostSortByIDLeTime, HostSortByOwner:
		return true
	}
	return false
}

func (e HostSortBy) String() string {
	return string(e)
}

func (e *HostSortBy) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = HostSortBy(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid HostSortBy", str)
	}
	return nil
}

func (e HostSortBy) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type MetStatus string

const (
	MetStatusUnmet   MetStatus = "UNMET"
	MetStatusMet     MetStatus = "MET"
	MetStatusPending MetStatus = "PENDING"
	MetStatusStarted MetStatus = "STARTED"
)

var AllMetStatus = []MetStatus{
	MetStatusUnmet,
	MetStatusMet,
	MetStatusPending,
	MetStatusStarted,
}

func (e MetStatus) IsValid() bool {
	switch e {
	case MetStatusUnmet, MetStatusMet, MetStatusPending, MetStatusStarted:
		return true
	}
	return false
}

func (e MetStatus) String() string {
	return string(e)
}

func (e *MetStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = MetStatus(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid MetStatus", str)
	}
	return nil
}

func (e MetStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type ProjectSettingsAccess string

const (
	ProjectSettingsAccessEdit ProjectSettingsAccess = "EDIT"
	ProjectSettingsAccessView ProjectSettingsAccess = "VIEW"
)

var AllProjectSettingsAccess = []ProjectSettingsAccess{
	ProjectSettingsAccessEdit,
	ProjectSettingsAccessView,
}

func (e ProjectSettingsAccess) IsValid() bool {
	switch e {
	case ProjectSettingsAccessEdit, ProjectSettingsAccessView:
		return true
	}
	return false
}

func (e ProjectSettingsAccess) String() string {
	return string(e)
}

func (e *ProjectSettingsAccess) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ProjectSettingsAccess(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ProjectSettingsAccess", str)
	}
	return nil
}

func (e ProjectSettingsAccess) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type ProjectSettingsSection string

const (
	ProjectSettingsSectionGeneral              ProjectSettingsSection = "GENERAL"
	ProjectSettingsSectionAccess               ProjectSettingsSection = "ACCESS"
	ProjectSettingsSectionVariables            ProjectSettingsSection = "VARIABLES"
	ProjectSettingsSectionGithubAndCommitQueue ProjectSettingsSection = "GITHUB_AND_COMMIT_QUEUE"
	ProjectSettingsSectionNotifications        ProjectSettingsSection = "NOTIFICATIONS"
	ProjectSettingsSectionPatchAliases         ProjectSettingsSection = "PATCH_ALIASES"
	ProjectSettingsSectionWorkstation          ProjectSettingsSection = "WORKSTATION"
	ProjectSettingsSectionTriggers             ProjectSettingsSection = "TRIGGERS"
	ProjectSettingsSectionPeriodicBuilds       ProjectSettingsSection = "PERIODIC_BUILDS"
	ProjectSettingsSectionPlugins              ProjectSettingsSection = "PLUGINS"
)

var AllProjectSettingsSection = []ProjectSettingsSection{
	ProjectSettingsSectionGeneral,
	ProjectSettingsSectionAccess,
	ProjectSettingsSectionVariables,
	ProjectSettingsSectionGithubAndCommitQueue,
	ProjectSettingsSectionNotifications,
	ProjectSettingsSectionPatchAliases,
	ProjectSettingsSectionWorkstation,
	ProjectSettingsSectionTriggers,
	ProjectSettingsSectionPeriodicBuilds,
	ProjectSettingsSectionPlugins,
}

func (e ProjectSettingsSection) IsValid() bool {
	switch e {
	case ProjectSettingsSectionGeneral, ProjectSettingsSectionAccess, ProjectSettingsSectionVariables, ProjectSettingsSectionGithubAndCommitQueue, ProjectSettingsSectionNotifications, ProjectSettingsSectionPatchAliases, ProjectSettingsSectionWorkstation, ProjectSettingsSectionTriggers, ProjectSettingsSectionPeriodicBuilds, ProjectSettingsSectionPlugins:
		return true
	}
	return false
}

func (e ProjectSettingsSection) String() string {
	return string(e)
}

func (e *ProjectSettingsSection) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ProjectSettingsSection(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ProjectSettingsSection", str)
	}
	return nil
}

func (e ProjectSettingsSection) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type RequiredStatus string

const (
	RequiredStatusMustFail    RequiredStatus = "MUST_FAIL"
	RequiredStatusMustFinish  RequiredStatus = "MUST_FINISH"
	RequiredStatusMustSucceed RequiredStatus = "MUST_SUCCEED"
)

var AllRequiredStatus = []RequiredStatus{
	RequiredStatusMustFail,
	RequiredStatusMustFinish,
	RequiredStatusMustSucceed,
}

func (e RequiredStatus) IsValid() bool {
	switch e {
	case RequiredStatusMustFail, RequiredStatusMustFinish, RequiredStatusMustSucceed:
		return true
	}
	return false
}

func (e RequiredStatus) String() string {
	return string(e)
}

func (e *RequiredStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = RequiredStatus(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid RequiredStatus", str)
	}
	return nil
}

func (e RequiredStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type SortDirection string

const (
	SortDirectionAsc  SortDirection = "ASC"
	SortDirectionDesc SortDirection = "DESC"
)

var AllSortDirection = []SortDirection{
	SortDirectionAsc,
	SortDirectionDesc,
}

func (e SortDirection) IsValid() bool {
	switch e {
	case SortDirectionAsc, SortDirectionDesc:
		return true
	}
	return false
}

func (e SortDirection) String() string {
	return string(e)
}

func (e *SortDirection) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = SortDirection(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid SortDirection", str)
	}
	return nil
}

func (e SortDirection) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type SpawnHostStatusActions string

const (
	SpawnHostStatusActionsStart     SpawnHostStatusActions = "START"
	SpawnHostStatusActionsStop      SpawnHostStatusActions = "STOP"
	SpawnHostStatusActionsTerminate SpawnHostStatusActions = "TERMINATE"
)

var AllSpawnHostStatusActions = []SpawnHostStatusActions{
	SpawnHostStatusActionsStart,
	SpawnHostStatusActionsStop,
	SpawnHostStatusActionsTerminate,
}

func (e SpawnHostStatusActions) IsValid() bool {
	switch e {
	case SpawnHostStatusActionsStart, SpawnHostStatusActionsStop, SpawnHostStatusActionsTerminate:
		return true
	}
	return false
}

func (e SpawnHostStatusActions) String() string {
	return string(e)
}

func (e *SpawnHostStatusActions) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = SpawnHostStatusActions(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid SpawnHostStatusActions", str)
	}
	return nil
}

func (e SpawnHostStatusActions) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type TaskQueueItemType string

const (
	TaskQueueItemTypeCommit TaskQueueItemType = "COMMIT"
	TaskQueueItemTypePatch  TaskQueueItemType = "PATCH"
)

var AllTaskQueueItemType = []TaskQueueItemType{
	TaskQueueItemTypeCommit,
	TaskQueueItemTypePatch,
}

func (e TaskQueueItemType) IsValid() bool {
	switch e {
	case TaskQueueItemTypeCommit, TaskQueueItemTypePatch:
		return true
	}
	return false
}

func (e TaskQueueItemType) String() string {
	return string(e)
}

func (e *TaskQueueItemType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = TaskQueueItemType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid TaskQueueItemType", str)
	}
	return nil
}

func (e TaskQueueItemType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type TaskSortCategory string

const (
	TaskSortCategoryName       TaskSortCategory = "NAME"
	TaskSortCategoryStatus     TaskSortCategory = "STATUS"
	TaskSortCategoryBaseStatus TaskSortCategory = "BASE_STATUS"
	TaskSortCategoryVariant    TaskSortCategory = "VARIANT"
)

var AllTaskSortCategory = []TaskSortCategory{
	TaskSortCategoryName,
	TaskSortCategoryStatus,
	TaskSortCategoryBaseStatus,
	TaskSortCategoryVariant,
}

func (e TaskSortCategory) IsValid() bool {
	switch e {
	case TaskSortCategoryName, TaskSortCategoryStatus, TaskSortCategoryBaseStatus, TaskSortCategoryVariant:
		return true
	}
	return false
}

func (e TaskSortCategory) String() string {
	return string(e)
}

func (e *TaskSortCategory) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = TaskSortCategory(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid TaskSortCategory", str)
	}
	return nil
}

func (e TaskSortCategory) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type TestSortCategory string

const (
	TestSortCategoryBaseStatus TestSortCategory = "BASE_STATUS"
	TestSortCategoryStatus     TestSortCategory = "STATUS"
	TestSortCategoryStartTime  TestSortCategory = "START_TIME"
	TestSortCategoryDuration   TestSortCategory = "DURATION"
	TestSortCategoryTestName   TestSortCategory = "TEST_NAME"
)

var AllTestSortCategory = []TestSortCategory{
	TestSortCategoryBaseStatus,
	TestSortCategoryStatus,
	TestSortCategoryStartTime,
	TestSortCategoryDuration,
	TestSortCategoryTestName,
}

func (e TestSortCategory) IsValid() bool {
	switch e {
	case TestSortCategoryBaseStatus, TestSortCategoryStatus, TestSortCategoryStartTime, TestSortCategoryDuration, TestSortCategoryTestName:
		return true
	}
	return false
}

func (e TestSortCategory) String() string {
	return string(e)
}

func (e *TestSortCategory) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = TestSortCategory(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid TestSortCategory", str)
	}
	return nil
}

func (e TestSortCategory) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
