package handlers

import (
	"bytes"
	"context"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/outbox"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
	"nusantarare/modul/nbfacin/backend/services"
)

type tahapLampiranH map[string]models.TahapKasus

func (t tahapLampiranH) TahapKasus(_ context.Context, id string) (models.TahapKasus, error) {
	k, ada := t[id]
	if !ada {
		return k, repository.ErrKasusTidakAda
	}
	return k, nil
}

// lampiranTiruanH - satu dokumen tersimpan "20261008010405007" (OFFER) + yang diunggah lewat uji.
type lampiranTiruanH struct {
	dok   map[string]models.DokumenPolis
	objek map[string]bool
}

func (l *lampiranTiruanH) KategoriLampiran(context.Context) ([]string, error) {
	return []string{"ACCEPTANCE / APPROVAL", "OFFER"}, nil
}
func (l *lampiranTiruanH) JumlahLampiran(_ context.Context, idPega string) (map[string]int, error) {
	out := map[string]int{}
	for _, d := range l.dok {
		if repository.IDPegaProduksi(d.IDPega) == repository.IDPegaProduksi(idPega) { // dua bentuk (idpegaproduksi.go)
			out[d.Kategori2]++
		}
	}
	return out, nil
}
func (l *lampiranTiruanH) isi(d models.DokumenPolis) models.DokumenPolis {
	d.Tersimpan, d.Tanggal, d.StatusEfek = l.objek[d.TStorageID], "08-10-2026 13:04", outbox.StatusEfekAntre
	return d
}
func (l *lampiranTiruanH) DaftarLampiran(_ context.Context, idPega, kategori, _, _ string) ([]models.DokumenPolis, error) {
	var out []models.DokumenPolis
	for _, d := range l.dok {
		if repository.IDPegaProduksi(d.IDPega) == repository.IDPegaProduksi(idPega) && d.Kategori2 == kategori {
			out = append(out, l.isi(d))
		}
	}
	return out, nil
}
func (l *lampiranTiruanH) AmbilLampiran(_ context.Context, idPega, id, _, _ string) (models.DokumenPolis, error) {
	d, ada := l.dok[id]
	if !ada || repository.IDPegaProduksi(d.IDPega) != repository.IDPegaProduksi(idPega) {
		return d, repository.ErrLampiranTidakAda
	}
	return l.isi(d), nil
}
func (l *lampiranTiruanH) KunciLampiran(context.Context, *db.Tx, string) (models.DokumenPolis, error) {
	return models.DokumenPolis{}, repository.ErrLampiranTidakAda
}
func (l *lampiranTiruanH) SisipLampiran(_ context.Context, _ *db.Tx, d models.DokumenPolis, saat time.Time) (string, error) {
	d.ID = models.IDDokumenPolis(saat)
	l.dok[d.ID] = d
	return d.ID, nil
}
func (l *lampiranTiruanH) HapusLampiran(_ context.Context, _ *db.Tx, _, id string) error {
	delete(l.dok, id)
	return nil
}
func (l *lampiranTiruanH) LinkStorage(_ context.Context, _ *db.Tx, imageID string) (models.LinkStorage, bool, error) {
	return models.LinkStorage{ImageID: imageID}, l.objek[imageID], nil
}
func (l *lampiranTiruanH) SimpanObjek(context.Context, *db.Tx, models.ObjekStorage) error { return nil }
func (l *lampiranTiruanH) PerbaruiObjek(context.Context, *db.Tx, models.ObjekStorage) error {
	return nil
}
func (l *lampiranTiruanH) HapusObjek(context.Context, *db.Tx, string) error { return nil }
func (l *lampiranTiruanH) AppStorage(context.Context) (string, error)       { return "", nil }
func (l *lampiranTiruanH) TokenStorageBerlaku(context.Context, *db.Tx, string, time.Time, time.Duration) (string, time.Duration, error) {
	return "", 0, nil
}
func (l *lampiranTiruanH) PungutEfekRujukan(context.Context, *db.Tx, string, string, time.Time) (outbox.BarisEfekKeluar, error) {
	return outbox.BarisEfekKeluar{}, outbox.ErrEfekTidakAda
}

// antreanDiamH - efek diantre, tidak pernah dijalankan (status tetap "menunggu").
type antreanDiamH struct{ n int }

func (a *antreanDiamH) Antre(context.Context, *db.Tx, string, string, string, time.Time) error {
	a.n++
	return nil
}
func (a *antreanDiamH) Pungut(context.Context, *db.Tx, time.Time) (outbox.BarisEfekKeluar, error) {
	return outbox.BarisEfekKeluar{}, outbox.ErrEfekTidakAda
}
func (a *antreanDiamH) PungutRujukan(context.Context, *db.Tx, string, time.Time) (outbox.BarisEfekKeluar, error) {
	return outbox.BarisEfekKeluar{}, outbox.ErrEfekTidakAda
}
func (a *antreanDiamH) Tuntaskan(context.Context, *db.Tx, string, string, time.Time, string, time.Time) error {
	return nil
}

