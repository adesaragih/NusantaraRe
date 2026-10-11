// Label modul Benefit - bahasa Inggris, VERBATIM dari section Pega `InboxBenefit`
// (`D:\NUSARE DEV\Menu Benefit\InboxBenefit.xml`, nomor = baris XML): judul "INSURANCE BENEFIT" b313 / b5616, form
// "Number / ID" b964 / b987 (disabled) dan "Benefit" b1143 / b1166 (wajib), tombol Save b1723 / b1772 dan Cancel b1985 /
// b2035, grid "ID" b3431 / "Benefit" b3570 / kolom tombol tanpa judul b3708 dengan Edit b4168 / b4220. Teks lain (saring,
// halaman, pesan) = teks antarmuka aplikasi, bukan label Pega.

/** Nama menu - `M_NAV_MENU.LABEL` baris modul ini (migrasi inti 945, keputusan work owner 08-10-2026 K4). */
export const MENU_BN = { kelompok: 'Benefit' } as const

export const BN = {
  judul: 'INSURANCE BENEFIT',
  nomorId: 'Number / ID',
  id: 'ID',
  benefit: 'Benefit',
  save: 'Save',
  menyimpan: 'Saving…',
  batal: 'Cancel',
  edit: 'Edit',

  idOtomatis: 'Generated when saved',
  modeUbah: (id: string) => `Editing ID ${id}`,
  galatBenefit: 'Benefit is required',
  galatPanjang: (batas: number) => `Benefit is longer than ${batas} characters`,
  tersimpan: (id: string) => `Benefit ${id} saved.`,
  saringID: 'Filter ID',
  saringBenefit: 'Filter Benefit',
  memuat: 'Loading benefits…',
  kosong: 'No benefit yet.',
  tidakCocok: 'No benefit matches the filter.',
  jumlah: (n: number) => `${n} benefit${n === 1 ? '' : 's'}`,
  halaman: (h: number, dari: number) => `Page ${h} of ${dari}`,
  sebelum: 'Previous',
  sesudah: 'Next',
} as const
