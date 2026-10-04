// Uji definisi medan layar - syarat VERBATIM section Pega (INVENTARIS bab 5;
// sel + wadah `pyContainerVisibleWhen` dibaca ulang 2026-10-03).

import { describe, expect, it } from 'vitest'

import type { Halaman } from './api'
import {
  MEDAN_ADMIN_UANG,
  MEDAN_ADMIN_UMUM,
  MEDAN_ATASAN_TOTAL,
  MEDAN_ATASAN_UANG,
  MEDAN_ATASAN_UMUM,
  medanTampil,
  saldoNegatif,
  type Medan,
} from './medan'

const P = 'PolicyTreatyIn.'
const hal = (nilai: Record<string, string>): Halaman => ({ nilai, daftar: {} })
const label = (ms: { label: string }[]) => ms.map((m) => m.label)
const tampil = (ms: Medan[], nilai: Record<string, string>) => label(medanTampil(ms, hal(nilai)))
const cari = (ms: Medan[], jalur: string, lbl?: string) => {
  const m = ms.find((x) => x.jalur === P + jalur && (lbl === undefined || x.label === lbl))
  if (!m) throw new Error(`medan ${jalur} tidak ada`)
  return m
}

describe('medan layar NB Treaty In - syarat tampil', () => {
  it('saldo: "Balance Due To You" bila BalanceDueTo < 0, "Due To Us" bila >= 0 (kosong = 0)', () => {
    expect(saldoNegatif(hal({ [P + 'BalanceDueTo']: '-0.01' }))).toBe(true)
    expect(saldoNegatif(hal({ [P + 'BalanceDueTo']: '-0' }))).toBe(false)
    expect(saldoNegatif(hal({}))).toBe(false)
    const neg = tampil(MEDAN_ADMIN_UANG, { [P + 'BalanceDueTo']: '-5' })
    expect(neg).toContain('Balance Due To You')
    expect(neg).not.toContain('Balance Due To Us')
    const pos = tampil(MEDAN_ADMIN_UANG, { [P + 'BalanceDueTo']: '5' })
    expect(pos).toContain('Balance Due To Us')
    expect(pos).toContain('Balance Before Tax')
  })

  it('admin FlagRetroTreaty tampil bila ClaimType != "XOL Retro" (DetailPolicyTreatyIn pyVisible)', () => {
    expect(tampil(MEDAN_ADMIN_UMUM, {})).toContain('FlagRetroTreaty')
    expect(tampil(MEDAN_ADMIN_UMUM, { [P + 'ClaimType']: 'XOL Retro' })).not.toContain('FlagRetroTreaty')
  })

  it('admin FlagPPH dan Type Tax di wadah ".ClaimType != XOL Retro"; Type Tax juga bila FlagPPH true', () => {
    expect(tampil(MEDAN_ADMIN_UMUM, {})).toContain('FlagPPH')
    expect(tampil(MEDAN_ADMIN_UMUM, {})).not.toContain('Type Tax')
    expect(tampil(MEDAN_ADMIN_UMUM, { [P + 'FlagPPH']: 'true' })).toContain('Type Tax')
    const retro = tampil(MEDAN_ADMIN_UMUM, { [P + 'FlagPPH']: 'true', [P + 'ClaimType']: 'XOL Retro' })
    expect(retro).not.toContain('FlagPPH')
    expect(retro).not.toContain('Type Tax')
  })

  it('Survey Report disembunyikan untuk NonProportional', () => {
    expect(tampil(MEDAN_ADMIN_UMUM, { 'Quotation.ProportionalType': 'NonProportional' })).not.toContain('Survey Report')
  })

  it('Quartal / YearOfQuartal / U/Y hanya Proportional, Layer* hanya NonProportional (wadah, kedua layar)', () => {
    for (const ms of [MEDAN_ADMIN_UMUM, MEDAN_ATASAN_UMUM]) {
      const prop = tampil(ms, { [P + 'QuotationData.ProportionalType']: 'Proportional' })
      expect(prop).toEqual(expect.arrayContaining(['.Quartal', 'YearOfQuartal', 'U/Y']))
      expect(prop).not.toContain('LayerType')
      const np = tampil(ms, { [P + 'QuotationData.ProportionalType']: 'NonProportional' })
      expect(np).not.toContain('.Quartal')
      expect(np).not.toContain('U/Y')
      expect(np).toEqual(expect.arrayContaining(['LayerType', 'Layer', 'LayerPartType', 'LayerPart']))
    }
  })

  it('bagian uang: admin bila IsNewPolicyNonProp != 1 && IsNewPolicyListFormat != 1; atasan bila IsNewPolicyNonProp != 1', () => {
    expect(tampil(MEDAN_ADMIN_UANG, {})).toContain('Premi Ogp')
    expect(tampil(MEDAN_ADMIN_UANG, { [P + 'IsNewPolicyNonProp']: '1' })).toEqual([])
    expect(tampil(MEDAN_ADMIN_UANG, { [P + 'IsNewPolicyListFormat']: '1' })).toEqual([])
    expect(tampil(MEDAN_ATASAN_UANG, { [P + 'IsNewPolicyListFormat']: '1' })).toContain('Premi Ogp')
    expect(tampil(MEDAN_ATASAN_UANG, { [P + 'IsNewPolicyNonProp']: '1' })).toEqual([])
    expect(tampil(MEDAN_ATASAN_TOTAL, { [P + 'IsNewPolicyNonProp']: '1' })).toEqual([])
  })

  it('atasan: Ceding Company, Marketing Officer, Type Tax NOTBLANK', () => {
    const kosong = tampil(MEDAN_ATASAN_UMUM, {})
    for (const l of ['Ceding Company', 'Marketing Officer', 'Type Tax']) expect(kosong).not.toContain(l)
    const isi = tampil(MEDAN_ATASAN_UMUM, {
      [P + 'CedingCoName']: 'UJI-CEDING',
      [P + 'MarketingOfficer']: 'UJI-MO',
      [P + 'TypeTax']: 'Inclusive',
    })
    for (const l of ['Ceding Company', 'Marketing Officer', 'Type Tax']) expect(isi).toContain(l)
  })

  it('atasan: No Polis dan Production Date di wadah "pyWorkPage.FlagViewPolicy = 1" (nol pengisi di korpus)', () => {
    expect(tampil(MEDAN_ATASAN_UMUM, {})).not.toContain('No Polis')
    expect(tampil(MEDAN_ATASAN_UMUM, {})).not.toContain('Production Date')
    expect(tampil(MEDAN_ATASAN_UMUM, { FlagViewPolicy: '1' })).toEqual(expect.arrayContaining(['No Polis', 'Production Date']))
  })

  it('elemen mati (1=2 / NEVER) tidak dibangun: Class Of Business dan Due To tidak ada di layar admin', () => {
    const semua = label(MEDAN_ADMIN_UMUM)
    expect(semua).not.toContain('Class Of Business')
    expect(semua).not.toContain('Due To Us / You')
  })
})

