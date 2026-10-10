// Definisi medan layar endorsemen = sel section korpus `EDM Treaty In` (dibaca 06-10-2026): label, urutan, syarat
// tampil, action set. Harapan diketik dari XML (`DetailPolicyTreatyInAddendum`, kelima section Prop).

import { describe, expect, it } from 'vitest'

import type { Halaman } from './api'
import {
  AKSI_SPREADING,
  MEDAN_KANAN,
  MEDAN_KIRI,
  deretQ,
  medanRemark,
  medanTampil,
  nonPropBaru,
  tampilTabData,
  tataTab,
  varianTab,
  type Medan,
  type TataTab,
} from './medan'

const hal = (nilai: Record<string, string>): Halaman => ({ nilai, daftar: {} })
const label = (ms: Medan[], h: Halaman) => medanTampil(ms, h).map((m) => m.label)
const semua = (t: TataTab) => [...t.gross, ...[...t.kiri, ...t.kanan].flatMap((k) => k.medan)]

describe('header DetailPolicyTreatyInAddendum', () => {
  it('kolom kiri S3-S5: Proportional - Treaty Group / Source Of Business tampil, Type Tax hanya bila With Tax', () => {
    const h = hal({ 'Quotation.ProportionalType': 'Proportional' })
    expect(label(MEDAN_KIRI, h)).toEqual([
      'Old Master ID',
      'Current Master ID',
      'Treaty Group',
      'Source Of Business',
      'No Polis',
      'No EDM',
      'EDM Type',
      'Overiding Commision',
      'With Tax',
    ])
    expect(label(MEDAN_KIRI, hal({ 'PolicyTreatyIn.FlagPPH': 'true' }))).toContain('Type Tax')
    // pyWorkPage.Quotation.ProportionalType = 'NonProportional': Treaty Group hilang; Source Of Business tanpa syarat
    // seperti NB (keputusan work owner 07-10-2026 "Class Of Business GANTI JADII Source Of Business, SAMAIN DENGAN NB")
    const np = label(MEDAN_KIRI, hal({ 'Quotation.ProportionalType': 'NonProportional' }))
    expect(np).not.toContain('Treaty Group')
    expect(np).toContain('Source Of Business')
  })

  it('Class Of Business dibuang; Source Of Business (.SOBName) SEKALI di header, di kolom kiri', () => {
    for (const ProportionalType of ['Proportional', 'NonProportional']) {
      const h = hal({ 'Quotation.ProportionalType': ProportionalType })
      const header = [...label(MEDAN_KIRI, h), ...label(MEDAN_KANAN, h)]
      expect(header).not.toContain('Class Of Business')
      expect(header.filter((l) => l === 'Source Of Business')).toHaveLength(1)
    }
    expect(MEDAN_KIRI.find((m) => m.label === 'Source Of Business')?.jalur).toBe('PolicyTreatyIn.SOBName')
  })

  it('kolom kiri: jalur sel VERBATIM (Old Master ID = .OldData.NoOffer, No Polis / No EDM pyWorkPage.PolicyTreatyIn)', () => {
    expect(MEDAN_KIRI.slice(0, 2).map((m) => m.jalur)).toEqual([
      'PolicyTreatyIn.OldData.NoOffer',
      'PolicyTreatyIn.NoOffer',
    ])
    expect(MEDAN_KIRI.find((m) => m.label === 'No EDM')?.jalur).toBe('PolicyTreatyIn.EDMNo')
  })

  it('hanya sel Auto yang dapat diisi: Overiding Commision, With Tax, Type Tax (kiri), Marketing Officer (kanan)', () => {
    const isian = [...MEDAN_KIRI, ...MEDAN_KANAN].filter((m) => m.jenis !== 'tampil' && !m.kunci).map((m) => m.label)
    expect(isian).toEqual(['Overiding Commision', 'With Tax', 'Type Tax', 'Marketing Officer'])
    expect(MEDAN_KIRI.find((m) => m.label === 'With Tax')?.aksi).toEqual([{ aksi: 'RemoveTypeTax' }])
    expect(MEDAN_KANAN.find((m) => m.label === 'Marketing Officer')?.aksi).toEqual([{ aksi: 'CheckDataMkt' }])
  })

  it('kolom kanan S7-S9: Proportional = deret Q / U/Y + Currency; NonProp lama = UW Year + deret Layer', () => {
    const prop = hal({
      'Quotation.ProportionalType': 'Proportional',
      'PolicyTreatyIn.QuotationData.ProportionalType': 'Proportional',
    })
    expect(label(MEDAN_KANAN, prop)).toEqual([
      'Statement Date',
      'Statement Period',
      'To',
      'Ceding Company',
      'Treaty Type',
      'Q',
      '/',
      'U/Y',
      'Marketing Officer',
      'Proportional Type',
      'Currency',
    ])
    const npLama = hal({
      'Quotation.ProportionalType': 'NonProportional',
      'PolicyTreatyIn.QuotationData.ProportionalType': 'NonProportional',
      'PolicyTreatyIn.IsNewPolicyNonProp': '0',
    })
    const l = label(MEDAN_KANAN, npLama)
    expect(l).toContain('UW Year')
    expect(l).not.toContain('Q')
    expect(l).toContain('Of')
    // NonProp baru: Currency dan Layer hilang
    const npBaru = label(MEDAN_KANAN, {
      ...npLama,
      nilai: { ...npLama.nilai, 'PolicyTreatyIn.IsNewPolicyNonProp': '1' },
    })
    expect(npBaru).not.toContain('Currency')
    expect(npBaru).not.toContain('Of')
  })

  it('deret Q dan Layer satu baris; label nama properti Layer dibuang', () => {
    const d = deretQ(medanTampil(MEDAN_KANAN, hal({ 'Quotation.ProportionalType': 'Proportional' })))
    expect(d.filter(Array.isArray).map((x) => (x as Medan[]).map((m) => m.label))).toEqual([['Q', '/', 'U/Y']])
    const layer = deretQ(
      medanTampil(
        MEDAN_KANAN,
        hal({
          'Quotation.ProportionalType': 'NonProportional',
          'PolicyTreatyIn.QuotationData.ProportionalType': 'NonProportional',
        }),
      ),
    )
    expect(layer.filter(Array.isArray).map((x) => (x as Medan[]).map((m) => m.label))).toEqual([['', '', 'Of', '']])
  })

  it('Remark: dapat diisi hanya di posisi admin (disabled bila PositionNote != ReasTreatyInAdmin)', () => {
    expect(medanRemark(true).kunci).toBe(false)
    expect(medanRemark(false).kunci).toBe(true)
  })
})

