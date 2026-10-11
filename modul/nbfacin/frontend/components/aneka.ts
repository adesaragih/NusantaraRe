// Fungsi murni tab Object keluarga Aneka (tiket 73 P1) - dipakai `TabObjectAneka` dan `FormObjekAneka`, diuji terpisah.
//
// Sumber: `Section\InputDtlObjectAneka_FacIn.xml` (form objek, 9 kontainer), `InputDtlObjectLocation_FacIn.xml` (form
// lokasi), `InputOkupasiAneka_FacIn.xml` (grid objek per okupasi), `DataTransform\SetObjectDescHE_DT.xml`, lookup kapal
// `BrowseShip_RD` (pemetaan SET `InputDtlObjectAneka_FacIn.xml` L10480-10570). Ringkasan: docs/issues/73 bagian S1-S2.

import type {
  AlamatBaru,
  AviationHullAneka,
  BarisKapal,
  BarisRisk,
  BrigunaAneka,
  LokasiAneka,
  MarineHullAneka,
  MemorandumAneka,
  ObjekAneka,
  ObjekFire,
  SaringRisk,
  TampilanAneka,
  VehicleHeAneka,
} from '../api'
import { kosongkanRisk, objekBaru, saringDari, terapkanAlamatBaru, terapkanRisk } from './TabObject'

const brigunaKosong = (): BrigunaAneka => ({
  loanAccountNumber: '', loanName: '', job: '', depositAccountNo: '', idCardNo: '', institution: '', loanType: '',
  dateOfBirth: '', salary: '',
})

const marineHullKosong = (): MarineHullAneka => ({
  typeOfVessel: '', nrtrt: '', lengthm: '', classification: '', grtrt: '', breadthm: '', construction: '', dwtrt: '', depthm: '',
})

const aviationHullKosong = (): AviationHullAneka => ({
  registrationMarks: '', maxPassengers: '', pilots: '', standardUses: '', specialUses: '', specialRentalUses: '',
})

const vehicleHeKosong = (): VehicleHeAneka => ({
  objectNameHe: '', objectIdHe: '', brandName: '', brand: '', typeName: '', type: '', chassisNumber: '', engineNumber: '',
})

const memorandumKosong = (): MemorandumAneka => ({
  suretyType: '', applicationNumber: '', accumulation: '', collateral: '', insuredSiup: '', keyPerson: '', originalRate: '',
  principalsOutGo: '', typeOfWording: '', buildingStatus: '', supportingDocument: '', principalsData: '',
  conclusionsProposals: '', principalsBackground: '', summaryFinancial: '', principalsGuarantee: '',
})

/** Objek kosong untuk Add (addRow). TSI awal "0" - sel TSI L2995 `pyDefaultValue` 0. */
export function objekAnekaBaru(): ObjekAneka {
  return {
    objectName: '', section: '', noKontrak: '', startDate: '', endDate: '', noKtp: '', noNpwp: '', currency: '', tsi: '0',
    tsiOld: '', draftWordingId: '', maintenancePeriod: '', others: '', year: '', quantity: '', unitQuantity: '', serialNo: '',
    builtIn: '', flag: '', geographicalLimits: '', invoiceNo: '', noOfTree: '', areaHectar: '', limitIdemnityDay: '',
    limitIdemnityMonth: '', limitOfLiability: '', almRiskId: '',
    briguna: brigunaKosong(),
    marineHull: marineHullKosong(),
    aviationHull: aviationHullKosong(),
    vehicleHe: vehicleHeKosong(),
    memorandum: memorandumKosong(),
  }
}

/** Objek dari server dilengkapi medan yang absen, supaya layar tidak runtuh pada data yang tidak lengkap. */
export function rapikanObjekAneka(d: Partial<ObjekAneka>): ObjekAneka {
  const awal = objekAnekaBaru()
  return {
    ...awal,
    ...d,
    briguna: { ...awal.briguna, ...d.briguna },
    marineHull: { ...awal.marineHull, ...d.marineHull },
    aviationHull: { ...awal.aviationHull, ...d.aviationHull },
    vehicleHe: { ...awal.vehicleHe, ...d.vehicleHe },
    memorandum: { ...awal.memorandum, ...d.memorandum },
  }
}

/** Lokasi kosong untuk Add; No = nomor urut baris (`pxListSubscript`, pola FIRE). */
export function lokasiAnekaBaru(objectNo: string): LokasiAneka {
  return {
    objectNo, isTopRisk: false, roadType: '', roadName: '', buildingNo: '', zipCode: '', country: '', riskLocation: '',
    territory: '', city: '', district: '', province: '', riskAddressId: '', okupasi: [], objekGolf: [],
  }
}

export function rapikanLokasiAneka(d: Partial<LokasiAneka>): LokasiAneka {
  return {
    ...lokasiAnekaBaru(d.objectNo ?? ''),
    ...d,
    okupasi: (d.okupasi ?? []).map((o) => ({ ...o, objek: (o.objek ?? []).map(rapikanObjekAneka) })),
    objekGolf: (d.objekGolf ?? []).map(rapikanObjekAneka),
  }
}

// Blok alamat lokasi Aneka = blok Risk Address FIRE (`InputDtlObjectLocation_FacIn.xml` L7166-11143 sama dengan
// `ObjectDetails`): aturan Choose / Clear Risk Address FIRE dipakai ulang lewat ObjekFire sementara, tanpa mengubah kode FIRE.
const KUNCI_ALAMAT = [
  'roadType', 'roadName', 'buildingNo', 'zipCode', 'country', 'riskLocation', 'territory', 'city', 'district', 'province', 'riskAddressId',
] as const

