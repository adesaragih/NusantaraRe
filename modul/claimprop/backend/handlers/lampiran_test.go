package handlers_test

// Lampiran klaim (GCNMSaveAttachments -> InsertDocument_Act, master T_KATEGORI_DOC_KLAIM) dan gerbang "Send Claim to
// Committee" yang membaca cacahnya (AttachmentProtect_ACT S3 `AttachCategory.pxResults`).

import (
	"bytes"
	"encoding/json"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"nusantarare/inti/backend/penyimpanan"
	"nusantarare/modul/claimprop/backend/handlers"
	"nusantarare/modul/claimprop/backend/models"
	"nusantarare/modul/claimprop/backend/services"
	"nusantarare/modul/claimprop/backend/tiruan"
)

var _ services.PenyimpananBerkas = (*tiruan.Berkas)(nil)

// baruUjiLampiran - uji seam HTTP dengan penyimpanan tiruan dan master kategori PROP (kategoriifile.xls).
func baruUjiLampiran(t *testing.T) (*uji, *tiruan.Berkas) {
	g, a := tiruan.Baru(), acuanUji()
	for _, k := range []string{"ADU", "DLA", "File", "Invoice", "LOD", "SPGR", "Salvage"} {
		g.KategoriDok[k] = k
	}
	b := g.BerkasBaru()
	jam := func() time.Time { return time.Date(2026, 5, 20, 10, 0, 0, 0, models.Jakarta) }
	l := services.Baru(g, a, jam, false).DenganPenyimpanan(b)
	return &uji{t: t, srv: handlers.Router(l, true), g: g, a: a}, b
}

// unggah - POST multipart `kategori` + `berkas` (nama -> isi).
func (u *uji) unggah(id, pelaku, peran, kategori string, berkas map[string]string) (int, map[string]any) {
	u.t.Helper()
	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)
	_ = mw.WriteField("kategori", kategori)
	for nama, isi := range berkas {
		f, err := mw.CreateFormFile("berkas", nama)
		if err != nil {
			u.t.Fatal(err)
		}
		_, _ = f.Write([]byte(isi))
	}
	_ = mw.Close()
	r := httptest.NewRequest(http.MethodPost, handlers.Prefix+"/kasus/"+id+"/lampiran", &buf)
	r.Header.Set("Content-Type", mw.FormDataContentType())
	r.Header.Set("X-Pelaku", pelaku)
	if peran != "" {
		r.Header.Set("X-Peran", peran)
	}
	w := httptest.NewRecorder()
	u.srv.ServeHTTP(w, r)
	var out map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &out)
	return w.Code, out
}

func (u *uji) cacah(id, pelaku string) map[string]int {
	u.t.Helper()
	kode, out := u.minta(http.MethodGet, handlers.Prefix+"/kasus/"+id+"/lampiran", pelaku, "", nil)
	u.wajib(kode, http.StatusOK, out, "daftar lampiran")
	c := map[string]int{}
	for _, k := range out["kategori"].([]any) {
		m := k.(map[string]any)
		c[m["id"].(string)] = int(m["countAttach"].(float64))
	}
	return c
}

