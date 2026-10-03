// API modul NB FacIn - tiket 20/21.

import { minta } from '../../../inti/frontend/klien'

/** Isian layar coverage kargo yang menjadi masukan rumus (butir 61). */
export interface IsianPremiCargo {
  /** `.Currency.Name`. */
  mataUang: string
  /** `.Rate` - persen, teks desimal bertitik apa adanya. */
  rate: string
  /** `.TSI` - teks desimal bertitik apa adanya. */
  tsi: string
}

/** Jawaban `POST /api/nbfacin/premi`. Premi berupa TEKS desimal - tanpa float. */
export interface HasilPremi {
  premi: string
  mataUang: string
  asalRumus: string
}

/**
 * Badan permintaan hitung premi MARINE CARGO. Angka TIDAK diubah menjadi `number`:
 * dikirim sebagai teks, diurai desimal di backend (CLAUDE.md §7). Master policy tidak
 * ada di blok section ini → selalu `false`.
 */
export function badanPremiCargo(isian: IsianPremiCargo) {
  return {
    liniBisnis: 'MARINE CARGO',
    mataUang: isian.mataUang.trim(),
    tsi: isian.tsi.trim(),
    rate: isian.rate.trim(),
    masterPolicy: false,
  }
}

export function hitungPremiCargo(isian: IsianPremiCargo): Promise<HasilPremi> {
  return minta<HasilPremi>('/api/nbfacin/premi', { metode: 'POST', badan: badanPremiCargo(isian) })
}

/** Satu baris `T_M_ACCOUNT` (popup ChooseAccount, tiket 26 C-8). Teks apa adanya. */
export interface BarisAccount {
  id: string
  insuredId: string
  insuredName: string
  groupBusinessId: string
  groupBusiness: string
}

/** Jawaban `GET /api/nbfacin/account` - kontrak yang disepakati dengan backend (tiket 27, sesi c3). */
export interface HalamanAccount {
  baris: BarisAccount[]
  total: number
  halaman: number
  ukuran: number
}

/** Cari account: `cari` "mengandung" (jawaban work owner 02-10-2026); `halaman` mulai 1. */
export function cariAccount(cari: string, halaman: number): Promise<HalamanAccount> {
  return minta<HalamanAccount>('/api/nbfacin/account', { kueri: { cari: cari.trim(), halaman } })
}

/** Satu baris `BUSINESS` untuk saran Class Of Business (tiket 26 C-11 / tiket 28). */
export interface BarisClassOfBusiness {
  id: string
  note: string
}

/**
 * Saran Class Of Business untuk Group Business terpilih: `BUSINESS.NOTE` dengan
 * `BUSINESSGROUPID = groupBusinessId` - pola `BrowseBusiness_RD` (`[dugaan]` untuk form ini).
 */
export function daftarClassOfBusiness(groupBusinessId: string): Promise<{ baris: BarisClassOfBusiness[] }> {
  return minta<{ baris: BarisClassOfBusiness[] }>('/api/nbfacin/class-of-business', { kueri: { groupBusinessId } })
}

/** Isian form Opportunity yang dikirim saat `Create opportunity` (tiket 29). Teks apa adanya. */
export interface IsianOpportunity {
  /** Bentuk kabel inti `DD-MM-YYYY`. */
  estimatedClosingDate: string
  businessProspectName: string
  accountId: string
  insuredId: string
  groupBusinessId: string
  groupBusiness: string
  classOfBusiness: string
  typeOfInward: string
  /** Kosong bila Type Of Inward bukan Facultative (medan itu tidak tampil). */
  typeOfFacultative: string
  phase: string
  stage: string
  opportunitySource: string
  businessStatus: string
  description: string
}

/** Jawaban `POST /api/nbfacin/opportunity`: nomor case baru, mis. `NB-184352`. */
export interface HasilOpportunity {
  caseId: string
}

export function buatOpportunity(isian: IsianOpportunity): Promise<HasilOpportunity> {
  return minta<HasilOpportunity>('/api/nbfacin/opportunity', { metode: 'POST', badan: isian })
}

