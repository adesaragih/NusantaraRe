// Tampilan polis NB Fac In HANYA-BACA keluarga Aneka (perintah work owner 11-10-2026: "tampilan view polis yang belum
// tergenerate tolong diperbaiki. cek tampilan nya di nb facin!"). Rupa = tab Object / Coverage / Spreading Aneka layar
// Inward Facultative `modul/nbfacin/frontend/components/{TabObjectAneka,FormObjekAneka,aneka,TabCoverageAneka,
// FormCoverageAneka,TabSpreadingAneka}.tsx` DISALIN (pola, bukan impor - batas modul), dibaca dari halaman polis klaim
// (`OfferFacIn.LocationList(n).Property.RiskLocation.(OccupationList(m).)AnekaList`). Label VERBATIM
// `modul/nbfacin/frontend/labels.ts`. Berkas ini hanya fungsi murni (diuji tanpa DOM); komponennya
// `components/PolisNBAneka.tsx`.
//
// Lini dan varian tampilan = port `modul/nbfacin/backend/services/lini.go` (`LiniKasus`, `tampilanAneka`) atas When
// registry nbfacin (`services/rules/registry_gen.go`, asal korpus NB FacIn\When) - fungsi murni BusinessType /
// BusinessCode / BusinessName `OfferFacIn.QuotationData`. When klaim (`backend/models/lini.go`, korpus Claim Fac In) TIDAK
// dipakai: keduanya berbeda (NB memasukkan CustomBond ke IsAneka; IsElectronicEquipment NB membandingkan " ElectronicEquipment "
// berspasi sehingga tidak pernah benar - verbatim).
//
// Total yang nbfacin HITUNG (Grand Total per mata uang tab Coverage `totalaneka.go`, total spreading per lokasi
// `spreadinganeka.go`) dihitung ulang di sini dengan rumus yang sama, eksak tanpa float.

import { jumlahDesimal } from '../../../inti/frontend/lib/desimal'
import type { Baris, Halaman } from './api'
import type { ItemDasarNB } from './polisNBBayar'
import {
  LOKASI_NB,
  NOL,
  SERATUS,
  anakNB,
  bagiE,
  bandingE,
  benarNB,
  bulatE,
  cariTotal,
  daftarNB,
  eksak,
  kaliE,
  nilaiNB,
  persenDari,
  ringkasSpreadingNB,
  tambahE,
  teksE,
  type Eksak,
  type HasilSpreadingNB,
  type LokasiTotalE,
  type TotalE,
} from './polisNB'

// ---------------------------------------------------------------- lini (When nbfacin)

/** Masukan When lini: `.Quotation.BusinessType` / `.BusinessCode` / `.BusinessName`. */
export interface BisnisNB {
  type: string
  code: string
  name: string
}

export function bisnisNB(h: Halaman): BisnisNB {
  const q = (p: string) => (h.nilai['OfferFacIn.QuotationData.' + p] ?? '').trim()
  return { type: q('BusinessType'), code: q('BusinessCode'), name: q('BusinessName') }
}

type When = (k: BisnisNB) => boolean
const jenis =
  (...v: string[]): When =>
  (k) =>
    v.includes(k.type)
const kode =
  (...v: string[]): When =>
  (k) =>
    v.includes(k.code)
const nama =
  (...v: string[]): When =>
  (k) =>
    v.includes(k.name)
const atau =
  (...w: When[]): When =>
  (k) =>
    w.some((f) => f(k))
const dan =
  (...w: When[]): When =>
  (k) =>
    w.every((f) => f(k))

const ANEKA = jenis('Aneka')
const isHE = jenis('HE')
const isGolf = jenis('GolfInsurance')
const isGrowingTrees = jenis('GrowingTrees')
const isBondingKBG = jenis('BondingKBG')
const isBonding = atau(jenis('Bonding'), isBondingKBG)
const isCustomBonds = jenis('CustomBond')
const isBondingAndCustomBonds = atau(jenis('Bonding'), isBondingKBG, isCustomBonds)
const isAsuransiKredit = kode('10053', '10243', '10244', '10245')
const isEnvironmental = dan(kode('10169'), ANEKA)
const isCrime = dan(kode('10168'), ANEKA)
const isCitCis = kode('10040', '10041', '10103', '10104')
const isLiability = atau(
  jenis('Liability'),
  dan(ANEKA, kode('10021')), // IsProductsLiability
  dan(ANEKA, kode('10023')), // IsProfessionalLiability
  jenis('Workmen'), // IsWorkmenCompensation
  kode('10048', '10184'),
  nama('BILLBOARD/NEON SIGN'), // IsBillboardNeon
)
const isWhetherIndex = dan(ANEKA, atau(kode('10165'), nama('WHETHER INDEX')))
const isFidelity = jenis('Fidelity')
const isEar = jenis('Ear')
const isGlass = jenis('Glass')
const isCar = jenis('Car')
const isCivilEngineering = dan(ANEKA, atau(kode('10166'), nama('CIVIL ENGINEERING COMPLETED RISKS')))
const isBillboardNeonSyariah = atau(jenis('BillboardNeonSyariah'), kode('10106'))
const isBoiler = jenis('Boiler')
const isCMIAneka = dan(ANEKA, atau(kode('10187'), nama('COMPREHENSIVE MACHINERIES INSURANCE')))
const isPortRisk = kode('10188')
const isBuilderRisk = atau(kode('10085'), nama('BUILDER RISK'))
const isMarineHull = jenis('MarineHull')
const isMarineHullOffshore = dan(ANEKA, atau(kode('10157'), nama('MARINE HULL OFFSHORE')))
const isAllRisk = jenis('AllRisk')
const isBurglary = jenis('Burglary')
const isObjectWithQuantityYear = jenis('MBD', 'Boiler', 'ContractorsPM', 'LandRig')
const isAviationHull = jenis('AviationHull')
const isYieldShortfall = dan(ANEKA, kode('10167'))
/** VERBATIM nbfacin `IsElectronicEquipment` (" ElectronicEquipment " berspasi) - tidak pernah benar. */
const isElectronicEquipment = jenis(' ElectronicEquipment ')
const isEngineering = atau(jenis('Car', 'Ear', 'MBD'), kode('10166', '10032'))
const isCreditBriguna = kode('10235')
const isMBD = jenis('MBD')
const isContractorsPM = jenis('ContractorsPM')
const isCMI = nama('COMPREHENSIVE MACHINERIES INSURANCE')
const isObjectSectionAneka = kode('24', '15', '32', '17', '18', '27', '07', 'SV', 'SX', '90', '91', 'T9')
const isAnekaNotBonding = atau(
  isLiability,
  isMarineHull,
  isGrowingTrees,
  jenis('Exclusion'),
  jenis('Maintenance'),
  isCar,
  isElectronicEquipment,
  isAviationHull,
  isEar,
  isAllRisk,
  isGlass,
  isFidelity,
  isBillboardNeonSyariah,
  isBurglary,
  jenis('CIT'),
  jenis('CIS'),
  isBoiler,
  isMBD,
  ANEKA,
  isHE,
  jenis('LandRig'),
  isContractorsPM,
)
const isAneka = atau(isAnekaNotBonding, isYieldShortfall, isCrime, isEnvironmental, isBondingAndCustomBonds)
const isMarineCargo = jenis('MarineCargo')

