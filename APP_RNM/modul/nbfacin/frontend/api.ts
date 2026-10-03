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
