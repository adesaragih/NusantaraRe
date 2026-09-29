package repository

// Skema WARISAN Treaty Contract Out dan pengurai nilainya - keputusan tco4.
//
// Untuk apa berkas ini: modul ini menulis dan membaca tabel warisan
// `POOLDATA` LANGSUNG (keputusan work owner 29-09-2026, tco4: "khusus modul
// treaty contract out tidak ada tabel baru sama sekali"). Di sini: nama
// tabel, kolom VERBATIM urutan parameter procedure penulisnya, tipe
// DEKLARASI `[data DBA]`, dan pengurai teks warisan (tanggal, desimal,
// identitas) yang dipakai tepi repository.
//
// ⛔ Nol `CREATE TABLE` / `CREATE SEQUENCE` untuk modul ini (penjaga
// `TestTCONolTabelBaru`). Migrasi 300-307 (tco1) dibuang.
//
// Peta lengkap tabel -> kolom -> tipe -> RDB: `.scratch/treaty-contract-out/
// STRUKTUR-TABEL-TREATY-CONTRACT-OUT.md`.

import (
	"database/sql"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/pkg/utils"
)

// Nama enam tabel WARISAN yang modul ini tulis dan baca (tco4).
const (
	warisanTahunTCO     = "TREATYYEAR"
	warisanKontrakTCO   = "TREATYCONTRACT"
	warisanReinsurerTCO = "TREATYREINSURER"
	warisanSecurityTCO  = "MTREATYSECURITY"
	warisanBusinessTCO  = "TREATYBUSINESS"
	warisanKlausulTCO   = "PROPORTIONALARRG"
)

// Nama tabel yang repository modul TULIS dan BACA - tabel warisan itu
// sendiri (tco4). Nama logis dipertahankan supaya pemanggil tidak berubah.
const (
	TabelTahunTCO     = warisanTahunTCO
	TabelKontrakTCO   = warisanKontrakTCO
	TabelReinsurerTCO = warisanReinsurerTCO
	TabelSecurityTCO  = warisanSecurityTCO
	TabelBusinessTCO  = warisanBusinessTCO
	TabelKlausulTCO   = warisanKlausulTCO
)

// Sequence WARISAN yang procedure penulis pakai `[data DBA]`
// (dba-procedures.md): `'1' || lpad(seq.nextval, lebar, '0')`. Nama Oracle
// tak berkutip tidak peka huruf - `TreatyYear_seq` = `TREATYYEAR_SEQ`.
//
// ⛔ MTREATYSECURITY TIDAK punya identitas (tanpa PK, `[data DBA]`): barisnya
// dikunci `(REAS_ID, TRIM(REAS_SECURITY))` seperti `UpdateMTreatySecurity`.
const (
	SeqTahunTCO     = "TREATYYEAR_SEQ"
	SeqKontrakTCO   = "TREATYCONTRACT_SEQ"
	SeqReinsurerTCO = "M_TREATYREINSURER_SEQ"
	SeqBusinessTCO  = "TREATY_BUSINESS_SEQ"
	SeqKlausulTCO   = "PROPORTIONALARRG_SEQ"

	LebarIdentitasTCO        = 6
	LebarIdentitasKlausulTCO = 7
)

