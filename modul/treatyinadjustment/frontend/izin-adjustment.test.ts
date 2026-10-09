// ⛔ APA YANG BOLEH DIISI SAAT ADJUSTMENT — dikunci terhadap EKSPOR, bukan
// terhadap ingatan.
//
// Permintaan pemilik proses 8 Oktober 2026: *"dipastikan juga untuk yang di
// masing-masing tab baca XML-nya dan pastikan sesuai dengan Adjustment Type
// dan Material Type … pastikan apa yang seharusnya bisa diinput saat adjust,
// apa yang tidak"*.
//
// Pola yang dibaca dari `Treaty In Adjustment/Section/*.xml` dan berlaku di
// SELURUH tab — 56 butir, dihitung di bawah, bukan ditulis tangan:
//
//   Material (1)      isi kontrak TERBUKA (Limits, Share, EGNPI, Reporting,
//                     Installment, Accumulation, Portfolio, Event Limits),
//                     tetapi Exclusions dan Special Conditions TERKUNCI.
//   Non Material (2)  kebalikannya: isi kontrak terkunci dan tombol
//                     Add/Delete/Update mati; Exclusions dan Special
//                     Conditions terbuka.
//
// ⭐ Pembalikan itu masuk akal dan karenanya layak dikunci: revisi material
// mengubah syarat material, revisi non-material mengubah yang bukan.
//
// ⚠️ Uji ini MENGHITUNG butirnya dari kerangka. Butir baru hasil pembangkitan
// ulang ikut teruji tanpa disentuh, dan hilangnya butir ketahuan dari
// cacahnya.

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import type { ButirKerangka } from './ekspor/jenis'
import { KERANGKA_RINCIAN, KERANGKA_TAB } from './ekspor/kerangka.gen'
import { OPSI_JENIS_REVISI, OPSI_MATERIAL } from './labelsPenyesuaian'
import { terkunci, tombolMati } from './komponen/aksiTombol'
import { syaratTerpenuhi } from './ekspor/syarat'

const PICKER = readFileSync(join(__dirname, 'komponen', 'PilihMaster.tsx'), 'utf8')

type Halaman = Record<string, string>

/** Semua butir sebuah pohon, serata apa pun sarangnya. */
function semua(x: unknown, keluar: ButirKerangka[] = []): ButirKerangka[] {
  if (Array.isArray(x)) {
    for (const e of x) semua(e, keluar)
    return keluar
  }
  if (x === null || typeof x !== 'object') return keluar
  const o = x as Record<string, unknown>
  if (typeof o.t === 'string') keluar.push(o as unknown as ButirKerangka)
  for (const v of Object.values(o)) semua(v, keluar)
  return keluar
}

const BUTIR: Array<{ tab: string; b: ButirKerangka }> = []
for (const [tab, isi] of Object.entries(KERANGKA_TAB)) {
  for (const b of semua(isi?.isi ?? [])) BUTIR.push({ tab, b })
}
for (const [tab, isi] of Object.entries(KERANGKA_RINCIAN)) {
  for (const b of semua(isi ?? [])) BUTIR.push({ tab: `rincian:${tab}`, b })
}

const MATERIAL = /EDMMaterialType/
/** Mode SUNTING — supaya yang diuji benar-benar Material Type, bukan ViewState. */
const sunting = (lebih: Halaman): Halaman => ({ ViewState: '0', IsEditData: '0', ...lebih })