/** Isian blok General layar Inward Facultative (section Periode) - tiket 30/31. Tanggal kabel DD-MM-YYYY. */
export interface GeneralInward {
  reffNumber: string
  qqName: string
  beginDate: string
  offeringDate: string
  endDate: string
  policyType: string
  /** `.QuotationData.MOID` - id marketing officer terpilih. */
  marketingId: string
  day: string
  typeFacultative: string
  /**
   * Kode SOB terpilih di popup Change SOB (`.QuotationData.SourceOfBusiness`, mis. `G0000075`; kolom pasangan
   * DDL FACINPRODUCTION `SOBID`). Ikut dikirim Save for later; backend memeriksanya ke tabel AGENT (tiket 33).
   */
  sourceOfBusinessId: string
  /** Nama SOB (`.QuotationData.SobName`) - tampil-saja, diisi server dari tabel AGENT menurut kode. */
  sourceOfBusiness: string
  /**
   * Daftar ceding (`.Quotation.CedingCoList`: `.CedingCo` = AGENT.ID, `.CedingCoName` = AGENT.CLIENTNAME), urut
   * sesuai urutan pilih - diisi popup Change Ceding Co (tiket 34).
   */
  cedingList: BarisCeding[]
  /** Gabungan nama ceding berpemisah `;` (`.QuotationData.CedingCoName`, SetCedingCo_Act) - tampil-saja. */
  cedingCoName: string
  groupName: string
  oldPolicyNumber: string
}

/** Medan General yang dikirim tombol Save for later (tanpa medan tampil-saja). */
export type SimpanGeneral = Omit<GeneralInward, 'sourceOfBusiness' | 'cedingList' | 'cedingCoName' | 'groupName' | 'oldPolicyNumber'> & {
  /** Kode ceding (AGENT.ID) berurutan; nama diambil server dari AGENT (tiket 34). */
  cedingIds: string[]
}

/** Satu ceding terpilih. */
export interface BarisCeding {
  id: string
  name: string
}

/** Jawaban `GET /api/nbfacin/kasus/{caseId}` (tiket 31, backend sesi c3). */
export interface KasusNB {
  caseId: string
  position: string
  statusWork: string
  opportunity: IsianOpportunity
  insuredName: string
  general: GeneralInward
}

export function ambilKasus(caseId: string): Promise<KasusNB> {
  return minta<KasusNB>(`/api/nbfacin/kasus/${encodeURIComponent(caseId)}`)
}

export function simpanGeneral(caseId: string, general: SimpanGeneral): Promise<KasusNB> {
  return minta<KasusNB>(`/api/nbfacin/kasus/${encodeURIComponent(caseId)}/general`, { metode: 'PUT', badan: general })
}

/** Pilihan Marketing Name (`BrowseMarketingOfficer_RD`). */
export function daftarMarketing(): Promise<{ baris: { id: string; nama: string }[] }> {
  return minta<{ baris: { id: string; nama: string }[] }>('/api/nbfacin/marketing-officer')
}

/** Satu baris daftar case NB di portal (tiket 32, backend sesi c3). Teks apa adanya. */
export interface BarisCaseNB {
  caseId: string
  name: string
  groupBusiness: string
  insuredName: string
  marketing: string
  status: string
}

/** Jawaban `GET /api/nbfacin/opportunity` - daftar case NB Fac In, terbaru dulu. */
export interface HalamanCaseNB {
  baris: BarisCaseNB[]
  total: number
  halaman: number
  ukuran: number
}

/** `cari` "mengandung" atas case id atau nama (placeholder Pega "NB-1234 or Name"); `halaman` mulai 1. */
export function daftarCaseNB(cari: string, halaman: number): Promise<HalamanCaseNB> {
  return minta<HalamanCaseNB>('/api/nbfacin/opportunity', { kueri: { cari: cari.trim(), halaman } })
}

/** Satu baris tabel AGENT untuk popup Change SOB (tiket 33). Teks apa adanya. */
export interface BarisSOB {
  /** Kode SOB, mis. `G0000075` - disimpan ke `.QuotationData.SourceOfBusiness` `[dugaan]`. */
  id: string
  clientId: string
  /** Ditampilkan sebagai Source of business (`.QuotationData.SobName`). */
  name: string
}