describe('cabang badan: S11 AddPremi vs S17 tab', () => {
  it('IsNewPolicyNonProp 1 = AddPremi; 0 atau kosong = tab; nilai lain = keduanya tidak tampil', () => {
    const x = (v: string) => hal({ 'PolicyTreatyIn.IsNewPolicyNonProp': v })
    expect([nonPropBaru(x('1')), tampilTabData(x('1'))]).toEqual([true, false])
    expect([nonPropBaru(x('0')), tampilTabData(x('0'))]).toEqual([false, true])
    expect([nonPropBaru(x('')), tampilTabData(x(''))]).toEqual([false, true])
    expect([nonPropBaru(x('2')), tampilTabData(x('2'))]).toEqual([false, false])
  })

  it('tab Old Data = PropOldData bila OldData.EDMNo kosong, PropOldData2 bila terisi; New Data menurut posisi', () => {
    expect(varianTab('Old Data', hal({}), true)).toBe('lama')
    expect(varianTab('Old Data', hal({ 'PolicyTreatyIn.OldData.EDMNo': 'UJI/E01' }), true)).toBe('lama2')
    expect(varianTab('New Data', hal({}), true)).toBe('baruAdmin')
    expect(varianTab('New Data', hal({}), false)).toBe('baru')
    expect(varianTab('Value Difference', hal({}), false)).toBe('selisih')
  })
})

