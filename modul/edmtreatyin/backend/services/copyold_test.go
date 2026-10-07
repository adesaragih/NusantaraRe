package services_test

// Copy Old (perintah work owner 07-10-2026 "BUATKAN TOMBOL COPY OLD SAMA SEPERTI MASTER PRODUCTNAME LIFE, KHUSUS BUAT
// SUPERUSER") di atas fixture pemuat (`pemuat_test.go`: rantai polis A tiga generasi + percabangan, polis B keutuhan,
// polis C tanpa generasi NB). Seam HTTP lewat `handlers.Router` - gerbang superadmin = pemegang menu Kelola User
// dengan menu EDM Treaty In ber-hak PENUH (pola Bordereaux). Fixture UJI-.

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/menu"
	"nusantarare/modul/edmtreatyin/backend/handlers"
	"nusantarare/modul/edmtreatyin/backend/models"
	"nusantarare/modul/edmtreatyin/backend/services"
	"nusantarare/modul/edmtreatyin/backend/tiruan"
)

var superUji = inti.Pelaku{AkunID: "UJI-SUPER"}

func cariLama(daftar []models.DokumenLama, id string) (models.DokumenLama, bool) {
	for _, d := range daftar {
		if d.ID == id {
			return d, true
		}
	}
	return models.DokumenLama{}, false
}

func TestCopyOldDaftar(t *testing.T) {
	g := siapkan(t)
	// baris json_polis Utility1 aplikasi baru (screenshot WO 07-10-2026: EDM Number kosong, "the old JSON cannot be
	// read"): bukan dokumen Pega lama - tidak tampil di popup
	g.dok["Z1"] = models.BarisJSONPolis{IDPega: "EDMT-990001", NoPolis: polisA, NoEndors: models.NomorEDM(polisA, 1), ProdKe: "1"}
	// saringan WO 07-10-2026: dokumen yang kasus Pega-nya tidak ada / belum berproduksi tidak tampil
	g.dok["Z2"] = edmUji["A3"].baris(t)
	g.dok["Z2"] = models.BarisJSONPolis{IDPega: models.KelasKerjaEDM + " EDMT-990099", NoPolis: "UJI-QP.LAIN", ProdKe: "1",
		NoEndors: models.NomorEDM("UJI-QP.LAIN", 1), DataJSON: g.dok["Z2"].DataJSON}
	g.tanpaProduksi = map[string]bool{"Z2": true}
	l := services.Baru(g, nil)
	daftar, err := l.DaftarDokumenLama(context.Background(), superUji, true)
	if err != nil {
		t.Fatal(err)
	}
	var ids []string
	for _, d := range daftar {
		ids = append(ids, d.ID)
	}
	// urutan generasi (NOPOLIS lalu PRODKE), bukan urutan ROWID
	if strings.Join(ids, ",") != "EDMT-990001,EDMT-990002,EDMT-990009,EDMT-990003,EDMT-990011,EDMT-990012,EDMT-990021" {
		t.Fatalf("daftar %v", ids)
	}
	boleh := map[string]bool{}
	for _, d := range daftar {
		boleh[d.ID] = d.BolehDisalin
	}
	harap := map[string]bool{"EDMT-990001": true, "EDMT-990002": true, "EDMT-990009": false, "EDMT-990003": true,
		"EDMT-990011": true, "EDMT-990012": true, "EDMT-990021": false}
	if !reflect.DeepEqual(boleh, harap) {
		t.Fatalf("boleh disalin %v, harap %v", boleh, harap)
	}
	d, _ := cariLama(daftar, "EDMT-990021")
	if len(d.Alasan) != 1 || d.Alasan[0] != models.AlasanSalinLama(models.ErrGenerasiSebelumnyaTidakAda) {
		t.Fatalf("polis C tanpa generasi NB: %v", d.Alasan)
	}
	d, _ = cariLama(daftar, "EDMT-990009")
	if len(d.Alasan) != 1 || d.Alasan[0] != models.AlasanSalinLama(models.ErrPercabangan) {
		t.Fatalf("percabangan: %v", d.Alasan)
	}
	d, _ = cariLama(daftar, "EDMT-990002")
	if d.NoPolis != polisA || d.ProdKe != 2 || d.EDMNo != models.NomorEDM(polisA, 2) || d.EDMType != "3" {
		t.Fatalf("baris popup %+v", d)
	}
	if len(g.Panggil) != 0 || len(g.Selisih) != 0 {
		t.Fatalf("daftar Copy Old menulis: %v", g.Panggil)
	}
}

