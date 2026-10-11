package models

// Proteksi "Process Endorsement" (keputusan work owner 08-10-2026 - BUKAN
// korpus; `SetPremi_EDM` Pega menerima Save tanpa perubahan): kasus Perubahan
// Data tanpa peserta `Old` yang ditandai Delete DAN tanpa peserta `New` hasil
// CSV ditolak - tidak ada yang di-endorse. Kasus Batal tidak pernah ditolak
// (seluruh peserta `Old` selalu ditandai Batal).

// PesanSudahDiproses - Process Endorsement adalah langkah FINAL (keputusan work
// owner 08-10-2026, mengganti perilaku Pega R31): sesudahnya peserta tidak dapat
// ditambah (CSV) maupun dihapus (centang). Sama persis di labels.ts.
const PesanSudahDiproses = "This endorsement has been processed. Participants can no longer be added or deleted."

// PesanTanpaPerubahan - pesan layar dan server (sama persis di labels.ts).
const PesanTanpaPerubahan = "No changes to process. Mark at least one participant for Delete or add participants from CSV."

// AdaPerubahan - boleh diproses? `ditandai` = peserta `Old` yang ditandai
// Delete, `baru` = peserta `New` kasus.
func AdaPerubahan(edmType string, ditandai, baru int) bool {
	if edmType != EdmTypePerubahanData {
		return true
	}
	return ditandai > 0 || baru > 0
}

// AdaPilihan - pilihan layar menandai sesuatu (sebelum dicocokkan ke basis
// data): centang satu per satu, atau DELETE ALL.
func (p PilihanHapus) AdaPilihan() bool {
	return p.Semua || len(p.Pilih) > 0
}
