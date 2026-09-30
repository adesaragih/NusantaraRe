//go:build db

// Seam HTTP lampiran tahun treaty (tiket 12) terhadap skema uji Oracle NYATA.
//
// ⛔ Penyimpanan berkasnya STUB LOKAL di bawah UNGGAHAN_DIR uji - satu-
// satunya yang dipalsukan di seluruh konteks ini. Rekam, outbox bersama
// `T_LOG_SERVICE_RNM`, master kategori, dan jejak adalah Oracle sungguhan.
package handlers_test

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/treatycontractout/backend/models"
	"nusantarare/modul/treatycontractout/backend/repository"
	"nusantarare/modul/treatycontractout/backend/services"
)

// mintaMultipart mengirim satu berkas + kategori.
func (u *ujiTCO) mintaMultipart(t *testing.T, jalur, nama, isi, kategori string, beridentitas bool) (int, string) {
	t.Helper()
	var badan bytes.Buffer
	mp := multipart.NewWriter(&badan)
	if nama != "" {
		f, err := mp.CreateFormFile("berkas", nama)
		if err != nil {
			t.Fatal(err)
		}
		_, _ = f.Write([]byte(isi))
	}
	_ = mp.WriteField("kategori", kategori)
	_ = mp.Close()
	r, err := http.NewRequest(http.MethodPost, u.srv.URL+jalur, &badan)
	if err != nil {
		t.Fatal(err)
	}
	r.Header.Set("Content-Type", mp.FormDataContentType())
	if beridentitas {
		r.Header.Set("X-Pelaku", "UJI-ADMIN")
	}
	res, err := http.DefaultClient.Do(r)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = res.Body.Close() }()
	b, _ := io.ReadAll(res.Body)
	return res.StatusCode, string(b)
}

// ambilMentah mengembalikan badan dan header mentah (unduhan).
func (u *ujiTCO) ambilMentah(t *testing.T, jalur string) (int, http.Header, []byte) {
	t.Helper()
	r, _ := http.NewRequest(http.MethodGet, u.srv.URL+jalur, nil)
	r.Header.Set("X-Pelaku", "UJI-ADMIN")
	res, err := http.DefaultClient.Do(r)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = res.Body.Close() }()
	b, _ := io.ReadAll(res.Body)
	return res.StatusCode, res.Header, b
}

type lampiranJSON struct {
	ID, IDTreatyYear, FileName, Category, Status, Galat, UserID string
}

