// Uji login sungguhan di sisi layar - M_LOGIN_GO (keputusan work owner 01-10-2026).

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { afterEach, describe, expect, it, vi } from 'vitest'

import { ApiFailure, ambilSesiSaya, gantiSandiLogin, keluarLogin, masukLogin } from '../klien'
import { LOGIN, PANJANG_MIN_SANDI, PERAN } from '../labels'
import { PERISTIWA_SESI_BERAKHIR, sesiDariProfil } from '../store/sesi'
import { periksaSandiBaru, pesanGagalGanti } from './GantiSandi'
import { pesanGagalLogin } from './Login'

const tertangkap: { url: string; init: RequestInit }[] = []

function pasangFetch(status: number, badan: string): void {
  vi.stubGlobal('fetch', (url: string, init: RequestInit) => {
    tertangkap.push({ url, init })
    const isi = status === 204 ? null : badan
    return Promise.resolve(new Response(isi, { status, headers: { 'Content-Type': 'application/json' } }))
  })
}

afterEach(() => {
  tertangkap.length = 0
  vi.unstubAllGlobals()
  vi.unstubAllEnvs()
})

const PROFIL = `{"akunId":"UJI-ADMIN","nama":"Uji Admin","peran":["ReasLifeAdmin","ReasFacInAdmin"],` +
  `"organisasi":"RNM","divisi":"TECH","unit":"CLM","wajibGantiSandi":false}`

describe('klien /api/auth/*', () => {
  it('login, profil, ganti sandi, logout memakai jalur dan metode backend', async () => {
    pasangFetch(200, PROFIL)
    const p = await masukLogin('UJI-ADMIN', 'Sandi-Uji-0001')
    expect(p.akunId).toBe('UJI-ADMIN')
    await ambilSesiSaya()
    await gantiSandiLogin('Sandi-Uji-0001', 'Sandi-Uji-0002')
    pasangFetch(204, '')
    await keluarLogin()
    expect(tertangkap.map((t) => `${t.init.method ?? 'GET'} ${new URL(t.url, 'http://x').pathname}`)).toEqual([
      'POST /api/auth/login',
      'GET /api/auth/saya',
      'POST /api/auth/ganti-sandi',
      'POST /api/auth/logout',
    ])
    expect(JSON.parse(String(tertangkap[0]?.init.body))).toEqual({ akun: 'UJI-ADMIN', sandi: 'Sandi-Uji-0001' })
    expect(JSON.parse(String(tertangkap[2]?.init.body))).toEqual({ sandiLama: 'Sandi-Uji-0001', sandiBaru: 'Sandi-Uji-0002' })
  })

  it('401 di tengah pemakaian memicu PERISTIWA_SESI_BERAKHIR - kecuali mode stub', async () => {
    const jendela = new EventTarget()
    let terpicu = 0
    jendela.addEventListener(PERISTIWA_SESI_BERAKHIR, () => {
      terpicu++
    })
    vi.stubGlobal('window', jendela)
    vi.stubEnv('VITE_AUTH_STUB', '')
    pasangFetch(401, '{"galat":"belum login atau sesi sudah berakhir"}')
    await ambilSesiSaya().catch(() => undefined)
    expect(terpicu).toBe(1)
    vi.stubEnv('VITE_AUTH_STUB', 'true')
    await ambilSesiSaya().catch(() => undefined)
    expect(terpicu).toBe(1)
  })
})

describe('pesan layar', () => {
  it('pesan gagal login dipilih dari status, bukan teks backend', () => {
    expect(pesanGagalLogin(new ApiFailure(401, { code: 'DITOLAK_BACKEND', message: 'akun atau sandi salah' }))).toBe(LOGIN.salah)
    expect(pesanGagalLogin(new ApiFailure(423, { code: 'DITOLAK_BACKEND', message: 'terkunci' }))).toBe(LOGIN.terkunci)
    expect(pesanGagalLogin(new ApiFailure(503, { code: 'DITOLAK_BACKEND', message: 'SESI_RAHASIA' }))).toBe(LOGIN.tidakTersedia)
    expect(pesanGagalLogin(new TypeError('Failed to fetch'))).toBe(LOGIN.gagal)
  })

  it('aturan sandi baru sama dengan backend - dihitung KARAKTER', () => {
    expect(PANJANG_MIN_SANDI).toBe(10)
    expect(periksaSandiBaru('123456789', '123456789')).toBe(LOGIN.minimal)
    expect(periksaSandiBaru('ééééééééé', 'ééééééééé')).toBe(LOGIN.minimal)
    expect(periksaSandiBaru('1234567890', '1234567899')).toBe(LOGIN.tidakSama)
    expect(periksaSandiBaru('1234567890', '1234567890')).toBeNull()
  })

  it('ganti sandi: sandi lama salah dikenali, 401 bukan galat medan', () => {
    expect(pesanGagalGanti(new ApiFailure(400, { code: 'DITOLAK_BACKEND', message: 'sandi lama salah' }))).toBe(LOGIN.lamaSalah)
    // Kalimat backend DIPETAKAN ke label berbahasa Inggris (layar login berbahasa
    // Inggris, permintaan work owner 01-10-2026); yang tak dikenal = pesan umum.
    expect(pesanGagalGanti(new ApiFailure(400, { code: 'DITOLAK_BACKEND', message: 'sandi minimal 10 karakter' }))).toBe(LOGIN.minimal)
    expect(pesanGagalGanti(new ApiFailure(400, { code: 'DITOLAK_BACKEND', message: 'sandi maksimal 72 byte' }))).toBe(LOGIN.terlaluPanjang)
    expect(pesanGagalGanti(new ApiFailure(400, { code: 'DITOLAK_BACKEND', message: 'sandi baru sama dengan sandi lama' }))).toBe(LOGIN.samaDenganLama)
    expect(pesanGagalGanti(new ApiFailure(400, { code: 'DITOLAK_BACKEND', message: 'kalimat lain' }))).toBe(LOGIN.gagal)
    expect(pesanGagalGanti(new ApiFailure(401, { code: 'DITOLAK_BACKEND', message: 'sesi sudah berakhir' }))).toBe(LOGIN.gagal)
  })
})

