// Uji modul Accounts: nama menu = migrasi inti 908, label form VERBATIM tangkapan layar, aturan wajib isi, jalur API
// (akun hanya diinput sekali: nol View/Edit/hapus), dan CSS modul bertema Kelola User di bawah akar `.accounts`.

import { readdirSync, readFileSync, statSync } from 'node:fs'
import { join } from 'node:path'

import { afterEach, describe, expect, it, vi } from 'vitest'

import { PREFIX_ACC, ambilDaftar, ambilPilihan, cariOrganisasi, tambah } from './api'
import { adaGalat, atauStrip, isianKosong, periksa } from './aturan'
import { ACC, MENU_ACC } from './labels'

const AKAR = __dirname

describe('label Accounts', () => {
  it('nama menu = M_NAV_MENU.LABEL migrasi inti 908, kelompok MASTER', () => {
    const sql = readFileSync(join(AKAR, '../../../inti/backend/migrations/908_m_nav_menu_accounts.sql'), 'utf8')
    expect(sql).toContain(`'accounts', '${MENU_ACC.kelompok}', 'MASTER', 'accounts'`)
  })

  it('label form VERBATIM tangkapan layar Pega; Territory tidak ada', () => {
    expect([ACC.insuredName, ACC.orgId, ACC.groupBusiness, ACC.owner, ACC.description, ACC.wajib]).toEqual([
      'Insured Name',
      'Org ID',
      'Group Business',
      'Owner',
      'Description',
      'Value cannot be blank',
    ])
    const sumber = berkas(AKAR, '.ts').concat(berkas(AKAR, '.tsx')).map((p) => readFileSync(p, 'utf8'))
    expect(sumber.filter((s) => s.includes("'Territory'"))).toEqual([])
  })
})

describe('aturan form', () => {
  it('wajib isi Insured Name dan Group Business, pesan VERBATIM', () => {
    const g = periksa(isianKosong())
    expect(g).toEqual({ insured: 'Value cannot be blank', groupBusiness: 'Value cannot be blank' })
    expect(adaGalat(g)).toBe(true)
    expect(adaGalat(periksa({ insuredId: 'X', groupBusinessId: '10001', description: '' }))).toBe(false)
    expect(adaGalat(periksa({ insuredId: ' ', groupBusinessId: '10001', description: '' }))).toBe(true)
  })

  it('sel kosong bertanda strip', () => {
    expect(atauStrip(' ')).toBe('—')
    expect(atauStrip('ACC-1')).toBe('ACC-1')
  })
})

interface Rekam {
  url: string
  metode: string
  badan: unknown
}

describe('jalur API', () => {
  afterEach(() => {
    vi.unstubAllGlobals()
  })

  it('jalur, metode, dan badan SAMA dengan rute backend', async () => {
    const rekam: Rekam[] = []
    vi.stubGlobal(
      'fetch',
      vi.fn(async (url: string, init: RequestInit) => {
        rekam.push({ url, metode: init.method ?? 'GET', badan: typeof init.body === 'string' ? JSON.parse(init.body) : init.body })
        return new Response('{}', { status: 200 })
      }),
    )
    const isi = { insuredId: 'X', groupBusinessId: '1', description: 'UJI' }
    await ambilDaftar('', 2, 20)
    await ambilDaftar('uji pt', 1, 20)
    await ambilPilihan()
    await cariOrganisasi('dua')
    await tambah(isi)
    expect(rekam.map((r) => [r.metode, r.url])).toEqual([
      ['GET', `${PREFIX_ACC}?halaman=2&ukuran=20`],
      ['GET', `${PREFIX_ACC}?cari=uji+pt&halaman=1&ukuran=20`],
      ['GET', `${PREFIX_ACC}/pilihan`],
      ['GET', `${PREFIX_ACC}/organisasi?cari=dua`],
      ['POST', PREFIX_ACC],
    ])
    expect(rekam[4]!.badan).toEqual(isi)
  })

  it('akun hanya diinput sekali: nol PUT/PATCH/DELETE, nol tombol View/Edit/Delete', () => {
    const sumber = berkas(AKAR, '.ts').concat(berkas(AKAR, '.tsx')).map((p) => readFileSync(p, 'utf8')).join(' ')
    expect(sumber).not.toMatch(/metode: '(PUT|PATCH|DELETE)'/)
    expect(Object.values(ACC).filter((v) => typeof v === 'string' && /^(View|Edit|Delete)$/.test(v))).toEqual([])
  })
})

function berkas(dir: string, akhiran: string): string[] {
  return readdirSync(dir).flatMap((n) => {
    const p = join(dir, n)
    if (statSync(p).isDirectory()) return berkas(p, akhiran)
    return n.endsWith(akhiran) && !n.endsWith('.test.ts') ? [p] : []
  })
}

describe('gaya modul Accounts', () => {
  const css = readFileSync(join(AKAR, 'accounts.css'), 'utf8')
  const aturan = css.replace(/\/\*[\s\S]*?\*\//g, '')

  it('satu-satunya berkas CSS adalah accounts.css; rute.tsx mengimpornya dan membungkus halaman dengan akarnya', () => {
    expect(berkas(AKAR, '.css').map((p) => p.slice(AKAR.length + 1))).toEqual(['accounts.css'])
    const rute = readFileSync(join(AKAR, 'rute.tsx'), 'utf8')
    expect(rute).toContain("import './accounts.css'")
    expect(rute).toContain('<div className="accounts">')
    expect(css).toMatch(/^\.accounts \{\s*display: contents;\s*\}/m)
  })

  it('tema Kelola User: halaman berakar .accounts__akar dengan token --acc-* terang dan gelap, tanpa kelas inti kelola-user', () => {
    const halaman = readFileSync(join(AKAR, 'pages/Accounts.tsx'), 'utf8')
    expect(halaman).toContain('<section className="inbox accounts__akar">')
    expect(halaman).toContain('className="inbox__tabel accounts__tabel"')
    expect(css).toMatch(/^\.accounts \.accounts__akar \{[^}]*--acc-latar: #eef1f6;/m)
    expect(css).toMatch(/^:root\[data-theme="dark"\] \.accounts \.accounts__akar \{[^}]*--acc-latar: #1b2130;/m)
    expect(aturan).not.toMatch(/kelola-user/)
  })

  it('SETIAP pemilih di bawah akar .accounts - nol pemilih global; nol properti yang memerangkap Modal', () => {
    const pemilih = [...aturan.matchAll(/([^{}]+)\{/g)]
      .flatMap((m) => m[1]!.split(',').map((s) => s.trim()))
      .filter((p) => !p.startsWith('@'))
    expect(pemilih.length).toBeGreaterThan(5)
    for (const p of pemilih) expect(p, p).toMatch(/^(:root\[data-theme="dark"\] )?\.accounts(\s|$)/)
    expect(aturan).not.toMatch(/(^|[\s;{])(-webkit-)?(backdrop-filter|filter|transform|translate|rotate|scale|perspective|will-change|contain)\s*:/m)
  })
})
