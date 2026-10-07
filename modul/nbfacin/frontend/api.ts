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
  /**
   * Risk Scoring di General (`Section\Periode.xml`: `.ScoringRisk.FinalScore` + `.ScoringRisk.NoteFinalScore`, baca-saja)
   * = nilai TERSIMPAN `T_SCORINGRISK` (FINAL_SCORE, NOTE_FINAL_SCORE), kosong bila Scoring Risk belum pernah di-Save
   * (tiket 54, backend sesi c3).
   */
  riskScoring: { finalScore: string; noteFinalScore: string }
  /**
   * Tab Loss Record boleh dibuka (tiket 55): akun login ber-`JOB_POSITION` = "IT" (keputusan work owner 06-10-2026
   * "untuk sementara untuk POSITION=IT saja"; Pega `When\IsIT.xml`). Backend sesi c3: dibaca lunak - galat baca /
   * tanpa sesi = false, GET kasus tidak pernah gagal karenanya.
   */
  aksesLossRecord: boolean
  /**
   * Tab Correspondence boleh dibuka (tiket 56): aturan SAMA dengan `aksesLossRecord` - akun login ber-`JOB_POSITION` = "IT"
   * (keputusan work owner 06-10-2026 "ini visibility menu nya juga cuma it ya"; Pega S95 selalu tampil). Backend sesi c3:
   * satu pembacaan lunak bersama `aksesLossRecord`; endpoint korespondensi menolak 401 / 403.
   */
  aksesKorespondensi: boolean
  /**
   * Case boleh diubah (work owner 07-10-2026: case yang sudah keluar dari tahap marketing baca-saja untuk marketing):
   * aturan sama dengan `TampilanSubmit.bolehSubmit` (posisi tahap marketing Offer + akun memegang workbasket marketing).
   * false = layar baca-saja (`BacaSajaKasus`) dan backend menolak tulisan (409 tahap / 403 tanpa workbasket marketing;
   * backend sesi c3 `services/bacasaja.go`, satu sumber dengan Submit dan Scoring).
   */
  bolehUbah: boolean
  /** Sebab baca-saja (teks singkat); '' bila boleh ubah. */
  alasanBacaSaja: string
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
  /**
   * `.NBStatus` (sel 95 grid `GetListOpportunityF`) = `T_GENERAL_POLIS.NB_STATUS` (tiket 57); belum pernah di-Submit =
   * "NEW NB" (`Flow\InputInwardFacultativeOffer.xml` awal flow).
   */
  nbStatus: string
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
  /** `.Property.PropertyItemList` (tiket 39), urut = PropertyItemNo. */
  items: ItemObjek[]
  /** `.Property.OccupationList` (tiket 40). */
  occupations: OkupasiObjek[]
  /** `.FEAList` baris objek (tiket 41) - Fire Extinguisher Availability. */
  fea: BarisFEA[]
  /** `.Property.ListCauseOfLoss` (tiket 42). */
  lossRecords: CatatanKerugian[]
  /** Loss ratio objek - BACA-SAJA, dikirim server; diabaikan saat PUT (tiket 42). */
  lossRatio: LossRatio
  /** `.Property.TotalTSIPremiGrossList` - total per mata uang; BACA-SAJA (tiket 43). */
  totalPerCurrency?: TotalCoverage[]
  /** `.Property.ListCauseOfLossClaim` - BACA-SAJA (Loss Record Internal); diabaikan saat PUT (tiket 42). */
  internalLossRecords: KlaimInternal[]
}

/** Satu catatan kerugian - `ASM-FW-GISFW-Data-CauseOfLoss`. Uang = teks desimal (ADR-0003). */
export interface CatatanKerugian {
  /** `.DateOfLoss` - kabel `DD-MM-YYYY`. */
  dateOfLoss: string
  /** `.CoinsData.CoinsName` - kolom grid "Insured Name". */
  coinsName: string
  /** `.LossObject`. */
  lossObject: string
  /** `.Currency`. */
  currency: string
  /** `.Amount` "Total of Loss" - uang, tanpa isian di layar (baca-saja). */
  amount: string
  /** `.Claim` "Total Claim (100%)" - uang. */
  claim: string
  /** `.PreventionOfLoss` - uang. */
  preventionOfLoss: string
  /** `.CauseOfLoss`. */
  causeOfLoss: string
  /** `.Remarks`. */
  remarks: string
  /** `.Detail` "Loss Detail". */
  detail: string
}

/** `.LossRatio1Year*` / `.LossRatio35Year*` baris objek (T_LOCATIONLIST). Teks desimal. */
export interface LossRatio {
  oneYearAmount: string
  oneYearPercent: string
  threeFiveYearAmount: string
  threeFiveYearPercent: string
}

/** Satu baris Loss Record Internal - `ASM-FW-GISFW-Data-CauseOfLossClaim`. Uang = teks desimal. */
export interface KlaimInternal {
  /** `.DateOfLoss` - kabel `DD-MM-YYYY` (kolom Year = tahunnya). */
  dateOfLoss: string
  locationNo: string
  location: string
  currency: string
  premium: string
  /** `.EstiamsiKlaim` (O/S Claim). */
  osClaim: string
  /** `.AkseptasiKlaim` (Acccepted Claim). */
  acceptedClaim: string
  /** `.InccuredKlaim` (Incurred Claim). */
  incurredClaim: string
  /** `.LossRatio` - pecahan (Pega menampilkan dengan simbol %). */
  lossRatio: string
  remark: string
}

/** Satu baris FEA - `ASM-FW-GISFW-Data-OfferFacIn-OfferFEAList`. Jumlah unit = teks angka. */
export interface BarisFEA {
  /** `.APAR`. */
  apar: string
  /** `.Sprinkler`. */
  sprinkler: string
  /** `.SmokeDetector`. */
  smokeDetector: string
  /** `.Hydrant`. */
  hydrant: string
  /** `.DataFEA.PrivateTruckBrigade`. */
  privateTruckBrigade: string
  /** `.DataFEA.PrivateFireBrigade`. */
  privateFireBrigade: string
  /** `.DataFEA.TeamSOPSafety`. */
  teamSopSafety: string
  /** `.DataFEA.TeamSOPRiskManagement`. */
  teamSopRiskManagement: string
  /** `.InfoFEA`. */
  info: string
}

/** Satu baris Occupation - `ASM-FW-GISFW-Data-Occupation`. */
export interface OkupasiObjek {
  /** `.OccupationId` = OCCUPATION.OLDID. */
  occupationId: string
  /** `.OccupationName` = OCCUPATION.NAME. */
  occupationName: string
  /** `.TableOfLimit.Category` - romawi "I"/"II"/"III" dari KDRiskExposure (`SetDataOccupation`). */
  category: string
  /** `.TableOfLimit.Description` - Class of Construction. */
  constructionClass: string
  /** `.TableOfLimit.PctLimit` - persen batas, TEKS apa adanya (Pega berkoma desimal). */
  pctLimit: string
}

/**
 * Satu baris Object Item - `ASM-FW-GISFW-Data-PropertyItem`. `tsi` = UANG, teks desimal bertitik apa adanya (ADR-0003,
 * ADR-0016) - tidak pernah lewat float.
 */
