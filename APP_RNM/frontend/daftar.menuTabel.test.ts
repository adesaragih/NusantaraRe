import { readdirSync, readFileSync } from 'node:fs'
import { join } from 'node:path'

import { afterEach, describe, expect, it, vi } from 'vitest'

import { AKAR_APLIKASI } from '../inti/frontend/uji/sumber'
import { ambilMenu } from '../inti/frontend/klien'
import { FOLDER_KORPUS } from './katalogKorpus'
import { HALAMAN_BERANDA } from '../inti/frontend/lib/daftarMenu'
import { ENTRI_MENU, MODUL_FRONTEND } from './daftar'

// Penjaga DUA ARAH: isi M_NAV_MENU (migrasi `inti` 900-949) ↔ `frontend/daftar.ts`
// (brief menu 30-09-2026 §3). SELURUH migrasi maju `inti/migrations/` dibaca,
// bukan 900 saja: menu berikutnya lahir di 901+ (panduan deploy bab 6).
//
// ⛔ Kenapa statik: sidebar dirakit dari `GET /api/menu` dan DIPOTONG dengan
// rute frontend - baris tabel tanpa rute tidak tampil, rute tanpa baris
// tabel juga tidak. Keduanya diam di layar (hanya satu baris konsol). Uji
// inilah yang membuatnya berbunyi sebelum sampai ke layar.

const FOLDER_MIGRASI_INTI = join(AKAR_APLIKASI, 'inti', 'backend', 'migrations')
const SQL = readdirSync(FOLDER_MIGRASI_INTI)
  .filter((n) => n.endsWith('.sql') && !n.endsWith('_down.sql'))
  .sort()
  .map((n) => readFileSync(join(FOLDER_MIGRASI_INTI, n), 'utf8'))
  .join('\n')

const POLA_KELOMPOK =
  /SELECT \{skema\}\.SEQ_M_NAV_MENU\.NEXTVAL, NULL, '([^']+)', '([^']+)', '([A-Z]+)', '([^']+)', \d+, '([01])' FROM DUAL/g
const POLA_BUTIR =
  /SELECT \{skema\}\.SEQ_M_NAV_MENU\.NEXTVAL, k\.ID, '([^']+)', '([^']+)', k\.GROUPMENU, k\.MODUL, \d+, k\.DIMIGRASI\s+FROM \{skema\}\.M_NAV_MENU k\s+WHERE k\.KODE = '([^']+)'/g

const kelompok = [...SQL.matchAll(POLA_KELOMPOK)].map((m) => ({ kode: m[1]!, label: m[2]!, dimigrasi: m[5] === '1' }))
const butir = [...SQL.matchAll(POLA_BUTIR)].map((m) => ({ kode: m[1]!, label: m[2]!, induk: m[3]! }))
const rute = ENTRI_MENU.filter((e) => e.modul !== HALAMAN_BERANDA)

describe('isi menu migrasi inti ↔ daftar.ts, dua arah', () => {
  it('setiap INSERT terbaca penjaga ini', () => {
    // INSERT berbentuk lain adalah baris yang uji ini tidak lihat.
    expect(kelompok.length + butir.length).toBe((SQL.match(/^INSERT INTO /gm) ?? []).length)
    expect(kelompok).toHaveLength(20)
    expect(butir).toHaveLength(5)
  })

  it('setiap KODE butir isi menu punya rute di daftar.ts', () => {
    for (const b of butir) expect(rute.map((e) => e.modul as string), b.kode).toContain(b.kode)
  })

  it('setiap butir daftar.ts punya baris di isi menu', () => {
    for (const e of rute) expect(butir.map((b) => b.kode), e.modul).toContain(e.modul)
  })

  it('LABEL butir = label frontend, VERBATIM', () => {
    for (const e of rute) expect(butir.find((b) => b.kode === e.modul)?.label, e.modul).toBe(e.label)
  })

  it('induk butir = modul pemilik rutenya', () => {
    for (const e of rute) expect(butir.find((b) => b.kode === e.modul)?.induk, e.modul).toBe(e.pemilik)
  })

  // Struktur tim satu folder per modul (30-09-2026): nama kelompok setiap modul
  // tinggal di `menu.ts`-nya sendiri. Ia harus nama folder korpus, sama dengan
  // LABEL barisnya di isi menu, dan sama untuk setiap butir modul itu.
  it('KELOMPOK setiap modul = LABEL kelompoknya di isi menu, dan ada di FOLDER_KORPUS', () => {
    expect(MODUL_FRONTEND.length).toBeGreaterThanOrEqual(4)
    for (const m of MODUL_FRONTEND) {
      expect(Object.values(FOLDER_KORPUS), m.nama).toContain(m.kelompok)
      expect(kelompok.find((k) => k.kode === m.nama)?.label, m.nama).toBe(m.kelompok)
      for (const b of m.menu) expect(b.kelompok, `${m.nama}/${b.modul}`).toBe(m.kelompok)
    }
  })

  it('LABEL kelompok = FOLDER_KORPUS (20 folder korpus)', () => {
    expect(new Set(kelompok.map((k) => k.label))).toEqual(new Set(Object.values(FOLDER_KORPUS)))
    expect(Object.values(FOLDER_KORPUS)).toHaveLength(20)
    expect(Object.values(FOLDER_KORPUS)).toContain('Treaty In')
    expect(Object.values(FOLDER_KORPUS)).toContain('Treaty In Adjustment')
  })
})

