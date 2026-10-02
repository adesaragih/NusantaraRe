import { afterEach, describe, expect, it, vi } from 'vitest'

import { KELOLA_USER } from '../labels'
import { KODE_MENU_KELOLA_USER, modulUntukAkun } from '../lib/daftarMenu'
import {
  ambilDaftarPengguna,
  ambilPilihanPengguna,
  aturSandiPengguna,
  buatPengguna,
  bukaKunciPengguna,
  hapusPengguna,
  setelAktifPengguna,
  ubahPengguna,
  type PilihanKelola,
  type RingkasAkun,
} from './api'
import {
  alihkan,
  badanBaru,
  badanSandi,
  badanUbah,
  isianKosong,
  kelompokMenu,
  periksaIsian,
  saringDaftar,
  setelDivisi,
  setelOrganisasi,
  tabGalat,
  teksJenjang,
  unitUntuk,
} from './aturan'

const PILIHAN: PilihanKelola = {
  organisasi: [{ kode: 'RNM', nama: 'Nusantara Re' }, { kode: 'LAIN', nama: 'Lain' }],
  divisi: [
    { kode: 'TECH', nama: 'Teknik', induk: 'RNM' },
    { kode: 'FIN', nama: 'Keuangan', induk: 'RNM' },
    { kode: 'L1', nama: 'Lain 1', induk: 'LAIN' },
  ],
  unit: [
    { kode: 'CLM', nama: 'Klaim', induk: 'TECH' },
    { kode: 'TAX', nama: 'Pajak', induk: 'FIN' },
  ],
  workbasket: [{ kode: 'ReasLifeAdmin', nama: 'Admin' }],
  menu: [
    { kode: 'premiumlistlife', label: 'PremiumList Life', golongan: 'TREATY', dimigrasi: true },
    { kode: 'claimfacin', label: 'Claim Fac In', golongan: 'KLAIM', dimigrasi: false },
    { kode: 'claimlife', label: 'Claim Life', golongan: 'KLAIM', dimigrasi: true },
    { kode: KODE_MENU_KELOLA_USER, label: 'Kelola User', golongan: 'ADMIN', dimigrasi: true },
  ],
}

describe('jenjang Organisasi → Divisi → Unit', () => {
  it('pilihan bawah hanya milik pilihan atasnya', () => {
    expect(unitUntuk(PILIHAN, 'TECH').map((u) => u.kode)).toEqual(['CLM'])
    expect(unitUntuk(PILIHAN, '')).toEqual([])
  })

  it('mengganti tingkat atas mengosongkan tingkat bawah yang bukan miliknya', () => {
    const isi = { ...isianKosong(), organisasi: 'RNM', divisi: 'TECH', unit: 'CLM' }
    expect(setelOrganisasi(isi, 'RNM', PILIHAN)).toMatchObject({ divisi: 'TECH', unit: 'CLM' })
    expect(setelOrganisasi(isi, 'LAIN', PILIHAN)).toMatchObject({ organisasi: 'LAIN', divisi: '', unit: '' })
    expect(setelDivisi(isi, 'FIN', PILIHAN)).toMatchObject({ divisi: 'FIN', unit: '' })
    expect(setelOrganisasi(isi, '', PILIHAN)).toMatchObject({ organisasi: '', divisi: '', unit: '' })
  })
})

