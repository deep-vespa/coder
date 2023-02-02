package authzquery_test

import (
	"testing"

	"github.com/google/uuid"

	"github.com/coder/coder/coderd/database"
	"github.com/coder/coder/coderd/database/dbgen"
	"github.com/coder/coder/coderd/rbac"
)

func (suite *MethodTestSuite) TestGroup() {
	suite.Run("DeleteGroupByID", func() {
		suite.RunMethodTest(func(t *testing.T, db database.Store) MethodCase {
			g := dbgen.Group(t, db, database.Group{})
			return methodCase(inputs(g.ID), asserts(g, rbac.ActionDelete))
		})
	})
	suite.Run("DeleteGroupMemberFromGroup", func() {
		suite.RunMethodTest(func(t *testing.T, db database.Store) MethodCase {
			g := dbgen.Group(t, db, database.Group{})
			m := dbgen.GroupMember(t, db, database.GroupMember{
				GroupID: g.ID,
			})
			return methodCase(inputs(database.DeleteGroupMemberFromGroupParams{
				UserID:  m.UserID,
				GroupID: g.ID,
			}), asserts(g, rbac.ActionUpdate))
		})
	})
	suite.Run("GetGroupByID", func() {
		suite.RunMethodTest(func(t *testing.T, db database.Store) MethodCase {
			g := dbgen.Group(t, db, database.Group{})
			return methodCase(inputs(g.ID), asserts(g, rbac.ActionRead))
		})
	})
	suite.Run("GetGroupByOrgAndName", func() {
		suite.RunMethodTest(func(t *testing.T, db database.Store) MethodCase {
			g := dbgen.Group(t, db, database.Group{})
			return methodCase(inputs(database.GetGroupByOrgAndNameParams{
				OrganizationID: g.OrganizationID,
				Name:           g.Name,
			}), asserts(g, rbac.ActionRead))
		})
	})
	suite.Run("GetGroupMembers", func() {
		suite.RunMethodTest(func(t *testing.T, db database.Store) MethodCase {
			g := dbgen.Group(t, db, database.Group{})
			_ = dbgen.GroupMember(t, db, database.GroupMember{})
			return methodCase(inputs(g.ID), asserts(g, rbac.ActionRead))
		})
	})
	suite.Run("InsertAllUsersGroup", func() {
		suite.RunMethodTest(func(t *testing.T, db database.Store) MethodCase {
			o := dbgen.Organization(t, db, database.Organization{})
			return methodCase(inputs(o.ID), asserts(rbac.ResourceGroup.InOrg(o.ID), rbac.ActionCreate))
		})
	})
	suite.Run("InsertGroupMember", func() {
		suite.RunMethodTest(func(t *testing.T, db database.Store) MethodCase {
			g := dbgen.Group(t, db, database.Group{})
			return methodCase(inputs(database.InsertGroupMemberParams{
				UserID:  uuid.New(),
				GroupID: g.ID,
			}), asserts(g, rbac.ActionUpdate))
		})
	})
	suite.Run("UpdateGroupByID", func() {
		suite.RunMethodTest(func(t *testing.T, db database.Store) MethodCase {
			g := dbgen.Group(t, db, database.Group{})
			return methodCase(inputs(database.UpdateGroupByIDParams{
				Name: "new-name",
				ID:   g.ID,
			}), asserts(g, rbac.ActionUpdate))
		})
	})
}
