// Pembantu uji - HASIL BERSIH isi M_NAV_MENU dari berkas migrasi (900 + 901 +
// slot menu modul), dan menu `GET /api/menu` yang backend susun darinya.
//
// Menu datar, keputusan work owner 30-09-2026 (`PROMPT-MENU-DATAR-PER-GROUPMENU.md`):
// 901 membuang lima butir anak 900; yang tersisa satu baris per modul. Penjaga
// otoritatif atas bentuk SQL-nya adalah skema tiruan Go
// (`inti/backend/penjaga/menu_test.go`); yang di sini hanya membaca hasilnya
// supaya uji frontend membandingkan dengan menu yang SAMA dengan backend -
// bukan isi 900 saja.

import { existsSync, readdirSync, readFileSync } from 'node:fs'
import { join } from 'node:path'

import type { MenuTabel } from '../lib/daftarMenu'
import { AKAR_APLIKASI } from './sumber'

/** Satu baris modul M_NAV_MENU sesudah seluruh migrasi menu. */
export interface BarisMenuBersih {
  kode: string
  label: string
  golongan: string
  modul: string
  urutan: number
  dimigrasi: boolean
}

/** Satu berkas migrasi maju yang menyebut M_NAV_MENU. */
export interface BerkasMigrasiMenu {
  nama: string
  isi: string
}

/**
 * Urutan golongan backend - DIBACA dari `inti/backend/menu/menu.go`
 * (`var Golongan`), bukan salinan tangan (temuan /code-review).
 */
export const GOLONGAN_MENU: readonly string[] = (() => {
  const go = readFileSync(join(AKAR_APLIKASI, 'inti', 'backend', 'menu', 'menu.go'), 'utf8')
  const isi = /var Golongan = \[\]string\{([^}]*)\}/.exec(go)?.[1]
  if (isi === undefined) throw new Error('menuBersih: var Golongan tidak terbaca dari inti/backend/menu/menu.go')
  return [...isi.matchAll(/"([A-Z][A-Z ]*)"/g)].map((m) => m[1]!)
})()

/** Setiap migrasi maju yang menyebut M_NAV_MENU, urutan nama berkas = urutan pelari. */
export function berkasMenu(): BerkasMigrasiMenu[] {
  const folder = [
    join(AKAR_APLIKASI, 'inti', 'backend', 'migrations'),
    ...readdirSync(join(AKAR_APLIKASI, 'modul'), { withFileTypes: true })
      .filter((d) => d.isDirectory())
      .map((d) => join(AKAR_APLIKASI, 'modul', d.name, 'backend', 'migrations'))
      .filter((d) => existsSync(d)),
  ]
  return folder
    .flatMap((d) =>
      readdirSync(d)
        .filter((n) => n.endsWith('.sql') && !n.endsWith('_down.sql'))
        .map((n) => ({ nama: n, isi: readFileSync(join(d, n), 'utf8') })),
    )
    .filter((b) => b.isi.includes('M_NAV_MENU'))
    .sort((a, b) => (a.nama < b.nama ? -1 : a.nama > b.nama ? 1 : 0))
}

/** Pernyataan satu berkas - pemisah baris `/`, baris komentar `--` dibuang (pelari Go). */
function pernyataan(isi: string): string[] {
  const out: string[] = []
  let kini: string[] = []
  for (const baris of isi.split('\n')) {
    if (baris.trim() === '/') {
      if (kini.join('').trim() !== '') out.push(kini.join('\n').trim())
      kini = []
      continue
    }
    if (!baris.trimStart().startsWith('--')) kini.push(baris)
  }
  if (kini.join('').trim() !== '') out.push(kini.join('\n').trim())
  return out
}

const POLA_KELOMPOK =
  /^INSERT INTO \{skema\}\.M_NAV_MENU \([^)]*\)\s+SELECT \{skema\}\.SEQ_M_NAV_MENU\.NEXTVAL, NULL, '([^']+)', '([^']+)', '([A-Z]+)', '([^']+)', (\d+), '([01])' FROM DUAL/
// Baris modul DI LUAR korpus (`MODUL_LUAR_KORPUS`): bentuk datar sesudah 901, tanpa PARENT_ID (migrasi inti 906 dst.).
// Keabsahannya - hanya `modulLuarKorpus` - dijaga skema tiruan Go; di sini hanya diterapkan.
const POLA_KELOMPOK_DATAR =
  /^INSERT INTO \{skema\}\.M_NAV_MENU \(ID, KODE, LABEL, GROUPMENU, MODUL, URUTAN, DIMIGRASI\)\s+SELECT \{skema\}\.SEQ_M_NAV_MENU\.NEXTVAL, '([^']+)', '([^']+)', '([A-Z]+)', '([^']+)', (\d+), '([01])' FROM DUAL/
