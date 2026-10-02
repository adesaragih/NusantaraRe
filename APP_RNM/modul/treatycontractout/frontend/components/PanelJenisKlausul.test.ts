// Uji panel jenis klausul — tiket 08 Treaty Contract Out.

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import type { AturanKlausul, JenisKlausul, Klausul } from '../api'
import {
  aturanAnak,
  aturanInduk,
  barisSubjenis,
  formKlausulDari,
  formKlausulKosong,
  jenisBerkurs,
  tampilMedanKlausul,
  keMasukKlausul,
  labelMedan,
  pemilihReinsType,
  addTampil,
  tabSubjenis,
  rencanaKonversi,
} from './PanelJenisKlausul'

const KODE = readFileSync(join(__dirname, 'PanelJenisKlausul.tsx'), 'utf8')
  .split('\n')
  .filter((b) => !b.trimStart().startsWith('//') && !b.trimStart().startsWith('*'))
  .join('\n')

function aturan(p: Partial<AturanKlausul>): AturanKlausul {
  return { jenis: 'EPI', anak: false, subjenis: '', medan: ['ReinsTypeID', 'Line', 'Rp', 'Usd'],
    wajib: ['ReinsTypeID', 'Rp', 'Usd'], turunan: null, ditahan: '', berkurs: false, konversi: '', sumber: 'SaveTreatyArrEPI_Act', ...p }
}

function klausul(p: Partial<Klausul>): Klausul {
  return { id: '10000001', treatyYear: '2026', treatyYearId: '1000001', treatyGroupId: '10001', treatyDescId: '10009',
    treatyDescName: 'UJI EPI', reinsTypeId: '10003', reinsTypeName: 'UJI QS', parentReinsTypeId: '00', subjenis: '',
    medan: {}, kurs: '', userId: 'UJI-ADMIN', tglUpdate: '', ...p }
}

describe('label medan VERBATIM per jenis', () => {
  it('bawaan, lalu penimpaan per jenis dan per subjenis', () => {
    expect(labelMedan(aturan({}), 'ReinsTypeID')).toBe('ReinsType')
    expect(labelMedan(aturan({ jenis: 'MinLOL' }), 'Pct')).toBe('Minimum LOL (%)')
    expect(labelMedan(aturan({ jenis: 'MinLOLMB' }), 'Pct')).toBe('Minimum LOL MB (%)')
    expect(labelMedan(aturan({ jenis: 'ExclutionTreaty', subjenis: 'Occupation' }), 'Usd')).toBe('TSI Less Than (USD)')
    expect(labelMedan(aturan({ jenis: 'ExclutionTreaty', subjenis: 'Object' }), 'Pct')).toBe('TSI BI >')
    expect(labelMedan(aturan({ jenis: 'CoinsPanel' }), 'CoIns_Min')).toBe('From')
  })
})

describe('form klausul', () => {
  it('kosong dan dari baris — hanya medan aturan', () => {
    const a = aturan({})
    expect(formKlausulKosong(a)).toEqual({ id: '', medan: { ReinsTypeID: '', Line: '', Rp: '', Usd: '' } })
    const f = formKlausulDari(a, klausul({ medan: { ReinsTypeID: '10003', Rp: '1000.5', Usd: '70', Kurs: 'X' } }))
    expect(f).toEqual({ id: '10000001', medan: { ReinsTypeID: '10003', Line: '', Rp: '1000.5', Usd: '70' } })
  })
  it('anak: Rp/Usd TURUNAN tidak dikirim; induk disebut', () => {
    const a = aturan({ jenis: 'EpiList', anak: true, medan: ['ReinsTypeID', 'Pct'], turunan: ['Rp', 'Usd'] })
    const m = keMasukKlausul(a, '10009', { id: '', medan: { ReinsTypeID: '10005', Pct: ' 25 ', Rp: '9' } }, '10003')
    expect(m).toEqual({ id: '', descId: '10009', anak: true, subjenis: '', parentReinsTypeId: '10003',
      medan: { ReinsTypeID: '10005', Pct: '25' } })
  })
  it('induk tidak menyebut induk', () => {
    expect(keMasukKlausul(aturan({}), '10009', formKlausulKosong(aturan({})), '00').parentReinsTypeId).toBe('')
  })
  it('baris exclusion disaring per subjenis', () => {
    const d = [klausul({ id: '1', subjenis: 'Occupation' }), klausul({ id: '2', subjenis: 'Clause' })]
    expect(barisSubjenis(d, 'Clause').map((k) => k.id)).toEqual(['2'])
    expect(barisSubjenis(d, '').length).toBe(2)
  })
  it('aturan induk/anak dari jenis', () => {
    const j: JenisKlausul = { id: '10009', descName: 'UJI EPI', isXol: '0', statusAktif: '1', catatan: '',
      aturan: [aturan({}), aturan({ jenis: 'EpiList', anak: true })] }
    expect(aturanInduk(j).map((a) => a.jenis)).toEqual(['EPI'])
    expect(aturanAnak(j)?.jenis).toBe('EpiList')
  })
})

