package services

import (
	"context"
	"errors"
	"reflect"
	"strconv"
	"strings"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
)

// halamanTiruan - PembacaHalamanObjek tanpa basis data; argumen terakhir direkam.
type halamanTiruan struct {
	jumlah  map[string]int
	halaman models.HalamanObjek
	err     error
	arg     *[]any
	// tulis - argumen PerbaruiObjek terakhir; errTulis - galat PerbaruiObjek.
	tulis    *PerubahanObjek
	errTulis error
}

func (h halamanTiruan) PerbaruiObjek(_ context.Context, _ *db.Tx, id string, ubah []models.LokasiUbah, tambah []models.ObjekFire,
	hapus []string) (int, error) {
	if h.errTulis != nil {
		return 0, h.errTulis
	}
	if h.tulis != nil {
		*h.tulis = PerubahanObjek{Ubah: ubah, Tambah: tambah, Hapus: hapus}
	}
	return h.jumlah[id] + len(tambah) - len(hapus), nil
}

func (h halamanTiruan) JumlahLokasi(_ context.Context, id string) (int, error) {
	if h.err != nil {
		return 0, h.err
	}
	n, ok := h.jumlah[id]
	if !ok {
		return 0, repository.ErrKasusTidakAda
	}
	return n, nil
}

func (h halamanTiruan) BacaObjekHalaman(_ context.Context, id string, halaman, ukuran int, cari string) (models.HalamanObjek, error) {
	if h.arg != nil {
		*h.arg = []any{id, halaman, ukuran, cari}
	}
	if _, ok := h.jumlah[id]; !ok {
		return models.HalamanObjek{}, repository.ErrKasusTidakAda
	}
	return h.halaman, h.err
}

// TestBacaObjekHalaman - tiket 75 F3: batas halaman / ukuran (400 sebelum basis data), 503 tanpa pembaca, 404, olahan baca
// sama dengan GET utuh (tanggal kerugian DD-MM-YYYY, total per mata uang).
func TestBacaObjekHalaman(t *testing.T) {
	ctx := context.Background()
	for _, c := range [][2]int{{0, 10}, {1, 0}, {1, 501}} {
		if _, err := Baru(nil).BacaObjekHalaman(ctx, "UJI-NB-1", c[0], c[1], ""); !errors.Is(err, ErrMasukanObjek) {
			t.Errorf("%v: %v", c, err)
		}
	}
	if _, err := Baru(nil).BacaObjekHalaman(ctx, "UJI-NB-1", 1, 500, ""); !errors.Is(err, ErrObjekTanpaDatabase) {
		t.Errorf("tanpa pembaca: %v", err)
	}
	var arg []any
	o := models.ObjekFire{ObjectNo: "7", LossRecords: []models.CatatanKerugian{{DateOfLoss: "20260105"}},
		Items: []models.ItemObjek{{Currency: "IDR", Coverages: []models.CoverageObjek{{TSI: d("10")}}}}}
	svc := Baru(nil).DenganHalamanObjek(halamanTiruan{jumlah: map[string]int{"UJI-NB-1": 1201}, arg: &arg,
		halaman: models.HalamanObjek{Total: 3, Lokasi: []models.ObjekFire{o}, IDLokasi: []string{"91"}}})
	h, err := svc.BacaObjekHalaman(ctx, "UJI-NB-1", 2, 50, "jl")
	if err != nil || h.Total != 3 || h.IDLokasi[0] != "91" || h.Lokasi[0].LossRecords[0].DateOfLoss != "05-01-2026" ||
		len(h.Lokasi[0].TotalPerCurrency) != 1 || arg[1] != 2 || arg[2] != 50 || arg[3] != "jl" {
		t.Errorf("%v %+v %v", err, h, arg)
	}
	if _, err := svc.BacaObjekHalaman(ctx, "UJI-NB-9", 1, 10, ""); !errors.Is(err, ErrKasusTidakAda) {
		t.Errorf("404: %v", err)
	}
	if n := svc.JumlahLokasi(ctx, "UJI-NB-1"); n == nil || *n != 1201 {
		t.Errorf("jumlah: %v", n)
	}
	if svc.JumlahLokasi(ctx, "UJI-NB-9") != nil || Baru(nil).JumlahLokasi(ctx, "UJI-NB-1") != nil {
		t.Error("jumlah tak terbaca bukan nil")
	}
}

