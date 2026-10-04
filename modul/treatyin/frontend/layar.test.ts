// Penjaga layar Treaty In — ronde layar 1 (3 Oktober 2026).
//
// ⛔ Yang dijaga di sini bukan rupa layarnya, melainkan bahwa isinya tetap
// SALINAN. Tiap label, pesan galat, dan daftar tab datang dari ekspor Pega
// 2026-09; uji ini gagal pada hari seseorang memperhalus salah satunya.
//
// ⚠️ Uji ini membaca BERKAS SUMBER, bukan merender. Sebabnya: yang hendak
// dijaga adalah teks dan susunannya, dan merender menambah ketergantungan
// (jsdom, komponen inti) tanpa menambah satu pun hal yang dijaga.

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import {
  DAFTAR_KONTRAK,
  FORM_KONTRAK,
  KOLOM_DAFTAR,
  REPORTING_PERIOD,
  TAB_NON_PROPORSIONAL,
  TAB_PROPORSIONAL,
  KOLOM_LIMITS,
  KOLOM_SHARE,
  KOLOM_EVENT_LIMITS,
  KOLOM_RNM_SHARE,
  KOLOM_COIN_SCALE,
  KOLOM_IN2_SEMUA,
  KOLOM_IN2_TIDAK_DIPAKAI,
  JENIS_LIMITS,
  JENIS_SHARE,
  JENIS_EVENT_LIMITS,
  JENIS_RNM_SHARE,
  JENIS_COIN_SCALE,
} from './labels'
import { aksiUntuk, MEDAN_WARISAN } from './pages/DaftarKontrakTreatyIn'
import {
  NON_PROPORSIONAL,
  PROPORSIONAL,
  selAngka,
  tabUntuk,
  tataLetakKolom,
} from './pages/FormKontrakTreatyIn'

const AKAR = __dirname
const LABELS = readFileSync(join(AKAR, 'labels.ts'), 'utf8')
const MENU = readFileSync(join(AKAR, 'menu.ts'), 'utf8')
const FORM = readFileSync(join(AKAR, 'pages', 'FormKontrakTreatyIn.tsx'), 'utf8')
const RUTE = readFileSync(join(AKAR, 'rute.tsx'), 'utf8')

describe('layar daftar kontrak', () => {
  it('kesembilan kolom ada, dalam urutan layar lama', () => {
    expect(KOLOM_DAFTAR.map((k) => k.label)).toEqual([
      'ID',
      'Contract Name',
      'Reinsurance Type',
      'Source of Business',
      'Ceding',
      'Commencement',
      'Termination',
      'Position To',
      'Status Accept',
    ])
  })

  // ⛔ Pokok layar ini: tombolnya BERBEDA menurut keadaan baris, dan
  // meratakannya menghapus satu-satunya petunjuk bahwa sebuah kontrak sudah
  // tidak dapat disunting.
  it('baris Resolve Complete kehilangan Edit dan mendapat Copy + Revision', () => {
    expect(aksiUntuk('Resolve Complete')).toEqual(['View', 'Copy', 'Revision'])
  })

  it('baris yang masih dapat disunting mendapat Edit + View', () => {
    for (const keadaan of ['Open', 'Pending-Approval', 'DRAFT', '']) {
      expect(aksiUntuk(keadaan)).toEqual(['Edit', 'View'])
    }
  })

  it('kedua susunan tombol BERBEDA - penjaganya menggigit', () => {
    expect(aksiUntuk('Resolve Complete')).not.toEqual(aksiUntuk('Open'))
  })

  // ⚠️ RALAT 3 Oktober 2026. Uji ini dulu menuntut petunjuk kosong MENYEBUT
  // tiket 59 — benar ketika layar membaca `KONTRAK` (model baru, nol baris).
  // Sejak layar membaca `TREATY_IN` (1.854 baris) ia **tidak lagi benar**:
  // layar ini tidak menunggu tiket 59, dan petunjuk yang menunjuk tiket yang
  // salah membuat orang menagih pekerjaan yang tidak mengubah apa pun.
  it('kosongnya TIDAK lagi menunjuk tiket 59', () => {
    expect(DAFTAR_KONTRAK.kosongPetunjuk).not.toContain('tiket 59')
    expect(DAFTAR_KONTRAK.kosongPetunjuk).toContain('1.854')
  })

  // Dan tiket 59 TETAP disebut — di tempat yang benar: keterangan sumber,
  // yang menjelaskan daftar MANA yang menunggunya.
  it('tiket 59 disebut sebagai pemindah model baru, bukan sebagai penahan layar ini', () => {
    expect(DAFTAR_KONTRAK.catatanSumber).toContain('tiket 59')
    expect(DAFTAR_KONTRAK.catatanSumber).toContain('POOLDATA.TREATY_IN')
  })
})

