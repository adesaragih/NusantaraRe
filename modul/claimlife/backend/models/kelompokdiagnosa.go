package models

// Daftar pilihan GROUP DIAGNOSE - grid diagnosa `ClaimLifeDetailGCNM` b5860 (`pxDropdown` b5863). Keputusan work owner
// 12-10-2026: daftar di bawah diberikan WO VERBATIM dan berurutan ("jika tabelnya belum ada buatkan").
//
// ⛔ SATU sumber. Tabel master grup diagnosa TIDAK ada di repo maupun korpus (rule properti `.GROUPDIAGNOSE` kelas
// `Data-DiagnoseLife` tidak diekspor - OQ-L, butir bf), dan membuat tabel = migrasi, sedangkan jatah Claim Life
// (001-029) habis. Sampai rentang baru tersedia daftar ini hidup di sini; layar membacanya lewat
// `GET /api/klaim-life/kelompok-diagnosa`, tidak menyalinnya. DDL usulan: docs/PARITAS-LAYAR-DAN-AKSI.md (12-10-2026).
//
// Nilai yang dipilih disimpan apa adanya di `T_CLAIMLF_DIAGNOSE.GROUP_DIAGNOSE VARCHAR2(255)` (migrasi 018); teks
// terpanjang 99 byte - muat, tidak dipotong.
var kelompokDiagnosa = []string{
	"Certain infectious and parasitic diseases",
	"Neoplasms",
	"Diseases of the blood and blood-forming organs and certain disorders involving the immune mechanism",
	"Endocrine, nutritional and metabolic diseases",
	"Mental and behavioural disorders",
	"Diseases of the nervous system",
	"Diseases of the eye and adnexa",
	"Diseases of the ear and mastoid process",
	"Diseases of the circulatory system",
	"Diseases of the respiratory system",
	"Diseases of the digestive system",
	"Diseases of the skin and subcutaneous tissue",
	"Diseases of the musculoskeletal system and connective tissue",
	"Diseases of the genitourinary system",
	"Pregnancy, childbirth and the puerperium",
	"Certain conditions originating in the perinatal period",
	"Congenital malformations, deformations and chromosomal abnormalities",
	"Symptoms, signs and abnormal clinical and laboratory findings, not elsewhere classified",
	"Injury, poisoning and certain other consequences of external causes",
	"External causes of morbidity and mortality",
	"Factors influencing health status and contact with health services",
	"Codes for special purposes",
	"Unknown cause of death/ Sudden death",
	"Unknown cause of death/ Death on arrival",
	"Multi organ failure",
	"Multi organ disease",
}

// KelompokDiagnosa mengembalikan SALINAN daftar pilihan, urutan WO.
func KelompokDiagnosa() []string {
	return append([]string(nil), kelompokDiagnosa...)
}

// KelompokDiagnosaSah menjawab apakah `nilai` boleh disimpan ke GROUP_DIAGNOSE baris yang kini berisi `tersimpan`:
// kosong, salah satu daftar, ATAU nilai warisan yang sama dengan yang sudah tersimpan (tidak diubah - Save atas baris
// lama tidak gagal hanya karena nilainya di luar daftar).
func KelompokDiagnosaSah(nilai, tersimpan string) bool {
	if nilai == "" || nilai == tersimpan {
		return true
	}
	for _, k := range kelompokDiagnosa {
		if k == nilai {
			return true
		}
	}
	return false
}
