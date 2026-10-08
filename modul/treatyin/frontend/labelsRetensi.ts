// Label tab **Maximum Retention** (cabang NON-PROPORSIONAL) — DISALIN dari
// ekspor, bukan dikarang.
//
// Sumbernya dua berkas:
//
//   Section/TreatyInTabsNonProportional.xml   wadah TABBED ke-1
//   Section/MaxRetention.xml                  rincian satu baris
//
// ⭐ Judul panel SAMA dengan nama tabnya (`Maximum Retention`) — itu bunyi
// ekspornya, dan tangkapan layar pemilik proses memperlihatkan keduanya
// tertulis dua kali di layar yang sama.
//
// ⚠️ Label total panel TIDAK diulang di sini: ia sudah hidup di `labels.ts`
// sebagai `TOTAL_RETENSI`, lengkap dengan pengukuran offset dan catatan
// tombol kembar yang mati. Dua tempat untuk satu teks adalah cara termudah
// keduanya berbeda diam-diam.

/** Kepala kolom grid `Maximum Retention` — urut ekspor. */
export const KOLOM_GRID_RETENSI = ['Treaty Group', 'Currency', 'Amount'] as const

export const RETENSI = {
  judul: 'Maximum Retention',
  panel: 'Maximum Retention',

  // Rincian baris — `Section/MaxRetention.xml`.
  treatyGroup: 'Treaty Group',
  /**
   * ⭐ Di rinciannya, `Amount` adalah label SATU BARIS yang memuat DUA medan:
   * pilihan mata uang lalu kotak nilai. `pyLabelFieldValue = Amount`
   * menempel pada `.Currency`, dan `.Amount` di sebelahnya tanpa label
   * sendiri. Tangkapan layar pemilik proses memperlihatkan persis itu.
   */
  jumlah: 'Amount',
  mataUang: 'Currency',
  keterangan: 'Note',

  // Tombol — label `pyLabel` apa adanya.
  tambah: 'Add',
  hapus: 'Delete',
  // ⛔ `Update Total` SENGAJA tidak ada di sini: ia sudah hidup di
  // `labels.ts` sebagai `TOTAL_RETENSI.perbarui`, bersama pengukuran offset
  // tombolnya. Menyalinnya ke sini membuat dua tempat untuk satu teks.

  // Teks layar kita sendiri — ditandai supaya tidak tertukar dengan ekspor.
  barisBaru: '(belum dipilih)',
  petunjukKosong: 'Belum ada baris Maximum Retention.',
} as const

/**
 * ⚠️ `Amount` di GRID gambar 26 berbunyi `3.500.000.000` — NOL desimal —
 * sementara panel `Total Retention Amount` tepat di bawahnya berbunyi
 * `3.500.000.000,00`. Nilai yang sama, dua presisi, dua tempat. Itu bukan
 * kekeliruan tangkapan layar; keduanya dipertahankan.
 */
export const DESIMAL_RETENSI = {
  jumlah: 0,
  nilaiTotal: 2,
} as const

/**
 * ⛔ `Update Total` MATI bila `TreatyIn.EDMMaterialType = 2`.
 *
 * `pyDisabledWhen` @202558, terukur di keempat sel tombol sejenis. Dicatat
 * sebagai fungsi supaya syaratnya satu tempat, bukan tersebar sebagai
 * perbandingan harfiah di tiap komponen.
 */
export function totalTerkunci(edmJenisMaterial: string): boolean {
  return edmJenisMaterial === '2'
}
