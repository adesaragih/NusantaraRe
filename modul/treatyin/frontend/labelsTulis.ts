// Label tombol tulis form Treaty In — Save, Submit, Actions, Decline offer.
// ⛔ Teks tombol dan modal DISALIN dari ekspor (`TreatyInActionButtons`,
// `TreatyInAction`, `TreatyInfoSubmit`, `TreatyInDeclineConfirmation`).

export const TOMBOL_TULIS = {
  simpan: 'Save', // `TreatyInActionButtons` @20066
  tutup: 'Close', // @36946
  aksi: 'Actions', // @54595 → local action `TreatyInAction`
  menyimpan: 'Menyimpan…',
  /** Modal `TreatyInAction`. */
  judulAksi: 'Action',
  dariID: 'Of ID :', // `TreatyIn.ID` pxDisplayText @25725
  pilihan: 'Status', // `TreatyIn.ChooseStatusAkseptasi` pxRadioButtons @57957 (tanpa label di ekspor)
  komentar: 'Comment', // `TreatyIn.Comment` pxTextArea @65533
  batal: 'Cancel', // @91045
  kirim: 'Submit', // @112654 → `TreatyInAkseptasi_Act`
  /** Modal `TreatyInDeclineConfirmation`. */
  judulTolak: 'Decline offer',
  tanyaTolak: 'Are you sure you want to DECLINE this offer?', // @14478
  tolak: 'Decline', // @65981 → `TreatyInDeclineConfirmation_postact`
  /** Properti terkirim tanpa kolom di tabel pendaratan — dilaporkan, tidak ditelan. */
  takTersimpan: 'Belum punya kolom di tabel, jadi TIDAK tersimpan:',
} as const

/**
 * Pilihan `ChooseStatusAkseptasi` — nilai yang `Akseptasi_DT` periksa.
 * ⚠️ Daftarnya `associated` (rule properti tidak diekspor); teksnya nilai itu
 * sendiri.
 */
export const PILIHAN_AKSEPTASI = ['Accept', 'Reject', 'Decline'] as const
