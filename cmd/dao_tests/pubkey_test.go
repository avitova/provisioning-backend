//go:build integration
// +build integration

// To override application configuration for integration tests, copy local.yaml into this directory.

package main

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/RHEnVision/provisioning-backend/internal/config"
	"github.com/RHEnVision/provisioning-backend/internal/dao"
	_ "github.com/RHEnVision/provisioning-backend/internal/dao/sqlx"
	"github.com/RHEnVision/provisioning-backend/internal/db"
	"github.com/RHEnVision/provisioning-backend/internal/logging"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/RHEnVision/provisioning-backend/internal/testing/identity"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func createPk() *models.Pubkey {
	return &models.Pubkey{
		AccountID: 1,
		Name:      "lzap-ed25519-2021",
		Body:      "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIEhnn80ZywmjeBFFOGm+cm+5HUwm62qTVnjKlOdYFLHN lzap",
	}
}

func InitTestEnvironment() error {
	config.Initialize()
	log.Logger = logging.InitializeStdout()

	err := db.Initialize("integration")
	if err != nil {
		panic(fmt.Errorf("database setup had failed: %v", err))
	}
	return nil
}

func Setup(t *testing.T, s string) (dao.PubkeyDao, context.Context, error) {
	err := db.Seed("dao_test")
	if err != nil {
		t.Errorf("Error purging the database: %v", err)
		return nil, nil, err
	}
	ctx := identity.WithTenant(t, context.Background())
	pkDao, err := dao.GetPubkeyDao(ctx)
	if err != nil {
		t.Errorf("%s test had failed: %v", s, err)
		return nil, nil, err
	}
	return pkDao, ctx, nil
}

func CleanUpDatabase(t *testing.T) {
	err := db.Seed("drop_integration")
	if err != nil {
		t.Errorf("Error purging the database: %v", err)
		return
	}

	err = db.Migrate("integration")

	if err != nil {
		t.Errorf("Error running migration: %v", err)
		return
	}
}

func TestCreatePubkey(t *testing.T) {
	CleanUpDatabase(t)
	pkDao, ctx, err := Setup(t, "Create pubkey")
	if err != nil {
		t.Errorf("Database setup had failed: %v", err)
		return
	}
	pk := createPk()
	err = pkDao.Create(ctx, pk)
	if err != nil {
		t.Errorf("Create pubkey test had failed: %v", err)
		return
	}

	pk2, err := pkDao.GetById(ctx, pk.ID)
	if err != nil {
		t.Errorf("Create pubkey test had failed: %v", err)
		return
	}

	assert.Equal(t, pk.Name, pk2.Name, "Create pubkey test had failed.")
}

func TestListPubkey(t *testing.T) {
	CleanUpDatabase(t)
	pkDao, ctx, err := Setup(t, "List pubkey")
	if err != nil {
		t.Errorf("Database setup had failed: %v", err)
		return
	}
	err = pkDao.Create(ctx, createPk())
	pubkeys, err := pkDao.List(ctx, 100, 0)
	if err != nil {
		t.Errorf("List pubkey test had failed: %v", err)
		return
	}
	assert.Equal(t, 2, len(pubkeys), "List Pubkey error.")
}

func TestUpdatePubkey(t *testing.T) {
	CleanUpDatabase(t)
	updatePk := &models.Pubkey{
		ID:        1,
		AccountID: 1,
		Name:      "avitova-ed25519-2021",
		Body:      "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIEhnn80ZywmjeBFFOGm+cm+5HUwm62qTVnjKlOdYFLHN avitova",
	}
	pkDao, ctx, err := Setup(t, "Update pubkey")
	if err != nil {
		t.Errorf("Database setup had failed. %s", err)
		return
	}
	err = pkDao.Create(ctx, createPk())
	err = pkDao.Update(ctx, updatePk)
	if err != nil {
		t.Errorf("Update pubkey test had failed. %s", err)
		return
	}

	pubkeys, err := pkDao.List(ctx, 10, 0)
	if err != nil {
		t.Errorf("Update pubkey test had failed. %s", err)
		return
	}
	assert.Equal(t, updatePk.Name, pubkeys[0].Name, "Update pubkey test had failed.")
}

func TestGetPubkeyById(t *testing.T) {
	CleanUpDatabase(t)
	pkDao, ctx, err := Setup(t, "Get pubkey")
	if err != nil {
		t.Errorf("Database setup had failed. %s", err)
		return
	}
	err = pkDao.Create(ctx, createPk())
	if err != nil {
		t.Errorf("Delete pubkey test had failed. %s", err)
		return
	}
	pubkey, err := pkDao.GetById(ctx, 1)
	if err != nil {
		t.Errorf("Get pubkey test had failed.")
		return
	}
	assert.Equal(t, "lzap-ed25519-2021", pubkey.Name, "Get Pubkey error: pubkey name does not match.")
	assert.Equal(t, "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIEhnn80ZywmjeBFFOGm+cm+5HUwm62qTVnjKlOdYFLHN lzap", pubkey.Body, "Get Pubkey error: pubkey body does not match.")

}

func TestDeletePubkeyById(t *testing.T) {
	CleanUpDatabase(t)
	pkDao, ctx, err := Setup(t, "Delete pubkey")
	if err != nil {
		t.Errorf("Database setup had failed")
		return
	}
	err = pkDao.Create(ctx, createPk())
	if err != nil {
		t.Errorf("Delete pubkey test had failed. %s", err)
		return
	}
	pubkeys, err := pkDao.List(ctx, 10, 0)
	if err != nil {
		t.Errorf("Delete pubkey test had failed")
		return
	}
	err = pkDao.Delete(ctx, 1)
	if err != nil {
		t.Errorf("Delete pubkey test had failed")
		return
	}
	pubkeysAfter, err := pkDao.List(ctx, 10, 0)
	if err != nil {
		t.Errorf("Delete pubkey test had failed")
		return
	}
	assert.Equal(t, len(pubkeys)-1, len(pubkeysAfter), "Delete Pubkey error.")
}

func TestMain(t *testing.M) {
	InitTestEnvironment()
	exitVal := t.Run()
	os.Exit(exitVal)
}
