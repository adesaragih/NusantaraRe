// Uji popup konfirmasi hapus — tiket 10 Treaty Contract Out.

import { existsSync, readdirSync, readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import { AKAR_APLIKASI } from '../../../../inti/frontend/uji/sumber'
import { HAPUS_TCO } from '../labels'
import { rincianDampak, teksBersama } from './KonfirmasiHapusTCO'

/**
 * Seluruh klien backend: dulu SATU `services/api.ts`, kini `inti/klien.ts` +
 * `modul/<nama>/api.ts` (refactor bentuk B). Pemeriksaan NEGATIF membaca
 * semuanya - membaca satu pecahan saja menyempitkan penjaga diam-diam.
 */
function semuaApi(): string {
  // Satu folder per modul (30-09-2026): `modul/<nama>/frontend/api.ts`.
  const akar = join(AKAR_APLIKASI, 'modul')
  const modul = readdirSync(akar, { withFileTypes: true })
    .filter((d) => d.isDirectory() && existsSync(join(akar, d.name, 'frontend', 'api.ts')))
    .map((d) => join(akar, d.name, 'frontend', 'api.ts'))
  return [join(AKAR_APLIKASI, 'inti', 'frontend', 'klien.ts'), ...modul].map((f) => readFileSync(f, 'utf8')).join('\n')
}

const baca = (b: string) => readFileSync(join(__dirname, b), 'utf8')
const d = { kontrak: 1, reinsurer: 2, security: 3, business: 1, klausulTetap: 5, bersama: 0 }

describe('popup konfirmasi hapus', () => {
  it('menyebut jumlah tiap jenis (AC 43)', () => {
    expect(rincianDampak(d, 'kontrak')).toEqual(['2 reinsurer', '3 security', '1 business'])
    expect(rincianDampak(d, 'reinsurer')).toEqual(['3 security'])
  })
  it('01-10-2026: catatan klausul DIBUANG dari popup (klausul tetap tidak dihapus server, AC 44)', () => {
    expect(HAPUS_TCO).not.toHaveProperty('klausulTetap')
    expect(baca('KonfirmasiHapusTCO.tsx')).not.toMatch(/dampak\.klausulTetap|HAPUS_TCO\.klausulTetap/)
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
    expect(teksBersama(1)).toMatch(/ARE deleted/)
    expect(teksBersama(2)).toMatch(/ARE deleted/)
    expect(baca('KonfirmasiHapusTCO.tsx')).toContain('role="alert"')
    expect(baca('KonfirmasiHapusTCO.tsx')).toContain('<strong>{dampak.bersama}</strong> {teksBersama(dampak.bersama)}')
  })
  it('01-10-2026: kalimat kontrak lain tunggal/jamak benar, business NULL ikut disebut', () => {
    expect(teksBersama(1)).toMatch(/^other contract uses .* its reinsurers/)
    expect(teksBersama(2)).toMatch(/^other contracts use .* their reinsurers/)
    expect(teksBersama(3)).toBe(teksBersama(2))
    for (const n of [1, 2]) {
      expect(teksBersama(n)).not.toMatch(/another/)
      expect(teksBersama(n)).toMatch(/belong to this treaty year or have no treaty year\.$/)
    }
  })
  it('OQ-TCO-21: DELETE kontrak mengirim cacah kontrak lain', () => {
    expect(readFileSync(join(__dirname, '..', 'api.ts'), 'utf8')).toContain('bersama: String(d.bersama)')
  })
  it('OQ-TCO-19: klien simpan kontrak utuh dibuang', () => {
    expect(semuaApi()).not.toMatch(/simpanKontrakUtuh|kontrak-utuh/)
  })
})
