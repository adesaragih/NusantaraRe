import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import { BUAT_EDM, INBOX_EDM } from '../labels'
import { KOLOM_INBOX } from './InboxEndorsementLife'

// Urutan kolom kotak masuk - `Create Date` tepat sesudah `Case ID` (keputusan work owner 07-10-2026).

const SUMBER = readFileSync(join(__dirname, 'InboxEndorsementLife.tsx'), 'utf8')

describe('kolom kotak masuk Endorsement Life', () => {
  it('urutan tampilan: 9 kolom, Create Date sesudah Case ID; empat kolom dihapus (07-10-2026)', () => {
    expect(KOLOM_INBOX).toEqual([
      INBOX_EDM.kolomCaseId,
      INBOX_EDM.kolomCreateDate,
      INBOX_EDM.kolomEndorsementNo,
      INBOX_EDM.kolomType,
      INBOX_EDM.kolomEdmType,
      INBOX_EDM.kolomPolicyNo,
      INBOX_EDM.kolomCeding,
      INBOX_EDM.kolomCreateOperator,
      INBOX_EDM.kolomStatus,
    ])
    // Status → Work Status (07-10-2026), isinya tetap STATUSS.
    expect(INBOX_EDM.kolomStatus).toBe('Work Status')
    // Policy No → PL_NUMBER (07-10-2026), isinya tetap OLD_POLICY_NO; label layar lain tetap "Policy No".
    expect(INBOX_EDM.kolomPolicyNo).toBe('PL_NUMBER')
    expect(BUAT_EDM.policyNo).toBe('Policy No')
    for (const buang of ['SOB', 'Policy Holder', 'Marketing Name', 'EDM Type Batal']) {
      expect(KOLOM_INBOX).not.toContain(buang)
    }
    expect(SUMBER).toContain('keputusan work owner 07-10-2026')
  })

  it('sel <td> sejajar dengan kepala (9 sel, urutan sama)', () => {
    const baris = SUMBER.slice(SUMBER.indexOf('<tr key={b.caseId}>'), SUMBER.indexOf('</tr>', SUMBER.indexOf('<tr key={b.caseId}>')))
    const urut = [
      'onBuka(b.caseId)',
      'selTanggal(b.createDate)',
      'sel(b.endorsementNo)',
      'sel(b.type)',
      'labelEdmType(b.edmType)',
      'sel(b.policyNo)',
      'sel(b.ceding)',
      'sel(b.createOperator)',
      'sel(b.status)',
    ]
    let dari = 0
    for (const jejak of urut) {
      const i = baris.indexOf(jejak, dari)
      expect(i, jejak).toBeGreaterThanOrEqual(dari)
      dari = i + jejak.length
    }
    expect(baris.split('<td>').length - 1).toBe(KOLOM_INBOX.length)
  })
})
