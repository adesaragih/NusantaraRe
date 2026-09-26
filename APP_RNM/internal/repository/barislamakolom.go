package repository

// Daftar kolom tabel datar warisan - satu sumber untuk semua pemakainya.
//
// Untuk apa berkas ini: OS_AKSEPTASI_KLAIM_LIFE punya 62 kolom - rule warisan
// MENULIS 55 di antaranya - dan tiga tempat
// berbeda perlu menyebutnya dalam urutan yang sama - pembaca (AmbilBarisLama),
// penulis fixture (skemauji.IsiBarisLama), dan tabel tiruan. Menyalin daftar
// itu tiga kali adalah cara paling pasti untuk membuat urutan SELECT dan urutan
// Scan berselisih tanpa ada yang menyadarinya. Karena itu daftarnya ditulis
// SEKALI di sini, berpasangan langsung dengan medan strukturnya.
//
// Dibaca sesudah: migrasidata.go.
//
// ⚠️ [data DBA] Tipe fisik kolom tabel SUNGGUHAN belum dikonfirmasi. Yang
// [terverifikasi] dari korpus hanyalah NAMA kelima puluh lima kolomnya - dari
// rule ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL!RNM!UPDATEOSAKSEPTASICLAIMLIFE_SQL
// bertipe Rule-Connect-SQL. Penggolongan angka dan tanggal di bawah adalah
// pembacaan atas NAMA kolomnya, bukan atas katalog Oracle. Bila kelak sebuah
// kolom ternyata VARCHAR2 di produksi, TO_CHAR berformat angka akan menjawab
// ORA-01722 dan yang berubah adalah dua daftar di bawah - bukan kode pembacanya.

import (
	"database/sql"
	"fmt"
	"strings"

	"nusantarare/pkg/utils"
)

// medanLama memasangkan satu nama kolom dengan tempat nilainya ditampung.
type medanLama struct {
	Kolom string
	Nilai *string
}

// KolomWarisan memasangkan satu kolom OS_AKSEPTASI_KLAIM_LIFE dengan tipenya.
type KolomWarisan struct {
	Kolom string
	Tipe  string
}

