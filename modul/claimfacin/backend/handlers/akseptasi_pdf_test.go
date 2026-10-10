package handlers_test

// Uji dokumen akseptasi tombol "Acceptation" (SaveAcceptation 12 `PrintPDFAccep_MultiAksep`, stream
// AcceptanceNotePDF): PDF diunggah sesudah aksi tersimpan dan dicatat DOCUMENT_CLAIM; klik ulang (IsPrintAccept 1) tanpa
// dokumen kedua; lini tanpa stream -> info OQ; unggah gagal -> akseptasi tetap tersimpan + PesanDokumenGagal.

import (
	"bytes"
	"compress/zlib"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"nusantarare/modul/claimfacin/backend/handlers"
	"nusantarare/modul/claimfacin/backend/models"
	"nusantarare/modul/claimfacin/backend/services"
	"nusantarare/modul/claimfacin/backend/tiruan"
)

const noAksepUji = "UJI-AKS.03.2026.00001"

// sampaiAkseptasi - adjustment 1 item 1 diterima komite tahap 2 (fixture, IsPrintAccept masih kosong), siap
// "Acceptation".
func (u *uji) sampaiAkseptasi() string {
	u.t.Helper()
	u.a.Adjuster["UJI-ADJ1"] = "UJI ADJUSTER"
	u.a.Roster = []models.AnggotaKomite{
		{ID: "1", OperatorID: "UJI-K1", Jabatan: "UJI SPV", Degree: "1", LimitBottom: "-9999999999999", LimitTop: "57750000"},
	}
	id := u.sampaiSurveyor()
	itemAdj := models.KunciPanel(models.PanelItemAdj, models.DaftarItem(1), 1)
	isi := map[string]string{models.CD + "ConsultantID": "UJI-ADJ1", models.CD + "AppointedADJID": "UJI-ADJ1"}
	for _, p := range []services.PermintaanAksi{
		{Aksi: "CountTotalEstimasi", Konteks: itemAdj, Masukan: isi},
		{Aksi: "SetAdjTypePayment", Konteks: panelAdj(1), Masukan: map[string]string{jAdj(1, "PaymentType"): "1"}},
		{Aksi: "CheckCurrency", Konteks: panelAdj(1), Masukan: map[string]string{jAdj(1, "UploadLOD"): idr}},
		{Aksi: "SetGrossAdjustment", Konteks: panelAdj(1),
			Masukan: map[string]string{jAdj(1, "GrossAdjustment"): "20000000"}},
	} {
		kode, out := u.aksiT(id, p)
		u.wajib(kode, http.StatusOK, out, p.Aksi)
	}
	h := u.g.Halaman(id)
	b := h.AmbilDaftar(models.DaftarAdj(1, 1))[0]
	b["AcceptanceStatus"], b["AcceptedNo"], b["IsApproved"], b["IsKomite"] = "1", noAksepUji, "1", "1"
	b["AcceptedDate"] = "2026-03-05 10:00:00"
	u.g.SetelHalaman(id, h)
	return id
}

// dokumenAkseptasi - baris DOCUMENT_CLAIM kategori AcceptanceNote.
func (u *uji) dokumenAkseptasi() []models.BarisDokumenKlaim {
	var out []models.BarisDokumenKlaim
	for _, d := range u.g.Dokumen {
		if d.Kategori1 == models.KategoriDokumenAkseptasi {
			out = append(out, d)
		}
	}
	return out
}

