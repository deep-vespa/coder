package authzquery_test

import (
	"testing"

	"github.com/coder/coder/coderd/database"
	"github.com/coder/coder/coderd/database/dbgen"
	"github.com/coder/coder/coderd/rbac"
)

func (suite *MethodTestSuite) TestFile() {
	suite.Run("GetFileByHashAndCreator", func() {
		suite.RunMethodTest(func(t *testing.T, db database.Store) MethodCase {
			f := dbgen.File(t, db, database.File{})
			return methodCase(inputs(database.GetFileByHashAndCreatorParams{
				Hash:      f.Hash,
				CreatedBy: f.CreatedBy,
			}), asserts(f, rbac.ActionRead))
		})
	})
	suite.Run("GetFileByID", func() {
		suite.RunMethodTest(func(t *testing.T, db database.Store) MethodCase {
			f := dbgen.File(t, db, database.File{})
			return methodCase(inputs(f.ID), asserts(f, rbac.ActionRead))
		})
	})
	suite.Run("InsertFile", func() {
		suite.RunMethodTest(func(t *testing.T, db database.Store) MethodCase {
			u := dbgen.User(t, db, database.User{})
			return methodCase(inputs(database.InsertFileParams{
				CreatedBy: u.ID,
			}), asserts(rbac.ResourceFile.WithOwner(u.ID.String()), rbac.ActionCreate))
		})
	})
}
