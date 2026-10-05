// Label modul Accounts - bahasa Inggris. Label form VERBATIM tangkapan layar Pega SFAGIS Account kiriman work owner
// 04-10-2026 (TIDAK ada di korpus XML): Insured Name, Org ID, Group Business, Owner, Description, dan pesan
// "Value cannot be blank". "Territory" dihapus; "Create Date" ditambah; View dan Edit dicabut - akun hanya diinput
// sekali (keputusan work owner 04-10-2026).

/** Nama menu - `M_NAV_MENU.LABEL` baris modul ini (migrasi inti 908, keputusan work owner 04-10-2026). */
export const MENU_ACC = { kelompok: 'Accounts' } as const

export const ACC = {
  judul: 'Accounts',
  tambah: 'Add',
  cari: 'Search ACC-n, Insured Name, Org ID, or Group Business',
  kosong: 'No account yet.',
  tidakCocok: 'No account matches the search.',
  memuat: 'Loading accounts…',

  kolomId: 'ID',
  insuredName: 'Insured Name',
  orgId: 'Org ID',
  groupBusiness: 'Group Business',
  owner: 'Owner',
  createDate: 'Create Date',
  description: 'Description',

  simpan: 'Save',
  menyimpan: 'Saving…',
  batal: 'Cancel',
  judulBaru: 'Add account',
  ownerOtomatis: 'Your user account, set when saved',
  orgOtomatis: 'Filled from Insured Name',
  /** VERBATIM tangkapan layar Pega. */
  wajib: 'Value cannot be blank',
  pilih: '— select —',
  tersimpan: (id: string) => `Account ${id} saved.`,
} as const
