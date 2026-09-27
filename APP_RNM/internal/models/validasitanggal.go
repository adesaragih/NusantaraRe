package models

// Dua validasi tanggal berambang produk - kelompok Detail & Tutup.
//
// Untuk apa berkas ini: Pega memeriksa dua jarak hari terhadap ambang yang
// disimpan di produk, dan menuliskan hasilnya sebagai PENANDA - kosong bila
// sah, tanggal yang melanggar bila tidak.
//
//	`Activity/ValidasiClaimReceived_Act.xml`  DATE_OF_LOSS -> CLAIM_RECEIVED_DATE
//	  ambang `MAXEXPIREDCLAIM`, penanda `.MAXCLAIM_RECEIVED`
//	`Activity/ValidasiSTNC_Act.xml`           EFFECTIVE_DATE -> RECEIVED_DATE
//	  ambang `MAXDATARECEIVE`, penanda `.STNC` (kolom `STNC_CLAIM`, migrasi 002)
//
// Keduanya memakai satu fungsi, sebab ATURANNYA sama: selisih hari
// dibandingkan dengan ambang, hasilnya penanda.
//
// ⚠️ Tetapi keduanya TIDAK sama persis, dan bedanya bukan gaya
// penulisan: ValidasiClaimReceived_Act b561 menyimpan selisihnya ke
// `TempDetail.CARI2` - properti pada halaman `TempDetail`, halaman yang
// SAMA dengan yang dipakai `SaveInsuredClaim_Act` untuk menumpuk peserta
// (`TempDetail.pxResults(<APPEND>)`) - lalu b582 membacanya kembali.
// ValidasiSTNC_Act b598 memakai `Local.DateDif`, local sejati.
//
// Nilai antara yang diparkir di halaman bersama dapat tertimpa di antara
// penulisan dan pembacaannya. Dilaporkan OQ-G3. Di sini tidak ada
// masalah itu: selisihnya peubah lokal Go dan tidak pernah meninggalkan
// fungsi ini.
//
// ⛔ AMBANGNYA BELUM PUNYA SUMBER. Ia dibaca `RDBList/GetProductName.xml:84`:
//
//	SELECT * FROM POOLDATA.PRODUCTINWARD_LIFE
//	 WHERE ID = {pyWorkPage.PolicyDataLife.ProductNameID}
//
// `PolicyDataLife` menunggu modul PremiumList Life (keputusan av), dan
// `PRODUCTINWARD_LIFE` berstatus `[data DBA]` OQ-001 - tabel warisan yang
// tidak boleh dimigrasi. Karena itu aturannya ditulis di sini sebagai fungsi
// murni yang MENERIMA ambang, dan pembacanya menyusul bersama modul itu.
// Menebak ambangnya berarti melarang klaim memakai angka karangan.
//
// Dibaca sesudah: tahap.go.

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"nusantarare/pkg/utils"
)

// Pesan kedua validasi, VERBATIM dari rule-nya.
//
// ⚠️ CACAT RULE WARISAN, dilaporkan OQ-G dan tidak ditiru hurufnya: di Pega
// kedua kalimat ini disetel ke `local.errmsg` (b451 dan b479) lalu TIDAK
// PERNAH DIPASANG - kedua activity itu tidak punya langkah
// `Property-Set-Messages`, tidak seperti `ValidasiDOL_Act` yang punya (b842).
// Menurut XML apa adanya, pemakai tidak pernah melihat sebab penolakannya;
// yang tampak hanya satu medan yang tiba-tiba terisi tanggal.
//
// Yang ditiru MAKSUDnya: kalimatnya ada di rule, dan kalimat itulah yang
// dipakai - supaya orang dapat mencari teks yang sama di kedua sistem.
const (
	// b451 `ValidasiClaimReceived_Act`.
	PesanMaxClaimTidakSah = "Max Claim invalid"
	// b479 `ValidasiSTNC_Act`.
	PesanSTNCTidakSah = "STNC invalid"
)

