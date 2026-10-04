//go:build db

package repository_test

// Bukti Oracle untuk tiket berskema yang belum punya buktinya — ronde 8.
//
// ⛔ Berkas ini TIDAK memperkenalkan cara baru. Ia memakai `siapkan(t)`,
// `wajibTolak`, `wajibTerima`, dan `dasar(...)` dari `invarian_db_test.go`:
// satu transaksi per test, diakhiri `Rollback`, nol `Commit`, nol `DROP`,
// nol `DELETE` di luar transaksi.
//
// ⭐ TIAP TIKET MENDAPAT DUA UJI, dan keduanya wajib.
//
// Uji negatif saja TIDAK CUKUP, dan papan mencatat sebabnya dengan tiga
// korban: ketiga kekeliruan lingkup di modul ini - `INV-12` tanpa mata uang,
// `INV-08`, `INV-09` - LULUS setiap uji negatif yang ditulis untuknya.
// Constraint yang menolak TERLALU BANYAK menolak duplikat dengan benar; yang
// menangkapnya baris SAH yang ia tolak. Karena itu tiap kunci alami di bawah
// diuji dua arah: kembar ditolak, dan baris yang HANYA berbeda pada ruas
// terakhir kuncinya DITERIMA.
//
// ⚠️ Pengenal 9000xxx dan kode berawalan ZZ: bila sebuah `Rollback` gagal,
// barisnya harus mudah dikenali sebagai milik uji.

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"testing"
)

// Pernyataan penyisip yang dipakai berulang. Nilainya minimal - hanya kolom
// WAJIB ISI, sebab yang diuji constraint-nya, bukan kelengkapan datanya.
const (
	insKelompok = `INSERT INTO {skema}.KELOMPOK_TREATY (ID_KELOMPOK_TREATY,KODE,NAMA,AKTIF) VALUES (%d,'%s','Uji','1')`
	insKelas    = `INSERT INTO {skema}.KELAS_BISNIS (ID_KELAS_BISNIS,KODE,NAMA,AKTIF) VALUES (%d,'%s','Uji','1')`

	insRetensi = `INSERT INTO {skema}.RETENSI_CEDANT (ID_RETENSI_CEDANT,ID_VERSI_KONTRAK,ID_KELOMPOK_TREATY,NILAI_RETENSI,KODE_MATA_UANG)
	              VALUES (%d,%d,%d,1000,'%s')`
	insEgnpi = `INSERT INTO {skema}.EGNPI (ID_EGNPI,ID_VERSI_KONTRAK,ID_KELOMPOK_TREATY,NILAI_EGNPI,KODE_MATA_UANG)
	            VALUES (%d,%d,%d,2000,'%s')`
	insPortofolio = `INSERT INTO {skema}.PORTOFOLIO (ID_PORTOFOLIO,ID_VERSI_KONTRAK,ARAH_PORTOFOLIO,JENIS_PORTOFOLIO)
	                 VALUES (%d,%d,'%s','%s')`
	insPeriodeLapor = `INSERT INTO {skema}.PERIODE_PELAPORAN (ID_PERIODE_PELAPORAN,ID_VERSI_KONTRAK,PERIODE,TANGGAL_AWAL,BATAS_PENYERAHAN,BATAS_KONFIRMASI,BATAS_PELUNASAN)
	                   VALUES (%d,%d,'%s',DATE '2026-01-01',DATE '2026-02-01',DATE '2026-03-01',DATE '2026-04-01')`
	insPeriodeAkum = `INSERT INTO {skema}.PERIODE_AKUMULASI (ID_PERIODE_AKUMULASI,ID_VERSI_KONTRAK,PERIODE,TANGGAL_LAPOR)
	                  VALUES (%d,%d,'%s',DATE '2026-06-30')`
	insTermin = `INSERT INTO {skema}.TERMIN (ID_TERMIN,ID_VERSI_KONTRAK,NOMOR_TERMIN,PERSEN_TERMIN,KODE_MATA_UANG,TANGGAL_JATUH_TEMPO)
	             VALUES (%d,%d,%d,25,'%s',DATE '2026-03-31')`
	insSkala = `INSERT INTO {skema}.SKALA_KOASURANSI (ID_SKALA_KOASURANSI,ID_VERSI_KONTRAK,PERSEN_LIMIT,PERSEN_BAGIAN)
	            VALUES (%d,%d,%s,10)`
	insDokumen = `INSERT INTO {skema}.DOKUMEN_KONTRAK (ID_DOKUMEN_KONTRAK,ID_VERSI_KONTRAK,ID_DOKUMEN,TANGGAL_LAMPIR)
	              VALUES (%d,%d,%d,DATE '2026-05-05')`
	insLayer = `INSERT INTO {skema}.LAYER (ID_LAYER,ID_VERSI_KONTRAK,NOMOR_LAYER,JENIS_LAYER,LIMIT,KODE_MATA_UANG,DEDUCTIBLE,MDP_DIGABUNG,TANPA_HITUNG_PREMI_PEMULIHAN)
	            VALUES (%d,%d,%d,'XOL',1000000,'ZZU1',50000,'0','0')`
	insDetailProp = `INSERT INTO {skema}.DETAIL_PROPORSIONAL (ID_DETAIL_PROPORSIONAL,ID_LAYER,ID_KELOMPOK_TREATY,JENIS_TREATY)
	                 VALUES (%d,%d,%d,'QUOTA_SHARE')`
	insJejak = `INSERT INTO {skema}.JEJAK_PERUBAHAN (ID_JEJAK_PERUBAHAN,ID_VERSI_KONTRAK,WAKTU_PERUBAHAN,PELAKU,PERAN_PELAKU,RUAS_YANG_BERUBAH)
	            VALUES (%d,%d,SYSDATE,'UJI','ADMIN','%s')`
	insArsip = `INSERT INTO {skema}.ARSIP_MUATAN_KELUAR (ID_ARSIP_MUATAN_KELUAR,ID_KONTRAK,TUJUAN,DIKIRIM_PADA,BERHASIL,MUATAN)
	            VALUES (%d,%d,'PEGA_TREATY_IN',SYSDATE,'1','%s')`
	setVersiDasar = `UPDATE {skema}.VERSI_KONTRAK SET ID_VERSI_KONTRAK_DASAR = %d WHERE ID_VERSI_KONTRAK = %d`
)

// pondasi menyiapkan mata uang acuan, kontrak, dan SATU versi — bahan yang
// hampir setiap uji di bawah butuhkan.
func pondasi(t *testing.T, ctx context.Context, db *sql.Tx, skema string) {
	t.Helper()
	dasar(t, ctx, db, skema)
	wajibTerima(t, ctx, db, skema, "versi pertama",
		fmt.Sprintf(insVersi, 9000001, 9000001, 1, 9000001))
}

// ---------------------------------------------------------------- tiket 20

// ------------------------------------------------------------ tiket 22, 23

