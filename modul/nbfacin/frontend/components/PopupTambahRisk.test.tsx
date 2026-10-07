// Paritas popup Add alamat risiko (tiket 37) - section Pega `InputRiskAddress` + tangkapan layar work owner
// 03-10-2026. Data uji sintetis.

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import { OPSI_TITLE_RISK, POPUP_TAMBAH_RISK as T } from '../labels'
import PopupTambahRisk, { alamatKosong, alamatLengkap, terapkanRW } from './PopupTambahRisk'

const HTML = renderToStaticMarkup(<PopupTambahRisk onTutup={() => {}} onTersimpan={() => {}} />)
const SUMBER = readFileSync(join(__dirname, 'PopupTambahRisk.tsx'), 'utf8').replace(/\r\n/g, '\n')

describe('PopupTambahRisk - keadaan awal (gambar 1)', () => {
  it('medan tampil: Country, Zip Code, Title, Address - Province/City/District/Territory belum', () => {
    const label = [...HTML.matchAll(/<label class="field__label"[^>]*>([^<]*)/g)].map((m) => m[1])
    expect(label).toEqual([T.country.label, T.zipCode.label, T.title.label, T.address.label])
  })

  it('Title = DESA terpilih; Save = kirim form; Close = tombol batal modal', () => {
    expect(HTML).toMatch(/<option value="DESA" selected="">DESA<\/option>/)
    expect(HTML).toMatch(new RegExp(`<button type="submit" class="btn btn--primary">${T.simpan.label}</button>`))
    expect(HTML).toContain(`>${T.tutup.label}</button>`)
  })
})

describe('PopupTambahRisk - aturan', () => {
  const R = {
    zipCode: '99999', territoryName: 'UJI KEL', districtName: 'UJI KEC', cityName: 'UJI KOTA', provinceName: 'UJI PROV',
    nationName: 'UJI NEGARA',
  }

  it('memilih saran Zip Code mengisi enam medan dari RW; Title dan Address tetap', () => {
    const a = terapkanRW({ ...alamatKosong(), title: 'GANG', address: 'UJI' }, R)
    expect([a.postalCode, a.nationName, a.provinceName, a.cityName, a.districtName, a.territoryName]).toEqual([
      '99999', 'UJI NEGARA', 'UJI PROV', 'UJI KOTA', 'UJI KEC', 'UJI KEL',
    ])
    expect([a.title, a.address]).toEqual(['GANG', 'UJI'])
  })

  it('alamat kosong: Title = pilihan pertama; Zip Code dan Address wajib (I-4)', () => {
    expect(alamatKosong().title).toBe(OPSI_TITLE_RISK[0])
    expect(alamatLengkap(alamatKosong())).toBe(false)
    expect(alamatLengkap({ ...alamatKosong(), postalCode: '99999' })).toBe(false)
    expect(alamatLengkap({ ...alamatKosong(), postalCode: '99999', address: 'UJI' })).toBe(true)
  })

  it('syarat tampil berantai dari isi medan sebelumnya', () => {
    expect(SUMBER).toContain("{a.nationName !== '' && <Field label={T.province.label}")
    expect(SUMBER).toContain("{a.provinceName !== '' && <Field label={T.city.label}")
    expect(SUMBER).toContain("{a.cityName !== '' && <Field label={T.district.label}")
    expect(SUMBER).toContain("{a.districtName !== '' && <Field label={T.territory.label}")
  })

  it('saran Zip Code hanya saat mengetik, >= 3 karakter, sesudah jeda; jawaban lama dibuang', () => {
    expect(SUMBER).toContain('if (!ketikZip || a.postalCode.trim().length < MIN_ZIP)')
    expect(SUMBER).toContain('const MIN_ZIP = 3')
    expect(SUMBER).toContain('if (nomor === nomorPermintaan.current) setSaran(h.baris)')
    expect(SUMBER).toContain('return () => window.clearTimeout(jadwal)')
  })

  it('Save: tidak lengkap -> tidak memanggil server; lengkap -> simpanAlamatBaru lalu onTersimpan(isian, id)', () => {
    expect(SUMBER).toMatch(/setCobaSimpan\(true\)\s*if \(!alamatLengkap\(a\)\) return/)
    expect(SUMBER).toMatch(/const h = await simpanAlamatBaru\(a\)\s*onTersimpan\(a, h\.id\)/)
  })
})
