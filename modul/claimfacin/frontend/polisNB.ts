// Tampilan polis NB Fac In HANYA-BACA di Claim Fac In (perintah work owner 10-10-2026: "SAAT CHOOSE, TAMPILAN POLIS NYA
// MENGIKUTI TAMPILAN POLIS JUGA. IKUTI MODUL NB FACIN. TAMPILAN READ ONLY" + "DI TAB POLICY DETAIL JUGA BEGITU"; pilihan
// "salin tampilan di Claim Fac In"). Rupa = layar Inward Facultative `modul/nbfacin/frontend/pages/InwardFacultative.tsx`
// dan tab-tabnya, DISALIN (pola, bukan impor - batas modul), dibaca dari halaman polis `OfferFacIn` kasus klaim (JSON_POLIS
// atau tabel datar nbfacin). Label VERBATIM `modul/nbfacin/frontend/labels.ts` (asal korpus NB FacIn). Berkas ini hanya
// fungsi murni atas halaman (diuji tanpa DOM); komponennya `components/PolisNB.tsx`.
//
// Total yang di nbfacin DIHITUNG (bukan disimpan Pega: TotalTSIList, TotalTSIPremiGrossList, TotalTSIPremiSpreadRNM,
// ringkasan TempTotal* / TotalSpreading*) dihitung ulang di sini dengan rumus nbfacin yang sama, eksak tanpa float. Nilai
// yang Pega simpan di polis (TotalTSINusaReSpreading, TotalPremiNusaRe, ClaimEstimation baris) ditampilkan apa adanya.

import { jumlahDesimal } from '../../../inti/frontend/lib/desimal'
import type { Baris, Halaman } from './api'
import { tampilAngka } from './nilai'

// ---------------------------------------------------------------- jalur halaman polis

export const AKAR_NB = 'OfferFacIn'
const Q = AKAR_NB + '.QuotationData.'
export const LOKASI_NB = AKAR_NB + '.LocationList'
export const CEDING_NB = Q + 'CedingCoList'
export const CEDANT_NB = AKAR_NB + '.CedingCedantList'

/** Jalur daftar anak baris ke-n (berbasis 1) - bentuk `models.JalurAnak`. */
export const anakNB = (daftar: string, n: number, sub: string) => `${daftar}(${n}).${sub}`

export const itemNB = (lokasi: number) => anakNB(LOKASI_NB, lokasi, 'Property.PropertyItemList')
export const coverageNB = (lokasi: number, item: number) => anakNB(itemNB(lokasi), item, 'CoverageList')
export const deductibleNB = (lokasi: number, item: number, cov: number) =>
  anakNB(coverageNB(lokasi, item), cov, 'DeductibleList')
export const spreadingNB = (lokasi: number, item: number, cov: number) =>
  anakNB(coverageNB(lokasi, item), cov, 'SpreadingList')

/** Baris daftar halaman (kosong bila tidak ada). */
export function daftarNB(h: Halaman, j: string): Baris[] {
  return h.daftar[j] ?? []
}

/** Nilai skalar halaman polis (`OfferFacIn.<j>`). */
export const nilaiNB = (h: Halaman, j: string) => h.nilai[AKAR_NB + '.' + j] ?? ''

/** Boolean Pega (teks "true" / "false"); tabel datar dapat membawa 1 / 0. */
export const benarNB = (v: string | undefined) => v === 'true' || v === '1'

// ---------------------------------------------------------------- label (VERBATIM nbfacin/labels.ts)

/** `PERIODE` (Section Periode) - label medan blok General. */
export const UMUM_NB = {
  judul: 'General',
  reffNumber: 'Reff. number',
  businessStatus: 'Business status',
  insuredName: 'Insured name',
  qqName: 'QQ name',
  beginDate: 'Begin date',
  offeringDate: 'Offering date',
  policyType: 'Policy Type',
  riskScoring: 'Risk Scoring',
  classOfBusiness: 'Class of business',
  typeFacultative: 'Type facultative',
  sourceOfBusiness: 'Source of business',
  cedingCoName: 'Ceding co name',
  groupName: 'Group Name',
  endDate: 'End date',
  followingPolicyNumber: 'Following Policy Number',
  oldPolicyNumber: 'Old Policy Number',
  marketingName: 'Marketing Name',
  day: 'Day',
} as const

/** `PILIHAN_PERIODE` (tangkapan layar work owner NB 03-10-2026). */
export const PILIHAN_POLICY_TYPE = ['Individual Policy', 'Master Policy'] as const
export const PILIHAN_DAY = ['365', '366'] as const

/** `TAB_DETAIL` (centang Show Detail nbfacin dibuang - tab langsung tampil, work owner 11-10-2026). */
export const TAB_NB = [
  'Object',
  'Coverage',
  'Clauses',
  'Spreading',
  'Inw Fac Cedant Panels',
  'Premium Deduction / Brokerage Fee',
  'Loss Record',
  'Payment',
  'Scoring Risk',
  'Correspondence',
] as const
export type TabNB = (typeof TAB_NB)[number]

/** `TEKS_INWARD.isiTab` / `kosong`. */
export const TEKS_NB = {
  isiTab: 'Isi tab',
  kosong: 'No items',
  spreadingKosong: 'Belum ada spreading.',
  /** `TEKS_INWARD.cobBelum` nbfacin - isi tab COB di luar FIRE / Aneka / Marine Cargo. */
  cobBelum: (tab: string, cob: string) =>
    `Tab ${tab} untuk Class of Business ${cob} belum tersedia di sistem baru (saat ini baru FIRE).`,
} as const

export const KOLOM_OBJEK_NB = ['Top Risk', 'No.', 'Object Name', 'Location'] as const
export const UKURAN_OBJEK_NB = 10
export const SUBTAB_OBJEK_NB = [
  'Object Address',
  'Surrounding Risk',
  'Object Item',
  'Occupation',
  'FEA',
  'Loss Record',
  'Loss Record Internal',
] as const
export type SubTabObjekNB = (typeof SUBTAB_OBJEK_NB)[number]

