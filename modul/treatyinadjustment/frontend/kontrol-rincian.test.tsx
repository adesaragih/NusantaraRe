// ⛔ NOL `pxAutoComplete` DI LAYAR — semuanya dirender dropdown.
//
// Permintaan pemilik proses 8 Oktober 2026 datang dua kali: mula-mula *"di
// non prop maximum retention perbaiki, jangan auto complete, buat dropdown
// seperti yang lainnya"*, lalu — setelah panel itu saja yang diubah — *"cek
// di tab lainnya juga, masih banyak ternyata auto complete"*.
//
// ⚠️ 82 butir ber-`pxAutoComplete` di kerangka, dan 42 di antaranya SEL GRID.
// Percobaan pertama hanya menyentuh medan panel rincian, jadi sebagian besar
// tetap autocomplete. Uji ini karena itu menghitung dari kerangka dan
// merender KEDUANYA — medan dan sel.

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import type { SisiPenyesuaian } from './api'
import type { ButirKerangka } from './ekspor/jenis'
import { KERANGKA_RINCIAN, KERANGKA_TAB } from './ekspor/kerangka.gen'
import { formatTampil, RenderKerangka } from './komponen/KerangkaTab'

const KT = readFileSync(join(__dirname, 'komponen', 'KerangkaTab.tsx'), 'utf8')

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

const BUTIR: ButirKerangka[] = []
for (const isi of Object.values(KERANGKA_TAB)) semua(isi?.isi ?? [], BUTIR)
for (const isi of Object.values(KERANGKA_RINCIAN)) semua(isi ?? [], BUTIR)

const sisi: SisiPenyesuaian = {
  medan: { TreatyGroup: 'PROPERTY', Currency: 'IDR', Amount: '1000', Note: '' },
  larik: {},
}
const konteks = {
  sisi,
  akar: sisi,
  halaman: { ViewState: '0' },
  ubah: true,
  ubahMedan: () => undefined,
  ubahLarik: () => undefined,
  opsi: () => [{ value: 'PROPERTY', label: 'PROPERTY' }],
}

describe('`pxAutoComplete` dirender dropdown — di SELURUH kerangka', () => {
  it('⛔ `formatTampil` mengganti autocomplete, dan hanya itu', () => {
    expect(formatTampil('pxAutoComplete')).toBe('pxDropdown')
    for (const lain of ['pxDropdown', 'pxNumber', 'pxTextInput', 'pxTextArea', 'pxDateTime', 'pxCheckbox', '']) {
      expect(formatTampil(lain), lain).toBe(lain)
    }
  })

  // ⚠️ Dihitung dari kerangka: butir baru hasil pembangkitan ulang ikut
  // terjaga, dan menyusutnya cacah ini ketahuan.
  it('⛔ sel GRID ikut — di situlah sebagian besarnya berada', () => {
    const selOtomatis = BUTIR.filter((b) => b.t === 'grid' && b.format.includes('pxAutoComplete'))
    expect(selOtomatis.length).toBeGreaterThanOrEqual(20)
    // Dan penggantinya dipasang di jalur sel, bukan hanya di medan.
    expect(KT).toContain("const fmt = formatTampil(g.format[i] ?? '')")
    expect(KT).toContain('const fmt = formatTampil(m.format)')
  })

  // ⚠️ Kontrolnya combobox sejak 8 Oktober 2026 (*"saat mengetik seperti
  // mencari"*), jadi yang dicari BUKAN `<select>` lagi melainkan
  // `role="combobox"` — dan tetap BUKAN `<datalist>`, yang menerima ketikan
  // bebas sebagai nilai.
  it('⛔ panel Maximum Retention: combobox, bukan isian ber-`list`', () => {
    const html = renderToStaticMarkup(<RenderKerangka isi={KERANGKA_RINCIAN.MaxRetention ?? []} k={konteks} />)
    expect(html).toContain('role="combobox"')
    expect(html).not.toContain('<datalist')
    expect(html).not.toMatch(/<input[^>]*\slist=/)
  })

  it('⛔ panel lain yang dahulu autocomplete pun dropdown sekarang', () => {
    for (const nama of ['CoBList', 'DetailEGNPI', 'DetailLimitsOldData']) {
      const html = renderToStaticMarkup(<RenderKerangka isi={KERANGKA_RINCIAN[nama] ?? []} k={konteks} />)
      expect(html, nama).not.toContain('<datalist')
      expect(html, nama).not.toMatch(/<input[^>]*\slist=/)
    }
  })

  // ⭐ Sapuan terakhir: satu pun butir kerangka tidak boleh tersisa sebagai
  // autocomplete sesudah melewati `formatTampil`.
  it('⭐ nol butir tersisa sebagai autocomplete sesudah diterjemahkan', () => {
    const sisa: string[] = []
    for (const b of BUTIR) {
      if (b.t === 'medan' && formatTampil(b.format) === 'pxAutoComplete') sisa.push(`medan ${b.kunci}`)
      if (b.t === 'grid') {
        b.format.forEach((f, i) => {
          if (formatTampil(f) === 'pxAutoComplete') sisa.push(`sel ${b.kunci[i] ?? '?'}`)
        })
      }
    }
    expect(sisa).toEqual([])
  })
})