/** Jawaban `GET /api/nbfacin/sob` - backend sesi c3 (tiket 33). */
export interface HalamanSOB {
  baris: BarisSOB[]
  total: number
  halaman: number
  ukuran: number
}

/** Syarat work owner 03-10-2026: StatusActive=1, AgentType2 != 'LIFE INSURANCE', ClientID terisi; cari tidak peka huruf. */
export function cariSOB(cari: string, halaman: number): Promise<HalamanSOB> {
  return minta<HalamanSOB>('/api/nbfacin/sob', { kueri: { cari: cari.trim(), halaman } })
}

/**
 * Satu objek tab Object (FIRE) - baris `.LocationList` (tiket 35). Semua teks apa adanya; Number of Floor teks
 * angka. Medan Risk Address diisi fitur Choose Risk Address (tahap berikut) - di sini tampil-saja, kecuali
 * Building No.
 */
export interface ObjekFire {
  /** `.Property.ObjectNo`. */
  objectNo: string
  /** `.Property.ObjectType`. */
  objectType: string
  /** `.Property.ObjectName` - kolom grid "Object Name". */
  objectName: string
  /** `.Property.IsMaterialDamage`. */
  isMaterialDamage: boolean
  /** `.Property.IsTopRisk`. */
  isTopRisk: boolean
  /** `.Property.RoadType` (Type). */
  roadType: string
  /** `.Property.RoadName` (Address). */
  roadName: string
  /** `.Property.BuildingNo`. */
  buildingNo: string
  /** `.Property.RiskLocation.ASMZipCode`. */
  zipCode: string
  /** `.Property.Country`. */
  country: string
  /** `.Property.RiskLocation.ASMAddress` (Risk Location) - kolom grid "Location". */
  riskLocation: string
  /** `.Property.RiskLocation.ASMRW` (Territory). */
  territory: string
  /** `.Property.RiskLocation.ASMCity`. */
  city: string
  /** `.Property.RiskLocation.ASMDistrict`. */
  district: string
  /** `.Property.Province`. */
  province: string
  /** `.Property.AlmRiskID` (Risk Address ID). */
  riskAddressId: string
  /** `.Property.BuildingConstruction.*` - teks apa adanya (Roof/Wall/Floor Type = kode). */
  numberOfFloor: string
  roofType: string
  wallType: string
  floorType: string
  partitionType: string
  supportWallType: string
  otherType: string
  /** `.Property.Ownership` (Surrounding Risk, tiket 38). */
  ownership: string
  /** `.Property.IsProductionProcessFlag`. */
  isProductionProcess: boolean
  /** `.Property.IsHotWorkProcessFlag`. */
  isHotWorkProcess: boolean
  /** `.Property.IsFlammableItemFlag`. */
  isFlammableItem: boolean
  /** `.Property.SurroundingRisk` (tiket 38). */
  surroundingRisk: SurroundingRisk
}

/** Satu sisi Surrounding Risk - `.{Front|Left|Back|Right}{Occupation|Construction|Distance|Note}`. */
export interface SisiRisiko {
  /** OCCUPATION.OLDID terpilih (autocomplete menyimpan `.OldID`). */
  occupation: string
  construction: string
  /** Teks angka, ≥ 0, 2 desimal (pxNumber). */
  distance: string
  /** Diisi OCCUPATION.NAME saat memilih Occupation. */
  note: string
}

/** `.Property.SurroundingRisk`. */
export interface SurroundingRisk {
  front: SisiRisiko
  left: SisiRisiko
  back: SisiRisiko
  right: SisiRisiko
  housekeepingStatus: string
  floodAreaStatus: string
  floodArea: string
  housekeepingRemark: string
}

/** Jawaban baca / simpan tab Object. */
export interface DaftarObjek {
  baris: ObjekFire[]
}

/** `GET /api/nbfacin/kasus/{caseId}/objek` - urut sesuai `.LocationList`. */
export function ambilObjek(caseId: string): Promise<DaftarObjek> {
  return minta<DaftarObjek>(`/api/nbfacin/kasus/${encodeURIComponent(caseId)}/objek`)
}