// kolomWarisanTCO adalah kolom tiap tabel warisan, VERBATIM, urutan parameter
// procedure penulisnya (dba-procedures.md, SaveMaster*_SQL.xml). Kolom yang
// SAMA namanya dipakai tabel baru - hilir membacanya dengan nama ini.
var kolomWarisanTCO = map[string][]string{
	warisanTahunTCO: {"ID", "TREATYYEAR", "UNDERWRITINGYEAR", "TREATYGROUPID",
		"TREATYGROUPNAME", "USERID", "TGLUPDATE", "PROPORTION", "STARTDATE", "ENDDATE"},
	warisanKontrakTCO: {"ID", "IDTREATYYEAR", "REINSTYPEID", "REINSTYPENAME",
		"TREATYSTARTDATE", "TREATYENDDATE", "USERID", "TGLUPDATE"},
	warisanReinsurerTCO: {"ID", "TREATYYEAR", "TREATYGROUPID", "TREATYGROUPNAME",
		"REINSTYPEID", "REINSTYPENAME", "REINSURERID", "CLIENTID", "NAME", "RICOMM",
		"PCTSHARE", "IUDATE", "USERID", "STARTDATE", "ENDDATE", "STATUSON", "STDRATING",
		"OPERATORNAME", "TGLUPDATE"},
	// MTREATYSECURITY warisan TANPA ID: urutan posisional InsertToMTreatySecurity.
	warisanSecurityTCO: {"THN_TREATY", "TOP_ID", "TP_TREATY", "REAS_ID", "PCT_SHARE",
		"USER_ID", "REAS_SECURITY"},
	warisanBusinessTCO: {"ID", "ISACTIVE", "TREATYYEAR", "TREATYYEARID", "TREATYGROUPID",
		"TREATYGROUPNAME", "REINSTYPEID", "REINSTYPENAME", "BIZCODE", "BIZNAME", "USERID",
		"TGLUPDATE"},
	warisanKlausulTCO: {"ID", "TREATYYEAR", "TREATYYEARID", "TREATYGROUPID",
		"TREATYGROUPNAME", "TREATYDESCID", "TREATYDESCNAME", "REINSTYPEID", "REINSTYPENAME",
		"LAYER", "LAYERPART", "LAYERPARTTYPE", "LAYERTYPE", "KURS", "TGLUPDATE", "USERID",
		"LINE", "PCT", "PCTME", "YDCF", "METHOD", "TERRITORIALLIMIT", "PARENTREINSTYPEID",
		"SPREADINGORDER", "RP", "USD", "ID_OCCUPATION", "OCCUPATION", "ID_CLAUSE", "CLAUSE",
		"TREATYLIMIT", "COINS_MIN", "COINS_MAX", "MORERP", "MOREUSD"},
}

// KolomWarisanTCO mengembalikan kolom warisan sebuah tabel, urutan tetap.
func KolomWarisanTCO(tabel string) []string {
	return append([]string(nil), kolomWarisanTCO[tabel]...)
}

// `PROPORTIONALLIST` dan `OBJECT` ada di DDL warisan tetapi tidak di-set procedure mana
// pun (AC 70). Tidak dibawa; isi hidupnya dicacah dan dilaporkan.

// TipeWarisan menyebut tipe DEKLARASI kolom warisan `[data DBA]`.
//
// ⚠️ Yang diketahui DBA hanya sebagian (dba-procedures.md bab DDL). Kolom
// yang tidak disebut diperlakukan VARCHAR2 - bila kenyataannya berbeda,
// pembacaan lewat TO_CHAR berformat GAGAL TERANG di Oracle (ORA-01722 /
// ORA-01481), bukan diam-diam salah baca.
type TipeWarisan int

const (
	// WarisanTeks - VARCHAR2; dibaca apa adanya.
	WarisanTeks TipeWarisan = iota
	// WarisanAngka - NUMBER; dibaca TO_CHAR TM9 ber-NLS.
	WarisanAngka
	// WarisanTanggal - DATE; dibaca TO_CHAR 'YYYY-MM-DD HH24:MI:SS'.
	WarisanTanggal
	// WarisanChar - CHAR(n); bertabur spasi di ekor, dibaca RTRIM.
	WarisanChar
)

// tipeWarisanTCO mendaftar kolom yang BUKAN VARCHAR2 di warisan.
var tipeWarisanTCO = map[string]map[string]TipeWarisan{
	warisanKontrakTCO:   {"TREATYSTARTDATE": WarisanTanggal, "TREATYENDDATE": WarisanTanggal},
	warisanReinsurerTCO: {"RICOMM": WarisanAngka, "PCTSHARE": WarisanAngka},
	warisanSecurityTCO: {"TP_TREATY": WarisanChar, "REAS_ID": WarisanChar,
		"USER_ID": WarisanChar, "REAS_SECURITY": WarisanChar},
	warisanKlausulTCO: {"TREATYLIMIT": WarisanAngka, "COINS_MIN": WarisanAngka,
		"COINS_MAX": WarisanAngka, "MORERP": WarisanAngka, "MOREUSD": WarisanAngka,
		"TGLUPDATE": WarisanTanggal},
}

// TipeWarisanTCO menjawab tipe deklarasi satu kolom warisan.
func TipeWarisanTCO(tabel, kolom string) TipeWarisan {
	if t, ada := tipeWarisanTCO[tabel][kolom]; ada {
		return t
	}
	return WarisanTeks
}

