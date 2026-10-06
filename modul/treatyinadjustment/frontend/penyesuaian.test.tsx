// Uji layar Adjustment — Old Data ‖ New Data, isi tab dari kerangka bangkitan.
//
// Dua jenis uji: (1) kerangka bangkitan dan peta tulisan tangan dijaga
// bentuk serta kelengkapannya, (2) `SisiForm` DIRENDER statis
// (`react-dom/server`) atas dokumen tiruan. Nol Oracle, nol jaringan.
//
// ⚠️ Render statis BUKAN pengganti melihat layar di peramban.

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import type { SisiPenyesuaian } from './api'
import { KUNCI_BERGOLONGAN, segmenAkhir } from './ekspor/golongan'
import type { ButirKerangka, GridKerangka } from './ekspor/jenis'
import { DIBUANG, GRID_KURS, KERANGKA_INCLUDE, KERANGKA_TAB } from './ekspor/kerangka.gen'
import { PENILAI_SYARAT } from './ekspor/syarat'
import { padankan } from './komponen/medan'
import SisiForm, { selNilai, type ModeLayar } from './komponen/SisiPenyesuaian'
import {
  JENIS_DAFTAR,
  KOLOM_DAFTAR,
  LEBAR_DAFTAR,
  MEDAN_KANAN_BARU,
  MEDAN_KANAN_LAMA,
  MEDAN_KIRI_BARU,
  MEDAN_KIRI_LAMA,
  PENYESUAIAN,
  TAB_BARU_NP,
  TAB_BARU_P,
  TAB_LAMA_NP,
  TAB_LAMA_P,
} from './labelsPenyesuaian'
import { riwayatDari, simpanTampil } from './pages/PenyesuaianKontrak'

const AKAR = __dirname
const HALAMAN = readFileSync(join(AKAR, 'pages', 'PenyesuaianKontrak.tsx'), 'utf8')
const KOMPONEN = ['SisiPenyesuaian.tsx', 'KerangkaTab.tsx', 'medan.tsx']
  .map((f) => readFileSync(join(AKAR, 'komponen', f), 'utf8'))
  .join('\n')
const CSS = readFileSync(join(AKAR, 'treatyinadjustment.css'), 'utf8')
const REPO = ['warisan_penyesuaian.go', 'kunci_kerangka_gen.go']
  .map((f) => readFileSync(join(AKAR, '..', 'backend', 'repository', f), 'utf8'))
  .join('\n')

/** Daftar string di dalam `var <nama> = []string{ … }` berkas Go. */
function daftarGo(nama: string): string[] {
  const m = new RegExp(`var ${nama} = \\[\\]string\\{([\\s\\S]*?)\\n?\\}`).exec(REPO)
  if (m === null) throw new Error(`${nama} tidak terbaca dari repository`)
  return [...(m[1] ?? '').matchAll(/"([^"]+)"/g)].map((x) => x[1] ?? '')
}
const MEDAN_GO = [...daftarGo('MedanDibaca'), ...daftarGo('MedanKerangka')]
const LARIK_GO = [...daftarGo('LarikDibaca'), ...daftarGo('LarikKerangka')]

/** Seluruh butir kerangka — tab, include, dan kurs. */
function semuaButir(): ButirKerangka[] {
  const out: ButirKerangka[] = []
  const jalan = (bs: readonly ButirKerangka[]) => {
    for (const b of bs) {
      out.push(b)
      if (b.t === 'blok') jalan(b.anak)
    }
  }
  for (const k of Object.values(KERANGKA_TAB)) jalan(k.isi)
  for (const v of Object.values(KERANGKA_INCLUDE)) jalan(v)
  jalan(Object.values(GRID_KURS))
  return out
}
const BUTIR = semuaButir()
const GRID = BUTIR.filter((b): b is GridKerangka => b.t === 'grid')

