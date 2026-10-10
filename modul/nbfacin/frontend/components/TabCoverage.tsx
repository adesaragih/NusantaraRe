// Tab Coverage FIRE layar Inward Facultative (tiket 43, tahap C1).
//
// Port berantai: `CoverageList.xml` (grid objek `.LocationList`: No. · Object Name · Location; baris dibuka) ->
// `PropertyItemListCoverage.xml` (grid item: Object Item Type · Currency · TSI Object Item · Total Gross Premium ·
// ‰ Total Net Rate; grid total per mata uang) -> `InputCoverageFire.xml` (grid `.CoverageList`: Coverage · ‰ Standard
// Rate · Premi; Tambah = `AddCoverageAutoFire` - lima coverage FIRE bila masih kosong; Hapus) -> `FormCoverage`.
// Save = Dtl sel `Save` (`SaveFacIn_Act`): seluruh objek dikirim `PUT …/objek` (sama dengan tab Object).
//
// Total Gross Premium item = Σ Premium coverage (`CountPremi_ACT` langkah total) - dijumlah eksak (`jumlahDesimal`).
// Total per mata uang dan ‰ Total Net Rate = dari server (pembagian di backend, tanpa float).
//
// Tahap C2 net rate (tiket 44): ‰ Total Net Rate dapat diisi bila item memuat kelima OLDID `OLDID_NET_RATE`
// (`CekNetRate_ACT`); perubahan memanggil `POST …/hitung-net-rate` (`CalculateNetRate_ACT`) - NetRate dan premi tiap
// coverage dihitung backend (jeda 500 ms). Item tanpa kelima coverage: Total Net Rate tampil baca-saja.
//
// Keputusan agent (tiket 43): P-5 perubahan di tab Coverage dan tab Object disimpan masing-masing lewat tombol Save
// tab itu; pindah tab tanpa Save membuang perubahan (sama dengan pola tab sebelumnya). P-6 tombol Copy Coverage /
// Copy Deductible / Copy Coverage From dan grid Summary = tahap berikut.

import { useEffect, useRef, useState } from 'react'

import { Gagal, Memuat } from '../../../../inti/frontend/components/ui/dasar'
import { desimalSah, jumlahDesimal } from '../../../../inti/frontend/lib/desimal'
import { angka4 } from '../angka4'
import {
  ambilObjek,
  coverageOtomatis,
  hitungCoverage,
  hitungNetRate,
  simpanAkumulasi,
  simpanObjek,
  type CoverageObjek,
  type ObjekFire,
} from '../api'
import {
  GRID_COV_ITEM,
  GRID_COV_OBJEK,
  GRID_COV_TOTAL,
  GRID_COVERAGE,
  GRID_OBJEK,
  OLDID_NET_RATE,
  RINGKASAN_COV_ITEM,
  RINGKASAN_COV_MATA_UANG,
  SIMPAN_COVERAGE,
  TEKS_FORM_OPPORTUNITY,
  TEKS_INWARD,
  TEKS_OBJEK,
  SALIN_COVERAGE,
} from '../labels'
import FormCoverage, { adaGalatCoverage, angkaSah, coverageBaru } from './FormCoverage'
import { adaGalatDeductible } from './GridDeductible'
import { dariKetikan } from './IsianAngka4'
import PopupSalinCoverage from './PopupSalinCoverage'
import {
  JENIS_ITEM_BI,
  salinCoverageDari,
  salinCoverageLokasi,
  salinCoverageSemua,
  ringkasSalinDeductible,
  salinDeductible,
  type PosisiItem,
} from './salinCoverage'
import { rapikanObjek } from './TabObject'

/** Total Gross Premium item = Σ Premium coverage (eksak); tanpa coverage -> nilai server. */
export function totalPremiItem(coverages: CoverageObjek[] | undefined, dariServer: string | undefined): string {
  if (!coverages || coverages.length === 0) return dariServer ?? ''
  return jumlahDesimal(coverages.map((c) => c.premium)).total
}

/** Satu baris ringkasan per Object Item Type + Currency (`TempTotalItem.pxResults`). */
export interface RingkasanItemCoverage {
  itemType: string
  currency: string
  tsi: string
  premium: string
}