func TestLampiranUnggahDanCacah(t *testing.T) {
	u, b := baruUjiLampiran(t)
	id := u.buat()
	kode, out := u.unggah(id, admin, "", "LOD", map[string]string{"UJI-LOD.PDF": "%PDF-UJI"})
	u.wajib(kode, http.StatusOK, out, "unggah LOD oleh pemegang (pembuat, Outstanding)")
	if c := u.cacah(id, admin); c["LOD"] != 1 || c["DLA"] != 0 {
		t.Fatalf("AttachCategory: %v", c)
	}
	m := b.Unggahan[0]
	if m.Folder != "Claim" || m.Durasi != 1800 || m.Ext != "pdf" || m.NamaFile != "UJI-LOD.PDF" || m.Pengguna != admin {
		t.Fatalf("InsertGoogleStorage_Act: %+v", m)
	}
	d := u.g.Dokumen[0]
	if len(d.ID) != 17 || d.IDPega != id || d.Kategori1 != "LOD" || d.MIME != "pdf" || d.StorageID != "UJI-IMG-1" ||
		d.Operator != admin || len(u.g.Storage) != 1 {
		t.Fatalf("InsertDocument_Act S3-S5: %+v", d)
	}
	// bukan pemegang 403; kategori di luar master 422; tanpa berkas 422; nol baris baru
	kode, out = u.unggah(id, lain, "", "DLA", map[string]string{"UJI.pdf": "x"})
	u.wajib(kode, http.StatusForbidden, out, "bukan pemegang")
	kode, out = u.unggah(id, admin, "", "UJI-BUKAN", map[string]string{"UJI.pdf": "x"})
	u.wajib(kode, http.StatusUnprocessableEntity, out, "kategori di luar master")
	kode, out = u.unggah(id, admin, "", "DLA", nil)
	u.wajib(kode, http.StatusUnprocessableEntity, out, "tanpa berkas")
	b.Gagal = penyimpanan.ErrStorageGagal
	kode, out = u.unggah(id, admin, "", "DLA", map[string]string{"UJI.pdf": "x"})
	u.wajib(kode, http.StatusBadGateway, out, "layanan penyimpanan gagal")
	b.Gagal = errors.Join(penyimpanan.ErrBerkasDitolak)
	kode, out = u.unggah(id, admin, "", "DLA", map[string]string{"UJI.exe": "x"})
	u.wajib(kode, http.StatusUnprocessableEntity, out, "jenis berkas ditolak (GetMimeType)")
	if len(u.g.Dokumen) != 1 || len(u.g.Storage) != 1 {
		t.Fatalf("galat tidak boleh menulis baris: %d dokumen, %d objek", len(u.g.Dokumen), len(u.g.Storage))
	}
}