/** `TampilanAneka` nbfacin (`lini.go` `tampilanAneka`). */
export interface TampilanAnekaNB {
  kontainer: string[]
  kolomTahunJumlah: boolean
  he: boolean
  golf: boolean
  tanamanHasil: boolean
  engineering: boolean
  bonding: boolean
  briguna: boolean
  topRisk: boolean
  sectionLol: boolean
  jumlahSatuan: boolean
  currencyCoverage: boolean
  mbd: boolean
}

/** C1..C9 - When yang di-OR-kan (`kontainerAneka` nbfacin). */
const KONTAINER_ANEKA: readonly [string, When][] = [
  [
    'C1',
    atau(
      isAsuransiKredit,
      isEnvironmental,
      isCrime,
      isGolf,
      isCitCis,
      isLiability,
      isWhetherIndex,
      isFidelity,
      isEar,
      isGlass,
      isBonding,
      isCar,
      isCivilEngineering,
      isBillboardNeonSyariah,
      isCustomBonds,
      isBoiler,
      isCMIAneka,
    ),
  ],
  ['C2', atau(isPortRisk, isBuilderRisk, isMarineHull, isMarineHullOffshore)],
  ['C3', atau(isAllRisk, isBurglary)],
  ['C4', isObjectWithQuantityYear],
  ['C5', atau(isObjectWithQuantityYear, isHE)],
  ['C6', isHE],
  ['C7', isAviationHull],
  ['C8', atau(isGrowingTrees, isYieldShortfall)],
  ['C9', isElectronicEquipment],
]

/** `BusinessCode = 20` (pembanding angka Pega). */
function samaDuaPuluh(k: string): boolean {
  return k === '20' || (k !== '' && Number(k) === 20)
}

export function tampilanAnekaNB(k: BisnisNB): TampilanAnekaNB {
  const kosong: TampilanAnekaNB = {
    kontainer: [],
    kolomTahunJumlah: false,
    he: false,
    golf: false,
    tanamanHasil: false,
    engineering: false,
    bonding: false,
    briguna: false,
    topRisk: false,
    sectionLol: false,
    jumlahSatuan: false,
    currencyCoverage: false,
    mbd: false,
  }
  if (!atau(isAnekaNotBonding, isHE, isGolf, isGrowingTrees, isBonding, isCustomBonds)(k)) return kosong
  return {
    kontainer: KONTAINER_ANEKA.filter(([, w]) => w(k)).map(([n]) => n),
    kolomTahunJumlah: isObjectWithQuantityYear(k),
    he: isHE(k),
    golf: isGolf(k),
    tanamanHasil: atau(isGrowingTrees, isYieldShortfall)(k),
    engineering: isEngineering(k),
    bonding: isBonding(k),
    briguna: isCreditBriguna(k),
    topRisk: atau(isMBD, isHE, isContractorsPM, isElectronicEquipment, isCMI, isMarineHull, isLiability)(k),
    sectionLol: isObjectSectionAneka(k),
    jumlahSatuan: !samaDuaPuluh(k.code),
    currencyCoverage: atau(
      isPortRisk,
      isBuilderRisk,
      isObjectWithQuantityYear,
      isAllRisk,
      isMarineHull,
      isMarineHullOffshore,
      isMBD,
      isHE,
      isAviationHull,
      isBurglary,
    )(k),
    mbd: isMBD(k),
  }
}

/** Lini kasus menurut layar Inward Facultative (`InwardFacultative.tsx` L329-336): FIRE > Aneka > Marine Cargo. */
export interface LiniNB {
  fire: boolean
  aneka: boolean
  /** `lini.aneka` nbfacin (IsAneka) - grid okupasi (`tampilOkupasi`). */
  anekaWhen: boolean
  cargo: boolean
  tampilan: TampilanAnekaNB
}

export function liniNB(h: Halaman, fire: boolean): LiniNB {
  const k = bisnisNB(h)
  const aneka = !fire && atau(isAnekaNotBonding, isGolf, isBondingAndCustomBonds)(k)
  return {
    fire,
    aneka,
    anekaWhen: isAneka(k),
    cargo: !fire && !aneka && isMarineCargo(k),
    tampilan: tampilanAnekaNB(k),
  }
}

/** Grid okupasi tampil (`tampilOkupasi` nbfacin): hanya Golf murni (bukan IsAneka) yang tanpa grid okupasi. */
export const tampilOkupasiNB = (anekaWhen: boolean, t: TampilanAnekaNB) => anekaWhen || !t.golf

// ---------------------------------------------------------------- jalur halaman