/** Satu baris ringkasan per Currency (`TempTotal.pxResults`). */
export interface RingkasanMataUangCoverage {
  currency: string
  tsi: string
  premium: string
  rate: string
}

/** Teks desimal -> bilangan bulat + jumlah digit pecahan. */
function keBulat(t: string): [bigint, number] {
  const negatif = t.startsWith('-')
  const [bulat = '0', pecahan = ''] = (negatif || t.startsWith('+') ? t.slice(1) : t).split('.')
  const n = BigInt((bulat === '' ? '0' : bulat) + pecahan)
  return [negatif ? -n : n, pecahan.length]
}

/** a / b dibulatkan setengah-ke-atas menjauhi nol (apd RoundHalfUp backend). */
function bagiBulat(a: bigint, b: bigint): bigint {
  const negatif = a < 0n !== b < 0n
  const pa = a < 0n ? -a : a
  const pb = b < 0n ? -b : b
  const q = (2n * pa + pb) / (2n * pb)
  return negatif ? -q : q
}

/**
 * Rate ringkasan per Currency = Premium / TSI, TSI 0 -> 0 - mengikuti layar Pega DEV (work owner 05-10-2026 "ikutin
 * sesuai pega saja": TSI 50.000.000.000, premi 9.391.000.000 -> 0,1878). Korpus `SumTSIPremiSpreadedRNM_FIRE_Act`
 * menulis `@divide(Premium, TSI, 20) * 1000` (= 187,82) - beda versi server, sejenis OQ-011; yang dipakai layar DEV.
 * Eksak tanpa float: bagi 20 desimal setengah-ke-atas, simpan 8 desimal - sama dengan `TotalPerMataUang` backend.
 */
export function rateRingkasan(premium: string, tsi: string): string {
  if (!desimalSah(premium.trim()) || !desimalSah(tsi.trim())) return ''
  const [p, sp] = keBulat(premium.trim())
  const [t, st] = keBulat(tsi.trim())
  if (t === 0n) return '0'
  const q20 = bagiBulat(p * 10n ** BigInt(st + 20), t * 10n ** BigInt(sp))
  const r8 = bagiBulat(q20, 10n ** 12n)
  const negatif = r8 < 0n
  const digit = (negatif ? -r8 : r8).toString().padStart(9, '0')
  return (negatif ? '-' : '') + digit.slice(0, -8) + '.' + digit.slice(-8)
}

/**
 * Ringkasan seluruh lokasi (`SumTSIPremiSpreadedRNM_FIRE_Act`): per item, TSI = `.TSIObjectItem` dan Premium =
 * `.TotalGrossPremi` (= Σ Premium coverage); dikelompokkan per (ItemType, Currency) dan per Currency, urut kemunculan.
 */
export function ringkasanCoverage(objek: ObjekFire[]): { item: RingkasanItemCoverage[]; mataUang: RingkasanMataUangCoverage[] } {
  const item: { itemType: string; currency: string; tsi: string[]; premium: string[] }[] = []
  const uang: { currency: string; tsi: string[]; premium: string[] }[] = []
  for (const x of objek) {
    for (const it of x.items) {
      const premi = totalPremiItem(it.coverages, it.totalGrossPremi)
      let a = item.find((r) => r.itemType === it.itemType && r.currency === it.currency)
      if (!a) item.push((a = { itemType: it.itemType, currency: it.currency, tsi: [], premium: [] }))
      a.tsi.push(it.tsi)
      a.premium.push(premi)
      let b = uang.find((r) => r.currency === it.currency)
      if (!b) uang.push((b = { currency: it.currency, tsi: [], premium: [] }))
      b.tsi.push(it.tsi)
      b.premium.push(premi)
    }
  }
  return {
    item: item.map((r) => ({ itemType: r.itemType, currency: r.currency, tsi: jumlahDesimal(r.tsi).total, premium: jumlahDesimal(r.premium).total })),
    mataUang: uang.map((r) => {
      const tsi = jumlahDesimal(r.tsi).total
      const premium = jumlahDesimal(r.premium).total
      return { currency: r.currency, tsi, premium, rate: rateRingkasan(premium, tsi) }
    }),
  }
}

