// Label modul Cause Of Loss Life - bahasa Inggris, VERBATIM dari section Pega `InboxCauseofLossLife`
// (`D:\NUSARE DEV\Menu Cause Of Loss\InboxCauseofLossLife.xml`, nomor = baris XML): judul "CAUSE OF LOSS" b343
// (`pyCaption CAUSE OF LOSS` b6533), form "ID" b823 (disabled b872) dan "Cause of Loss" b1001 (wajib b996), tombol Save
// b5366 dan Cancel b5624, grid "ID" b2989 / "Cause of Loss" b3127 / kolom tombol tanpa judul b3245 dengan Edit b3747.
// Teks lain (halaman, pesan) = teks antarmuka aplikasi, bukan label Pega.

/** Nama menu - `M_NAV_MENU.LABEL` baris modul ini (slot menu modul 955, keputusan work owner 08-10-2026 K5). */
export const MENU_COL = { kelompok: 'Cause Of Loss Life' } as const

export const COL = {
  judul: 'CAUSE OF LOSS',
  id: 'ID',
  causeOfLoss: 'Cause of Loss',
  save: 'Save',
  menyimpan: 'Saving…',
  batal: 'Cancel',
  edit: 'Edit',

  idOtomatis: 'Generated when saved',
  modeUbah: (id: string) => `Editing ID ${id}`,
  galatWajib: 'Cause of Loss is required',
  galatPanjang: (batas: number) => `Cause of Loss is longer than ${batas} characters`,
  tersimpan: (id: string) => `Cause of Loss ${id} saved.`,
  memuat: 'Loading causes of loss…',
  kosong: 'No cause of loss yet.',
  jumlah: (n: number) => `${n} cause${n === 1 ? '' : 's'} of loss`,
  halaman: (h: number, dari: number) => `Page ${h} of ${dari}`,
  sebelum: 'Previous',
  sesudah: 'Next',
} as const
