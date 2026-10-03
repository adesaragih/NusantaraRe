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
  /** Tampil-saja (diisi fitur Change SOB / Change Ceding Co / Following, tahap berikutnya). */
  sourceOfBusiness: string
  cedingCoName: string
  groupName: string
  oldPolicyNumber: string
}

/** Medan General yang dikirim tombol Save for later (tanpa medan tampil-saja). */
export type SimpanGeneral = Omit<GeneralInward, 'sourceOfBusiness' | 'cedingCoName' | 'groupName' | 'oldPolicyNumber'>

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