describe('periksa isian', () => {
  const sah = { ...isianKosong(), akunId: 'UJI.USER-1', nama: 'Uji', sandi: 'Sandi-Uji-01', ulangiSandi: 'Sandi-Uji-01' }

  it('akun baru: username, nama, sandi awal minimal 10 karakter dan sama', () => {
    expect(periksaIsian(sah, true, 'ADMIN')).toBeNull()
    expect(periksaIsian({ ...sah, akunId: 'ada spasi' }, true, 'ADMIN')).toBe(KELOLA_USER.galatAkun)
    // Titik saja - jalurnya dibersihkan peramban (`login.AkunSah`).
    expect(periksaIsian({ ...sah, akunId: '..' }, true, 'ADMIN')).toBe(KELOLA_USER.galatAkun)
    expect(periksaIsian({ ...sah, akunId: 'a..b' }, true, 'ADMIN')).toBeNull()
    expect(periksaIsian({ ...sah, nama: '  ' }, true, 'ADMIN')).toBe(KELOLA_USER.galatNama)
    expect(periksaIsian({ ...sah, sandi: '123456789', ulangiSandi: '123456789' }, true, 'ADMIN')).toBe(KELOLA_USER.galatSandi)
    expect(periksaIsian({ ...sah, ulangiSandi: 'lain-lain-01' }, true, 'ADMIN')).toBe(KELOLA_USER.galatUlangi)
  })

  it('ubah: password kosong = tidak diganti; terisi = diperiksa seperti saat buat (tab Security)', () => {
    const ubah = { ...sah, sandi: '', ulangiSandi: '', akunId: 'LAIN' }
    expect(periksaIsian(ubah, false, 'ADMIN')).toBeNull()
    expect(periksaIsian({ ...ubah, sandi: 'pendek', ulangiSandi: 'pendek' }, false, 'ADMIN')).toBe(KELOLA_USER.galatSandi)
    expect(periksaIsian({ ...ubah, sandi: 'Sandi-Uji-02', ulangiSandi: 'Sandi-Uji-03' }, false, 'ADMIN')).toBe(KELOLA_USER.galatUlangi)
    expect(periksaIsian({ ...ubah, ulangiSandi: 'Sandi-Uji-02' }, false, 'ADMIN')).toBe(KELOLA_USER.galatSandi)
    expect(periksaIsian({ ...ubah, sandi: 'Sandi-Uji-02', ulangiSandi: 'Sandi-Uji-02' }, false, 'ADMIN')).toBeNull()
    // Galat password membuka tab Security; galat lain tab Profil.
    expect(tabGalat(KELOLA_USER.galatSandi)).toBe('security')
    expect(tabGalat(KELOLA_USER.galatUlangi)).toBe('security')
    expect(tabGalat(KELOLA_USER.galatNama)).toBe('profil')
  })

  it('ubah: Kelola User tidak dapat dicabut dari diri sendiri', () => {
    const ubah = { ...sah, sandi: '', ulangiSandi: '', akunId: 'ADMIN', menu: ['claimlife'] }
    expect(periksaIsian(ubah, false, 'ADMIN')).toBe(KELOLA_USER.menuDiriSendiri)
    expect(periksaIsian({ ...ubah, menu: [KODE_MENU_KELOLA_USER] }, false, 'ADMIN')).toBeNull()
    expect(periksaIsian(ubah, false, 'ADMIN-LAIN')).toBeNull()
  })
})

describe('badan permintaan', () => {
  it('sandi hanya di badan buat dan Security; ubah tanpa username dan sandi', () => {
    const isi = { ...isianKosong(), akunId: ' UJI ', nama: ' Uji ', sandi: 'Sandi-Uji-01', ulangiSandi: 'Sandi-Uji-01', menu: ['claimlife'] }
    // User baru bawaannya wajib ganti password ("Change Password Next Login" tercentang).
    expect(isianKosong().wajibGanti).toBe(true)
    expect(badanBaru(isi)).toEqual({
      akunId: 'UJI', nama: 'Uji', sandi: 'Sandi-Uji-01', wajibGanti: true, organisasi: '', divisi: '', unit: '', workbasket: [], menu: ['claimlife'],
    })
    expect(badanBaru({ ...isi, wajibGanti: false }).wajibGanti).toBe(false)
    const ubah = badanUbah(isi)
    expect(Object.keys(ubah).sort()).toEqual(['divisi', 'menu', 'nama', 'organisasi', 'unit', 'workbasket'])
  })

  it('Security saat ubah: dikirim hanya bila password diisi atau centangnya berubah', () => {
    const isi = { ...isianKosong(), sandi: '', wajibGanti: false }
    expect(badanSandi(isi, false)).toBeNull()
    expect(badanSandi({ ...isi, wajibGanti: true }, false)).toEqual({ sandi: '', wajibGanti: true })
    expect(badanSandi({ ...isi, sandi: 'Sandi-Uji-09' }, false)).toEqual({ sandi: 'Sandi-Uji-09', wajibGanti: false })
  })

  it('kotak centang: urut, tanpa ganda', () => {
    expect(alihkan(['b', 'a'], 'c', true)).toEqual(['a', 'b', 'c'])
    expect(alihkan(['a', 'b'], 'a', true)).toEqual(['a', 'b'])
    expect(alihkan(['a', 'b'], 'a', false)).toEqual(['b'])
  })
})