// kolomWarisan adalah SELURUH 62 kolom OS_AKSEPTASI_KLAIM_LIFE beserta tipenya,
// berurut sama dengan katalog.
//
// ⭐ SUMBER: `[data DBA]` .scratch\claim-life\TIPE-KOLOM-OS-AKSEPTASI-KLAIM-LIFE.md,
// dibaca dari ALL_TAB_COLUMNS instance pengembangan 26 September 2026 - bukan
// lagi tebakan atas nama kolom. TestPetaTipeCocokDenganKatalog membandingkan
// tabel ini dengan dokumen itu baris demi baris, sehingga keduanya tidak dapat
// berselisih diam-diam.
//
// ⚠️ Tabelnya 62 kolom; rule warisan hanya MENULIS 55 di antaranya. Ketujuh
// sisanya tetap didaftar di sini supaya tabel TIRUAN di skema uji berbentuk
// sama dengan tabel sungguhan - kalau tidak, test pulang-pergi menguji bentuk
// yang tidak pernah ada di produksi.
//
// ⛔ Tiga di antaranya dulu DITEBAK teks dan ternyata bukan: STS_REJECT
// NUMBER(38,0), CLAIM_RETRO NUMBER, WPC DATE. Tebakan itulah yang membuat
// tiruan lama bertipe VARCHAR2 semua, sehingga jebakan NLS justru lolos uji.
var kolomWarisan = []KolomWarisan{
	{"CASEID", "VARCHAR2(100)"},
	{"NO_CLAIM", "VARCHAR2(100)"},
	{"NO_ACCEPTATION", "VARCHAR2(100)"},
	{"STS_REJECT", "NUMBER(38,0)"},
	{"ACCEPTATION_DATE", "DATE"},
	{"POLICY_NO", "VARCHAR2(100)"},
	{"POLICY_HOLDER", "VARCHAR2(1000)"},
	{"CERTIFICATE_NO", "VARCHAR2(100)"},
	{"NAME_OF_INSURED", "VARCHAR2(1000)"},
	{"SEX", "VARCHAR2(50)"},
	{"DOB", "DATE"},
	{"AGE", "NUMBER(38,0)"},
	{"PLAN", "VARCHAR2(255)"},
	{"BEGIN_DATE", "DATE"},
	{"LAPSE_DATE", "DATE"},
	{"EXPIRED_DATE", "DATE"},
	{"STATUS", "VARCHAR2(100)"},
	{"CURRENCY", "VARCHAR2(100)"},
	{"STS_KONVERSI", "CHAR(1)"},
	{"TGL_KONVERSI", "DATE"},
	{"WPC", "DATE"},
	{"PL_NUMBER", "VARCHAR2(100)"},
	{"DISEASE", "VARCHAR2(1000)"},
	{"ICD_CODE", "VARCHAR2(10)"},
	{"NOTES", "VARCHAR2(1000)"},
	{"CEDINGCO", "VARCHAR2(100)"},
	{"CEDINGCONAME", "VARCHAR2(1000)"},
	{"SOB", "VARCHAR2(100)"},
	{"SOBNAME", "VARCHAR2(1000)"},
	{"BUSINESSID", "VARCHAR2(100)"},
	{"BUSINESSNAME", "VARCHAR2(100)"},
	{"KETERANGAN", "VARCHAR2(1000)"},
	{"EM_PERCENT", "NUMBER"},
	{"SUM_INSURED", "NUMBER"},
	{"CEDING_RETENTION", "NUMBER"},
	{"SUM_REASURED", "NUMBER"},
	{"SHARE_NUSANTARA_RE", "NUMBER"},
	{"CLAIM_AMOUNT", "NUMBER"},
	{"SHARE_RETRO", "NUMBER"},
	{"CLAIM_RETRO", "NUMBER"},
	{"RETROID", "VARCHAR2(100)"},
	{"RETRONAME", "VARCHAR2(100)"},
	{"SECURITYREINSURERID", "VARCHAR2(100)"},
	{"SECURITYREINSURER", "VARCHAR2(100)"},
	{"TYPECEDING", "VARCHAR2(50)"},
	{"TYPE", "VARCHAR2(10)"},
	{"CONFIRMATION_DATE", "DATE"},
	{"CLAIM_RECEIVED_DATE", "DATE"},
	{"COMPLETE_DATE", "DATE"},
	{"ID", "VARCHAR2(100)"},
	{"NAME_OF_BANK", "VARCHAR2(100)"},
	{"IDBANK", "VARCHAR2(100)"},
	{"ACCOUNTNO", "VARCHAR2(100)"},
	{"PRODUCTNAMEID", "VARCHAR2(100)"},
	{"PRODUCTNAME", "VARCHAR2(1000)"},
	{"CREATEOPNAME", "VARCHAR2(100)"},
	{"RETROCEDED_SHARE", "NUMBER"},
	{"LAYER_1", "VARCHAR2(10)"},
	{"LAYER_2", "VARCHAR2(10)"},
	{"LAYER_3", "VARCHAR2(10)"},
	{"LAYER_4", "VARCHAR2(10)"},
	{"INDEXLIST", "VARCHAR2(10)"},
}

// kolomAngkaLama menyebut kolom yang isinya bilangan.
//
// Ia dibaca lewat TO_CHAR ber-argumen NLS supaya tiba sebagai TEKS dengan titik
// sebagai pemisah desimal, apa pun setelan sesi (ADR-U-0003, ADR-U-0016).
// Tanpa itu, sesi ber-NLS koma menyerahkan "1234,56" dan seluruh pembacaan
// uang rusak diam-diam.
//
// Isinya DITURUNKAN dari kolomWarisan, bukan ditulis ulang: satu daftar yang
// salah lebih baik daripada dua daftar yang berselisih.
var kolomAngkaLama = golonganWarisan("NUMBER")

// kolomTanggalLama menyebut kolom yang isinya tanggal.
//
// Dibaca dengan bentuk yang sama persis dengan yang dikenal utils.ParseTanggal,
// supaya tidak ada tanggal yang bergantung pada NLS_DATE_FORMAT sesi.
var kolomTanggalLama = golonganWarisan("DATE")

// golonganWarisan mengumpulkan kolom yang tipenya berawalan awalan tertentu.
func golonganWarisan(awalan string) map[string]bool {
	hasil := map[string]bool{}
	for _, k := range kolomWarisan {
		if strings.HasPrefix(k.Tipe, awalan) {
			hasil[k.Kolom] = true
		}
	}
	return hasil
}

// tipeWarisan memetakan nama kolom ke tipe katalognya.
var tipeWarisan = func() map[string]string {
	m := map[string]string{}
	for _, k := range kolomWarisan {
		m[k.Kolom] = k.Tipe
	}
	return m
}()

// NamaKolomTabelWarisan mengeluarkan keenam puluh dua nama kolom, berurut sama
// dengan katalog.
//
// Dipakai skema uji untuk membuat tabel tiruan. Berbeda dari
// NamaKolomBarisLama, yang hanya kelima puluh lima kolom yang DITULIS rule.
func NamaKolomTabelWarisan() []string {
	out := make([]string, 0, len(kolomWarisan))
	for _, k := range kolomWarisan {
		out = append(out, k.Kolom)
	}
	return out
}

