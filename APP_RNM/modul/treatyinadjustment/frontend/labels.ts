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
