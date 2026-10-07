// Label layar modul EDM Treaty In. Asal pola: `modul/nbtreatyin/frontend/labels.ts` (06-10-2026).
//
// ⛔ Label medan, tombol, dan judul kolom VERBATIM dari korpus `D:\XML\RNM_BRD\EDM Treaty In\` (Section / Harness /
// FlowAction / Flow) - section dan sel disebut di setiap butir. Label generik Pega ("Text Input", "Text Area",
// "Checkbox") diganti nama propertinya / caption kotak centangnya (pola NB): label itu bawaan kontrol, bukan teks
// yang pernah dibaca pengguna sebagai nama medan.

export const JUDUL = {
  /** LABEL `Section/SFAPortal_Endorsement_Treaty` S126 (pyCaption "Addendum Treaty"). */
  portal: 'Addendum Treaty',
  /** Tombol "Create New Addendum Treaty": showHarness `TreatyCreateEdm` target newDocument, judul "Create EDM". */
  buat: 'Create EDM',
  /** Flow `InputAddendumTreatyIn` Assignment2 (WB ReasTreatyInAdmin). */
  admin: 'Input Realitation',
  /** Flow `InputAddendumTreatyIn` Assignment4 / Assignment7 (WB ReasTreatyInSecHead). */
  secHead: 'Acceptance by Sec Treaty',
  /** Flow `InputAddendumTreatyIn` Assignment1 (WB ReasTreatyInDeptHead). */
  deptHead: 'Acceptance by Dept. Head',
  /** Tombol "Choose Business" `DetailPolicyTreatyInAddendum` S6: showHarness popup `pyWindowName`. */
  pilihBisnis: 'Business And SOB List',
  /** FlowAction `PolicyTreatyInDeclineConfirm` pyLabel. */
  tolak: 'Confirm Decline NB',
  /** FlowAction `ShowPolicyNoTreaty` pyLabel. */
  nomorPolis: 'Show PolicyNo',
} as const

/** Judul layar kasus per posisi (nama assignment Flow `InputAddendumTreatyIn`). */
export const JUDUL_POSISI: Record<string, string> = {
  ReasTreatyInAdmin: JUDUL.admin,
  ReasTreatyInSecHead: JUDUL.secHead,
  ReasTreatyInDeptHead: JUDUL.deptHead,
}

export const TOMBOL = {
  /** `SFAPortal_Endorsement_Treaty` S127 pyButtonLabel. */
  create: 'Create New Addendum Treaty',
  /** `SFAPortal_Endorsement_Treaty` S116 pyButtonLabel. */
  filter: 'Filter',
  /** `SFAPortal_Endorsement_Treaty` S119 pyButtonLabel + pyTooltip. */
  refresh: 'Refresh',
  refreshTooltip: 'Refresh EDM Grid',
  /** `TreatyCreateEdm` S8 pyButtonLabel. */
  buat: 'Create',
  /** `DetailPolicyTreatyInAddendum` S6. */
  chooseBusiness: 'Choose Business',
  /** `BusinessAndSOBListEDM` kolom 1. */
  choose: 'Choose',
  /** `DetailPolicyTreatyInAddendum` S19, `DetailPolicyTreatyInPropNewData2` S7. */
  save: 'Save',
  /** `DetailPolicyTreatyInAddendum` S18 / S19. */
  submit: 'Submit',
  /** `PolicyTreatyInDeclineConfirm`. */
  yes: 'Yes',
  no: 'No',
  /** `ShowPolicyNoTreaty_SC`. */
  ok: 'OK',
  /** `DetailPolicyTreatyInPropNewData2` S24. */
  hitungSelisih: 'Calculate Value Difference',
  kembali: 'Kembali',
} as const

/** Kolom tambahan tab Resolved (keputusan work owner 07-10-2026 "keluarin tanggal produksinya, DD-MM-YYYY") - BUKAN
 *  kolom RD `InboxEDM_RD2`; ditaruh di ujung supaya urutan kolom XML tetap. */
export const KOLOM_RESOLVED = { tglProd: 'Production Date' } as const

/** Kolom grid portal `Section/SFAPortal_Endorsement_Treaty` (RD `InboxEDM_RD2`) - LABEL baris judul, urut VERBATIM.
 *  ⛔ Kolom 1 dan 4 SAMA-SAMA "EDM Number" di XML (`A.pyID` vs `A.PolicyTreatyIn.EDMNo`) - dibiarkan kembar. */