export const okupasiAnekaNB = (lokasi: number) => anakNB(LOKASI_NB, lokasi, 'Property.RiskLocation.OccupationList')
export const objekAnekaNB = (lokasi: number, okupasi: number) => anakNB(okupasiAnekaNB(lokasi), okupasi, 'AnekaList')
export const golfNB = (lokasi: number) => anakNB(LOKASI_NB, lokasi, 'Property.RiskLocation.AnekaList')
export const coverageAnekaNB = (objek: string, n: number) => anakNB(objek, n, 'CoverageList')

/** Objek Aneka (jalur daftar + baris) satu lokasi: objek per okupasi lalu objek Golf. */
export interface ObjekLokasiNB {
  /** Jalur daftar objek (`...AnekaList`) dan nomor baris (berbasis 1). */
  daftar: string
  n: number
  b: Baris
}

function objekOkupasi(h: Halaman, lokasi: number): ObjekLokasiNB[] {
  const out: ObjekLokasiNB[] = []
  daftarNB(h, okupasiAnekaNB(lokasi)).forEach((_, o) => {
    const d = objekAnekaNB(lokasi, o + 1)
    daftarNB(h, d).forEach((b, n) => out.push({ daftar: d, n: n + 1, b }))
  })
  return out
}

function objekGolf(h: Halaman, lokasi: number): ObjekLokasiNB[] {
  const d = golfNB(lokasi)
  return daftarNB(h, d).map((b, n) => ({ daftar: d, n: n + 1, b }))
}

// ---------------------------------------------------------------- label (VERBATIM nbfacin/labels.ts)

/** `GRID_LOKASI_ANEKA`. */
export const GRID_LOKASI_ANEKA_NB = ['No', 'Top Risk', 'Location'] as const
/** `GRID_OKUPASI_ANEKA`. */
export const GRID_OKUPASI_ANEKA_NB = ['Occupation Code', 'Occupation Name'] as const
/** `GRID_GOLF_ANEKA`. */
export const GRID_GOLF_ANEKA_NB = ['Object Name', 'Currency', 'TSI'] as const
/** `GRID_OBJEK_ANEKA`. */
export const GRID_OBJEK_ANEKA_NB = {
  objectName: 'Object Name',
  section: 'Section',
  limitOfLiability: 'Limit of Liability',
  year: 'Year',
  quantity: 'Quantity',
  unitQuantity: 'Unit Quantity',
  chassisNumber: 'Chassis Number',
  engineNumber: 'Engine Number',
  currency: 'Currency',
  tsi: 'TSI',
} as const
/** `OPSI_SECTION_ANEKA`. */
export const OPSI_SECTION_ANEKA_NB = [
  { value: '1', label: 'Section I' },
  { value: '2', label: 'Section II' },
  { value: '3', label: 'Warranty' },
  { value: '4', label: 'No Section' },
] as const
/** `FORM_OBJEK_ANEKA`. */
export const FORM_OBJEK_ANEKA_NB = {
  objectDescription: 'Object Description',
  objectName: 'Object Name',
  vesselName: 'Vessel Name',
  modelTipe: 'Model dan Tipe',
  section: 'Section',
  noKontrak: 'Contract No',
  startDate: 'Start Date',
  endDate: 'End Date',
  noKtp: 'KTP No',
  noNpwp: 'NPWP No',
  currency: 'Currency',
  tsi: 'TSI',
  vesselPrice: 'Vessel Price (TSI)',
  draftWordingId: 'Draft Wording ID',
  maintenancePeriod: 'Maintenance Period',
  others: 'Others (Note)',
  productionYear: 'Production Year',
  builtIn: 'Place of Manufacture',
  flag: 'Flag',
  quantity: 'Quantity',
  invoiceNo: 'Invoice No',
  chassisNumber: 'Chassis Number',
  engineNumber: 'Engine Number',
  serialNo: 'Serial No',
  year: 'Year',
  unit: 'Unit',
  objectNameHe: 'Object Name',
  merk: 'Merk',
  type: 'Type',
  tahun: 'Tahun',
  batasGeografis: 'Batas Geografis',
  yearOfPlanting: 'Year of Planting',
  noOfTree: 'No of Trees',
  areaHectar: 'Area ( Hectar )',
  limitIdemnityDay: 'Limit Idemnity Day',
  limitIdemnityMonth: 'Limit Idemnity Month',
  briguna: {
    loanAccountNumber: 'Loan Account No',
    loanName: 'Loan Name',
    job: 'Job',
    depositAccountNo: 'Deposit Account No',
    idCardNo: 'ID Card Number',
    institution: 'Institution',
    loanType: 'Loan Type',
    dateOfBirth: 'Date Of Birth',
    salary: 'Salary',
  },
  marineHull: {
    typeOfVessel: 'Type of Vessel',
    nrtrt: 'NRT(RT)',
    lengthm: 'Length (M)',
    classification: 'Classification',
    grtrt: 'GRT(RT)',
    breadthm: 'Breadth (M)',
    construction: 'Construction',
    dwtrt: 'DWT(RT)',
    depthm: 'Depth (M)',
  },
  aviationHull: {
    registrationMarks: 'Kode Registrasi',
    maxPassengers: 'Maksimal Penumpang',
    pilots: 'Pilot',
    standardUses: 'Penggunaan Standar',
    specialUses: 'Penggunaan Khusus',
    specialRentalUses: 'Penggunaan Khusus Rental',
  },
} as const
/** `FORM_MEMORANDUM`. */
export const FORM_MEMORANDUM_NB = {
  judul: 'Akseptasi Memorandum',
  applicationNumber: 'Permohanan Ke',
  accumulation: 'Akumulasi',
  collateral: 'Colateral',
  insuredSiup: 'S.I.U.P',
  keyPerson: 'Key Person ( PIC )',
  originalRate: 'Rate Baku ( % )',
  principalsOutGo: 'Out Go ( % )',
  typeOfWording: 'Jenis Wording',
  buildingStatus: 'Status Bangunan',
  supportingDocument: 'Dok Pendukung',
  principalsData: 'Data Prinsipal',
  conclusionsProposals: 'Kesimpulan Dan Usulan',
  principalsBackground: 'Latar Belakang',
  summaryFinancial: 'Summary Laporan Keuangan',
  principalsGuarantee: 'Penjaminan',
} as const
/** `TEKS_ANEKA` (teks layar nbfacin). */
export const TEKS_ANEKA_NB = {
  kosongLokasi: 'Belum ada lokasi.',
  kosongOkupasi: 'Belum ada okupasi.',
  kosongObjek: 'Belum ada objek.',
  tanpaKontainer: 'Form objek untuk Class of Business ini tidak tersedia di Pega (tanpa kontainer).',
  belumDiisi: '(belum diisi)',
  grup: {
    objek: 'Object',
    kapal: 'Vessel Details',
    kendaraan: 'Vehicle / Equipment Details',
    pesawat: 'Aircraft Details',
    periode: 'Period',
    nilai: 'Sum Insured',
    identitas: 'Identity',
    briguna: 'Loan Details',
    catatan: 'Note',
  },
} as const
/** `GRID_COV_LOKASI_ANEKA` / `GRID_COV_OKUPASI_ANEKA` / `GRID_COV_OBJEK_ANEKA` / `GRID_COV_TOTAL_ANEKA` / `GRID_COVERAGE_ANEKA`. */
export const GRID_COV_ANEKA_NB = {
  lokasi: ['No.', 'Location'],
  okupasi: ['Occupation Name'],
  objek: ['Object Name', 'Currency', 'TSI'],
  total: ['Currency', 'Grand Total Gross TSI', 'Grand Total Gross Premium', 'Grand Total Gross Rate'],
  coverage: ['Coverage', 'Gross Premium'],
  kosongCoverage: 'Belum ada coverage.',
} as const
/** `FORM_COVERAGE_ANEKA`. */
export const FORM_COVERAGE_ANEKA_NB = {
  coverageNote: 'Coverage Note',
  calculateMethod: 'Calculate Method',
  proRate: '% Pro Rate',
  shortPeriod: '% Short Period',
  day: 'Day',
  currency: 'Currency',
  tsi: 'TSI',
  loading: '% Loading',
  scale: '% Scale',
  rate: '% Rate',
  indemnity: 'Indemnity',
  unit: 'Indemnity Unit',
  pctIndemnity: '% Indemnity',
  premium: 'Premium',
  pctDiscount: '% Discount',
  discount: 'Discount',
  lol: 'Limit of Liability',
  pctLol: '% Limit of Liability / % Sub Limit / % First Loss Limit',
  subLimit: 'Sub Limit',
} as const
/** `OPSI_METODE_ANEKA`. */
export const OPSI_METODE_ANEKA_NB = [
  { value: '1', label: 'PRO RATA' },
  { value: '2', label: 'SHORT PERIOD' },
  { value: '3', label: 'FIX RATE' },
] as const
/** `OPSI_DAY_ANEKA`. */
export const OPSI_DAY_ANEKA_NB = ['365', '366', '360'] as const
/** `OPSI_UNIT_INDEMNITY`. */
export const OPSI_UNIT_INDEMNITY_NB = [
  { value: '0', label: 'Day (s)' },
  { value: '1', label: 'Month (s)' },
  { value: '2', label: 'Year (s)' },
] as const
/** `SPREADING_ANEKA`. */
export const SPREADING_ANEKA_NB = {
  kolomLokasi: ['No', 'Location'],
  kolomOkupasi: ['Occupation Name'],
  kolomObjek: ['Object Name', 'Currency', 'TSI'],
  kolomCoverage: ['Coverage', 'Gross Premium'],
  kolomObjekGolf: ['Object Name'],
  kolomCoverageGolf: ['Coverage Note', '% Loading', '% Rate', 'TSI', 'Premi'],
  detailCoverage: {
    tsi: 'Gross TSI',
    tsiNusantaraRe: 'TSI Nusantara Re',
    premium: 'Gross Premium',
    premiNusantaraRe: 'Premium Nusantara Re',
  },
} as const

