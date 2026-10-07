package models

// Uji penulis laporan pemuat dokumen lama (tiket 22): berkas CSV ARSIP medan
// tanpa kolom (F3, WO 04-10-2026: arsip audit pemuatan, bukan penampung -
// POLIS_ID, JALUR, NILAI, KEPUTUSAN), berkas CSV galat (AC 58), dan
// ringkasan yang dicetak (K15; AC 59 = nol medan BELUM DIPUTUSKAN).

import (
	"bytes"
	"encoding/csv"
	"errors"
	"strings"
	"testing"
)

func bacaCSV(t *testing.T, b *bytes.Buffer) [][]string {
	t.Helper()
	baris, err := csv.NewReader(bytes.NewReader(b.Bytes())).ReadAll()
	if err != nil {
		t.Fatal(err)
	}
	return baris
}

// hitungKeputusan - cacah baris data arsip per nilai KEPUTUSAN.
func hitungKeputusan(baris [][]string) map[string]int {
	m := map[string]int{}
	for _, b := range baris[1:] {
		m[b[3]]++
	}
	return m
}

func TestLaporanArsipMedanTanpaKolomBerkasCSV(t *testing.T) { // F3; K17, AC 57, 59
	var arsip, gal bytes.Buffer
	l, err := LaporanPemuatBaru(&arsip, &gal, true)
	if err != nil {
		t.Fatal(err)
	}
	d := strings.Replace(dokumenUjiProp, `"UJIMedanFiktif": "UJI-nilai-lewat"`,
		`"UJIMedanFiktif": "UJI-a, \"b\"\nbaris kedua"`, 1)
	h, err := PecahDokumenLama(barisUji(d))
	if err != nil {
		t.Fatal(err)
	}
	if err := l.Berhasil(h, UsulanDisalin); err != nil {
		t.Fatal(err)
	}
	if err := l.Tutup(); err != nil {
		t.Fatal(err)
	}
	baris := bacaCSV(t, &arsip)
	if strings.Join(baris[0], ",") != "POLIS_ID,JALUR,NILAI,KEPUTUSAN" {
		t.Fatalf("kepala arsip %q", baris[0])
	}
	// Setiap medan daun yang TIDAK masuk kolom tercatat beserta nilainya dan
	// keputusannya - yang dibuang menurut keputusan tertulis juga (arsip audit,
	// tidak hilang diam-diam). Fixture: pxObjClass 4 simpul, satu nomor baris
	// Pega, Show, TotalPremium; dua medan fiktif tanpa keputusan.
	harapCacah := map[string]int{
		"dibuang: internal_pega": 4, "dibuang: nourut": 1, "dibuang: keadaan_layar": 1, "dibuang: turunan": 1,
		KeputusanBelumDiputuskan: 2,
	}
	cacah := hitungKeputusan(baris)
	if len(cacah) != len(harapCacah) {
		t.Errorf("keputusan arsip %v, harap %v", cacah, harapCacah)
	}
	for k, n := range harapCacah {
		if cacah[k] != n {
			t.Errorf("KEPUTUSAN %q: %d baris, harap %d", k, cacah[k], n)
		}
	}
	ada := map[string]bool{}
	for _, b := range baris[1:] {
		ada[strings.Join(b, "|")] = true
	}
	for _, b := range [][]string{
		{"UJI-77", "PolicyTreatyIn.UJIMedanFiktif", "UJI-a, \"b\"\nbaris kedua", KeputusanBelumDiputuskan}, // nilai utuh
		{"UJI-77", "PolicyTreatyIn.QuotationData.UJIFiktifQuotation", "", KeputusanBelumDiputuskan},
		{"UJI-77", "PolicyTreatyIn.Show", "true", "dibuang: keadaan_layar"},
		{"UJI-77", "PolicyTreatyIn.TotalPremium", "592629512.880000276", "dibuang: turunan"},
	} {
		if !ada[strings.Join(b, "|")] {
			t.Errorf("baris arsip %q tidak ada di %q", b, baris)
		}
	}
	r := l.Ringkasan()
	if r.Dimuat != 1 || r.MedanBelumDiputuskan != 2 || r.Selesai() {
		t.Errorf("ringkasan %+v; selesai harus false selama medan belum diputuskan > 0 (AC 59)", r)
	}
	if r.DiabaikanPerAlasan[AlasanInternalPega] != 4 || r.DiabaikanPerAlasan[AlasanTurunan] != 1 {
		t.Errorf("dibuang per alasan %v", r.DiabaikanPerAlasan)
	}
	if !strings.Contains(r.Teks(), "BELUM DIPUTUSKAN") || !strings.Contains(r.Teks(), "BELUM SELESAI") {
		t.Errorf("teks ringkasan:\n%s", r.Teks())
	}
}