function sisi(medan: Record<string, string>, larik: Record<string, Record<string, string>[]> = {}): SisiPenyesuaian {
  return { medan, larik }
}

const AKAR_NP = sisi(
  {
    ID: '1000080/R02', OLDID: '1000080/R01', ProportionType: 'NonProportional', EDMState: '2',
    EDMMaterialType: '1', TreatyContractName: 'KONTRAK BARU', ContractRefNo: 'REF-BARU',
    TeritorialScope: 'INDONESIA', Commencement: '20180101', Termination: '',
    TreatyYear: '2018', AccountingModeNonProp: 'loss', Ceding: 'CEDANT A', LeadingReinsSource: 'SOB A',
    IsProRate: 'true', ProRateDays: '365', ProRateTotalDays: '395', ProRatePercent: '92.4050632911',
    Exclusions: 'PENGECUALIAN BARU', SpecialConditions: 'SYARAT BARU', RSMDLimit: '500000000000', CurrencyRSMD: 'IDR',
    RNMShare: '10', BrokeragePercent: '2.5', FacultativeShare: '0',
    'ValueDifference.RNMShare': '-2.5', 'ValueDifference.BrokeragePercent': '0',
    'ValueBeforeProrate.RNMShare': '-3', StatusAkseptasi: 'Accept', TotalLimitsROL: '87.56',
  },
  {
    CurrencyList: [{ Currency: 'USD', CurrencyID: '10001', Conversion: '14250.5', PeriodStart: '20180101', PeriodEnd: '20181231' }],
    Limits: [{ LayerType: 'Layer', Layer: '1', Limit: '1500000000', Deductible: '250000000.5', Limit2: '100000', Deductible2: '' }],
    EGNPI: [{ TreatyGroup: 'PROPERTY', AsDate: '20191231', Proportion: '54.57908460091369152500', Currency: 'IDR', Amount: '43994125000', AmountIDR: '43994125000.0' }],
    Retention: [],
    Share: [{ LayerType: 'layer', Layer: '1', 'RnmLimitListDisplay(1).Currency': 'IDR', 'RnmLimitListDisplay(1).Value': '1000000', RNMShare: '10' }],
    Installment: [{ Currency: 'IDR' }],
  },
)
const OLD_NP = sisi(
  {
    TreatyContractName: 'KONTRAK LAMA', TeritorialScope: 'INDONESIA', Commencement: '20170101',
    TreatyYear: '2017', AccountingModeNonProp: 'loss', Ceding: 'CEDANT A', LeadingReinsSource: 'SOB A',
    FacultativeShare: '0', RNMShare: '10',
  },
  {
    CurrencyList: [{ Currency: 'USD', Conversion: '14250.5', PeriodStart: '20170101', PeriodEnd: '20171231' }],
    Limits: [{ LayerType: 'Layer', Layer: '1', Currency: 'IDR', Limit: '1500000000', Deductible: '250000000.5' }],
    EGNPI: [{ TreatyGroup: 'PROPERTY', AsDate: '20181231', Proportion: '54.57908460091369152500', Currency: 'IDR', Amount: '43994125000', AmountIDR: '43994125000.0' }],
  },
)

function renderLama(tabAwal?: string, akar = AKAR_NP, lama = OLD_NP): string {
  const cabang = akar.medan.ProportionType ?? ''
  const np = cabang === 'NonProportional'
  return renderToStaticMarkup(
    <SisiForm
      judul={PENYESUAIAN.panelLama}
      sisi={lama}
      akar={akar}
      bacaSaja
      mode="0"
      cabang={cabang}
      medanKiri={MEDAN_KIRI_LAMA}
      medanKanan={MEDAN_KANAN_LAMA}
      bagian={np ? 'TreatyInTabsNonProportionalOldData' : 'TreatyInTabsProportionalOldData'}
      tab={tabAwal !== undefined ? [tabAwal] : np ? TAB_LAMA_NP : TAB_LAMA_P}
    />,
  )
}

