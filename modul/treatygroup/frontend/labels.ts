// Label modul Treaty Group - bahasa Inggris. Caption kolom mengikuti `BrowseTreatyGroup_RD` Pega (ID, TreatyGroupName,
// TreatyGroupSOAName, OJKBusinessName) dan keputusan work owner 05-10-2026.

/** Nama menu - `M_NAV_MENU.LABEL` baris modul ini (migrasi inti 917). */
export const MENU_TG = { kelompok: 'Treaty Group' } as const

export const TG = {
  judul: 'Treaty Group',
  tambah: 'Add',
  cari: 'Search ID or name',
  semuaOjk: 'All OJK Business',
  memuat: 'Loading treaty groups…',
  kosong: 'No treaty group yet.',
  tidakCocok: 'No treaty group matches the search.',
  jumlah: (n: number) => `${n} treaty group${n === 1 ? '' : 's'}`,

  no: 'No',
  id: 'ID',
  oldId: 'Old ID',
  ojk: 'OJK Business',
  ojkIdn: 'OJK Business (IDN)',
  nama: 'Treaty Group Name',
  soa: 'SOA Name',
  coa: 'COA',
  pengubah: 'Last Edited By',
  tglUbah: 'Last Edited',
  aksi: 'Action',
  edit: 'Edit',

  judulTambah: 'Add Treaty Group',
  judulUbah: (id: string) => `Edit Treaty Group ${id}`,
  kartuCoa: 'COA & SOA',
  jejak: (id: string, oleh: string, kapan: string) =>
    [`ID ${id}`, oleh === '' ? '' : `Last edited by ${oleh}`, kapan].filter((s) => s !== '').join(' · '),
  idOtomatis: 'Generated when saved',
  pilih: '— select —',
  nilaiLama: (s: string) => `${s} (old value)`,
  coaOtomatis: 'Follows the OJK Business',
  coaKosong: 'No COA in this OJK Business yet',
  catatan: 'Names are saved in capital letters. COA is not typed: it follows the OJK Business.',
  galatOjk: 'OJK Business is required',
  galatNama: 'Treaty Group Name is required',
  anak: 'Business Groups',
  anakKosong: 'No business group under this treaty group.',
  alias: 'Alias Name',
  save: 'Save',
  menyimpan: 'Saving…',
  batal: 'Cancel',
  tutup: 'Close',
  tersimpan: (id: string) => `Treaty group ${id} saved.`,
} as const