// Tiket 22 — INV-08: kelompok treaty + mata uang unik di dalam satu versi.
func TestTiket22RetensiCedantGandaDitolak(t *testing.T) {
	db, skema, ctx := siapkan(t)
	pondasi(t, ctx, db, skema)
	wajibTerima(t, ctx, db, skema, "kelompok treaty", fmt.Sprintf(insKelompok, 9000001, "ZZK1"))

	wajibTerima(t, ctx, db, skema, "retensi pertama",
		fmt.Sprintf(insRetensi, 9000001, 9000001, 9000001, "IDR"))
	wajibTolak(t, ctx, db, skema, "UQ_RETENSI_CEDANT",
		fmt.Sprintf(insRetensi, 9000002, 9000001, 9000001, "IDR"))
	wajibTolak(t, ctx, db, skema, "FK_RETENSI_CEDANT_2",
		fmt.Sprintf(insRetensi, 9000003, 9000001, 9009999, "IDR"))
}

// Uji POSITIF tiket 22 — kunci alaminya BERTIGA ruas, dan mata uang ada di
// dalamnya. Retensi kelompok yang sama dalam dua mata uang adalah dua fakta
// yang sah; `UQ` tanpa mata uang akan menolak yang kedua.
func TestTiket22RetensiDuaMataUangDiterima(t *testing.T) {
	db, skema, ctx := siapkan(t)
	pondasi(t, ctx, db, skema)
	wajibTerima(t, ctx, db, skema, "kelompok treaty", fmt.Sprintf(insKelompok, 9000001, "ZZK1"))

	wajibTerima(t, ctx, db, skema, "retensi IDR",
		fmt.Sprintf(insRetensi, 9000001, 9000001, 9000001, "IDR"))
	wajibTerima(t, ctx, db, skema, "retensi USD, kelompok SAMA",
		fmt.Sprintf(insRetensi, 9000002, 9000001, 9000001, "USD"))
}

// Tiket 23 — INV-09, bentuknya sama dengan 22 dan diuji terpisah sebab
// tabelnya terpisah: kekeliruan pada satu tidak membuktikan yang lain.
func TestTiket23EgnpiGandaDitolak(t *testing.T) {
	db, skema, ctx := siapkan(t)
	pondasi(t, ctx, db, skema)
	wajibTerima(t, ctx, db, skema, "kelompok treaty", fmt.Sprintf(insKelompok, 9000001, "ZZK1"))

	wajibTerima(t, ctx, db, skema, "egnpi pertama",
		fmt.Sprintf(insEgnpi, 9000001, 9000001, 9000001, "IDR"))
	wajibTolak(t, ctx, db, skema, "UQ_EGNPI",
		fmt.Sprintf(insEgnpi, 9000002, 9000001, 9000001, "IDR"))
}

// Uji POSITIF tiket 23 — dua mata uang diterima, DAN `ID_KELAS_BISNIS` boleh
// kosong. Kolom nullable yang diam-diam menjadi wajib adalah cacat yang tidak
// terlihat sampai baris pertama tanpa kelas bisnis ditolak.
func TestTiket23EgnpiDuaMataUangDanKelasBisnisKosongDiterima(t *testing.T) {
	db, skema, ctx := siapkan(t)
	pondasi(t, ctx, db, skema)
	wajibTerima(t, ctx, db, skema, "kelompok treaty", fmt.Sprintf(insKelompok, 9000001, "ZZK1"))

	wajibTerima(t, ctx, db, skema, "egnpi IDR tanpa kelas bisnis",
		fmt.Sprintf(insEgnpi, 9000001, 9000001, 9000001, "IDR"))
	wajibTerima(t, ctx, db, skema, "egnpi USD, kelompok SAMA",
		fmt.Sprintf(insEgnpi, 9000002, 9000001, 9000001, "USD"))
}

// ---------------------------------------------------------------- tiket 24

// Tiket 24 — INV-66: arah + jenis portofolio unik di dalam satu versi.
func TestTiket24PortofolioGandaDitolak(t *testing.T) {
	db, skema, ctx := siapkan(t)
	pondasi(t, ctx, db, skema)

	wajibTerima(t, ctx, db, skema, "portofolio masuk",
		fmt.Sprintf(insPortofolio, 9000001, 9000001, "MASUK", "PREMI"))
	wajibTolak(t, ctx, db, skema, "UQ_PORTOFOLIO",
		fmt.Sprintf(insPortofolio, 9000002, 9000001, "MASUK", "PREMI"))
}

// Uji POSITIF tiket 24 — kuncinya BERTIGA ruas. Arah yang sama dengan jenis
// berbeda, dan jenis yang sama dengan arah berbeda, keduanya sah; kunci yang
// hanya memakai arah akan menolak keduanya.
func TestTiket24PortofolioEmpatKombinasiDiterima(t *testing.T) {
	db, skema, ctx := siapkan(t)
	pondasi(t, ctx, db, skema)

	n := 9000001
	for _, p := range [][2]string{{"MASUK", "PREMI"}, {"MASUK", "KLAIM"}, {"KELUAR", "PREMI"}, {"KELUAR", "KLAIM"}} {
		wajibTerima(t, ctx, db, skema, "portofolio "+p[0]+"/"+p[1],
			fmt.Sprintf(insPortofolio, n, 9000001, p[0], p[1]))
		n++
	}
}

// ------------------------------------------------------------ tiket 25, 26

// Tiket 25 — INV-10: periode unik di dalam satu versi.
func TestTiket25PeriodePelaporanGandaDitolak(t *testing.T) {
	db, skema, ctx := siapkan(t)
	pondasi(t, ctx, db, skema)

	wajibTerima(t, ctx, db, skema, "periode Q1",
		fmt.Sprintf(insPeriodeLapor, 9000001, 9000001, "2026-Q1"))
	wajibTolak(t, ctx, db, skema, "UQ_PERIODE_PELAPORAN",
		fmt.Sprintf(insPeriodeLapor, 9000002, 9000001, "2026-Q1"))
}

// Uji POSITIF tiket 25 — empat triwulan pada satu versi diterima. Kontrak
// berperiode triwulanan adalah keadaan biasa, dan kunci yang terlalu sempit
// akan membuatnya mustahil.
func TestTiket25EmpatPeriodePelaporanDiterima(t *testing.T) {
	db, skema, ctx := siapkan(t)
	pondasi(t, ctx, db, skema)

	for i, p := range []string{"2026-Q1", "2026-Q2", "2026-Q3", "2026-Q4"} {
		wajibTerima(t, ctx, db, skema, "periode "+p,
			fmt.Sprintf(insPeriodeLapor, 9000001+i, 9000001, p))
	}
}

