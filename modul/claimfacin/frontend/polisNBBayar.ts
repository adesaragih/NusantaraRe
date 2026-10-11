// Tampilan polis NB Fac In HANYA-BACA - tab Premium Deduction / Brokerage Fee, Payment, dan Inw Fac Cedant Panels untuk
// FIRE, keluarga Aneka, dan Marine Cargo (perintah work owner 11-10-2026 "tampilan view polis yang belum tergenerate tolong
// diperbaiki. cek tampilan nya di nb facin!"). Rupa = `modul/nbfacin/frontend/components/{TabDeduction,TabPayment,
// TabCedant}.tsx` (pola, bukan impor). Label VERBATIM `modul/nbfacin/frontend/labels.ts`.
//
// Keputusan work owner 11-10-2026: Premium Deduction / Payment = NILAI TERSIMPAN tabel NB (T_GENERAL_POLIS PAY_*,
// T_CURRENCYLIST PAY_* / SUM_TOTAL_PAYMENT, T_LISTINSTALLMENT), bukan rumus NB yang menghitung ulang. Kepala PPH / PPN tanpa
// tarif (tarif master TAX terbaru dapat berbeda dari tarif saat nilai disimpan; nbfacin `kepalaPajak` pun tanpa tarif bila
// kosong). Rincian Cedant per mata uang TIDAK disimpan NB (T_CEDING_CURRENCYLIST hanya data impor Pega) - dihitung seperti
// nbfacin (`ringkasanRNM` dasar x ShareCeding / 100).

import type { Baris, Halaman } from './api'
import {
  AKAR_NB,
  LOKASI_NB,
  NOL,
  SERATUS,
  anakNB,
  bagiE,
  coverageNB,
  daftarNB,
  eksak,
  itemNB,
  kaliE,
  nilaiNB,
  persenDari,
  tambahE,
  teksE,
  type Eksak,
} from './polisNB'

// ---------------------------------------------------------------- label (VERBATIM nbfacin/labels.ts)

/** `DEDUCTION`. */
export const DEDUCTION_NB = {
  pctPremiumDeduction: '% Premium Deduction',
  pctBrokerageFee: '% Brokerage Fee',
  excludePpn: 'Exclude PPN ?',
  kolomMataUang: 'Currency',
  totalPremiumDeduction: 'Total Premium Deduction',
  totalPremiumDeduction2: 'Total Premium Deduction 2',
  totalBrokerageFee: 'Total Brokerage Fee',
  pph: 'PPH',
  ppn: 'PPN',
  kolomLokasi: ['No.', 'Object Name', 'Location'],
  kolomItem: ['Object Item Type', 'Currency', 'TSI Object Item', 'Total Gross Premium', 'Total Premium RNM'],
  kolomCoverage: ['Coverage', '‰ Standard Rate', 'Gross Premium', 'Premium RNM'],
  detailCoverage: {
    tsiLiability: 'Gross TSI Liability',
    premium: 'Gross Premium',
    tsiNusantaraRe: 'TSI Nusantara Re',
    premiNusantaraRe: 'Premium Nusantara Re',
  },
  kolomSpreading: ['Type Treaty', '% Share', 'TSI Spreaded', 'Limit Of Liability', 'Premium Spreaded'],
} as const
/** `DEDUCTION_ANEKA`. */
export const DEDUCTION_ANEKA_NB = {
  kolomLokasi: ['No', 'Location'],
  kolomOkupasi: ['Occupation Name'],
  kolomObjek: ['Object Name', 'Currency', 'TSI'],
  kolomObjekGolf: ['Object Name'],
  kolomCoverage: ['Coverage', 'Gross Premium'],
  detailCoverage: {
    tsi: 'TSI Liability',
    tsiNusantaraRe: 'TSI Nusantara Re',
    premium: 'Gross Premium',
    premiNusantaraRe: 'Premium Nusantara Re',
  },
} as const
/** `DEDUCTION_CARGO`. */
export const DEDUCTION_CARGO_NB = {
  kolomCargo: ['No.', 'Trading', 'Goods', 'Packing', 'Conveyance'],
  kolomCoverage: ['Coverage', 'TSI', 'Gross Premium'],
} as const
/** `PAYMENT`. */
export const PAYMENT_NB = {
  judul: 'Payment Information',
  installment: 'Installment',
  kolomMataUang: ['Currency', 'Premium RNM', 'Total Payment'],
  judulCicilan: 'Installment Data Information',
  no: 'No',
  dueDate: 'Due Date',
  persen: '% Installment',
  premium: 'Premium RNM',
  deduction: 'Deduction',
  deduction2: 'Deduction 2',
  brokerage: 'Total Brokerage Fee',
  pph: 'PPH',
  ppn: 'PPN',
  discount: 'Discount',
  totalPayment: 'Total Payment',
} as const
/** `TEKS_PAYMENT` (pesan Σ % Installment). */
export const TEKS_PAYMENT_NB = {
  persenLebih: "Total % Installment Can't Be More Than 100%",
  persenKurang: 'Total % Installment Must Be 100%',
} as const
/** `CEDANT` (detail baris + kolom mata uang) / `TEKS_CEDANT.mataUangTotal`. */
export const CEDANT_DETAIL_NB = {
  kolomCeding: 'Ceding',
  kolomShare: '% Share',
  dari: 'of',
  kolomMataUang: ['Currency', 'TSI Cedant', 'Premium Cedant'],
  mataUangTotal: 'IDR',
} as const