function renderBaru(mode: ModeLayar, tabAwal?: string, akar = AKAR_NP): string {
  const cabang = akar.medan.ProportionType ?? ''
  const np = cabang === 'NonProportional'
  return renderToStaticMarkup(
    <SisiForm
      judul={PENYESUAIAN.panelBaru}
      sisi={akar}
      akar={akar}
      bacaSaja={false}
      mode={mode}
      cabang={cabang}
      medanKiri={MEDAN_KIRI_BARU}
      medanKanan={MEDAN_KANAN_BARU}
      bagian={np ? 'TreatyInTabsNonProportional' : 'TreatyInTabsProportional'}
      tab={tabAwal !== undefined ? [tabAwal] : np ? TAB_BARU_NP : TAB_BARU_P}
    />,
  )
}

/** Seluruh kontrol isian di HTML: tag pembukanya. */
function kontrol(html: string): string[] {
  return [...html.matchAll(/<(input|textarea|select)\b[^>]*>/g)].map((m) => m[0])
}
const bebas = (k: string[]) => k.filter((t) => !/\breadOnly=""|\breadonly=""|\bdisabled=""/i.test(t))
const judulTab = (html: string) =>
  [...html.matchAll(/role="tab"[^>]*>([^<]+)</g)].map((m) => (m[1] ?? '').replace(/&amp;/g, '&'))

describe('kerangka bangkitan — bentuk dan kelengkapan', () => {
  it('⛔ setiap grid: ketujuh larik kolomnya SAMA panjang', () => {
    expect(GRID.length).toBeGreaterThan(60)
    for (const g of GRID) {
      const n = g.kolom.length
      expect([g.kunci, g.lebar, g.desimal, g.format, g.syaratSel, g.atSel].map((x) => x.length), `${g.prop} @${g.at}`).toEqual([n, n, n, n, n, n])
    }
    expect([LEBAR_DAFTAR.length, JENIS_DAFTAR.length]).toEqual([KOLOM_DAFTAR.length, KOLOM_DAFTAR.length])
  })

  it('kerangka memuat TEPAT tab tiap Section sisi × cabang', () => {
    const tab = (bagian: string) =>
      Object.keys(KERANGKA_TAB).filter((k) => k.startsWith(`${bagian}#`)).map((k) => k.slice(bagian.length + 1))
    expect(tab('TreatyInTabsNonProportional')).toEqual([...TAB_BARU_NP])
    expect(tab('TreatyInTabsNonProportionalOldData')).toEqual([...TAB_LAMA_NP])
    expect(tab('TreatyInTabsProportional')).toEqual([...TAB_BARU_P])
    expect(tab('TreatyInTabsProportionalOldData')).toEqual([...TAB_LAMA_P])
  })

  it('⛔ setiap teks syarat di kerangka punya penilai — yang tak dikenal menggagalkan', () => {
    const syarat = new Set<string>()
    for (const k of Object.values(KERANGKA_TAB)) k.syarat.forEach((s) => syarat.add(s))
    for (const b of BUTIR) {
      b.syarat.forEach((s) => syarat.add(s))
      if (b.t === 'grid') b.syaratSel.forEach((s) => s !== null && syarat.add(s))
    }
    expect([...syarat].filter((s) => PENILAI_SYARAT[s] === undefined)).toEqual([])
  })

  it('⛔ nol butir MATI di kerangka — penjaga mati dibuang dan DICATAT', () => {
    const mati = /^\s*(1\s*==?\s*2|3\s*==?\s*4|never|false)\s*$/i
    for (const b of BUTIR) for (const s of b.syarat) expect(mati.test(s), `${b.t} @${b.at}: ${s}`).toBe(false)
    expect(DIBUANG.length).toBeGreaterThan(50)
    expect(DIBUANG.every((d) => d.alasan !== '')).toBe(true)
  })

  it('⛔ nol nama orang di kerangka — butir bersyarat identitas operator dibuang', () => {
    const teks = readFileSync(join(AKAR, 'ekspor', 'kerangka.gen.ts'), 'utf8')
    expect(teks).not.toMatch(/pxInsName|pyUserName|pyUserIdentifier/)
  })

  it('⛔ setiap kunci yang kerangka dan kepala baca ADA di daftar-izin repository', () => {
    const medan = new Set<string>()
    const larik = new Set<string>(['CurrencyList', 'CommentList'])
    for (const m of [...MEDAN_KIRI_LAMA, ...MEDAN_KANAN_LAMA, ...MEDAN_KIRI_BARU, ...MEDAN_KANAN_BARU]) medan.add(m.kunci)
    for (const b of BUTIR) {
      if (b.t === 'medan') medan.add(b.kunci)
      if (b.t === 'grid') larik.add(b.larik)
    }
    for (const k of ['ProRateDays', 'ProRateTotalDays', 'ProRatePercent', 'StatusAkseptasi', 'IsEditData', 'EDMState', 'EDMMaterialType']) medan.add(k)
    expect([...medan].filter((k) => !MEDAN_GO.includes(k))).toEqual([])
    expect([...larik].filter((k) => !LARIK_GO.includes(k))).toEqual([])
  })

  it('⛔ setiap sel `pxNumber` punya golongan tertulis — yang lain didaftarkan, bukan ditebak', () => {
    const tanpa = new Set<string>()
    for (const b of BUTIR) {
      if (b.t === 'medan' && b.format === 'pxNumber' && !KUNCI_BERGOLONGAN.has(segmenAkhir(b.kunci))) tanpa.add(b.kunci)
      if (b.t === 'grid') b.kunci.forEach((k, i) => b.format[i] === 'pxNumber' && !KUNCI_BERGOLONGAN.has(segmenAkhir(k)) && tanpa.add(k))
    }
    expect([...tanpa]).toEqual([])
  })
})