// bentukTanggalWarisanTCO adalah bentuk teks tanggal yang DIKENAL di warisan.
//
// `[terverifikasi]`/`[dugaan]` per baris:
//   - YYYYMMDD: bentuk properti Date Pega; `SaveTreatyContract_Act.xml` b1040
//     memotong `substring(.TreatyStartDate, 0, 4)` sebagai tahun.
//   - YYYYMMDDTHHMMSS.mmm GMT: stempel DateTime Pega (`TREATYEXCHANGEYEARLY`
//     dibaca `TO_TIMESTAMP_TZ(STARTDATE, 'YYYYMMDD"T"HH24MISS.FF3 TZR')`,
//     `GetMasterKursList.xml`). ⛔ RALAT 29-09-2026 (tco4): stempel ini
//     GMT, dan tanggal kalendernya dibaca di zona Asia/Jakarta - Pega sendiri
//     memformatnya begitu (`SaveTreatyContract_Act` b1479 `@FormatDateTime(…,
//     "dd/MM/yyyy","Asia/Jakarta")`; `SetTanggalTreatyContract` b567 menambah
//     8 jam sebelum memakainya). `…T170000.000 GMT` = 00:00 WIB hari BERIKUT.
//     Dugaan lama "tanpa pergeseran zona" dibantah XML.
//   - DD/MM/YYYY: bentuk yang procedure `to_date(…,'DD/MM/YYYY')` terima
//     (`PEGA_TREATYCONTRACT`, `SaveTreatyContract_Act.xml` b1040/b1479).
//   - YYYY-MM-DD[ HH24:MI:SS]: bentuk utils.ParseTanggal.
//   - DD-MON-RR / DD-MON-YYYY: `[dugaan]` NLS_DATE_FORMAT bawaan Oracle saat
//     `SYSDATE` ditulis ke kolom VARCHAR2 (TGLUPDATE).
//
// Bentuk lain DILAPORKAN, tidak ditebak.
var bentukTanggalWarisanTCO = []struct {
	pola   *regexp.Regexp
	layout string
	nama   string
}{
	{regexp.MustCompile(`^\d{8}$`), "20060102", "YYYYMMDD"},
	{regexp.MustCompile(`^\d{8}T\d{6}(\.\d{1,3})?( GMT)?$`), "20060102", "stempel Pega"},
	{regexp.MustCompile(`^\d{2}/\d{2}/\d{4}$`), "02/01/2006", "DD/MM/YYYY"},
	{regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`), "2006-01-02", "YYYY-MM-DD"},
	{regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}$`), "2006-01-02 15:04:05", "YYYY-MM-DD HH24:MI:SS"},
	{regexp.MustCompile(`^\d{2}-[A-Za-z]{3}-\d{2}$`), "02-Jan-06", "DD-MON-RR"},
	{regexp.MustCompile(`^\d{2}-[A-Za-z]{3}-\d{4}$`), "02-Jan-2006", "DD-MON-YYYY"},
}

// UraiTanggalWarisanTCO membaca teks tanggal warisan.
//
// Teks kosong adalah KOSONG (waktu nol, ok=true) - bukan galat (ADR-U-0027).
// Hasilnya tanggal kalender (UTC, 00:00). Bentuk tanggal dibaca apa adanya;
// stempel Pega (GMT) dibaca di zona Jakarta lebih dulu (lihat daftar bentuk).
func UraiTanggalWarisanTCO(teks string) (time.Time, bool) {
	t := strings.TrimSpace(teks)
	if t == "" {
		return time.Time{}, true
	}
	for _, b := range bentukTanggalWarisanTCO {
		if !b.pola.MatchString(t) {
			continue
		}
		if b.nama == "stempel Pega" {
			w, ok := UraiWaktuWarisanTCO(t)
			if !ok {
				return time.Time{}, false
			}
			y, m, d := w.In(zonaJakartaTCO).Date()
			return time.Date(y, m, d, 0, 0, 0, 0, time.UTC), true
		}
		hasil, err := time.Parse(b.layout, t)
		if err != nil {
			return time.Time{}, false
		}
		return hasil, true
	}
	return time.Time{}, false
}

// polaRibuanTitikTCO - `1.000.000`: titik pemisah ribuan, lihat UraiDesimalWarisanTCO.
var polaRibuanTitikTCO = regexp.MustCompile(`^-?[1-9]\d{0,2}(\.\d{3}){2,}$`)

