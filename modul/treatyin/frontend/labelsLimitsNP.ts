// TAB LIMITS CABANG NON-PROPORSIONAL — dibaca dari ekspor
// `D:\XML_NURE\Treaty In\Section`.
//
//   `TreatyInTabsNonProportional.xml` @883398  tab Limits (TABBED):
//     grid `TreatyIn.Limits` (expandPane → `Layers`) @919339
//     Summary of Limit @1058846 · Summary of MDP (MATI, `Never`) @1201663
//     Total All Layers @1289784
//   `Layers.xml`   rincian satu layer
//   `CoBList.xml`  rincian satu Treaty Group (expandPane `TreatyGroupList`)
//
// Offset `@n` = posisi sesudah `<pyIncludedRuleXML>` dibuang dengan
// MENGHITUNG SARANG. Label tanpa kepala memang tanpa kepala di ekspor.

/** Kolom grid satu baris: label, kunci, dan desimal (`pyDecimalPlaces`). */
export interface KolomNP {
  label: string
  kunci: string
  desimal: number | null
}

export const LIMITS_NP = {
  judul: 'Limits',
  // Grid layer @919339 — kepala `Layers` di atas lima sel tanpa kepala.
  layers: 'Layers',
  of: 'of', // sel `.pyTemplateRichTextEditor` @981637 — teks "of"
  tambahLayer: 'add Layer', // @961609 — huruf kecil, ejaan ekspor
  hapus: 'Delete',
  // Layers.xml
  partOf: 'Part of', // @34564
  treatyGroup: 'Treaty Group', // @83515
  tambahGrup: 'Add Treaty Group', // @87799
  rolProfile: 'ROL Profile', // @98698
  kelasBisnis: 'Class of Business', // CoBList @52277
  egnpiLayer: 'Egnpi this layer', // @189458
  cover: 'Cover', // sel `.Cover` @258063 — label sel "Text Input", tak bermakna
  relasi: 'Currency Relation', // @263937
  pilihRelasi: 'Choose Relation', // pyNoSelectionText
  limit100: '100 % Limit', // @305129 / @494823
  agregat: 'Agregate Year Limit', // @364118 / @553773
  deductible: 'Deductible', // @422495 / @612160
  noRIP: 'No Reinstatement Premium Calculation', // @682523
  reinstatement: 'Reinstatement', // @704974
  adjRate: 'Adjustment Rate %', // @890974
  premiEarned: 'Premium Earned', // @922374
  nilai: 'Value',
  mdpMinPct: 'Min Premium %', // @966791
  mdpMinAmt: 'Min Premium Amt', // @996052
  mdpPct: 'MDP %', // @1061948
  mdp: 'MDP', // @1093318
  tambahMDP: 'Add MDP', // @1003824 / @1101086
  hapusMDP: 'Remove', // @1030655 / @1127917
  combineMDP: '(*) Combine MDP', // @1159215
  rol: 'ROL %', // @1166406
  // Summary / Total
  ringkasan: 'Summary of Limit',
  totalSemua: 'Total All Layers',
  totalROL: 'Total ROL',
  perbaruiTotal: 'Update Total', // @1638126
  perbaruiList: 'Update Value in List', // @1652282 — hanya kontrak revisi
  pilihKosong: 'Choose',
  tanpaBaris: 'No items',
  // Judul kelompok tampilan — pengelompokan layar, bukan label ekspor.
  layer: 'Layer',
  part: 'Part',
  bagianLimit: '100 % Limit · Agregate · Deductible',
  mataUang1: 'Currency 1',
  mataUang2: 'Currency 2',
  bagianPremi: 'Premium · MDP · ROL',
  belumDipilih: 'Not selected yet',
  layerBaru: 'New layer',
} as const

/** Grid layer — kolom dan desimal ekspor @943503–@1013782. */
export const KOLOM_LAYER: readonly KolomNP[] = [
  { label: '100% Limits ( IDR )', kunci: 'Limit', desimal: 2 },
  { label: 'Deductible ( IDR )', kunci: 'Deductible', desimal: 2 },
  { label: '100% Limits ( USD )', kunci: 'Limit2', desimal: 2 },
  { label: 'Deductible ( USD )', kunci: 'Deductible2', desimal: 2 },
]

/** Grid Reinstatement @770158 — urutan dan kepala ekspor. */
export const KOLOM_REINSTATEMENT: readonly KolomNP[] = [
  { label: 'Reinstatement No.', kunci: 'ReinstatementValue', desimal: null },
  { label: '% Additional Premium', kunci: 'ReinstatementPct', desimal: 2 },
  { label: 'Note', kunci: 'ReinstatementNote', desimal: null },
  { label: 'Reinstatement Premium Amount IDR (MDP x % Add Premium)', kunci: 'AdditionalAmount1', desimal: 2 },
  { label: 'Reinstatement Premium Amount USD (MDP x % Add Premium)', kunci: 'AdditionalAmount2', desimal: 2 },
  { label: 'Reinstatement %', kunci: 'AdditionalPct', desimal: 2 },
  { label: 'Reinstatement Amount IDR', kunci: 'ReinstatementAmount1', desimal: 2 },
  // Tampil bila `.Limit2 != 0` @801378.
  { label: 'Reinstatement Amount USD', kunci: 'ReinstatementAmount2', desimal: 2 },
]

/** Grid Summary of Limit @1085096. */
export const KOLOM_RINGKASAN: readonly KolomNP[] = [
  { label: 'Note', kunci: 'Note', desimal: null },
  { label: '100% Limit (IDR)', kunci: 'Limit', desimal: 2 },
  { label: '100% Limit (USD)', kunci: 'Limit2', desimal: 2 },
  { label: 'MDP (IDR)', kunci: 'MDP', desimal: 2 },
  { label: 'MDP (USD)', kunci: 'MDP2', desimal: 2 },
  { label: 'Agregate Limit (IDR)', kunci: 'AggregateLimit', desimal: 2 },
  { label: 'Agregate Limit (USD)', kunci: 'AggregateLimit2', desimal: 2 },
  { label: 'Deductible (IDR)', kunci: 'Deductible', desimal: 2 },
  { label: 'Deductible (USD)', kunci: 'Deductible2', desimal: 2 },
]

/** Keempat grid Total All Layers — kepala kolom pertama, kunci akar. */
export const GRID_TOTAL: readonly { judul: string; kunci: string }[] = [
  { judul: 'Total 100% Limit', kunci: 'TotalLimitIOONP' },
  { judul: 'Total Deductible', kunci: 'TotalLimitDeductblNP' },
  { judul: 'Total Premium Earned', kunci: 'TotalLimitPremiEarnNP' },
  { judul: 'Total MDP', kunci: 'TotalLimitMDPNP' },
]

/**
 * `TreatyMasterInEDM` — syarat tombol `Update Value in List`: kontrak
 * revisi (`EDMState` 1, 2, atau 3).
 */
export const kontrakRevisi = (edmState: string) => ['1', '2', '3'].includes(edmState.trim())