describe('§0 — kolom disalin per SEL dari ekspor, sarang dibuang dengan benar', () => {
  const grid = (kunci: string, larik: string) =>
    KERANGKA_TAB[kunci]?.isi.find((b): b is GridKerangka => b.t === 'grid' && b.larik === larik)

  it('⛔ Limits Old = DELAPAN kolom (@589778), kepala Layers hanya di kolom pertama', () => {
    const g = grid('TreatyInTabsNonProportionalOldData#Limits', 'Limits')
    expect(g?.kolom).toEqual(['Layers', '', '', '', '', 'Currency', '100% Limits', 'Deductible'])
    expect(g?.kunci).toEqual(['LayerType', 'Layer', 'pyTemplateRichTextEditor', 'LayerPartType', 'LayerPart', 'Currency', 'Limit', 'Deductible'])
    expect(g?.lebar).toEqual([101, 101, 101, 101, 104, 122, 182, 182])
  })

  it('⛔ Limits New NP berpasangan IDR/USD, tanpa Currency — beda dari Old DIPERTAHANKAN', () => {
    const g = grid('TreatyInTabsNonProportional#Limits', 'Limits')
    expect(g?.kolom.slice(5)).toEqual(['100% Limits ( IDR )', 'Deductible ( IDR )', '100% Limits ( USD )', 'Deductible ( USD )'])
    expect(g?.kunci).not.toContain('Currency')
  })

  it('⛔ kolom tombol (jalur tulis) tidak ikut, dan tercatat dibuang', () => {
    for (const g of GRID) expect(g.format).not.toContain('pxButton')
    expect(DIBUANG.some((d) => d.alasan === 'kolom tombol (jalur tulis)')).toBe(true)
  })
})

