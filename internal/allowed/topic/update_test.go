package topic_test

import (
	"testing"

	"github.com/OpenSlides/openslides-permission-service/pkg/definitions"

	"github.com/OpenSlides/openslides-permission-service/internal/allowed"
	"github.com/OpenSlides/openslides-permission-service/internal/allowed/topic"

	"github.com/OpenSlides/openslides-permission-service/internal/tests"
)

func assertUpdateFailWithError(t *testing.T, params *allowed.IsAllowedParams) {
	allowed, addition, err := topic.Update(params)
	if nil != addition {
		t.Errorf("Expected to fail without an addition: %s", addition)
	}
	if nil == err {
		t.Errorf("Expected to fail with an error")
	}

	if allowed {
		t.Errorf("Expected to fail with allowed=false")
	}
}

func assertUpdateIsNotAllowed(t *testing.T, params *allowed.IsAllowedParams) {
	allowed, addition, err := topic.Update(params)
	if nil != addition {
		t.Errorf("Expected to fail without an addition: %s", addition)
	}
	if nil != err {
		t.Errorf("Expected to fail without an error (error: %s)", err)
	}

	if allowed {
		t.Errorf("Expected to fail with allowed=false")
	}
}

func assertUpdateIsAllowed(t *testing.T, params *allowed.IsAllowedParams) {
	allowed, addition, err := topic.Update(params)
	if nil != addition {
		t.Errorf("Expected to fail without an addition: %s", addition)
	}
	if nil != err {
		t.Errorf("Expected to fail without an error (error: %s)", err)
	}

	if !allowed {
		t.Errorf("Expected to be allowed")
	}
}

func TestUpdateUnknownUser(t *testing.T) {
	dp := tests.NewTestDataProvider()
	data := definitions.FqfieldData{
		"id": "1",
	}
	params := &allowed.IsAllowedParams{UserId: 1, Data: data, DataProvider: dp.GetDataprovider()}

	assertUpdateFailWithError(t, params)
}

func TestUpdateSuperadminRole(t *testing.T) {
	dp := tests.NewTestDataProvider()
	data := definitions.FqfieldData{} // No meeting id needed, it is always possible.
	dp.AddUserWithSuperadminRole(1)
	params := &allowed.IsAllowedParams{UserId: 1, Data: data, DataProvider: dp.GetDataprovider()}

	assertUpdateIsAllowed(t, params)
}

func TestUpdateNoId(t *testing.T) {
	dp := tests.NewTestDataProvider()
	data := definitions.FqfieldData{}
	dp.AddUserWithAdminGroupToMeeting(1, 1)
	params := &allowed.IsAllowedParams{UserId: 1, Data: data, DataProvider: dp.GetDataprovider()}

	assertUpdateFailWithError(t, params)
}

func TestUpdateUserNotInMeeting(t *testing.T) {
	dp := tests.NewTestDataProvider()
	addBasicTopic(dp)
	data := definitions.FqfieldData{
		"id": "1",
	}
	dp.AddUser(1)
	params := &allowed.IsAllowedParams{UserId: 1, Data: data, DataProvider: dp.GetDataprovider()}

	assertUpdateIsNotAllowed(t, params)
}

func TestUpdateAdminUser(t *testing.T) {
	dp := tests.NewTestDataProvider()
	addBasicTopic(dp)
	data := definitions.FqfieldData{
		"id": "1",
	}
	dp.AddUserWithAdminGroupToMeeting(1, 1)
	params := &allowed.IsAllowedParams{UserId: 1, Data: data, DataProvider: dp.GetDataprovider()}

	assertUpdateIsAllowed(t, params)
}

func TestUpdateUser(t *testing.T) {
	dp := tests.NewTestDataProvider()
	addBasicTopic(dp)
	data := definitions.FqfieldData{
		"id": "1",
	}
	dp.AddUserToMeeting(1, 1)
	dp.AddPermissionToGroup(1, "agenda.can_manage")
	params := &allowed.IsAllowedParams{UserId: 1, Data: data, DataProvider: dp.GetDataprovider()}

	assertUpdateIsAllowed(t, params)
}

func TestUpdateUserNoPermissions(t *testing.T) {
	dp := tests.NewTestDataProvider()
	addBasicTopic(dp)
	data := definitions.FqfieldData{
		"id": "1",
	}
	dp.AddUserToMeeting(1, 1)
	params := &allowed.IsAllowedParams{UserId: 1, Data: data, DataProvider: dp.GetDataprovider()}

	assertUpdateIsNotAllowed(t, params)
}

func TestUpdateInvaldFields(t *testing.T) {
	dp := tests.NewTestDataProvider()
	dp.AddUserWithSuperadminRole(1)
	data := definitions.FqfieldData{
		"not_allowed": "some value",
	}
	params := &allowed.IsAllowedParams{UserId: 1, Data: data, DataProvider: dp.GetDataprovider()}

	assertUpdateFailWithError(t, params)
}

func TestUpdateDisabledAnonymous(t *testing.T) {
	dp := tests.NewTestDataProvider()
	addBasicTopic(dp)
	data := definitions.FqfieldData{
		"id": "1",
	}
	params := &allowed.IsAllowedParams{UserId: 0, Data: data, DataProvider: dp.GetDataprovider()}

	assertUpdateIsNotAllowed(t, params)
}

func TestUpdateEnabledAnonymous(t *testing.T) {
	dp := tests.NewTestDataProvider()
	addBasicTopic(dp)
	dp.EnableAnonymous()
	data := definitions.FqfieldData{
		"id": "1",
	}
	params := &allowed.IsAllowedParams{UserId: 0, Data: data, DataProvider: dp.GetDataprovider()}

	assertUpdateIsNotAllowed(t, params)
}

func TestUpdateEnabledAnonymousWithPermissions(t *testing.T) {
	dp := tests.NewTestDataProvider()
	addBasicTopic(dp)
	dp.EnableAnonymous()
	dp.AddPermissionToGroup(1, "agenda.can_manage")
	data := definitions.FqfieldData{
		"id": "1",
	}
	params := &allowed.IsAllowedParams{UserId: 0, Data: data, DataProvider: dp.GetDataprovider()}

	assertUpdateIsAllowed(t, params)
}
