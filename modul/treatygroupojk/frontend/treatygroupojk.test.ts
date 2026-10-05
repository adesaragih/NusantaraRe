import { readdirSync, readFileSync, statSync } from 'node:fs'
import { join } from 'node:path'

import { afterEach, describe, expect, it, vi } from 'vitest'

import { ambil, ambilDaftar, tambah, ubah } from './api'
import { isianDari, periksaIsian, rapikanIsian } from './aturan'
import { MENU_TGO, TGO } from './labels'

describe('label Treaty Group OJK', () => {
  it('nama menu = M_NAV_MENU.LABEL migrasi inti 916, golongan MASTER', () => {
    const sql = readFileSync(`${__dirname}/../../../inti/backend/migrations/916_m_nav_menu_treatygroupojk.sql`, 'utf8')
    expect(sql).toContain(`'treatygroupojk', '${MENU_TGO.kelompok}', 'MASTER', 'treatygroupojk'`)
  })
})

describe('aturan Treaty Group OJK', () => {
  it('isian: kosong untuk Add, terisi untuk Edit; spasi tepi dibuang', () => {
    expect(isianDari(null)).toEqual({ name: '', nameIdn: '' })
    const o = { id: '01', name: 'UJI', nameIdn: 'UJI IDN', orderNo: '1' }
    expect(isianDari(o)).toEqual({ name: 'UJI', nameIdn: 'UJI IDN' })
    expect(rapikanIsian({ name: ' uji ', nameIdn: ' idn ' })).toEqual({ name: 'uji', nameIdn: 'idn' })
  })

  it('wajib: Name dan Name (IDN)', () => {
    expect(periksaIsian({ name: ' ', nameIdn: 'X' })).toBe(TGO.galatNama)
    expect(periksaIsian({ name: 'X', nameIdn: '' })).toBe(TGO.galatNamaIdn)
    expect(periksaIsian({ name: 'X', nameIdn: 'X' })).toBeNull()
  })

  it('Order No tidak tampil di daftar, form, maupun View (perintah work owner 05-10-2026: "ORDERNO hide aja")', () => {
    for (const berkas of ['pages/TreatyGroupOjk.tsx', 'components/FormOjk.tsx']) {
      const isi = readFileSync(`${__dirname}/${berkas}`, 'utf8')
      expect(isi, berkas).not.toMatch(/orderNo|TGO\.order/)
    }
  })
})

describe('klien Treaty Group OJK', () => {
  afterEach(() => {
    vi.unstubAllGlobals()
  })

  it('rute dan metode sama dengan handlers/rute.go; tanpa hapus', async () => {
    const panggil: string[] = []
    vi.stubGlobal(
      'fetch',
      vi.fn((url: string, init: RequestInit) => {
        panggil.push(`${String(init.method ?? 'GET')} ${url.replace(/^https?:\/\/[^/]+/, '')}`)
        return Promise.resolve(new Response('{}', { status: 200, headers: { 'Content-Type': 'application/json' } }))
      }),
    )
    await ambilDaftar(' uji ')
    await ambilDaftar('')
    await ambil('01')
    await tambah({ name: 'UJI', nameIdn: 'UJI' })
    await ubah('01', { name: 'UJI', nameIdn: 'UJI' })
    expect(panggil).toEqual([
      'GET /api/treaty-group-ojk?q=uji',
      'GET /api/treaty-group-ojk',
      'GET /api/treaty-group-ojk/01',
      'POST /api/treaty-group-ojk',
      'PUT /api/treaty-group-ojk/01',
    ])
  })
})

describe('tanpa popup View (perintah work owner 05-10-2026: "yang view nya popup ... ga usah di buat, toh bisa liat langsung dari list tabelnya")', () => {
  it('tidak ada mode lihat, tombol View, atau komponen Lihat*', () => {
    const berkas = (dir: string): string[] =>
      readdirSync(dir).flatMap((n) => {
        const p = join(dir, n)
        if (statSync(p).isDirectory()) return berkas(p)
        return /\.tsx?$/.test(n) && !/\.test\.tsx?$/.test(n) ? [p] : []
      })
    for (const p of berkas(__dirname)) {
      const isi = readFileSync(p, 'utf8')
      expect(isi, p).not.toMatch(/jenis: 'lihat'|[A-Z]+\.view\b|judulLihat|function Lihat|from '\.\/Lihat/)
      expect(p, p).not.toMatch(/Lihat[A-Z][a-z]*\.tsx$/)
    }
  })
})