func TestCopyOldProcessCopyUrutGenerasi(t *testing.T) {
	g := siapkan(t)
	l := services.Baru(g, nil)
	ctx := context.Background()
	// urutan permintaan sengaja terbalik: generasi 3 dulu
	j, err := l.SalinDokumenLama(ctx, superUji, true,
		[]string{"EDMT-990003", "EDMT-990012", "EDMT-990002", "EDMT-990011", "EDMT-990001", "EDMT-990021", "EDMT-990001", " "})
	if err != nil {
		t.Fatal(err)
	}
	status := map[string]string{}
	var urut []string
	for _, h := range j.Hasil {
		status[h.ID] = h.Status
		urut = append(urut, h.ID)
	}
	harap := map[string]string{
		"EDMT-990001": models.SalinDisalin, "EDMT-990002": models.SalinDisalin, "EDMT-990003": models.SalinDisalin,
		// B1 kehilangan baris spreading generasi NB; B2 lalu tidak punya generasi sebelumnya
		"EDMT-990011": models.SalinDitolak, "EDMT-990012": models.SalinDitolak,
		"EDMT-990021": models.SalinDitolak,
	}
	if j.Disalin != 3 || !reflect.DeepEqual(status, harap) {
		t.Fatalf("hasil %+v", j)
	}
	if strings.Join(urut, ",") != "EDMT-990001,EDMT-990002,EDMT-990003,EDMT-990011,EDMT-990012,EDMT-990021" {
		t.Fatalf("Process Copy wajib berurutan generasi: %v", urut)
	}
	for _, h := range j.Hasil {
		for _, p := range h.Pesan {
			if strings.Contains(p, "UJI-QP") {
				t.Fatalf("pesan layar memuat nomor polis: %q", p)
			}
		}
	}
	x := g.Generasi["EDMT-990002"]
	if x == nil || x.OldPolisID != "EDMT-990001" || x.ProdKe != 2 || g.Selisih["EDMT-990002"] == nil {
		t.Fatalf("generasi 2 tersalin lewat pemuat: %+v", x)
	}
	// sesudah disalin: hilang dari popup; percabangan generasi 2 tetap ditolak
	daftar, err := l.DaftarDokumenLama(ctx, superUji, true)
	if err != nil {
		t.Fatal(err)
	}
	for _, id := range []string{"EDMT-990001", "EDMT-990002", "EDMT-990003"} {
		if _, ada := cariLama(daftar, id); ada {
			t.Errorf("%s sudah tersalin tetapi masih di popup", id)
		}
	}
	if d, _ := cariLama(daftar, "EDMT-990009"); d.BolehDisalin {
		t.Error("percabangan atas generasi 1 yang sudah punya generasi 2 di tabel flat")
	}
	// salin ulang = sudah ada, nol tulisan
	j, err = l.SalinDokumenLama(ctx, superUji, true, []string{"EDMT-990002"})
	if err != nil || len(j.Hasil) != 1 || j.Hasil[0].Status != models.SalinSudahAda || j.Disalin != 0 {
		t.Fatalf("salin ulang %+v %v", j, err)
	}
	if _, err := l.SalinDokumenLama(ctx, superUji, true, []string{" "}); err == nil {
		t.Fatal("tanpa pilihan wajib ditolak")
	}
}

