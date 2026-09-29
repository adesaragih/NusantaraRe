package repository

// Konversi data warisan Treaty Contract Out - bagian MURNI (tiket 01).
//
// Untuk apa berkas ini: mengubah baris enam tabel warisan `POOLDATA.TREATYYEAR`,
// `TREATYCONTRACT`, `TREATYREINSURER`, `MTREATYSECURITY`, `TREATYBUSINESS`,
// `PROPORTIONALARRG` (seluruhnya sudah dibaca sebagai TEKS) menjadi nilai
// bertipe untuk tabel `T_*` yang baru, sambil MELAPORKAN apa pun yang tidak
// dapat dipindahkan dengan yakin - bukan menebak, bukan mendiamkan.
//
// Fungsi di sini tidak menyentuh basis data. Pembaca dan penulisnya ada di
// tco_pindah.go; bentuk tiruan untuk uji ada di skemauji.
//
// Aturan yang dijaga:
//   AC 51/52  uang dan persen menjadi desimal presisi arbitrer, nol float.
//   AC 53     STARTDATE, ENDDATE, TREATYSTARTDATE, TREATYENDDATE, TGLUPDATE
//             menjadi DATE.
//   AC 64     nol pembacaan tabel dokumen `M_*` - sumbernya HANYA relasional.
//   AC 66/67  nilai yang tidak terurai DILAPORKAN, tidak diam-diam jadi nol
//             atau kosong; desimal yang melampaui NUMBER(38,8) dilaporkan,
//             bukan dibulatkan diam-diam oleh Oracle.
//   AC 68     baris MTREATYSECURITY mendapat PK surrogate; rujukan yatim ke
//             reinsurer dilaporkan.
//   AC 70     kolom warisan PROPORTIONALLIST dan OBJECT tidak dibawa; isi
//             hidupnya dicacah dan dilaporkan.
//   ADR-U-0022 konversi tipe SEKALI saat masuk; kode dan pengenal tetap teks.

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/pkg/utils"
)

// Nama enam tabel WARISAN - sumber, dibaca saja.
//
// ⛔ Keenam nama ini hanya boleh muncul di pembaca (`bacaWarisanTCO`) dan di
// pencacah kolom mati; penjaga `TestTCOWarisanHanyaDibaca` menolak fungsi
// yang menyebutnya bersama INSERT/UPDATE/DELETE.
const (
	warisanTahunTCO     = "TREATYYEAR"
	warisanKontrakTCO   = "TREATYCONTRACT"
	warisanReinsurerTCO = "TREATYREINSURER"
	warisanSecurityTCO  = "MTREATYSECURITY"
	warisanBusinessTCO  = "TREATYBUSINESS"
	warisanKlausulTCO   = "PROPORTIONALARRG"
)

// Nama enam tabel BARU (keputusan tco1: awalan T_).
const (
	TabelTahunTCO     = "T_TREATYYEAR"
	TabelKontrakTCO   = "T_TREATYCONTRACT"
	TabelReinsurerTCO = "T_TREATYREINSURER"
	TabelSecurityTCO  = "T_MTREATYSECURITY"
	TabelBusinessTCO  = "T_TREATYBUSINESS"
	TabelKlausulTCO   = "T_PROPORTIONALARRG"
	TabelJejakTCO     = "T_TREATYCO_JEJAK"
	// TabelLampiranTCO - lampiran tahun treaty, fitur BARU tiket 12 (tanpa warisan).
	TabelLampiranTCO = "T_TREATYYEAR_LAMPIRAN"
)

// Sequence tiap tabel baru, beserta LEBAR digit di belakang awalan '1'
// (spec §7: enam digit, kecuali klausul tujuh).
const (
	SeqTahunTCO     = "SEQ_T_TREATYYEAR"
	SeqKontrakTCO   = "SEQ_T_TREATYCONTRACT"
	SeqReinsurerTCO = "SEQ_T_TREATYREINSURER"
	SeqSecurityTCO  = "SEQ_T_MTREATYSECURITY"
	SeqBusinessTCO  = "SEQ_T_TREATYBUSINESS"
	SeqKlausulTCO   = "SEQ_T_PROPORTIONALARRG"
	SeqJejakTCO     = "SEQ_T_TREATYCO_JEJAK"
	// SeqLampiranTCO - tiket 12.
	SeqLampiranTCO = "SEQ_T_TREATYYEAR_LAMPIRAN"

	LebarIdentitasTCO        = 6
	LebarIdentitasKlausulTCO = 7
	// LebarIdentitasLampiranTCO - tiket 12, tanpa padanan warisan (keputusan kami).
	LebarIdentitasLampiranTCO = 9
)

