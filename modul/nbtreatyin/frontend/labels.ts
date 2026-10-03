// Label layar modul NB Treaty In.
//
// ⛔ Label medan, tombol, dan judul kolom VERBATIM dari section Pega
// (INVENTARIS-XML.md bab 5: `Section/DetailPolicyTreatyIn`,
// `DetailDeptHeadTreatyIn_UW`, `ListSuggest`, `SpreadingRiskList`,
// `BusinessAndSOBList`, `SFAPortal_OpportunitiesList`,
// `PolicyTreatyInDeclineConfirm`, `ShowPolicyNoTreaty_SC`). Label generik Pega
// ("Text Input", "Formatted Text", "Checkbox") diganti nama propertinya - label
// itu bawaan kontrol, bukan teks yang pernah dibaca pengguna sebagai nama medan.

export const JUDUL = {
  portal: 'NB Treaty In',
  admin: 'Input Realitation',
  secHead: 'Acceptance by Head. Treaty',
  deptHead: 'Acceptance by Dept. Head',
  pilihBisnis: 'Choose Business',
  sumberBisnis: 'Source Of Business',
  tolak: 'Decline',
  nomorPolis: 'Policy No',
  riwayat: 'History',
  catatan: 'Suggest',
} as const

/** Judul layar per posisi kasus (nama assignment Flow). */
export const JUDUL_POSISI: Record<string, string> = {
  ReasTreatyInAdmin: JUDUL.admin,
  ReasTreatyInSecHead: JUDUL.secHead,
  ReasTreatyInDeptHead: JUDUL.deptHead,
}

export const TOMBOL = {
  create: 'Create opportunity',
  filter: 'Filter',
  chooseBusiness: 'Choose Business',
  selectSOB: 'Select Source Of Business',
  choose: 'Choose',
  save: 'Save',
  submit: 'Submit',
  yes: 'Yes',
  no: 'No',
  ok: 'OK',
  add: 'Add',
  delete: 'Delete',
  enableDisable: 'Enable / Disable Input Type',
  kembali: 'Kembali',
} as const

/** Kolom daftar portal (`SFAPortal_OpportunitiesList`). */
export const KOLOM_PORTAL = {
  id: 'No',
  bisnis: 'Business',
  tertanggung: 'Insured',
  marketing: 'Marketing',
  status: 'NB Status',
  posisi: 'Position',
  nopol: 'No Polis',
} as const

export const PORTAL = {
  filter: 'Filter Term for Opportunity',
  kosong: 'Belum ada berkas realisasi treaty yang terbuka.',
  kosongPetunjuk: 'Tombol Create membuat berkas baru di antrean admin.',
  hanyaBaca: 'Berkas ini menunggu di antrean lain atau sudah selesai - hanya-baca.',
} as const

/** Judul kolom TreeGrid `Section/SourceHierarki` (pyCaption, `.ClientName`). */
export const KOLOM_SOB = 'Source of Business Name'

/** Kolom grid popup `BusinessAndSOBList` (nama kolom view, apa adanya). */
export const KOLOM_BISNIS = [
  'TREATYID',
  'TREATYCONTRACTNAME',
  'CLASSOFBUSINESS',
  'SOB',
  'CEDING',
  'PROPORTIONTYPE',
  'TREATYTYPE',
  'TREATYGROUP',
  'TREATYYEAR',
  'LIMITCURRENCY',
  'LIMITVALUE',
  'RETENTIONCURRENCY',
  'RETENTIONVALUE',
  'EPICURRENCY',
  'EPIVALUE',
  'LAYERTYPE',
  'LAYER',
  'LAYERPARTTYPE',
  'LAYERPART',
  'MDPCURRENCY',
  'MDPVALUE',
  'NETPREMICURRENCY',
  'NETPREMIVALUE',
  'SHARECURRENCY',
  'SHAREVALUE',
] as const

export const BAGIAN = {
  umum: 'General',
  uang: 'Premium & Claim',
  spreading: 'Spreading Risk',
  angsuran: 'Installment',
  usulan: 'Suggest',
} as const

/** Kolom grid SpreadingRiskList. */
export const KOLOM_SPREADING = {
  treatyType: 'Treaty Type',
  share: '%Share',
  premium: 'Premium Spreaded',
  claimPct: '%Share Claim',
  claim: 'Claim Spreaded',
  totalShare: 'Total %Share',
  totalPremium: 'Total Premium',
  totalShareClaim: 'Total %Share Claim',
  totalClaim: 'Total Claim',
} as const

/** Kolom grid ListInstallment. */
export const KOLOM_ANGSURAN = {
  no: 'Installment No',
  dueDate: 'Due Date',
  pct: 'Installment %',
  premium: 'Premium',
  total: 'Payment Total',
  installment: 'Installment',
} as const

/** Kolom grid SuggestList (ListSuggest). */
export const KOLOM_USULAN = {
  tanggal: 'Date',
  operator: 'Operator',
  putusan: 'Approval',
  catatan: 'Suggest',
} as const

/** `PolicyTreatyInDeclineConfirm` - hanya dua tombol di XML; kalimatnya
 *  milik label rule yang tidak terekspor. */
export const KONFIRMASI_TOLAK = 'Decline?'

/**
 * Label pilihan Approval. Nilai "1"/"0" VERBATIM `DecisionTable/isApproved`;
 * teksnya VERBATIM `SaveViewSuggest` CARI9 (`1` -> "Accept", `0` -> "Reject").
 * Pilihan "associated values" radio aslinya tinggal di rule Property yang tidak
 * terekspor.
 */
export const PILIHAN_APPROVAL = [
  { value: '1', label: 'Accept' },
  { value: '0', label: 'Reject' },
] as const

export const PESAN = {
  tersimpan: 'Tersimpan.',
  terkirim: 'Berkas dikirim.',
  opsiTerbuka:
    'Pilihan medan ini didefinisikan di rule Property yang tidak ada di korpus; nilai diketik apa adanya.',
} as const