export interface ItemObjek {
  /** `.ItemTypeID` = V_JN_OBJ_ITEM.MJOI_KODE. */
  itemTypeId: string
  /** `.ItemType` = V_JN_OBJ_ITEM.JN_OBJ_ITEM (kolom grid "Object Item Type"). */
  itemType: string
  /** `.PropertiItemNote` (KETERANGAN dari `GetObjectItem`). */
  note: string
  /** `.PropertyYear`. */
  propertyYear: string
  /** `.Unit`. */
  unit: string
  /** `.Condition`. */
  condition: string
  /** `.Currency` (CURRENCY.CURRENCY). */
  currency: string
  /** `.TSIObjectItem` - uang, teks desimal. */
  tsi: string
  /** `.Year` (Year of Planting). */
  yearOfPlanting: string
  /** `.NoOfTree`. */
  noOfTree: string
  /** `.AreaHectar`. */
  areaHectar: string
  /** `.Remark`. */
  remark: string
  /** `.IsAdjustableFlag`. */
  isAdjustable: boolean
  /** `.PctAdjust2` (dropdown, tampil bila tidak Adjustable). */
  pctAdjust2: string
  /** `.PctAdjustOther` (angka, tampil bila Adjustable). */
  pctAdjustOther: string
  /** `.CoverageList` (tab Coverage, tiket 43). Boleh absen di data lama. */
  coverages?: CoverageObjek[]
  /** `.TotalGrossPremi` - Σ Premium coverage; BACA-SAJA (dihitung server). */
  totalGrossPremi?: string
  /** `.TotalNetRate` ‰ - BACA-SAJA. */
  totalNetRate?: string
}

/**
 * Satu coverage - `ASM-FW-GISFW-Data-Coverage` (tiket 43). Uang / rate / persen = teks desimal bertitik (ADR-0003,
 * ADR-0034). Medan bertanda "server" dihitung backend (rumus `CountPremi_ACT`), diabaikan saat PUT.
 */
export interface CoverageObjek {
  /** `.Coverage` (ID COVERAGE_FACIN). */
  coverage: string
  /** `.OLDID` - kode tampil kolom "Coverage". */
  oldId: string
  /** `.CoverageNote` - nama coverage. */
  coverageNote: string
  /** `.CoverageBasis` "1".."5". */
  coverageBasis: string
  /** `.Day`. */
  day: string
  /** `.TSI` - server (= TSI Object Item). */
  tsi: string
  /** `.Indemnity`. */
  indemnity: string
  /** `.Rate` ‰ Gross Rate. */
  rate: string
  /** `.RateOJK` ‰ Standard Rate (lookup tarif = tahap C4). */
  rateOjk: string
  /** `.FirstLoss` %. */
  firstLoss: string
  /** `.DiscountPercentage` %. */
  discountPercentage: string
  /** `.TSILiability` - server. */
  tsiLiability: string
  /** `.NetRate` ‰. */
  netRate: string
  /** `.LimitofLiability`. */
  limitOfLiability: string
  /** `.PctLoL` %. */
  pctLol: string
  /** `.ProRatePercent` - server (periode polis). */
  proRatePercent: string
  /** `.IndemnityPercentage` %. */
  indemnityPercentage: string
  /** `.FirstScale` %. */
  firstScale: string
  /** `.Sublimit` %. */
  sublimit: string
  /** `.LostLimit` %. */
  lostLimit: string
  /** `.EmlPml` %. */
  emlPml: string
  /** `.Discount` (uang). */
  discount: string
  /** `.Premium` Gross Premium (uang) - server bila mode "percent". */
  premium: string
  /** `.Conditions`. */
  conditions: string
  /** `.DeductibleList` (tahap C3, tiket 45). Boleh absen di data lama. */
  deductibles?: Deductible[]
  /** `.Unit` Indemnity Unit (kode; tiket 46). Boleh absen. */
  unit?: string
  /** `.AccumulationCode` (tiket 46; Choose Accumulation Code / Copy Accumulation). Boleh absen. */
  accumulationCode?: string
  /** `.AccumulationDescription` = Accumulation Address (tiket 46). Boleh absen. */
  accumulationDescription?: string
  /** `.Zone` = ZONES.DESCRIPTION terpilih (tiket 52; hanya coverage `COVERAGE_ZONA`). Boleh absen di data lama. */
  zone?: string
  /** `.ConstructionEarthquake` = CONSTRUCTION.CONSTRUCTION (tiket 52; hanya `COVERAGE_KONSTRUKSI`). Boleh absen. */
  constructionEarthquake?: string
}

/** Satu deductible - `ASM-FW-GISFW-Data-Deductible`. Persen / uang = teks desimal. */
export interface Deductible {
  /** `.TypeDeductible` (kode, `DDL\TypeDeductible.xml`). */
  typeDeductible: string
  /** `.PctDeductible` %. */
  pctDeductible: string
  /** `.MinMax` 1 Min / 2 Max / 3 Or. */
  minMax: string
  /** `.Currency`. */
  currency: string
  /** `.TypeDeductible2` (tampil bila MinMax = 3). */
  typeDeductible2: string
  /** `.PctDeductible2` %. */
  pctDeductible2: string
  /** `.Condition` (kode, `DDL\Condition.xml` Deductible). */
  condition: string
  /** `.Amount` (uang). */
  amount: string
  /** `.InputCondition` (bila Condition = 5 Other). */
  inputCondition: string
  /** `.TimeExcess` hari. */
  timeExcess: string
}

/** Mode hitung `CountPremi_ACT`: dari rate ("percent") atau dari premi ("amount", rate dihitung balik). */
export type ModeHitung = 'percent' | 'amount'

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
  /** `.KDRiskExposure` ("01"/"02"/"03") - dipakai sub-tab Occupation (tiket 40); boleh absen di jawaban lama. */
  kdRiskExposure?: string
}

/** `GET /api/nbfacin/occupation?cari=` - saran Occupation Surrounding Risk (tiket 38). */
export function cariOccupation(cari: string): Promise<{ baris: BarisOccupation[] }> {
  return minta<{ baris: BarisOccupation[] }>('/api/nbfacin/occupation', { kueri: { cari: cari.trim() } })
}

/** Satu jenis item objek - view V_JN_OBJ_ITEM (RD `BrowseV_JN_OBJ_ITEM`, ISACTIVE = 1). */
export interface JenisItem {
  /** MJOI_KODE. */
  kode: string
  /** JN_OBJ_ITEM. */
  nama: string
  /** KETERANGAN (`GetObjectItem`) -> Object Item Note. */
  keterangan: string
}

/** `GET /api/nbfacin/jenis-item-objek` - pilihan Object Item Type (tiket 39). */
export function daftarJenisItem(): Promise<{ baris: JenisItem[] }> {
  return minta<{ baris: JenisItem[] }>('/api/nbfacin/jenis-item-objek')
}

/** `GET /api/nbfacin/mata-uang` - pilihan Currency (RD `BrowseCurrency_RD`, tanpa "ITL"). */
export function daftarMataUang(): Promise<{ baris: string[] }> {
  return minta<{ baris: string[] }>('/api/nbfacin/mata-uang')
}

/** Satu baris TABLEOFLIMIT (RD `BrowseTableOfLimit_RD`). */
export interface BarisTableOfLimit {
  description: string
  /** Teks apa adanya. */
  pctLimit: string
}

/**
 * `GET /api/nbfacin/kasus/{caseId}/table-of-limit?category=` - popup Choose Class of Construction (tiket 40). Tahun dan
 * Bizcode (`pyWorkPage.OfferFacIn.CurrentYear`, `QuotationData.BusinessCode`) diturunkan server dari case.
 */