describe('medan layar NB Treaty In - dapat diisi', () => {
  it('AC 52: atasan tidak mengisi satu pun medan General / uang (DueTo, FlagPPH, No Offer Slip pyDisabled)', () => {
    const terbuka = [...MEDAN_ATASAN_UMUM, ...MEDAN_ATASAN_UANG, ...MEDAN_ATASAN_TOTAL].filter(
      (m) => m.jenis !== 'tampil' && !m.kunci,
    )
    expect(terbuka.map((m) => m.jalur)).toEqual([])
    expect(cari(MEDAN_ATASAN_UMUM, 'FlagPPH').kunci).toBe(true)
  })

  it('layar atasan: label VERBATIM (tanpa "(%)")', () => {
    expect(label(MEDAN_ATASAN_UANG)).toContain('Deduction In A (OGP)')
    expect(label(MEDAN_ATASAN_UANG)).not.toContain('(%) Deduction In A (OGP)')
  })
})

describe('medan layar NB Treaty In - aksi sel (action set XML, berurutan)', () => {
  it('persen/hasil OGP-ONP admin: Count*_Act(Data) LALU CountOGPONP_Act', () => {
    const harap: [string, string, string][] = [
      ['RiCommOgp', 'CountResult1', 'Pct'],
      ['ResultOgp1', 'CountResult1', 'Amount'],
      ['OveriddingCommOgp', 'CountResult2Ogp', 'Pct'],
      ['ResultOgp2', 'CountResult2Ogp', 'Amount'],
      ['RiCommOnp', 'CountResult1Onp', 'Pct'],
      ['ResultOnp1', 'CountResult1Onp', 'Amount'],
      ['OveriddingCommOnp', 'CountResult2Onp', 'Pct'],
      ['ResultOnp2', 'CountResult2Onp', 'Amount'],
    ]
    for (const [m, aksi, param] of harap) {
      expect(cari(MEDAN_ADMIN_UANG, m).aksi).toEqual([{ aksi, param }, { aksi: 'CountOGPONP' }])
    }
  })

  it('Gross* -> CalculatePremi_Act(Action); Premi/Claim/Salvage/Excess/Deduction -> CountOGPONP_Act; Outstanding tanpa aksi', () => {
    expect(cari(MEDAN_ADMIN_UANG, 'GrossPremium').aksi).toEqual([{ aksi: 'CalculatePremi', param: 'PREMIUM' }])
    expect(cari(MEDAN_ADMIN_UANG, 'GrossClaim').aksi).toEqual([{ aksi: 'CalculatePremi', param: 'CLAIM' }])
    for (const m of ['PremiOgp', 'PremiOnp', 'Claim', 'SalvageValue', 'ExcessLoss', 'Deduction1', 'Deduction2']) {
      expect(cari(MEDAN_ADMIN_UANG, m).aksi).toEqual([{ aksi: 'CountOGPONP' }])
    }
    expect(cari(MEDAN_ADMIN_UANG, 'OutstandingClaim').aksi).toBeUndefined()
  })

  // Audit silang putaran 3: action set sel `Section/DetailPolicyTreatyIn` -
  // `.EndDate` change -> refresh ProtectDate; `.FlagPPH` change -> runActivity
  // RemoveTypeTax_ACT; `.QuotationData.MOID` change -> refresh CheckDataMkt;
  // `.IDCurrency` change -> refresh SetCurrency_act(CURR=.IDCurrency).
  it('EndDate, FlagPPH, Marketing Officer, Currency admin: satu aksi XML masing-masing', () => {
    expect(cari(MEDAN_ADMIN_UMUM, 'EndDate').aksi).toEqual([{ aksi: 'ProtectDate' }])
    expect(cari(MEDAN_ADMIN_UMUM, 'FlagPPH').aksi).toEqual([{ aksi: 'RemoveTypeTax' }])
    expect(cari(MEDAN_ADMIN_UMUM, 'QuotationData.MOID').aksi).toEqual([{ aksi: 'CheckDataMkt' }])
    expect(cari(MEDAN_ADMIN_UMUM, 'IDCurrency').aksi).toEqual([{ aksi: 'SetCurrency' }])
  })
})

