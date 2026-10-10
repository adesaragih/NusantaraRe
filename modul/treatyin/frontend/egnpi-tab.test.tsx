// Uji tab EGNPI (non-prop) — bentuk layar lawan ekspor Pega.
//
// Dirender statis (`react-dom/server`) — BUKAN pengganti melihat layar di
// peramban. RUMUSNYA diuji di services (`hitung_egnpi_test.go`); yang
// ditanyakan di sini hanya: apakah yang ekspor perlihatkan ADA di layar,
// dan apakah yang ekspor kunci TETAP terkunci.

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import type { BarisEgnpi, OpsiLimits } from './api'
import TabEgnpi from './components/TabEgnpi'
import { DESIMAL_EGNPI, EGNPI, KOLOM_GRID_EGNPI } from './labelsEgnpi'

const AKAR = __dirname
const SRC = readFileSync(join(AKAR, 'components', 'TabEgnpi.tsx'), 'utf8')
const CSS = readFileSync(join(AKAR, 'treatyin.css'), 'utf8')

const OPSI: OpsiLimits = {
  jenisTreaty: [],
  kelompokTreaty: [{ id: '10007', nama: 'PROPERTY', kembar: false }],
  mataUang: [
    { id: '10026', nama: 'IDR', kembar: false },
    { id: '10001', nama: 'USD', kembar: false },
  ],
}

const BARIS: BarisEgnpi[] = [
  {
    ID: '1',
    TreatyGroup: 'PROPERTY',
    TreatyGroupID: '10007',
    AsDate: '01-01-2025',
    Proportion: '41.63535247256876597',
    Currency: 'USD',
    CurrencyID: '10001',
    Amount: '11080000000',
    AmountIDR: '137849315068',
    ClassOfBusiness: '',
    ClassOfBusinessID: '',
    Note: 'catatan baris',
  },
]

function render(mode: 'lihat' | 'ubah') {
  return renderToStaticMarkup(
    <TabEgnpi baris={BARIS} kurs={[]} retensi={[]} opsiAwal={OPSI} mode={mode} />,
  )
}

