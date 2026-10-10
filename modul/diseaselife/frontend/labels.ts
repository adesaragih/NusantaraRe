// Label modul Disease Life - bahasa Inggris, VERBATIM dari section Pega `InboxDisease`
// (`D:\NUSARE DEV\Menu Disease\InboxDisease.xml`, nomor = baris XML): judul "DISEASE" b371 (`pyCaption DISEASE` b6223),
// form "Number / ID" b1021 (disabled b1070), "ICD Code" b1200 dan "Disease" b1470 (wajib b1465), tombol Save b2102 dan
// Cancel b2365, grid "ID" b3758 / "ICD Code" b3894 / "Disease" b4033 / kolom tombol tanpa judul b4154 dengan Edit b4829.
// Teks lain (saring, halaman, pesan) = teks antarmuka aplikasi, bukan label Pega.

/** Nama menu - `M_NAV_MENU.LABEL` baris modul ini (slot menu modul 951, keputusan work owner 08-10-2026 D4). */
export const MENU_DSL = { kelompok: 'Disease Life' } as const

export const DSL = {
  judul: 'DISEASE',
  nomorId: 'Number / ID',
  id: 'ID',
  icdCode: 'ICD Code',
  disease: 'Disease',
  save: 'Save',
  menyimpan: 'Saving…',
  batal: 'Cancel',
  edit: 'Edit',

  idOtomatis: 'Generated when saved',
  modeUbah: (id: string) => `Editing ID ${id}`,
  galatICD: 'ICD Code is required',
  galatDisease: 'Disease is required',
  galatPanjang: (medan: string, batas: number) => `${medan} is longer than ${batas} characters`,
  tersimpan: (id: string) => `Disease ${id} saved.`,
  saringICD: 'Filter ICD Code',
  saringDisease: 'Filter Disease',
  memuat: 'Loading diseases…',
  kosong: 'No disease yet.',
  tidakCocok: 'No disease matches the filter.',
  jumlah: (n: number) => `${n} disease${n === 1 ? '' : 's'}`,
  halaman: (h: number, dari: number) => `Page ${h} of ${dari}`,
  sebelum: 'Previous',
  sesudah: 'Next',
} as const