// Tiket 26 — INV-11, tabel terpisah dari 25 dan diuji terpisah.
func TestTiket26PeriodeAkumulasiGandaDitolak(t *testing.T) {
	db, skema, ctx := siapkan(t)
	pondasi(t, ctx, db, skema)

	wajibTerima(t, ctx, db, skema, "akumulasi H1",
		fmt.Sprintf(insPeriodeAkum, 9000001, 9000001, "2026-H1"))
	wajibTolak(t, ctx, db, skema, "UQ_PERIODE_AKUMULASI",
		fmt.Sprintf(insPeriodeAkum, 9000002, 9000001, "2026-H1"))
}

// Uji POSITIF tiket 26 — dan ia memeriksa dua hal sekaligus: periode kedua
// diterima, dan kedua kolom opsionalnya (`HARI_BATAS_PENYERAHAN`,
// `BATAS_PENYERAHAN`) boleh kosong.
func TestTiket26AkumulasiKeduaDanKolomOpsionalKosongDiterima(t *testing.T) {
	db, skema, ctx := siapkan(t)
	pondasi(t, ctx, db, skema)

	wajibTerima(t, ctx, db, skema, "akumulasi H1",
		fmt.Sprintf(insPeriodeAkum, 9000001, 9000001, "2026-H1"))
	wajibTerima(t, ctx, db, skema, "akumulasi H2 tanpa batas penyerahan",
		fmt.Sprintf(insPeriodeAkum, 9000002, 9000001, "2026-H2"))
}

// ---------------------------------------------------------------- tiket 27

// Tiket 27 — INV-12, DAN ia yang koreksinya paling mahal dibayar: kunci
// alaminya semula "nomor termin" saja, dan itu menolak data yang SAH.
func TestTiket27TerminGandaDitolak(t *testing.T) {
	db, skema, ctx := siapkan(t)
	pondasi(t, ctx, db, skema)

	wajibTerima(t, ctx, db, skema, "termin 1 IDR",
		fmt.Sprintf(insTermin, 9000001, 9000001, 1, "IDR"))
	wajibTolak(t, ctx, db, skema, "UQ_TERMIN",
		fmt.Sprintf(insTermin, 9000002, 9000001, 1, "IDR"))
}

// ⭐ Uji POSITIF tiket 27 — INILAH uji yang gagal pada rumusan `INV-12` yang
// lama, dan papan menuntutnya DIJALANKAN, bukan diargumentasikan: termin
// nomor 1 dalam IDR dan termin nomor 1 dalam USD, pada versi yang sama,
// KEDUANYA sah.
func TestTiket27TerminNomorSamaDuaMataUangDiterima(t *testing.T) {
	db, skema, ctx := siapkan(t)
	pondasi(t, ctx, db, skema)

	wajibTerima(t, ctx, db, skema, "termin 1 IDR",
		fmt.Sprintf(insTermin, 9000001, 9000001, 1, "IDR"))
	wajibTerima(t, ctx, db, skema, "termin 1 USD — nomor SAMA, mata uang beda",
		fmt.Sprintf(insTermin, 9000002, 9000001, 1, "USD"))
}

// ---------------------------------------------------------------- tiket 28

// Tiket 28 — INV-13: persen limit unik di dalam satu versi.
func TestTiket28SkalaKoasuransiGandaDitolak(t *testing.T) {
	db, skema, ctx := siapkan(t)
	pondasi(t, ctx, db, skema)

	wajibTerima(t, ctx, db, skema, "skala 25%",
		fmt.Sprintf(insSkala, 9000001, 9000001, "25"))
	wajibTolak(t, ctx, db, skema, "UQ_SKALA_KOASURANSI",
		fmt.Sprintf(insSkala, 9000002, 9000001, "25"))
}

// Uji POSITIF tiket 28 — ko-asuransi memang BEBERAPA BARIS, dan itu seluruh
// pokok tiketnya. Tiga tingkat pada satu versi diterima.
func TestTiket28TigaBarisSkalaKoasuransiDiterima(t *testing.T) {
	db, skema, ctx := siapkan(t)
	pondasi(t, ctx, db, skema)

	for i, p := range []string{"25", "50", "75"} {
		wajibTerima(t, ctx, db, skema, "skala "+p+"%",
			fmt.Sprintf(insSkala, 9000001+i, 9000001, p))
	}
}

// ---------------------------------------------------------------- tiket 30

// Tiket 30 — INV-67: dokumen unik di dalam satu versi.
func TestTiket30DokumenGandaDitolak(t *testing.T) {
	db, skema, ctx := siapkan(t)
	pondasi(t, ctx, db, skema)

	wajibTerima(t, ctx, db, skema, "dokumen pertama",
		fmt.Sprintf(insDokumen, 9000001, 9000001, 7001))
	wajibTolak(t, ctx, db, skema, "UQ_DOKUMEN_KONTRAK",
		fmt.Sprintf(insDokumen, 9000002, 9000001, 7001))
}

// Uji POSITIF tiket 30 — dokumen yang SAMA boleh dilampirkan ke versi yang
// BERBEDA. `ADR-0027`: dokumen DIRUJUK, tidak dimiliki — dan rujukan yang
// hanya boleh sekali di seluruh tabel akan menolak lampiran ulang pada versi
// berikutnya, yang justru keadaan normal sebuah addendum.
func TestTiket30DokumenSamaPadaDuaVersiDiterima(t *testing.T) {
	db, skema, ctx := siapkan(t)
	pondasi(t, ctx, db, skema)
	wajibTerima(t, ctx, db, skema, "versi kedua",
		fmt.Sprintf(insVersi, 9000002, 9000001, 2, 9000001))

	wajibTerima(t, ctx, db, skema, "dokumen pada versi 1",
		fmt.Sprintf(insDokumen, 9000001, 9000001, 7001))
	wajibTerima(t, ctx, db, skema, "dokumen SAMA pada versi 2",
		fmt.Sprintf(insDokumen, 9000002, 9000002, 7001))
}

// ---------------------------------------------------------------- tiket 34

// Tiket 34 — INV-06: kelompok treaty unik di dalam satu layer.
func TestTiket34DetailProporsionalGandaDitolak(t *testing.T) {
	db, skema, ctx := siapkan(t)
	pondasi(t, ctx, db, skema)
	wajibTerima(t, ctx, db, skema, "kelompok treaty", fmt.Sprintf(insKelompok, 9000001, "ZZK1"))
	wajibTerima(t, ctx, db, skema, "layer", fmt.Sprintf(insLayer, 9000001, 9000001, 1))

	wajibTerima(t, ctx, db, skema, "ketentuan pertama",
		fmt.Sprintf(insDetailProp, 9000001, 9000001, 9000001))
	wajibTolak(t, ctx, db, skema, "UQ_DETAIL_PROPORSIONAL",
		fmt.Sprintf(insDetailProp, 9000002, 9000001, 9000001))
}

