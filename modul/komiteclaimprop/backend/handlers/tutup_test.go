package handlers_test

// Kasus komite Close Without Payment (TT 4, `KomitePost_Close`; perintah work owner 10-10-2026): kasus tanpa baris
// adjustment, satu tingkat. Disetujui -> OS close STS 4, CLAIMREJECTED, konversi "4", klaim induk Resolved-Completed;
// ditolak -> klaim induk tetap terbuka.

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"nusantarare/inti/backend/kontrak"
	"nusantarare/modul/komiteclaimprop/backend/models"
	"nusantarare/modul/komiteclaimprop/backend/tiruan"
)

const (
	komiteTutupUji = "TKMT-UJI009"
	pemutusTutup   = "UJI-K9"
)

// siapTutup - fixture klaim induk + kasus komite TT 4 satu tingkat (Claim Dept. Head) di atasnya.
func siapTutup(t *testing.T, produksi bool) *uji {
	t.Helper()
	u := siap(t, 1, produksi)
	nilai := u.g.Klaim.Nilai(tiruan.KlaimUji)
	nilai["ClaimData.Remark_Close"] = "UJI-REMARKS-CWP"
	nilai["TreatyInMaster.ID"] = "UJI-MASTER-ID"
	nilai[kontrak.JalurPembuatKlaimTreaty] = tiruan.PembuatUji
	nilai[kontrak.JalurNamaPembuatKlaimTreaty] = "UJI Admin"
	daftar := map[string][]map[string]string{}
	for _, j := range []string{"ClaimData.AdjustmentList", "ClaimData.SpreadingRisk", "ClaimData.EstimationList",
		"ClaimData.InterestList", "ClaimData.ListClaimAmount", "ClaimData.SuggestList"} {
		daftar[j] = u.g.Klaim.Daftar(tiruan.KlaimUji, j)
	}
	u.g.Klaim.Setel(tiruan.KlaimUji, nilai, daftar)
	u.g.LahirkanTutup(komiteTutupUji, tiruan.KlaimUji, "UJI-KRONOLOGI", tiruan.PembuatUji, "UJI Admin",
		[]models.Anggota{{OperatorID: pemutusTutup, Jabatan: "Claim Dept. Head"}}, tiruan.SaatUji)
	return u
}

func (u *uji) putusTutup(pelaku string, kep models.Keputusan, mau int) {
	u.t.Helper()
	w := u.minta("POST", "/kasus/"+komiteTutupUji+"/putuskan", pelaku, kep)
	if w.Code != mau {
		u.t.Fatalf("putuskan close %s: %d (mau %d) %s", pelaku, w.Code, mau, w.Body.String())
	}
}

func riwayatKlaim(u *uji) string {
	var out []string
	for _, b := range u.g.Klaim.Daftar(tiruan.KlaimUji, "ClaimData.SuggestList") {
		out = append(out, b["CommentSuggest"]+"|"+b["IsCedingConfirm"])
	}
	return strings.Join(out, ";")
}

