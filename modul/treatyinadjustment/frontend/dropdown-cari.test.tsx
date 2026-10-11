// ⛔ SELURUH DROPDOWN DAPAT DIKETIK UNTUK MENCARI.
//
// Permintaan pemilik proses 8 Oktober 2026 untuk Treaty In dan Adjustment:
// *"kondisi saat melakukan pengetikannya seharusnya terlihat layaknya
// melakukan mencari, kemudian data yang keluar adalah yang 100% mirip dengan
// yang diketik"*.
//
// ⚠️ BUKAN kembali ke `<datalist>` yang dibuang sehari sebelumnya: di sana
// ketikan apa pun lolos menjadi nilai. Combobox menyaring daftarnya, tetapi
// nilai HANYA berubah saat satu butir dipilih.

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import { cocokSaring } from '../../../inti/frontend/components/ui/pilihSaring'
import type { SisiPenyesuaian } from './api'
import { KERANGKA_RINCIAN } from './ekspor/kerangka.gen'
import { RenderKerangka } from './komponen/KerangkaTab'

const KT = readFileSync(join(__dirname, 'komponen', 'KerangkaTab.tsx'), 'utf8')

const sisi: SisiPenyesuaian = { medan: { TreatyGroup: 'PROPERTY', Currency: 'IDR' }, larik: {} }
const html = renderToStaticMarkup(
  <RenderKerangka
    isi={KERANGKA_RINCIAN.MaxRetention ?? []}
    k={{
      sisi,
      akar: sisi,
      halaman: { ViewState: '0' },
      ubah: true,
      ubahMedan: () => undefined,
      ubahLarik: () => undefined,
      opsi: () => [{ value: 'PROPERTY', label: 'PROPERTY' }],
    }}
  />,
)

describe('dropdown = combobox bercari', () => {
  it('⛔ kotaknya menerima ketikan dan menyaring, bukan `<select>` kaku', () => {
    expect(html).toContain('role="combobox"')
    expect(html).toContain('aria-autocomplete="list"')
    // ⚠️ Petunjuk kotaknya datang dari `TEKS_UI` dan BERGANTUNG BAHASA;
    // uji ini merender tanpa `BahasaUI.Provider`, jadi yang diperiksa
    // adalah ADANYA petunjuk, bukan bunyinya.
    expect(html).toMatch(/placeholder="[^"]+"/)
    expect(html).not.toContain('<select')
  })

  it('⛔ BUKAN `<datalist>` — ketikan bebas tidak boleh lolos jadi nilai', () => {
    expect(html).not.toContain('<datalist')
    expect(html).not.toMatch(/<input[^>]*\slist=/)
  })

  // ⭐ "100% mirip dengan yang diketik": pencocokan BAGIAN TEKS, tanpa huruf
  // besar-kecil — pemakai mengetik potongan yang diingatnya, bukan selalu
  // huruf pertama.
  it('⛔ saringannya mencocoki bagian teks, bukan hanya awalan', () => {
    const o = { value: '1', label: '2025 QS 181M TRT', keterangan: '10500' }
    expect(cocokSaring(o, '')).toBe(true)
    expect(cocokSaring(o, '2025')).toBe(true)
    expect(cocokSaring(o, '181')).toBe(true)
    expect(cocokSaring(o, 'qs')).toBe(true)
    expect(cocokSaring(o, 'TRT')).toBe(true)
    // Kode/ID ikut dicari — di sanalah nomor kontrak berada.
    expect(cocokSaring(o, '10500')).toBe(true)
    // Yang TIDAK ada memang tidak keluar.
    expect(cocokSaring(o, 'SURPLUS')).toBe(false)
    expect(cocokSaring(o, '999')).toBe(false)
  })

  it('⭐ saringan tanpa jeda — daftarnya sudah di klien', () => {
    expect(KT).toContain('jedaMs={0}')
  })

  it('⛔ sel grid: label disembunyikan, namanya pindah ke `aria-label`', () => {
    expect(KT).toContain('sembunyikanLabel')
    // Medan panel tetap berlabel; yang disembunyikan hanya sel.
    expect(html).toContain('field__label')
  })
})