describe('strip tab - DUA himpunan, dan keduanya tidak sama', () => {
  it('proporsional sebelas tab, non-proporsional dua belas', () => {
    expect(TAB_PROPORSIONAL).toHaveLength(11)
    expect(TAB_NON_PROPORSIONAL).toHaveLength(12)
  })

  // ⛔ Uji yang paling berguna di berkas ini. Perkiraan ronde ini: "keduanya
  // memuat kesepuluh namanya". Ekspor menyatakan sebaliknya, dan uji ini
  // membekukan temuan itu supaya tidak diratakan diam-diam nanti.
  it('hanya enam nama yang dimiliki KEDUA cabang', () => {
    const sama = TAB_PROPORSIONAL.filter((t) => (TAB_NON_PROPORSIONAL as readonly string[]).includes(t))
    expect(sama).toEqual([
      'Limits',
      'Share',
      'Retro',
      'Exclusions',
      'Special Conditions',
      'Information & Submit',
    ])
  })

  it('lima tab HANYA ada di proporsional, enam HANYA di non-proporsional', () => {
    const hanyaProp = TAB_PROPORSIONAL.filter((t) => !(TAB_NON_PROPORSIONAL as readonly string[]).includes(t))
    const hanyaNon = TAB_NON_PROPORSIONAL.filter((t) => !(TAB_PROPORSIONAL as readonly string[]).includes(t))
    expect(hanyaProp).toEqual([
      'Reporting Period',
      'Portfolio',
      'Co-Ins Scale',
      'Accumulation',
      'Achievement In IDR',
    ])
    expect(hanyaNon).toEqual([
      'Maximum Retention',
      'Event Limits',
      'EGNPI',
      'RNM Share',
      'Installment',
      'Value Difference',
    ])
  })

  it('Reporting Period hanya ada di cabang proporsional', () => {
    expect(TAB_PROPORSIONAL).toContain('Reporting Period')
    expect(TAB_NON_PROPORSIONAL).not.toContain('Reporting Period')
  })

  it('radio memilih himpunannya', () => {
    expect(tabUntuk(PROPORSIONAL)).toBe(TAB_PROPORSIONAL)
    expect(tabUntuk(NON_PROPORSIONAL)).toBe(TAB_NON_PROPORSIONAL)
  })
})

describe('teks yang disalin apa adanya', () => {
  // Pesan galat dirangkai Pega dari DUA rule: `pyCaption Start Date, Due`
  // dan `pyCaption . Must Not Be Empty`. Ejaannya tidak diperhalus.
  it('pesan galat Reporting Period persis seperti ekspor', () => {
    expect(REPORTING_PERIOD.galatKosong).toBe('Start Date, Due. Must Not Be Empty')
  })

  it('ejaan ekspor dipertahankan, termasuk yang tampak keliru', () => {
    // `Teritorial` satu R - begitu ekspornya, dan memperbaikinya membuat
    // pencarian atas layar lama tidak menemukan apa pun.
    expect(FORM_KONTRAK.lingkupWilayah).toBe('Teritorial Scope')
    // `TreatyIn.Bordeaux` adalah nama PROPERTINYA; labelnya `Bordereaux`.
    expect(FORM_KONTRAK.bordereaux).toBe('Bordereaux')
    expect(LABELS).toContain('TreatyIn.Bordeaux')
  })

  it('tiap medan membawa jejak posisi bitanya', () => {
    const jejak = LABELS.match(/@\d{4,}/g) ?? []
    expect(jejak.length).toBeGreaterThan(30)
  })
})

describe('blok bergaris mati tidak dibangun', () => {
  // Dua blok `pyContainerVisibleWhen = 1=2` ditemukan di
  // `TreatyInTabsProportional.xml`. Keduanya disebut di kepala berkas form
  // supaya yang membaca kode tahu apa yang SENGAJA tidak ada.
  it('keduanya disebut di kepala FormKontrakTreatyIn.tsx', () => {
    expect(FORM).toContain('1=2')
    expect(FORM).toContain('@1095716')
    expect(FORM).toContain('@1300558')
  })

  it('salinan HIDUP Account Reporting Period yang dipakai', () => {
    expect(FORM).toContain('@45828')
    expect(REPORTING_PERIOD.judul).toBe('Account Reporting Period')
  })
})

describe('L-4 dicabut, dan pencabutannya tercatat', () => {
  it('menu.ts menyebut pencabutannya berikut cacah berkasnya', () => {
    expect(MENU).toContain('L-4')
    expect(MENU).toContain('DICABUT')
    expect(MENU).toMatch(/50 berkas `Section\/` dan 3 `Harness\/`/)
  })

  it('halaman awal modul kini daftar kontrak', () => {
    expect(MENU).toContain("HALAMAN_AWAL_TREATYIN: HalamanTreatyIn = 'treatyin-kontrak'")
  })
})

// ===========================================================================
// RONDE LAYAR 2 — susunan, bukan medan
// ===========================================================================

