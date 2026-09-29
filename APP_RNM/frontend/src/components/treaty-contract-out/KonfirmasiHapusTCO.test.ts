// Uji popup konfirmasi hapus — tiket 10 Treaty Contract Out.

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import { HAPUS_TCO } from '../../assets/labels.treaty-contract-out'
import { rincianDampak } from './KonfirmasiHapusTCO'

const baca = (b: string) => readFileSync(join(__dirname, b), 'utf8')
const d = { kontrak: 1, reinsurer: 2, security: 3, business: 1, klausulTetap: 5, bersama: 0 }

describe('popup konfirmasi hapus', () => {
  it('menyebut jumlah tiap jenis (AC 43)', () => {
    expect(rincianDampak(d, 'kontrak')).toEqual(['2 reinsurer', '3 security', '1 business'])
    expect(rincianDampak(d, 'reinsurer')).toEqual(['3 security'])
  })
  it('menyatakan eksplisit klausul TIDAK terhapus (AC 44)', () => {
    expect(HAPUS_TCO.klausulTetap).toMatch(/TIDAK ikut terhapus/)
    expect(baca('KonfirmasiHapusTCO.tsx')).toContain('{dampak.klausulTetap} {HAPUS_TCO.klausulTetap}')
  })
  it('Batal tidak menghapus; Ya mengirim jumlah yang DILIHAT', () => {
    for (const panel of ['PanelKontrakTahun.tsx', 'PanelReinsurerKombinasi.tsx']) {
      const kode = baca(panel)
      expect(kode, panel).toMatch(/onBatal=\{\(\) => \{\s*setKonfirmasi\(null\)\s*\}\}/)
      expect(kode, panel).toMatch(/hapus(Kontrak|Reinsurer)\(tahun(ID|\.id)[^)]*konfirmasi\.dampak\)/)
    }
  })
  it('temuan /code-review: angka dimuat ulang sesudah galat; kombinasi bersama dinyatakan', () => {
    expect((baca('PanelKontrakTahun.tsx').match(/ambilDampakHapusKontrak\(/g) ?? []).length).toBe(2)
    expect((baca('PanelReinsurerKombinasi.tsx').match(/ambilDampakHapusReinsurer\(/g) ?? []).length).toBe(2)
    expect(baca('KonfirmasiHapusTCO.tsx')).toContain('dampak.bersama > 0')
  })
  it('OQ-TCO-21: kontrak lain terdampak disebut sebagai peringatan dan ikut dikonfirmasi', () => {
    expect(HAPUS_TCO.bersama).toMatch(/IKUT terhapus/)
    expect(baca('KonfirmasiHapusTCO.tsx')).toContain('role="alert"')
  })
  it('OQ-TCO-21: DELETE kontrak mengirim cacah kontrak lain', () => {
    expect(readFileSync(join(__dirname, '..', '..', 'services', 'api.ts'), 'utf8')).toContain('bersama: String(d.bersama)')
  })
  it('OQ-TCO-19: klien simpan kontrak utuh dibuang', () => {
    expect(readFileSync(join(__dirname, '..', '..', 'services', 'api.ts'), 'utf8')).not.toMatch(/simpanKontrakUtuh|kontrak-utuh/)
  })
})
