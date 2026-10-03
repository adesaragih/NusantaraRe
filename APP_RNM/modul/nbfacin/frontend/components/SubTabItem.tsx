// Sub-tab Object Item baris objek FIRE (tiket 39).
//
// Port `NB FacIn\Section\PropertyItemList.xml` (layout NB `!IsEDM`): grid `.Property.PropertyItemList` (kolom Object
// Item Type · Condition · Year · Unit(s) · Currency · TSI Object Item; Add di kepala, Delete per baris; baris dibuka
// = form `PropertyItemFacIn_Section`) dan grid Total `.Property.TotalTSIList` (Currency · Total TSI per mata uang).
//
// - Object Item Type (RD `BrowseV_JN_OBJ_ITEM`, nilai MJOI_KODE): memilih mengisi `.ItemType` dan Object Item Note
//   (`SetObjItemType_Act` -> `GetObjectItem`: JN_OBJ_ITEM, KETERANGAN).
// - Currency (RD `BrowseCurrency_RD`, tanpa ITL). TSI Object Item = UANG: teks desimal, dijumlah per mata uang
//   dengan `jumlahDesimal` (BigInt, ADR-0003) - sama dengan `AddPaymentCurency_ACT` (Σ TSIObjectItem per Currency,
//   baris total hanya bila > 0; tanpa konversi kurs, tanpa pembulatan). Tampil 4 desimal (grid Pega).
// - Adjustable (`ResetPct_Adjustment`): dicentang -> Adjustment Pct. angka (`PctAdjustOther`, sah 60..100,
//   `ValidateAdjustPct`); tidak -> dropdown `PctAdjust2` dan PctAdjustOther kembali 100.
//
// Keputusan agent (tiket 39): K-1 hitung ulang premi / spreading saat TSI berubah (`CountPremi_ACT`,
// `cekSpreadingFactIn`) = tab Coverage, bukan di sini. K-2 urutan baris Total = urutan muncul mata uang di daftar
// item (Pega: urutan tabel CURRENCY). K-3 pesan Unit / TSI minus = teks sistem baru (field value Pega tidak ada di
// korpus). K-4 daftar Condition = `DDL\Condition.xml`; Adjustment Pct. (PctAdjust2) = `DDL\PctAdjust2.xml` ("100").
// K-7 PctAdjust2 item baru = "100" (satu-satunya pilihan; data contoh selalu 100).
// A133 (backend c3): Currency wajib - bertanda wajib, menahan Save; pesannya tampil sesudah Save dicoba.

import { useEffect, useState } from 'react'

import { Area, Field, Pilih, type Opsi } from '../../../../inti/frontend/components/ui/dasar'
import { desimalSah, jumlahDesimal } from '../../../../inti/frontend/lib/desimal'
import { formatNumber } from '../../../../inti/frontend/lib/format'
import { daftarJenisItem, daftarMataUang, type ItemObjek, type JenisItem } from '../api'
import IsianUang from './IsianUang'
import {
  CONDITION_KOSONG,
  FORM_ITEM as F,
  GRID_ITEM,
  GRID_OBJEK,
  ITEM_KOSONG,
  OPSI_CONDITION,
  OPSI_PCT_ADJUST,
  PCT_ADJUST_AWAL,
  TEKS_INWARD,
  TEKS_ITEM,
  TEKS_OBJEK,
} from '../labels'

/** Desimal tampilan grid (pxNumber 4 dp). */
const DESIMAL_GRID = 4

/** Item kosong (Add). `PctAdjustOther` awal 0 (nilai awal sel 36); `PctAdjust2` awal "100" (K-7). */
export function itemBaru(): ItemObjek {
  return {
    itemTypeId: '', itemType: '', note: '', propertyYear: '', unit: '', condition: '', currency: '', tsi: '',
    yearOfPlanting: '', noOfTree: '', areaHectar: '', remark: '', isAdjustable: false, pctAdjust2: PCT_ADJUST_AWAL, pctAdjustOther: '0',
  }
}

/** Memilih Object Item Type (`SetObjItemType_Act`). */
export function pilihJenis(i: ItemObjek, j: JenisItem | undefined, kode: string): ItemObjek {
  return { ...i, itemTypeId: kode, itemType: j?.nama ?? '', note: j?.keterangan ?? '' }
}

/** Centang Adjustable (`ResetPct_Adjustment`: tidak dicentang -> PctAdjustOther = 100). */
export function setelAdjustable(i: ItemObjek, v: boolean): ItemObjek {
  return v ? { ...i, isAdjustable: true } : { ...i, isAdjustable: false, pctAdjustOther: '100' }
}