describe('GET /api/menu', () => {
  afterEach(() => {
    vi.unstubAllGlobals()
  })

  function jawab(badan: string, status = 200): void {
    vi.stubGlobal(
      'fetch',
      vi.fn(() => Promise.resolve(new Response(badan, { status }))),
    )
  }

  it('membaca pohon dari jalur yang dipasang cmd/api', async () => {
    jawab('{"golongan":[{"kode":"KLAIM","kelompok":[]}]}')
    await expect(ambilMenu()).resolves.toEqual({ golongan: [{ kode: 'KLAIM', kelompok: [] }] })
    expect(String(vi.mocked(fetch).mock.calls[0]?.[0])).toContain('/api/menu')
    const rakit = readFileSync(join(AKAR_APLIKASI, 'cmd', 'api', 'rakit.go'), 'utf8')
    expect(rakit).toContain('mux.HandleFunc("GET /api/menu"')
  })

  it('bentuk yang tak dikenal GAGAL - bukan menu kosong diam-diam', async () => {
    for (const badan of ['{}', '{"golongan":{}}', '{"golongan":[{"kode":1}]}']) {
      jawab(badan)
      await expect(ambilMenu(), badan).rejects.toThrow('GET /api/menu')
    }
  })

  it('galat backend diteruskan apa adanya', async () => {
    jawab('{"galat":"tabel M_NAV_MENU belum ada - migrasi 900 belum dijalankan (-migrate, oleh work owner)"}', 503)
    await expect(ambilMenu()).rejects.toMatchObject({ status: 503 })
  })
})

describe('sidebar dan palet dari GET /api/menu', () => {
  const SRC = __dirname
  const app = readFileSync(join(SRC, 'App.tsx'), 'utf8')
  const shell = readFileSync(join(AKAR_APLIKASI, 'inti', 'frontend', 'components', 'Shell.tsx'), 'utf8')

  it('App membacanya dan meneruskannya ke Shell', () => {
    expect(app).toContain('ambilMenu()')
    expect(app).toContain('menuTabel={menuTabel}')
  })

  it('Shell memotongnya dengan rute, mencatat baris tanpa rute, dan memakai daftar yang SAMA untuk palet', () => {
    expect(shell).toContain('susunMenu(menuTabel.menu, menu)')
    expect(shell).toContain('console.warn(')
    expect(shell).toContain('menu={tersusun?.entri ?? berandaSaja}')
  })

  it('kepala golongan, galat menu, dan menu kosong tampil di sidebar', () => {
    expect(shell).toContain('shell__golongan-judul')
    expect(shell).toContain('<Gagal galat={menuTabel.galat} />')
    // Tabel tanpa satu pun baris yang dapat tampil: DIKATAKAN, bukan sidebar
    // yang hanya berisi Beranda tanpa sebab.
    expect(shell).toContain('tersusun !== null && tersusun.golongan.length === 0')
    expect(shell).toContain('{KERANGKA.menuKosong}')
  })

  it("kelompok DIMIGRASI '0' dirender 'belum dimigrasi' menurut DIMIGRASI-nya", () => {
    expect(shell).toContain('{!k.dimigrasi ? (')
  })
})
