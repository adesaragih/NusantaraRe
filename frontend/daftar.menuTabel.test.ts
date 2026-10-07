import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { afterEach, describe, expect, it, vi } from 'vitest'

import { AKAR_APLIKASI } from '../inti/frontend/uji/sumber'
import { berkasMenu, menuBersih, type BarisMenuBersih } from '../inti/frontend/uji/menuBersih'
import { ambilMenu } from '../inti/frontend/klien'
import { FOLDER_KORPUS, LABEL_MENU, LABEL_TAMPIL, MODUL_LUAR_KORPUS } from './katalogKorpus'
import { MODUL_FRONTEND } from './daftar'

// Penjaga DUA ARAH: HASIL BERSIH M_NAV_MENU (900 + 901 + slot menu modul) ↔
// `frontend/daftar.ts` - menu datar, keputusan work owner 30-09-2026
// (`PROMPT-MENU-DATAR-PER-GROUPMENU.md` §4). Bahwa hanya 900, 901, dan slot
// menu yang boleh menyebut M_NAV_MENU dijaga `inti/backend/penjaga/rentang_test.go`.
//
// ⛔ Kenapa statik: sidebar dirakit dari `GET /api/menu` dan DIPOTONG dengan
// modul frontend - baris dimigrasi tanpa modul frontend tidak tampil, modul
// frontend tanpa baris juga tidak. Keduanya diam di layar (hanya satu baris
// konsol). Uji inilah yang membuatnya berbunyi sebelum sampai ke layar.
//
// Menggantikan uji "lima butir" (KODE butir ↔ rute, LABEL butir, induk butir).

/** Bentuk modul frontend yang dibandingkan - subset `ModulFrontend`. */
interface ModulUji {
  nama: string
  kelompok: string
  halaman: readonly string[]
  halamanAwal: string
}

/**
 * Selisih dua arah. Kosong = sepakat.
 *
 *   - setiap baris `DIMIGRASI = '1'` punya modul frontend dengan HALAMAN_AWAL
 *     yang benar-benar halamannya, dan nama modulnya = LABEL barisnya
 *   - setiap modul frontend terdaftar punya baris tabel, dan baris itu
 *     `DIMIGRASI = '1'` (modul berlayar tidak boleh tampil "belum dimigrasi")
 */
function selisihMenuModul(baris: readonly BarisMenuBersih[], modul: readonly ModulUji[]): string[] {
  const alasan: string[] = []
  for (const b of baris.filter((x) => x.dimigrasi)) {
    const m = modul.find((x) => x.nama === b.kode)
    if (m === undefined) {
      alasan.push(`baris ${b.kode} DIMIGRASI='1' tanpa modul frontend`)
      continue
    }
    if (!m.halaman.includes(m.halamanAwal)) alasan.push(`modul ${m.nama}: HALAMAN_AWAL ${m.halamanAwal} bukan halamannya`)
    if (m.kelompok !== b.label) alasan.push(`modul ${m.nama}: nama ${m.kelompok}, LABEL baris ${b.label}`)
  }
  for (const m of modul) {
    const b = baris.find((x) => x.kode === m.nama)
    if (b === undefined) alasan.push(`modul frontend ${m.nama} tanpa baris M_NAV_MENU`)
    else if (!b.dimigrasi) alasan.push(`modul frontend ${m.nama}: barisnya DIMIGRASI='0'`)
  }
  return alasan
}

const BERSIH = menuBersih()

