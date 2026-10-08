// Penyusun rupa layar kasus Claim Prop (pola Kelola User, keputusan work owner 08-10-2026): kartu bertitel, isi
// form-grid, label pendamping menempel ke medannya. Pohon tata server tidak diubah isinya.

import { describe, expect, it } from 'vitest'

import type { Tata } from '../api'
import { jumlahHalaman, potongHalaman, susunIsi, susunLayar } from './susun'

const medan = (jalur: string, label = ''): Tata => ({ jenis: 'medan', jalur, label, kendali: 'teks' })
const label = (teks: string): Tata => ({ jenis: 'label', label: teks })
const tombol = (id: string, teks = id): Tata => ({ jenis: 'tombol', id, label: teks, aksi: id })
const bagian = (judul: string, ...anak: Tata[]): Tata => ({ jenis: 'bagian', label: judul, anak })
const grid = (jalur: string): Tata => ({ jenis: 'grid', jalur })
const letak = (l: NonNullable<Tata['letak']>, ...anak: Tata[]): Tata => ({ jenis: 'bagian', letak: l, anak })

describe('susunIsi - isi satu kartu', () => {
  it('label tepat sebelum medan tanpa label menjadi label medan itu (Q, /, U/Y)', () => {
    const u = susunIsi([medan('PolicyNo', 'Policy No'), label('Q'), medan('Quater'), label('/'), medan('Year')])
    expect(u.map((x) => (x.jenis === 'medan' ? x.t.label : x.jenis))).toEqual(['Policy No', 'Q', '/'])
  })

  it('label "%" sesudah medan menjadi satuan medan itu, bukan judul', () => {
    const u = susunIsi([medan('RNMShareP', 'RNM Share'), label('%'), medan('Lain', 'Lain')])
    expect(u).toHaveLength(2)
    expect(u[0]).toMatchObject({ jenis: 'medan', t: { label: 'RNM Share', satuan: '%' } })
  })

  it('bagian tanpa label diratakan; bagian berlabel menjadi sub-bagian', () => {
    const u = susunIsi([bagian('', medan('A', 'A'), bagian('Spreading List', grid('X'))), medan('B', 'B')])
    expect(u.map((x) => x.jenis)).toEqual(['medan', 'sub', 'medan'])
  })

  it('tombol sesudah medan menempel ke medan itu; ikon tanpa label menempel walau di akhir', () => {
    const u = susunIsi([
      medan('CauseOfLoss', 'Cause of Loss'),
      tombol('ChooseCauseOfLoss', 'Choose Cause of Loss'),
      medan('Consultant', 'Consultant ID'),
      tombol('Ikon', ''),
    ])
    expect(u).toHaveLength(2)
    expect(u[0]).toMatchObject({ jenis: 'medan', tombol: [{ id: 'ChooseCauseOfLoss' }] })
    expect(u[1]).toMatchObject({ jenis: 'medan', tombol: [{ id: 'Ikon' }] })
  })

  // work owner 08-10-2026 ("masa add nya begitu"): ikon Add berlabel di samping Consultant ID / Adjuster ID jatuh ke
  // baris sendiri saat Consultant Name / Adjuster Name tersembunyi - ikon tetap menempel walau berlabel dan di akhir.
  it('tombol ikon berlabel di akhir isi tetap menempel ke medan sebelumnya', () => {
    const u = susunIsi([
      medan('ConsultantID', 'Consultant ID'),
      { ...tombol('AdjusterConsultantBaru1', 'Add'), ikon: 'tambah' },
    ])
    expect(u).toHaveLength(1)
    expect(u[0]).toMatchObject({ jenis: 'medan', tombol: [{ id: 'AdjusterConsultantBaru1' }] })
  })

  it('tombol di awal = baris alat; tombol berlabel di akhir = baris aksi, tidak menempel', () => {
    const u = susunIsi([tombol('ChooseMaster'), medan('IDMaster', 'Treaty ID'), tombol('Save'), tombol('PrintPLA')])
    expect(u.map((x) => x.jenis)).toEqual(['tombol', 'medan', 'tombol'])
    expect(u[0]).toMatchObject({ akhir: false })
    expect(u[2]).toMatchObject({ akhir: true, tombol: [{ id: 'Save' }, { id: 'PrintPLA' }] })
  })

  it('label sisa = judul baris; grid = baris penuh', () => {
    expect(susunIsi([label('Claim Estimation'), grid('G')]).map((x) => x.jenis)).toEqual(['judul', 'grid'])
  })
})

