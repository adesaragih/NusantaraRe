// Tab Clauses kasus FIRE (tiket 47) - data uji sintetis.

import { existsSync, readFileSync } from 'node:fs'
import { join } from 'node:path'
import { describe, expect, it } from 'vitest'

import type { HasilKlausa, KlausaKasus } from '../api'
import { KLAUSA as K } from '../labels'
import { adaArgumen, gantiArgumen, kodeBahasa, tambahKlausa } from './TabClauses'

const KORPUS = 'D:\\migrasi\\RNM\\NB FacIn\\Section\\'

/** Isi `rowdata` sel ber-`pyCellId` tertentu (rowdata terdalam). */
function blokSel(xml: string, sel: string): string[] {
  const hasil: string[] = []
  const tumpukan: { mulai: number; cocok: boolean }[] = []
  for (const m of xml.matchAll(/<rowdata\b[^>]*>|<\/rowdata>|<pyCellId>(\d+)<\/pyCellId>/g)) {
    const i = m.index ?? 0
    if (m[0].startsWith('<rowdata')) tumpukan.push({ mulai: i, cocok: false })
    else if (m[0] === '</rowdata>') {
      const atas = tumpukan.pop()
      if (atas?.cocok) hasil.push(xml.slice(atas.mulai, i + m[0].length))
    } else if (m[1] === sel) {
      const atas = tumpukan.at(-1)
      if (atas) atas.cocok = true
    }
  }
  return hasil
}
const ada = (xml: string, u: { sel: string; tag: string; label: string }) =>
  blokSel(xml, u.sel).some((b) => b.includes(`<${u.tag}>${u.label}</${u.tag}>`))

describe.skipIf(!existsSync(KORPUS + 'InputClauseFire_FacIn.xml'))('label tab Clauses = korpus', () => {
  it('Choose Clause (InputDtlClause_FacIn sel 209)', () => {
    expect(ada(readFileSync(KORPUS + 'InputDtlClause_FacIn.xml', 'utf-8'), K.pilihKlausa)).toBe(true)
  })
  it.each([K.clauseCode, K.totalArgumen, K.title, K.description, K.lihatArgumen].map((u) => [u.label, u] as const))(
    '%s (InputClauseFire_FacIn)',
    (_, u) => {
      expect(ada(readFileSync(KORPUS + 'InputClauseFire_FacIn.xml', 'utf-8'), u)).toBe(true)
    },
  )
  it.each([K.language, K.keyword, K.cari, K.submit].map((u) => [u.label, u] as const))('%s (ChooseClauseFire)', (_, u) => {
    expect(ada(readFileSync(KORPUS + 'ChooseClauseFire.xml', 'utf-8'), u)).toBe(true)
  })
  it('kolom grid popup dan grid argumen', () => {
    const pilih = readFileSync(KORPUS + 'ChooseClauseFire.xml', 'utf-8')
    for (const k of K.kolomPilih) expect(pilih).toContain(`<pyValue>${k}</pyValue>`)
    const arg = readFileSync(KORPUS + 'InputClauseFire_ViewDtl.xml', 'utf-8')
    for (const k of K.kolomArgumen) expect(arg).toContain(`<pyValue>${k}</pyValue>`)
  })
})

const hasil = (id: string, language = 'Indonesia'): HasilKlausa => ({
  id,
  title: `JUDUL ${id}`,
  info: `INFO ${id}`,
  text: `ISI ${id} _&1`,
  language,
  argumentCount: '1',
})

describe('TabClauses - aturan Pega', () => {
  it('Submit menambah klausa terpilih yang belum ada; ganda dilewati (SearchClauseFireSQL_PostAct)', () => {
    const awal: KlausaKasus[] = tambahKlausa([], [hasil('UJI1')])
    const akhir = tambahKlausa(awal, [hasil('UJI1'), hasil('UJI2', 'Inggris'), hasil('UJI2')])
    expect(akhir.map((k) => k.clauseCode)).toEqual(['UJI1', 'UJI2'])
    expect(akhir[1]).toMatchObject({
      clauseTitle: 'JUDUL UJI2',
      clauseDescription: 'INFO UJI2',
      clauseLanguage: '1',
      clauseLanguageId: 'Inggris',
      clauseContent: 'ISI UJI2 _&1',
      clauseContentTemp: 'ISI UJI2 _&1',
      argumentList: [],
    })
  })
  it('kode bahasa: Indonesia 0, Inggris 1, lainnya 2', () => {
    expect([kodeBahasa('Indonesia'), kodeBahasa('Inggris'), kodeBahasa('Lain')]).toEqual(['0', '1', '2'])
  })
  it('argumen mengganti _&<nomor> di isi asli (ReplaceClauseArgumentFireAct), urutan Pega', () => {
    const arg = [
      { argumentNumber: '1', argumentDescription: 'A', argumentValue: 'SATU' },
      { argumentNumber: '2', argumentDescription: 'B', argumentValue: 'DUA' },
    ]
    expect(gantiArgumen('Nilai _&1 dan _&2, lagi _&1', arg)).toBe('Nilai SATU dan DUA, lagi SATU')
    expect(gantiArgumen('tanpa argumen', [])).toBe('tanpa argumen')
  })
  it('Total Argument tampil bila > 0', () => {
    expect([adaArgumen('0'), adaArgumen(''), adaArgumen('3'), adaArgumen('10')]).toEqual([false, false, true, true])
  })
})

describe('Choose Clause - klausa yang sudah ada tercentang (SearchClauseFireSQL_PreAct langkah 6)', () => {
  const SUMBER = readFileSync(join(__dirname, 'TabClauses.tsx'), 'utf8').replace(/\r\n/g, '\n')
  it('tercentang dan dikunci; Submit hanya menambah yang baru', () => {
    expect(SUMBER).toContain('checked={pilih.has(h.id) || sudahAda.has(h.id)}')
    expect(SUMBER).toContain('disabled={sudahAda.has(h.id)}')
    expect(SUMBER).toContain('sudahAda={new Set(daftar.map((k) => k.clauseCode))}')
  })
})

describe('Isi Klasula tanpa argumen (K47-6, keterangan work owner 05-10-2026)', () => {
  const SUMBER = readFileSync(join(__dirname, 'TabClauses.tsx'), 'utf8').replace(/\r\n/g, '\n')
  it('tanpa grid dan tanpa catatan; tombol Cancel (bawaan Modal) + Submit; isi tidak berubah', () => {
    expect(SUMBER).not.toContain('tanpaArgumen')
    expect(SUMBER).toContain('{K.submit.label}')
    expect(SUMBER).toContain('argumen !== null && argumen.length > 0 && (')
    expect(gantiArgumen('ISI _&1 TETAP', [])).toBe('ISI _&1 TETAP')
  })
})

describe('Language Choose Clause - ganti bahasa langsung cari ulang (05-10-2026)', () => {
  const SUMBER = readFileSync(join(__dirname, 'TabClauses.tsx'), 'utf8').replace(/\r\n/g, '\n')
  it('ganti Language mengosongkan centang dan mencari ulang dengan bahasa baru; kartu menampilkan bahasa klausa', () => {
    expect(SUMBER).toMatch(/setBahasa\(v\)\s*setPilih\(new Map\(\)\)\s*setHalaman\(1\)\s*setSaring\(\{ bahasa: v, kata \}\)/)
    expect(SUMBER).toContain("{k.clauseLanguageId !== '' && <span className=\"nbf-akum__jumlah\">{k.clauseLanguageId}</span>}")
  })
})