/** Tanda desimal eksak: -1 negatif, 0 nol, 1 positif (lewat `jumlahDesimal`, tanpa float). */
export function tandaDesimal(teks: string): -1 | 0 | 1 {
  const t = jumlahDesimal([teks]).total
  if (/^-?0*(\.0*)?$/.test(t)) return 0
  return t.startsWith('-') ? -1 : 1
}

/** a - b sebagai teks desimal eksak (a dan b desimal tanpa tanda). */
const selisih = (a: string, b: string) => jumlahDesimal([a, `-${b}`]).total

/** Galat per item: Currency kosong (A133), Unit <= 0, TSI bukan desimal / minus, % adjustment di luar 60..100. */
export function galatItem(i: ItemObjek): { currency?: string; unit?: string; tsi?: string; pct?: string } {
  const g: { currency?: string; unit?: string; tsi?: string; pct?: string } = {}
  if (i.currency.trim() === '') g.currency = TEKS_ITEM.currencyWajib
  const unit = i.unit.trim()
  if (unit !== '' && !(desimalSah(unit) && tandaDesimal(unit) > 0)) g.unit = TEKS_ITEM.unit
  const tsi = i.tsi.trim()
  if (tsi !== '') {
    if (!desimalSah(tsi)) g.tsi = TEKS_ITEM.tsiBukanAngka
    else if (tandaDesimal(tsi) < 0) g.tsi = TEKS_ITEM.tsiMinus
  }
  if (i.isAdjustable) {
    const p = i.pctAdjustOther.trim()
    const sah = /^\d+(\.\d+)?$/.test(p) && tandaDesimal(selisih(p, '60')) >= 0 && tandaDesimal(selisih('100', p)) >= 0
    if (!sah) g.pct = TEKS_ITEM.pctAdjust
  }
  return g
}

/** Ada item bergalat. */
export function adaGalatItem(items: ItemObjek[]): boolean {
  return items.some((i) => Object.keys(galatItem(i)).length > 0)
}

/** Total TSI per mata uang (`AddPaymentCurency_ACT`): hanya mata uang dengan total > 0. */
export function totalPerMataUang(items: ItemObjek[]): { currency: string; total: string }[] {
  const urutan: string[] = []
  for (const i of items) if (i.currency !== '' && !urutan.includes(i.currency)) urutan.push(i.currency)
  return urutan
    .map((c) => ({ currency: c, total: jumlahDesimal(items.filter((i) => i.currency === c).map((i) => i.tsi)).total }))
    .filter((t) => tandaDesimal(t.total) > 0)
}

function FormItem({
  i,
  ubah,
  jenis,
  mataUang,
  tandaiWajib,
}: {
  i: ItemObjek
  ubah: (i: ItemObjek) => void
  jenis: JenisItem[]
  mataUang: string[]
  /** Save sudah dicoba - pesan wajib ditampilkan. */
  tandaiWajib: boolean
}) {
  const set = (k: keyof ItemObjek) => (v: string) => ubah({ ...i, [k]: v })
  const g = galatItem(i)
  const opsiJenis: Opsi[] = jenis.map((j) => ({ value: j.kode, label: j.nama }))
  const opsiUang: Opsi[] = mataUang.map((m) => ({ value: m, label: m }))
  return (
    <div className="nbf-objek__isi">
      <div className="nbf-opp__kolom">
        <div className="nbf-opp__tumpuk">
          <Pilih
            label={F.itemType.label}
            value={i.itemTypeId}
            onChange={(v) => ubah(pilihJenis(i, jenis.find((j) => j.kode === v), v))}
            opsi={opsiJenis}
            kosong={ITEM_KOSONG}
          />
          <Area label={F.note.label} value={i.note} onChange={set('note')} baris={2} />
          <Field label={F.year.label} value={i.propertyYear} onChange={set('propertyYear')} />
          <Field label={F.unit.label} value={i.unit} onChange={set('unit')} error={g.unit} />
          <Pilih label={F.condition.label} value={i.condition} onChange={set('condition')} opsi={OPSI_CONDITION} kosong={CONDITION_KOSONG} />
          <Pilih
            label={F.currency.label}
            value={i.currency}
            onChange={set('currency')}
            opsi={opsiUang}
            kosong={ITEM_KOSONG}
            required
            error={tandaiWajib ? g.currency : undefined}
          />
          <IsianUang label={F.tsi.label} value={i.tsi} onChange={set('tsi')} error={g.tsi} />
        </div>
        <div className="nbf-opp__tumpuk">
          <Field label={F.yearOfPlanting.label} value={i.yearOfPlanting} onChange={set('yearOfPlanting')} />
          <Field label={F.noOfTree.label} value={i.noOfTree} onChange={set('noOfTree')} />
          <Field label={F.areaHectar.label} value={i.areaHectar} onChange={set('areaHectar')} />
          <Area label={`${F.remark.label} ${i.itemType}`.trim()} value={i.remark} onChange={set('remark')} baris={2} />
          <label className="nbf-inward__pilihan">
            <input type="checkbox" checked={i.isAdjustable} onChange={(e) => ubah(setelAdjustable(i, e.target.checked))} /> {F.adjustable.label}
          </label>
          {i.isAdjustable ? (
            <Field label={F.pctAdjust.label} value={i.pctAdjustOther} onChange={set('pctAdjustOther')} error={g.pct} />
          ) : (
            <Pilih label={F.pctAdjust.label} value={i.pctAdjust2} onChange={set('pctAdjust2')} opsi={OPSI_PCT_ADJUST} />
          )}
        </div>
      </div>
    </div>
  )
}

