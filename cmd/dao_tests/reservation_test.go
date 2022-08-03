//go:build integration
// +build integration

package main

import (
	"context"
	"testing"

	"github.com/RHEnVision/provisioning-backend/internal/dao"
	_ "github.com/RHEnVision/provisioning-backend/internal/dao/sqlx"
	"github.com/RHEnVision/provisioning-backend/internal/db"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/RHEnVision/provisioning-backend/internal/testing/identity"
	"github.com/stretchr/testify/assert"
)

func createNoopReservation() *models.NoopReservation {
	return &models.NoopReservation{
		Reservation: models.Reservation{
			ID:        10,
			Provider:  models.ProviderTypeNoop,
			AccountID: 1,
			Status:    "Created",
		},
	}
}

func createInstancesReservation(reservationId int64) *models.ReservationInstance {
	return &models.ReservationInstance{
		ReservationID: reservationId,
		InstanceID:    "1",
	}
}

func createAWSReservation() *models.AWSReservation {
	return &models.AWSReservation{
		Reservation: models.Reservation{
			Provider:  models.ProviderTypeAWS,
			AccountID: 1,
			Status:    "Created",
		},
		PubkeyID: 1,
	}
}

func SetupReservation(t *testing.T, s string) (dao.ReservationDao, context.Context, error) {
	err := db.Seed("dao_test")
	if err != nil {
		t.Errorf("Error purging the database: %v", err)
		return nil, nil, err
	}
	ctx := identity.WithTenant(t, context.Background())
	reservationDao, err := dao.GetReservationDao(ctx)
	if err != nil {
		t.Errorf("%s test failed: %v", s, err)
		return nil, nil, err
	}
	return reservationDao, ctx, nil
}

func TestCreateNoop(t *testing.T) {
	CleanUpDatabase(t)
	reservationDao, ctx, err := SetupReservation(t, "Create noop reservation")
	if err != nil {
		t.Errorf("Database setup failed: %v", err)
		return
	}
	err = reservationDao.CreateNoop(ctx, createNoopReservation())
	if err != nil {
		t.Errorf("createNoop failed: %v", err)
		return
	}

	reservations, err := reservationDao.List(ctx, 10, 0)
	if err != nil {
		t.Errorf("list failed: %v", err)
		return
	}

	assert.Equal(t, 1, len(reservations), "CreateNoop error:.")
}

func TestCreateAWS(t *testing.T) {
	CleanUpDatabase(t)
	reservationDao, ctx, err := SetupReservation(t, "Create AWS reservation")
	if err != nil {
		t.Errorf("Database setup failed: %v", err)
		return
	}
	err = reservationDao.CreateAWS(ctx, createAWSReservation())
	if err != nil {
		t.Errorf("createAWS failed: %v", err)
		return
	}

	reservations, err := reservationDao.List(ctx, 10, 0)
	if err != nil {
		t.Errorf("list failed: %v", err)
		return
	}

	assert.Equal(t, 1, len(reservations), "Create AWS reservation error.")
}

func TestCreateInstance(t *testing.T) {
	CleanUpDatabase(t)
	reservationDao, ctx, err := SetupReservation(t, "Create Instance reservation")
	if err != nil {
		t.Errorf("Database setup failed: %v", err)
		return
	}

	reservation := createAWSReservation()
	err = reservationDao.CreateAWS(ctx, reservation)
	if err != nil {
		t.Errorf("createAWS failed: %v", err)
		return
	}

	err = reservationDao.CreateInstance(ctx, createInstancesReservation(reservation.ID))
	if err != nil {
		t.Errorf("createInstance failed: %v", err)
		return
	}

	reservations, err := reservationDao.ListInstances(ctx, 10, 0)
	if err != nil {
		t.Errorf("list failed: %v", err)
		return
	}

	assert.Equal(t, 1, len(reservations), "Create Instances reservation error.")
}

func TestListReservation(t *testing.T) {
	CleanUpDatabase(t)
	reservationDao, ctx, err := SetupReservation(t, "List reservation")
	if err != nil {
		t.Errorf("Database setup failed: %v", err)
		return
	}
	err = reservationDao.CreateAWS(ctx, createAWSReservation())
	if err != nil {
		t.Errorf("createAWS failed: %v", err)
		return
	}
	err = reservationDao.CreateNoop(ctx, createNoopReservation())
	if err != nil {
		t.Errorf("createNoop failed: %v", err)
		return
	}

	reservations, err := reservationDao.List(ctx, 10, 0)
	if err != nil {
		t.Errorf("list failed: %v", err)
		return
	}
	assert.Equal(t, 2, len(reservations), "List reservation error.")
}

