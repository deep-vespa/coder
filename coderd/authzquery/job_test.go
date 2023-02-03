package authzquery_test

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"

	"github.com/coder/coder/coderd/database"
	"github.com/coder/coder/coderd/database/dbgen"
	"github.com/coder/coder/coderd/rbac"
)

func (suite *MethodTestSuite) TestProvsionerJob() {
	suite.Run("Build/GetProvisionerJobByID", func() {
		suite.RunMethodTest(func(t *testing.T, db database.Store) MethodCase {
			w := dbgen.Workspace(t, db, database.Workspace{})
			j := dbgen.ProvisionerJob(t, db, database.ProvisionerJob{
				Type: database.ProvisionerJobTypeWorkspaceBuild,
			})
			_ = dbgen.WorkspaceBuild(t, db, database.WorkspaceBuild{JobID: j.ID, WorkspaceID: w.ID})
			return methodCase(values(j.ID), asserts(w, rbac.ActionRead), values(j))
		})
	})
	suite.Run("TemplateVersion/GetProvisionerJobByID", func() {
		suite.RunMethodTest(func(t *testing.T, db database.Store) MethodCase {
			j := dbgen.ProvisionerJob(t, db, database.ProvisionerJob{
				Type: database.ProvisionerJobTypeTemplateVersionImport,
			})
			tpl := dbgen.Template(t, db, database.Template{})
			v := dbgen.TemplateVersion(t, db, database.TemplateVersion{
				TemplateID: uuid.NullUUID{UUID: tpl.ID, Valid: true},
				JobID:      j.ID,
			})
			return methodCase(values(j.ID), asserts(v.RBACObject(tpl), rbac.ActionRead), values(j))
		})
	})
	suite.Run("TemplateVersionDryRun/GetProvisionerJobByID", func() {
		suite.RunMethodTest(func(t *testing.T, db database.Store) MethodCase {
			tpl := dbgen.Template(t, db, database.Template{})
			v := dbgen.TemplateVersion(t, db, database.TemplateVersion{
				TemplateID: uuid.NullUUID{UUID: tpl.ID, Valid: true},
			})
			j := dbgen.ProvisionerJob(t, db, database.ProvisionerJob{
				Type: database.ProvisionerJobTypeTemplateVersionDryRun,
				Input: must(json.Marshal(struct {
					TemplateVersionID uuid.UUID `json:"template_version_id"`
				}{TemplateVersionID: v.ID})),
			})
			return methodCase(values(j.ID), asserts(v.RBACObject(tpl), rbac.ActionRead), values(j))
		})
	})
	suite.Run("Build/UpdateProvisionerJobWithCancelByID", func() {
		suite.RunMethodTest(func(t *testing.T, db database.Store) MethodCase {
			tpl := dbgen.Template(t, db, database.Template{AllowUserCancelWorkspaceJobs: true})
			w := dbgen.Workspace(t, db, database.Workspace{TemplateID: tpl.ID})
			j := dbgen.ProvisionerJob(t, db, database.ProvisionerJob{
				Type: database.ProvisionerJobTypeWorkspaceBuild,
			})
			_ = dbgen.WorkspaceBuild(t, db, database.WorkspaceBuild{JobID: j.ID, WorkspaceID: w.ID})
			return methodCase(values(database.UpdateProvisionerJobWithCancelByIDParams{ID: j.ID}), asserts(w, rbac.ActionUpdate), values())
		})
	})
	suite.Run("TemplateVersion/UpdateProvisionerJobWithCancelByID", func() {
		suite.RunMethodTest(func(t *testing.T, db database.Store) MethodCase {
			j := dbgen.ProvisionerJob(t, db, database.ProvisionerJob{
				Type: database.ProvisionerJobTypeTemplateVersionImport,
			})
			tpl := dbgen.Template(t, db, database.Template{})
			v := dbgen.TemplateVersion(t, db, database.TemplateVersion{
				TemplateID: uuid.NullUUID{UUID: tpl.ID, Valid: true},
				JobID:      j.ID,
			})
			return methodCase(values(database.UpdateProvisionerJobWithCancelByIDParams{ID: j.ID}),
				asserts(v.RBACObject(tpl), []rbac.Action{rbac.ActionRead, rbac.ActionUpdate}), values())
		})
	})
	suite.Run("TemplateVersionDryRun/UpdateProvisionerJobWithCancelByID", func() {
		suite.RunMethodTest(func(t *testing.T, db database.Store) MethodCase {
			tpl := dbgen.Template(t, db, database.Template{})
			v := dbgen.TemplateVersion(t, db, database.TemplateVersion{
				TemplateID: uuid.NullUUID{UUID: tpl.ID, Valid: true},
			})
			j := dbgen.ProvisionerJob(t, db, database.ProvisionerJob{
				Type: database.ProvisionerJobTypeTemplateVersionDryRun,
				Input: must(json.Marshal(struct {
					TemplateVersionID uuid.UUID `json:"template_version_id"`
				}{TemplateVersionID: v.ID})),
			})
			return methodCase(values(database.UpdateProvisionerJobWithCancelByIDParams{ID: j.ID}),
				asserts(v.RBACObject(tpl), []rbac.Action{rbac.ActionRead, rbac.ActionUpdate}), values())
		})
	})
	suite.Run("GetProvisionerJobsByIDs", func() {
		suite.RunMethodTest(func(t *testing.T, db database.Store) MethodCase {
			a := dbgen.ProvisionerJob(t, db, database.ProvisionerJob{})
			b := dbgen.ProvisionerJob(t, db, database.ProvisionerJob{})
			return methodCase(values([]uuid.UUID{a.ID, b.ID}), asserts(), values(a, b))
		})
	})
	suite.Run("GetProvisionerLogsByIDBetween", func() {
		suite.RunMethodTest(func(t *testing.T, db database.Store) MethodCase {
			w := dbgen.Workspace(t, db, database.Workspace{})
			j := dbgen.ProvisionerJob(t, db, database.ProvisionerJob{
				Type: database.ProvisionerJobTypeWorkspaceBuild,
			})
			_ = dbgen.WorkspaceBuild(t, db, database.WorkspaceBuild{JobID: j.ID, WorkspaceID: w.ID})
			return methodCase(values(database.GetProvisionerLogsByIDBetweenParams{
				JobID: j.ID,
			}), asserts(w, rbac.ActionRead), values([]database.ProvisionerJobLog{}))
		})
	})
}