const CSS = readFileSync(join(AKAR, 'treatyin.css'), 'utf8')
const DAFTAR = readFileSync(join(AKAR, 'pages', 'DaftarKontrakTreatyIn.tsx'), 'utf8')
const ACUAN = readFileSync(join(AKAR, 'pages', 'AcuanTreatyIn.tsx'), 'utf8')

describe('form dua kolom yang mengalir sendiri', () => {
  // ⛔ Pokoknya: kolom kanan LEBIH PANJANG. Grid pengisi baris tidak dapat
  // menghasilkan itu tanpa sel kosong yang disisipkan tangan, dan sel kosong
  // itu bergeser diam-diam pada medan berikutnya yang ditambahkan.
  it('kolom kanan lebih panjang daripada kiri', () => {
    expect(tataLetakKolom.kiri).toHaveLength(5)
    expect(tataLetakKolom.kanan).toHaveLength(7)
    expect(tataLetakKolom.kanan.length).toBeGreaterThan(tataLetakKolom.kiri.length)
  })

  it('urutannya urutan rujukan', () => {
    expect(tataLetakKolom.kiri[0]).toBe('Treaty Contract Name')
    expect(tataLetakKolom.kanan[0]).toBe('Commencement')
    expect(tataLetakKolom.kanan.slice(-3)).toEqual([
      'Ceding',
      'RNM as Treaty Leader',
      'Source of Business',
    ])
  })

  it('nol medan hilang atau kembar di antara kedua kolom', () => {
    const semua = [...tataLetakKolom.kiri, ...tataLetakKolom.kanan]
    expect(new Set(semua).size).toBe(semua.length)
    expect(semua).toHaveLength(12)
  })

  // Pengisi `<span />` yang dipakai sebelumnya DICABUT; uji ini menjaganya
  // tidak kembali diam-diam.
  //
  // ⚠️ Komentar dibuang lebih dulu. Kepala berkas itu MENYEBUT `<span />`
  // untuk menjelaskan kenapa ia dicabut, dan sapuan yang memindai teks
  // mentah menemukan penjelasannya sendiri — pelajaran yang sama yang sudah
  // membuat `tanpaKomentarGo` lahir di penjaga arsip.
  it('form tidak lagi memakai sel pengisi', () => {
    const kode = FORM.replace(/\{\/\*[\s\S]*?\*\/\}/g, '').replace(/^\s*\/\/.*$/gm, '')
    expect(kode).toContain('trin__dwikolom')
    expect(kode).not.toContain('<span />')
  })
})