describe('Material Type menentukan apa yang boleh diisi — diadu dengan ekspor', () => {
  /** Medan yang terkunci oleh `… || EDMMaterialType = 2` (isi kontrak). */
  const medanNonMaterial = BUTIR.filter(
    ({ b }) => b.t === 'medan' && Array.isArray(b.baca) && b.baca.some((s) => MATERIAL.test(s) && /EDMMaterialType\s*=\s*2/.test(s)),
  )
  /** Medan yang terkunci oleh `… || EDMMaterialType = 1` (Exclusions, Special Conditions). */
  const medanMaterial = BUTIR.filter(
    ({ b }) => b.t === 'medan' && Array.isArray(b.baca) && b.baca.some((s) => MATERIAL.test(s) && /EDMMaterialType\s*=\s*1/.test(s)),
  )
  const tombolNonMaterial = BUTIR.filter(
    ({ b }) => b.t === 'tombol' && (b.nonaktif ?? []).some((s) => /EDMMaterialType\s*=\s*2/.test(s)),
  )

  it('⛔ cacahnya tidak menyusut diam-diam', () => {
    // Berubah bila kerangka dibangkitkan ulang — periksa SEBABNYA sebelum
    // memperbarui angkanya. Menyusut = medan kehilangan penjaganya.
    expect(medanNonMaterial.length).toBeGreaterThanOrEqual(21)
    expect(medanMaterial.length).toBeGreaterThanOrEqual(4)
    expect(tombolNonMaterial.length).toBeGreaterThanOrEqual(19)
  })

  it('⭐ Material (1): isi kontrak TERBUKA — inilah yang boleh diisi saat adjust', () => {
    for (const { tab, b } of medanNonMaterial) {
      if (b.t !== 'medan') continue
      expect(terkunci(b.baca, sunting({ EDMMaterialType: '1' })), `${tab} · ${b.label || b.kunci}`).toBe(false)
    }
  })

  it('⛔ Non Material (2): isi kontrak TERKUNCI', () => {
    for (const { tab, b } of medanNonMaterial) {
      if (b.t !== 'medan') continue
      expect(terkunci(b.baca, sunting({ EDMMaterialType: '2' })), `${tab} · ${b.label || b.kunci}`).toBe(true)
    }
  })

  it('⭐ PEMBALIKANNYA: Exclusions / Special Conditions terbuka pada Non Material, terkunci pada Material', () => {
    for (const { tab, b } of medanMaterial) {
      if (b.t !== 'medan') continue
      const nama = `${tab} · ${b.label || b.kunci}`
      expect(terkunci(b.baca, sunting({ EDMMaterialType: '1' })), nama).toBe(true)
      expect(terkunci(b.baca, sunting({ EDMMaterialType: '2' })), nama).toBe(false)
    }
    // Dan memang hanya kedua tab itu yang berperilaku terbalik.
    const tabs = new Set(medanMaterial.map(({ tab }) => tab.split('#')[1] ?? tab))
    expect([...tabs].sort()).toEqual(['Exclusions', 'Special Conditions'])
  })

  it('⛔ tombol Add / Delete / Update mati pada Non Material, hidup pada Material', () => {
    for (const { tab, b } of tombolNonMaterial) {
      if (b.t !== 'tombol') continue
      const nama = `${tab} · ${b.label}`
      expect(tombolMati(b, sunting({ EDMMaterialType: '2' })), nama).toBe(true)
      expect(tombolMati(b, sunting({ EDMMaterialType: '1' })), nama).toBe(false)
    }
  })

  // ⚠️ `IsEditData = 1` mengunci TERPISAH dari Material Type: ia keadaan
  // layar (data sudah dikunci), bukan jenis revisi. Keduanya bergabung
  // dengan `||`, jadi yang satu tidak boleh menutupi yang lain.
  it('⛔ `IsEditData = 1` tetap mengunci walau Material Type membuka', () => {
    const dua = BUTIR.filter(
      ({ b }) => b.t === 'medan' && Array.isArray(b.baca) && b.baca.some((s) => /IsEditData/.test(s) && MATERIAL.test(s)),
    )
    expect(dua.length).toBeGreaterThan(0)
    for (const { tab, b } of dua) {
      if (b.t !== 'medan' || !Array.isArray(b.baca)) continue
      const buka = b.baca.some((s: string) => /EDMMaterialType\s*=\s*1/.test(s)) ? '2' : '1'
      expect(terkunci(b.baca, sunting({ EDMMaterialType: buka, IsEditData: '1' })), `${tab} · ${b.label || b.kunci}`).toBe(true)
    }
  })
})

describe('Adjustment Type (EDMState) menentukan tab mana yang ada', () => {
  // ⚠️ Syaratnya melekat pada TAB-nya, bukan pada butir di dalamnya —
  // seluruh tab muncul atau tidak sama sekali.
  it('⭐ Value Difference hanya pada revisi NON-premi yang Material', () => {
    const bersyarat = Object.entries(KERANGKA_TAB).filter(
      ([, isi]) => (isi?.syarat ?? []).some((s) => /EDMState/.test(s) && /EDMMaterialType/.test(s)),
    )
    expect(bersyarat.map(([nama]) => nama.split('#')[1])).toEqual(['Value Difference'])
    for (const [nama, isi] of bersyarat) {
      const syarat = isi?.syarat ?? []
      // Addendum Premi (EDMState 3) tidak punya tab ini, dan revisi
      // Non Material pun tidak.
      expect(syaratTerpenuhi(syarat, { EDMState: '1', EDMMaterialType: '1' }), nama).toBe(true)
      expect(syaratTerpenuhi(syarat, { EDMState: '3', EDMMaterialType: '1' }), nama).toBe(false)
      expect(syaratTerpenuhi(syarat, { EDMState: '1', EDMMaterialType: '2' }), nama).toBe(false)
    }
  })

  // ⛔ Kode yang dikirim picker HARUS kode ekspor, bukan urutan radio —
  // seluruh matriks di atas dinilai atas kode ini. Tertukar sedikit saja dan
  // tab yang salah terbuka tanpa galat apa pun.
  it('⛔ picker mengirim kode ekspor: Internal 1 · External 2 · Addendum Premi 3', () => {
    expect(OPSI_JENIS_REVISI.map((o) => o.nilai)).toEqual(['1', '2'])
    expect([...OPSI_MATERIAL]).toEqual(['1', '2'])
    // Addendum Premi memintas kedua radio: `EDMState 3`, dan Material `1`
    // karena tab Value Difference-nya memang tidak berlaku di sana.
    expect(PICKER).toContain("internalType: jenis === 'premi' ? '3' : edmState")
    expect(PICKER).toContain("materialType: jenis === 'premi' ? '1' : material")
  })
})
