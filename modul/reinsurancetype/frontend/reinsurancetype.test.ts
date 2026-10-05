import { readdirSync, readFileSync, statSync } from 'node:fs'
import { join } from 'node:path'

import { afterEach, describe, expect, it, vi } from 'vitest'

import { ambil, ambilDaftar, tambah, ubah, type Jenis } from './api'
import { FLAG, GROUP_TYPE, isianDari, kelasFlag, opsiDengan, periksaIsian, rapikanIsian, TYPE } from './aturan'
import { LABEL_TYPE, MENU_RT, RT } from './labels'

const lama: Jenis = {
  id: '10033',
  name: 'UJI POOL SURPLUS',
  type: '',
  soaName: '',
  code: '00',
  flag: '',
  noUrut: '',
  groupType: '',
  userId: '',
  tglUpdate: '',
  diubah: '',
}

describe('label Reinsurance Type', () => {
  it('nama menu = M_NAV_MENU.LABEL migrasi inti 921, golongan MASTER TREATY, langsung menyala', () => {
    const sql = readFileSync(`${__dirname}/../../../inti/backend/migrations/921_m_nav_menu_reinsurancetype.sql`, 'utf8')
    expect(sql).toContain(`'reinsurancetype', '${MENU_RT.kelompok}', 'MASTER TREATY', 'reinsurancetype', 8, '1'`)
  })

  it('label Type = Prompt value Pega', () => {
    expect(LABEL_TYPE).toEqual({ '1': 'Own Retention', '2': 'Treaty Out', '3': 'Facultative', '4': 'Treaty In' })
    expect(RT.labelType('2')).toBe('Treaty Out')
    expect(RT.labelType('9')).toBe('9')
    expect(RT.labelType('')).toBe('')
    expect([TYPE, FLAG, GROUP_TYPE]).toEqual([
      ['1', '2', '3', '4'],
      ['active', 'inactive'],
      ['OR', 'QS', 'RI', 'SPL'],
    ])
  })
})

describe('aturan Reinsurance Type', () => {
  it('isian: Add Flag active dan Code 00; Edit membawa nilai warisan apa adanya', () => {
    expect(isianDari(null)).toMatchObject({ flag: 'active', code: '00', type: '' })
    expect(isianDari({ ...lama, flag: '1' })).toMatchObject({ flag: '1', type: '' })
    expect(rapikanIsian({ ...isianDari(null), name: ' uji ', code: ' 30 ' })).toMatchObject({ name: 'uji', code: '30' })
  })

  it('periksa: Type dan Flag wajib untuk baris baru; kosong warisan boleh tetap', () => {
    const baru = { ...isianDari(null), name: 'UJI', type: '1' }
    expect(periksaIsian(baru, null)).toBeNull()
    expect(periksaIsian({ ...baru, name: ' ' }, null)).toBe(RT.galatNama)
    expect(periksaIsian({ ...baru, type: '' }, null)).toBe(RT.galatType)
    expect(periksaIsian({ ...baru, flag: '' }, null)).toBe(RT.galatFlag)
    expect(periksaIsian({ ...baru, code: '3A' }, null)).toBe(RT.galatCode)
    expect(periksaIsian({ ...baru, noUrut: 'x' }, null)).toBe(RT.galatNoUrut)
    expect(periksaIsian(isianDari(lama), lama)).toBeNull()
  })

  it('opsi: nilai warisan di luar pilihan tampil sebagai nilai lama; lencana Flag', () => {
    expect(opsiDengan(FLAG, '1')).toContainEqual({ value: '1', label: RT.nilaiLama('1') })
    expect(opsiDengan(FLAG, 'active')).toHaveLength(2)
    expect(opsiDengan(TYPE, '', RT.labelType)[0]).toEqual({ value: '1', label: 'Own Retention' })
    expect(kelasFlag('active')).toContain('lencana--aktif')
    expect(kelasFlag('inactive')).toContain('lencana--nonaktif')
    expect(kelasFlag('1')).toBe('reinsurancetype__lencana')
  })

  it('Code tidak tampil di form maupun View (perintah work owner 05-10-2026: "Code HAPUS"); nilainya tetap terkirim apa adanya', () => {
    for (const berkas of ['components/FormJenis.tsx', 'pages/ReinsuranceType.tsx']) {
      expect(readFileSync(`${__dirname}/${berkas}`, 'utf8'), berkas).not.toMatch(/RT\.code|isi\.code|j\.code/)
    }
    expect(isianDari({ ...lama, code: '30' }).code).toBe('30')
    expect(isianDari(null).code).toBe('00')
  })

  it('daftar tanpa kolom Code, No Urut, dan Group Type (perintah work owner 05-10-2026: "hapus tampilan ini")', () => {
    const halaman = readFileSync(`${__dirname}/pages/ReinsuranceType.tsx`, 'utf8')
    expect(halaman).not.toMatch(/RT\.(code|noUrut|groupType)|j\.(code|noUrut|groupType)/)
  })
})

describe('klien Reinsurance Type', () => {
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
    const isi = isianDari(null)
    await ambilDaftar(' qs ', '2', 'active')
    await ambilDaftar('', '', '')
    await ambil('10007')
    await tambah(isi)
    await ubah('10007', isi)
    expect(panggil).toEqual([
      'GET /api/reinsurance-type?q=qs&type=2&flag=active',
      'GET /api/reinsurance-type',
      'GET /api/reinsurance-type/10007',
      'POST /api/reinsurance-type',
      'PUT /api/reinsurance-type/10007',
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
