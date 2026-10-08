// Label modul R/I Rate Life - bahasa Inggris, VERBATIM dari section Pega `InboxSummaryRIRate` (RIRate.xml) bila ada:
// judul "R/I RATE SUMMARY" b382, "R/I RATE NAME" b1242, kolom "MODIFY OPERATOR" b9748 / "MODIFY DATE" b9894, tombol
// Save b1809 / Cancel b2109 / Upload CSV b2949 / View Upload b3467 / Simpan Upload b4511 / Edit b11109 / Detail b11398 /
// Delete b12238, popup "Rate Detail" b11446, label format b5305. Rate Detail: section `InboxRIRate` (View Detail.xml).

/** Nama menu - `M_NAV_MENU.LABEL` baris modul ini (migrasi inti 922). */
export const MENU_RR = { kelompok: 'R/I Rate Life' } as const

export const RR = {
  judul: 'R/I RATE SUMMARY',
  id: 'ID',
  nama: 'R/I RATE NAME',
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
  format: 'Format excel : USEDBY, CONTRACT, GENDER, AGE, RATE',
  judulDetail: 'Rate Detail',
  // Section `InboxRIRate` (View Detail.xml): judul b367, tombol Clear Field b1064, placeholder CONTRACT/AGE b2145 dan
  // RATE b2583.
  judulRincian: 'R/I RATE DETAIL',
  bersihkan: 'Clear Field',
  contohBulat: '0',
  contohRate: '0,0000',
  idRateOtomatis: 'Generated when saved',
  modeUbahRate: (id: string) => `Editing rate ID ${id}`,
  galatContract: 'CONTRACT is required',
  galatRate: 'RATE is required',
  rateTersimpan: (id: string) => `Rate ${id} saved.`,
  tutup: 'Close',

  idOtomatis: 'Generated when saved',
  modeUbah: (id: string) => `Editing ID ${id}`,
  saringID: 'Filter ID',
  saringNama: 'Filter R/I RATE NAME',
  memuat: 'Loading R/I rates…',
  kosong: 'No R/I rate summary yet.',
  tidakCocok: 'No R/I rate summary matches the filter.',
  jumlah: (n: number) => `${n} summar${n === 1 ? 'y' : 'ies'}`,
  halaman: (h: number, dari: number) => `Page ${h} of ${dari}`,
  sebelum: 'Previous',
  sesudah: 'Next',
  galatNama: 'R/I RATE NAME is required',
  tersimpan: (id: string) => `R/I rate ${id} saved.`,
  terhapus: (id: string, n: number) => `R/I rate ${id} deleted with ${n} rate row${n === 1 ? '' : 's'}.`,

  judulHapus: (id: string) => `Delete R/I rate ${id}`,
  tanyaHapus: (nama: string, n: number) =>
    `Delete "${nama}" and its ${n} rate row${n === 1 ? '' : 's'}? This cannot be undone.`,
  menghitung: 'Counting rate rows…',

  usedby: 'USEDBY',
  contract: 'CONTRACT',
  gender: 'GENDER',
  age: 'AGE',
  rate: 'RATE',
  barisKe: 'Row',
  pesan: 'Message',
  kosongDetail: 'No rate rows for this R/I rate.',

  pilihBerkas: 'CSV file',
  berkasDipilih: (nama: string, kb: number) => `${nama} (${kb} KB) is ready. Use View Upload to check it, then Simpan Upload.`,
  belumAdaBerkas: 'Choose a CSV file with Upload CSV first.',
  galatBerkas: 'Choose a .csv file of at most 4 MB.',
  catatanUnggah:
    'The first row is the header. Separator ; (a RATE may use a decimal comma) or , (a RATE with a decimal comma must be quoted, e.g. "0,5"). GENDER U, M, or F; AGE and CONTRACT 0-120 (CONTRACT may be empty).',
  ringkasanBaru: 'new',
  barisSah: (n: number) => `${n} valid row${n === 1 ? '' : 's'}`,
  barisGalat: (n: number) => `${n} row${n === 1 ? '' : 's'} with errors`,
  tujuan: 'Target R/I RATE NAME',
  siapSimpan: 'No errors. Simpan Upload will save every row in one transaction.',
  tidakSiap: 'Fix the errors in the file and upload it again; nothing is saved while any row has an error.',
  unggahTersimpan: (n: number, baru: number) =>
    `${n} rate row${n === 1 ? '' : 's'} saved${baru > 0 ? `, ${baru} new R/I rate summar${baru === 1 ? 'y' : 'ies'}` : ''}.`,
} as const