describe('susunIsi - letak layout XML', () => {
  it('bagian berletak (dua / sebaris) tidak diratakan; label Q / U/Y di dalam baris sebaris menempel ke medannya', () => {
    const qy = letak('sebaris', label('Q'), medan('Quater'), label('U/Y'), medan('TreatyYear'))
    qy.label = 'Quarter/Year'
    const u = susunIsi([
      medan('PolicyNo', 'Policy No'),
      qy,
      letak('dua', bagian('', medan('A', 'A')), bagian('', medan('B', 'B'))),
    ])
    expect(u.map((x) => x.jenis)).toEqual(['medan', 'letak', 'letak'])
    expect(susunIsi(qy.anak ?? []).map((x) => (x.jenis === 'medan' ? x.t.label : x.jenis))).toEqual(['Q', 'U/Y'])
  })
})

describe('susunLayar - kartu tingkat layar', () => {
  it('Outstanding Claim: kepala tengah, kartu Claim Treaty, kartu Claim Information, kartu tab membawa RNM Share', () => {
    const b = susunLayar([
      letak('judul', label('Outstanding Claim'), label('Claim No      ..........')),
      bagian(
        'Claim Treaty',
        letak('sebaris', tombol('ChooseMaster')),
        letak('dua', bagian('', medan('IDMaster', 'Treaty ID'))),
      ),
      label('Claim Information'),
      letak('dua', bagian('', medan('PolicyNo', 'Policy No')), bagian('', medan('DateOfLoss', 'Date of Loss'))),
      medan('Location', 'Location of Loss'),
      letak('sebaris', medan('RNMShareP', 'RNM Share'), label('%')),
      letak('tab', bagian('Interest', grid('I')), bagian('Estimation', grid('E'))),
      tombol('Save'),
      tombol('SaveToIssueRNM'),
      bagian('Claim History', grid('R')),
    ])
    expect(b.map((x) => (x.jenis === 'panel' ? x.judul || 'tab' : x.jenis))).toEqual([
      'kepala',
      'Claim Treaty',
      'Claim Information',
      'tab',
      'aksi',
      'Claim History',
    ])
    const info = b[2]
    expect(info?.jenis === 'panel' && info.isi.map((x) => x.letak ?? x.jalur)).toEqual(['dua', 'Location'])
    const tab = b[3]
    expect(tab?.jenis === 'panel' && tab.isi.map((x) => x.letak)).toEqual(['sebaris', 'tab'])
  })

  it('bagian berlabel kecil di tengah kartu terbuka (Information) tetap di kartu itu', () => {
    const b = susunLayar([
      label('Claim Information'),
      tombol('SummaryOutstanding'),
      bagian('Information', medan('Message')),
      tombol('ChoosePolicy'),
      medan('PolicyNo', 'Policy No'),
      bagian('Interest', grid('I')),
    ])
    expect(b.map((x) => (x.jenis === 'panel' ? x.judul : 'aksi'))).toEqual(['Claim Information', 'Interest'])
  })

  it('bagian tanpa label dan tanpa letak di tingkat layar diratakan ke kartu yang terbuka', () => {
    const b = susunLayar([label('Acceptation Claim'), bagian('', medan('NoClaim', 'Claim No')), bagian('Claim Treaty')])
    expect(b[0]).toMatchObject({ jenis: 'panel', judul: 'Acceptation Claim', isi: [{ jalur: 'NoClaim' }] })
  })
})

describe('paging grid (pyGridPaginator)', () => {
  it('lima baris per halaman; tanpa paging = semua baris', () => {
    const baris = [1, 2, 3, 4, 5, 6, 7]
    expect(potongHalaman(baris, 5, 1)).toEqual([1, 2, 3, 4, 5])
    expect(potongHalaman(baris, 5, 2)).toEqual([6, 7])
    expect(potongHalaman(baris, undefined, 1)).toEqual(baris)
    expect(jumlahHalaman(7, 5)).toBe(2)
    expect(jumlahHalaman(0, 5)).toBe(1)
  })
})