export function cariTableOfLimit(caseId: string, category: string): Promise<{ baris: BarisTableOfLimit[] }> {
  return minta<{ baris: BarisTableOfLimit[] }>(`/api/nbfacin/kasus/${encodeURIComponent(caseId)}/table-of-limit`, { kueri: { category } })
}

/** Satu baris total per mata uang objek (`.Property.TotalTSIPremiGrossList`). Teks desimal. */
export interface TotalCoverage {
  currency: string
  tsi: string
  premium: string
  /** Premium / TSI × 1000 (‰). */
  rate: string
}

/** Satu baris COVERAGE_FACIN (RD `BrowseCoverageFacIn_RD`, Type FIRE). */
export interface BarisCoverage {
  id: string
  oldId: string
  nama: string
}

/** `GET /api/nbfacin/coverage?cari=` - popup Choose Coverage (tiket 43). */
export function cariCoverage(cari: string): Promise<{ baris: BarisCoverage[] }> {
  return minta<{ baris: BarisCoverage[] }>('/api/nbfacin/coverage', { kueri: { cari: cari.trim() } })
}

/** `GET /api/nbfacin/coverage-otomatis` - lima coverage awal FIRE (`AddCoverageAutoFire`). */
export function coverageOtomatis(): Promise<{ baris: BarisCoverage[] }> {
  return minta<{ baris: BarisCoverage[] }>('/api/nbfacin/coverage-otomatis')
}

/**
 * `POST /api/nbfacin/kasus/{caseId}/hitung-coverage` - hitung satu coverage (`CountPremi_ACT`) tanpa menyimpan: TSI dari
 * item, Prorate dari periode polis case.
 */
export function hitungCoverage(
  caseId: string,
  badan: {
    coverage: CoverageObjek
    tsi: string
    mode: ModeHitung
    /** `Param.DiscountStatus` - medan diskon yang diubah: "% Discount" = percent, "Discount" = amount. */
    modeDiskon?: ModeHitung
    /** Item pemilik coverage (`CountPremi_ACT` langkah 11-16, PctAdjustment). */
    isAdjustable?: boolean
    pctAdjustOther?: string
    /**
     * `Param.CallAct="rate"` (tiket 52): Coverage / Zone diubah -> `CountPremiRate_ACT`; coverage ber-OLDID "4.2 PRGBI"
     * (dicocokkan harfiah seperti Pega `.OLDID`) mengambil ‰ Gross Rate dari `POOLDATA.ZONEOJK` menurut Zone. Absen = hitung
     * biasa.
     */
    callAct?: 'rate'
    /** Hanya bersama `callAct`: `.Property.Province` objek (4.3 = 0.5 untuk DKI JAKARTA / BANTEN / JAWA BARAT, selain itu 0.45). */
    province?: string
    /** Hanya bersama `callAct`: OccupationId okupasi objek (FLEXAS = max OCCUPATION.RATELOWER). */
    occupationIds?: string[]
  },
): Promise<CoverageObjek> {
  return minta<CoverageObjek>(`/api/nbfacin/kasus/${encodeURIComponent(caseId)}/hitung-coverage`, { metode: 'POST', badan })
}

/**
 * `POST /api/nbfacin/kasus/{caseId}/hitung-net-rate` - tahap C2 (`CalculateNetRate_ACT`): bagi ‰ Total Net Rate ke
 * setiap coverage (NetRate = Rate / ΣRate × TotalNetRate) lalu hitung ulang premi. Tanpa menyimpan.
 */
export function hitungNetRate(
  caseId: string,
  badan: { tsi: string; totalNetRate: string; coverages: CoverageObjek[]; isAdjustable?: boolean; pctAdjustOther?: string },
): Promise<{ coverages: CoverageObjek[]; totalNetRate: string; flagNetRate: boolean }> {
  return minta(`/api/nbfacin/kasus/${encodeURIComponent(caseId)}/hitung-net-rate`, { metode: 'POST', badan })
}

/**
 * Saringan popup Choose Accumulation (`SearchRiskAccumCov`; kosong = tidak menyaring). Backend menjalankan
 * `GetDataAccumulation_act`: City / District / Policy No terisi DAN Zip Code kosong -> jalur SQL (grid sel 48), selain
 * itu RD `SearchRiskAccumulation_RD` (grid sel 90).
 */
export interface SaringAkumulasi {
  /** `SearchAccumulation.ID` (Accumulation Code). */
  id: string
  /** `SearchAccumulation.CenterTransStatus` (Policy No). */
  policyNo: string
  /** `SearchAccumulation.Note` (Road; RD Contains tanpa peka huruf). */
  note: string
  /** `SearchAccumulation.PostalCode` (Zip Code; awal = zip lokasi risiko; memilih Area menimpanya). */
  postalCode: string
  /** `SearchAccumulation.SyariahStatus` (= NationInitial Country terpilih). */
  syariahStatus: string
  /** `SearchAccumulation.ProvinceID` (Province terpilih). */
  provinceId: string
  /** `InputFilter.City` (ID City terpilih). */
  cityId: string
  /** `InputFilter.District` (ID District terpilih). */
  districtId: string
  /** `SearchAccumulation.CZone`. */
  czone: string
  /** `SearchAccumulation.Keyword`. */
  keyword: string
}

/** Satu baris hasil (grid sel 90 `.ID` / `.AccumulationName` / `.Note`, atau sel 48 `.CARI1` / `.CARI2` / `.CARI3`). */
export interface BarisAkumulasi {
  id: string
  /** `.AccumulationName` (kolom Type). */
  accumulationName: string
  /** `.Note` (kolom Description; ditulis ke AccumulationDescription). */
  note: string
}

/** `GET /api/nbfacin/akumulasi` - tombol Filter (tiket 46). */
export function cariAkumulasi(saring: SaringAkumulasi): Promise<{ baris: BarisAkumulasi[] }> {
  const kueri: Record<string, string> = {}
  for (const [k, v] of Object.entries(saring)) if (v.trim() !== '') kueri[k] = v.trim()
  return minta<{ baris: BarisAkumulasi[] }>('/api/nbfacin/akumulasi', { kueri })
}

/** Jenis saran autocomplete popup accumulation - satu per RD sel 27-35 `SearchRiskAccumCov`. */
export type JenisSaranAkumulasi = 'nation' | 'province' | 'accumtype' | 'city' | 'district' | 'area' | 'czone'

/** Satu saran: `label` = teks yang ditulis ke kotak; `id` / `ekstra` = nilai yang disalin Pega saat memilih. */
export interface SaranAkumulasi {
  id: string
  label: string
  /** nation: `.NationInitial`; area: `.ZipCode`. */
  ekstra?: string
}

/**
 * `GET /api/nbfacin/akumulasi/saran/{jenis}?q=&induk=` - `induk` = saringan berantai Pega: province <- nama Country,
 * city <- ProvinceID, district <- nama City, area <- nama District.
 */
export function saranAkumulasi(jenis: JenisSaranAkumulasi, q: string, induk = ''): Promise<{ baris: SaranAkumulasi[] }> {
  const kueri: Record<string, string> = {}
  if (q.trim() !== '') kueri.q = q.trim()
  if (induk.trim() !== '') kueri.induk = induk.trim()
  return minta<{ baris: SaranAkumulasi[] }>(`/api/nbfacin/akumulasi/saran/${jenis}`, { kueri })
}

