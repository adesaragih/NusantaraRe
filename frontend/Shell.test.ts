// Penjaga Shell dan menu — F0.3.
//
// ⛔ Yang dijaga di sini adalah janji yang paling mudah dilanggar tanpa
// berbunyi: menu yang TIDAK ADA di sistem lama. Menambah satu butir menu
// terasa tidak berbahaya, dan itulah sebabnya ia harus menjadi kegagalan uji
// alih-alih keputusan sepi.

import { existsSync, readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import { AKAR_APLIKASI, folderKorpusBelumDimigrasi } from '../inti/frontend/uji/sumber'
import { MENU, MODUL_LAIN_TERLARANG } from '../inti/frontend/labels'
import { LABEL_MENU } from './katalogKorpus'
import { MENU_TCO } from '../modul/treatycontractout/frontend/labels'
import { ENTRI_MENU, MODUL_FRONTEND } from './daftar'
import { KELOMPOK_CLAIMLIFE } from '../modul/claimlife/frontend/menu'

const SUMBER = readFileSync(join(AKAR_APLIKASI, 'inti', 'frontend', 'components', 'Shell.tsx'), 'utf8')
const APP = readFileSync(join(__dirname, 'App.tsx'), 'utf8')
/** Sejak butir bg palet hidup di berkasnya sendiri. */
const PALET = readFileSync(join(AKAR_APLIKASI, 'inti', 'frontend', 'components', 'PaletMenu.tsx'), 'utf8')

const KORPUS = 'D:\\XML\\RNM_BRD\\Claim Life'
const adaKorpus = existsSync(KORPUS)

/** Buang komentar: prosa yang MENYEBUT sebuah nama tidak boleh tertuduh. */
function tanpaKomentar(teks: string): string {
  return teks
    .replace(/\/\*[\s\S]*?\*\//g, '')
    .split('\n')
    .filter((b) => !b.trimStart().startsWith('//') && !b.trimStart().startsWith('*'))
    .join('\n')
}

const kode = tanpaKomentar(SUMBER)

/** Nama kelompok yang Shell render, dibaca dari labelnya. */
const KELOMPOK_SIDEBAR = Object.values(LABEL_MENU)

describe('menu hanya yang berbukti korpus', () => {
  it('kelompok sidebar TEPAT dua puluh', () => {
    // ⛔ DIRALAT 28-09-2026 (butir bg). Uji ini dulu berbunyi "butir
    // sidebar TEPAT dua" dan membaca literal ${b}const BUTIR${b} di Shell.tsx.
    // Sejak bg, butirnya diturunkan dari ${b}ENTRI_MENU${b} dan kelompoknya
    // tujuh belas - nama folder korpus. Yang dijaga BERPINDAH, bukan hilang.
    // ⛔ DELAPAN BELAS sejak 28-09-2026: kelompok `Treaty Contract Out`
    // ditambahkan sesi modul itu (folder korpus 20; ralat §8 PROMPT-EKSEKUSI).
    // ⛔ DUA PULUH sejak brief menu M_NAV_MENU (30-09-2026): isi awal tabel
    // memuat satu kelompok per folder korpus, Treaty In dan Treaty In
    // Adjustment ikut (`frontend/daftar.menuTabel.test.ts` menjaga LABEL-nya).
    // Dua puluh folder korpus + Marketing Officer (modul di luar korpus, migrasi inti 906).
    // + Company Detail (modul di luar korpus, migrasi inti 907).
    // + Accounts (modul di luar korpus, migrasi inti 908).
    expect(KELOMPOK_SIDEBAR).toHaveLength(24)
    // Dan seluruhnya disebut di Shell, supaya tidak ada kelompok yang
    // terdaftar di label tetapi tidak dirender.
    for (const nama of KELOMPOK_SIDEBAR) {
      expect(Object.values(LABEL_MENU)).toContain(nama)
    }
  })

  it('satu entri per modul terdaftar, berlabel nama folder korpus - nol butir navigasi', () => {
    // ⛔ MENU DATAR (keputusan work owner 30-09-2026): "1 modul 1 menu" - butir
    // inbox, register, premiumlist, komite, tco-tahun DICABUT; tombol modul
    // membuka halaman awalnya. Beranda TIDAK dihitung - ia kerangka aplikasi.
    // Menggantikan "kelima butir menu lama tetap, dan seluruhnya berbukti".
    // Kelola User (01-10-2026) juga TIDAK dihitung - ia menu aplikasi, bukan
    // modul korpus (`daftar.kelolauser.test.ts`).
    const modul = ENTRI_MENU.filter((e) => e.modul !== 'beranda' && e.modul !== 'kelolauser')
    // Lima sejak tiket 03 Treaty Contract Out. Sempat tujuh (tiket 04, 08);
    // tco5 [keputusan work owner 29-09-2026]: kelompok Treaty Contract Out
    // SATU butir "Treaty Contract Out" - ReinsType dan Description popup form
    // kontrak (InputTreatyContract b20778, b22196), bukan menu.
    //
    // ⛔ DIUBAH struktur tim satu folder per modul (30-09-2026): dulu
    // `toHaveLength(5)` atas SELURUH menu - angka di berkas milik tim inti
    // yang harus disunting setiap butir menu baru, dan dua modul yang
    // menambah butir bersamaan berkonflik di baris itu. Yang dijaga tidak
    // hilang, ia BERPINDAH: butir baru hanya lahir dari baris M_NAV_MENU di
    // slot menu modulnya (`frontend/daftar.menuTabel.test.ts` menagih kedua
    // arah), dan berkas slot `9*` di folder migrasi modul tetap ditinjau tim
    // inti lewat pull request - menu tetap keputusan yang ditinjau, bukan
    // keputusan sepi.
    const label = modul.map((e) => e.label)
    expect(modul).toHaveLength(MODUL_FRONTEND.length)
    for (const l of label) expect(Object.values(LABEL_MENU), l).toContain(l)
    // Label butir navigasi lama tidak tampil lagi sebagai menu.
    for (const lama of ['Inbox Claim Life', 'Register', 'PremiumList', 'Inbox Komite']) {
      expect(label, lama).not.toContain(lama)
    }
    expect(label).not.toContain(MENU_TCO.inboxTreatyContractReinsType)
    expect(label).not.toContain(MENU_TCO.inboxTreatyContractDescription)
    expect(ENTRI_MENU.filter((e) => e.kelompok === MENU_TCO.kelompok)).toHaveLength(1)
    // Halaman yang dibuka DARI DALAM modul tidak pernah menjadi entri menu.
    expect(label).toHaveLength(new Set(modul.map((e) => e.pemilik)).size)
  })

  it('modul yang belum dimigrasi berdiri sebagai tombol NONAKTIF, tanpa halaman', () => {
    // ⚠️ Berdiri, bukan disembunyikan. Aplikasi yang menampilkan empat
    // modul dari dua puluh tampak lengkap padahal tidak.
    //
    // Struktur tim satu folder per modul (30-09-2026): dulu "enam belas" -
    // kini daftarnya pernyataan `Status` di MODUL.md setiap modul.
    const berbutir = new Set(ENTRI_MENU.map((e) => e.kelompok))
    const kosong = KELOMPOK_SIDEBAR.filter((n) => !berbutir.has(n))
    expect([...kosong].sort()).toEqual(folderKorpusBelumDimigrasi())
    expect(kosong.length).toBeGreaterThan(0)
    expect(kode).toContain('KETERANGAN_BELUM_DIMIGRASI')
    expect(kode).toContain('aria-disabled="true"')
  })

  it.each(MODUL_LAIN_TERLARANG)('tidak ada butir menu bernama %s', (nama) => {
    // ⛔ Sebelas modul lain dibuka lewat menu HANYA bila korpusnya
    // membuktikan menunya ada. Belum ada yang membuktikannya.
    //
    // ⚠️ Yang diperiksa LABEL MENUnya, bukan seluruh berkas. Ronde
    // pertama memeriksa seluruh teks dan menuduh `ShellProps` karena
    // memuat "Prop". Penjaga yang menuduh nama tipe akan membuat
    // orang mengganti nama tipenya - bukan memperbaiki menunya.
    for (const label of Object.values(MENU)) {
      expect(label).not.toContain(nama)
    }
  })

  it('nol BelumTersedia di menu', () => {
    // `BelumTersedia` sah DI DALAM halaman; sebagai BUTIR MENU ia menjanjikan
    // layar yang sistem lama tidak punya.
    expect(kode).not.toContain('BelumTersedia')
  })

  it('Detail bukan butir menu — di Pega ia flow action dari dalam kasus', () => {
    // ${b}outstanding${b} dan ${b}detail${b} ada di tipe ${b}Halaman${b} tetapi TIDAK di
    // ${b}ENTRI_MENU${b}: keduanya dibuka DARI DALAM kasus, bukan dari navigasi.
    const modul = ENTRI_MENU.map((e) => e.modul as string)
    expect(modul).not.toContain('detail')
    expect(modul).not.toContain('outstanding')
  })
})

describe('bukti XML label menu', () => {
  it.runIf(adaKorpus)('MENU.register VERBATIM Register_Flow.xml baris 155', () => {
    const alur = readFileSync(join(KORPUS, 'Flow', 'Register_Flow.xml'), 'utf8')
    const baris155 = alur.split('\n')[154] ?? ''
    expect(baris155).toContain(`<pyLabel>${MENU.register}</pyLabel>`)
  })

  it.runIf(adaKorpus)('nama kelompok dari pyWorkTypeName baris 270', () => {
    const alur = readFileSync(join(KORPUS, 'Flow', 'Register_Flow.xml'), 'utf8')
    const baris270 = alur.split('\n')[269] ?? ''
    // Korpus menulisnya tanpa spasi (`ClaimLife`); menu menampilkannya
    // dengan spasi, dan itu satu-satunya penyimpangan yang diizinkan.
    expect(baris270).toContain('<pyWorkTypeName>ClaimLife</pyWorkTypeName>')
    expect(KELOMPOK_CLAIMLIFE.replace(' ', '')).toBe('ClaimLife')
  })

  it('Inbox ditandai tidak ada di korpus', () => {
    const labels = readFileSync(join(AKAR_APLIKASI, 'inti', 'frontend', 'labels.ts'), 'utf8')
    const blok = labels.slice(labels.indexOf('export const MENU'))
    // Klaimnya harus berdiri di komentar tepat di atas MENU.
    const kepala = labels.slice(labels.lastIndexOf('/**', labels.indexOf('export const MENU')), labels.indexOf('export const MENU'))
    expect(kepala).toContain('[tidak ada di korpus]')
    expect(blok).toContain('inbox:')
  })
})

describe('perilaku shell yang ditiru referensi', () => {
  it.each([
    ['Ctrl/Cmd + K membuka palet', "e.key.toLowerCase() === 'k'"],
    ['Esc menutup', "e.key === 'Escape'"],
    ['klik luar menutup menu profil', 'mousedown'],
    ['sidebar dapat dilipat', 'setTerlipat'],
    ['laci ponsel', 'setLaciBuka'],
  ])('%s', (_nama, tanda) => {
    expect(kode).toContain(tanda)
  })

  it.each([
    ['fokus kembali sesudah palet ditutup', 'fokusSebelum'],
    ['panah tidak memindahkan karet', 'preventDefault'],
    ['sorotan melingkar', '% hasil.length'],
    ['daftar kosong menyebut kuerinya', 'palet__kosong'],
  ])('palet: %s', (_nama, tanda) => {
    // ⛔ BERPINDAH BERKAS 28-09-2026 (butir bg). Palet dulu hidup
    // sebagai fungsi di dalam Shell.tsx; sejak bg ia komponen tersendiri
    // yang membaca `lib/daftarMenu.ts`. Penjaganya ikut pindah - yang
    // dijaga sama, tempatnya yang berbeda.
    expect(PALET).toContain(tanda)
  })

  it('palet membaca daftarMenu, BUKAN DOM sidebar', () => {
    // ⛔ Sejak bg empat belas kelompok terlipat, dan `KelompokMenu`
    // melepas anak kelompok yang terlipat dari DOM. Palet yang membaca
    // DOM akan menjawab "tidak ada" untuk menu yang ADA.
    expect(PALET).toContain("from '../lib/daftarMenu'")
    expect(PALET).not.toContain('querySelector')
    expect(PALET).not.toContain('getElementsBy')
  })

  it('PagarGalat membungkus isi halaman', () => {
    // Satu galat render di satu halaman tidak boleh memutihkan seluruh layar.
    expect(kode).toContain('<PagarGalat>')
  })
})

describe('nol halaman dirender di luar Shell', () => {
  it('App merender halaman hanya sebagai anak Shell', () => {
    const app = tanpaKomentar(APP)
    // ⛔ Login sungguhan (M_LOGIN_GO, keputusan work owner 01-10-2026): yang
    // boleh di luar Shell HANYA layar login dan ganti sandi - tidak satu pun
    // halaman modul.
    const luar = app.slice(0, app.indexOf('<Shell'))
    expect(luar).toContain('<Login')
    expect(luar).toContain('<GantiSandi')
    expect(luar).not.toContain('<RegisterKlaim')
    expect(luar).not.toContain('<KlaimLife')
    expect(luar).not.toContain('<InboxClaimLife')
    // Refactor bentuk B paket 7: halaman dirender rute tiap modul, dan rute
    // itu hanya dipasang DI DALAM Shell.
    expect(luar).not.toContain('<m.Rute')
    expect(app.slice(app.indexOf('<Shell'))).toContain('<m.Rute')
  })

  it('login sungguhan: sesi dari /api/auth/saya, stub tetap dari env', () => {
    // Keputusan work owner 01-10-2026 menggantikan perintah 27-09-2026 (tanpa
    // form login): kini ada akun dan sandi sungguhan (M_LOGIN_GO). Mode stub
    // tetap untuk pengembangan, tanpa layar login.
    expect(APP).toContain('ambilSesiSaya()')
    expect(APP).toContain('bolehMasukStub()')
    expect(APP).toContain('pelakuStub()')
    // 401 di tengah pemakaian kembali ke layar login.
    expect(APP).toContain('PERISTIWA_SESI_BERAKHIR')
  })

  it('tombol Keluar hanya di login sungguhan — mode stub tidak punya yang dikeluari', () => {
    expect(kode).toContain('onKeluar ? (')
    expect(APP).toMatch(/onKeluar=\{\s*stub\s*\?\s*undefined/)
    expect(kode).not.toContain('sesi.hapus')
  })

  it('bilah sesi sementara F0.2 sudah dibuang', () => {
    expect(APP).not.toContain('bilah-sesi')
  })
})

// Menu profil TANPA daftar peran - permintaan work owner 01-10-2026 ("buang
// aja, ga perlu"): nama akun, lalu Ganti sandi dan Keluar.
describe('menu profil', () => {
  it('tanpa daftar peran', () => {
    expect(SUMBER).not.toContain('shell__profil-peran')
    expect(SUMBER).not.toContain('masuk.peran.map(')
    expect(SUMBER).toContain('<p className="shell__profil-akun">{masuk.akunID}</p>')
  })

  // Pemicu menu profil di topbar: nama saja - "tulisan dibawah namanya dihapus"
  // (work owner 01-10-2026, dulu "Admin Klaim Jiwa +2").
  it('pemicu topbar tanpa keterangan peran di bawah nama', () => {
    expect(SUMBER).not.toContain('sebutanPeran')
    expect(SUMBER).not.toContain('PERAN_ID')
    expect(SUMBER).toContain('<strong>{masuk.nama ?? masuk.akunID}</strong>\n                </span>')
  })
})
