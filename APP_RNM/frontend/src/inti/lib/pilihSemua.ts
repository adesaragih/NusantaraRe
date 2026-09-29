// Penjungkit `Select All` — ditiru dari `SelectAllClaimLife_act`.
//
// Sumber, dibaca sebagai pohon langkah utuh:
//
//   `Activity/SelectAllClaimLife_act.xml`
//     b223 `Property-Set`
//       b247 Select.CARI1 = @if(Select.CARI1=="","true",
//                               @if(Select.CARI1=="true","false","true"))
//     b353 langkah ber-ULANG(EMBEDDED) b550 atas
//          `pyWorkPage.ClaimData.PremiumListSummary.PremiumListDetail`
//       b402 `Property-Set`
//         b429 .IsAccept = Select.CARI1
//
// ⛔ LAYAR MANA: tombolnya berdiri di `Section/InputOSClaimLife.xml`
// (b16633, b16710) — layar **Outstanding**, BUKAN Register. Diperiksa: nol
// `Select All` dan nol `IsAccept` di `InputRegisterClaimLife.xml`. Karena itu
// modul ini belum punya pemanggil; ia menunggu grid peserta layar Outstanding
// (kelompok Detail & Tutup). Aturannya ditulis sekarang, bersama ujinya,
// supaya ia tidak lahir dari ingatan ketika gridnya tiba.

/**
 * Keadaan penjungkit, disimpan sebagai TEKS persis seperti `Select.CARI1`.
 *
 * ⛔ Tiga keadaan, bukan dua. `''` berarti tombolnya belum pernah ditekan,
 * dan Pega membedakannya dari `'false'`: tekanan pertama pada daftar yang
 * seluruhnya belum tercentang harus MENCENTANG, bukan melepas. Boolean tidak
 * dapat membedakan keduanya tanpa penunjuk kedua.
 */
export type KeadaanPilihSemua = '' | 'true' | 'false'

/**
 * jungkitPilihSemua mengembalikan keadaan berikutnya.
 *
 * b247 apa adanya: kosong → `'true'`; `'true'` → `'false'`; selain itu →
 * `'true'`. Cabang terakhir sengaja menampung nilai tak dikenal, sama seperti
 * `@if` bersarang Pega — nilai asing berarti "belum tercentang".
 */
export function jungkitPilihSemua(sekarang: KeadaanPilihSemua | string): KeadaanPilihSemua {
  if (sekarang === '') return 'true'
  if (sekarang === 'true') return 'false'
  return 'true'
}

/**
 * terpilihSemua menjawab apakah keadaan itu berarti seluruh baris tercentang.
 *
 * Dipakai layar untuk menggambar centangnya. `IsAccept` di Pega menyimpan
 * teks yang SAMA dengan `Select.CARI1` (b429), jadi pembacaannya satu aturan.
 */
export function terpilihSemua(keadaan: KeadaanPilihSemua | string): boolean {
  return keadaan === 'true'
}
