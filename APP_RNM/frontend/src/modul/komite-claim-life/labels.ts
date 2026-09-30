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
  /**
   * Tiket 09 — judul VERBATIM grid `.KomiteList` `Section/ShowTransfer.xml`
   * b29727: `Committee` (`.IDKomite` = jabatan), `Status` (`.KomiteAproval`),
   * `Date Approve`, `Comment`.
   */
  committee: 'Committee',
  status: 'Status',
  dateApprove: 'Date Approve',
  comment: 'Comment',
  /** `KOMITE_OPERATORID` — pengenal akun; kosakata kami. */
  anggota: 'Anggota',
  /**
   * OQ-K-05 (GILIRAN-17) — keputusan tingkat sebelum langkah 5.1 menimpanya.
   * Kosakata kami: Pega tidak menampilkan riwayat yang ditimpanya.
   */
  asli: 'Sebelum ditimpa Tolak akhir',
  eskalasi: 'Eskalasi',
  giliranAnda: 'Giliran Anda memutuskan.',
  bukanGiliran: 'Bukan giliran Anda — kasus ini menunggu tingkat lain atau sudah selesai.',
} as const

/**
 * Layar keputusan — tiket 02, `Section/ShowTransfer.xml`.
 *
 * ⚠️ Pilihan `.AcceptStatus` (b32607, `pyListSource associated`) milik aturan
 * properti yang TIDAK diekspor; `1 = Setuju`, `2 = Tolak` adalah
 * `[keputusan work owner]`, bukan teks korpus. Kalimat konfirmasi dan tombol
 * VERBATIM.
 */
export const KEPUTUSAN_KOMITE = {
  /** VERBATIM — label di atas dropdown `.AcceptStatus`. */
  konfirmasi: 'Are you sure to accept this document?',
  pilih: 'Pilih keputusan',
  setuju: 'Setuju',
  tolak: 'Tolak',
  /** `.KomiteComment` b31001 — `pyRequired false`. */
  komentar: 'Comment',
  /** VERBATIM b34722. */
  submit: 'Submit',
  /** VERBATIM b33880. */
  cancel: 'Cancel',
  wajibPilih: 'Keputusan wajib dipilih.',
} as const

/**
 * Efek keluar dan laporan "perlu intervensi" — tiket 08. ⚠️ Tidak ada di
 * korpus (perilaku baru, ADR-0015); kosakata kami.
 */
export const EFEK_KOMITE = {
  judul: 'Efek keluar',
  jenis: 'Efek',
  keadaan: 'Keadaan',
  percobaan: 'Percobaan',
  sejak: 'Sejak',
  laporan: 'Perlu intervensi hari ini',
  laporanKosong: 'Laporan hari ini: tidak ada efek keluar yang perlu intervensi.',
} as const

/**
 * Eskalasi — tiket 03. ⚠️ Tidak ada di korpus (penyimpangan sadar, ADR-0014);
 * kosakata kami.
 */
export const ESKALASI_KOMITE = {
  tombol: 'Eskalasi naik satu tingkat',
  keterangan:
    'Melewati anggota tingkat berjalan yang berhalangan. Tingkat yang dilewati ' +
    'tidak mencatat keputusan apa pun.',
} as const