describe('hasil bersih M_NAV_MENU ↔ daftar.ts, dua arah', () => {
  it('setiap INSERT terbaca penjaga ini', () => {
    // INSERT berbentuk lain adalah baris yang uji ini tidak lihat.
    expect(BERSIH.insertTerbaca).toBe(BERSIH.insert)
    // Dan setiap pernyataan lain pun dikenal - bukan diabaikan diam-diam.
    expect(BERSIH.takDikenal).toEqual([])
    expect(berkasMenu().map((b) => b.nama)).toEqual(expect.arrayContaining(['900_m_nav_menu.sql', '901_m_nav_menu_datar.sql']))
  })

  it('dua puluh baris modul, nol butir anak (901)', () => {
    // Dua puluh folder korpus + modul di luar korpus (Marketing Officer, migrasi inti 906 - keputusan work owner 03-10-2026).
    expect(BERSIH.baris).toHaveLength(20 + Object.keys(MODUL_LUAR_KORPUS).length)
    expect(BERSIH.butir).toEqual([])
    for (const b of BERSIH.baris) expect(b.kode, b.label).toBe(b.modul)
  })

  it('dua arah: baris DIMIGRASI ↔ modul frontend dengan HALAMAN_AWAL', () => {
    expect(MODUL_FRONTEND.length).toBeGreaterThanOrEqual(4)
    expect(selisihMenuModul(BERSIH.baris, MODUL_FRONTEND)).toEqual([])
  })

  it('uji gigit: baris dimigrasi tanpa modul frontend', () => {
    const tanpaSatu = MODUL_FRONTEND.filter((m) => m.nama !== 'premiumlistlife')
    expect(selisihMenuModul(BERSIH.baris, tanpaSatu)).toEqual(['baris premiumlistlife DIMIGRASI=\'1\' tanpa modul frontend'])
  })

  it('uji gigit: modul frontend tanpa baris, atau berbaris DIMIGRASI=0', () => {
    const tiruan = { nama: 'modultiruan', kelompok: 'Modul Tiruan', halaman: ['t'], halamanAwal: 't' }
    expect(selisihMenuModul(BERSIH.baris, [...MODUL_FRONTEND, tiruan])).toEqual(['modul frontend modultiruan tanpa baris M_NAV_MENU'])
    // Contoh modul berbaris DIMIGRASI='0' - claimprop (edmtreatyin menyala 06-10-2026, slot 970).
    const claimprop = { nama: 'claimprop', kelompok: 'Claim Prop', halaman: ['claimprop'], halamanAwal: 'claimprop' }
    expect(selisihMenuModul(BERSIH.baris, [...MODUL_FRONTEND, claimprop])).toEqual(["modul frontend claimprop: barisnya DIMIGRASI='0'"])
  })

  it('uji gigit: HALAMAN_AWAL di luar halaman modul, dan nama ≠ LABEL', () => {
    const rusak = MODUL_FRONTEND.map((m) => (m.nama === 'claimlife' ? { ...m, halamanAwal: 'bukan', kelompok: 'ClaimLife' } : m))
    expect(selisihMenuModul(BERSIH.baris, rusak)).toEqual([
      'modul claimlife: HALAMAN_AWAL bukan bukan halamannya',
      'modul claimlife: nama ClaimLife, LABEL baris Claim Life',
    ])
  })

  it('LABEL = FOLDER_KORPUS (20 folder korpus), kecuali nama tampilan LABEL_TAMPIL', () => {
    expect(new Set(BERSIH.baris.map((b) => b.label))).toEqual(new Set(Object.values(LABEL_MENU)))
    // Keputusan work owner 03-10-2026 - SATU-SATUNYA nama tampilan; pasangan Go `labelTampilDisetujui` memuat yang sama.
    expect(LABEL_TAMPIL).toEqual({ masterContractRetroLife: 'Contract Retro Life', masterProductNameLife: 'Product Name Life' })
    for (const [k, v] of Object.entries(LABEL_TAMPIL)) expect(FOLDER_KORPUS[k as keyof typeof FOLDER_KORPUS]).toBe(`Master ${v}`)
    expect(Object.values(FOLDER_KORPUS)).toHaveLength(20)
    expect(MODUL_LUAR_KORPUS).toEqual({
      marketingOfficer: 'Marketing Officer',
      companyDetail: 'Company Detail',
      accounts: 'Accounts',
      masterNation: 'Nation',
      masterProvince: 'Province',
      masterCity: 'City',
      masterDistrict: 'District',
      masterCzone: 'CZone',
      masterAccumulatedType: 'Accumulated Type',
      masterAccumulation: 'Accumulation',
      masterObjectItemType: 'Object Item Type',
      aggregate: 'Aggregate',
      bordereaux: 'Bordereaux',
      adjusterConsultant: 'Adjuster Consultant',
      treatyGroupOjk: 'Treaty Group OJK',
      treatyGroup: 'Treaty Group',
      businessGroup: 'Business Group',
      treatyExchangeYearly: 'Treaty Exchange Yearly',
      treatyDescription: 'Treaty Description',
      reinsuranceType: 'Reinsurance Type',
      riRateLife: 'R/I Rate Life',
    })
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

  it('membaca golongan → modul dari jalur yang dipasang cmd/api', async () => {
    const badan = '{"golongan":[{"kode":"KLAIM","modul":[{"kode":"claimlife","label":"Claim Life","modul":"claimlife","urutan":2,"dimigrasi":true}]}]}'
    jawab(badan)
    await expect(ambilMenu()).resolves.toEqual(JSON.parse(badan))
    expect(String(vi.mocked(fetch).mock.calls[0]?.[0])).toContain('/api/menu')
    const rakit = readFileSync(join(AKAR_APLIKASI, 'cmd', 'api', 'rakit.go'), 'utf8')
    expect(rakit).toContain('mux.HandleFunc("GET /api/menu"')
  })

  it('bentuk yang tak dikenal - termasuk pohon lama - GAGAL, bukan menu kosong diam-diam', async () => {
    for (const badan of ['{}', '{"golongan":{}}', '{"golongan":[{"kode":1}]}', '{"golongan":[{"kode":"KLAIM","kelompok":[]}]}']) {
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

  it('Shell memotongnya dengan modul frontend, mencatat baris tanpa modul, dan memakai daftar yang SAMA untuk palet', () => {
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

  it("modul DIMIGRASI '0' dirender nonaktif 'belum dimigrasi'", () => {
    expect(shell).toContain('{tujuan === null ? (')
    expect(shell).toContain('aria-disabled="true"')
  })
})
