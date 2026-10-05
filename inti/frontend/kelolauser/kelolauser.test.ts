import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { afterEach, describe, expect, it, vi } from 'vitest'

import { KELOLA_USER, LOGIN } from '../labels'
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
  periksaGanda,
  periksaIsian,
  isianDari,
  periksaKontak,
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
    const isi = {
      ...isianKosong(), akunId: ' UJI ', nama: ' Uji ', sandi: 'Sandi-Uji-01', ulangiSandi: 'Sandi-Uji-01', menu: ['claimlife'],
      email: ' uji@nusantara.example ', telepon: ' 0812 3456 7890 ', nik: ' UJI-001 ', jabatan: ' Analyst ',
    }
    // User baru bawaannya wajib ganti password ("Change Password Next Login" tercentang).
    expect(isianKosong().wajibGanti).toBe(true)
    expect(badanBaru(isi)).toEqual({
      akunId: 'UJI', nama: 'Uji', sandi: 'Sandi-Uji-01', wajibGanti: true, organisasi: '', divisi: '', unit: '', workbasket: [], menu: ['claimlife'],
      menuLihat: [],
      // Kontak (Kelola User 03-10-2026): spasi tepi dibuang.
      email: 'uji@nusantara.example', telepon: '0812 3456 7890', nik: 'UJI-001', jabatan: 'Analyst',
    })
    expect(badanBaru({ ...isi, wajibGanti: false }).wajibGanti).toBe(false)
    const ubah = badanUbah(isi)
    expect(Object.keys(ubah).sort()).toEqual(['divisi', 'email', 'jabatan', 'menu', 'menuLihat', 'nama', 'nik', 'organisasi', 'telepon', 'unit', 'workbasket'])
  })

  it('hak View only (migrasi 914): terbaca dari akun, dan menu yang dicabut tidak ikut terkirim', () => {
    const r = {
      akunId: 'UJI-K', nama: 'Uji', organisasi: '', divisi: '', unit: '', aktif: true, terkunci: false, wajibGantiSandi: false,
      loginTerakhir: '', email: '', telepon: '', nik: '', jabatan: '', contactId: 'CON-1009',
      workbasket: [], menu: ['accounts', 'aggregate'], menuLihat: ['accounts', 'aggregate'],
    }
    const isi = isianDari(r)
    expect(isi.menuLihat).toEqual(['accounts', 'aggregate'])
    expect(badanUbah({ ...isi, menu: ['aggregate'] }).menuLihat).toEqual(['aggregate'])
    // Backend lama tanpa medan menuLihat = semua Full.
    expect(isianDari({ ...r, menuLihat: undefined }).menuLihat).toEqual([])
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
      email: '', telepon: '', nik: '', jabatan: '', contactId: '',
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
    const kontak = { email: '', telepon: '', nik: '', jabatan: '' }
    await buatPengguna({ akunId: 'U/1', nama: 'U', sandi: 'Sandi-Uji-01', wajibGanti: true, organisasi: '', divisi: '', unit: '', workbasket: [], menu: [], menuLihat: [], ...kontak })
    await ubahPengguna('U/1', { nama: 'U', organisasi: '', divisi: '', unit: '', workbasket: [], menu: [], menuLihat: [], ...kontak })
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

describe('kontak akun - Email, Phone Number, Employee ID (NIK), Position (Kelola User 03-10-2026)', () => {
  const sah = { ...isianKosong(), akunId: 'UJI-K', nama: 'Uji', sandi: 'Sandi-Uji-01', ulangiSandi: 'Sandi-Uji-01' }

  it('opsional: semua kosong sah; terisi sah diterima', () => {
    expect(periksaIsian(sah, true, 'ADMIN')).toBeNull()
    expect(
      periksaIsian({ ...sah, email: 'uji.user@nusantara.example', telepon: '+62 812-3456-7890', nik: 'UJI/2026.001-A', jabatan: 'Underwriter' }, true, 'ADMIN'),
    ).toBeNull()
    expect(periksaKontak({ email: '', telepon: '02112345', nik: '', jabatan: '' })).toBeNull()
  })

  it('format tidak sah ditolak dengan pesan berbahasa Inggris yang sama dengan backend', () => {
    expect(periksaIsian({ ...sah, email: 'uji@nusantara' }, true, 'ADMIN')).toBe(KELOLA_USER.galatEmail)
    expect(periksaIsian({ ...sah, email: 'uji user@nusantara.example' }, true, 'ADMIN')).toBe(KELOLA_USER.galatEmail)
    expect(periksaIsian({ ...sah, telepon: '1234567' }, true, 'ADMIN')).toBe(KELOLA_USER.galatTelepon)
    expect(periksaIsian({ ...sah, telepon: '1234567890123456' }, true, 'ADMIN')).toBe(KELOLA_USER.galatTelepon)
    expect(periksaIsian({ ...sah, telepon: '0812+34567890' }, true, 'ADMIN')).toBe(KELOLA_USER.galatTelepon)
    expect(periksaIsian({ ...sah, nik: 'UJI#001' }, true, 'ADMIN')).toBe(KELOLA_USER.galatNIK)
    expect(periksaIsian({ ...sah, nik: '1'.repeat(31) }, true, 'ADMIN')).toBe(KELOLA_USER.galatNIK)
    expect(periksaIsian({ ...sah, jabatan: 'J'.repeat(151) }, true, 'ADMIN')).toBe(KELOLA_USER.galatJabatan)
    for (const p of [KELOLA_USER.galatEmail, KELOLA_USER.galatTelepon, KELOLA_USER.galatNIK, KELOLA_USER.galatJabatan]) {
      expect(p).toMatch(/^(Email|Phone Number|Employee ID|Position) /)
    }
  })

  it('label berbahasa Inggris', () => {
    expect([KELOLA_USER.email, KELOLA_USER.telepon, KELOLA_USER.nik, KELOLA_USER.jabatan]).toEqual([
      'Email', 'Phone Number', 'Employee ID (NIK)', 'Position',
    ])
  })

  it('form ubah terisi dari akun tersimpan', () => {
    const r = {
      akunId: 'UJI-K', nama: 'Uji', organisasi: '', divisi: '', unit: '', aktif: true, terkunci: false, wajibGantiSandi: false,
      loginTerakhir: '', email: 'k@x.example', telepon: '08123456789', nik: 'UJI-9', jabatan: 'Staff', contactId: 'CON-1009',
      workbasket: [], menu: [],
    }
    const isi = isianDari(r)
    expect([isi.email, isi.telepon, isi.nik, isi.jabatan]).toEqual(['k@x.example', '08123456789', 'UJI-9', 'Staff'])
  })
})

describe('identitas akun - Contact ID, username dan email sudah terdaftar (migrasi 905, 03-10-2026)', () => {
  const akun = (akunId: string, email: string, contactId: string): RingkasAkun => ({
    akunId, nama: `Nama ${akunId}`, organisasi: '', divisi: '', unit: '', aktif: true, terkunci: false, wajibGantiSandi: false,
    loginTerakhir: '', email, telepon: '', nik: '', jabatan: '', contactId,
  })
  const daftar = [akun('UJI-ADMIN', 'uji.admin@nusantara.example', 'CON-1001'), akun('UJI-DUA', 'uji.dua@nusantara.example', 'CON-1002')]
  const isi = (akunId: string, email: string) => ({ ...isianKosong(), akunId, email })

  it('username sudah terdaftar: sama persis atau beda huruf saja', () => {
    for (const id of ['UJI-ADMIN', ' uji-admin ', 'Uji-Dua']) {
      expect(periksaGanda(isi(id, ''), null, daftar)).toBe(KELOLA_USER.galatUsernameTerdaftar)
    }
    expect(periksaGanda(isi('UJI-BARU', ''), null, daftar)).toBeNull()
  })

  it('email sudah terdaftar: email akun lain, tanpa beda huruf; milik sendiri boleh', () => {
    expect(periksaGanda(isi('UJI-BARU', ' UJI.ADMIN@nusantara.example '), null, daftar)).toBe(KELOLA_USER.galatEmailTerdaftar)
    expect(periksaGanda(isi('UJI-ADMIN', 'Uji.Admin@nusantara.example'), 'UJI-ADMIN', daftar)).toBeNull()
    expect(periksaGanda(isi('UJI-ADMIN', 'UJI.Dua@nusantara.example'), 'UJI-ADMIN', daftar)).toBe(KELOLA_USER.galatEmailTerdaftar)
    expect(periksaGanda(isi('UJI-BARU', ''), null, daftar)).toBeNull()
  })

  it('daftar belum termuat: tidak menebak, backend tetap menjaga', () => {
    expect(periksaGanda(isi('UJI-ADMIN', ''), null, null)).toBeNull()
  })

  it('pesan ganda berbahasa Inggris dan SAMA dengan jawaban 409 backend; galat ganda tinggal di tab Profil', () => {
    expect(KELOLA_USER.galatUsernameTerdaftar).toBe('Username is already registered')
    expect(KELOLA_USER.galatEmailTerdaftar).toBe('Email is already registered to another account')
    expect(tabGalat(KELOLA_USER.galatUsernameTerdaftar)).toBe('profil')
    expect(tabGalat(KELOLA_USER.galatEmailTerdaftar)).toBe('profil')
  })

  it('cari menemukan Contact ID dan email', () => {
    expect(saringDaftar(daftar, 'con-1002').map((x) => x.akunId)).toEqual(['UJI-DUA'])
    expect(saringDaftar(daftar, 'uji.admin@').map((x) => x.akunId)).toEqual(['UJI-ADMIN'])
  })

  it('label: Contact ID; layar login tetap username saja (login lewat email dibatalkan work owner)', () => {
    expect(KELOLA_USER.contactId).toBe('Contact ID')
    expect(LOGIN.akun).toBe('Username')
    expect(LOGIN.salah).toBe('Incorrect username or password.')
  })
})

describe('tampilan daftar user', () => {
  // Permintaan work owner 03-10-2026: tombol aksi baris (Ubah, Nonaktifkan, Hapus) berderet ke samping,
  // bukan bertumpuk - sel `.table__actions` sesempit isinya, jadi wadahnya tidak boleh membungkus.
  it('tombol aksi baris berderet ke samping', () => {
    const css = readFileSync(join(__dirname, '..', 'styles.css'), 'utf8')
    const aturan = /\.kelola-user__aksi \{([^}]*)\}/.exec(css)?.[1] ?? ''
    expect(aturan).toContain('flex-wrap: nowrap;')
    expect(aturan).not.toContain('flex-wrap: wrap;')
  })
})
