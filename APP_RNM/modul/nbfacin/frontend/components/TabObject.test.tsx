// Paritas tab Object FIRE (tiket 35) - section Pega `ObjectList` / `Property` / `ObjectDetails` + tangkapan layar
// work owner 03-10-2026. Data uji sintetis.

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import { GRID_OBJEK, OBJECT_TYPE_LAINNYA, OPSI_FLOOR_TYPE, OPSI_ROOF_TYPE, OPSI_WALL_TYPE, SIMPAN_OBJEK, SUBTAB_OBJEK, TEKS_INWARD } from '../labels'
import TabObject, { gantiObjectType, lantaiMinus, objekBaru } from './TabObject'

const HTML = renderToStaticMarkup(<TabObject caseId="NB-1" />)
const SUMBER = readFileSync(join(__dirname, 'TabObject.tsx'), 'utf8').replace(/\r\n/g, '\n')

describe('TabObject - keadaan awal', () => {
  it('kolom grid: (buka), Top Risk, No., Object Name, Location, lalu tombol Tambah', () => {
    const kolom = [...HTML.matchAll(/<th[^>]*>([\s\S]*?)<\/th>|<th[^>]*\/>/g)].map((m) => (m[1] ?? '').replace(/<[^>]+>/g, ''))
    expect(kolom).toEqual(['', ...GRID_OBJEK.kolom.map((k) => k.label), GRID_OBJEK.tambah])
  })

  it('belum ada objek = No items; tombol Save di kanan', () => {
    expect(HTML).toContain(`<td colSpan="6">${TEKS_INWARD.kosong}</td>`)
    expect(HTML).toContain(`<div class="nbf-objek__kaki"><button type="button" class="btn btn--primary">${SIMPAN_OBJEK.label}</button></div>`)
  })
})

describe('TabObject - aturan', () => {
  it('objek baru: Object No. diisi, Material Damage tercentang, Top Risk tidak, Building Construction = Lain-lain (14 / 9 / KELAS I)', () => {
    const o = objekBaru('3')
    expect(o.objectNo).toBe('3')
    expect(o.isMaterialDamage).toBe(true)
    expect(o.isTopRisk).toBe(false)
    expect(o.objectType).toBe('')
    expect([o.roofType, o.wallType, o.floorType]).toEqual(['14', '9', 'KELAS I'])
    expect(OPSI_ROOF_TYPE.find((x) => x.value === o.roofType)?.label).toBe('Lain-lain')
    expect(OPSI_WALL_TYPE.find((x) => x.value === o.wallType)?.label).toBe('Lain-lain')
    expect(OPSI_FLOOR_TYPE.find((x) => x.value === o.floorType)?.label).toBe('Lain-lain')
  })

  it('daftar Roof 14 / Wall 9 / Floor 3 pilihan (tangkapan layar); nilai Roof/Wall = nomor urut (G-9)', () => {
    expect(OPSI_ROOF_TYPE).toHaveLength(14)
    expect(OPSI_WALL_TYPE).toHaveLength(9)
    expect(OPSI_FLOOR_TYPE.map((x) => x.label)).toEqual(['Keramik', 'Kayu', 'Lain-lain'])
    expect(OPSI_ROOF_TYPE[0]).toEqual({ value: '1', label: 'Dak Beton' })
    expect(OPSI_WALL_TYPE[0]).toEqual({ value: '1', label: 'Batu bata' })
  })

  it('Object Type berubah: Object Name = Object Type; Others -> Object Name kosong (G-2)', () => {
    expect(gantiObjectType(objekBaru('1'), 'Dwelling House').objectName).toBe('Dwelling House')
    const lain = gantiObjectType({ ...objekBaru('1'), objectName: 'UJI' }, OBJECT_TYPE_LAINNYA)
    expect(lain.objectType).toBe(OBJECT_TYPE_LAINNYA)
    expect(lain.objectName).toBe('')
  })

  it('Number of Floor minus = galat; kosong / nol / positif sah', () => {
    expect(lantaiMinus('-1')).toBe(true)
    expect(lantaiMinus('0')).toBe(false)
    expect(lantaiMinus('3')).toBe(false)
    expect(lantaiMinus('')).toBe(false)
  })

  it('Tambah: baris di akhir dengan Object No. = urutan, langsung terbuka di Object Address', () => {
    expect(SUMBER).toContain('setBaris((b) => [...b, { kunci, data: objekBaru(String(b.length + 1)) }])')
    expect(SUMBER).toContain('setTerbuka((t) => ({ ...t, [kunci]: SUBTAB_OBJEK[0] }))')
  })

  it('Hapus membuang baris itu; baris terbuka menampilkan tujuh sub-tab, Object Address terisi, lainnya BelumTersedia', () => {
    expect(SUMBER).toContain('setBaris((b) => b.filter((x) => x.kunci !== kunci))')
    expect(SUMBER).toContain('<StripTab tab={SUBTAB_OBJEK}')
    expect(SUMBER).toMatch(/sub === SUBTAB_OBJEK\[0\] \? \(\s*<ObjectAddress/)
    expect(SUBTAB_OBJEK).toHaveLength(7)
  })

  it('Save: Object Type kosong atau lantai minus menahan simpan dan membuka baris salah; sah -> PUT seluruh daftar', () => {
    expect(SUMBER).toMatch(/if \(typeKosong \|\| adaLantaiMinus\) \{[\s\S]*?return\s*\}/)
    expect(SUMBER).toMatch(/simpanObjek\(\s*caseId,\s*baris\.map\(\(x\) => x\.data\),?\s*\)/)
    expect(SUMBER).toContain("galatType={cobaSimpan && x.data.objectType === ''}")
  })

  it('Choose / Clear Risk Address nonaktif (G-1); nilai tersimpan di luar daftar tetap tampil (data lama)', () => {
    expect(SUMBER).toMatch(/disabled>\s*\{A\.chooseRisk\.label\}/)
    expect(SUMBER).toMatch(/disabled>\s*\{A\.clearRisk\.label\}/)
    expect(SUMBER).toContain('opsi={denganTersimpan(OPSI_ROOF_TYPE, o.roofType)} kosong={BANGUNAN_KOSONG}')
  })
})