// urutanWarisanTCO adalah urutan pemindahan: induk sebelum anak (FK).
var urutanWarisanTCO = []string{
	warisanTahunTCO, warisanKontrakTCO, warisanReinsurerTCO,
	warisanSecurityTCO, warisanBusinessTCO, warisanKlausulTCO,
}

// tabelBaruDariWarisanTCO memetakan nama warisan ke nama baru.
var tabelBaruDariWarisanTCO = map[string]string{
	warisanTahunTCO:     TabelTahunTCO,
	warisanKontrakTCO:   TabelKontrakTCO,
	warisanReinsurerTCO: TabelReinsurerTCO,
	warisanSecurityTCO:  TabelSecurityTCO,
	warisanBusinessTCO:  TabelBusinessTCO,
	warisanKlausulTCO:   TabelKlausulTCO,
}

// sequenceDariWarisanTCO memetakan nama warisan ke sequence baru dan lebarnya.
var sequenceDariWarisanTCO = map[string]struct {
	Nama  string
	Lebar int
}{
	warisanTahunTCO:     {SeqTahunTCO, LebarIdentitasTCO},
	warisanKontrakTCO:   {SeqKontrakTCO, LebarIdentitasTCO},
	warisanReinsurerTCO: {SeqReinsurerTCO, LebarIdentitasTCO},
	warisanBusinessTCO:  {SeqBusinessTCO, LebarIdentitasTCO},
	warisanKlausulTCO:   {SeqKlausulTCO, LebarIdentitasKlausulTCO},
}

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

// TipeTujuan menyebut tipe kolom di tabel BARU.
type TipeTujuan int

const (
	// TujuanTeks - VARCHAR2.
	TujuanTeks TipeTujuan = iota
	// TujuanDesimal - NUMBER(38,8): uang, persen, kurs.
	TujuanDesimal
	// TujuanTanggal - DATE.
	TujuanTanggal
)

// tipeTujuanTCO mendaftar kolom baru yang BUKAN teks (STRUKTUR-TABEL-
// TREATY-CONTRACT-OUT.md). Kolom lain tetap teks.
var tipeTujuanTCO = map[string]map[string]TipeTujuan{
	warisanTahunTCO: {"TGLUPDATE": TujuanTanggal, "STARTDATE": TujuanTanggal,
		"ENDDATE": TujuanTanggal},
	warisanKontrakTCO: {"TREATYSTARTDATE": TujuanTanggal, "TREATYENDDATE": TujuanTanggal,
		"TGLUPDATE": TujuanTanggal},
	warisanReinsurerTCO: {"RICOMM": TujuanDesimal, "PCTSHARE": TujuanDesimal,
		"STARTDATE": TujuanTanggal, "ENDDATE": TujuanTanggal, "TGLUPDATE": TujuanTanggal},
	warisanSecurityTCO: {"PCT_SHARE": TujuanDesimal},
	warisanBusinessTCO: {"TGLUPDATE": TujuanTanggal},
	warisanKlausulTCO: {"KURS": TujuanDesimal, "TGLUPDATE": TujuanTanggal,
		"PCT": TujuanDesimal, "PCTME": TujuanDesimal, "RP": TujuanDesimal,
		"USD": TujuanDesimal, "TREATYLIMIT": TujuanDesimal, "COINS_MIN": TujuanDesimal,
		"COINS_MAX": TujuanDesimal, "MORERP": TujuanDesimal, "MOREUSD": TujuanDesimal},
}

// TipeTujuanTCO menjawab tipe kolom di tabel baru, dinamai dengan nama
// warisannya.
func TipeTujuanTCO(tabel, kolom string) TipeTujuan {
	if t, ada := tipeTujuanTCO[tabel][kolom]; ada {
		return t
	}
	return TujuanTeks
}