describe('medan layar NB Treaty In - penyajian (AC 85, K3, K14)', () => {
  const PERSEN = ['RiCommOgp', 'OveriddingCommOgp', 'RiCommOnp', 'OveriddingCommOnp']

  it('AC 85 + K3: setiap angka uang membawa kode mata uang - termasuk Deduction1/2 (pxCurrency)', () => {
    for (const ms of [MEDAN_ADMIN_UANG, MEDAN_ATASAN_UANG]) {
      for (const m of ms) {
        const persen = PERSEN.some((x) => m.jalur === P + x)
        expect(m.mataUang === undefined, `${m.jalur} mataUang`).toBe(persen)
      }
    }
    expect(cari(MEDAN_ADMIN_UANG, 'Deduction1').mataUang).toBe(P + 'Currency')
    expect(cari(MEDAN_ATASAN_UANG, 'Deduction2').mataUang).toBe(P + 'Currency')
  })

  it('pyDecimalPlaces terbaca: pxNumber 2 desimal; pxCurrency tak terbaca (pola inti)', () => {
    expect(cari(MEDAN_ADMIN_UANG, 'GrossPremium').sajian).toEqual({ desimal: 2, formatSaatSunting: true })
    expect(cari(MEDAN_ADMIN_UANG, 'RiCommOgp').sajian).toEqual({ desimal: 2, formatSaatSunting: true })
    expect(cari(MEDAN_ADMIN_UANG, 'PremiOgp').sajian).toEqual({ formatSaatSunting: true })
    expect(cari(MEDAN_ADMIN_UANG, 'NetPremium').sajian).toEqual({})
    expect(cari(MEDAN_ATASAN_UANG, 'GrossPremium').sajian).toEqual({ desimal: 2 })
    expect(cari(MEDAN_ATASAN_TOTAL, 'TotalPremium').sajian).toEqual({ desimal: 2 })
  })

  it('atasan Quartal / YearOfQuartal: pyFormatType number, 0 desimal, tanpa pemisah ribuan', () => {
    expect(cari(MEDAN_ATASAN_UMUM, 'Quartal').sajian).toEqual({ desimal: 0, ribuan: false })
    expect(cari(MEDAN_ATASAN_UMUM, 'YearOfQuartal').sajian).toEqual({ desimal: 0, ribuan: false })
    expect(cari(MEDAN_ADMIN_UMUM, 'YearOfQuartal').sajian).toBeUndefined()
  })
})
