// Tab Scoring Risk kasus FIRE (tiket 54) - data uji sintetis.

import { existsSync, readFileSync } from 'node:fs'
import { join } from 'node:path'
import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import type { KolomScoring, TampilanScoring } from '../api'
import { FAKTOR_SCORING, HASIL_SCORING as H, SCORING as S, TEKS_SCORING as T } from '../labels'
import TabScoring, { nadaKelas, persenLossRatioTerkunci, pilihScoring, subtotalFaktor, terkunciPilihanLain, terkunciSilang } from './TabScoring'

const NB = 'D:\\migrasi\\RNM\\NB FacIn\\'
const SUMBER = readFileSync(join(__dirname, 'TabScoring.tsx'), 'utf8').replace(/\r\n/g, '\n')

/** Teks tampil -> bentuk di XML ekspor Pega (`<` `>` `&` di-entity). */
const xml = (s: string) => s.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;')
const utuh = (korpus: string, s: string) => new RegExp(`<(\\w+)>${xml(s).replace(/[.*+?^${}()|[\]\\]/g, '\\$&')}</\\1>`).test(korpus)

describe.skipIf(!existsSync(NB + 'Section\\ScoringRisk.xml'))('label = korpus', () => {
  const korpus = ['ScoringRisk', 'ScoringRiskForm1', 'ScoringRiskForm2', 'ScoringRiskForm3', 'ScoringRiskForm4']
    .map((f) => readFileSync(`${NB}Section\\${f}.xml`, 'utf-8'))
    .join('')

  it('kepala, kolom, kelas, hasil (ScoringRisk.xml) - teks elemen utuh', () => {
    const teks = [
      S.judul, S.insured, S.topRisk, S.dominantRisk, S.periode, S.lokasi1, S.lokasi2, S.status, S.skor, S.total, S.hasil,
      S.klasifikasi, S.rekomendasi, S.persetujuan, S.dari, S.karena, S.lossRatioInput, ...S.kolom,
      ...S.kelas.flatMap((k) => [k.rentang, k.nama]),
      ...[...H.rekomendasi, ...H.dari, ...H.karena].map((r) => r.label),
    ]
    for (const t of teks) expect(utuh(korpus, t), t).toBe(true)
  })

  it('faktor, bobot, sub-faktor, pilihan (ScoringRiskForm1..4; spasi ganda dan < > dipertahankan)', () => {
    for (const f of FAKTOR_SCORING) {
      for (const t of [f.faktor, f.bobot]) expect(utuh(korpus, t), t).toBe(true)
      for (const s of f.sub) for (const t of [s.label, s.bobot, ...s.opsi].filter((x) => x !== '')) expect(utuh(korpus, t), t).toBe(true)
    }
  })

  it('jumlah pilihan per sub-faktor = jumlah ChechBox kolom Top Risk di form', () => {
    const halaman: Record<string, string> = {
      occupation: 'Occupation', riot: 'RiskLocation.Riot', flood: 'RiskLocation.Flood', lossRatio: 'LossRatio.LossRatio',
    }
    for (const [k, h] of Object.entries(halaman)) {
      const n = new Set(korpus.match(new RegExp(`<pyValue>\\.DataScoringRiskList\\(1\\)\\.${h.replace('.', '\\.')}\\.ChechBox(\\d+)</pyValue>`, 'g')) ?? []).size
      const sub = FAKTOR_SCORING.flatMap((f) => f.sub).find((s) => s.kunci === k)
      expect(n, k).toBe(sub?.opsi.length)
    }
  })

  it('% Loss Ratio terkunci bila ChechBox1/2/3/10 (ScoringRiskForm2 pyDisabledWhen)', () => {
    const f2 = readFileSync(NB + 'Section\\ScoringRiskForm2.xml', 'utf-8')
    for (const n of [1, 2, 3, 10]) expect(f2).toContain(`.DataScoringRiskList(1).LossRatio.LossRatio.ChechBox${n} = "true"`)
  })

  it('kunci silang Growing Trees (S54-2): hanya PFE kolom 1+2, Hydrant / Sprinkler kolom 1', () => {
    const f2 = readFileSync(NB + 'Section\\ScoringRiskForm2.xml', 'utf-8')
    const silang = [...f2.matchAll(/<pyDisabledWhen>([^<]*ClassConstruction\.ChechBox4[^<]*)<\/pyDisabledWhen>/g)]
      .map((m) => m[1] ?? '')
      .filter((d) => !/^\.DataScoringRiskList\(\d\)\.ObjectConditions/.test(d))
      .map((d) => (d.match(/^\.DataScoringRiskList\((\d)\)\.FEA\.(\w+)\./) ?? []).slice(1).join(' '))
    expect(silang.sort()).toEqual(['1 HydrantSystem', '1 PortableFireExtinguisher', '1 SprinklerSystem', '2 PortableFireExtinguisher'])
  })
})