// BarisWarisanTCO adalah satu baris tabel warisan, seluruh nilainya TEKS
// persis seperti dikembalikan Oracle (NULL menjadi teks kosong).
type BarisWarisanTCO struct {
	Tabel string
	Nilai map[string]string
}

// BarisSiapTCO adalah satu baris yang siap ditulis ke tabel baru.
//
// Nilai per kolom: string (teks), time.Time (tanggal), *apd.Decimal
// (desimal), atau nil (kosong). Kanonik menyimpan bentuk TEKS tiap kolom
// untuk rekonsiliasi: tanggal `YYYY-MM-DD HH24:MI:SS`, desimal
// `utils.FormatDecimal`, teks apa adanya.
type BarisSiapTCO struct {
	Warisan string
	Tabel   string
	// ID kosong pada MTREATYSECURITY: surrogate diberikan saat menulis.
	ID      string
	Kolom   []string
	Nilai   map[string]any
	Kanonik map[string]string
}

// Jenis temuan konversi, di samping yang sudah ada di migrasidata.go.
const (
	TemuanPresisiMelampaui      = "desimal melampaui NUMBER(38,8)"
	TemuanRujukanYatim          = "rujukan ke induk tidak ditemukan"
	TemuanIdentitasTakBerbentuk = "identitas tidak berbentuk '1' + digit"
	TemuanKolomMatiBerisi       = "kolom warisan yang tidak dibawa berisi"
)

// temuanMemblokirTCO adalah jenis temuan yang MEMBATALKAN pemindahan.
//
// Yang tidak memblokir - identitas tak berbentuk dan kolom mati berisi -
// tetap dilaporkan: yang pertama hanya menyingkirkan baris itu dari
// penyelarasan sequence, yang kedua keputusan work owner (AC 70).
var temuanMemblokirTCO = map[string]bool{
	TemuanTanggalTakTerurai: true,
	TemuanUangTakTerurai:    true,
	TemuanPresisiMelampaui:  true,
	TemuanRujukanYatim:      true,
	TemuanIdentitasKosong:   true,
}

// LaporanMigrasiTCO menceritakan apa yang benar-benar terjadi.
type LaporanMigrasiTCO struct {
	// Masuk mencacah baris warisan per tabel warisan.
	Masuk map[string]int
	// Ditulis mencacah baris yang ditulis per tabel baru.
	Ditulis map[string]int
	// DibacaKembali mencacah baris tabel baru saat rekonsiliasi.
	DibacaKembali map[string]int
	// SequenceBerikut adalah nilai NEXTVAL yang dipasang per sequence.
	SequenceBerikut map[string]int64
	// KolomMatiBerisi mencacah baris PROPORTIONALARRG yang kolom matinya terisi.
	KolomMatiBerisi map[string]int
	Temuan          []Temuan
	// Selisih adalah hasil rekonsiliasi: kosong berarti seluruh nilai pulang
	// persis seperti yang dikonversi.
	Selisih []string
}

func laporanKosongTCO() LaporanMigrasiTCO {
	return LaporanMigrasiTCO{
		Masuk: map[string]int{}, Ditulis: map[string]int{}, DibacaKembali: map[string]int{},
		SequenceBerikut: map[string]int64{}, KolomMatiBerisi: map[string]int{},
	}
}

// Memblokir menyatakan ada temuan yang membatalkan pemindahan.
func (l LaporanMigrasiTCO) Memblokir() bool {
	for _, t := range l.Temuan {
		if temuanMemblokirTCO[t.Jenis] {
			return true
		}
	}
	return false
}

// Bersih menyatakan nol temuan dan nol selisih.
func (l LaporanMigrasiTCO) Bersih() bool { return len(l.Temuan) == 0 && len(l.Selisih) == 0 }

