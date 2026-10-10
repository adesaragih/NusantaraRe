package models

// Untuk apa berkas ini: SATU FUNGSI PERHITUNGAN TURUNAN klaim Non Prop (pola Claim Prop `models/turunan.go`).
//
// ⚠️ PENYIMPANGAN SADAR (PARITAS): di Pega nilai di bawah ditulis activity saat medan tertentu berubah dan hidup di
// clipboard; di sini dihitung ulang dari basis tersimpan sesudah setiap aksi dan setiap halaman dimuat, tanpa kolom:
//   - TotalInterestInsured / TotalSumInsuredIDR       CountTotalInterest_Act langkah 1-2;
//   - ListTotalEstimation (Summary / Total XOL)       CountLossAllocation_act langkah 19 (`SusunSummaryXOL`);
//   - SpreadingAdjustment(QS) tingkat klaim           DeleteAkseptasi_Act langkah 3 (`SusunSpreadingAkseptasiKlaim`);
//   - AdjustmentList(n).CommentLOD                    InputAkseptasi_PreAct langkah 3.1 (`TandaiAkseptasiLama`);
//   - AdjustmentList(n).CurencyAdjustment             AddAkseptasiCNP_Act langkah 9-10 (`MataUangAkseptasi`);
//   - AdjustmentList(n).FlagCurrency                  `[inferensi]` - tidak punya penulis di korpus (SetCurrency_Act /
//     SetPayableTreatyNP_Act hanya menulis `Local.FlagCurrency` / `FlagCurrency.CARI15`); bank kedua tampil bila
//     SetPayableTreatyNP_Act mengisinya (`FlagCurrency.CARI15 = 1`, ada akseptasi bermata uang lebih dari satu);
//   - AdjustmentList(n).SpreadingAdjustment rekening  SetAccoutNo_Act (`SetAccountNo`).
//
// Medan tanpa penulis di korpus Non Prop dibiarkan kosong (persis XML): `ClaimData.TotalListClaimAmount(IDR)` dan
// `.TotalEstimasi` (kaki grid).

import "github.com/cockroachdb/apd/v3"

// HitungTurunan menghitung seluruh medan turunan halaman klaim.
func HitungTurunan(h *Halaman) error {
	if err := hitungInterest(h); err != nil {
		return err
	}
	if err := SusunSummaryXOL(h); err != nil {
		return err
	}
	if err := SusunSpreadingAkseptasiKlaim(h); err != nil {
		return err
	}
	TandaiAkseptasiLama(h)
	multi := MultiMataUang(h)
	for i, b := range h.AmbilDaftar(DaftarAdjustment) {
		n := i + 1
		h.SetelDaftar(JalurAdj(n, AnakMataUangAdj), MataUangAkseptasi(h, n))
		if multi {
			b["FlagCurrency"] = "1"
		} else {
			delete(b, "FlagCurrency")
		}
		if err := SetAccountNo(h, n); err != nil {
			return err
		}
	}
	return nil
}

// hitungInterest = CountTotalInterest_Act: TotalInterestInsured per `.Currency` (nama mata uang, urut kemunculan
// pertama) = Σ TSIPerObject; TotalSumInsuredIDR = Σ TSIPerObject x KursObjectItem. Tanpa interest keduanya terhapus
// (langkah 1).
func hitungInterest(h *Halaman) error {
	var kal Kalkulator
	d := h.AmbilDaftar(DaftarInterest)
	var tot []Baris
	idr := apd.New(0, 0)
	for _, b := range d {
		var t Baris
		for _, o := range tot {
			if o["Currency"] == b["Currency"] {
				t = o
			}
		}
		if t == nil {
			t = Baris{"Currency": b["Currency"], "Value": "0"}
			tot = append(tot, t)
		}
		t["Value"] = Teks(kal.Tambah(kal.B(t, "Value"), kal.B(b, "TSIPerObject")))
		idr = kal.Tambah(idr, kal.Kali(kal.B(b, "TSIPerObject"), kal.B(b, "KursObjectItem")))
	}
	if err := kal.Galat(); err != nil {
		return err
	}
	h.SetelDaftar(DaftarTotalTSI, tot)
	if len(d) == 0 {
		h.Hapus(CD + "TotalSumInsuredIDR")
		return nil
	}
	h.Setel(CD+"TotalSumInsuredIDR", Teks(idr))
	return nil
}
