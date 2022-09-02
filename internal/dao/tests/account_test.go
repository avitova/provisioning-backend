//go:build integration
// +build integration

package tests

import (
	"context"
	"database/sql"
	"testing"

	"github.com/RHEnVision/provisioning-backend/internal/dao"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/RHEnVision/provisioning-backend/internal/testing/identity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createAccount() *models.Account {
	return &models.Account{
		OrgID:         "2",
		AccountNumber: sql.NullString{String: "100", Valid: true},
	}
}

func createAccountWithNullAccountNumber() *models.Account {
	return &models.Account{
		OrgID:         "2",
		AccountNumber: sql.NullString{},
	}
}

func setupAccount(t *testing.T) (dao.AccountDao, context.Context) {
	setup()
	ctx := identity.WithTenant(t, context.Background())
	accDao, err := dao.GetAccountDao(ctx)
	if err != nil {
		panic(err)
	}
	return accDao, ctx
}

func teardownAccount(_ *testing.T) {
	teardown()
}

func TestCreateAccount(t *testing.T) {
	accDao, ctx := setupAccount(t)
	defer teardownAccount(t)
	acc := createAccount()
	err := accDao.Create(ctx, acc)
	require.NoError(t, err)
	account, err := accDao.GetById(ctx, 2)
	require.NoError(t, err)

	assert.Equal(t, acc, account)
}

func TestCreateAccountWithNullAccountNumber(t *testing.T) {
	accDao, ctx := setupAccount(t)
	defer teardownAccount(t)
	acc := createAccountWithNullAccountNumber()
	err := accDao.Create(ctx, acc)
	require.NoError(t, err)
	account, err := accDao.GetById(ctx, 2)
	require.NoError(t, err)

	assert.Equal(t, acc, account)
}

func TestListAccount(t *testing.T) {
	accDao, ctx := setupAccount(t)
	defer teardownAccount(t)
	acc := createAccount()
	err := accDao.Create(ctx, acc)
	accounts, err := accDao.List(ctx, 100, 0)
	require.NoError(t, err)

	assert.Equal(t, 2, len(accounts))
	require.Contains(t, accounts, acc)
}

func TestGetByIdAccount(t *testing.T) {
	accDao, ctx := setupAccount(t)
	defer teardownAccount(t)
	account, err := accDao.GetById(ctx, 1)
	require.NoError(t, err)

	assert.Equal(t, "1", account.OrgID)
	assert.Equal(t, "1", account.AccountNumber.String)
}

func TestGetByAccountNumber(t *testing.T) {
	accDao, ctx := setupAccount(t)
	defer teardownAccount(t)
	account, err := accDao.GetByAccountNumber(ctx, "1")
	require.NoError(t, err)

	assert.Equal(t, "1", account.OrgID)
}

func TestGetByOrgId(t *testing.T) {
	accDao, ctx := setupAccount(t)
	defer teardownAccount(t)
	account, err := accDao.GetByOrgId(ctx, "1")
	require.NoError(t, err)

	assert.Equal(t, int64(1), account.ID)
	assert.Equal(t, "1", account.AccountNumber.String)
}

func TestGetOrCreateByIdentityGet(t *testing.T) {
	accDao, ctx := setupAccount(t)
	defer teardownAccount(t)
	account, err := accDao.GetOrCreateByIdentity(ctx, "1", "1")
	require.NoError(t, err)

	assert.Equal(t, "1", account.OrgID)
	assert.Equal(t, "1", account.AccountNumber.String)
}

func TestGetOrCreateByIdentityAccountCreate(t *testing.T) {
	accDao, ctx := setupAccount(t)
	defer teardownAccount(t)
	accountsBefore, err := accDao.List(ctx, 100, 0)
	require.NoError(t, err)
	_, err = accDao.GetOrCreateByIdentity(ctx, "2", "100")
	require.NoError(t, err)
	accountsAfter, err := accDao.List(ctx, 100, 0)
	require.NoError(t, err)
	account, err := accDao.GetByOrgId(ctx, "2")
	require.NoError(t, err)

	assert.Equal(t, len(accountsBefore)+1, len(accountsAfter))
	assert.Equal(t, "2", account.OrgID)
	assert.Equal(t, "100", account.AccountNumber.String)
	assert.Equal(t, true, account.AccountNumber.Valid)
}