describe('sesi dari profil', () => {
  it('peran = workbasket yang dikenal layar; nama ikut', () => {
    const s = sesiDariProfil(JSON.parse(PROFIL) as { akunId: string; nama: string; peran: string[] })
    expect(s).toEqual({ akunID: 'UJI-ADMIN', nama: 'Uji Admin', peran: [PERAN.admin] })
  })
})

describe('sandi tidak pernah disimpan di layar', () => {
  it.each(['Login.tsx', 'GantiSandi.tsx'])('%s tanpa penyimpanan peramban', (berkas) => {
    const isi = readFileSync(join(__dirname, berkas), 'utf8')
    for (const terlarang of ['localStorage', 'sessionStorage', 'console.', 'document.cookie']) {
      expect(isi, terlarang).not.toContain(terlarang)
    }
  })
})

describe('tampilan loginbaru.html', () => {
  const login = readFileSync(join(__dirname, 'Login.tsx'), 'utf8')
  const css = readFileSync(join(__dirname, '..', 'styles.css'), 'utf8')

  // ⛔ Animasi loginbaru.html kembali (permintaan work owner 03-10-2026, asap
  // kopi ikut bergerak) - HANYA di lapisan SVG terpisah. Animasi di dalam SVG
  // utama menggambar ulang kartu kaca tiap bingkai: layar login dan ganti sandi
  // berkedip (02-10-2026).
  it('animasi hanya di lapisan terpisah; ilustrasi utama diam', () => {
    const ilustrasi = readFileSync(join(__dirname, 'IlustrasiLogin.tsx'), 'utf8')
    const svg = [...ilustrasi.matchAll(/<svg\b[^>]*>/g)].map((m) => m[0])
    expect(svg).toHaveLength(5)
    // SVG pertama = ilustrasi utama: tanpa kelas, jadi tanpa animasi.
    expect(svg[0]).not.toContain('className')
    for (const s of svg.slice(1)) expect(s).toContain('className="halaman-masuk__lapis halaman-masuk__')
    expect(svg.filter((s) => s.includes('halaman-masuk__melayang'))).toHaveLength(2)
    expect(svg.filter((s) => s.includes('halaman-masuk__uap'))).toHaveLength(2)
    // Satu-satunya pemilih halaman login (`.halaman-masuk*`) ber-`animation`: melayang dan uap.
    const blok = css.slice(css.indexOf('HALAMAN LOGIN M_LOGIN_GO'))
    const beranimasi = [...blok.replace(/\/\*[\s\S]*?\*\//g, '').matchAll(/([^{}]+)\{([^{}]*)\}/g)]
      .filter((m) => /(^|[;\s])animation\s*:/.test(m[2] ?? '') && !/animation\s*:\s*none/.test(m[2] ?? ''))
      .map((m) => (m[1] ?? '').trim())
      .filter((p) => p.startsWith('.halaman-masuk'))
    expect(beranimasi.sort()).toEqual(['.halaman-masuk__melayang', '.halaman-masuk__uap'])
    expect(css).toContain('will-change: transform, opacity;')
    expect(css).toMatch(/@media \(prefers-reduced-motion: reduce\) \{\s*\.halaman-masuk__melayang,\s*\.halaman-masuk__uap \{\s*animation: none;/)
  })

  it('teks layar login berbahasa Inggris', () => {
    // Bahasa Inggris - permintaan work owner 01-10-2026.
    expect(LOGIN.judul).toBe('Hello Again!')
    expect(LOGIN.sub).toBe("Welcome back, it's great to see you again!")
    expect(LOGIN.isianAkun).toBe('Enter your username')
    expect(LOGIN.masuk).toBe('Login')
  })

  it('isian dikenali pengelola sandi peramban', () => {
    expect(login).toContain('autoComplete="username"')
    expect(login).toContain('autoComplete="current-password"')
    expect(login).toContain('aria-pressed={terlihat}')
  })

  it('logo bertulisan hitam lewat gambar latar, bukan atribut sumber', () => {
    expect(css).toContain("url('./assets/logo-login.png')")
    expect(login).not.toMatch(/<img/)
    // PNG 450x73 (IHDR): logo yang dikirim work owner, tulisan diubah hitam.
    const png = readFileSync(join(__dirname, '..', 'assets', 'logo-login.png'))
    expect(png.readUInt32BE(16)).toBe(450)
    expect(png.readUInt32BE(20)).toBe(73)
  })
})
