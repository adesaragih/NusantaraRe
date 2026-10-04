// Label Company Detail: nama menu = baris M_NAV_MENU migrasi inti 907; caption layar mengikuti screenshot Pega SFAGIS
// (layarnya tidak ada di korpus XML) - field yang dihapus work owner 04-10-2026 tidak boleh muncul lagi.

import { readFileSync } from 'node:fs'

import { describe, expect, it } from 'vitest'

import { CD, MENU_CD } from './labels'

describe('label Company Detail', () => {
  it('nama menu = M_NAV_MENU.LABEL migrasi inti 907', () => {
    const sql = readFileSync(`${__dirname}/../../../inti/backend/migrations/907_m_nav_menu_companydetail.sql`, 'utf8')
    expect(sql).toContain(`'companydetail', '${MENU_CD.kelompok}', 'MASTER', 'companydetail'`)
  })

  it('caption layar Pega Company Detail', () => {
    expect([CD.companyDetail, CD.npwp, CD.parent, CD.country, CD.title, CD.organizationName, CD.businessField, CD.note]).toEqual([
      'Company Detail', 'NPWP', 'Parent organization', 'COUNTRY', 'Title', 'Organization Name', 'Business Field', 'Note',
    ])
    expect([CD.name, CD.position, CD.gender, CD.email, CD.dateOfBirth, CD.phoneNumber]).toEqual([
      'Name', 'Position', 'Gender', 'Email', 'Date of birth', 'Phone number',
    ])
    expect([CD.type, CD.address, CD.phoneAndFax, CD.create]).toEqual(['Type', 'Address', 'Phone and Fax', 'Create'])
  })

  it('field yang dihapus work owner tidak ada di label', () => {
    const teks = JSON.stringify(CD)
    for (const dihapus of ['Client Status', 'Established Date', 'Number of Employees', 'Owner']) {
      expect(teks, dihapus).not.toContain(dihapus)
    }
  })
})