func TestLaporanGalatBerkasCSVDanTanggalAmbiguDihitung(t *testing.T) { // AC 58, K15
	var tak, gal bytes.Buffer
	l, err := LaporanPemuatBaru(&tak, &gal, false)
	if err != nil {
		t.Fatal(err)
	}
	d := strings.Replace(dokumenUjiProp, `"StartDate": "20171001"`, `"StartDate": "05/06/2017"`, 1)
	b := barisUji(d)
	h, err := PecahDokumenLama(b)
	if err != nil {
		t.Fatal(err)
	}
	if err := l.Gagal(b, h.Galat); err != nil {
		t.Fatal(err)
	}
	if err := l.Gagal(barisUji(`{`), []GalatDokumen{{Err: ErrDokumenRusak}}); err != nil {
		t.Fatal(err)
	}
	l.Lewat(ErrBukanTreatyIn)
	l.Lewat(ErrBarisAplikasiBaru)
	if r := l.Ringkasan(); r.BarisAplikasiBaru != 1 || r.BukanTreatyIn != 1 || r.DokumenGagal != 2 {
		t.Fatalf("baris aplikasi baru dihitung terpisah, bukan galat: %+v", r)
	}
	if err := l.Tutup(); err != nil {
		t.Fatal(err)
	}
	baris := bacaCSV(t, &gal)
	if len(baris) != 3 || strings.Join(baris[0], ",") != "IDPEGA,NOPOLIS,JALUR,NILAI,SEBAB" {
		t.Fatalf("CSV galat %q", baris)
	}
	if baris[1][0] != "ASM-FW-GISFW-WORK-NB UJI-77" || baris[1][2] != "PolicyTreatyIn.StartDate" ||
		baris[1][3] != "05/06/2017" || !strings.Contains(baris[1][4], "ambigu") {
		t.Errorf("baris galat tanggal %q", baris[1])
	}
	if len(bacaCSV(t, &tak)) != 1 {
		t.Error("dokumen gagal tidak menulis arsip medan")
	}
	r := l.Ringkasan()
	if r.DokumenGagal != 2 || r.TanggalAmbigu != 1 || r.BukanTreatyIn != 1 || r.Selesai() {
		t.Errorf("ringkasan %+v", r)
	}
	if r.GalatPerJenis[JenisGalat(ErrTanggalAmbigu)] != 1 || r.GalatPerJenis[JenisGalat(ErrDokumenRusak)] != 1 {
		t.Errorf("galat per jenis %v", r.GalatPerJenis)
	}
}

// F3: kode keluar pemuat bukan nol HANYA bila masih ada medan belum
// diputuskan (atau galat) - medan yang dibuang menurut keputusan tertulis
// tetap diarsipkan tanpa menahan selesai.
func TestRingkasanSelesaiHanyaBilaNolGalatDanNolBelumDiputuskan(t *testing.T) { // AC 59
	var arsip, gal bytes.Buffer
	l, _ := LaporanPemuatBaru(&arsip, &gal, true)
	d := strings.Replace(strings.Replace(dokumenUjiProp, `"UJIMedanFiktif": "UJI-nilai-lewat",`, ``, 1),
		`"UJIFiktifQuotation": "",`, ``, 1)
	h, err := PecahDokumenLama(barisUji(d))
	if err != nil || len(h.BelumDiputuskan) != 0 {
		t.Fatalf("%v %+v", err, h.BelumDiputuskan)
	}
	_ = l.Berhasil(h, UsulanDisalin)
	l.SudahDimuat()
	_ = l.Tutup()
	r := l.Ringkasan()
	// pyID fixture `UJI-77` bukan ruang nomor `AwalanKasus`: SEQ_WORK_POLIS tidak tersentuh.
	if !r.Selesai() || r.AngkaKasusTerbesar != 0 || r.SudahDimuat != 1 {
		t.Errorf("ringkasan %+v", r)
	}
	if n := len(bacaCSV(t, &arsip)) - 1; n != 7 {
		t.Errorf("arsip tetap memuat 7 medan yang dibuang, dapat %d", n)
	}
	if !strings.Contains(r.Teks(), "STATUS: SELESAI") {
		t.Errorf("teks ringkasan:\n%s", r.Teks())
	}
	if strings.Contains(r.Teks(), "Nomor kasus terbesar") {
		t.Errorf("pyID di luar awalan %s tidak dilaporkan:\n%s", AwalanKasus, r.Teks())
	}
}

