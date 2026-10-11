package handlers_test

// Kasus komite Reject Claim (TT 3, `KomitePost_Reject`; perintah work owner 10-10-2026): disetujui -> OS per baris
// estimasi (SaveRejectTreatyIn_Act_KMT, Type 2), CLAIMREJECTED, konversi "2", klaim induk Resolved-Completed; ditolak ->
// klaim induk tetap terbuka. Keduanya: PDF CommitteReject_CC (S12, stream ClaimTreaty diekspor work owner 10-10-2026).

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"nusantarare/inti/backend/kontrak"
	"nusantarare/modul/komiteclaimprop/backend/models"
	"nusantarare/modul/komiteclaimprop/backend/tiruan"
)

const komiteTolakUji = "TKMT-UJI010"

// siapTolak - fixture klaim induk berestimasi dua baris + kasus komite TT 3 satu tingkat (Claim Dept. Head).
func siapTolak(t *testing.T) *uji {
	t.Helper()
	u := siap(t, 1, true)
	nilai := u.g.Klaim.Nilai(tiruan.KlaimUji)
	nilai["ClaimData.Remark_Close"] = "UJI-REMARKS-REJECT"
	nilai[kontrak.JalurPembuatKlaimTreaty] = tiruan.PembuatUji
	nilai[kontrak.JalurNamaPembuatKlaimTreaty] = "UJI Admin"
	daftar := map[string][]map[string]string{}
	for _, j := range []string{"ClaimData.AdjustmentList", "ClaimData.SpreadingRisk", "ClaimData.InterestList",
		"ClaimData.ListClaimAmount", "ClaimData.SuggestList"} {
		daftar[j] = u.g.Klaim.Daftar(tiruan.KlaimUji, j)
	}
	daftar["ClaimData.AdjustmentList"] = nil // reject hanya sebelum Claim No / akseptasi
	daftar["ClaimData.EstimationList"] = []map[string]string{
		{"Currency": "IDR", "CurrencyID": "UJI-IDR", "EstimationValue": "250", "GrossEstimationPct": "1000", "PrintFaceClaim": "1"},
		{"Currency": "USD", "CurrencyID": "UJI-USD", "EstimationValue": "5", "GrossEstimationPct": "20", "PrintFaceClaim": "1"},
	}
	u.g.Klaim.Setel(tiruan.KlaimUji, nilai, daftar)
	u.g.LahirkanTanpaAdj(komiteTolakUji, tiruan.KlaimUji, models.TransferTolak, "UJI-KRONOLOGI-REJECT",
		tiruan.PembuatUji, "UJI Admin", []models.Anggota{{OperatorID: pemutusTutup, Jabatan: "Claim Dept. Head"}},
		tiruan.SaatUji)
	return u
}

func (u *uji) putusTolak(pelaku string, kep models.Keputusan, mau int) {
	u.t.Helper()
	w := u.minta("POST", "/kasus/"+komiteTolakUji+"/putuskan", pelaku, kep)
	if w.Code != mau {
		u.t.Fatalf("putuskan reject %s: %d (mau %d) %s", pelaku, w.Code, mau, w.Body.String())
	}
}

func TestTolakDisetujuiMenutupKlaim(t *testing.T) {
	u := siapTolak(t)
	osAwal, efekAwal := len(u.g.OS), len(u.g.Efek)
	u.putusTolak(pemutusTutup, setuju("UJI-SETUJU-REJECT"), http.StatusOK)

	if k := u.g.Kasus[komiteTolakUji]; !k.Tertutup() || k.Count != 2 {
		t.Fatalf("kepala kasus komite reject %+v", k)
	}
	baru := u.g.OS[osAwal:]
	if len(baru) != 2 {
		t.Fatalf("baris OS reject %d, mau 2 (per baris estimasi)", len(baru))
	}
	for i, mau := range [][3]string{{"UJI-IDR", "250", "1000"}, {"UJI-USD", "5", "20"}} {
		o := baru[i]
		if o.StsReject != models.StsOSTolak || o.CurrencyID != mau[0] || o.Value != mau[1] || o.GrossValue != mau[2] ||
			o.MasterID != "" || !strings.Contains(o.DataJSON, `"Type":"2"`) || !strings.Contains(o.DataJSON, `"Value":"`+mau[1]+`"`) {
			t.Fatalf("baris OS reject %d: %+v", i, o)
		}
	}
	if len(u.g.Ditolak) != 1 || u.g.Ditolak[0].Remark != "UJI-REMARKS-REJECT" {
		t.Fatalf("CLAIMREJECTED %+v", u.g.Ditolak)
	}
	if len(u.g.Klaim.Ditutup) != 1 {
		t.Fatalf("klaim induk ditutup %v", u.g.Klaim.Ditutup)
	}
	if !strings.Contains(riwayatKlaim(u), models.TeksTolakDisetujui+"Claim Dept. Head|Claim Dept. Head") {
		t.Fatalf("riwayat klaim %q", riwayatKlaim(u))
	}
	konversi := ""
	for _, e := range u.g.Efek[efekAwal:] {
		if strings.HasPrefix(e, "konversi-klaim:") {
			konversi = e
		}
	}
	if !strings.Contains(konversi, `"STS_REJECT":"2"`) {
		t.Fatalf("efek konversi %q", konversi)
	}
	if len(u.g.Riwayat) != 0 {
		t.Fatalf("HISTORYAKSEPTASIPEGA tidak untuk TT 3: %+v", u.g.Riwayat)
	}
	dokumenTolak(t, u)
}