// Send Claim to Committee: tanpa lampiran ditolak; sesudah LOD + DLA + SPGR (adjustment Type 1 di bawah batas Direktur
// Utama) kasus komite TKMT- lahir dengan tangga roster EMAILKOMITE PROP.
func TestKirimKomiteSesudahLampiranLengkap(t *testing.T) {
	u, _ := baruUjiLampiran(t)
	id := sampaiAdjustment(u)
	u.a.Roster = []tiruan.AnggotaRoster{
		{AnggotaKomite: models.AnggotaKomite{ID: "1", OperatorID: "UJI-K1", Jabatan: "UJI-JABATAN-1", Degree: "1"},
			Batas: "-1", Sts: models.STSKlaimProp},
		{AnggotaKomite: models.AnggotaKomite{ID: "2", OperatorID: "UJI-K2", Jabatan: "UJI-JABATAN-2", Degree: "2"},
			Batas: "999999999999", Sts: models.STSKlaimProp},
	}
	isian := map[string]string{models.CD + "Payable": "2",
		models.JalurAnak(models.DaftarAdjustment, 1, "DataCommitteeTreaty.CircumCauseOfLoss"): "UJI",
		models.JalurAnak(models.DaftarAdjustment, 1, "DataCommitteeTreaty.Remarks"):           "UJI",
		models.JalurAnak(models.DaftarAdjustment, 1, "NameOfBank"):                            "UJI BANK",
		models.JalurAnak(models.DaftarAdjustment, 1, "NoAccount"):                             "UJI-REK"}
	isi := func(aksi string, m map[string]string) {
		t.Helper()
		kode, out := u.aksi(id, teknik, models.WorkbasketAcceptation, aksi, 1, "", m)
		u.wajib(kode, http.StatusOK, out, aksi)
	}
	isi("SetNameCurrency", map[string]string{models.JalurAnak(models.DaftarAdjustment, 1, "CurrencyID"): matauang})
	isi("CountGrossAdjTreaty", map[string]string{models.JalurAnak(models.DaftarAdjustment, 1, "Type"): "1",
		models.JalurAnak(models.DaftarAdjustment, 1, "GrossAdjustment"): "100"})
	kirim := func() (int, map[string]any) {
		return u.aksi(id, teknik, models.WorkbasketAcceptation, "AddKomiteTreatyChild", 1, "", isian)
	}
	// tanpa lampiran: tombol pembuka "Send to Committe" (BukaKomite, AttachmentProtect_ACT) menampilkan pesan lampiran;
	// "Send Claim to Committee" dipanggil langsung ditolak 422 dengan pesan yang sama (AC 58)
	kode, out := u.aksi(id, teknik, models.WorkbasketAcceptation, "BukaKomite", 1, "", isian)
	u.wajib(kode, http.StatusOK, out, "Send to Committe tanpa lampiran")
	if p, _ := json.Marshal(out); !strings.Contains(string(p), "Please Upload Attachment LOD") ||
		!strings.Contains(string(p), "Please Upload Attachment DLA") || !strings.Contains(string(p), "Please Upload Attachment SPGR") {
		t.Fatalf("pesan lampiran AttachmentProtect_ACT: %s", p)
	}
	kode, out = kirim()
	u.wajib(kode, http.StatusUnprocessableEntity, out, "Send Claim to Committee tanpa lampiran")
	if !strings.Contains(strings.Join(teks(out["pesan"]), ";"), "Please Upload Attachment LOD") {
		t.Fatalf("pesan lampiran: %v", out["pesan"])
	}
	for _, k := range []string{"LOD", "DLA", "SPGR"} {
		kode, out = u.unggah(id, teknik, models.WorkbasketAcceptation, k, map[string]string{"UJI-" + k + ".pdf": "%PDF"})
		u.wajib(kode, http.StatusOK, out, "unggah "+k)
	}
	kode, out = kirim()
	u.wajib(kode, http.StatusOK, out, "Send Claim to Committee sesudah lampiran lengkap")
	var komite string
	for kid := range u.g.Kasus {
		if strings.HasPrefix(kid, models.AwalanKomite) {
			komite = kid
		}
	}
	if komite == "" || len(u.g.Komite[komite]) != 1 || u.g.Komite[komite][0].OperatorID != "UJI-K1" {
		t.Fatalf("kasus komite / tangga: %q %+v", komite, u.g.Komite)
	}
	if b := u.g.Halaman(id).AmbilDaftar(models.DaftarAdjustment)[0]; b[models.PropKomiteID] != komite || b["IsKomite"] != "1" {
		t.Fatalf("baris adjustment tertaut: %+v", b)
	}
}