/** Tombol buka/tutup baris. */
function Buka({ buka, onKlik }: { buka: boolean; onKlik: () => void }) {
  return (
    <button
      type="button"
      className={'btn btn--ghost btn--sm nbf-buka' + (buka ? ' nbf-buka--terbuka' : '')}
      aria-label={TEKS_OBJEK.bukaBaris}
      aria-expanded={buka}
      onClick={onKlik}
    >
      ▸
    </button>
  )
}

const balik = (daftar: string[], k: string) => (daftar.includes(k) ? daftar.filter((x) => x !== k) : [...daftar, k])

/** `CekNetRate_ACT`: FlagNetRate = kelima OLDID net rate ada di daftar coverage item (cocok harfiah). */
export function bolehNetRate(coverages: CoverageObjek[] | undefined): boolean {
  const ada = new Set((coverages ?? []).map((c) => c.oldId))
  return OLDID_NET_RATE.every((k) => ada.has(k))
}

/** Isian ‰ Total Net Rate (teks desimal); hitung ulang lewat backend sesudah jeda. */
function IsianNetRate({ nilai, onUbah }: { nilai: string; onUbah: (v: string) => void }) {
  return (
    <input
      className="field__input nbf-angka"
      type="text"
      inputMode="decimal"
      aria-label={GRID_COV_ITEM[4].label}
      value={nilai}
      onChange={(e) => onUbah(dariKetikan(e.target.value))}
    />
  )
}

/**
 * `hanyaAkumulasi` = tahap Underwriting (tiket 60, `KasusNB.ubahAkumulasi`): layar dibungkus baca-saja dengan pengecualian
 * `data-ubah-akumulasi` (Choose / Copy Accumulation + popup + Save); Save menulis Accumulation saja (`simpanAkumulasi`).
 */
