package repository

// Larik TOTAL yang layar pegang HANYA di penampung halaman — dibaca kembali
// dari `T_TREATY_TOTAL` supaya isian yang di-Save tampil lagi tanpa
// menekan Refresh / Update Total.
//
// ⭐ Laporan pemakai 8 Oktober 2026: tab Share Prop sesudah Save — rincian
// terisi, tetapi `Total Share RNM Limit`, `Total Value Spreading OR` dan
// `… R/I` "No items" sampai Refresh ditekan. Sebabnya: Save MENULIS larik
// ini (`LarikGabung` `T_TREATY_TOTAL`, migrasi `448`), tetapi pembaca
// `T_TREATY_TOTAL` hanya mengembalikan keempat total Limits Non-Prop
// (`BacaLimitsAkarPendaratan`). Tab yang memegangnya di penampung lalu
// menyemai larik kosong.
//
// ⛔ `TotalRetentionAmountNP` TIDAK di sini: tab Maximum Retention sudah
// menyemainya dari kontrak (`totalRetensi`). Total Share Non-Prop
// (`ShareNP.Total`) DIHITUNG ULANG saat kontrak dibuka (`SiapkanShareNP` →
// `NPSetTotalShare`), jadi tidak dibaca dari sini.

import "context"

// LarikTotalPenampung - `JENIS` `T_TREATY_TOTAL` yang tabnya membaca dari
// penampung halaman (`useProperti`), ejaan Pega.
var LarikTotalPenampung = []string{
	// Tab Share Proportional (`TabShareProp`).
	"TotalShareRnmProp", "TotalSpreadedRnmProp", "TotalSpreadedRnmRIProp",
	// Tab EGNPI Non-Prop — panel total per mata uang (`TabEgnpi`).
	"TotalEgnpiAmountNP",
	// Tab Installment Non-Prop — `Total` (`TabAngsuran`).
	"TotalInstallmentNP",
}

// BacaTotalPenampung membaca larik `LarikTotalPenampung` satu kontrak.
//
// ⛔ Larik yang NOL baris tidak dimasukkan: kunci yang tidak ada membiarkan
// tab menyemai bawaannya sendiri, persis seperti sebelum Save pertama.
func (g *Gudang) BacaTotalPenampung(ctx context.Context, masterID string) (map[string][]map[string]any, error) {
	dipakai := map[string]bool{}
	for _, n := range LarikTotalPenampung {
		dipakai[n] = true
	}
	baris, err := g.bacaEntri(ctx, "T_TREATY_TOTAL", masterID)
	if err != nil {
		return nil, err
	}
	p, _ := entriPeta("T_TREATY_TOTAL")
	out := map[string][]map[string]any{}
	for _, b := range baris {
		j := b.Nilai[kolomJenis]
		if dipakai[j] {
			out[j] = append(out[j], simpulPeta(b, p))
		}
	}
	return out, nil
}