// ---------------------------------------------------------------- nilai tersimpan

/** Medan tingkat polis `OfferFacIn.PolicyData.Payment.<k>` (PAY_* T_GENERAL_POLIS). */
export const bayarPolisNB = (h: Halaman, k: string) => nilaiNB(h, 'PolicyData.Payment.' + k)

/** Baris CurrencyList (`OfferFacIn.CurrencyList`) - medan PAY_* di `.Policy.Payment.*`. */
export const MATA_UANG_NB = AKAR_NB + '.CurrencyList'
export const bayarUangNB = (b: Baris, k: string) => b['Policy.Payment.' + k] ?? ''

/** ListInstallment mata uang n: tabel datar `CurrencyList(n).ListInstallment`, Pega `.Policy.Payment.ListInstallment`. */
export function cicilanNB(h: Halaman, n: number): Baris[] {
  const datar = daftarNB(h, anakNB(MATA_UANG_NB, n, 'ListInstallment'))
  return datar.length > 0 ? datar : daftarNB(h, anakNB(MATA_UANG_NB, n, 'Policy.Payment.ListInstallment'))
}

/** Pesan Σ % Installment satu mata uang (`pesanPersenCicilan` nbfacin / SetValidateInstallment_Act). */
export function pesanPersenCicilanNB(cicilan: readonly Baris[]): string {
  if (cicilan.length === 0) return ''
  const sah = cicilan.map((c) => (c['InstallmentPercentage'] ?? '').trim()).filter((s) => /^[+-]?\d+(\.\d+)?$/.test(s))
  const total = sah.map(eksak).reduce(tambahE, NOL)
  const banding = total.n * 10n ** 0n - 100n * 10n ** BigInt(total.s)
  if (banding > 0n) return TEKS_PAYMENT_NB.persenLebih
  return banding === 0n ? '' : TEKS_PAYMENT_NB.persenKurang
}

// ---------------------------------------------------------------- Cedant (ringkasanRNM nbfacin)

/** Satu item dasar `ringkasanRNM` (pseudo-item lini: FIRE PropertyItem, Aneka `kasusSpreadingAneka`, Cargo per coverage). */
export interface ItemDasarNB {
  currency: string
  /** TSI item (kosong = 0). */
  tsi: string
  /** Premium / Discount coverage item (premi NR = PS x Premium / 100 - PS x Discount / 100, `hitungNR`). */
  coverage: { premium: string; discount: string }[]
}

