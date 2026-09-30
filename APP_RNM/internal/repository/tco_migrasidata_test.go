package repository

// Uji konversi data warisan Treaty Contract Out - TANPA Oracle (tiket 01).
//
// Fixture SINTETIS: nol nama orang, nol nomor polis nyata, nol potongan data
// produksi; seluruh nilai berawalan UJI.

import (
	"strings"
	"testing"

	"nusantarare/pkg/utils"
)

func barisUji(tabel string, nilai map[string]string) BarisWarisanTCO {
	return BarisWarisanTCO{Tabel: tabel, Nilai: nilai}
}

// contohWarisanTCO memberi satu baris wajar per tabel, saling merujuk.
func contohWarisanTCO() map[string][]BarisWarisanTCO {
	return map[string][]BarisWarisanTCO{
		warisanTahunTCO: {barisUji(warisanTahunTCO, map[string]string{
			"ID": "1000001", "TREATYYEAR": "2026", "UNDERWRITINGYEAR": "2026",
			"TREATYGROUPID": "10001", "TREATYGROUPNAME": "UJI GRUP", "USERID": "UJI-OP",
			"TGLUPDATE": "15-SEP-26", "PROPORTION": "P", "STARTDATE": "20260101", "ENDDATE": "20261231",
		})},
		warisanKontrakTCO: {barisUji(warisanKontrakTCO, map[string]string{
			"ID": "1000007", "IDTREATYYEAR": "1000001", "REINSTYPEID": "10002",
			"REINSTYPENAME": "UJI QS", "TREATYSTARTDATE": "2026-01-01 00:00:00",
			"TREATYENDDATE": "2026-12-31 00:00:00", "USERID": "UJI-OP", "TGLUPDATE": "01/02/2026",
		})},
		warisanReinsurerTCO: {barisUji(warisanReinsurerTCO, map[string]string{
			"ID": "1000003", "TREATYYEAR": "2026", "TREATYGROUPID": "10001", "REINSTYPEID": "10002",
			"REINSURERID": "R1", "CLIENTID": "C1", "NAME": "UJI REINS A", "RICOMM": "12.5",
			"PCTSHARE": "33.333", "IUDATE": "UJI-APA-ADANYA", "STARTDATE": "20260101T170000.000 GMT",
			"ENDDATE": "", "OPERATORNAME": "UJI-OP", "TGLUPDATE": "",
		})},
		warisanSecurityTCO: {barisUji(warisanSecurityTCO, map[string]string{
			"THN_TREATY": "2026", "TOP_ID": "", "TP_TREATY": "  ", "REAS_ID": "1000003",
			"PCT_SHARE": "50,5", "USER_ID": "", "REAS_SECURITY": "SEC-A     ",
		})},
		warisanBusinessTCO: {barisUji(warisanBusinessTCO, map[string]string{
			"ID": "1000002", "ISACTIVE": "1", "TREATYYEAR": "2026", "TREATYYEARID": "",
			"TREATYGROUPID": "10001", "REINSTYPEID": "10002", "BIZCODE": "UJI-BIZ", "BIZNAME": "UJI BISNIS",
		})},
		warisanKlausulTCO: {barisUji(warisanKlausulTCO, map[string]string{
			"ID": "10000009", "TREATYYEAR": "2026", "TREATYYEARID": "1000001", "TREATYGROUPID": "10001",
			"TREATYDESCID": "10001", "REINSTYPEID": "10002", "PARENTREINSTYPEID": "00",
			"PCT": "12,5", "RP": "1000000000.12345678", "USD": ".5", "TGLUPDATE": "2026-01-02 03:04:05",
		})},
	}
}

func TestUraiTanggalWarisanTCOMengenalBentukYangDikenal(t *testing.T) {
	kasus := map[string]string{
		"20260131":                "2026-01-31 00:00:00",
		"20260131T170000.000 GMT": "2026-01-31 00:00:00",
		"31/01/2026":              "2026-01-31 00:00:00",
		"2026-01-31":              "2026-01-31 00:00:00",
		"2026-01-31 10:20:30":     "2026-01-31 10:20:30",
		"31-JAN-26":               "2026-01-31 00:00:00",
		"31-Jan-2026":             "2026-01-31 00:00:00",
		" 20260131 ":              "2026-01-31 00:00:00",
	}
	for masuk, mau := range kasus {
		hasil, ok := UraiTanggalWarisanTCO(masuk)
		if !ok {
			t.Errorf("%q tidak terurai", masuk)
			continue
		}
		if dapat := utils.FormatTanggalWaktu(hasil); dapat != mau {
			t.Errorf("%q -> %q, mau %q", masuk, dapat, mau)
		}
	}
}

