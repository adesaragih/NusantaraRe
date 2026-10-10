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
  /** showHarness tombol `Survey Report` `pyWindowName`. */
  surveiHistoris: 'Historical Survey Report',
  sumberBisnis: 'Source Of Business',
  /** FlowAction `PolicyTreatyInDeclineConfirm` pyLabel. */
  tolak: 'Confirm Decline NB',
  /** FlowAction `ShowPolicyNoTreaty` pyLabel. */
  nomorPolis: 'Show PolicyNo',
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
  /** Tombol sel 22 `Section/DetailPolicyTreatyIn` (LABEL). */
  surveyReport: 'Survey Report',
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
  // keputusan work owner 06-10-2026 - bukan kolom GetListOpportunity
  jenis: 'Type',
  // perintah work owner 10-10-2026 "ganti yang di inbox group bisnis jadi treaty bisnis" (dulu "Group Business")
  bisnis: 'Treaty Business',
  tertanggung: 'Insured Name',
  marketing: 'Marketing',
  status: 'Status',
  // keputusan work owner 06-10-2026 - bukan kolom GetListOpportunity
  pembuat: 'User Create',
  tanggal: 'Tanggal Create',
} as const

/** Kolom tambahan tab Resolved portal (perintah work owner 07-10-2026: "KALO DAH RESOLVE TAMBAHIN KOLOM NOPOLISNYA" dan "SEKALIAN KELUARIN TANGGAL PRODUKSINYA AJA DD-MM-YYYY"); label sama dengan portal EDM Treaty In. */
export const KOLOM_PORTAL_SELESAI = {
  noPolis: 'Policy Number',
  tglProduksi: 'Production Date',
} as const

/** Kolom daftar kotak masuk Beranda (keputusan work owner 06-10-2026: portal + kolom yang ada datanya). */
export const KOLOM_BERANDA = {
  id: 'Offer No',
  jenis: 'Type',
  tertanggung: 'Insured Name',
  bisnis: 'Treaty Business', // perintah work owner 10-10-2026 (dulu "Group Business")
  ceding: 'Ceding Company',
  mulai: 'Inception Date',
  marketing: 'Marketing',
  status: 'Status Inbox',
  pembuat: 'User Create',
  tanggal: 'Tanggal Create',
  sejak: 'Time Since Last Update',
} as const

export const PORTAL = {
  /** LABEL `SFAPortalOpportunitiesHeader` C[1.1]. */
  judul: 'Opportunity',
  /** `.FilterTermForOpportunity` pyLabelFor (nama aksesibel) dan pyPlaceholder. */
  filter: 'Filter Term for Opportunity',
  // menyimpang dari pyPlaceholder XML ("NB-1234 or Name"): menyebut kolom yang dicari (perintah work owner 07-10-2026)
  placeholder: 'Search NB no, master ID, policy no, insured, business, ceding, marketing...',
  /** Nama aksesibel ikon pengosong C[1.2] (pxIcon `pyImage webwb/pyiconclearfield.png`,
   *  `pyIncludeLabel=false` - tanpa teks tampil). */
  bersihkan: 'Clear field',
  /** Teks daftar kosong - komponen `Kosong` inti (tiket 11: dasar unsur bawaan inti). */
  kosong: 'Belum ada berkas realisasi treaty yang terbuka.',
  kosongPetunjuk: 'Tombol Create membuat berkas baru di antrean admin.',
  /** Daftar kosong di posisi switch Resolved. */
  kosongSelesai: 'Belum ada berkas yang selesai.',
  hanyaBaca: 'Berkas ini menunggu di antrean lain atau sudah selesai - hanya-baca.',
} as const

/** Switch portal di atas daftar (keputusan work owner 06-10-2026: "switch untuk lihat yang lagi proses atau resolve,
 *  default ke proses"). Resolved = Resolved-Completed / Resolved-Rejected. */
export const STATUS_PORTAL = ['In Progress', 'Resolved'] as const
export type StatusPortal = (typeof STATUS_PORTAL)[number]

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

/** Prompt values property `.ClaimType` (pxDropdown pyListSource `associated`; screenshot work owner 05-10-2026):
 *  nilai standar disimpan, teks prompt ditampilkan. */
