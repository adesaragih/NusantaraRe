// Label modul PremiumList Life — tiket 01.
//
// ⛔ VERBATIM dari korpus `D:\XML\RNM_BRD\PremiumList Life\`, dengan path +
// baris di sebelahnya. Label yang diketik dari ingatan adalah label yang
// bergeser, dan yang bergeser tidak berbunyi: keduanya benar menurut dirinya
// sendiri.

// `LABEL_MENU_PREMIUMLIST` (butir menu `PremiumList`) DIBUANG - menu datar,
// keputusan work owner 30-09-2026: tombol modul berlabel `M_NAV_MENU.LABEL`.

/**
 * Kolom kotak masuk — `ReportDefinition/InboxPremiumList.xml`.
 *
 * ⛔ SEBELAS dari tiga belas. Dua yang tidak ada:
 *
 *   `.pzInsKey` b961 `Instance Handle Key` — pengenal internal Pega, bukan
 *     kolom untuk manusia.
 *   `A.KetentuanUnderwriting` b891 — **tidak punya kolom** di migrasi
 *     050–056 mana pun. Ia TIDAK ditampilkan, bukan diisi teks kosong:
 *     sel kosong di layar terbaca "memang kosong", bukan "kami tidak punya
 *     datanya".
 */
export const KOLOM_INBOX_POLIS = {
  /** b721 `pyFieldLabel`. */
  caseId: 'Case ID',
  /** b735. */
  tglCreate: 'Create Date/Time',
  /** b750. */
  createOpName: 'Create Operator Name',
  /** b764. */
  statusWork: 'Work Status',
  /** b778. */
  cedingCoName: 'CedingCoName',
  /** b794. */
  policyHolderName: 'PolicyHolderName',
  /** b807 — ⚠️ jalur propertinya `PremiumListSummary`, kolomnya di detail. */
  plNumber: 'PL_NUMBER',
  /** b821. */
  riSlipRnm: 'RISLIPRNM',
  /** b836. */
  type: 'Type',
  /** b849. */
  marketingName: 'MarketingName',
  /** b863. */
  sobName: 'SobName',
  /** b878. */
  dateReceived: 'DateReceived',
} as const

/**
 * Tombol portal — `Section/PremiumList.xml`.
 *
 * ⚠️ DUA baris untuk satu label, dan keduanya nyata: `pyLabel` pada kendali
 * (b3273, b3921) dan `pyRuleName pyButtonLabel` pada rule labelnya (b16964,
 * b16472). Yang dipakai layar adalah `pyLabel`.
 */
export const TOMBOL_POLIS = {
  /** b3273 `pyLabel` → `runActivity` b3283 → `CreateInputLife` b3291. */
  inputOffer: 'Input Offer',
  /** b3921 `pyLabel` → `runActivity` b3931 → `CreateInputLife` b3939. */
  inputPremium: 'Input Premium',
} as const

/**
 * Keputusan penawaran — VERBATIM `pyExpression` konektor
 * `InputPolicyHolder.xml`.
 *
 * ⛔ `Reject` hanya ada pada `Decision2` (b2306), yaitu sesudah Input Premium
 * Detail. Di tahap penawaran ia TIDAK punya jalur — lihat ralat AC 2 tiket 01.
 */
export const KEPUTUSAN_POLIS = {
  confirm: 'Confirm',
  reject: 'Reject',
  decline: 'Decline',
} as const

/**
 * Popup konfirmasi tombol Decision — mencegah tombol tertekan tanpa sengaja
 * (keputusan work owner 03-10-2026). Kalimatnya menyebut AKIBAT, bukan hanya
 * "apakah Anda yakin", supaya orang tahu apa yang ia setujui.
 */
export const KONFIRMASI_KEPUTUSAN = {
  judulConfirm: 'Confirm this case?',
  judulDecline: 'Decline this case?',
  /** Confirm di Life Business Offering. */
  confirmPenawaran: 'This will confirm the offer and move the case to the next step.',
  /** Confirm di Premium List Detail. */
  confirmDetail: "This will create the PL Number and finish the case. You can't undo this.",
  decline: "This will reject and close the case. You can't undo this.",
  /** Confirm terkunci: form Premium List Detail belum disimpan (03-10-2026). */
  belumTersimpan: 'You have unsaved changes. Press Save Data before Confirm.',
  ya: 'Yes, continue',
} as const

