package models

// Untuk apa berkas ini: ATURAN MURNI LAIN di sekitar layar realisasi - pra-proses
// flow action (Data Transform), penanda penempatan keluar,
// status riwayat, dan bentuk nomor polis. Seluruhnya port satu rule; nama rule
// dan langkahnya ditulis di atas tiap fungsi (INVENTARIS-XML.md bab 6-7).

import (
	"strings"
	"time"
)

// ---------------------------------------------------------------- pra-proses

// SalinQuotation meniru `pyWorkPage.PolicyTreatyIn.QuotationData = pyWorkPage.Quotation`
// (salinan seluruh halaman, termasuk properti yang kosong di sumber).
func SalinQuotation(h *Halaman) {
	h.pastikan()
	for k := range h.Nilai {
		if strings.HasPrefix(k, HalamanPolis+".QuotationData.") {
			delete(h.Nilai, k)
		}
	}
	for k, v := range h.Nilai {
		if strings.HasPrefix(k, HalamanQuotation+".") {
			h.Nilai[HalamanPolis+".QuotationData."+strings.TrimPrefix(k, HalamanQuotation+".")] = v
		}
	}
}

// GeserTanggalProduksi = `InputPolicyTreatyInPre_Act` langkah 9 ("Set
// Production Date kalau diatas tanggal 25"):
//
//	syarat  `@substring(pyWorkPage.PolicyTreatyIn.StatementDate,6,2)>25`
//	        benar -> jalankan; salah -> lewati (hari = batas TIDAK digeser)
//	ProductionDate = @addCalendar(StatementDate,'0','1','0','0','0','0','0')
//	ProductionDate = @substring(ProductionDate,0,6) + "01" + @substring(ProductionDate,8)
//
// Hasilnya tanggal 1 bulan berikut; `@substring(..,8)` membawa BAGIAN WAKTU
// StatementDate (`THHmmss.SSS GMT`) apa adanya - jamnya dipertahankan.
// ⛔ RALAT putaran 2: port pertama menulis pukul 00:00:00 - tidak ada di rule.
//
// ⛔ PENYIMPANGAN SADAR. Rule menanam `>25`; `GeneratePolicyNoTreaty_Act`
// langkah 5.3 membaca hari yang sama dari `POOLDATA.TANGGAL_CLOSING`
// (`GETTanggalClosing_SQL` -> `Local.TglProd`). Dua aturan hidup berdampingan
// - pola yang SAMA dengan `InsertJsonPolisLife_Act` b1170 lawan
// `PROC_GENERATE_SEQUENCE_NUMBER` di PremiumList Life, tempat `[keputusan work
// owner]` menetapkan "ikuti yang dari DB" (penjaga
// `TestNolAmbangTutupBukuTertanam`). Preseden itu diterapkan di sini:
// `hariClosing` dibaca pemanggil dari tabel (`penomor.HariClosing`), dan
// pembandingnya tetap `>` seperti rule.
//
// ⚠️ Catatan: `@substring(StatementDate,6,2)` membaca hari dari teks internal
// DateTime Pega (GMT); di sini hari dibaca di zona waktu `statement` itu sendiri
// (jam aplikasi). Keduanya sama kecuali pukul 00:00-06:59 WIB bila jam
// aplikasi berzona WIB - dicatat, tidak ditiru (bergantung zona JVM Pega, yang
// tidak ada di korpus).
func GeserTanggalProduksi(statement time.Time, hariClosing int) time.Time {
	if statement.Day() > hariClosing {
		return time.Date(statement.Year(), statement.Month()+1, 1,
			statement.Hour(), statement.Minute(), statement.Second(), statement.Nanosecond(), statement.Location())
	}
	return statement
}

// ---------------------------------------------------------------- penempatan keluar

// ---------------------------------------------------------------- riwayat

// StatusRiwayat = `Activity/InsertHistoryAkseptasiPega` langkah 2-3 untuk
// kasus treaty: IsApproved "1" -> ACCEPT, "0" -> REJECT, selain itu kosong
// (keenam penanda EmailType* milik modul lain, nol di halaman treaty).
func StatusRiwayat(isApproved string) string {
	switch isApproved {
	case "1":
		return "ACCEPT"
	case "0":
		return "REJECT"
	}
	return ""
}

// ---------------------------------------------------------------- nomor polis

// TanggalProduksiNomor = `GeneratePolicyNoTreaty_Act` langkah 5.1-5.3.
//
//	5.1 ProductionDate = @CurrentDateTime()        (kotak When tak dicentang)
//	5.2 StatementDate di depan hari ini -> ProductionDate = StatementDate
//	5.3 hari ProductionDate (Asia/Jakarta) > tanggal closing ->
//	    tanggal 1 bulan berikut, pukul 05:00 GMT
func TanggalProduksiNomor(sekarang, statement time.Time, hariClosing int) time.Time {
	prod := sekarang
	// 5.2 `@toInt(@DateTime.DateTimeDifference(StatementDate,@CurrentDateTime(),D))<0`
	// - selisih (sekarang - statement) dalam HARI, dipotong ke nol oleh @toInt:
	// benar hanya bila StatementDate sedikitnya satu hari penuh di depan.
	if !statement.IsZero() && sekarang.Sub(statement) <= -24*time.Hour {
		prod = statement
	}
	lokal := prod.In(zonaJakarta())
	if lokal.Day() > hariClosing {
		bulan := time.Date(lokal.Year(), lokal.Month(), 1, 5, 0, 0, 0, time.UTC)
		return bulan.AddDate(0, 1, 0)
	}
	return prod
}
