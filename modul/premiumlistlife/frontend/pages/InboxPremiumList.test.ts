import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import { KOLOM_INBOX_POLIS, TOMBOL_POLIS } from '../labels'
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
  it('DELAPAN kolom: sebelas dari InboxPremiumList.xml + tanggal terima, dikurangi empat', () => {
    // Tiga belas terdaftar di RD; dua sengaja tidak ada. Sejak 02-10-2026
    // PolicyHolderName, MarketingName, SobName, dan DateReceived juga tidak
    // ditampilkan (keputusan work owner).
    expect(KOLOM_EKSPOR_POLIS).toHaveLength(8)
    const kunci = KOLOM_EKSPOR_POLIS.map((k) => String(k.kunci))
    for (const k of ['policyHolderName', 'marketingName', 'sobName', 'dateReceived']) {
      expect(kunci).not.toContain(k)
    }
  })

  it('urutan kolom keputusan work owner 02-10-2026', () => {
    expect(KOLOM_EKSPOR_POLIS.map((k) => String(k.kunci))).toEqual([
      'caseId',
      'tglCreate',
      'plNumber',
      'riSlipRnm',
      'type',
      'cedingCoName',
      'createOpName',
      'statusWork',
    ])
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

describe('kasus tertutup tidak dapat dibuka (02-10-2026)', () => {
  it('hanya tiga tahap aktif yang dapat dibuka', async () => {
    const { kasusBisaDibuka } = await import('../api')
    expect(kasusBisaDibuka('Input Offer Life')).toBe(true)
    expect(kasusBisaDibuka('Input Premium Detail')).toBe(true)
    expect(kasusBisaDibuka('Input Premium Summary')).toBe(true)
    expect(kasusBisaDibuka('Resolved-Completed')).toBe(false)
    expect(kasusBisaDibuka('Resolved-Rejected')).toBe(false)
    expect(kasusBisaDibuka('')).toBe(false)
  })

  it('baris tertutup tanpa onClick, dan rute tidak merender layar keputusannya', () => {
    expect(SUMBER).toContain('kasusBisaDibuka(b.statusWork) ? (')
    expect(SUMBER).toContain('className="pl-inbox__tutup"')
    const rute = readFileSync(join(__dirname, '..', 'rute.tsx'), 'utf8')
    expect(rute).toContain("polis.id !== '' && kasusBisaDibuka(polis.tahap) && (")
  })
})
