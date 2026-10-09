// Label layar modul Treaty In Adjustment.

export const MENU_TREATYINADJUSTMENT = {
  rantaiVersi: 'Treaty In Adjustment — Version Chain',
} as const

export const RANTAI_VERSI = {
  pilihKontrak: 'Contract',
  // Judul kartu SEBELUM ada kontrak yang dipilih; sesudahnya judulnya nama kontrak itu.
  panelPilih: 'Select a contract',
  kolomNomor: 'No.',
  kolomNama: 'Contract Name',
  kolomKeadaan: 'Keadaan',
  kolomJenis: 'Addendum Type',
  kolomMaterial: 'Materialitas',
  kolomBerlaku: 'Valid From',
  kolomDasar: 'Dasar',
  belumDinomori: 'not numbered yet',
  versiPertama: 'first version',
  // Dipisah karena `Kosong` punya dua medan: `pesan` = keadaan, `petunjuk` =
  // siapa yang akan mengisinya.
  kosong: 'No contracts recorded yet.',
  kosongPetunjuk: '',
  keterangan: '',
} as const

// ===========================================================================
// PANEL ATTACHMENT + HISTORY
// Sumber: `Section/WorkAttachments.xml` dan `ShowAttachmentTreaty.xml` —
// keduanya ada di ekspor modul INI, bukan dipinjam dari Treaty In.
// ===========================================================================

/** Kolom panel Attachment — `pyCaption` di `WorkAttachments.xml` apa adanya. */
export const KOLOM_LAMPIRAN = ['Category', 'Count', 'Upload file', 'View File'] as const

/** Kolom daftar berkas di dalam kategori. */
export const KOLOM_BERKAS_LAMPIRAN = ['File Name', 'Type', 'Uploaded', 'By'] as const

/** Kolom panel History — Date · PIC · Approval · Comment. */
export const KOLOM_HISTORY = ['Date', 'PIC', 'Approval', 'Comment'] as const

export const LAMPIRAN = {
  judul: 'Attachment',
  judulHistory: 'History',
  /** `pyCaption Recommended safe substitute should be . or _` — disalin apa adanya. */
  spanduk: 'Recommended safe substitute should be . or _',
  unduhSemua: 'Download All', // pyButtonLabel DOWNLOAD ALL
  segarkan: 'Refresh', // pyButtonLabel REFRESH
  /** Teks kosong layar lama, apa adanya. */
  tanpaIsi: 'No items',
  petunjukLampiran: '',
  petunjukHistory:
    'Approval history for this contract. Empty means no entries yet.',
  /** Pengenal kontrak warisan yang panelnya tampilkan — sementara, lihat layar. */
  labelPengenal: 'Contract ID',
} as const