func TestAcceptationMenyimpanDokumenAkseptasi(t *testing.T) {
	u := baruUji(t)
	id := u.sampaiAkseptasi()
	kode, out := u.aksiT(id, services.PermintaanAksi{Aksi: "Acceptation", Konteks: panelAdj(1)})
	u.wajib(kode, http.StatusOK, out, "acceptation")
	if out["info"] != nil {
		t.Fatalf("info Acceptation lini Fire (stream AcceptanceNotePDF diekspor): %v", out["info"])
	}
	nama := models.NamaBerkasAkseptasi(id, noAksepUji)
	dok := u.dokumenAkseptasi()
	if len(dok) != 1 || dok[0].IDPega != models.KunciInstans(id) || dok[0].NamaFile != nama || dok[0].MIME != "pdf" ||
		dok[0].Operator != teknik || dok[0].StorageID == "" {
		t.Fatalf("DOCUMENT_CLAIM AcceptanceNote: %+v", dok)
	}
	if len(u.g.Storage) == 0 || u.g.Storage[len(u.g.Storage)-1].ImageID != dok[0].StorageID {
		t.Fatalf("T_STORAGE_IMAGE: %+v", u.g.Storage)
	}
	m := u.b.Unggahan[len(u.b.Unggahan)-1]
	if m.NamaFile != nama || m.Folder != "Claim" || m.Ext != "pdf" || m.Durasi != 1800 || m.Pengguna != teknik ||
		!bytes.HasPrefix(m.Isi, []byte("%PDF-")) {
		t.Fatalf("unggahan: %+v", m)
	}
	teks := teksPDF(t, m.Isi)
	for _, s := range []string{"(ACCEPTED CLAIM INSURANCE)", "(Accepted No : " + noAksepUji + ")", "(UJI-RNM-F.001)",
		"(01/01/2026- 31/12/2026)", "(Accepted Claim)", "(UJI QS)", "(UJI FAC RETRO)", "(Jakarta, 05 March 2026)"} {
		if !strings.Contains(teks, s) {
			t.Errorf("PDF tanpa %s", s)
		}
	}
	// tampil di daftar lampiran klaim (DOCUMENT_CLAIM per IDPEGA, semua kategori) dan terbuka lewat View File
	kode, out = u.minta(http.MethodGet, handlers.Prefix+"/kasus/"+id+"/lampiran", teknik, perTeknik, nil)
	u.wajib(kode, http.StatusOK, out, "daftar lampiran")
	ada := false
	for _, l := range out["lampiran"].([]any) {
		if m := l.(map[string]any); m["namaFile"] == nama && m["kategori"] == models.KategoriDokumenAkseptasi &&
			m["id"] == dok[0].ID {
			ada = true
		}
	}
	if !ada {
		t.Fatalf("PDF akseptasi tidak di daftar lampiran: %v", out["lampiran"])
	}
	r := httptest.NewRequest(http.MethodGet, handlers.Prefix+"/kasus/"+id+"/lampiran/"+dok[0].ID+"/isi", nil)
	r.Header.Set("X-Pelaku", teknik)
	r.Header.Set("X-Peran", perTeknik)
	w := httptest.NewRecorder()
	u.srv.ServeHTTP(w, r)
	if w.Code != http.StatusOK || !bytes.Equal(w.Body.Bytes(), m.Isi) {
		t.Fatalf("View File PDF akseptasi: HTTP %d, %d byte", w.Code, w.Body.Len())
	}
	// klik ulang: tombol nonaktif (IsPrintAccept 1 && DirectToKasir 'false') - tanpa dokumen kedua
	kode, out = u.aksiT(id, services.PermintaanAksi{Aksi: "Acceptation", Konteks: panelAdj(1)})
	u.wajib(kode, http.StatusConflict, out, "acceptation ulang")
	if n := len(u.dokumenAkseptasi()); n != 1 {
		t.Fatalf("klik ulang menambah dokumen: %d", n)
	}
}

func TestAcceptationLiniTanpaStream(t *testing.T) {
	// S17 lini MBU memakai stream AcceptanceNotePDFMBU yang tidak diekspor korpus: akseptasi tersimpan, tanpa berkas.
	u := baruUji(t)
	id := u.sampaiAkseptasi()
	// OfferFacIn dibaca ulang dari JSON_POLIS setiap muat (muatPolis): lini diganti di dokumen polis
	u.a.Dokumen[tiruan.KunciPolis(polisUji, "1")] = []byte(strings.Replace(polisFire, `"BusinessType": "Fire"`,
		`"BusinessType": "MBUCar"`, 1))
	kode, out := u.aksiT(id, services.PermintaanAksi{Aksi: "Acceptation", Konteks: panelAdj(1)})
	u.wajib(kode, http.StatusOK, out, "acceptation MBU")
	if out["info"] != models.InfoTanpaStream(models.StreamAkseptasiMBU) {
		t.Fatalf("info: %v", out["info"])
	}
	if b := u.g.Halaman(id).AmbilDaftar(models.DaftarAdj(1, 1))[0]; b["IsPrintAccept"] != "1" || len(u.dokumenAkseptasi()) != 0 {
		t.Fatalf("akseptasi MBU: adj %v dokumen %d", b, len(u.dokumenAkseptasi()))
	}
}

func TestAcceptationUnggahGagalAkseptasiTetap(t *testing.T) {
	// [penyimpangan sadar] pola Komite Claim Prop: unggah gagal sesudah commit -> akseptasi TETAP tersimpan, pelaku diberi
	// PesanDokumenGagal, tanpa baris DOCUMENT_CLAIM.
	u := baruUji(t)
	id := u.sampaiAkseptasi()
	u.b.Gagal = errors.New("UJI penyimpanan tak terjangkau")
	kode, out := u.aksiT(id, services.PermintaanAksi{Aksi: "Acceptation", Konteks: panelAdj(1)})
	u.wajib(kode, http.StatusOK, out, "acceptation unggah gagal")
	if out["info"] != models.PesanDokumenGagal {
		t.Fatalf("info: %v", out["info"])
	}
	if b := u.g.Halaman(id).AmbilDaftar(models.DaftarAdj(1, 1))[0]; b["IsPrintAccept"] != "1" || len(u.dokumenAkseptasi()) != 0 {
		t.Fatalf("akseptasi: adj %v dokumen %d", b, len(u.dokumenAkseptasi()))
	}
}

// awalStream - penanda awal isi stream PDF yang ditulis fpdf.
const awalStream = "stream\n"

// teksPDF - isi stream FlateDecode PDF (konten halaman) sebagai teks.
func teksPDF(t *testing.T, b []byte) string {
	t.Helper()
	var out strings.Builder
	for {
		i := bytes.Index(b, []byte(awalStream))
		if i < 0 {
			return out.String()
		}
		b = b[i+len(awalStream):]
		j := bytes.Index(b, []byte("endstream"))
		if j < 0 {
			t.Fatal("stream tanpa endstream")
		}
		if r, err := zlib.NewReader(bytes.NewReader(b[:j])); err == nil {
			isi, _ := io.ReadAll(r)
			out.Write(isi)
		}
		b = b[j+len("endstream"):]
	}
}
