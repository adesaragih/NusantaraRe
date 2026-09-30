import { existsSync, readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import { AKAR_APLIKASI, folderKorpusBelumDimigrasi } from '../inti/frontend/uji/sumber'
import { menuTabelDariMigrasi } from '../inti/frontend/uji/menuBersih'
import { MENU } from '../inti/frontend/labels'
import { susunMenu } from '../inti/frontend/lib/daftarMenu'
import { FOLDER_KORPUS } from './katalogKorpus'
import { ENTRI_MENU, MODUL_FRONTEND } from './daftar'

// Menu DATAR - [keputusan work owner 30-09-2026] (`PROMPT-MENU-DATAR-PER-GROUPMENU.md`):
// "menu jangan ada model seperti child ... 1 modul 1 menu". Di bawah kepala
// GROUPMENU: satu tombol per modul, tanpa kelompok yang dilipat, tanpa anak,
// tanpa panah buka-tutup. Klik tombol modul membuka halaman awalnya.
//
// Menggantikan uji penanda `datar` Treaty Contract Out (`butirDatar`): yang
// dulu satu pengecualian kini aturan untuk semua modul.

const SHELL = readFileSync(join(AKAR_APLIKASI, 'inti', 'frontend', 'components', 'Shell.tsx'), 'utf8')
const NAV = SHELL.slice(SHELL.indexOf('<nav className="shell__nav"'), SHELL.indexOf('</nav>'))
const SEMUA = susunMenu(menuTabelDariMigrasi(null), ENTRI_MENU)
const TOMBOL = SEMUA.golongan.flatMap((g) => g.modul)

describe('sidebar satu tombol per modul, dikelompokkan GROUPMENU', () => {
  it('tepat satu tombol per modul yang dimigrasi dan aktif - kini empat', () => {
    const aktif = TOMBOL.filter((t) => t.halaman !== null)
    expect(aktif.map((t) => t.kode).sort()).toEqual(MODUL_FRONTEND.map((m) => m.nama).sort())
    expect(aktif).toHaveLength(4)
    expect(new Set(TOMBOL.map((t) => t.kode)).size).toBe(TOMBOL.length)
  })

  it('Treaty Contract Out tetap SATU tombol (tco5)', () => {
    expect(TOMBOL.filter((t) => t.label === FOLDER_KORPUS.treatyContractOut)).toHaveLength(1)
  })

  it("tombol nonaktif untuk DIMIGRASI='0' - tepat folder korpus yang belum dimigrasi", () => {
    const nonaktif = TOMBOL.filter((t) => t.halaman === null).map((t) => t.label)
    expect(nonaktif.sort()).toEqual(folderKorpusBelumDimigrasi())
    expect(NAV).toContain('aria-disabled="true"')
    expect(NAV).toContain('{KETERANGAN_BELUM_DIMIGRASI}')
  })

  it('kepala bagian = GROUPMENU, urutan TREATY, FACULTATIVE, KLAIM, MASTER', () => {
    expect(SEMUA.golongan.map((g) => g.kode)).toEqual(['TREATY', 'FACULTATIVE', 'KLAIM', 'MASTER'])
  })

  it('tidak ada elemen buka-tutup kelompok di sidebar', () => {
    for (const terlarang of ['KelompokMenu', 'aria-expanded', 'kelompok__panah', 'kelompok__judul', 'IkonChevron', 'kelompok__isi']) {
      expect(NAV, terlarang).not.toContain(terlarang)
    }
    expect(existsSync(join(AKAR_APLIKASI, 'inti', 'frontend', 'components', 'KelompokMenu.tsx'))).toBe(false)
    expect(existsSync(join(AKAR_APLIKASI, 'inti', 'frontend', 'lib', 'lipatMenu.ts'))).toBe(false)
  })
})

describe('klik tombol modul membuka halaman awalnya', () => {
  it.each([
    ['claimlife', 'inbox'],
    ['premiumlistlife', 'premiumlist'],
    ['komiteclaimlife', 'komite'],
    ['treatycontractout', 'tco-tahun'],
  ])('%s → %s', (modul, halaman) => {
    expect(TOMBOL.find((t) => t.kode === modul)?.halaman).toBe(halaman)
    expect(MODUL_FRONTEND.find((m) => m.nama === modul)?.halamanAwal).toBe(halaman)
  })

  it('tombol memanggil pilih() dengan halaman tombol itu', () => {
    expect(NAV).toContain('const tujuan = m.halaman')
    expect(NAV).toContain('pilih(tujuan)')
  })

  it('tombol modul menyala selama halaman modul itu tampil (mis. Outstanding milik Claim Life)', () => {
    expect(TOMBOL.find((t) => t.kode === 'claimlife')?.halamanModul).toContain('outstanding')
    expect(NAV).toContain('const aktif = m.halamanModul.includes(halaman as H)')
  })

  it('tombol Register di Inbox Claim Life tetap membuka Register', () => {
    const inbox = readFileSync(join(AKAR_APLIKASI, 'modul', 'claimlife', 'frontend', 'pages', 'InboxClaimLife.tsx'), 'utf8')
    const rute = readFileSync(join(AKAR_APLIKASI, 'modul', 'claimlife', 'frontend', 'rute.tsx'), 'utf8')
    expect(inbox).toContain('onClick={onRegister}')
    expect(inbox).toContain('{MENU.register}')
    expect(MENU.register).toBe('Register')
    const panggil = rute.slice(rute.indexOf('onRegister={() => {'))
    expect(panggil.slice(0, 80)).toContain("onPindah('register')")
  })
})
