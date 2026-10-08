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

import type { BarisBersarang, SisiPenyesuaian } from './api'
import { KUNCI_BERGOLONGAN, segmenAkhir } from './ekspor/golongan'
import type { AksiTombol, Baca, ButirKerangka, GridKerangka, TombolKerangka } from './ekspor/jenis'
import { DIBUANG, GRID_KURS, KERANGKA_INCLUDE, KERANGKA_RINCIAN, KERANGKA_TAB } from './ekspor/kerangka.gen'
import { PENILAI_SYARAT } from './ekspor/syarat'
import { konteksBaris, RenderKerangka, type KonteksKerangka } from './komponen/KerangkaTab'
import { padankan } from './komponen/medan'
import { saringAgen } from './komponen/PilihAgen'
import SisiForm, { selNilai, TAB_DISEMBUNYIKAN, type ModeLayar } from './komponen/SisiPenyesuaian'
import {
  JENIS_DAFTAR,
  KOLOM_DAFTAR,
  LEBAR_DAFTAR,
  MEDAN_KANAN_BARU,
  MEDAN_KANAN_LAMA,
  MEDAN_KIRI_BARU,
  MEDAN_KIRI_LAMA,
  PENYESUAIAN,
  POLIS_MASTER,
  TAB_BARU_NP,
  TAB_BARU_NP_ADJ,
  TAB_BARU_P,
  TAB_LAMA_NP,
  TAB_LAMA_P,
} from './labelsPenyesuaian'
import { idMasterPolis } from './komponen/PanelPolisMaster'
import type { JenisTulis } from './komponen/aksiTombol'
import { actionsTampil, riwayatDari, simpanTampil, sisiKiriman } from './pages/PenyesuaianKontrak'

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

/**
 * Butir Section RINCIAN BARIS — ikatannya RELATIF atas baris, jadi TIDAK
 * ikut penjaga medan/larik akar yang backend baca (`MEDAN_GO`/`LARIK_GO`).
 */
function butirRincian(): ButirKerangka[] {
  const out: ButirKerangka[] = []
  const jalan = (bs: readonly ButirKerangka[]) => {
    for (const b of bs) {
      out.push(b)
      if (b.t === 'blok') jalan(b.anak)
    }
  }
  for (const v of Object.values(KERANGKA_RINCIAN)) jalan(v)
  return out
}
const BUTIR_RINCIAN = butirRincian()

