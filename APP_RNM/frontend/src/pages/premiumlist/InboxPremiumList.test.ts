import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import { KOLOM_INBOX_POLIS, TOMBOL_POLIS } from '../../assets/labels.premiumlist'
import { KOLOM_EKSPOR_POLIS, sel } from './InboxPremiumList'

// Uji kotak masuk PremiumList — tiket 01 bagian 2.

const BERKAS = readFileSync(join(__dirname, 'InboxPremiumList.tsx'), 'utf8')

/**
 * Sumber TANPA komentar.
 *
 * ⛔ Prosa yang MENERANGKAN sebuah kolom - hal biasa di repositori ini -
 * jangan sampai dituduh sebagai kolomnya. Penjaga yang menuduh hal yang
 * benar akan dilonggarkan orang, bukan dipatuhi; pelajaran yang sama sudah
 * dibayar di penjaga nama tabel telanjang sisi Go.
 */
const SUMBER = BERKAS.split('\n')
  .filter((b) => {
    const t = b.trimStart()
    return !t.startsWith('//') && !t.startsWith('*')
  })
  .join('\n')

describe('kolom kotak masuk', () => {
  it('SEBELAS kolom dari InboxPremiumList.xml, plus tanggal terima', () => {
    // Tiga belas terdaftar di RD; dua sengaja tidak ada.
    expect(KOLOM_EKSPOR_POLIS).toHaveLength(12)
  })

  it('label datang dari labels.premiumlist, tidak diketik ulang', () => {
    const label = KOLOM_EKSPOR_POLIS.map((k) => k.label)
    expect(label).toContain(KOLOM_INBOX_POLIS.caseId)
    expect(label).toContain(KOLOM_INBOX_POLIS.plNumber)
    expect(SUMBER).not.toContain("label: 'Case ID'")
  })

  it('KetentuanUnderwriting TIDAK ditampilkan', () => {
    // ⛔ b891 terdaftar di RD tetapi tidak punya kolom di migrasi mana pun.
    // Sel kosong di layar terbaca "memang kosong", bukan "kami tidak punya
    // datanya".
    expect(KOLOM_EKSPOR_POLIS.map((k) => String(k.kunci))).not.toContain(
      'ketentuanUnderwriting',
    )
    expect(SUMBER.toLowerCase()).not.toContain('ketentuan')
  })

  it('pzInsKey TIDAK ditampilkan — pengenal internal Pega', () => {
    expect(SUMBER).not.toContain('pzInsKey')
  })

  it('ekspor memakai kolom tabel yang SAMA', () => {
    // Satu daftar merender kepala tabel, sel, dan ekspor — jadi ketiganya
    // tidak dapat menyimpang.
    expect(SUMBER).toContain('KOLOM_EKSPOR_POLIS.map((k) => (')
  })
})

describe('sel', () => {
  it('kosong ditandai, bukan dibiarkan kosong', () => {
    expect(sel('')).toBe('—')
    expect(sel('   ')).toBe('—')
    expect(sel('QP')).toBe('QP')
  })
})

describe('dua tombol portal', () => {
  it('membuat kasus lewat CreateInputLife, bukan berdiri sebagai BelumTersedia', () => {
    // ⛔ GILIRAN-13 butir bn: SEQ_WORK_POLIS dan FLAG_ONGOING_POLICY (057)
    // kini ada, jadi kedua tombol MEMBUAT kasus (POST /api/polis-life).
    expect(SUMBER).not.toContain('BelumTersedia')
    expect(SUMBER).toContain('await buatKasusPolis(flag)')
    expect(SUMBER).toContain('buat(FLAG_POLIS.inputOffer)')
    expect(SUMBER).toContain('buat(FLAG_POLIS.inputPremium)')
    expect(TOMBOL_POLIS.inputOffer).toBe('Input Offer')
    expect(TOMBOL_POLIS.inputPremium).toBe('Input Premium')
  })

  it('kasus yang lahir langsung DIBUKA di tahapnya', () => {
    // CreateInputLife b982 "ASSIGN-WORKLIST <pzInsKey>!InputPolicyHolder":
    // Pega langsung menyerahkan assignment tahap pertamanya.
    expect(SUMBER).toContain('onBuka(hasil.caseId, hasil.statusWork)')
  })
})