// Uji POSITIF tiket 34 — DUA kelompok di dalam SATU layer diterima. Itu
// seluruh pokok tiketnya: ketentuan proporsional per kelompok DI DALAM layer.
func TestTiket34DuaKelompokDalamSatuLayerDiterima(t *testing.T) {
	db, skema, ctx := siapkan(t)
	pondasi(t, ctx, db, skema)
	wajibTerima(t, ctx, db, skema, "kelompok 1", fmt.Sprintf(insKelompok, 9000001, "ZZK1"))
	wajibTerima(t, ctx, db, skema, "kelompok 2", fmt.Sprintf(insKelompok, 9000002, "ZZK2"))
	wajibTerima(t, ctx, db, skema, "layer", fmt.Sprintf(insLayer, 9000001, 9000001, 1))

	wajibTerima(t, ctx, db, skema, "ketentuan kelompok 1",
		fmt.Sprintf(insDetailProp, 9000001, 9000001, 9000001))
	wajibTerima(t, ctx, db, skema, "ketentuan kelompok 2, layer SAMA",
		fmt.Sprintf(insDetailProp, 9000002, 9000001, 9000002))
}

// ---------------------------------------------------------------- tiket 39

// Tiket 39 — jejak perubahan: yatim ditolak, dan `tolak` pada induknya
// terbukti. `ERD.md` §2.3 menandainya `[hapus: tolak]`, dan sebabnya tertulis:
// jejak yang dapat dihapus bersama bendanya bukan jejak.
func TestTiket39JejakYatimDitolakDanVersiBerjejakTidakDapatDihapus(t *testing.T) {
	db, skema, ctx := siapkan(t)
	pondasi(t, ctx, db, skema)

	wajibTolak(t, ctx, db, skema, "FK_JEJAK_PERUBAHAN_1",
		fmt.Sprintf(insJejak, 9000001, 9009999, "NAMA_KONTRAK"))

	wajibTerima(t, ctx, db, skema, "jejak pada versi yang ada",
		fmt.Sprintf(insJejak, 9000002, 9000001, "NAMA_KONTRAK"))
	// ⭐ Inilah yang membuktikan `tolak`, bukan `ikut hapus`.
	wajibTolak(t, ctx, db, skema, "FK_JEJAK_PERUBAHAN_1",
		fmt.Sprintf(`DELETE FROM {skema}.VERSI_KONTRAK WHERE ID_VERSI_KONTRAK = %d`, 9000001))
}

// Uji POSITIF tiket 39 — peristiwa yang SAMA dua kali pada versi yang sama
// DITERIMA. Entitas ini sengaja TANPA kunci alami (`Z00_KUNCI_ALAMI.sql`:
// "peristiwa yang sama dapat terjadi dua kali"), dan `UNIQUE` yang
// disisipkan diam-diam akan menolak baris kedua ini.
func TestTiket39DuaJejakRuasSamaPadaVersiSamaDiterima(t *testing.T) {
	db, skema, ctx := siapkan(t)
	pondasi(t, ctx, db, skema)

	wajibTerima(t, ctx, db, skema, "jejak pertama",
		fmt.Sprintf(insJejak, 9000001, 9000001, "NAMA_KONTRAK"))
	wajibTerima(t, ctx, db, skema, "jejak kedua, ruas SAMA",
		fmt.Sprintf(insJejak, 9000002, 9000001, "NAMA_KONTRAK"))
}

// ---------------------------------------------------------------- tiket 42

// Tiket 42 — arsip: yatim ditolak, dan kontrak berarsip TIDAK dapat dihapus.
// Migrasi 425 memilih `tolak` dengan sebab tertulis: arsip yang lenyap
// bersama kontraknya berhenti menjadi arsip tepat saat ia paling dibutuhkan.
func TestTiket42ArsipYatimDitolakDanKontrakBerarsipTidakDapatDihapus(t *testing.T) {
	db, skema, ctx := siapkan(t)
	pondasi(t, ctx, db, skema)

	wajibTolak(t, ctx, db, skema, "FK_ARSIP_MUATAN_KELUAR_1",
		fmt.Sprintf(insArsip, 9000001, 9009999, "{}"))

	// ⚠️ Kontrak KEDUA, TANPA versi — dan itu bukan kerapian, itu syarat.
	// Percobaan pertama memakai kontrak 9000001 dan GAGAL dengan benar:
	// penghapusannya ditolak `FK_VERSI_KONTRAK_1` LEBIH DULU, sebab ia punya
	// versi. Uji itu akan "lulus" sebagai penolakan sambil nol membuktikan
	// tentang arsip — persis uji yang lulus secara kebetulan yang
	// `wajibTolak` ada untuk menangkapnya. Yang membuktikan FK arsip hanya
	// kontrak yang arsipnya SATU-SATUNYA anaknya.
	wajibTerima(t, ctx, db, skema, "kontrak tanpa versi", fmt.Sprintf(insKontrak, 9000002))
	wajibTerima(t, ctx, db, skema, "arsip pada kontrak itu",
		fmt.Sprintf(insArsip, 9000002, 9000002, "{}"))
	wajibTolak(t, ctx, db, skema, "FK_ARSIP_MUATAN_KELUAR_1",
		fmt.Sprintf(`DELETE FROM {skema}.KONTRAK WHERE ID_KONTRAK = %d`, 9000002))
}

// Uji POSITIF tiket 42 — muatan yang TIDAK dapat diurai tetap tersimpan, dan
// dua pengiriman ke kontrak yang sama diterima. Arsip yang hanya menerima
// yang sudah benar tidak melestarikan apa pun yang berguna; arsip yang hanya
// menerima satu pengiriman bukan arsip.
func TestTiket42MuatanRusakDanPengirimanKeduaDiterima(t *testing.T) {
	db, skema, ctx := siapkan(t)
	pondasi(t, ctx, db, skema)

	wajibTerima(t, ctx, db, skema, "muatan rusak",
		fmt.Sprintf(insArsip, 9000001, 9000001, `{"TreatyIn": <<rusak`))
	wajibTerima(t, ctx, db, skema, "pengiriman kedua",
		fmt.Sprintf(insArsip, 9000002, 9000001, "{}"))
}

// ------------------------------------------------- tiket 01 papan Adjustment

// Tiket 01 — `ID_VERSI_KONTRAK_DASAR`: rujukan ke versi yang menjadi dasarnya.
// `ERD.md` §2.2 menandainya `[hapus: tolak]`, sebab menghapus versi dasar
// membuat seluruh baris selisih kehilangan artinya.
func TestTiket01VersiDasarTakAdaDitolakDanVersiDasarTidakDapatDihapus(t *testing.T) {
	db, skema, ctx := siapkan(t)
	pondasi(t, ctx, db, skema)
	wajibTerima(t, ctx, db, skema, "versi kedua",
		fmt.Sprintf(insVersi, 9000002, 9000001, 2, 9000001))

	wajibTolak(t, ctx, db, skema, "FK_VERSI_KONTRAK_DASAR",
		fmt.Sprintf(setVersiDasar, 9009999, 9000002))

	wajibTerima(t, ctx, db, skema, "versi 2 berdasar versi 1",
		fmt.Sprintf(setVersiDasar, 9000001, 9000002))
	// ⭐ `tolak` terbukti: versi yang menjadi DASAR tidak dapat dihapus.
	wajibTolak(t, ctx, db, skema, "FK_VERSI_KONTRAK_DASAR",
		fmt.Sprintf(`DELETE FROM {skema}.VERSI_KONTRAK WHERE ID_VERSI_KONTRAK = %d`, 9000001))
}

