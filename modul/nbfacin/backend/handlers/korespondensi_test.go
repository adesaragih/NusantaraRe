package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
	"nusantarare/modul/nbfacin/backend/services"
)

type korespondensiTiruanH struct{ daftar []models.Korespondensi }

func (k *korespondensiTiruanH) DaftarKorespondensi(context.Context, string) ([]models.Korespondensi, error) {
	return append([]models.Korespondensi{}, k.daftar...), nil
}
func (k *korespondensiTiruanH) TambahKorespondensi(_ context.Context, _ *db.Tx, _ string, x models.Korespondensi, _ string) (string, error) {
	x.No = "7"
	k.daftar = append(k.daftar, x)
	return "7", nil
}
func (k *korespondensiTiruanH) UbahKorespondensi(_ context.Context, _ *db.Tx, _, no string, isi models.IsiKorespondensi) error {
	for i := range k.daftar {
		if k.daftar[i].No == no {
			if k.daftar[i].Terkirim {
				return repository.ErrKorespondensiTerkirim
			}
			k.daftar[i].Subject, k.daftar[i].Message, k.daftar[i].Penerima = isi.Subject, isi.Message, isi.Penerima
			return nil
		}
	}
	return repository.ErrKorespondensiTidakAda
}
func (k *korespondensiTiruanH) HapusKorespondensi(context.Context, *db.Tx, string, string) error {
	k.daftar = nil
	return nil
}
func (k *korespondensiTiruanH) TandaiTerkirim(context.Context, *db.Tx, string, string, string) error {
	k.daftar[0].Terkirim = true
	return nil
}
func (k *korespondensiTiruanH) AkunLogin(context.Context, string) (string, string, bool, error) {
	return "UJI PENGIRIM", "uji.p@contoh.test", true, nil
}
func (k *korespondensiTiruanH) SaranPenerima(_ context.Context, kata string) ([]models.SaranAkumulasi, error) {
	return []models.SaranAkumulasi{{ID: "UJI-A", Label: "UJI " + kata, Ekstra: "uji.a@contoh.test"}}, nil
}

type surelTiruanH struct{ galat error }

func (s surelTiruanH) KirimSurel(context.Context, models.Surel) error { return s.galat }

// TestKorespondensiHandler - bentuk persis kontrak; Add / Save / Send / Delete; 204, 400, 409, 502, 503 (tanpa SMTP);
// jalur ditolak (tab hanya JOB_POSITION IT): keenam rute 401 tanpa sesi / 403 bukan IT, sebelum badan diurai.
func TestKorespondensiHandler(t *testing.T) {
	kt := &korespondensiTiruanH{}
	svc := services.Baru(nil).DenganKorespondensi(kt, kt).DenganJabatan(jabatanTiruanH{"UJI-USER": "IT", "UJI-UW": "Underwriter"}).
		DenganTransaksi(tanpaOracle)
	jalur := "/api/nbfacin/kasus/UJI-NB-1/korespondensi"
	kode, isi := minta(t, svc, "POST", jalur, "", "UJI-USER")
	if kode != 200 || !strings.Contains(isi, `"no":"7"`) || !strings.Contains(isi, `"sender":"UJI PENGIRIM","senderMail":"uji.p@contoh.test","subject":"","message":"","penerima":[],"terkirim":false}`) {
		t.Fatalf("tambah: %d %s", kode, isi)
	}
	badan := `{"subject":"UJI","message":"<p style=x>halo</p><script>x</script>","penerima":[{"recipient":"UJI A","recipientMail":"uji.a@contoh.test"}]}`
	kode, isi = minta(t, svc, "PUT", jalur+"/7", badan, "UJI-USER")
	var simpan struct{ Message string }
	if err := json.Unmarshal([]byte(isi), &simpan); err != nil || kode != 200 || simpan.Message != "<p>halo</p>" {
		t.Errorf("simpan: %d %s", kode, isi)
	}
	if kode, isi := minta(t, svc, "POST", jalur+"/7/kirim", badan, "UJI-USER"); kode != 503 || !strings.Contains(isi, "SMTP belum dikonfigurasi") {
		t.Errorf("503: %d %s", kode, isi)
	}
	svc.DenganPengirimSurel(surelTiruanH{galat: errors.New("UJI rincian server")})
	if kode, isi := minta(t, svc, "POST", jalur+"/7/kirim", badan, "UJI-USER"); kode != 502 || strings.Contains(isi, "UJI rincian server") {
		t.Errorf("502 (rincian server tidak bocor ke layar): %d %s", kode, isi)
	}
	svc.DenganPengirimSurel(surelTiruanH{galat: repository.ErrSMTPTanpaTLS})
	if kode, isi := minta(t, svc, "POST", jalur+"/7/kirim", badan, "UJI-USER"); kode != 502 || !strings.Contains(isi, "tidak mendukung TLS") {
		t.Errorf("502 tanpa TLS: %d %s", kode, isi)
	}
	svc.DenganPengirimSurel(surelTiruanH{})
	if kode, isi := minta(t, svc, "POST", jalur+"/7/kirim", badan, "UJI-USER"); kode != 200 || !strings.Contains(isi, `"terkirim":true`) {
		t.Errorf("kirim: %d %s", kode, isi)
	}
	if kode, _ := minta(t, svc, "PUT", jalur+"/7", badan, "UJI-USER"); kode != 409 {
		t.Errorf("409: %d", kode)
	}
	if kode, _ := minta(t, svc, "PUT", jalur+"/8", badan, "UJI-USER"); kode != 404 {
		t.Errorf("404: %d", kode)
	}
	if kode, _ := minta(t, svc, "PUT", jalur+"/7", `{"subject":"x","rasio":1}`, "UJI-USER"); kode != 400 {
		t.Errorf("400 kunci asing: %d", kode)
	}
	if kode, _ := minta(t, svc, "PUT", jalur+"/7", badan, ""); kode != 401 {
		t.Errorf("401: %d", kode)
	}
	if kode, isi := minta(t, svc, "GET", jalur, "", "UJI-USER"); kode != 200 || !strings.HasPrefix(isi, `{"daftar":[{"no":"7"`) || !strings.HasSuffix(isi, `"pesan":[]}`) {
		t.Errorf("daftar: %d %s", kode, isi)
	}
	for _, r := range [][2]string{{"GET", jalur}, {"POST", jalur}, {"PUT", jalur + "/7"}, {"DELETE", jalur + "/7"},
		{"POST", jalur + "/7/kirim"}, {"GET", "/api/nbfacin/korespondensi/penerima/saran?q=a"}} {
		for pelaku, mau := range map[string]int{"": 401, "UJI-UW": 403} {
			if kode, _ := minta(t, svc, r[0], r[1], "bukan json", pelaku); kode != mau {
				t.Errorf("%s %s %q: %d, mau %d", r[0], r[1], pelaku, kode, mau)
			}
		}
	}
	if kode, isi := minta(t, svc, "DELETE", jalur+"/7", "", "UJI-USER"); kode != 204 || isi != "" {
		t.Errorf("hapus: %d %q", kode, isi)
	}
	if kode, isi := minta(t, svc, "GET", "/api/nbfacin/korespondensi/penerima/saran?q=a", "", "UJI-USER"); kode != 200 ||
		isi != `{"baris":[{"id":"UJI-A","label":"UJI a","ekstra":"uji.a@contoh.test"}]}` {
		t.Errorf("saran: %d %s", kode, isi)
	}
}
