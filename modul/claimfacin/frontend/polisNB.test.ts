// Tampilan polis NB Fac In hanya-baca - fungsi murni (`polisNB.ts`) dan rupa komponen (`components/PolisNB.tsx`,
// dirender react-dom/server). Data uji sintetis UJI-*.

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { createElement } from 'react'
import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import type { Halaman } from './api'
import PolisNB, { tabTampilNB } from './components/PolisNB'
import { susunIsi } from './components/susun'
import {
  angkaNB,
  coverageNB,
  hitungSpreadingNB,
  itemNB,
  LOKASI_NB,
  premiItemNB,
  ratePerMilNB,
  ringkasanCoverageNB,
  spreadingNB,
  TAB_NB,
  tahunNB,
  tanggalNB,
  totalCoverageLokasiNB,
  totalTSIItemNB,
  umumNB,
} from './polisNB'

const kosong = (): Halaman => ({ nilai: {}, daftar: {} })

describe('tanggal polis NB (TanggalDMY nbfacin)', () => {
  it('stempel Pega GMT digeser ke WIB: 17:00 GMT = hari berikutnya, termasuk lintas bulan / tahun / kabisat', () => {
    expect(tanggalNB('20250101T050000.000 GMT')).toBe('01/01/2025')
    expect(tanggalNB('20241231T170000.000 GMT')).toBe('01/01/2025')
    expect(tanggalNB('20240228T170000.000 GMT')).toBe('29/02/2024')
    expect(tanggalNB('20250228T170000.000 GMT')).toBe('01/03/2025')
    expect(tanggalNB('2025-04-30T18:30:00Z')).toBe('01/05/2025')
  })

  it('tanggal biasa (YYYYMMDD, YYYY-MM-DD halaman, DD-MM-YYYY kabel) tidak digeser; teks lain apa adanya', () => {
    expect(tanggalNB('20261009')).toBe('09/10/2026')
    expect(tanggalNB('2025-01-01')).toBe('01/01/2025')
    expect(tanggalNB('2025-01-01 00:00:00')).toBe('01/01/2025')
    expect(tanggalNB('09-10-2026')).toBe('09/10/2026')
    expect(tanggalNB('')).toBe('')
    expect(tanggalNB('UJI')).toBe('UJI')
    expect(tahunNB('20231231T170000.000 GMT')).toBe('2024')
  })
})

describe('angka polis NB', () => {
  // work owner 11-10-2026 "UNTUK CLAIM FACIN ... TAMPILAN 4 DIGIT BELAKANG KOMA": semua angka tampilan polis, termasuk
  // Loss Ratio yang di nbfacin berbatas 2 / 3 desimal
  it('SELALU 4 desimal pemisah Indonesia; kosong tetap kosong', () => {
    expect(angkaNB('1500000000')).toBe('1.500.000.000,0000')
    expect(angkaNB('12.5')).toBe('12,5000')
    expect(angkaNB('0.123456')).toBe('0,1235')
    expect(angkaNB(undefined)).toBe('')
  })
})

describe('blok General', () => {
  it('Business status kode 1 = New Business; Ceding co name = gabungan CedingCoList ";" (kosong -> CedingCoName)', () => {
    const h = kosong()
    h.nilai['OfferFacIn.QuotationData.StatusBusiness'] = '1'
    h.nilai['OfferFacIn.QuotationData.CedingCoName'] = 'UJI GABUNG'
    h.nilai['OfferFacIn.PolicyData.StartDateTime'] = '20241231T170000.000 GMT'
    h.nilai['OfferFacIn.Following'] = 'UJI-LAMA'
    expect(umumNB(h)).toMatchObject({
      statusBisnis: 'New Business',
      ceding: 'UJI GABUNG',
      mulai: '01/01/2025',
      lama: 'UJI-LAMA',
    })
    h.daftar['OfferFacIn.QuotationData.CedingCoList'] = [{ CedingCoName: 'UJI A' }, { CedingCoName: 'UJI B' }]
    expect(umumNB(h).ceding).toBe('UJI A;UJI B')
    h.nilai['OfferFacIn.QuotationData.StatusBusiness'] = 'Renewal'
    expect(umumNB(h).statusBisnis).toBe('Renewal')
    // Risk Scoring = skor tersimpan + catatan kelas (lencana)
    h.nilai['OfferFacIn.ScoringRisk.FinalScore'] = '80.5'
    h.nilai['OfferFacIn.ScoringRisk.NoteFinalScore'] = 'Above Average'
    expect(umumNB(h)).toMatchObject({ riskScoring: '80.5', catatanSkor: 'Above Average' })
  })
})