// Uji POSITIF tiket 01 — kolomnya BOLEH KOSONG, dan itu bukan kelalaian:
// versi PERTAMA sebuah kontrak memang tidak punya dasar. Constraint
// `NOT NULL` yang terlanjur dipasang akan membuat kontrak baru mustahil.
func TestTiket01VersiPertamaTanpaDasarDiterima(t *testing.T) {
	db, skema, ctx := siapkan(t)
	pondasi(t, ctx, db, skema)

	n := cacah(t, ctx, db, skema, fmt.Sprintf(
		`SELECT COUNT(*) FROM {skema}.VERSI_KONTRAK WHERE ID_VERSI_KONTRAK = %d AND ID_VERSI_KONTRAK_DASAR IS NULL`,
		9000001))
	if n != 1 {
		t.Errorf("versi pertama tanpa dasar: cacah %d, mau 1", n)
	}
}

// ------------------------------------------------- tiket 54 (melepas 49, 55, 56)

const insCatatan = `INSERT INTO {skema}.CATATAN_PERSETUJUAN (ID_CATATAN_PERSETUJUAN,ID_VERSI_KONTRAK,WAKTU_KEPUTUSAN,NAMA_PEMUTUS,DISETUJUI,ALASAN)
	             VALUES (%d,%d,SYSDATE,'%s','%s',%s)`

// Tiket 54 — catatan yatim ditolak, DAN versi bercatat tidak dapat dihapus.
//
// ⭐ Penolakan kedua itulah pokoknya, dan ia yang membedakan entitas ini dari
// sepuluh saudara `VERSI_KONTRAK` lainnya: kesepuluhnya `ikut hapus`, yang
// ini `tolak`. `ERD.md` §2.3 menuliskan sebabnya: "jejak yang dapat dihapus
// bersama bendanya bukan jejak."
func TestTiket54CatatanYatimDitolakDanVersiBercatatTidakDapatDihapus(t *testing.T) {
	db, skema, ctx := siapkan(t)
	pondasi(t, ctx, db, skema)

	wajibTolak(t, ctx, db, skema, "FK_CATATAN_PERSETUJUAN_1",
		fmt.Sprintf(insCatatan, 9000001, 9009999, "UJI", "1", "NULL"))

	wajibTerima(t, ctx, db, skema, "catatan pada versi yang ada",
		fmt.Sprintf(insCatatan, 9000002, 9000001, "UJI", "1", "NULL"))
	// ⭐ `tolak` terbukti, dan bukan `ikut hapus`.
	wajibTolak(t, ctx, db, skema, "FK_CATATAN_PERSETUJUAN_1",
		fmt.Sprintf(`DELETE FROM {skema}.VERSI_KONTRAK WHERE ID_VERSI_KONTRAK = %d`, 9000001))
}

// Uji POSITIF tiket 54 — DUA keputusan oleh orang yang SAMA pada versi yang
// SAMA diterima, dan `ALASAN` boleh kosong.
//
// ⛔ Inilah uji yang gagal bila seseorang kelak menyisipkan `UNIQUE` di sini.
// §10.20 menyatakan entitas ini sengaja TANPA kunci alami: dua keputusan pada
// versi yang sama, oleh orang yang sama, pada hari yang sama adalah keadaan
// yang SAH — yang membedakan barisnya urutan waktu, bukan sebuah nilai.
func TestTiket54DuaKeputusanOrangSamaPadaVersiSamaDiterima(t *testing.T) {
	db, skema, ctx := siapkan(t)
	pondasi(t, ctx, db, skema)

	wajibTerima(t, ctx, db, skema, "keputusan pertama — menyetujui, tanpa alasan",
		fmt.Sprintf(insCatatan, 9000001, 9000001, "SITI", "1", "NULL"))
	wajibTerima(t, ctx, db, skema, "keputusan kedua — orang SAMA, versi SAMA",
		fmt.Sprintf(insCatatan, 9000002, 9000001, "SITI", "0", "'dikembalikan'"))
}

// ----------------------------------------------- tiket 71, 72, 73 (perkakas)

const (
	insKorelasi = `INSERT INTO {skema}.MIGRASI_KORELASI (ID_MIGRASI_KORELASI,KUNCI_PEGA,ID_PEGA,ID_WARISAN,ID_KONTRAK_BARU,DIPINDAHKAN_PADA)
	               VALUES (%d,'%s','%s','%s',%s,SYSDATE)`
	insPendaratan = `INSERT INTO {skema}.MIGRASI_PENDARATAN (ID_MIGRASI_PENDARATAN,KUNCI_WARISAN,MUATAN,MENDARAT_PADA)
	                 VALUES (%d,'%s','%s',SYSDATE)`
	insDitolak = `INSERT INTO {skema}.MIGRASI_NILAI_DITOLAK (ID_MIGRASI_NILAI_DITOLAK,ID_MIGRASI_PENDARATAN,JALUR_SIMPUL,NILAI_MENTAH,SEBAB_DITOLAK)
	              VALUES (%d,%d,'%s','%s','bukan angka')`
)

// ⭐ Uji POSITIF tiket 71, dan ia yang membuktikan keputusan "nol kunci asing"
// sungguh berlaku: baris korelasi yang menunjuk kontrak TIDAK ADA tetap
// DITERIMA, dan ia BERTAHAN sesudah kontraknya dihapus.
//
// Kunci asing akan menolak yang pertama dan menghapus yang kedua — dan
// keduanya menghancurkan guna tabel ini.
func TestTiket71KorelasiBertahanTanpaKunciAsing(t *testing.T) {
	db, skema, ctx := siapkan(t)
	dasar(t, ctx, db, skema)

	wajibTerima(t, ctx, db, skema, "korelasi ke kontrak yang TIDAK ADA",
		fmt.Sprintf(insKorelasi, 9000001, "PEGA-KEY-1", "TI-1", "OLD-1", "9009999"))

	wajibTerima(t, ctx, db, skema, "korelasi ke kontrak yang ada",
		fmt.Sprintf(insKorelasi, 9000002, "PEGA-KEY-2", "TI-2", "OLD-2", "9000001"))
	wajibTerima(t, ctx, db, skema, "kontrak dihapus",
		fmt.Sprintf(`DELETE FROM {skema}.KONTRAK WHERE ID_KONTRAK = %d`, 9000001))

	n := cacah(t, ctx, db, skema,
		`SELECT COUNT(*) FROM {skema}.MIGRASI_KORELASI WHERE ID_MIGRASI_KORELASI IN (9000001,9000002)`)
	if n != 2 {
		t.Errorf("korelasi bertahan: cacah %d, mau 2 — jejak asal-usul hilang bersama barisnya", n)
	}
}