/** `OBJECT_ADDRESS` (Section ObjectDetails). */
export const ALAMAT_NB = {
  objectNo: 'Object No.',
  objectType: 'Object Type',
  materialDamage: 'Material Damage',
  topRisk: 'Top Risk',
  objectName: 'Object Name',
  judulRisk: 'Risk Address',
  type: 'Type',
  address: 'Address',
  buildingNo: 'Building No.',
  zipCode: 'Zip Code',
  country: 'Country',
  riskLocation: 'Risk Location',
  territory: 'Territory',
  city: 'City',
  district: 'District',
  province: 'Province',
  riskAddressId: 'Risk Address ID',
  judulBangunan: 'Building Construction',
  numberOfFloor: 'Number of Floor',
  roofType: 'Roof Type',
  wallType: 'Wall Type',
  floorType: 'Floor Type',
  partitionType: 'Partition Type',
  supportWallType: 'Support Wall Type',
  otherType: 'Other Type',
} as const

/** Object Name hanya tampil bila Object Type = `OBJECT_TYPE_LAINNYA`. */
export const OBJECT_TYPE_LAINNYA_NB = 'Others'

export interface OpsiNB {
  value: string
  label: string
}

const urutKode = (labels: string[]): OpsiNB[] => labels.map((label, i) => ({ value: String(i + 1), label }))

/** `OPSI_ROOF_TYPE` / `OPSI_WALL_TYPE` / `OPSI_FLOOR_TYPE` (PromptList DDL RoofType / WallType / FloorType). */
export const OPSI_ROOF_NB = urutKode([
  'Dak Beton',
  'Genteng Beton',
  'Bilik',
  'Genteng Tanah Liat',
  'Sirap',
  'Seng Gelombang',
  'Seng Lembaran',
  'Aluminium Gelombang',
  'Aluminium Lembaran',
  'Kaca',
  'Plastik / Policarbon Lembaran',
  'Plastik / Policarbon Gelombang',
  'Daun',
  'Lain-lain',
])
export const OPSI_WALL_NB = urutKode([
  'Batu bata',
  'Kayu / Papan',
  'Semi Permanen',
  'Bilik',
  'Batako',
  'Panel Beton',
  'Seng / Plat',
  'Kaca',
  'Lain-lain',
])
export const OPSI_FLOOR_NB: OpsiNB[] = [
  { value: 'KELAS III', label: 'Keramik' },
  { value: 'KELAS II', label: 'Kayu' },
  { value: 'KELAS I', label: 'Lain-lain' },
]

/** `SISI_SEKITAR` / `MEDAN_SISI` / `LAIN_SEKITAR` (Section SurroundingRisk). */
export const SISI_NB = [
  { awalan: 'Front', judul: 'Front' },
  { awalan: 'Left', judul: 'Left' },
  { awalan: 'Back', judul: 'Back' },
  { awalan: 'Right', judul: 'Right' },
] as const
export const MEDAN_SISI_NB = {
  occupation: 'Occupation',
  construction: 'Construction',
  distance: 'Distance (meter)',
  note: 'Note',
} as const
export const LAIN_SEKITAR_NB = {
  judul: 'Other Description',
  ownership: 'Ownership',
  housekeepingStatus: 'House keeping Status',
  floodAreaStatus: 'Flood Area Status',
  floodArea: 'Flood Area',
  housekeepingRemark: 'Housekeeping Remark',
  productionProcess: 'Production Process',
  hotWork: 'Job With a Chance of Fire',
  flammable: 'Flammable Item',
  faktor:
    'Choose the factors below and give a description for any of it. Please Skip this step if the factor you need is not available.',
} as const
/** Flood Area tampil bila Flood Area Status = `FLOOD_STATUS_TAMPIL`. */
export const FLOOD_STATUS_TAMPIL_NB = '0'
export const OPSI_CONSTRUCTION_NB: OpsiNB[] = [
  { value: 'Reinforce concrete', label: 'I' },
  {
    value:
      'Building with all structural members and at least 80% of the walls constructed with non - combustible materials, roof covering to be non - combustible, but roof supports may be of wood or other combustible materials',
    label: 'II',
  },
  { value: 'All other buildings', label: 'III' },
]
export const OPSI_OWNERSHIP_NB: OpsiNB[] = [
  { value: '2', label: 'Not Informed' },
  { value: '0', label: 'Own' },
  { value: '1', label: 'Rent' },
]
export const OPSI_HOUSEKEEPING_NB: OpsiNB[] = [
  { value: '0', label: 'Not Informed' },
  { value: '1', label: 'Good' },
  { value: '2', label: 'Fair' },
  { value: '3', label: 'Poor' },
]
export const OPSI_FLOOD_STATUS_NB: OpsiNB[] = [
  { value: '2', label: 'Not Informed' },
  { value: '0', label: 'Yes' },
  { value: '1', label: 'No' },
]
export const OPSI_FLOOD_AREA_NB: OpsiNB[] = [
  { value: '1', label: 'Low' },
  { value: '2', label: 'Medium' },
  { value: '3', label: 'High' },
  { value: '4', label: 'Very High' },
]

/** `GRID_ITEM` / `FORM_ITEM` (Section PropertyItemList). */
export const KOLOM_ITEM_NB = [
  'Object Item Type',
  'Condition',
  'Year',
  'Unit(s)',
  'Currency',
  'TSI Object Item',
] as const
export const KOLOM_TOTAL_ITEM_NB = ['Currency', 'Total TSI'] as const
export const FORM_ITEM_NB = {
  itemType: 'Object Item Type',
  note: 'Object Item Note',
  year: 'Year',
  unit: 'Unit(s)',
  condition: 'Condition',
  currency: 'Currency',
  tsi: 'TSI Object Item (All Unit)',
  yearOfPlanting: 'Year of Planting',
  noOfTree: 'No of Trees',
  areaHectar: 'Area ( Hectar )',
  remark: 'Remark of',
  adjustable: 'Adjustable',
  pctAdjust: 'Adjustment Pct. %',
} as const
export const OPSI_CONDITION_NB: OpsiNB[] = [
  { value: '1', label: 'Good' },
  { value: '2', label: 'Fair' },
  { value: '3', label: 'Poor' },
]