func TestUraiTanggalWarisanTCOKosongDanTakDikenal(t *testing.T) {
	if hasil, ok := UraiTanggalWarisanTCO("   "); !ok || !hasil.IsZero() {
		t.Errorf("kosong harus KOSONG dan ok, dapat %v %v", hasil, ok)
	}
	for _, buruk := range []string{"2026", "31-13-2026", "20261340", "besok", "1/2/2026"} {
		if _, ok := UraiTanggalWarisanTCO(buruk); ok {
			t.Errorf("%q seharusnya TIDAK terurai - ia harus dilaporkan, bukan ditebak", buruk)
		}
	}
}

func TestUraiDesimalWarisanTCO(t *testing.T) {
	baik := map[string]string{
		"12.5": "12.5", "12,5": "12.5", ".5": "0.5", "1000000000.12345678": "1000000000.12345678",
		" 7 ": "7", "0": "0", "-3.25": "-3.25",
	}
	for masuk, mau := range baik {
		d, catatan, ok := UraiDesimalWarisanTCO(masuk)
		if !ok {
			t.Errorf("%q ditolak: %s", masuk, catatan)
			continue
		}
		if dapat := utils.FormatDecimal(d); dapat != mau {
			t.Errorf("%q -> %q, mau %q", masuk, dapat, mau)
		}
	}
	if d, _, ok := UraiDesimalWarisanTCO(""); !ok || d != nil {
		t.Error("kosong harus KOSONG (nil) dan ok")
	}
	buruk := []string{"1.000,5", "abc", "NaN", "Infinity", "1e400"}
	for _, b := range buruk {
		if _, _, ok := UraiDesimalWarisanTCO(b); ok {
			t.Errorf("%q seharusnya ditolak", b)
		}
	}
}

// ⛔ Oracle membulatkan NUMBER(38,8) DIAM-DIAM. Yang melampaui delapan angka
// di belakang koma harus ditolak di sini, sebelum satu digit pun hilang.
func TestUraiDesimalWarisanTCOMenolakPresisiMelampaui(t *testing.T) {
	_, catatan, ok := UraiDesimalWarisanTCO("0.123456789")
	if ok || !strings.Contains(catatan, "belakang koma") {
		t.Errorf("sembilan desimal harus ditolak dengan sebab; ok=%v catatan=%q", ok, catatan)
	}
	// Nol di ekor BUKAN presisi: 0.500000000 tetap 0.5.
	if _, _, ok := UraiDesimalWarisanTCO("0.500000000"); !ok {
		t.Error("nol di ekor tidak melampaui presisi")
	}
	_, catatan, ok = UraiDesimalWarisanTCO("1234567890123456789012345678901")
	if ok || !strings.Contains(catatan, "depan koma") {
		t.Errorf("tiga puluh satu digit bulat harus ditolak; ok=%v catatan=%q", ok, catatan)
	}
}

func TestEkorIdentitasTCO(t *testing.T) {
	if n, ok := EkorIdentitasTCO("1000042", LebarIdentitasTCO); !ok || n != 42 {
		t.Errorf("1000042 -> %d %v, mau 42 true", n, ok)
	}
	if n, ok := EkorIdentitasTCO("10000042", LebarIdentitasKlausulTCO); !ok || n != 42 {
		t.Errorf("10000042 lebar 7 -> %d %v, mau 42 true", n, ok)
	}
	for _, buruk := range []string{"10000042", "2000001", "UJI-1", "", "100000"} {
		if _, ok := EkorIdentitasTCO(buruk, LebarIdentitasTCO); ok {
			t.Errorf("%q seharusnya bukan identitas lebar 6", buruk)
		}
	}
}

func cari(baris []BarisSiapTCO, warisan string) BarisSiapTCO {
	for _, b := range baris {
		if b.Warisan == warisan {
			return b
		}
	}
	return BarisSiapTCO{}
}