export const KOLOM_PORTAL = [
  { kunci: 'id', judul: 'EDM Number' },
  { kunci: 'noOffer', judul: 'Offer No' },
  { kunci: 'noPolis', judul: 'Policy Number' },
  { kunci: 'edmNo', judul: 'EDM Number' },
  { kunci: 'sobName', judul: 'SOB' },
  { kunci: 'cedingCoName', judul: 'Ceding' },
  { kunci: 'edmType', judul: 'EDM Type' },
  { kunci: 'proportionalType', judul: 'Proportional Type' },
  { kunci: 'marketingName', judul: 'Marketing' },
  { kunci: 'nbStatus', judul: 'Status' },
] as const

/** Kolom daftar kotak masuk Beranda: kolom portal + kolom yang ada datanya, label tambahan = label Beranda NB
 *  Treaty In (keputusan work owner 06-10-2026 untuk Beranda). */
export const KOLOM_BERANDA = {
  id: 'EDM Number',
  noOffer: 'Offer No',
  noPolis: 'Policy Number',
  edmNo: 'EDM Number',
  sob: 'SOB',
  ceding: 'Ceding',
  edmType: 'EDM Type',
  jenis: 'Proportional Type',
  marketing: 'Marketing',
  status: 'Status',
  tanggal: 'Tanggal Create',
  sejak: 'Time Since Last Update',
} as const

export const PORTAL = {
  /** `.FilterTermForEndorsement` pyLabelFor (nama aksesibel). */
  filter: 'Filter Term for Endorsement',
  /** XML pyPlaceholder "Policy Number"; ⛔ perintah work owner 07-10-2026 (pencarian diperluas: nomor kasus / EDM No,
   *  nomor polis, insured, SOB, ceding, marketing, treaty group, class of business - `repository.kolomCariPortal`). */
  placeholder: 'Search EDM no, master ID, policy no, insured, business, ceding, marketing...',
  /** Nama aksesibel tombol pengosong kapsul saring (pola NB). */
  bersihkan: 'Clear field',
  kosong: 'Belum ada berkas endorsemen treaty yang terbuka.',
  kosongPetunjuk: 'Tombol Create New Addendum Treaty membuat berkas baru di antrean admin.',
  hanyaBaca: 'Berkas ini menunggu di antrean lain atau sudah selesai - hanya-baca.',
} as const

/** Switch portal In Progress / Resolved - aturan portal NB Treaty In berlaku untuk EDM (keputusan work owner
 *  07-10-2026); bawaan In Progress. */
export const STATUS_PORTAL = ['In Progress', 'Resolved'] as const
export type StatusPortal = (typeof STATUS_PORTAL)[number]

/** DT `TreatyEDMListType` VERBATIM (screenshot Pega, work owner 07-10-2026) - SAMA dengan
 *  `backend/models/edm_buat.go` `LabelJenisEDM` (dijaga `labels.test.ts`). Kode di luar daftar tampil apa adanya. */
export const LABEL_JENIS_EDM: Readonly<Record<string, string>> = {
  '1': 'Internal',
  '2': 'External',
  '3': 'Adjustment Premium',
  '4': 'Cancel Input',
}
export const labelJenisEDM = (kode: string): string => LABEL_JENIS_EDM[kode] ?? kode

/** Medan layar `Section/TreatyCreateEdm` S4 (pyCaption). */
export const BUAT = {
  noPolis: 'No Polis Treaty',
  noMaster: 'No Master Treaty',
  jenis: 'Source of Change',
} as const

/** Kolom grid popup `Section/BusinessAndSOBListEDM` (S1 dan S6 sama) - LABEL baris judul, urut VERBATIM.
 *  `kolom` = nama kolom baris jawaban `POST /bisnis` (`models.Baris`). Kolom 1 = tombol Choose tanpa judul. */
export const KOLOM_BISNIS = [
  { kolom: 'ID', judul: 'ID Revision' },
  { kolom: 'OLDID', judul: 'Previous ID' },
  { kolom: 'TreatyContractName', judul: 'Treaty Contract Name' },
  { kolom: 'ProportionType', judul: 'Proportion Type' },
  { kolom: 'LeadingReinsSource', judul: 'SOB' },
  { kolom: 'Ceding', judul: 'Ceding' },
  { kolom: 'Commencement', judul: 'Start Date' },
  { kolom: 'Termination', judul: 'End Date' },
] as const

