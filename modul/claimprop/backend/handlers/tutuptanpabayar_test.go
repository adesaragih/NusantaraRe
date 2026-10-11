package handlers_test

// Close Without Payment (`SendCloseClaimToKomite`, perintah work owner 10-10-2026): tombol "Yes" pop-up
// PreventRejectClaimProp bercentang melahirkan kasus komite TKMT- satu tingkat tanpa adjustment.

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"time"

	"nusantarare/modul/claimprop/backend/handlers"
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
	// 5.5 Komite.CARI1: layar sesudah aksi membawa ID komite (pop-up berganti pesan sukses) dan centangnya
	nilai, _ := out["halaman"].(map[string]any)["nilai"].(map[string]any)
	if nilai[models.PropKomiteBaru] != komite || nilai[models.PropCentangTutupTanpaBayar] != "true" {
		t.Fatalf("Komite.CARI1 %v / centang %v di layar sesudah CWP", nilai[models.PropKomiteBaru],
			nilai[models.PropCentangTutupTanpaBayar])
	}
	if b, _ := json.Marshal(out["modal"].(map[string]any)["tutupKlaim"]); !strings.Contains(string(b),
		models.JudulSuksesTutup) {
		t.Fatalf("pop-up tanpa pesan sukses: %s", b)
	}
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

	// penyimpangan sadar (perintah work owner 10-10-2026): selama kasus komite close menunggu, CWP / Close Claim
	// ditolak dengan nomornya, dan pop-up menampilkan kartu "already submitted"
	sudah := models.PesanSudahDiajukan(models.TransferTutup, komite)
	kode, out = u.minta(http.MethodGet, handlers.Prefix+"/kasus/"+id, teknik, models.WorkbasketAcceptation, nil)
	u.wajib(kode, http.StatusOK, out, "buka kasus ber-CWP menunggu")
	if n := nilaiLayar(out); n[models.PropKomiteMenunggu] != komite || n[models.PropJenisKomiteMenunggu] != models.TransferTutup {
		t.Fatalf("penanda pengajuan menunggu: %v / %v", n[models.PropKomiteMenunggu], n[models.PropJenisKomiteMenunggu])
	}
	if b, _ := json.Marshal(out["modal"].(map[string]any)["tutupKlaim"]); !strings.Contains(string(b), models.JudulMenungguTutup) {
		t.Fatalf("pop-up Close Claim tanpa kartu already submitted: %s", b)
	}
	kode, out = kerja("SendCloseClaimToKomite", cwp)
	u.wajib(kode, http.StatusUnprocessableEntity, out, "CWP kedua")
	if strings.Join(teks(out["pesan"]), ";") != sudah {
		t.Fatalf("pesan CWP kedua: %v", out["pesan"])
	}
	if n := len(komiteTutup(u, id)); n != 1 {
		t.Fatalf("kasus komite close %d sesudah kiriman kedua", n)
	}

	tutup := map[string]string{models.PropCentangTutupTanpaBayar: "false", "TempCommiteClaim.Remarks": "UJI-TUTUP"}
	kode, out = kerja("CloseClaimProp", tutup)
	u.wajib(kode, http.StatusUnprocessableEntity, out, "Close Claim selama CWP menunggu")
	if strings.Join(teks(out["pesan"]), ";") != sudah || u.g.Kasus[id].Tertutup() {
		t.Fatalf("Close Claim selama CWP menunggu: %v", out["pesan"])
	}

	// komite sudah memutus (kasus komite selesai): Close Claim terbuka lagi
	if err := u.g.TutupKomiteTerbuka(context.Background(), nil, id, "", time.Now()); err != nil {
		t.Fatal(err)
	}
	kode, out = kerja("CloseClaimProp", tutup)
	u.wajib(kode, http.StatusOK, out, "tutup klaim")
	if !u.g.Kasus[id].Tertutup() {
		t.Fatalf("klaim %v belum tertutup", u.g.Kasus[id].StatusWork)
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
