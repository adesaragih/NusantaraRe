// Label Aggregate: nama menu = baris M_NAV_MENU migrasi inti 911; caption layar mengikuti section Pega folder korpus
// `Aggregate` (`GridDasbordAgg`, `ShowAggregateList`, `ChooseMasterID`).

import { readFileSync } from 'node:fs'

import { describe, expect, it } from 'vitest'

import { KOLOM_GRID } from './aturan'
import { AG, LABEL_KOLOM, MENU_AG } from './labels'

describe('label Aggregate', () => {
  it('nama menu = M_NAV_MENU.LABEL migrasi inti 911, golongan MASTER TREATY', () => {
    const sql = readFileSync(`${__dirname}/../../../inti/backend/migrations/911_m_nav_menu_aggregate.sql`, 'utf8')
    expect(sql).toContain(`'aggregate', '${MENU_AG.kelompok}', 'MASTER TREATY', 'aggregate'`)
  })

  it('caption daftar GridDasbordAgg', () => {
    expect([AG.addData, AG.tanggalInput, AG.cedingCode, AG.cedingName, AG.treatyType, AG.asAt, AG.uwYear, AG.hapus]).toEqual([
      'Add Data', 'Tanggal Input', 'Ceding Code', 'Ceding Name', 'Treaty Type', 'As At', 'UW Year', 'Delete',
    ])
  })

  it('tombol ShowAggregateList dan kolom ChooseMasterID', () => {
    expect([AG.masterId, AG.template, AG.uploadCsv, AG.save, AG.close]).toEqual(['Master ID', 'Template', 'Upload CSV', 'Save', 'Close'])
    expect([AG.treatyContractName, AG.reinsuranceType, AG.ceding, AG.sob, AG.treatyGroup, AG.rnmShare, AG.treatyYear]).toEqual([
      'Treaty Contract Name', 'Reinsurance Type', 'Ceding', 'SOB', 'Treaty Group', 'RNM Share', 'Treaty Year',
    ])
  })

  it('setiap kolom grid TempCSV punya caption; caption khas Pega dipertahankan apa adanya', () => {
    for (const k of KOLOM_GRID) expect(LABEL_KOLOM[k.nama], k.nama).toBeTruthy()
    expect(LABEL_KOLOM['ASSESMENT_ZONE']).toBe('Assesment Zone')
    expect(LABEL_KOLOM['CEDING_NAME']).toBe('Ceding Names')
    expect(LABEL_KOLOM['TO_USD']).toBe('To_USD')
    expect([LABEL_KOLOM['NOR_BUILDINGS'], LABEL_KOLOM['BUILDINGS']]).toEqual(['NoR Buildings', 'IA Buildings'])
    expect([LABEL_KOLOM['TOTAL_IN_AMOUNT_IN_USD'], LABEL_KOLOM['RNM_VALUE_IN_USD']]).toEqual(['Total In Amount (USD)', 'RNM Value (USD)'])
  })
})
