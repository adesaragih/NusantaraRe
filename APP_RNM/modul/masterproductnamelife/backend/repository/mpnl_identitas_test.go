package repository

// Audit 02-10-2026: sequence identitas yang tertinggal dari data (keadaan DEV: nilai berikut 200,
// `100202` sudah ada) tidak lagi menggagalkan produk baru - ID terpakai dilewati, tidak ditimpa.

import (
	"errors"
	"reflect"
	"testing"
)

// urutan - sequence tiruan mulai `mulai`, mencatat berapa kali diambil.
func urutan(mulai int64) (func() (int64, error), *int) {
	n, ambil := mulai, 0
	return func() (int64, error) {
		ambil++
		v := n
		n++
		return v, nil
	}, &ambil
}

func dipakai(ids ...string) func(string) (bool, error) {
	ada := map[string]bool{}
	for _, id := range ids {
		ada[id] = true
	}
	return func(id string) (bool, error) { return ada[id], nil }
}

func TestPilihIdentitasBebasMelewatiIDTerpakai(t *testing.T) {
	// Keadaan DEV 02-10-2026: 100200 dan 100201 bebas, 100202 terpakai.
	berikut, _ := urutan(200)
	id, lewat, err := PilihIdentitasBebas(berikut, dipakai("100202"))
	if err != nil || id != "100200" || len(lewat) != 0 {
		t.Fatalf("nomor bebas dipakai langsung: %q %v %v", id, lewat, err)
	}
	berikut, ambil := urutan(202)
	id, lewat, err = PilihIdentitasBebas(berikut, dipakai("100202", "100203"))
	if err != nil || id != "100204" || !reflect.DeepEqual(lewat, []string{"100202", "100203"}) || *ambil != 3 {
		t.Fatalf("ID terpakai dilewati ke nomor bebas berikut: %q %v %v (diambil %d)", id, lewat, err, *ambil)
	}
}

func TestPilihIdentitasBebasBatasLompatanGagalTerang(t *testing.T) {
	berikut, ambil := urutan(1)
	_, lewat, err := PilihIdentitasBebas(berikut, func(string) (bool, error) { return true, nil })
	if !errors.Is(err, ErrIdentitasBentrok) || len(lewat) != MaksLewatiIdentitas+1 || *ambil != MaksLewatiIdentitas+1 {
		t.Fatalf("lebih dari %d lompatan = bentrok terang: %v (dilewati %d, diambil %d)", MaksLewatiIdentitas, err, len(lewat), *ambil)
	}
}

func TestPilihIdentitasBebasLebarDanGalatTetapTerang(t *testing.T) {
	berikut, _ := urutan(100000)
	if _, _, err := PilihIdentitasBebas(berikut, dipakai()); !errors.Is(err, ErrIdentitasMelampauiLebar) {
		t.Errorf("nomor 6 digit tidak dipotong: %v", err)
	}
	gagal := errors.New("uji: pembacaan gagal")
	berikut, _ = urutan(5)
	if _, _, err := PilihIdentitasBebas(berikut, func(string) (bool, error) { return false, gagal }); !errors.Is(err, gagal) {
		t.Errorf("galat pemeriksaan diteruskan: %v", err)
	}
	if _, _, err := PilihIdentitasBebas(func() (int64, error) { return 0, gagal }, dipakai()); !errors.Is(err, gagal) {
		t.Errorf("galat sequence diteruskan: %v", err)
	}
}
