// Label modul Reinsurance Type - bahasa Inggris. Caption kolom mengikuti `BrowseReinsuranceType_RD` Pega (ID, Note,
// Type, SOANote, Code, Flag, NoUrut, GroupType); label Type = Prompt value Pega (dari work owner 05-10-2026).

/** Nama menu - `M_NAV_MENU.LABEL` baris modul ini (migrasi inti 921). */
export const MENU_RT = { kelompok: 'Reinsurance Type' } as const

/** Label Type - Prompt value Pega. */
export const LABEL_TYPE: Readonly<Record<string, string>> = {
  '1': 'Own Retention',
  '2': 'Treaty Out',
  '3': 'Facultative',
  '4': 'Treaty In',
}

export const RT = {
  judul: 'Reinsurance Type',
  tambah: 'Add',
  cari: 'Search ID, name, SOA name, or code',
  semuaType: 'All Types',
  semuaFlag: 'All Flags',
  memuat: 'Loading reinsurance types…',
  kosong: 'No reinsurance type yet.',
  tidakCocok: 'No reinsurance type matches the search.',
  jumlah: (n: number) => `${n} reinsurance type${n === 1 ? '' : 's'}`,

  no: 'No',
  id: 'ID',
  nama: 'Name',
  type: 'Type',
  soa: 'SOA Name',
  code: 'Code',
  flag: 'Flag',
  noUrut: 'No Urut',
  groupType: 'Group Type',
  pengubah: 'Last Edited By',
  tglUbah: 'Last Edited',
  aksi: 'Action',
  edit: 'Edit',

  judulTambah: 'Add Reinsurance Type',
  judulUbah: (id: string) => `Edit Reinsurance Type ${id}`,
  idOtomatis: 'Generated when saved',
  pilih: '— select —',
  tanpaGroup: '— none —',
  kosongNilai: '(empty)',
  nilaiLama: (s: string) => `${s} (old value)`,
  // Nama saja, tanpa kode depan (perintah work owner 05-10-2026: "kode depan nya hapus"); kode tak dikenal apa adanya.
  labelType: (t: string) => LABEL_TYPE[t] ?? t,
  /** Flag untuk dibaca: `active` -> Active, `inactive` -> Inactive; nilai warisan lain apa adanya. Yang tersimpan tetap huruf kecil. */
  labelFlag: (f: string) => (f === 'active' ? 'Active' : f === 'inactive' ? 'Inactive' : f),
  catatan: 'Name and SOA Name are saved in capital letters. To stop using a type, set Flag to inactive.',
  galatNama: 'Name is required',
  galatType: 'Type is required',
  galatFlag: 'Flag is required',
  galatCode: 'Code must be a number',
  galatNoUrut: 'No Urut must be a number',
  save: 'Save',
  menyimpan: 'Saving…',
  batal: 'Cancel',
  tutup: 'Close',
  tersimpan: (id: string) => `Reinsurance type ${id} saved.`,
} as const
