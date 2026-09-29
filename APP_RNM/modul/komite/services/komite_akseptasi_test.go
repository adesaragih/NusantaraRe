package services

// Nomor akseptasi Komite - tiket 04a. TANPA Oracle.

import (
	"nusantarare/inti"
	"nusantarare/modul/komite/repository"
	"os"
	"strings"
	"testing"
	"time"
)

func sumberKomiteAkseptasi(t *testing.T) string {
	t.Helper()
	isi, err := os.ReadFile("komite_akseptasi.go")
	if err != nil {
		t.Fatal(err)
	}
	var kode []string
	for _, b := range strings.Split(string(isi), "\n") {
		if !strings.HasPrefix(strings.TrimSpace(b), "//") {
			kode = append(kode, b)
		}
	}
	return strings.Join(kode, "\n")
}

// TestPenghitungKomiteBukanSequenceClaimLife - AC 14 spec, keputusan o.
func TestPenghitungKomiteBukanSequenceClaimLife(t *testing.T) {
	kode := sumberKomiteAkseptasi(t)
	for _, wajib := range []string{"penghitung.AwalanProduksi(", "penghitung.UrutNomorBerikut(",
		"models.ClassPenghitungKomiteLife", "models.JenisPenghitungKomite(", "models.NomorAkseptasiKomite("} {
		if !strings.Contains(kode, wajib) {
			t.Errorf("jalur nomor Komite tidak memakai %s", wajib)
		}
	}
	for _, terlarang := range []string{"UrutAkseptasiBerikut(", "ACCEPTATIONNOLIFE_SEQ",
		"Generate_NoAccept_KMT", "RakitNomorAkseptasi(", `"RNML-`} {
		if strings.Contains(kode, terlarang) {
			t.Errorf("jalur nomor Komite memuat %q - penghitung/awalan milik jalur lain atau literal", terlarang)
		}
	}
}

// TestKeunikanDiperiksaSebelumStempel - OQ-K-04a, tabrakan gagal terang.
func TestKeunikanDiperiksaSebelumStempel(t *testing.T) {
	kode := sumberKomiteAkseptasi(t)
	iLama := strings.Index(kode, ".NomorAkseptasiDipakai(ctx, tx, nomor)")
	iAdj := strings.Index(kode, ".NomorAkseptasiDipakaiDiAdjustment(ctx, tx, nomor)")
	iStempel := strings.Index(kode, "baca.PerbaruiStatusBaris(")
	if iLama < 0 || iAdj < 0 || iStempel < 0 {
		t.Fatal("pemeriksa keunikan atau stempel hilang")
	}
	// Pemeriksa ada di terbitkanNomor, yang dipanggil SEBELUM stempel.
	iPanggil := strings.Index(kode, "p.terbitkanNomor(ctx, tx,")
	if iPanggil < 0 || iPanggil > iStempel {
		t.Error("nomor tidak diterbitkan (dan diperiksa) sebelum baris distempel")
	}
}

// TestRekamAkhirSatuJalurDalamAkseptasi - tiket 04b, AC 16 spec.
//
// Refactor bentuk B (30-09-2026): kode status kini `kontrak.Kode*` (kosakata
// dibagi Claim Life dan Komite), dulu `models.Kode*`.
func TestRekamAkhirSatuJalurDalamAkseptasi(t *testing.T) {
	kode := sumberKomiteAkseptasi(t)
	iHeader := strings.Index(kode, "baca.CerminkanHeader(")
	iRekam := strings.Index(kode, "p.rekamAkhir(ctx, tx, kasus, kontrak.KodeAksep, nomor, saat)")
	iJejak := strings.Index(kode, "p.jejak.Rekam(")
	if iHeader < 0 || iRekam < 0 || iJejak < 0 || !(iHeader < iRekam && iRekam < iJejak) {
		t.Error("urutan akseptasi harus: stempel → rekam akhir → jejak")
	}
	if n := strings.Count(kode, "RekamAkhirWarisan("); n != 1 {
		t.Errorf("rekam akhir dipanggil di %d tempat, mau 1 (satu jalur berparameter status)", n)
	}
}

