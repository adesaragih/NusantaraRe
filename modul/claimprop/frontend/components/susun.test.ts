// Penyusun rupa layar kasus Claim Prop (pola Kelola User, keputusan work owner 08-10-2026): kartu bertitel, isi
// form-grid, label pendamping menempel ke medannya. Pohon tata server tidak diubah isinya.

import { describe, expect, it } from 'vitest'

import type { Tata } from '../api'
import { susunIsi, susunLayar } from './susun'

const medan = (jalur: string, label = ''): Tata => ({ jenis: 'medan', jalur, label, kendali: 'teks' })
const label = (teks: string): Tata => ({ jenis: 'label', label: teks })
const tombol = (id: string, teks = id): Tata => ({ jenis: 'tombol', id, label: teks, aksi: id })
const bagian = (judul: string, ...anak: Tata[]): Tata => ({ jenis: 'bagian', label: judul, anak })
const grid = (jalur: string): Tata => ({ jenis: 'grid', jalur })

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

describe('susunLayar - kartu tingkat layar', () => {
  it('Outstanding Claim: label judul membuka kartu, bagian berlabel kartu sendiri, tombol akhir baris aksi', () => {
    const b = susunLayar([
      label('Outstanding Claim'),
      label('Claim No      ..........'),
      bagian('Claim Treaty', tombol('ChooseMaster'), medan('IDMaster', 'Treaty ID')),
      label('Claim Information'),
      tombol('ChoosePolicy'),
      medan('PolicyNo', 'Policy No'),
      label('Q'),
      medan('Quater'),
      bagian('Interest', bagian('', grid('I'))),
      bagian('Estimation', bagian('', grid('E'))),
      tombol('Save'),
      tombol('SaveToIssueRNM'),
      bagian('Claim History', grid('R')),
    ])
    expect(b.map((x) => (x.jenis === 'panel' ? x.judul : 'aksi'))).toEqual([
      'Outstanding Claim',
      'Claim Treaty',
      'Claim Information',
      'Interest',
      'Estimation',
      'aksi',
      'Claim History',
    ])
    const info = b[2]
    expect(
      info?.jenis === 'panel' && susunIsi(info.isi).map((x) => (x.jenis === 'medan' ? x.t.label : x.jenis)),
    ).toEqual(['tombol', 'Policy No', 'Q'])
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

  it('bagian tanpa label di tingkat layar diratakan ke kartu yang terbuka (Claim No terisi)', () => {
    const b = susunLayar([label('Acceptation Claim'), bagian('', medan('NoClaim', 'Claim No')), bagian('Claim Treaty')])
    expect(b[0]).toMatchObject({ jenis: 'panel', judul: 'Acceptation Claim', isi: [{ jalur: 'NoClaim' }] })
  })
})
