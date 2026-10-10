package handlers_test

// Uji seam HTTP temuan review 10-10-2026: Submit pop-up memeriksa tombol PEMBUKA-nya (modal dievaluasi untuk setiap
// baris), dan penghapusan adjustment yang sudah bertaut komite = 409, bukan galat server.

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"nusantarare/inti/backend/penyimpanan"

	"nusantarare/modul/claimfacin/backend/handlers"
	"nusantarare/modul/claimfacin/backend/models"
	"nusantarare/modul/claimfacin/backend/services"
)

func TestGeneratePLAMenuntutPrintPLATerbuka(t *testing.T) {
	u := baruUji(t)
	id := u.sampaiEstimasi()
	h := u.g.Halaman(id)
	ob := h.AmbilDaftar(models.DaftarObjek)[0]
	ob["IsFacretro"], ob["PlaStatus"] = "1", "1" // PLA sudah dicetak: Print PLA nonaktif
	u.g.SetelHalaman(id, h)
	kode, out := u.aksiP(id, admin, "", services.PermintaanAksi{Aksi: "GeneratePLA", Konteks: "pla:1",
		Masukan: map[string]string{models.JalurObjek(1) + ".RemarksPLA": "UJI"}})
	u.wajib(kode, http.StatusConflict, out, "GeneratePLA atas objek yang Print PLA-nya nonaktif")
	galatMemuat(t, out, "Print PLA objek 1 tertutup")
}

func TestKirimKomiteMenuntutSendToCommitteTerbuka(t *testing.T) {
	u := baruUji(t)
	u.a.Roster = []models.AnggotaKomite{{ID: "1", OperatorID: "UJI-K1", Jabatan: "UJI SPV", Degree: "1",
		LimitBottom: "-9999999999999", LimitTop: "57750000"}}
	id := u.adjustmentFinal()
	u.lampirkan(id, "LOD", "DLA", "SPGR")
	h := u.g.Halaman(id)
	b := h.AmbilDaftar(models.DaftarAdj(1, 1))[0]
	b["NameOfBank"], b["NoAccount"], b["IDOfBank"] = "UJI BANK", "0001", "UJI-B1"
	h.SetelDaftar(models.DaftarDiAdj(1, 1, 1, models.AnakAdjSpread), nil) // spreading kosong: Send to Committe nonaktif
	u.g.SetelHalaman(id, h)
	modal := models.KunciPanel(models.ModalKomite, models.DaftarAdj(1, 1), 1)
	kode, out := u.aksiT(id, services.PermintaanAksi{Aksi: "KirimKomite", Konteks: modal,
		Masukan: map[string]string{jAdj(1, "DataCommitteFacin.Remarks"): "UJI"}})
	u.wajib(kode, http.StatusConflict, out, "KirimKomite tanpa Spreading Adjustment")
	galatMemuat(t, out, "Send to Committe adjustment 1 tertutup")
	if len(u.g.Komite) != 0 {
		t.Fatalf("kasus komite lahir: %v", u.g.Komite)
	}
}

func TestHapusAdjustmentBerkomiteKonflik(t *testing.T) {
	u := baruUji(t)
	id := u.adjustmentFinal()
	h := u.g.Halaman(id)
	b := h.AmbilDaftar(models.DaftarAdj(1, 1))[0]
	b[models.PropKomiteID], b["IsKomite"] = "KMT-UJI1", "" // tertaut komite, tombol hapus tampil
	u.g.SetelHalaman(id, h)
	kode, out := u.aksiT(id, services.PermintaanAksi{Aksi: "DisableSendComite",
		Konteks: models.KunciPanel(models.PanelItemAdj, models.DaftarItem(1), 1), Indeks: 1})
	u.wajib(kode, http.StatusConflict, out, "hapus adjustment berkomite")
	galatMemuat(t, out, "sedang di komite")
}

func galatMemuat(t *testing.T, out map[string]any, teks string) {
	t.Helper()
	if g, _ := out["galat"].(string); !strings.Contains(g, teks) {
		t.Fatalf("galat %q, mau memuat %q", g, teks)
	}
}

// Hapus lampiran: baris DOCUMENT_CLAIM, objek storage, dan catatan T_STORAGE_IMAGE dalam satu transaksi - storage menolak
// = baris tetap utuh (temuan review: dahulu objek dihapus sebelum transaksi).
func TestHapusLampiranSatuTransaksi(t *testing.T) {
	for _, gagal := range []bool{true, false} {
		u := baruUji(t)
		id := u.buat()
		u.g.Dokumen = append(u.g.Dokumen, models.BarisDokumenKlaim{ID: "UJI-DOK-1", IDPega: models.KunciInstans(id),
			NamaFile: "UJI.pdf", MIME: "pdf", Kategori1: "LOD", StorageID: "UJI-IMG-1",
			Tanggal: time.Date(2026, 3, 5, 9, 0, 0, 0, models.Jakarta)})
		u.g.Storage = append(u.g.Storage, penyimpanan.Objek{ImageID: "UJI-IMG-1"})
		if gagal {
			u.b.GagalHapus = errors.New("UJI storage menolak")
		}
		kode, out := u.minta(http.MethodPost, handlers.Prefix+"/kasus/"+id+"/lampiran/UJI-DOK-1/hapus", admin, "", nil)
		if gagal {
			if kode == http.StatusOK || len(u.g.Dokumen) != 1 || len(u.g.Storage) != 1 {
				t.Fatalf("storage gagal: HTTP %d %v, dokumen %d, storage %d", kode, out, len(u.g.Dokumen), len(u.g.Storage))
			}
			continue
		}
		if kode != http.StatusOK && kode != http.StatusNoContent || len(u.g.Dokumen) != 0 || len(u.g.Storage) != 0 ||
			len(u.b.Dihapus) != 1 {
			t.Fatalf("hapus: HTTP %d %v, dokumen %d, storage %d, dihapus %v", kode, out, len(u.g.Dokumen),
				len(u.g.Storage), u.b.Dihapus)
		}
	}
}

