package services_test

import (
	"testing"

	"nusantarare/modul/treatyin/backend/services"
)

// ⛔ CACAT YANG DILAPORKAN PEMILIK PROSES 7 Oktober 2026:
// kotak `Valid From` / `Valid Until` grid Rate of Exchange KOSONG, padahal
// `TREATYEXCHANGEYEARLY.STARTDATE` berisi `20180101T000000.000 GMT`.
//
// Sebabnya: layar diberi bentuk TAMPIL `dd/mm/yy` (tahun DUA digit, garis
// miring), sementara kotak tanggal menerima bentuk kabel `DD-MM-YYYY`.
// Tidak terbaca, dan kosongnya tanpa satu pun galat.
func TestTanggalKursPunyaBentukTersimpanDanTampil(t *testing.T) {
	const stempel = "20180101T000000.000 GMT"

	// Bentuk TAMPIL — untuk dibaca.
	if got := services.TanggalTampil(stempel); got != "01/01/18" {
		t.Fatalf("TanggalTampil = %q, mau 01/01/18", got)
	}
	// Bentuk TERSIMPAN — untuk kotak tanggal.
	if got := services.TanggalWIB(stempel); got != "20180101" {
		t.Fatalf("TanggalWIB = %q, mau 20180101", got)
	}
	// ⛔ Keduanya BERBEDA, dan itu seluruh sebab medan `…Asli` ada.
	if services.TanggalTampil(stempel) == services.TanggalWIB(stempel) {
		t.Fatal("kedua bentuk sama — medan `…Asli` jadi tidak berguna")
	}
}

// ⛔ Dan dari bentuk TAMPIL, tahun empat digitnya TIDAK dapat dipulihkan —
// itulah mengapa `…Asli` wajib diisi SEBELUM medan tampil ditimpa.
func TestBentukTampilKehilanganAbad(t *testing.T) {
	if services.TanggalWIB("01/01/18") == "20180101" {
		t.Fatal("bentuk tampil ternyata dapat dipulihkan — urutan pengisian jadi tidak mengikat")
	}
}