/** `GRID_OKUPASI` / `FORM_OKUPASI`. */
export const KOLOM_OKUPASI_NB = ['Occupation ID', 'Occupation Name', 'Class Of Construction'] as const
export const FORM_OKUPASI_NB = {
  occupationId: 'Occupation ID',
  occupationName: 'Occupation Name',
  konstruksi: 'Class of Construction',
} as const

/** `GRID_FEA` / `FORM_FEA`. */
export const KOLOM_FEA_NB = [
  'APAR',
  'Sprinkler',
  'Smoke Detector & Alarm',
  'Hydrant',
  'Private Truck Brigade',
  'Others Info',
] as const
export const FORM_FEA_NB = {
  apar: 'APAR (Unit)',
  sprinkler: 'Sprinkler (Unit)',
  smokeDetector: 'Smoke Detector & Alarm (Unit)',
  hydrant: 'Hydrant (Unit)',
  privateTruckBrigade: 'Private Truck Brigade (Unit)',
  privateFireBrigade: 'Private Team Fire Brigade',
  teamSopSafety: 'Team & SOP Safety',
  teamSopRiskManagement: 'Team & SOP Risk Management',
  info: 'Others Info',
} as const

/** `GRID_KERUGIAN` / `FORM_KERUGIAN` / `LABEL_REMARKS` / `LOSS_RATIO` / `GRID_KLAIM_INTERNAL`. */
export const KOLOM_KERUGIAN_NB = [
  'Date of Loss',
  'Insured Name',
  'Loss Object',
  'Currency',
  'Total of Loss',
  'Total Claim',
] as const
export const FORM_KERUGIAN_NB = {
  insuredName: 'Insured Name',
  dateOfLoss: 'Date of Loss',
  lossObject: 'Loss Object',
  currency: 'Currency',
  claim: 'Total Claim (100%)',
  preventionOfLoss: 'Prevention Of Loss',
  causeOfLoss: 'Cause of Loss',
  remarks: 'Remarks',
  detail: 'Loss Detail',
} as const
export const LOSS_RATIO_NB = [
  { label: 'LR 1 Year', kunci: 'LossRatio1YearAmount' },
  { label: '%LR 1 Years', kunci: 'LossRatio1YearPercent' },
  { label: 'LR 3 - 5 Years', kunci: 'LossRatio35YearAmount' },
  { label: '%LR 3 - 5 Years', kunci: 'LossRatio35YearPercent' },
] as const
export const KOLOM_KLAIM_INTERNAL_NB = [
  'Year',
  'Date of Loss',
  'Location No',
  'Location',
  'Currency',
  'Premium',
  'O/S Claim',
  'Acccepted Claim',
  'Incurred Claim',
  'Loss Ratio',
  'Remark',
] as const

/** `GRID_COV_OBJEK` / `GRID_COV_ITEM` / `GRID_COV_TOTAL` / `RINGKASAN_COV_*` / `GRID_COVERAGE` / `FORM_COV`. */
export const KOLOM_COV_OBJEK_NB = ['No.', 'Object Name', 'Location'] as const
export const KOLOM_COV_ITEM_NB = [
  'Object Item Type',
  'Currency',
  'TSI Object Item',
  'Total Gross Premium',
  '‰ Total Net Rate',
] as const
export const KOLOM_COV_TOTAL_NB = ['Currency', 'Total TSI', 'Total Premi', 'Total Gross Rate'] as const
export const RINGKASAN_COV_ITEM_NB = ['Object Item Type', 'Currency', 'Total TSI', 'Total Premium'] as const
export const RINGKASAN_COV_UANG_NB = ['Currency', 'Total TSI', 'Total Premium', 'Rate'] as const
export const KOLOM_COVERAGE_NB = ['Coverage', '‰ Standard Rate', 'Premi'] as const
export const FORM_COV_NB = {
  coverageBasis: 'Coverage Basis',
  coverage: 'Coverage',
  accumulationCode: 'Accumulation Code',
  accumulationAddress: 'Accumulation Address',
  conditions: 'Conditions',
  day: 'Days',
  tsi: 'TSI',
  indemnity: 'Indemnity',
  rate: '‰ Gross Rate',
  firstLoss: '% First Loss',
  discountPercentage: '% Discount',
  tsiLiability: 'First Loss',
  netRate: '‰ Net Rate',
  limitOfLiability: 'Limit of Liability',
  pctLol: '% Limit of Liability',
  proRate: '% Pro Rate',
  unit: 'Indemnity Unit',
  indemnityPercentage: '% Indemnity',
  firstScale: '% First Scale',
  sublimit: '% Sub Limit',
  lostLimit: '% Loss Limit',
  emlPml: '% EML/PML',
  discount: 'Discount',
  premium: 'Gross Premium',
} as const
export const PILIHAN_DAY_COV_NB = ['365', '366', '360'] as const
export const AKUMULASI_KOSONG_NB = '---'
export const OPSI_INDEMNITY_UNIT_NB: OpsiNB[] = [
  { value: '0', label: 'Day (s)' },
  { value: '1', label: 'Month (s)' },
  { value: '2', label: 'Year (s)' },
]
export const OPSI_COVERAGE_BASIS_NB: OpsiNB[] = [
  { value: '1', label: 'Sum Insured Basis' },
  { value: '2', label: 'First Loss Basis' },
  { value: '3', label: 'EML / PML Basis' },
  { value: '4', label: 'Sub Limit Basis' },
  { value: '5', label: 'Layering Basis' },
]