/** `PUT /api/nbfacin/kasus/{caseId}/objek` - tombol Save tab Object; daftar diganti utuh. */
export function simpanObjek(caseId: string, baris: ObjekFire[]): Promise<DaftarObjek> {
  return minta<DaftarObjek>(`/api/nbfacin/kasus/${encodeURIComponent(caseId)}/objek`, { metode: 'PUT', badan: { baris } })
}

/** Saringan popup Choose Risk Address - nama medan = properti Pega yang diikat sel 78-84. */
export interface SaringRisk {
  address: string
  zipCode: string
  country: string
  province: string
  city: string
  district: string
  territory: string
}

/** Satu baris RISKADDRESS (RD `BrowseRisksAddress_RD`, kelas Int-RISKADDRESS). Teks apa adanya. */
export interface BarisRisk {
  /** `.ID` -> `.Property.AlmRiskID`. */
  id: string
  /** `.Title` (Type, mis. "JL.", "OTHERS"). */
  title: string
  address: string
  nationName: string
  provinceName: string
  cityName: string
  districtName: string
  territoryName: string
  postalCode: string
}

export interface HalamanRisk {
  baris: BarisRisk[]
  total: number
  halaman: number
  ukuran: number
}

/** `GET /api/nbfacin/risk-address` - minimal satu saringan terisi (tiket 36). */
export function cariRiskAddress(saring: SaringRisk, halaman: number): Promise<HalamanRisk> {
  const kueri: Record<string, string | number> = { halaman }
  for (const [k, v] of Object.entries(saring)) if (v.trim() !== '') kueri[k] = v.trim()
  return minta<HalamanRisk>('/api/nbfacin/risk-address', { kueri })
}

/** Satu baris RW untuk saran Zip Code (RD `BrowseRW_RD`, STS_AKTIF = "1"). */
export interface BarisRW {
  zipCode: string
  /** `.Note` -> Territory. */
  territoryName: string
  districtName: string
  cityName: string
  provinceName: string
  nationName: string
}

/** `GET /api/nbfacin/rw?zipCode=` - saran Zip Code popup Add (tiket 37). */
export function cariZipCode(zipCode: string): Promise<{ baris: BarisRW[] }> {
  return minta<{ baris: BarisRW[] }>('/api/nbfacin/rw', { kueri: { zipCode: zipCode.trim() } })
}

/** Isian popup Add - urutan kolom yang ditulis ke RISKADDRESS (urutan parameter `InsertUpdateRISKADDRESS` Pega). */
export interface AlamatBaru {
  nationName: string
  provinceName: string
  districtName: string
  cityName: string
  territoryName: string
  title: string
  address: string
  postalCode: string
}

/**
 * `POST /api/nbfacin/risk-address` - Save popup Add: INSERT RISKADDRESS dari Go, TANPA memanggil prosedur
 * (ADR-0043; register butir 81 menggantikan W-1); `id` = ID baru dengan ekspresi prosedur Pega
 * `GETCURRENTSITE || LPAD(RISKADDRESS_SEQ.NEXTVAL, 12, '0')`.
 */
export function simpanAlamatBaru(a: AlamatBaru): Promise<{ id: string }> {
  return minta<{ id: string }>('/api/nbfacin/risk-address', { metode: 'POST', badan: a })
}

/** Satu baris OCCUPATION untuk saran Occupation (data page `D_BrowseOccupationFacInFIRE`, TYPE = 'FIRE'). */
export interface BarisOccupation {
  /** `.OldID` - nilai yang disimpan. */
  oldId: string
  /** `.Name` - masuk ke Note sisi itu. */
  name: string
}

/** `GET /api/nbfacin/occupation?cari=` - saran Occupation Surrounding Risk (tiket 38). */
export function cariOccupation(cari: string): Promise<{ baris: BarisOccupation[] }> {
  return minta<{ baris: BarisOccupation[] }>('/api/nbfacin/occupation', { kueri: { cari: cari.trim() } })
}
