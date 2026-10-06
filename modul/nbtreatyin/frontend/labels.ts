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
  /** showHarness tombol `Choose Business` (`DetailPolicyTreatyIn`) `pyWindowName` - judul jendela popup. */
  pilihBisnis: 'Business And SOB List',
  sumberBisnis: 'Source Of Business',
  /** FlowAction `PolicyTreatyInDeclineConfirm` pyLabel. */
  tolak: 'Confirm Decline NB',
  /** FlowAction `ShowPolicyNoTreaty` pyLabel. */
  nomorPolis: 'Show PolicyNo',
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
  /** pyNoSelectionText dropdown `.TreatyType` grid spreading. */
  pilihKosong: 'Choose',
} as const

/** Kolom daftar portal (`SFAPortal_OpportunitiesList`, LABEL judul grid `GetListOpportunity`).
 *  ⛔ RALAT W6 audit silang P3: kolom "Position" / "No Polis" (tanpa sel XML) DIBUANG.
 *  Kolom XML ke-2 "Name" (`.Name` pxLink `openWorkByHandle(.pzInsKey)`) tidak dirender:
 *  `.Name` ditulis NOL rule korpus dan tak berkolom - tautannya pindah ke "Offer No"
 *  (`.TextNoQuotation` = pengenal kasus yang sama; RALAT tiket 11). */
export const KOLOM_PORTAL = {
  id: 'Offer No',
  bisnis: 'Group Business',
  tertanggung: 'Insured Name',
  marketing: 'Marketing',
  status: 'Status',
} as const

export const PORTAL = {
  /** LABEL `SFAPortalOpportunitiesHeader` C[1.1]. */
  judul: 'Opportunity',
  /** `.FilterTermForOpportunity` pyLabelFor (nama aksesibel) dan pyPlaceholder. */
  filter: 'Filter Term for Opportunity',
  placeholder: 'NB-1234 or Name',
  /** Nama aksesibel ikon pengosong C[1.2] (pxIcon `pyImage webwb/pyiconclearfield.png`,
   *  `pyIncludeLabel=false` - tanpa teks tampil). */
  bersihkan: 'Clear field',
  /** Teks daftar kosong - komponen `Kosong` inti (tiket 11: dasar unsur bawaan inti). */
  kosong: 'Belum ada berkas realisasi treaty yang terbuka.',
  kosongPetunjuk: 'Tombol Create membuat berkas baru di antrean admin.',
  hanyaBaca: 'Berkas ini menunggu di antrean lain atau sudah selesai - hanya-baca.',
} as const

/** Judul kolom TreeGrid `Section/SourceHierarki` (pyCaption, `.ClientName`). */
export const KOLOM_SOB = 'Source of Business Name'

/** Kolom grid AKTIF popup `Section/BusinessAndSOBList` (RD `BrowseTreatyJoinEDM`), urutan sel:
 *  `kolom` = properti sel baris 2 (kolom view), `judul` = LABEL sel baris 1 di atasnya, VERBATIM -
 *  termasuk dua judul KOSONG (`.LAYER`, `.LAYERPART`) dan "Insured Name" untuk `.CEDING`.
 *  ⛔ RALAT audit silang P3 W1: bunyi lama "nama kolom view, apa adanya" (judul = nama kolom) keliru. */
export const KOLOM_BISNIS = [
  { kolom: 'TREATYID', judul: 'Treaty Offer ID' },
  { kolom: 'TREATYCONTRACTNAME', judul: 'Contract Name' },
  { kolom: 'CLASSOFBUSINESS', judul: 'Class of Business' },
  { kolom: 'SOB', judul: 'Source Of Business' },
  { kolom: 'CEDING', judul: 'Insured Name' },
  { kolom: 'PROPORTIONTYPE', judul: 'Proportion Type' },
  { kolom: 'TREATYTYPE', judul: 'Treaty Type' },
  { kolom: 'TREATYGROUP', judul: 'Treaty Group' },
  { kolom: 'TREATYYEAR', judul: 'Treaty Year' },
  { kolom: 'LIMITCURRENCY', judul: 'Currency' },
  { kolom: 'LIMITVALUE', judul: 'Limit' },
  { kolom: 'RETENTIONCURRENCY', judul: 'Currency' },
  { kolom: 'RETENTIONVALUE', judul: 'Retention' },
  { kolom: 'EPICURRENCY', judul: 'Currency' },
  { kolom: 'EPIVALUE', judul: 'EPI' },
  { kolom: 'LAYERTYPE', judul: 'Layer' },
  { kolom: 'LAYER', judul: '' },
  { kolom: 'LAYERPARTTYPE', judul: 'Part of' },
  { kolom: 'LAYERPART', judul: '' },
  { kolom: 'MDPCURRENCY', judul: 'Currency' },
  { kolom: 'MDPVALUE', judul: 'MDP' },
  { kolom: 'NETPREMICURRENCY', judul: 'Currency' },
  { kolom: 'NETPREMIVALUE', judul: 'Net Premium' },
  { kolom: 'SHARECURRENCY', judul: 'Currency' },
  { kolom: 'SHAREVALUE', judul: 'Share RNM Value' },
] as const

/** Judul di dalam layar realisasi yang TAMPIL di XML. ⛔ RALAT W6 audit silang P3: judul
 *  buatan "General", "Premium & Claim", "Spreading Risk", "Suggest" DIBUANG - wadahnya
 *  `NOHEADER` / `pyIncludeHeader=false` (S2 berjudul "General" pun tanpa header). */
export const BAGIAN = {
  /** wadah S45 / S107 grid `.ListInstallment` (pyIncludeHeader=true, pyTitle). */
  angsuran: 'Installment Data Information',
  /** LABEL sel `pyIncludeLabel=true` (Heading 4) wadah S24/S93 dan S25/S94. */
  ogp: 'OGP',
  onp: 'ONP',
} as const

/** Kolom grid SpreadingRiskList - LABEL judul dan kaki grid (`SpreadingRiskList` S2,
 *  sama di DetailPolicyTreatyIn dan DetailDeptHeadTreatyIn_UW). */
export const KOLOM_SPREADING = {
  treatyType: 'Type Treaty',
  share: '% Share',
  premium: 'Premium',
  claimPct: '% Share',
  claim: 'Claim',
  totalShare: 'Total',
  totalPremium: 'Total Premium',
  totalShareClaim: 'Total %Share Claim',
  totalClaim: 'Total Claim',
} as const

/** Kolom grid ListInstallment - LABEL judul grid S45. */
export const KOLOM_ANGSURAN = {
  no: 'No',
  dueDate: 'Due Date',
  pct: '% Installment',
  premium: 'Premium Nusantara Re',
  total: 'Total Payment',
  installment: 'Installment',
} as const

/** Kolom grid SuggestList (ListSuggest). */
export const KOLOM_USULAN = {
  tanggal: 'Date',
  operator: 'PIC',
  putusan: 'Approval',
  catatan: 'Suggest',
} as const

/** LABEL (Heading 2) `Section/PolicyTreatyInDeclineConfirm` - VERBATIM. */
export const KONFIRMASI_TOLAK = 'Are you sure you want to decline this NB'

/** LABEL `Section/ShowPolicyNoTreaty_SC` di antara pyID dan PolicyNo - VERBATIM. */
export const NOMOR_DIAKSEP = 'telah diaksep menjadi'

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