describe('kabel', () => {
  it('Cancel membuang isian panel INI saja (AC 29) — state form per GridAturan', () => {
    expect(KODE).toMatch(/function GridAturan[\s\S]*useState<FormKlausul \| null>/)
    expect(KODE).toContain('KLAUSUL_TCO.cancel')
  })
  it('jenis ditahan tampil dengan alasannya, tanpa form', () => {
    expect(KODE).toContain('aturan.ditahan !== ')
  })
  it('uang/persen tidak menjadi angka JavaScript', () => {
    // `(?<![A-Za-z])`: `formatNumber(` (pemformat TEKS bersama, bekerja pada
    // digit tanpa float - `inti/frontend/lib/format.ts`) bukan `Number(` JavaScript.
    expect(KODE).not.toMatch(/(?<![A-Za-z])Number\(|parseFloat|toFixed/)
    // Sasaran impornya, bukan kedalaman `../`: letak folder berganti (struktur
    // tim satu folder per modul), pemformat bersamanya tidak.
    expect(KODE).toMatch(/from '(\.\.\/)+inti\/frontend\/lib\/format'/)
  })
  it('Show Child dan Close Child dari label', () => {
    expect(KODE).toContain('KLAUSUL_TCO.showChild')
    expect(KODE).toContain('KLAUSUL_TCO.closeChild')
  })
})

describe('kurs (tiket 11)', () => {
  it('rencana konversi mengikuti HitungRpUsd_depan dan CalculateTSIExcludeTreaty', () => {
    expect(rencanaKonversi('RpKeUsd', 'Rp')).toEqual({ dari: 'Rp', ke: 'Usd', skala: '8' })
    expect(rencanaKonversi('RpKeUsd', 'Usd')).toBeNull()
    expect(rencanaKonversi('DuaArah', 'Rp')).toEqual({ dari: 'Rp', ke: 'Usd', skala: '4' })
    expect(rencanaKonversi('DuaArah', 'Usd')).toEqual({ dari: 'Usd', ke: 'Rp', skala: '8' })
    expect(rencanaKonversi('', 'Rp')).toBeNull()
  })
  it('jenis berkurs bila salah satu aturannya berkurs', () => {
    const j: JenisKlausul = { id: '10009', descName: 'UJI EPI', isXol: '0', statusAktif: '1', catatan: '',
      aturan: [aturan({ berkurs: true, konversi: 'RpKeUsd', turunan: ['Usd'] }), aturan({ jenis: 'EpiList', anak: true, berkurs: true })] }
    expect(jenisBerkurs(j)).toBe(true)
    expect(jenisBerkurs({ ...j, aturan: [aturan({})] })).toBe(false)
  })
  it('induk berkurs: Usd TURUNAN tidak dikirim', () => {
    const a = aturan({ berkurs: true, konversi: 'RpKeUsd', turunan: ['Usd'] })
    const m = keMasukKlausul(a, '10009', { id: '', medan: { ReinsTypeID: '10003', Line: '', Rp: '1000', Usd: '0.06451613' } }, '00')
    expect(Object.keys(m.medan)).not.toContain('Usd')
  })
  it('kurs dibaca dari server; tanpa kurs Add nonaktif; konversi di server', () => {
    expect(KODE).toContain('ambilKursTahun(tahunID)')
    expect(KODE).toContain('disabled={aturan.berkurs && !kursAda}')
    expect(KODE).toContain('konversiKurs(tahunID, r.dari, nilai, r.skala)')
  })
})

describe('nol catatan kurs di layar (keputusan work owner 30-09-2026)', () => {
  it('catatan baris kembar tidak dirender', () => {
    expect(KODE).not.toMatch(/catatanMasterKurs|catatanKembar|barisMasterKembar/)
  })
})

describe('tampilan desimal berpemisah ribuan (keputusan work owner 30-09-2026)', () => {
  it('medan desimal: titik ribuan, koma desimal, TANPA digit dibuang', () => {
    expect(tampilMedanKlausul('Rp', '1000000')).toBe('1.000.000')
    expect(tampilMedanKlausul('Usd', '1234567.12345678')).toBe('1.234.567,12345678')
    expect(tampilMedanKlausul('TreatyLimit', '250000000.50')).toBe('250.000.000,5')
    expect(tampilMedanKlausul('Pct', '33.3333333333')).toBe('33,3333333333')
    expect(tampilMedanKlausul('CoIns_Min', '')).toBe('')
  })
  it('medan bukan desimal apa adanya — Line, Layer, kode, teks', () => {
    expect(tampilMedanKlausul('Line', '1000')).toBe('1000')
    expect(tampilMedanKlausul('Layer', '2')).toBe('2')
    expect(tampilMedanKlausul('ID_Occupation', '100123')).toBe('100123')
  })
  it('grid, Total Pct, dan kurs memakai pemformat yang sama; catatan pengembang tidak tampil', () => {
    expect(KODE).toContain("tampilMedanKlausul(m, k.medan[m] ?? '')")
    expect(KODE).toContain('formatNumber(daftar.totalPct, DESIMAL_TAK_DIBATASI)')
    expect(KODE).toContain('formatNumber(kurs.kurs, DESIMAL_TAK_DIBATASI)')
    expect(KODE).not.toMatch(/turunanServer|catatanRpKeUsd|catatanDuaArah/)
  })
})

// ReinsType Treaty Limit ikut XML: `pxAutoComplete` di grid induk
// (`GridTreatyArrangementTreatyLimit.xml` b3025, RD induk) dan anak
// (`GridTreatyArrTreatyLimitList.xml` b2892, jenis porsi, tanpa induk
// [keputusan work owner 30-09-2026]). Anak SETIAP jenis memakai pilihan yang
// sama [keputusan work owner 02-10-2026]; induk jenis lain tetap dropdown.
describe('pemilih ReinsType per aturan', () => {
  it('Treaty Limit induk: dapat difilter, daftar induk', () => {
    expect(pemilihReinsType(aturan({ jenis: 'TreatyLimit' }))).toBe('saring-induk')
  })
  it('anak Treaty Limit: dapat difilter, daftar porsi - dari penanda ATURAN', () => {
    expect(pemilihReinsType(aturan({ jenis: 'TreatyLimitChild', anak: true, pilihanReins: 'anak-treaty-limit' }))).toBe('saring-anak')
    // Penandanya yang menentukan, bukan nama jenis.
    expect(pemilihReinsType(aturan({ jenis: 'TreatyLimitChild', anak: true }))).toBe('dropdown')
  })
  it('induk jenis lain tetap dropdown daftar induk', () => {
    for (const jenis of ['EPI', 'PLA', 'CashLossLimit', 'FacIn', 'ExGratia', 'ClaimCoorp', 'Ricomm']) {
      expect(pemilihReinsType(aturan({ jenis })), jenis).toBe('dropdown')
    }
  })
  it('anak SETIAP jenis = anak Treaty Limit: dapat difilter, jenis porsi saja', () => {
    for (const jenis of ['PLAList', 'CashLossLimitList', 'FacInList', 'ExGratiaChildList', 'EpiList', 'ClaimCoorpChild']) {
      expect(pemilihReinsType(aturan({ jenis, anak: true, pilihanReins: 'anak-treaty-limit' })), jenis).toBe('saring-anak')
    }
  })
  it('pemilih anak tidak menerima induk - daftar porsi saja [keputusan work owner 02-10-2026]', () => {
    expect(KODE).toContain("anak={pemilih === 'saring-anak'}")
    expect(KODE).not.toContain('anakTreatyLimitDari')
  })
})

// Minimum LOL, Max Coins Panel, Minimum LOL MB: satu baris per tahun - `Add` hanya bila `ID == ''`
// (`GridTreatyArrangementMinLOL.xml` b2232 …) [keputusan work owner 02-10-2026].
describe('Add jenis satu baris', () => {
  it('hilang begitu jenis berisi satu baris, dan selama daftar belum dimuat', () => {
    const satu = aturan({ jenis: 'MinLOL', satuBaris: true })
    expect(addTampil(satu, true, 0)).toBe(true)
    expect(addTampil(satu, true, 1)).toBe(false)
    expect(addTampil(satu, false, 0)).toBe(false)
  })
  it('jenis lain tidak dibatasi', () => {
    expect(addTampil(aturan({ jenis: 'EPI' }), true, 3)).toBe(true)
  })
  it('tombol Add dirender lewat penjaga itu', () => {
    expect(KODE).toContain('{addTampil(aturan, daftar !== null, baris.length) && (')
  })
})

// 10013 Exclusion Treaty [keputusan work owner 02-10-2026]: subjenis sebagai tab; ID Occupation /
// ID Clause satu dropdown yang dapat difilter, tanpa kotak Search terpisah.
describe('10013 Exclusion Treaty', () => {
  const exclusion: JenisKlausul = {
    id: '10013', descName: 'UJI EXCLUSION', isXol: '0', statusAktif: '', catatan: '',
    aturan: ['Occupation', 'Clause', 'Object', 'Periode'].map((subjenis) => aturan({ jenis: 'ExclutionTreaty', subjenis })),
  }
  it('subjenis tampil sebagai tab, urut aturan', () => {
    expect(tabSubjenis(exclusion)).toEqual(['Occupation', 'Clause', 'Object', 'Periode'])
  })
  it('jenis berinduk satu tanpa tab', () => {
    const epi: JenisKlausul = { ...exclusion, id: '10009', aturan: [aturan({}), aturan({ jenis: 'EpiList', anak: true })] }
    expect(tabSubjenis(epi)).toEqual([])
  })
  it('hanya grid tab aktif yang dirender, lewat StripTab inti', () => {
    expect(KODE).toContain('<StripTab tab={tab} aktif={tabAktif} onPilih={setTabAktif} />')
    expect(KODE).toContain(".filter((a) => tab.length === 0 || a.subjenis === tabAktif)")
  })
  it('ID Occupation / ID Clause: PilihMasterKlausul, nol kotak Search', () => {
    expect(KODE).toContain('<PilihMasterKlausul')
    expect(KODE).not.toMatch(/cariPilihan\b|setCari|KLAUSUL_TCO\.cariPilihan/)
  })
})