describe('kerangka tab (lima section Prop)', () => {
  const h = hal({})
  it('awalan properti: OldData / OldData.TreatyDifference / data baru / TreatyDifference', () => {
    expect(tataTab('lama').gross[0]?.jalur).toBe('PolicyTreatyIn.OldData.GrossPremium')
    expect(tataTab('lama2').gross[0]?.jalur).toBe('PolicyTreatyIn.OldData.TreatyDifference.GrossPremium')
    expect(tataTab('baru').gross[0]?.jalur).toBe('PolicyTreatyIn.GrossPremium')
    expect(tataTab('baruAdmin').gross[0]?.jalur).toBe('PolicyTreatyIn.GrossPremium')
    expect(tataTab('selisih').gross[0]?.jalur).toBe('PolicyTreatyIn.TreatyDifference.GrossPremium')
  })

  it('PropOldData / PropValueDifference: label "% Deduction", Net Premium, empat Balance, PPH / PPN; hanya-baca', () => {
    for (const v of ['lama', 'lama2', 'selisih'] as const) {
      const t = tataTab(v)
      expect(t.kiri.map((k) => k.judul)).toEqual(['OGP'])
      expect(label(t.kiri[0]!.medan, h)).toEqual([
        'Premi Ogp',
        '% Deduction In A (OGP)',
        'Deduction In A (OGP)',
        '% Deduction In B (OGP)',
        'Deduction In B (OGP)',
        'Claim',
        'Salvage',
        'Excess Loss',
        'Net Premium',
        // BalanceDueTo kosong = 0: >= 0
        'Balance Before Tax',
        'Balance Before Withholding Tax (PPH 2.2)',
        'Balance Due To Us',
      ])
      expect(t.kanan.flatMap((k) => label(k.medan, h))).toEqual([
        'Premi Onp',
        '% Deduction In A (ONP)',
        'Deduction In A (ONP)',
        '% Deduction In B (ONP)',
        'Deduction In B (ONP)',
        'Deduction1',
        'Deduction2',
        'PPH 2%',
        'PPN 2.2%',
      ])
      expect(semua(t).every((m) => m.jenis === 'tampil' && m.aksi === undefined)).toBe(true)
      expect(t.totalBerlabel).toEqual([])
      expect(t.installment?.jalur).toBe(t.awalan + 'Installment')
    }
  })

  it('kejanggalan XML ditiru: Balance Before Tax / PPH bersyarat .BalanceDueTo DATA BARU; Due To You/Us awalan sendiri', () => {
    const t = tataTab('lama')
    const h1 = hal({ 'PolicyTreatyIn.OldData.BalanceDueTo': '-5', 'PolicyTreatyIn.BalanceDueTo': '7' })
    expect(label(t.kiri[0]!.medan, h1)).toEqual(
      expect.arrayContaining(['Balance Due To You', 'Balance Before Tax', 'Balance Before Withholding Tax (PPH 2.2)']),
    )
    expect(label(t.kiri[0]!.medan, h1)).not.toContain('Balance Due To Us')
    const h2 = hal({ 'PolicyTreatyIn.OldData.BalanceDueTo': '5', 'PolicyTreatyIn.BalanceDueTo': '-7' })
    expect(label(t.kiri[0]!.medan, h2)).not.toContain('Balance Before Tax')
    expect(label(t.kiri[0]!.medan, h2)).toContain('Balance Due To Us')
  })

  it('PropNewData (atasan): tanpa Outstanding Claim / Balance Before* / PPH / PPN; total berlabel; tanpa Installment', () => {
    const t = tataTab('baru')
    const l = semua(t).map((m) => m.label)
    for (const x of [
      'Outstanding Claim',
      'Balance Before Tax',
      'Balance Before Withholding Tax (PPH 2.2)',
      'PPH 2%',
      'PPN 2.2%',
    ])
      expect(l).not.toContain(x)
    expect(l).toContain('Net Premium')
    expect(t.totalBerlabel.map((m) => m.label)).toEqual([
      'Total %Share',
      'Total Premium',
      'Total %Share Claim',
      'Total Claim',
    ])
    expect(t.installment).toBeNull()
    expect(semua(t).every((m) => m.jenis === 'tampil')).toBe(true)
  })

  it('PropNewData2 (admin): label "(%)", Total Premium Before Claim, Outstanding Claim, tombol Save / Calculate', () => {
    const t = tataTab('baruAdmin')
    expect(t.kiri.map((k) => k.judul)).toEqual(['OGP', undefined])
    expect(t.kanan.map((k) => k.judul)).toEqual(['ONP', undefined, undefined])
    expect(label(t.kiri[1]!.medan, h)).toEqual([
      'Claim',
      'Outstanding Claim',
      'Salvage',
      'Excess Loss',
      'Total Premium Before Claim',
      'Balance Before Tax',
      'Balance Before Withholding Tax (PPH 2.2)',
      'Balance Due To Us',
    ])
    expect(label(t.kiri[0]!.medan, h)[1]).toBe('(%) Deduction In A (OGP)')
    // WO 08-10-2026 ("COBA CEK TOMBOLITU, APAKAH MASIH DIPERLUKAN? KALAU SUDAH TIDAK DIHAPUS AJA!"): tombol "Calculate Value Difference" (PropNewData2 S24)
    // dibuang - setiap sel yang mengubah data baru menghitung ulang tab Value Difference sendiri
    expect('hitungSelisih' in t).toBe(false)
    expect(t.installment?.aksi).toEqual([{ aksi: 'FillPaymentInstallment' }, { aksi: 'EDMTCalculateTreatyDifference' }])
    expect(AKSI_SPREADING).toEqual([{ aksi: 'CountSpreading' }, { aksi: 'EDMTCalculateTreatyDifference' }])
  })

  it('PropNewData2 action set = NB MEDAN_ADMIN_UANG + EDMTCalculateTreatyDifference bila sel me-refresh Value Difference', () => {
    const aksi = Object.fromEntries(
      semua(tataTab('baruAdmin'))
        .filter((m) => m.jenis === 'angka')
        .map((m) => [
          m.jalur.replace('PolicyTreatyIn.', ''),
          (m.aksi ?? []).map((a) => a.aksi + (a.param ? `:${a.param}` : '')),
        ]),
    )
    const S = 'EDMTCalculateTreatyDifference'
    expect(aksi).toEqual({
      GrossPremium: [],
      PremiOgp: ['CountOGPONP', S],
      // XML: RiCommOgp tanpa refresh otherSection PropValueDifference
      RiCommOgp: ['CountResult1:Pct', 'CountOGPONP'],
      ResultOgp1: ['CountResult1:Amount', 'CountOGPONP', S],
      OveriddingCommOgp: ['CountResult2Ogp:Pct', 'CountOGPONP', S],
      ResultOgp2: ['CountResult2Ogp:Amount', 'CountOGPONP', S],
      PremiOnp: ['CountOGPONP', S],
      RiCommOnp: ['CountResult1Onp:Pct', 'CountOGPONP', S],
      ResultOnp1: ['CountResult1Onp:Amount', 'CountOGPONP', S],
      OveriddingCommOnp: ['CountResult2Onp:Pct', 'CountOGPONP', S],
      ResultOnp2: ['CountResult2Onp:Amount', 'CountOGPONP', S],
      Claim: ['CountOGPONP', S],
      OutstandingClaim: [],
      SalvageValue: ['CountOGPONP', S],
      ExcessLoss: ['CountOGPONP', S],
      Deduction1: ['CountOGPONP', S],
      Deduction2: ['CountOGPONP', S],
    })
  })

  it('bagian uang 4 desimal, nol tampil "0" (perintah work owner NB 06-10-2026)', () => {
    for (const v of ['lama', 'baru', 'baruAdmin', 'selisih'] as const)
      for (const m of semua(tataTab(v))) expect(m.sajian).toMatchObject({ desimal: 4, nolPolos: true })
  })
})
