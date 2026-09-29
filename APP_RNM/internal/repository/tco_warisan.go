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

// kolomMatiKlausulTCO ada di DDL warisan tetapi tidak di-set procedure mana
// pun (AC 70). Tidak dibawa; isi hidupnya dicacah dan dilaporkan.
var kolomMatiKlausulTCO = []string{"PROPORTIONALLIST", "OBJECT"}

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
//     `GetMasterKursList.xml`). `[dugaan]` bagian tanggalnya adalah tanggal
//     kalender yang dimaksud - tanpa pergeseran zona (AC 67).
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
// Hasilnya tanpa lokasi (UTC) dan tanpa pergeseran: angka yang tertulis
// adalah angka yang disimpan.
func UraiTanggalWarisanTCO(teks string) (time.Time, bool) {
	t := strings.TrimSpace(teks)
	if t == "" {
		return time.Time{}, true
	}
	for _, b := range bentukTanggalWarisanTCO {
		if !b.pola.MatchString(t) {
			continue
		}
		bahan := t
		if b.nama == "stempel Pega" {
			bahan = t[:8]
		}
		hasil, err := time.Parse(b.layout, bahan)
		if err != nil {
			return time.Time{}, false
		}
		return hasil, true
	}
	return time.Time{}, false
}

// Batas NUMBER(38,8): delapan angka di belakang koma, tiga puluh di depan.
const (
	skalaTujuanTCO   = 8
	digitBulatTujuan = 30
)

// UraiDesimalWarisanTCO membaca teks uang/persen warisan.
//
// Koma sebagai pemisah desimal DITERIMA - `[terverifikasi]` existing
// menjalankan `@replaceAll(.Pct, ",", ".")` (`TreatyTestChildTotal_Act.xml`
// b571, `SetErrorMessageReinsurer.xml` b270). Teks yang memuat koma DAN titik
// sekaligus ditolak: menafsirkannya berarti menebak mana ribuan mana desimal.
//
// Nilai yang melampaui NUMBER(38,8) - lebih dari delapan angka di belakang
// koma atau lebih dari tiga puluh di depannya - DITOLAK di sini, sebab Oracle
// akan membulatkannya DIAM-DIAM saat disimpan (AC 66).
//
// Teks kosong adalah KOSONG (nil, ok=true).
func UraiDesimalWarisanTCO(teks string) (*apd.Decimal, string, bool) {
	t := strings.TrimSpace(teks)
	if t == "" {
		return nil, "", true
	}
	if strings.Contains(t, ",") {
		if strings.Contains(t, ".") {
			return nil, "memuat koma dan titik sekaligus", false
		}
		t = strings.ReplaceAll(t, ",", ".")
	}
	d, err := utils.ParseDecimal(t)
	if err != nil {
		return nil, err.Error(), false
	}
	ringkas := new(apd.Decimal).Set(d)
	ringkas.Reduce(ringkas)
	if ringkas.Exponent < -skalaTujuanTCO {
		return nil, fmt.Sprintf("lebih dari %d angka di belakang koma", skalaTujuanTCO), false
	}
	if ringkas.NumDigits()+int64(ringkas.Exponent) > digitBulatTujuan {
		return nil, fmt.Sprintf("lebih dari %d angka di depan koma", digitBulatTujuan), false
	}
	return d, "", true
}

// polaEkorIdentitasTCO mengenali identitas warisan `'1' + digit`.
var polaEkorIdentitasTCO = regexp.MustCompile(`^1(\d+)$`)

// EkorIdentitasTCO membaca nomor urut dari identitas berbentuk '1' + lpad.
//
// `1000001` dengan lebar 6 memberi 1; `10000001` dengan lebar 7 memberi 1.
// Bentuk lain - termasuk lebar yang tidak sesuai - bukan identitas sequence.
func EkorIdentitasTCO(id string, lebar int) (int64, bool) {
	m := polaEkorIdentitasTCO.FindStringSubmatch(strings.TrimSpace(id))
	if m == nil || len(m[1]) != lebar {
		return 0, false
	}
	n, err := strconv.ParseInt(m[1], 10, 64)
	if err != nil {
		return 0, false
	}
	return n, true
}