// fmtTanggalOracle adalah bentuk tanggal yang diminta dari Oracle. Ia cocok
// dengan utils.TanggalWaktu, dan ParseTanggal juga menerima bentuk pendeknya.
const fmtTanggalOracle = `TO_CHAR(%s, 'YYYY-MM-DD HH24:MI:SS')`

// ekspresiBacaLama menyusun potongan SELECT untuk satu kolom.
//
// Kolom teks dibaca apa adanya; kolom angka dan tanggal dibungkus TO_CHAR
// supaya jenis datanya tidak lagi bergantung pada setelan sesi.
func ekspresiBacaLama(kolom string) string {
	switch {
	case kolomAngkaLama[kolom]:
		return fmt.Sprintf(fmtDesimal, kolom)
	case kolomTanggalLama[kolom]:
		return fmt.Sprintf(fmtTanggalOracle, kolom)
	default:
		return kolom
	}
}

// medanBarisLama memasangkan kelima puluh lima kolom dengan medan strukturnya.
//
// Urutannya mengikat: pemanggil menyusun SELECT dan Scan dari daftar yang sama,
// sehingga keduanya tidak mungkin berselisih.
func medanBarisLama(b *BarisLama) []medanLama {
	return []medanLama{
		{"CASEID", &b.CASEID}, {"NO_CLAIM", &b.NO_CLAIM},
		{"POLICY_NO", &b.POLICY_NO}, {"POLICY_HOLDER", &b.POLICY_HOLDER},
		{"CERTIFICATE_NO", &b.CERTIFICATE_NO},
		{"NAME_OF_INSURED", &b.NAME_OF_INSURED}, {"SEX", &b.SEX},
		{"DOB", &b.DOB}, {"AGE", &b.AGE}, {"PLAN", &b.PLAN},
		{"BEGIN_DATE", &b.BEGIN_DATE}, {"LAPSE_DATE", &b.LAPSE_DATE},
		{"EXPIRED_DATE", &b.EXPIRED_DATE}, {"STATUS", &b.STATUS},
		{"EM_PERCENT", &b.EM_PERCENT}, {"CURRENCY", &b.CURRENCY},
		{"SUM_INSURED", &b.SUM_INSURED}, {"CEDING_RETENTION", &b.CEDING_RETENTION},
		{"SUM_REASURED", &b.SUM_REASURED},
		{"SHARE_NUSANTARA_RE", &b.SHARE_NUSANTARA_RE},
		{"CLAIM_AMOUNT", &b.CLAIM_AMOUNT}, {"WPC", &b.WPC},
		{"PL_NUMBER", &b.PL_NUMBER}, {"DISEASE", &b.DISEASE},
		{"ICD_CODE", &b.ICD_CODE}, {"NOTES", &b.NOTES},
		{"CEDINGCO", &b.CEDINGCO}, {"CEDINGCONAME", &b.CEDINGCONAME},
		{"SOB", &b.SOB}, {"SOBNAME", &b.SOBNAME},
		{"BUSINESSID", &b.BUSINESSID}, {"BUSINESSNAME", &b.BUSINESSNAME},
		{"SHARE_RETRO", &b.SHARE_RETRO}, {"KETERANGAN", &b.KETERANGAN},
		{"ACCEPTATION_DATE", &b.ACCEPTATION_DATE},
		{"NO_ACCEPTATION", &b.NO_ACCEPTATION}, {"STS_REJECT", &b.STS_REJECT},
		{"CLAIM_RETRO", &b.CLAIM_RETRO}, {"RETROID", &b.RETROID},
		{"RETRONAME", &b.RETRONAME},
		{"SECURITYREINSURERID", &b.SECURITYREINSURERID},
		{"SECURITYREINSURER", &b.SECURITYREINSURER},
		{"TYPECEDING", &b.TYPECEDING}, {"TYPE", &b.TYPE},
		{"CONFIRMATION_DATE", &b.CONFIRMATION_DATE},
		{"CLAIM_RECEIVED_DATE", &b.CLAIM_RECEIVED_DATE},
		{"COMPLETE_DATE", &b.COMPLETE_DATE}, {"ID", &b.ID},
		{"NAME_OF_BANK", &b.NAME_OF_BANK}, {"IDBANK", &b.IDBANK},
		{"ACCOUNTNO", &b.ACCOUNTNO}, {"CREATEOPNAME", &b.CREATEOPNAME},
		{"PRODUCTNAMEID", &b.PRODUCTNAMEID}, {"PRODUCTNAME", &b.PRODUCTNAME},
		{"RETROCEDED_SHARE", &b.RETROCEDED_SHARE},
	}
}