// ---------------------------------------------------------------- form objek (FormObjekAneka nbfacin)

/** Jenis tampilan satu medan form objek. */
export type JenisMedanAneka = 'teks' | 'angka' | 'tanggal' | 'section'

/** Medan form objek: kunci baris halaman AnekaList (Pega) - `.MarineHull` / `.VehicleHE` / `.AviationHull` bertitik. */
export interface MedanAneka {
  kunci: string
  label: string
  jenis: JenisMedanAneka
  wajib?: boolean
  syarat?: 'bonding' | 'engineering' | 'jumlahSatuan'
}

const F = FORM_OBJEK_ANEKA_NB
const MH = F.marineHull
const AH = F.aviationHull

/** `KONTAINER` FormObjekAneka - urutan sel Pega (`InputDtlObjectAneka_FacIn.xml`). */
export const KONTAINER_FORM_ANEKA: Record<string, MedanAneka[]> = {
  C1: [
    { kunci: 'ObjectName', label: F.objectDescription, jenis: 'teks', wajib: true },
    { kunci: 'Section', label: F.section, jenis: 'section' },
    { kunci: 'NO_KONTRAK', label: F.noKontrak, jenis: 'teks' },
    { kunci: 'START_DATE', label: F.startDate, jenis: 'tanggal' },
    { kunci: 'END_DATE', label: F.endDate, jenis: 'tanggal' },
    { kunci: 'NO_KTP', label: F.noKtp, jenis: 'teks' },
    { kunci: 'NO_NPWP', label: F.noNpwp, jenis: 'teks' },
    { kunci: 'Currency.Name', label: F.currency, jenis: 'teks' },
    { kunci: 'TSI', label: F.tsi, jenis: 'angka' },
    { kunci: 'DraftWordingID', label: F.draftWordingId, jenis: 'teks', wajib: true, syarat: 'bonding' },
    { kunci: 'MaintenancePeriod', label: F.maintenancePeriod, jenis: 'angka', syarat: 'engineering' },
  ],
  C2: [
    { kunci: 'ObjectName', label: F.vesselName, jenis: 'teks' },
    { kunci: 'Currency.Name', label: F.currency, jenis: 'teks' },
    { kunci: 'Year', label: F.productionYear, jenis: 'teks', wajib: true },
    { kunci: 'TSI', label: F.vesselPrice, jenis: 'angka' },
    { kunci: 'BuiltIn', label: F.builtIn, jenis: 'teks' },
    { kunci: 'Flag', label: F.flag, jenis: 'teks' },
    { kunci: 'MarineHull.TypeOfVessel', label: MH.typeOfVessel, jenis: 'teks' },
    { kunci: 'MarineHull.Classification', label: MH.classification, jenis: 'teks' },
    { kunci: 'MarineHull.Construction', label: MH.construction, jenis: 'teks' },
    { kunci: 'MarineHull.NRTRT', label: MH.nrtrt, jenis: 'teks' },
    { kunci: 'MarineHull.GRTRT', label: MH.grtrt, jenis: 'teks' },
    { kunci: 'MarineHull.DWTRT', label: MH.dwtrt, jenis: 'teks' },
    { kunci: 'MarineHull.LENGTHM', label: MH.lengthm, jenis: 'teks' },
    { kunci: 'MarineHull.BREADTHM', label: MH.breadthm, jenis: 'teks' },
    { kunci: 'MarineHull.DEPTHM', label: MH.depthm, jenis: 'teks' },
  ],
  C3: [
    { kunci: 'ObjectName', label: F.objectDescription, jenis: 'teks' },
    { kunci: 'Currency.Name', label: F.currency, jenis: 'teks' },
    { kunci: 'TSI', label: F.tsi, jenis: 'angka' },
    { kunci: 'Quantity', label: F.quantity, jenis: 'teks' },
    { kunci: 'InvoiceNo', label: F.invoiceNo, jenis: 'teks' },
  ],
  C4: [
    { kunci: 'ObjectName', label: F.objectName, jenis: 'teks' },
    { kunci: 'Section', label: F.section, jenis: 'section' },
    { kunci: 'Currency.Name', label: F.currency, jenis: 'teks' },
    { kunci: 'TSI', label: F.tsi, jenis: 'angka' },
    { kunci: 'VehicleHE.ChassisNumber', label: F.chassisNumber, jenis: 'teks' },
    { kunci: 'VehicleHE.EngineNumber', label: F.engineNumber, jenis: 'teks' },
    { kunci: 'SerialNo', label: F.serialNo, jenis: 'teks' },
    { kunci: 'MaintenancePeriod', label: F.maintenancePeriod, jenis: 'angka', syarat: 'engineering' },
  ],
  C5: [
    { kunci: 'Year', label: F.year, jenis: 'teks' },
    { kunci: 'Quantity', label: F.quantity, jenis: 'teks', syarat: 'jumlahSatuan' },
    { kunci: 'UnitQuantity', label: F.unit, jenis: 'teks', syarat: 'jumlahSatuan' },
  ],
  C6: [
    { kunci: 'ObjectName', label: F.objectDescription, jenis: 'teks' },
    { kunci: 'VehicleHE.ObjectNameHE', label: F.objectNameHe, jenis: 'teks' },
    { kunci: 'Currency.Name', label: F.currency, jenis: 'teks' },
    { kunci: 'TSI', label: F.tsi, jenis: 'angka' },
    { kunci: 'VehicleHE.ChassisNumber', label: F.chassisNumber, jenis: 'teks' },
    { kunci: 'VehicleHE.EngineNumber', label: F.engineNumber, jenis: 'teks' },
    { kunci: 'VehicleHE.BrandName', label: F.merk, jenis: 'teks' },
    { kunci: 'VehicleHE.TypeName', label: F.type, jenis: 'teks' },
  ],
  C7: [
    { kunci: 'ObjectName', label: F.modelTipe, jenis: 'teks' },
    { kunci: 'Year', label: F.tahun, jenis: 'teks' },
    { kunci: 'GeographicalLimits', label: F.batasGeografis, jenis: 'teks' },
    { kunci: 'AviationHull.RegistrationMarks', label: AH.registrationMarks, jenis: 'teks' },
    { kunci: 'AviationHull.MaxPassengers', label: AH.maxPassengers, jenis: 'teks' },
    { kunci: 'AviationHull.Pilots', label: AH.pilots, jenis: 'teks' },
    { kunci: 'AviationHull.StandardUses', label: AH.standardUses, jenis: 'teks' },
    { kunci: 'AviationHull.SpecialUses', label: AH.specialUses, jenis: 'teks' },
    { kunci: 'AviationHull.SpecialRentalUses', label: AH.specialRentalUses, jenis: 'teks' },
  ],
  C8: [
    { kunci: 'ObjectName', label: F.objectName, jenis: 'teks' },
    { kunci: 'Currency.Name', label: F.currency, jenis: 'teks' },
    { kunci: 'TSI', label: F.tsi, jenis: 'angka' },
    { kunci: 'Year', label: F.yearOfPlanting, jenis: 'teks' },
    { kunci: 'NoOfTree', label: F.noOfTree, jenis: 'angka' },
    { kunci: 'AreaHectar', label: F.areaHectar, jenis: 'angka' },
  ],
  C9: [
    { kunci: 'ObjectName', label: F.objectName, jenis: 'teks' },
    { kunci: 'Section', label: F.section, jenis: 'section' },
    { kunci: 'Quantity', label: F.quantity, jenis: 'teks' },
    { kunci: 'Currency.Name', label: F.currency, jenis: 'teks' },
    { kunci: 'TSI', label: F.tsi, jenis: 'angka' },
    { kunci: 'Year', label: F.year, jenis: 'teks' },
    { kunci: 'LimitIdemnityDay', label: F.limitIdemnityDay, jenis: 'angka' },
    { kunci: 'LimitIdemnityMonth', label: F.limitIdemnityMonth, jenis: 'angka' },
  ],
}