function sisi(medan: Record<string, string>, larik: Record<string, BarisBersarang[]> = {}): SisiPenyesuaian {
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

function renderBaru(mode: ModeLayar, tabAwal?: string, akar = AKAR_NP, tulis?: (j: JenisTulis) => void): string {
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
      tulis={tulis}
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
  it('⛔ setiap grid: larik kolomnya SAMA panjang — juga grid rincian', () => {
    expect(GRID.length).toBeGreaterThan(60)
    const semua = [...GRID, ...BUTIR_RINCIAN.filter((b): b is GridKerangka => b.t === 'grid')]
    for (const g of semua) {
      const n = g.kolom.length
      expect(
        [g.kunci, g.lebar, g.desimal, g.format, g.syaratSel, g.atSel, g.baca, g.tombol, g.tombolKepala, g.pilihan, g.aksiUbah].map((x) => x.length),
        `${g.prop} @${g.at}`,
      ).toEqual(Array(11).fill(n))
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
    const baca = (x: Baca | null | undefined) => {
      if (Array.isArray(x)) x.forEach((s) => syarat.add(s))
    }
    // ⭐ Syarat AKSI (`pyActionConditions`) ikut — tombol dan perilaku `change`.
    const aksi = (xs: readonly AksiTombol[] | null | undefined) => {
      for (const a of xs ?? []) if (a.syarat !== undefined) syarat.add(a.syarat)
    }
    const tombol = (t: TombolKerangka | null) => {
      if (t === null) return
      t.syarat.forEach((s) => syarat.add(s))
      t.nonaktif?.forEach((s) => syarat.add(s))
      aksi(t.aksi)
    }
    for (const b of [...BUTIR, ...BUTIR_RINCIAN]) {
      b.syarat.forEach((s) => syarat.add(s))
      if (b.t === 'grid') {
        b.syaratSel.forEach((s) => s !== null && syarat.add(s))
        b.baca.forEach(baca)
        b.tombol.forEach(tombol)
        b.tombolKepala.forEach(tombol)
        b.aksiUbah.forEach(aksi)
      }
      if (b.t === 'medan') {
        baca(b.baca)
        aksi(b.aksiUbah)
      }
      if (b.t === 'tombol') tombol(b)
    }
    expect([...syarat].filter((s) => PENILAI_SYARAT[s] === undefined)).toEqual([])
  })

  it('⛔ nol butir MATI di kerangka — penjaga mati dibuang dan DICATAT', () => {
    const mati = /^\s*(1\s*==?\s*2|3\s*==?\s*4|never|false)\s*$/i
    for (const b of [...BUTIR, ...BUTIR_RINCIAN]) for (const s of b.syarat) expect(mati.test(s), `${b.t} @${b.at}: ${s}`).toBe(false)
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

  it('⛔ sel `pxNumber` Section RINCIAN yang BELUM bergolongan — daftar persis, tampil apa adanya', () => {
    const tanpa = new Set<string>()
    for (const b of BUTIR_RINCIAN) {
      if (b.t === 'medan' && b.format === 'pxNumber' && !KUNCI_BERGOLONGAN.has(segmenAkhir(b.kunci))) tanpa.add(b.kunci)
      if (b.t === 'grid') b.kunci.forEach((k, i) => b.format[i] === 'pxNumber' && !KUNCI_BERGOLONGAN.has(segmenAkhir(k)) && tanpa.add(k))
    }
    // Belum digolongkan: butuh bukti tampilan (gambar Limits/Share Pega) —
    // sementara tampil APA ADANYA, tidak ditebak format angkanya.
    expect([...tanpa].sort()).toEqual([
      // `Amount2` — rincian `ShareRetro` (Retro, §17), lewat Premium Adjustment (7 Oktober 2026).
      'AdditionalAmount1', 'AdditionalAmount2', 'AdditionalPct', 'AdjRate', 'Amount2', 'BROKERAGE', 'CASHCALL', 'CessionPct',
      'Deduction', 'DeductionPct', 'IncuredClaim', 'LossRatio', 'LowerBand', 'MDPMinPct', 'MDPPct', 'NETPREMIUM',
      'OutstandingClaim', 'PREMIUM', 'PaidClaim', 'Pct', 'Period', 'Periode', 'QSPct', 'Quarter', 'RICOMM', 'ROLPct',
      'ReinstatementAmount1', 'ReinstatementAmount2', 'ReinstatementPct', 'ReinstatementValue', 'ReisuredParticipant',
      'RetentionPct', 'SpreadingTotalPct', 'SpreadingTotalPctXOL', 'Surplus', 'Total', 'UpperBand',
    ])
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
    // Kolom DATA saja — kolom tombol (Delete) kini ikut di ujung grid.
    const data = g?.kolom.filter((_, i) => g.tombol[i] === null)
    expect(data?.slice(5)).toEqual(['100% Limits ( IDR )', 'Deductible ( IDR )', '100% Limits ( USD )', 'Deductible ( USD )'])
    expect(g?.kunci).not.toContain('Currency')
  })

  it('⭐ kolom tombol IKUT bersama aksinya — dulu dibuang sebagai "jalur tulis" (dibalik 7 Oktober 2026)', () => {
    // Add/Delete grid adalah tata letak Pega, dan Delete/Add tidak menulis
    // apa pun ke basis data: barisnya tinggal di layar sampai Save.
    for (const g of GRID) {
      g.format.forEach((f, i) => {
        expect(f === 'pxButton', `${g.prop} kolom ${String(i)}`).toBe(g.tombol[i] !== null)
      })
    }
    const hapus = GRID.flatMap((g) => g.tombol).filter((t) => t?.aksi.some((a) => a.aksi === 'deleteRow'))
    expect(hapus.length).toBeGreaterThan(5)
    expect(DIBUANG.some((d) => d.alasan === 'kolom tombol (jalur tulis)')).toBe(false)
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

  it('⛔ angka yang SAMA tampil IDENTIK di Old dan New (mode View)', () => {
    // Mode View: di mode Edit sel New yang dapat disunting memegang nilai
    // MENTAH — uji `mode sunting` di bawah.
    const lama = renderLama('Limits')
    const baru = renderBaru('1', 'Limits')
    for (const t of ['1.500.000.000,00', '250.000.000,50']) {
      expect(lama).toContain(t)
      expect(baru).toContain(t)
    }
    for (const html of [renderLama('EGNPI'), renderBaru('1', 'EGNPI')]) {
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
    expect(renderBaru('1', 'Share')).toContain('1.000.000')
  })
})

describe('strip tab dan syarat — bercabang, dari ekspor', () => {
  // ⛔ Retro DISEMBUNYIKAN dari strip (keputusan pemilik proses 7 Oktober 2026).
  const tanpaRetro = (xs: readonly string[]) => xs.filter((t) => !TAB_DISEMBUNYIKAN.has(t))

  it('Old NP dan New NP punya strip BERBEDA', () => {
    expect(judulTab(renderLama())).toEqual(tanpaRetro(TAB_LAMA_NP))
    expect(judulTab(renderBaru('0'))).toEqual(tanpaRetro(TAB_BARU_NP))
  })

  it('⛔ Value Difference TIDAK dirender bila EDMState = 3 atau material bukan 1', () => {
    expect(judulTab(renderBaru('0', undefined, sisi({ ...AKAR_NP.medan, EDMState: '3' })))).not.toContain('Value Difference')
    expect(judulTab(renderBaru('0', undefined, sisi({ ...AKAR_NP.medan, EDMMaterialType: '2' })))).not.toContain('Value Difference')
  })

  it('⛔ Retro DISEMBUNYIKAN di ketiga cabang — syarat ekspornya tetap dibangkitkan', () => {
    expect(judulTab(renderBaru('0'))).not.toContain('Retro')
    const p = sisi({ ...AKAR_NP.medan, ProportionType: 'Proportional', IsMultipleRetro: 'true' })
    expect(judulTab(renderBaru('0', undefined, p))).toEqual(tanpaRetro(TAB_BARU_P))
    // Syarat tampilnya tetap di kerangka: NP tanpa syarat tingkat tab, P bila IsMultipleRetro.
    expect(KERANGKA_TAB['TreatyInTabsNonProportional#Retro']?.syarat).toEqual([])
    expect(KERANGKA_TAB['TreatyInTabsProportional#Retro']?.syarat.join(' ')).toContain('IsMultipleRetro')
  })

  it('⭐ cabang Adjust Premium (EDMState 3): strip Actual GNPI … Information & Submit, tanpa Actual Retro', () => {
    const akar = sisi({ ...AKAR_NP.medan, EDMState: '3' }, { 'ActualValue.EGNPI': [{ TreatyGroup: 'FIRE', Amount: '500' }] })
    const html = renderToStaticMarkup(
      <SisiForm
        judul={PENYESUAIAN.panelBaru}
        sisi={akar}
        akar={akar}
        bacaSaja={false}
        mode="0"
        cabang="NonProportional"
        medanKiri={MEDAN_KIRI_BARU}
        medanKanan={MEDAN_KANAN_BARU}
        bagian="TreatyInTabsNonProportionalAdjustPremi"
        tab={TAB_BARU_NP_ADJ}
      />,
    )
    expect(judulTab(html)).toEqual(['Actual GNPI', 'Actual Limits', 'Actual Share', 'Premium Adjustment', 'Information & Submit'])
    // Isi tab pertama dari ekspor: grid Actual GNPI berisi baris ActualValue.EGNPI.
    expect(html).toContain('FIRE')
    // Halaman memilih Section ini bila Non-Prop dan EDMState 3 (@566980).
    const halaman = readFileSync(join(__dirname, 'pages/PenyesuaianKontrak.tsx'), 'utf8')
    expect(halaman).toContain("const adjustPremi = np && (p.baru.medan.EDMState ?? '') === '3'")
  })

  it('⛔ tab Retro yang diminta langsung jatuh ke tab pertama — isinya (§17) tidak dirender', () => {
    const html = renderBaru('0', 'Retro', sisi({ ...AKAR_NP.medan, IsMultipleRetro: 'true' }))
    expect(html).not.toContain(PENYESUAIAN.retroJarang)
  })

  it('⛔ kolom bersyarat sel: Auto Calculate hilang di View, ada di Edit', () => {
    const p = sisi({ ...AKAR_NP.medan, ProportionType: 'Proportional' }, { ReportingPeriodList: [{ Period: '1', AutoCalculate: 'true', InitialDate: '20180101' }] })
    expect(renderBaru('0', 'Reporting Period', p)).toContain('Auto Calculate')
    expect(renderBaru('1', 'Reporting Period', p)).not.toContain('Auto Calculate')
  })

  it('Installment: grid Currency berincian (▸, tertutup semula) — bukan lagi catatan "tidak diekspor"', () => {
    const html = renderBaru('0', 'Installment', sisi(AKAR_NP.medan, { Installment: [{ Currency: 'IDR', InstallmentList: [] }] }))
    expect(html).toContain('aria-expanded="false"')
    expect(html).toContain('▸')
    // ⛔ 8 Oktober 2026 — `rincianHilang` DIKOSONGKAN atas permintaan
    // pemilik proses, jadi `not.toContain` atasnya selalu lulus (tiap teks
    // memuat ''). Yang dijaga kini sumbernya: catatan itu nol isi, sehingga
    // mustahil tampil di mana pun.
    expect(PENYESUAIAN.rincianHilang).toBe('')
    expect(html).not.toContain('% Installment')
  })

  it('Event Limits: medan dari ekspor, angka tanpa desimal dinyatakan → aturan lama', () => {
    const html = renderBaru('1', 'Event Limits')
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

  it('IsProRate bukan true → After saja, tanpa blok Pro Rate — dan TANPA judulnya', () => {
    const html = renderBaru('0', 'Value Difference', sisi({ ...AKAR_NP.medan, IsProRate: 'false' }, AKAR_NP.larik))
    expect(html).not.toContain('Before Pro Rate Calculation')
    // ⛔ RALAT 7 Oktober 2026: blok `TreatyIn.IsProRate != true` (@47293)
    // ber-`pyIncludeHeader = false` — judulnya tidak tampil di Pega. Yang
    // berjudul hanya pasangan Before/After di cabang Pro Rate (`true`).
    expect(html).not.toContain('After Pro Rate Calculation')
    expect(html).toContain('value="-2,5"')
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

  // ⭐ RALAT 7 Oktober 2026 — uji ini dulu menuntut Save/Actions/Submit MATI.
  // Jalur tulisnya kini ADA (`/api/treaty-in/penyesuaian/*`), jadi yang
  // dijaga berbalik: tombol hidup menurut syarat tampil ekspor, dan mati
  // HANYA selama tombol tulis berjalan.
  it('⭐ Save dan Actions HIDUP menurut syarat tampil; Close selalu', () => {
    const i = HALAMAN.indexOf('function DeretTombol')
    const blok = HALAMAN.slice(i, HALAMAN.indexOf('function Detail(', i))
    expect((blok.match(/disabled=\{sibuk\}/g) ?? []).length).toBe(2)
    expect(blok).not.toMatch(/disabled title/)
    expect(blok).toContain('onClick={onTutup}')
    expect(blok).toMatch(/\{simpanTampil\(mode, medan\) && \(/)
    expect(blok).toMatch(/\{aksiTampil && \(/)
  })

  it('⭐ Submit EDM tab Information & Submit: mati tanpa jalur tulis, HIDUP bila form memberinya', () => {
    const akar = sisi({ ...AKAR_NP.medan, EDMEffective: '20210101' }, AKAR_NP.larik)
    expect(renderBaru('0', 'Information & Submit', akar)).toMatch(/<button[^>]*disabled[^>]*>Submit</)
    const hidup = renderBaru('0', 'Information & Submit', akar, () => undefined)
    expect(hidup).toMatch(/<button[^>]*class="btn btn--primary btn--sm"[^>]*>Submit</)
    expect(hidup).not.toMatch(/<button[^>]*disabled[^>]*>Submit</)
    expect(hidup).not.toMatch(/<button[^>]*disabled[^>]*>Decline offer</)
    // `pyDisabledWhen` `TreatyIn.EDMEffective = ''` tetap berlaku.
    const tanpaEfektif = sisi({ ...AKAR_NP.medan, EDMEffective: '' }, AKAR_NP.larik)
    expect(renderBaru('0', 'Information & Submit', tanpaEfektif, () => undefined)).toMatch(/<button[^>]*disabled[^>]*>Submit</)
  })

  it('Actions hanya bagi pemegang posisi penyetuju, status Accept/Reject', () => {
    const m = { Position: 'ReasTreatyInSecHead', StatusAkseptasi: 'Accept' }
    expect(actionsTampil(m, ['ReasTreatyInSecHead'])).toBe(true)
    expect(actionsTampil({ ...m, StatusAkseptasi: 'Reject' }, ['ReasTreatyInSecHead'])).toBe(true)
    expect(actionsTampil({ ...m, StatusAkseptasi: '' }, ['ReasTreatyInSecHead'])).toBe(false)
    expect(actionsTampil(m, ['ReasTreatyInAdmin'])).toBe(false)
    expect(actionsTampil({ ...m, Position: 'ReasTreatyInAdmin' }, ['ReasTreatyInAdmin'])).toBe(false)
  })

  it('⛔ grid Rate of Exchange dikirim HANYA bila berubah', () => {
    const asal = sisi({ ID: 'X' }, { CurrencyList: [{ Currency: 'USD', Conversion: '1' }] })
    expect(sisiKiriman(asal, asal).larik).not.toHaveProperty('CurrencyList')
    const ubah = sisi({ ID: 'X' }, { CurrencyList: [{ Currency: 'USD', Conversion: '2' }] })
    expect(sisiKiriman(ubah, asal).larik.CurrencyList).toEqual([{ Currency: 'USD', Conversion: '2' }])
  })

  it('⛔ History dari CommentList dokumen, bukan T_VIEW_COMMENT', () => {
    const r = riwayatDari(sisi({}, { CommentList: [{ Date: '20201008T024457.341 GMT', OperatorName: 'OP', IsApproved: 'Accept', Suggest: 'ok' }] }))
    expect(r).toEqual([{ tanggal: '20201008T024457.341 GMT', operator: 'OP', disetujui: 'Accept', catatan: 'ok' }])
  })

  it('⛔ nol impor dari modul tetangga', () => {
    for (const s of [HALAMAN, KOMPONEN]) expect(s).not.toMatch(/from '.*modul\/treatyin\//)
  })
})

describe('⭐ mode sunting panel New — 7 Oktober 2026', () => {
  const dengan = (medan: Record<string, string>, larik: Record<string, Record<string, string>[]> = {}) =>
    sisi({ ...AKAR_NP.medan, ...medan }, { ...AKAR_NP.larik, ...larik })

  // ⛔ RALAT 7 Oktober 2026 — UJI INI DULU MENEGASKAN HAL YANG KINI SALAH.
  //
  // Bentuk pertamanya menuntut nilai MENTAH di kotak isian, mengikuti
  // catatan `KerangkaTab.tsx`: *"memformat di tiap ketukan membuat koma
  // desimal mustahil diketik."* Pengamatannya benar, kesimpulannya terlalu
  // jauh — yang merusak pengetikan adalah `formatNumber` (ia membulatkan
  // dan membuang nol ekor), bukan pemformatan itu sendiri.
  //
  // Pemilik proses 7 Oktober 2026: *"saat meng input juga sudah otomatis
  // langsung ada pemisahnya"*. `FieldAngka` memformat bagian BULAT saja dan
  // membiarkan ekor desimal apa adanya selama kotaknya dipegang;
  // `inti/frontend/lib/angkaKetik.test.ts` memaku keempat keadaan
  // tengah-pengetikan yang dahulu hilang.
  it('⭐ Edit: sel angka grid BERPEMISAH RIBUAN di kotak isian', () => {
    const html = renderBaru('0', 'EGNPI')
    expect(html).toContain('43.994.125.000')
    expect(html).not.toContain('value="43994125000"')
  })

  it('Edit: tombol Add di kepala grid dan Delete per baris tampil (sel 67/72)', () => {
    const html = renderBaru('0', 'Maximum Retention', dengan({}, { Retention: [{ TreatyGroup: 'PROPERTY', Currency: 'IDR', Amount: '10' }] }))
    expect(html).toMatch(/<th[^>]*><button[^>]*>Add<\/button><\/th>/)
    expect(html).toContain('aria-label="Delete 1"')
  })

  it('⛔ View: nol tombol Add/Delete, nol kotak isian di grid', () => {
    const html = renderBaru('1', 'Maximum Retention', dengan({}, { Retention: [{ TreatyGroup: 'PROPERTY', Currency: 'IDR', Amount: '10' }] }))
    expect(html).not.toContain('>Add<')
    expect(html).not.toContain('aria-label="Delete 1"')
  })

  it('⛔ Material Type 2 MENGUNCI Add (pyDisabledWhen sel 67) — tampil tetapi mati', () => {
    const html = renderBaru('0', 'Maximum Retention', dengan({ EDMMaterialType: '2' }))
    expect(html).toMatch(/<button[^>]*disabled=""[^>]*>Add<\/button>/)
  })

  it('⛔ IsEditData = 1 menyembunyikan Add Retention — syarat tampilnya `IsEditData !=\'1\'`', () => {
    // Dua Add tanpa kunci: grid kurs (`ViewState !='1'`) dan Retention. Dengan
    // `IsEditData = 1` tinggal satu — milik kurs, yang syaratnya berbeda.
    const cacah = (h: string) => h.split('>Add<').length - 1
    expect(cacah(renderBaru('0', 'Maximum Retention'))).toBe(2)
    expect(cacah(renderBaru('0', 'Maximum Retention', dengan({ IsEditData: '1' })))).toBe(1)
  })

  it('⛔ panel Old TETAP baca-saja walau mode Edit', () => {
    const html = renderLama('EGNPI')
    expect(html).not.toContain('value="43994125000"')
    expect(html).not.toContain('>Add<')
    expect(html).not.toContain('>Delete<')
  })

  it('Co-Ins Scale (Prop) kembali tampil — `pyIsVisibilityOption = ALWAYS` menimpa `1=2`', () => {
    const p = dengan({ ProportionType: 'Proportional' }, { CoInScale: [{ CoInShare: 'A', PctLimit: '50' }] })
    const html = renderBaru('0', 'Co-Ins Scale', p)
    expect(html).toContain('Co-Insurance Share')
    expect(html).toContain('% Treaty Limit')
    // Tombol ikon tanpa `pyLabel`: teks bawaan Pega.
    expect(html).toContain('>Tambah<')
    expect(html).toContain('>Hapus<')
  })
})

describe('⭐ Choose Ceding / Choose Source of Business — 7 Oktober 2026', () => {
  it('tombol tampil di mode Edit panel New (sel 27/29, `ViewState != 1`)', () => {
    const html = renderBaru('0')
    expect(html).toContain('>Choose Ceding<')
    expect(html).toContain('>Choose Source of Business<')
    // Nilai tetap terbaca sebagai teks di samping tombolnya.
    expect(html).toContain('CEDANT A')
  })

  it('⛔ mode View dan panel Old: nol tombol Choose', () => {
    expect(renderBaru('1')).not.toContain('>Choose Ceding<')
    expect(renderLama()).not.toContain('>Choose Ceding<')
  })

  it('saringan `.ClientName Contains` kata kunci yang SUDAH huruf besar', () => {
    const daftar = [
      { id: '1', nama: 'PT ASURANSI ALFA' },
      { id: '2', nama: 'REASURANSI BETA' },
    ]
    expect(saringAgen(daftar, 'ASURANSI').map((a) => a.id)).toEqual(['1', '2'])
    expect(saringAgen(daftar, 'BETA').map((a) => a.id)).toEqual(['2'])
    expect(saringAgen(daftar, '')).toHaveLength(2)
  })
})

describe('⭐ daftar penyesuaian — 7 Oktober 2026', () => {
  const DAFTAR = readFileSync(join(AKAR, 'pages', 'PenyesuaianKontrak.tsx'), 'utf8')
  it('⛔ tombol `Show/Hide filter` tidak dirender — `pyVisible = NEVER` (sel 112)', () => {
    // Label tombolnya tidak lagi ada, dan halaman daftar tidak memakainya.
    expect(DAFTAR).not.toContain('PENYESUAIAN.saring')
    expect(Object.keys(PENYESUAIAN)).not.toContain('saring')
  })
})

describe('⭐ Existing Policy for Master ID — 7 Oktober 2026', () => {
  it('pengenal yang dicari: `@if(EDMState="", ID, OLDID)` (FetchTreatyExistingProduction)', () => {
    expect(idMasterPolis({ EDMState: '2' }, '1000080/R02', '1000080/R01')).toBe('1000080/R01')
    expect(idMasterPolis({}, '1000080', '')).toBe('1000080')
  })

  it('dua kolom dari sel 31/32, dipasang di antara kepala dan panel Old/New', () => {
    expect([...POLIS_MASTER.kolom]).toEqual(['Policy No', 'Pega ID'])
    const DETAIL = readFileSync(join(AKAR, 'pages', 'PenyesuaianKontrak.tsx'), 'utf8')
    const kepala = DETAIL.indexOf('<Kepala p={p} />')
    expect(kepala).toBeGreaterThan(0)
    expect(DETAIL.indexOf('<PanelPolisMaster')).toBeGreaterThan(kepala)
    expect(DETAIL.indexOf('<PanelPolisMaster')).toBeLessThan(DETAIL.indexOf('tria__bandingan'))
  })
})

describe('rincian baris (expandPane) — Section dari FlowAction ekspor', () => {
  it('⭐ setiap grid expandPane menunjuk Section rincian yang ADA; nol rincian hilang', () => {
    const grid = [...GRID, ...BUTIR_RINCIAN.filter((b): b is GridKerangka => b.t === 'grid')]
    const rujukan = grid.filter((g) => g.rincian !== undefined).map((g) => g.rincian ?? '')
    expect(new Set(rujukan)).toEqual(new Set(Object.keys(KERANGKA_RINCIAN)))
    expect(grid.filter((g) => g.rincianHilang !== undefined)).toEqual([])
    // Installment New/Old, Limits NP (Layers), Share, Retention, EGNPI, dan CoBList di dalam Layers.
    for (const n of ['Installments', 'Installments_ReadOnly', 'Layers', 'Share', 'MaxRetention', 'DetailEGNPI', 'CoBList']) {
      expect(KERANGKA_RINCIAN[n], n).toBeDefined()
    }
  })

  it('⛔ judul blok ber-pyIncludeHeader=false TIDAK tampil (hiddden, hidden reference, Testing); Exclusions tetap', () => {
    const judul = new Set<string>()
    for (const b of [...BUTIR, ...BUTIR_RINCIAN]) if (b.t === 'blok') judul.add(b.judul)
    for (const j of ['hiddden', 'hidden, reference', 'Testing']) expect(judul.has(j), j).toBe(false)
    for (const j of ['Exclusions', 'Summarry of RNM Share', 'Deduction Details']) expect(judul.has(j), j).toBe(true)
  })

  const BARIS_IDR = {
    Currency: 'IDR',
    AmountTotal: '644674819.59',
    PctTotal: '100',
    InstallmentList: [
      { Installment: '1', DueDate: '20240118', WPC: '60', PaymentDate: '20240318', Currency: 'IDR', InstallmentPct: '25', Amount: '161168704.8975' },
    ],
  }
  const konteks = (ubah: boolean, ganti: (b: BarisBersarang) => void = () => undefined): KonteksKerangka =>
    konteksBaris(
      { sisi: sisi(AKAR_NP.medan), akar: sisi(AKAR_NP.medan), halaman: { ...AKAR_NP.medan, ViewState: ubah ? '0' : '1' }, ubah },
      BARIS_IDR,
      ganti,
    )

  it('Section Installments: grid InstallmentList + % Total + Total, nilai dari BARIS', () => {
    const html = renderToStaticMarkup(<RenderKerangka isi={KERANGKA_RINCIAN.Installments ?? []} k={konteks(false)} />)
    for (const x of ['Installment', 'Due Date', 'WPC (In Days)', 'Payment Date', '% Installment', 'Amount', '% Total', 'Total']) expect(html).toContain(x)
    expect(html).toContain('161.168.704,90')
    expect(html).toContain('644.674.819,59')
  })

  it('mode Edit: % Installment dan Amount dapat disunting; Installment dan Payment Date tidak (Read-only)', () => {
    const html = renderToStaticMarkup(<RenderKerangka isi={KERANGKA_RINCIAN.Installments ?? []} k={konteks(true)} />)
    expect(html).toContain('aria-label="% Installment"')
    expect(html).toContain('aria-label="Amount"')
    expect(html).not.toContain('aria-label="Installment"')
    expect(html).not.toContain('aria-label="Payment Date"')
    // Section Old (`Installments_ReadOnly`) — semua sel Read-only.
    const lama = renderToStaticMarkup(<RenderKerangka isi={KERANGKA_RINCIAN.Installments_ReadOnly ?? []} k={konteks(true)} />)
    expect(lama).not.toContain('aria-label="% Installment"')
  })

  it('konteksBaris: suntingan dan hasil rumus MENGGANTI baris; syarat relatif membaca kunci bertitik', () => {
    const diganti: BarisBersarang[] = []
    const k = konteks(true, (b) => diganti.push(b))
    expect(k.sisi.medan.Currency).toBe('IDR')
    expect(k.sisi.larik.InstallmentList).toHaveLength(1)
    expect(k.halaman['.Currency']).toBe('IDR')
    k.ubahMedan?.('PctTotal', '99')
    k.terapkan?.({ medan: { AmountTotal: '1' }, larik: { InstallmentList: [] }, pesan: [] })
    expect(diganti[0]).toMatchObject({ Currency: 'IDR', PctTotal: '99' })
    expect(diganti[1]).toMatchObject({ Currency: 'IDR', AmountTotal: '1', InstallmentList: [] })
    // Rincian di dalam rincian: kunci bertitik baris luar tidak terbawa.
    const dalam = konteksBaris(k, { Layer: '2' }, () => undefined)
    expect(dalam.halaman['.Currency']).toBeUndefined()
    expect(dalam.halaman['.Layer']).toBe('2')
    expect(dalam.panel).toBe(k.panel)
  })
})