func TestTutupDisetujuiMenutupKlaim(t *testing.T) {
	u := siapTutup(t, true)
	osAwal, efekAwal := len(u.g.OS), len(u.g.Efek)
	u.putusTutup("UJI-LAIN", setuju("UJI"), http.StatusForbidden)
	// Subjectivity / Propose tidak tampil untuk TT 4: dikirim pun diabaikan
	kep := setuju("UJI-SETUJU-CWP")
	kep.IsSubjectivity, kep.SubjectivityNote, kep.UsulTutup = true, "1", true
	u.putusTutup(pemutusTutup, kep, http.StatusOK)

	k := u.g.Kasus[komiteTutupUji]
	if !k.Tertutup() || k.Count != 2 || k.Subjectivity != models.UsulTidak || k.UsulTutup != models.UsulTidak {
		t.Fatalf("kepala kasus komite close %+v", k)
	}
	if got := u.g.Tangga[komiteTutupUji]; len(got) != 1 || got[0].Keputusan != models.KeputusanSetuju ||
		got[0].Komentar != "UJI-SETUJU-CWP" {
		t.Fatalf("tangga %+v", got)
	}
	if len(u.g.OS) != osAwal+1 {
		t.Fatalf("baris OS %d, mau %d", len(u.g.OS), osAwal+1)
	}
	os := u.g.OS[len(u.g.OS)-1]
	if os.StsReject != models.StsOSTutup || os.MasterID != "UJI-MASTER-ID" || os.CaseID != tiruan.KlaimUji ||
		!strings.Contains(os.DataJSON, `"IDMasterTreaty":"UJI-MASTER-ID"`) {
		t.Fatalf("baris OS close %+v", os)
	}
	if len(u.g.Ditolak) != 1 {
		t.Fatalf("CLAIMREJECTED %d baris", len(u.g.Ditolak))
	}
	if d := u.g.Ditolak[0]; d.Label != models.LabelKlaimTreaty || d.StatusWork != models.StatusKlaimTerbuka ||
		d.Kelas != models.KelasKlaim || d.Remark != "UJI-REMARKS-CWP" || d.PembuatID != tiruan.PembuatUji ||
		d.PengubahID != pemutusTutup || d.ID != tiruan.KlaimUji {
		t.Fatalf("baris CLAIMREJECTED %+v", d)
	}
	if len(u.g.Klaim.Ditutup) != 1 || u.g.Klaim.Ditutup[0] != tiruan.KlaimUji {
		t.Fatalf("klaim induk ditutup %v", u.g.Klaim.Ditutup)
	}
	if !strings.Contains(riwayatKlaim(u), models.TeksTutupDisetujui+"Claim Dept. Head|Claim Dept. Head") {
		t.Fatalf("riwayat klaim %q", riwayatKlaim(u))
	}
	if len(u.g.Riwayat) != 0 {
		t.Fatalf("HISTORYAKSEPTASIPEGA ditulis: %+v", u.g.Riwayat)
	}
	if _, ada := u.g.JSONKlaim[tiruan.KlaimUji]; !ada {
		t.Fatal("JSON_KLAIM tidak ditulis")
	}
	konversi := ""
	for _, e := range u.g.Efek[efekAwal:] {
		if strings.HasPrefix(e, "konversi-klaim:") {
			konversi = e
		}
	}
	if !strings.Contains(konversi, `"STS_REJECT":"4"`) {
		t.Fatalf("efek konversi %q (efek %v)", konversi, u.g.Efek[efekAwal:])
	}
	// kasus selesai: Submit kedua 409
	u.putusTutup(pemutusTutup, setuju("UJI"), http.StatusConflict)
}

func TestTutupDitolakKlaimTetapTerbuka(t *testing.T) {
	u := siapTutup(t, true)
	osAwal, efekAwal := len(u.g.OS), len(u.g.Efek)
	u.putusTutup(pemutusTutup, models.Keputusan{AcceptStatus: models.KeputusanTolak, Comment: "UJI-TOLAK"},
		http.StatusOK)
	if !u.g.Kasus[komiteTutupUji].Tertutup() {
		t.Fatal("kasus komite close belum selesai sesudah ditolak")
	}
	if len(u.g.Klaim.Ditutup) != 0 || len(u.g.OS) != osAwal || len(u.g.Ditolak) != 0 {
		t.Fatalf("penolakan menulis close: ditutup %v OS %d ditolak %d", u.g.Klaim.Ditutup, len(u.g.OS), len(u.g.Ditolak))
	}
	for _, e := range u.g.Efek[efekAwal:] {
		if strings.HasPrefix(e, "konversi-klaim:") {
			t.Fatalf("konversi diantre saat ditolak: %q", e)
		}
	}
	if !strings.Contains(riwayatKlaim(u), models.TeksTutupDitolak+"Claim Dept. Head") {
		t.Fatalf("riwayat klaim %q", riwayatKlaim(u))
	}
}

func TestBukaKasusTutupWajahClose(t *testing.T) {
	u := siapTutup(t, false)
	w := u.minta("GET", "/kasus/"+komiteTutupUji, pemutusTutup, nil)
	if w.Code != http.StatusOK {
		t.Fatalf("buka kasus close: %d %s", w.Code, w.Body.String())
	}
	var ly models.Layar
	if err := json.Unmarshal(w.Body.Bytes(), &ly); err != nil {
		t.Fatal(err)
	}
	if strings.Join(ly.Judul, " ") != models.JudulKomite+" "+models.JudulTutup || !ly.BolehKerja || ly.Isian.Adjustment {
		t.Fatalf("judul / kerja / isian %v %v %v", ly.Judul, ly.BolehKerja, ly.Isian.Adjustment)
	}
	teks := map[string]string{}
	for _, b := range ly.Bagian {
		switch b.Kunci {
		case "totalEstimasi", "riwayatAdjustment", "deductible", "spreading", "bayar", "bank":
			t.Errorf("bagian TT 2 %q tampil di wajah CLOSE", b.Kunci)
		case "estimasi":
			if len(b.Grid) != 1 {
				t.Errorf("grid estimasi %d, mau hanya Estimation List", len(b.Grid))
			}
		case "teksKomite":
			for _, m := range b.Medan {
				teks[m.Label] = m.Nilai
			}
		}
	}
	if teks["Circumstances"] != "UJI-KRONOLOGI" || teks["Remarks"] != "UJI-REMARKS-CWP" {
		t.Fatalf("teks komite close %v", teks)
	}
}