// UraiDesimalWarisanTCO membaca teks uang/persen warisan.
//
// Koma sebagai pemisah desimal DITERIMA - `[terverifikasi]` existing
// menjalankan `@replaceAll(.Pct, ",", ".")` (`TreatyTestChildTotal_Act.xml`
// b571, `SetErrorMessageReinsurer.xml` b270). Teks yang memuat koma DAN titik
// sekaligus ditolak: menafsirkannya berarti menebak mana ribuan mana desimal.
//
// ⚠️ tco4: TANPA batas skala. Batas NUMBER(38,8) milik tabel `T_*` yang
// dibuang; kolom warisan yang dibaca di sini VARCHAR2 dan boleh menyimpan
// `33.3333333333` yang diketik di Pega - menolaknya membuat seluruh layar
// selingkup gagal terbuka (temuan /code-review lanjutan 3).
//
// Titik pemisah RIBUAN (`1.000.000`) DITERIMA hanya dalam bentuk yang tidak
// mungkin desimal: dua titik atau lebih, kelompok pertama 1-3 angka tanpa nol
// di depan, setiap kelompok sesudahnya tepat tiga angka - desimal hanya punya
// satu titik, jadi ini bukan tebakan. `[data DEV 29-09-2026]`
// `PROPORTIONALARRG` Limit MB: `RP`/`USD` "1.000.000" sedangkan
// `PCTME`/`MORERP` baris yang sama 1000000; tanpanya satu baris mematikan
// seluruh grid (500). ⚠️ `1.000` / `500.000` (SATU titik) tetap DESIMAL - bila
// kolom uang warisan juga memuat ribuan bertitik satu, ia terbaca 1000x lebih
// kecil tanpa galat (OQ-TCO-27, terbuka).
//
// Teks kosong adalah KOSONG (nil, ok=true).
func UraiDesimalWarisanTCO(teks string) (*apd.Decimal, string, bool) {
	t := strings.TrimSpace(teks)
	if t == "" {
		return nil, "", true
	}
	if polaRibuanTitikTCO.MatchString(t) {
		t = strings.ReplaceAll(t, ".", "")
	}
	if strings.Contains(t, ",") {
		if strings.Contains(t, ".") {
			return nil, "memuat koma dan titik sekaligus", false
		}
		t = strings.ReplaceAll(t, ",", ".")
	}
	// Notasi eksponen bukan bentuk desimal Pega (`@toDecimal` / ketikan layar).
	if strings.ContainsAny(t, "eE") {
		return nil, "notasi eksponen bukan bentuk desimal warisan", false
	}
	d, err := utils.ParseDecimal(t)
	if err != nil {
		return nil, err.Error(), false
	}
	return d, "", true
}

// ---------------------------------------------------------------------------
// Tepi TULIS - bentuk teks yang Pega tulis ke kolom VARCHAR2 warisan (tco4).
// Peta kolom -> bentuk: STRUKTUR-TABEL-TREATY-CONTRACT-OUT.md bab "Bentuk nilai".
// ---------------------------------------------------------------------------

// zonaJakartaTCO - WIB. Tetap +7 (Indonesia tanpa DST sejak 1964); zona tetap
// dipakai supaya biner tidak bergantung basis data zona sistem operasi.
var zonaJakartaTCO = time.FixedZone("WIB", 7*60*60)

// polaStempelPegaTCO - `YYYYMMDDTHHMMSS[.mmm][ GMT]`.
var polaStempelPegaTCO = regexp.MustCompile(`^(\d{8}T\d{6})(?:\.(\d{1,3}))?(?: GMT)?$`)

// UraiWaktuWarisanTCO membaca stempel Pega sebagai INSTAN (dalam WIB); bentuk
// tanggal lain dibaca apa adanya (UTC). Kosong = waktu nol, ok.
func UraiWaktuWarisanTCO(teks string) (time.Time, bool) {
	t := strings.TrimSpace(teks)
	if t == "" {
		return time.Time{}, true
	}
	m := polaStempelPegaTCO.FindStringSubmatch(t)
	if m == nil {
		return UraiTanggalWarisanTCO(t)
	}
	w, err := time.ParseInLocation("20060102T150405", m[1], time.UTC)
	if err != nil {
		return time.Time{}, false
	}
	if m[2] != "" {
		ms, _ := strconv.Atoi((m[2] + "00")[:3])
		w = w.Add(time.Duration(ms) * time.Millisecond)
	}
	// Ditampilkan dalam WIB - jam yang Pega tampilkan (temuan /code-review).
	return w.In(zonaJakartaTCO), true
}

