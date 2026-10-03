import { describe, expect, it } from 'vitest'

import { menuTabelDariMigrasi } from '../inti/frontend/uji/menuBersih'
import { KELOLA_USER } from '../inti/frontend/labels'
import {
  daftarPalet,
  HALAMAN_KELOLA_USER,
  KODE_MENU_KELOLA_USER,
  susunMenu,
  type MenuTabel,
} from '../inti/frontend/lib/daftarMenu'
import { ENTRI_MENU, halamanAktif, MODUL_BACKEND } from './daftar'

// Kelola User (keputusan work owner 01-10-2026) - menu APLIKASI di golongan
// ADMIN, bukan baris M_NAV_MENU (tabel itu tetap dua puluh baris, satu per
// folder modul korpus). Ia tampil HANYA bila `GET /api/menu` mengirimnya -
// yaitu bagi akun yang memegang KODE-nya (`menu.SaringMenuUntukAkun`).

const ADMIN: MenuTabel = {
  golongan: [
    ...menuTabelDariMigrasi(null).golongan,
    {
      kode: 'ADMIN',
      modul: [{ kode: KODE_MENU_KELOLA_USER, label: KELOLA_USER.judul, modul: KODE_MENU_KELOLA_USER, urutan: 1, dimigrasi: true }],
    },
  ],
}

describe('menu aplikasi Kelola User', () => {
  it('satu entri, berpemilik KODE menunya, di luar MODUL_AKTIF', () => {
    const entri = ENTRI_MENU.filter((e) => e.modul === HALAMAN_KELOLA_USER)
    expect(entri).toEqual([
      { modul: HALAMAN_KELOLA_USER, label: KELOLA_USER.judul, kelompok: KELOLA_USER.judul, pemilik: KODE_MENU_KELOLA_USER },
    ])
    expect(MODUL_BACKEND[HALAMAN_KELOLA_USER]).toBeNull()
    expect(halamanAktif(HALAMAN_KELOLA_USER, ['claimlife'])).toBe(true)
  })

  it('tampil di golongan ADMIN - dan di palet - hanya bila backend mengirimnya', () => {
    const dengan = susunMenu(ADMIN, ENTRI_MENU)
    const admin = dengan.golongan.find((g) => g.kode === 'ADMIN')
    expect(admin?.modul).toEqual([
      { kode: KODE_MENU_KELOLA_USER, label: KELOLA_USER.judul, halaman: HALAMAN_KELOLA_USER, halamanModul: [HALAMAN_KELOLA_USER] },
    ])
    expect(dengan.golongan.at(-1)?.kode).toBe('ADMIN')
    expect(daftarPalet(dengan.entri).map((h) => h.modul)).toContain(HALAMAN_KELOLA_USER)
    expect(dengan.tanpaRute).toEqual([])

    const tanpa = susunMenu(menuTabelDariMigrasi(null), ENTRI_MENU)
    expect(tanpa.golongan.map((g) => g.kode)).not.toContain('ADMIN')
    expect(daftarPalet(tanpa.entri).map((h) => h.modul)).not.toContain(HALAMAN_KELOLA_USER)
  })
})