describe('§2 — desimal per kolom dari ekspor', () => {
  it('⛔ desimal dinyatakan → nol di ekor DIPERTAHANKAN sampai presisinya', () => {
    expect(selNilai('uang', '1500000000', 2)).toBe('1.500.000.000,00')
    expect(selNilai('uang', '1', 2)).toBe('1,00')
    expect(selNilai('persen', '54.57908460091369152500', 2)).toBe('54,58')
    expect(selNilai('persen', '100.00', 4)).toBe('100,0000')
    expect(selNilai('persen', '92.4050632911', 4)).toBe('92,4051')
  })

  it('desimal TIDAK dinyatakan → aturan lama, nol di ekor dibuang', () => {
    expect(selNilai('uang', '43994125000.0')).toBe('43.994.125.000')
    expect(selNilai('uang', '250000000.12345')).toBe('250.000.000,1235')
    expect(selNilai('persen', '87.5600')).toBe('87,56')
  })

  it('⛔ teks bukan angka apa adanya; tanpa % tempelan; tanggal DD-MM-YYYY', () => {
    expect(selNilai('persen', '>=30% up to < 50%', 2)).toBe('>=30% up to < 50%')
    expect(selNilai('uang', 'TBA', 2)).toBe('TBA')
    expect(selNilai('persen', '10', 2)).toBe('10,00')
    expect(selNilai('tanggal', '20180101')).toBe('01-01-2018')
    expect(selNilai('tanggal', 'bukan tanggal')).toBe('bukan tanggal')
    expect(padankan('-2,5', 2)).toBe('-2,50')
    expect(padankan('IDR', 2)).toBe('IDR')
  })

  it('⛔ pemadanan di lapis modul — format.ts DIPANGGIL, tidak ditulis ulang', () => {
    expect(KOMPONEN).toContain("from '../../../../inti/frontend/lib/format'")
    expect(KOMPONEN).not.toMatch(/function formatNumber|function formatPersen|function formatDate|toLocaleString/)
  })

  it('⛔ angka yang SAMA tampil IDENTIK di Old dan New', () => {
    const lama = renderLama('Limits')
    const baru = renderBaru('0', 'Limits')
    for (const t of ['1.500.000.000,00', '250.000.000,50']) {
      expect(lama).toContain(t)
      expect(baru).toContain(t)
    }
    for (const html of [renderLama('EGNPI'), renderBaru('0', 'EGNPI')]) {
      expect(html).toContain('54,58')
      expect(html).toContain('43.994.125.000')
    }
    expect(lama).toContain('14.250,5')
    expect(baru).toContain('14.250,5')
  })
})

describe('§5 — kotak tanggal diisi nilai TERSIMPAN, bukan terjemahan tampil', () => {
  const p = () => sisi({ ...AKAR_NP.medan, ProportionType: 'Proportional', ReportingStart: '20180315' })

  it('⛔ kepala dan tab: `type="date"` berisi YYYY-MM-DD dari nilai tersimpan', () => {
    expect(renderLama()).toMatch(/type="date"[^>]*value="2017-01-01"/)
    expect(renderBaru('0')).toMatch(/type="date"[^>]*value="2018-01-01"/)
    expect(renderBaru('0', 'Reporting Period', p())).toMatch(/type="date"[^>]*value="2018-03-15"/)
  })

  it('⛔ nol bentuk bergaris miring masuk ke kotak tanggal', () => {
    for (const html of [renderLama(), renderBaru('0'), renderBaru('0', 'Reporting Period', p())]) {
      for (const t of kontrol(html).filter((x) => x.includes('type="date"'))) {
        expect(t).toMatch(/value="(\d{4}-\d{2}-\d{2})?"/)
      }
    }
    // Cabang pxDateTime memberi nilai MENTAH, bukan hasil `selNilai`.
    expect(KOMPONEN).toMatch(/case 'pxDateTime':[\s\S]{0,200}<TanggalBacaSaja label=\{label\} nilai=\{v\} \/>/)
  })

  it('⛔ tanggal kosong: value="" + CSS modul, bukan dd/mm/yyyy', () => {
    expect(renderBaru('0')).toMatch(/type="date"[^>]*value=""/)
    expect(CSS).toMatch(/\.treatyinadjustment input\[type="date"\]\[value=""\]:not\(:focus\)/)
  })
})

