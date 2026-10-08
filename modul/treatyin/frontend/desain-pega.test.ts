// Penjaga ronde "layar disamakan dengan dokumen desain" (5 Oktober 2026).
//
// ⛔ Sumber setiap angka di berkas ini adalah salah satu dari 43 tangkapan
// layar `D:\NUSANTARA RE APP\design-treaty-in-gambar\`, dan nomor gambarnya
// ditulis di tiap uji. Uji yang tidak dapat menyebut gambarnya tidak punya
// tempat di sini.

import { readdirSync, readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import {
  KOLOM_AKUMULASI,
  KOLOM_EGNPI,
  KOLOM_LIHAT_BERKAS,
  KOLOM_POLIS_PRODUKSI,
  KOLOM_PORTOFOLIO,
  PORTOFOLIO,
  KOLOM_RETENSI,
  JENIS_EGNPI,
  JENIS_POLIS_PRODUKSI,
  JENIS_RETENSI,
  LAMPIRAN,
  NAMA_KATEGORI_LAMPIRAN_NON_PROP,
  NAMA_KATEGORI_LAMPIRAN_PROP,
  POLIS_PRODUKSI,
  SUB_TAB_SHARE,
  TAB_NON_PROPORSIONAL,
  TAB_PROPORSIONAL,
  golongan,
  JENIS_LIMITS,
  JENIS_SHARE,
  JENIS_RNM_SHARE,
  desimalPadan,
  GRID_BERTOMBOL_TAMBAH,
  CO_INS_SCALE,
  KOLOM_COIN_SCALE,
  JENIS_COIN_SCALE,
  TOMBOL_TAMBAH_BELUM_BERGRID,
  TOTAL_SHARE,
  GRID_TOTAL_RNM_SHARE,
  KOLOM_KIND_OF_TREATY_SHARE,
  INFO_SUBMIT,
} from './labels'
import { padankanDesimal, selAngka } from './pages/FormKontrakTreatyIn'

const AKAR = join(__dirname)
// ⭐ Layar ini DIPECAH 5 Oktober 2026: halaman + `components/`. `FORM`
// karena itu membaca KEDUANYA, supaya tiap pernyataan di bawah tetap
// menanyakan hal yang sama — "apakah kode layar modul ini memuat ini" —
// tanpa satu pun disunting. Itu bukti pemindahannya murni.
const SELA = String.fromCharCode(10)
const FORM =
  readFileSync(join(AKAR, 'pages', 'FormKontrakTreatyIn.tsx'), 'utf8') +
  readdirSync(join(AKAR, 'components'))
    .filter((f) => f.endsWith('.tsx') || f.endsWith('.ts'))
    .map((f) => readFileSync(join(AKAR, 'components', f), 'utf8'))
    .join(SELA)
const CSS = readFileSync(join(AKAR, 'treatyin.css'), 'utf8')

describe('§2 panel Existing Policy for Master ID', () => {
  it('judul dan keempat kolomnya disalin dari gambar 01/26', () => {
    expect(POLIS_PRODUKSI.judul).toBe('Existing Policy for Master ID')
    expect([...KOLOM_POLIS_PRODUKSI]).toEqual(['Policy No', 'Pega ID', 'Quarter', 'Quarter Year'])
    expect(POLIS_PRODUKSI.tanpaBaris).toBe('No items')
  })

  it('keempat kolomnya TEKS — pengenal tidak pernah diformat', () => {
    expect(JENIS_POLIS_PRODUKSI).toHaveLength(KOLOM_POLIS_PRODUKSI.length)
    expect(JENIS_POLIS_PRODUKSI.every((j) => j === 'teks')).toBe(true)
  })

  it('panelnya dirender di kepala layar, kedua cabang', () => {
    // ⛔ Di LUAR cabang `tabTampil` — ia bukan isi tab mana pun.
    expect(FORM).toContain('<PanelPolisProduksi baris={warisan?.polisProduksi ?? []} />')
    const iPanel = FORM.indexOf('<PanelPolisProduksi baris=')
    const iTab = FORM.indexOf("tabTampil === 'Reporting Period'")
    expect(iPanel).toBeGreaterThan(0)
    expect(iPanel).toBeLessThan(iTab)
  })

  it('lebarnya dibatasi dan turun penuh di layar sempit', () => {
    expect(CSS).toContain('.treatyin .trin__polis')
    expect(CSS).toContain('max-width: 34rem')
    expect(CSS).toContain('@media (max-width: 52rem)')
  })
})

describe('§0 judul kolom disalin dari gambarnya, bukan dari kunci dokumen', () => {
  it('Portfolio — gambar 02', () => {
    expect([...KOLOM_PORTOFOLIO]).toEqual([
      'Portfolio Type',
      'Premium / Loss Type',
      'Description',
    ])
  })

  it('Accumulation — gambar 19', () => {
    expect([...KOLOM_AKUMULASI]).toEqual([
      'Period',
      'Reporting Date',
      'Submission Days',
      'Submission Due',
    ])
  })

  it('EGNPI — urut layar gambar 29, bukan abjad', () => {
    // ⭐ Enam kolom pertama persis urutan gambarnya.
    expect(KOLOM_EGNPI.slice(0, 6)).toEqual([
      'Treaty Group',
      'As Date',
      'Proportion %',
      'Currency',
      'Amount',
      'Amount in IDR',
    ])
    // ⛔ Dua kolom yang tidak terlihat di gambar TIDAK dihapus.
    expect(KOLOM_EGNPI).toContain('Class of Business')
    expect(KOLOM_EGNPI).toContain('Note')
    expect(JENIS_EGNPI).toHaveLength(KOLOM_EGNPI.length)
  })

  it('Maximum Retention — urut layar gambar 26', () => {
    expect(KOLOM_RETENSI.slice(0, 3)).toEqual(['Treaty Group', 'Currency', 'Amount'])
    expect(KOLOM_RETENSI).toContain('Class of Business')
    expect(KOLOM_RETENSI).toContain('Note')
    expect(JENIS_RETENSI).toHaveLength(KOLOM_RETENSI.length)
    // Kolom uangnya pindah posisi, dan golongannya ikut pindah.
    expect(golongan(JENIS_RETENSI[KOLOM_RETENSI.indexOf('Amount')]!)).toBe('uang')
  })
})

describe('§6 RNM Share adalah SUB-TAB, dan nol tab dihapus', () => {
  it('strip sub-tab ada di dalam tab Share — gambar 16/17/34', () => {
    expect([...SUB_TAB_SHARE]).toEqual(['RNM Share'])
    expect(FORM).toContain('function SubTabShare(')
    expect(FORM).toContain('<StripTab tab={daftar} aktif={tampil} onPilih={setSub} />')
  })

  it('⛔ RNM Share TETAP di daftar tab — §6 melarang menghapusnya', () => {
    expect(TAB_NON_PROPORSIONAL).toContain('RNM Share')
    expect(TAB_NON_PROPORSIONAL).toHaveLength(12)
    expect(TAB_PROPORSIONAL).toHaveLength(11)
  })

  it('cabang tab RNM Share masih merender sesuatu', () => {
    // Tab yang dapat dipilih dan tidak merender apa pun adalah layar rusak.
    expect(FORM).toContain("tabTampil === 'RNM Share'")
    expect(FORM).toContain('<PanelRnmShare')
  })

  it('gridnya SATU bentuk, dipakai dua tempat', () => {
    // ⛔ Dua salinan grid yang sama berarti satu yang basi.
    expect(FORM.split('function PanelRnmShare(')).toHaveLength(2)
    expect(FORM.split('KOLOM_RNM_SHARE}').length - 1).toBe(1)
  })
})

describe('§4 permintaan perubahan — View File menjadi pop-up', () => {
  it('memakai Modal dari inti, bukan modal kedua', () => {
    expect(FORM).toContain('Modal,')
    expect(FORM).toContain('judul={LAMPIRAN.judulLihatBerkas}')
    expect(FORM).not.toContain('function ModalLihatBerkas(')
  })

  it('isinya DUA kolom, persis gambar 25', () => {
    expect([...KOLOM_LIHAT_BERKAS]).toEqual(['File Name', 'Type'])
  })

  // ⭐ RALAT 8 Oktober 2026 — Upload HIDUP (keputusan pemakai: unggahan masuk
  // `M_ATTACHMENTTREATY_2`), bersyarat `TreatyIn.ViewState !='1' ||
  // TreatyIn.RevisionState='1'` dan kontrak ber-ID.
  it('tombol View hidup, tombol Upload hidup bersyarat ekspor', () => {
    expect(LAMPIRAN.lihatBerkas).toBe('View')
    expect(LAMPIRAN.unggah).toBe('Upload')
    expect(FORM).toMatch(/\{bisaUnggah && \(\s*<button/)
    expect(FORM).toContain("disabled={sibuk || idKontrak === '' || k.kode === ''}")
    expect(FORM).toContain('judul={LAMPIRAN.judulUnggah}')
    expect(LAMPIRAN.judulUnggah).toBe('ASM Attach Content')
    expect(LAMPIRAN.lampirkan).toBe('Attach')
    expect(FORM).toContain('bisaUnggah={bisaUbah || revisi}')
    const iLihat = FORM.indexOf('{LAMPIRAN.lihatBerkas}')
    expect(FORM.slice(iLihat - 220, iLihat)).toContain('setBerkasDilihat')
  })

  it('penyaringnya memakai KODE kategori, bukan namanya', () => {
    // Nama kategori yang belum dipastikan tampil sebagai kodenya; menyaring
    // menurut yang tampil berarti menyaring menurut dua hal berbeda.
    expect(FORM).toContain('b.kodeKategori === kategori.kode')
  })
})

describe('§8 katalog kategori lampiran BERCABANG — gambar 24 lawan 42', () => {
  it('sebelas nama, dan tepat satu berbeda', () => {
    expect(NAMA_KATEGORI_LAMPIRAN_PROP).toHaveLength(11)
    expect(NAMA_KATEGORI_LAMPIRAN_NON_PROP).toHaveLength(11)
    const beda = NAMA_KATEGORI_LAMPIRAN_PROP.filter(
      (n) => !(NAMA_KATEGORI_LAMPIRAN_NON_PROP as readonly string[]).includes(n),
    )
    expect(beda).toEqual(['Pega Proportional Calculation /Perhitungan Pega Proportional'])
  })

  it('⛔ daftar nama TIDAK menutup pertanyaan kode↔nama', () => {
    // Urutan layarnya alfabetis, dan ronde 4 Oktober 2026 melarang keras
    // menyimpulkan kode dari urutan abjad.
    expect(LAMPIRAN.namaBelumBerumah).toContain('PERTANYAAN-TERBUKA-KODE-KATEGORI-LAMPIRAN.md')
  })
})

describe('§2 Treaty Year — tahun polos, TIDAK diformat', () => {
  it('nilainya lewat apa adanya, nol pemformat', () => {
    // Gambar 01: `Treaty Year 2025`. `selAngka` akan memberi `2.025`.
    const i = FORM.indexOf('value={tahunTreaty}')
    expect(i).toBeGreaterThan(0)
    expect(FORM.slice(i - 200, i + 60)).not.toContain('selAngka')
  })

  // ⭐ RALAT 6 Oktober 2026 — `Treaty Year` TURUNAN, tetapi hanya pada satu
  // peristiwa, dan ekspor yang menyebut peristiwanya.
  //
  // `DataTransform/TreatyInSetTreatyYear.xml` menyetelnya pada event `change`
  // milik `TreatyIn.Commencement` (`pyEvent change` + `postValue` + `refresh`,
  // terbaca di `Section/TreatyInNONProportional.xml` sesudah
  // `pyIncludedRuleXML` bersarang dibuang).
  //
  // ⛔ Keputusan LAMA — "dibaca dari kolom, jangan dihitung" — TIDAK dicabut,
  // ia dipersempit: pada SAAT MEMBUKA kontrak kolomnya tetap yang dipakai,
  // sebab 5 dari 1.854 kontrak sengaja bertahun berbeda dari tahun mulainya,
  // dan menghitung ulang saat membuka akan menimpa kelimanya.
  it('⭐ Treaty Year disetel ulang pada perubahan Commencement, sesuai ekspor', () => {
    expect(FORM).toContain('const ubahMulai =')
    expect(FORM).toContain('setTahunTreaty(tahunDariMulai(v))')
    expect(FORM).toContain('setBerakhir(akhirSetahunSesudah(v))')
    expect(FORM).toContain('onChange={ubahMulai}')
  })

  it('dugaan 05/1 sebagai turunan dicabut di komentarnya', () => {
    expect(FORM).toContain('tahun polos')
    expect(FORM).not.toContain('PERTANYAAN TERBUKA: tangkapan layar Pega memperlihatkan `05/1`')
  })
})

describe('§24 desimal PER KOLOM, nol di ekor dipertahankan', () => {
  it('nol ditambahkan sampai presisi kolomnya', () => {
    // Gambar 01: `Value to IDR` berbunyi `1,00` dan `16.000,00`.
    expect(selAngka(['uang', 2], '1')).toBe('1,00')
    expect(selAngka(['uang', 2], '16000')).toBe('16.000,00')
    expect(selAngka(['uang', 2], '16000.5')).toBe('16.000,50')
  })

  it('⛔ NOL desimal berarti nol — bukan "tidak diketahui"', () => {
    // Gambar 29: `Amount in IDR` berbunyi `137.849.315.068`, tanpa koma.
    expect(selAngka(['uang', 0], '137849315068')).toBe('137.849.315.068')
    expect(selAngka(['uang', 0], '137849315068.00')).toBe('137.849.315.068')
  })

  it('⛔ yang SUDAH lebih panjang tidak dipotong — memotong itu membulatkan', () => {
    // ⚠️ DUA batas yang berbeda, dan keduanya berlaku berurutan:
    //
    //   `formatNumber` memotong di `DESIMAL_UANG` (4) — batas ATAS uang,
    //                  aturan lama yang §24 tidak sentuh;
    //   `padankanDesimal` menambah nol sampai presisi kolom (2) — dan
    //                  karena yang tiba sudah 4 desimal, ia tidak berbuat
    //                  apa-apa.
    //
    // Jadi hasilnya 4 desimal, bukan 2: pemadanan MENAMBAH, tidak pernah
    // memotong. Memotong berarti membulatkan, dan membulatkan menghilangkan
    // digit berarti.
    expect(selAngka(['uang', 2], '1500.12345')).toBe('1.500,1235')
    expect(padankanDesimal('1.500,12345', 2)).toBe('1.500,12345')
    expect(padankanDesimal('1.500,1235', 2)).toBe('1.500,1235')
  })

  it('tanpa presisi kolom, aturan ronde 69 tetap berlaku — nol di ekor dibuang', () => {
    // ⛔ Inilah yang menjaga §24 TIDAK bocor ke kolom yang tidak terbaca.
    expect(selAngka('uang', '1')).toBe('1')
    expect(selAngka('uang', '2484250')).toBe('2.484.250')
    expect(padankanDesimal('2.484.250', null)).toBe('2.484.250')
  })

  it('persen: nol masuk SEBELUM tanda %', () => {
    expect(selAngka(['persenShare', 2], '25')).toBe('25,00%')
    expect(padankanDesimal('100%', 4)).toBe('100,0000%')
  })

  it('⛔ teks bukan-angka lewat apa adanya — nol padanan', () => {
    for (const t of ['', '>=30% up to < 50%', 'n/a', '-', 'IDR']) {
      expect(padankanDesimal(t, 2)).toBe(t)
    }
  })

  it('kolom yang desimalnya terbaca membawa ANGKA, bukan hanya golongan', () => {
    // Gambar 29: dua kolom di baris yang SAMA, presisi berbeda.
    expect(desimalPadan(JENIS_EGNPI[KOLOM_EGNPI.indexOf('Amount')]!)).toBe(2)
    expect(desimalPadan(JENIS_EGNPI[KOLOM_EGNPI.indexOf('Amount in IDR')]!)).toBe(0)
    // Gambar 26: `Amount` di grid Maximum Retention nol desimal.
    expect(desimalPadan(JENIS_RETENSI[KOLOM_RETENSI.indexOf('Amount')]!)).toBe(0)
  })

  it('⚠️ kolom yang TIDAK terbaca di gambar tetap null — bukan ditebak', () => {
    // Nol dari 43 gambar memperlihatkan presisi kolom-kolom ini.
    for (const j of JENIS_LIMITS) expect(desimalPadan(j)).toBeNull()
    for (const j of JENIS_SHARE) expect(desimalPadan(j)).toBeNull()
    for (const j of JENIS_RNM_SHARE) expect(desimalPadan(j)).toBeNull()
  })

  it('⛔ format.ts dipanggil, tidak ditulis ulang — pemadanan di lapis modul', () => {
    expect(FORM).toContain("from '../../../../inti/frontend/lib/format'")
    expect(FORM).toContain('export function padankanDesimal(')
    // Pemadanan TIDAK menguraikan angka: nol pemformat kedua.
    expect(FORM).not.toMatch(/function formatNumber|function formatPersen/)
  })
})

describe('§1 tab Limits — dua cabang, dua susunan dari ekspor', () => {
  it('puncaknya BERBEDA per cabang — Kind of Treaty lawan Layers', () => {
    expect(FORM).toContain("judul={KOLOM_KIND_OF_TREATY[0]?.label ?? ''}")
    expect(FORM).toContain('{LIMITS_NP.layers}')
    expect(FORM).toContain('<TabLimitsNonProp')
    expect(FORM).toContain('<TabLimitsProp')
  })

  it('⛔ PohonLimits dicabut — grid layer dari tabel pendaratan, bukan pengelompokan teks', () => {
    // Pohon lama mengelompokkan layer menurut `.Layer`, sehingga dua layer
    // bernomor sama (layer dan sublayer) tergabung menjadi satu.
    expect(FORM).not.toContain('function PohonLimits(')
    expect(FORM).not.toContain('const k = nonProp ? b.layer : b.jenisTreaty')
  })

  it('⭐ tombol tab Limits HIDUP di mode Edit — rumus di services', () => {
    expect(FORM).toContain('{LIMITS_NP.tambahLayer}')
    expect(FORM).toContain('{LIMITS_NP.tambahGrup}')
    expect(FORM).toContain('{LIMITS_NP.perbaruiTotal}')
    expect(FORM).toContain('hitungLimitNP(')
  })
})

describe('§1 tombol Add — tiap grid yang gambarnya memperlihatkannya', () => {
  it('daftarnya berkunci NAMA TAB, supaya tab berikutnya tidak lupa', () => {
    // ⛔ RALAT 6 Oktober 2026 — `Portfolio` MASUK daftar.
    //
    // Daftar ini dibaca dari gambar, dan gambar `02` tangkapan layar
    // mode-BACA. Ekspor `Section/TreatyInTabsProportional.xml` menunjukkan
    // tombolnya ADA: sel `Add` 110 dan `Delete` 115, `pyVisible` = `OTHER`
    // dengan `pyCondition` = `TreatyIn.IsEditData!='1'`.
    expect([...GRID_BERTOMBOL_TAMBAH].sort()).toEqual(['Accumulation', 'EGNPI', 'Portfolio'])
  })

  it('⛔ MEKANISMENYA DIGANTI 6 Oktober 2026 — tombol mengikuti MODE', () => {
    // ⭐ Permintaan pemakai yang lebih baru mengalahkan daftar ini:
    //   "apabila di klik edit maka tiap function itu bisa digunakan semuanya
    //    seperti add pada tiap table, kalau view baru tidak bisa"
    //
    // Jadi `Add` tidak lagi dipasang per tab dari daftar gambar; ia muncul
    // pada SETIAP grid di mode `ubah` dan pada NOL grid di mode `lihat` —
    // dan ekspor mendukungnya: tombol-tombol itu dijaga
    // `TreatyIn.ViewState !='1'`.
    //
    // ⛔ Daftarnya TIDAK dihapus: ia tetap catatan grid mana yang gambar
    // mode-BACA perlihatkan bertombol, dan itu bukti terpisah dari mode.
    expect(FORM).toContain("mode === 'ubah'")
    expect(FORM).toContain('bisaTambah')
  })

  it('⛔ dua grid yang gambarnya TIDAK punya Add tetap tanpa tombol', () => {
    // Memasangnya di sana membuat layar berbeda dari layar lama.
    //
    // ⚠️ Dua, bukan tiga. `Portfolio` KELUAR dari daftar ini 6 Oktober
    // 2026 sebab ekspor membuktikan tombolnya ada — gambarnya yang
    // mode-BACA. Kedua yang tersisa BELUM diperiksa terhadap ekspor dengan
    // cara yang sama, dan itu dikatakan apa adanya: mereka tinggal di sini
    // atas dasar gambar saja, bukan atas dasar bukti bahwa tombolnya nihil.
    for (const tab of ['Co-Ins Scale', 'Maximum Retention']) {
      expect(GRID_BERTOMBOL_TAMBAH).not.toContain(tab)
    }
  })

  it('yang gambarnya punya tombol tetapi GRIDNYA belum ada: didaftarkan', () => {
    expect(TOMBOL_TAMBAH_BELUM_BERGRID).toHaveLength(3)
    expect(TOMBOL_TAMBAH_BELUM_BERGRID.join(' ')).toContain('Reinsurer')
    expect(TOMBOL_TAMBAH_BELUM_BERGRID.join(' ')).toContain('Update Value')
  })
})

describe('§1.2 pemilih bercari — papan tik, pola ditiru bukan diimpor', () => {
  const DD = readFileSync(join(AKAR, 'components', 'DropdownWarisan.tsx'), 'utf8')

  it('⛔ pola DITIRU, kode TIDAK diimpor dari modul lain', () => {
    // ⚠️ Namanya MEMANG disebut — di komentar, sebagai asal polanya. Yang
    // dilarang IMPORNYA, dan itu yang diuji: nol baris `import` menyebut
    // modul mana pun selain `inti` dan modul ini sendiri.
    expect(DD).toContain('DropdownCari.tsx')
    for (const ln of DD.split(String.fromCharCode(10))) {
      if (ln.trimStart().startsWith('import ')) {
        // Yang DIIZINKAN: paket (`react`), `inti`, dan berkas modul ini
        // sendiri. Selain itu impor lintas modul, dan itu dilarang.
        const m = /from '([^']+)'/.exec(ln)
        if (m !== null) {
          const jalur = m[1]!
          const sah =
            !jalur.startsWith('.') ||
            jalur.includes('/inti/frontend/') ||
            /^\.\.?\/[^.]/.test(jalur)
          expect(`${jalur} sah=${String(sah)}`).toContain('sah=true')
        }
      }
    }
  })

  it('kelima tombol papan tik ditangani', () => {
    for (const k of ['ArrowDown', 'ArrowUp', 'PageDown', 'PageUp', 'Enter', 'Escape', 'Tab']) {
      expect(DD).toContain(`'${k}'`)
    }
  })

  it('⛔ geseran DIJEPIT, tidak melingkar', () => {
    // Melingkar membuat yang menekan terus tidak pernah tahu ia di ujung.
    expect(DD).toContain('Math.min(cocok.length - 1, Math.max(0, a + langkah))')
  })

  it('baris aktif disebut kepada pembaca layar dan digulirkan', () => {
    expect(DD).toContain('aria-activedescendant')
    expect(DD).toContain('scrollIntoView')
    expect(DD).toContain('aria-controls')
  })

  it('Escape menutup TANPA memilih', () => {
    const i = DD.indexOf("e.key === 'Escape'")
    const blok = DD.slice(i, DD.indexOf("e.key === 'Tab'", i))
    expect(blok).toContain('setBuka(false)')
    expect(blok).toContain('setAktif(-1)')
    // ⛔ Nol pemilihan di cabang itu.
    expect(blok).not.toContain('pilih(')
  })

  it('⚠️ ketikan baru MENGGANTI baris aktif — baris pertama daftar BARU tersorot', () => {
    // Baris ke-3 daftar LAMA bukan baris ke-3 daftar baru. Sejak 7 Oktober
    // 2026 yang tersorot baris pertama (paling mirip), seperti autocomplete
    // Pega; kotak kosong → nol baris tersorot.
    const i = DD.indexOf('setKetik(e.target.value)')
    expect(DD.slice(i, i + 800)).toContain("setAktif(e.target.value.trim() === '' ? -1 : 0)")
  })

  it('⭐ Enter tanpa baris tersorot memilih yang paling mirip — bukan diam lalu kembali kosong', () => {
    const i = DD.indexOf("} else if (e.key === 'Enter')")
    const blok = DD.slice(i, DD.indexOf("e.key === 'Escape'", i))
    expect(blok).toContain("cocok[aktif] ?? (q !== '' ? cocok[0] : undefined)")
    expect(blok).toContain('e.preventDefault()')
  })

  it('⭐ nama KEMBAR tetap dibedakan pengenalnya — beda yang milik kita', () => {
    expect(DD).toContain('b.kembar ? `${b.nama} — ${b.id}` : b.nama')
  })

  it('⚠️ beda yang DISENGAJA dari modul contoh dicatat: dapat diketik', () => {
    expect(DD).toContain('MENOLAK ketikan bebas')
    expect(DD).toContain('HARUS')
  })

  it('baris aktif punya tanda yang BUKAN hover', () => {
    expect(DD).toContain('trin__ketikpilih-baris--aktif')
    expect(CSS).toContain('.treatyin .trin__ketikpilih-baris--aktif')
  })
})

// ---------------------------------------------------------------------
// §25 tab Portfolio — bentuknya dari EKSPOR, pilihannya dari PENGUKURAN
// ---------------------------------------------------------------------
//
// ⛔ Berkas ini biasanya menuntut tiap uji menyebut GAMBARNYA. Blok ini
// pengecualian yang dinyatakan: gambar `02` tangkapan layar mode-BACA,
// dan mode-baca tidak memperlihatkan kendali apa pun — nol dropdown, nol
// tombol. Jadi sumbernya di sini `Section/TreatyInTabsProportional.xml`,
// disebut sel demi sel, dan gambar `02` tetap menjadi sumber JUDULNYA.
describe('§25 tab Portfolio — dua daftar pilihan, satu area teks', () => {
  const TP = readFileSync(join(AKAR, 'components', 'TabPortofolio.tsx'), 'utf8')

  it('⛔ kolom 1 memakai TypePortfolio, kolom 2 memakai Type — TIDAK tertukar', () => {
    // Pasangan sel di ekspor: 107↔112 (`.TypePortfolio`), 108↔113 (`.Type`).
    // Saksi kedua: judul kolom 2 MENYEBUT kedua nilainya, dan yang bernilai
    // `Premium`/`Loss` adalah kunci `Type`.
    expect([...KOLOM_PORTOFOLIO]).toEqual(['Portfolio Type', 'Premium / Loss Type', 'Description'])
    // ⭐ 7 Oktober 2026: baris kini EJAAN PEGA (penampung halaman, `halaman.tsx`).
    const k1 = TP.indexOf('TypePortfolio: b.jenisPortfolio')
    const k2 = TP.indexOf('Type: b.jenis')
    expect(k1).toBeGreaterThan(-1)
    expect(k2).toBeGreaterThan(k1)
    // ⛔ Dan `TypePortfolio` BENAR-BENAR dirender di kolom pertama.
    expect(TP.indexOf('nilai={b.TypePortfolio}')).toBeLessThan(TP.indexOf('nilai={b.Type}'))
  })

  it('pilihannya DIUKUR, bukan dikarang — dan hanya itu nilainya', () => {
    // Sapuan 6 Oktober 2026 atas seluruh 1.855 dokumen `M_TREATY_IN`:
    // 1.925 elemen `Portfolio[]` di 845 dokumen, nol gagal urai, nol kosong.
    //   TypePortfolio  Withdrawal 1.578 · Assumption 347
    //   Type           Premium    1.032 · Loss       893
    expect([...PORTOFOLIO.opsiArah]).toEqual(['Withdrawal', 'Assumption'])
    expect([...PORTOFOLIO.opsiJenis]).toEqual(['Premium', 'Loss'])
  })

  it('⛔ keduanya <select>, bukan kotak teks', () => {
    expect(TP).toContain('<select')
    expect(TP).toContain('opsi={PORTOFOLIO.opsiArah}')
    expect(TP).toContain('opsi={PORTOFOLIO.opsiJenis}')
    // ⛔ Nol `<input type="text">` tersisa di tab ini — itulah cacat yang
    // ronde ini tutup.
    expect(TP).not.toContain("type=\"text\"")
  })

  it('Description area teks — sel 114 ber-`pyFormat` pxTextArea', () => {
    expect(TP).toContain('<textarea')
  })

  it('⚠️ nilai di luar daftar TIDAK dijatuhkan diam-diam', () => {
    // Daftarnya terukur dari yang TERPAKAI, bukan dibaca dari aturan
    // properti (yang tidak ikut diekspor). Menjatuhkan nilai asing akan
    // mengubah kontrak yang sekadar dibuka.
    expect(TP).toContain('const asing =')
  })

  it('Add dan Delete mengikuti satu sakelar: IsEditData', () => {
    // Sel 110 dan 115: `pyCondition` = `TreatyIn.IsEditData!='1'`.
    // Sel 112·113·114: `pyReadOnlyCondition` = `TreatyIn.IsEditData= 1`.
    expect(TP).toContain('IsEditData')
    expect(TP).toContain('const bisaUbah = mode === ')
    expect(TP).toContain('{bisaUbah && (')
  })

  it('baris baru lahir KOSONG bertiga — seperti Activity-nya', () => {
    // `Activity/TreatyInPropAdd.xml`, langkah 1 (`param.Type=="portfolio"`):
    // satu-satunya penetapan adalah
    // `TreatyIn.Portfolio(<APPEND>).Description = ""`.
    expect(TP).toContain("{ TypePortfolio: '', Type: '', Description: '' }")
  })

  it('⛔ EDMMaterialType DICATAT, tidak dikarang', () => {
    // `pyDisabledWhen` = `TreatyIn.EDMMaterialType = 2` ada di keempat sel,
    // tetapi medannya belum terbaca layar ini. Menyalakannya berarti
    // menebak nilainya.
    expect(TP).toContain('EDMMaterialType')
  })
})

// ---------------------------------------------------------------------
// §26 FUNGSI YANG APLIKASI PUNYA TETAPI PEGA TIDAK — 6 Oktober 2026
// ---------------------------------------------------------------------
//
// ⛔ Setiap pencabutan di bawah DIADU ke ekspor lebih dulu, bukan diambil
// dari laporan. Bukti per sel ada di komentar tiap berkas; uji ini menjaga
// supaya pencabutannya tidak dibalikkan diam-diam oleh suntingan berikutnya.
describe('§26 Add/Delete yang Pega tidak punya — dicabut', () => {
  const SHARE = readFileSync(join(AKAR, 'components', 'TabShare.tsx'), 'utf8')
  const HAL = readFileSync(join(AKAR, 'pages', 'FormKontrakTreatyIn.tsx'), 'utf8')

  it('⛔ kedua grid Share per layer tanpa Add/Delete — sel 382/398 `1=2`', () => {
    // `Section/TreatyInTabsNonProportional.xml`, grid `TreatyIn.Share`
    // (tab `RNM Share < Title < Share`): `Add` sel 382 dan `Delete` sel 398,
    // keduanya `pyCondition` = `1=2`.
    expect((SHARE.match(/bisaTambah=\{false\}/g) ?? []).length).toBe(2)
    expect(SHARE).toContain('1=2')
  })

  it('⛔ Installment tanpa Add/Delete — nol sel tombol di 52 seksi', () => {
    // ⭐ 7 Oktober 2026: tab berumus `TabAngsuran.tsx` — barisnya LAHIR dari
    // Update Value; nol tombol tambah/hapus baris.
    expect(HAL).toContain('<TabAngsuran')
    const ANG = readFileSync(join(AKAR, 'components', 'TabAngsuran.tsx'), 'utf8')
    expect(ANG).not.toMatch(/TombolTambah|TombolHapus|>Add<|>Delete</)
    expect(ANG).toContain("nilai('update')")
  })

  it('⭐ grid yang Add/Delete-nya HIDUP di Pega TIDAK ikut dicabut', () => {
    // Kontrol: tanpa ini, mencabut seluruh Add/Delete akan meloloskan dua
    // uji di atas. Kelimanya hidup di ekspor —
    //   Retention sel 67/72 · EGNPI · Accumulation sel 319/325 ·
    //   Portfolio sel 110/115 · Limits sel 151/154.
    // ⭐ 7 Oktober 2026: Accumulation kini tab BERUMUS sendiri
    // (`TabAkumulasi.tsx`) — Add = DataTransform `TreatyInAddAccumulation`,
    // Delete per baris; keduanya hidup di mode ubah.
    expect(HAL).toContain('<TabAkumulasi')
    const AK = readFileSync(join(AKAR, 'components', 'TabAkumulasi.tsx'), 'utf8')
    expect(AK).toContain('barisAkumulasiBaru(periode, isi.length)')
    expect(AK).toContain('{AKUMULASI.hapus}')
  })

  // ⛔ EGNPI PINDAH RUMAH 6 Oktober 2026, dan pernyataannya ikut pindah.
  //
  // Ia tidak lagi `kolom={KOLOM_EGNPI}` di halaman: ekspor memberinya tab
  // sendiri berikut dua tombol berumus, jadi `TabEgnpi.tsx` yang memegang
  // Add/Delete-nya. Yang dijaga TETAP SAMA — keduanya hidup — hanya
  // tempat bertanyanya yang berubah.
  //
  // ⚠️ Bentuk lama uji ini lulus selama EGNPI berupa grid umum; begitu
  // tabnya lahir ia menunjuk ke tempat yang sudah kosong. Dicatat supaya
  // ronde berikutnya tidak mengembalikan EGNPI menjadi grid agar lulus.
  // ⛔ Maximum Retention PINDAH RUMAH 6 Oktober 2026, sebab yang sama
  // dengan EGNPI: ekspor memberinya tab berumus, bukan grid umum.
  it('⭐ Add/Delete Maximum Retention hidup di tab sendiri', () => {
    const RT = readFileSync(join(AKAR, 'components', 'TabRetensi.tsx'), 'utf8')
    expect(HAL).not.toContain('kolom={KOLOM_RETENSI}')
    expect(RT).toContain('{RETENSI.tambah}')
    expect(RT).toContain('{RETENSI.hapus}')
    expect(RT).not.toContain('bisaTambah={false}')
  })

  it('⭐ Add/Delete EGNPI hidup di tab sendiri', () => {
    const EG = readFileSync(join(AKAR, 'components', 'TabEgnpi.tsx'), 'utf8')
    expect(HAL).not.toContain('kolom={KOLOM_EGNPI}')
    expect(EG).toContain('{EGNPI.tambah}')
    expect(EG).toContain('{EGNPI.hapus}')
    expect(EG).not.toContain('bisaTambah={false}')
  })
})

describe('§27 Co-Ins Scale — bentuk Pega (gambar 18 + tangkapan layar mode ubah)', () => {
  const CO = readFileSync(join(AKAR, 'components', 'TabCoInsScale.tsx'), 'utf8')
  const KEPALA = CO.slice(CO.indexOf('<thead>'), CO.indexOf('</thead>'))
  const BADAN = CO.slice(CO.indexOf('<tbody>'), CO.indexOf('</tbody>'))

  it('⛔ grid TETAP dirender, dan tambah/hapus-nya TETAP — sel 262/266 hidup', () => {
    // ⛔ RALAT DI RONDE YANG SAMA. Grid ini sempat dicabut atas dasar
    // `pyContainerVisibleWhen = 1=2`; wadahnya `pyIsVisibilityOption =
    // ALWAYS` menimpanya, dan gambar 18 memperlihatkan gridnya TAMPIL.
    // Kedua tombolnya ber-`pyLabel` kosong — sapuan berdasarkan label
    // melewatkannya; yang menemukannya API `addRow`/`deleteRow`.
    expect(CO).toContain('<table')
    expect(CO).toContain('CO_INS_SCALE.tambah')
    expect(CO).toContain('CO_INS_SCALE.hapus')
  })

  it('⛔ TANPA judul di atas grid — wadahnya `NOHEADER`; tetap di dalam KARTU', () => {
    expect(CO).not.toContain('<Panel')
    expect(CO).not.toContain('<TabGridWarisan')
    // Kartu tanpa judul — tanpanya isi tab menempel ke kartu Attachment.
    expect(CO).toContain('<section className="panel trin__coin"')
  })

  it('tombol tambah DUDUK di sel kepala ketiga (sel 262 `pyCellHeader`), hapus per baris', () => {
    expect(KEPALA).toContain('CO_INS_SCALE.tambah')
    expect(KEPALA.match(/<th /g)?.length).toBe(3)
    expect(BADAN).toContain('CO_INS_SCALE.hapus')
    // Bukan lagi deret tombol + petunjuk di atas tabel.
    expect(CO).not.toContain('GRID_TAMBAH')
  })

  it('teks tombol dari Pega yang berjalan, bukan dari XML — `pyLabel` keduanya kosong', () => {
    expect(CO_INS_SCALE.tambah).toBe('Tambah')
    expect(CO_INS_SCALE.hapus).toBe('Hapus')
  })

  it('kosong berbunyi "No items" — `GridNoResultsOnLoad`, tanpa kartu kosong buatan kita', () => {
    expect(CO_INS_SCALE.kosong).toBe('No items')
    expect(CO).toContain('CO_INS_SCALE.kosong')
    expect(CO).not.toContain('<Kosong')
  })

  it('dua kolom, bukan empat — `.CoInShare` dan `.PctLimit` saja; `%` di kanan', () => {
    expect([...KOLOM_COIN_SCALE]).toEqual(['Co-Insurance Share', '% Treaty Limit'])
    expect(JENIS_COIN_SCALE.length).toBe(KOLOM_COIN_SCALE.length)
    expect(CO_INS_SCALE.simbolPersen).toBe('%')
  })

  it('grid DI ATAS, kedua medan DI BAWAH dengan label di KIRI', () => {
    expect(CO.indexOf('<table')).toBeLessThan(CO.indexOf('CO_INS_SCALE.nonGroup'))
    expect(CO).toContain('trin__kolom trin__coin-medan')
  })

  it('label disalin dari `pyLabelFieldValue` sel 277/278, lengkap dengan `Panel`', () => {
    expect(CO_INS_SCALE.nonGroup).toBe('Max Co-Insurance Panel (Non Group)')
    expect(CO_INS_SCALE.group).toBe('Max Co-Insurance Panel (Group)')
  })

  it('⛔ nol catatan buatan kita di layar — Pega tidak punya', () => {
    // Catatan "belum terjangkau" dicabut atas permintaan pemakai. Celahnya
    // (173 · 87 dari 1.855 dokumen, nol kolom pendaratan) tetap tercatat di
    // kepala `TabCoInsScale.tsx`, dan kekosongan itu TIDAK disebut "tidak
    // ada di dokumen".
    expect(CO).not.toContain('belumTerjangkau')
    expect(CO).not.toContain('MedanTakAda')
    expect(CO).not.toContain('takAdaDiWarisan')
    expect(CO).toContain('173')
  })

  it("dapat diubah HANYA di mode Edit — pyReadOnlyCondition `ViewState ='1'` / `IsEditData='1'`", () => {
    expect(CO).toContain("const bisaUbah = mode === 'ubah'")
    expect(CO).toContain('readOnly={!bisaUbah}')
  })

  it('CSS: kolom data `pyWidth` 206 · 206, kolom tombol selebar isinya — nol gulir di grid nyaris kosong', () => {
    expect(CSS).toContain('.treatyin .trin__coin-bungkus')
    expect(CSS).toMatch(/\.trin__coin \.trin__tabel th:nth-child\(2\) \{\s*width: 206px/)
    expect(CSS).toMatch(/\.trin__coin \.trin__tabel td:last-child \{\s*width: auto;\s*white-space: nowrap/)
    // ⛔ Bukan lagi 40/40/20 % dari 513px — kolom tombol lebih sempit dari tombolnya.
    expect(CSS).not.toContain('max-width: 513px')
  })
})


describe('§28 Choose Ceding / Source of Business — NAMA DAN PENGENAL', () => {
  const HAL = readFileSync(join(AKAR, 'pages', 'FormKontrakTreatyIn.tsx'), 'utf8')

  it('⛔ memilih mengisi KEDUANYA — `TreatyInSetReinsured` langkah 1.1/1.2, 2.1/2.2', () => {
    // Dahulu hanya `setCedant(b.nama)`: memilih ulang meninggalkan pengenal
    // LAMA di samping nama BARU.
    const c = HAL.indexOf('setCedant(b.nama)')
    expect(c).toBeGreaterThan(-1)
    expect(HAL.slice(c, c + 80)).toContain('setIdCedant(b.id)')
    const s = HAL.indexOf('setAsalBisnis(b.nama)')
    expect(s).toBeGreaterThan(-1)
    expect(HAL.slice(s, s + 80)).toContain('setIdAsalBisnis(b.id)')
  })

  it('pengenal awal dari kolom TREATY_IN — 0 dari 1.854 bernama tanpa pengenal', () => {
    expect(HAL).toContain('setIdCedant(k.idCedant)')
    expect(HAL).toContain('setIdAsalBisnis(k.idAsalBisnis)')
  })

  it('⛔ pengenal TERSEMBUNYI — Pega hanya menampilkan namanya (sel 23 · 28)', () => {
    expect(HAL).toContain('<input type="hidden" name="cedingId" value={idCedant} />')
    expect(HAL).toContain('<input type="hidden" name="leadingReinsSourceId" value={idAsalBisnis} />')
  })
})

describe('§29 LimitCalculation tab Limits Prop — dari Activity, diukur ulang', () => {
  const LP = readFileSync(join(AKAR, 'components', 'TabLimitsProp.tsx'), 'utf8')

  it('⭐ QS% dan Lines memicu rumus dengan parameter persis sel 37 dan 42', () => {
    // `DetailLimits` sel 37 `.QSPct` → (qs, '', true); sel 42 `.Surplus` →
    // (surplus, '', true). Keduanya `autocalculate = true` — satu-satunya
    // pemicu yang tidak ambigu.
    expect(LP).toContain("hitung('qs')")
    expect(LP).toContain("hitung('surplus')")
    expect(LP).toContain("jalankanLC(d, mode, '', true)")
  })

  it('⛔ dipicu saat LEPAS FOKUS, bukan di tiap ketukan', () => {
    // Pega: peristiwa `change`. Tiap ketukan akan membagi limit dengan `4`
    // sebelum `40` selesai diketik.
    const i = LP.indexOf("hitung('qs')")
    expect(LP.slice(Math.max(0, i - 120), i)).toContain('onBlur')
  })

  it('⛔ hanya medan yang Activity TULIS yang digabung kembali', () => {
    expect(LP).toContain('DITULIS_LIMIT_CALCULATION[mode]')
  })

  it('⛔ FetchQSfromMaster TIDAK dibangun dari ronde ini', () => {
    // Ia hidup lewat sel 3 `.TreatyGroupID` (`pyVisible` ALWAYS menimpa
    // `1=2`), tetapi bergantung dua RD induk `PROPORTIONALARRG` dan layar
    // Share per pohon yang belum ada. Tidak dibuat setengah.
    expect(LP).not.toContain('FetchQSfromMaster(')
  })
})

describe('§30 Event Limits Non-Prop — SATU set per kontrak, dari properti akar', () => {
  const EV = readFileSync(join(AKAR, 'components', 'TabEventLimits.tsx'), 'utf8')

  it('⛔ tidak lagi diulang per layer', () => {
    expect(EV).not.toContain('layer.map(')
    expect(EV).not.toContain('BarisLayerWarisan')
  })

  it('keempat baris, urut ekspor, berkunci properti AKAR', () => {
    for (const k of ['RSMDLimit', 'Earthquake', 'FloodJab', 'FloodNation']) {
      expect(EV).toContain(`nilai: '${k}'`)
    }
  })

  it('⛔ kekosongannya DIKATAKAN "belum terjangkau" — 49 kontrak punya nilainya', () => {
    expect(EV).toContain('EVENT_LIMITS.belumTerjangkau')
    expect(EV).not.toContain('takAdaDiWarisan')
  })
})

// ===========================================================================
// §3 Tab `Share` cabang PROPORSIONAL — bentuknya dari ekspor
// ===========================================================================
//
// `Section/TreatyInShareProp.xml` menyebut seluruhnya apa adanya; tangkapan
// layar Pega milik pemilik proses 6 Oktober 2026 memperlihatkan bentuk yang
// sama. Bentuk sebelumnya di aplikasi: satu grid datar berjudul NAMA KOLOM
// BASIS DATA (`LAYER`, `CESSIONPCT`, …) — bukan bentuk ini sama sekali.
describe('§3 tab Share proporsional', () => {
  const SRC = readFileSync(join(__dirname, 'components', 'TabShareProp.tsx'), 'utf8')

  it('panel `Total Share` punya Refresh dan ketiga medannya', () => {
    expect(TOTAL_SHARE.judul).toBe('Total Share')
    expect(TOTAL_SHARE.segarkan).toBe('Refresh')
    expect(TOTAL_SHARE.persenRnmShare).toBe('% RNM Share')
    expect(TOTAL_SHARE.persenBrokerage).toBe('% Brokerage')
    expect(TOTAL_SHARE.opsi).toBe('Option')
  })

  it('ketiga grid total berjudul persis seperti ekspor', () => {
    expect([...GRID_TOTAL_RNM_SHARE]).toEqual([
      'Total Share RNM Limit',
      'Total Value Spreading OR',
      'Total Value Spreading R/I',
    ])
    expect([...KOLOM_KIND_OF_TREATY_SHARE]).toEqual(['Kind of Treaty'])
  })

  // ⭐ EKSPORNYA DATANG 7 Oktober 2026 — rule Property `OptionLimit`,
  // kelas `ASM-FW-GISFW-Int-TREATY_IN`, Table type `Prompt List`:
  //   1 → `Of Cession to R/I`   2 → `Of 100% Limit`
  //
  // ⛔ Bentuk sebelumnya memaku label `2` sebagai `"2"` — kode apa adanya,
  // bukan karangan — dan MENUNGGUNYA TERNYATA BENAR: pasangan yang masuk
  // akal bagi `Of Cession to R/I` adalah sesuatu tentang cession, sedangkan
  // yang sebenarnya `Of 100% Limit`. Tebakan apa pun akan meleset.
  it('⭐ kedua label Option dari rule Property, bukan tebakan', () => {
    expect(SRC).toContain("{ value: '1', label: 'Of Cession to R/I' }")
    expect(SRC).toContain("{ value: '2', label: 'Of 100% Limit' }")
    // ⛔ Dan nol sisa penampung lama.
    expect(SRC).not.toContain("label: '2'")
  })

  // ⚠️ Kosong di ketiga grid total berarti SUMBERNYA belum ada, bukan
  // kontraknya yang kosong. Keduanya terlihat sama di layar.
  // ⛔ DIBALIK 8 Oktober 2026 — petunjuknya DICABUT, bukan diperbaiki.
  //
  // Pemilik proses: *"hapus semua komen komen yg ada seperti ini di
  // aplikasi"*. Kalimat lama menyebut `tabel pendaratan` — istilah internal
  // yang nol artinya bagi pemakai, dan membocorkan bentuk basis data ke
  // layar. Grid kosong kini berhenti pada `No items`.
  it('⛔ kosongnya ketiga total TIDAK lagi dijelaskan di layar', () => {
    expect(TOTAL_SHARE.petunjukTotal).toBe('')
  })
})

// ===========================================================================
// §4 Tab `Information & Submit` — FORM, bukan grid riwayat
// ===========================================================================
//
// ⛔ Cacat yang diperbaiki 6 Oktober 2026: tab ini merender grid riwayat
// (Date · Operator · Approved · Suggest). Di Pega ia FORM.
//
// ⚠️ Dan yang dirender bukan sekadar salah — ia SALINAN: riwayat sudah punya
// panelnya sendiri di kaki layar, membaca tabel yang sama. Layar menampilkan
// daftar yang sama dua kali, dan tab yang seharusnya tempat MENGIRIM justru
// tempat membaca.
describe('§4 tab Information & Submit', () => {
  const SRC = readFileSync(join(__dirname, 'components', 'TabInfoSubmit.tsx'), 'utf8')

  it('medannya dua Text area, persis `TreatyInfoSubmit.xml`', () => {
    expect(INFO_SUBMIT.infoTambahan).toBe('Additional Information')
    expect(INFO_SUBMIT.komentar).toBe('Comment')
    expect((SRC.match(/<Area\b/g) ?? []).length).toBe(2)
  })

  it('tombolnya `Submit` dan `Decline offer`', () => {
    expect(INFO_SUBMIT.kirim).toBe('Submit')
    expect(INFO_SUBMIT.tolak).toBe('Decline offer')
  })

  // ⛔ `pyCondition` kedua tombol `TreatyIn.ViewState != '1'` — tombol yang
  // Pega sembunyikan di mode lihat tidak boleh muncul di sini.
  // ⭐ 7 Oktober 2026: KECUALI Submit revisi (cell 21, `RevisionState='1'`
  // saja) — `revisi-submit.test.tsx` merendernya di kedua mode.
  it("⛔ kedua tombol tidak dirender di mode lihat", () => {
    expect(SRC).toContain("const bisaUbah = mode === 'ubah'")
    expect(SRC).toContain('const viewState1 = !bisaUbah || revisi')
    expect(SRC).toContain('{(!viewState1 || revisi) && (')
    expect(SRC).toContain('{!viewState1 && (')
  })

  // ⛔ Syarat KEDUA milik `Submit` saja: kontrak yang sudah tuntas tidak
  // dapat dikirim ulang.
  it("⛔ `Submit` hilang ketika status `Resolve Complete`", () => {
    expect(SRC).toContain("statusAkseptasi !== 'Resolve Complete'")
  })

  // ⚠️ Keduanya MATI, dan sebabnya bukan kelalaian — ia terhalang keputusan
  // sasaran tulis. Tombol hidup yang tidak menyimpan apa pun adalah cara
  // tercepat kehilangan suntingan tanpa seorang pun tahu.
  // ⛔ DIBALIK 8 Oktober 2026 — lihat catatan di atas. Kalimat lama
  // menyebut nama Activity Pega dan nama tabel yang dilarang; keduanya
  // urusan pengembang, bukan urusan pemakai.
  it('⛔ matinya tombol TIDAK lagi dijelaskan di layar', () => {
    expect(INFO_SUBMIT.petunjukTombol).toBe('')
  })

  // ⛔ Grid riwayat TIDAK boleh kembali ke tab ini.
  it('⛔ tab ini tidak lagi merender grid riwayat', () => {
    const i = FORM.indexOf("tabTampil === 'Information & Submit'")
    expect(i).toBeGreaterThan(0)
    expect(FORM.slice(i, i + 900)).not.toContain('KOLOM_CATATAN')
  })
})
