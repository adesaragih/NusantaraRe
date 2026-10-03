package models

// Uji penulis laporan pemuat dokumen lama (tiket 22): berkas CSV medan tak
// dikenal (K17: POLIS_ID, JALUR, NILAI - bukan tabel), berkas CSV galat
// (AC 58), dan ringkasan yang dicetak (K15, AC 59).

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

func TestLaporanMedanTakDikenalBerkasCSV(t *testing.T) { // K17, AC 57
	var tak, gal bytes.Buffer
	l, err := LaporanPemuatBaru(&tak, &gal, true)
	if err != nil {
		t.Fatal(err)
	}
	d := strings.Replace(dokumenUjiProp, `"UJIMedanFiktif": "UJI-nilai-lewat"`,
		`"UJIMedanFiktif": "UJI-a, \"b\"\nbaris kedua"`, 1)
	h, err := PecahDokumenLama(barisUji(d))
	if err != nil {
		t.Fatal(err)
	}
	if err := l.Berhasil(h); err != nil {
		t.Fatal(err)
	}
	if err := l.Tutup(); err != nil {
		t.Fatal(err)
	}
	baris := bacaCSV(t, &tak)
	harap := [][]string{
		{"POLIS_ID", "JALUR", "NILAI"},
		{"UJI-77", "PolicyTreatyIn.QuotationData.UJIFiktifQuotation", ""},
		{"UJI-77", "PolicyTreatyIn.UJIMedanFiktif", "UJI-a, \"b\"\nbaris kedua"}, // nilai utuh, tidak dibuang
	}
	if len(baris) != len(harap) {
		t.Fatalf("CSV %q", baris)
	}
	for i := range harap {
		if strings.Join(baris[i], "|") != strings.Join(harap[i], "|") {
			t.Errorf("baris %d = %q, harap %q", i, baris[i], harap[i])
		}
	}
	r := l.Ringkasan()
	if r.Dimuat != 1 || r.MedanTakDikenal != 2 || r.Selesai() {
		t.Errorf("ringkasan %+v; selesai harus false selama medan tak dikenal > 0 (AC 59)", r)
	}
	if !strings.Contains(r.Teks(), "Medan tak dikenal") || !strings.Contains(r.Teks(), "BELUM SELESAI") {
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
		t.Error("dokumen gagal tidak menulis medan tak dikenal")
	}
	r := l.Ringkasan()
	if r.DokumenGagal != 2 || r.TanggalAmbigu != 1 || r.BukanTreatyIn != 1 || r.Selesai() {
		t.Errorf("ringkasan %+v", r)
	}
	if r.GalatPerJenis[JenisGalat(ErrTanggalAmbigu)] != 1 || r.GalatPerJenis[JenisGalat(ErrDokumenRusak)] != 1 {
		t.Errorf("galat per jenis %v", r.GalatPerJenis)
	}
}

func TestRingkasanSelesaiHanyaBilaNolGalatDanNolTakDikenal(t *testing.T) { // AC 59
	var tak, gal bytes.Buffer
	l, _ := LaporanPemuatBaru(&tak, &gal, true)
	d := strings.Replace(strings.Replace(dokumenUjiProp, `"UJIMedanFiktif": "UJI-nilai-lewat",`, ``, 1),
		`"UJIFiktifQuotation": "",`, ``, 1)
	h, err := PecahDokumenLama(barisUji(d))
	if err != nil || len(h.TakDikenal) != 0 {
		t.Fatalf("%v %+v", err, h.TakDikenal)
	}
	_ = l.Berhasil(h)
	l.SudahDimuat()
	_ = l.Tutup()
	r := l.Ringkasan()
	// pyID fixture `UJI-77` bukan ruang nomor `AwalanKasus`: SEQ_WORK_POLIS tidak tersentuh.
	if !r.Selesai() || r.AngkaKasusTerbesar != 0 || r.SudahDimuat != 1 {
		t.Errorf("ringkasan %+v", r)
	}
	if strings.Contains(r.Teks(), "Nomor kasus terbesar") {
		t.Errorf("pyID di luar awalan %s tidak dilaporkan:\n%s", AwalanKasus, r.Teks())
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
		_ = l.Berhasil(h)
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
