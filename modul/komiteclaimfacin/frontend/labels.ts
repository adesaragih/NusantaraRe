// Label modul Komite Claim Fac In. Label LAYAR KASUS (judul, ubin, judul bagian, medan, judul grid, kolom, isian, tombol)
// datang dari server VERBATIM Section ShowTransfer (`models.Layar`). Berkas ini hanya memuat teks perangkat yang tidak
// punya label di server (ditandai). Pola `modul/komiteclaimnonprop/frontend/labels.ts` (disalin, bukan impor).

/** Nama modul - nama folder korpus VERBATIM (`Komite Claim FacIn`, tanpa spasi; `LayarModul.kelompok`). TANPA menu. */
export const KORPUS_KCFI = { kelompok: 'Komite Claim FacIn' } as const

/** Teks perangkat layar kasus komite yang tidak punya label di server (ditandai). */
export const KCFI = {
  memuat: 'Loading…',
  /** `[tidak ada di korpus]` - layar dibuka bukan oleh pemegang (mis. kasus sudah berpindah tingkat). */
  hanyaLihat: 'Read only - this case is waiting for another committee member.',
  memproses: 'Submitting…',
  /** `[tidak ada di korpus]` - View more details: menu Claim Fac In tidak dipegang akun ini (teks sama dengan Claim Prop). */
  berkasTakTerpasang: 'This file cannot be opened here: the menu of its module is not assigned to your account.',
  /** `[tidak ada di korpus]` - kembali ke inbox Claim Fac In sesudah pesan Submit (kasus dibuka dari tabel komite). */
  kembali: 'Back',
  /** pxLink kolom Spreading Adjustment (DetailAdjustmentFac `VIS .TreatyType = '10015'`, local action ShowRetro). */
  lihatRetro: 'View Retro',
  /** Tombol penutup pop-up ShowRetro (sama dengan jendela berkas App). */
  tutup: 'Close',
} as const

/**
 * Tata letak layar komite (pola Komite Claim Prop / Non Prop, keputusan work owner 09-10-2026) - lencana keadaan kepala,
 * kartu keputusan, dan grid kosong `[tidak ada di korpus]`. Judul kartu, label medan / grid / isian tetap dari server.
 */
export const TATA_KCFI = {
  statusGiliranAnda: 'Waiting for your decision',
  statusMenunggu: 'Waiting for',
  statusSelesai: 'Completed',
  keputusanAnda: 'Your Decision',
  tanpaData: 'No data',
  /** aria-label kolom panah grid berincian (expand pane). */
  rincian: 'Details',
} as const