describe('lebar yang disetel membawa asal-usulnya', () => {
  // Perbandingan 193 : 349 : 80 adalah satu-satunya lebar yang ekspor
  // sungguh bedakan di layar ini.
  it('grid Rate of Exchange memakai perbandingan dari ekspor', () => {
    expect(FORM).toContain('trin__kol-mata-uang')
    expect(CSS).toMatch(/\.trin__kol-mata-uang\s*\{\s*width:\s*27%/)
    expect(CSS).toMatch(/\.trin__kol-nilai\s*\{\s*width:\s*48%/)
    expect(CSS).toContain('193')
    expect(CSS).toContain('349')
  })

  // ⛔ RALAT yang dijaga: lebar kolom daftar BUKAN dari ekspor. Kesembilan
  // sel baris saring `pyWidth=130`, seragam. Komentar yang mengaku
  // sebaliknya sudah dicabut, dan uji ini menahannya tetap tercabut.
  it('lebar kolom daftar TIDAK mengaku berasal dari ekspor', () => {
    expect(CSS).not.toMatch(/tidak seragam - mengikuti ekspor/)
    expect(CSS).toContain('pyWidth=130')
    expect(CSS).toContain('keputusan rupa KITA')
  })

  /**
   * ⛔ RALAT 3 Oktober 2026. Uji ini dulu menuntut `white-space: nowrap` atas
   * SELURUH sel `.trin__tabel`, membekukan pembacaan `pyWrap=false` yang benar
   * tentang Pega dan salah tentang kita.
   *
   * Yang sungguh terjadi di layar: dengan nama kontrak sepanjang "WHOLE ACCOUNT
   * RISK AND CATASTROPHE EXCESS OF LOSS REINSURANCE 2025", tabelnya melar jauh
   * melewati kartunya dan HANYA DUA dari sembilan kolom terlihat — tujuh
   * sisanya di balik gulir mendatar. Uji lama akan menolak setiap perbaikannya.
   *
   * Pega merender pada lebar tetap; tabel ini tidak. `pyWrap=false` karena itu
   * dihormati di tempat ia berarti — pengenal, jenis, dan tanggal tidak boleh
   * terpotong — dan tidak di tiga kolom teks panjang.
   */
  it('nowrap HANYA pada kolom pendek, tiga kolom teks membungkus', () => {
    // `sifatProporsi` DICABUT dari daftar ini 3 Oktober 2026: dua kata yang
    // meluber ke kolom sebelahnya, bukan satu kata yang rusak bila dipenggal.
    for (const k of ['id', 'tanggalMulai', 'tanggalBerakhir']) {
      expect(CSS).toMatch(new RegExp(`trin__k--${k}[^}]*white-space: nowrap`, 's'))
    }
    // Tiga kolom teks panjang TIDAK boleh nowrap.
    for (const k of ['namaKontrak', 'idAsalBisnis', 'idCedant', 'sifatProporsi']) {
      expect(CSS).not.toMatch(new RegExp(`trin__k--${k} \\{[^}]*white-space: nowrap`, 's'))
    }
    // Dan nowrap menyeluruh tidak boleh kembali.
    expect(CSS).not.toMatch(/\.trin__tabel td,\s*\n\s*\.treatyin \.trin__tabel th \{\s*\n\s*white-space: nowrap/)
  })

  /**
   * Kesembilan kolom harus MUAT. `table-layout: fixed` + lebar persen adalah
   * yang menjamin itu; tanpa keduanya kolom melebar mengikuti isi terpanjang.
   */
  it('kesembilan kolom punya lebar persen dan tabelnya fixed', () => {
    expect(CSS).toMatch(/\.treatyin \.trin__tabel \{[^}]*table-layout: fixed/s)
    const kunci = ['id', 'namaKontrak', 'sifatProporsi', 'idAsalBisnis', 'idCedant',
      'tanggalMulai', 'tanggalBerakhir', 'posisiKe', 'keadaanSiklusHidup']
    for (const k of kunci) {
      expect(CSS).toMatch(new RegExp(`trin__k--${k}[^}]*width: \\d+%`, 's'))
    }
  })

  /**
   * ⛔ Nama kelas lebar HARUS sama dengan `kunci` kolomnya. Enam di antaranya
   * pernah salah eja (`jenisReasuransi` lawan `sifatProporsi`, `mulai` lawan
   * `tanggalMulai`, …), sehingga aturannya mati tanpa suara: CSS yang tidak
   * cocok tidak melempar galat, ia hanya tidak berlaku.
   */
  it('tiap kunci kolom punya aturan lebarnya, dan ejaannya cocok', () => {
    const kunci = [...LABELS.matchAll(/kunci: '([a-zA-Z]+)'/g)].map((m) => m[1])
    expect(kunci.length).toBe(9)
    for (const k of kunci) {
      expect(CSS).toContain(`trin__k--${k}`)
    }
  })
})

describe('daftar: saring per kolom dan penomoran halaman', () => {
  it('kesembilan kolom punya ikon saringnya', () => {
    expect(DAFTAR).toContain('trin__ikon-saring')
    expect(DAFTAR).toContain('aria-pressed')
  })

  it('baris saring berisi kotak, bukan sel kosong', () => {
    expect(DAFTAR).toContain('trin__baris-saring')
    expect(DAFTAR).toMatch(/aria-label=\{`Saring \$\{k\.label\}`\}/)
  })

  // ⚠️ RALAT 3 Oktober 2026 — syaratnya berubah bersama sumber datanya.
  // Dengan 25 baris per halaman dari 1.854, `tampil.length` SELALU 25 dan
  // penomoran akan mengira hanya ada satu halaman. Yang menentukan `total`
  // dari server.
  it('penomoran halaman memakai total dari SERVER, bukan panjang halaman', () => {
    expect(DAFTAR).toContain('trin__kepala-kanan')
    expect(DAFTAR).toContain('<Halaman')
    expect(DAFTAR).toMatch(/total=\{total\}/)
    expect(DAFTAR).not.toMatch(/total=\{tampil\.length\}/)
  })

  // ⛔ Server sudah memotong halamannya. Memotong lagi di peramban membuat
  // halaman kedua menampilkan halaman kosong.
  it('nol pemotongan halaman kedua kali di peramban', () => {
    expect(DAFTAR).not.toContain('.slice((halaman - 1) * UKURAN_HALAMAN')
  })

  // Ukuran halaman keputusan kita, bukan salinan — dan ia menyatakannya.
  it('ukuran halaman dinyatakan sebagai keputusan kita', () => {
    expect(DAFTAR).toContain('UKURAN_HALAMAN')
    expect(DAFTAR).toContain('KEPUTUSAN KITA')
  })
})

describe('tab belum dibangun dibedakan dari data kosong', () => {
  // ⛔ Dua pernyataan yang berbeda: "belum ada KODE" versus "belum ada DATA".
  // Yang salah membacanya menagih pemuatan data padahal yang kurang layarnya.
  it('memakai .trin__belum, bukan komponen Kosong', () => {
    expect(FORM).toContain('trin__belum')
    const potong = FORM.slice(FORM.indexOf('tabTampil === '))
    expect(potong).not.toContain('<Kosong')
  })

  it('bentuknya berbeda di CSS - tepi putus-putus, bukan ikon', () => {
    expect(CSS).toMatch(/\.trin__belum \{[\s\S]*?border: 1px dashed/)
  })
})

describe('AcuanTreatyIn ikut bahasa visual yang sama', () => {
  it('memakai trin__tabel dan merender kepala kolom walau kosong', () => {
    expect(ACUAN).toContain('className="trin__tabel"')
    expect(ACUAN).toContain('colSpan={bersusun ? 4 : 3}')
  })
})

describe('pin .form-grid tidak merusak lebar telepon', () => {
  // ⛔ RALAT yang dijaga. Pin dua-kolom berkekhususan lebih tinggi daripada
  // aturan inti `@media (max-width: 860px) { .form-grid { 1fr } }`, jadi
  // tanpa media query sendiri ia menimpanya — dan Reporting Period bertahan
  // dua kolom di lebar telepon, persis cacat yang aturan inti itu cegah.
  it('pin hanya berlaku dari 861px ke atas', () => {
    expect(CSS).toMatch(/@media \(min-width: 861px\) \{\s*\n\s*\.treatyin \.form-grid \{/)
  })

  it('dwikolom runtuh satu kolom di bawah 900px', () => {
    expect(CSS).toMatch(/@media \(max-width: 899px\) \{\s*\n\s*\.treatyin \.trin__dwikolom \{\s*\n\s*grid-template-columns: minmax\(0, 1fr\)/)
  })
})

// ===========================================================================
// DATA NYATA — `POOLDATA.TREATY_IN`
// ===========================================================================

describe('layar daftar membaca tabel warisan', () => {
  it('memanggil jalur warisan, bukan jalur model baru', () => {
    expect(DAFTAR).toContain('ambilDaftarWarisan')
    expect(DAFTAR).not.toContain('ambilDaftarKontrak(')
  })

  // ⛔ Satu halaman per permintaan: 1.854 baris tidak ditarik sekaligus.
  it('memuat ulang setiap kali halamannya berganti', () => {
    expect(DAFTAR).toMatch(/\}, \[halaman\]\)/)
  })

  // Pemetaan kolom layar -> medan warisan berdiri sendiri supaya dapat diuji.
  it('kesembilan kolom terpetakan ke medan warisan', () => {
    const kunci = KOLOM_DAFTAR.map((k) => k.kunci)
    expect(Object.keys(MEDAN_WARISAN).sort()).toEqual([...kunci].sort())
    expect(MEDAN_WARISAN.idCedant).toBe('cedant')
    expect(MEDAN_WARISAN.idAsalBisnis).toBe('asalBisnis')
    expect(MEDAN_WARISAN.keadaanSiklusHidup).toBe('statusAkseptasi')
    expect(MEDAN_WARISAN.posisiKe).toBe('posisiKe')
  })

  // ⭐ Aturan tombol TIDAK ditulis ulang — ia aturan yang sama, kini diberi
  // makan oleh `STATUSAKSEPTASI` yang nyata. Sapuan menemukan EMPAT nilai,
  // bukan satu: `Resolve Complete` 1.820 · `Accept` 12 · `Decline` 11 ·
  // NULL 11. Ketiga yang terakhir 34 baris nyata, bukan kasus teoretis.
  it('keempat nilai STATUSAKSEPTASI nyata memilih tombolnya', () => {
    expect(aksiUntuk('Resolve Complete')).toEqual(['View', 'Copy', 'Revision'])
    for (const lain of ['Accept', 'Decline', '']) {
      expect(aksiUntuk(lain)).toEqual(['Edit', 'View'])
    }
  })

  it('baris memakai statusAkseptasi, bukan medan model baru', () => {
    expect(DAFTAR).toContain('aksiUntuk(b.statusAkseptasi)')
  })

  // ⛔ Pengenal TEKS. `TREATY_IN.ID` adalah VARCHAR2(100).
  it('pengenal kontrak bertipe teks di seluruh jalurnya', () => {
    expect(DAFTAR).toContain('onBuka: (id: string) => void')
    expect(FORM).toContain('idKontrak: string')
    expect(RUTE).toContain('useState<string | null>(null)')
  })
})

// ===========================================================================
// FORM MEMBACA KONTRAK WARISAN NYATA
// ===========================================================================

describe('form membaca kontrak warisan', () => {
  it('memanggil jalur warisan, bukan jalur model baru', () => {
    expect(FORM).toContain('ambilKontrakWarisan')
    expect(FORM).not.toContain('ambilKontrak(')
  })

  it('memuat ulang ketika pengenalnya berganti, dan melewati kontrak baru', () => {
    expect(FORM).toMatch(/\}, \[idKontrak\]\)/)
    expect(FORM).toContain("if (idKontrak === '')")
  })

  // ⭐ Radio disambungkan ke NILAI NYATA. Kontrak non-proporsional yang
  // membuka strip proporsional memperlihatkan sebelas tab yang tidak satu
  // pun miliknya.
  it('radio jenis diisi dari kontraknya, bukan dibiarkan bawaan', () => {
    expect(FORM).toMatch(/setJenis\(k\.sifatProporsi === NON_PROPORSIONAL/)
  })

  // ⛔ Medan MATI, bukan kotak kosong — dan ia bukan kasus langka:
  // `ContractRefNo` tidak ada di 1.112 dari 1.854 kontrak, `TreatyLeader`
  // di 1.195, `BordereauxNote` di 836.
  it('ketiga medan yang kuncinya bisa hilang memakai MedanTakAda', () => {
    expect(FORM).toContain('function MedanTakAda')
    for (const kunci of ['ContractRefNo', 'BordereauxNote', 'TreatyLeader']) {
      expect(FORM).toContain(`adaKunci('${kunci}')`)
    }
  })

  it('keterangan medan mati menyatakan SISTEM LAMA, bukan "belum diisi"', () => {
    expect(FORM_KONTRAK.takAdaDiWarisan).toContain('sistem lama')
    expect(FORM_KONTRAK.takAdaDiWarisan).not.toContain('belum diisi')
  })

  // ⭐ `TREATYYEAR` adalah kolom — dibaca, bukan dipotong dari tanggal.
  it('Treaty Year dibaca dari kolomnya', () => {
    expect(FORM).toContain('warisan?.tahunTreaty')
    expect(FORM).not.toContain("mulai.slice(0, 4)")
  })

  // ⚠️ RALAT: nilai pilihan adalah yang TERSIMPAN, dan ada DUA per pilihan.
  it('pilihan membawa kedua nilai nyata, huruf kecil seperti tersimpan', () => {
    expect(FORM_KONTRAK.bordereauxNilai).toEqual(['reporting', 'nonreporting'])
    expect(FORM_KONTRAK.caraPembukuanNilai).toEqual(['underwriting', 'accounting'])
  })

  // Terjemahan TIDAK ditulis dua kali — ia dipakai ulang dari services.
  it('nol penerjemah kedua di frontend', () => {
    expect(FORM).not.toMatch(/NonProportional'\s*\?\s*'Non Proportional'/)
    expect(FORM).not.toContain('/19')
  })
})

describe('empat tab dari M_TREATY_IN2 — satu tabel, empat proyeksi', () => {
  it('keempatnya membaca `warisan.layer`, bukan empat sumber berbeda', () => {
    for (const tab of ['Limits', 'Share', 'Event Limits', 'RNM Share']) {
      expect(FORM).toContain(`tabTampil === '${tab}'`)
    }
    // ⛔ Satu seam. Empat pembacaan berarti empat kueri untuk satu baca, dan
    // empat kesempatan agar yang satu melihat layer yang lain tidak.
    const pakai = FORM.match(/warisan\?\.layer \?\? \[\]/g) ?? []
    expect(pakai.length).toBe(4)
  })

  it('nol tabel baru dibuat untuk keempatnya', () => {
    for (const salah of [
      'M_TREATYIN_LIMITS',
      'M_TREATYIN_SHARE',
      'M_TREATYIN_EVENTLIMITS',
      'M_TREATYIN_RNMSHARE',
    ]) {
      expect(FORM).not.toContain(salah)
      expect(LABELS).not.toContain(salah)
    }
  })

  it('tiap kolom punya golongan angkanya, dan panjangnya sejajar', () => {
    expect(JENIS_LIMITS.length).toBe(KOLOM_LIMITS.length)
    expect(JENIS_SHARE.length).toBe(KOLOM_SHARE.length)
    expect(JENIS_EVENT_LIMITS.length).toBe(KOLOM_EVENT_LIMITS.length)
    expect(JENIS_RNM_SHARE.length).toBe(KOLOM_RNM_SHARE.length)
    expect(JENIS_COIN_SCALE.length).toBe(KOLOM_COIN_SCALE.length)
  })

  it('LAYER tidak pernah diformat — ia pengenal, dan nol di depannya bermakna', () => {
    for (const [kolom, jenis] of [
      [KOLOM_LIMITS, JENIS_LIMITS],
      [KOLOM_SHARE, JENIS_SHARE],
      [KOLOM_EVENT_LIMITS, JENIS_EVENT_LIMITS],
      [KOLOM_RNM_SHARE, JENIS_RNM_SHARE],
    ] as const) {
      const i = kolom.indexOf('LAYER' as never)
      expect(i).toBeGreaterThanOrEqual(0)
      expect(jenis[i]).toBe('teks')
    }
  })

  it('kesepuluh kolom kepala disebut namanya, bukan dibuang diam-diam', () => {
    expect(KOLOM_IN2_TIDAK_DIPAKAI.length).toBe(10)
    // ⛔ Dan tidak satu pun dari kesepuluhnya menyelinap ke grid: ia akan
    // mengulang nilai yang sama di setiap baris tanpa menambah keterangan.
    const semuaKolom = [
      ...KOLOM_LIMITS,
      ...KOLOM_SHARE,
      ...KOLOM_EVENT_LIMITS,
      ...KOLOM_RNM_SHARE,
    ]
    for (const k of KOLOM_IN2_TIDAK_DIPAKAI) {
      expect(semuaKolom).not.toContain(k)
    }
  })

  it('⭐ ke-41 kolom TERHITUNG HABIS: 31 di tab, 10 kepala, NOL hilang', () => {
    // ⛔ Penjumlahan ini yang memaksa kolom yang hilang TERLIHAT. Kolom yang
    // lenyap dari layar tidak menimbulkan satu pun galat — hanya sebuah
    // angka yang tidak ada lagi di mana pun.
    const dipakai = new Set<string>([
      ...KOLOM_LIMITS,
      ...KOLOM_SHARE,
      ...KOLOM_EVENT_LIMITS,
      ...KOLOM_RNM_SHARE,
    ])
    expect(KOLOM_IN2_SEMUA.length).toBe(41)
    expect(dipakai.size).toBe(31)
    expect(KOLOM_IN2_TIDAK_DIPAKAI.length).toBe(10)
    expect(dipakai.size + KOLOM_IN2_TIDAK_DIPAKAI.length).toBe(41)

    // Dan keduanya benar-benar MENUTUPI ke-41, tanpa tumpang tindih.
    const tertutup = new Set<string>([...dipakai, ...KOLOM_IN2_TIDAK_DIPAKAI])
    expect(tertutup.size).toBe(41)
    for (const k of KOLOM_IN2_SEMUA) {
      expect(tertutup.has(k)).toBe(true)
    }
  })
})

describe('angka: titik ribuan, koma desimal, dan batas desimalnya', () => {
  it('uang memakai 4 desimal dan titik ribuan', () => {
    expect(selAngka('uang', '5700000000')).toBe('5.700.000.000')
    expect(selAngka('uang', '1234.56789')).toBe('1.234,5679')
  })

  it('persen biasa dipangkas 2 desimal dan diberi tanda', () => {
    expect(selAngka('persen', '12.1000')).toBe('12,1%')
    expect(selAngka('persen', '33.333')).toBe('33,33%')
  })

  it('⛔ persen SHARE dipotong pada DELAPAN desimal — keputusan §13', () => {
    // Tiga share 33,333 berjumlah TEPAT 100; tiga share 33,33 tidak. Delapan
    // desimal menjaga sifat itu tanpa mengaku lebih teliti daripada
    // penyimpanannya, yang juga 8 (`NUMBER(38,8)`).
    expect(selAngka('persenShare', '33.333')).toBe('33,333%')

    // ⚠️ `99.999999999999900` TETAP tampil `100%` pada 8 desimal — dan itu
    // DIUKUR, bukan diandaikan. Pembulatannya setengah-ke-atas dan
    // limpahannya merambat, jadi selisihnya baru terlihat pada 13 desimal
    // ke atas. Salah satu dari dua alasan pemilik proses menolak 2 desimal
    // karena itu TIDAK terpenuhi oleh 8; alasan keduanya (sama dengan batas
    // penyimpanan) terpenuhi. Dicatat di KEPUTUSAN §13.1 dan ditagih balik.
    expect(selAngka('persenShare', '99.999999999999900')).toBe('100%')
    expect(selAngka('persenShare', '99.9999999')).toBe('99,9999999%')

    // ⛔ Tiga puluh desimal DIPOTONG pada delapan, tidak dibiarkan utuh.
    expect(selAngka('persenShare', '2.825601535925207120348922139444')).toBe('2,82560154%')
  })

  it('persen BUKAN-share tetap 2 desimal — ia tidak dijumlahkan menjadi 100', () => {
    expect(selAngka('persen', '33.333')).toBe('33,33%')
    expect(selAngka('persen', '12.1000')).toBe('12,1%')
  })

  it('⛔ pita tidak diberi tanda % kedua walau salah golong', () => {
    // `formatPersen` menempelkan `%` pada teks bukan-angka apa adanya,
    // menghasilkan `>=30% up to < 50%%`. Satu huruf salah di JENIS_* sudah
    // cukup untuk memicunya.
    expect(selAngka('persenShare', '>=30% up to < 50%')).toBe('>=30% up to < 50%')
    expect(selAngka('persen', '1/04/2023')).toBe('1/04/2023')
  })

  it('nol di ekor tetap dibuang, dan uang tetap 4 desimal', () => {
    expect(selAngka('persenShare', '12.10000000')).toBe('12,1%')
    expect(selAngka('persenShare', '50')).toBe('50%')
    expect(selAngka('uang', '1234.56789')).toBe('1.234,5679')
  })



  it('teks bukan-angka dikembalikan APA ADANYA, tidak berpura-pura nol', () => {
    // `CoInShare` adalah PITA, bukan bilangan — 702 barisnya berbentuk ini.
    expect(selAngka('teks', '>=30% up to < 50%')).toBe('>=30% up to < 50%')
    expect(selAngka('uang', '1/04/2023')).toBe('1/04/2023')
  })

  it('pengenal tidak kehilangan nol di depannya', () => {
    expect(selAngka('teks', '007')).toBe('007')
    expect(selAngka('teks', '1A')).toBe('1A')
  })

  it('kosong tetap kosong — bukan 0', () => {
    for (const j of ['uang', 'persen', 'persenShare', 'teks'] as const) {
      expect(selAngka(j, '')).toBe('')
    }
  })

  it('format.ts dipanggil, tidak ditulis ulang', () => {
    expect(FORM).toContain("from '../../../../inti/frontend/lib/format'")
    // ⛔ Nol pemformat kedua di dalam modul ini.
    expect(FORM).not.toMatch(/function formatNumber|function formatPersen/)
    expect(LABELS).not.toMatch(/function formatNumber|function formatPersen/)
  })
})

describe('tab Co-Ins Scale dan petunjuk kosong yang membedakan sebabnya', () => {
  it('Co-Ins Scale membaca tabel pendaratan kesembilan', () => {
    expect(FORM).toContain("tabTampil === 'Co-Ins Scale'")
    expect(FORM).toContain('warisan?.skalaKoasuransi')
  })

  it('⛔ petunjuk kosong keempat tab layer TIDAK mengaku tahu sebabnya', () => {
    // Pada 510 kontrak, "kontrak ini memang tidak punya" KELIRU — datanya
    // ada di dokumen, tabelnya yang tidak mencakupnya.
    expect(FORM_KONTRAK.petunjukLayer).toContain('1.340 dari 1.854')
    expect(FORM_KONTRAK.petunjukLayer).toContain('510')
    expect(FORM_KONTRAK.petunjukLayer).not.toContain('memang tidak punya.')
  })

  it('petunjuk Event Limits menyebut tiga batas yang TIDAK ada sumbernya', () => {
    for (const batas of ['Flood Jabodetabek', 'Flood Nationwide', 'RSMD']) {
      expect(FORM_KONTRAK.petunjukEventLimits).toContain(batas)
    }
  })

  it('keempat tab layer memakai petunjuknya sendiri, bukan petunjuk umum', () => {
    for (const tab of ['Limits', 'Share', 'RNM Share']) {
      const i = FORM.indexOf(`tabTampil === '${tab}'`)
      expect(i).toBeGreaterThan(0)
      const blok = FORM.slice(i, i + 1400)
      expect(blok).toContain('petunjukLayer')
    }
  })
})

describe('dua tab TEKS — Jalan B, keputusan §15', () => {
  it('keduanya membaca dokumen, bukan tabel pendaratan', () => {
    expect(FORM).toContain("tabTampil === 'Exclusions'")
    expect(FORM).toContain("tabTampil === 'Special Conditions'")
    expect(FORM).toContain('warisan?.pengecualian')
    expect(FORM).toContain('warisan?.syaratKhusus')
    // ⛔ Nol tabel pendaratan dibuat untuk keduanya — itulah Jalan B.
    for (const salah of ['M_TREATYIN_TEKS', 'M_TREATYIN_EXCLUSIONS', 'M_TREATYIN_SPECIALCONDITIONS']) {
      expect(FORM).not.toContain(salah)
      expect(LABELS).not.toContain(salah)
    }
  })

  it('teksnya DAPAT DIGULIR, bukan satu baris terpotong', () => {
    expect(FORM).toContain('trin__teks')
    // Isinya mencapai 23.453 aksara; sel yang memotongnya menyembunyikan
    // teks tanpa memberi tahu ada yang disembunyikan.
    expect(CSS).toMatch(/\.trin__teks \{[\s\S]*?overflow-y: auto/)
    expect(CSS).toMatch(/\.trin__teks \{[\s\S]*?white-space: pre-wrap/)
  })

  it('⛔ kosong memakai Kosong (belum ada DATA), bukan .trin__belum', () => {
    const i = FORM.indexOf('function TabTeksPanjang')
    expect(i).toBeGreaterThan(0)
    const blok = FORM.slice(i, i + 1800)
    expect(blok).toContain('<Kosong')
    expect(blok).not.toContain('trin__belum')
  })

  it('peringatan ejaan lain menyatakan isinya BERBEDA, bukan salinan', () => {
    expect(FORM_KONTRAK.ejaanLainBerisi).toContain('BERBEDA')
    expect(FORM_KONTRAK.ejaanLainBerisi).toContain('bukan salinan')
    expect(FORM).toContain('ejaanLainBerisi')
  })

  it('petunjuk kosong menyebut ejaan cabang seberang TIDAK dipakai', () => {
    expect(FORM_KONTRAK.petunjukTeksPengecualian).toContain('TIDAK dipakai')
    expect(FORM_KONTRAK.petunjukTeksSyarat).toContain('SpecialConditionsp')
  })

  it('⛔ Value Difference TETAP belum dibangun — ia objek, bukan teks', () => {
    expect(FORM).not.toContain("tabTampil === 'Value Difference'")
  })
})