/** Badan `POST /api/nbfacin/akumulasi` - Add New accumulation (port `SaveAccumulation_Act` + RDBMASTERACCUMULATION). */
export interface AkumulasiBaru {
  /** `InputAccumulation.Accumulation` = ID tipe akumulasi (ACCUMULATEDTYPE.ID). */
  accumulation: string
  /** `InputAccumulation.AccumulationType` = nama tipe. */
  accumulationType: string
  /** Accumulation Description (`.Note`; disimpan huruf besar). */
  note: string
  keyword: string
  scopeArea: string
  /** Cresta Zone (`.CZone` = CZONE.CODE) dan ID-nya (`.CZoneID`). */
  cZone: string
  cZoneId: string
  provinceId: string
  /** Primary Zip Code (`.PostalCode`). */
  zipCode: string
  /** NationInitial dari zip terpilih (`InputData.CARI34` -> awalan ID). */
  negara: string
}

/** `POST /api/nbfacin/akumulasi` -> ID baru + Note tersimpan. */
export function tambahAkumulasi(isi: AkumulasiBaru): Promise<{ id: string; note: string }> {
  return minta<{ id: string; note: string }>('/api/nbfacin/akumulasi', { metode: 'POST', badan: isi })
}

/** `GET /api/nbfacin/akumulasi/czone?zip=` - Cresta Zone awal dari zip (`GetCzone_Act`, `GetAccumulationZipcode_SQL`). */
export function czoneDariZip(zip: string): Promise<{ cZoneId: string; cZone: string }> {
  return minta<{ cZoneId: string; cZone: string }>('/api/nbfacin/akumulasi/czone', { kueri: { zip: zip.trim() } })
}

/** Satu baris Choose Zip Code (`BrowseRiskAddressZipCode_RD`). */
export interface BarisZipAkumulasi {
  zipCode: string
  city: string
  province: string
  nation: string
  /** NationInitial (`SetCountryID_Act`) - awalan ID akumulasi. */
  nationInitial: string
}

/** Satu halaman Choose Zip Code (paging server, 15 baris - permintaan work owner 04-10-2026). */
export interface HalamanZipAkumulasi {
  baris: BarisZipAkumulasi[]
  total: number
  halaman: number
  ukuran: number
}

/** `GET /api/nbfacin/akumulasi/zipcode?provinceName=&q=&halaman=`. */
export function cariZipAkumulasi(provinceName: string, q: string, halaman: number): Promise<HalamanZipAkumulasi> {
  const kueri: Record<string, string | number> = { halaman }
  if (provinceName.trim() !== '') kueri.provinceName = provinceName.trim()
  if (q.trim() !== '') kueri.q = q.trim()
  return minta<HalamanZipAkumulasi>('/api/nbfacin/akumulasi/zipcode', { kueri })
}

/** Satu argumen klausa (`ASM-FW-GISFW-Data-Argument`; SearchClauseArgFireSQL atas M_ARGCLAUSEFIRE). */
export interface ArgumenKlausa {
  argumentNumber: string
  argumentDescription: string
  /** Isi argumen (awal = DefaultValue). */
  argumentValue: string
}

/** Satu baris `pyWorkPage.ClauseList` kasus (`ASM-FW-GISFW-Data-Clause`, `SearchClauseFireSQL_PostAct`). */
export interface KlausaKasus {
  clauseCode: string
  clauseTitle: string
  clauseDescription: string
  /** Kode bahasa "0" Indonesia / "1" Inggris / "2" lain (`SearchClauseFireSQL_PostAct`). */
  clauseLanguage: string
  clauseLanguageId: string
  /** Isi klausa sesudah argumen diterapkan. */
  clauseContent: string
  /** Isi klausa asli (`.Text` master) - dasar penggantian argumen. */
  clauseContentTemp: string
  argumentCount: string
  argumentList: ArgumenKlausa[]
}

/** Satu hasil popup Choose Clause (`SearchClauseListOutput`, kelas Int-CLAUSE). */
export interface HasilKlausa {
  id: string
  title: string
  /** `.Info` (kolom Description). */
  info: string
  /** `.Text` - isi klausa. */
  text: string
  language: string
  /** Jumlah argumen (RetrieveArgumentNumberClauseSQL). */
  argumentCount: string
}

/** `GET /api/nbfacin/klausa?bahasa=&q=&halaman=` - popup Choose Clause (tiket 47). */
export function cariKlausa(
  bahasa: string,
  q: string,
  halaman: number,
): Promise<{ baris: HasilKlausa[]; total: number; halaman: number; ukuran: number }> {
  const kueri: Record<string, string | number> = { halaman }
  if (bahasa.trim() !== '') kueri.bahasa = bahasa.trim()
  if (q.trim() !== '') kueri.q = q.trim()
  return minta('/api/nbfacin/klausa', { kueri })
}

/** `GET /api/nbfacin/klausa/{id}/argumen` - argumen + nilai bawaan (ViewClauseArgFireSQL). */
export function argumenKlausa(id: string): Promise<{ baris: ArgumenKlausa[] }> {
  return minta(`/api/nbfacin/klausa/${encodeURIComponent(id)}/argumen`)
}

/** `GET /api/nbfacin/kasus/{caseId}/klausa` - ClauseList kasus. */
export function ambilKlausaKasus(caseId: string): Promise<{ baris: KlausaKasus[] }> {
  return minta(`/api/nbfacin/kasus/${encodeURIComponent(caseId)}/klausa`)
}

/** `PUT /api/nbfacin/kasus/{caseId}/klausa` - simpan ClauseList kasus (Save tab, `SaveFacIn_Act`). */
export function simpanKlausaKasus(caseId: string, baris: KlausaKasus[]): Promise<{ baris: KlausaKasus[] }> {
  return minta(`/api/nbfacin/kasus/${encodeURIComponent(caseId)}/klausa`, { metode: 'PUT', badan: { baris } })
}

/** Satu pilihan Type Treaty (`GetTreatyName` -> SpreadingList1). */
export interface TreatySpreading {
  id: string
  name: string
}

/** Satu baris template Copy Spreading. */
export interface TemplateSpreading {
  treatyType: string
  treatyName: string
  sharePercentage: string
}

/** Satu baris `CoverageList(k).SpreadingList` (`ASM-FW-GISFW-Data-SpreadingRisk`). Angka = teks desimal. */
export interface BarisSpreading {
  treatyType: string
  treatyName: string
  sharePercentage: string
  tsiGrossSpreaded: string
  tsiSpreaded: string
  claimEstimation: string
  premiumSpreaded: string
}

export interface CoverageSpreading {
  oldId: string
  coverageBasis: string
  rate: string
  premium: string
  discount: string
  tsi: string
  tsiLiability: string
  tsiNusantaraRe: string
  premiNusantaraRe: string
  spreading: BarisSpreading[]
}

export interface ItemSpreading {
  itemType: string
  currency: string
  tsi: string
  totalGrossPremi: string
  totalPremiumRnm: string
  coverages: CoverageSpreading[]
}

/** Satu baris total per lokasi (`Property.TotalTSIPremiSpreadRNM`) / ringkasan per treaty (`TotalSpreadingCurrency`). */
export interface TotalSpreading {
  currency: string
  treatyType: string
  treatyName: string
  sharePercentage: string
  claimSpreaded: string
  tsiSpreaded: string
  claimEstimation: string
  /** Hanya ringkasan per treaty (Total First Loss). */
  claimAmountIdr?: string
  premiumSpreaded: string
}

export interface LokasiSpreading {
  objectNo: string
  objectName: string
  location: string
  items: ItemSpreading[]
  total: TotalSpreading[]
}

