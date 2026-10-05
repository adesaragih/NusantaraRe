// Label modul Treaty Description - bahasa Inggris. Caption kolom mengikuti `BrowseTreatyDesc_RD` Pega (ID, Description
// Name) dan rancangan sendiri (keputusan work owner 05-10-2026: "buat versi kamu").

/** Nama menu - `M_NAV_MENU.LABEL` baris modul ini (migrasi inti 920). */
export const MENU_TD = { kelompok: 'Treaty Description' } as const

export const TD = {
  judul: 'Treaty Description',
  tambah: 'Add',
  cari: 'Search ID or description name',
  semua: 'All',
  memuat: 'Loading treaty descriptions…',
  kosong: 'No treaty description yet.',
  tidakCocok: 'No treaty description matches the search.',
  jumlah: (n: number) => `${n} description${n === 1 ? '' : 's'}`,

  no: 'No',
  id: 'ID',
  nama: 'Description Name',
  jenis: 'Type',
  nonXol: 'Non XOL',
  xol: 'XOL',
  status: 'Status',
  aktif: 'Active',
  nonaktif: 'Inactive',
  aksi: 'Action',
  view: 'View',
  edit: 'Edit',

  judulTambah: 'Add Treaty Description',
  judulUbah: (id: string) => `Edit Treaty Description ${id}`,
  judulLihat: (id: string) => `Treaty Description ${id}`,
  klasifikasi: 'Classification',
  idOtomatis: 'Generated when saved',
  catatan: 'Names are saved in capital letters. Descriptions are never deleted: set the status to Inactive instead.',
  galatNama: 'Description Name is required',
  galatPanjang: 'Description Name is longer than 100 characters',
  save: 'Save',
  menyimpan: 'Saving…',
  batal: 'Cancel',
  tutup: 'Close',
  tersimpan: (id: string) => `Treaty description ${id} saved.`,
} as const
