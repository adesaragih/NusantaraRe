package handlers

import (
	"context"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
	"nusantarare/modul/nbfacin/backend/services"
)

// halamanTiruanH - PembacaHalamanObjek tanpa basis data: case UJI-NB-1 ber-1201 lokasi, satu lokasi per halaman.
type halamanTiruanH struct{}

func (halamanTiruanH) JumlahLokasi(_ context.Context, id string) (int, error) {
	if id != "UJI-NB-1" {
		return 0, repository.ErrKasusTidakAda
	}
	return 1201, nil
}

func (halamanTiruanH) BacaObjekHalaman(_ context.Context, id string, _, _ int, cari string) (models.HalamanObjek, error) {
	if id != "UJI-NB-1" {
		return models.HalamanObjek{}, repository.ErrKasusTidakAda
	}
	return models.HalamanObjek{Total: 1201, Lokasi: []models.ObjekFire{{ObjectNo: "7", ObjectName: cari}}, IDLokasi: []string{"91"}}, nil
}

func (halamanTiruanH) PerbaruiObjek(_ context.Context, _ *db.Tx, _ string, ubah []models.LokasiUbah, tambah []models.ObjekFire,
	hapus []string) (int, error) {
	for _, u := range ubah {
		if u.ID == "404" {
			return 0, repository.ErrLokasiBukanKasus
		}
		if u.ID == "777" {
			return 0, repository.ErrLokasiBerSpreading
		}
	}
	return 1201 + len(tambah) - len(hapus), nil
}

// TestObjekHalamanKabel - tiket 75 F3 lewat HTTP: `{total, halaman, ukuran, lokasi: [{id, …objek}]}`; halaman bukan angka /
// ukuran > 500 -> 400; case lain 404; GET kasus memuat jumlahLokasi (null tanpa pembaca).
func TestObjekHalamanKabel(t *testing.T) {
	svc := layananKasus().DenganHalamanObjek(halamanTiruanH{})
	kode, isi := minta(t, svc, "GET", "/api/nbfacin/kasus/UJI-NB-1/objek/halaman?halaman=3&ukuran=1&cari=UJI", "", "")
	if kode != 200 || !strings.HasPrefix(isi, `{"total":1201,"halaman":3,"ukuran":1,"lokasi":[{"id":"91","objectNo":"7","objectType":"","objectName":"UJI"`) {
		t.Errorf("%d %s", kode, isi)
	}
	for jalur, mau := range map[string]int{"UJI-NB-1/objek/halaman?halaman=x&ukuran=1": 400, "UJI-NB-1/objek/halaman?halaman=1&ukuran=501": 400,
		"UJI-NB-9/objek/halaman?halaman=1&ukuran=10": 404} {
		if kode, isi := minta(t, svc, "GET", "/api/nbfacin/kasus/"+jalur, "", ""); kode != mau {
			t.Errorf("%s: %d %s", jalur, kode, isi)
		}
	}
	if _, isi := minta(t, svc, "GET", "/api/nbfacin/kasus/UJI-NB-1", "", ""); !strings.Contains(isi, `"jumlahLokasi":1201`) {
		t.Errorf("kasus: %s", isi)
	}
	if _, isi := minta(t, layananKasus(), "GET", "/api/nbfacin/kasus/UJI-NB-1", "", ""); !strings.Contains(isi, `"jumlahLokasi":null`) {
		t.Errorf("kasus tanpa pembaca: %s", isi)
	}
}

// TestPerbaruiObjekKabel - PATCH …/objek: `{ubah: [{id, …}], tambah, hapus}` -> `{jumlahLokasi}`; JSON rusak 400; lokasi
// berubah / ber-spreading 409 berpesan.
func TestPerbaruiObjekKabel(t *testing.T) {
	svc := layananKasus().DenganHalamanObjek(halamanTiruanH{})
	jalur := "/api/nbfacin/kasus/UJI-NB-1/objek"
	kode, isi := minta(t, svc, "PATCH", jalur, `{"ubah":[{"id":"91","objectType":"UJI"}],"tambah":[{"objectType":"UJI"},{"objectType":"UJI"}],"hapus":["92","93","94"]}`, "UJI-USER")
	if kode != 200 || isi != `{"jumlahLokasi":1200}` {
		t.Errorf("%d %s", kode, isi)
	}
	if kode, _ := minta(t, svc, "PATCH", jalur, `{"ubah":`, "UJI-USER"); kode != 400 {
		t.Errorf("JSON rusak: %d", kode)
	}
	for id, pesan := range map[string]string{"404": "Lokasi objek sudah berubah", "777": "sudah memiliki spreading"} {
		kode, isi := minta(t, svc, "PATCH", jalur, `{"ubah":[{"id":"`+id+`","objectType":"UJI"}]}`, "UJI-USER")
		if kode != 409 || !strings.Contains(isi, pesan) {
			t.Errorf("%s: %d %s", id, kode, isi)
		}
	}
}

// TestAmbangHalamanObjekSamaDenganApiTs - ambang mode per halaman backend = frontend api.ts AMBANG_HALAMAN_OBJEK (0f, 10-10-2026).
func TestAmbangHalamanObjekSamaDenganApiTs(t *testing.T) {
	isi, err := os.ReadFile(filepath.Join("..", "..", "frontend", "api.ts"))
	if err != nil {
		t.Fatal(err)
	}
	m := regexp.MustCompile(`(?m)^export const AMBANG_HALAMAN_OBJEK = (\d+)\s*$`).FindSubmatch(isi)
	if m == nil {
		t.Fatal("api.ts tanpa AMBANG_HALAMAN_OBJEK")
	}
	if n, _ := strconv.Atoi(string(m[1])); n != services.AmbangHalamanObjek {
		t.Errorf("api.ts %d, services %d", n, services.AmbangHalamanObjek)
	}
}