// TestTolakAkhirDuaTingkatBarisSatuJalur - tiket 05 Komite, AC 15 spec.
func TestTolakAkhirDuaTingkatBarisSatuJalur(t *testing.T) {
	kode := sumberKomiteAkseptasi(t)
	i := strings.Index(kode, "func (p penyelesaiAkhirOracle) Tolak(")
	if i < 0 {
		t.Fatal("Tolak tidak ditemukan")
	}
	badan := kode[i:]
	for _, wajib := range []string{"periksaGiliran(kasus, pelaku.AkunID)",
		"kontrak.KodeOutstanding, kontrak.KodeDitolak", "baca.CabutPenandaDipilih(",
		"baca.CerminkanHeader(ctx, tx, klaimID, kontrak.KodeDitolak", "p.rekamAkhir(ctx, tx, kasus, kontrak.KodeDitolak",
		"p.jejak.Rekam("} {
		if !strings.Contains(badan, wajib) {
			t.Errorf("Tolak tidak memuat %q", wajib)
		}
	}
	if strings.Contains(badan, "ErrPenyelesaianAkhirBelumAda") {
		t.Error("Tolak masih gagal terang - tiket 05 belum tersambung")
	}
	// ⛔ OQ-K-05 DITUTUP (GILIRAN-17): 5.1 DITIRU - tangga dibaca di dalam
	// transaksi, keputusan lama tiap tingkat dijejaki, lalu SATU UPDATE
	// bersyarat menimpanya; semuanya SEBELUM 5.3. SQL-nya tinggal di
	// repository (penjaga komite_statik), bukan di sini.
	urut := []string{"komite.TanggaSebelumDitimpa(ctx, tx, kasus.Baris.KasusID)", "jejakTimpaTangga(",
		"komite.TimpaTanggaTolakAkhir(ctx, tx, kasus.Baris.KasusID, saat)", "baca.PerbaruiStatusBaris("}
	lalu := -1
	for _, s := range urut {
		j := strings.Index(badan, s)
		if j < 0 || j < lalu {
			t.Errorf("urutan 5.1 salah di %q (urutan: %v)", s, urut)
		}
		lalu = j
	}
	if strings.Contains(kode, "KOMITE_APPROVAL = '2'") {
		t.Error("SQL penimpaan tangga ditulis di services; tempatnya repository")
	}
}

// OQ-K-05 (GILIRAN-17): jejak "ditimpa" hanya untuk tingkat yang MEMUTUS -
// tingkat yang dilewati eskalasi (NULL) tidak pernah memberi keputusan
// (tiket 03), dan tingkat yang menunggu tidak ada saat tingkat akhir menolak.
func TestJejakTimpaTanggaHanyaTingkatBerkeputusan(t *testing.T) {
	saat := time.Date(2026, 9, 29, 10, 0, 0, 0, time.UTC)
	tangga := []repository.AnggotaKasus{
		{Urut: 1, OperatorID: "UJI-A", Approval: "1", Komentar: "UJI setuju"},
		{Urut: 2, OperatorID: "UJI-B", Approval: ""},
		{Urut: 3, OperatorID: "UJI-C", Approval: "2", Komentar: "UJI tolak akhir"},
	}
	kasus := repository.KasusKomite{AdjID: "UJI-ADJ", Baris: repository.BarisInboxKomite{KasusID: "KMTLF-UJI", KlaimID: "UJI-K"}}
	j := jejakTimpaTangga(kasus, tangga, inti.Pelaku{AkunID: "UJI-C"}, saat)
	if len(j) != 2 {
		t.Fatalf("jejak %d, mau 2 (tingkat 1 dan 3): %+v", len(j), j)
	}
	if j[0].Dari != awalanJejakTingkat+"1" || j[0].Ke != awalanJejakTimpa+"Setuju (KMTLF-UJI)" || j[0].Komentar != "UJI setuju" {
		t.Errorf("tingkat 1: %+v", j[0])
	}
	if j[1].Dari != awalanJejakTingkat+"3" || j[1].Komentar != "UJI tolak akhir" || j[1].AdjustmentID != "UJI-ADJ" {
		t.Errorf("tingkat 3: %+v", j[1])
	}
	for _, c := range j {
		if len(c.Ke) > 64 || len(c.Dari) > 64 {
			t.Errorf("DARI/KE melebihi VARCHAR2(64): %+v", c)
		}
	}
}
