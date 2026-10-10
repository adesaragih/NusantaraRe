// Label tab **EGNPI** (cabang NON-PROPORSIONAL) — DISALIN dari ekspor,
// bukan dikarang.
//
// Sumbernya dua berkas:
//
//   Section/TreatyInTabsNonProportional.xml   wadah TABBED ke-3 `EGNPI`
//   Section/DetailEGNPI.xml                   rincian satu baris
//
// ⚠️ GRID DAN RINCIAN MEMAKAI KATA YANG BERBEDA untuk medan yang SAMA:
// kepala kolom grid berbunyi `As Date`, sementara rinciannya `As At`.
// Keduanya ditulis apa adanya di bawah — menyeragamkannya berarti memilih
// salah satu tanpa dasar, dan yang membaca layar ini berikutnya tidak akan
// tahu bahwa ekspornya memang berbeda.

/** Kepala kolom grid `Estimate Gross Net Premium Income` — urut ekspor. */
export const KOLOM_GRID_EGNPI = [
  'Treaty Group',
  'As Date',
  'Proportion %',
  'Currency',
  'Amount',
  'Amount in IDR',
] as const

export const EGNPI = {
  judul: 'EGNPI',
  panel: 'Estimate Gross Net Premium Income',

  // Rincian baris — `Section/DetailEGNPI.xml`.
  treatyGroup: 'Treaty Group',
  /** ⚠️ Rinciannya `As At`; kepala kolom gridnya `As Date`. Lihat kepala berkas. */
  asAt: 'As At',
  mataUang: 'Currency',
  jumlah: 'Amount',
  jumlahIDR: 'Amount in IDR',
  proporsi: 'Proportion %',
  keterangan: 'Note',

  // Tombol — label `pyLabel` apa adanya.
  tambah: 'Add',
  hapus: 'Delete',
  perbaruiTotal: 'Update Total',
  perbaruiNilai: 'Update EGNPI Value',

  // Panel bawah.
  totalPerMataUang: 'Total EGNPI Amount',
  nilai: 'Value',
  totalIDR: 'Total Amount in IDR',
  satuanIDR: 'IDR',
  /** Petunjuk kotak `Proportion %` — tangkapan layar Pega. */
  satuanPersen: '%',
  totalProporsi: 'Total Proportion %',

  // Teks layar kita sendiri — ditandai supaya tidak tertukar dengan ekspor.
  barisBaru: '(not selected)',
  tanpaBaris: 'No items',
  petunjukKosong: 'No EGNPI rows yet.',
} as const

/**
 * ⛔ DESIMAL — dibaca dari gambar `29` (§24 `labels.ts`), bukan ditebak.
 *
 * Dua kolom di BARIS YANG SAMA berbeda, dan itu justru buktinya:
 *
 *   `Amount`         `137.849.315.068,00`   2 desimal
 *   `Amount in IDR`  `137.849.315.068`      0 desimal
 *
 * `Proportion %` tampil `100,00` di grid; rinciannya menyimpan 20 desimal.
 */
export const DESIMAL_EGNPI = {
  proporsi: 2,
  jumlah: 2,
  jumlahIDR: 0,
  /** Panel `Total Amount in IDR` — serupa `Amount in IDR`. */
  totalIDR: 0,
  /** Panel `Total Proportion %`. */
  totalProporsi: 2,
  /** Grid `Total EGNPI Amount` — nilainya `Amount`, jadi 2. */
  nilaiPerMataUang: 2,
} as const