/** Blok `IsCreditBriguna` (C1). */
export const MEDAN_BRIGUNA_NB: MedanAneka[] = [
  { kunci: 'LoanAccountNumber', label: F.briguna.loanAccountNumber, jenis: 'teks' },
  { kunci: 'LoanName', label: F.briguna.loanName, jenis: 'teks' },
  { kunci: 'Job', label: F.briguna.job, jenis: 'teks' },
  { kunci: 'DepositAccountNo', label: F.briguna.depositAccountNo, jenis: 'teks' },
  { kunci: 'IDCardNo', label: F.briguna.idCardNo, jenis: 'teks' },
  { kunci: 'Institution', label: F.briguna.institution, jenis: 'teks' },
  { kunci: 'LoanType', label: F.briguna.loanType, jenis: 'teks' },
  { kunci: 'DateOfBirth', label: F.briguna.dateOfBirth, jenis: 'tanggal' },
  { kunci: 'Salary', label: F.briguna.salary, jenis: 'angka' },
]

const M = FORM_MEMORANDUM_NB
/** Sub-form `InputDtlMemorandumAccept` (C1, Bonding). */
export const MEDAN_MEMORANDUM_NB: MedanAneka[] = (
  [
    'applicationNumber',
    'accumulation',
    'collateral',
    'insuredSiup',
    'keyPerson',
    'originalRate',
    'principalsOutGo',
    'typeOfWording',
    'buildingStatus',
    'supportingDocument',
    'principalsData',
    'conclusionsProposals',
    'principalsBackground',
    'summaryFinancial',
    'principalsGuarantee',
  ] as const
).map((k) => ({ kunci: 'MemorandumAccept.' + k.charAt(0).toUpperCase() + k.slice(1), label: M[k], jenis: 'teks' }))

