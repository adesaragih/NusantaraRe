// Fungsi murni tab Object Aneka (tiket 73 P1). Data uji sintetis - tidak ada nama orang, nomor identitas, atau data nasabah.

import { describe, expect, it } from 'vitest'

import type { BarisKapal, TampilanAneka } from '../api'
import {
  deskripsiHe,
  kolomGridObjek,
  kosongkanRiskAneka,
  lokasiAnekaBaru,
  objekAnekaBaru,
  rapikanLokasiAneka,
  rapikanObjekAneka,
  saringDariAneka,
  tampilOkupasi,
  terapkanKapal,
  terapkanRiskAneka,
  ubahHe,
  varianGridObjek,
} from './aneka'
import { objekBaru, terapkanRisk } from './TabObject'

const TAMPIL: TampilanAneka = {
  kontainer: ['C1'], kolomTahunJumlah: false, he: false, golf: false, tanamanHasil: false, engineering: false, bonding: false,
  briguna: false, topRisk: false, sectionLol: false, jumlahSatuan: true, currencyCoverage: false, mbd: false,
}

const RISK = {
  id: 'UJI-R1', title: 'JL', address: 'UJI ALAMAT', nationName: 'UJI NEGARA', provinceName: 'UJI PROV', cityName: 'UJI KOTA',
  districtName: 'UJI KEC', territoryName: 'UJI KEL', postalCode: '00000',
}

describe('aneka - objek & lokasi', () => {
  it('objek baru: TSI awal "0" (sel TSI L2995 pyDefaultValue 0), anak lengkap kosong', () => {
    const o = objekAnekaBaru()
    expect(o.tsi).toBe('0')
    expect(o.marineHull.nrtrt).toBe('')
    expect(o.vehicleHe.typeName).toBe('')
    expect(o.memorandum.suretyType).toBe('')
  })

  it('rapikan: data server tak lengkap dilengkapi; larik okupasi / objek / Golf selalu ada', () => {
    const o = rapikanObjekAneka({ objectName: 'UJI', marineHull: { nrtrt: '12' } as never })
    expect(o.objectName).toBe('UJI')
    expect(o.marineHull.nrtrt).toBe('12')
    expect(o.marineHull.grtrt).toBe('')
    const l = rapikanLokasiAneka({ objectNo: '1', okupasi: [{ occupationId: 'X', occupationName: 'Y' } as never] })
    expect(l.objekGolf).toEqual([])
    expect(l.okupasi[0]?.objek).toEqual([])
  })

  it('Choose / Clear Risk Address = aturan FIRE (SetRiskIdDT_FacIn / ClearRiskLocation_act); Building No. tetap saat Clear', () => {
    const l = { ...lokasiAnekaBaru('1'), buildingNo: '7' }
    const pilih = terapkanRiskAneka(l, RISK)
    const fire = terapkanRisk(objekBaru('1'), RISK)
    expect(pilih.riskLocation).toBe(fire.riskLocation)
    expect(pilih.riskAddressId).toBe('UJI-R1')
    expect(pilih.buildingNo).toBe('')
    const bersih = kosongkanRiskAneka({ ...pilih, buildingNo: '9' })
    expect(bersih.riskLocation).toBe('')
    expect(bersih.buildingNo).toBe('9')
    expect(saringDariAneka(pilih).address).toBe('UJI ALAMAT')
  })
})

describe('aneka - lookup kapal & HE', () => {
  it('pilih kapal: Name/ID/FlagName/YearMake1/NRT/GRT/DWT/Classification/Construction (L10480-10570); medan lain tetap', () => {
    const k: BarisKapal = {
      id: 'S-1', nama: 'UJI KAPAL', flag: 'UJI BENDERA', tahun: '2001', nrt: '10', grt: '20', dwt: '30', klasifikasi: 'K', konstruksi: 'B',
    }
    const awal = { ...objekAnekaBaru(), tsi: '5', marineHull: { ...objekAnekaBaru().marineHull, typeOfVessel: 'TETAP' } }
    const o = terapkanKapal(awal, k)
    expect([o.objectName, o.almRiskId, o.flag, o.year, o.tsi]).toEqual(['UJI KAPAL', 'S-1', 'UJI BENDERA', '2001', '5'])
    expect(o.marineHull).toMatchObject({ nrtrt: '10', grtrt: '20', dwtrt: '30', classification: 'K', construction: 'B', typeOfVessel: 'TETAP' })
  })

  it('SetObjectDescHE_DT L155: ObjectNameHE + " " + TypeName + " Merk : " + BrandName (masing-masing bila ada)', () => {
    const v = { ...objekAnekaBaru().vehicleHe, objectNameHe: 'UJI ALAT' }
    expect(deskripsiHe({ ...v, typeName: 'T1', brandName: 'M1' })).toBe('UJI ALAT T1 Merk : M1')
    expect(deskripsiHe(v)).toBe('UJI ALAT ')
    expect(deskripsiHe({ ...v, brandName: 'M1' })).toBe('UJI ALAT  Merk : M1')
    expect(deskripsiHe({ ...v, typeName: '  ' })).toBe('UJI ALAT ')
    // On-change: ubah medan HE menghitung ulang Object Name.
    expect(ubahHe(objekAnekaBaru(), { objectNameHe: 'A', typeName: 'B' }).objectName).toBe('A B')
  })
})

describe('aneka - grid objek (InputOkupasiAneka_FacIn)', () => {
  it('varian: K4 Growing/Yield menang, lalu K3 HE, selain itu K2', () => {
    expect(varianGridObjek(TAMPIL)).toBe('K2')
    expect(varianGridObjek({ ...TAMPIL, he: true })).toBe('K3')
    expect(varianGridObjek({ ...TAMPIL, he: true, tanamanHasil: true })).toBe('K4')
  })

  it('kolom K2 / K3 / K4; Section + Limit of Liability hanya bila sectionLol (K73-4)', () => {
    expect(kolomGridObjek(TAMPIL)).toEqual(['objectName'])
    expect(kolomGridObjek({ ...TAMPIL, kolomTahunJumlah: true })).toEqual(['objectName', 'year', 'quantity', 'unitQuantity'])
    expect(kolomGridObjek({ ...TAMPIL, he: true })).toEqual(['objectName', 'year', 'chassisNumber', 'engineNumber'])
    expect(kolomGridObjek({ ...TAMPIL, tanamanHasil: true })).toEqual(['objectName', 'currency', 'tsi'])
    expect(kolomGridObjek({ ...TAMPIL, sectionLol: true })).toEqual(['objectName', 'section', 'limitOfLiability'])
  })

  it('grid okupasi: tersembunyi hanya untuk Golf murni (bukan IsAneka)', () => {
    expect(tampilOkupasi(true, TAMPIL)).toBe(true)
    expect(tampilOkupasi(false, { ...TAMPIL, golf: true })).toBe(false)
    expect(tampilOkupasi(true, { ...TAMPIL, golf: true })).toBe(true)
  })
})
