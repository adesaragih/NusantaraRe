// Label modul Business Group - bahasa Inggris. Caption kolom mengikuti `BrowseBusinessGroup_RD` Pega (ID, Note,
// AliasName, TOPID, TreatyName) dan keputusan work owner 05-10-2026.

/** Nama menu - `M_NAV_MENU.LABEL` baris modul ini (migrasi inti 918). */
export const MENU_BG = { kelompok: 'Business Group' } as const

export const BG = {
  judul: 'Business Group',
  tambah: 'Add',
  cari: 'Search ID, name, or alias',
  semuaTreaty: 'All Treaty Groups',
  memuat: 'Loading business groups…',
  kosong: 'No business group yet.',
  tidakCocok: 'No business group matches the search.',
  jumlah: (n: number) => `${n} business group${n === 1 ? '' : 's'}`,

  no: 'No',
  id: 'ID',
  nama: 'Name',
  alias: 'Alias Name',
  treaty: 'Treaty Group',
  aksi: 'Action',
  view: 'View',
  edit: 'Edit',

  judulTambah: 'Add Business Group',
  judulUbah: (id: string) => `Edit Business Group ${id}`,
  judulLihat: (id: string) => `Business Group ${id}`,
  rincian: 'Details',
  idOtomatis: 'Generated when saved',
  pilih: '— select —',
  nilaiLama: (s: string) => `${s} (old value)`,
  aliasOtomatis: 'Same as Name when left empty',
  catatan: 'Names are saved in capital letters. Business groups ending with SYARIAH are not managed here.',
  galatTreaty: 'Treaty Group is required',
  galatNama: 'Name is required',
  galatSyariah: 'Business groups ending with SYARIAH are not managed here',
  save: 'Save',
  menyimpan: 'Saving…',
  batal: 'Cancel',
  tutup: 'Close',
  tersimpan: (id: string) => `Business group ${id} saved.`,
} as const