// Uji NEGATIF tiket 71 — kolom wajibnya memang wajib. Baris tanpa cap waktu
// pemindahan tidak menjembatani apa pun yang dapat ditelusuri.
func TestTiket71KorelasiTanpaCapWaktuDitolak(t *testing.T) {
	db, skema, ctx := siapkan(t)
	wajibTolak(t, ctx, db, skema, "DIPINDAHKAN_PADA",
		`INSERT INTO {skema}.MIGRASI_KORELASI (ID_MIGRASI_KORELASI,ID_WARISAN) VALUES (9000001,'OLD-1')`)
}

// ⭐ Uji POSITIF tiket 72, dan ia yang menentukan: muatan yang TIDAK DAPAT
// DIURAI SAMA SEKALI tetap mendarat. Tabel yang hanya menerima yang sudah
// benar tidak melestarikan apa pun yang berguna — justru muatan rusak itulah
// yang paling perlu ditunjukkan nanti.
func TestTiket72MuatanRusakTetapMendarat(t *testing.T) {
	db, skema, ctx := siapkan(t)

	wajibTerima(t, ctx, db, skema, "muatan rusak di tengah",
		fmt.Sprintf(insPendaratan, 9000001, "OLD-1", `{"TreatyIn": <<rusak`))
	wajibTerima(t, ctx, db, skema, "dokumen warisan KEDUA berkunci sama",
		fmt.Sprintf(insPendaratan, 9000002, "OLD-1", "{}"))
}

// Uji NEGATIF tiket 72 — pendaratan tanpa muatan tidak mengarsipkan apa pun.
func TestTiket72PendaratanTanpaMuatanDitolak(t *testing.T) {
	db, skema, ctx := siapkan(t)
	wajibTolak(t, ctx, db, skema, "MUATAN",
		`INSERT INTO {skema}.MIGRASI_PENDARATAN (ID_MIGRASI_PENDARATAN,KUNCI_WARISAN,MENDARAT_PADA) VALUES (9000001,'OLD-1',SYSDATE)`)
}

// Tiket 73 — nilai ditolak yatim DITOLAK, dan pendaratan yang punya nilai
// ditolak TIDAK DAPAT dihapus.
//
// ⚠️ Penolakan kedua MENGOREKSI tiket 73, yang menyebut "ikut hapus". ERD
// baris 38 menulis ON DELETE-nya "di Go" — tidak meresepkan aturan basis
// data — dan `tolak` yang dipilih: catatan forensik yang lenyap bersama
// induknya berhenti menjadi catatan forensik.
func TestTiket73NilaiDitolakYatimDanPendaratanTakDapatDihapus(t *testing.T) {
	db, skema, ctx := siapkan(t)

	wajibTolak(t, ctx, db, skema, "FK_MIGRASI_NILAI_DITOLAK_1",
		fmt.Sprintf(insDitolak, 9000001, 9009999, ".Limits[0].MDP", "abc"))

	wajibTerima(t, ctx, db, skema, "pendaratan",
		fmt.Sprintf(insPendaratan, 9000001, "OLD-1", "{}"))
	wajibTerima(t, ctx, db, skema, "nilai ditolak pada pendaratan itu",
		fmt.Sprintf(insDitolak, 9000002, 9000001, ".Limits[0].MDP", "abc"))
	wajibTolak(t, ctx, db, skema, "FK_MIGRASI_NILAI_DITOLAK_1",
		fmt.Sprintf(`DELETE FROM {skema}.MIGRASI_PENDARATAN WHERE ID_MIGRASI_PENDARATAN = %d`, 9000001))
}

// ⭐ Uji POSITIF tiket 73 — nilai mentah yang BUKAN ANGKA SAMA SEKALI
// tersimpan apa adanya. Itu yang membuktikan "tidak dibulatkan": kolom
// bertipe angka akan menolaknya, dan pembulatan diam-diam mengubah angka yang
// pernah dibukukan menjadi angka yang tidak pernah ada.
func TestTiket73NilaiBukanAngkaTersimpanApaAdanya(t *testing.T) {
	db, skema, ctx := siapkan(t)
	wajibTerima(t, ctx, db, skema, "pendaratan",
		fmt.Sprintf(insPendaratan, 9000001, "OLD-1", "{}"))

	wajibTerima(t, ctx, db, skema, "nilai mentah berupa teks",
		fmt.Sprintf(insDitolak, 9000001, 9000001, ".A", "tidak diketahui"))
	wajibTerima(t, ctx, db, skema, "jalur simpul KEDUA pada pendaratan sama",
		fmt.Sprintf(insDitolak, 9000002, 9000001, ".B", "-"))

	n := cacah(t, ctx, db, skema,
		`SELECT COUNT(*) FROM {skema}.MIGRASI_NILAI_DITOLAK WHERE NILAI_MENTAH = 'tidak diketahui'`)
	if n != 1 {
		t.Errorf("nilai mentah bukan-angka: cacah %d, mau 1", n)
	}
}

// -------------------------------------------------- tiket 38 (cabang penyebaran)

const (
	insJenisReas = `INSERT INTO {skema}.JENIS_REASURANSI (ID_JENIS_REASURANSI,KODE,NAMA,AKTIF) VALUES (%d,'%s','Uji','1')`
	insBagian    = `INSERT INTO {skema}.BAGIAN (ID_BAGIAN,ID_LAYER,PERSEN_BAGIAN_NURE,PREMI_BRUTO) VALUES (%d,%d,50,1000000)`
	insPenyebar  = `INSERT INTO {skema}.PENYEBARAN (ID_PENYEBARAN,ID_BAGIAN,ID_DETAIL_PROPORSIONAL,ID_JENIS_REASURANSI,PERSEN_PENYEBARAN)
	                VALUES (%d,%s,%s,%d,100)`
	insRincian = `INSERT INTO {skema}.RINCIAN_PENYEBARAN (ID_RINCIAN_PENYEBARAN,ID_PENYEBARAN,ID_JENIS_REASURANSI,PERSEN_RINCIAN)
	              VALUES (%d,%d,%d,50)`
	insNilaiPenyebaran = `INSERT INTO {skema}.NILAI_PENYEBARAN (ID_NILAI_PENYEBARAN,ID_RINCIAN_PENYEBARAN,NILAI,KODE_MATA_UANG)
	                      VALUES (%d,%d,1000,'%s')`
)