/** Popup hasil Confirm yang menerbitkan PL Number (03-10-2026). */
export const HASIL_NOMOR_PL = {
  judul: 'PL Number created',
  labelNomor: 'PL Number',
  labelWpc: 'WPC',
  kalimat: 'The case is finished (Resolved-Completed).',
  ok: 'OK',
} as const

/**
 * Judul kolom grid Premium List Detail — `Section/PL_Detail_Sec.xml`.
 *
 * ⛔ NOL judul yang dikarang. Section itu tidak memuat satu pun `pyCaption`
 * untuk medannya: Pega menampilkan label properti masing-masing, yang tidak
 * ikut ter-ekspor ke korpus. Karena itu yang di bawah adalah NAMA KOLOM apa
 * adanya, dirapikan hanya dengan mengganti garis bawah menjadi spasi —
 * bukan terjemahan yang kami karang sendiri.
 *
 * ⚠️ Nama kolomnya sendiri DATANG DARI SERVER; peta ini hanya memperindah
 * yang datang. Kolom yang tidak ada di sini tetap tampil, dengan namanya.
 */
export const JUDUL_KOLOM_PESERTA: Record<string, string> = {
  BEGIN_DATE: 'Begin Date',
  ENTRY_AGE: 'Entry Age',
  PLAN: 'Plan',
  STNC: 'STNC',
  EXPIRED_DATE: 'Expired Date',
  CURRENT_AGE: 'Current Age',
  PERIOD_MM: 'Period MM',
  POLICY_HOLDER: 'Policy Holder',
  EFFECTIVE_DATE: 'Effective Date',
  WPC: 'WPC',
  CERTIFICATE_NO: 'Certificate No',
  CURRENCY: 'Currency',
  SUM_REASURED: 'Sum Reasured',
  DEDUCTION: 'Deduction',
  BROKERAGE_FEE: 'Brokerage Fee',
  SUM_AT_RISK_RETRO: 'Sum At Risk Retro',
  RETROCEDED_SHARE: 'Retroceded Share',
  SUM_INSURED: 'Sum Insured',
  SHARE_NUSANTARA_RE_GROSS: 'Share Nusantara Re Gross',
  EM_PERCENT: 'EM Percent',
  SUM_AT_RISK_GROSS: 'Sum At Risk Gross',
  GROSS_PREMIUM: 'Gross Premium',
  CEDING_RETENTION: 'Ceding Retention',
  RATE: 'Rate',
  FACTOR: 'Factor',
  RISK: 'Risk',
  NET_PREMIUM: 'Net Premium',
  DEDUCTION_REFUND: 'Deduction Refund',
  BROKERAGE_FEE_REFUND: 'Brokerage Fee Refund',
  GROSS_PREMIUM_REFUND: 'Gross Premium Refund',
  NET_PREMIUM_REFUND: 'Net Premium Refund',
  RI_ADMIN_FEE_REFUND_RETRO: 'RI Admin Fee Refund Retro',
  SHARE_RETRO: 'Share Retro',
  GROSS_PREMIUM_REFUND_RETRO: 'Gross Premium Refund Retro',
  NET_PREMIUM_REFUND_RETRO: 'Net Premium Refund Retro',
  RI_ADMIN_FEE_RETRO: 'RI Admin Fee Retro',
  GROSS_PREMIUM_RETRO: 'Gross Premium Retro',
  NET_PREMIUM_RETRO: 'Net Premium Retro',
  // Kolom tambahan di luar PL_Detail_Sec (keputusan work owner 02-10-2026,
  // `models.KolomGridTambahan`).
  NAME_OF_INSURED: 'Name Of Insured',
  DOB: 'DOB',
  GROSS_VALUATION_BEGIN_DATE: 'Gross Valuation Begin Date',
  GROSS_VALUATION_EXPIRED_DATE: 'Gross Valuation Expired Date',
  RI_ADMIN_FEE: 'RI Admin Fee',
}

/**
 * KUNCI URUT kolom grid peserta — urutan kolom berkas CSV unggahan ceding
 * (permintaan work owner 02-10-2026). Bukan daftar kolom: kolom yang tidak
 * dimuat grid server dilewati, dan kolom server yang tidak ada di sini tetap
 * tampil sesudahnya (`susunKolom`).
 */