describe('panel Old Data — baca-saja, dan itu seluruh bedanya', () => {
  it('⛔ NOL medan dapat disunting — di kepala DAN di setiap tab', () => {
    for (const t of TAB_LAMA_NP) {
      const html = renderLama(t)
      expect(bebas(kontrol(html)), t).toEqual([])
      expect(html).not.toContain('<select')
    }
    const akarP = sisi({ ...AKAR_NP.medan, ProportionType: 'Proportional' })
    for (const t of TAB_LAMA_P) expect(bebas(kontrol(renderLama(t, akarP, OLD_NP))), t).toEqual([])
  })

  it('kepala New: Edit membuka medan yang boleh disunting, View menguncinya', () => {
    const akar = sisi({ ...AKAR_NP.medan, EDMMaterialType: '2' }, AKAR_NP.larik)
    expect(bebas(kontrol(renderBaru('1', 'Exclusions', akar))).length).toBeLessThan(bebas(kontrol(renderBaru('0', 'Exclusions', akar))).length)
  })

  it('⛔ Contract Ref No Old mengikat AKAR (@40302); Exclusions Old NP juga (@1777411)', () => {
    expect(renderLama()).toContain('value="REF-BARU"')
    expect(renderLama('Exclusions')).toContain('PENGECUALIAN BARU')
  })

  it('⛔ label Share Old BERBEDA dari New — dari Section masing-masing', () => {
    const lama = renderLama('Share')
    const baru = renderBaru('0', 'Share')
    expect(lama).toContain('Fakultative Brokerage')
    expect(baru).not.toContain('Fakultative Brokerage')
    // `Brokerage From Other Retro` bersyarat `FacultativeShare >0` — 0 di sini.
    expect(baru).not.toContain('Brokerage From Other Retro')
  })

  it('grid RNM Share membaca larik bersarang berindeks', () => {
    expect(renderBaru('0', 'Share')).toContain('1.000.000')
  })
})

describe('strip tab dan syarat — bercabang, dari ekspor', () => {
  it('Old NP dan New NP punya strip BERBEDA', () => {
    expect(judulTab(renderLama())).toEqual([...TAB_LAMA_NP])
    expect(judulTab(renderBaru('0'))).toEqual([...TAB_BARU_NP])
  })

  it('⛔ Value Difference TIDAK dirender bila EDMState = 3 atau material bukan 1', () => {
    expect(judulTab(renderBaru('0', undefined, sisi({ ...AKAR_NP.medan, EDMState: '3' })))).not.toContain('Value Difference')
    expect(judulTab(renderBaru('0', undefined, sisi({ ...AKAR_NP.medan, EDMMaterialType: '2' })))).not.toContain('Value Difference')
  })

  it('⛔ Retro NP tanpa syarat tingkat tab; Retro P hanya bila IsMultipleRetro = "true"', () => {
    expect(judulTab(renderBaru('0'))).toContain('Retro')
    expect(judulTab(renderBaru('0', undefined, sisi({ ...AKAR_NP.medan, ProportionType: 'Proportional' })))).not.toContain('Retro')
    expect(judulTab(renderBaru('0', undefined, sisi({ ...AKAR_NP.medan, ProportionType: 'Proportional', IsMultipleRetro: 'true' })))).toEqual([...TAB_BARU_P])
  })

  it('⛔ isi Retro TIDAK dibangun (§17) — wadah + syaratnya saja', () => {
    expect(renderBaru('0', 'Retro')).not.toContain(PENYESUAIAN.retroJarang)
    const dengan = renderBaru('0', 'Retro', sisi({ ...AKAR_NP.medan, IsMultipleRetro: 'true' }))
    expect(dengan).toContain(PENYESUAIAN.retroJarang)
    expect(dengan).toContain('Has Share To Retro:')
  })

  it('⛔ kolom bersyarat sel: Auto Calculate hilang di View, ada di Edit', () => {
    const p = sisi({ ...AKAR_NP.medan, ProportionType: 'Proportional' }, { ReportingPeriodList: [{ Period: '1', AutoCalculate: 'true', InitialDate: '20180101' }] })
    expect(renderBaru('0', 'Reporting Period', p)).toContain('Auto Calculate')
    expect(renderBaru('1', 'Reporting Period', p)).not.toContain('Auto Calculate')
  })

  it('Installment: grid Currency + catatan rincian baris yang tidak diekspor', () => {
    const html = renderBaru('0', 'Installment')
    expect(html).toContain(PENYESUAIAN.rincianBaris)
    expect(html).toContain('pyGridRowDetails')
  })

  it('Event Limits: medan dari ekspor, angka tanpa desimal dinyatakan → aturan lama', () => {
    const html = renderBaru('0', 'Event Limits')
    expect(html).toContain('RSMD Limit')
    expect(html).toContain('500.000.000.000')
  })
})

