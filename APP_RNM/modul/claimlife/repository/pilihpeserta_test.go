package repository

import "testing"

// Uji kedua aturan pemilihan nilai peserta - SaveInsuredClaim_Act.
//
// ⛔ Yang dikunci di sini bukan hanya hasilnya, melainkan ASIMETRI antara
// kedua aturan. Keduanya terlihat mirip dan seorang pembaca yang rapi akan
// tergoda menyeragamkannya; menyeragamkannya mengubah angka uang.

func TestUmurPeserta(t *testing.T) {
	// `Section/.../SaveInsuredClaim_Act.xml:661`
	//   @if(.AGE!="",.AGE,@if(.ENTRY_AGE!="",.ENTRY_AGE,.CURRENT_AGE))
	kasus := []struct {
		nama                      string
		age, entryAge, currentAge string
		mau                       string
	}{
		{"AGE terisi dipakai", "35", "30", "36", "35"},
		{"AGE kosong jatuh ke ENTRY_AGE", "", "30", "36", "30"},
		{"keduanya kosong jatuh ke CURRENT_AGE", "", "", "36", "36"},
		{"ketiganya kosong tetap kosong", "", "", "", ""},
		// ⛔ INI ASIMETRINYA. Untuk umur, "0" BUKAN kosong: bayi berumur nol
		// tahun adalah umur yang sah, dan XML hanya membandingkan dengan "".
		// Memperlakukan "0" sebagai kosong di sini akan diam-diam menaikkan
		// umur bayi menjadi ENTRY_AGE.
		{"umur NOL dipakai apa adanya, bukan dianggap kosong", "0", "30", "36", "0"},
		{"ENTRY_AGE nol pun dipakai", "", "0", "36", "0"},
	}
	for _, k := range kasus {
		t.Run(k.nama, func(t *testing.T) {
			if got := UmurPeserta(k.age, k.entryAge, k.currentAge); got != k.mau {
				t.Errorf("UmurPeserta(%q,%q,%q) = %q, mau %q",
					k.age, k.entryAge, k.currentAge, got, k.mau)
			}
		})
	}
}

func TestShareNusantaraReTeks(t *testing.T) {
	// `SaveInsuredClaim_Act.xml:601` dan `:1412` - ungkapan yang SAMA
	//   @if(.SHARE_NUSANTARA_RE=="0", GROSS,
	//       @if(.SHARE_NUSANTARA_RE=="", GROSS, SHARE))
	kasus := []struct {
		nama         string
		share, gross string
		mau          string
	}{
		{"share terisi dipakai", "12.5", "99", "12.5"},
		// ⛔ ASIMETRI: di sini "0" MEMANG dianggap kosong. Nol share berarti
		// baris polisnya belum dibagi, bukan bahwa bagian Nusantara Re nol -
		// dan pyDeleteMemo rule ini berbunyi "fix share nusantara re".
		{"share NOL jatuh ke GROSS", "0", "99", "99"},
		{"share kosong jatuh ke GROSS", "", "99", "99"},
		{"keduanya kosong tetap kosong", "", "", ""},
		{"share nol dan gross kosong tetap kosong", "0", "", ""},
		{"nol berdesimal BUKAN nol menurut XML", "0.00", "99", "0.00"},
	}
	for _, k := range kasus {
		t.Run(k.nama, func(t *testing.T) {
			if got := ShareNusantaraReTeks(k.share, k.gross); got != k.mau {
				t.Errorf("ShareNusantaraReTeks(%q,%q) = %q, mau %q",
					k.share, k.gross, got, k.mau)
			}
		})
	}
}
