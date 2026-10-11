// Tampilan polis NB Fac In HANYA-BACA Marine Cargo (perintah work owner 11-10-2026: "tampilan view polis yang belum
// tergenerate tolong diperbaiki. cek tampilan nya di nb facin!"). Rupa = tab Object / Coverage / Spreading Marine Cargo
// layar Inward Facultative `modul/nbfacin/frontend/components/{TabObjectCargo,TabCoverageCargo,TabSpreadingCargo}.tsx`
// DISALIN (pola, bukan impor), dibaca dari halaman polis klaim `OfferFacIn.CargoList`. Label VERBATIM
// `modul/nbfacin/frontend/labels.ts`. Berkas ini hanya fungsi murni; komponennya `components/PolisNBCargo.tsx`.
//
// Dua bentuk halaman: polis Pega (JSON_POLIS) menaruh data pelayaran di `CargoList(n).PolicyData.{SailDate, BLNumber,
// InvoiceNumber, InvoiceDate, MaxAge, MaxAgeDetail, Ship, LC, SurveyAgent}` (`InputDtlCargo_FacIn.xml`); polis tabel datar
// nbfacin melipat `.PolicyData` ke baris cargo (loader V-24b: `CargoList(n).SailDate`, `.Ship`, `.LC`). Pembaca memakai
// keduanya. Total yang nbfacin HITUNG (Grand Total per cargo `totalGrossCargo`, total spreading `totalCargo`) dihitung
// ulang dengan rumus yang sama, eksak tanpa float.

import type { Baris, Halaman } from './api'
import type { ItemDasarNB } from './polisNBBayar'
import {
  AKAR_NB,
  NOL,
  SERATUS,
  anakNB,
  bagiE,
  bandingE,
  benarNB,
  bulatE,
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
  type TotalE,
} from './polisNB'

export const CARGO_NB = AKAR_NB + '.CargoList'
export const coverageCargoNB = (cargo: number) => anakNB(CARGO_NB, cargo, 'CoverageList')

/** Nilai medan cargo: baris datar dulu, lalu `.PolicyData.<k>` Pega. */
export const medanCargo = (b: Baris, k: string) => b[k] || b['PolicyData.' + k] || ''

/** Daftar Additional Ship cargo n (`.Ship.AdditionalShip` datar / `.PolicyData.Ship.AdditionalShip` Pega). */
export function kapalTambahanNB(h: Halaman, cargo: number): Baris[] {
  const datar = daftarNB(h, anakNB(CARGO_NB, cargo, 'Ship.AdditionalShip'))
  return datar.length > 0 ? datar : daftarNB(h, anakNB(CARGO_NB, cargo, 'PolicyData.Ship.AdditionalShip'))
}

// ---------------------------------------------------------------- label (VERBATIM nbfacin/labels.ts)