// Pembeda urutan: hapus DOCUMENT_CLAIM gagal = objek storage TIDAK disentuh (urutan lama menghapus objek lebih dulu).
func TestHapusLampiranDBGagalStorageUtuh(t *testing.T) {
	u := baruUji(t)
	id := u.buat()
	u.g.Dokumen = append(u.g.Dokumen, models.BarisDokumenKlaim{ID: "UJI-DOK-1", IDPega: models.KunciInstans(id),
		NamaFile: "UJI.pdf", MIME: "pdf", Kategori1: "LOD", StorageID: "UJI-IMG-1",
		Tanggal: time.Date(2026, 3, 5, 9, 0, 0, 0, models.Jakarta)})
	u.g.Storage = append(u.g.Storage, penyimpanan.Objek{ImageID: "UJI-IMG-1"})
	u.g.GagalHapusDok = errors.New("UJI baris terkunci")
	kode, out := u.minta(http.MethodPost, handlers.Prefix+"/kasus/"+id+"/lampiran/UJI-DOK-1/hapus", admin, "", nil)
	if kode == http.StatusOK || len(u.b.Dihapus) != 0 || len(u.g.Storage) != 1 {
		t.Fatalf("HTTP %d %v, objek dihapus %v, storage %d", kode, out, u.b.Dihapus, len(u.g.Storage))
	}
}

// Worklist pembuat tak peka huruf, sama dengan Pemegang (strings.EqualFold): kasus tetap tampil bila akun masuk dengan
// huruf berbeda (temuan review 10-10-2026).
func TestWorklistPembuatTakPekaHuruf(t *testing.T) {
	u := baruUji(t)
	id := u.buat()
	kode, out := u.mintaDaftar("/kasus?daftar=saya", strings.ToLower(admin))
	if kode != http.StatusOK || !strings.Contains(out, id) {
		t.Fatalf("HTTP %d: %s", kode, out)
	}
}

// mintaDaftar - GET daftar (jawaban JSON larik) sebagai teks.
func (u *uji) mintaDaftar(jalur, pelaku string) (int, string) {
	u.t.Helper()
	r := httptest.NewRequest(http.MethodGet, handlers.Prefix+jalur, nil)
	r.Header.Set("X-Pelaku", pelaku)
	w := httptest.NewRecorder()
	u.srv.ServeHTTP(w, r)
	return w.Code, w.Body.String()
}

// Layar 422 = halaman sebelum penangan (+ pesan): hapus baris Spreading Adjustment yang membuat total < 100% dibatalkan,
// dan layar validasinya tetap memuat baris itu (temuan review: dahulu grid 2 baris dipadu halaman lokal 3 baris).
func TestLayarValidasiSebelumPenangan(t *testing.T) {
	u := baruUji(t)
	id := u.adjustmentFinal()
	h := u.g.Halaman(id)
	b := h.AmbilDaftar(models.DaftarAdj(1, 1))[0]
	b["ExGratia"] = "1"
	sp := models.DaftarDiAdj(1, 1, 1, models.AnakAdjSpread)
	h.SetelDaftar(sp, []models.Baris{{"TreatyType": "10003", "SharePercentage": "50"},
		{"TreatyType": "10015", "SharePercentage": "30"}, {"TreatyType": "10003", "SharePercentage": "20"}})
	u.g.SetelHalaman(id, h)
	kode, out := u.aksiT(id, services.PermintaanAksi{Aksi: "HapusSpreadAdj", Konteks: panelAdj(1), Indeks: 3})
	u.wajib(kode, http.StatusUnprocessableEntity, out, "hapus spreading yang membuat total < 100")
	ly, _ := out["layar"].(map[string]any)
	hal, _ := ly["halaman"].(map[string]any)
	daftar, _ := hal["daftar"].(map[string]any)
	if rows, _ := daftar[sp].([]any); len(rows) != 3 {
		t.Fatalf("layar validasi memuat %d baris spreading, mau 3 (keadaan tersimpan)", len(rows))
	}
	if n := len(u.g.Halaman(id).AmbilDaftar(sp)); n != 3 {
		t.Fatalf("tersimpan %d baris", n)
	}
}