const imageUjiH = "ABCDEF0123456789ABCDEF0123456789"

func layananLampiranH(t *testing.T, flag string) (*services.Service, *lampiranTiruanH, *antreanDiamH) {
	t.Helper()
	folder := t.TempDir()
	if err := os.MkdirAll(filepath.Join(folder, "nbfacin", "penyimpanan"), 0o750); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(folder, "nbfacin", "penyimpanan", imageUjiH), []byte("ISI-UJI"), 0o640); err != nil {
		t.Fatal(err)
	}
	repo := &lampiranTiruanH{dok: map[string]models.DokumenPolis{"20261008010405007": {ID: "20261008010405007",
		IDPega: "ASM-FW-GISFW-WORK UJI-NB-1", NamaFile: "Slip \"Final\".pdf", Kategori2: "OFFER", TStorageID: imageUjiH}},
		objek: map[string]bool{imageUjiH: true}}
	a := &antreanDiamH{}
	svc := services.Baru(nil).DenganTransaksi(tanpaOracle).
		DenganTahap(tahapLampiranH{"UJI-NB-1": {Position: "ReasFacInMarketing", FlagOnGoingPolicy: flag}}).
		DenganLampiranUji(repo, a, services.StorageLampiranLokal(folder), folder)
	return svc, repo, a
}

func mintaMultipart(t *testing.T, svc *services.Service, nama, isi, kategori, peran string) (int, string) {
	t.Helper()
	var b bytes.Buffer
	mw := multipart.NewWriter(&b)
	if nama != "" {
		f, _ := mw.CreateFormFile("berkas", nama)
		_, _ = f.Write([]byte(isi))
	}
	_ = mw.WriteField("kategori", kategori)
	_ = mw.Close()
	mux := http.NewServeMux()
	DaftarkanRute(mux, svc, true)
	r := httptest.NewRequest("POST", "/api/nbfacin/kasus/UJI-NB-1/lampiran/berkas", &b)
	r.Header.Set("Content-Type", mw.FormDataContentType())
	r.Header.Set("X-Pelaku", "UJI-MKT")
	r.Header.Set("X-Peran", peran)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	return w.Code, strings.TrimSpace(w.Body.String())
}

