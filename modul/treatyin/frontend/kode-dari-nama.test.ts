// Pencarian BALIK pengenal dari nama — tab Limits, dropdown `Treaty Type`.
//
// ---------------------------------------------------------------------
// ⛔ MENGAPA ADA, dan dari mana aturannya
// ---------------------------------------------------------------------
// Dokumen warisan menyimpan NAMA jenis treaty hampir selalu, pengenalnya
// hampir tidak pernah — terukur `TreatyTypeID` terisi pada 19 dari 1.360
// elemen `Limits[]`. Dropdown mengikat pengenal (`pyValue .ID` di
// `Section/LimitProportional.xml`), jadi tanpa pencarian ini 98,6% layer
// berbunyi `Choose` padahal namanya tertulis di kepala simpul di atasnya.
//
// ⭐ Polanya DISALIN dari ekspor, bukan dikarang:
// `Claim Non Prop/RDBList/GetReinsuranceTypeBYName_SQL.xml` mencari `ID`
// dari `NOTE` pada `REINSURANCETYPE`. Pega sendiri melakukannya.
//
// ⚠️ Yang TIDAK disalin saringan `type = '4'` rule itu: ia melayani
// spreading, sementara dropdown Limits menyaring `Flag = "active"` SAJA.
// Mencari di populasi yang lebih sempit daripada yang dropdown tawarkan akan
// gagal menemukan jenis yang dropdown-nya sendiri tampilkan.

import { describe, expect, it } from 'vitest'

import type { PilihanWarisan } from './api'
import { cariKodeDariNama } from './components/TabLimitsProp'

const daftar: PilihanWarisan[] = [
  { id: '10042', nama: 'QUOTA SHARE', kembar: false },
  { id: '10043', nama: 'SURPLUS', kembar: false },
  // Dua baris bernama sama — keduanya ditandai `kembar` oleh
  // `repository.tandaiKembar`.
  { id: '10050', nama: 'EXCESS OF LOSS', kembar: true },
  { id: '10051', nama: 'EXCESS OF LOSS', kembar: true },
]

describe('pengenal dari nama — hanya bila namanya TUNGGAL', () => {
  it('nama tunggal menemukan pengenalnya', () => {
    expect(cariKodeDariNama(daftar, 'QUOTA SHARE')).toBe('10042')
    expect(cariKodeDariNama(daftar, 'SURPLUS')).toBe('10043')
  })

  // ⛔ UJI PALING PENTING DI BERKAS INI.
  //
  // Nama kembar membuat pencarian memilih salah satu, dan layer akan menunjuk
  // jenis treaty yang KELIRU tanpa ada yang tahu. Tebakan yang terlihat benar
  // adalah tebakan yang paling mahal.
  it('⛔ nama KEMBAR tidak ditebak — dikembalikan null', () => {
    expect(cariKodeDariNama(daftar, 'EXCESS OF LOSS')).toBeNull()
  })

  it('nama yang tidak ada di daftar tidak dikarang', () => {
    expect(cariKodeDariNama(daftar, 'TIDAK ADA')).toBeNull()
  })

  it('nama kosong tidak mencari apa pun', () => {
    expect(cariKodeDariNama(daftar, '')).toBeNull()
    expect(cariKodeDariNama(daftar, '   ')).toBeNull()
  })

  // Nilai tersimpan membawa spasi di ujung pada sebagian dokumen; pembanding
  // yang tidak memangkasnya akan gagal pada baris yang sebenarnya cocok.
  it('spasi di ujung dipangkas sebelum dibandingkan', () => {
    expect(cariKodeDariNama(daftar, '  QUOTA SHARE  ')).toBe('10042')
  })

  it('daftar kosong tidak menemukan apa pun', () => {
    expect(cariKodeDariNama([], 'QUOTA SHARE')).toBeNull()
  })
})
