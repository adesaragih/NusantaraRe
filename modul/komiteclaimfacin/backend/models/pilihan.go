package models

// Untuk apa berkas ini: DAFTAR PILIHAN dan label tampilan. Korpus `Komite Claim FacIn` tidak mengekspor aturan properti
// (prompt values `associated`); `[inferensi]` daftar kelas yang sama dari ekspor korpus `Komite Claim Prop` (pola Komite
// Claim Non Prop) dan teks activity korpus ini. Nilai tersimpan tetap KODE; label hanya untuk tampilan.
//
//	AcceptStatus      .AcceptStatus (Work-Komite)         1 Approve, 2 Reject - SetDataForInformation_Act S3 / S5
//	KomiteAproval     .KomiteAproval (Data-Comitee)       1 Approved, 2 Reject, 0 Waiting - ekspor Komite Claim Prop
//	AcceptanceStatus  .AcceptanceStatus (Data-Adjustment) 0 Transfer to committee, 1 Approve, 2 Reject - idem
//	PaymentType       .PaymentType 1 Final, 2 Interim, 3 Salvage, 4 Adjuster Fee, 5 Adjustment - SetDataForInformation_Act
//	                  S13-S17 (6 / 7 tanpa label: kode apa adanya, OQ-CFI-18)
//	Payable           .Payable 1 Ceding Co Name, 2 Broker Name, 3 Others - Claim Fac In tahap 1

// Pilihan - satu opsi dropdown.
type Pilihan struct {
	Nilai string `json:"nilai"`
	Label string `json:"label"`
}

// LabelTerima - pilihan `.AcceptStatus` "Are you sure to accept this document?".
var LabelTerima = []Pilihan{{Nilai: KeputusanSetuju, Label: "Approve"}, {Nilai: KeputusanTolak, Label: "Reject"}}

// LabelKeputusanAnggota - kolom "Status" grid "List of Committee" (`.KomiteAproval`).
var LabelKeputusanAnggota = map[string]string{"0": "Waiting", "1": "Approved", "2": "Reject"}

// LabelStatusAdjustment - `.AcceptanceStatus` baris adjustment (kolom "Status" History Adjustment).
var LabelStatusAdjustment = map[string]string{"0": "Transfer to committee", "1": "Approve", "2": "Reject"}

// LabelJenisBayar - `.PaymentType` (kolom "Payment Type" History Adjustment).
var LabelJenisBayar = map[string]string{"1": "Final", "2": "Interim", "3": "Salvage", "4": "Adjuster Fee",
	"5": "Adjustment"}

// LabelPayable - `.Payable` "Payable To" (DetailAdjustmentFac LS18).
var LabelPayable = map[string]string{"1": "Ceding Co Name", "2": "Broker Name", "3": "Others"}

// labelKode - label kode; kode tanpa label tampil apa adanya.
func labelKode(peta map[string]string, v string) string {
	if l, ok := peta[v]; ok {
		return l
	}
	return v
}

// LabelStatusBaris - label `.AcceptanceStatus` (kode lain apa adanya).
func LabelStatusBaris(kode string) string { return labelKode(LabelStatusAdjustment, kode) }
