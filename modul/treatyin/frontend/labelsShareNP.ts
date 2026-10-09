// Label tab **Share** cabang NON-PROPORSIONAL — teks APA ADANYA dari ekspor:
// `Section/TreatyInTabsNonProportional.xml` @1695720–@3291797 (tab) dan
// `Section/Share.xml` (panel rincian baris, `pyEditAction = Share`).
//
// ⛔ Ejaan Pega dipertahankan walau keliru (`Summarry of RNM Share`,
// `Overiding Commision` di services): label adalah yang pemakai lama baca.

export const SHARE_NP = {
  judul: 'Share', // pyTitle @1705870
  nonProRate: 'Share Non Pro Rate (actual)', // @1712021 — tampil bila `IsProRate = True`
  persenRnm: '% RNM Share', // @1742644
  persenBrokerage: '% Brokerage', // @1751526
  acrossTheBoard: 'Share Across The Board', // pyCheckboxCaption @1771416
  shareKeRetro: 'Share to Other Retro', // @1798136
  brokerageRetro: 'Brokerage From Other Retro', // @1807238 — tampil bila `FacultativeShare > 0`
  placeholderAngka: '0,00',
  placeholderTeks: 'Text',
  perbaruiRingkasan: 'Update Summary', // @1832393
  tambah: 'Add',
  hapus: 'Delete',
  // Sub-tab
  shareKeRnm: 'Share to RNM :', // @2108188 — tampil bila `FacultativeShare > 0`
  persen: '%',
  bagianOf: 'Part of', // pyDefaultValue sel @2221902
  ringkasan: 'Summarry of RNM Share', // pyTitle @2324385 (ejaan Pega)
  totalSemua: 'Total All Layers RNM Share', // @2486274
  nilai: 'Value',
  perbaruiTotal: 'Update Total', // @3188300
  perbaruiNilai: 'Update Value in Share', // @3202223 — `TreatyMasterInEDM`
  tanpaBaris: 'No items',
  // `TreatyInNonAddItem` [8]: % RNM Share kosong/0 → "Value Cannot Be Empty", nol baris.
  petunjukKosong: 'Fill % RNM Share, then click Update Summary — one row per layer in Limits.',
  bukaRincian: 'Detail',
  // Panel rincian (`Share.xml`)
  kelasBisnis: 'Class of Business',
  treatyGroup: 'Treaty Group',
  cover: 'Cover',
  jenisLayer: 'Layer Type', // sel grid ekspor tanpa kepala — nama propertinya
  layer: 'Layer',
  deduksi: 'Deduction Details',
  spreadingType: 'Spreading Type',
  spreading: 'Spreading',
  totalPct: 'Total Pct', // kaki grid `.SpreadingListXOL` Share.xml @444655
  spreadingTotalPct: 'Spreading Total Pct', // pyLabelFieldValue `.SpreadingTotalPctXOL` Share.xml @473985
  totalSharePct: 'Total Share Pct',
  hapusBaris: 'Remove',
  pilihKosong: 'Choose',
  adaPesan: 'Needs attention',
} as const

/** Kolom grid `Reinsurer Name` (`TreatyIn.ShareReins`, @1879131). */
export const KOLOM_REINS = ['Reinsurer Name', 'Layer', '% Share'] as const
/** Kolom grid `Facultative Reinsurers` (`ShareFacultativeReinsurers`, @1975959). */
export const KOLOM_FAC_REINS = ['Facultative Reinsurers', 'Layer', '% Share'] as const

/**
 * Kepala grid `TreatyIn.Share` (@2149205): lima sel layer berkepala KOSONG,
 * lalu pasangan mata uang · nilai. Kepala ekspor hanya di atas sel NILAI.
 */
export const KEPALA_SHARE = ['', '', '', '', '', '', '100% Limit', '', '100% Limit', '', 'MDP', '', 'MDP', '% Share'] as const

/** Kolom "Summarry of RNM Share" (`LimitShareSummaryList`, @2339756). */
export const KOLOM_RINGKASAN_SHARE = [
  { label: 'Note', kunci: 'Note' },
  { label: '100% Limit (IDR)', kunci: 'Limit' },
  { label: '100% Limit (USD)', kunci: 'Limit2' },
  { label: 'MDP (IDR)', kunci: 'MDP' },
  { label: 'MDP (USD)', kunci: 'MDP2' },
  { label: 'Deduction (IDR)', kunci: 'Deductible' },
  { label: 'Deduction (USD)', kunci: 'Deductible2' },
  { label: 'Net Premi (IDR)', kunci: 'NetPremi' },
  { label: 'Net Premi (USD)', kunci: 'NetPremi2' },
] as const