/** Medan C1 "Others (Note)". */
export const MEDAN_CATATAN_NB: MedanAneka = { kunci: 'Others', label: F.others, jenis: 'teks' }

export type GrupAneka = keyof typeof TEKS_ANEKA_NB.grup
/** Kartu lebar (bertumpuk) dan kartu kecil (sebaris) - `GRUP_BESAR` / `GRUP_KECIL` nbfacin. */
export const GRUP_BESAR_ANEKA: readonly GrupAneka[] = ['objek', 'kapal', 'kendaraan', 'pesawat']
export const GRUP_KECIL_ANEKA: readonly GrupAneka[] = ['nilai', 'identitas', 'periode']

/** Medan C2 yang nbfacin tampilkan baca-saja (jenis 'tampil') - masuk kartu Vessel Details. */
const KAPAL_TAMPIL = new Set(['Year', 'BuiltIn', 'Flag'])

/** `grupDari` nbfacin. */
export function grupMedanAneka(m: MedanAneka, kontainer: string): GrupAneka {
  const induk = m.kunci.includes('.') && m.kunci !== 'Currency.Name' ? m.kunci.split('.')[0] : ''
  if (induk === 'MarineHull' || (kontainer === 'C2' && KAPAL_TAMPIL.has(m.kunci))) return 'kapal'
  if (induk === 'VehicleHE') return 'kendaraan'
  if (induk === 'AviationHull') return 'pesawat'
  if (m.kunci === 'START_DATE' || m.kunci === 'END_DATE') return 'periode'
  if (m.kunci === 'Currency.Name' || m.kunci === 'TSI') return 'nilai'
  if (m.kunci === 'NO_KTP' || m.kunci === 'NO_NPWP') return 'identitas'
  return 'objek'
}

/**
 * Medan unik yang tampil per kartu (`medanTampil` + `kelompokkan` nbfacin): urut kontainer, satu kunci sekali (kontainer
 * pertama menang), syarat dari `tampilanAneka`; Object Name / Description selalu pertama di kartu Object.
 */
export function kelompokMedanAneka(t: TampilanAnekaNB): Record<GrupAneka, MedanAneka[]> {
  const hasil = Object.fromEntries(Object.keys(TEKS_ANEKA_NB.grup).map((g) => [g, [] as MedanAneka[]])) as Record<
    GrupAneka,
    MedanAneka[]
  >
  const sudah = new Set<string>()
  for (const k of ['C1', 'C2', 'C3', 'C4', 'C5', 'C6', 'C7', 'C8', 'C9']) {
    if (!t.kontainer.includes(k)) continue
    for (const m of KONTAINER_FORM_ANEKA[k] ?? []) {
      if (m.syarat !== undefined && !t[m.syarat]) continue
      if (sudah.has(m.kunci)) continue
      sudah.add(m.kunci)
      hasil[grupMedanAneka(m, k)].push(m)
    }
  }
  hasil.objek.sort((a, b) => Number(b.kunci === 'ObjectName') - Number(a.kunci === 'ObjectName'))
  return hasil
}

/** Kolom grid objek okupasi (`kolomGridObjek` nbfacin): K4 tanaman, K3 HE, K2 lainnya. */
export function kolomObjekAnekaNB(t: TampilanAnekaNB): (keyof typeof GRID_OBJEK_ANEKA_NB)[] {
  if (t.tanamanHasil) return ['objectName', 'currency', 'tsi']
  const sl: (keyof typeof GRID_OBJEK_ANEKA_NB)[] = t.sectionLol ? ['section', 'limitOfLiability'] : []
  if (t.he) return ['objectName', ...sl, 'year', 'chassisNumber', 'engineNumber']
  return ['objectName', ...sl, ...(t.kolomTahunJumlah ? (['year', 'quantity', 'unitQuantity'] as const) : [])]
}

/** Kunci baris halaman per kolom grid objek. */
export const KUNCI_KOLOM_OBJEK_ANEKA: Record<keyof typeof GRID_OBJEK_ANEKA_NB, string> = {
  objectName: 'ObjectName',
  section: 'Section',
  limitOfLiability: 'LimitOfLiability',
  year: 'Year',
  quantity: 'Quantity',
  unitQuantity: 'UnitQuantity',
  chassisNumber: 'VehicleHE.ChassisNumber',
  engineNumber: 'VehicleHE.EngineNumber',
  currency: 'Currency.Name',
  tsi: 'TSI',
}

// ---------------------------------------------------------------- total Coverage (totalaneka.go nbfacin)

export interface TotalCoverageAnekaNB {
  currency: string
  tsi: string
  premium: string
  rate: string
}