/** Ringkasan per mata uang (`TotalSpreadAll`). */
export interface RingkasanMataUangSpreading {
  currency: string
  tsiTopRisk: string
  tsi: string
  lol: string
  premium: string
}

/** `GET /api/nbfacin/kasus/{caseId}/spreading` (tiket 48, kontrak usulan sesi 0f). */
export interface TampilanSpreading {
  percentShare: string
  treaty: TreatySpreading[]
  template: TemplateSpreading[]
  lokasi: LokasiSpreading[]
  ringkasanTreaty: TotalSpreading[]
  ringkasanMataUang: RingkasanMataUangSpreading[]
  pesan: string[]
}

const jalurSpreading = (caseId: string) => `/api/nbfacin/kasus/${encodeURIComponent(caseId)}/spreading`

export function ambilSpreading(caseId: string): Promise<TampilanSpreading> {
  return minta<TampilanSpreading>(jalurSpreading(caseId))
}

/** % Share RNM diubah: hitung TSI / premi Nusantara Re + spreading otomatis (`GetKapasitasTreaty`); tidak menyimpan. */
export function hitungShareSpreading(caseId: string, percentShare: string): Promise<TampilanSpreading> {
  return minta<TampilanSpreading>(`${jalurSpreading(caseId)}/hitung-share`, { metode: 'POST', badan: { percentShare } })
}

/** Copy To All Spreading (`CopyToAllSpreadingFire_ACT`); tidak menyimpan. */
export function salinSpreading(caseId: string, percentShare: string, template: TemplateSpreading[]): Promise<TampilanSpreading> {
  return minta<TampilanSpreading>(`${jalurSpreading(caseId)}/salin`, { metode: 'POST', badan: { percentShare, template } })
}

/** Save tab Spreading (`IsThereAnyObjectLocation_Act`): hitung ulang lalu simpan. */
export function simpanSpreading(caseId: string, isi: TampilanSpreading): Promise<TampilanSpreading> {
  return minta<TampilanSpreading>(jalurSpreading(caseId), {
    metode: 'PUT',
    badan: { percentShare: isi.percentShare, template: isi.template, lokasi: isi.lokasi },
  })
}

/** Satu baris `.OfferFacIn.CedingCedantList` (tiket 49). % Share = teks desimal. */
export interface BarisCedant {
  cedingCo: string
  cedingCoName: string
  shareCeding: string
}

/** `GET /api/nbfacin/kasus/{caseId}/cedant` (tiket 49, kontrak usulan sesi 0f). Angka = teks desimal. */
export interface TampilanCedant {
  /** `.OfferFacIn.QuotationData.SobName`. */
  sobName: string
  /** `.OfferFacIn.ShareCedantType`: "" | "0" (Gross) | "1" (Share RNM). */
  shareCedantType: string
  /** `.OfferFacIn.PercentShare` (baca-saja di tab ini). */
  percentShare: string
  /**
   * `.OfferFacIn.TotalTSINusaReSpreading` dalam IDR: tiap mata uang dikali kurs TOIDR (`TREATYEXCHANGEYEARLY`) lalu
   * dijumlah (work owner 05-10-2026, C-8 - Pega menjumlah angka tanpa kurs). "" bila ada mata uang tanpa kurs.
   */
  totalTsiRnm: string
  /** `.OfferFacIn.TotalPremiNusaRe` dalam IDR (aturan sama dengan `totalTsiRnm`). */
  totalPremiRnm: string
  /**
   * Mata uang yang kursnya ke IDR tidak ditemukan (KursPada Begin date, sama dengan Spreading; item tanpa mata uang =
   * "(kosong)") - bila berisi, kedua total "". Selalu larik (C-8, K49-11).
   */
  kursTidakAda: string[]
  /** Proteksi `ProtectShareCedant_Act` berlaku (IsGroup == "Group" atau FlagOnGoingPolicy == 2). */
  wajib: boolean
  cedant: BarisCedant[]
  /** Ceding Co tab General (`QuotationData.CedingCoList`, atau CedingCo tunggal) - sumber Add pertama. */
  cedingUmum: BarisCedant[]
  /** `.OfferFacIn.QuotationData.ShareOfCeding` tersimpan ("100%" / "N%"), "" bila kosong; teks "OF" di detail baris. */
  shareOfCeding: string
  /**
   * Dasar CurrencyList per mata uang pada % Share cedant 100 (`SetTSIPremiCedant_Act` langkah 4 atas
   * `OfferFacIn.CurrencyList`; langkah 5 `TotalSpreadAll` ber-label `//` = mati). tsiRnm = Σ TSI objek x PS / 100,
   * premiumRnm = Σ premi NR, tsiGross = Σ TSI objek, premiumGross = Σ premi NR x 100 / PS. [] bila % Share RNM kosong.
   */
  dasarMataUang: DasarCedant[]
}

/** Satu baris dasar CurrencyList cedant (teks desimal). */
export interface DasarCedant {
  currency: string
  tsiRnm: string
  premiumRnm: string
  tsiGross: string
  premiumGross: string
}

const jalurCedant = (caseId: string) => `/api/nbfacin/kasus/${encodeURIComponent(caseId)}/cedant`

export function ambilCedant(caseId: string): Promise<TampilanCedant> {
  return minta<TampilanCedant>(jalurCedant(caseId))
}

/** Save tab Cedant: simpan Share Cedant Type + grid; backend menjalankan SetShareOfCeding dan proteksi. */
export function simpanCedant(caseId: string, isi: { shareCedantType: string; cedant: BarisCedant[] }): Promise<TampilanCedant> {
  return minta<TampilanCedant>(jalurCedant(caseId), { metode: 'PUT', badan: isi })
}

/** Satu mata uang tab Premium Deduction / Brokerage Fee (`.OfferFacIn.CurrencyList`, tiket 50). Angka = teks desimal. */
export interface BarisDeduction {
  /** `.Name`. */
  currency: string
  /** `.Policy.Payment.Premium` = Σ PremiNusantaraRe mata uang itu (baca-saja; dasar semua hitungan). */
  premium: string
  /** `.Policy.Payment.Commision` - Total Premium Deduction. */
  commision: string
  /** `.Policy.Payment.Deduction2` - Total Premium Deduction 2. */
  deduction2: string
  /** `.Policy.Payment.BrokerageFee`. */
  brokerageFee: string
  /** `.Policy.Payment.PPh` / `.PPN` (baca-saja). */
  pph: string
  ppn: string
}

/** `GET /api/nbfacin/kasus/{caseId}/deduction` (tiket 50, kontrak usulan sesi 0f). */
export interface TampilanDeduction {
  /** `.Policy.Payment.RICommision` - % Premium Deduction. */
  pctPremiumDeduction: string
  /** `.Policy.Payment.PctBrokerageFee`. */
  pctBrokerageFee: string
  /** `.OfferFacIn.PPnCheck` - "Exclude PPN ?" (Pega selalu menyetelnya true saat menghitung). */
  ppnCheck: boolean
  /** Tarif `POOLDATA.TAX` baris TANGGAL terbaru (PPH / PPN, persen); "" bila tabel kosong. */
  tarifPph: string
  tarifPpn: string
  mataUang: BarisDeduction[]
  pesan: string[]
}

/** Pemicu hitung: isian persen atas (`SetRIComm_Act`), atau sel grid baris `idx`. */
export type AksiDeduction = 'persen' | 'commision' | 'deduction2' | 'brokerage'

