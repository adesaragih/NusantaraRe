// Copy Old - permintaan work owner 03-10-2026: "di samping tombol Add, tombol COPY OLD ... saat dibuka muncul popup,
// muncul semua list dari tabel lama yang belum dimigrasi, di setiap list bisa dicentang ... tombol PROCESS COPY untuk
// mengcopy yang dicentang lalu masuk ke table baru".

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

const DIALOG = readFileSync(join(__dirname, 'DialogCopyOld.tsx'), 'utf8')
const HALAMAN = readFileSync(join(__dirname, '..', 'pages', 'MasterProductNameLife.tsx'), 'utf8')

describe('tombol Copy Old dan popup-nya', () => {
  it('tombol Copy Old tepat di samping tombol Add', () => {
    const add = HALAMAN.indexOf('{GRID_MPNL.add}')
    const copyOld = HALAMAN.indexOf('{LAIN_MPNL.copyOld}')
    expect(add).toBeGreaterThan(-1)
    expect(copyOld).toBeGreaterThan(add)
    // Tidak ada elemen lain di antara keduanya selain penutup tombol Add.
    expect(HALAMAN.slice(add, copyOld)).not.toMatch(/<(span|div|Halaman)\b/)
    expect(HALAMAN).toContain('<DialogCopyOld')
  })

  it('popup memuat daftar produk lama, kotak centang per baris, dan Process Copy di kaki popup', () => {
    expect(DIALOG).toContain('ambilProdukLama()')
    expect(DIALOG).toContain('type="checkbox"')
    expect(DIALOG).toContain('disabled={!d.bolehDisalin}')
    // Tombol Process Copy di `aksi` Modal (kaki popup), mengirim hanya ID terpilih yang boleh disalin.
    const aksi = DIALOG.slice(DIALOG.indexOf('aksi={'), DIALOG.indexOf('>', DIALOG.indexOf('{LAIN_MPNL.prosesCopy}')))
    expect(aksi).toContain('{LAIN_MPNL.prosesCopy}')
    expect(DIALOG).toContain('salinProdukLama(idBolehDisalin(')
    // Sesudah Process Copy daftar dibaca ulang: yang tersalin hilang dari popup.
    expect(DIALOG.slice(DIALOG.indexOf('salinProdukLama(idBolehDisalin('))).toContain('muat()')
  })
})