/** Satu lokasi, dua item (IDR, USD), coverage + spreading. */
function polisFire(): Halaman {
  const h = kosong()
  h.nilai['OfferFacIn.PercentShare'] = '50'
  h.daftar[LOKASI_NB] = [{ 'Property.ObjectNo': '1', 'Property.IsTopRisk': 'false' }]
  h.daftar[itemNB(1)] = [
    { ItemType: 'UJI BANGUNAN', Currency: 'IDR', TSIObjectItem: '1000', TotalGrossPremi: '9' },
    { ItemType: 'UJI MESIN', Currency: 'IDR', TSIObjectItem: '500' },
    { ItemType: 'UJI STOK', Currency: 'USD', TSIObjectItem: '0' },
  ]
  h.daftar[coverageNB(1, 1)] = [
    { OLDID: 'FLEXAS', Premium: '1.5', TSI: '1000', LimitofLiability: '1000', PremiNusantaraRe: '0.75' },
    { OLDID: 'RSMDCC', Premium: '0.25', TSI: '1000', LimitofLiability: '1000', PremiNusantaraRe: '0.125' },
  ]
  h.daftar[coverageNB(1, 2)] = [{ OLDID: 'FLEXAS', Premium: '1', TSI: '500', LimitofLiability: '500' }]
  h.daftar[spreadingNB(1, 1, 1)] = [
    { TreatyType: 'T1', TreatyName: 'UJI QS', SharePercentage: '60', TSISpreaded: '300', PremiumSpreaded: '0.45' },
    { TreatyType: 'T2', TreatyName: 'UJI SURPLUS', SharePercentage: '40', TSISpreaded: '200', PremiumSpreaded: '0.3' },
  ]
  h.daftar[spreadingNB(1, 1, 2)] = [
    { TreatyType: 'T1', SharePercentage: '60', TSISpreaded: '300', PremiumSpreaded: '0.075' },
  ]
  h.daftar[spreadingNB(1, 2, 1)] = [
    { TreatyType: 'T1', SharePercentage: '50', TSISpreaded: '125', PremiumSpreaded: '0.25' },
  ]
  return h
}

describe('total Object Item / Coverage (rumus nbfacin)', () => {
  it('Total TSI per mata uang: hanya total > 0, urut kemunculan', () => {
    expect(totalTSIItemNB(polisFire().daftar[itemNB(1)]!)).toEqual([{ currency: 'IDR', total: '1500' }])
  })

  it('Total Gross Premium item = Σ Premium coverage; tanpa coverage = nilai tersimpan', () => {
    const h = polisFire()
    expect(premiItemNB(h.daftar[coverageNB(1, 1)]!, h.daftar[itemNB(1)]![0]!)).toBe('1.75')
    expect(premiItemNB([], { TotalGrossPremi: '9' })).toBe('9')
  })

  it('Rate = Premium / TSI × 1000 (20 desimal setengah-ke-atas, 8 desimal); TSI 0 = 0', () => {
    expect(ratePerMilNB('2.75', '1500')).toBe('1.83333333')
    expect(ratePerMilNB('1', '3')).toBe('333.33333333')
    expect(ratePerMilNB('5', '0')).toBe('0')
  })

  it('grid total lokasi + ringkasan item / mata uang', () => {
    const h = polisFire()
    expect(totalCoverageLokasiNB(h, 1)).toEqual([
      { currency: 'IDR', tsi: '1500', premium: '2.75', rate: '1.83333333' },
      { currency: 'USD', tsi: '0', premium: '0', rate: '0' },
    ])
    const r = ringkasanCoverageNB(h)
    expect(r.item.map((x) => [x.itemType, x.tsi, x.premium])).toEqual([
      ['UJI BANGUNAN', '1000', '1.75'],
      ['UJI MESIN', '500', '1'],
      ['UJI STOK', '0', '0'],
    ])
    expect(r.mataUang[0]).toEqual({ currency: 'IDR', tsi: '1500', premium: '2.75', rate: '1.83333333' })
  })
})

