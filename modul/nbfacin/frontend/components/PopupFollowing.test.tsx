// Popup Following Policy Number (tiket 70).

import { existsSync, readFileSync } from 'node:fs'
import { join } from 'node:path'
import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import { FOLLOWING as F } from '../labels'
import PopupFollowing, { tanggalFollowing } from './PopupFollowing'

const NB = 'D:\\migrasi\\RNM\\NB FacIn\\'
const DDL = 'D:\\migrasi\\RNM\\DDL\\FilterForViewRetro.xml'
const baca = (...p: string[]) => readFileSync(join(__dirname, ...p), 'utf8').replace(/\r\n/g, '\n')
const SUMBER = baca('PopupFollowing.tsx')
const HALAMAN = baca('..', 'pages', 'InwardFacultative.tsx')

describe('Following Policy - label korpus', () => {
  it.skipIf(!existsSync(NB + 'Section\\GridViewFollowingNB.xml'))('kolom & link grid = GridViewFollowingNB', () => {
    const xml = readFileSync(NB + 'Section\\GridViewFollowingNB.xml', 'utf-8')
    for (const k of F.kolom) expect(xml, k).toContain(`>${k}<`)
    for (const l of [F.lihat, F.salinSemua, F.salinNopolis]) expect(xml, l).toContain(`<pyLabel>${l}</pyLabel>`)
  })
  it.skipIf(!existsSync(NB + 'Section\\ViewOpenFollowingPolicy.xml'))('teks popup = ViewOpenFollowingPolicy', () => {
    const xml = readFileSync(NB + 'Section\\ViewOpenFollowingPolicy.xml', 'utf-8')
    for (const t of [F.daftar, F.viewOption, F.filter, F.placeholder, F.refresh, F.back, F.pilihKosong]) expect(xml, t).toContain(t)
  })
  it.skipIf(!existsSync(DDL))('View Option = FilterForViewRetro.xml (tanpa 3 No Slip - K70-3)', () => {
    const xml = readFileSync(DDL, 'utf-8')
    expect(xml).toContain('DATA-PORTAL!FILTERFORVIEWRETRO')
    for (const o of F.opsi) expect(xml).toMatch(new RegExp(`<pyLocalizedValue>${o.label.replace('.', '\\.')}</pyLocalizedValue>\\s*<pyStandardValue>${o.value}</pyStandardValue>`))
    expect(F.opsi.map((o) => o.value)).not.toContain('3')
  })
})

describe('PopupFollowing', () => {
  it('render awal: View Option "Choose" + 6 opsi, Filter berplaceholder, kepala 8 kolom', () => {
    const html = renderToStaticMarkup(<PopupFollowing caseId="NB-UJI" onTutup={() => {}} onTersalin={() => {}} tampilDetail={() => null} />)
    expect(html).toContain(`<option value="" selected="">${F.viewOption}: ${F.pilihKosong}</option>`)
    expect(html.match(/<option value="\d">/g)).toHaveLength(6)
    expect(html).toContain(`placeholder="${F.placeholder}"`)
    for (const k of F.kolom) expect(html).toContain(`<th scope="col">${k}</th>`)
  })
  it('filter kosong tidak mencari; View Details / Copy Nopolis disabled bila Policy Number kosong; Copy All berkonfirmasi', () => {
    expect(SUMBER).toMatch(/if \(o === '' \|\| f\.trim\(\) === ''\) \{\s*setHasil\(null\)/)
    expect(SUMBER).toContain("disabled={b.policyNo === ''} onClick={() => setDetail(b)}")
    expect(SUMBER).toContain("disabled={menyalin || b.policyNo === ''}")
    expect(SUMBER).toContain('onClick={() => setYakinSemua(b)}')
    expect(SUMBER).toContain('onClick={() => void salin(yakinSemua, true)}')
  })
  it('View Details = detail polis bersama (DetailPolis): polis sistem baru = layar case hanya lihat, polis lama Pega = ringkasan (work owner 09-10-2026 "PERBAIKI")', () => {
    expect(SUMBER).toContain('{tampilDetail(detail, () => setDetail(null))}')
    const HALAMAN_INI = readFileSync(join(__dirname, '..', 'pages', 'InwardFacultative.tsx'), 'utf8').replace(/\r\n/g, '\n')
    expect(HALAMAN_INI).toMatch(/tampilDetail=\{\(b, kembali\) => \(\s*<DetailPolis\s+baris=\{b\}\s+judul=""\s+onKembali=\{kembali\}/)
    expect(HALAMAN_INI).toContain('tampilKasus={(caseId) => <InwardFacultative kasus={{ caseId }} onBatal={kembali} hanyaLihat />}')
  })
  it('tanggal grid dd/mm/yyyy', () => {
    expect(tanggalFollowing('05-01-2026')).toBe('05/01/2026')
    expect(tanggalFollowing('')).toBe('')
  })
  it('halaman case: tombol Search membuka popup, aktif hanya bila bolehFollowing; sesudah salin case dimuat ulang', () => {
    expect(HALAMAN).toContain('disabled={hanyaLihat || dariServer?.bolehFollowing !== true}')
    expect(HALAMAN).toContain('onClick={() => setBukaFollowing(true)}')
    expect(HALAMAN).toMatch(/onTersalin=\{\(k\) => \{\s*terapkan\(k\)\s*setVersiDetail\(\(n\) => n \+ 1\)/)
  })
})