// String meringkas laporan untuk log.
func (l LaporanMigrasiTCO) String() string {
	var b strings.Builder
	for _, t := range urutanWarisanTCO {
		fmt.Fprintf(&b, "%s: masuk %d, ditulis %d, dibaca kembali %d\n",
			t, l.Masuk[t], l.Ditulis[tabelBaruDariWarisanTCO[t]],
			l.DibacaKembali[tabelBaruDariWarisanTCO[t]])
	}
	for _, k := range kolomMatiKlausulTCO {
		fmt.Fprintf(&b, "kolom mati %s berisi pada %d baris\n", k, l.KolomMatiBerisi[k])
	}
	seq := make([]string, 0, len(l.SequenceBerikut))
	for s := range l.SequenceBerikut {
		seq = append(seq, s)
	}
	sort.Strings(seq)
	for _, s := range seq {
		fmt.Fprintf(&b, "sequence %s berikutnya %d\n", s, l.SequenceBerikut[s])
	}
	for _, t := range l.Temuan {
		fmt.Fprintf(&b, "temuan: %s\n", t)
	}
	for _, s := range l.Selisih {
		fmt.Fprintf(&b, "selisih: %s\n", s)
	}
	return b.String()
}

// ---------------------------------------------------------------------------
// Pengurai nilai warisan
// ---------------------------------------------------------------------------

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

// ---------------------------------------------------------------------------
// Konversi
// ---------------------------------------------------------------------------

