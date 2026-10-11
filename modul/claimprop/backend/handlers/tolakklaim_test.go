package handlers_test

// Reject Claim (perintah work owner 10-10-2026): tombol di posisi admin (Outstanding) selama Claim No belum ada;
// sesudah Claim No ada diganti Close Claim posisi teknik. ClaimRejectTreaty -> SureRejectClaimTreaty ->
// `CheckRejectClaimTreaty2` melahirkan kasus komite TKMT- TT 3 satu tingkat tanpa adjustment.

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"nusantarare/modul/claimprop/backend/handlers"
	"nusantarare/modul/claimprop/backend/models"
)

// tombolLayar - tombol berID `id` di tata utama jawaban (nil = tidak tampil).
func tombolLayar(t *testing.T, out map[string]any, id string) *models.Tata {
	t.Helper()
	b, _ := json.Marshal(out["tata"])
	var ts []models.Tata
	if err := json.Unmarshal(b, &ts); err != nil {
		t.Fatal(err)
	}
	var cari func([]models.Tata) *models.Tata
	cari = func(xs []models.Tata) *models.Tata {
		for i := range xs {
			if xs[i].Jenis == models.JenisTombol && xs[i].ID == id {
				return &xs[i]
			}
			if x := cari(xs[i].Anak); x != nil {
				return x
			}
		}
		return nil
	}
	return cari(ts)
}

func nilaiLayar(out map[string]any) map[string]any {
	n, _ := out["halaman"].(map[string]any)["nilai"].(map[string]any)
	return n
}

