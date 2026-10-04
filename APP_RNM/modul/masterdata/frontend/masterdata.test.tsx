// Uji layar Master Data - data uji sintetis. Kamus label / rujukan diperiksa terhadap definisi master backend
// (`backend/models/daftar.go`) dua arah supaya kolom baru tidak tampil tanpa label.

import { readFileSync } from 'node:fs'
import { join } from 'node:path'
import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import type { MetaMaster } from './api'
import { jumlahHalaman } from './components/DaftarMaster'
import FormMaster, { badan, KUNCI_NEGARA, medanForm, panjangByte, periksa } from './components/FormMaster'
import { LABEL_KOLOM, RUJUKAN, TEKS } from './labels'

const DAFTAR_GO = readFileSync(join(__dirname, '..', 'backend', 'models', 'daftar.go'), 'utf8')

/** Awal deklarasi kolom jejak ubah (MD-7) - dibaca terpisah, tidak ikut blok master terakhir. */
const AWAL_AUDIT = DAFTAR_GO.indexOf('var KolomAudit')
const ISI_MASTER = AWAL_AUDIT >= 0 ? DAFTAR_GO.slice(0, AWAL_AUDIT) : DAFTAR_GO
/** Teks `var KolomAudit = []Kolom{ … }` (kosong bila tidak ada). */
const BLOK_AUDIT = AWAL_AUDIT >= 0 ? DAFTAR_GO.slice(AWAL_AUDIT, DAFTAR_GO.indexOf('\n}', AWAL_AUDIT) + 2) : ''

/** Blok tiap master di daftar.go (tanpa KolomAudit): kunci -> teks bloknya. */
function blokMaster(): Map<string, string> {
  const m = new Map<string, string>()
  const posisi = [...ISI_MASTER.matchAll(/Kunci: "([a-z]+)"/g)].map((x) => ({ kunci: x[1]!, i: x.index ?? 0 }))
  posisi.forEach((p, n) => m.set(p.kunci, ISI_MASTER.slice(p.i, posisi[n + 1]?.i ?? ISI_MASTER.length)))
  return m
}
/**
 * Kunci JSON kolom dalam satu blok (dua bentuk: `JSON: "x"` dan `kf(<nama kolom>, "x")` - nama kolom bisa berkutip
 * ganda atau backtick, mis. kf(`"GROUP"`, "group")).
 */
