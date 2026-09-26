package repository

// Daftar kolom tabel datar warisan - satu sumber untuk semua pemakainya.
//
// Untuk apa berkas ini: OS_AKSEPTASI_KLAIM_LIFE punya 55 kolom, dan tiga tempat
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

// kolomAngkaLama menyebut kolom yang isinya bilangan.
//
// Ia dibaca lewat TO_CHAR ber-argumen NLS supaya tiba sebagai TEKS dengan titik
// sebagai pemisah desimal, apa pun setelan sesi (ADR-U-0003, ADR-U-0016).
// Tanpa itu, sesi ber-NLS koma menyerahkan "1234,56" dan seluruh pembacaan
// uang rusak diam-diam.
var kolomAngkaLama = map[string]bool{
	"AGE": true, "EM_PERCENT": true, "SUM_INSURED": true,
	"CEDING_RETENTION": true, "SUM_REASURED": true, "SHARE_NUSANTARA_RE": true,
	"CLAIM_AMOUNT": true, "SHARE_RETRO": true, "RETROCEDED_SHARE": true,
}

// kolomTanggalLama menyebut kolom yang isinya tanggal.
//
// Dibaca dengan bentuk yang sama persis dengan yang dikenal utils.ParseTanggal,
// supaya tidak ada tanggal yang bergantung pada NLS_DATE_FORMAT sesi.
var kolomTanggalLama = map[string]bool{
	"DOB": true, "BEGIN_DATE": true, "LAPSE_DATE": true, "EXPIRED_DATE": true,
	"ACCEPTATION_DATE": true, "CONFIRMATION_DATE": true,
	"CLAIM_RECEIVED_DATE": true, "COMPLETE_DATE": true,
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
	switch {
	case kolom == "AGE":
		return "NUMBER(5)"
	case kolomAngkaLama[kolom]:
		return "NUMBER(38,8)"
	case kolomTanggalLama[kolom]:
		return "DATE"
	default:
		return "VARCHAR2(255)"
	}
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