export interface DasarCedantNB {
  currency: string
  tsiRnm: string
  premiumRnm: string
  tsiGross: string
  premiumGross: string
}

/** `ringkasanRNM` dasar: per mata uang item (urut kemunculan; tanpa mata uang dilewati); kosong bila % Share RNM kosong. */
export function dasarCedantNB(h: Halaman, items: readonly ItemDasarNB[]): DasarCedantNB[] {
  const psTeks = nilaiNB(h, 'PercentShare').trim()
  if (psTeks === '') return []
  const ps = eksak(psTeks)
  const per: { currency: string; tsi: Eksak; premi: Eksak }[] = []
  for (const it of items) {
    if (it.currency === '') continue
    let r = per.find((x) => x.currency === it.currency)
    if (!r) per.push((r = { currency: it.currency, tsi: NOL, premi: NOL }))
    for (const c of it.coverage) {
      const nr = tambahE(persenDari(ps, eksak(c.premium)), kaliE(persenDari(ps, eksak(c.discount)), { n: -1n, s: 0 }))
      r.premi = tambahE(r.premi, nr)
    }
    r.tsi = tambahE(r.tsi, eksak(it.tsi))
  }
  return per.map((r) => ({
    currency: r.currency,
    tsiRnm: teksE(persenDari(r.tsi, ps)),
    premiumRnm: teksE(r.premi),
    tsiGross: teksE(r.tsi),
    premiumGross: ps.n === 0n ? '' : teksE(bagiE(kaliE(r.premi, SERATUS), ps, 20)),
  }))
}

/** CurrencyList satu cedant (`mataUangCedant` nbfacin): ShareCeding x dasar (Share RNM: RNM; Gross: gross) / 100. */
export function mataUangCedantNB(
  share: string,
  type: string,
  dasar: readonly DasarCedantNB[],
): { currency: string; tsi: string; premium: string }[] {
  const s = share.trim()
  if (!/^[+-]?\d+(\.\d+)?$/.test(s) || (type !== '0' && type !== '1')) return []
  const kali = (x: string) => (x === '' ? '' : teksE(persenDari(eksak(s), eksak(x))))
  return dasar.map((d) => ({
    currency: d.currency,
    tsi: kali(type === '1' ? d.tsiRnm : d.tsiGross),
    premium: kali(type === '1' ? d.premiumRnm : d.premiumGross),
  }))
}

/** "of <x>" baris cedant (`teksShareOfCeding` nbfacin): Gross 100%, Share RNM = % Share RNM polis, selain itu tersimpan. */
export function shareOfCedingNB(type: string, percentShare: string, tersimpan: string): string {
  if (type === '0') return '100%'
  if (type === '1') {
    const p = percentShare.trim()
    return (/^\d+(\.\d+)?$/.test(p) ? p.replace(/(\.\d*?)0+$/, '$1').replace(/\.$/, '') : p) + '%'
  }
  return tersimpan
}

/** Pseudo-item `ringkasanRNM` FIRE: setiap PropertyItem lokasi (mata uang item, TSI Object Item, coverage item). */
export function itemDasarFireNB(h: Halaman): ItemDasarNB[] {
  const out: ItemDasarNB[] = []
  daftarNB(h, LOKASI_NB).forEach((_, l) => {
    daftarNB(h, itemNB(l + 1)).forEach((it, j) => {
      out.push({
        currency: it['Currency'] ?? '',
        tsi: it['TSIObjectItem'] ?? '',
        coverage: daftarNB(h, coverageNB(l + 1, j + 1)).map((c) => ({
          premium: c['Premium'] ?? '',
          discount: c['Discount'] ?? '',
        })),
      })
    })
  })
  return out
}
