// Penjaga bukti label NB FacIn - tiket 21.
//
// ⛔ Komentar bukti yang tidak pernah diperiksa adalah HIASAN. Berkas ini membuka
// section korpus dan memastikan setiap sel yang disebut `labels.ts` memang berlabel
// dan berproperti seperti tertulis. Korpus READ-ONLY - hanya dibaca.
//
// ⚠️ Bila korpus tidak terjangkau, test DILEWATI dengan pesan - bukan gagal.

import { existsSync, readFileSync } from 'node:fs'

import { describe, expect, it } from 'vitest'

import { MEDAN_COVERAGE_CARGO, TOMBOL_COVERAGE_CARGO } from './labels'

const SECTION = 'D:\\migrasi\\RNM\\NB FacIn\\Section\\InputCoverageCargo_FacIn.xml'
const ada = existsSync(SECTION)

/**
 * Isi elemen sel ber-`pyCellId` tertentu: `rowdata` terdalam yang terbuka saat
 * `<pyCellId>` itu ditemukan, sampai penutup PASANGANNYA (rowdata bersarang - mis.
 * aksi tombol - ikut di dalamnya). Pemindai kedalaman, bukan regex tak-rakus.
 */
function blokSel(xml: string, sel: string): string[] {
  const hasil: string[] = []
  const pola = /<rowdata\b[^>]*>|<\/rowdata>|<pyCellId>(\d+)<\/pyCellId>/g
  const tumpukan: { mulai: number; cocok: boolean }[] = []
  for (const m of xml.matchAll(pola)) {
    const i = m.index ?? 0
    if (m[0].startsWith('<rowdata')) {
      tumpukan.push({ mulai: i, cocok: false })
    } else if (m[0] === '</rowdata>') {
      const atas = tumpukan.pop()
      if (atas?.cocok) hasil.push(xml.slice(atas.mulai, i + m[0].length))
    } else if (m[1] === sel) {
      const atas = tumpukan.at(-1)
      if (atas) atas.cocok = true
    }
  }
  return hasil
}

describe.skipIf(!ada)('label NB FacIn = korpus InputCoverageCargo_FacIn', () => {
  const xml = ada ? readFileSync(SECTION, 'utf-8') : ''

  it.each(Object.entries(MEDAN_COVERAGE_CARGO))('medan %s', (_, m) => {
    const blok = blokSel(xml, m.sel)
    expect(blok.some((b) => b.includes(`<pyLabelFor>${m.label}</pyLabelFor>`) && b.includes(`<pyValue>${m.properti}</pyValue>`))).toBe(true)
  })

  it.each(Object.entries(TOMBOL_COVERAGE_CARGO))('tombol %s', (_, t) => {
    expect(blokSel(xml, t.sel).some((b) => b.includes(`<pyLabel>${t.label}</pyLabel>`))).toBe(true)
  })

  it('uji ini menggigit: label salah tidak ditemukan', () => {
    expect(blokSel(xml, '11').some((b) => b.includes('<pyLabelFor>Rate (‰)</pyLabelFor>'))).toBe(false)
  })
})

if (!ada) {
  it('korpus tidak terjangkau - bukti label dilewati', () => {
    console.warn(`korpus tidak ada di ${SECTION}; uji bukti label NB FacIn dilewati`)
  })
}