export default function TabCoverage({ caseId, hanyaAkumulasi = false }: { caseId: string; hanyaAkumulasi?: boolean }) {
  const [objek, setObjek] = useState<ObjekFire[] | null>(null)
  const [terbuka, setTerbuka] = useState<string[]>([])
  const [galat, setGalat] = useState<unknown>(null)
  const [menyimpan, setMenyimpan] = useState(false)
  const [tersimpan, setTersimpan] = useState(false)
  // Tiket 53: salin coverage sedang menghitung ulang premi; item tujuan popup Copy Coverage From.
  const [menyalin, setMenyalin] = useState(false)
  const [salinKe, setSalinKe] = useState<PosisiItem | null>(null)
  // Umpan balik Copy Deductible di bawah tombolnya (item sumber) - sebelumnya klik tanpa tanda apa pun.
  const [pesanDeductible, setPesanDeductible] = useState<(PosisiItem & { teks: string; ok: boolean }) | null>(null)
  const aktif = useRef(true)
  const jadwalNet = useRef<number | undefined>(undefined)
  const nomorNet = useRef(0)

  useEffect(() => {
    aktif.current = true
    ambilObjek(caseId).then(
      (h) => {
        if (aktif.current) setObjek(h.baris.map(rapikanObjek))
      },
      (err: unknown) => {
        if (aktif.current) setGalat(err)
      },
    )
    return () => {
      aktif.current = false
    }
  }, [caseId])

  function ubahCoverage(o: number, i: number, coverages: CoverageObjek[]) {
    setObjek((d) =>
      d === null
        ? d
        : d.map((x, a) => (a !== o ? x : { ...x, items: x.items.map((it, b) => (b !== i ? it : { ...it, coverages })) })),
    )
    setTersimpan(false)
  }

  /**
   * Copy Accumulation (tiket 46) - `CopyAccumulationCode_Act`: Accumulation Code / Description coverage sumber disalin ke
   * SETIAP coverage di SETIAP item objek (lokasi) yang sama (langkah 2 = PropertyItemList, 2.1 = CoverageList).
   */
  function salinAkumulasi(o: number, sumber: CoverageObjek) {
    const { accumulationCode, accumulationDescription } = sumber
    setObjek((d) =>
      d === null
        ? d
        : d.map((x, a) =>
            a !== o
              ? x
              : {
                  ...x,
                  items: x.items.map((it) =>
                    it.coverages ? { ...it, coverages: it.coverages.map((cv) => ({ ...cv, accumulationCode, accumulationDescription })) } : it,
                  ),
                },
          ),
    )
    setTersimpan(false)
  }

  /** ‰ Total Net Rate diubah: simpan isian, lalu bagi ke coverage di backend (`CalculateNetRate_ACT`). */
  function ubahNetRate(o: number, i: number, v: string) {
    if (objek === null) return
    const it = objek[o]!.items[i]!
    setObjek((d) =>
      d === null ? d : d.map((x, a) => (a !== o ? x : { ...x, items: x.items.map((y, b) => (b !== i ? y : { ...y, totalNetRate: v })) })),
    )
    setTersimpan(false)
    window.clearTimeout(jadwalNet.current)
    if (v.trim() !== '' && !angkaSah(v)) return
    jadwalNet.current = window.setTimeout(() => {
      const n = ++nomorNet.current
      hitungNetRate(caseId, {
        tsi: it.tsi,
        totalNetRate: v,
        coverages: it.coverages ?? [],
        isAdjustable: it.isAdjustable,
        pctAdjustOther: it.pctAdjustOther,
      }).then(
        (h) => {
          if (n === nomorNet.current) ubahCoverage(o, i, h.coverages)
        },
        (err: unknown) => {
          if (n === nomorNet.current) setGalat(err)
        },
      )
    }, 500)
  }

  /**
   * Tiket 53: pasang hasil salin, lalu hitung ulang premi setiap coverage item tujuan (CountPremi_ACT percent / percent,
   * TSI item tujuan - K53-1). Gagal satu -> galat tampil, data salinan tetap di layar (belum disimpan).
   */
  /** Copy Deductible To All / For BI: salin (CopyDeductibleToAllLocation_Act) lalu tampilkan apa yang tersalin. */
  function salinDeductibleDari(dari: PosisiItem, metode: 'ALL' | 'BI') {
    if (objek === null) return
    const r = ringkasSalinDeductible(objek, dari, metode)
    void pasangSalinan(salinDeductible(objek, dari, metode), [])
    const utama =
      r.sumber === 0
        ? SALIN_COVERAGE.deductibleKosong
        : r.coverage === 0
          ? r.biDilewati > 0
            ? ''
            : SALIN_COVERAGE.deductibleTanpaPadanan
          : SALIN_COVERAGE.deductibleTersalin(r.coverage, r.item)
    const teks = [utama, r.biDilewati > 0 ? SALIN_COVERAGE.deductibleBiDilewati(r.biDilewati) : ''].filter((x) => x !== '').join(' ')
    setPesanDeductible({ ...dari, teks, ok: r.sumber > 0 && r.coverage > 0 && r.biDilewati === 0 })
  }

  async function pasangSalinan(data: ObjekFire[], tujuan: PosisiItem[]) {
    setPesanDeductible(null)
    setObjek(data)
    setTersimpan(false)
    if (tujuan.length === 0) return
    setMenyalin(true)
    setGalat(null)
    try {
      const hasil = await Promise.all(
        tujuan.map(async ({ o, i }) => {
          const it = data[o]!.items[i]!
          const coverages = await Promise.all(
            (it.coverages ?? []).map((c) =>
              hitungCoverage(caseId, {
                coverage: c,
                tsi: it.tsi,
                mode: 'percent',
                modeDiskon: 'percent',
                isAdjustable: it.isAdjustable,
                pctAdjustOther: it.pctAdjustOther,
              }),
            ),
          )
          return { o, i, coverages }
        }),
      )
      setObjek((d) =>
        d === null
          ? d
          : d.map((x, a) => ({
              ...x,
              items: x.items.map((y, b) => {
                const h = hasil.find((z) => z.o === a && z.i === b)
                return h ? { ...y, coverages: h.coverages } : y
              }),
            })),
      )
    } catch (err) {
      setGalat(err)
    } finally {
      setMenyalin(false)
    }
  }

  async function tambah(o: number, i: number, ada: CoverageObjek[]) {
    if (ada.length > 0) {
      ubahCoverage(o, i, [...ada, coverageBaru()])
      setTerbuka((t) => [...t, `c${o}-${i}-${ada.length}`])
      return
    }
    try {
      const h = await coverageOtomatis()
      ubahCoverage(o, i, h.baris.map((b) => coverageBaru(b)))
    } catch (err) {
      setGalat(err)
    }
  }

  async function simpan() {
    if (
      objek === null ||
      objek.some((x) => x.items.some((it) => adaGalatCoverage(it.coverages ?? []) || (it.coverages ?? []).some((c) => adaGalatDeductible(c.deductibles))))
    )
      return
    setMenyimpan(true)
    setGalat(null)
    setTersimpan(false)
    setPesanDeductible(null)
    try {
      const h = await (hanyaAkumulasi ? simpanAkumulasi : simpanObjek)(caseId, objek)
      setObjek(h.baris.map(rapikanObjek))
      setTersimpan(true)
    } catch (err) {
      setGalat(err)
    } finally {
      setMenyimpan(false)
    }
  }

  if (objek === null) return galat ? <Gagal galat={galat} /> : <Memuat />
  const ringkasan = ringkasanCoverage(objek)

  return (
    <div className="nbf-objek nbf-cov-tab">
      <Gagal galat={galat} />
      {tersimpan && <div className="alert alert--ok">{TEKS_INWARD.tersimpan}</div>}
      <div className="nbf-cov-wrap">
        <table className="nbf-tabel">
          <thead>
            <tr>
              <th scope="col" />
              {GRID_COV_OBJEK.map((k) => (
                <th key={k.sel} scope="col">
                  {k.label}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {objek.length === 0 && (
              <tr>
                <td colSpan={4}>{TEKS_INWARD.kosong}</td>
              </tr>
            )}
            {objek.map((x, o) => [
              <tr key={`o${o}`} className={terbuka.includes(`o${o}`) ? 'nbf-baris--terbuka' : undefined}>
                <td>
                  <Buka buka={terbuka.includes(`o${o}`)} onKlik={() => setTerbuka((t) => balik(t, `o${o}`))} />
                </td>
                <td>{x.objectNo}</td>
                <td>{x.objectName}</td>
                <td>{x.riskLocation}</td>
              </tr>,
              terbuka.includes(`o${o}`) && (
                <tr key={`od${o}`} className="nbf-objek__detail">
                  <td colSpan={4}>
                    <div className="nbf-objek__isi nbf-lapis nbf-lapis--objek">
                      <div className="nbf-cov-wrap">
                        <table className="nbf-tabel">
                          <thead>
                            <tr>
                              <th scope="col" />
                              {GRID_COV_ITEM.map((k) => (
                                <th key={k.sel} scope="col">
                                  {k.label}
                                </th>
                              ))}
                            </tr>
                          </thead>
                          <tbody>
                            {x.items.length === 0 && (
                              <tr>
                                <td colSpan={6}>{TEKS_INWARD.kosong}</td>
                              </tr>
                            )}
                            {x.items.map((it, i) => {
                              const covs = it.coverages ?? []
                              return [
                                <tr key={`i${i}`} className={terbuka.includes(`i${o}-${i}`) ? 'nbf-baris--terbuka' : undefined}>
                                  <td>
                                    <Buka buka={terbuka.includes(`i${o}-${i}`)} onKlik={() => setTerbuka((t) => balik(t, `i${o}-${i}`))} />
                                  </td>
                                  <td>{it.itemType}</td>
                                  <td>{it.currency}</td>
                                  <td className="nbf-angka">{angka4(it.tsi)}</td>
                                  <td className="nbf-angka">{angka4(totalPremiItem(covs, it.totalGrossPremi))}</td>
                                  <td className="nbf-angka">
                                    {bolehNetRate(covs) ? (
                                      <IsianNetRate nilai={it.totalNetRate ?? ''} onUbah={(v) => ubahNetRate(o, i, v)} />
                                    ) : (
                                      angka4(it.totalNetRate ?? '')
                                    )}
                                  </td>
                                </tr>,
                                terbuka.includes(`i${o}-${i}`) && (
                                  <tr key={`id${i}`} className="nbf-objek__detail">
                                    <td colSpan={6}>
                                      <div className="nbf-cov-wrap nbf-lapis nbf-lapis--item">
                                        <table className="nbf-tabel">
                                          <thead>
                                            <tr>
                                              <th scope="col" />
                                              {GRID_COVERAGE.map((k) => (
                                                <th key={k.sel} scope="col">
                                                  {k.label}
                                                </th>
                                              ))}
                                              <th scope="col" className="table__actions">
                                                <button type="button" className="btn nbf-tema__tombol btn--ghost btn--sm" onClick={() => void tambah(o, i, covs)}>
                                                  {GRID_OBJEK.tambah}
                                                </button>
                                              </th>
                                            </tr>
                                          </thead>
                                          <tbody>
                                            {covs.length === 0 && (
                                              <tr>
                                                <td colSpan={5}>{TEKS_INWARD.kosong}</td>
                                              </tr>
                                            )}
                                            {covs.map((c, n) => [
                                              <tr key={`c${n}`} className={terbuka.includes(`c${o}-${i}-${n}`) ? 'nbf-baris--terbuka' : undefined}>
                                                <td>
                                                  <Buka
                                                    buka={terbuka.includes(`c${o}-${i}-${n}`)}
                                                    onKlik={() => setTerbuka((t) => balik(t, `c${o}-${i}-${n}`))}
                                                  />
                                                </td>
                                                <td>
                                                  <span className="nbf-cov-kode">{c.oldId}</span>
                                                </td>
                                                <td className="nbf-angka">{angka4(c.rateOjk)}</td>
                                                <td className="nbf-angka">{angka4(c.premium)}</td>
                                                <td className="table__actions">
                                                  <button
                                                    type="button"
                                                    className="btn nbf-tema__tombol btn--ghost btn--sm"
                                                    onClick={() => ubahCoverage(o, i, covs.filter((_, k) => k !== n))}
                                                  >
                                                    {GRID_OBJEK.hapus}
                                                  </button>
                                                </td>
                                              </tr>,
                                              terbuka.includes(`c${o}-${i}-${n}`) && (
                                                <tr key={`cd${n}`} className="nbf-objek__detail">
                                                  <td colSpan={5}>
                                                    <FormCoverage
                                                      caseId={caseId}
                                                      c={c}
                                                      tsiItem={it.tsi}
                                                      item={{ isAdjustable: it.isAdjustable, pctAdjustOther: it.pctAdjustOther, currency: it.currency }}
                                                      pertama={n === 0}
                                                      zipRisiko={x.zipCode}
                                                      lokasi={{ province: x.province, occupationIds: x.occupations.map((k) => k.occupationId) }}
                                                      onSalinAkumulasi={() => salinAkumulasi(o, c)}
                                                      ubah={(baru) => ubahCoverage(o, i, covs.map((y, k) => (k === n ? baru : y)))}
                                                    />
                                                  </td>
                                                </tr>
                                              ),
                                            ])}
                                          </tbody>
                                        </table>
                                      </div>
                                      {/* Tiket 53: tombol salin (InputCoverageFire); visibilitas = pyCondition Pega. */}
                                      <div className="nbf-cov-salin">
                                        {i === 0 && (
                                          <button
                                            type="button"
                                            className="btn nbf-tema__tombol btn--sm"
                                            disabled={menyalin}
                                            onClick={() => {
                                              const h = salinCoverageLokasi(objek, o)
                                              void pasangSalinan(h.objek, h.tujuan)
                                            }}
                                          >
                                            {SALIN_COVERAGE.coverage}
                                          </button>
                                        )}
                                        {o === 0 && i === 0 && (
                                          <button
                                            type="button"
                                            className="btn nbf-tema__tombol btn--sm"
                                            disabled={menyalin}
                                            onClick={() => {
                                              const h = salinCoverageSemua(objek)
                                              void pasangSalinan(h.objek, h.tujuan)
                                            }}
                                          >
                                            {SALIN_COVERAGE.semua}
                                          </button>
                                        )}
                                        <button type="button" className="btn nbf-tema__tombol btn--sm" disabled={menyalin} onClick={() => setSalinKe({ o, i })}>
                                          {SALIN_COVERAGE.dari}
                                        </button>
                                        {o === 0 && i === 0 && (
                                          <button
                                            type="button"
                                            className="btn nbf-tema__tombol btn--sm"
                                            disabled={menyalin}
                                            onClick={() => salinDeductibleDari({ o, i }, 'ALL')}
                                          >
                                            {SALIN_COVERAGE.deductibleSemua}
                                          </button>
                                        )}
                                        {o === 0 && it.itemTypeId === JENIS_ITEM_BI && (
                                          <button
                                            type="button"
                                            className="btn nbf-tema__tombol btn--sm"
                                            disabled={menyalin}
                                            onClick={() => salinDeductibleDari({ o, i }, 'BI')}
                                          >
                                            {SALIN_COVERAGE.deductibleBI}
                                          </button>
                                        )}
                                        {menyalin && <span className="muted">{SALIN_COVERAGE.menyalin}</span>}
                                      </div>
                                      {pesanDeductible?.o === o && pesanDeductible.i === i && (
                                        <div className={`alert ${pesanDeductible.ok ? 'alert--ok' : 'alert--warn'}`} role="status">
                                          {pesanDeductible.teks}
                                        </div>
                                      )}
                                    </td>
                                  </tr>
                                ),
                              ]
                            })}
                          </tbody>
                        </table>
                      </div>
                      <div className="nbf-cov-wrap">
                        <table className="nbf-tabel">
                          <thead>
                            <tr>
                              {GRID_COV_TOTAL.map((k) => (
                                <th key={k.sel} scope="col">
                                  {k.label}
                                </th>
                              ))}
                            </tr>
                          </thead>
                          <tbody>
                            {(x.totalPerCurrency ?? []).length === 0 && (
                              <tr>
                                <td colSpan={4}>{TEKS_INWARD.kosong}</td>
                              </tr>
                            )}
                            {(x.totalPerCurrency ?? []).map((t) => (
                              <tr key={t.currency}>
                                <td>{t.currency}</td>
                                <td className="nbf-angka">{angka4(t.tsi)}</td>
                                <td className="nbf-angka">{angka4(t.premium)}</td>
                                <td className="nbf-angka">{angka4(t.rate)}</td>
                              </tr>
                            ))}
                          </tbody>
                        </table>
                      </div>
                    </div>
                  </td>
                </tr>
              ),
            ])}
          </tbody>
        </table>
      </div>
      {/* Ringkasan seluruh coverage (SummaryCoverage_Section) - hijau seperti Pega, ikut berubah saat premi diubah. */}
      <section className="nbf-cov-ringkasan">
        <div className="nbf-cov-wrap">
          <table className="nbf-tabel">
            <thead>
              <tr>
                {RINGKASAN_COV_ITEM.map((k) => (
                  <th key={k} scope="col">
                    {k}
                  </th>
                ))}
              </tr>
            </thead>
            <tbody>
              {ringkasan.item.length === 0 && (
                <tr>
                  <td colSpan={4}>{TEKS_INWARD.kosong}</td>
                </tr>
              )}
              {ringkasan.item.map((r) => (
                <tr key={`${r.itemType}|${r.currency}`}>
                  <td>{r.itemType}</td>
                  <td>{r.currency}</td>
                  <td className="nbf-angka">{angka4(r.tsi)}</td>
                  <td className="nbf-angka">{angka4(r.premium)}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
        <div className="nbf-cov-wrap">
          <table className="nbf-tabel">
            <thead>
              <tr>
                {RINGKASAN_COV_MATA_UANG.map((k) => (
                  <th key={k} scope="col">
                    {k}
                  </th>
                ))}
              </tr>
            </thead>
            <tbody>
              {ringkasan.mataUang.length === 0 && (
                <tr>
                  <td colSpan={4}>{TEKS_INWARD.kosong}</td>
                </tr>
              )}
              {ringkasan.mataUang.map((r) => (
                <tr key={r.currency}>
                  <td>{r.currency}</td>
                  <td className="nbf-angka">{angka4(r.tsi)}</td>
                  <td className="nbf-angka">{angka4(r.premium)}</td>
                  <td className="nbf-angka">{angka4(r.rate)}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </section>
      {salinKe !== null && (
        <PopupSalinCoverage
          objek={objek}
          onTutup={() => setSalinKe(null)}
          onSalin={(dari) => {
            const h = salinCoverageDari(objek, dari, salinKe)
            setSalinKe(null)
            void pasangSalinan(h.objek, h.tujuan)
          }}
        />
      )}
      <div className="nbf-objek__kaki">
        <button
          type="button"
          className="btn nbf-tema__tombol nbf-tema__tombol--utama btn--primary"
          data-ubah-akumulasi={hanyaAkumulasi ? '' : undefined}
          onClick={() => void simpan()}
          disabled={menyimpan}
        >
          {menyimpan ? TEKS_FORM_OPPORTUNITY.menyimpan : SIMPAN_COVERAGE}
        </button>
      </div>
    </div>
  )
}