const kunciJSON = (blok: string) => [
  ...[...blok.matchAll(/JSON: "([A-Za-z0-9]+)"/g)].map((x) => x[1]!),
  ...[...blok.matchAll(/kf\((?:"[^"]*"|`[^`]*`), "([A-Za-z0-9]+)"\)/g)].map((x) => x[1]!),
]

describe('kamus label & rujukan = definisi master backend', () => {
  const blok = blokMaster()
  it('instrumen: delapan master terbaca, Nation memuat id / oldId / note / nationInitial', () => {
    expect([...blok.keys()]).toEqual(['nation', 'province', 'city', 'district', 'czone', 'accumulatedtype', 'accumulation', 'objectitemtype'])
    expect(kunciJSON(blok.get('nation')!).sort()).toEqual(['id', 'nationInitial', 'note', 'oldId'])
    expect(kunciJSON(blok.get('objectitemtype')!)).toContain('group')
    // Jejak ubah dibaca dari `var KolomAudit` sendiri, tidak bocor ke blok Object Item Type.
    expect(kunciJSON(BLOK_AUDIT)).toEqual(['createOp', 'tglCreate', 'updateOp', 'tglUpdate'])
    expect(kunciJSON(blok.get('objectitemtype')!)).not.toContain('createOp')
  })
  const semua = new Set([...[...blok.values()].flatMap(kunciJSON), ...kunciJSON(BLOK_AUDIT)])
  it('setiap kolom backend (termasuk jejak ubah) punya label', () => {
    expect([...semua].filter((k) => !(k in LABEL_KOLOM))).toEqual([])
  })
  it('setiap label (selain negara) dipakai kolom backend', () => {
    expect(Object.keys(LABEL_KOLOM).filter((k) => k !== KUNCI_NEGARA && !semua.has(k))).toEqual([])
  })
  it('setiap rujukan frontend = kolom master itu dan master tujuannya ada', () => {
    for (const [master, r] of Object.entries(RUJUKAN)) {
      for (const [kunci, { master: tujuan }] of Object.entries(r)) {
        expect(kunciJSON(blok.get(master)!)).toContain(kunci)
        expect(blok.has(tujuan)).toBe(true)
      }
    }
  })
})

const PROVINCE: MetaMaster = {
  kunci: 'province', judul: 'Province', idOtomatis: false,
  kolom: [
    { kunci: 'id', kolom: 'ID', lebar: 4000, wajib: true, turunan: false },
    { kunci: 'nationId', kolom: 'NATIONID', lebar: 4000, wajib: false, turunan: false },
    { kunci: 'note', kolom: 'NOTE', lebar: 5, wajib: true, turunan: false },
    { kunci: 'nationName', kolom: 'NATIONNAME', lebar: 4000, wajib: false, turunan: true },
  ],
}
const AKUMULASI: MetaMaster = {
  kunci: 'accumulation', judul: 'Accumulation', idOtomatis: true,
  kolom: [
    { kunci: 'id', kolom: 'ID', lebar: 4000, wajib: false, turunan: true },
    { kunci: 'note', kolom: 'NOTE', lebar: 4000, wajib: true, turunan: false },
    { kunci: 'zipCode', kolom: 'ZIPCODE', lebar: 4000, wajib: true, turunan: false },
  ],
}

describe('FormMaster', () => {
  it('Tambah: ID diisi pengguna; turunan tidak menjadi isian', () => {
    expect(medanForm(PROVINCE, true).map((k) => k.kunci)).toEqual(['id', 'nationId', 'note'])
  })
  it('Ubah: ID bukan isian (diambil dari rute, MD-9)', () => {
    expect(medanForm(PROVINCE, false).map((k) => k.kunci)).toEqual(['nationId', 'note'])
  })
  it('Accumulation Tambah: tanpa ID, meminta negara (MD-3)', () => {
    expect(medanForm(AKUMULASI, true).map((k) => k.kunci)).toEqual([KUNCI_NEGARA, 'note', 'zipCode'])
    expect(medanForm(AKUMULASI, false).map((k) => k.kunci)).toEqual(['note', 'zipCode'])
  })
  it('periksa: wajib dan lebar BYTE; badan dipangkas tanpa turunan', () => {
    expect(periksa(PROVINCE, true, { id: ' ', note: 'ABCDEF' })).toEqual({ id: TEKS.wajib, note: TEKS.terlaluPanjang(5) })
    expect(panjangByte('é')).toBe(2)
    expect(periksa(PROVINCE, true, { id: 'P1', note: 'éé' })).toEqual({})
    expect(badan(PROVINCE, true, { id: ' P1 ', nationId: 'N1', note: 'UJI', nationName: 'X' })).toEqual({ id: 'P1', nationId: 'N1', note: 'UJI' })
  })
  it('render Ubah: ID dan kolom turunan baca-saja; rujukan Nation = isian cari', () => {
    const html = renderToStaticMarkup(
      <FormMaster meta={PROVINCE} baris={{ aktif: true, id: 'P1', nationId: 'N1', note: 'UJI', nationName: 'NEGARA UJI' }} onTutup={() => {}} onTersimpan={() => {}} />,
    )
    expect(html).toContain(TEKS.judulUbah('Province'))
    expect(html).toMatch(/value="P1"[^>]*readOnly|readOnly[^>]*value="P1"/i)
    expect(html).toContain('NEGARA UJI')
    expect(html).toContain(`placeholder="${TEKS.cariRujukan}"`)
  })
  it('render Tambah Accumulation: ID otomatis + Negara', () => {
    const html = renderToStaticMarkup(<FormMaster meta={AKUMULASI} baris={null} onTutup={() => {}} onTersimpan={() => {}} />)
    expect(html).toContain(TEKS.idOtomatis)
    expect(html).toContain(LABEL_KOLOM.negara)
  })
})

describe('DaftarMaster', () => {
  it('jumlah halaman minimal 1', () => {
    expect(jumlahHalaman(0, 50)).toBe(1)
    expect(jumlahHalaman(101, 50)).toBe(3)
  })
})
