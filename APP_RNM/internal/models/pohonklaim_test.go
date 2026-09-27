package models

import (
	"encoding/json"
	"reflect"
	"sort"
	"strings"
	"testing"
)

// TestNamaJSONDokumenDikunci mengunci himpunan kunci JSON dokumen.
//
// ⛔ SISI GO DARI KONTRAK DUA SISI. Pasangannya `Dokumen` di
// `frontend/src/services/api.ts` dan `PanelDokumenPeserta.test.ts`.
//
// Tanpa tag JSON, Go mengirim `ID`/`NamaFile` berhuruf besar dan React
// membaca `undefined` - daftar dokumen yang tampil KOSONG padahal barisnya
// ada, tanpa satu pun galat. Itu bentuk cacat yang sudah terjadi empat kali
// di modul ini, dan tiap kali tiap sisi hijau sendirian.
//
// ⛔ Dikunci sebagai HIMPUNAN UTUH: daftar nama terlarang selalu kalah dari
// nama yang belum terpikirkan. Dan tidak satu pun kunci di bawah dapat
// memuat nama orang - `namaFile` adalah nama BERKAS unggahan.
func TestNamaJSONDokumenDikunci(t *testing.T) {
	b, err := json.Marshal(Dokumen{ID: 1, PesertaID: "UJI-PES-1", NamaFile: "UJI-berkas.pdf"})
	if err != nil {
		t.Fatal(err)
	}
	var peta map[string]json.RawMessage
	if err := json.Unmarshal(b, &peta); err != nil {
		t.Fatal(err)
	}
	kunci := make([]string, 0, len(peta))
	for k := range peta {
		kunci = append(kunci, k)
	}
	sort.Strings(kunci)
	mau := []string{"id", "kategori1", "kategori2", "mime", "namaFile",
		"paymentDate", "pesertaId", "tStorageId", "tanggal"}
	if !reflect.DeepEqual(kunci, mau) {
		t.Errorf("kunci JSON dokumen = %v, mau %v.\n"+
			"Ubah JUGA tipe Dokumen di api.ts - keduanya satu kontrak.", kunci, mau)
	}
	// ⛔ Tanggal NULL menyeberang sebagai null, bukan sebagai tahun 1 Masehi.
	if !strings.Contains(string(b), `"tanggal":null`) {
		t.Errorf("tanggal kosong tidak dikirim sebagai null: %s", b)
	}
}