/** `GRID_CARGO`. */
export const GRID_CARGO_NB = ['No.', 'Top Risk', 'SailDate', 'Goods', 'Ship Name', 'Conveyance'] as const
/** `FORM_CARGO` (tombol Choose / Select / Add / Delete tidak dibawa). */
export const FORM_CARGO_NB = {
  topRisk: 'Top Risk',
  trading: 'Trading',
  goods: 'Goods',
  conveyance: 'Conveyance',
  packing: 'Packing',
  grupAgent: 'Survey Agent',
  fullName: 'Full Name',
  address: 'Address',
  city: 'City/Town',
  country: 'Country Code',
  grupLc: 'LC',
  lcStatus: 'LC Status',
  lcBank: 'LC Bank',
  lcNumber: 'LC Number',
  lcCondition: 'LC Condition',
  lcRemark: 'LC Remark',
  lcEndDate: 'LC Expired Date',
  grupKapal: 'Ship Data',
  shipName: 'Ship Name',
  shipType: 'Ship Type',
  grt: 'GRT',
  nrt: 'NRT',
  dwt: 'DWT',
  age: 'Age',
  flagName: 'Flag Name',
  construction: 'Construction',
  yMake1: 'Year Make 1',
  remark: 'Remark',
  sailDate: 'Sail Date',
  fromRute: 'From Rute',
  toRute: 'To Rute',
  blNumber: 'No BL',
  invoiceNumber: 'Invoice Number',
  invoiceDate: 'Invoice Date',
  maxAge: 'Max Age ID',
  maxAgeDetail: 'Max Age Detail',
  grupKapalTambahan: 'Additional Ship',
  kolomKapalTambahan: ['Register ID', 'Ship Name', 'Voyage No'],
} as const
/** `TEKS_CARGO` / `TEKS_ANEKA.belumDiisi`. */
export const TEKS_CARGO_NB = { kosong: 'Belum ada objek cargo.', belumDiisi: '(belum diisi)' } as const
/** `GRID_CARGO_COVERAGE`. */
export const GRID_CARGO_COVERAGE_NB = ['No.', 'Trading', 'Goods', 'Packing', 'Conveyance'] as const
/** `GRID_COVERAGE_CARGO`. */
export const GRID_COVERAGE_CARGO_NB = {
  kolom: ['Coverage Note', 'Coverage Initial'],
  total: ['Currency', 'Grand Total Gross TSI', 'Grand Total Gross Premium', 'Grand Total Gross Rate'],
} as const
/** `FORM_COVERAGE_CARGO`. */
export const FORM_COVERAGE_CARGO_NB = {
  coverageNote: 'Coverage Note',
  coverageInitial: 'Coverage Initial',
  currency: 'Currency',
  rate: 'Rate (%)',
  limitOfLiability: 'Limit of Liability',
  discountPercentage: 'Discount (%)',
  tsi: 'TSI',
  premium: 'Premium',
  minPremium: 'Min Premium',
  discount: 'Discount',
  grupWarranty: 'Warranty Extension',
  kolomWarranty: ['Coverage', 'Deductible'],
  deskripsi: 'Deskripsi',
} as const
/** `SPREADING_CARGO`. */
export const SPREADING_CARGO_NB = {
  kolomCargo: ['No.', 'Trading', 'Goods', 'Packing', 'Conveyance'],
  kolomCoverage: ['Coverage', 'Standard Rate (%)', 'Gross Premium', 'Premium Nusantara Re'],
} as const

// ---------------------------------------------------------------- Grand Total (totalGrossCargo nbfacin)

export interface TotalGrossCargoNB {
  currency: string
  tsi: string
  premium: string
  rate: string
}

/** Grand Total per mata uang coverage satu cargo: Σ TSI, Σ Premium, Rate = Premium / TSI x 100 (20 desimal), 8 desimal. */
export function totalGrossCargoNB(cov: readonly Baris[]): TotalGrossCargoNB[] {
  const baris: { currency: string; tsi: Eksak; premium: Eksak }[] = []
  for (const c of cov) {
    const u = c['Currency.Name'] ?? ''
    let r = baris.find((x) => x.currency === u)
    if (!r) baris.push((r = { currency: u, tsi: NOL, premium: NOL }))
    r.tsi = tambahE(r.tsi, eksak(c['TSI']))
    r.premium = tambahE(r.premium, eksak(c['Premium']))
  }
  return baris.map((r) => ({
    currency: r.currency,
    tsi: teksE(bulatE(r.tsi, 8)),
    premium: teksE(bulatE(r.premium, 8)),
    rate: r.tsi.n === 0n ? '0' : teksE(bulatE(kaliE(bagiE(r.premium, r.tsi, 20), SERATUS), 8)),
  }))
}

// ---------------------------------------------------------------- Spreading (spreadingcargo.go nbfacin)