// TestLampiranHandler - kontrak tiket 67 (api.ts TampilanLampiran / KategoriLampiran / BarisLampiran dua arah), 401, POST
// multipart 201 {id} / 400 / 403 (case sudah jadi polis, K67-6) / 422, isi attachment + nosniff, ulangi 409 (tersimpan), DELETE 204 / 403.
func TestLampiranHandler(t *testing.T) {
	svc, repo, a := layananLampiranH(t, "0")
	jalur := "/api/nbfacin/kasus/UJI-NB-1/lampiran"
	kode, isi := mintaPeran(t, svc, "GET", jalur, "", "UJI-MKT", "ReasFacInMarketing")
	if kode != 200 || isi != `{"kategori":[{"kategori":"ACCEPTANCE / APPROVAL","jumlah":0},{"kategori":"OFFER","jumlah":1}],`+
		`"bolehUnggah":true,"bolehHapus":true}` {
		t.Fatalf("GET: %d %s", kode, isi)
	}
	if kode, _ := mintaPeran(t, svc, "GET", jalur, "", "", ""); kode != 401 {
		t.Errorf("401: %d", kode)
	}
	if kode, _ := mintaPeran(t, svc, "GET", "/api/nbfacin/kasus/UJI-NB-9/lampiran", "", "UJI-MKT", "ReasFacInMarketing"); kode != 404 {
		t.Errorf("404: %d", kode)
	}
	kode, isi = mintaPeran(t, svc, "GET", jalur+"/berkas?kategori=OFFER", "", "UJI-UW", "ReasFacInUnderwriting")
	if kode != 200 || isi != `{"baris":[{"id":"20261008010405007","namaFile":"Slip \"Final\".pdf","catatan":"OFFER",`+
		`"tanggal":"08-10-2026 13:04","status":""}]}` {
		t.Errorf("berkas: %d %s", kode, isi)
	}
	if kode, _ := mintaPeran(t, svc, "GET", jalur+"/berkas", "", "UJI-MKT", "ReasFacInMarketing"); kode != 400 {
		t.Errorf("berkas tanpa kategori: %d", kode)
	}
	// Unggah.
	kode, isi = mintaMultipart(t, svc, "b.docx", "ISI-DOCX-UJI", "offer", "ReasFacInMarketing")
	if kode != 201 || kunciJSON(t, isi) != "id" || len(repo.dok) != 2 || a.n != 1 {
		t.Errorf("unggah: %d %s", kode, isi)
	}
	for nama, c := range map[string]struct {
		file, kategori, peran string
		kode                  int
		pesan                 string
	}{
		"jenis":      {"b.exe", "OFFER", "ReasFacInMarketing", 400, "Jenis berkas tidak didukung: .exe"},
		"tanpa":      {"", "OFFER", "ReasFacInMarketing", 400, "Pilih berkas"},
		"kategori":   {"b.pdf", "LAIN", "ReasFacInMarketing", 422, ""},
		"tanpa kat.": {"b.pdf", "", "ReasFacInMarketing", 400, "Kategori wajib dipilih."},
	} {
		if kode, isi := mintaMultipart(t, svc, c.file, "ISI-PDF-UJI", c.kategori, c.peran); kode != c.kode || !strings.Contains(isi, c.pesan) {
			t.Errorf("unggah %s: %d %s", nama, kode, isi)
		}
	}
	if kode, _ := mintaPeran(t, svc, "POST", jalur+"/berkas", "bukan multipart", "UJI-MKT", "ReasFacInMarketing"); kode != 400 {
		t.Errorf("bukan multipart: %d", kode)
	}
	// Isi.
	mux := http.NewServeMux()
	DaftarkanRute(mux, svc, true)
	r := httptest.NewRequest("GET", jalur+"/berkas/20261008010405007/isi", nil)
	r.Header.Set("X-Pelaku", "UJI-UW")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	if w.Code != 200 || w.Body.String() != "ISI-UJI" || w.Header().Get("Content-Type") != "application/pdf" ||
		!strings.HasPrefix(w.Header().Get("Content-Disposition"), "attachment;") || w.Header().Get("X-Content-Type-Options") != "nosniff" {
		t.Errorf("isi: %d %q %v", w.Code, w.Body.String(), w.Header())
	}
	if kode, _ := mintaPeran(t, svc, "GET", jalur+"/berkas/9999/isi", "", "UJI-UW", ""); kode != 404 {
		t.Errorf("isi 404: %d", kode)
	}
	// Ulangi tersimpan 409; hapus 204.
	if kode, _ := mintaPeran(t, svc, "POST", jalur+"/berkas/20261008010405007/ulangi", "", "UJI-MKT", "ReasFacInMarketing"); kode != 409 {
		t.Errorf("ulangi: %d", kode)
	}
	if kode, isi := mintaPeran(t, svc, "DELETE", jalur+"/berkas/20261008010405007", "", "UJI-UW", ""); kode != 204 || isi != "" ||
		len(repo.dok) != 1 {
		t.Errorf("hapus: %d %s", kode, isi)
	}
	// Flag 1: hapus 403 bukan IT.
	svc, _, _ = layananLampiranH(t, "1")
	if kode, _ := mintaPeran(t, svc, "DELETE", jalur+"/berkas/20261008010405007", "", "UJI-MKT", "ReasFacInMarketing"); kode != 403 {
		t.Errorf("hapus flag 1: %d", kode)
	}
	// K67-6: unggah oleh bukan pemegang workbasket boleh (201); case sudah jadi polis -> 403.
	svc, _, _ = layananLampiranH(t, "0")
	if kode, isi := mintaMultipart(t, svc, "b.pdf", "ISI-PDF-UJI", "OFFER", "ReasFacInUnderwriting"); kode != 201 {
		t.Errorf("unggah bukan pemegang: %d %s", kode, isi)
	}
	svc.DenganTahap(tahapLampiranH{"UJI-NB-1": {Position: "ReasFacInMarketing", FlagOnGoingPolicy: "1", StatusWork: "Resolved-Completed"}})
	if kode, isi := mintaMultipart(t, svc, "b.pdf", "ISI-PDF-UJI", "OFFER", "ReasFacInMarketing"); kode != 403 {
		t.Errorf("unggah polis: %d %s", kode, isi)
	}
	// Kunci api.ts dua arah.
	for antarmuka, gob := range map[string][]string{"TampilanLampiran": kunciTagJSON(tampilanLampiranKabel{}),
		"KategoriLampiran": kunciTagJSON(kategoriLampiranKabel{}), "BarisLampiran": kunciTagJSON(barisLampiranKabel{})} {
		if ts := kunciInterfaceTS(t, antarmuka); !reflect.DeepEqual(ts, gob) {
			t.Errorf("%s: api.ts %v, handler %v", antarmuka, ts, gob)
		}
	}
}