// AC 54-61 lewat HTTP: unggah -> terkirim, daftar, unduh satu, unduh semua,
// keselarasan (rekam tanpa berkas terdeteksi), hapus.
func TestLampiranTahunTreatyLingkaranPenuh(t *testing.T) {
	u, bersihkan := serverTCO(t)
	defer bersihkan()
	u.isiKategori([]string{"CLAUSES", "R/I SLIP", "OTHERS"})

	kode, badan := u.minta(t, http.MethodPost, "/api/treaty-contract-out/tahun", badanTahun("2026-01-01", "2026-12-31"), true)
	if kode != http.StatusOK {
		t.Fatalf("POST tahun: %d %s", kode, badan)
	}
	var tahun tahunJSON
	_ = json.Unmarshal([]byte(badan), &tahun)
	dasar := "/api/treaty-contract-out/tahun/" + tahun.ID + "/lampiran"

	kode, badan = u.minta(t, http.MethodGet, "/api/treaty-contract-out/kategori-lampiran", nil, true)
	if kode != http.StatusOK || !strings.Contains(badan, `"total":3`) {
		t.Fatalf("kategori: %d %s", kode, badan)
	}

	// Gerbang: tanpa identitas 401, kategori asing 422, tanpa berkas 400 b376,
	// tahun lain 404 - dan tidak satu pun menulis.
	if k, b := u.mintaMultipart(t, dasar, "UJI-slip.pdf", "ISI", "CLAUSES", false); k != http.StatusUnauthorized {
		t.Errorf("tanpa identitas: %d %s", k, b)
	}
	if k, b := u.mintaMultipart(t, dasar, "UJI-slip.pdf", "ISI", "KARANGAN", true); k != http.StatusUnprocessableEntity {
		t.Errorf("kategori asing: %d %s", k, b)
	}
	if k, b := u.mintaMultipart(t, dasar, "", "", "CLAUSES", true); k != http.StatusBadRequest ||
		!strings.Contains(b, "No file attached") {
		t.Errorf("tanpa berkas: %d %s", k, b)
	}
	if k, b := u.mintaMultipart(t, "/api/treaty-contract-out/tahun/1999999/lampiran", "UJI-slip.pdf", "ISI",
		"CLAUSES", true); k != http.StatusNotFound {
		t.Errorf("tahun tidak ada: %d %s", k, b)
	}

	kode, badan = u.mintaMultipart(t, dasar, "UJI-kontrak.pdf", "ISI-UJI-LAMPIRAN", "r/i slip", true)
	if kode != http.StatusCreated {
		t.Fatalf("unggah: %d %s", kode, badan)
	}
	var hasil struct {
		Lampiran   lampiranJSON
		Peringatan string
	}
	_ = json.Unmarshal([]byte(badan), &hasil)
	l := hasil.Lampiran
	// tco4: ID = stempel YYYYMMDDHH24MISSFF3 (17 digit) di M_ATTACHMENTTREATY_2.
	if l.Status != "terkirim" || l.Category != "R/I SLIP" || l.IDTreatyYear != tahun.ID || l.UserID != "UJI-ADMIN" ||
		len(l.ID) != 17 || hasil.Peringatan != "" {
		t.Fatalf("hasil unggah: %+v", hasil)
	}
	// Baris warisan: TREATYID = TreatyYear + TreatyYearID, CATEGORY "File",
	// kategori pilihan di CATEGORY_ID, objek tercatat di T_STORAGE_IMAGE.
	var treatyID, kategoriPega, kategoriID, storageID string
	if err := u.sqlDBMentah().QueryRowContext(u.ctx, `SELECT TREATYID, CATEGORY, CATEGORY_ID, T_STORAGE_ID FROM `+
		u.skema+`.M_ATTACHMENTTREATY_2 WHERE ID = :1`, l.ID).Scan(&treatyID, &kategoriPega, &kategoriID, &storageID); err != nil {
		t.Fatal(err)
	}
	if treatyID != "2026"+tahun.ID || kategoriPega != "File" || kategoriID != "R/I SLIP" {
		t.Errorf("baris warisan: TREATYID %q CATEGORY %q CATEGORY_ID %q", treatyID, kategoriPega, kategoriID)
	}
	var objek int
	if err := u.sqlDBMentah().QueryRowContext(u.ctx, `SELECT COUNT(*) FROM `+u.skema+`.T_STORAGE_IMAGE
		 WHERE IMAGEID = :1 AND STORAGE = 'standard'`, storageID).Scan(&objek); err != nil || objek != 1 {
		t.Errorf("objek T_STORAGE_IMAGE: %d %v", objek, err)
	}
	// OQ-TCO-26 (lanjutan 4): `Update_T_Storage_SQL` terhadap Oracle - kedua
	// bentuk To_date (`DD/MM/YYYY` exp, `MM/DD/YYYY` DateTime) diterima.
	mo := repository.NewMasterLampiranTCO(u.db)
	if err := services.New(u.db).DalamTransaksi(u.ctx, func(tx *db.Tx) error {
		return mo.PerbaruiObjek(u.ctx, tx, models.ObjekPenyimpananTCO{ImageID: storageID, URLPublic: "UJI-URL-SEGAR",
			AppFolder: "UJI-FOLDER", Exp: "29/09/2026 10:00:00", TanggalUpload: "09/29/2026 09:00:00"})
	}); err != nil {
		t.Fatal(err)
	}
	var urlSegar, folder, exp, unggah string
	if err := u.sqlDBMentah().QueryRowContext(u.ctx, `SELECT URLPUBLIC, APPFOLDER, TO_CHAR(EXPDATE, 'YYYYMMDDHH24MISS'),
		 TO_CHAR(TANGGAL_UPLOAD, 'YYYYMMDDHH24MISS') FROM `+u.skema+`.T_STORAGE_IMAGE WHERE IMAGEID = :1`, storageID).
		Scan(&urlSegar, &folder, &exp, &unggah); err != nil {
		t.Fatal(err)
	}
	if urlSegar != "UJI-URL-SEGAR" || folder != "UJI-FOLDER" || exp != "20260929100000" || unggah != "20260929090000" {
		t.Errorf("Update_T_Storage_SQL: %q %q %q %q", urlSegar, folder, exp, unggah)
	}

	kode, badan = u.minta(t, http.MethodGet, dasar, nil, true)
	if kode != http.StatusOK || !strings.Contains(badan, `"total":1`) || !strings.Contains(badan, `"status":"terkirim"`) {
		t.Errorf("daftar: %d %s", kode, badan)
	}

	kode, kepala, isi := u.ambilMentah(t, dasar+"/"+l.ID+"/isi")
	if kode != http.StatusOK || string(isi) != "ISI-UJI-LAMPIRAN" ||
		!strings.HasPrefix(kepala.Get("Content-Disposition"), "attachment") {
		t.Errorf("unduh: %d %v %q", kode, kepala, isi)
	}
	kode, _, isi = u.ambilMentah(t, "/api/treaty-contract-out/tahun/1999999/lampiran/"+l.ID+"/isi")
	if kode != http.StatusNotFound {
		t.Errorf("unduh lewat tahun lain: %d %s", kode, isi)
	}

	kode, _, isi = u.ambilMentah(t, dasar+"/semua")
	if kode != http.StatusOK {
		t.Fatalf("unduh semua: %d %s", kode, isi)
	}
	arsip, err := zip.NewReader(bytes.NewReader(isi), int64(len(isi)))
	if err != nil || len(arsip.File) != 1 || arsip.File[0].Name != l.ID+"_UJI-kontrak.pdf" {
		t.Errorf("zip: %v %v", arsip, err)
	}

	kode, badan = u.minta(t, http.MethodGet, dasar+"/selaras", nil, true)
	if kode != http.StatusOK || !strings.Contains(badan, `"total":0`) {
		t.Errorf("selaras: %d %s", kode, badan)
	}

	// AC 61: berkas lenyap dari penyimpanan -> terdeteksi, unduh 409.
	simpan, _ := filepath.Glob(filepath.Join(u.unggahan, "treaty-contract-out", "penyimpanan", "*"))
	if len(simpan) != 1 {
		t.Fatalf("penyimpanan lokal berisi %v, mau 1 berkas", simpan)
	}
	_ = os.Remove(simpan[0])
	kode, badan = u.minta(t, http.MethodGet, dasar+"/selaras", nil, true)
	if kode != http.StatusOK || !strings.Contains(badan, `"total":1`) || !strings.Contains(badan, "record without a file") {
		t.Errorf("selaras sesudah berkas hilang: %d %s", kode, badan)
	}
	if kode, _, isi = u.ambilMentah(t, dasar+"/"+l.ID+"/isi"); kode != http.StatusConflict {
		t.Errorf("unduh rekam tanpa berkas: %d %s", kode, isi)
	}

	// Hapus tetap berhasil meski berkasnya sudah tidak ada.
	kode, badan = u.minta(t, http.MethodDelete, dasar+"/"+l.ID, nil, true)
	if kode != http.StatusOK || strings.Contains(badan, "could not") {
		t.Fatalf("hapus: %d %s", kode, badan)
	}
	kode, badan = u.minta(t, http.MethodGet, dasar, nil, true)
	if kode != http.StatusOK || !strings.Contains(badan, `"total":0`) {
		t.Errorf("daftar sesudah hapus: %d %s", kode, badan)
	}
}