/** `JUDUL_DEDUCTIBLE` / `FORM_DEDUCTIBLE` / `KOLOM_TIME_EXCESS` / `OPSI_*` deductible. */
export const JUDUL_DEDUCTIBLE_NB = 'Deductible'
export const KOLOM_DEDUCTIBLE_NB = [
  'Pct Deductible',
  'Type Deductible',
  'MinMax',
  'Currency',
  'Amount',
  'Condition',
  'Time Excess (Days)',
] as const
export const FORM_DEDUCTIBLE_NB = {
  typeDeductible: 'Type Deductible',
  pctDeductible: 'Pct Deductible',
  minMax: 'MinMax',
  currency: 'Currency',
  condition: 'Condition',
  amount: 'Amount',
  timeExcess: 'Time Excess (In Days)',
  /** `FORM_COVERAGE_CARGO.deskripsi` - Deductible Marine Cargo. */
  deskripsi: 'Deskripsi',
} as const
export const OPSI_MINMAX_NB: OpsiNB[] = [
  { value: '1', label: 'Min' },
  { value: '2', label: 'Max' },
  { value: '3', label: 'Or' },
]
export const OPSI_KONDISI_DEDUCTIBLE_NB: OpsiNB[] = [
  { value: '1', label: 'Any One Occurrence' },
  { value: '2', label: 'Any One Accident' },
  { value: '3', label: 'Each and Every Loss' },
  { value: '4', label: 'Each and Every Claim' },
  { value: '5', label: 'Other' },
]
/** Condition = Other -> teks `InputCondition`. */
export const KONDISI_LAIN_NB = '5'
export const OPSI_TYPE_DEDUCTIBLE_NB: OpsiNB[] = [
  { value: '1', label: '% Of Claim' },
  { value: '2', label: '% Of Loss' },
  { value: '3', label: '% Of TSI' },
  { value: '4', label: '% Of Approved Loss Value' },
  { value: '5', label: '% Of Recoverable Claim Amount' },
  { value: '6', label: '% Of Recoverable Amount' },
  { value: '7', label: 'In Amount' },
  { value: '0', label: 'NIL' },
]
export const OPSI_TYPE_DEDUCTIBLE2_NB: OpsiNB[] = [
  { value: '1', label: '% Of Claim' },
  { value: '2', label: '% Of Loss' },
  { value: '3', label: '% Of TSI' },
  { value: '4', label: '% Of TSI Whichever Is Higher' },
  { value: '0', label: 'NIL' },
]
/** Type Deductible "In Amount": Pct / MinMax disembunyikan; MinMax "Or": Type / Pct kedua tampil. */
export const TYPE_IN_AMOUNT_NB = '7'
export const MINMAX_OR_NB = '3'

/** `SPREADING` / `CEDANT` / `OPSI_SHARE_CEDANT_TYPE`. */
export const SPREADING_NB = {
  percentShare: '% Share RNM',
  kolomLokasi: ['No.', 'Object Name', 'Location'],
  kolomItem: ['Object Item Type', 'Currency', 'TSI Object Item', 'Total Gross Premium', 'Total Premium RNM'],
  kolomCoverage: ['Coverage', '‰ Standard Rate', 'Gross Premium', 'Premium RNM'],
  detailCoverage: {
    tsi: 'TSI 100%',
    tsiLiability: 'TSI Liability',
    tsiNusantaraRe: 'TSI Nusantara Re',
    premium: 'Gross Premium',
    premiNusantaraRe: 'Premium Nusantara Re',
  },
  kolomSpreading: [
    'Type Treaty',
    '% Share',
    'TSI 100% Spreaded',
    'TSI Spreaded (RNM)',
    'Limit Of Liability',
    'Premium Spreaded',
  ],
  kolomTotal: [
    'Currency',
    'Treaty Name',
    '(%) Share',
    'Total TSI Top Risk Spreaded',
    'Total TSI Spreaded',
    'Total Limit of Liability',
    'Total Premium Spreaded',
  ],
  kolomRingkasanTreaty: [
    'Currency',
    'Treaty Name',
    '(%) Share',
    'Total TSI Top Risk Spreaded',
    'Total TSI Spreaded',
    'Total Limit of Liability',
    'Total First Loss',
    'Total Premium Spreaded',
  ],
  kolomRingkasanMataUang: [
    'Currency',
    'Total TSI Top Risk Spreaded',
    'Total TSI Spreaded',
    'Total Limit of Liability',
    'Total Premium Spreaded',
  ],
} as const
export const CEDANT_NB_LABEL = {
  sob: 'Source of Business',
  shareCedantType: 'Share Cedant Type',
  percentShare: '% Share RNM',
  totalTsi: 'Total TSI RNM',
  totalPremi: 'Total Premi RNM',
  kolomCeding: 'Ceding',
  kolomShare: '% Share',
} as const
export const OPSI_SHARE_CEDANT_NB: OpsiNB[] = [
  { value: '0', label: 'Gross' },
  { value: '1', label: 'Share RNM' },
]

// ---------------------------------------------------------------- format

/** Label pilihan kode; kode tak dikenal tampil apa adanya. */
export function labelKode(opsi: readonly OpsiNB[], v: string | undefined): string {
  const s = v ?? ''
  return opsi.find((o) => o.value === s)?.label ?? s
}

/**
 * Uang / angka NB: titik ribuan, koma desimal, SELALU 4 desimal seperti angka layar Claim Fac In lainnya (`tampilAngka`;
 * work owner 11-10-2026 "TAMPILAN 4 DIGIT BELAKANG KOMA" - menggantikan batas 2 / 3 desimal Loss Ratio nbfacin).
 */
export const angkaNB = (v: string | undefined) => tampilAngka(v ?? '')