export const URUTAN_KOLOM_PESERTA: readonly string[] = [
  'POLICY_NO', 'POLICY_HOLDER', 'CERTIFICATE_NO', 'NAME_OF_INSURED', 'SEX', 'DOB', 'ENTRY_AGE', 'CURRENT_AGE',
  'PLAN', 'BEGIN_DATE', 'EXPIRED_DATE', 'GROSS_VALUATION_BEGIN_DATE', 'GROSS_VALUATION_EXPIRED_DATE',
  'PERIOD_MM', 'UW_STATUS', 'EM_PERCENT', 'CURRENCY', 'SUM_INSURED', 'CEDING_RETENTION',
  'SUM_REASURED', 'SHARE_NUSANTARA_RE_GROSS', 'SUM_AT_RISK_GROSS', 'GROSS_PREMIUM', 'DEDUCTION',
  'RI_ADMIN_FEE', 'BROKERAGE_FEE', 'NET_PREMIUM', 'FACTOR',
]

/** Kolom grid peserta yang selalu PALING KANAN, urut seperti ini (02-10-2026). */
export const KOLOM_PALING_KANAN = ['STNC', 'WPC'] as const

/** Label layar Premium List Detail — tiket 03. */
export const DETAIL_POLIS = {
  judul: 'Premium List Detail',
  nomor: 'PL Number',
  belumBernomor: 'Not yet numbered',
  terbitkan: 'Generate PL Number',
  /**
   * ⛔ Kalimatnya MENYEBUT apa yang harus dikerjakan lebih dahulu. "Tidak
   * dapat dinomori" membuat orang mencari kerusakan yang tidak ada.
   */
  perluPeserta:
    'Polis ini belum punya baris peserta, sehingga nomornya belum punya ' +
    'tempat tersimpan. Unggah rincian peserta lebih dahulu.',
  sudahBernomor: 'Polis ini sudah bernomor; tombolnya tidak perlu ditekan lagi.',
} as const

/**
 * Label layar unggah CSV peserta — tiket 04.
 *
 * ⚠️ Pesan penolakan TIDAK ada di sini: ia datang dari server, VERBATIM dari
 * `ValidasiUploadPL_act.xml`. Menyalinnya ke layar membuat dua sumber untuk
 * satu kalimat, dan yang satu akan tertinggal.
 */
export const UNGGAH_CSV = {
  judul: 'Upload CSV',
  // Tautan lipat aturan format CSV (03-10-2026).
  formatCsv: 'CSV format',
  pilihBerkas: 'CSV file',
  // Memeriksa CSV tanpa menyimpan (03-10-2026) — teks tombol keputusan work owner.
  tinjau: 'Validate CSV',
  // Simpan peserta + hitung summary (03-10-2026) — teks tombol keputusan work owner.
  simpan: 'Calculate CSV',
  /**
   * ⛔ Aturan pemisah DINYATAKAN DI MUKA, bukan hanya saat menolak. Korpus
   * menyebutnya enam kali di nama langkahnya ("SEPARATOR MENGGUNAKAN TITIK")
   * dan tidak pernah di pesan yang dilihat pemakai — sehingga orang baru tahu
   * aturannya setelah berkasnya ditolak.
   */
  aturanPemisah:
    'Columns: separated by a comma (,) or semicolon (;). ' +
    'Numbers: no thousands separator; decimals use a DOT (1234567.89), or a comma ' +
    '(1234567,89) in semicolon files. Not 1,234,567.89 or 1.234.567,89. ' +
    'Dates: dd/mm/yyyy.',
  perbaikiDulu:
    'Fix the rejected rows first, then validate again. While any row is ' +
    'rejected, no row is saved.',
  kolomBaris: 'Row',
  kolomKolom: 'Column',
  kolomPesan: 'Message',
  kolomSebab: 'Reason',
} as const

/**
 * Label layar rekap — tiket 05a bagian 2, `Section/ShowLifePremiumSummary.xml`.
 */
/**
 * Panel Summary di Premium List Detail (keputusan work owner 03-10-2026):
 * rekap yang dihitung saat Save Data / Save peserta CSV, dan yang dipakai
 * Confirm.
 */
/** Dua tab panel Participants (03-10-2026): rincian peserta dan Summary. */
export const TAB_PESERTA = ['rincian', 'summary'] as const
export type TabPeserta = (typeof TAB_PESERTA)[number]
export const LABEL_TAB_PESERTA: Record<TabPeserta, string> = {
  rincian: 'Participant Details',
  summary: 'Summary',
}