// cabangPenyebaran menyiapkan layer, bagian, dan dua jenis reasuransi.
func cabangPenyebaran(t *testing.T, ctx context.Context, db *sql.Tx, skema string) {
	t.Helper()
	pondasi(t, ctx, db, skema)
	wajibTerima(t, ctx, db, skema, "layer", fmt.Sprintf(insLayer, 9000001, 9000001, 1))
	wajibTerima(t, ctx, db, skema, "bagian", fmt.Sprintf(insBagian, 9000001, 9000001))
	wajibTerima(t, ctx, db, skema, "jenis reasuransi 1", fmt.Sprintf(insJenisReas, 9000001, "ZZR1"))
	wajibTerima(t, ctx, db, skema, "jenis reasuransi 2", fmt.Sprintf(insJenisReas, 9000002, "ZZR2"))
}

// ⭐ Tiket 38 — `CK_PENYEBARAN_INDUK`: TEPAT SATU induk terisi.
//
// Dua arah diuji, dan keduanya perlu: NOL induk ditolak, DUA induk ditolak.
// `CHECK` yang hanya menolak nol akan meloloskan baris yang menggantung pada
// dua tempat sekaligus — dan baris itu akan muncul dua kali di setiap rekap.
func TestTiket38PenyebaranTanpaIndukDanBerindukGandaDitolak(t *testing.T) {
	db, skema, ctx := siapkan(t)
	cabangPenyebaran(t, ctx, db, skema)
	wajibTerima(t, ctx, db, skema, "kelompok", fmt.Sprintf(insKelompok, 9000001, "ZZK1"))
	wajibTerima(t, ctx, db, skema, "detail proporsional",
		fmt.Sprintf(insDetailProp, 9000001, 9000001, 9000001))

	wajibTolak(t, ctx, db, skema, "CK_PENYEBARAN_INDUK",
		fmt.Sprintf(insPenyebar, 9000001, "NULL", "NULL", 9000001))
	wajibTolak(t, ctx, db, skema, "CK_PENYEBARAN_INDUK",
		fmt.Sprintf(insPenyebar, 9000002, "9000001", "9000001", 9000001))
}

// Uji POSITIF tiket 38 — KEDUA pelekatan sah, masing-masing sendirian.
// `CHECK` yang terlalu ketat (mis. hanya mengizinkan `ID_BAGIAN`) lulus kedua
// uji negatif di atas dan menolak separuh data yang sah.
func TestTiket38KeduaPelekatanPenyebaranDiterima(t *testing.T) {
	db, skema, ctx := siapkan(t)
	cabangPenyebaran(t, ctx, db, skema)
	wajibTerima(t, ctx, db, skema, "kelompok", fmt.Sprintf(insKelompok, 9000001, "ZZK1"))
	wajibTerima(t, ctx, db, skema, "detail proporsional",
		fmt.Sprintf(insDetailProp, 9000001, 9000001, 9000001))

	wajibTerima(t, ctx, db, skema, "penyebaran pada BAGIAN",
		fmt.Sprintf(insPenyebar, 9000001, "9000001", "NULL", 9000001))
	wajibTerima(t, ctx, db, skema, "penyebaran pada DETAIL_PROPORSIONAL",
		fmt.Sprintf(insPenyebar, 9000002, "NULL", "9000001", 9000001))
}

// Tiket 38 — INV-16: jenis reasuransi unik di dalam satu induk.
func TestTiket38PenyebaranJenisGandaDitolak(t *testing.T) {
	db, skema, ctx := siapkan(t)
	cabangPenyebaran(t, ctx, db, skema)

	wajibTerima(t, ctx, db, skema, "penyebaran pertama",
		fmt.Sprintf(insPenyebar, 9000001, "9000001", "NULL", 9000001))
	wajibTolak(t, ctx, db, skema, "UQ_PENYEBARAN",
		fmt.Sprintf(insPenyebar, 9000002, "9000001", "NULL", 9000001))
}

// Uji POSITIF — dua JENIS berbeda pada bagian yang sama diterima. Itu seluruh
// pokok penyebaran: satu bagian disebar ke beberapa jenis reasuransi.
func TestTiket38DuaJenisPadaSatuBagianDiterima(t *testing.T) {
	db, skema, ctx := siapkan(t)
	cabangPenyebaran(t, ctx, db, skema)

	wajibTerima(t, ctx, db, skema, "jenis 1",
		fmt.Sprintf(insPenyebar, 9000001, "9000001", "NULL", 9000001))
	wajibTerima(t, ctx, db, skema, "jenis 2, bagian SAMA",
		fmt.Sprintf(insPenyebar, 9000002, "9000001", "NULL", 9000002))
}

// Tiket 38 — INV-65 pada rincian, dan rantai `ikut hapus` sampai ke nilai.
//
// ⭐ Uji kaskade dua tingkat: menghapus `PENYEBARAN` harus membawa serta
// rincian DAN nilai di bawahnya. Rantai yang putus di tengah meninggalkan
// nilai yatim yang tidak terlihat sampai seseorang menjumlahkannya.
func TestTiket38RantaiIkutHapusDuaTingkat(t *testing.T) {
	db, skema, ctx := siapkan(t)
	cabangPenyebaran(t, ctx, db, skema)

	wajibTerima(t, ctx, db, skema, "penyebaran",
		fmt.Sprintf(insPenyebar, 9000001, "9000001", "NULL", 9000001))
	wajibTerima(t, ctx, db, skema, "rincian",
		fmt.Sprintf(insRincian, 9000001, 9000001, 9000001))
	wajibTolak(t, ctx, db, skema, "UQ_RINCIAN_PENYEBARAN",
		fmt.Sprintf(insRincian, 9000002, 9000001, 9000001))
	wajibTerima(t, ctx, db, skema, "nilai penyebaran",
		fmt.Sprintf(insNilaiPenyebaran, 9000001, 9000001, "IDR"))

	wajibTerima(t, ctx, db, skema, "hapus penyebaran",
		`DELETE FROM {skema}.PENYEBARAN WHERE ID_PENYEBARAN = 9000001`)
	for _, q := range []string{
		`SELECT COUNT(*) FROM {skema}.RINCIAN_PENYEBARAN WHERE ID_RINCIAN_PENYEBARAN = 9000001`,
		`SELECT COUNT(*) FROM {skema}.NILAI_PENYEBARAN WHERE ID_NILAI_PENYEBARAN = 9000001`,
	} {
		if n := cacah(t, ctx, db, skema, q); n != 0 {
			t.Errorf("rantai ikut hapus putus: cacah %d, mau 0 — %s", n, q)
		}
	}
}

// Uji POSITIF — jenis reasuransi yang sudah dipakai TIDAK dapat dihapus dari
// tabel acuannya (`ERD.md` §2.7, `tolak`).
func TestTiket38JenisReasuransiTerpakaiTidakDapatDihapus(t *testing.T) {
	db, skema, ctx := siapkan(t)
	cabangPenyebaran(t, ctx, db, skema)
	wajibTerima(t, ctx, db, skema, "penyebaran",
		fmt.Sprintf(insPenyebar, 9000001, "9000001", "NULL", 9000001))

	wajibTolak(t, ctx, db, skema, "FK_PENYEBARAN_3",
		`DELETE FROM {skema}.JENIS_REASURANSI WHERE ID_JENIS_REASURANSI = 9000001`)
}