// sampaiAdjustment - alur TestAlurPenuhSampaiResolved sampai satu baris adjustment di Input Acceptation (Teknik).
func sampaiAdjustment(u *uji) string {
	u.t.Helper()
	id := u.buat()
	langkah := func(pelaku, peran, aksi string, n int, param string, m map[string]string) {
		u.t.Helper()
		kode, out := u.aksi(id, pelaku, peran, aksi, n, param, m)
		u.wajib(kode, http.StatusOK, out, aksi)
	}
	langkah(admin, "", "SetValueToClaim", 0, masterUji+"|UJI-TG|UJI-COB", nil)
	langkah(admin, "", "CheckNoPolicy", 0, polisUji, nil)
	langkah(admin, "", "GetNameCauseofLoss", 0, "UJI-S1", nil)
	langkah(admin, "", "SetOutstanding", 0, "", map[string]string{
		models.CD + "PolicyData.StartDateTime": "2026-01-01", models.CD + "PolicyData.EndDateTime": "2026-12-31",
		models.CD + "DateOfLoss": "2026-05-01 08:00:00", models.CD + "ReportDate": "2026-05-02",
		models.CD + "DateReceived": "2026-05-03 09:00:00", models.CD + "ReporterName": "UJI-PELAPOR",
		models.CD + "ReporterTelp": "1", models.CD + "ReportAddress": "UJI-ALAMAT", models.CD + "Location": "uji-lokasi",
		models.CD + "Province": "UJI-PROVINSI", models.CD + "ConsultantID": "UJI-ADJ", models.CD + "AppointedADJID": "UJI-ADJ",
		models.CD + "ShareCeding": "100", models.CD + "Occupation": "uji-okupasi"})
	langkah(admin, "", "AddInterest", 0, "", nil)
	langkah(admin, "", "SetCurencyInterest", 1, "", map[string]string{
		models.JalurAnak(models.DaftarInterest, 1, "CurrencyID"): matauang,
		models.JalurAnak(models.DaftarInterest, 1, "ObjectName"): "UJI-OBJ"})
	langkah(admin, "", "CountTotalInsterest", 1, "", map[string]string{
		models.JalurAnak(models.DaftarInterest, 1, "TSIPerObject"): "1000"})
	langkah(admin, "", "AddListClaimAmount", 0, "", nil)
	langkah(admin, "", "AddLossAllocation", 0, "", nil)
	langkah(admin, "", "SetNameTreaty", 1, "", map[string]string{
		models.JalurAnak(models.DaftarLossAlloc, 1, "TreatyName"): "QUOTA SHARE"})
	langkah(admin, "", "CountPersen", 1, "", map[string]string{
		models.JalurAnak(models.DaftarLossAlloc, 1, "SharePercentage"): "100"})
	langkah(admin, "", "AddEstimation", 0, "", nil)
	langkah(admin, "", "CountEstimation", 1, "", map[string]string{models.JalurAnak(models.DaftarEstimasi, 1, "Type"): "1"})
	h := u.g.Halaman(id)
	h.SetelDaftar(models.DaftarSpreading, []models.Baris{{"CurrencyID": matauang, "Currency": "UJA"}})
	h.SetelDaftar(models.DaftarBreakQS, []models.Baris{{"CurrencyID": matauang, "Currency": "UJA"}})
	u.g.SetelHalaman(id, h)
	langkah(admin, "", "CountSpreading", 1, "", map[string]string{
		models.JalurAnak(models.DaftarSpreading, 1, "TreatyName"):      "UJI-QS",
		models.JalurAnak(models.DaftarSpreading, 1, "SharePercentage"): "100"})
	langkah(admin, "", "SaveOutstanding", 0, "", nil)
	langkah(admin, "", "CheckNopolicy", 0, "", nil)
	langkah(teknik, models.WorkbasketAcceptation, "AddAdjustment", 0, "", nil)
	return id
}

// Master kosong (536 belum dijalankan) dan nol dokumen: daftar dikirim sebagai array kosong, bukan null - panel
// frontend memanggil `.map` (08-10-2026: "Cannot read properties of null (reading 'map')").
func TestLampiranKosongBukanNull(t *testing.T) {
	g, a := tiruan.Baru(), acuanUji()
	jam := func() time.Time { return time.Date(2026, 5, 20, 10, 0, 0, 0, models.Jakarta) }
	u := &uji{t: t, srv: handlers.Router(services.Baru(g, a, jam, false).DenganPenyimpanan(g.BerkasBaru()), true), g: g, a: a}
	id := u.buat()
	r := httptest.NewRequest(http.MethodGet, handlers.Prefix+"/kasus/"+id+"/lampiran", nil)
	r.Header.Set("X-Pelaku", admin)
	w := httptest.NewRecorder()
	u.srv.ServeHTTP(w, r)
	if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), `"kategori":[]`) ||
		!strings.Contains(w.Body.String(), `"lampiran":[]`) {
		t.Fatalf("daftar kosong harus []: %d %s", w.Code, w.Body.String())
	}
}