describe('TabScoring - aturan', () => {
  const kolom = (pilihan: Record<string, number>): KolomScoring =>
    ({ indexLocation: '1', occupationNote: '', topRiskLocation: '', pilihan, skor: {}, persenLossRatio: '', total: '' }) as KolomScoring

  it('S54-1 (12-10-2026): sesudah satu dipilih, pilihan lain di sub & kolom yang sama terkunci sampai dilepas', () => {
    expect(terkunciPilihanLain(kolom({ riot: 2 }), 'riot', 3)).toBe(true)
    expect(terkunciPilihanLain(kolom({ riot: 2 }), 'riot', 2)).toBe(false)
    expect(terkunciPilihanLain(kolom({ riot: 2 }), 'flood', 1)).toBe(false)
    expect(terkunciPilihanLain(kolom({ riot: 0 }), 'riot', 1)).toBe(false)
    expect(terkunciPilihanLain(kolom({}), 'riot', 1)).toBe(false)
  })

  it('centang saling-eksklusif: pilih pindah, klik yang terpilih = kosong (S54-1)', () => {
    expect(pilihScoring(kolom({ riot: 2 }), 'riot', 3).pilihan.riot).toBe(3)
    expect(pilihScoring(kolom({ riot: 2 }), 'riot', 2).pilihan.riot).toBe(0)
    expect(pilihScoring(kolom({}), 'flood', 1).pilihan).toEqual({ flood: 1 })
  })

  it('% Loss Ratio terkunci hanya pada Nil 5 / 3 / 1 tahun dan No Information', () => {
    expect([1, 2, 3, 10].map((n) => persenLossRatioTerkunci(kolom({ lossRatio: n })))).toEqual([true, true, true, true])
    expect([0, 4, 9].map((n) => persenLossRatioTerkunci(kolom({ lossRatio: n })))).toEqual([false, false, false])
  })

  it('kunci silang: Class Construction = Growing Trees mengunci pilihan "Not Available (Growing Trees)" seperti Pega', () => {
    const gt = kolom({ classConstruction: 4 })
    expect(terkunciSilang(0, gt, 'portableFireExtinguisher', 3)).toBe(true)
    expect(terkunciSilang(1, gt, 'portableFireExtinguisher', 3)).toBe(true)
    expect(terkunciSilang(0, gt, 'hydrant', 2)).toBe(true)
    expect(terkunciSilang(0, gt, 'sprinkler', 2)).toBe(true)
    expect(terkunciSilang(1, gt, 'hydrant', 2)).toBe(false)
    expect(terkunciSilang(1, gt, 'sprinkler', 2)).toBe(false)
    expect(terkunciSilang(0, gt, 'portableFireExtinguisher', 2)).toBe(false)
    expect(terkunciSilang(0, kolom({ classConstruction: 1 }), 'portableFireExtinguisher', 3)).toBe(false)
  })

  it('hitungan di backend; baca-saja di luar marketing; tanpa float', () => {
    expect(SUMBER).toContain('hitungScoring(caseId, baru)')
    expect(SUMBER).toContain('const h = await simpanScoring(caseId, v)')
    expect(SUMBER).toContain('saatTersimpan?.({ finalScore: h.finalScore, noteFinalScore: h.noteFinalScore })')
    expect(SUMBER).toContain('const kunci = !v.bolehUbah ||')
    expect(SUMBER).toContain('{v.bolehUbah && (')
    expect(SUMBER).not.toMatch(/(?<![A-Za-z])Number\(|parseFloat|parseInt/)
  })
})

describe('TabScoring - tampilan', () => {
  const faktor = (no: string) => FAKTOR_SCORING.find((f) => f.no === no)!
  const kolom = (pilihan: Record<string, number>, skor: Record<string, string>, extra: Partial<KolomScoring> = {}): KolomScoring => ({
    indexLocation: '', occupationNote: '', topRiskLocation: '', pilihan, skor, persenLossRatio: '', total: '', ...extra,
  })
  const tampilan = (bolehUbah: boolean, dominan: boolean): TampilanScoring => ({
    bolehUbah, insuredName: 'UJI-TERTANGGUNG', status: 'New Business', mulai: '01-01-2026', selesai: '01-01-2027',
    kolom: [
      kolom({ riot: 2, flood: 2, classConstruction: 4 }, { riot: '1', flood: '5' }, { occupationNote: 'UJI-OKUPASI-1', topRiskLocation: 'UJI-LOKASI-1', total: '88.5' }),
      kolom({}, {}, dominan ? { occupationNote: 'UJI-OKUPASI-2', indexLocation: '2', total: '70' } : {}),
    ],
    finalScore: dominan ? '79.25' : '88.5', noteFinalScore: dominan ? 'Above Average' : 'Good',
    hasil: {
      acceptableNotApproval: true, acceptableWithApproval: false, notAcceptable: false, opDivHead: false, opDirector: false,
      techDirector: false, presDirector: false, dueOccupation: false, dueObjectConditions: false, dueLossRatio: false, dueOthers: false,
    },
    lokasiDominan: [{ indexLocation: '2', label: 'UJI-LOKASI-2' }], pesan: [],
  })

  it('nada kelas: Good / Above Average hijau, Average kuning, Below Average / Poor merah', () => {
    expect(S.kelas.map((k) => nadaKelas(k.nama))).toEqual(['baik', 'baik', 'sedang', 'buruk', 'buruk'])
    expect(nadaKelas('')).toBe('')
  })

  it('subtotal faktor eksak; kosong bila belum ada centang', () => {
    expect(subtotalFaktor(kolom({ riot: 2, flood: 2 }, { riot: '1', flood: '5', earthquake: '' }), faktor('2'))).toBe('6')
    expect(subtotalFaktor(kolom({ riot: 1 }, { riot: '-0.5' }), faktor('2'))).toBe('-0.5')
    expect(subtotalFaktor(kolom({}, {}), faktor('2'))).toBe('')
    expect(subtotalFaktor(undefined, faktor('2'))).toBe('')
  })

  it('render: ringkasan skor + kelas, kartu per faktor, baris terpilih disorot, hasil aktif disorot', () => {
    const html = renderToStaticMarkup(<TabScoring caseId="NB-UJI" awal={tampilan(true, true)} />)
    expect(html.match(/class="nbf-score__faktor"/g)).toHaveLength(FAKTOR_SCORING.length)
    expect(html).toContain('UJI-TERTANGGUNG')
    expect(html).toContain('data-nada="baik"')
    expect(html).toContain('>Above Average</span>')
    expect(html).toContain('class="nbf-score__baris--pilih"')
    expect(html).toContain('<li class="nbf-score__aktif" data-nada="baik"><span class="nbf-score__tanda">✓</span>Acceptable, not required approval</li>')
    expect(html).toContain(`>${T.simpan}<`)
    expect(html).not.toContain(T.tanpaDominan)
  })

  it('render baca-saja: tanpa Save, semua centang terkunci', () => {
    const html = renderToStaticMarkup(<TabScoring caseId="NB-UJI" awal={tampilan(false, false)} />)
    expect(html).toContain(T.bacaSaja)
    expect(html).not.toContain(`>${T.simpan}<`)
    const centang = html.match(/<input type="checkbox"[^>]*>/g) ?? []
    expect(centang.length).toBeGreaterThan(0)
    expect(centang.every((c) => c.includes('disabled'))).toBe(true)
  })

  it('kunci silang tampil di layar: Class Construction Growing Trees -> PFE pilihan 3 Top Risk terkunci', () => {
    const html = renderToStaticMarkup(<TabScoring caseId="NB-UJI" awal={tampilan(true, true)} />)
    const pfe = html.match(/<label title="[^"]*"><input type="checkbox" aria-label="Top Risk - Portable Fire Extinguisher - Not Available \(Growing Trees\)"[^>]*>/)
    expect(pfe?.[0]).toContain(T.terkunciGrowingTrees)
    expect(pfe?.[0]).toContain('disabled')
  })
})
