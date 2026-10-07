// Popup Choose Accumulation (tiket 46) - data uji sintetis.

import { readFileSync } from 'node:fs'
import { join } from 'node:path'
import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import { OPSI_KEYWORD, POPUP_AKUMULASI as P, TEKS_AKUMULASI } from '../labels'
import PopupAkumulasi, { ISIAN_KOSONG, jumlahKondisi, keSaring, potongHalaman, UKURAN_HALAMAN_AKUMULASI, zipDariId } from './PopupAkumulasi'

const SUMBER = readFileSync(join(__dirname, 'PopupAkumulasi.tsx'), 'utf8').replace(/\r\n/g, '\n')
const label = (html: string) => [...html.matchAll(/class="field__label"[^>]*>([^<]*)/g)].map((m) => m[1])

describe('zipDariId (SetDataAccum_Act langkah 2-3)', () => {
  it('format NEGARA-ZIP-NOMOR: substring(4,9) tanpa "-"', () => {
    expect(zipDariId('UJI-12345-000001')).toBe('12345')
  })
  it('karakter ke-3 bukan "-": substring(3,8)', () => {
    expect(zipDariId('AB012345X')).toBe('12345')
  })
  it('ID terlalu pendek -> kosong', () => {
    expect(zipDariId('UJI-1')).toBe('')
  })
})

describe('PopupAkumulasi', () => {
  const html = renderToStaticMarkup(<PopupAkumulasi zipRisiko="12345" onTutup={() => {}} onPilih={() => {}} />)

  it('urutan medan seperti SearchRiskAccumCov; Zip Code awal = zip lokasi risiko', () => {
    expect(label(html)).toEqual([
      P.accumulationCode.label, P.policyNo.label, P.road.label,
      P.zipCode.label, P.country.label, P.province.label, P.accumType.label,
      P.city.label, P.district.label, P.area.label, P.czone.label,
      P.keyword.label,
    ])
    for (const t of [P.judul, P.cari.label, P.kondisi, P.filter.label, P.bersih.label]) expect(html).toContain(t)
    expect(html).toContain('value="12345"')
  })

  it('Key Word: pilihan kosong bertulisan kosong lalu 19 kata kunci', () => {
    expect(html).toContain('<option value="" selected=""></option>')
    expect(OPSI_KEYWORD).toHaveLength(19)
  })

  it('Area terpilih menimpa Zip Code saat Filter (GetDataAccumulation_act langkah 4)', () => {
    expect(keSaring({ ...ISIAN_KOSONG, postalCode: '11111', areaZip: '22222' }).postalCode).toBe('22222')
    expect(keSaring({ ...ISIAN_KOSONG, postalCode: '11111' }).postalCode).toBe('11111')
  })

  it('saran berantai: Province <- Country, City <- ProvinceID, District <- City, Area <- District', () => {
    expect(SUMBER).toContain("saranAkumulasi('province', q, s.nation)")
    expect(SUMBER).toContain("saranAkumulasi('city', q, s.provinceId)")
    expect(SUMBER).toContain("saranAkumulasi('district', q, s.city)")
    expect(SUMBER).toContain("saranAkumulasi('area', q, s.district)")
  })

  it('Choose: zip beda -> pesan verbatim, popup tetap; sama -> onPilih(ID, Note)', () => {
    expect(TEKS_AKUMULASI.zipBeda('12345', '54321')).toBe('ZIpCode Harus Sama dengan ZIpCode  Yang di Object Item >>>> 12345 != 54321')
    expect(SUMBER).toMatch(/if \(dariId !== zipRisiko\) \{\s*setPesan\(TEKS_AKUMULASI\.zipBeda\(zipRisiko, dariId\)\)\s*return\s*\}\s*onPilih\(b\.id, b\.note\)/)
  })
})

describe('PopupAkumulasi - tampilan dirapikan (04-10-2026)', () => {
  const html = renderToStaticMarkup(<PopupAkumulasi zipRisiko="12345" onTutup={() => {}} onPilih={() => {}} />)
  it('kepala: keterangan + zip lokasi risiko; Filter = tombol kirim form (Enter menyaring)', () => {
    expect(html).toContain(TEKS_AKUMULASI.sub)
    expect(html).toContain(`${TEKS_AKUMULASI.zipLokasi}<strong>12345</strong>`)
    expect(html).toMatch(/<form class="nbf-akum__kartu"/)
    expect(html).toMatch(new RegExp(`<button type="submit" class="btn btn--primary">${P.filter.label}</button>`))
    expect(html).toContain(TEKS_AKUMULASI.petunjukAwal)
  })
  it('jumlah kondisi aktif menghitung teks yang terisi', () => {
    expect(jumlahKondisi(ISIAN_KOSONG)).toBe(0)
    expect(jumlahKondisi({ ...ISIAN_KOSONG, postalCode: '12345', city: 'UJI', cityId: '9', keyword: 'RESTORAN' })).toBe(3)
  })
  it('baris hasil diberi tanda zip sesuai / berbeda sebelum Choose', () => {
    expect(SUMBER).toContain('const cocok = zipDariId(b.id) === zipRisiko')
    expect(SUMBER).toContain('cocok ? TEKS_AKUMULASI.zipSesuai : TEKS_AKUMULASI.zipBerbeda')
  })
})

describe('PopupAkumulasi - hasil 10 per halaman (04-10-2026)', () => {
  const baris = Array.from({ length: 23 }, (_, i) => i + 1)
  it('ukuran 10; halaman 1 / 3 / di luar jangkauan dijepit', () => {
    expect(UKURAN_HALAMAN_AKUMULASI).toBe(10)
    expect(potongHalaman(baris, 1)).toEqual([1, 2, 3, 4, 5, 6, 7, 8, 9, 10])
    expect(potongHalaman(baris, 3)).toEqual([21, 22, 23])
    expect(potongHalaman(baris, 9)).toEqual([21, 22, 23])
    expect(potongHalaman([], 1)).toEqual([])
  })
  it('pencarian baru kembali ke halaman 1; pager hanya bila lebih dari 10 baris', () => {
    expect(SUMBER).toMatch(/setHasil\(\(await cariAkumulasi\(keSaring\(s\)\)\)\.baris\)\s*setHalaman\(1\)/)
    expect(SUMBER).toContain('hasil.length > UKURAN_HALAMAN_AKUMULASI')
    expect(TEKS_AKUMULASI.menampilkan(2, 10, 23)).toBe('Menampilkan 11–20 dari 23 data')
  })
})
