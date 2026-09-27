// Penjaga Shell dan menu — F0.3.
//
// ⛔ Yang dijaga di sini adalah janji yang paling mudah dilanggar tanpa
// berbunyi: menu yang TIDAK ADA di sistem lama. Menambah satu butir menu
// terasa tidak berbahaya, dan itulah sebabnya ia harus menjadi kegagalan uji
// alih-alih keputusan sepi.

import { existsSync, readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import { MENU, MODUL_LAIN_TERLARANG } from '../assets/labels'

const SUMBER = readFileSync(join(__dirname, 'Shell.tsx'), 'utf8')
const APP = readFileSync(join(__dirname, '..', 'App.tsx'), 'utf8')

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

describe('menu hanya yang berbukti korpus', () => {
  it('butir sidebar TEPAT dua', () => {
    // Daftarnya literal di Shell.tsx; cacahnya dikunci supaya butir ketiga
    // harus melewati uji ini lebih dulu.
    const daftar = kode.slice(kode.indexOf('const BUTIR'), kode.indexOf('export interface ShellProps'))
    expect(daftar.match(/halaman: '/g) ?? []).toHaveLength(2)
    expect(daftar).toContain('MENU.inbox')
    expect(daftar).toContain('MENU.register')
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
    const daftar = kode.slice(kode.indexOf('const BUTIR'), kode.indexOf('export interface ShellProps'))
    expect(daftar).not.toContain("'detail'")
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
    expect(MENU.kelompokClaimLife.replace(' ', '')).toBe('ClaimLife')
  })

  it('Inbox ditandai tidak ada di korpus', () => {
    const labels = readFileSync(join(__dirname, '..', 'assets', 'labels.ts'), 'utf8')
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
    ['fokus kembali sesudah palet ditutup', 'fokusSebelum'],
  ])('%s', (_nama, tanda) => {
    expect(kode).toContain(tanda)
  })

  it('PagarGalat membungkus isi halaman', () => {
    // Satu galat render di satu halaman tidak boleh memutihkan seluruh layar.
    expect(kode).toContain('<PagarGalat>')
  })
})

describe('nol halaman dirender di luar Shell', () => {
  it('App merender halaman hanya sebagai anak Shell', () => {
    const app = tanpaKomentar(APP)
    // ⛔ Sejak F0.6 satu-satunya yang boleh di luar Shell adalah
    // pernyataan "identitas tidak ada". Form login DIBUANG.
    const luar = app.slice(0, app.indexOf('<Shell'))
    expect(luar).toContain('<BelumTersedia')
    expect(luar).not.toContain('<RegisterKlaim')
    expect(luar).not.toContain('<KlaimLife')
    expect(luar).not.toContain('<InboxClaimLife')
  })

  it('nol form login — identitas datang dari env', () => {
    // `[perintah work owner 27-09-2026]`. Layar masuk tanpa sandi bukan
    // autentikasi; mempertahankannya hanya menambah langkah yang tidak
    // memutuskan apa pun.
    expect(APP).not.toContain('Masuk')
    expect(APP).toContain('pelakuStub()')
  })

  it('nol tombol Keluar di Shell — tanpa masuk tidak ada keluar', () => {
    expect(kode).not.toContain('onKeluar')
    expect(kode).not.toContain('sesi.hapus')
  })

  it('bilah sesi sementara F0.2 sudah dibuang', () => {
    expect(APP).not.toContain('bilah-sesi')
  })
})
