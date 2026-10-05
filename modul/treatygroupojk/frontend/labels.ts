// Label modul Treaty Group OJK - bahasa Inggris. Pega tidak punya layar master tabel ini; caption kolom mengikuti nama
// kolom `TREATYGROUPOJK` (keputusan work owner 05-10-2026).

/** Nama menu - `M_NAV_MENU.LABEL` baris modul ini (migrasi inti 916). */
export const MENU_TGO = { kelompok: 'Treaty Group OJK' } as const

export const TGO = {
  judul: 'Treaty Group OJK',
  tambah: 'Add',
  cari: 'Search ID or name',
  memuat: 'Loading treaty group OJK…',
  kosong: 'No treaty group OJK yet.',
  tidakCocok: 'No treaty group OJK matches the search.',
  jumlah: (n: number) => `${n} OJK business${n === 1 ? '' : 'es'}`,

  no: 'No',
  id: 'ID',
  nama: 'Name',
  namaIdn: 'Name (IDN)',
  aksi: 'Action',
  edit: 'Edit',

  judulTambah: 'Add Treaty Group OJK',
  judulUbah: (id: string) => `Edit Treaty Group OJK ${id}`,
  rincian: 'Names',
  idOtomatis: 'Generated when saved',
  catatanHurufBesar: 'Name and Name (IDN) are saved in capital letters.',
  galatNama: 'Name is required',
  galatNamaIdn: 'Name (IDN) is required',
  save: 'Save',
  menyimpan: 'Saving…',
  batal: 'Cancel',
  tutup: 'Close',
  tersimpan: (id: string) => `Treaty group OJK ${id} saved.`,
} as const