type AlamatLokasi = Pick<LokasiAneka, (typeof KUNCI_ALAMAT)[number]>

function alamatDari(o: AlamatLokasi): AlamatLokasi {
  return Object.fromEntries(KUNCI_ALAMAT.map((k) => [k, o[k]])) as AlamatLokasi
}

function keObjekFire(l: LokasiAneka): ObjekFire {
  return { ...objekBaru(l.objectNo), ...alamatDari(l) }
}

/** Choose Risk Address - `SetRiskIdDT_FacIn` (sama dengan FIRE `terapkanRisk`). */
export function terapkanRiskAneka(l: LokasiAneka, b: BarisRisk): LokasiAneka {
  return { ...l, ...alamatDari(terapkanRisk(keObjekFire(l), b)) }
}

/** Alamat baru dari popup Add (tiket 37) - sama dengan FIRE `terapkanAlamatBaru` (Building No. tidak disentuh). */
export function terapkanAlamatBaruAneka(l: LokasiAneka, a: AlamatBaru, id: string): LokasiAneka {
  return { ...l, ...alamatDari(terapkanAlamatBaru(keObjekFire(l), a, id)) }
}

/** Clear Risk Address - `ClearRiskLocation_act` (sama dengan FIRE `kosongkanRisk`; Building No. tidak ikut). */
export function kosongkanRiskAneka(l: LokasiAneka): LokasiAneka {
  return { ...l, ...alamatDari(kosongkanRisk(keObjekFire(l))) }
}

/** Saringan awal popup Risk Address (sama dengan FIRE `saringDari`). */
export function saringDariAneka(l: LokasiAneka): SaringRisk {
  return saringDari(keObjekFire(l))
}

/**
 * Pilih kapal (`BrowseShip_RD`, SET L10480-10570): Name → ObjectName, ID → AlmRiskID, FlagName → Flag, YearMake1 → Year,
 * NRT / GRT / DWT / Classification / Construction → `.MarineHull.*`. Medan lain tidak disentuh.
 */
export function terapkanKapal(o: ObjekAneka, k: BarisKapal): ObjekAneka {
  return {
    ...o,
    objectName: k.nama,
    almRiskId: k.id,
    flag: k.flag,
    year: k.tahun,
    marineHull: { ...o.marineHull, nrtrt: k.nrt, grtrt: k.grt, dwtrt: k.dwt, classification: k.klasifikasi, construction: k.konstruksi },
  }
}

/**
 * `SetObjectDescHE_DT` L155 (dipicu on-change Object Name HE / Merk / Type):
 * `.ObjectName = .VehicleHE.ObjectNameHE + " " + @if(PropertyHasValue(TypeName), TypeName, "") +
 * @if(PropertyHasValue(BrandName), " Merk : " + BrandName, "")`. Spasi sesudah ObjectNameHE selalu ada (persis Pega).
 */
export function deskripsiHe(v: VehicleHeAneka): string {
  const ada = (s: string) => s.trim() !== ''
  return v.objectNameHe + ' ' + (ada(v.typeName) ? v.typeName : '') + (ada(v.brandName) ? ' Merk : ' + v.brandName : '')
}

/** Ubah medan HE lalu jalankan `SetObjectDescHE_DT` (on-change Pega). */
export function ubahHe(o: ObjekAneka, v: Partial<VehicleHeAneka>): ObjekAneka {
  const vehicleHe = { ...o.vehicleHe, ...v }
  return { ...o, vehicleHe, objectName: deskripsiHe(vehicleHe) }
}

/** Grid objek di detail okupasi: K3 (HE), K4 (Growing Trees / Yield Shortfall), atau K2 (lainnya) - `InputOkupasiAneka_FacIn.xml`. */
export type VarianGridObjek = 'K2' | 'K3' | 'K4'

export function varianGridObjek(t: TampilanAneka): VarianGridObjek {
  if (t.tanamanHasil) return 'K4'
  return t.he ? 'K3' : 'K2'
}

/** Kolom grid objek (kunci `GRID_OBJEK_ANEKA`) per varian; Section / Limit of Liability hanya bila `sectionLol` (K73-4). */
export function kolomGridObjek(t: TampilanAneka): (keyof ObjekAneka | 'chassisNumber' | 'engineNumber')[] {
  const v = varianGridObjek(t)
  if (v === 'K4') return ['objectName', 'currency', 'tsi']
  const sl: (keyof ObjekAneka)[] = t.sectionLol ? ['section', 'limitOfLiability'] : []
  if (v === 'K3') return ['objectName', ...sl, 'year', 'chassisNumber', 'engineNumber']
  // K2: Year / Quantity / Unit Quantity bila `IsObjectWithQuantityYear || IsHE` (K2 sudah !IsHE).
  return ['objectName', ...sl, ...(t.kolomTahunJumlah ? (['year', 'quantity', 'unitQuantity'] as const) : [])]
}

/** Nilai satu sel grid objek. */
export function nilaiKolom(o: ObjekAneka, k: string): string {
  if (k === 'chassisNumber' || k === 'engineNumber') return o.vehicleHe[k]
  const v = (o as unknown as Record<string, unknown>)[k]
  return typeof v === 'string' ? v : ''
}

/**
 * Grid okupasi tampil - kontainer L15288 `IsAneka || IsGrowingTrees || IsMBD || … || IsHE || IsMarineHull || IsLiability || …`:
 * untuk case yang lolos gerbang tab Object Aneka, hanya Golf murni (isGolfInsurance, bukan IsAneka) yang tanpa grid okupasi.
 */
export function tampilOkupasi(aneka: boolean, t: TampilanAneka): boolean {
  return aneka || !t.golf
}