/**
 * Port `totalCargo` nbfacin (`SumTSIPremiSpreadedRNM_MARINECARGO_Act` langkah 1), lalu `ringkasSpreadingNB`: total per mata
 * uang coverage (urut kemunculan) x jenis treaty; cargo ber-FlagDelete "1" dilewati; LoL = LimitofLiability x PS / 100
 * hanya cargo PERTAMA, spread x share / 100; TSI top risk = TSISpreaded cargo Top Risk; share beda -> 0.
 */
export function hitungSpreadingCargoNB(h: Halaman): HasilSpreadingNB {
  const ps = eksak(nilaiNB(h, 'PercentShare'))
  const cargo = daftarNB(h, CARGO_NB)
  const urut: string[] = []
  cargo.forEach((_, i) => {
    for (const c of daftarNB(h, coverageCargoNB(i + 1))) {
      const u = c['Currency.Name'] ?? ''
      if (!urut.includes(u)) urut.push(u)
    }
  })
  const namaTreaty = new Map<string, string>()
  const hasil: TotalE[] = []
  for (const mu of urut) {
    const awal = hasil.length
    cargo.forEach((x, i) => {
      if ((x['FlagDelete'] ?? '').trim() === '1') return
      daftarNB(h, coverageCargoNB(i + 1)).forEach((c, k) => {
        if ((c['Currency.Name'] ?? '') !== mu) return
        const tsiLoL = i === 0 ? persenDari(eksak(c['LimitofLiability']), ps) : NOL
        for (const b of daftarNB(h, anakNB(coverageCargoNB(i + 1), k + 1, 'SpreadingList'))) {
          const jenis = b['TreatyType'] ?? ''
          if (!namaTreaty.has(jenis) && (b['TreatyName'] ?? '') !== '') namaTreaty.set(jenis, b['TreatyName'] ?? '')
          const share = eksak(b['SharePercentage'])
          const lolSpread = persenDari(tsiLoL, share)
          const top = benarNB(x['IsTopRisk']) ? eksak(b['TSISpreaded']) : NOL
          let n = -1
          for (let j = awal; j < hasil.length; j++) {
            if (hasil[j]!.treatyType === jenis) {
              n = j
              break
            }
          }
          if (n < 0) {
            hasil.push({
              currency: mu,
              treatyType: jenis,
              share,
              claimSpreaded: top,
              tsiSpreaded: eksak(b['TSISpreaded']),
              claimEstimation: lolSpread,
              claimAmountIDR: NOL,
              premiumSpreaded: eksak(b['PremiumSpreaded']),
            })
            continue
          }
          const t = hasil[n]!
          t.premiumSpreaded = tambahE(t.premiumSpreaded, eksak(b['PremiumSpreaded']))
          t.claimEstimation = tambahE(t.claimEstimation, lolSpread)
          t.tsiSpreaded = tambahE(t.tsiSpreaded, eksak(b['TSISpreaded']))
          t.claimSpreaded = tambahE(t.claimSpreaded, top)
          if (bandingE(share, t.share) !== 0) t.share = NOL
        }
      })
    })
  }
  // total per mata uang x treaty sudah unik: satu "lokasi" tanpa top risk = masukan ringkasan nbfacin apa adanya
  return ringkasSpreadingNB([{ top: false, total: hasil }], namaTreaty, false)
}

/** Pseudo-item `ringkasanRNM` Marine Cargo (`kasusSpreadingCargo`): satu item per coverage cargo ber-FlagDelete != 1, mata uang
 * coverage, TSI = TSILiability. */
export function itemDasarCargoNB(h: Halaman): ItemDasarNB[] {
  const out: ItemDasarNB[] = []
  daftarNB(h, CARGO_NB).forEach((x, i) => {
    if ((x['FlagDelete'] ?? '').trim() === '1') return
    for (const c of daftarNB(h, coverageCargoNB(i + 1))) {
      out.push({
        currency: c['Currency.Name'] ?? '',
        tsi: c['TSILiability'] ?? '',
        coverage: [{ premium: c['Premium'] ?? '', discount: c['Discount'] ?? '' }],
      })
    }
  })
  return out
}