// View File: GetBase64Attachment -> GetUrlGoogleStorage_Act (Durasi 1800); isi dibaca backend.
func TestLampiranViewFile(t *testing.T) {
	u, _ := baruUjiLampiran(t)
	id := u.buat()
	kode, out := u.unggah(id, admin, "", "DLA", map[string]string{"UJI-DLA.pdf": "%PDF-UJI-ISI"})
	u.wajib(kode, http.StatusOK, out, "unggah DLA")
	lid := u.g.Dokumen[0].ID
	ambil := func(lid, pelaku string) *httptest.ResponseRecorder {
		r := httptest.NewRequest(http.MethodGet, handlers.Prefix+"/kasus/"+id+"/lampiran/"+lid+"/isi", nil)
		r.Header.Set("X-Pelaku", pelaku)
		w := httptest.NewRecorder()
		u.srv.ServeHTTP(w, r)
		return w
	}
	// siapa pun yang boleh membuka kasus boleh melihat (bukan hanya pemegang)
	w := ambil(lid, lain)
	if w.Code != http.StatusOK || w.Body.String() != "%PDF-UJI-ISI" ||
		!strings.Contains(w.Header().Get("Content-Disposition"), "UJI-DLA.pdf") {
		t.Fatalf("View File: %d %q %q", w.Code, w.Body.String(), w.Header().Get("Content-Disposition"))
	}
	if w := ambil("UJI-TIDAK-ADA", admin); w.Code != http.StatusNotFound {
		t.Fatalf("lampiran tak ada: %d", w.Code)
	}
}

// Add attachment: satu kiriman banyak berkas, kategori PER BERKAS (`.pyCategory`, GCNMSaveAttachments S1.2 bila
// TempInputParam.pyCategory kosong). Satu kategori di luar master = seluruh kiriman ditolak, nol berkas terunggah.
func TestLampiranMultiKategoriPerBerkas(t *testing.T) {
	u, b := baruUjiLampiran(t)
	id := u.buat()
	kirim := func(berkas [][3]string) (int, map[string]any) {
		var buf bytes.Buffer
		mw := multipart.NewWriter(&buf)
		for _, x := range berkas {
			f, _ := mw.CreateFormFile("berkas", x[0])
			_, _ = f.Write([]byte(x[1]))
			_ = mw.WriteField("kategoriBerkas", x[2])
		}
		_ = mw.Close()
		r := httptest.NewRequest(http.MethodPost, handlers.Prefix+"/kasus/"+id+"/lampiran", &buf)
		r.Header.Set("Content-Type", mw.FormDataContentType())
		r.Header.Set("X-Pelaku", admin)
		w := httptest.NewRecorder()
		u.srv.ServeHTTP(w, r)
		var out map[string]any
		_ = json.Unmarshal(w.Body.Bytes(), &out)
		return w.Code, out
	}
	kode, out := kirim([][3]string{{"UJI-A.pdf", "a", "LOD"}, {"UJI-B.csv", "b", "File"}, {"UJI-C.pdf", "c", "UJI-BUKAN"}})
	u.wajib(kode, http.StatusUnprocessableEntity, out, "satu kategori di luar master")
	if len(b.Unggahan) != 0 || len(u.g.Dokumen) != 0 {
		t.Fatalf("kategori diperiksa sebelum unggah: %d unggahan", len(b.Unggahan))
	}
	kode, out = kirim([][3]string{{"UJI-A.pdf", "a", "LOD"}, {"UJI-B.csv", "b", "File"}})
	u.wajib(kode, http.StatusOK, out, "dua berkas dua kategori")
	if c := u.cacah(id, admin); c["LOD"] != 1 || c["File"] != 1 {
		t.Fatalf("AttachCategory: %v", c)
	}
	if u.g.Dokumen[1].Kategori1 != "File" || u.g.Dokumen[1].MIME != "csv" {
		t.Fatalf("berkas kedua: %+v", u.g.Dokumen[1])
	}
}

