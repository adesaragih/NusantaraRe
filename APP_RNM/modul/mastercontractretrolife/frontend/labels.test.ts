// Penjaga bukti label dan tombol Master Contract Retro Life - paket 9.
//
// ⛔ Komentar bukti yang tidak pernah diperiksa adalah HIASAN. Berkas ini membuka korpus XML pada baris
// yang tiap label sebut dan memastikan teksnya memang ada di sana, lalu memastikan:
//   - setiap kunci label berbukti XML ATAU terdaftar `[tidak ada di korpus]` - tidak keduanya, tidak nol;
//   - setiap tombol yang layar render bertekskan label berbukti `pyLabel` (atau `Yes` popup penyimpangan
//     sadar 3/5, atau ikon tutup harness tanpa teks);
//   - ke-31 tombol unik korpus (PARITAS §0) SEMUANYA dirender - tidak ada aksi XML yang dilewati.
// Korpus READ-ONLY - hanya dibaca. Bila korpus tidak terjangkau, uji korpus DILEWATI, bukan gagal.

import { existsSync, readdirSync, readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import * as LABEL from './labels'

const KORPUS = 'D:\\XML\\RNM_BRD\\Master Contract Retro Life'
const adaKorpus = existsSync(KORPUS)

const GRID = 'Section\\GridRetrocessionLife.xml'
const TAHUN = 'Section\\InputRetrocessionLife.xml'
const FORM = 'Section\\InputDtlRetrocessionLife.xml'
const KONTRAK = 'Section\\InputRetroLimitReinsurers.xml'
const REAS = 'Section\\InputSecurityLifeReinsurers.xml'
const SEC = 'Section\\InputSecurityReinsurerLife.xml'
const BIZ = 'Section\\InputBusinessLifeReinsurers.xml'
const RATE = 'Section\\ViewRate.xml'

type Bukti = readonly [kunci: string, berkas: string, baris: number, tag: string]

/** Kunci `OBJEK.medan` → baris tag korpus yang memuat teksnya. */
const BUKTI: readonly Bukti[] = [
  ['TAHUN_MCRL.judul', GRID, 1082, 'pyValue'],
  ['TAHUN_MCRL.labelSelAdd', TAHUN, 8927, 'pyLabelFieldValue'],
  ['TAHUN_MCRL.add', TAHUN, 9007, 'pyLabel'],
  ['TAHUN_MCRL.tooltipAdd', TAHUN, 9005, 'pyTooltip'],
  ['TAHUN_MCRL.kolomId', TAHUN, 9588, 'pyValue'],
  ['TAHUN_MCRL.kolomUnderwritingYear', TAHUN, 9750, 'pyValue'],
  ['TAHUN_MCRL.kolomTransactionYear', TAHUN, 9922, 'pyValue'],
  ['TAHUN_MCRL.kolomStartDate', TAHUN, 10094, 'pyValue'],
  ['TAHUN_MCRL.kolomEndDate', TAHUN, 10266, 'pyValue'],
  ['TAHUN_MCRL.edit', TAHUN, 11774, 'pyLabel'],
  ['TAHUN_MCRL.tooltipEdit', TAHUN, 11772, 'pyTooltip'],
  ['TAHUN_MCRL.reinsType', TAHUN, 12117, 'pyLabel'],
  ['TAHUN_MCRL.inputNewData', FORM, 892, 'pyValue'],
  ['TAHUN_MCRL.formId', FORM, 1834, 'pyLabelFieldValue'],
  ['TAHUN_MCRL.formUnderwritingYear', FORM, 2018, 'pyLabelFieldValue'],
  ['TAHUN_MCRL.formTransactionYear', FORM, 2277, 'pyLabelFieldValue'],
  ['TAHUN_MCRL.formStartDate', FORM, 2522, 'pyLabelFieldValue'],
  ['TAHUN_MCRL.formEndDate', FORM, 2812, 'pyLabelFieldValue'],
  ['TAHUN_MCRL.formModifiedDate', FORM, 3633, 'pyLabelFieldValue'],
  ['TAHUN_MCRL.formInputor', FORM, 3819, 'pyLabelFieldValue'],
  ['TAHUN_MCRL.save', FORM, 5059, 'pyLabel'],
  ['TAHUN_MCRL.cancel', FORM, 5336, 'pyLabel'],
  ['KONTRAK_MCRL.judul', KONTRAK, 666, 'pyValue'],
  ['KONTRAK_MCRL.idTreatyYear', KONTRAK, 1093, 'pyLabelFieldValue'],
  ['KONTRAK_MCRL.formInputor', KONTRAK, 2108, 'pyLabelFieldValue'],
  ['KONTRAK_MCRL.formModifiedDate', KONTRAK, 2537, 'pyLabelFieldValue'],
  ['KONTRAK_MCRL.formReinsType', KONTRAK, 3278, 'pyLabelFieldValue'],
  ['KONTRAK_MCRL.formTreatyStart', KONTRAK, 3642, 'pyLabelFieldValue'],
  ['KONTRAK_MCRL.formTreatyEnd', KONTRAK, 3932, 'pyLabelFieldValue'],
  ['KONTRAK_MCRL.formMinIdr', KONTRAK, 4222, 'pyLabelFieldValue'],
  ['KONTRAK_MCRL.formMaxIdr', KONTRAK, 4486, 'pyLabelFieldValue'],
  ['KONTRAK_MCRL.formMinUsd', KONTRAK, 4709, 'pyLabelFieldValue'],
  ['KONTRAK_MCRL.formMaxUsd', KONTRAK, 4974, 'pyLabelFieldValue'],
  ['KONTRAK_MCRL.save', KONTRAK, 5889, 'pyLabel'],
  ['KONTRAK_MCRL.cancel', KONTRAK, 6176, 'pyLabel'],
  ['KONTRAK_MCRL.add', KONTRAK, 8711, 'pyLabel'],
  ['KONTRAK_MCRL.kolomId', KONTRAK, 9370, 'pyValue'],
  ['KONTRAK_MCRL.kolomReinsType', KONTRAK, 9519, 'pyValue'],
  ['KONTRAK_MCRL.kolomTreatyStart', KONTRAK, 9681, 'pyValue'],
  ['KONTRAK_MCRL.kolomTreatyEnd', KONTRAK, 9843, 'pyValue'],
  ['KONTRAK_MCRL.kolomMinIdr', KONTRAK, 10009, 'pyValue'],
  ['KONTRAK_MCRL.kolomMaxIdr', KONTRAK, 10175, 'pyValue'],
  ['KONTRAK_MCRL.kolomMinUsd', KONTRAK, 10341, 'pyValue'],
  ['KONTRAK_MCRL.kolomMaxUsd', KONTRAK, 10507, 'pyValue'],
  ['KONTRAK_MCRL.edit', KONTRAK, 12738, 'pyLabel'],
  ['KONTRAK_MCRL.businessList', KONTRAK, 13113, 'pyLabel'],
  ['KONTRAK_MCRL.reinsurerList', KONTRAK, 14148, 'pyLabel'],
  ['KONTRAK_MCRL.delete', KONTRAK, 15275, 'pyLabel'],
  ['REINSURER_MCRL.judul', REAS, 637, 'pyValue'],
  ['REINSURER_MCRL.idTreatyYear', REAS, 1066, 'pyLabelFieldValue'],
  ['REINSURER_MCRL.idReinsType', REAS, 1280, 'pyLabelFieldValue'],
  ['REINSURER_MCRL.reinsType', REAS, 1483, 'pyLabelFieldValue'],
  ['REINSURER_MCRL.formInputor', REAS, 2539, 'pyLabelFieldValue'],
  ['REINSURER_MCRL.formReinsurerName', REAS, 3782, 'pyLabelFieldValue'],
  ['REINSURER_MCRL.formShare', REAS, 4821, 'pyLabelFieldValue'],
  ['REINSURER_MCRL.formDiscount', REAS, 5111, 'pyLabelFieldValue'],
  ['REINSURER_MCRL.formOvrComm', REAS, 5403, 'pyLabelFieldValue'],
  ['REINSURER_MCRL.save', REAS, 5770, 'pyLabel'],
  ['REINSURER_MCRL.cancel', REAS, 6050, 'pyLabel'],
  ['REINSURER_MCRL.add', REAS, 8797, 'pyLabel'],
  ['REINSURER_MCRL.kolomId', REAS, 9481, 'pyValue'],
  ['REINSURER_MCRL.kolomReinsurerName', REAS, 9625, 'pyValue'],
  ['REINSURER_MCRL.kolomShare', REAS, 9781, 'pyValue'],
  ['REINSURER_MCRL.kolomDiscount', REAS, 9939, 'pyValue'],
  ['REINSURER_MCRL.kolomOvrComm', REAS, 10097, 'pyValue'],
  ['REINSURER_MCRL.kolomInputor', REAS, 10253, 'pyValue'],
  ['REINSURER_MCRL.kolomUpdateDate', REAS, 10408, 'pyValue'],
  ['REINSURER_MCRL.edit', REAS, 12216, 'pyLabel'],
  ['REINSURER_MCRL.securityReinsurer', REAS, 12555, 'pyLabel'],
  ['REINSURER_MCRL.delete', REAS, 13532, 'pyLabel'],
  ['REINSURER_MCRL.totalShare', REAS, 14404, 'pyValue'],
  ['SECURITY_MCRL.judul', SEC, 659, 'pyValue'],
  ['SECURITY_MCRL.idReinsurer', SEC, 1088, 'pyLabelFieldValue'],
  ['SECURITY_MCRL.reinsurerName', SEC, 1306, 'pyLabelFieldValue'],
  ['SECURITY_MCRL.pctShare', SEC, 1522, 'pyLabelFieldValue'],
  ['SECURITY_MCRL.formInputor', SEC, 2603, 'pyLabelFieldValue'],
  ['SECURITY_MCRL.formSecurityReinsurerName', SEC, 3856, 'pyLabelFieldValue'],
  ['SECURITY_MCRL.formShare', SEC, 4898, 'pyLabelFieldValue'],
  ['SECURITY_MCRL.save', SEC, 5266, 'pyLabel'],
  ['SECURITY_MCRL.cancel', SEC, 5537, 'pyLabel'],
  ['SECURITY_MCRL.add', SEC, 8287, 'pyLabel'],
  ['SECURITY_MCRL.kolomId', SEC, 8978, 'pyValue'],
  ['SECURITY_MCRL.kolomReinsurerName', SEC, 9122, 'pyValue'],
  ['SECURITY_MCRL.kolomShare', SEC, 9278, 'pyValue'],
  ['SECURITY_MCRL.kolomInputor', SEC, 9434, 'pyValue'],
  ['SECURITY_MCRL.kolomUpdateDate', SEC, 9589, 'pyValue'],
  ['SECURITY_MCRL.edit', SEC, 10893, 'pyLabel'],
  ['SECURITY_MCRL.delete', SEC, 11214, 'pyLabel'],
  ['BUSINESS_MCRL.judul', BIZ, 684, 'pyValue'],
  ['BUSINESS_MCRL.idTreatyYear', BIZ, 1103, 'pyLabelFieldValue'],
  ['BUSINESS_MCRL.idReinsType', BIZ, 1320, 'pyLabelFieldValue'],
  ['BUSINESS_MCRL.reinsType', BIZ, 1535, 'pyLabelFieldValue'],
  ['BUSINESS_MCRL.formInputor', BIZ, 2573, 'pyLabelFieldValue'],
  ['BUSINESS_MCRL.formModifiedDate', BIZ, 2982, 'pyLabelFieldValue'],
  ['BUSINESS_MCRL.formBusinessCode', BIZ, 3712, 'pyLabelFieldValue'],
  ['BUSINESS_MCRL.formBusinessName', BIZ, 3898, 'pyLabelFieldValue'],
  ['BUSINESS_MCRL.formRiRate', BIZ, 4351, 'pyLabelFieldValue'],
  ['BUSINESS_MCRL.viewRateForm', BIZ, 4849, 'pyLabel'],
  ['BUSINESS_MCRL.save', BIZ, 5846, 'pyLabel'],
  ['BUSINESS_MCRL.cancel', BIZ, 6064, 'pyLabel'],
  ['BUSINESS_MCRL.add', BIZ, 8599, 'pyLabel'],
  ['BUSINESS_MCRL.kolomBusinessName', BIZ, 9128, 'pyValue'],
  ['BUSINESS_MCRL.kolomRiRate', BIZ, 9286, 'pyValue'],
  ['BUSINESS_MCRL.kolomInputor', BIZ, 9440, 'pyValue'],
  ['BUSINESS_MCRL.kolomUpdateDate', BIZ, 9595, 'pyValue'],
  ['BUSINESS_MCRL.edit', BIZ, 10938, 'pyLabel'],
  ['BUSINESS_MCRL.delete', BIZ, 11264, 'pyLabel'],
  ['BUSINESS_MCRL.viewRate', BIZ, 11581, 'pyLabel'],
  ['BUSINESS_MCRL.copyToAll', BIZ, 12135, 'pyLabel'],
  ['RATE_MCRL.judul', RATE, 861, 'pyValue'],
  ['RATE_MCRL.kolomId', RATE, 1178, 'pyValue'],
  ['RATE_MCRL.kolomUsedBy', RATE, 1324, 'pyValue'],
  ['RATE_MCRL.kolomGender', RATE, 1470, 'pyValue'],
  ['RATE_MCRL.kolomContract', RATE, 1616, 'pyValue'],
  ['RATE_MCRL.kolomAge', RATE, 1762, 'pyValue'],
  ['RATE_MCRL.kolomRate', RATE, 1908, 'pyValue'],
  ['RATE_MCRL.submit', 'FlowAction\\ViewRate.xml', 19, 'pySubmitLabel'],
  ['RATE_MCRL.cancel', 'FlowAction\\ViewRate.xml', 20, 'pyCancelLabel'],
]

/** Kunci `[tidak ada di korpus]` - masing-masing beralasan di `labels.ts`. */
const BUKAN_KORPUS: readonly string[] = [
  'MENU_MCRL.kelompok', // nama FOLDER korpus, dibuktikan terpisah di bawah
  'UMUM_MCRL.tutup',
  'UMUM_MCRL.kosong',
  'REINSURER_MCRL.totalBukan100',
  'SECURITY_MCRL.kolomEksposur',
  'HAPUS_MCRL.pertanyaan',
  'HAPUS_MCRL.ikutTerhapus',
  'HAPUS_MCRL.tanpaAnak',
  'HAPUS_MCRL.reinsurer',
  'HAPUS_MCRL.security',
  'HAPUS_MCRL.business',
  'HAPUS_MCRL.memuat',
  'HAPUS_MCRL.ya',
  'SALIN_MCRL.dasar',
  'SALIN_MCRL.akanDitulis',
  'SALIN_MCRL.nol',
  'SALIN_MCRL.memuat',
  'SALIN_MCRL.ya',
]

/** Nilai label untuk kunci `OBJEK.medan`. */
function nilai(kunci: string): string {
  const [objek, medan] = kunci.split('.') as [string, string]
  const o = (LABEL as unknown as Record<string, Record<string, string>>)[objek]
  const v = o?.[medan]
  if (v === undefined) throw new Error(`label ${kunci} tidak ada`)
  return v
}

/** Setiap kunci label yang diekspor `labels.ts`. */
function semuaKunci(): string[] {
  return Object.entries(LABEL as unknown as Record<string, Record<string, string>>).flatMap(([objek, isi]) =>
    Object.keys(isi).map((m) => `${objek}.${m}`),
  )
}

const xml = (t: string): string => t.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;')

describe.skipIf(!adaKorpus)('label Master Contract Retro Life berbukti barisnya', () => {
  const baris = new Map<string, string[]>()
  const baca = (berkas: string, nomor: number): string => {
    if (!baris.has(berkas)) baris.set(berkas, readFileSync(join(KORPUS, berkas), 'utf8').split('\n'))
    return (baris.get(berkas)?.[nomor - 1] ?? '').trim()
  }

  it.each(BUKTI.map((b) => [...b]))('%s = baris %s:%i <%s>', (kunci, berkas, nomor, tag) => {
    expect(baca(berkas as string, nomor as number)).toBe(`<${tag}>${xml(nilai(kunci as string))}</${tag}>`)
  })

  it('nama menu = nama folder korpus', () => {
    const induk = readdirSync(join(KORPUS, '..'))
    expect(induk).toContain(LABEL.MENU_MCRL.kelompok)
  })

  it('ikon tutup harness = pxIconCancel di keempat harness', () => {
    const ikon: Array<[string, number]> = [
      ['Harness\\InboxRetroLimitReinsurers.xml', 694],
      ['Harness\\InboxRetroLifeReinsurersList.xml', 692],
      ['Harness\\InboxSecurityReinsurerLife.xml', 697],
      ['Harness\\InboxBusinessLifeReinsurers.xml', 691],
    ]
    for (const [berkas, nomor] of ikon) expect(baca(berkas, nomor), berkas).toBe('<pyFormat>pxIconCancel</pyFormat>')
  })

  it('sensus tombol korpus: 33 sel pxButton di 8 section, 31 unik (PARITAS §0)', () => {
    const section = readdirSync(join(KORPUS, 'Section')).filter((f) => f.endsWith('.xml'))
    const cacah = section.reduce(
      (n, f) => n + readFileSync(join(KORPUS, 'Section', f), 'utf8').split('\n').filter((l) => l.trim() === '<pyFormat>pxButton</pyFormat>').length,
      0,
    )
    expect(section).toHaveLength(8)
    expect(cacah).toBe(33)
    // Dua sel ekstra = `Save`/`Cancel` salinan `InputDtlRetrocessionLife` di `InputRetrocessionLife`.
    expect(BUKTI.filter((b) => b[3] === 'pyLabel')).toHaveLength(31)
  })
})

describe('setiap label berbukti XML atau dinyatakan bukan korpus', () => {
  it('dua daftar menutup seluruh kunci, tanpa tumpang tindih', () => {
    const berbukti = new Set(BUKTI.map((b) => b[0]))
    const bukan = new Set(BUKAN_KORPUS)
    for (const k of berbukti) expect(bukan.has(k), k).toBe(false)
    expect(new Set([...berbukti, ...bukan])).toEqual(new Set(semuaKunci()))
    expect(berbukti.size).toBe(BUKTI.length)
  })
})

// ---------------------------------------------------------------------------
// Tombol: yang dirender ⊆ berbukti, dan yang berbukti ⊆ yang dirender.
// ---------------------------------------------------------------------------

/** Isi setiap `<button …>…</button>` di berkas TSX modul. */
function tombolDirender(): { berkas: string; isi: string }[] {
  const akar = __dirname
  const berkas = ['pages', 'components'].flatMap((d) =>
    readdirSync(join(akar, d))
      .filter((f) => f.endsWith('.tsx'))
      .map((f) => join(d, f)),
  )
  return berkas.flatMap((f) =>
    [...readFileSync(join(akar, f), 'utf8').matchAll(/<button\b[\s\S]*?<\/button>/g)].map((m) => ({ berkas: f, isi: m[0] })),
  )
}

/** Kunci label teks tombol - anak `{X_MCRL.k}` tepat sebelum `</button>`; atribut (`title`, `aria-label`) tidak dihitung. */
function kunciTeks(isi: string): string[] {
  const m = /\{([A-Z]+_MCRL)\.(\w+)\}\s*<\/button>$/.exec(isi)
  return m ? [`${m[1]}.${m[2]}`] : []
}

describe('tombol layar = tombol korpus', () => {
  const tombol = tombolDirender()
  // `pyLabel` = teks `pxButton` section; `pySubmitLabel` = tombol dialog FlowAction `ViewRate`.
  const berbuktiTombol = new Set(BUKTI.filter((b) => b[3] === 'pyLabel' || b[3] === 'pySubmitLabel').map((b) => b[0]))
  // Penyimpangan sadar 3 (popup hapus) dan 5 (pratinjau salin): `Yes`; `Cancel` popup = label XML panelnya.
  const tambahanSah = new Set(['HAPUS_MCRL.ya', 'SALIN_MCRL.ya'])

  it('terbaca: ada tombol di layar', () => {
    expect(tombol.length).toBeGreaterThan(25)
  })

  it('setiap tombol yang dirender bertekskan label berbukti pyLabel (atau ikon tutup harness)', () => {
    const langgar: string[] = []
    for (const t of tombol) {
      const kunci = kunciTeks(t.isi)
      if (kunci.length === 0) {
        if (!t.isi.includes('aria-label={UMUM_MCRL.tutup}')) langgar.push(`${t.berkas}: tombol tanpa teks label`)
        continue
      }
      for (const k of kunci) if (!berbuktiTombol.has(k) && !tambahanSah.has(k)) langgar.push(`${t.berkas}: ${k}`)
    }
    expect(langgar).toEqual([])
  })

  it('ke-31 tombol korpus SEMUANYA dirender - tidak ada aksi XML yang dilewati', () => {
    const dirender = new Set(tombol.flatMap((t) => kunciTeks(t.isi)))
    // `Save`/`Cancel` form tahun dan kontrak/reinsurer/... berbeda kunci - masing-masing dihitung.
    const hilang = [...berbuktiTombol].filter((k) => !dirender.has(k))
    expect(hilang).toEqual([])
  })
})