const jalurDeduction = (caseId: string) => `/api/nbfacin/kasus/${encodeURIComponent(caseId)}/deduction`

const badanDeduction = (isi: TampilanDeduction) => ({
  pctPremiumDeduction: isi.pctPremiumDeduction,
  pctBrokerageFee: isi.pctBrokerageFee,
  ppnCheck: isi.ppnCheck,
  mataUang: isi.mataUang.map((b) => ({ currency: b.currency, commision: b.commision, deduction2: b.deduction2, brokerageFee: b.brokerageFee })),
})

export function ambilDeduction(caseId: string): Promise<TampilanDeduction> {
  return minta<TampilanDeduction>(jalurDeduction(caseId))
}

/** Hitung ulang (tanpa simpan): `SetRIComm_Act` / `SetTotalPremium_Act` / `ProtekDeduction2_Act` / `TotalBrokerageFee_Act`. */
export function hitungDeduction(caseId: string, aksi: AksiDeduction, idx: number, isi: TampilanDeduction): Promise<TampilanDeduction> {
  return minta<TampilanDeduction>(`${jalurDeduction(caseId)}/hitung`, { metode: 'POST', badan: { aksi, idx, ...badanDeduction(isi) } })
}

/** Save tab (`SaveFacIn_Act`): hitung ulang lalu simpan. */
export function simpanDeduction(caseId: string, isi: TampilanDeduction): Promise<TampilanDeduction> {
  return minta<TampilanDeduction>(jalurDeduction(caseId), { metode: 'PUT', badan: badanDeduction(isi) })
}

/** Satu cicilan `.Policy.Payment.ListInstallment` (tiket 51). Angka = teks desimal; Due Date kabel `DD-MM-YYYY`. */
export interface CicilanPayment {
  installmentNo: string
  dueDate: string
  installmentPercentage: string
  premium: string
  /** `.RICommision` - kolom "Deduction". */
  riCommision: string
  deduction2: string
  brokerageFee: string
  pph: string
  ppn: string
  discount: string
  paymentTotal: string
}

/** Satu mata uang `.OfferFacIn.CurrencyList` di tab Payment. */
export interface MataUangPayment {
  currency: string
  /** `.Policy.Payment.Premium` - Premium RNM. */
  premium: string
  /** `.SumTotalPayment` - Total Payment. */
  sumTotalPayment: string
  cicilan: CicilanPayment[]
}

/** `GET /api/nbfacin/kasus/{caseId}/payment` (tiket 51, kontrak usulan sesi 0f). */
export interface TampilanPayment {
  /** `.Policy.Payment.Installment` (teks; "" bila belum diisi). */
  installment: string
  /** Tarif `POOLDATA.TAX` terbaru untuk kepala kolom PPH / PPN; "" bila kosong. */
  tarifPph: string
  tarifPpn: string
  mataUang: MataUangPayment[]
  /** Pesan proteksi (`ProtectPremiPolicy_Act`) / peringatan dari backend. */
  pesan: string[]
}

/** Pemicu hitung: jumlah cicilan (`FillPaymentInstallment`) atau % cicilan baris `idx` mata uang `currency`. */
export type AksiPayment = 'installment' | 'persen'

const jalurPayment = (caseId: string) => `/api/nbfacin/kasus/${encodeURIComponent(caseId)}/payment`

const badanPayment = (isi: TampilanPayment) => ({
  installment: isi.installment,
  mataUang: isi.mataUang.map((m) => ({
    currency: m.currency,
    cicilan: m.cicilan.map((c) => ({ installmentNo: c.installmentNo, dueDate: c.dueDate, installmentPercentage: c.installmentPercentage })),
  })),
})

export function ambilPayment(caseId: string): Promise<TampilanPayment> {
  return minta<TampilanPayment>(jalurPayment(caseId))
}

/** Hitung ulang tanpa simpan: `FillPaymentInstallment` (installment) / `SetValidateInstallment_Act` (persen). */
export function hitungPayment(caseId: string, aksi: AksiPayment, currency: string, idx: number, isi: TampilanPayment): Promise<TampilanPayment> {
  return minta<TampilanPayment>(`${jalurPayment(caseId)}/hitung`, { metode: 'POST', badan: { aksi, currency, idx, ...badanPayment(isi) } })
}

/** Save tab Payment: hitung ulang (`SaveFillPaymentInstallment` + `CountPaymentInstallment_Act`), proteksi, simpan. */
export function simpanPayment(caseId: string, isi: TampilanPayment): Promise<TampilanPayment> {
  return minta<TampilanPayment>(jalurPayment(caseId), { metode: 'PUT', badan: badanPayment(isi) })
}

/**
 * `GET /api/nbfacin/zona?coverage=&oldId=&q=` (tiket 52, kontrak usulan) - RD `BrowseZones_RD` atas view `POOLDATA.ZONES`:
 * `OLDID` = OLDID coverage, zona X/Y/Z (RSMD) atau 1-6 (gempa) menurut kode coverage (`SetPropertyCoverageDT_FacIn` langkah
 * 10-11). Saran: id = ZONES.ID, label = ZONES.DESCRIPTION.
 */
export function saranZona(coverage: string, oldId: string, q: string): Promise<{ baris: SaranAkumulasi[] }> {
  const p = new URLSearchParams({ coverage, oldId, q })
  return minta<{ baris: SaranAkumulasi[] }>(`/api/nbfacin/zona?${p.toString()}`)
}

/**
 * `GET /api/nbfacin/konstruksi-gempa?q=` (tiket 52, kontrak usulan) - `SearchClassConstructionReas_SQL`
 * (`select distinct construction from construction order by CONSTRUCTION`). Saran: id = label = CONSTRUCTION.
 */
export function saranKonstruksiGempa(q: string): Promise<{ baris: SaranAkumulasi[] }> {
  return minta<{ baris: SaranAkumulasi[] }>(`/api/nbfacin/konstruksi-gempa?${new URLSearchParams({ q }).toString()}`)
}

/** Satu kolom Scoring Risk (Top Risk / Dominant Risk) - `.ScoringRisk.DataScoringRiskList(n)` (tiket 54). */
export interface KolomScoring {
  occupationNote: string
  topRiskLocation: string
  /** `.IndexLocation` = ObjectNo lokasi (Dominant Risk); "" untuk kolom Top Risk / tanpa Dominant Risk. */
  indexLocation: string
  /** Kunci sub-faktor -> nomor ChechBox terpilih (1..n), 0 = tidak dicentang. */
  pilihan: Record<string, number>
  /** `.LossRatio.LossRatio.PercentLossRatio` (teks desimal). */
  persenLossRatio: string
  /** Kunci sub-faktor -> ScorePerFactor (teks desimal). */
  skor: Record<string, string>
  /** `.Score` kolom. */
  total: string
}

/** `.ScoringRisk.ScoringResult.*` - rekomendasi dan persetujuan. */
export interface HasilScoring {
  acceptableNotApproval: boolean
  acceptableWithApproval: boolean
  notAcceptable: boolean
  opDivHead: boolean
  opDirector: boolean
  techDirector: boolean
  presDirector: boolean
  dueOccupation: boolean
  dueObjectConditions: boolean
  dueLossRatio: boolean
  dueOthers: boolean
}