export const PILIHAN_CLAIM_TYPE: { value: string; label: string }[] = [
  { value: 'SOA', label: 'SOA' },
  { value: 'CashLoss', label: 'Cash Loss' },
  { value: 'XOL', label: 'XOL' },
  { value: 'XOL Retro', label: 'XOL Retro' },
]

/** Prompt values property `.StatementType` (pxDropdown pyListSource `associated`; screenshot work owner 06-10-2026). */
export const PILIHAN_STATEMENT_TYPE: { value: string; label: string }[] = [
  { value: 'SOA', label: 'Statement of Account' },
  { value: 'LPC', label: 'Loss Participation Clause' },
  { value: 'PC', label: 'Profit Commission' },
  { value: 'SC', label: 'Sliding Scale' },
]

/** Prompt values property `.TypeTax` (pxRadioButtons pyListSource `associated`; screenshot work owner 06-10-2026). */
export const PILIHAN_TYPE_TAX: { value: string; label: string }[] = [
  { value: 'Inclusive', label: 'Inclusive' },
  { value: 'Exclusive', label: 'Exclusive' },
]

/** Prompt values property `.QuotationData.IsSurveyReport` (pxRadioButtons; screenshot work owner 06-10-2026). */
export const PILIHAN_SURVEY_REPORT: { value: string; label: string }[] = [
  { value: 'Yes', label: 'Yes' },
  { value: 'No', label: 'No' },
]

/** Prompt values property `.ClaimPaymentType` (sumber sama dengan `PILIHAN_CLAIM_TYPE`). */
export const PILIHAN_CLAIM_PAYMENT_TYPE: { value: string; label: string }[] = [
  { value: 'Claim', label: 'Claim' },
  { value: 'AdjusterFee', label: 'Adjuster Fee' },
  { value: 'Salvage', label: 'Salvage' },
  { value: 'Adjustment', label: 'Adjustment' },
  { value: 'Retro', label: 'XOL Retro' },
]

/** Judul grid `Section/HistoricalSurveyReportDtl` (LABEL baris 1) dan medan Insured Name di atasnya. */
export const KOLOM_SURVEI = {
  tertanggung: 'Insured Name',
  tanggal: 'Date of Survey',
  oleh: 'Surveyed by (Ceding Company)',
  lossPrevention: 'Loss Prevention',
  remarks: 'Remarks',
} as const

/** Prompt values property `.Remarks` (pxDropdown pyListSource `associated`; screenshot work owner 06-10-2026). */
export const PILIHAN_REMARKS_SURVEI: { value: string; label: string }[] = [
  { value: '1', label: 'Satisfied' },
  { value: '0', label: 'Unsatisfied' },
]

export const PESAN = {
  tersimpan: 'Tersimpan.',
  terkirim: 'Berkas dikirim.',
  /** `.DateofSurvey` wajib (`Section/InputHistoricalSurveyReportDtl`). */
  surveiTanpaTanggal: 'Date of Survey wajib diisi pada baris',
  opsiTerbuka: 'Pilihan medan ini didefinisikan di rule Property yang tidak ada di korpus; nilai diketik apa adanya.',
} as const

/** `[tidak ada di korpus]` - tombol Copy Old di samping Create dan popupnya (perintah work owner 07-10-2026: "nb
 *  ttreatyin tobol copy untuk data lama mana?" -> "langsung anda kerjakan!"); teks sama dengan Copy Old EDM Treaty In. */
export const COPY_OLD = {
  tombol: 'Copy Old',
  keterangan:
    'Old NB Treaty In policies (JSON) that are not in the new tables yet. Tick the ones to copy, then press Process Copy.',
  prosesCopy: 'Process Copy',
  tutup: 'Close',
  cari: 'Search',
  pilihSemua: 'Select all',
  pilihBaris: 'Select',
  dipilih: 'selected',
  kosong: 'All old NB Treaty In policies are already in the new tables.',
  kolom: {
    id: 'NB Number',
    noOffer: 'Master ID',
    noPolis: 'Policy Number',
    insured: 'Insured Name',
    bisnis: 'Group Business',
    sob: 'SOB',
    ceding: 'Ceding',
    tglProd: 'Production Date',
    catatan: 'Notes',
  },
  status: {
    disalin: 'Copied',
    sudahAda: 'Already in the new tables',
    ditolak: 'Cannot be copied',
    gagal: 'Failed',
  },
} as const