func TestKonversiWarisanTCOMemindahkanSeluruhNilai(t *testing.T) {
	keluar, lap := KonversiWarisanTCO(contohWarisanTCO())
	if len(keluar) != 6 {
		t.Fatalf("baris siap %d, mau 6", len(keluar))
	}
	if lap.Memblokir() {
		t.Fatalf("laporan memblokir padahal fixture wajar:\n%s", lap)
	}
	if !lap.Bersih() {
		t.Errorf("laporan tidak bersih:\n%s", lap)
	}

	tahun := cari(keluar, warisanTahunTCO)
	if tahun.Tabel != TabelTahunTCO || tahun.ID != "1000001" {
		t.Errorf("tahun: tabel %s id %s", tahun.Tabel, tahun.ID)
	}
	if tahun.Kanonik["STARTDATE"] != "2026-01-01 00:00:00" || tahun.Kanonik["TGLUPDATE"] != "2026-09-15 00:00:00" {
		t.Errorf("tanggal tahun: %v", tahun.Kanonik)
	}
	// Teks tetap teks, termasuk kode tahun dan PROPORTION.
	if tahun.Nilai["TREATYYEAR"] != "2026" || tahun.Nilai["PROPORTION"] != "P" {
		t.Errorf("teks tahun berubah: %v", tahun.Nilai)
	}

	reins := cari(keluar, warisanReinsurerTCO)
	if reins.Kanonik["PCTSHARE"] != "33.333" || reins.Kanonik["RICOMM"] != "12.5" {
		t.Errorf("desimal reinsurer: %v", reins.Kanonik)
	}
	if reins.Kanonik["STARTDATE"] != "2026-01-01 00:00:00" {
		t.Errorf("stempel Pega tidak bergeser zona: %q", reins.Kanonik["STARTDATE"])
	}
	// ENDDATE kosong tetap KOSONG, bukan tanggal karangan (ADR-U-0027).
	if reins.Nilai["ENDDATE"] != nil || reins.Kanonik["ENDDATE"] != "" {
		t.Errorf("ENDDATE kosong menjadi %v", reins.Nilai["ENDDATE"])
	}
	// IUDATE dibawa apa adanya sebagai teks.
	if reins.Nilai["IUDATE"] != "UJI-APA-ADANYA" {
		t.Errorf("IUDATE berubah: %v", reins.Nilai["IUDATE"])
	}

	sec := cari(keluar, warisanSecurityTCO)
	if sec.ID != "" {
		t.Errorf("security warisan tidak punya ID; surrogate lahir saat menulis, dapat %q", sec.ID)
	}
	if sec.Kanonik["REAS_SECURITY"] != "SEC-A" || sec.Kanonik["TP_TREATY"] != "" {
		t.Errorf("spasi ekor CHAR harus dibuang: %v", sec.Kanonik)
	}
	if sec.Kanonik["PCT_SHARE"] != "50.5" {
		t.Errorf("koma desimal PCT_SHARE: %q", sec.Kanonik["PCT_SHARE"])
	}

	klausul := cari(keluar, warisanKlausulTCO)
	if klausul.Kanonik["RP"] != "1000000000.12345678" || klausul.Kanonik["USD"] != "0.5" ||
		klausul.Kanonik["PCT"] != "12.5" {
		t.Errorf("desimal klausul: %v", klausul.Kanonik)
	}
	if klausul.Kanonik["PARENTREINSTYPEID"] != "00" {
		t.Error("sentinel 00 harus dibawa apa adanya")
	}
	// Sembilan kolom khusus induk kosong pada baris anak -> nil, bukan "".
	for _, k := range []string{"ID_OCCUPATION", "OCCUPATION", "TREATYLIMIT", "MOREUSD"} {
		if klausul.Nilai[k] != nil {
			t.Errorf("kolom %s harus nil pada fixture kosong, dapat %v", k, klausul.Nilai[k])
		}
	}

	// Sequence berikutnya = ekor terbesar + 1, per tabel; klausul lebar 7.
	mau := map[string]int64{
		SeqTahunTCO: 2, SeqKontrakTCO: 8, SeqReinsurerTCO: 4, SeqBusinessTCO: 3, SeqKlausulTCO: 10,
	}
	for seq, n := range mau {
		if lap.SequenceBerikut[seq] != n {
			t.Errorf("%s berikutnya %d, mau %d", seq, lap.SequenceBerikut[seq], n)
		}
	}
	if _, ada := lap.SequenceBerikut[SeqSecurityTCO]; ada {
		t.Error("sequence security tidak diselaraskan dari warisan - warisan tidak punya ID")
	}
}

