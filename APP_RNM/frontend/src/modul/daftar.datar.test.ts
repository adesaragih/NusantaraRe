import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import { MODUL } from '../inti/labels'
import { butirDatar, butirKelompok } from '../inti/lib/daftarMenu'
import { ENTRI_MENU } from './daftar'

// Menu DATAR - [keputusan work owner 30-09-2026] untuk Treaty Contract Out:
// navbar tidak bermodel kelompok-beranak "Treaty Contract Out ▸ Treaty
// Contract Out"; satu tombol langsung. Modul lain TIDAK berubah.

const butir = (nama: string) => butirKelompok(ENTRI_MENU, nama)

describe('menu datar', () => {
  it('Treaty Contract Out: satu butir, tampil datar', () => {
    const d = butirDatar(butir(MODUL.treatyContractOut))
    expect(d?.halaman).toBe('tco-tahun')
  })
  it('modul lain tetap bertingkat - penanda tidak menular', () => {
    for (const nama of [MODUL.claimLife, MODUL.premiumListLife, MODUL.komiteClaimLife]) {
      expect(butirDatar(butir(nama)), nama).toBeUndefined()
      expect(butir(nama).some((b) => 'datar' in b), nama).toBe(false)
    }
  })
  it('penanda diabaikan bila kelompoknya berbutir lebih dari satu', () => {
    const t = butir(MODUL.treatyContractOut)[0]!
    expect(butirDatar([t, { ...t, halaman: 'tco-kontrak' }])).toBeUndefined()
  })
  it('Shell merender butir datar sebagai satu tombol, tanpa KelompokMenu', () => {
    const shell = readFileSync(join(__dirname, '..', 'inti', 'components', 'Shell.tsx'), 'utf8')
    expect(shell).toContain('const datar = butirDatar(butir)')
    const cabang = shell.slice(shell.indexOf('{datar !== undefined ? ('), shell.indexOf(': butir.length === 0 ? ('))
    expect(cabang).toContain('pilih(datar.halaman)')
    expect(cabang).not.toContain('<KelompokMenu')
  })
})
