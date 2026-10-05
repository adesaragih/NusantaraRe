// Label modul Treaty Exchange Yearly - bahasa Inggris. Caption kolom mengikuti `BrowseTreatyExchangeYearly_RD` Pega
// (TreatyYear, Currency, StartDate, EndDate, ToIDR, ToUSD, Quarter, UserID, DateIU).

/** Nama menu - `M_NAV_MENU.LABEL` baris modul ini (migrasi inti 919). */
export const MENU_TEY = { kelompok: 'Treaty Exchange Yearly' } as const

export const TEY = {
  judul: 'Treaty Exchange Yearly',
  tambah: 'Add',
  cari: 'Search ID or currency',
  semuaTahun: 'All Treaty Years',
  memuat: 'Loading exchange rates…',
  kosong: 'No exchange rate yet.',
  tidakCocok: 'No exchange rate matches the search.',
  jumlah: (n: number) => `${n} exchange rate${n === 1 ? '' : 's'}`,

  no: 'No',
  id: 'ID',
  tahun: 'Treaty Year',
  mataUang: 'Currency',
  quarter: 'Quarter',
  mulai: 'Start Date',
  akhir: 'End Date',
  toIdr: 'To IDR',
  toUsd: 'To USD',
  pengubah: 'Last Edited By',
  tglUbah: 'Last Edited',
  aksi: 'Action',
  edit: 'Edit',

  judulTambah: 'Add Exchange Rate',
  judulUbah: (id: string) => `Edit Exchange Rate ${id}`,
  idOtomatis: 'Generated when saved',
  pilih: '— select —',
  nilaiLama: (s: string) => `${s} (old value)`,
  labelQuarter: (q: string) => (q === '0' ? '0 - Yearly' : `${q} - Q${q}`),
  catatan:
    'Rates use a dot for decimals and no thousand separators. Start and End Date follow the treaty year (1 July - 30 June) and can be changed.',
  galatTahun: 'Treaty Year must be 4 digits',
  galatMataUang: 'Currency is required',
  galatTanggal: 'Start Date and End Date are required',
  galatUrutan: 'End Date must not be before Start Date',
  galatIdr: 'To IDR must be a number (dot for decimals, no thousand separators)',
  galatUsd: 'To USD must be a number (dot for decimals, no thousand separators)',
  save: 'Save',
  menyimpan: 'Saving…',
  batal: 'Cancel',
  tutup: 'Close',
  tersimpan: (id: string) => `Exchange rate ${id} saved.`,
} as const