// TestHasilSubmitKunciTS - tiket 64b: kunci HasilSubmit (pesan, position, nbStatus, pesanBinding) dua arah dengan api.ts;
// `galat` hanya di jawaban 422 (bukan medan HasilSubmit api.ts).
func TestHasilSubmitKunciTS(t *testing.T) {
	gob := kunciTagJSON(hasilSubmitKabel{})
	if len(gob) == 0 || gob[0] != "galat" {
		t.Fatalf("instrumen %v", gob)
	}
	if ts := kunciInterfaceTS(t, "HasilSubmit"); !reflect.DeepEqual(ts, gob[1:]) {
		t.Errorf("api.ts %v, handler %v", ts, gob[1:])
	}
}

// TestKunciPeriodeHandler - tiket 64d K1: GET kasus tahap Binding "kunciPeriode":true; PUT general yang mengubah Begin date
// -> 409; tanggal sama (kosong = tersimpan kosong) -> 200.
func TestKunciPeriodeHandler(t *testing.T) {
	svc := layananKasus().DenganTahap(tahapTiruanH{models.TahapKasus{Position: "ReasFacInMarketing", FlagOnGoingPolicy: "2"}})
	jalur := "/api/nbfacin/kasus/UJI-NB-1"
	if kode, isi := mintaPeran(t, svc, "GET", jalur, "", "UJI-MKT", "ReasFacInMarketing"); kode != 200 || !strings.Contains(isi, `"kunciPeriode":true,`) {
		t.Errorf("GET: %d %s", kode, isi)
	}
	if kode, isi := mintaPeran(t, svc, "PUT", jalur+"/general", `{"beginDate":"01-11-2026","offeringDate":"01-10-2026"}`, "UJI-MKT",
		"ReasFacInMarketing"); kode != 409 || !strings.Contains(isi, "Begin Date dan End Date tidak dapat diubah di tahap Binding") {
		t.Errorf("PUT ubah: %d %s", kode, isi)
	}
	if kode, isi := mintaPeran(t, svc, "PUT", jalur+"/general", `{"qqName":"UJI QQ","offeringDate":"01-10-2026"}`, "UJI-MKT",
		"ReasFacInMarketing"); kode != 200 {
		t.Errorf("PUT sama: %d %s", kode, isi)
	}
}

// TestWajibTypeCedantHandler - tiket 64d K2: GET cedant tahap Binding "wajibType":true walau gerbang ProtectShareCedant
// ("wajib") tidak berlaku; PUT tipe kosong -> 400 pesan persis.
func TestWajibTypeCedantHandler(t *testing.T) {
	c := &cedantTiruanH{ada: models.KasusCedant{SobName: "UJI SOB", Cedant: []models.BarisCedant{}}}
	sp := kasusSpreadingUji()
	sp.PercentShare = angkaUji("50")
	svc := services.Baru(nil).DenganCedant(c).DenganSpreading(&spreadingTiruanH{tersimpan: sp}).DenganTransaksi(tanpaOracleKlausa).
		DenganTahap(tahapTiruanH{models.TahapKasus{Position: "ReasFacInMarketing", FlagOnGoingPolicy: "2"}})
	kode, isi := mintaPeran(t, svc, "GET", "/api/nbfacin/kasus/S1/cedant", "", "UJI-MKT", "ReasFacInMarketing")
	if kode != 200 || !strings.Contains(isi, `"wajib":false,"wajibType":true`) {
		t.Errorf("GET: %d %s", kode, isi)
	}
	kode, isi = mintaPeran(t, svc, "PUT", "/api/nbfacin/kasus/S1/cedant", `{"shareCedantType":"","cedant":[]}`, "UJI-MKT", "ReasFacInMarketing")
	if kode != 400 || !strings.Contains(isi, "Share Cedant Type wajib diisi.") {
		t.Errorf("PUT: %d %s", kode, isi)
	}
}

// TestGalatPolicy409 - tiket 69a: selesai Policy / FacRetro -> 409 dengan pesan apa adanya.
func TestGalatPolicy409(t *testing.T) {
	for err, pesan := range map[error]string{
		services.ErrBusinessOldIDTidakAda: "BusinessOldId case tidak dapat ditentukan dari tabel BUSINESS (Class of Business / Group Business); nomor polis tidak dapat dibentuk.",
		services.ErrFacRetroBelumTersedia: "Case FacRetro belum dapat diselesaikan di sistem baru (tiket 65).",
		services.ErrIDKapalFacIn:          "ID kapal cargo lebih dari 20 karakter - tidak muat kolom IDSHIP FACINOFFER / FACINPRODUCTION",
	} {
		w := httptest.NewRecorder()
		tulisGalatSubmit(w, err)
		if w.Code != 409 || !strings.Contains(w.Body.String(), `"`+pesan+`"`) {
			t.Errorf("%v: %d %s", err, w.Code, w.Body.String())
		}
	}
}
