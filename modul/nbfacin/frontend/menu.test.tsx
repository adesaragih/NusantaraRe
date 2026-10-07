// Halaman depan NB FacIn - tiket 25/26: tombol sidebar membuka portal, portal membuka form,
// coverage kargo (tiket 21) tetap terdaftar tetapi bukan halaman awal.

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import type { Sesi } from '../../../inti/frontend/store/sesi'
import { KEPALA_PORTAL, SARING_PORTAL } from './labels'
import { HALAMAN_AWAL_NBFACIN, HALAMAN_NBFACIN, PENDAFTARAN_MENU } from './menu'
import { RuteNbFacIn } from './rute'

const render = (halaman: string) =>
  renderToStaticMarkup(<RuteNbFacIn halaman={halaman} masuk={{} as Sesi} onPindah={() => {}} />)

describe('halaman depan NB FacIn', () => {
  it('halaman awal = portal Opportunity', () => {
    expect(HALAMAN_AWAL_NBFACIN).toBe('nbfacin-portal')
    expect(PENDAFTARAN_MENU.halamanAwal).toBe('nbfacin-portal')
  })

  it('empat halaman terdaftar; coverage kargo tidak dihapus', () => {
    expect([...HALAMAN_NBFACIN]).toEqual(['nbfacin-portal', 'nbfacin-opportunity', 'nbfacin-inward', 'nbfacin-coverage-cargo'])
  })

  it('rute: tiap halaman merender layarnya sendiri saja', () => {
    expect(render('nbfacin-portal')).toContain(`>${KEPALA_PORTAL.buat.label}</button>`)
    expect(render('nbfacin-opportunity')).toContain('Business Prospect Name')
    expect(render('nbfacin-opportunity')).not.toContain(SARING_PORTAL.placeholder.label)
    expect(render('beranda')).toBe('')
  })

  it('Create opportunity memindah ke form Opportunity', () => {
    const rute = readFileSync(join(__dirname, 'rute.tsx'), 'utf8')
    expect(rute).toContain("onBuat={() => onPindah('nbfacin-opportunity')}")
  })

  it('sesudah Create berhasil: case disimpan di rute lalu pindah ke Inward Facultative (tiket 30)', () => {
    const rute = readFileSync(join(__dirname, 'rute.tsx'), 'utf8')
    expect(rute).toMatch(/onDibuat=\{\(k\) => \{\s*setKasus\(k\)\s*onPindah\('nbfacin-inward'\)/)
    // Tanpa case (mis. dibuka langsung), halaman Inward tidak dirender.
    expect(rute).toContain("{halaman === 'nbfacin-inward' && kasus && <InwardFacultative")
    expect(render('nbfacin-inward')).toBe('')
  })

  it('klik Name di portal membuka case tersebut di layar Inward Facultative (tiket 32)', () => {
    const rute = readFileSync(join(__dirname, 'rute.tsx'), 'utf8')
    expect(rute).toMatch(/onBuka=\{\(caseId\) => \{\s*setKasus\(\{ caseId \}\)\s*onPindah\('nbfacin-inward'\)/)
  })
})
