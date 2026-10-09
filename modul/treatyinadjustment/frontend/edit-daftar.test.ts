// Tombol `Edit` daftar Adjustment — `Section/InputTreatyInAdjustment.xml`
// @782051 (syarat tampil) dan DT `TreatyInSetEdit` (8 Oktober 2026).

import { describe, expect, it } from 'vitest'

import type { Penyesuaian } from './api'
import { editTampil, terapkanSetEdit } from './pages/PenyesuaianKontrak'

describe('Edit @782051 — syarat tampil', () => {
  const admin = ['ReasTreatyInAdmin']
  it('workbasket Admin, Position Admin/kosong, status bukan Decline/Resolve Complete', () => {
    expect(editTampil({ kodePosisi: '', statusAkseptasi: '' }, admin)).toBe(true)
    expect(editTampil({ kodePosisi: 'ReasTreatyInAdmin', statusAkseptasi: 'Reject' }, admin)).toBe(true)
  })
  it('tidak tampil di luar syarat', () => {
    expect(editTampil({ kodePosisi: '', statusAkseptasi: '' }, ['ReasTreatyInSecHead'])).toBe(false)
    expect(editTampil({ kodePosisi: 'ReasTreatyInSecHead', statusAkseptasi: 'Accept' }, admin)).toBe(false)
    expect(editTampil({ kodePosisi: '', statusAkseptasi: 'Resolve Complete' }, admin)).toBe(false)
    expect(editTampil({ kodePosisi: '', statusAkseptasi: 'Decline' }, admin)).toBe(false)
  })
})

describe('DT TreatyInSetEdit', () => {
  const p = (medan: Record<string, string>): Penyesuaian => ({
    id: '1000080/R01',
    idAsal: '1000080',
    baru: { medan, larik: {} },
    lama: { medan: {}, larik: {} },
  })
  it('Position Admin, PositionUsername operator, StatusAkseptasi/Comment kosong, IsEditData 0, ViewState 0', () => {
    const h = terapkanSetEdit(p({ Position: 'ReasTreatyInSecHead', StatusAkseptasi: 'Reject', Comment: 'x', IsEditData: '1' }), 'ADESAMUEL')
    expect(h.baru.medan).toMatchObject({
      ViewState: '0', Position: 'ReasTreatyInAdmin', IsEditData: '0', PositionUsername: 'ADESAMUEL', StatusAkseptasi: '', Comment: '',
    })
  })
  // ⛔ DIBALIK 8 Oktober 2026. Dahulu uji ini menuntut `RevisionState 1 →
  // ViewState 1`, dan tuntutan itu yang mengunci layar pemilik proses.
  //
  // `TreatyInSetEdit[2]` memang bercabang atas `TreatyIn.RevisionState`, tapi
  // properti itu disetel `SetTreatyIn_Act[7]` yang berprasyarat
  // `param.revisionstate==1` — parameter TOMBOL. `Section/InputTreatyInOffer`
  // memasangkannya: Edit mengirim kosong, View mengirim 1. Jadi pada tombol
  // Edit cabang itu TIDAK PERNAH menyala, berapa pun isi kolomnya.
  it('⛔ RevisionState tersimpan TIDAK mengunci — tombol Edit selalu ViewState 0', () => {
    expect(terapkanSetEdit(p({ RevisionState: '1' }), 'X').baru.medan.ViewState).toBe('0')
  })
})
