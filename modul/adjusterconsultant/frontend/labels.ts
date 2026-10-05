// Label modul Adjuster Consultant - bahasa Inggris. Caption kolom mengikuti layar Pega `MstAdjusterConsultant`
// (Claim Fac In / Claim Prop): ID, Name, Address, Telp No, Edit Date; sisanya keputusan work owner 05-10-2026.

/** Nama menu - `M_NAV_MENU.LABEL` baris modul ini (migrasi inti 915, keputusan work owner 05-10-2026). */
export const MENU_ADJ = { kelompok: 'Adjuster Consultant' } as const

export const ADJ = {
  judul: 'Adjuster Consultant',
  tambah: 'Add',
  cari: 'Search ID, name, address, or telp',
  semua: 'All',
  aktif: 'Active',
  nonaktif: 'Inactive',
  memuat: 'Loading adjuster consultants…',
  kosong: 'No adjuster consultant yet.',
  tidakCocok: 'No adjuster consultant matches the search.',
  jumlah: (n: number) => `${n} adjuster consultant${n === 1 ? '' : 's'}`,

  // Kolom (MstAdjusterConsultant) dan tambahan.
  no: 'No',
  id: 'ID',
  nama: 'Name',
  alamat: 'Address',
  telp: 'Telp No',
  pengubah: 'Last Edited By',
  tglUbah: 'Edit Date',
  status: 'Status',
  aksi: 'Action',
  view: 'View',
  edit: 'Edit',
  deactivate: 'Deactivate',
  activate: 'Activate',

  // Form.
  judulTambah: 'Add Adjuster Consultant',
  judulUbah: (id: string) => `Edit Adjuster Consultant ${id}`,
  judulLihat: (id: string) => `Adjuster Consultant ${id}`,
  kontak: 'Contact',
  jejak: (id: string, oleh: string, kapan: string) =>
    [`ID ${id}`, oleh === '' ? '' : `Last edited by ${oleh}`, kapan].filter((s) => s !== '').join(' · '),
  idOtomatis: 'Generated when saved',
  catatanHurufBesar: 'Name and Address are saved in capital letters.',
  galatNama: 'Name is required',
  save: 'Save',
  menyimpan: 'Saving…',
  batal: 'Cancel',
  tutup: 'Close',
  tersimpan: (id: string) => `Adjuster consultant ${id} saved.`,

  // Activate / Deactivate (pengganti hapus, keputusan work owner 05-10-2026).
  judulNonaktif: 'Deactivate adjuster consultant',
  judulAktif: 'Activate adjuster consultant',
  kalimatNonaktif: (nama: string) =>
    `Deactivate ${nama}? It stays in the list as Inactive and can be activated again at any time.`,
  kalimatAktif: (nama: string) => `Activate ${nama} again?`,
  memproses: 'Processing…',
  dinonaktifkan: (id: string) => `Adjuster consultant ${id} deactivated.`,
  diaktifkan: (id: string) => `Adjuster consultant ${id} activated.`,
} as const