describe('total spreading (port SumTSIPremiSpreadedRNM_FIRE_Act nbfacin)', () => {
  it('per lokasi per (treaty, mata uang): TSI hanya coverage pertama, premi semua coverage, LoL = LoL × %Share RNM × share', () => {
    const s = hitungSpreadingNB(polisFire())
    const t1 = s.lokasi[0]!.find((t) => t.treatyType === 'T1')!
    // TSI: coverage 1 item 1 (300) + coverage 1 item 2 (125); coverage 2 item 1 tidak dijumlah
    expect(t1.tsiSpreaded).toBe('425')
    expect(t1.premiumSpreaded).toBe('0.775')
    // LoL: 1000 × 50% × 60% + 500 × 50% × 50% (coverage pertama saja)
    expect(t1.claimEstimation).toBe('425')
    // share beda (60 vs 50) -> 0 ("Pct Share Beda")
    expect(t1.sharePercentage).toBe('0')
    expect(t1.treatyName).toBe('UJI QS')
    const t2 = s.lokasi[0]!.find((t) => t.treatyType === 'T2')!
    expect([t2.sharePercentage, t2.tsiSpreaded, t2.claimEstimation]).toEqual(['40', '200', '200'])
  })

  it('ringkasan per treaty dan per mata uang; tanpa top risk TSI Top Risk = 0', () => {
    const s = hitungSpreadingNB(polisFire())
    expect(s.treaty.map((t) => [t.treatyName, t.tsiSpreaded])).toEqual([
      ['UJI QS', '425'],
      ['UJI SURPLUS', '200'],
    ])
    expect(s.mataUang).toEqual([{ currency: 'IDR', tsiTopRisk: '0', tsi: '625', lol: '625', premium: '1.075' }])
  })

  it('lokasi top risk: TSI Top Risk = TSISpreaded coverage pertama; LoL ringkasan hanya dari lokasi top risk', () => {
    const h = polisFire()
    h.daftar[LOKASI_NB] = [{ 'Property.IsTopRisk': 'true' }]
    const s = hitungSpreadingNB(h)
    expect(s.lokasi[0]!.find((t) => t.treatyType === 'T1')!.claimSpreaded).toBe('425')
    expect(s.mataUang[0]!.tsiTopRisk).toBe('625')
  })

  it('First Loss 0 < x < 100 pada coverage pertama: TSI spreaded = TSI × %Share RNM × share', () => {
    const h = polisFire()
    h.daftar[coverageNB(1, 1)]![0]!['FirstLoss'] = '20'
    const t1 = hitungSpreadingNB(h).lokasi[0]!.find((t) => t.treatyType === 'T1')!
    // 1000 × 50% × 60% = 300 (item 1) + 125 (item 2); Total First Loss = 1000 × 20% × 50% × 60% = 60
    expect(t1.tsiSpreaded).toBe('425')
    expect(t1.claimAmountIDR).toBe('60')
  })
})

