package models

// Untuk apa berkas ini: SALINAN SuggestList DOKUMEN LAMA ENDORSEMEN ke riwayat produksi
// `POOLDATA.HISTORYAKSEPTASIPRODUCTION` - tiket EDM 10. Asal: adaptasi `modul/nbtreatyin/backend/models/usulanlama.go`
// (keputusan work owner F3 NB 04-10-2026: "SuggestList dokumen lama disalin ke HISTORYAKSEPTASIPRODUCTION, dengan
// penjaga dobel menurut IDPEGA").
//
// ⚠️ EDM: korpus EDM Treaty In TIDAK memuat `SaveViewSuggest` / `InsertViewSuggest_SQL` - SuggestList EDM di Pega
// hanya tinggal di work object (dan di DATA_JSON json_polis). Jalur biasa EDM sudah memakai ketetapan NB K4
// (`usulan.go`: HISTORYAKSEPTASIPRODUCTION ber-IDPEGA = `KunciInstans(id)`, dibaca balik `repository.bacaUsulan`),
// sehingga catatan lama disalin ke tempat dan kunci yang SAMA supaya kasus lama yang dibuka di sistem baru
// memperlihatkan catatannya. Pemetaan = `petaUsulan` jalur biasa (`SaveViewSuggest` langkah 2.1.2), dengan sumber
// halaman kerja yang di dokumen lama tidak ada:
//
//	CARI1 TYPE_POLIS   awalan pyID tanpa tanda hubung ("EDMT")
//	CARI7 TYPE         QuotationData.BusinessFac dokumen; BUSINESS_CODE QuotationData.BusinessCode
//	AKSES_LOGIN        selalu NULL - baris SuggestList dokumen tanpa anggota operator (data guide `$.SuggestList`:
//	                   Date, Suggest, IsApproved, pxObjClass, OperatorName); tidak dikarang
//	PIC (CARI4)        .OperatorName baris - kosong = NULL, dihitung
//	TGL_INP (CARI5)    .Date baris dibaca `BacaTanggalLama` (format lain = galat dokumen, K15)

import "strings"

// pmMedanUsulanLama - anggota baris SuggestList dokumen lama yang punya tujuan di riwayat produksi (atau penanda
// langkah 2.1). Anggota lain - selain pxObjClass (internal Pega) - belum diputuskan.
var pmMedanUsulanLama = map[string]string{
	"Date":         "TGL_INP (CARI5)",
	"IsApproved":   "APPROVAL (CARI9)",
	"OperatorName": "PIC (CARI4)",
	"Suggest":      "KETERANGAN (CARI10)",
	"IsSave":       "penanda langkah 2.1 (.IsSave == \"\")",
}

// pmPolaUsulanLama - awalan pola anggota baris SuggestList.
const pmPolaUsulanLama = DaftarUsulan + "()."

// pmAnggotaUsulanLama - nama anggota bila pola medan adalah anggota baris SuggestList yang disalin.
func pmAnggotaUsulanLama(pola string) (string, bool) {
	nama, ok := strings.CutPrefix(pola, pmPolaUsulanLama)
	if !ok {
		return "", false
	}
	_, dikenal := pmMedanUsulanLama[nama]
	return nama, dikenal
}

// UsulanDokumenLamaEDM memetakan baris SuggestList dokumen lama yang `IsSave`-nya kosong (langkah 2.1) menjadi
// baris riwayat produksi. `id` = pyID kasus (`PyIDKasus`); `.Date` sudah dibaca pemecah. Halaman tidak
// diubah.
func UsulanDokumenLamaEDM(id string, h *Halaman) []UsulanProduksi {
	q := HalamanPolis + ".QuotationData."
	k := kasusUsulan{
		TypePolis:    strings.ReplaceAll(strings.TrimRight(id, "0123456789"), "-", ""),
		BusinessFac:  h.Ambil(q + "BusinessFac"),
		BusinessCode: h.Ambil(q + "BusinessCode"),
	}
	var out []UsulanProduksi
	for _, b := range h.AmbilDaftar(DaftarUsulan) {
		if b["IsSave"] != "" {
			continue
		}
		u := petaUsulan(b, k)
		u.AksesLogin = "" // NULL: baris dokumen lama tanpa anggota operator (kepala berkas)
		out = append(out, u)
	}
	return out
}
