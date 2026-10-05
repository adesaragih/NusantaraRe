// Label modul Bordereaux - bahasa Inggris. Caption mengikuti section Pega folder korpus `Bordereaux`:
// `PortalBordereaux_Sec` + `PortalShowHeader` (daftar, filter), `InputBordereaux` (form, tab Details/Summary/Submit,
// History), `ShowMasterTreatyIn_Sec` (popup Choose Master Treaty). Pesan penolakan datang dari backend.

/** Nama menu - `M_NAV_MENU.LABEL` baris modul ini (migrasi inti 913, keputusan work owner 04-10-2026). */
export const MENU_BDX = { kelompok: 'Bordereaux' } as const

export const BDX = {
  judul: 'Bordereaux',
  inputData: 'Input Data',
  searchFilter: 'Search / Filter',
  applyFilter: 'Apply Filter',
  reset: 'Reset',
  memuat: 'Loading bordereaux…',
  kosong: 'No bordereaux yet.',
  tidakCocok: 'No bordereaux matches the filter.',
  jumlah: (n: number) => `${n} bordereaux`,
  halaman: (h: number, t: number) => `Page ${h} of ${t}`,
  sebelumnya: '‹ Previous',
  berikutnya: 'Next ›',

  // Chart daftar (dirancang sendiri - PieChartBordereaux tidak ada di ekspor; work owner 04-10-2026).
  chart: 'Bordereaux by Business',
  semuaBusiness: 'All Business',
  jejak: 'Chart level',
  totalBerkas: 'Total bordereaux',
  jumlahBusiness: 'Number of Business',
  jumlahCeding: 'Number of Cedings',
  petunjukBuka: 'Click a business to see its cedings',
  memuatChart: 'Loading chart…',
  chartKosong: 'No bordereaux yet.',
  berkas: (n: number) => `${n} bordereaux`,

  // Filter dan kolom daftar (PortalBordereaux_Sec).
  id: 'ID',
  type: 'Type',
  business: 'Business',
  filterBusiness: 'BUSINESS',
  reffNoSoa: 'Reff No SOA',
  filterReffNoSoa: 'Reff No Soa',
  reffNoBdx: 'Reff No BDX',
  cedingCo: 'Ceding Co',
  treatyName: 'Treaty Name',
  reportStart: 'Bordereaux Report Start',
  reportEnd: 'Bordereaux Report End',
  position: 'Position',
  status: 'Status',
  semua: 'All',
  aksi: 'Action',
  edit: 'Edit',
  view: 'View',
  hapus: 'Delete',

  // Form (InputBordereaux).
  typeBusiness: 'Type Business',
  chooseMasterTreaty: 'Choose Master Treaty',
  masterTreaty: 'Master Treaty',
  masterKosong: 'No Master Treaty chosen yet.',
  subrogasiBonding: 'BONDING (Subrogation is always Bonding)',
  masterId: 'Master ID',
  sobName: 'SOB Name',
  bdxStartEnd: 'Bordereaux Start / End',
  reffNoOfSoa: 'Reff No of SOA',
  reffNoOfBdx: 'Reff No of Bordereaux',
  tabDetails: 'Details',
  tabSummary: 'Summary',
  tabSubmit: 'Submit',
  unggahPetunjuk: 'Please upload file with .csv format',
  template: 'Template',
  uploadCsv: 'Upload CSV',
  membaca: 'Reading CSV…',
  detailKosong: 'No detail row yet. Choose Type and Business, then Upload CSV.',
  pilihTypeDulu: 'Choose Type and Type Business first.',
  baris: (n: number) => `${n} row${n === 1 ? '' : 's'}`,
  /** Kolom nomor kepala format Excel bordereaux 2025. */
  no: 'No',
  close: 'Close',
  save: 'Save',
  menyimpan: 'Saving…',
  tersimpan: (id: string) => `Bordereaux ${id} saved.`,
  memuatBerkas: 'Loading bordereaux…',

  // Summary (CountSummaryBdx).
  currency: 'Currency',
  ringkasan: {
    PREMIUM: ['Total Premium Reinsurer', 'Total Premium NusantaraRe'],
    CLAIM: ['Total_Claims_Reinsurer', 'Total_Claims_RNM'],
    SUBROGATION: ['Total_Claims_Subrogation_Reinsurer', 'Total_Claims_Subrogation_RNM'],
  } as Record<string, [string, string]>,
  ringkasanKosong: 'No summary yet.',

  // Submit (ActionSubmit, AkseptasiBdx_DT).
  statusPersetujuan: 'Status',
  approve: 'Approve',
  reject: 'Reject',
  comment: 'Comment',
  submit: 'Submit',
  mengirim: 'Submitting…',
  terkirim: 'Submitted.',
  submitPembuat: 'Send this bordereaux to the Checker.',
  submitTidakMenunggu: 'This bordereaux is not waiting for you.',
  submitSelesai: 'This bordereaux is Resolve-Complete.',
  history: 'History',
  date: 'Date',
  pic: 'PIC',
  approval: 'Approval',
  historyKosong: 'No history yet.',

  // Popup Choose Master Treaty (ShowMasterTreatyIn_Sec).
  cedant: 'Cedant',
  cariCedant: 'Type a cedant name',
  contractName: 'Contract Name',
  reinsuranceType: 'Reinsurance Type',
  sourceOfBusiness: 'Source of Business',
  ceding: 'Ceding',
  choose: 'Choose',
  mencari: 'Searching…',
  treatyKosong: 'No Master Treaty for this cedant.',
  batal: 'Cancel',

  // Hapus (konfirmasi dirancang sendiri - XML ConfrimDelete tidak ada; keputusan work owner).
  judulHapus: 'Delete bordereaux',
  kalimatHapus: (id: string) =>
    `Delete bordereaux ${id}? Its detail rows, attachments, and history are deleted too. This cannot be undone.`,
  menghapus: 'Deleting…',
  terhapus: (id: string) => `Bordereaux ${id} deleted.`,

  // Copy Old Data (keputusan work owner 04-10-2026; hanya superadmin; bentuk sama dengan Copy Old Company Detail).
  copyOld: 'Copy Old Data',
  copyOldKeterangan:
    'Bordereaux from the old application whose header, detail rows, or approval history are still only in the old JSON. ' +
    'Checked rows are copied one bordereaux at a time. Data already in the tables is never overwritten, and the old JSON is only read.',
  copyOldKosong: 'All old bordereaux are already in the tables.',
  cariLama: 'Search ID, ceding, or treaty',
  pilihSemua: 'Select all',
  pilihBaris: 'Select',
  dipilih: 'selected',
  akanDisalin: 'Will copy',
  lamaHeader: 'header',
  lamaDetail: (n: number) => `${n} detail row${n === 1 ? '' : 's'}`,
  lamaRiwayat: (n: number) => `${n} history row${n === 1 ? '' : 's'}`,
  lamaTakTerbaca: 'old JSON to check',
  prosesCopy: 'Process Copy',
  menyalin: 'Copying…',
  menyalinKemajuan: (selesai: number, total: number) => `Copying… ${selesai}/${total}`,
  statusDisalin: 'copied',
  statusSudahAda: 'already in the tables',
  statusDitolak: 'rejected',
  statusGagal: 'failed',
} as const

/** Nama status untuk dibaca (STATUSAKSEP kosong = draf). */
export function teksStatus(s: string): string {
  return s === '' ? 'Draft' : s
}
