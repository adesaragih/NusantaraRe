// Tombol salin di bawah grid coverage item FIRE (tiket 53) - port `NB FacIn\Section\InputCoverageFire.xml`
// (`ASM-FW-GISFW-DATA-PROPERTYITEM!INPUTCOVERAGEFIRE`) dan activity-nya. Fungsi murni atas daftar objek di layar; premi
// coverage tujuan dihitung ulang sesudahnya lewat `hitung-coverage` (CountPremi_ACT percent / percent).
//
// - Copy Coverage (`CopyCoverageObjItem_Act`): sumber = item ke-1 lokasi itu; tujuan = SEMUA item lokasi itu (item 1
//   ikut); CoverageList tujuan DIGANTI (Property-Remove lalu Page-Copy); TotalNetRate / FlagNetRate disalin (langkah 6.4).
// - Copy Coverage To All (`CopyCoverageToAllLocation_Act`): sumber = lokasi 1 / item 1; tujuan = semua item semua lokasi;
//   diganti; TotalNetRate disalin. (Pega langsung Obj-Save; di sini tersimpan lewat Save tab - keputusan agent K53-2.)
// - Copy Coverage From -> "Copy to Another Object Item" (`CopyCoverage_Act`): sumber = lokasi / item yang dipilih di
//   popup; tujuan = item yang sedang dibuka; diganti; TotalNetRate TIDAK disalin.
// - Copy Deductible To All / For BI (`CopyDeductibleToAllLocation_Act`, Method "ALL" / "BI"): tiap coverage tujuan,
//   DeductibleList DIGANTI dengan deductible coverage sumber berkode `.Coverage` sama (tanpa padanan -> kosong);
//   "ALL" melewati item Business Interruption 10000992, "BI" hanya item 10000992. Premi tidak dihitung ulang.
//
// Keputusan agent K53-1: TSI coverage tujuan = TSI Object Item tujuan (maksud langkah 6.6.1 `CopyCoverageObjItem_Act`);
// Pega `CountPremi_ACT` langkah 10 dapat memakai TSI item SUMBER karena indeks salinan ikut tersalin `[dugaan]`.

import type { CoverageObjek, ObjekFire } from '../api'

/** Kode jenis item Business Interruption (`CopyDeductibleToAllLocation_Act` "jangan / hanya copy ke Business Interuption"). */
export const JENIS_ITEM_BI = '10000992'

/** Posisi item: indeks lokasi (objek) dan indeks item. */
export interface PosisiItem {
  o: number
  i: number
}

const salinan = <T>(v: T): T => JSON.parse(JSON.stringify(v)) as T

/** Coverage sumber untuk item tujuan: salinan utuh, TSI = TSI Object Item tujuan (K53-1). */
function untukItem(sumber: CoverageObjek[], tsiTujuan: string): CoverageObjek[] {
  return sumber.map((c) => ({ ...salinan(c), tsi: tsiTujuan }))
}

/** Hasil salin: daftar objek baru + item tujuan yang premi coverage-nya perlu dihitung ulang. */
export interface HasilSalin {
  objek: ObjekFire[]
  tujuan: PosisiItem[]
}

/** Copy Coverage: item 1 lokasi `o` -> semua item lokasi `o`. */
export function salinCoverageLokasi(objek: ObjekFire[], o: number): HasilSalin {
  const lok = objek[o]
  const sumber = lok?.items[0]
  if (!lok || !sumber) return { objek, tujuan: [] }
  const cov = sumber.coverages ?? []
  const tujuan: PosisiItem[] = []
  const baru = objek.map((x, a) =>
    a !== o
      ? x
      : {
          ...x,
          items: x.items.map((it, i) => {
            tujuan.push({ o, i })
            return { ...it, coverages: untukItem(cov, it.tsi), totalNetRate: sumber.totalNetRate }
          }),
        },
  )
  return { objek: baru, tujuan }
}

/** Copy Coverage To All: lokasi 1 / item 1 -> semua item semua lokasi. */
export function salinCoverageSemua(objek: ObjekFire[]): HasilSalin {
  const sumber = objek[0]?.items[0]
  if (!sumber) return { objek, tujuan: [] }
  const cov = sumber.coverages ?? []
  const tujuan: PosisiItem[] = []
  const baru = objek.map((x, o) => ({
    ...x,
    items: x.items.map((it, i) => {
      tujuan.push({ o, i })
      return { ...it, coverages: untukItem(cov, it.tsi), totalNetRate: sumber.totalNetRate }
    }),
  }))
  return { objek: baru, tujuan }
}

/** Copy Coverage From: item `dari` -> item `ke` (diganti; TotalNetRate tidak disalin). */
export function salinCoverageDari(objek: ObjekFire[], dari: PosisiItem, ke: PosisiItem): HasilSalin {
  const sumber = objek[dari.o]?.items[dari.i]
  const target = objek[ke.o]?.items[ke.i]
  if (!sumber || !target) return { objek, tujuan: [] }
  const cov = untukItem(sumber.coverages ?? [], target.tsi)
  const baru = objek.map((x, o) =>
    o !== ke.o ? x : { ...x, items: x.items.map((it, i) => (i !== ke.i ? it : { ...it, coverages: cov })) },
  )
  return { objek: baru, tujuan: [ke] }
}

/** Copy Deductible To All ("ALL") / For BI ("BI"): deductible per kode coverage dari item `dari` ke item lain. */
export function salinDeductible(objek: ObjekFire[], dari: PosisiItem, metode: 'ALL' | 'BI'): ObjekFire[] {
  const sumber = objek[dari.o]?.items[dari.i]?.coverages ?? []
  const deductibleUntuk = (kode: string) =>
    salinan(sumber.filter((c) => c.coverage === kode).flatMap((c) => c.deductibles ?? []))
  return objek.map((x) => ({
    ...x,
    items: x.items.map((it) => {
      const bi = it.itemTypeId === JENIS_ITEM_BI
      if ((metode === 'ALL' && bi) || (metode === 'BI' && !bi) || !it.coverages) return it
      return { ...it, coverages: it.coverages.map((c) => ({ ...c, deductibles: deductibleUntuk(c.coverage) })) }
    }),
  }))
}
