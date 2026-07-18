package main

import (
	"database/sql"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	randSource = rand.NewSource(time.Now().UnixNano())

	randRange = rand.New(randSource)
)

func getTestParcel() Parcel {
	return Parcel{
		Client:    1000,
		Status:    ParcelStatusRegistered,
		Address:   "test",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

func TestAddGetDelete(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()
	store := NewParcelStore(db)

	parcel := getTestParcel()

	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotEmpty(t, id)
	parcel.Number = id

	got, err := store.Get(id)
	require.NoError(t, err)
	assert.Equal(t, parcel, got)

	err = store.Delete(id)
	require.NoError(t, err)

	_, err = store.Get(id)
	require.ErrorIs(t, err, sql.ErrNoRows)
}

func TestSetAddress(t *testing.T) {

	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()
	store := NewParcelStore(db)

	parcel := getTestParcel()
	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	newAddress := "new test address"

	err = store.SetAddress(id, newAddress)
	require.NoError(t, err)

	p, err := store.Get(id)
	require.NoError(t, err)
	assert.Equal(t, newAddress, p.Address)

}

func TestSetStatus(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()
	store := NewParcelStore(db)

	parcel := getTestParcel()
	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	err = store.SetStatus(id, ParcelStatusRegistered)
	require.NoError(t, err)

	p, err := store.Get(id)
	require.NoError(t, err)
	assert.Equal(t, ParcelStatusRegistered, p.Status)
}

func TestGetByClient(t *testing.T) {

	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()
	store := NewParcelStore(db)

	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}
	parcelMap := map[int]Parcel{}

	client := randRange.Intn(10_000_000)
	parcels[0].Client = client
	parcels[1].Client = client
	parcels[2].Client = client

	for i := 0; i < len(parcels); i++ {
		id, err := store.Add(parcels[i])
		require.NoError(t, err)
		parcelMap[id] = parcels[i]

		parcels[i].Number = id

		parcelMap[id] = parcels[i]
	}

	storedParcels, err := store.GetByClient(client)
	require.NoError(t, err)
	assert.Len(t, storedParcels, len(parcelMap))

	for _, parcel := range storedParcels {
		origin, ok := parcelMap[parcel.Number]
		assert.True(t, ok)
		assert.Equal(t, origin, parcel)
	}
}