export const PANEL_SUMMARY = {
  judul: 'Summary',
  kosong:
    'No summary yet. It is calculated when participants are saved or Save Data is ' +
    'pressed, once Type is filled in.',
} as const

export const SUMMARY_POLIS = {
  judul: 'Summary Premium Life',
  submit: 'Submit',
  tanpaRekap: 'Belum ada rekap: polis ini belum punya baris peserta.',
  /** Tiket 05b — `finishAssignment` → `END52`, `pyWorkStatus` b947. */
  ditutup: 'Kasus ditutup: Resolved-Completed.',
  /** Tiket 06 — efek keluar sesudah simpan. */
  efekDilewati: 'Kiriman ke Arasapas dilewati: lingkungan ini bukan produksi.',
  efekGagal: 'Rekap tersimpan, tetapi efek keluar gagal:',
  efekTerantre: 'Kegagalannya terantre untuk dicoba ulang.',
  efekTidakTerantre: 'Sebagian kegagalan TIDAK dapat diantre — hubungi tim operasi.',
  /** ⛔ `.COB` baris rekap tidak ditetapkan rule mana pun di korpus. */
  cobTanpaSumber:
    'Kolom COB tidak terisi: tidak satu pun rule PremiumList Life menetapkan ' +
    '.COB pada baris rekap mata uang.',
} as const

/** Satu kolom grid rekap: judul VERBATIM dan properti sumbernya. */
export interface KolomRekap {
  judul: string
  kolom: string
}

/**
 * Grid rekap per `Type`, VERBATIM.
 *
 * `[terverifikasi]` Empat layout ber-`pyContainerVisibleWhen` `.Type = 'QR'`
 * / `'QP'` / `'TP'` / `'TR'` di `ShowLifePremiumSummary.xml` - judul grid
 * `Summary` di b4018 / b8761 / b15196 / b20690 (nomor baris =
 * `sed -e 's/></>
</g'`); judul dan properti kolom disalin urut dokumen.
 *
 * ⚠️ KEANEHAN WARISAN, disalin apa adanya: pada QP, TP, dan TR judul
 * `PREMIUM DEDUCTION` berdiri di atas properti `.COMMISSION`. Pada QR judul
 * `DEDUCTION` berdiri di atas `.DEDUCTION` dan kolom `.COMMISSION` tidak ada.
 */