/**
 * Kesembilan grid "Total All Layers RNM Share", tata letak ekspor: tiga
 * kolom × lima baris — kolom kiri RNM, tengah OR, kanan R/I. Sel kosong =
 * tempat yang di ekspor memang kosong.
 */
export const GRID_TOTAL_SHARE: readonly (readonly ({ judul: string; kunci: string } | null)[])[] = [
  [
    { judul: 'Total RNM Limit (RNM Share)', kunci: 'TotalShareRnmNP' },
    { judul: 'Total OR Limit', kunci: 'TotalSpreadedRnmProp' },
    { judul: 'Total R/I Limit', kunci: 'TotalSpreadedRnmRIProp' },
  ],
  [{ judul: 'Total Gross Min Premium', kunci: 'TotalShareGrossMinNP' }, null, null],
  [{ judul: 'Total Gross Premium (MDP)', kunci: 'TotalShareGrossNP' }, null, null],
  [{ judul: 'Total Deduction', kunci: 'TotalShareDeductionNP' }, null, null],
  [
    { judul: 'Total Net Premium', kunci: 'TotalShareNetNP' },
    { judul: 'Total OR Net Premium', kunci: 'TotalSpreadedNetPremi' },
    { judul: 'Total R/I Net Premium', kunci: 'TotalSpreadedNetPremiRI' },
  ],
]

/** Grid `Deduction Details` panel rincian (`Share.xml` @171027). */
export const KOLOM_DEDUKSI_SHARE = ['Description', 'Currency', 'Deduction', 'or', 'Deduction %', 'Auto Calculate %'] as const

/** Grid spreading — bernama (`Reins Type · Pct`) dan manual (`Reins Type · Pct Share`). */
export const KOLOM_SPREADING = ['Reins Type', 'Pct'] as const
export const KOLOM_SPREADING_MANUAL = ['Reins Type', 'Pct Share'] as const

/**
 * Kelima belas grid blok `hidden, reference` panel rincian (`Share.xml`
 * @516234). ⚠️ Judul blok itu menyesatkan: wadahnya `NOHEADER` (judul tidak
 * tampil) dan tidak bersyarat sendiri — ia TAMPIL bersama blok Spreading
 * (`.SpreadingTypeXOL != ''`). Tata letak ekspor: dasar · OR · R/I per baris.
 */
export const GRID_RINCIAN_SPREADING: readonly (readonly { judul: string; kunci: string }[])[] = [
  [
    { judul: 'RNM Limit', kunci: 'RnmLimitList' },
    { judul: 'Spreading OR Limit', kunci: 'RNMSpreadedListXOL' },
    { judul: 'Spreading R/I Limit', kunci: 'RNMSpreadedListRIXOL' },
  ],
  [
    { judul: 'Gross Min Premium', kunci: 'GrossPremiumMinList' },
    { judul: 'Spreading OR Min Premium', kunci: 'RNMSpreadedListGrossMinXOL' },
    { judul: 'Spreading R/I Min Premium', kunci: 'RNMSpreadedListGrossRIMinXOL' },
  ],
  [
    { judul: 'Gross Premium (MDP)', kunci: 'GrossPremiumList' },
    { judul: 'Spreading OR MDP', kunci: 'RNMSpreadedListGrossXOL' },
    { judul: 'Spreading R/I MDP', kunci: 'RNMSpreadedListGrossRIXOL' },
  ],
  [
    { judul: 'Deductions', kunci: 'DeductionTotalList' },
    { judul: 'Spreading OR Deductions', kunci: 'RNMSpreadedListDeductXOL' },
    { judul: 'Spreading R/I Deductions', kunci: 'RNMSpreadedListDeductRIXOL' },
  ],
  [
    { judul: 'Net Premium', kunci: 'NetPremiumList' },
    { judul: 'Spreading OR Net Premi', kunci: 'RNMSpreadedListNetXOL' },
    { judul: 'Spreading R/I Net Premium', kunci: 'RNMSpreadedListNetRIXOL' },
  ],
]
