// Paritas baca-saja form produk dengan XML (audit 02-10-2026).
//
// `InboxProductName.xml`: enam medan pemilih master ber-`pyReadOnly` true, `pyEditOptions` Read-only, dan
// `pyReadOnlyCondition` KOSONG = SELALU baca-saja (Ceding b4040, SOB b4428, R/I Risk Name b7362, Cause Of
// Loss b10626, Policy Holder b17062, Currency b28105): nilainya hanya dari daftar master, tidak diketik. Sel
// `PLAN LIST` `Bussines` (`.Name` b33504) dan `Benefit` b33658 juga selalu baca-saja (diisi autocomplete `Plan Name`).
//
// Keputusan work owner 02-10-2026 ("perubahan pada tampilan untuk semua Choose ubah jadi dropdown saja"): tombol
// `Choose*` + popup FlowAction diganti dropdown master `DropdownMaster` - ketujuh master, termasuk R/I Rate baris
// `PLAN LIST`. Nilai tetap HANYA dari daftar master.

import { existsSync, readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import { INWARD_MPNL, UMUM_MPNL } from '../labels'

const KODE = readFileSync(join(__dirname, 'FormProduk.tsx'), 'utf8')

/** Potongan satu elemen `<DropdownMaster ... />` yang memuat `penanda`. */
function dropdown(penanda: string): string {
  const i = KODE.indexOf(penanda)
  if (i < 0) return ''
  const awal = KODE.lastIndexOf('<DropdownMaster', i)
  return KODE.slice(awal, KODE.indexOf('/>', i) + 2)
}

describe('form produk - ketujuh pemilih master berupa dropdown (keputusan work owner 02-10-2026)', () => {
  it('setiap master: label verbatim, jenis RD, mode lihat, dan penerima set*_DT yang menyalin ID + nama', () => {
    const medan = [
      ['UMUM_MPNL.ceding', "jenis=\"ceding\"", 'ubahUmum({ ceding: v.nama, cedingId: v.id })'],
      ['UMUM_MPNL.sob', "jenis=\"sob\"", 'ubahUmum({ sobName: v.nama, sobId: v.id })'],
      ['UMUM_MPNL.riRisk', "jenis=\"ri-risk\"", 'ubahUmum({ riRisk: v.nama, riRiskId: v.id })'],
      ['UMUM_MPNL.causeOfLoss', "jenis=\"penyebab\"", 'ubahUmum({ cause: v.nama, causeId: v.id })'],
      ['INWARD_MPNL.policyHolder', "jenis=\"pemegang-polis\"", 'ubahInward({ policyHolderName: v.nama, policyHolder: v.id })'],
      ['INWARD_MPNL.currency', "jenis=\"mata-uang\"", 'ubahInward({ currency: v.nama, currencyId: v.id })'],
      ['PLAN_MPNL.riRate', "jenis=\"ri-rate\"", 'ganti<BarisPlan>(x.planList, i, { riRate: v.nama, riRateId: v.id })'],
    ] as const
    for (const [label, jenis, pilih] of medan) {
      const d = dropdown(`labelAria={${label}}`)
      expect(d, label).not.toBe('')
      expect(d, label).toContain(jenis)
      expect(d, label).toContain('lihat={lihat}')
      expect(d, label).toContain(pilih)
      // Medan form: dropdown di dalam baris berlabel `Medan` yang sama (mode lihat menampilkan teks namanya).
      if (label !== 'PLAN_MPNL.riRate') {
        const i = letakMedan(label)
        const j = KODE.indexOf('<DropdownMaster', i)
        expect(i, label).toBeGreaterThan(-1)
        expect(KODE.slice(i, j), label).not.toContain('</Medan>')
        expect(KODE.slice(j, KODE.indexOf('/>', j)), label).toContain(`labelAria={${label}}`)
      }
    }
    expect(KODE.split('<DropdownMaster').length - 1).toBe(medan.length)
    // `SetRIRate` 1 b249: kepala kolom nama pemilih R/I Rate `RIRate Name`.
    expect(dropdown('labelAria={PLAN_MPNL.riRate}')).toContain('kolomNama={PEMILIH_MPNL.kolomRiRateName}')
  })

  it('tombol Choose* dan popup pemilih tidak dirender lagi; nilai tetap tidak dapat diketik', () => {
    expect(KODE).not.toMatch(/\{(?:UMUM|INWARD|PLAN)_MPNL\.choose\w*\}/)
    expect(KODE).not.toContain('PemilihMaster')
    expect(KODE).not.toContain('bukaPemilih')
    // Tidak ada isian ketik (autocomplete master) atau dropdown Cause Of Loss yang menerima teks bebas.
    expect(KODE).not.toContain('<Saran<NilaiMaster>')
    expect(KODE).not.toContain("cariMaster('penyebab'")
  })

  it('Policy Holder lewat dropdown tidak menjalankan SetTreatyName_Act (setPolicyHolder_DT b2416)', () => {
    expect(dropdown('labelAria={INWARD_MPNL.policyHolder}')).not.toContain('namaTreaty')
  })

  it('mode lihat: medan master tampil baca-saja seperti sebelumnya, sel R/I Rate PLAN LIST tetap teks', () => {
    const DROPDOWN = readFileSync(join(__dirname, 'DropdownMaster.tsx'), 'utf8')
    expect(readFileSync(join(__dirname, 'DropdownCari.tsx'), 'utf8')).toContain('if (lihat) return <span>{nilai}</span>')
    // Baris form: mode lihat = teks nilai, TIDAK pernah isian atau anak (dropdown, pilihan, area).
    const MEDAN = readFileSync(join(__dirname, 'Medan.tsx'), 'utf8')
    expect(MEDAN).toMatch(/if \(lihat\) \{\s*return \(\s*<div className="mpnl-medan mpnl-medan--lihat">/)
    // Daftar dibaca dengan saringan `Search` dan batas + 1 (potongan dinyatakan).
    expect(DROPDOWN).toContain('cariMaster(jenis, kata, BATAS_DROPDOWN + 1)')
  })

  it('sel Bussines dan Benefit PLAN LIST tampil sebagai teks, tidak dapat diketik', () => {
    expect(KODE).toContain('<td>{b.name}</td>')
    expect(KODE).toContain('<td>{b.benefit}</td>')
    expect(KODE).not.toMatch(/label=\{PLAN_MPNL\.(bussines|benefit)\}/)
  })
})

describe('form produk - tanda wajib (keputusan work owner 02-10-2026)', () => {
  it('empat medan ber-pyRequired true di XML bertanda wajib *', () => {
    // Product Name b3585, Premium Factor (%) b25362, Annuity Interest (%) b26054, Premium Refund Factor (%) b26268.
    for (const label of ['UMUM_MPNL.productName', 'INWARD_MPNL.premiumFactor', 'INWARD_MPNL.annuityInterest', 'INWARD_MPNL.premiumRefundFactor']) {
      const i = KODE.indexOf(`label={${label}}`)
      expect(i, label).toBeGreaterThan(-1)
      const medan = KODE.slice(i, KODE.indexOf('/>', i))
      expect(medan, label).toMatch(/\brequired\b/)
    }
  })
})

// ---------------------------------------------------------------------------
// Tata letak form = wadah XML `InboxProductName` (foto layar Pega work owner 02-10-2026): kartu TREATY NAME (b2934)
// di kiri berisi medan sisi umum, LIEN CLAUSE, DOCUMENT CLAIM; kartu INWARD (b16621) di kanan - Policy Holder,
// Insured, pasangan Addendum/Amandement, lalu kolom kiri b21723 dan kolom kanan b27004.
// ---------------------------------------------------------------------------

const KORPUS = 'D:\\XML\\RNM_BRD\\Master Product Name Life\\Section\\InboxProductName.xml'
const xml = (t: string): string => t.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;')

/** Kunci `OBJEK.k` setiap `<Medan label={OBJEK.k}` di potongan kode, urut kemunculan. */
function urutanMedan(potong: string): string[] {
  return [...potong.matchAll(/<Medan\s+label=\{((?:UMUM|INWARD)_MPNL\.\w+)\}/g)].map((m) => m[1] ?? '')
}

/** Letak `<Medan label={kunci}` di kode - spasi/baris baru di antara `<Medan` dan `label` diterima. */
function letakMedan(kunci: string): number {
  return KODE.search(new RegExp(`<Medan\\s+label=\\{${kunci.replace('.', '\\.')}\\}`))
}

/** Potongan kode di antara dua penanda. */
function antara(dari: string, sampai: string): string {
  const i = KODE.indexOf(dari)
  const j = KODE.indexOf(sampai, i + 1)
  expect(i, dari).toBeGreaterThan(-1)
  expect(j, sampai).toBeGreaterThan(i)
  return KODE.slice(i, j)
}

describe.skipIf(!existsSync(KORPUS))('tata letak form = wadah XML InboxProductName', () => {
  const baris = readFileSync(KORPUS, 'utf8').split('\n').map((b) => b.trim())
  /** Kunci label medan objek yang `pyLabelFieldValue`-nya ada di baris [dari, sampai), urut baris XML. */
  const urutKorpus = (objek: Record<string, string>, nama: string, dari: number, sampai: number): string[] =>
    Object.entries(objek)
      .map(([k, v]) => [`${nama}.${k}`, baris.indexOf(`<pyLabelFieldValue>${xml(v)}</pyLabelFieldValue>`, dari - 1) + 1] as const)
      .filter(([, n]) => n >= dari && n < sampai)
      .sort((a, b) => a[1] - b[1])
      .map(([k]) => k)

  it('kartu TREATY NAME: medan sisi umum urut XML b2934–b12201, lalu LIEN CLAUSE dan DOCUMENT CLAIM', () => {
    const kiri = antara('{UMUM_MPNL.judul}', '{INWARD_MPNL.judul}')
    expect(urutanMedan(kiri)).toEqual(urutKorpus(UMUM_MPNL, 'UMUM_MPNL', 2934, 12201))
    expect(kiri.indexOf('{LIEN_MPNL.judul}')).toBeGreaterThan(kiri.lastIndexOf('<Medan label='))
    expect(kiri.indexOf('{DOKUMEN_MPNL.judul}')).toBeGreaterThan(kiri.indexOf('{LIEN_MPNL.judul}'))
  })

  it('kartu INWARD: baris atas, kolom kiri b21723 dan kolom kanan b27004 urut XML', () => {
    expect(urutanMedan(antara('{INWARD_MPNL.judul}', 'b21723 kolom kiri'))).toEqual(urutKorpus(INWARD_MPNL, 'INWARD_MPNL', 16621, 21723))
    expect(urutanMedan(antara('b21723 kolom kiri', 'b27004 kolom kanan'))).toEqual(urutKorpus(INWARD_MPNL, 'INWARD_MPNL', 21723, 27004))
    expect(urutanMedan(antara('b27004 kolom kanan', 'wadah b30995'))).toEqual(urutKorpus(INWARD_MPNL, 'INWARD_MPNL', 27004, 30995))
    // `OF SUM REASURED` b24880 menempel di Nusantara Re Share.
    const share = letakMedan('INWARD_MPNL.rnmShare')
    expect(share).toBeGreaterThan(-1)
    expect(KODE.slice(share, letakMedan('INWARD_MPNL.rnmLimit'))).toContain('{INWARD_MPNL.ofSumReasured}')
  })
})

describe('jenis tampilan medan = tipe kolom flat (backend mpnl_flat.go)', () => {
  const GO = readFileSync(join(__dirname, '..', '..', 'backend', 'repository', 'mpnl_flat.go'), 'utf8')
  const induk = GO.slice(GO.indexOf('var KolomFlatInduk'), GO.indexOf('type AnakFlat'))
  const kunciLabel = new Map<string, string>([
    ...Object.entries(UMUM_MPNL).map(([k, v]) => [v, `UMUM_MPNL.${k}`] as const),
    ...Object.entries(INWARD_MPNL).map(([k, v]) => [v, `INWARD_MPNL.${k}`] as const),
  ])
  /** Elemen `<Medan label={kunci} ...` sampai `/>` atau `<Medan` berikutnya (yang lebih dulu). */
  const medan = (kunci: string): string => {
    const i = letakMedan(kunci)
    if (i < 0) return ''
    const ujung = [KODE.indexOf('/>', i), KODE.indexOf('<Medan', i + 1)].filter((n) => n > i)
    return KODE.slice(i, Math.min(...ujung))
  }

  it('desimal/bulat tampil angka, tanggal tampil tanggal; kolom teks tidak', () => {
    const kolom = [...induk.matchAll(/\b(teks|desimal|bulat|tanggal)\("[A-Z_]+", (?:\d+, )?"([^"]+)"/g)]
    expect(kolom.length).toBeGreaterThan(40)
    let diperiksa = 0
    for (const [, jenis, label] of kolom) {
      const kunci = kunciLabel.get(label ?? '')
      if (kunci === undefined) continue
      const el = medan(kunci)
      if (el === '') continue
      diperiksa++
      const mau = jenis === 'desimal' || jenis === 'bulat' ? 'jenis="angka"' : jenis === 'tanggal' ? 'jenis="tanggal"' : null
      if (mau === null) expect(el, kunci).not.toMatch(/jenis="(angka|tanggal)"/)
      else expect(el, kunci).toContain(mau)
    }
    // Setiap kolom induk berlabel yang dirender sebagai `Medan` ikut diperiksa (tidak ada yang terlewat diam-diam).
    expect(diperiksa).toBe([...kolom].filter(([, , l]) => kunciLabel.has(l ?? '') && letakMedan(kunciLabel.get(l ?? '') ?? '') > -1).length)
    expect(diperiksa).toBeGreaterThan(35)
  })
})


// ---------------------------------------------------------------------------
// Keputusan work owner 03-10-2026: "jika view tidak tambah/edit/delete, saat klik edit baru bisa" - di mode lihat
// SEMUA aksi ubah tersembunyi (menyimpang dari XML: `Add`/`Delete` grid dan checkbox `On Retention` tidak ber-`ro`).
// ---------------------------------------------------------------------------

const LAMPIRAN = readFileSync(join(__dirname, 'PanelLampiran.tsx'), 'utf8')

/**
 * Teks `{label}` berada DI DALAM satu blok `{!lihat && …}` (hanya dirender di mode sunting): dari setiap `{!lihat && `
 * sebelum label, kurung kurawalnya belum tertutup saat label dicapai.
 */
function dijagaLihat(kode: string, label: string): boolean {
  const i = kode.indexOf(`{${label}}`)
  if (i < 0) return false
  for (let j = kode.lastIndexOf('{!lihat && ', i); j >= 0; j = kode.lastIndexOf('{!lihat && ', j - 1)) {
    let dalam = 0
    let tertutup = false
    for (let k = j; k < i && !tertutup; k++) {
      if (kode[k] === '{') dalam++
      else if (kode[k] === '}' && --dalam === 0) tertutup = true
    }
    if (!tertutup) return true
  }
  return false
}

describe('mode lihat: nol tambah/edit/delete - baru ada sesudah Edit (keputusan work owner 03-10-2026)', () => {
  it('PLAN LIST mode lihat: sel Plan Name teks (dropdown merender teks di mode lihat)', () => {
    const i = KODE.indexOf('<DropdownCari<JenisPlan>')
    expect(KODE.slice(i, KODE.indexOf('/>', KODE.indexOf('onPilih=', i)))).toContain('lihat={lihat}')
  })

  it('PLAN LIST: Add dan Delete hanya di mode sunting; View Rate tetap', () => {
    expect(dijagaLihat(KODE, 'PLAN_MPNL.add')).toBe(true)
    expect(dijagaLihat(KODE, 'PLAN_MPNL.delete')).toBe(true)
    expect(dijagaLihat(KODE, 'PLAN_MPNL.viewRate')).toBe(false)
  })

  it('FINANCIAL UNDERWRITING / UNDERWRITING LIMIT: Add, Copy row, Delete hanya di mode sunting', () => {
    expect(KODE).toContain('{!lihat && <div className="aksi-baris">{tambah}</div>}')
    expect(dijagaLihat(KODE, 'LAIN_MPNL.salinBaris')).toBe(true)
    expect(KODE).toMatch(/\{!lihat && \(\s*<td className="table__actions">[\s\S]*?\{hapus\(i\)\}\s*<\/td>\s*\)\}/)
  })

  it('checkbox On Retention mati di mode lihat; tombol Copy hanya di mode sunting; Generate tetap', () => {
    expect(KODE).toMatch(/checked=\{u\.isOrs\}\s*disabled=\{lihat\}/)
    expect(dijagaLihat(KODE, 'TOMBOL_MPNL.copy')).toBe(true)
    expect(dijagaLihat(KODE, 'TOMBOL_MPNL.generate')).toBe(false)
  })

  it('lampiran: Add attachment, Retry, Delete hanya di mode sunting; Refresh dan Download All tetap', () => {
    expect(KODE).toContain('<PanelLampiran produkId={p.id} lihat={lihat} />')
    for (const k of ['LAMPIRAN_MPNL.add', 'LAMPIRAN_MPNL.delete', 'LAIN_MPNL.ulangi']) expect(dijagaLihat(LAMPIRAN, k), k).toBe(true)
    for (const k of ['LAMPIRAN_MPNL.refresh', 'LAMPIRAN_MPNL.downloadAll']) expect(dijagaLihat(LAMPIRAN, k), k).toBe(false)
  })
})

describe('Document List dipilih dari daftar (permintaan work owner 03-10-2026)', () => {
  it('kolom Document List DOCUMENT CLAIM membawa PILIHAN_DOKUMEN_KLAIM', () => {
    expect(KODE).toContain("kolom={[[DOKUMEN_MPNL.documentList, 'document', PILIHAN_DOKUMEN_KLAIM]]}")
  })

  it('sel berpilihan di mode sunting = pilihan (nilai di luar daftar tetap tampil, tidak dibuang); mode lihat = teks', () => {
    const grid = KODE.slice(KODE.indexOf('function GridSederhana'), KODE.indexOf('function GridBerangka'))
    expect(grid).toMatch(/opsi !== undefined && !lihat \? \(\s*<>\s*<PilihanMedan/)
    expect(grid).toContain('<SelIsi')
    // Nama utuh nilai terpilih tampil di bawah kotak pilihan (kotak bawaan tidak membungkus teks panjang).
    expect(grid).toContain('<div className="mpnl-sel-teks">')
    const MEDAN = readFileSync(join(__dirname, 'Medan.tsx'), 'utf8')
    expect(MEDAN).toContain("const asing = value !== '' && !opsi.some((o) => o.value === value)")
  })
})

describe('Plan Name = dropdown (keputusan work owner 03-10-2026 "tolong ubah jadi model dropdown")', () => {
  // Menyimpang dari XML b33121 (pxAutoComplete, isian bebas b33137): nilai hanya dari daftar `BrowseProductTypeLife_RD`.
  const i = KODE.indexOf('<DropdownCari<JenisPlan>')
  const dropdown = KODE.slice(i, KODE.indexOf('/>', KODE.indexOf('onPilih=', i)))

  it('kolom daftar = pyAdditionalFields XML: ID, CoverName, Business, Benefit', () => {
    expect(i).toBeGreaterThan(-1)
    expect(dropdown).toMatch(
      /judul: \[SARAN_PLAN_MPNL\.kolomId, SARAN_PLAN_MPNL\.kolomCoverName, SARAN_PLAN_MPNL\.kolomBusiness, SARAN_PLAN_MPNL\.kolomBenefit\],\s*isi: \(t\) => \[t\.id, t\.coverName, t\.business, t\.benefit\]/,
    )
  })

  it('memilih: .CoverName → .Plan, .ID → .PlanID, .Business → .Name, .Benefit → .Benefit (b33216, b33250, b33283, b33315)', () => {
    expect(dropdown).toContain('{ plan: t.coverName, planId: t.id, name: t.business, benefit: t.benefit }')
  })

  it('tidak ada lagi isian ketik bebas (autocomplete Saran dihapus)', () => {
    expect(KODE).not.toContain('<Saran')
    expect(existsSync(join(__dirname, 'Saran.tsx'))).toBe(false)
  })
})