const panjangBulan = (th: number, bl: number) =>
  [31, (th % 4 === 0 && th % 100 !== 0) || th % 400 === 0 ? 29 : 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31][bl - 1] ??
  31

/** Tanggal + 1 hari (teks, tanpa `Date`). */
function hariBerikut(th: number, bl: number, hr: number): [number, number, number] {
  if (hr < panjangBulan(th, bl)) return [th, bl, hr + 1]
  if (bl < 12) return [th, bl + 1, 1]
  return [th + 1, 1, 1]
}

const dua = (n: number) => String(n).padStart(2, '0')
const dmy = (th: number, bl: number, hr: number) => `${dua(hr)}/${dua(bl)}/${th}`

/**
 * Teks tanggal polis -> `dd/mm/yyyy` (TanggalDMY nbfacin). Stempel Pega `YYYYMMDDTHHMMSS.mmm GMT` (StartDateTime,
 * EndDateTime, DateOfLoss) dibaca UTC lalu digeser ke WIB (UTC+7) sebelum diambil tanggalnya - sama dengan nbfacin
 * `services/kasus.go` (data lama `T170000.000 GMT` = 00:00 WIB hari BERIKUTNYA). `YYYYMMDD` (OfferingDate) dan tanggal
 * halaman `YYYY-MM-DD` tanggal biasa. Tidak terbaca = apa adanya.
 */
export function tanggalNB(v: string | undefined): string {
  const s = (v ?? '').trim()
  if (s === '') return ''
  const utc =
    /^(\d{4})(\d{2})(\d{2})T(\d{2})(\d{2})(\d{2})(?:\.\d+)? GMT$/.exec(s) ??
    /^(\d{4})-(\d{2})-(\d{2})T(\d{2}):(\d{2})(?::(\d{2})(?:\.\d+)?)?Z$/.exec(s)
  if (utc) {
    const [th, bl, hr, jam] = [Number(utc[1]), Number(utc[2]), Number(utc[3]), Number(utc[4])]
    return jam + 7 >= 24 ? dmy(...hariBerikut(th, bl, hr)) : dmy(th, bl, hr)
  }
  const polos = /^(\d{4})-?(\d{2})-?(\d{2})(?:$|[ T])/.exec(s)
  if (polos) return dmy(Number(polos[1]), Number(polos[2]), Number(polos[3]))
  const kabel = /^(\d{2})-(\d{2})-(\d{4})$/.exec(s)
  if (kabel) return `${kabel[1]}/${kabel[2]}/${kabel[3]}`
  return s
}

/** Tahun tanggal (kolom Year Loss Record Internal). */
export function tahunNB(v: string | undefined): string {
  const t = tanggalNB(v)
  return /^\d{2}\/\d{2}\/\d{4}$/.test(t) ? t.slice(6) : ''
}

// ---------------------------------------------------------------- blok General

export interface UmumPolisNB {
  reff: string
  statusBisnis: string
  insured: string
  qq: string
  mulai: string
  penawaran: string
  jenisPolis: string
  /** Risk Scoring (`Section\Periode.xml`: `.ScoringRisk.FinalScore` + lencana `.NoteFinalScore`, tersimpan). */
  riskScoring: string
  catatanSkor: string
  kelas: string
  typeFac: string
  sob: string
  ceding: string
  grup: string
  akhir: string
  lama: string
  marketing: string
  hari: string
}

/** Kode Business status polis (Section Periode, `kodeStatusBisnis` backend) -> teks layar NB. */
const STATUS_BISNIS_NB: Record<string, string> = { '1': 'New Business' }

/** Blok General layar Inward Facultative dari halaman polis. */
export function umumNB(h: Halaman): UmumPolisNB {
  const q = (p: string) => h.nilai[Q + p] ?? ''
  const nama = daftarNB(h, CEDING_NB)
    .map((b) => b['CedingCoName'] ?? '')
    .filter((n) => n !== '')
  const status = q('StatusBusiness')
  return {
    reff: q('NoOfferSlip'),
    statusBisnis: STATUS_BISNIS_NB[status] ?? status,
    insured: q('InsuredName'),
    qq: q('QQName'),
    mulai: tanggalNB(nilaiNB(h, 'PolicyData.StartDateTime')),
    penawaran: tanggalNB(nilaiNB(h, 'PolicyData.OfferingDate')),
    jenisPolis: q('PolicyType'),
    riskScoring: nilaiNB(h, 'ScoringRisk.FinalScore'),
    catatanSkor: nilaiNB(h, 'ScoringRisk.NoteFinalScore'),
    kelas: q('BusinessName'),
    typeFac: q('TypeFacultative'),
    sob: q('SobName'),
    ceding: nama.length > 0 ? nama.join(';') : q('CedingCoName'),
    grup: q('GroupName'),
    akhir: tanggalNB(nilaiNB(h, 'PolicyData.EndDateTime')),
    lama: nilaiNB(h, 'Following'),
    marketing: q('MarketingName'),
    hari: q('EDMDay'),
  }
}

// ---------------------------------------------------------------- aritmetika eksak (BigInt)

/** Bilangan desimal eksak: nilai = n / 10^s. */
export interface Eksak {
  n: bigint
  s: number
}

export const NOL: Eksak = { n: 0n, s: 0 }

/** Teks desimal bertitik -> eksak; kosong / tak terbaca = 0 (`nolBila` nbfacin). */
export function eksak(t: string | undefined): Eksak {
  const m = /^([+-]?)(\d*)(?:\.(\d*))?$/.exec((t ?? '').trim())
  if (!m || ((m[2] ?? '') === '' && (m[3] ?? '') === '')) return NOL
  const pecahan = m[3] ?? ''
  const n = BigInt((m[2] || '0') + pecahan)
  return { n: m[1] === '-' ? -n : n, s: pecahan.length }
}