// StempelPegaTCO - bentuk `@getCurrentTimeStamp()`: `YYYYMMDDTHHMMSS.mmm GMT`
// (UTC). Waktu nol = teks kosong (NULL).
//
// `[terverifikasi]` `SaveTreatyYear_Act` b328, `SaveTreatyContract_Act` b1458.
func StempelPegaTCO(w time.Time) string {
	if w.IsZero() {
		return ""
	}
	u := w.UTC()
	return u.Format("20060102T150405") + fmt.Sprintf(".%03d GMT", u.Nanosecond()/int(time.Millisecond))
}

// TanggalYYYYMMDDTCO - tanggal kalender -> `YYYYMMDD`, bentuk `TREATYYEAR.
// STARTDATE/ENDDATE` warisan. OQ-TCO-01 DITUTUP dari data DEV (brief lanjutan
// 4): 182/182 baris berbentuk delapan angka - bukan stempel Pega (dugaan kuat
// lanjutan 3 dibantah data). Waktu nol = teks kosong (NULL).
func TanggalYYYYMMDDTCO(tgl time.Time) string {
	if tgl.IsZero() {
		return ""
	}
	return tgl.Format("20060102")
}

// polaYYYYMMDDTCO - delapan angka, satu-satunya bentuk tanggal tahun treaty.
var polaYYYYMMDDTCO = regexp.MustCompile(`^\d{8}$`)

// tanggalTahunWarisanTeks - `TREATYYEAR.STARTDATE/ENDDATE`: HANYA `YYYYMMDD`.
// Bentuk lain (termasuk stempel `…T…GMT`) GAGAL dengan galat berkata-kata -
// bentuk asing di kolom ini berarti penulis lain yang perlu diketahui, bukan
// ditebak. Kosong = kosong.
func tanggalTahunWarisanTeks(v sql.NullString, kolom string) (time.Time, error) {
	t, ok := uraiYYYYMMDDTCO(v.String)
	if !ok {
		return time.Time{}, fmt.Errorf("repository: kolom %s bernilai %q: bukan YYYYMMDD (bentuk TREATYYEAR warisan)",
			kolom, v.String)
	}
	return t, nil
}

// uraiYYYYMMDDTCO - `YYYYMMDD` -> tanggal (UTC 00:00); kosong = nol, ok.
func uraiYYYYMMDDTCO(teks string) (time.Time, bool) {
	t := strings.TrimSpace(teks)
	if t == "" {
		return time.Time{}, true
	}
	if !polaYYYYMMDDTCO.MatchString(t) {
		return time.Time{}, false
	}
	hasil, err := time.Parse("20060102", t)
	return hasil, err == nil
}

// TulisDesimalWarisanTCO - desimal ke kolom teks warisan: titik, tanpa
// pemisah ribuan, tanpa eksponen, nol di ekor dibuang. nil = NULL.
//
// `[dugaan kuat]` bentuk hasil `@toDecimal` Pega (`HitungRpUsd_depan`
// b293-b294, b382-b383) - OQ-TCO-23. Pembaca menerima titik ATAU koma.
func TulisDesimalWarisanTCO(d *apd.Decimal) any {
	if d == nil {
		return nil
	}
	r := new(apd.Decimal).Set(d)
	r.Reduce(r)
	// Text('f') menulis bentuk biasa walau Reduce menghasilkan eksponen positif.
	return r.Text('f')
}

// tanggalWarisanTeks - kolom VARCHAR2 bertanggal -> tanggal kalender; teks
// yang tidak dikenal GAGAL TERANG, tidak menjadi kosong diam-diam.
func tanggalWarisanTeks(v sql.NullString, kolom string) (time.Time, error) {
	t, ok := UraiTanggalWarisanTCO(v.String)
	if !ok {
		return time.Time{}, fmt.Errorf("repository: kolom %s bernilai %q: bentuk tanggal tidak dikenal", kolom, v.String)
	}
	return t, nil
}

// waktuWarisanTeks - kolom VARCHAR2 cap waktu (`TGLUPDATE`) -> instan.
func waktuWarisanTeks(v sql.NullString, kolom string) (time.Time, error) {
	t, ok := UraiWaktuWarisanTCO(v.String)
	if !ok {
		return time.Time{}, fmt.Errorf("repository: kolom %s bernilai %q: bentuk waktu tidak dikenal", kolom, v.String)
	}
	return t, nil
}