// dokumenTolak - KomitePost_Reject S12.3-S12.7: satu PDF "Reject Claim "+NoClaim+".PDF" (folder Claim, Durasi 1800,
// MIME pdf) diunggah lalu T_STORAGE_IMAGE + DOCUMENT_CLAIM dicatat; kategori tetap CloseClaim (VERBATIM).
func dokumenTolak(t *testing.T, u *uji) {
	t.Helper()
	if len(u.b.Unggahan) != 1 {
		t.Fatalf("PDF reject: unggahan %d", len(u.b.Unggahan))
	}
	m := u.b.Unggahan[0]
	if m.Folder != "Claim" || m.Durasi != 1800 || m.Ext != "pdf" || m.NamaFile != "Reject Claim UJI-K-0001.PDF" ||
		!bytes.HasPrefix(m.Isi, []byte("%PDF-")) {
		t.Fatalf("InsertGoogleStorage_Act: %s %s %d %s", m.Folder, m.NamaFile, m.Durasi, m.Ext)
	}
	if len(u.g.Dokumen) != 1 || u.g.Dokumen[0].Kategori1 != models.KategoriDokumenTutup ||
		u.g.Dokumen[0].NamaFile != m.NamaFile || u.g.Dokumen[0].IDPega != tiruan.KlaimUji {
		t.Fatalf("DOCUMENT_CLAIM: %+v", u.g.Dokumen)
	}
}

func TestTolakDitolakKlaimTetapTerbuka(t *testing.T) {
	u := siapTolak(t)
	osAwal := len(u.g.OS)
	u.putusTolak(pemutusTutup, models.Keputusan{AcceptStatus: models.KeputusanTolak, Comment: "UJI-TOLAK"}, http.StatusOK)
	if !u.g.Kasus[komiteTolakUji].Tertutup() || len(u.g.Klaim.Ditutup) != 0 || len(u.g.OS) != osAwal ||
		len(u.g.Ditolak) != 0 {
		t.Fatalf("penolakan menulis reject: ditutup %v OS %d ditolak %d", u.g.Klaim.Ditutup, len(u.g.OS), len(u.g.Ditolak))
	}
	if !strings.Contains(riwayatKlaim(u), models.TeksTolakDitolak+"Claim Dept. Head") {
		t.Fatalf("riwayat klaim %q", riwayatKlaim(u))
	}
	dokumenTolak(t, u) // S12 tanpa gerbang AcceptStatus
}

func TestBukaKasusTolakWajahReject(t *testing.T) {
	u := siapTolak(t)
	w := u.minta("GET", "/kasus/"+komiteTolakUji, pemutusTutup, nil)
	if w.Code != http.StatusOK {
		t.Fatalf("buka kasus reject: %d %s", w.Code, w.Body.String())
	}
	var ly models.Layar
	if err := json.Unmarshal(w.Body.Bytes(), &ly); err != nil {
		t.Fatal(err)
	}
	if strings.Join(ly.Judul, " ") != models.JudulKomite+" "+models.JudulTolak || !ly.BolehKerja || ly.Isian.Adjustment {
		t.Fatalf("judul / kerja / isian %v %v %v", ly.Judul, ly.BolehKerja, ly.Isian.Adjustment)
	}
	for _, b := range ly.Bagian {
		if b.Kunci == "teksKomite" {
			teks := map[string]string{}
			for _, m := range b.Medan {
				teks[m.Label] = m.Nilai
			}
			if teks["Circumstances"] != "UJI-KRONOLOGI-REJECT" || teks["Remarks"] != "UJI-REMARKS-REJECT" {
				t.Fatalf("teks komite reject %v", teks)
			}
			return
		}
	}
	t.Fatal("bagian teksKomite tidak ada")
}
