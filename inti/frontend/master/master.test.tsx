// Uji layar master generik (inti) - data uji sintetis.

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import { klienMaster, type KlienMaster, type MetaMaster } from './api'
import { jumlahHalaman, nomorHalaman } from './DaftarMaster'
import FormMaster, { badan, KUNCI_NEGARA, medanForm, panjangByte, periksa } from './FormMaster'
import { LABEL_KOLOM_MASTER, TEKS_MASTER as T } from './labels'
import { opsiRujukan } from './PilihRujukan'

const PROVINCE: MetaMaster = {
  kunci: 'province', judul: 'Province', idOtomatis: false,
  kolom: [
    { kunci: 'id', kolom: 'ID', lebar: 4000, wajib: true, turunan: false },
    { kunci: 'nationId', kolom: 'NATIONID', lebar: 4000, wajib: false, turunan: false },
    { kunci: 'note', kolom: 'NOTE', lebar: 5, wajib: true, turunan: false },
    { kunci: 'nationName', kolom: 'NATIONNAME', lebar: 4000, wajib: false, turunan: true },
    { kunci: 'createOp', kolom: 'CREATE_OP', lebar: 64, wajib: false, turunan: true },
  ],
  rujukan: [{ kunci: 'nationId', judul: 'Nation', nilai: 'id', nama: 'note' }],
}
const AKUMULASI: MetaMaster = {
  kunci: 'accumulation', judul: 'Accumulation', idOtomatis: true,
  kolom: [
    { kunci: 'id', kolom: 'ID', lebar: 4000, wajib: false, turunan: true },
    { kunci: 'note', kolom: 'NOTE', lebar: 4000, wajib: true, turunan: false },
    { kunci: 'zipCode', kolom: 'ZIPCODE', lebar: 4000, wajib: true, turunan: false },
  ],
  rujukan: [],
}
const KLIEN = klienMaster('/api/master-uji') as KlienMaster

describe('FormMaster (inti)', () => {
  it('Tambah: ID diisi pengguna; turunan tidak menjadi isian', () => {
    expect(medanForm(PROVINCE, true).map((k) => k.kunci)).toEqual(['id', 'nationId', 'note'])
  })
  it('Ubah: ID bukan isian (diambil dari rute)', () => {
    expect(medanForm(PROVINCE, false).map((k) => k.kunci)).toEqual(['nationId', 'note'])
  })
  it('ID otomatis (Accumulation): Tambah meminta negara, tanpa ID', () => {
    expect(medanForm(AKUMULASI, true).map((k) => k.kunci)).toEqual([KUNCI_NEGARA, 'note', 'zipCode'])
    expect(medanForm(AKUMULASI, false).map((k) => k.kunci)).toEqual(['note', 'zipCode'])
  })
  it('periksa: wajib dan lebar BYTE; badan dipangkas tanpa turunan', () => {
    expect(periksa(PROVINCE, true, { id: ' ', note: 'ABCDEF' })).toEqual({ id: T.wajib, note: T.terlaluPanjang(5) })
    expect(panjangByte('é')).toBe(2)
    expect(periksa(PROVINCE, true, { id: 'P1', note: 'éé' })).toEqual({})
    expect(badan(PROVINCE, true, { id: ' P1 ', nationId: 'N1', note: 'UJI', nationName: 'X', createOp: 'Y' })).toEqual({
      id: 'P1',
      nationId: 'N1',
      note: 'UJI',
    })
  })
  it('render Ubah: ID dan kolom turunan baca-saja; kolom rujukan = combobox', () => {
    const html = renderToStaticMarkup(
      <FormMaster
        klien={KLIEN}
        meta={PROVINCE}
        baris={{ aktif: true, id: 'P1', nationId: 'N1', note: 'UJI', nationName: 'NEGARA UJI', createOp: 'uji' }}
        onTutup={() => {}}
        onTersimpan={() => {}}
      />,
    )
    expect(html).toContain(T.judulUbah('Province'))
    expect(html).toContain('NEGARA UJI')
    expect(html).toContain(`${LABEL_KOLOM_MASTER.createOp} (${T.turunan})`)
    expect(html).toContain('role="combobox"')
  })
  it('render Tambah ID otomatis: keterangan ID + Negara', () => {
    const html = renderToStaticMarkup(<FormMaster klien={KLIEN} meta={AKUMULASI} baris={null} onTutup={() => {}} onTersimpan={() => {}} />)
    expect(html).toContain(T.idOtomatis)
    expect(html).toContain(LABEL_KOLOM_MASTER.negara)
  })
})

describe('PilihRujukan / DaftarMaster', () => {
  it('opsi rujukan: nilai disimpan, nama ditampilkan, nilai sebagai keterangan', () => {
    const r = PROVINCE.rujukan[0]!
    expect(opsiRujukan(r, { id: 'INA', note: 'INDONESIA' })).toEqual({ value: 'INA', label: 'INDONESIA', keterangan: 'INA' })
    expect(opsiRujukan(r, { id: 'X1', note: '' })).toEqual({ value: 'X1', label: 'X1', keterangan: undefined })
  })
  it('jumlah halaman minimal 1 (20 per halaman)', () => {
    expect(jumlahHalaman(0, 20)).toBe(1)
    expect(jumlahHalaman(20, 20)).toBe(1)
    expect(jumlahHalaman(41, 20)).toBe(3)
  })
  it('nomor halaman: 1, terakhir, dua di kiri-kanan halaman aktif; jeda = null', () => {
    expect(nomorHalaman(1, 1)).toEqual([1])
    expect(nomorHalaman(1, 3)).toEqual([1, 2, 3])
    expect(nomorHalaman(6, 20)).toEqual([1, null, 4, 5, 6, 7, 8, null, 20])
    expect(nomorHalaman(20, 20)).toEqual([1, null, 18, 19, 20])
    expect(nomorHalaman(4, 20)).toEqual([1, 2, 3, 4, 5, 6, null, 20])
  })
  it('teks rentang tampil', () => {
    expect(T.menampilkan(2, 20, 135)).toBe('Menampilkan 21–40 dari 135 data')
    expect(T.menampilkan(7, 20, 135)).toBe('Menampilkan 121–135 dari 135 data')
  })
})

describe('kamus label = kolom definisi master backend (dua arah)', () => {
  const GO = readFileSync(join(__dirname, '..', '..', 'backend', 'master', 'models', 'daftar.go'), 'utf8')
  /** Kunci JSON kolom: `JSON: "x"` dan `kf(<nama kolom>, "x")` (nama berkutip ganda atau backtick). */
  const kunci = new Set([
    ...[...GO.matchAll(/JSON: "([A-Za-z0-9]+)"/g)].map((x) => x[1]!),
    ...[...GO.matchAll(/kf\((?:"[^"]*"|`[^`]*`), "([A-Za-z0-9]+)"\)/g)].map((x) => x[1]!),
  ])
  it('instrumen: kunci yang diketahui terbaca (id, note, group, createOp)', () => {
    for (const k of ['id', 'note', 'group', 'createOp']) expect(kunci).toContain(k)
  })
  it('setiap kolom backend punya label', () => {
    expect([...kunci].filter((k) => !(k in LABEL_KOLOM_MASTER))).toEqual([])
  })
  it('setiap label (selain negara) dipakai kolom backend', () => {
    expect(Object.keys(LABEL_KOLOM_MASTER).filter((k) => k !== KUNCI_NEGARA && !kunci.has(k))).toEqual([])
  })
})