// KonversiWarisanTCO mengubah baris warisan menjadi baris siap tulis.
//
// Masukannya dikelompokkan per tabel warisan. Keluarannya berurut induk lalu
// anak, dan laporan memuat setiap nilai yang tidak dapat dipindahkan dengan
// yakin. Fungsi ini MURNI.
func KonversiWarisanTCO(masuk map[string][]BarisWarisanTCO) ([]BarisSiapTCO, LaporanMigrasiTCO) {
	lap := laporanKosongTCO()
	var keluar []BarisSiapTCO

	identitas := map[string]map[string]bool{}
	for _, tabel := range urutanWarisanTCO {
		identitas[tabel] = map[string]bool{}
	}

	for _, tabel := range urutanWarisanTCO {
		baris := masuk[tabel]
		lap.Masuk[tabel] = len(baris)
		lebar := LebarIdentitasTCO
		if tabel == warisanKlausulTCO {
			lebar = LebarIdentitasKlausulTCO
		}
		for i, b := range baris {
			sumber := fmt.Sprintf("%s#%d", tabel, i+1)
			if id := strings.TrimSpace(b.Nilai["ID"]); id != "" {
				sumber = tabel + "/" + id
			}
			siap := BarisSiapTCO{
				Warisan: tabel, Tabel: tabelBaruDariWarisanTCO[tabel],
				Kolom: KolomWarisanTCO(tabel), Nilai: map[string]any{}, Kanonik: map[string]string{},
			}
			for _, kolom := range siap.Kolom {
				teks := b.Nilai[kolom]
				if TipeWarisanTCO(tabel, kolom) == WarisanChar {
					// CHAR(n) bertabur spasi di ekor: itu artefak penyimpanan,
					// bukan data. trim() di kunci warisan mengakuinya
					// (UpdateMTreatySecurity.xml). Yang dibuang HANYA spasi ekor.
					teks = strings.TrimRight(teks, " ")
				}
				switch TipeTujuanTCO(tabel, kolom) {
				case TujuanTanggal:
					t, ok := UraiTanggalWarisanTCO(teks)
					if !ok {
						lap.Temuan = append(lap.Temuan, Temuan{Jenis: TemuanTanggalTakTerurai,
							Sumber: sumber, Medan: kolom, Nilai: teks,
							Catatan: "bentuk tidak dikenal; lihat bentukTanggalWarisanTCO"})
						continue
					}
					if t.IsZero() {
						siap.Nilai[kolom] = nil
						siap.Kanonik[kolom] = ""
					} else {
						siap.Nilai[kolom] = t
						siap.Kanonik[kolom] = utils.FormatTanggalWaktu(t)
					}
				case TujuanDesimal:
					d, catatan, ok := UraiDesimalWarisanTCO(teks)
					if !ok {
						jenis := TemuanUangTakTerurai
						if strings.Contains(catatan, "angka di") {
							jenis = TemuanPresisiMelampaui
						}
						lap.Temuan = append(lap.Temuan, Temuan{Jenis: jenis,
							Sumber: sumber, Medan: kolom, Nilai: teks, Catatan: catatan})
						continue
					}
					if d == nil {
						siap.Nilai[kolom] = nil
						siap.Kanonik[kolom] = ""
					} else {
						siap.Nilai[kolom] = d
						// Kanonik dalam bentuk RINGKAS: "12.50" dan "12.5" satu nilai.
						ringkas := new(apd.Decimal).Set(d)
						ringkas.Reduce(ringkas)
						siap.Kanonik[kolom] = utils.FormatDecimal(ringkas)
					}
				default:
					if teks == "" {
						siap.Nilai[kolom] = nil
					} else {
						siap.Nilai[kolom] = teks
					}
					siap.Kanonik[kolom] = teks
				}
			}

			if tabel != warisanSecurityTCO {
				siap.ID = strings.TrimSpace(b.Nilai["ID"])
				if siap.ID == "" {
					lap.Temuan = append(lap.Temuan, Temuan{Jenis: TemuanIdentitasKosong,
						Sumber: sumber, Medan: "ID", Catatan: "baris tanpa identitas tidak dapat ditulis"})
				} else {
					identitas[tabel][siap.ID] = true
					if _, ok := EkorIdentitasTCO(siap.ID, lebar); !ok {
						lap.Temuan = append(lap.Temuan, Temuan{Jenis: TemuanIdentitasTakBerbentuk,
							Sumber: sumber, Medan: "ID", Nilai: siap.ID,
							Catatan: fmt.Sprintf("dibawa apa adanya; tidak ikut menyelaraskan sequence (lebar %d)", lebar)})
					}
				}
			}
			keluar = append(keluar, siap)
		}
	}

	// Rujukan yatim: FK yang akan ditolak Oracle dilaporkan LEBIH DULU dengan
	// menyebut barisnya, bukan dibiarkan gagal sebagai ORA-02291 tanpa nama.
	for _, s := range keluar {
		switch s.Warisan {
		case warisanKontrakTCO:
			if induk := s.Kanonik["IDTREATYYEAR"]; induk != "" && !identitas[warisanTahunTCO][induk] {
				lap.Temuan = append(lap.Temuan, Temuan{Jenis: TemuanRujukanYatim,
					Sumber: warisanKontrakTCO + "/" + s.ID, Medan: "IDTREATYYEAR", Nilai: induk,
					Catatan: "tahun treaty dengan ID itu tidak ada di TREATYYEAR"})
			}
		case warisanSecurityTCO:
			induk := s.Kanonik["REAS_ID"]
			if induk == "" || !identitas[warisanReinsurerTCO][induk] {
				lap.Temuan = append(lap.Temuan, Temuan{Jenis: TemuanRujukanYatim,
					Sumber: warisanSecurityTCO + "/" + s.Kanonik["REAS_SECURITY"], Medan: "REAS_ID",
					Nilai: induk, Catatan: "reinsurer dengan ID itu tidak ada di TREATYREINSURER"})
			}
		}
	}

	// Nilai sequence berikutnya: satu lebih dari ekor identitas terbesar yang
	// berbentuk. Tabel tanpa baris berbentuk mulai dari 1.
	for _, tabel := range urutanWarisanTCO {
		seq, ada := sequenceDariWarisanTCO[tabel]
		if !ada {
			continue
		}
		var maks int64
		for id := range identitas[tabel] {
			if n, ok := EkorIdentitasTCO(id, seq.Lebar); ok && n > maks {
				maks = n
			}
		}
		lap.SequenceBerikut[seq.Nama] = maks + 1
	}
	return keluar, lap
}

// KunciKanonikTCO merakit satu teks dari seluruh kolom kanonik sebuah baris,
// urutan kolom tetap. Dipakai rekonsiliasi MTREATYSECURITY, yang di warisan
// tidak punya identitas: barisnya dibandingkan sebagai HIMPUNAN GANDA.
func KunciKanonikTCO(kolom []string, kanonik map[string]string) string {
	bagian := make([]string, 0, len(kolom))
	for _, k := range kolom {
		bagian = append(bagian, k+"="+kanonik[k])
	}
	return strings.Join(bagian, "|")
}