/** Judul tab layout group `Section/DetailPolicyTreatyInAddGeneralEditable` S2108 (pyTitle S2109-S2113). */
export const TAB_DATA = ['Old Data', 'New Data', 'Value Difference'] as const
export type TabData = (typeof TAB_DATA)[number]

/** Judul wadah yang TAMPIL (`pyIncludeHeader=true`) - selebihnya NOHEADER, tanpa judul buatan (pola NB W6). */
export const BAGIAN = {
  /** `DetailPolicyTreatyInAddendum` S13; `PropNewData2` S20; `PropNewData` / `PropOldData*` / `PropValueDifference`. */
  angsuran: 'Installment Data Information',
  /** LABEL `pyIncludeLabel=true` Heading 4 di kelima section Prop. */
  ogp: 'OGP',
  onp: 'ONP',
  /** `DetailPolicyTreatyInAddPremi` S2 / S7 / S12. */
  premiLama: 'Previous Premium',
  premiBaru: 'Current Premium',
  premiSelisih: 'Total Difference',
} as const

/** Kolom grid `.SpreadingRiskList` kelima section Prop dan `AddPremi` S17 (LABEL judul + kaki "Total"). */
export const KOLOM_SPREADING = {
  treatyType: 'Type Treaty',
  share: '% Share',
  premium: 'Premium',
  claimPct: '% Share',
  claim: 'Claim',
  total: 'Total',
} as const

/** Kolom grid `.ListInstallment` section Prop (LABEL judul). "Installment" = sel `.Installment` (label generik
 *  "Text Input" diganti nama properti). */
export const KOLOM_ANGSURAN = {
  no: 'No',
  dueDate: 'Due Date',
  pct: '% Installment',
  premium: 'Premium Nusantara Re',
  total: 'Total Payment',
  installment: 'Installment',
} as const

/** Kolom grid `.ListInstallment` cabang NonProp baru (`DetailPolicyTreatyInAddendum` S13). */
export const KOLOM_ANGSURAN_NP = {
  currency: 'Currency',
  total: 'Total',
} as const

/** Kolom grid `.SuggestList` (`Section/ListSuggestEDM`). */
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

/** Pilihan Approval `.IsApproved` (pxRadioButtons): nilai "1"/"0" `DecisionTable/isApproved`, teks = NB Treaty In
 *  (`SaveViewSuggest` CARI9) - properti yang sama. */
export const PILIHAN_APPROVAL = [
  { value: '1', label: 'Accept' },
  { value: '0', label: 'Reject' },
] as const

/** Prompt values property `.TypeTax` (pxRadioButtons `associated`; screenshot work owner 06-10-2026 untuk sel NB
 *  yang sama - sel EDM salinannya). */
export const PILIHAN_TYPE_TAX: { value: string; label: string }[] = [
  { value: 'Inclusive', label: 'Inclusive' },
  { value: 'Exclusive', label: 'Exclusive' },
]

export const PESAN = {
  tersimpan: 'Tersimpan.',
  terkirim: 'Berkas dikirim.',
} as const

/** `[tidak ada di korpus]` - tombol Copy Old di samping Create dan popupnya (perintah work owner 07-10-2026 "SAMA
 *  SEPERTI MASTER PRODUCTNAME LIFE, KHUSUS BUAT SUPERUSER"); teks mengikuti Copy Old Product Name Life. */
export const COPY_OLD = {
  tombol: 'Copy Old',
  keterangan:
    'Old EDM Treaty In documents (JSON) that are not in the new tables yet. Tick the ones to copy, then press Process Copy. Generations are copied in order; generation 1 needs its NB policy in the new tables.',
  prosesCopy: 'Process Copy',
  tutup: 'Close',
  cari: 'Search',
  pilihSemua: 'Select all',
  pilihBaris: 'Select',
  dipilih: 'selected',
  kosong: 'All old EDM Treaty In documents are already in the new tables.',
  kolom: {
    id: 'EDM Number',
    noPolis: 'Policy Number',
    edmNo: 'EDM No',
    prodKe: 'Generation',
    edmType: 'EDM Type',
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