/** Rate persen = `@if(TSI = 0, 0, @divide(Premium, TSI, 20) * 100)`, 8 desimal (`akumulatorTotal.hasil`). */
function ratePersen(premium: Eksak, tsi: Eksak): string {
  if (tsi.n === 0n) return '0'
  return teksE(bulatE(kaliE(bagiE(premium, tsi, 20), SERATUS), 8))
}

/** 10 When langkah 3b `SumTotalTSIPremiGross_Act` (`khususTotalAneka`): TSI dan premi per COVERAGE, mata uang coverage. */
const khususTotalAneka = atau(
  isObjectWithQuantityYear,
  isAllRisk,
  isMarineHull,
  isMarineHullOffshore,
  isMBD,
  isHE,
  isAviationHull,
  isBurglary,
  isPortRisk,
  isBuilderRisk,
)

/** Grand Total per mata uang satu lokasi tab Coverage Aneka (`totalAnekaLokasi`): urut kemunculan; premi mentah. */
export function totalCoverageAnekaNB(h: Halaman, lokasi: number): TotalCoverageAnekaNB[] {
  const k = bisnisNB(h)
  const baris: { currency: string; tsi: Eksak; premium: Eksak }[] = []
  const ke = (c: string) => {
    let r = baris.find((x) => x.currency === c)
    if (!r) baris.push((r = { currency: c, tsi: NOL, premium: NOL }))
    return r
  }
  const tambahObjek = (o: ObjekLokasiNB, perCoverage: boolean) => {
    const cov = daftarNB(h, coverageAnekaNB(o.daftar, o.n))
    if (!perCoverage) {
      const t = ke(o.b['Currency.Name'] ?? '')
      t.tsi = tambahE(t.tsi, eksak(o.b['TSI']))
      for (const c of cov) t.premium = tambahE(t.premium, eksak(c['Premium']))
      return
    }
    for (const c of cov) {
      const t = ke(c['Currency.Name'] ?? '')
      t.tsi = tambahE(t.tsi, eksak(c['TSI']))
      t.premium = tambahE(t.premium, eksak(c['Premium']))
    }
  }
  if (atau(isAneka, isBondingAndCustomBonds)(k)) {
    const khusus = khususTotalAneka(k)
    for (const o of objekOkupasi(h, lokasi)) tambahObjek(o, khusus)
  } else if (isGolf(k)) {
    for (const o of objekGolf(h, lokasi)) tambahObjek(o, false)
  }
  return baris.map((r) => ({
    currency: r.currency,
    tsi: teksE(bulatE(r.tsi, 8)),
    premium: teksE(bulatE(r.premium, 8)),
    rate: ratePersen(r.premium, r.tsi),
  }))
}

/** Total TSI objek-objek (penjumlah desimal inti) - dipakai sel TSI daftar. */
export const jumlahTSI = (bs: readonly Baris[]) => jumlahDesimal(bs.map((b) => b['TSI'] ?? '')).total

// ---------------------------------------------------------------- Spreading (spreadinganeka.go nbfacin)

const liniMataUangCoverage24 = atau(
  isMBD,
  isHE,
  isBurglary,
  isObjectWithQuantityYear,
  isLiability,
  isAllRisk,
  isMarineHull,
  isMarineHullOffshore,
  isAviationHull,
  isBuilderRisk,
  isPortRisk,
)
const liniMataUangCoverage25 = atau(liniMataUangCoverage24, isElectronicEquipment, isCMI)
const liniTopRisk25 = atau(
  isMBD,
  isObjectWithQuantityYear,
  isLiability,
  isElectronicEquipment,
  isCMI,
  isMarineHull,
  isHE,
)

/** `kunciMataUang`: mata uang coverage bila lini ber-mata uang coverage dan terisi, selain itu mata uang objek. */
const kunciMataUang = (o: Baris, c: Baris, perCoverage: boolean) =>
  perCoverage && (c['Currency.Name'] ?? '').trim() !== '' ? (c['Currency.Name'] ?? '') : (o['Currency.Name'] ?? '')

/**
 * Port `totalLokasiAneka` + `totalLokasiGolf` + `totalAneka` nbfacin, lalu `ringkasSpreadingNB` (totalMataUang + ringkasan
 * FIRE). Aneka: premi selalu dijumlah; TSI hanya coverage pertama (baris BARU memakai TSISpreaded baris berjalan); LoL =
 * LimitofLiability x PS / 100 coverage pertama ber-LoL > 0, spread x share / 100 (tanpa top risk: semua lokasi; dengan top
 * risk: hanya lokasi top risk); TSI top risk hanya cabang ber-top-risk, coverage pertama, lokasi top risk, lini liniTopRisk25;
 * share beda -> 0. Golf: TSI dan premi setiap coverage, mata uang objek.
 */