func TestCopyOldHanyaSuperadmin(t *testing.T) {
	l := services.Baru(siapkan(t), nil)
	ctx := context.Background()
	if _, err := l.DaftarDokumenLama(ctx, inti.Pelaku{AkunID: "UJI-ADMIN"}, false); err == nil {
		t.Fatal("bukan superadmin membuka Copy Old")
	}
	if _, err := l.SalinDokumenLama(ctx, inti.Pelaku{AkunID: "UJI-ADMIN"}, false, []string{"EDMT-990001"}); err == nil {
		t.Fatal("bukan superadmin menyalin")
	}
	if l.HakPortal(inti.Pelaku{AkunID: "UJI-ADMIN"}, false).CopyOld || !l.HakPortal(superUji, true).CopyOld {
		t.Fatal("tombol Copy Old hanya bagi superadmin")
	}
	// gudang tanpa pembaca JSON_POLIS (tiruan biasa): tombol tidak tampil walau superadmin
	if services.Baru(tiruan.Baru(), nil).HakPortal(superUji, true).CopyOld {
		t.Fatal("tanpa pembaca JSON_POLIS tombol tidak tampil")
	}
}

// minta - permintaan HTTP sesi stub `akun` dengan menu akun `menuAkun` dan menu ber-hak LIHAT `lihat`.
func minta(h http.Handler, metode, jalur, badan, akun string, lihat []string, menuAkun ...string) *httptest.ResponseRecorder {
	r := httptest.NewRequest(metode, handlers.Prefix+jalur, strings.NewReader(badan))
	r.Header.Set("X-Pelaku", akun)
	ctx := inti.DenganMenuLihat(inti.DenganAksesMenu(r.Context(), menuAkun), lihat)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r.WithContext(ctx))
	return w
}

func TestCopyOldHTTP(t *testing.T) {
	h := handlers.Router(services.Baru(siapkan(t), nil), true)
	edm := handlers.KodeMenu
	if w := minta(h, "GET", "/hak", "", "UJI-ADMIN", nil, edm); w.Code != 200 || !strings.Contains(w.Body.String(), `"copyOld":false`) {
		t.Fatalf("hak admin %d %s", w.Code, w.Body.String())
	}
	if w := minta(h, "GET", "/lama", "", "UJI-ADMIN", nil, edm); w.Code != http.StatusForbidden {
		t.Fatalf("bukan superadmin %d", w.Code)
	}
	if w := minta(h, "GET", "/hak", "", "UJI-SUPER", nil, edm, menu.KodeKelolaUser); !strings.Contains(w.Body.String(), `"copyOld":true`) {
		t.Fatalf("hak superadmin %s", w.Body.String())
	}
	// View only berlaku juga bagi superadmin (keputusan work owner 05-10-2026)
	if w := minta(h, "GET", "/lama", "", "UJI-SUPER", []string{edm}, edm, menu.KodeKelolaUser); w.Code != http.StatusForbidden {
		t.Fatalf("superadmin View only %d", w.Code)
	}
	w := minta(h, "GET", "/lama", "", "UJI-SUPER", nil, edm, menu.KodeKelolaUser)
	var daftar []models.DokumenLama
	if w.Code != 200 || json.Unmarshal(w.Body.Bytes(), &daftar) != nil || len(daftar) != 7 {
		t.Fatalf("daftar %d %s", w.Code, w.Body.String())
	}
	if w := minta(h, "POST", "/lama/salin", `{"ids":[]}`, "UJI-SUPER", nil, edm, menu.KodeKelolaUser); w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("Process Copy tanpa pilihan %d", w.Code)
	}
	w = minta(h, "POST", "/lama/salin", `{"ids":["EDMT-990001"]}`, "UJI-SUPER", nil, edm, menu.KodeKelolaUser)
	var j models.JawabanSalinLama
	if w.Code != 200 || json.Unmarshal(w.Body.Bytes(), &j) != nil || j.Disalin != 1 {
		t.Fatalf("Process Copy %d %s", w.Code, w.Body.String())
	}
}
