// Label modul Plan - bahasa Inggris, VERBATIM dari section Pega `InboxProductType`
// (`D:\NUSARE DEV\Menu Plan\InboxProductType.xml`, nomor = baris XML): judul "Plan" b331, form "Plan Name" b812 / b834,
// "Business" b1079, "Benefit" b1432 / b1455, tombol Save b6354 / b6403 dan New b6615 / b6667, grid "Plan Name" b2978 /
// "Business" b3117 / "Benefit" b3255 / kolom tombol tanpa judul b3371 berisi Edit b3975 / b4027. Kolom autocomplete
// Business: OLDID b1200, Note b1231, ID b1265. Teks lain = teks antarmuka aplikasi, bukan label Pega.

/** Nama menu - `M_NAV_MENU.LABEL` baris modul ini (migrasi inti 949, keputusan work owner 08-10-2026 K6). */
export const MENU_PL = { kelompok: 'Plan' } as const

export const PL = {
  judul: 'Plan',
  planName: 'Plan Name',
  business: 'Business',
  benefit: 'Benefit',
  save: 'Save',
  menyimpan: 'Saving…',
  baru: 'New',
  edit: 'Edit',
  oldId: 'OLDID',
  note: 'Note',
  id: 'ID',

  modeUbah: (nama: string) => `Editing ${nama}`,
  galatWajib: (medan: string) => `${medan} is required`,
  galatPanjang: (medan: string, batas: number) => `${medan} is longer than ${batas} characters`,
  tersimpan: (nama: string) => `Plan ${nama} saved.`,
  pilihDariDaftar: 'Type to search, then choose from the list.',
  tanpaPilihan: 'No match in the list.',
  memuat: 'Loading plans…',
  kosong: 'No plan yet.',
  jumlah: (n: number) => `${n} plan${n === 1 ? '' : 's'}`,
  halaman: (h: number, dari: number) => `Page ${h} of ${dari}`,
  sebelum: 'Previous',
  sesudah: 'Next',
} as const