func TestRejectClaimPosisiAdmin(t *testing.T) {
	u := baruUji(t)
	rosterTutupUji(u)
	id := u.buat()
	kode, out := u.minta(http.MethodGet, handlers.Prefix+"/kasus/"+id, admin, "", nil)
	u.wajib(kode, http.StatusOK, out, "buka kasus")
	if b := tombolLayar(t, out, "RejectClaim"); b == nil || b.Nonaktif || tombolLayar(t, out, "CloseClaim") != nil {
		t.Fatalf("Claim No kosong: Reject Claim %+v, Close Claim %+v", b, tombolLayar(t, out, "CloseClaim"))
	}
	if n := nilaiLayar(out); n[models.JalurTanggalTolak] == "" || n[models.JalurPICTolak] == "" {
		t.Fatalf("Date Reject / PIC Name kosong: %v / %v", n[models.JalurTanggalTolak], n[models.JalurPICTolak])
	}
	kerja := func(aksi string, m map[string]string) (int, map[string]any) {
		return u.aksi(id, admin, "", aksi, 0, "", m)
	}

	// Submit ClaimRejectTreaty: Remarks wajib; isian yang diketik dibawa ke layar konfirmasi
	kode, out = kerja("SureRejectClaimT", nil)
	u.wajib(kode, http.StatusUnprocessableEntity, out, "Submit tanpa Remarks")
	isian := map[string]string{"TempCommiteClaim.Remarks": "UJI-TOLAK", models.PropKronologiTutup: "UJI-KRONOLOGI"}
	kode, out = kerja("SureRejectClaimT", isian)
	u.wajib(kode, http.StatusOK, out, "Submit ClaimRejectTreaty")
	if n := nilaiLayar(out); n["TempCommiteClaim.Remarks"] != "UJI-TOLAK" || n[models.PropKronologiTutup] != "UJI-KRONOLOGI" {
		t.Fatalf("isian pop-up hilang sesudah Submit: %v", n)
	}

	// CheckRejectClaimTreaty2 4-6: pesan TERAKHIR yang tampil. IsAnyAcceptation diturunkan dari baris adjustment
	// ber-AcceptanceStatus 1 (CheckAnyAcceptationProp).
	h := u.g.Halaman(id)
	h.SetelDaftar(models.DaftarAdjustment, []models.Baris{{"ID": "UJI-ADJ-1", "AcceptanceStatus": "1"}})
	h.SetelDaftar(models.DaftarEstimasi, []models.Baris{{"PrintFaceClaim": "1"}, {"PrintFaceClaim": ""}})
	u.g.SetelHalaman(id, h)
	kode, out = kerja("CheckRejectClaimTreaty2", isian)
	u.wajib(kode, http.StatusUnprocessableEntity, out, "estimasi belum terkirim")
	if p := strings.Join(teks(out["pesan"]), ";"); p != models.PesanHapusEstimasi(2) {
		t.Fatalf("pesan %q", p)
	}
	h.SetelDaftar(models.DaftarEstimasi, []models.Baris{{"PrintFaceClaim": "1"}})
	u.g.SetelHalaman(id, h)
	kode, out = kerja("CheckRejectClaimTreaty2", isian)
	u.wajib(kode, http.StatusUnprocessableEntity, out, "sudah ada akseptasi")
	if p := strings.Join(teks(out["pesan"]), ";"); p != models.PesanTolakSudahAkseptasi {
		t.Fatalf("pesan %q", p)
	}
	if n := len(komiteTutup(u, id)); n != 0 {
		t.Fatalf("kasus komite lahir walau ditolak: %d", n)
	}
	h.SetelDaftar(models.DaftarAdjustment, nil)
	u.g.SetelHalaman(id, h)

	kode, out = kerja("CheckRejectClaimTreaty2", isian)
	u.wajib(kode, http.StatusOK, out, "Reject Claim")
	kt := komiteTutup(u, id)
	if len(kt) != 1 {
		t.Fatalf("kasus komite reject %v, mau satu", kt)
	}
	komite := kt[0]
	if u.g.KomiteTransfer[komite] != models.TransferTolak || u.g.KomiteKronologi[komite] != "UJI-KRONOLOGI" {
		t.Fatalf("TRANSFER_TYPE %q, Chronology %q", u.g.KomiteTransfer[komite], u.g.KomiteKronologi[komite])
	}
	if a := u.g.Komite[komite]; len(a) != 1 || a[0].OperatorID != wbDeptHeadUji {
		t.Fatalf("tangga komite reject %+v", a)
	}
	if nilaiLayar(out)[models.PropKomiteBaru] != komite {
		t.Fatalf("Komite.CARI1 %v", nilaiLayar(out)[models.PropKomiteBaru])
	}
	if u.g.Kasus[id].Tertutup() {
		t.Fatal("klaim tertutup sebelum komite memutus")
	}
	h = u.g.Halaman(id)
	if h.Ambil(models.CD+"Remark_Close") != "UJI-TOLAK" || !riwayatMemuat(h, models.TeksMintaTolak+komite) {
		t.Fatalf("Remark_Close %q / riwayat Request reject claim", h.Ambil(models.CD+"Remark_Close"))
	}
	// penyimpangan sadar (perintah work owner 10-10-2026): selama kasus komite reject menunggu, pop-up Reject Claim
	// hanya menampilkan kartu "already submitted" bernomor, dan Submit / Yes ditolak dengan nomornya
	sudah := models.PesanSudahDiajukan(models.TransferTolak, komite)
	kode, out = u.minta(http.MethodGet, handlers.Prefix+"/kasus/"+id, admin, "", nil)
	u.wajib(kode, http.StatusOK, out, "buka kasus ber-reject menunggu")
	if n := nilaiLayar(out); n[models.PropKomiteMenunggu] != komite || n[models.PropJenisKomiteMenunggu] != models.TransferTolak {
		t.Fatalf("penanda pengajuan menunggu: %v / %v", n[models.PropKomiteMenunggu], n[models.PropJenisKomiteMenunggu])
	}
	if b, _ := json.Marshal(out["modal"].(map[string]any)["tolakKlaim"]); !strings.Contains(string(b), models.JudulMenungguTolak) {
		t.Fatalf("pop-up Reject Claim tanpa kartu already submitted: %s", b)
	}
	for _, aksi := range []string{"SureRejectClaimT", "CheckRejectClaimTreaty2"} {
		kode, out = kerja(aksi, isian)
		u.wajib(kode, http.StatusUnprocessableEntity, out, aksi+" kedua")
		if p := strings.Join(teks(out["pesan"]), ";"); p != sudah {
			t.Fatalf("pesan %s kedua %q", aksi, p)
		}
	}

	// sesudah Claim No ada: Reject Claim diganti Close Claim; aksi reject tertutup
	h = u.g.Halaman(id)
	h.Setel(models.CD+"NoClaim", "UJI-K-0001")
	u.g.SetelHalaman(id, h)
	kode, out = u.minta(http.MethodGet, handlers.Prefix+"/kasus/"+id, admin, "", nil)
	u.wajib(kode, http.StatusOK, out, "buka kasus ber-Claim No")
	if tombolLayar(t, out, "RejectClaim") != nil || tombolLayar(t, out, "CloseClaim") == nil {
		t.Fatal("Claim No ada: Reject Claim harus hilang, Close Claim tampil")
	}
	if kode, out = kerja("SureRejectClaimT", isian); kode == http.StatusOK {
		t.Fatalf("Submit reject lolos sesudah Claim No ada: %v", out)
	}
}

// Close Claim posisi admin (Claim No ada) = pop-up PreventRejectClaimProp posisi teknik: klaim Resolved-Completed.
func TestCloseClaimPosisiAdminSesudahNomor(t *testing.T) {
	u := baruUji(t)
	id := u.buat()
	h := u.g.Halaman(id)
	h.Setel(models.CD+"NoClaim", "UJI-K-0001")
	u.g.SetelHalaman(id, h)
	kode, out := u.aksi(id, admin, "", "CloseClaimProp", 0, "", map[string]string{"TempCommiteClaim.Remarks": "UJI"})
	u.wajib(kode, http.StatusOK, out, "Close Claim posisi admin")
	if !u.g.Kasus[id].Tertutup() {
		t.Fatal("klaim belum tertutup")
	}
}

// Hanya pemegang (pembuat, posisi admin) yang boleh reject.
func TestRejectClaimBukanPemegang(t *testing.T) {
	u := baruUji(t)
	rosterTutupUji(u)
	id := u.buat()
	kode, out := u.aksi(id, lain, "", "CheckRejectClaimTreaty2", 0, "", map[string]string{"TempCommiteClaim.Remarks": "UJI"})
	if kode == http.StatusOK || len(komiteTutup(u, id)) != 0 {
		t.Fatalf("bukan pemegang lolos: %d %v", kode, out)
	}
}
