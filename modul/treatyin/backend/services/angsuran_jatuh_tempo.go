package services

import "time"

// jadwalJatuhTempo - `DueDate` tiap baris yang `TreatyInSetValueInstallment`
// bentuk (langkah 8), `n` baris.
//
// ⛔ MENYIMPANG DARI EKSPOR, atas keputusan pemakai 9 Oktober 2026. Pega
// menulis `.InstallmentList(<LAST>).DueDate = @CurrentDate("dd/MM/yyyy",
// "Asia/Jakarta")` - SEMUA baris bertanggal hari ini. Pemakai: *"due date nya
// tidak per hari ini, seharusnya dari Commencement"*, dan memilih jadwal
// DIBAGI RATA: baris 1 = Commencement, baris berikutnya maju
// (lama periode / jumlah installment) bulan. Contoh 01/07/2026-30/06/2027,
// 4 installment: 01/07/2026, 01/10/2026, 01/01/2027, 01/04/2027.
//
// Lama periode = bulan penuh dari Commencement sampai SEHARI SESUDAH
// Termination (01/07/2026-30/06/2027 = 12 bulan). Geser ke-k dibulatkan
// setengah ke atas (12 bulan / 5 = 0, 2, 5, 7, 10); hari dipotong ke akhir
// bulan (`tambahBulan`).
//
// Cadangan, dinyatakan:
//   - Commencement kosong / tak terbaca → hari ini, seperti ekspor;
//   - Termination kosong / tak terbaca / sebelum Commencement / periode
//     kurang dari sebulan → semua baris = Commencement.
func jadwalJatuhTempo(commencement, termination string, n int, kini time.Time) []string {
	out := make([]string, n)
	awal, ok := tanggalMasukan(commencement)
	if !ok {
		for i := range out {
			out[i] = kini.Format("20060102")
		}
		return out
	}
	bulan := 0
	if akhir, ok := tanggalMasukan(termination); ok && !akhir.Before(awal) {
		bulan = selisihBulan(awal, akhir.AddDate(0, 0, 1))
	}
	for k := range out {
		geser := 0
		if bulan > 0 && n > 0 {
			geser = (2*k*bulan + n) / (2 * n)
		}
		out[k] = tambahBulan(awal, geser).Format("20060102")
	}
	return out
}