const POLA_BUTIR = /^INSERT INTO \{skema\}\.M_NAV_MENU \([^)]*\)\s+SELECT \{skema\}\.SEQ_M_NAV_MENU\.NEXTVAL, k\.ID, '([^']+)'/
const HAPUS_BUTIR = "EXECUTE IMMEDIATE 'DELETE FROM {skema}.M_NAV_MENU WHERE PARENT_ID IS NOT NULL'"
// Bentuk slot menu SESUDAH 901 saja (`WHERE KODE = '<modul>'`, berjangkar).
const POLA_UBAH = /^UPDATE \{skema\}\.M_NAV_MENU SET DIMIGRASI = '([01])', TGL_UBAH = SYSDATE\s+WHERE KODE = '([^']+)'$/
// Nama tampilan baris modul di slot menunya (keputusan work owner 03-10-2026). Keabsahannya - hanya baris
// `labelTampilDisetujui`, teks persis - dijaga `inti/backend/penjaga`; di sini hanya diterapkan.
// Baris modul pindah golongan (langkah inti 909 MASTER TREATY): GROUPMENU dan URUTAN sekaligus.
const POLA_UBAH_GOLONGAN =
  /^UPDATE \{skema\}\.M_NAV_MENU SET GROUPMENU = '([^']+)', URUTAN = (\d+), TGL_UBAH = SYSDATE\s+WHERE KODE = '([^']+)'$/
// CHECK GROUPMENU dibuat ulang (909) - tidak mengubah baris.
const TAMBAH_CEK_GOLONGAN = 'ALTER TABLE {skema}.M_NAV_MENU ADD CONSTRAINT CK_M_NAV_MENU_GROUPMENU CHECK'
const POLA_UBAH_LABEL = /^UPDATE \{skema\}\.M_NAV_MENU SET LABEL = '([^']+)', TGL_UBAH = SYSDATE\s+WHERE KODE = '([^']+)'$/
// Baris modul dipensiunkan (920 - Master Data dipecah delapan modul): DELETE berjangkar satu KODE.
const POLA_HAPUS_BARIS = /^DELETE FROM \{skema\}\.M_NAV_MENU WHERE KODE = '([^']+)'$/
// Hak akun atas menu (`M_LOGIN_GO_MENU`, 920 menyalin hak lalu membuangnya) - bukan baris menu, tidak mengubah hasil.
const HAK_MENU = /^(INSERT INTO|DELETE FROM) \{skema\}\.M_LOGIN_GO_MENU\b/

/** Hasil bersih: baris modul, butir yang tersisa, dan cacah INSERT yang terbaca. */
export interface MenuBersih {
  baris: BarisMenuBersih[]
  butir: string[]
  /** Cacah pernyataan INSERT di berkas menu, dan yang terbaca pola di atas. */
  insert: number
  insertTerbaca: number
  /** Pernyataan yang tidak dikenal - harus kosong (bukan diabaikan diam-diam). */
  takDikenal: string[]
}

/** Menerapkan berkas menu berurutan, seperti pelari. */
export function menuBersih(berkas: readonly BerkasMigrasiMenu[] = berkasMenu()): MenuBersih {
  const hasil: MenuBersih = { baris: [], butir: [], insert: 0, insertTerbaca: 0, takDikenal: [] }
  for (const b of berkas) {
    for (const p of pernyataan(b.isi)) {
      if (HAK_MENU.test(p)) continue
      if (p.startsWith('INSERT')) hasil.insert++
      const k = POLA_KELOMPOK.exec(p) ?? POLA_KELOMPOK_DATAR.exec(p)
      if (k !== null) {
        hasil.insertTerbaca++
        if (!hasil.baris.some((x) => x.kode === k[1])) {
          hasil.baris.push({ kode: k[1]!, label: k[2]!, golongan: k[3]!, modul: k[4]!, urutan: Number(k[5]), dimigrasi: k[6] === '1' })
        }
        continue
      }
      const bt = POLA_BUTIR.exec(p)
      if (bt !== null) {
        hasil.insertTerbaca++
        hasil.butir.push(bt[1]!)
        continue
      }
      if (p.startsWith('CREATE ')) continue
      // Blok berpelindung katalog 901: hanya DELETE butir yang mengubah baris.
      if (p.startsWith('DECLARE') && p.includes('EXECUTE IMMEDIATE')) {
        if (p.includes(HAPUS_BUTIR)) hasil.butir = []
        continue
      }
      if (p.startsWith(TAMBAH_CEK_GOLONGAN)) continue
      const g = POLA_UBAH_GOLONGAN.exec(p)
      if (g !== null) {
        for (const x of hasil.baris) {
          if (x.kode === g[3]) {
            x.golongan = g[1]!
            x.urutan = Number(g[2])
          }
        }
        continue
      }
      const u = POLA_UBAH.exec(p)
      if (u !== null) {
        for (const x of hasil.baris) if (x.kode === u[2]) x.dimigrasi = u[1] === '1'
        continue
      }
      const l = POLA_UBAH_LABEL.exec(p)
      if (l !== null) {
        for (const x of hasil.baris) if (x.kode === l[2]) x.label = l[1]!
        continue
      }
      const h = POLA_HAPUS_BARIS.exec(p)
      if (h !== null) {
        hasil.baris = hasil.baris.filter((x) => x.kode !== h[1])
        continue
      }
      hasil.takDikenal.push(`${b.nama}: ${p.slice(0, 80)}`)
    }
  }
  return hasil
}

/**
 * `GET /api/menu` seperti backend menyusunnya dari baris bersih (`menu.Susun`):
 * golongan menurut `GOLONGAN_MENU`, modul menurut URUTAN, modul dimigrasi di
 * luar `aktif` tidak dikirim (`null` = semua dipasang).
 */
export function menuTabelDariMigrasi(aktif: readonly string[] | null = null, baris = menuBersih().baris): MenuTabel {
  return {
    golongan: GOLONGAN_MENU.map((g) => ({
      kode: g,
      modul: baris
        .filter((b) => b.golongan === g && (!b.dimigrasi || aktif === null || aktif.includes(b.modul)))
        .sort((a, b) => a.urutan - b.urutan)
        .map((b) => ({ kode: b.kode, label: b.label, modul: b.modul, urutan: b.urutan, dimigrasi: b.dimigrasi })),
    })).filter((g) => g.modul.length > 0),
  }
}
