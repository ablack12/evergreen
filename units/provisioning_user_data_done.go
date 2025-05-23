package units

import (
	"context"
	"fmt"
	"time"

	"github.com/evergreen-ci/evergreen"
	"github.com/evergreen-ci/evergreen/model/event"
	"github.com/evergreen-ci/evergreen/model/host"
	"github.com/evergreen-ci/utility"
	"github.com/mongodb/amboy"
	"github.com/mongodb/amboy/job"
	"github.com/mongodb/amboy/registry"
	"github.com/mongodb/grip"
	"github.com/mongodb/grip/message"
	"github.com/mongodb/jasper/options"
	"github.com/pkg/errors"
)

const (
	userDataDoneJobName = "user-data-done"
)

type userDataDoneJob struct {
	HostID   string `bson:"host_id" json:"host_id" yaml:"host_id"`
	job.Base `bson:"base" json:"base" yaml:"base"`

	env      evergreen.Environment
	settings *evergreen.Settings
	host     *host.Host
}

func init() {
	registry.AddJobType(userDataDoneJobName, func() amboy.Job {
		return makeUserDataDoneJob()
	})
}

func makeUserDataDoneJob() *userDataDoneJob {
	j := &userDataDoneJob{
		Base: job.Base{
			JobType: amboy.JobType{
				Name:    userDataDoneJobName,
				Version: 0,
			},
		},
	}
	return j
}

// NewUserDataDoneJob creates a job that checks if the host is done provisioning
// with user data (if bootstrapped with user data). This check only applies to
// spawn hosts, since hosts running agents check into the server to verify their
// liveliness.
func NewUserDataDoneJob(env evergreen.Environment, hostID string, ts time.Time) amboy.Job {
	j := makeUserDataDoneJob()
	j.HostID = hostID
	j.env = env
	j.SetID(fmt.Sprintf("%s.%s.%s", userDataDoneJobName, j.HostID, ts.Format(TSFormat)))
	j.SetScopes([]string{fmt.Sprintf("%s.%s", userDataDoneJobName, hostID)})
	j.SetEnqueueAllScopes(true)
	j.UpdateRetryInfo(amboy.JobRetryOptions{
		Retryable:   utility.TruePtr(),
		MaxAttempts: utility.ToIntPtr(50),
		WaitUntil:   utility.ToTimeDurationPtr(20 * time.Second),
	})
	return j
}

func (j *userDataDoneJob) Run(ctx context.Context) {
	defer j.MarkComplete()

	defer func() {
		if j.HasErrors() && j.IsLastAttempt() {
			event.LogHostProvisionFailed(ctx, j.HostID, j.Error().Error())
			grip.Error(message.WrapError(j.Error(), message.Fields{
				"message":         "provisioning failed",
				"operation":       "user-data-done",
				"current_attempt": j.RetryInfo().CurrentAttempt,
				"distro":          j.host.Distro.Id,
				"provider":        j.host.Provider,
				"job":             j.ID(),
				"host_id":         j.host.Id,
				"host_tag":        j.host.Tag,
				"reason":          j.Error().Error(),
			}))
		}
	}()

	if err := j.populateIfUnset(ctx); err != nil {
		j.AddRetryableError(err)
		return
	}

	if j.host.Status != evergreen.HostStarting {
		j.UpdateRetryInfo(amboy.JobRetryOptions{
			NeedsRetry: utility.TruePtr(),
		})
		return
	}

	path := j.host.UserDataProvisioningDoneFile()

	if output, err := j.host.RunJasperProcess(ctx, j.env, &options.Create{
		Args: []string{
			j.host.Distro.ShellBinary(),
			"-l", "-c",
			fmt.Sprintf("ls %s", path)}}); err != nil {
		grip.Debug(message.WrapError(err, message.Fields{
			"message": "host was checked but is not yet ready",
			"output":  output,
			"host_id": j.host.Id,
			"distro":  j.host.Distro.Id,
			"job":     j.ID(),
		}))
		j.AddRetryableError(err)
		return
	}

	if j.host.IsVirtualWorkstation {
		if err := attachVolume(ctx, j.env, j.host); err != nil {
			grip.Error(message.WrapError(err, message.Fields{
				"message": "can't attach volume",
				"host_id": j.host.Id,
				"distro":  j.host.Distro.Id,
				"job":     j.ID(),
			}))
			j.AddError(err)
			j.AddError(j.host.SetStatus(ctx, evergreen.HostProvisionFailed, evergreen.User, "decommissioning host after failing to mount volume"))

			terminateJob := NewHostTerminationJob(j.env, j.host, HostTerminationOptions{
				TerminateIfBusy:   true,
				TerminationReason: "failed to mount volume",
			})

			j.AddError(EnqueueTerminateHostJob(ctx, j.env, terminateJob))

			return
		}
		if err := writeIcecreamConfig(ctx, j.env, j.host); err != nil {
			grip.Error(message.WrapError(err, message.Fields{
				"message": "can't write icecream config file",
				"host_id": j.host.Id,
				"distro":  j.host.Distro.Id,
				"job":     j.ID(),
			}))
		}
	}

	if j.host.ProvisionOptions != nil && j.host.ProvisionOptions.SetupScript != "" {
		// Run the spawn host setup script in a separate job to avoid forcing
		// this job to wait for task data to be loaded.
		j.AddError(j.env.RemoteQueue().Put(ctx, NewHostSetupScriptJob(j.env, j.host)))
	}

	j.finishJob(ctx)
}

func (j *userDataDoneJob) finishJob(ctx context.Context) {
	if err := j.host.SetUserDataHostProvisioned(ctx); err != nil {
		grip.Error(message.WrapError(err, message.Fields{
			"message": "could not mark host that has finished running user data as done provisioning",
			"host_id": j.host.Id,
			"distro":  j.host.Distro.Id,
			"job":     j.ID(),
		}))
		j.AddRetryableError(err)
		return
	}
}

func (j *userDataDoneJob) populateIfUnset(ctx context.Context) error {
	if j.host == nil {
		h, err := host.FindOneId(ctx, j.HostID)
		if err != nil {
			return errors.Wrapf(err, "finding host '%s'", j.HostID)
		}
		if h == nil {
			return errors.Errorf("host '%s' not found", j.HostID)
		}
		j.host = h
	}

	if j.env == nil {
		j.env = evergreen.GetEnvironment()
	}
	if j.settings == nil {
		j.settings = j.env.Settings()
	}

	return nil
}
