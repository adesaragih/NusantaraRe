// Uji panel data polis — A3 kelompok Register.

import { describe, expect, it } from 'vitest'

import { REGISTER } from '../assets/labels'
import { medanPolis } from './PanelDataPolis'

describe('himpunan medan PERSIS section', () => {
  const medan = medanPolis({})

  it('sebelas medan data, tidak kurang', () => {
    // ⛔ Medan yang DIHILANGKAN dari layar sama merusaknya dengan medan yang
    // dikarang - dan yang pertama tidak berbunyi sama sekali. Cacahnya
    // dikunci supaya menghilangkan satu menjadi kegagalan uji.
    expect(medan).toHaveLength(12)
  })

  it('setiap label datang dari REGISTER, bukan diketik ulang', () => {
    const sah = new Set<string>(Object.values(REGISTER))
    for (const m of medan) {
      expect(sah.has(m.label), `label ${m.label}`).toBe(true)
    }
  })

  it('label tidak ganda', () => {
    const lihat = new Set(medan.map((m) => m.label))
    expect(lihat.size).toBe(medan.length)
  })

  it('yang MENUNGGU MODUL dinyatakan, bukan dihilangkan', () => {
    const belum = medan.filter((m) => m.belumBersumber === true)
    // ⛔ SEBELAS, bukan sepuluh: `Type` ikut menunggu sejak ralat
    // pembacaan 27-09-2026 - ia `pyReadOnly` true di Pega dan terikat
    // `.PolicyDataLife.Type`, bukan isian bebas.
    expect(belum).toHaveLength(11)
    for (const m of belum) {
      expect(m.nilai).toBe('—')
    }
  })

  it('yang PUNYA sumber terisi dari nilainya', () => {
    const isi = medanPolis({ namaTertanggung: 'UJI-TERTANGGUNG' })
    const nama = isi.find((m) => m.label === REGISTER.namaTertanggung)
    expect(nama?.nilai).toBe('UJI-TERTANGGUNG')
    expect(nama?.belumBersumber).toBeUndefined()
  })

  it('Type TETAP menunggu walau nilainya diberikan', () => {
    // ⛔ Ia read-only di Pega; nilai yang layar punya hari ini datang
    // dari isian sementara kita, BUKAN dari PolicyDataLife.
    const isi = medanPolis({ type: 'QP' })
    const type = isi.find((m) => m.label === REGISTER.type)
    expect(type?.belumBersumber).toBe(true)
  })

  it('nilai kosong tampil sebagai tanda pisah, bukan string kosong', () => {
    // Sel kosong tidak dapat dibedakan dari sel yang gagal dimuat.
    const isi = medanPolis({ namaTertanggung: '' })
    expect(isi.find((m) => m.label === REGISTER.namaTertanggung)?.nilai).toBe('—')
  })
})