// ------------------------------------- tabel WARISAN `TREATY_IN` (baca saja)

// ⭐ SATU uji yang membuktikan jalurnya membaca tabel SUNGGUHAN.
//
// Yang dibuktikan tiga hal, dan ketiganya tidak dapat dibuktikan
// `gudangTiruan`: tabelnya ada, ia berisi, dan kesembilan kolom layar
// terbaca dengan nama yang ditulis kode.
//
// ⛔ BACA SAJA, dan ia tetap memakai `siapkan(t)` yang mengakhiri dengan
// `Rollback` — meski uji ini nol menulis. Satu jalan keluar untuk semua
// lebih murah daripada dua, dan transaksi yang tidak pernah di-commit tidak
// dapat berubah menjadi yang menulis karena kelalaian berikutnya.
func TestWarisanTreatyInTerbacaDariOracle(t *testing.T) {
	db, skema, ctx := siapkan(t)

	n := cacah(t, ctx, db, skema, `SELECT COUNT(*) FROM {skema}.TREATY_IN`)
	if n <= 0 {
		t.Fatalf("TREATY_IN cacah %d; layar daftar membacanya, jadi nol baris berarti "+
			"layar menampilkan kosong di atas data yang ada", n)
	}
	t.Logf("TREATY_IN: %d baris", n)

	// Kesembilan kolom layar dibaca dengan nama yang SAMA dengan yang
	// `repository/warisan_daftar.go` tulis. Salah satu nama yang meleset
	// akan menjadi ORA-00904 di sini, bukan layar kosong di produksi.
	q := `SELECT ID, TREATYCONTRACTNAME, PROPORTIONTYPE, LEADINGREINSSOURCE,
	      CEDING, COMMENCEMENT, TERMINATION, POSITIONUSERNAME, STATUSAKSEPTASI
	      FROM {skema}.TREATY_IN
	      ORDER BY TO_NUMBER(ID DEFAULT 0 ON CONVERSION ERROR) DESC, ID DESC
	      OFFSET 0 ROWS FETCH NEXT 3 ROWS ONLY`
	baris, err := db.QueryContext(ctx, strings.ReplaceAll(q, "{skema}", skema))
	if err != nil {
		t.Fatalf("membaca TREATY_IN: %v", err)
	}
	defer func() { _ = baris.Close() }()

	terbaca := 0
	var idPertama string
	for baris.Next() {
		var id, nk, pt, lrs, cd, cm, tm, pu, sa sql.NullString
		if err := baris.Scan(&id, &nk, &pt, &lrs, &cd, &cm, &tm, &pu, &sa); err != nil {
			t.Fatalf("memindai baris: %v", err)
		}
		if terbaca == 0 {
			idPertama = id.String
		}
		terbaca++
	}
	if err := baris.Err(); err != nil {
		t.Fatal(err)
	}
	if terbaca != 3 {
		t.Errorf("terbaca %d baris, mau 3", terbaca)
	}
	// Urut MENURUN: baris pertama harus pengenal TERBESAR. Urutan menaik
	// lulus setiap uji yang hanya menghitung barisnya.
	var terbesar string
	if err := db.QueryRowContext(ctx,
		strings.ReplaceAll(`SELECT MAX(TO_NUMBER(ID DEFAULT 0 ON CONVERSION ERROR)) FROM {skema}.TREATY_IN`, "{skema}", skema),
	).Scan(&terbesar); err != nil {
		t.Fatal(err)
	}
	if idPertama != terbesar {
		t.Errorf("baris teratas %q, mau %q — urutannya tidak menurun", idPertama, terbesar)
	}
	t.Logf("baris teratas: %s", idPertama)
}

// ⭐ Satu kontrak warisan SUNGGUHAN terbaca utuh — kolom DAN dokumennya.
//
// Yang dibuktikan di sini dan tidak dapat dibuktikan `gudangTiruan`: `JOIN`
// antara `TREATY_IN` dan `M_TREATY_IN` cocok, `CLOB`-nya terbaca sebagai
// teks oleh driver, dan JSON-nya terurai menjadi kelima kunci yang form
// pakai.
//
// ⛔ BACA SAJA, dan tetap memakai `siapkan(t)` yang mengakhiri `Rollback`.
func TestWarisanSatuKontrakTerbacaUtuh(t *testing.T) {
	db, skema, ctx := siapkan(t)

	// Kontrak yang PASTI punya kelima kunci — dipilih oleh datanya sendiri,
	// bukan pengenal yang ditanam di sini. Pengenal yang ditanam akan
	// menjadi uji yang merah pada hari baris itu berubah.
	var id string
	q := `SELECT ID FROM (
	        SELECT m.ID FROM {skema}.M_TREATY_IN m
	        WHERE INSTR(m.JSONDATA,'"ContractRefNo":') > 0
	          AND INSTR(m.JSONDATA,'"TreatyLeader":') > 0
	          AND INSTR(m.JSONDATA,'"Bordeaux":') > 0
	        ORDER BY m.ID)
	      WHERE ROWNUM = 1`
	if err := db.QueryRowContext(ctx, strings.ReplaceAll(q, "{skema}", skema)).Scan(&id); err != nil {
		t.Fatalf("mencari kontrak berkunci lengkap: %v", err)
	}
	t.Logf("kontrak uji: %s", id)

	var nama, jenis, mulai, dok sql.NullString
	qb := `SELECT t.TREATYCONTRACTNAME, t.PROPORTIONTYPE, t.COMMENCEMENT, m.JSONDATA
	       FROM {skema}.TREATY_IN t LEFT JOIN {skema}.M_TREATY_IN m ON m.ID = t.ID
	       WHERE t.ID = :1`
	if err := db.QueryRowContext(ctx, strings.ReplaceAll(qb, "{skema}", skema), id).
		Scan(&nama, &jenis, &mulai, &dok); err != nil {
		t.Fatalf("membaca kontrak %s: %v", id, err)
	}
	if !dok.Valid || len(dok.String) < 1000 {
		t.Fatalf("JSONDATA kontrak %s terbaca %d bita; dokumen terpendek yang terukur 1.764",
			id, len(dok.String))
	}
	if nama.String == "" || jenis.String == "" {
		t.Errorf("kolom kosong: nama=%q jenis=%q", nama.String, jenis.String)
	}
	// Kelima kunci yang form pakai ADA di dokumen ini.
	for _, kunci := range []string{
		`"Bordeaux":`, `"AccountingMode":`, `"BordereauxNote":`,
		`"ContractRefNo":`, `"TreatyLeader":`,
	} {
		if !strings.Contains(dok.String, kunci) {
			t.Errorf("kunci %s tidak ada di dokumen %s", kunci, id)
		}
	}
	t.Logf("dokumen %s: %d bita, kelima kunci ada", id, len(dok.String))
}
