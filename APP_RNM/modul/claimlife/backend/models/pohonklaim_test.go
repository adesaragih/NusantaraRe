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
// `modul/claimlife/frontend/api.ts` dan `PanelDokumenPeserta.test.ts`.
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

// TestPengenalDokumenMenyeberangSebagaiTeks menjaga cacat ketujuh tetap mati.
//
// ⛔ Pengenal dokumen cap waktu 17 angka (`InsertDocument_Act.xml` b648),
// yaitu ≈2,0e16. `Number.MAX_SAFE_INTEGER` di JavaScript 9.007.199.254.740.991
// ≈9,0e15. Tanpa `,string`, `JSON.parse` membulatkannya diam-diam dan tautan
// unduh menunjuk dokumen yang tidak ada - tanpa satu pun galat di kedua sisi.
//
// ⚠️ Uji ini memeriksa BENTUK KAWATNYA, bukan tipe Go-nya. Tipe Go boleh
// tetap int64 - yang harus tidak berubah adalah kutip di sekeliling nilainya.
func TestPengenalDokumenMenyeberangSebagaiTeks(t *testing.T) {
	const capWaktu = 20260927103000123
	b, err := json.Marshal(Dokumen{ID: capWaktu, PesertaID: "UJI-PES-1"})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(b), `"id":"20260927103000123"`) {
		t.Errorf("pengenal dokumen tidak berkutip: %s\n"+
			"Tanpa tag `,string`, JavaScript membulatkannya menjadi "+
			"20260927103000124 dan tautan unduhnya menunjuk dokumen "+
			"yang tidak ada.", b)
	}
	// ⛔ Dan angkanya memang di luar jangkauan aman JavaScript - kalau suatu
	// hari tidak lagi, uji ini yang harus dibaca ulang, bukan dihapus.
	const maxAman = 9007199254740991
	if capWaktu <= maxAman {
		t.Errorf("cap waktu %d ternyata masih aman di JavaScript (batas %d); "+
			"alasan tag `,string` perlu dibaca ulang", capWaktu, maxAman)
	}
}