func TestUpdateReservationIDForAWS(t *testing.T) {
	CleanUpDatabase(t)
	reservationDao, ctx, err := SetupReservation(t, "Update status reservation")
	if err != nil {
		t.Errorf("Database setup failed. %s", err)
		return
	}

	reservation := createAWSReservation()

	err = reservationDao.CreateAWS(ctx, reservation)
	if err != nil {
		t.Errorf("createAWS failed: %v", err)
		return
	}
	var count int

	err = db.DB.Get(&count, "SELECT COUNT(*) FROM aws_reservation_details")
	if err != nil {
		t.Errorf("count records in aws_reservation_details has failed: %v", err)
		return
	}
	assert.Equal(t, 1, count, "Number of aws reservations mismatch.")

	err = reservationDao.UpdateReservationIDForAWS(ctx, reservation.ID, "2")
	if err != nil {
		t.Errorf("UpdateReservationIDForAWS failed %s", err)
		return
	}

	var awsReservationId string
	err = db.DB.Get(&awsReservationId, "SELECT aws_reservation_id FROM aws_reservation_details WHERE reservation_id = $1", reservation.ID)
	if err != nil {
		t.Errorf("select aws_reservation_id from aws_reservation_details has failed: %v", err)
		return
	}
	assert.Equal(t, "2", awsReservationId, "Update reservation id  error: aws reservation id does not match.")

}

func TestUpdateStatusReservation(t *testing.T) {
	CleanUpDatabase(t)
	reservationDao, ctx, err := SetupReservation(t, "Update status reservation")
	if err != nil {
		t.Errorf("Database setup failed. %s", err)
		return
	}
	err = reservationDao.CreateNoop(ctx, createNoopReservation())
	if err != nil {
		t.Errorf("createNoop failed. %s", err)
		return
	}

	reservationsBefore, err := reservationDao.List(ctx, 10, 0)
	if err != nil {
		t.Errorf("list failed %s", err)
		return
	}

	err = reservationDao.UpdateStatus(ctx, reservationsBefore[0].ID, "Edited")
	if err != nil {
		t.Errorf("update status failed %s", err)
		return
	}

	reservationsAfter, err := reservationDao.List(ctx, 10, 0)
	if err != nil {
		t.Errorf("second list failed %s", err)
		return
	}
	assert.Equal(t, "Edited", reservationsAfter[0].Status, "Update status reservation error: status does not match.")
}

func TestDeleteReservation(t *testing.T) {
	CleanUpDatabase(t)
	reservationDao, ctx, err := SetupReservation(t, "Delete reservation")
	if err != nil {
		t.Errorf("Database setup failed: %v", err)
		return
	}
	err = reservationDao.CreateNoop(ctx, createNoopReservation())
	if err != nil {
		t.Errorf("createNoop failed: %v", err)
		return
	}

	reservationsBefore, err := reservationDao.List(ctx, 10, 0)
	if err != nil {
		t.Errorf("list failed: %v", err)
		return
	}

	err = reservationDao.Delete(ctx, reservationsBefore[0].ID)
	if err != nil {
		t.Errorf("delete failed: %v", err)
		return
	}
	reservationsAfter, err := reservationDao.List(ctx, 10, 0)
	if err != nil {
		t.Errorf("second list failed: %v", err)
		return
	}
	assert.Equal(t, len(reservationsBefore)-1, len(reservationsAfter), "Delete reservation error.")
}

func TestFinishReservation(t *testing.T) {
	CleanUpDatabase(t)
	reservationDao, ctx, err := SetupReservation(t, "Delete reservation")
	if err != nil {
		t.Errorf("Database setup failed: %v", err)
		return
	}
	err = reservationDao.CreateNoop(ctx, createNoopReservation())
	if err != nil {
		t.Errorf("createNoop failed: %v", err)
		return
	}

	reservationsBefore, err := reservationDao.List(ctx, 10, 0)
	if err != nil {
		t.Errorf("list failed: %v", err)
		return
	}

	err = reservationDao.Finish(ctx, reservationsBefore[0].ID, true, "Finished")
	if err != nil {
		t.Errorf("finish failed: %v", err)
		return
	}

	reservationsAfter, err := reservationDao.List(ctx, 10, 0)
	if err != nil {
		t.Errorf("second list failed: %v", err)
		return
	}

	assert.Equal(t, reservationsBefore[0].ID, reservationsAfter[0].ID, "Finish reservation error.")
	assert.Equal(t, true, reservationsAfter[0].Success.Bool, "Finish reservation error: success value does not match.")
	assert.Equal(t, "Finished", reservationsAfter[0].Status, "Finish reservation error: status does not match.")
}