// bentukPenanda adalah bentuk tanggal penanda, `@FormatDateTime` b582/b627.
//
// ⚠️ `dd/MM/YYYY` di Pega berarti hari-bulan-tahun berdigit tetap. Ia BUKAN
// bentuk yang dipakai basis data (`YYYY-MM-DD`): penanda ini dibaca orang,
// bukan diurai mesin.
const bentukPenanda = "02/01/2006"

// SelisihHari mengembalikan jarak hari dari `dari` ke `ke`.
//
// Padanan `@DateTimeDifference(dari, ke, "D")` b561/b598.
//
// ⛔ Dihitung dari tanggalnya saja, sesudah jamnya dibuang. `@addCalendar`
// dengan seluruh argumen nol (b498, b519, b526, b556) memang tidak menggeser
// apa pun - ia hanya menormalkan nilainya menjadi tanggal. Menghitung selisih
// berjam membuat dua tanggal yang berselisih 25 jam terbaca 1 hari dan yang
// berselisih 23 jam terbaca 0.
//
// ⚠️ Hasilnya boleh NEGATIF, dan itu disengaja - lihat PenandaBatasHari.
func SelisihHari(dari, ke string) (int, error) {
	a, err := utils.ParseTanggal(strings.TrimSpace(dari))
	if err != nil {
		return 0, fmt.Errorf("models: tanggal awal %q: %w", dari, err)
	}
	b, err := utils.ParseTanggal(strings.TrimSpace(ke))
	if err != nil {
		return 0, fmt.Errorf("models: tanggal akhir %q: %w", ke, err)
	}
	hariA := time.Date(a.Year(), a.Month(), a.Day(), 0, 0, 0, 0, time.UTC)
	hariB := time.Date(b.Year(), b.Month(), b.Day(), 0, 0, 0, 0, time.UTC)
	return int(hariB.Sub(hariA).Hours() / 24), nil
}

// PenandaBatasHari menjawab penanda kedua validasi itu.
//
// Kosong berarti SAH. Terisi berarti tidak sah, dan isinya adalah tanggal
// `ke` dalam bentuk `dd/MM/yyyy` - tanggal yang melanggarnya.
//
// ⛔ Kenapa penanda dan bukan bool: itulah yang ditulis Pega ke medannya
// (b582, b627), dan medan itu yang dibaca layar serta disimpan. Menggantinya
// dengan bool membuang keterangan MANA tanggal yang melanggar, dan layar
// tidak lagi dapat menunjukkannya tanpa menghitung ulang.
//
// ⛔ `batas` KOSONG menjawab galat, bukan nol. Ambangnya datang dari produk
// dan selama modul PremiumList Life belum ada ia memang belum diketahui;
// menganggapnya nol akan menandai hampir setiap klaim tidak sah (ADR-U-0027).
func PenandaBatasHari(dari, ke, batas string) (string, error) {
	b := strings.TrimSpace(batas)
	if b == "" {
		// ⛔ Pesannya TIDAK menyebut nama tabel warisan. Pesan galat
		// mendarat di log, dan nama objek warisan tidak pernah masuk log
		// atau artefak. Keterangan lengkapnya ada di komentar kepala
		// berkas ini, tempat yang memang untuk itu.
		return "", fmt.Errorf(
			"models: ambang hari belum diketahui; sumbernya data produk, " +
				"yang menunggu modul PremiumList Life")
	}
	ambang, err := strconv.Atoi(b)
	if err != nil {
		return "", fmt.Errorf("models: ambang hari %q bukan bilangan bulat: %w", batas, err)
	}
	selisih, err := SelisihHari(dari, ke)
	if err != nil {
		return "", err
	}
	// ⚠️ Perbandingannya `<=` apa adanya, TANPA lantai bawah - selisih
	// negatif karena itu selalu lolos. Itu lubang warisan yang ditiru dengan
	// sengaja dan dilaporkan (OQ-G); `TestLubangSelisihNegatif` menjaganya
	// supaya ia tidak berubah diam-diam, dan akan menagih keputusan bila
	// kelak seseorang menambalnya.
	if selisih <= ambang {
		return "", nil
	}
	t, err := utils.ParseTanggal(strings.TrimSpace(ke))
	if err != nil {
		return "", err
	}
	return t.Format(bentukPenanda), nil
}
