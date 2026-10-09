// Label modul R/I Risk - bahasa Inggris, VERBATIM dari section Pega bila ada.
// `InboxSummaryRIRisk`: judul "R/I RISK SUMMARY" b375, form "ID" b1042 / "R/I RISK NAME" b1225, kolom "R/I RISK NAME"
// b7147 / "MODIFY OPERATOR" b7295 / "MODIFY DATE" b7396, tombol Upload CSV b3156 / View Upload b3676 / Simpan Upload
// b4690 / Edit b8589 / Detail b8869 / Delete b9624, Save b1796 / Cancel b2063 (di XML wadahnya `1=2` b1511 - DITAMPILKAN
// atas keputusan work owner 08-10-2026). Label format b5495 tersembunyi (`1=2` b5263) - tidak dirender.
// `InboxRIRisk`: judul "R/I RISK DETAIL" b382, Clear Field b1125, medan ID b1316 / R/I RISK NAME b1696 / CONTRACT
// b1889 / YEAR b2171 / MONTH b2381 (pyLabelPreview; pyLabelFieldValue b2358 bertuliskan "YEAR" - salin-tempel Pega,
// [penyimpangan sadar - menunggu WO]) / RISK (PERMIL) b2548, placeholder 0 b1936/b2219/b2406, Save b3132 / Cancel
// b3399, grid ID b7635 / R/I RISK NAME b7789 / CONTRACT b7956 / YEAR b8123 / MONTH b8277 / RISK (PERMIL) b8423, EDIT
// b9784.

/** Nama menu - `M_NAV_MENU.LABEL` baris modul ini (migrasi inti 941, keputusan work owner 08-10-2026 K4). */
export const MENU_RK = { kelompok: 'R/I Risk' } as const

export const RK = {
  judul: 'R/I RISK SUMMARY',
  id: 'ID',
  nama: 'R/I RISK NAME',
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
  /** Kolom berkas CSV (urutan label format b5495). */
  kolomUsedby: 'USEDBY',
  unggah: 'Upload CSV',
  lihatUnggah: 'View Upload',
  simpanUnggah: 'Simpan Upload',
  /** Hidden in Pega (`1=2` b5263) - kept for the CSV header check only, not rendered. */
  format: 'Format excel : USEDBY, CONTRACT, YEAR, MONTH, RISK',

  judulRincian: 'R/I RISK DETAIL',
  bersihkan: 'Clear Field',
  editRincian: 'EDIT',
  contract: 'CONTRACT',
  year: 'YEAR',
  month: 'MONTH',
  risk: 'RISK (PERMIL)',
  contohBulat: '0',
  idDetailOtomatis: 'Generated when saved',
  modeUbahDetail: (id: string) => `Editing detail ID ${id}`,
  galatContract: 'CONTRACT is required',
  galatRisk: 'RISK is required',
  galatAngka: (medan: string, digit: number) => `${medan} must be a whole number of at most ${digit} digits`,
  galatRiskNegatif: 'RISK must not be negative',
  detailTersimpan: (id: string) => `R/I risk detail ${id} saved.`,
  tutup: 'Close',
  kosongDetail: 'No R/I risk detail rows for this R/I risk.',

  idOtomatis: 'Generated when saved',
  modeUbah: (id: string) => `Editing ID ${id}`,
  galatNama: 'R/I RISK NAME is required',
  tersimpan: (id: string) => `R/I risk ${id} saved.`,
  saringID: 'Filter ID',
  saringNama: 'Filter R/I RISK NAME',
  memuat: 'Loading R/I risks…',
  kosong: 'No R/I risk summary yet.',
  tidakCocok: 'No R/I risk summary matches the filter.',
  jumlah: (n: number) => `${n} summar${n === 1 ? 'y' : 'ies'}`,
  halaman: (h: number, dari: number) => `Page ${h} of ${dari}`,
  sebelum: 'Previous',
  sesudah: 'Next',
  terhapus: (id: string, n: number) => `R/I risk ${id} deleted with ${n} detail row${n === 1 ? '' : 's'}.`,

  judulHapus: (id: string) => `Delete R/I risk ${id}`,
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
    'The first row is the header USEDBY, CONTRACT, YEAR, MONTH, RISK. Separator ; (a RISK may use a decimal comma) or , (a RISK with a decimal comma must be quoted, e.g. "0,5"). CONTRACT a whole number, YEAR and MONTH empty or a whole number of at most 4 digits, RISK a number that is not negative.',
  ringkasanBaru: 'new',
  barisSah: (n: number) => `${n} valid row${n === 1 ? '' : 's'}`,
  barisGalat: (n: number) => `${n} row${n === 1 ? '' : 's'} with errors`,
  tujuan: 'Target R/I RISK NAME',
  siapSimpan: 'No errors. Simpan Upload will save every row in one transaction.',
  tidakSiap: 'Fix the errors in the file and upload it again; nothing is saved while any row has an error.',
  unggahTersimpan: (n: number, baru: number) =>
    `${n} detail row${n === 1 ? '' : 's'} saved${baru > 0 ? `, ${baru} new R/I risk summar${baru === 1 ? 'y' : 'ies'}` : ''}.`,
} as const
