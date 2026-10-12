// Uji Inbox Claim Life — F0.4.
//
// Yang diuji: tab mana yang tampil untuk peran mana, label kolom yang
// VERBATIM korpus, dan bahwa nomor tahap sejalan dengan `models.Tahap` di Go.

import { existsSync, readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import { PERAN } from '../../../../inti/frontend/labels'
import { TAHAP } from '../labels'
import { TAHAP_NOMOR } from '../api'
import { CARI_INBOX, JEDA_CARI_MS, KOLOM_EKSPOR, tabTampil, tabUntuk, tanggalTampil } from './InboxClaimLife'

const KORPUS = 'D:\\XML\\RNM_BRD\\Claim Life'
const adaKorpus = existsSync(KORPUS)
const SUMBER = readFileSync(join(__dirname, 'InboxClaimLife.tsx'), 'utf8')
const TAHAP_GO = join(__dirname, '..', '..', 'backend', 'models', 'tahap.go')

describe('tab per peran', () => {
  it('Admin melihat KEDUA antrian Admin', () => {
    // `Register_Flow.xml`: Input Register 1508 `Current operator`,
    // Outstanding 1351 `ToWorklist` — keduanya ReasLifeAdmin.
    const t = tabUntuk([PERAN.admin])
    expect(t.map((x) => x.judul)).toEqual([TAHAP.inputRegister, TAHAP.outstandingClaim])
  })

  it('Medical Advisor hanya Medical Check', () => {
    expect(tabUntuk([PERAN.medis]).map((x) => x.judul)).toEqual([TAHAP.medicalCheck])
  })

  it('SPV hanya Claim Analis', () => {
    expect(tabUntuk([PERAN.spv]).map((x) => x.judul)).toEqual([TAHAP.claimAnalis])
  })

  it('peran GANDA melihat gabungan, berurut tangga', () => {
    // ⛔ Bukan salah satu. Layar yang memaksa memilih satu peran
    // menyembunyikan separuh pekerjaan orang itu.
    const t = tabUntuk([PERAN.spv, PERAN.admin])
    expect(t.map((x) => x.judul)).toEqual([
      TAHAP.inputRegister,
      TAHAP.outstandingClaim,
      TAHAP.claimAnalis,
    ])
  })

  it('tanpa peran: nol tab, bukan tab kosong', () => {
    // Tab kosong yang tidak pernah dapat berisi hanyalah tempat orang
    // menunggu sia-sia.
    expect(tabUntuk([])).toHaveLength(0)
  })
})

describe('nomor tahap sejalan dengan models.Tahap di Go', () => {
  it('keempatnya 1..4 berurut tangga', () => {
    // ⛔ Angka di kabel. Bila urutan `models.Tahap` di Go berubah, tab akan
    // meminta antrian yang salah - dan layarnya tetap tampil normal.
    expect(TAHAP_NOMOR).toEqual({
      inputRegister: 1,
      outstanding: 2,
      medicalCheck: 3,
      claimAnalis: 4,
    })
  })

  // Satu modul, satu folder (30-09-2026): konstanta Go modul ini tinggal di
  // `backend/` di sebelah folder frontend-nya - jalur relatif, bukan absolut.
  it.runIf(existsSync(TAHAP_GO))(
    'urutannya sama dengan konstanta Go',
    () => {
      const go = readFileSync(TAHAP_GO, 'utf8')
      // `TahapTidakDikenal = iota` (0), lalu keempatnya berurut.
      const urut = ['TahapInputRegister', 'TahapOutstanding', 'TahapMedicalCheck', 'TahapClaimAnalis']
      const letak = urut.map((n) => go.indexOf(`\t${n}\n`))
      expect(letak.every((i) => i > 0)).toBe(true)
      for (let i = 1; i < letak.length; i++) {
        expect(letak[i]!).toBeGreaterThan(letak[i - 1]!)
      }
    },
  )
})

describe('label kolom VERBATIM korpus', () => {
  const kolom: Array<[string, number, string]> = [
    ['Case ID', 721, 'caseId'],
    ['Create Date/Time', 736, 'tglCreate'],
    ['Create Operator Name', 751, 'createOpName'],
    ['Work Status', 765, 'status'],
  ]

  it.runIf(adaKorpus).each(kolom)(
    'label %s ada di InboxPremiumList.xml baris %i',
    (label, baris) => {
      const rd = readFileSync(join(KORPUS, 'ReportDefinition', 'InboxPremiumList.xml'), 'utf8')
      expect(rd.split('\n')[baris - 1] ?? '').toContain(`<pyFieldLabel>${label}</pyFieldLabel>`)
    },
  )

  it.each(kolom)('layar memakai label %s apa adanya', (label) => {
    expect(SUMBER).toContain(`'${label}'`)
  })

  it.runIf(adaKorpus)('Claim No VERBATIM InputOSClaimLife.xml:951', () => {
    const sec = readFileSync(join(KORPUS, 'Section', 'InputOSClaimLife.xml'), 'utf8')
    expect(sec.split('\n')[950] ?? '').toContain('Claim No')
    expect(SUMBER).toContain("'Claim No'")
  })
})

describe('halaman dan urutan dari XML', () => {
  it.runIf(adaKorpus)('50 per halaman, maksimum 500, urut DESC', () => {
    const rd = readFileSync(join(KORPUS, 'ReportDefinition', 'InboxPremiumList.xml'), 'utf8')
    const b = rd.split('\n')
    expect(b[591] ?? '').toContain('<pyPageSize>50</pyPageSize>')
    expect(b[941] ?? '').toContain('<pyMaxRecords>500</pyMaxRecords>')
    expect(b[732] ?? '').toContain('<pySortType>DESC</pySortType>')
  })
})

// Ekspor xlsx — butir bg, 28-09-2026.
describe('ekspor xlsx hanya kolom yang tampil', () => {
  it('kolom ekspor SAMA dan seurutan dengan kolom tabel', () => {
    // ⛔ Ekspor yang memuat kolom yang layarnya tidak tampilkan membuat
    // berkas dan layar menjawab pertanyaan yang berbeda — dan yang
    // memegang berkasnya tidak akan tahu mana yang benar.
    // Urutan keputusan work owner 11-10-2026.
    const urutTabel = ['caseId', 'nomorKlaim', 'plNumber', 'createOpName', 'updateOpName', 'tglCreate', 'status']
    expect(KOLOM_EKSPOR.map((k) => k.kunci)).toEqual(urutTabel)
  })

  it('label ekspor dari KOLOM, tidak diketik ulang', () => {
    // Label yang sama di dua tempat adalah dua tempat untuk bergeser.
    const sumber = readFileSync(join(__dirname, 'InboxClaimLife.tsx'), 'utf8')
    expect(sumber).toContain('label: KOLOM.caseId')
    expect(sumber).not.toContain("label: 'Case ID'")
  })

  it('tombol ekspor mati saat nol baris', () => {
    // Tombol yang menghasilkan berkas kosong mengajari orang mengabaikannya.
    const sumber = readFileSync(join(__dirname, 'InboxClaimLife.tsx'), 'utf8')
    expect(sumber).toContain('disabled={hal === null || hal.baris.length === 0}')
  })
})

describe('inbox dirapikan (keputusan work owner 10-10-2026)', () => {
  const SUMBER_INBOX = readFileSync(join(__dirname, 'InboxClaimLife.tsx'), 'utf8').replace(/\r\n/g, '\n')

  it('tab Input Register disembunyikan selama kosong; muncul lagi bila ada kasus Send Back to Register', () => {
    const admin = tabUntuk([PERAN.admin])
    expect(tabTampil(admin, {}).map((t) => t.judul)).toEqual([TAHAP.outstandingClaim])
    expect(tabTampil(admin, { [TAHAP_NOMOR.inputRegister]: 0 }).map((t) => t.judul)).toEqual([TAHAP.outstandingClaim])
    expect(tabTampil(admin, { [TAHAP_NOMOR.inputRegister]: 2 }).map((t) => t.judul)).toEqual([TAHAP.inputRegister, TAHAP.outstandingClaim])
    const semua = tabUntuk([PERAN.admin, PERAN.medis, PERAN.spv])
    expect(tabTampil(semua, {}).map((t) => t.judul)).toEqual([TAHAP.outstandingClaim, TAHAP.medicalCheck, TAHAP.claimAnalis])
  })

  it('tab bawaan = tab tampil pertama (Outstanding Claim); label Inggris saja, badge di setiap tab', () => {
    expect(SUMBER_INBOX).toContain('useState<NomorTahap | null>(tabTampil(semuaTab, {})[0]?.nomor ?? null)')
    expect(SUMBER_INBOX).not.toContain('judulID')
    expect(SUMBER_INBOX).not.toContain('inbox__tab-id')
    expect(SUMBER_INBOX).toContain('{jumlah[t.nomor] !== undefined && <span className="inbox__lencana">{jumlah[t.nomor]}</span>}')
  })

  it('teks kosong dan paginasi seperti inbox lain; Register primary + Export xlsx ghost', () => {
    expect(SUMBER_INBOX).toContain('<Kosong pesan="No cases in this queue." />')
    expect(SUMBER_INBOX).toContain('<Halaman halaman={hal.halaman} ukuran={hal.ukuran} total={hal.total} onPindah={setHalaman} />')
    for (const id of ['Sebelumnya', 'Berikutnya', 'Tidak ada kasus', 'Tidak ada antrian']) expect(SUMBER_INBOX).not.toContain(id)
    expect(SUMBER_INBOX).toContain('className="btn btn--primary" onClick={onRegister}')
    expect(SUMBER_INBOX).toContain('className="btn btn--ghost"')
    expect(SUMBER_INBOX).toContain('Export xlsx')
  })

  it('waktu buat dd/mm/yyyy HH:mm; kosong tanda pisah', () => {
    expect(tanggalTampil('2026-10-10T07:05:00')).toBe('10/10/2026 07:05')
    expect(tanggalTampil('')).toBe('—')
  })
})

describe('perbaikan DEV 10-10-2026', () => {
  it('Claim No kosong tampil tanda pisah (nomor lahir di Save to RNM); modul berbahasa UI Inggris seperti endorsementlife', () => {
    const src = readFileSync(join(__dirname, 'InboxClaimLife.tsx'), 'utf8')
    expect(src).toContain("<td>{b.nomorKlaim === '' ? '—' : b.nomorKlaim}</td>")
    const rute = readFileSync(join(__dirname, '..', 'rute.tsx'), 'utf8')
    expect(rute).toContain('<BahasaUI.Provider value="en">')
    expect(rute).toContain('<div className="claimlife">')
  })
})

describe('kolom PL Number dan Update Operator Name (11-10-2026)', () => {
  it('label, urutan header tabel = urutan ekspor, Policy No hilang', () => {
    expect(KOLOM_EKSPOR.map((k) => k.label)).toEqual([
      'Case ID', 'Claim No', 'PL Number', 'Create Operator Name', 'Update Operator Name', 'Create Date/Time', 'Work Status',
    ])
    const sumber = readFileSync(join(__dirname, 'InboxClaimLife.tsx'), 'utf8').replace(/\r\n/g, '\n')
    const kepala = [...sumber.matchAll(/<th>\{KOLOM\.([a-zA-Z]+)\}<\/th>/g)].map((m) => m[1])
    expect(kepala).toEqual(KOLOM_EKSPOR.map((k) => k.kunci))
    expect(sumber).not.toContain("'Policy No'")
    expect(sumber).toContain("<td>{b.plNumber === '' ? '—' : b.plNumber}</td>")
  })
})

describe('kotak cari inbox (work owner 12-10-2026)', () => {
  const src = readFileSync(join(__dirname, 'InboxClaimLife.tsx'), 'utf8').replace(/\r\n/g, '\n')
  it('placeholder, jeda 300 ms, kata kunci ke server per antrian aktif, kembali ke halaman 1, tombol Clear', () => {
    expect(CARI_INBOX).toBe('Search case ID, claim no, PL number, insured name')
    expect(JEDA_CARI_MS).toBe(300)
    expect(src).toContain('setHal(await ambilKotakMasuk(t, h, undefined, c))')
    expect(src).toContain('void muat(aktif, halaman, cari)')
    expect(src).toContain('if (lama !== k) setHalaman(1)')
    expect(src).toContain("onClick={() => setKata('')}")
    expect(src).toContain('<input\n          className="field__input claimlife__cari"')
  })
  it('badge tetap jumlah seluruh antrian; Export xlsx = baris yang tampil (ikut kata kunci)', () => {
    expect(src).toContain("if (aktif !== null && hal !== null && cari === '')")
    // Hitungan badge tanpa kata kunci.
    expect(src).toContain('void ambilKotakMasuk(t.nomor, 1)')
    expect(src).toContain('KOLOM_EKSPOR, hal.baris)')
    expect(src).toContain('<Kosong pesan="No cases match the search." />')
  })
})

