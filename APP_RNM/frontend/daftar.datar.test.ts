import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import { AKAR_APLIKASI } from '../inti/frontend/uji/sumber'
import { FOLDER_KORPUS } from './katalogKorpus'
import { butirDatar, type ButirSidebar } from '../inti/frontend/lib/daftarMenu'
import { ENTRI_MENU } from './daftar'

// Menu DATAR - [keputusan work owner 30-09-2026] untuk Treaty Contract Out:
// navbar tidak bermodel kelompok-beranak "Treaty Contract Out ▸ Treaty
// Contract Out"; satu tombol langsung. Modul lain TIDAK berubah.
//
// Sejak menu dari tabel M_NAV_MENU (30-09-2026) butirnya dari `GET /api/menu`;
// penanda `datar` tetap milik rute frontend dan dibawa `susunMenu`
// (`inti/lib/daftarMenu.test.ts`).

/** Butir satu kelompok menurut rute frontend, dalam bentuk sidebar. */
const butir = (nama: string): ButirSidebar[] =>
  ENTRI_MENU.filter((e) => e.kelompok === nama).map((e) => ({
    halaman: e.modul,
    label: e.label,
    pemilik: e.pemilik,
    ...(e.datar === true ? { datar: true as const } : {}),
  }))

describe('menu datar', () => {
  it('Treaty Contract Out: satu butir, tampil datar', () => {
    const d = butirDatar(butir(FOLDER_KORPUS.treatyContractOut))
    expect(d?.halaman).toBe('tco-tahun')
  })
  it('modul lain tetap bertingkat - penanda tidak menular', () => {
    for (const nama of [FOLDER_KORPUS.claimLife, FOLDER_KORPUS.premiumListLife, FOLDER_KORPUS.komiteClaimLife]) {
      expect(butirDatar(butir(nama)), nama).toBeUndefined()
      expect(butir(nama).some((b) => 'datar' in b), nama).toBe(false)
    }
  })
  it('penanda diabaikan bila kelompoknya berbutir lebih dari satu', () => {
    const t = butir(FOLDER_KORPUS.treatyContractOut)[0]!
    expect(butirDatar([t, { ...t, halaman: 'tco-kontrak' }])).toBeUndefined()
  })
  it('Shell merender butir datar sebagai satu tombol, tanpa KelompokMenu', () => {
    const shell = readFileSync(join(AKAR_APLIKASI, 'inti', 'frontend', 'components', 'Shell.tsx'), 'utf8')
    expect(shell).toContain('const datar = butirDatar(k.butir)')
    const awal = shell.indexOf(': datar !== undefined ? (')
    const cabang = shell.slice(awal, shell.indexOf(') : (', awal))
    expect(cabang).toContain('pilih(datar.halaman)')
    expect(cabang).not.toContain('<KelompokMenu')
  })
})