func TestKonversiWarisanTCOMelaporkanBukanMenebak(t *testing.T) {
	masuk := contohWarisanTCO()
	masuk[warisanTahunTCO][0].Nilai["STARTDATE"] = "kapan-kapan"
	masuk[warisanKlausulTCO][0].Nilai["RP"] = "1.000.000"
	masuk[warisanKlausulTCO][0].Nilai["PCT"] = "0.123456789"
	masuk[warisanSecurityTCO][0].Nilai["REAS_ID"] = "9999999"
	masuk[warisanKontrakTCO][0].Nilai["IDTREATYYEAR"] = "1000999"
	masuk[warisanBusinessTCO] = append(masuk[warisanBusinessTCO],
		barisUji(warisanBusinessTCO, map[string]string{"ID": "", "BIZCODE": "TANPA-ID"}))

	_, lap := KonversiWarisanTCO(masuk)
	if !lap.Memblokir() {
		t.Fatal("lima cacat harus memblokir pemindahan")
	}
	jenis := map[string]int{}
	for _, tm := range lap.Temuan {
		jenis[tm.Jenis]++
	}
	mau := map[string]int{
		TemuanTanggalTakTerurai: 1, TemuanUangTakTerurai: 1, TemuanPresisiMelampaui: 1,
		TemuanRujukanYatim: 2, TemuanIdentitasKosong: 1,
	}
	for j, n := range mau {
		if jenis[j] != n {
			t.Errorf("temuan %q = %d, mau %d\n%s", j, jenis[j], n, lap)
		}
	}
	// Setiap temuan menyebut SUMBER dan MEDAN supaya dapat ditelusuri.
	for _, tm := range lap.Temuan {
		if tm.Sumber == "" || tm.Medan == "" {
			t.Errorf("temuan tanpa sumber/medan: %s", tm)
		}
	}
}

func TestKonversiWarisanTCOIdentitasTakBerbentukTidakMemblokir(t *testing.T) {
	masuk := contohWarisanTCO()
	masuk[warisanTahunTCO][0].Nilai["ID"] = "UJI-TAHUN"
	masuk[warisanKontrakTCO][0].Nilai["IDTREATYYEAR"] = "UJI-TAHUN"
	keluar, lap := KonversiWarisanTCO(masuk)
	if lap.Memblokir() {
		t.Fatalf("identitas tak berbentuk hanya informasi:\n%s", lap)
	}
	ada := false
	for _, tm := range lap.Temuan {
		if tm.Jenis == TemuanIdentitasTakBerbentuk && tm.Nilai == "UJI-TAHUN" {
			ada = true
		}
	}
	if !ada {
		t.Error("identitas tak berbentuk harus tetap dilaporkan")
	}
	if cari(keluar, warisanTahunTCO).ID != "UJI-TAHUN" {
		t.Error("identitas dibawa apa adanya")
	}
	if lap.SequenceBerikut[SeqTahunTCO] != 1 {
		t.Errorf("tanpa identitas berbentuk sequence mulai 1, dapat %d", lap.SequenceBerikut[SeqTahunTCO])
	}
}

func TestKunciKanonikTCOStabil(t *testing.T) {
	k := KunciKanonikTCO([]string{"A", "B"}, map[string]string{"B": "2", "A": "1"})
	if k != "A=1|B=2" {
		t.Errorf("kunci = %q", k)
	}
}

// Penjaga instrumen: pembaca mengenal SELURUH kolom yang procedure tulis.
func TestKolomWarisanTCOSesuaiProcedure(t *testing.T) {
	cacah := map[string]int{
		warisanTahunTCO: 10, warisanKontrakTCO: 8, warisanReinsurerTCO: 19,
		warisanSecurityTCO: 7, warisanBusinessTCO: 12, warisanKlausulTCO: 35,
	}
	for tabel, n := range cacah {
		if len(KolomWarisanTCO(tabel)) != n {
			t.Errorf("%s: %d kolom, mau %d (parameter procedure)", tabel, len(KolomWarisanTCO(tabel)), n)
		}
	}
	// Sembilan kolom khusus induk adalah sembilan kolom TERAKHIR klausul.
	kol := KolomWarisanTCO(warisanKlausulTCO)
	induk := []string{"ID_OCCUPATION", "OCCUPATION", "ID_CLAUSE", "CLAUSE", "TREATYLIMIT",
		"COINS_MIN", "COINS_MAX", "MORERP", "MOREUSD"}
	for i, k := range induk {
		if kol[26+i] != k {
			t.Errorf("kolom ke-%d = %s, mau %s", 27+i, kol[26+i], k)
		}
	}
}
