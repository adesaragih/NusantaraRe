// Label modul R/I Comm Life - bahasa Inggris, VERBATIM dari section Pega `InboxSummaryRIComm` bila ada: judul
// "R/I COMM SUMMARY" b382, form "USEDBY" b1243, kolom "R/I COMM NAME" b7220 / "MODIFY OPERATOR" b7377 / "MODIFY DATE"
// b7521, tombol Save b1799 / Cancel b2092 / Upload CSV b3189 / View Upload b3706 / Simpan Upload b4771 / Edit b8754 /
// Detail b9042 / Delete b9680, label format b5569. Detail: section `InboxRIComm` - judul "R/I COMM DETAIL" b349, Clear
// Field b1100, CONTRACT b1873 / YEAR b2154 / COMM b2359, placeholder 0 b1927/b2205.

/** Nama menu - `M_NAV_MENU.LABEL` baris modul ini (migrasi inti 925). */
export const MENU_RC = { kelompok: 'R/I Comm Life' } as const

export const RC = {
  judul: 'R/I COMM SUMMARY',
  id: 'ID',
  usedby: 'USEDBY',
  nama: 'R/I COMM NAME',
  operator: 'MODIFY OPERATOR',
  tanggal: 'MODIFY DATE',
  save: 'Save',
  menyimpan: 'Saving…',
  batal: 'Cancel',
  edit: 'Edit',
  detail: 'Detail',
  hapus: 'Delete',
  menghapus: 'Deleting…',
  aksi: 'Action',
  unggah: 'Upload CSV',
  lihatUnggah: 'View Upload',
  simpanUnggah: 'Simpan Upload',
  format: 'Format excel : USEDBY, CONTRACT, YEAR, COMM',

  judulRincian: 'R/I COMM DETAIL',
  bersihkan: 'Clear Field',
  contract: 'CONTRACT',
  year: 'YEAR',
  comm: 'COMM',
  contohBulat: '0',
  idDetailOtomatis: 'Generated when saved',
  modeUbahDetail: (id: string) => `Editing detail ID ${id}`,
  galatContract: 'CONTRACT is required',
  galatYear: 'YEAR is required',
  galatComm: 'COMM is required',
  detailTersimpan: (id: string) => `R/I comm detail ${id} saved.`,
  tutup: 'Close',
  kosongDetail: 'No R/I comm detail rows for this R/I comm.',

  idOtomatis: 'Generated when saved',
  modeUbah: (id: string) => `Editing ID ${id}`,
  saringID: 'Filter ID',
  saringNama: 'Filter R/I COMM NAME',
  memuat: 'Loading R/I comms…',
  kosong: 'No R/I comm summary yet.',
  tidakCocok: 'No R/I comm summary matches the filter.',
  jumlah: (n: number) => `${n} summar${n === 1 ? 'y' : 'ies'}`,
  halaman: (h: number, dari: number) => `Page ${h} of ${dari}`,
  sebelum: 'Previous',
  sesudah: 'Next',
  galatNama: 'USEDBY is required',
  tersimpan: (id: string) => `R/I comm ${id} saved.`,
  terhapus: (id: string, n: number) => `R/I comm ${id} deleted with ${n} detail row${n === 1 ? '' : 's'}.`,

  judulHapus: (id: string) => `Delete R/I comm ${id}`,
  tanyaHapus: (nama: string, n: number) =>
    `Delete "${nama}" and its ${n} detail row${n === 1 ? '' : 's'}? This cannot be undone.`,
  menghitung: 'Counting detail rows…',

  barisKe: 'Row',
  pesan: 'Message',

  pilihBerkas: 'CSV file',
  berkasDipilih: (nama: string, kb: number) => `${nama} (${kb} KB) is ready. Use View Upload to check it, then Simpan Upload.`,
  belumAdaBerkas: 'Choose a CSV file with Upload CSV first.',
  galatBerkas: 'Choose a .csv file of at most 4 MB.',
  catatanUnggah:
    'The first row is the header. Separator ; (a COMM may use a decimal comma) or , (a COMM with a decimal comma must be quoted, e.g. "0,5"). CONTRACT a whole number, YEAR a whole number of at most 4 digits, COMM a number with at most 8 decimals.',
  ringkasanBaru: 'new',
  barisSah: (n: number) => `${n} valid row${n === 1 ? '' : 's'}`,
  barisGalat: (n: number) => `${n} row${n === 1 ? '' : 's'} with errors`,
  tujuan: 'Target R/I COMM NAME',
  siapSimpan: 'No errors. Simpan Upload will save every row in one transaction.',
  tidakSiap: 'Fix the errors in the file and upload it again; nothing is saved while any row has an error.',
  unggahTersimpan: (n: number, baru: number) =>
    `${n} detail row${n === 1 ? '' : 's'} saved${baru > 0 ? `, ${baru} new R/I comm summar${baru === 1 ? 'y' : 'ies'}` : ''}.`,
} as const