// NamaKolomBarisLama mengembalikan kelima puluh lima nama kolom, berurut.
//
// Dipakai penulis fixture supaya daftar kolomnya tidak ditulis ulang.
func NamaKolomBarisLama() []string {
	var b BarisLama
	medan := medanBarisLama(&b)
	out := make([]string, 0, len(medan))
	for _, m := range medan {
		out = append(out, m.Kolom)
	}
	return out
}

// ekspresiSelectLama menyusun bagian SELECT untuk kelima puluh lima kolom.
func ekspresiSelectLama() string {
	var b BarisLama
	medan := medanBarisLama(&b)
	ekspresi := make([]string, 0, len(medan))
	for _, m := range medan {
		ekspresi = append(ekspresi, ekspresiBacaLama(m.Kolom))
	}
	return strings.Join(ekspresi, ", ")
}

// tujuanScanLama menyiapkan penampung Scan untuk satu baris.
//
// Nilainya ditampung sebagai NullString lebih dulu, sebab kolom mana pun boleh
// NULL. NULL menjadi teks kosong - yang di BarisLama memang berarti "tidak ada
// nilai", bukan nol (ADR-U-0027).
func tujuanScanLama(n int) ([]sql.NullString, []any) {
	nilai := make([]sql.NullString, n)
	tujuan := make([]any, n)
	for i := range nilai {
		tujuan[i] = &nilai[i]
	}
	return nilai, tujuan
}

// salinKeBarisLama memindahkan hasil Scan ke medan strukturnya.
func salinKeBarisLama(b *BarisLama, nilai []sql.NullString) {
	for i, m := range medanBarisLama(b) {
		if i < len(nilai) {
			*m.Nilai = nilai[i].String
		}
	}
}

// TipeKolomBarisLama menyebut tipe Oracle untuk satu kolom tabel tiruan.
//
// ⚠️ [data DBA] Ini tipe TABEL TIRUAN, bukan tipe tabel sungguhan - lihat
// catatan kepala berkas. Ia sengaja dibuat setia pada dugaan terbaik atas
// bentuk aslinya: kalau tiruannya berisi teks semua, test pulang-pergi tidak
// menguji apa pun tentang TO_CHAR, dan jebakan NLS justru lolos.
func TipeKolomBarisLama(kolom string) string {
	if t, ada := tipeWarisan[kolom]; ada {
		return t
	}
	// Kolom di luar katalog tidak boleh dikarang bentuknya. Pemanggilnya
	// hanya skema uji, dan daftarnya datang dari kolomWarisan juga, jadi
	// cabang ini hanya tercapai bila seseorang memanggilnya dengan nama
	// karangan - dan teks ini yang akan muncul di DDL-nya.
	return ""
}

// PenampungTulisLama menyusun penampung VALUES untuk satu kolom.
//
// Kolom tanggal dibungkus TO_DATE dengan bentuk yang sama persis dengan yang
// dipakai saat membaca, sehingga tulis dan baca bertemu di satu bentuk dan
// tidak satu pun bergantung pada NLS_DATE_FORMAT sesi.
func PenampungTulisLama(kolom string, urutan int) string {
	if kolomTanggalLama[kolom] {
		return fmt.Sprintf(`TO_DATE(:%d, 'YYYY-MM-DD HH24:MI:SS')`, urutan)
	}
	return fmt.Sprintf(":%d", urutan)
}

// NilaiBarisLama mengeluarkan kelima puluh lima nilai, berurut sama dengan
// NamaKolomBarisLama.
//
// Teks kosong menjadi NULL, bukan menjadi nol dan bukan menjadi teks kosong
// (ADR-U-0027). Tanggal dinormalkan lebih dulu ke satu bentuk supaya TO_DATE
// tidak pernah menemui bentuk yang tidak dikenalnya.
func NilaiBarisLama(b BarisLama) []any {
	medan := medanBarisLama(&b)
	out := make([]any, 0, len(medan))
	for _, m := range medan {
		teks := strings.TrimSpace(*m.Nilai)
		if teks == "" {
			out = append(out, nil)
			continue
		}
		if kolomTanggalLama[m.Kolom] {
			t, err := utils.ParseTanggal(teks)
			if err != nil {
				// Fixture yang tanggalnya tidak terbaca ditulis sebagai NULL,
				// bukan ditebak. Test yang memerlukannya akan gagal dengan
				// terang, dan itu memang yang diinginkan.
				out = append(out, nil)
				continue
			}
			out = append(out, utils.FormatTanggalWaktu(t))
			continue
		}
		out = append(out, teks)
	}
	return out
}
