import { readdirSync, readFileSync, statSync } from 'node:fs'
import { join } from 'node:path'

import { afterEach, describe, expect, it, vi } from 'vitest'

import { ambil, ambilDaftar, ambilPilihan, tambah, ubah, type Grup } from './api'
import { coaTampil, isianDari, kodeNama, opsiOjk, periksaIsian, rapikanIsian } from './aturan'
import { MENU_TG, TG } from './labels'

const grup: Grup = {
  id: '10007',
  oldId: '01',
  ojkId: '01',
  ojkName: 'UJI PROPERTY',
  ojkNameIdn: 'UJI HARTA BENDA',
  orderNo: '1',
  name: 'UJI PROPERTY',
  soaName: '',
  tglUpdate: '',
  diubah: '',
  userId: '',
  coaId: '10013',
  coaName: 'UJI FIRE',
}
const ojk = [
  { id: '01', name: 'UJI PROPERTY', nameIdn: 'UJI HARTA BENDA', orderNo: '7', coaId: '10013', coaName: 'UJI FIRE' },
  { id: '03', name: 'UJI CARGO', nameIdn: 'UJI PENGANGKUTAN', orderNo: '5', coaId: '', coaName: '' },
  { id: '09', name: 'UJI ENGINEERING', nameIdn: 'UJI REKAYASA', orderNo: '2', coaId: '10019', coaName: 'UJI CAR' },
]

describe('label Treaty Group', () => {
  it('nama menu = M_NAV_MENU.LABEL migrasi inti 917, golongan MASTER', () => {
    const sql = readFileSync(`${__dirname}/../../../inti/backend/migrations/917_m_nav_menu_treatygroup.sql`, 'utf8')
    expect(sql).toContain(`'treatygroup', '${MENU_TG.kelompok}', 'MASTER', 'treatygroup'`)
  })
})

describe('aturan Treaty Group', () => {
  it('isian tanpa COAID (opsi A); spasi tepi dibuang; OJK dan nama wajib', () => {
    expect(isianDari(null)).toEqual({ ojkId: '', name: '', soaName: '' })
    expect(isianDari(grup)).toEqual({ ojkId: '01', name: 'UJI PROPERTY', soaName: '' })
    expect(rapikanIsian({ ojkId: ' 01 ', name: ' uji ', soaName: ' s ' })).toEqual({
      ojkId: '01',
      name: 'uji',
      soaName: 's',
    })
    expect(periksaIsian({ ojkId: '', name: 'X', soaName: '' })).toBe(TG.galatOjk)
    expect(periksaIsian({ ojkId: '01', name: ' ', soaName: '' })).toBe(TG.galatNama)
    expect(periksaIsian({ ojkId: '01', name: 'X', soaName: '' })).toBeNull()
  })

  it('kode dan nama; OJK lama yang hilang dari daftar tetap terbaca', () => {
    expect(kodeNama('10013', 'UJI FIRE')).toBe('10013 - UJI FIRE')
    expect(kodeNama('10013', '')).toBe('10013')
    expect(kodeNama('', 'X')).toBe('')
    expect(opsiOjk(ojk, '01').map((o) => o.value)).toEqual(['01', '03', '09'])
    expect(opsiOjk(ojk, '99', 'UJI LAMA')).toContainEqual({ value: '99', label: TG.nilaiLama('99 - UJI LAMA') })
  })

  it('COA ikut OJK yang dipilih (perintah work owner 05-10-2026); OJK tidak diganti = COA baris tersimpan', () => {
    const lain = { ...grup, ojkId: '09', coaId: '10777', coaName: 'UJI LAMA' }
    expect(coaTampil(ojk, '01', null)).toBe('10013 - UJI FIRE')
    expect(coaTampil(ojk, '03', null)).toBe('')
    expect(coaTampil(ojk, '', null)).toBe('')
    expect(coaTampil(ojk, '09', lain)).toBe('10777 - UJI LAMA')
    expect(coaTampil(ojk, '01', lain)).toBe('10013 - UJI FIRE')
  })

  it('Order No tidak tampil di daftar, form, maupun View (perintah work owner 05-10-2026) - tetap disalin backend', () => {
    for (const berkas of ['pages/TreatyGroup.tsx', 'components/FormGrup.tsx']) {
      const isi = readFileSync(`${__dirname}/${berkas}`, 'utf8')
      expect(isi, berkas).not.toMatch(/orderNo|orderTampil|TG\.order/)
    }
  })
})

describe('klien Treaty Group', () => {
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
    await ambilDaftar(' uji ', '01')
    await ambilDaftar('', '')
    await ambilPilihan()
    await ambil('10007')
    await tambah({ ojkId: '01', name: 'UJI', soaName: '' })
    await ubah('10007', { ojkId: '01', name: 'UJI', soaName: '' })
    expect(panggil).toEqual([
      'GET /api/treaty-group?q=uji&ojk=01',
      'GET /api/treaty-group',
      'GET /api/treaty-group/pilihan',
      'GET /api/treaty-group/10007',
      'POST /api/treaty-group',
      'PUT /api/treaty-group/10007',
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
