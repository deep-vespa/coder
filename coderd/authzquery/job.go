package authzquery

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"golang.org/x/xerrors"

	"github.com/coder/coder/coderd/database"
	"github.com/coder/coder/coderd/rbac"
	"github.com/coder/coder/coderd/util/slice"
)

func (q *AuthzQuerier) UpdateProvisionerJobWithCancelByID(ctx context.Context, arg database.UpdateProvisionerJobWithCancelByIDParams) error {
	job, err := q.database.GetProvisionerJobByID(ctx, arg.ID)
	if err != nil {
		return err
	}

	switch job.Type {
	case database.ProvisionerJobTypeWorkspaceBuild:
		build, err := q.database.GetWorkspaceBuildByJobID(ctx, arg.ID)
		if err != nil {
			return err
		}
		workspace, err := q.database.GetWorkspaceByID(ctx, build.WorkspaceID)
		if err != nil {
			return err
		}

		template, err := q.database.GetTemplateByID(ctx, workspace.TemplateID)
		if err != nil {
			return err
		}

		// Template can specify if cancels are allowed.
		// Would be nice to have a way in the rbac rego to do this.
		if !template.AllowUserCancelWorkspaceJobs {
			// Only owners can cancel workspace builds
			actor, ok := ActorFromContext(ctx)
			if !ok {
				return NoActorError
			}
			if !slice.Contains(actor.Roles.Names(), rbac.RoleOwner()) {
				return xerrors.Errorf("only owners can cancel workspace builds")
			}
		}

		err = q.authorizeContext(ctx, rbac.ActionUpdate, workspace)
		if err != nil {
			return err
		}
	case database.ProvisionerJobTypeTemplateVersionDryRun, database.ProvisionerJobTypeTemplateVersionImport:
		// Authorized call to get template version.
		templateVersion, err := authorizedTemplateVersionFromJob(ctx, q, job)
		if err != nil {
			return err
		}

		if templateVersion.TemplateID.Valid {
			template, err := q.database.GetTemplateByID(ctx, templateVersion.TemplateID.UUID)
			if err != nil {
				return err
			}
			err = q.authorizeContext(ctx, rbac.ActionUpdate, templateVersion.RBACObject(template))
			if err != nil {
				return err
			}
		} else {
			err = q.authorizeContext(ctx, rbac.ActionUpdate, templateVersion.RBACObjectNoTemplate())
			if err != nil {
				return err
			}
		}
	default:
		return xerrors.Errorf("unknown job type: %q", job.Type)
	}
	return q.database.UpdateProvisionerJobWithCancelByID(ctx, arg)
}

func (q *AuthzQuerier) GetProvisionerJobByID(ctx context.Context, id uuid.UUID) (database.ProvisionerJob, error) {
	job, err := q.database.GetProvisionerJobByID(ctx, id)
	if err != nil {
		return database.ProvisionerJob{}, err
	}

	switch job.Type {
	case database.ProvisionerJobTypeWorkspaceBuild:
		// Authorized call to get workspace build. If we can read the build, we
		// can read the job.
		_, err := q.GetWorkspaceBuildByJobID(ctx, id)
		if err != nil {
			return database.ProvisionerJob{}, err
		}
	case database.ProvisionerJobTypeTemplateVersionDryRun, database.ProvisionerJobTypeTemplateVersionImport:
		// Authorized call to get template version.
		_, err := authorizedTemplateVersionFromJob(ctx, q, job)
		if err != nil {
			return database.ProvisionerJob{}, err
		}
	default:
		return database.ProvisionerJob{}, xerrors.Errorf("unknown job type: %q", job.Type)
	}

	return job, nil
}

func authorizedTemplateVersionFromJob(ctx context.Context, q *AuthzQuerier, job database.ProvisionerJob) (database.TemplateVersion, error) {
	switch job.Type {
	case database.ProvisionerJobTypeTemplateVersionDryRun:
		// TODO: This is really unfortunate that we need to inspect the json
		// payload. We should fix this.
		tmp := struct {
			TemplateVersionID uuid.UUID `json:"template_version_id"`
		}{}
		err := json.Unmarshal(job.Input, &tmp)
		if err != nil {
			return database.TemplateVersion{}, xerrors.Errorf("dry-run unmarshal: %w", err)
		}
		// Authorized call to get template version.
		tv, err := q.GetTemplateVersionByID(ctx, tmp.TemplateVersionID)
		if err != nil {
			return database.TemplateVersion{}, err
		}
		return tv, nil
	case database.ProvisionerJobTypeTemplateVersionImport:
		// Authorized call to get template version.
		tv, err := q.GetTemplateVersionByJobID(ctx, job.ID)
		if err != nil {
			return database.TemplateVersion{}, err
		}
		return tv, nil
	default:
		return database.TemplateVersion{}, xerrors.Errorf("unknown job type: %q", job.Type)
	}
}