describe('tab EGNPI — bentuk dari ekspor', () => {
  it('⭐ panel dan kedua tombol berumus ada', () => {
    const html = render('ubah')
    expect(html).toContain(EGNPI.panel)
    // `Update Total` → TreatyInNPSetTotal(egnpi)
    expect(html).toContain(EGNPI.perbaruiTotal)
    // `Update EGNPI Value` → TreatyInEGNPIListValue
    expect(html).toContain(EGNPI.perbaruiNilai)
  })

  it('⭐ ketiga panel total ekspor ada', () => {
    const html = render('ubah')
    expect(html).toContain(EGNPI.totalPerMataUang)
    expect(html).toContain(EGNPI.totalIDR)
    expect(html).toContain(EGNPI.totalProporsi)
  })

  // ⛔ Tanpa `Add`, nol baris EGNPI dapat lahir — dan tab Limits membaca
  // baris ini lewat `TotalEgnpi`. Jadi tombol ini bukan hiasan.
  it('⛔ Add hanya di mode ubah', () => {
    expect(render('ubah')).toContain(EGNPI.tambah)
    expect(render('lihat')).not.toContain(EGNPI.tambah)
  })

  it('⛔ Delete hanya di mode ubah', () => {
    expect(render('ubah')).toContain(EGNPI.hapus)
    expect(render('lihat')).not.toContain(EGNPI.hapus)
  })

  // ⭐ Kepala kartu memuat apa yang di Pega menjadi KOLOM grid. Tiap kolom
  // ekspor harus punya tempatnya di layar — sebagai kepala kartu atau
  // sebagai medan rincian.
  it('⭐ keenam kolom grid ekspor terwakili', () => {
    const html = render('lihat')
    for (const k of KOLOM_GRID_EGNPI) {
      // `As Date` di grid berbunyi `As At` di rinciannya — keduanya sah.
      const cari = k === 'As Date' ? EGNPI.asAt : k
      expect(html, k).toContain(cari)
    }
  })

  // ⛔ RALAT 7 Oktober 2026 — UJI INI DULU MENUNTUT HAL YANG SALAH.
  //
  // Bentuk pertamanya menuntut `As At` dan `Note` terkunci SELAMANYA, atas
  // dasar `pyReadOnly = true` di `Section/DetailEGNPI.xml`. Sel yang SAMA
  // memuat `pyReadOnlyCondition = TreatyIn.ViewState = 1`, dan syarat itulah
  // yang berlaku — `ViewState 1` = mode lihat, jadi keduanya AKTIF di mode
  // Edit. Aturannya terkodekan di `alat-baca-ekspor/baca.py` (`hanya_baca`).
  it('⭐ As At dan Note AKTIF di mode ubah, terkunci di mode lihat', () => {
    expect(render('ubah')).toContain('catatan baris')
    // Mode lihat mengunci; mode ubah tidak.
    expect(SRC).toContain('readOnly={!bisaUbah}')
    expect(SRC).toContain('onUbah({ ...b, Note: v })')
    expect(SRC).toContain('onUbah({ ...b, AsDate: v })')
  })

  // ⛔ PERMINTAAN PEMILIK PROSES 7 Oktober 2026 — `As At` INPUTAN TANGGAL.
  // Ekspor pun berkata begitu: `pyFormat = pxDateTime`, dan tangkapan layar
  // memperlihatkan ikon kalender. Bentuk sebelumnya kotak teks biasa.
  // ⭐ 7 Oktober 2026 — tanggal yang DAPAT DIKETIK + ikon kalender
  // (`TanggalRedup` → `TanggalKetik.tsx`), permintaan pemakai.
  it('⛔ As At medan tanggal (ketik atau kalender), bukan kotak teks biasa', () => {
    expect(SRC).toContain('<TanggalRedup')
    expect(SRC).toContain('label={EGNPI.asAt}')
    const html = render('ubah')
    expect(html).toContain(`aria-label="Pick from the calendar — ${EGNPI.asAt}"`)
    expect(html).toContain('type="date"')
  })

  // ⛔ PERMINTAAN PEMILIK PROSES — pemilihnya DROPDOWN yang sama dengan
  // `Ceding` dan `Source of Business`, bukan autocomplete ketik-bebas.
  // Komponennya SAMA (`DropdownWarisan`), bukan tiruannya.
  it('⛔ Treaty Group dan Currency memakai DropdownWarisan', () => {
    expect(SRC).toContain("import DropdownWarisan from './DropdownWarisan'")
    expect(SRC).not.toContain('IsianAuto')
    expect((SRC.match(/<DropdownWarisan/g) ?? []).length).toBe(2)
  })

  // ⛔ `ambil` adalah TANGGUNGAN efek `DropdownWarisan`; fungsi baru tiap
  // render akan mengambil ulang daftarnya tanpa henti.
  it('⛔ fungsi `ambil` dibekukan useCallback', () => {
    expect(SRC).toContain('useCallback')
    expect(SRC).toContain('const ambilKelompokTreaty = useCallback(')
    expect(SRC).toContain('const ambilMataUang = useCallback(')
  })

  // ⭐ `Amount` dan `Amount in IDR` SEPASANG: satuan di kiri, nilai di
  // kanan tanpa label sendiri — bentuk ekspor dan tangkapan layar.
  it('⭐ Amount dan Amount in IDR dirender berpasangan', () => {
    // ⭐ 8 Oktober 2026 (agen N, tata Pega): kedua baris kini `Inline grid
    // 30 70` @1324 = `TataPegaBlok tata="t3070"` [satuan berlabel | nilai],
    // bukan lagi kelas `trin__egnpi-pasangan` (1fr:2fr). Empat t3070 di
    // berkas: dua baris rincian ini + dua di blok total.
    expect((SRC.match(/<TataPegaBlok tata="t3070">/g) ?? []).length).toBe(4)
    // Satuan baris kedua TETAP `IDR`, hanya-baca.
    expect(SRC).toContain('value={EGNPI.satuanIDR} readOnly')
    expect(CSS).toContain('.trin__egnpi-pasangan')
  })

  // ⚠️ Dua desimal BERBEDA di baris yang sama — gambar 29:
  //   Amount 137.849.315.068,00 (2) · Amount in IDR 137.849.315.068 (0)
  it('⚠️ Amount dan Amount in IDR memakai desimal yang berbeda', () => {
    expect(DESIMAL_EGNPI.jumlah).toBe(2)
    expect(DESIMAL_EGNPI.jumlahIDR).toBe(0)
    const html = render('lihat')
    expect(html).toContain('11.080.000.000,00')
    expect(html).toContain('137.849.315.068')
  })

  // ⛔ Aturan B pemilik proses: mengetik tidak menyentuh basis data. Satu-
  // satunya rute yang tab ini panggil ber-awalan `/hitung/`, yang
  // menyatakan dirinya nol tulis.
  it('⛔ nol rumus di layar — seluruhnya lewat /hitung/egnpi', () => {
    expect(SRC).toContain('hitungEgnpi(')
    // Nol aritmetika di berkas layar: tanda kali/bagi hanya boleh muncul
    // di komentar, bukan di kode.
    const kode = SRC.split('\n')
      .filter((b) => !b.trimStart().startsWith('//') && !b.trimStart().startsWith('*'))
      .join('\n')
    expect(kode).not.toContain(' * 100')
    expect(kode).not.toContain('@divide')
  })

  // ⭐ 8 Oktober 2026 — bentuk Pega: grid tetap tampil (kepala + `Add`),
  // dan baris kosongnya satu sel `No items`, seperti tangkapan layar Pega.
  it('⭐ nol baris → grid Pega berisi "No items" dan tombol Add', () => {
    const html = renderToStaticMarkup(
      <TabEgnpi baris={[]} kurs={[]} retensi={[]} opsiAwal={OPSI} mode="ubah" />,
    )
    expect(html).toContain(`<td colSpan="8" class="trin__kosong-pega">${EGNPI.tanpaBaris}</td>`)
    expect(html).toContain(`>${EGNPI.tambah}</button>`)
    expect(html).not.toContain(EGNPI.petunjukKosong)
  })
})

// ⭐ Laporan pemakai 8 Oktober 2026: "set idr atau yg lain tidak bisa malah
// hilang". Konversi yang dipicu pemilihan mata uang berjalan atas baris LAMA
// (keadaan belum dirender) lalu jawabannya menimpa pilihan itu. Kini baris
// BARU dikirim bersama pemicunya.
describe('pilih mata uang Amount — konversi atas baris yang sudah memuat pilihannya', () => {
  it('onPilih membentuk baris baru dan menyerahkannya ke konversi', () => {
    expect(SRC).toContain('const baru = { ...b, Currency: o.nama, CurrencyID: o.id }')
    expect(SRC).toContain('onKonversi(baru)')
    expect(SRC).toContain("jalankan('konversi', i, baru === undefined ? rows : rows.map((x, j) => (j === i ? baru : x)))")
    expect(SRC).not.toMatch(/onUbah\(\{ \.\.\.b, Currency: o\.nama, CurrencyID: o\.id \}\)\s*onKonversi\(\)/)
  })
})