// TestPerbaruiObjek - PATCH tiket 75 F3: isian (400 sebelum basis data), 503, kunci unggah (409), penyiapan sama dengan PUT
// (loss ratio / coverage per baris) lalu ubah / tambah / hapus diteruskan ke repository; galat repository -> 409 / 404.
func TestPerbaruiObjek(t *testing.T) {
	ctx, akun := context.Background(), inti.Pelaku{AkunID: "UJI-USER"}
	sah := models.ObjekFire{ObjectType: "UJI"}
	banyak := func(n int) []string {
		x := make([]string, n)
		for i := range x {
			x[i] = strconv.Itoa(i + 1)
		}
		return x
	}
	for nama, c := range map[string]struct {
		p     PerubahanObjek
		pesan string
	}{
		"kosong":          {PerubahanObjek{}, "tidak ada perubahan"},
		"id bukan angka":  {PerubahanObjek{Hapus: []string{"1 OR 1=1"}}, "tidak sah"},
		"id ganda":        {PerubahanObjek{Ubah: []models.LokasiUbah{{ID: "5", Objek: sah}}, Hapus: []string{"5"}}, "ganda"},
		"hapus 501":       {PerubahanObjek{Hapus: banyak(501)}, "paling banyak 500 lokasi diubah"},
		"tambah 501":      {PerubahanObjek{Tambah: make([]models.ObjekFire, 501)}, "paling banyak 500 lokasi ditambah"},
		"ubah tanpa tipe": {PerubahanObjek{Ubah: []models.LokasiUbah{{ID: "5"}}}, "ubah: "},
		"tambah lantai":   {PerubahanObjek{Tambah: []models.ObjekFire{{ObjectType: "U", NumberOfFloor: "x"}}}, "tambah: "},
	} {
		if _, err := Baru(nil).PerbaruiObjek(ctx, akun, "UJI-NB-1", c.p); !errors.Is(err, ErrMasukanObjek) || !strings.Contains(err.Error(), c.pesan) {
			t.Errorf("%s: %v", nama, err)
		}
	}
	p := PerubahanObjek{Ubah: []models.LokasiUbah{{ID: "91", Objek: sah}}, Tambah: []models.ObjekFire{sah, sah}, Hapus: []string{"92"}}
	if _, err := Baru(nil).PerbaruiObjek(ctx, inti.Pelaku{}, "UJI-NB-1", p); !errors.Is(err, inti.ErrTanpaIdentitas) {
		t.Errorf("tanpa identitas: %v", err)
	}
	if _, err := Baru(nil).PerbaruiObjek(ctx, akun, "UJI-NB-1", p); !errors.Is(err, ErrObjekTanpaDatabase) {
		t.Errorf("tanpa DB: %v", err)
	}
	var tulis PerubahanObjek
	svc := Baru(nil).DenganTransaksi(tanpaTx).DenganHalamanObjek(halamanTiruan{jumlah: map[string]int{"UJI-NB-1": 1201}, tulis: &tulis})
	n, err := svc.PerbaruiObjek(ctx, akun, "UJI-NB-1", p)
	if err != nil || n != 1202 || len(tulis.Ubah) != 1 || tulis.Ubah[0].ID != "91" || len(tulis.Tambah) != 2 ||
		!reflect.DeepEqual(tulis.Hapus, []string{"92"}) || tulis.Ubah[0].Objek.LossRatio.OneYearAmount == nil {
		t.Errorf("%v %d %+v", err, n, tulis)
	}
	for galatRepo, mau := range map[error]error{repository.ErrLokasiBukanKasus: ErrLokasiObjekBerubah,
		repository.ErrLokasiBerSpreading: ErrLokasiBerSpreading, repository.ErrKasusTidakAda: ErrKasusTidakAda} {
		svc := Baru(nil).DenganTransaksi(tanpaTx).DenganHalamanObjek(halamanTiruan{jumlah: map[string]int{"UJI-NB-1": 1}, errTulis: galatRepo})
		if _, err := svc.PerbaruiObjek(ctx, akun, "UJI-NB-1", p); !errors.Is(err, mau) {
			t.Errorf("%v: %v", galatRepo, err)
		}
	}
	u := barUnggahTiruan(nil)
	_, _ = u.BuatUnggah(ctx, nil, models.UnggahCsv{CaseID: "UJI-NB-1", Status: StatusUnggahMenyimpan, Fase: FaseUnggahSimpan})
	kunci := Baru(nil).DenganTransaksi(tanpaTx).DenganHalamanObjek(halamanTiruan{jumlah: map[string]int{"UJI-NB-1": 1}}).DenganUnggahCsv(u, nil)
	if _, err := kunci.PerbaruiObjek(ctx, akun, "UJI-NB-1", p); !errors.Is(err, ErrUnggahBerjalan) {
		t.Errorf("kunci unggah: %v", err)
	}
}
