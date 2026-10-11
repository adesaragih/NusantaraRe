package models

// Untuk apa berkas ini: DAFTAR PILIHAN properti (prompt values `pyLocalList`). Korpus `Komite Claim Non Prop` tidak
// mengekspor aturan properti; `[inferensi]` daftar yang sama dengan ekspor korpus `Komite Claim Prop` (work owner
// 08-10-2026, kelas induk yang sama `ASM-FW-GCNMFW-Data-*`). Nilai tersimpan tetap KODE; label hanya untuk tampilan.
//
//	AcceptStatus.xml       ASM-FW-GCNMFW-Work-KomiteTreaty .AcceptStatus      1 Approve, 2 Reject
//	AcceptanceStatus.xml   ASM-FW-GCNMFW-Data-Adjustment   .AcceptanceStatus  0 Transfer to committee, 1 Approve, 2 Reject
//	KomiteAproval.xml      ASM-FW-GCNMFW-Data-Comitee      .KomiteAproval     1 Approved, 2 Reject, 0 Waiting
//	Payable.xml            ASM-FW-GCNMFW-Data-Adjustment   .Payable           1 Ceding Co Name, 2 Broker Name, 3 Others
//	SubjectivityNote.xml   ASM-FW-GCNMFW-Work-KomiteTreaty .SubjectivityNote  1..7 (versi 01-01-24)

// PilihanSubjectivityNote - `.SubjectivityNote` (dropdown, versi 01-01-24), urutan VERBATIM.
var PilihanSubjectivityNote = []Pilihan{
	{Nilai: "1", Label: "Treaty Leader Approval"}, {Nilai: "2", Label: "Top 5 Settled Claims"},
	{Nilai: "3", Label: "Invoice(s)"}, {Nilai: "4", Label: "Salvage Document"},
	{Nilai: "5", Label: "Premium Payment Receipt"}, {Nilai: "6", Label: "Reinstatement Evidence"},
	{Nilai: "7", Label: "Others"},
}

// LabelKeputusanAnggota - kolom "Status" grid "Committe Accept Status" (`.KomiteAproval`).
var LabelKeputusanAnggota = map[string]string{"0": "Waiting", "1": "Approved", "2": "Reject"}

// LabelStatusAdjustment - `.AcceptanceStatus` baris adjustment (kolom "Status" History Adjustment).
var LabelStatusAdjustment = map[string]string{"0": "Transfer to committee", "1": "Approve", "2": "Reject"}

// LabelPayable - `.Adjustment.Payable` "Payable To".
var LabelPayable = map[string]string{"1": "Ceding Co Name", "2": "Broker Name", "3": "Others"}

// LabelStatusBaris - label `.AcceptanceStatus` (kode lain apa adanya).
func LabelStatusBaris(kode string) string { return labelKode(LabelStatusAdjustment, kode) }

// SubjectivityNoteSah - kode termasuk daftar pilihan `.SubjectivityNote`.
func SubjectivityNoteSah(kode string) bool {
	for _, p := range PilihanSubjectivityNote {
		if p.Nilai == kode {
			return true
		}
	}
	return false
}
