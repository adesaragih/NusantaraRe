package models

// Untuk apa berkas ini: SALINAN SuggestList DOKUMEN LAMA ke riwayat produksi
// `POOLDATA.HISTORYAKSEPTASIPRODUCTION` - keputusan work owner F3
// (04-10-2026, PROMPT-NB-TREATY-IN-PUTARAN-3.md bab 2): "`SuggestList`
// dokumen lama disalin ke `HISTORYAKSEPTASIPRODUCTION` (tempat grilling untuk
// Suggest, sejalan K4), dengan penjaga dobel menurut `IDPEGA`".
//
// Pemetaannya SAMA dengan jalur biasa (`UsulanBelumTersimpan`):
// `Activity/SaveViewSuggest` langkah 2 "UNTUK TREATY" - kalang
// `pyWorkPage.PolicyTreatyIn.SuggestList`, per baris bila `.IsSave == ""`
// (2.1), langkah 2.1.2 CARI1..CARI10 -> `RDBList/InsertViewSuggest_SQL`.
// Bedanya hanya sumber nilai halaman kerja yang di dokumen lama tidak ada:
//
//	CARI1 @replaceAll(pyWorkIDPrefix,"-","") - awalan pyID kasus lama (IDPEGA)
//	CARI7 Quotation.BusinessFac   - dokumen: QuotationData.BusinessFac, salinan
//	BUSINESS_CODE Quotation.BusinessCode       halaman Quotation (Page-Copy
//	                                           GeneratePolicyNoTreaty_Act 10)
//	AKSES_LOGIN OperatorID.pyUserIdentifier - operator yang menjalankan
//	              SaveViewSuggest (sesi, bukan isi baris). Catatan lama SELALU
//	              NULL: baris SuggestList tidak punya anggota operator -
//	              `docs/dataguide-json-polis.json` hanya `$.SuggestList.Date`,
//	              `.Suggest`, `.IsApproved`, `.pxObjClass`, `.OperatorName`, dan
//	              penulis barisnya `DataTransform/AddToListCommentsPolicyTreatyIn_DT`
//	              hanya menulis `.Suggest`, `.IsApproved`, `.Date`, `.OperatorName`
//	              (tinjauan spec P3 (c)3; RALAT tiket 22). Tidak dikarang.
//	PIC (CARI4) .OperatorName baris - kosong = NULL, dihitung
//	TGL_INP (CARI5) .Date baris - cap waktu Pega dibaca `BacaTanggalLama`
//	              (jam dinding Asia/Jakarta, 24 jam - F2 butir 2); format lain
//	              = galat dokumen (K15), tidak ditebak
//	DIV, B2B, PERCENT_RNM - tanpa sumber (OperatorID.pyOrgDivision, OfferFacIn)
//	              = NULL, sama dengan jalur biasa
//
// Penjaga dobel menurut IDPEGA ada di repository (`SalinUsulanLama`):
// IDPEGA yang sudah punya baris riwayat produksi tidak disalin lagi.

import "strings"

// medanUsulanLama - anggota baris SuggestList dokumen lama yang punya tujuan
// di riwayat produksi (atau penanda langkah 2.1). Anggota lain - selain
// pxObjClass (internal Pega) - "belum diputuskan".
var medanUsulanLama = map[string]string{
	"Date":         "TGL_INP (CARI5)",
	"IsApproved":   "APPROVAL (CARI9)",
	"OperatorName": "PIC (CARI4)",
	"Suggest":      "KETERANGAN (CARI10)",
	"IsSave":       "penanda langkah 2.1 (.IsSave == \"\")",
}

// polaUsulanLama - awalan pola anggota baris SuggestList.
const polaUsulanLama = DaftarUsulan + "()."

// anggotaUsulanLama - nama anggota bila pola medan adalah anggota baris
// SuggestList yang disalin ke riwayat produksi.
func anggotaUsulanLama(pola string) (string, bool) {
	nama, ok := strings.CutPrefix(pola, polaUsulanLama)
	if !ok {
		return "", false
	}
	_, dikenal := medanUsulanLama[nama]
	return nama, dikenal
}

// kasusUsulan - nilai halaman kerja yang dibaca langkah 2.1.2 (sama untuk
// setiap baris SuggestList kasus itu) - jalur biasa (`UsulanBelumTersimpan`)
// dan dokumen lama (`UsulanDokumenLama`).
type kasusUsulan struct {
	// TypePolis - CARI1 `@replaceAll(pyWorkIDPrefix,"-","")`.
	TypePolis string
	// BusinessFac - CARI7 `Quotation.BusinessFac` (dokumen: QuotationData).
	BusinessFac string
	// BusinessCode - `Quotation.BusinessCode` (dokumen: QuotationData).
	BusinessCode string
}

// petaUsulan = `SaveViewSuggest` langkah 2.1.2 untuk satu baris SuggestList.
// Tanggal saja (dokumen lama, dibaca `BacaTanggalLama`) = tengah malam, sama
// dengan konversi kolom tanggal repository. AKSES_LOGIN = `.OperatorID` baris
// yang ditulis `TambahCatatan` jalur biasa (lihat `UsulanBelumTersimpan`).
func petaUsulan(b Baris, k kasusUsulan) UsulanProduksi {
	tgl := b["Date"]
	if len(tgl) == len("2006-01-02") {
		tgl += " 00:00:00"
	}
	return UsulanProduksi{
		TypePolis:    k.TypePolis,
		Posisi:       PosisiUsulanProduksi,
		PIC:          b["OperatorName"],
		TglInp:       tgl,
		Type:         k.BusinessFac,
		Putaran:      PutaranUsulanProduksi,
		Approval:     approvalUsulan(b["IsApproved"]),
		Keterangan:   potongKarakter(b["Suggest"], PanjangKeterangan),
		AksesLogin:   b["OperatorID"],
		BusinessCode: k.BusinessCode,
	}
}

// UsulanDokumenLama memetakan baris SuggestList dokumen lama yang `IsSave`-nya
// kosong (langkah 2.1) menjadi baris riwayat produksi. `id` = pyID kasus lama
// (`IDKasusDariIDPega`); `.Date` baris sudah dibaca `BacaTanggalLama` oleh
// pemecah. Halaman tidak diubah: penanda IsSave milik jalur biasa.
func UsulanDokumenLama(id string, h *Halaman) []UsulanProduksi {
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
		u.AksesLogin = "" // NULL: baris dokumen lama tanpa anggota operator (lihat kepala berkas)
		out = append(out, u)
	}
	return out
}