export function hitungSpreadingAnekaNB(h: Halaman): HasilSpreadingNB {
  const k = bisnisNB(h)
  const ps = eksak(nilaiNB(h, 'PercentShare'))
  const lokasi = daftarNB(h, LOKASI_NB)
  const adaTop = lokasi.some((l) => benarNB(l['Property.IsTopRisk']))
  const perCoverage = adaTop ? liniMataUangCoverage25(k) : liniMataUangCoverage24(k)
  const liniTop = adaTop && liniTopRisk25(k)
  const aneka = atau(isAneka, isBondingAndCustomBonds)(k)
  const golf = isGolf(k)
  const namaTreaty = new Map<string, string>()
  const catatNama = (b: Baris) => {
    const j = b['TreatyType'] ?? ''
    if (!namaTreaty.has(j) && (b['TreatyName'] ?? '') !== '') namaTreaty.set(j, b['TreatyName'] ?? '')
  }
  const lokasiE: LokasiTotalE[] = lokasi.map((lb, i) => {
    const top = benarNB(lb['Property.IsTopRisk'])
    const total: TotalE[] = []
    const baru = (uang: string, b: Baris, tsi: Eksak, lol: Eksak, tsiTop: Eksak): TotalE => ({
      currency: uang,
      treatyType: b['TreatyType'] ?? '',
      share: eksak(b['SharePercentage']),
      claimSpreaded: tsiTop,
      tsiSpreaded: tsi,
      claimEstimation: lol,
      claimAmountIDR: NOL,
      premiumSpreaded: eksak(b['PremiumSpreaded']),
    })
    if (aneka) {
      for (const o of objekOkupasi(h, i + 1)) {
        daftarNB(h, coverageAnekaNB(o.daftar, o.n)).forEach((c, m) => {
          const pertama = m === 0
          const uang = kunciMataUang(o.b, c, perCoverage)
          let tsiLoL = NOL
          const lolCov = eksak(c['LimitofLiability'])
          if (pertama && bandingE(lolCov, NOL) > 0) tsiLoL = persenDari(lolCov, ps)
          for (const b of daftarNB(h, anakNB(coverageAnekaNB(o.daftar, o.n), m + 1, 'SpreadingList'))) {
            catatNama(b)
            const share = eksak(b['SharePercentage'])
            let tsiObj = NOL
            let lolSpread = NOL
            let tsiTop = NOL
            if (!adaTop) {
              tsiObj = eksak(b['TSISpreaded'])
              lolSpread = persenDari(tsiLoL, share)
            } else {
              if (pertama) tsiObj = eksak(b['TSISpreaded'])
              if (top) lolSpread = persenDari(tsiLoL, share)
              if (pertama && top && liniTop) tsiTop = eksak(b['TSISpreaded'])
            }
            const n = cariTotal(total, b['TreatyType'] ?? '', uang)
            if (n >= 0) {
              const x = total[n]!
              x.premiumSpreaded = tambahE(x.premiumSpreaded, eksak(b['PremiumSpreaded']))
              if (pertama) {
                x.tsiSpreaded = tambahE(x.tsiSpreaded, tsiObj)
                x.claimEstimation = tambahE(x.claimEstimation, lolSpread)
              }
              if (adaTop && top && liniTop) x.claimSpreaded = tambahE(x.claimSpreaded, tsiTop)
              if (bandingE(share, x.share) !== 0) x.share = NOL
              continue
            }
            // baris BARU: TSISpreaded baris berjalan (verbatim), LoL spread cabangnya
            total.push(baru(uang, b, eksak(b['TSISpreaded']), lolSpread, adaTop && top && liniTop ? tsiTop : NOL))
          }
        })
      }
    }
    // Golf = daftar total TERPISAH (`totalAneka`: append(aneka[i], golf[i]...)) - tidak dilebur dengan total Aneka.
    const totalGolf: TotalE[] = []
    if (golf) {
      for (const o of objekGolf(h, i + 1)) {
        const uang = o.b['Currency.Name'] ?? ''
        daftarNB(h, coverageAnekaNB(o.daftar, o.n)).forEach((_, m) => {
          for (const b of daftarNB(h, anakNB(coverageAnekaNB(o.daftar, o.n), m + 1, 'SpreadingList'))) {
            catatNama(b)
            const n = cariTotal(totalGolf, b['TreatyType'] ?? '', uang)
            if (n >= 0) {
              const x = totalGolf[n]!
              x.tsiSpreaded = tambahE(x.tsiSpreaded, eksak(b['TSISpreaded']))
              x.premiumSpreaded = tambahE(x.premiumSpreaded, eksak(b['PremiumSpreaded']))
              if (bandingE(eksak(b['SharePercentage']), x.share) !== 0) x.share = NOL
              continue
            }
            totalGolf.push(baru(uang, b, eksak(b['TSISpreaded']), NOL, NOL))
          }
        })
      }
    }
    return { top, total: [...total, ...totalGolf] }
  })
  return ringkasSpreadingNB(lokasiE, namaTreaty, adaTop)
}

// ---------------------------------------------------------------- dasar Cedant (kasusSpreadingAneka nbfacin)

const liniMataUangCount = atau(
  isMBD,
  isBurglary,
  isAviationHull,
  isHE,
  isMarineHullOffshore,
  isMarineHull,
  isAllRisk,
  isObjectWithQuantityYear,
  isBuilderRisk,
  isPortRisk,
)
const liniTSICoverageCount = atau(
  isMBD,
  isObjectWithQuantityYear,
  isAllRisk,
  isMarineHull,
  isMarineHullOffshore,
  isHE,
  isAviationHull,
  isBurglary,
)

/**
 * Pseudo-item `ringkasanRNM` keluarga Aneka (`kasusSpreadingAneka`): lini liniMataUangCount = satu item per coverage (mata
 * uang coverage, kosong -> objek; TSI coverage hanya liniTSICoverageCount), lini lain dan Golf = satu item per objek.
 */
export function itemDasarAnekaNB(h: Halaman): ItemDasarNB[] {
  const k = bisnisNB(h)
  const out: ItemDasarNB[] = []
  const cov = (o: ObjekLokasiNB) => daftarNB(h, coverageAnekaNB(o.daftar, o.n))
  const covDasar = (cs: readonly Baris[]) =>
    cs.map((c) => ({ premium: c['Premium'] ?? '', discount: c['Discount'] ?? '' }))
  daftarNB(h, LOKASI_NB).forEach((_, i) => {
    if (atau(isAneka, isBondingAndCustomBonds)(k)) {
      for (const o of objekOkupasi(h, i + 1)) {
        if (!liniMataUangCount(k)) {
          out.push({ currency: o.b['Currency.Name'] ?? '', tsi: o.b['TSI'] ?? '', coverage: covDasar(cov(o)) })
          continue
        }
        for (const c of cov(o)) {
          out.push({
            currency: kunciMataUang(o.b, c, true),
            tsi: liniTSICoverageCount(k) ? (c['TSI'] ?? '') : '',
            coverage: covDasar([c]),
          })
        }
      }
    }
    if (isGolf(k)) {
      for (const o of objekGolf(h, i + 1)) {
        out.push({ currency: o.b['Currency.Name'] ?? '', tsi: o.b['TSI'] ?? '', coverage: covDasar(cov(o)) })
      }
    }
  })
  return out
}