describe('rupa PolisNB (salinan Inward Facultative nbfacin)', () => {
  const h = polisFire()
  h.nilai['OfferFacIn.QuotationData.NoOfferSlip'] = 'UJI-OFR-1'
  h.nilai['OfferFacIn.QuotationData.PolicyType'] = '2'
  const html = renderToStaticMarkup(createElement(PolisNB, { h, fire: true }))

  it('blok General dua kolom berlabel verbatim, hanya-baca (tanpa input teks / tombol)', () => {
    for (const l of [
      'General',
      'Reff. number',
      'Business status',
      'Insured name',
      'Class of business',
      'Following Policy Number',
      'Day',
    ]) {
      expect(html).toContain(`>${l}<`)
    }
    expect(html).toContain('UJI-OFR-1')
    expect(html).not.toMatch(/<input(?![^>]*type="(radio|checkbox)")/)
    // satu-satunya tombol = tab detail (tanpa Add / Choose / Save / Upload)
    expect([...html.matchAll(/<button[^>]*>/g)].every((m) => m[0].includes('role="tab"'))).toBe(true)
    // kode Policy Type di luar pilihan tidak hilang
    expect(html).toContain('claimfacin__nb-mentah">2<')
  })

  it('tab detail langsung tampil - pemisah kosong dan centang Show Detail dibuang (work owner 11-10-2026)', () => {
    expect(html).not.toContain('Show Detail')
    expect(html).not.toContain('claimfacin__nb-pemisah')
    expect(html).toContain('role="tablist"')
    expect(TAB_NB).toHaveLength(10)
    // FIRE tanpa akses IT: Loss Record / Correspondence tidak tampil (nbfacin `tabTampil`)
    for (const t of tabTampilNB(true, false)) expect(html).toContain(`>${t.replace('&', '&amp;')}<`)
    expect(html).not.toContain('>Loss Record<')
    expect(html).not.toContain('>Correspondence<')
  })

  it('tab tampil per lini seperti nbfacin (work owner 11-10-2026 "cek tampilan nya di nb facin")', () => {
    expect(tabTampilNB(true, true)).toEqual([...TAB_NB])
    expect(tabTampilNB(true, false)).toEqual(TAB_NB.filter((t) => t !== 'Loss Record' && t !== 'Correspondence'))
    // Clauses / Scoring Risk hanya FIRE
    expect(tabTampilNB(false, false)).toEqual([
      'Object',
      'Coverage',
      'Spreading',
      'Inw Fac Cedant Panels',
      'Premium Deduction / Brokerage Fee',
      'Payment',
    ])
    const sumber = readFileSync(join(__dirname, 'components', 'PolisNB.tsx'), 'utf8')
    expect(sumber).not.toContain('TAB_BERISI_NB')
  })

  it('polis Aneka: tab Object berisi grid lokasi Aneka, bukan BelumTersedia', () => {
    const aneka = {
      nilai: {
        'OfferFacIn.QuotationData.BusinessType': 'Liability',
        'OfferFacIn.QuotationData.BusinessName': 'UJI LIABILITY',
      },
      daftar: {
        'OfferFacIn.LocationList': [{ 'Property.ObjectNo': '1', 'Property.RiskLocation.ASMAddress': 'UJI JALAN' }],
      },
    } as Halaman
    const out = renderToStaticMarkup(createElement(PolisNB, { h: aneka, fire: false }))
    expect(out).toContain('>Top Risk<')
    expect(out).toContain('UJI JALAN')
    expect(out).not.toContain('Isi tab Object')
    // Marine Cargo: grid cargo nbfacin
    const cargo = {
      nilai: { 'OfferFacIn.QuotationData.BusinessType': 'MarineCargo' },
      daftar: { 'OfferFacIn.CargoList': [{ TradingNote: 'UJI DAGANG', 'PolicyData.Ship.NM_SHIP': 'UJI KAPAL' }] },
    } as Halaman
    const outCargo = renderToStaticMarkup(createElement(PolisNB, { h: cargo, fire: false }))
    expect(outCargo).toContain('>SailDate<')
    expect(outCargo).toContain('UJI KAPAL')
    // COB di luar FIRE / Aneka / Marine Cargo = pesan nbfacin
    const mbu = {
      nilai: {
        'OfferFacIn.QuotationData.BusinessType': 'MBUCar',
        'OfferFacIn.QuotationData.BusinessName': 'UJI MOTOR',
      },
      daftar: {},
    } as Halaman
    expect(renderToStaticMarkup(createElement(PolisNB, { h: mbu, fire: false }))).toContain(
      'Tab Object untuk Class of Business UJI MOTOR belum tersedia di sistem baru (saat ini baru FIRE).',
    )
  })

  it('susunIsi meneruskan unsur polisNB tanpa diratakan', () => {
    expect(susunIsi([{ jenis: 'polisNB', kendali: 'fire' }])).toEqual([
      { jenis: 'polisNB', t: { jenis: 'polisNB', kendali: 'fire' } },
    ])
  })
})
