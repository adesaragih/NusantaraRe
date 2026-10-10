package handlers_test

// Close Without Payment (`SendCloseClaimToKomite`, perintah work owner 10-10-2026): tombol "Yes" pop-up
// PreventRejectClaimProp bercentang melahirkan kasus komite TKMT- satu tingkat tanpa adjustment.

import (
	"net/http"
	"strings"
	"testing"

	"nusantarare/modul/claimprop/backend/models"
	"nusantarare/modul/claimprop/backend/tiruan"
)

const wbDeptHeadUji = "UJI-WB-DEPTHEAD"

func rosterTutupUji(u *uji) {
	u.a.Roster = []tiruan.AnggotaRoster{
		{AnggotaKomite: models.AnggotaKomite{ID: "1", OperatorID: wbDeptHeadUji, Jabatan: models.JabatanTutupTanpaBayar,
			Degree: "1"}, Batas: "-9999999999999", Sts: models.STSKlaimProp},
		{AnggotaKomite: models.AnggotaKomite{ID: "2", OperatorID: "UJI-WB-TECHDIV", Jabatan: "Technic Div. Head",
			Degree: "2"}, Batas: "189750001", Sts: models.STSKlaimProp},
		// jabatan sama di lini lain tidak boleh terpilih
		{AnggotaKomite: models.AnggotaKomite{ID: "3", OperatorID: "UJI-WB-FACIN", Jabatan: models.JabatanTutupTanpaBayar,
			Degree: "2"}, Batas: "-9999999999999", Sts: "FACIN"},
	}
}

// komiteTutup - kasus komite Close Without Payment klaim `id` (KomiteAdj kosong).
func komiteTutup(u *uji, id string) []string {
	var out []string
	for kid, klaim := range u.g.KomiteKlaim {
		if klaim == id && u.g.KomiteAdj[kid] == "" {
			out = append(out, kid)
		}
	}
	return out
}

func TestTutupTanpaBayarMelahirkanKomiteSatuTingkat(t *testing.T) {
	u := baruUji(t)
	rosterTutupUji(u)
	id := sampaiAdjustment(u)
	cwp := map[string]string{models.PropCentangTutupTanpaBayar: "true", "TempCommiteClaim.Remarks": "UJI-CWP",
		models.PropKronologiTutup: "UJI-KRONOLOGI"}
	kerja := func(aksi string, m map[string]string) (int, map[string]any) {
		return u.aksi(id, teknik, models.WorkbasketAcceptation, aksi, 0, "", m)
	}

	// langkah 3-4: adjustment berstatus kosong menolak
	kode, out := kerja("SendCloseClaimToKomite", cwp)
	u.wajib(kode, http.StatusUnprocessableEntity, out, "CWP dengan adjustment tertunda")
	if !strings.Contains(strings.Join(teks(out["pesan"]), ";"), models.PesanKomiteMasihJalan) {
		t.Fatalf("pesan CWP: %v", out["pesan"])
	}
	if n := len(komiteTutup(u, id)); n != 0 {
		t.Fatalf("kasus komite close lahir walau ditolak: %d", n)
	}
	kode, out = u.aksi(id, teknik, models.WorkbasketAcceptation, "DeleteAjsutment", 1, "", nil)
	u.wajib(kode, http.StatusOK, out, "hapus adjustment")

	// Yes biasa ditolak selama Close Without Payment dicentang (tombol itu tidak tampil)
	kode, out = kerja("CloseClaimProp", cwp)
	if kode == http.StatusOK {
		t.Fatalf("CloseClaimProp lolos dengan centang CWP: %v", out)
	}

	kode, out = kerja("SendCloseClaimToKomite", cwp)
	u.wajib(kode, http.StatusOK, out, "Close Without Payment")
	kt := komiteTutup(u, id)
	if len(kt) != 1 {
		t.Fatalf("kasus komite close %v, mau satu", kt)
	}
	komite := kt[0]
	if a := u.g.Komite[komite]; len(a) != 1 || a[0].OperatorID != wbDeptHeadUji || a[0].Jabatan != models.JabatanTutupTanpaBayar {
		t.Fatalf("tangga komite close %+v", a)
	}
	if u.g.KomiteKronologi[komite] != "UJI-KRONOLOGI" {
		t.Fatalf("Chronology %q", u.g.KomiteKronologi[komite])
	}
	if u.g.Kasus[id].Tertutup() {
		t.Fatal("klaim tertutup sebelum komite memutus")
	}
	h := u.g.Halaman(id)
	if h.Ambil(models.CD+"Remark_Close") != "UJI-CWP" {
		t.Fatalf("Remark_Close %q", h.Ambil(models.CD+"Remark_Close"))
	}
	if !riwayatMemuat(h, models.TeksMintaTutupTanpaBayar+komite) {
		t.Fatal("riwayat Request close claim without payment tidak tercatat")
	}

	// penyimpangan sadar: kiriman kedua selama kasus komite close menunggu ditolak
	kode, out = kerja("SendCloseClaimToKomite", cwp)
	u.wajib(kode, http.StatusUnprocessableEntity, out, "CWP kedua")
	if !strings.Contains(strings.Join(teks(out["pesan"]), ";"), models.PesanTutupTanpaBayarBerjalan) {
		t.Fatalf("pesan CWP kedua: %v", out["pesan"])
	}
	if n := len(komiteTutup(u, id)); n != 1 {
		t.Fatalf("kasus komite close %d sesudah kiriman kedua", n)
	}

	// Close Claim langsung: klaim Resolved-Completed, kasus komite terbuka ikut ditutup (CloseAllSubCases true)
	kode, out = kerja("CloseClaimProp", map[string]string{models.PropCentangTutupTanpaBayar: "false",
		"TempCommiteClaim.Remarks": "UJI-TUTUP"})
	u.wajib(kode, http.StatusOK, out, "tutup klaim")
	if !u.g.Kasus[id].Tertutup() || !u.g.Kasus[komite].Tertutup() {
		t.Fatalf("klaim %v / komite %v belum tertutup", u.g.Kasus[id].StatusWork, u.g.Kasus[komite].StatusWork)
	}
}

func TestTutupTanpaBayarTanpaRosterDeptHead(t *testing.T) {
	u := baruUji(t)
	u.a.Roster = []tiruan.AnggotaRoster{{AnggotaKomite: models.AnggotaKomite{ID: "2", OperatorID: "UJI-WB-TECHDIV",
		Jabatan: "Technic Div. Head", Degree: "2"}, Batas: "-9999999999999", Sts: models.STSKlaimProp}}
	id := sampaiAdjustment(u)
	kode, out := u.aksi(id, teknik, models.WorkbasketAcceptation, "DeleteAjsutment", 1, "", nil)
	u.wajib(kode, http.StatusOK, out, "hapus adjustment")
	kode, out = u.aksi(id, teknik, models.WorkbasketAcceptation, "SendCloseClaimToKomite", 0, "", map[string]string{
		models.PropCentangTutupTanpaBayar: "true", "TempCommiteClaim.Remarks": "UJI-CWP"})
	u.wajib(kode, http.StatusUnprocessableEntity, out, "CWP tanpa Claim Dept. Head di roster")
	if n := len(komiteTutup(u, id)); n != 0 {
		t.Fatalf("kasus komite close lahir tanpa penyetuju: %d", n)
	}
}