const skalakan = (a: Eksak, s: number) => a.n * 10n ** BigInt(s - a.s)

export function tambahE(a: Eksak, b: Eksak): Eksak {
  const s = Math.max(a.s, b.s)
  return { n: skalakan(a, s) + skalakan(b, s), s }
}

export const kaliE = (a: Eksak, b: Eksak): Eksak => ({ n: a.n * b.n, s: a.s + b.s })

/** a / b dibulatkan setengah-ke-atas menjauhi nol pada `skala` desimal; b = 0 -> 0 (`bagiSkala` nbfacin). */
export function bagiE(a: Eksak, b: Eksak, skala: number): Eksak {
  if (b.n === 0n) return NOL
  // a/b = (a.n / 10^a.s) / (b.n / 10^b.s); hasil x 10^skala
  const pembilang = a.n * 10n ** BigInt(skala + b.s)
  const penyebut = b.n * 10n ** BigInt(a.s)
  const negatif = pembilang < 0n !== penyebut < 0n
  const pa = pembilang < 0n ? -pembilang : pembilang
  const pb = penyebut < 0n ? -penyebut : penyebut
  const q = (2n * pa + pb) / (2n * pb)
  return { n: negatif ? -q : q, s: skala }
}

export const bulatE = (a: Eksak, skala: number) => (a.s <= skala ? a : bagiE(a, { n: 1n, s: 0 }, skala))

export function bandingE(a: Eksak, b: Eksak): number {
  const s = Math.max(a.s, b.s)
  const x = skalakan(a, s)
  const y = skalakan(b, s)
  return x === y ? 0 : x < y ? -1 : 1
}

/** Teks bertitik tanpa nol ekor. */
export function teksE(a: Eksak): string {
  const negatif = a.n < 0n
  const digit = (negatif ? -a.n : a.n).toString().padStart(a.s + 1, '0')
  const bulat = digit.slice(0, digit.length - a.s)
  const pecahan = digit.slice(digit.length - a.s).replace(/0+$/, '')
  const t = bulat + (pecahan === '' ? '' : '.' + pecahan)
  return negatif && t !== '0' ? '-' + t : t
}

export const SERATUS: Eksak = { n: 100n, s: 0 }
const SEPULUH_RIBU: Eksak = { n: 10000n, s: 0 }

/** `persenDari` nbfacin: x × y / 100 (20 desimal). */
export const persenDari = (x: Eksak, y: Eksak) => bagiE(kaliE(x, y), SERATUS, 20)

/** `simpan8` nbfacin: setengah-ke-atas 8 desimal. */
export const simpan8 = (a: Eksak) => teksE(bulatE(a, 8))

// ---------------------------------------------------------------- Object Item / Coverage (dihitung seperti nbfacin)

/** Total TSI per mata uang (`totalPerMataUang` nbfacin, `AddPaymentCurency_ACT`): hanya total > 0, urut kemunculan. */
export function totalTSIItemNB(items: readonly Baris[]): { currency: string; total: string }[] {
  const urutan: string[] = []
  for (const i of items) {
    const c = i['Currency'] ?? ''
    if (c !== '' && !urutan.includes(c)) urutan.push(c)
  }
  return urutan
    .map((c) => ({
      currency: c,
      total: jumlahDesimal(items.filter((i) => i['Currency'] === c).map((i) => i['TSIObjectItem'] ?? '')).total,
    }))
    .filter((t) => bandingE(eksak(t.total), NOL) > 0)
}

/** Total Gross Premium item = Σ Premium coverage (`totalPremiItem` nbfacin); tanpa coverage -> nilai tersimpan. */
export function premiItemNB(coverages: readonly Baris[], item: Baris): string {
  if (coverages.length === 0) return item['TotalGrossPremi'] ?? ''
  return jumlahDesimal(coverages.map((c) => c['Premium'] ?? '')).total
}

/** Total Premium RNM item = Σ PremiNusantaraRe coverage (`tampilan` spreading nbfacin; tidak disimpan). */
export function premiRNMItemNB(coverages: readonly Baris[]): string {
  const ada = coverages.map((c) => c['PremiNusantaraRe'] ?? '').filter((v) => v.trim() !== '')
  return ada.length === 0 ? '' : simpan8(ada.map(eksak).reduce(tambahE, NOL))
}

/**
 * Rate = `@if(TSI = 0, 0, @divide(Premium, TSI, 20) * 1000)` (`ratePerMil` nbfacin): bagi 20 desimal setengah-ke-atas,
 * x 1000, 8 desimal.
 */
export function ratePerMilNB(premium: string, tsi: string): string {
  const t = eksak(tsi)
  if (t.n === 0n) return '0'
  const q20 = bagiE(eksak(premium), t, 20)
  return teksE(bulatE(kaliE(q20, { n: 1000n, s: 0 }), 8))
}

export interface TotalCoverageNB {
  currency: string
  tsi: string
  premium: string
  rate: string
}

/** Grid total lokasi tab Coverage (`TotalFireLokasi` nbfacin, TotalTSIPremiGrossList): per Currency, urut kemunculan. */
export function totalCoverageLokasiNB(h: Halaman, lokasi: number): TotalCoverageNB[] {
  const out: { currency: string; tsi: string[]; premium: string[] }[] = []
  daftarNB(h, itemNB(lokasi)).forEach((it, j) => {
    const c = it['Currency'] ?? ''
    let r = out.find((x) => x.currency === c)
    if (!r) out.push((r = { currency: c, tsi: [], premium: [] }))
    r.tsi.push(it['TSIObjectItem'] ?? '')
    r.premium.push(premiItemNB(daftarNB(h, coverageNB(lokasi, j + 1)), it))
  })
  return out.map((r) => {
    const tsi = jumlahDesimal(r.tsi).total
    const premium = jumlahDesimal(r.premium).total
    return { currency: r.currency, tsi, premium, rate: ratePerMilNB(premium, tsi) }
  })
}

