// Asal: salinan `modul/nbtreatyin/frontend/asalKasus.test.ts` (06-10-2026), disesuaikan rute EDM Treaty In.
// Tombol Back kembali ke tempat berkas dibuka (Beranda atau portal).

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import { tujuanKembali } from './asalKasus'

describe('tujuanKembali', () => {
  it('Back berkas dari Beranda = Beranda; dari portal = portal', () => {
    expect(tujuanKembali('beranda', undefined, true)).toBe('beranda')
    expect(tujuanKembali('portal', undefined, true)).toBe('portal')
  })

  it('sesudah Submit (ada pesan) tetap ke portal, tempat pesannya tampil', () => {
    expect(tujuanKembali('beranda', 'UJI terkirim', true)).toBe('portal')
  })

  it('tanpa jalan ke Beranda = portal', () => {
    expect(tujuanKembali('beranda', undefined, false)).toBe('portal')
  })
})

describe('rute EDM Treaty In mencatat asal berkas', () => {
  const rute = readFileSync(join(__dirname, 'rute.tsx'), 'utf8')

  it('dibuka dari Beranda = asal beranda; dari portal atau menu = asal portal', () => {
    const efekBeranda = rute.slice(rute.indexOf('if (ketukBuka !== undefined'), rute.indexOf('}, [ketukBuka, idBuka])'))
    expect(efekBeranda).toContain("setAsal('beranda')")
    const bukaPortal = rute.slice(rute.indexOf('<PortalEDMTreatyIn'), rute.indexOf('<LayarKasus'))
    expect(bukaPortal).toContain("setAsal('portal')")
    const efekMenu = rute.slice(rute.indexOf('const ketukLalu'), rute.indexOf('}, [ketukMenu])'))
    expect(efekMenu).toContain("setAsal('portal')")
  })

  it('Back memanggil onBeranda bila tujuannya Beranda', () => {
    expect(rute).toContain('tujuanKembali(asal, p, onBeranda !== undefined)')
    expect(rute).toContain('onBeranda?.()')
  })
})