// F3: SuggestList dokumen lama disalin ke HISTORYAKSEPTASIPRODUCTION; PIC kosong
// ditulis apa adanya (NULL) dan DICATAT; AKSES_LOGIN selalu NULL (tanpa anggota
// sumber); dokumen yang IDPEGA-nya sudah punya baris riwayat produksi dilewati
// salinannya (penjaga dobel); dokumen tanpa baris tidak dicatat.
func TestRingkasanSalinanUsulanLama(t *testing.T) {
	var arsip, gal bytes.Buffer
	l, _ := LaporanPemuatBaru(&arsip, &gal, true)
	h, err := PecahDokumenLama(barisUji(dokumenUjiUsulan))
	if err != nil || len(h.Usulan) != 2 {
		t.Fatalf("%v %+v", err, h.Usulan)
	}
	_ = l.Berhasil(h, UsulanDisalin)
	_ = l.Berhasil(h, UsulanDilewati) // jalankan berikutnya: IDPEGA sudah punya baris
	_ = l.Tutup()
	r := l.Ringkasan()
	_ = l.Berhasil(HasilPecah{}, UsulanTanpaBaris)
	if r = l.Ringkasan(); r.Usulan != (RingkasanUsulan{Disalin: 2, TanpaPIC: 1, DokumenSudahAda: 1}) {
		t.Errorf("ringkasan usulan %+v", r)
	}
	if !r.Selesai() {
		t.Errorf("AKSES_LOGIN/PIC kosong dicatat, bukan galat: %+v", r)
	}
	for _, s := range []string{"HISTORYAKSEPTASIPRODUCTION", "AKSES_LOGIN selalu NULL", "PIC kosong"} {
		if !strings.Contains(r.Teks(), s) {
			t.Errorf("teks ringkasan tanpa %q:\n%s", s, r.Teks())
		}
	}
}

// Nomor kasus terbesar = ruang nomor `AwalanKasus` (SEQ_WORK_POLIS bersama,
// models/kasus.go). pyID dirakit dari `RakitIDKasus` - perilaku yang diuji
// memang ruang nomor itu; fixture sendiri tetap `UJI-`.
func TestRingkasanMencetakNomorKasusTerbesarBerawalanKasus(t *testing.T) {
	var tak, gal bytes.Buffer
	l, _ := LaporanPemuatBaru(&tak, &gal, true)
	for _, urut := range []string{"770", "77"} {
		b := barisUji(dokumenUjiProp)
		b.IDPega = strings.ToUpper(KelasDeret) + " " + RakitIDKasus(urut)
		h, err := PecahDokumenLama(b)
		if err != nil {
			t.Fatal(err)
		}
		_ = l.Berhasil(h, UsulanDisalin)
	}
	_ = l.Tutup()
	r := l.Ringkasan()
	if r.AngkaKasusTerbesar != 770 {
		t.Errorf("angka kasus terbesar %d, harap 770", r.AngkaKasusTerbesar)
	}
	if !strings.Contains(r.Teks(), RakitIDKasus("770")) {
		t.Errorf("nomor kasus terbesar wajib dicetak (SEQ_WORK_POLIS):\n%s", r.Teks())
	}
}

func TestJenisGalatPenyimpanan(t *testing.T) {
	if JenisGalat(errors.New("UJI-ORA")) != JenisGalatSimpan {
		t.Error("galat lain = ditolak penyimpanan")
	}
	if JenisGalat(ErrIDKasusDipakai) == JenisGalatSimpan {
		t.Error("ID kasus dipakai punya jenisnya sendiri")
	}
}