/** Ringkasan seluruh lokasi tab Coverage (`ringkasanCoverage` nbfacin, SumTSIPremiSpreadedRNM_FIRE_Act). */
export function ringkasanCoverageNB(h: Halaman): {
  item: { itemType: string; currency: string; tsi: string; premium: string }[]
  mataUang: TotalCoverageNB[]
} {
  const item: { itemType: string; currency: string; tsi: string[]; premium: string[] }[] = []
  const uang: { currency: string; tsi: string[]; premium: string[] }[] = []
  daftarNB(h, LOKASI_NB).forEach((_, l) => {
    daftarNB(h, itemNB(l + 1)).forEach((it, j) => {
      const premi = premiItemNB(daftarNB(h, coverageNB(l + 1, j + 1)), it)
      const jenis = it['ItemType'] ?? ''
      const c = it['Currency'] ?? ''
      let a = item.find((r) => r.itemType === jenis && r.currency === c)
      if (!a) item.push((a = { itemType: jenis, currency: c, tsi: [], premium: [] }))
      a.tsi.push(it['TSIObjectItem'] ?? '')
      a.premium.push(premi)
      let b = uang.find((r) => r.currency === c)
      if (!b) uang.push((b = { currency: c, tsi: [], premium: [] }))
      b.tsi.push(it['TSIObjectItem'] ?? '')
      b.premium.push(premi)
    })
  })
  return {
    item: item.map((r) => ({
      itemType: r.itemType,
      currency: r.currency,
      tsi: jumlahDesimal(r.tsi).total,
      premium: jumlahDesimal(r.premium).total,
    })),
    mataUang: uang.map((r) => {
      const tsi = jumlahDesimal(r.tsi).total
      const premium = jumlahDesimal(r.premium).total
      return { currency: r.currency, tsi, premium, rate: ratePerMilNB(premium, tsi) }
    }),
  }
}

// ---------------------------------------------------------------- Spreading (SumTSIPremiSpreadedRNM_FIRE_Act, nbfacin)

/** Baris total spreading (`models.TotalSpreading` nbfacin), teks bertitik 8 desimal. */
export interface TotalSpreadingNB {
  currency: string
  treatyType: string
  treatyName: string
  sharePercentage: string
  claimSpreaded: string
  tsiSpreaded: string
  claimEstimation: string
  claimAmountIDR: string
  premiumSpreaded: string
}

export interface RingkasanUangSpreadingNB {
  currency: string
  tsiTopRisk: string
  tsi: string
  lol: string
  premium: string
}

export interface TotalE {
  currency: string
  treatyType: string
  share: Eksak
  claimSpreaded: Eksak
  tsiSpreaded: Eksak
  claimEstimation: Eksak
  claimAmountIDR: Eksak
  premiumSpreaded: Eksak
}

/** `catatanTerorisme` nbfacin (SumTSIPremiSpreadedRNM_FIRE_Act .4.4). */
const CATATAN_TERORISME = 'TERRORISM & SABOTAGE'

export const cariTotal = (t: readonly TotalE[], jenis: string, uang: string) =>
  t.findIndex((x) => x.treatyType === jenis && x.currency === uang)

const keTotalNB = (t: TotalE, nama: string): TotalSpreadingNB => ({
  currency: t.currency,
  treatyType: t.treatyType,
  treatyName: nama,
  sharePercentage: simpan8(t.share),
  claimSpreaded: simpan8(t.claimSpreaded),
  tsiSpreaded: simpan8(t.tsiSpreaded),
  claimEstimation: simpan8(t.claimEstimation),
  claimAmountIDR: simpan8(t.claimAmountIDR),
  premiumSpreaded: simpan8(t.premiumSpreaded),
})

export interface HasilSpreadingNB {
  /** Total per lokasi (TotalTSIPremiSpreadRNM), sejajar LocationList. */
  lokasi: TotalSpreadingNB[][]
  /** Ringkasan per treaty (TotalSpreadingCurrency). */
  treaty: TotalSpreadingNB[]
  /** Ringkasan per mata uang (TotalSpreadAll). */
  mataUang: RingkasanUangSpreadingNB[]
}

/**
 * Port `totalLokasi` + `totalMataUang` + `ringkasan` (tanpa pemeriksaan batas treaty / kurs - itu pesan Save, bukan
 * tampilan) `modul/nbfacin/backend/services/spreading.go`. Nama treaty: nbfacin memakai master GetTreatyName; polis klaim
 * membawa TreatyName di baris SpreadingList, jadi nama = TreatyName baris pertama ber-TreatyType itu.
 */