/** `GET /api/nbfacin/kasus/{caseId}/scoring` (tiket 54, kontrak usulan sesi 0f). */
export interface TampilanScoring {
  /** Tahap marketing & bukan on-going (Pega "Write"); false = baca-saja ("Read"). */
  bolehUbah: boolean
  insuredName: string
  /** "New Business" / "Renewal". */
  status: string
  /** Period (kabel DD-MM-YYYY). */
  mulai: string
  selesai: string
  /** [Top Risk, Dominant Risk]. */
  kolom: KolomScoring[]
  finalScore: string
  /** "Good" / "Above Average" / "Average" / "Below Average" / "Poor" / "". */
  noteFinalScore: string
  hasil: HasilScoring
  /** Pilihan Dominant Risk: lokasi (ObjectNo) + okupasi pertamanya. */
  lokasiDominan: { indexLocation: string; label: string }[]
  pesan: string[]
}

const jalurScoring = (caseId: string) => `/api/nbfacin/kasus/${encodeURIComponent(caseId)}/scoring`

const badanScoring = (isi: TampilanScoring) => ({
  kolom: isi.kolom.map((k) => ({ indexLocation: k.indexLocation, pilihan: k.pilihan, persenLossRatio: k.persenLossRatio })),
})

export function ambilScoring(caseId: string): Promise<TampilanScoring> {
  return minta<TampilanScoring>(jalurScoring(caseId))
}

/** Hitung ulang tanpa simpan (`SetScore_Act` / `SetLossRatio_Act` / `ScoringResult`, isi otomatis kolom). */
export function hitungScoring(caseId: string, isi: TampilanScoring): Promise<TampilanScoring> {
  return minta<TampilanScoring>(`${jalurScoring(caseId)}/hitung`, { metode: 'POST', badan: badanScoring(isi) })
}

/** Save tab Scoring Risk. */
export function simpanScoring(caseId: string, isi: TampilanScoring): Promise<TampilanScoring> {
  return minta<TampilanScoring>(jalurScoring(caseId), { metode: 'PUT', badan: badanScoring(isi) })
}

/** Satu baris grid Loss Record tingkat case - `.OfferFacIn.ListCauseOfLoss(n)` (kelas Data-CauseOfLoss, tiket 55). */
export interface BarisLossRecord {
  /** `.UWYEAR` - DIHITUNG backend dari Date of Loss (`Get5YearAgo`): tahun awal jendela ke-1..5 sebelum periode polis. */
  uwYear: string
  classOfBusiness: string
  /** `.DOL` - kabel `DD-MM-YYYY`; kosong -> backend mengisi Begin date - 1 tahun (`SaveLossRecord` 7.6). */
  dateOfLoss: string
  /** `.CauseOfLoss` - backend menyimpan huruf besar (`@toUpperCase`). */
  causeOfLoss: string
  currency: string
  /** `.LOSS_AMOUNT` - teks desimal. */
  lossAmount: string
  /** `.LOSS_AMOUNT_IDR` - DIHITUNG backend: IDR = Loss Amount; lain = × kurs TREATYEXCHANGEYEARLY per Begin date. */
  lossAmountIdr: string
}

/** Satu baris "Total Loss" per Class Of Business + Currency (`TotalLossCOB`, `SaveLossRecord` langkah 8). */
export interface TotalLossRecord {
  classOfBusiness: string
  currency: string
  lossAmount: string
  lossAmountIdr: string
}

/** Loss ratio (persen, teks desimal) - `.OfferFacIn.LossRatio{1,3,5}YearPercent` / `.NetLossRatio{1,3,5}Year`. */
export interface RasioLossRecord {
  gross1: string
  gross3: string
  gross5: string
  net1: string
  net3: string
  net5: string
}

/** `GET /api/nbfacin/kasus/{caseId}/loss-record` (tiket 55, kontrak usulan sesi 0f). */
export interface TampilanLossRecord {
  /** `.OfferFacIn.PremiPAR` - bawaan backend bila belum diisi (keputusan 2: Σ premi × kurs Begin date). */
  premiPar: string
  premiMb: string
  premiPl: string
  premiOther: string
  /** Urut Date of Loss terbaru dulu (`SaveLossRecord` langkah 9). */
  baris: BarisLossRecord[]
  totalCob: TotalLossRecord[]
  rasio: RasioLossRecord
  pesan: string[]
}

const jalurLossRecord = (caseId: string) => `/api/nbfacin/kasus/${encodeURIComponent(caseId)}/loss-record`

const badanLossRecord = (isi: TampilanLossRecord) => ({
  premiPar: isi.premiPar,
  premiMb: isi.premiMb,
  premiPl: isi.premiPl,
  premiOther: isi.premiOther,
  baris: isi.baris.map((b) => ({
    classOfBusiness: b.classOfBusiness,
    dateOfLoss: b.dateOfLoss,
    causeOfLoss: b.causeOfLoss,
    currency: b.currency,
    lossAmount: b.lossAmount,
  })),
})

export function ambilLossRecord(caseId: string): Promise<TampilanLossRecord> {
  return minta<TampilanLossRecord>(jalurLossRecord(caseId))
}

/** Hitung ulang tanpa simpan (Uw Year, Loss Amount IDR, total, loss ratio). */
export function hitungLossRecord(caseId: string, isi: TampilanLossRecord): Promise<TampilanLossRecord> {
  return minta<TampilanLossRecord>(`${jalurLossRecord(caseId)}/hitung`, { metode: 'POST', badan: badanLossRecord(isi) })
}

/** Save Data tab Loss Record - daftar diganti utuh, dihitung ulang (keputusan 2). */
export function simpanLossRecord(caseId: string, isi: TampilanLossRecord): Promise<TampilanLossRecord> {
  return minta<TampilanLossRecord>(jalurLossRecord(caseId), { metode: 'PUT', badan: badanLossRecord(isi) })
}

/** `GET /api/nbfacin/class-of-business/saran?q=` - `BUSINESS.NOTE` semua grup (`BrowseBusiness_RD`, nilai `.Note`). */
export function saranClassOfBusiness(q: string): Promise<{ baris: SaranAkumulasi[] }> {
  return minta<{ baris: SaranAkumulasi[] }>('/api/nbfacin/class-of-business/saran', { kueri: { q } })
}

/** Satu penerima email korespondensi - `.RecipientList(n)` (tiket 56). */
export interface PenerimaKorespondensi {
  /** `.Recipient` - nama akun terpilih, atau teks bebas. */
  recipient: string
  /** `.RecipientMail` - diisi dari akun terpilih (`M_LOGIN_GO.EMAIL`) atau diketik. */
  recipientMail: string
}

/** Satu email korespondensi - `.OfferFacIn.CorrespondenceList(n)` (tiket 56, kontrak usulan sesi 0f). */
export interface Korespondensi {
  /** Kunci stabil dari server. */
  no: string
  /** `.Date` - waktu Add, WIB, kabel `DD-MM-YYYY HH:MM`. */
  tanggal: string
  /** `.Sender` / `.SenderMail` - akun login saat Add (`M_LOGIN_GO.NAME` / `EMAIL`); baca-saja. */
  sender: string
  senderMail: string
  subject: string
  /** `.Message` - HTML yang SUDAH dibersihkan server (daftar tag `TAG_HTML_KORESPONDENSI`, tanpa atribut). */
  message: string
  penerima: PenerimaKorespondensi[]
  /** Status: false = "Not Sent" (`.IsEditable` 1, dapat diubah), true = "Sent" (terkunci). */
  terkirim: boolean
}

/** `GET /api/nbfacin/kasus/{caseId}/korespondensi`. */
export interface TampilanKorespondensi {
  daftar: Korespondensi[]
  pesan: string[]
}