export const GRID_REKAP: Record<string, readonly KolomRekap[]> = {
  QR: [
    { judul: 'COB', kolom: 'COB' },
    { judul: 'PL NUMBER', kolom: 'PL_NUMBER' },
    { judul: 'CURRENCY', kolom: 'CURRENCY' },
    { judul: 'PREMIUM', kolom: 'PREMIUM' },
    { judul: 'DEDUCTION', kolom: 'DEDUCTION' },
    { judul: 'BROKERAGE FEE', kolom: 'BROKERAGE_FEE' },
    { judul: 'RI ADMIN FEE', kolom: 'RI_ADMIN_FEE' },
    { judul: 'TAX', kolom: 'TAX' },
    { judul: 'PROF COMM', kolom: 'PROF_COMM' },
    { judul: 'CLAIM', kolom: 'CLAIM' },
    { judul: 'BALANCE', kolom: 'BALANCE' },
  ],
  QP: [
    { judul: 'COB', kolom: 'COB' },
    { judul: 'PL NUMBER', kolom: 'PL_NUMBER' },
    { judul: 'CURRENCY', kolom: 'CURRENCY' },
    { judul: 'PREMIUM', kolom: 'PREMIUM' },
    { judul: 'PREMIUM DEDUCTION', kolom: 'COMMISSION' },
    { judul: 'BROKERAGE FEE', kolom: 'BROKERAGE_FEE' },
    { judul: 'OVR COMM', kolom: 'OVR_COMM' },
    { judul: 'TAX', kolom: 'TAX' },
    { judul: 'GROSS_PREMIUM_REFUND', kolom: 'GROSS_PREMIUM_REFUND' },
    { judul: 'DEDUCTION REFUND', kolom: 'DEDUCTION_REFUND' },
    { judul: 'RI ADMIN FEE REFUND', kolom: 'RI_ADMIN_FEE_REFUND' },
    { judul: 'BROKERAGE FEE REFUND', kolom: 'BROKERAGE_FEE_REFUND' },
    { judul: 'TAX REFUND', kolom: 'TAX_REFUND' },
    { judul: 'CLAIM_AMOUNT', kolom: 'CLAIM_AMOUNT' },
    { judul: 'NET_PREMIUM_REFUND', kolom: 'NET_PREMIUM_REFUND' },
    { judul: 'BALANCE', kolom: 'BALANCE' },
  ],
  TP: [
    { judul: 'COB', kolom: 'COB' },
    { judul: 'PL NUMBER', kolom: 'PL_NUMBER' },
    { judul: 'CURRENCY', kolom: 'CURRENCY' },
    { judul: 'SHARE RETRO', kolom: 'SHARE_RETRO' },
    { judul: 'GROSS PREMIUM RETRO', kolom: 'GROSS_PREMIUM_RETRO' },
    { judul: 'PREMIUM DEDUCTION', kolom: 'COMMISSION' },
    { judul: 'BROKERAGE FEE RETRO', kolom: 'BROKERAGE_FEE_RETRO' },
    { judul: 'DISCOUNT PREMIUM RETRO', kolom: 'DISCOUNT_PREMIUM_RETRO' },
    { judul: 'RI ADMIN FEE RETRO', kolom: 'RI_ADMIN_FEE_RETRO' },
    { judul: 'TAX', kolom: 'TAX' },
    { judul: 'PROF COMM', kolom: 'PROF_COMM' },
    { judul: 'CLAIM', kolom: 'CLAIM' },
    { judul: 'BALANCE', kolom: 'BALANCE' },
  ],
  TR: [
    { judul: 'COB', kolom: 'COB' },
    { judul: 'PL NUMBER', kolom: 'PL_NUMBER' },
    { judul: 'CURRENCY', kolom: 'CURRENCY' },
    { judul: 'SHARE RETRO', kolom: 'SHARE_RETRO' },
    { judul: 'GROSS PREMIUM REFUND RETRO', kolom: 'GROSS_PREMIUM_REFUND_RETRO' },
    { judul: 'PREMIUM DEDUCTION', kolom: 'COMMISSION' },
    { judul: 'BROKERAGE FEE REFUND RETRO', kolom: 'BROKERAGE_FEE_REFUND_RETRO' },
    { judul: 'DISCOUNT PREMIUM REFUND RETRO', kolom: 'DISCOUNT_PREMIUM_REFUND_RETRO' },
    { judul: 'RI ADMIN FEE REFUND RETRO', kolom: 'RI_ADMIN_FEE_REFUND_RETRO' },
    { judul: 'TAX', kolom: 'TAX' },
    { judul: 'PROF COMM', kolom: 'PROF_COMM' },
    { judul: 'CLAIM', kolom: 'CLAIM' },
    { judul: 'BALANCE', kolom: 'BALANCE' },
  ],
}

/**
 * Layar Input Offer — `Section/InputOfferLife.xml` (tiket 01 bagian 3).
 *
 * ⛔ VERBATIM `pyLabel` tiap sel, termasuk DUA radio yang sama-sama berlabel
 * `Status` (satu tampil saat bendera "0", satu saat "1"). Dibaca dari salinan
 * korpus `kelvin\PremiumListLife (Done)\` 30-09-2026.
 */