describe('Value Difference — dua Section yang TIDAK disatukan', () => {
  it('IsProRate = true → Before (ValueBeforeProrate) DAN After (ValueDifference)', () => {
    const html = renderBaru('0', 'Value Difference')
    expect(html).toContain('Before Pro Rate Calculation')
    expect(html).toContain('After Pro Rate Calculation')
    expect(html).toContain('value="-3"')
    expect(html).toContain('value="-2,5"')
    expect(html).toContain('Pro Rate Percentage')
  })

  it('IsProRate bukan true → After saja, tanpa blok Pro Rate', () => {
    const html = renderBaru('0', 'Value Difference', sisi({ ...AKAR_NP.medan, IsProRate: 'false' }, AKAR_NP.larik))
    expect(html).not.toContain('Before Pro Rate Calculation')
    expect(html).toContain('After Pro Rate Calculation')
    expect(html).not.toContain('Pro Rate Percentage')
  })
})

describe('deret tombol dan History', () => {
  it('Save tampil hanya di Edit dan bila belum Resolve Complete; kunci tak ada = tampil', () => {
    expect(simpanTampil('0', { StatusAkseptasi: 'Accept' })).toBe(true)
    expect(simpanTampil('0', {})).toBe(true)
    expect(simpanTampil('0', { StatusAkseptasi: 'Resolve Complete' })).toBe(false)
    expect(simpanTampil('1', { StatusAkseptasi: 'Accept' })).toBe(false)
  })

  it('⛔ tombol tulis MATI, bukan hilang; Close hidup', () => {
    const i = HALAMAN.indexOf('function DeretTombol')
    const blok = HALAMAN.slice(i, HALAMAN.indexOf('/** Mode detail', i))
    expect((blok.match(/disabled/g) ?? []).length).toBe(2)
    expect(blok).toContain('onClick={onTutup}')
    // Tombol di dalam tab (Update Total, Submit, …) juga mati.
    expect(renderBaru('0', 'Information & Submit')).toMatch(/<button[^>]*disabled[^>]*>Submit</)
  })

  it('⛔ History dari CommentList dokumen, bukan T_VIEW_COMMENT', () => {
    const r = riwayatDari(sisi({}, { CommentList: [{ Date: '20201008T024457.341 GMT', OperatorName: 'OP', IsApproved: 'Accept', Suggest: 'ok' }] }))
    expect(r).toEqual([{ tanggal: '20201008T024457.341 GMT', operator: 'OP', disetujui: 'Accept', catatan: 'ok' }])
  })

  it('⛔ nol impor dari modul tetangga', () => {
    for (const s of [HALAMAN, KOMPONEN]) expect(s).not.toMatch(/from '.*modul\/treatyin\//)
  })
})