export default function SubTabItem({
  items,
  ubah,
  tandaiWajib = false,
}: {
  items: ItemObjek[]
  ubah: (items: ItemObjek[]) => void
  /** Save tab Object sudah dicoba. */
  tandaiWajib?: boolean
}) {
  const [terbuka, setTerbuka] = useState<number[]>([])
  const [jenis, setJenis] = useState<JenisItem[]>([])
  const [mataUang, setMataUang] = useState<string[]>([])

  useEffect(() => {
    let batal = false
    daftarJenisItem().then(
      (h) => {
        if (!batal) setJenis(h.baris)
      },
      () => {},
    )
    daftarMataUang().then(
      (h) => {
        if (!batal) setMataUang(h.baris)
      },
      () => {},
    )
    return () => {
      batal = true
    }
  }, [])

  const total = totalPerMataUang(items)

  return (
    <div className="nbf-objek__isi">
      <div className="table-wrap">
        <table className="nbf-tabel">
          <thead>
            <tr>
              <th scope="col" />
              {GRID_ITEM.kolom.map((k) => (
                <th key={k.sel} scope="col">
                  {k.label}
                </th>
              ))}
              <th scope="col" className="table__actions">
                <button
                  type="button"
                  className="btn btn--ghost btn--sm"
                  onClick={() => {
                    ubah([...items, itemBaru()])
                    setTerbuka((t) => [...t, items.length])
                  }}
                >
                  {GRID_OBJEK.tambah}
                </button>
              </th>
            </tr>
          </thead>
          <tbody>
            {items.length === 0 && (
              <tr>
                <td colSpan={8}>{TEKS_INWARD.kosong}</td>
              </tr>
            )}
            {items.map((i, n) => [
              <tr key={`b-${n}`}>
                <td>
                  <button
                    type="button"
                    className="btn btn--ghost btn--sm"
                    aria-label={TEKS_OBJEK.bukaBaris}
                    aria-expanded={terbuka.includes(n)}
                    onClick={() => setTerbuka((t) => (t.includes(n) ? t.filter((x) => x !== n) : [...t, n]))}
                  >
                    {terbuka.includes(n) ? '▾' : '▸'}
                  </button>
                </td>
                <td>{i.itemType}</td>
                <td>{OPSI_CONDITION.find((o) => o.value === i.condition)?.label ?? i.condition}</td>
                <td>{i.propertyYear}</td>
                <td>{i.unit}</td>
                <td>{i.currency}</td>
                <td className="nbf-angka">{formatNumber(i.tsi, DESIMAL_GRID)}</td>
                <td className="table__actions">
                  <button
                    type="button"
                    className="btn btn--ghost btn--sm"
                    onClick={() => {
                      ubah(items.filter((_, x) => x !== n))
                      setTerbuka((t) => t.filter((x) => x !== n).map((x) => (x > n ? x - 1 : x)))
                    }}
                  >
                    {GRID_OBJEK.hapus}
                  </button>
                </td>
              </tr>,
              terbuka.includes(n) && (
                <tr key={`d-${n}`} className="nbf-objek__detail">
                  <td colSpan={8}>
                    <FormItem
                      i={i}
                      ubah={(baru) => ubah(items.map((x, k) => (k === n ? baru : x)))}
                      jenis={jenis}
                      mataUang={mataUang}
                      tandaiWajib={tandaiWajib}
                    />
                  </td>
                </tr>
              ),
            ])}
          </tbody>
        </table>
      </div>

      <div className="table-wrap">
        <table className="nbf-tabel">
          <thead>
            <tr>
              {GRID_ITEM.total.map((k) => (
                <th key={k.sel} scope="col">
                  {k.label}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {total.length === 0 && (
              <tr>
                <td colSpan={2}>{TEKS_INWARD.kosong}</td>
              </tr>
            )}
            {total.map((t) => (
              <tr key={t.currency}>
                <td>{t.currency}</td>
                <td className="nbf-angka">{formatNumber(t.total, DESIMAL_GRID)}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  )
}
