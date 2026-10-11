// "This name is on Agent Negative List" — `Activity/TreatyInCheckCedingBlacklist`
// [2]/[3]: pesan MEDAN pada Ceding / Source of Business, tidak menghalangi.

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import type { HasilDaftarNegatifAgen } from './api'
import { PesanMedanAgen, pesanUntukPengenal } from './components/PesanDaftarNegatif'

const PESAN = 'This name is on Agent Negative List'

const hasil: HasilDaftarNegatifAgen = {
  cedant: { id: 'G1', ditemukan: true, statusAktif: 'inactive', daftarNegatif: true, pesan: PESAN },
  asalBisnis: { id: 'S1', ditemukan: true, statusAktif: '1', daftarNegatif: false },
}

describe('TreatyInCheckCedingBlacklist — pesan medan', () => {
  it('pesan hanya milik medan yang negatif', () => {
    expect(pesanUntukPengenal(hasil, 'G1', 'S1')).toEqual({ cedant: PESAN })
  })

  it('⛔ memilih agen lain melepas pesannya (Choose tidak menjalankan pemeriksaan)', () => {
    expect(pesanUntukPengenal(hasil, 'G2', 'S1')).toEqual({})
  })

  it('belum diperiksa = tanpa pesan', () => {
    expect(pesanUntukPengenal(null, 'G1', 'S1')).toEqual({})
  })

  it('pesan dirender di bawah medan sebagai teks galat medan', () => {
    expect(renderToStaticMarkup(<PesanMedanAgen pesan={PESAN} />)).toBe(
      `<p class="field__error" role="alert">${PESAN}</p>`,
    )
    expect(renderToStaticMarkup(<PesanMedanAgen />)).toBe('')
  })

  it('⭐ form memeriksa HANYA di mode ubah kontrak lama (tombol Edit), sekali sesudah memuat', () => {
    const isi = readFileSync(join(__dirname, 'pages', 'FormKontrakTreatyIn.tsx'), 'utf8')
    expect(isi).toContain("useDaftarNegatifAgen(bisaUbah && idKontrak !== '' && !memuat, idKontrak, idCedant, idAsalBisnis)")
    expect(isi).toContain('<PesanMedanAgen pesan={pesanAgen.cedant} />')
    expect(isi).toContain('<PesanMedanAgen pesan={pesanAgen.asalBisnis} />')
  })
})
