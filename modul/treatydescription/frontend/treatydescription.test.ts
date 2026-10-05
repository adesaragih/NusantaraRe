import { readFileSync } from 'node:fs'

import { afterEach, describe, expect, it, vi } from 'vitest'

import { ambil, ambilDaftar, tambah, ubah, type Desc } from './api'
import { isianDari, labelJenis, labelStatus, periksaIsian, rapikanIsian } from './aturan'
import { MENU_TD, TD } from './labels'

const baris: Desc = {
  id: '10001',
  descName: 'UJI TREATY LIMIT',
  isXol: '0',
  statusAktif: '',
  aktif: true,
}

describe('label Treaty Description', () => {
  it('nama menu = M_NAV_MENU.LABEL migrasi inti 920, golongan MASTER TREATY', () => {
    const sql = readFileSync(
      `${__dirname}/../../../inti/backend/migrations/920_m_nav_menu_treatydescription.sql`,
      'utf8',
    )
    expect(sql).toContain(`'treatydescription', '${MENU_TD.kelompok}', 'MASTER TREATY', 'treatydescription'`)
  })
})

describe('aturan Treaty Description', () => {
  it('isian bawaan Add: Non XOL, Active; Edit: status NULL dibaca Active', () => {
    expect(isianDari(null)).toEqual({ descName: '', isXol: '0', statusAktif: '1' })
    expect(isianDari(baris)).toEqual({ descName: 'UJI TREATY LIMIT', isXol: '0', statusAktif: '1' })
    expect(isianDari({ ...baris, isXol: '1', statusAktif: '0', aktif: false })).toEqual({
      descName: 'UJI TREATY LIMIT',
      isXol: '1',
      statusAktif: '0',
    })
  })

  it('nama dirapikan huruf besar; wajib; maks. 100 byte', () => {
    expect(rapikanIsian({ descName: '  uji pla ', isXol: '1', statusAktif: '0' })).toEqual({
      descName: 'UJI PLA',
      isXol: '1',
      statusAktif: '0',
    })
    expect(periksaIsian({ descName: ' ', isXol: '0', statusAktif: '1' })).toBe(TD.galatNama)
    expect(periksaIsian({ descName: 'A'.repeat(101), isXol: '0', statusAktif: '1' })).toBe(TD.galatPanjang)
    expect(periksaIsian({ descName: 'A'.repeat(100), isXol: '0', statusAktif: '1' })).toBeNull()
  })

  it('label jenis dan status', () => {
    expect(labelJenis('0')).toBe('Non XOL')
    expect(labelJenis('1')).toBe('XOL')
    expect(labelStatus(true)).toBe('Active')
    expect(labelStatus(false)).toBe('Inactive')
  })

  it('hanya kolom TREATYDESC: nol data tabel lain di layar (perintah work owner 05-10-2026)', () => {
    for (const berkas of [
      'pages/TreatyDescription.tsx',
      'components/FormDesc.tsx',
      'api.ts',
      'labels.ts',
      'aturan.ts',
    ]) {
      const isi = readFileSync(`${__dirname}/${berkas}`, 'utf8')
      expect(isi, berkas).not.toMatch(/dipakai|PROPORTIONALARRG|Used in|Treaty Contract Out/)
    }
  })

  it('tanpa tombol hapus di layar mana pun', () => {
    for (const berkas of ['pages/TreatyDescription.tsx', 'components/FormDesc.tsx', 'api.ts']) {
      const isi = readFileSync(`${__dirname}/${berkas}`, 'utf8')
      expect(isi, berkas).not.toMatch(/DELETE|hapus\(|Delete/)
    }
  })
})

describe('klien Treaty Description', () => {
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
    await ambilDaftar({ q: ' limit ', xol: '1', status: '0' })
    await ambilDaftar({ q: '', xol: '', status: '' })
    await ambil('10001')
    await tambah({ descName: 'UJI', isXol: '0', statusAktif: '1' })
    await ubah('10001', { descName: 'UJI', isXol: '0', statusAktif: '1' })
    expect(panggil).toEqual([
      'GET /api/treaty-description?q=limit&xol=1&status=0',
      'GET /api/treaty-description',
      'GET /api/treaty-description/10001',
      'POST /api/treaty-description',
      'PUT /api/treaty-description/10001',
    ])
  })
})