export function hitungSpreadingNB(h: Halaman): HasilSpreadingNB {
  const ps = eksak(nilaiNB(h, 'PercentShare'))
  const lokasiBaris = daftarNB(h, LOKASI_NB)
  const adaTop = lokasiBaris.some((l) => benarNB(l['Property.IsTopRisk']))
  const namaTreaty = new Map<string, string>()
  const lokasiE: { top: boolean; total: TotalE[] }[] = []
  lokasiBaris.forEach((lb, l) => {
    const top = benarNB(lb['Property.IsTopRisk'])
    const total: TotalE[] = []
    daftarNB(h, itemNB(l + 1)).forEach((it, j) => {
      const uang = it['Currency'] ?? ''
      daftarNB(h, coverageNB(l + 1, j + 1)).forEach((c, m) => {
        const pertama = m === 0
        let pctFL = NOL
        const fl = eksak(c['FirstLoss'])
        if (bandingE(fl, NOL) > 0 && bandingE(fl, SERATUS) < 0) pctFL = fl
        const tsiLoL = persenDari(eksak(c['LimitofLiability']), ps)
        const tsiGross = persenDari(eksak(c['TSI']), ps)
        const tsiFL = bagiE(kaliE(kaliE(eksak(c['TSI']), pctFL), ps), SEPULUH_RIBU, 20)
        const teror = (c['CoverageNote'] ?? '') === CATATAN_TERORISME
        for (const b of daftarNB(h, spreadingNB(l + 1, j + 1, m + 1))) {
          const jenis = b['TreatyType'] ?? ''
          if (!namaTreaty.has(jenis) && (b['TreatyName'] ?? '') !== '') namaTreaty.set(jenis, b['TreatyName'] ?? '')
          const share = eksak(b['SharePercentage'])
          let tsiObj = NOL
          let tsiTop = NOL
          const flSpread = persenDari(tsiFL, share)
          const lolSpread = persenDari(tsiLoL, share)
          if (!adaTop || pertama) tsiObj = eksak(b['TSISpreaded'])
          if (bandingE(pctFL, NOL) > 0 && pertama) tsiObj = persenDari(tsiGross, share)
          if (pertama && top) tsiTop = eksak(b['TSISpreaded'])
          const n = cariTotal(total, jenis, uang)
          if (n >= 0) {
            const t = total[n]!
            t.premiumSpreaded = tambahE(t.premiumSpreaded, eksak(b['PremiumSpreaded']))
            // .4.8.2 (coverage pertama) dan .4.8.3 (terorisme) langkah terpisah: coverage pertama ber-catatan terorisme
            // dijumlah DUA kali (verbatim nbfacin).
            for (const kena of [pertama, teror]) {
              if (!kena) continue
              t.tsiSpreaded = tambahE(t.tsiSpreaded, tsiObj)
              t.claimEstimation = tambahE(t.claimEstimation, lolSpread)
              t.claimAmountIDR = tambahE(t.claimAmountIDR, flSpread)
            }
            if (pertama && top) t.claimSpreaded = tambahE(t.claimSpreaded, tsiTop)
            if (bandingE(share, t.share) !== 0) t.share = NOL // "Pct Share Beda"
            continue
          }
          total.push({
            currency: uang,
            treatyType: jenis,
            share,
            claimSpreaded: tsiTop,
            tsiSpreaded: tsiObj,
            claimEstimation: lolSpread,
            claimAmountIDR: flSpread,
            premiumSpreaded: eksak(b['PremiumSpreaded']),
          })
        }
      })
    })
    lokasiE.push({ top, total })
  })
  return ringkasSpreadingNB(lokasiE, namaTreaty, adaTop)
}

/** Total per lokasi (`LokasiSpreading` nbfacin: IsTopRisk + Total) bahan `ringkasSpreadingNB`. */
export interface LokasiTotalE {
  top: boolean
  total: TotalE[]
}

/**
 * Tahap `totalMataUang` + `ringkasan` nbfacin (`spreading.go`) atas total per lokasi - dipakai FIRE (`hitungSpreadingNB`) dan
 * lini lain yang hanya berbeda akumulasi per lokasinya (Aneka `totalAneka`, Cargo `totalCargo`).
 */
export function ringkasSpreadingNB(
  lokasiE: readonly LokasiTotalE[],
  namaTreaty: ReadonlyMap<string, string>,
  adaTop: boolean,
): HasilSpreadingNB {
  // totalMataUang
  const perUang: TotalE[] = []
  for (const l of lokasiE) {
    for (const t of l.total) {
      const n = cariTotal(perUang, t.treatyType, t.currency)
      if (n < 0) {
        perUang.push({ ...t, claimEstimation: adaTop && !l.top ? NOL : t.claimEstimation })
        continue
      }
      const x = perUang[n]!
      x.premiumSpreaded = tambahE(x.premiumSpreaded, t.premiumSpreaded)
      x.tsiSpreaded = tambahE(x.tsiSpreaded, t.tsiSpreaded)
      x.claimSpreaded = tambahE(x.claimSpreaded, t.claimSpreaded)
      x.claimAmountIDR = tambahE(x.claimAmountIDR, t.claimAmountIDR)
      if (!adaTop || l.top) x.claimEstimation = tambahE(x.claimEstimation, t.claimEstimation)
      if (bandingE(x.share, t.share) !== 0) x.share = NOL
    }
  }
  // ringkasan: nama treaty terbawa dari baris sebelumnya bila jenisnya tidak dikenal (19.5.5 verbatim)
  let namaSebelum = ''
  const treaty = perUang.map((t) => {
    const n = namaTreaty.get(t.treatyType) ?? ''
    if (n !== '') namaSebelum = n
    return keTotalNB(t, namaSebelum)
  })
  const semua: { currency: string; tsiTopRisk: Eksak; tsi: Eksak; lol: Eksak; premium: Eksak }[] = []
  for (const t of perUang) {
    const lol = bandingE(t.claimAmountIDR, NOL) > 0 ? t.claimAmountIDR : t.claimEstimation
    const s = semua.find((x) => x.currency === t.currency)
    if (!s) {
      semua.push({
        currency: t.currency,
        tsiTopRisk: t.claimSpreaded,
        tsi: t.tsiSpreaded,
        lol,
        premium: t.premiumSpreaded,
      })
      continue
    }
    s.tsi = tambahE(s.tsi, t.tsiSpreaded)
    s.premium = tambahE(s.premium, bulatE(t.premiumSpreaded, 4)) // 22.1 (dijumlah); baris baru 22.3 nilai utuh
    s.tsiTopRisk = tambahE(s.tsiTopRisk, t.claimSpreaded)
    s.lol = tambahE(s.lol, lol)
  }
  return {
    lokasi: lokasiE.map((l) => l.total.map((t) => keTotalNB(t, namaTreaty.get(t.treatyType) ?? ''))),
    treaty,
    mataUang: semua.map((s) => ({
      currency: s.currency,
      tsiTopRisk: simpan8(s.tsiTopRisk),
      tsi: simpan8(s.tsi),
      lol: simpan8(s.lol),
      premium: simpan8(s.premium),
    })),
  }
}
