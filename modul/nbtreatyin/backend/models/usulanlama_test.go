package models

// Uji murni salinan SuggestList DOKUMEN LAMA ke riwayat produksi (F3,
// keputusan work owner 04-10-2026: "SuggestList dokumen lama disalin ke
// HISTORYAKSEPTASIPRODUCTION"). Pemetaan = `Activity/SaveViewSuggest`
// langkah 2 "UNTUK TREATY" + `RDBList/InsertViewSuggest_SQL`; nilai harapan
// dihitung tangan dari XML. Fixture fiktif UJI-.

import (
	"errors"
	"testing"
)

// dokumenUjiUsulan - dokumen lama berisi tiga catatan SuggestList: dua belum
// pernah tersimpan (IsSave kosong - kasus treaty tidak pernah lolos syarat
// langkah 2 `BusinessFac == "F"`), satu bertanda IsSave "Yes".
const dokumenUjiUsulan = `{
 "pxObjClass": "ASM-FW-GISFW-Data-PolicyTreatyIn",
 "PolicyNo": "UJI-QP.T1.10.2017.00001",
 "QuotationData": {"ProportionalType": "Proportional", "BusinessFac": "T", "BusinessCode": "UJI-B01"},
 "SuggestList": [
  {"pxObjClass": "UJI-kelas", "Date": "20171002T020000.000 GMT", "IsApproved": "1",
   "OperatorName": "UJI-PENGGUNA A", "Suggest": "UJI-catatan satu"},
  {"Date": "20171003T100000.000 GMT", "IsApproved": "0", "OperatorName": "",
   "Suggest": "UJI-catatan dua", "OperatorID": "UJI-AKUN-B"},
  {"Date": "20171004T000000.000 GMT", "IsApproved": "", "OperatorName": "UJI-PENGGUNA C",
   "Suggest": "UJI-sudah tersimpan", "IsSave": "Yes"}
 ]
}`

func TestSuggestListLamaDisalinMenurutSaveViewSuggest(t *testing.T) { // F3; AC 39-44
	h, err := PecahDokumenLama(barisUji(dokumenUjiUsulan))
	if err != nil || len(h.Galat) > 0 {
		t.Fatalf("%v %+v", err, h.Galat)
	}
	// Anggota SuggestList punya tujuan (riwayat produksi): bukan "belum diputuskan".
	if len(h.BelumDiputuskan) != 0 {
		t.Errorf("medan belum diputuskan %+v", h.BelumDiputuskan)
	}
	if len(h.Usulan) != 2 {
		t.Fatalf("hanya baris ber-IsSave kosong yang disalin (langkah 2.1 `.IsSave==\"\"`), dapat %d: %+v", len(h.Usulan), h.Usulan)
	}
	// Langkah 2.1.2 dihitung tangan:
	//   CARI1 @replaceAll(pyWorkIDPrefix,"-","") - awalan pyID kasus lama "UJI-77" -> "UJI"
	//   CARI3 "Policy" · CARI4 .OperatorName · CARI7 Quotation.BusinessFac
	//   (dokumen: QuotationData = salinan Quotation, GeneratePolicyNoTreaty_Act 10)
	//   CARI5 .Date: 2017-10-02 02:00 GMT = 09:00 WIB (24 jam, F2 butir 2)
	//   CARI8 2 · CARI9 "1" -> "Accept", "0" -> "Reject" · CARI10 .Suggest
	//   AKSES_LOGIN: OperatorID baris bila ada, tidak dikarang; BUSINESS_CODE Quotation.BusinessCode
	harap := []UsulanProduksi{
		{TypePolis: "UJI", Posisi: "Policy", PIC: "UJI-PENGGUNA A", TglInp: "2017-10-02 09:00:00", Type: "T",
			Putaran: "2", Approval: "Accept", Keterangan: "UJI-catatan satu", BusinessCode: "UJI-B01"},
		{TypePolis: "UJI", Posisi: "Policy", PIC: "", TglInp: "2017-10-03 17:00:00", Type: "T",
			Putaran: "2", Approval: "Reject", Keterangan: "UJI-catatan dua", AksesLogin: "UJI-AKUN-B", BusinessCode: "UJI-B01"},
	}
	for i, u := range h.Usulan {
		if u != harap[i] {
			t.Errorf("usulan %d = %+v\nharap      %+v", i+1, u, harap[i])
		}
	}
	if h.Diabaikan[AlasanInternalPega] < 1 {
		t.Errorf("pxObjClass baris SuggestList = internal Pega: %v", h.Diabaikan)
	}
}

func TestSuggestListLamaTanggalAmbiguMenggagalkanDokumen(t *testing.T) { // K15, AC 58
	d := `{"pxObjClass": "ASM-FW-GISFW-Data-PolicyTreatyIn", "PolicyNo": "UJI-QP.T1.10.2017.00001",
	 "SuggestList": [{"Date": "05/06/2017", "Suggest": "UJI-x"}]}`
	h, err := PecahDokumenLama(barisUji(d))
	if err != nil {
		t.Fatal(err)
	}
	if len(h.Galat) != 1 || h.Galat[0].Jalur != "PolicyTreatyIn.SuggestList(1).Date" || !errors.Is(h.Galat[0].Err, ErrTanggalAmbigu) {
		t.Errorf("galat %+v", h.Galat)
	}
}