/** Isi yang dikirim Save / Send Email. */
export interface IsiKorespondensi {
  subject: string
  message: string
  penerima: PenerimaKorespondensi[]
}

const jalurKorespondensi = (caseId: string) => `/api/nbfacin/kasus/${encodeURIComponent(caseId)}/korespondensi`

export function ambilKorespondensi(caseId: string): Promise<TampilanKorespondensi> {
  return minta<TampilanKorespondensi>(jalurKorespondensi(caseId))
}

/** Add (`CorrespondenceAdd`): draf baru - tanggal dan pengirim diisi server dari akun login. */
export function tambahKorespondensi(caseId: string): Promise<Korespondensi> {
  return minta<Korespondensi>(jalurKorespondensi(caseId), { metode: 'POST' })
}

/** Save editor (`SaveFacIn_Act`) - 409 bila sudah terkirim. */
export function simpanKorespondensi(caseId: string, no: string, isi: IsiKorespondensi): Promise<Korespondensi> {
  return minta<Korespondensi>(`${jalurKorespondensi(caseId)}/${encodeURIComponent(no)}`, { metode: 'PUT', badan: isi })
}

/** Remove - hanya draf (409 bila sudah terkirim). */
export function hapusKorespondensi(caseId: string, no: string): Promise<void> {
  return minta<void>(`${jalurKorespondensi(caseId)}/${encodeURIComponent(no)}`, { metode: 'DELETE' })
}

/**
 * Send Email (`SendCorrespondence`, dibetulkan): isi disimpan dulu, lalu dikirim lewat SMTP konfigurasi - To = penerima,
 * CC + Reply-To = pengirim. Berhasil -> terkirim (terkunci). SMTP belum dikonfigurasi (503) / gagal (502) -> isi tetap
 * tersimpan, status tetap "Not Sent".
 */
export function kirimKorespondensi(caseId: string, no: string, isi: IsiKorespondensi): Promise<Korespondensi> {
  return minta<Korespondensi>(`${jalurKorespondensi(caseId)}/${encodeURIComponent(no)}/kirim`, { metode: 'POST', badan: isi })
}

/** `GET /api/nbfacin/korespondensi/penerima/saran?q=` - akun aktif ber-EMAIL: id = LOGIN_ID, label = NAME, ekstra = EMAIL. */
export function saranPenerimaKorespondensi(q: string): Promise<{ baris: SaranAkumulasi[] }> {
  return minta<{ baris: SaranAkumulasi[] }>('/api/nbfacin/korespondensi/penerima/saran', { kueri: { q } })
}

/** Kode Email Type tahap marketing (`.EmailType`, `SaveViewSuggest`): 1 Accept, 2 Reject, 4 Banding, 7 Decline (tiket 57). */
export type KodeEmailType = '1' | '2' | '4' | '7'

/** Satu baris grid riwayat akseptasi (`TempViewSuggest`: Date · PIC · Approval · Status · Comment), terbaru dulu. */
export interface RiwayatAkseptasi {
  /** Kabel `DD-MM-YYYY HH:MM` WIB. */
  tanggal: string
  pic: string
  /** Label Approval (Accept / Reject / Ask / Banding / Decline / Revise / …), sudah dipetakan server. */
  approval: string
  /** `.IsCedingConfirm` apa adanya (Offer / Binding / Accepted / Policy). */
  status: string
  comment: string
}

/** Satu baris T/C popup Cause of Decline (`POOLDATA.CAUSEOFDECLINE` `DESCRIPTION IS NULL`) + pilihan Option-nya. */
export interface BarisSebabTc {
  termCondition: string
  /** Pilihan kolom Option (`.Description1`) untuk T/C ini; kosong = tidak ada daftar. */
  opsi: string[]
}

/** Pilihan pengguna atas satu baris T/C. */
export interface PilihanSebabTc {
  termCondition: string
  /** Kolom centang: 'I' (`.CheckBox1`), 'II' (`.CheckBox2`), '' = tidak dicentang. */
  prioritas: '' | 'I' | 'II'
  description1: string
  description2: string
  /** Teks desimal persen (5 desimal di Pega); '' = kosong. */
  pctRnm: string
  pctSob: string
}

/** `GET /api/nbfacin/kasus/{caseId}/submit` (tiket 57, kontrak usulan sesi 0f). */
export interface TampilanSubmit {
  /**
   * Submit boleh: posisi case tahap marketing Offer DAN akun login memegang workbasket posisi itu (Pega: tombol aktif bila
   * `TempEmail.CARI28 == .PositionNote`, CARI28 = workbasket operator yang sama dengan PositionNote).
   */
  bolehSubmit: boolean
  /** Sebab tidak boleh (tampil di layar); '' bila boleh. */
  alasan: string
  riwayat: RiwayatAkseptasi[]
  /** Pilihan "Banding to" (`GetAksepBanding_SQL`); kosong bila case belum pernah ke approver non-marketing. */
  opsiBanding: string[]
  sebabTc: BarisSebabTc[]
  /** Daftar Non - T/C (`TERMCONDITION2 IS NOT NULL`). */
  sebabNonTc: string[]
  /**
   * % RNM otomatis saat T/C "Rate" / "R/I Comm, Deduction" dicentang prioritas I (`SetCODRate`): rate =
   * (Total Net Rate / 10) × (100 − RI Commission) / 100 dibulatkan 4 desimal; riComm = RI Commission. '' = tidak tersedia.
   */
  pctRnmOtomatis: { rate: string; riComm: string }
}

/** Isi Submit / periksa. */
export interface IsiSubmit {
  emailType: KodeEmailType
  comment: string
  bandingTo: string
  /** Decline: pilihan T/C (hanya baris yang dicentang) dan Non - T/C terpilih ('' = tidak ada). */
  sebabTc: PilihanSebabTc[]
  sebabNonTc: string
}

/** Jawaban Submit / periksa. */
export interface HasilSubmit {
  /** Pesan proteksi (`Protection_Act` dkk., verbatim). Kosong = lolos / berhasil. */
  pesan: string[]
  /** Sesudah Submit berhasil: posisi baru case dan status teksnya (`NBStatus`). */
  position: string
  nbStatus: string
}

const jalurSubmit = (caseId: string) => `/api/nbfacin/kasus/${encodeURIComponent(caseId)}/submit`

export function ambilSubmit(caseId: string): Promise<TampilanSubmit> {
  return minta<TampilanSubmit>(jalurSubmit(caseId))
}

/** Proteksi Accept tanpa submit (Pega `Protection_Act` saat radio dipilih); Reject / Banding / Decline -> pesan kosong. */
export function periksaSubmit(caseId: string, emailType: KodeEmailType): Promise<HasilSubmit> {
  return minta<HasilSubmit>(`${jalurSubmit(caseId)}/periksa`, { metode: 'POST', badan: { emailType } })
}

/**
 * Submit tahap marketing (flow action `InwardFacultative` + `SetValidateDate_PostAct` + `Decision3`): Accept -> Team Leader;
 * Reject -> tetap di marketing; Banding -> jabatan Banding to; Decline -> selesai. Proteksi gagal: 200 + `pesan` (diminta ke c3)
 * atau 422 `{"galat": pesan digabung "; "}` - keduanya ditangani `BlokEmailType` (`pesanProteksiDari`).
 */
export function ajukanSubmit(caseId: string, isi: IsiSubmit): Promise<HasilSubmit> {
  return minta<HasilSubmit>(jalurSubmit(caseId), { metode: 'POST', badan: isi })
}
