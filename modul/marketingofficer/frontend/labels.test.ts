// Caption form Pega `InputMarketingOfficer` dibaca dari korpus (baca saja), bukan diketik ulang dari ingatan.
// Mesin tanpa korpus melewati uji ini.

import { existsSync, readFileSync } from 'node:fs'

import { describe, expect, it } from 'vitest'

import { MENU_MO, MO } from './labels'

const SECTION = 'D:\\XML\\RNM_BRD\\NB FacIn\\Section\\InputMarketingOfficer.xml'

describe('label Marketing Officer', () => {
  it.skipIf(!existsSync(SECTION))('caption form Pega VERBATIM ada di section InputMarketingOfficer', () => {
    const xml = readFileSync(SECTION, 'utf8')
    for (const caption of [MO.code, MO.namaMarketing, MO.setLeader, MO.branch, MO.subBranch, MO.leader, MO.active]) {
      expect(xml, caption).toContain(`pyCaption ${caption}<`)
    }
  })

  it('nama menu = M_NAV_MENU.LABEL migrasi inti 906', () => {
    const sql = readFileSync(`${__dirname}/../../../inti/backend/migrations/906_m_nav_menu_marketingofficer.sql`, 'utf8')
    expect(sql).toContain(`'marketingofficer', '${MENU_MO.kelompok}', 'MASTER', 'marketingofficer'`)
  })
})
