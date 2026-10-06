// Label layar modul Treaty In Adjustment.

export const MENU_TREATYINADJUSTMENT = {
  rantaiVersi: 'Treaty In Adjustment — Rantai Versi',
} as const

export const RANTAI_VERSI = {
  pilihKontrak: 'Kontrak',
  // Judul kartu SEBELUM ada kontrak yang dipilih; sesudahnya judulnya nama kontrak itu.
  panelPilih: 'Pilih kontrak',
  kolomNomor: 'No. urut',
  kolomNama: 'Nama kontrak',
  kolomKeadaan: 'Keadaan',
  kolomJenis: 'Jenis addendum',
  kolomMaterial: 'Materialitas',
  kolomBerlaku: 'Berlaku sejak',
  kolomDasar: 'Dasar',
  belumDinomori: 'belum dinomori',
  versiPertama: 'versi pertama',
  // Dipisah karena `Kosong` punya dua medan: `pesan` = keadaan, `petunjuk` =
  // siapa yang akan mengisinya.
  kosong: 'Belum ada kontrak tercatat.',
  kosongPetunjuk: 'Pemindahan kepala kontrak warisan adalah tiket 59.',
  keterangan:
    'Rantai versi sebuah kontrak, baca-saja. Kolom "Dasar" adalah rujukan eksplisit ke versi berlaku terakhir saat versi itu dibuat (tiket 01); kosong berarti versi pertama. Kolom "No. urut" boleh kosong sampai penomoran ulang baris warisan selesai (tiket 10).',
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
  petunjukLampiran:
    'Lampiran dibaca dari M_ATTACHMENTTREATY_2, tabel sistem lama. Kosong berarti kontrak ini memang belum punya berkas.',
  petunjukHistory:
    'Riwayat persetujuan kontrak ini. Kosong berarti belum ada catatan.',
  /**
   * ⛔ Penanda kategori yang pasangan kode↔namanya BELUM dipastikan —
   * diperlakukan SAMA PERSIS seperti di modul Treaty In. Menebak di salah
   * satu modul saja sudah cukup untuk menaruh berkas di kategori yang salah.
   */
  kategoriBelumPasti: 'nama kategori belum dipastikan',
  namaBelumBerumah:
    'Empat nama kategori ada di layar lama tetapi kodenya belum dipastikan, jadi keempatnya belum ditampilkan sebagai nama: Binding, signed share Email · Claim Data · Info Pack · Letter of Acknowledgment / LOA. Lihat treatyin/docs/PERTANYAAN-TERBUKA-KODE-KATEGORI-LAMPIRAN.md.',
  /** Pengenal kontrak warisan yang panelnya tampilkan — sementara, lihat layar. */
  labelPengenal: 'Pengenal kontrak sistem lama',
} as const
