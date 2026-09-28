/**
 * Label modul Komite Claim Life — tiket 01.
 *
 * ⚠️ Tiga kolom pertama Inbox Komite DIPINJAM dari kolom Pega-standar
 * `InboxPremiumList` (`pyFieldLabel` b721/b735/b764): worklist `KomiteRouter`
 * tidak punya section inbox sendiri di korpus. Sisanya kolom komite
 * (klaim induk, tingkat, nilai + `CURRENCY`, status baris).
 */
export const KOLOM_INBOX_KOMITE = {
  /** dipinjam `InboxPremiumList` b721. */
  kasusId: 'Case ID',
  /** dipinjam b735 — di sini waktu pembaruan terakhir work object. */
  tglUpdate: 'Update Date/Time',
  /** dipinjam b764. */
  statusWork: 'Work Status',
  nomorKlaim: 'Klaim induk',
  tingkat: 'Tingkat',
  nilaiKlaim: 'Nilai klaim',
  mataUang: 'CURRENCY',
  statusBaris: 'Status baris',
} as const

export const INBOX_KOMITE = {
  judul: 'Inbox Komite',
  kosong: 'Tidak ada kasus komite yang menunggu keputusan Anda.',
  /** ⛔ Menyebut ATURANNYA, supaya daftar pendek tidak dikira data hilang. */
  aturan:
    'Hanya kasus yang tingkat berjalannya milik Anda: anggota pertama di tangga ' +
    'yang belum memutuskan.',
} as const

export const KASUS_KOMITE = {
  judul: 'Kasus Komite',
  kembali: 'Kembali ke Inbox Komite',
  tangga: 'Tangga persetujuan',
  urut: 'Tingkat',
  jabatan: 'Jabatan',
  approval: 'Keputusan',
  komentar: 'Komentar',
  tanggal: 'Tanggal',
  giliranAnda: 'Giliran Anda memutuskan.',
  bukanGiliran: 'Bukan giliran Anda — kasus ini menunggu tingkat lain atau sudah selesai.',
  /** Layar keputusan `ShowTransfer` menyusul (tiket 02). */
  keputusanMenyusul: 'Layar keputusan (ShowTransfer) dibangun di tiket 02.',
} as const
