package loader

import "strings"

// kolomSkema - satu baris lembar `Kolom`. medan = FIELD ASLI ("" untuk kolom
// sistem); wajib = NULL "NOT NULL"; turunan = SIKLUS "turunan rancangan" (tidak ada
// di ekspor Pega - FIELD ASLI-nya nama rekaan, bukan nama medan di dokumen).
type kolomSkema struct {
	nama, tipe, medan string
	wajib, turunan    bool
}

// kolomPay - V-22/V-22b: kolom PAY_ ber-FIELD ASLI rekaan "Pay"+medan Payment; medan
// itu dicari dengan awalan lipatan, jadi FIELD ASLI-nya dipetakan apa adanya.
func (k kolomSkema) kolomPay() bool {
	return k.turunan && strings.HasPrefix(k.nama, "PAY_") && strings.HasPrefix(k.medan, "Pay")
}

// angka - kolom bertipe NUMBER: nilainya diurai eksak, bukan disimpan sebagai teks.
func (k kolomSkema) angka() bool { return strings.HasPrefix(k.tipe, "NUMBER") }

// bawaanUnknown - 24 pembawa mata uang K-069 usulan 10: NOT NULL DEFAULT 'UNKNOWN'.
func (k kolomSkema) bawaanUnknown() bool { return strings.Contains(k.tipe, "DEFAULT 'UNKNOWN'") }

// kandidatPenunjuk - satu baris penunjuk lembar `Kandidat Hapus` (workbook
// `Claude outputs`): KATEGORI (R1 indeks-diri, R2 duplikat induk, R3 penunjuk
// leluhur, R3b penunjuk campuran) dan STATUS DI SKEMA INI.
type kandidatPenunjuk struct{ kategori, status string }

// jalurSkema - satu baris lembar `Jalur Sumber`: jalur SESUDAH keputusan V-22/V-24b/
// V-30 (bukan jalur mentah), tanpa ruas kode mata uang (BAHAN §0).
type jalurSkema struct{ jalur, tabel, induk string }