// View File pola NB Treaty In: View Office Online (xls...pptx bertanda objek) dan Delete (pemegang, kasus terbuka).
func TestLampiranOfficeDanDelete(t *testing.T) {
	u, b := baruUjiLampiran(t)
	id := u.buat()
	kode, out := u.unggah(id, admin, "", "File", map[string]string{"UJI-DATA.xlsx": "x"})
	u.wajib(kode, http.StatusOK, out, "unggah xlsx")
	kode, out = u.unggah(id, admin, "", "LOD", map[string]string{"UJI-LOD.pdf": "%PDF"})
	u.wajib(kode, http.StatusOK, out, "unggah pdf")
	xlsx, pdf := u.g.Dokumen[0].ID, u.g.Dokumen[1].ID
	jalur := handlers.Prefix + "/kasus/" + id + "/lampiran/"
	kode, out = u.minta(http.MethodGet, jalur+xlsx+"/office", lain, "", nil)
	u.wajib(kode, http.StatusOK, out, "View Office Online xlsx (bukan pemegang boleh melihat)")
	if out["url"] != "UJI-URL-UJI-IMG-1" {
		t.Fatalf("URL penampil: %v", out)
	}
	kode, out = u.minta(http.MethodGet, jalur+pdf+"/office", admin, "", nil)
	u.wajib(kode, http.StatusUnprocessableEntity, out, "View Office Online pdf")
	kode, out = u.minta(http.MethodPost, jalur+pdf+"/hapus", lain, "", nil)
	u.wajib(kode, http.StatusForbidden, out, "Delete bukan pemegang")
	kode, out = u.minta(http.MethodPost, jalur+pdf+"/hapus", admin, "", nil)
	u.wajib(kode, http.StatusOK, out, "Delete pemegang")
	if len(u.g.Dokumen) != 1 || u.g.Dokumen[0].ID != xlsx || len(b.Dihapus) != 1 || b.Dihapus[0] != "UJI-IMG-2" ||
		len(u.g.Storage) != 1 {
		t.Fatalf("sesudah Delete: dokumen %+v, objek dihapus %v, catatan %d", u.g.Dokumen, b.Dihapus, len(u.g.Storage))
	}
	if c := u.cacah(id, admin); c["LOD"] != 0 || c["File"] != 1 {
		t.Fatalf("AttachCategory sesudah Delete: %v", c)
	}
}

// Change Category (layar Pega View File): dokumen terpilih pindah kategori dalam satu transaksi; pemegang saja.
func TestLampiranPindahKategori(t *testing.T) {
	u, _ := baruUjiLampiran(t)
	id := u.buat()
	kode, out := u.unggah(id, admin, "", "File", map[string]string{"UJI-A.pdf": "a", "UJI-B.pdf": "b"})
	u.wajib(kode, http.StatusOK, out, "unggah dua berkas")
	ids := []string{u.g.Dokumen[0].ID, u.g.Dokumen[1].ID}
	jalur := handlers.Prefix + "/kasus/" + id + "/lampiran/kategori"
	kode, out = u.minta(http.MethodPost, jalur, lain, "", map[string]any{"kategori": "LOD", "ids": ids})
	u.wajib(kode, http.StatusForbidden, out, "bukan pemegang")
	kode, out = u.minta(http.MethodPost, jalur, admin, "", map[string]any{"kategori": "UJI-BUKAN", "ids": ids})
	u.wajib(kode, http.StatusUnprocessableEntity, out, "kategori di luar master")
	kode, out = u.minta(http.MethodPost, jalur, admin, "", map[string]any{"kategori": "LOD", "ids": []string{ids[0], "UJI-X"}})
	u.wajib(kode, http.StatusNotFound, out, "dokumen tak dikenal")
	if c := u.cacah(id, admin); c["File"] != 2 || c["LOD"] != 0 {
		t.Fatalf("penolakan tidak boleh memindah apa pun: %v", c)
	}
	kode, out = u.minta(http.MethodPost, jalur, admin, "", map[string]any{"kategori": "LOD", "ids": ids})
	u.wajib(kode, http.StatusOK, out, "Change Category")
	if c := u.cacah(id, admin); c["File"] != 0 || c["LOD"] != 2 {
		t.Fatalf("AttachCategory sesudah Change Category: %v", c)
	}
}
