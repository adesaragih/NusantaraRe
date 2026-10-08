// Asal: salinan `modul/nbtreatyin/frontend/asalKasus.ts` (06-10-2026), kelas `nbti__` -> `edmt__`.
// Asal berkas yang sedang terbuka - permintaan work owner 06-10-2026: "tombol back nya bisa ngetrack darimana
// bukanya, dari beranda atau NB TREATY IN nya". Berkas dari kotak masuk Beranda kembali ke Beranda; berkas dari
// portal kembali ke portal.

/** Tempat berkas dibuka. */
export type AsalKasus = 'portal' | 'beranda'

/**
 * Tujuan sesudah layar kasus ditutup. Hanya tombol Back (tanpa pesan) yang mengikuti asal; sesudah Submit (ada
 * pesan) tetap ke portal, tempat pesan "terkirim" tampil. Tanpa jalan ke Beranda = portal.
 */
export function tujuanKembali(asal: AsalKasus, pesan: string | undefined, adaBeranda: boolean): AsalKasus {
  return asal === 'beranda' && pesan === undefined && adaBeranda ? 'beranda' : 'portal'
}