describe('daftar dan menu', () => {
  it('menu dikelompokkan per golongan menurut urutan server', () => {
    expect(kelompokMenu(PILIHAN.menu).map((g) => [g.golongan, g.menu.map((m) => m.kode)])).toEqual([
      ['TREATY', ['premiumlistlife']],
      ['KLAIM', ['claimfacin', 'claimlife']],
      ['ADMIN', [KODE_MENU_KELOLA_USER]],
    ])
  })

  it('saringan daftar mencari username dan nama, setiap kata', () => {
    const a = (akunId: string, nama: string): RingkasAkun => ({
      akunId, nama, organisasi: '', divisi: '', unit: '', aktif: true, terkunci: false, wajibGantiSandi: false, loginTerakhir: '',
    })
    const daftar = [a('SUPERADMIN', 'Super Admin'), a('UJI-1', 'Budi Klaim')]
    expect(saringDaftar(daftar, 'klaim').map((x) => x.akunId)).toEqual(['UJI-1'])
    expect(saringDaftar(daftar, 'admin super').map((x) => x.akunId)).toEqual(['SUPERADMIN'])
    expect(saringDaftar(daftar, ' ')).toHaveLength(2)
    expect(teksJenjang({ organisasi: 'RNM', divisi: '', unit: '' })).toBe('RNM')
    expect(teksJenjang({ organisasi: '', divisi: '', unit: '' })).toBe('—')
  })

  it('modul untuk akun: modul aktif yang menunya dipegang; null = tanpa saringan', () => {
    expect(modulUntukAkun(['claimlife', 'premiumlistlife'], ['premiumlistlife', 'kelolauser'], [])).toEqual(['premiumlistlife'])
    expect(modulUntukAkun(null, ['claimlife'], ['claimlife', 'premiumlistlife'])).toEqual(['claimlife'])
    expect(modulUntukAkun(['claimlife'], null, [])).toEqual(['claimlife'])
    expect(modulUntukAkun(null, null, ['claimlife'])).toBeNull()
  })
})

describe('klien /api/admin', () => {
  afterEach(() => {
    vi.unstubAllGlobals()
  })

  it('jalur, metode, dan badan setiap panggilan', async () => {
    const tertangkap: { url: string; init: RequestInit }[] = []
    vi.stubGlobal(
      'fetch',
      vi.fn((url: string, init: RequestInit) => {
        tertangkap.push({ url, init })
        return Promise.resolve(init.method === 'DELETE' ? new Response(null, { status: 204 }) : new Response('[]', { status: 200 }))
      }),
    )
    await ambilDaftarPengguna()
    await ambilPilihanPengguna()
    await buatPengguna({ akunId: 'U/1', nama: 'U', sandi: 'Sandi-Uji-01', wajibGanti: true, organisasi: '', divisi: '', unit: '', workbasket: [], menu: [] })
    await ubahPengguna('U/1', { nama: 'U', organisasi: '', divisi: '', unit: '', workbasket: [], menu: [] })
    await setelAktifPengguna('U/1', false)
    await bukaKunciPengguna('U/1')
    await aturSandiPengguna('U/1', { sandi: '', wajibGanti: true })
    await expect(hapusPengguna('U/1')).resolves.toBeUndefined()
    expect(tertangkap.map((t) => `${t.init.method ?? 'GET'} ${t.url}`)).toEqual([
      'GET /api/admin/pengguna',
      'GET /api/admin/pilihan-pengguna',
      'POST /api/admin/pengguna',
      'PUT /api/admin/pengguna/U%2F1',
      'POST /api/admin/pengguna/U%2F1/aktif',
      'POST /api/admin/pengguna/U%2F1/buka-kunci',
      'POST /api/admin/pengguna/U%2F1/sandi',
      'DELETE /api/admin/pengguna/U%2F1',
    ])
    expect(JSON.parse(String(tertangkap[4]?.init.body))).toEqual({ aktif: false })
    expect(tertangkap[5]?.init.body).toBeUndefined()
    expect(JSON.parse(String(tertangkap[6]?.init.body))).toEqual({ sandi: '', wajibGanti: true })
  })
})