export const LABEL_PENAWARAN = {
  judul: 'Life Business Offering',
  noOffer: 'Confirmation Number',
  cedingCoName: 'Ceding Name',
  policyHolderName: 'Policy Holder',
  pilihCeding: 'Choose Ceding Name',
  /**
   * ⚠️ `[dugaan]` Tombol pembuka `PolicyHolder_Harness` ber-`pyLabel` bawaan
   * `Button` — kendali tanpa teks. Kalimat ini karangan layar baru supaya
   * tombolnya dapat dibaca; bukan label korpus.
   */
  pilihPemegang: 'Choose Policy Holder',
  typeCeding: 'System Reinsurance',
  jenisAsuransi: 'Reinsurance Type',
  businessCode: 'Class of Business',
  dateReceived: 'Email Received Date',
  status: 'Status',
  description: 'Comment',
  batasUsiaPeserta: 'Age Limit',
  periodePertanggungan: 'Coverage Period',
  sumInsured: 'Sum Insured',
  tanggalPenawaran: 'Offering Date',
  tanggalRespon: 'Response Date',
  tanggalKonfirmasi: 'Confirmation Date',
  tbc: 'Input TBC',
  tanggalTbc: 'Max TBC',
  statusUpdate: 'Status Update',
  keteranganMarketing: 'Marketing Note',
  qqName: 'Insured Name',
  jenisUsaha: 'Occupation',
  ketentuanUnderwriting: 'Underwriting Policy',
  tanggalKonfirmasiBalik: 'Re-Confirmation Date',
  tanggalRealisasi: 'Realization Date',
  tanggalBind: 'Binding Date',
  statusFinal: 'Final Status',
  simpan: 'Save Offer',
  cari: 'Search',
  pilih: 'Choose',
  kolomId: 'ID',
  kolomNama: 'Name',
  /** Grid `PolicyHolder_Section` kolom ketiga (`.BU_Note`). */
  kolomBisnis: 'Business',
} as const

/** Judul grid riwayat penawaran — `InputOfferLife.xml` (`.OfferFacIn.ViewSuggest`). */
export const KOLOM_RIWAYAT_PENAWARAN = {
  dateSuggest: 'Date',
  picSuggest: 'PIC',
  isCedingConfirm: 'Status',
  initialSuggest: 'Position',
  commentSuggest: 'Comment',
} as const

/**
 * Layar Input Premium Detail — `Section/ShowLifePremiumDetail.xml` (tiket 03
 * bagian 2). VERBATIM `pyLabel` tiap sel; dibaca 01-10-2026.
 */
export const LABEL_DATA_POLIS = {
  judul: 'Input Life Premium Detail',
  pilihProduk: 'Choose Product Name',
  productName: 'Product Name',
  productNameId: 'Product Name ID',
  type: 'Type',
  typeCeding: 'System Reinsurance',
  riSlip: 'R/I SLIP RNM No.',
  proRateType: 'Premium Payment Method',
  marketing: 'Marketing Officer',
  sumInsured: 'Sum Insured',
  annuityInterest: 'Annuity Interest',
  premiumRefundFactor: 'Premium Refund Factor',
  noOffer: 'Confirmation Number',
  ketentuanUnderwriting: 'Underwriting Policy',
  ceding: 'Ceding',
  policyHolder: 'Policy Holder',
  jenisAsuransi: 'Reinsurance Type',
  businessCode: 'Class of Business',
  batasUsia: 'Age Limit',
  periode: 'Coverage Period',
  catatanBilling: 'BILLING NAME IS MANDATORY FOR TYPE TP & TR',
  billing: 'Billing Name',
  pilihBilling: 'Choose Billing Name',
  retro: 'Retrocessionaire',
  pilihRetro: 'Choose Retrocessionaire',
  dateReceived: 'Email Received Date',
  tanggalPenawaran: 'Offering Date',
  tanggalRespon: 'Response Date',
  tanggalKonfirmasi: 'Confirmation Date',
  tanggalKonfirmasiBalik: 'Re-Confirmation Date',
  tanggalRealisasi: 'Realization Date',
  tanggalBind: 'Binding Date',
  tanggalTbc: 'Max TBC',
  wpc: 'WPC',
  status: 'Status',
  statusUpdate: 'Status Update',
  keteranganMarketing: 'Marketing Note',
  simpan: 'Save Data',
} as const

/** Kolom grid popup Choose Product Name — `Section/ChooseProdName.xml`. */
export const KOLOM_PRODUK = {
  id: 'ID',
  inwardName: 'TREATY NAME',
  ceding: 'CEDING',
  sob: 'SOB',
  policyHolder: 'POLICY HOLDER',
} as const

/**
 * Teks pilihan kosong dropdown modul ini — bahasa Inggris (permintaan work owner
 * 01-10-2026). Diteruskan lewat prop `kosong` komponen `Pilih` inti, sehingga
 * bawaan inti (`-- pilih --`) untuk modul lain tidak berubah.
 */
export const TEKS_PILIH = '-- choose --'

/**
 * Teks tombol pilih yang menempel di kotak isian Premium List Detail (02-10-2026).
 * Nama lengkapnya (mis. "Choose Product Name") tetap di `aria-label` dan `title`.
 */
export const TEKS_TOMBOL_PILIH = 'Choose'
