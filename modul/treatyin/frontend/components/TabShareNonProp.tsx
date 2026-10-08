// Tab **Share** cabang NON-PROPORSIONAL — bentuk dan tombol dari ekspor
// (`labelsShareNP.ts`): panel Share (% RNM Share, % Brokerage, Share Across
// The Board, Share to Other Retro, Brokerage From Other Retro, Update
// Summary) · grid Reinsurer Name · grid Facultative Reinsurers · sub-tab
// RNM Share (grid per layer + panel rincian `Share`) · Summarry of RNM
// Share · Total All Layers RNM Share (sembilan grid) · Update Total.
//
// ⭐ Rumusnya hidup di services (`hitung_share_np.go`), disalin dari
// Activity dan diukur terhadap data Pega. Layar ini hanya mengirim isian
// dan memasang isi tab yang dikembalikan — nol rumus di sini.
//
// ⭐ Grid reasuradur: sel ber-`pyReadOnlyCondition TreatyIn.ViewState = 1`
// (bisa diisi di mode Edit — tangkapan layar Pega memperlihatkan isian
// aktif, dan 35 baris tersimpan terisi). Add = `TreatyInNonAddItem`
// `sharereins` / `sharefacname`, yang hanya menambah baris ber-ID kosong —
// tanpa rumus, jadi dikerjakan di layar.
//
// ⭐ Kunci sel dibaca dengan `hanya_baca()` pembaca bersama
// (`D:\XML_NURE\_migration-docs\alat-baca-ekspor`): `pyReadOnlyCondition`
// MENIMPA `pyReadOnly`. Di panel rincian, Layer Type · Layer · Layer Part
// Type · Layer Part · Cover · Deduction Details · Spreading Type terkunci
// HANYA bila `TreatyIn.ViewState = 1`; Class of Business dan sel grid
// spreading bernama ber-`Auto`. Semuanya dapat diisi di mode Edit.
//
// ⚠️ Simpangan tampilan yang disengaja:
//  - Dropdown `Spreading Type` panel rincian di Pega hanya tampil bila
//    nilainya SUDAH terisi; baris yang baru disusun Update Summary tidak
//    pernah dapat memilihnya. Di sini ia tampil pula di mode Edit saat
//    kosong (lihat S1 di `hitung_share_np.go`).
//  - Grid `TreatyIn.Share` Pega membuka panel rincian lewat klik baris
//    (`masterDetail` · `expandPane`); di sini lewat tombol `Detail` per
//    baris — klik baris tidak terbaca pembaca layar.
//
// ⛔ Hasil suntingan dan hitungan hidup di salinan layar ini; jalur Save
// menunggu keputusan pemilik proses.

import { useEffect, useRef, useState, type ReactNode } from 'react'

import { Field, Kosong, Panel, Pilih } from '../../../../inti/frontend/components/ui/dasar'
import {
  ambilIndukSpreading,
  ambilOpsiLimits,
  ambilReasuradurShare,
  hitungShareNP,
  type AksiShareNP,
  type BarisDeduksiShare,
  type BarisReinsShare,
  type BarisShareNP,
  type NilaiShare,
  type OpsiLimits,
  type OpsiPilihan,
  type PilihanWarisan,
  type ShareNP,
  type SimpulLimit,
  type SusunanSpreading,
} from '../api'
import { SUB_TAB_SHARE } from '../labels'
import { kontrakRevisi } from '../labelsLimitsNP'
import {
  GRID_RINCIAN_SPREADING,
  GRID_TOTAL_SHARE,
  KEPALA_SHARE,
  KOLOM_DEDUKSI_SHARE,
  KOLOM_FAC_REINS,
  KOLOM_REINS,
  KOLOM_RINGKASAN_SHARE,
  KOLOM_SPREADING,
  KOLOM_SPREADING_MANUAL,
  SHARE_NP,
} from '../labelsShareNP'
import type { ModeForm } from '../mode'
import IsianAuto from './IsianAuto'
import { StripTabNavigasi, TombolNavigasi } from './navigasi'
import { PemicuUbah, usePemicuUbah } from './pemicuUbah'
import { Bagian, IkonChevronKanan, KepalaBagian, TombolHapus, TombolTambah } from './limitsUI'
import { formatLimit } from './TabLimitsProp'

const uang = (v: string) => formatLimit('uang', 2, v)
const persen = (v: string) => formatLimit('persen', 2, v)
const lebihDariNol = (v: string) => Number.parseFloat(v.replace(',', '.')) > 0

function ganti<T>(larik: readonly T[], i: number, baru: T): T[] {
  return larik.map((x, j) => (j === i ? baru : x))
}

/** Aksi yang menyentuh SATU baris Share — pesannya hanya mengganti pesan baris itu. */
const AKSI_BARIS: readonly AksiShareNP[] = [
  'spreading-type',
  'rnm-baris',
  'deduksi',
  'spreading-tambah',
  'spreading-hapus',
  'spreading-pct',
]

/** Isi tab kosong — kontrak tanpa data Share. */
export const SHARE_NP_KOSONG: ShareNP = {
  RNMShare: '',
  BrokeragePercent: '',
  RNMShareAcrossTheBoard: 'true',
  FacultativeShare: '',
  FacultativeShareBrokerage: '',
  RnmShareDeducted: '',
  IsProRate: '',
  ShareReins: [],
  ShareFacultativeReinsurers: [],
  Share: [],
  FacultativeShareList: [],
  LimitShareSummaryList: [],
  LimitFacShareSummaryList: [],
  Total: {},
}

/** Medan teks/angka; `onLepas` = peristiwa `change` Pega (sesudah keluar medan). */
function Medan({
  label,
  nilai,
  bisaUbah,
  placeholder,
  onUbah,
  onLepas,
}: {
  label: string
  nilai: string
  bisaUbah: boolean
  placeholder?: string
  onUbah: (v: string) => void
  onLepas?: () => void
}) {
  // Peristiwa `change` Pega — hanya bila nilainya berubah (`pemicuUbah.tsx`).
  const pemicu = usePemicuUbah(nilai, bisaUbah ? onLepas : undefined)
  return (
    <div className="trin__limit-medan" onFocus={pemicu.masuk} onBlur={pemicu.keluar}>
      <Field label={label} value={nilai} readOnly={!bisaUbah} placeholder={bisaUbah ? placeholder : undefined} onChange={onUbah} />
    </div>
  )
}

/** Dropdown `associated` — baca-saja di luar mode Edit (label tampil). */
function PilihMedan({
  label,
  nilai,
  opsi,
  bisaUbah,
  onUbah,
}: {
  label: string
  nilai: string
  opsi: readonly OpsiPilihan[]
  bisaUbah: boolean
  onUbah: (v: string) => void
}) {
  if (!bisaUbah) {
    const t = opsi.find((o) => o.value === nilai)?.label ?? nilai
    return <Field label={label} value={t} readOnly onChange={() => undefined} />
  }
  return <Pilih label={label} value={nilai} kosong={SHARE_NP.pilihKosong} opsi={[...opsi]} onChange={onUbah} />
}

const OPSI_KOSONG: OpsiLimits = { jenisTreaty: [], kelompokTreaty: [], mataUang: [] }

/** Grid reasuradur — `ShareReins` atau `ShareFacultativeReinsurers`. */
function GridReins({
  kolom,
  baris,
  pilihan,
  bisaUbah,
  onUbah,
}: {
  kolom: readonly string[]
  baris: readonly BarisReinsShare[]
  pilihan: readonly PilihanWarisan[]
  bisaUbah: boolean
  onUbah: (b: BarisReinsShare[]) => void
}) {
  return (
    <div className="tl-share-reins">
      <KepalaBagian
        judul={kolom[0] ?? ''}
        jumlah={baris.length}
        aksi={
          bisaUbah && (
            <TombolTambah
              label={SHARE_NP.tambah}
              onClick={() => {
                // `TreatyInNonAddItem(sharereins | sharefacname)`: baris ber-ID kosong.
                onUbah([...baris, { ID: '', ReinsID: '', ReinsName: '', Layer: '', SharePct: '' }])
              }}
            />
          )
        }
      />
      <div className="table-wrap trin__limit-grid">
        <table className="trin__tabel">
          <thead>
            <tr>
              {kolom.map((k) => (
                <th key={k} scope="col">
                  {k}
                </th>
              ))}
              {bisaUbah && <th scope="col" className="tl-kolom-aksi" aria-label={SHARE_NP.hapus} />}
            </tr>
          </thead>
          <tbody>
            {baris.length === 0 && (
              <tr>
                <td colSpan={kolom.length + (bisaUbah ? 1 : 0)}>{SHARE_NP.tanpaBaris}</td>
              </tr>
            )}
            {baris.map((b, i) => (
              <tr key={i}>
                <td>
                  {/* `.ReinsName` pxAutoComplete `BrowseAgentNusaRe_RD`:
                      `.ClientName`, `.ID → .ReinsID`. */}
                  <IsianAuto
                    label=""
                    nilai={b.ReinsName}
                    pilihan={pilihan}
                    bisaUbah={bisaUbah}
                    onPilih={(nama, id) => {
                      onUbah(ganti(baris, i, { ...b, ReinsName: nama, ReinsID: id }))
                    }}
                  />
                </td>
                <td>
                  <Medan
                    label=""
                    nilai={b.Layer}
                    bisaUbah={bisaUbah}
                    placeholder={SHARE_NP.placeholderTeks}
                    onUbah={(v) => onUbah(ganti(baris, i, { ...b, Layer: v }))}
                  />
                </td>
                <td>
                  <Medan
                    label=""
                    nilai={bisaUbah ? b.SharePct : persen(b.SharePct)}
                    bisaUbah={bisaUbah}
                    placeholder={SHARE_NP.placeholderAngka}
                    onUbah={(v) => onUbah(ganti(baris, i, { ...b, SharePct: v }))}
                  />
                </td>
                {bisaUbah && (
                  <td>
                    <TombolHapus
                      label={SHARE_NP.hapus}
                      labelAkses={`${SHARE_NP.hapus} ${kolom[0] ?? ''} ${i + 1}`}
                      onClick={() => onUbah(baris.filter((_, j) => j !== i))}
                    />
                  </td>
                )}
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  )
}

/**
 * Grid baca-saja `Currency · Value`. Grid total berkepala `judul · Value`;
 * grid rincian panel (`satuKepala`) hanya berjudul, seperti di ekspor.
 */
function GridTotal({ judul, baris, satuKepala = false }: { judul: string; baris: readonly NilaiShare[]; satuKepala?: boolean }) {
  return (
    <div className="table-wrap tl-share-total">
      <table className="trin__tabel">
        <thead>
          <tr>
            {satuKepala ? (
              <th scope="colgroup" colSpan={2}>
                {judul}
              </th>
            ) : (
              <>
                <th scope="col">{judul}</th>
                <th scope="col">{SHARE_NP.nilai}</th>
              </>
            )}
          </tr>
        </thead>
        <tbody>
          {baris.length === 0 && (
            <tr>
              <td colSpan={2}>{SHARE_NP.tanpaBaris}</td>
            </tr>
          )}
          {baris.map((b, i) => (
            <tr key={i}>
              <td>{b.Currency}</td>
              <td className="tl-angka">{uang(b.Value)}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}

/** Sel pasangan mata uang · nilai ke-`i` sebuah larik. */
function selPasangan(daftar: readonly NilaiShare[], i: number): [string, string] {
  const v = daftar[i]
  return v === undefined ? ['', ''] : [v.Currency, uang(v.Value)]
}

/** Judul baris: `Layer 1 Part of Layer 1`. */
export function judulBarisShare(b: BarisShareNP): string {
  return [b.LayerType, b.Layer, SHARE_NP.bagianOf, b.LayerPartType, b.LayerPart].filter((x) => x !== '').join(' ')
}

/**
 * Panel rincian satu baris — flow action `Share` (`Section/Share.xml`).
 */
export function RincianShare({
  b,
  induk,
  indukManual = [],
  opsi = OPSI_KOSONG,
  pesan = [],
  bisaUbah,
  modeUbah,
  onUbah,
  hitung,
}: {
  b: BarisShareNP
  /** Isi dropdown `Spreading Type` — RD `ParentReinsMasterTrt` grup baris ini. */
  induk: readonly SusunanSpreading[]
  /**
   * Isi dropdown `Reins Type` spreading manual — RD yang sama dengan
   * `TreatyGroupID = TempSprd.TreatyGroupID`, yang tak pernah diisi untuk
   * baris Non-Prop: filter grup DILEWATI, induk SEMUA Treaty Group.
   */
  indukManual?: readonly SusunanSpreading[]
  /** Jenis layer · Cover · Treaty Group — opsi tab Limits. */
  opsi?: OpsiLimits
  /** Pesan Activity yang menempel pada baris ini (`pesanBaris`). */
  pesan?: readonly string[]
  /** Mode Edit DAN bukan kunci materialitas. */
  bisaUbah: boolean
  modeUbah: boolean
  onUbah: (b: BarisShareNP) => void
  /** Jalankan aksi atas baris ini; `dasar` = baris terbaru bila baru diubah. */
  hitung: (aksi: AksiShareNP, opsi?: { baris?: number; sts?: string; dasar?: BarisShareNP }) => void
}) {
  const bernama = b.SpreadingTypeXOL !== ''
  const ubahDeduksi = (i: number, d: BarisDeduksiShare) => {
    onUbah({ ...b, DeductionList: ganti(b.DeductionList, i, d) })
  }
  return (
    <div className="tl-rincian tl-share-rincian">
      {pesan.length > 0 && (
        <ul className="tl-pesan" role="alert">
          {pesan.map((p, i) => (
            <li key={i}>{p}</li>
          ))}
        </ul>
      )}
      <div className="form-grid">
        <Medan
          label={SHARE_NP.persenRnm}
          nilai={b.RNMShare}
          bisaUbah={bisaUbah}
          placeholder={SHARE_NP.placeholderAngka}
          onUbah={(v) => onUbah({ ...b, RNMShare: v })}
          // `TreatyInXOLAddSpreadingDetail(idx)`.
          onLepas={() => hitung('rnm-baris')}
        />
        {/* `.Cover` pxDropdown — terkunci hanya bila `ViewState = 1`. */}
        <PilihMedan
          label={SHARE_NP.cover}
          nilai={b.Cover}
          opsi={opsi.cover ?? []}
          bisaUbah={modeUbah}
          onUbah={(v) => onUbah({ ...b, Cover: v })}
        />
      </div>
      {/* Layer Type · Layer · "Part of" · Layer Part Type · Layer Part —
          terkunci hanya bila `ViewState = 1`, nol aksi. */}
      <div className="tl-share-layer">
        <PilihMedan
          label={SHARE_NP.jenisLayer}
          nilai={b.LayerType}
          opsi={opsi.jenisLayer ?? []}
          bisaUbah={modeUbah}
          onUbah={(v) => onUbah({ ...b, LayerType: v })}
        />
        <Medan label={SHARE_NP.layer} nilai={b.Layer} bisaUbah={modeUbah} onUbah={(v) => onUbah({ ...b, Layer: v })} />
        <span className="tl-share-layer__of">{SHARE_NP.bagianOf}</span>
        <PilihMedan
          label={SHARE_NP.jenisLayer}
          nilai={b.LayerPartType}
          opsi={opsi.jenisLayer ?? []}
          bisaUbah={modeUbah}
          onUbah={(v) => onUbah({ ...b, LayerPartType: v })}
        />
        <Medan label={SHARE_NP.layer} nilai={b.LayerPart} bisaUbah={modeUbah} onUbah={(v) => onUbah({ ...b, LayerPart: v })} />
      </div>

      <Bagian judul={SHARE_NP.kelasBisnis}>
        <div className="table-wrap trin__limit-grid">
          <table className="trin__tabel">
            <thead>
              <tr>
                <th scope="col">{SHARE_NP.treatyGroup}</th>
              </tr>
            </thead>
            <tbody>
              {b.TreatyGroupList.length === 0 && (
                <tr>
                  <td>{SHARE_NP.tanpaBaris}</td>
                </tr>
              )}
              {b.TreatyGroupList.map((g, i) => (
                <tr key={i}>
                  <td>
                    {/* `.TreatyGroup` pxAutoComplete `BrowseTreatyGroup_RD`
                        (`.ID → .TreatyGroupID`), `Auto`. Aksi ekspornya
                        (`SetIndexLayer_DT`, `TotalEgnpi`) milik tab Limits —
                        nol pengaruh ke hitungan Share; yang berubah di sini
                        isi dropdown Spreading Type (grup baris pertama). */}
                    <IsianAuto
                      label=""
                      nilai={g.TreatyGroup}
                      pilihan={opsi.kelompokTreaty}
                      bisaUbah={modeUbah}
                      onPilih={(nama, id) => {
                        onUbah({ ...b, TreatyGroupList: ganti(b.TreatyGroupList, i, { TreatyGroup: nama, TreatyGroupID: id }) })
                      }}
                    />
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </Bagian>

      <Bagian
        judul={SHARE_NP.deduksi}
        aksi={
          modeUbah && (
            <TombolTambah
              label={SHARE_NP.tambah}
              disabled={!bisaUbah}
              onClick={() => {
                onUbah({
                  ...b,
                  DeductionList: [...b.DeductionList, { Comment: '', Currency: '', CurrencyID: '', Deduction: '', DeductionPct: '' }],
                })
              }}
            />
          )
        }
      >
        <div className="table-wrap trin__limit-grid">
          <table className="trin__tabel">
            <thead>
              <tr>
                {KOLOM_DEDUKSI_SHARE.map((k) => (
                  <th key={k} scope="col">
                    {k}
                  </th>
                ))}
                {modeUbah && <th scope="col" className="tl-kolom-aksi" aria-label={SHARE_NP.hapusBaris} />}
              </tr>
            </thead>
            <tbody>
              {b.DeductionList.length === 0 && (
                <tr>
                  <td colSpan={KOLOM_DEDUKSI_SHARE.length + (modeUbah ? 1 : 0)}>{SHARE_NP.tanpaBaris}</td>
                </tr>
              )}
              {b.DeductionList.map((d, i) => (
                <tr key={i}>
                  <td>
                    <Medan label="" nilai={d.Comment} bisaUbah={bisaUbah} onUbah={(v) => ubahDeduksi(i, { ...d, Comment: v })} />
                  </td>
                  <td>
                    {/* `.Currency` pxAutoComplete `BrowseCurrency_RD` — dipanggil
                        tanpa parameter, jadi filter `Currency`/`ID` dilewati:
                        isinya = daftar mata uang tab Limits (`!= ITL`).
                        `change` → `CalculateDeduction(index)`. */}
                    <PemicuUbah
                      className="trin__limit-medan"
                      nilai={d.Currency}
                      aktif={bisaUbah}
                      aksi={() => hitung('deduksi', { baris: i, sts: '' })}
                    >
                      <IsianAuto
                        label=""
                        nilai={d.Currency}
                        pilihan={opsi.mataUang}
                        bisaUbah={bisaUbah}
                        onPilih={(nama, id) => ubahDeduksi(i, { ...d, Currency: nama, CurrencyID: id })}
                      />
                    </PemicuUbah>
                  </td>
                  <td>
                    <Medan
                      label=""
                      nilai={bisaUbah ? d.Deduction : uang(d.Deduction)}
                      bisaUbah={bisaUbah}
                      onUbah={(v) => ubahDeduksi(i, { ...d, Deduction: v })}
                      onLepas={() => hitung('deduksi', { baris: i, sts: 'val' })}
                    />
                  </td>
                  <td className="tl-share-atau">or</td>
                  <td>
                    <Medan
                      label=""
                      nilai={bisaUbah ? d.DeductionPct : persen(d.DeductionPct)}
                      bisaUbah={bisaUbah}
                      onUbah={(v) => ubahDeduksi(i, { ...d, DeductionPct: v })}
                      onLepas={() => hitung('deduksi', { baris: i, sts: 'pct' })}
                    />
                  </td>
                  <td>
                    <input
                      type="checkbox"
                      aria-label={KOLOM_DEDUKSI_SHARE[5]}
                      checked={d.DeductionPctCalculate === 'true'}
                      disabled={!bisaUbah}
                      onChange={(e) => ubahDeduksi(i, { ...d, DeductionPctCalculate: e.target.checked ? 'true' : 'false' })}
                    />
                  </td>
                  {modeUbah && (
                    <td>
                      <TombolHapus
                        label={SHARE_NP.hapusBaris}
                        labelAkses={`${SHARE_NP.hapusBaris} ${SHARE_NP.deduksi} ${i + 1}`}
                        disabled={!bisaUbah}
                        onClick={() => {
                          // deleteRow, lalu `CalculateDeduction(index)`.
                          const baru = { ...b, DeductionList: b.DeductionList.filter((_, j) => j !== i) }
                          onUbah(baru)
                          hitung('deduksi', { baris: i, sts: '', dasar: baru })
                        }}
                      />
                    </td>
                  )}
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </Bagian>

      <Bagian judul={SHARE_NP.spreading}>
        {(bernama || modeUbah) && (
          <div className="trin__limit-medan">
            {modeUbah ? (
              <label className="field">
                <span className="field__label">{SHARE_NP.spreadingType}</span>
                <select
                  className="field__input"
                  value={b.SpreadingTypeXOL}
                  onChange={(e) => {
                    // `FetchQSfromMasterXOL(ParentReinsTypeID = .SpreadingTypeXOL, IsUpdate = 0)`.
                    const baru = { ...b, SpreadingTypeXOL: e.target.value }
                    onUbah(baru)
                    hitung('spreading-type', { dasar: baru })
                  }}
                >
                  <option value="">{SHARE_NP.pilihKosong}</option>
                  {b.SpreadingTypeXOL !== '' && !induk.some((p) => p.reinsTypeName === b.SpreadingTypeXOL) && (
                    <option value={b.SpreadingTypeXOL}>{b.SpreadingTypeXOL}</option>
                  )}
                  {induk.map((p) => (
                    <option key={p.reinsTypeId} value={p.reinsTypeName}>
                      {p.reinsTypeName}
                    </option>
                  ))}
                </select>
              </label>
            ) : (
              <Field label={SHARE_NP.spreadingType} value={b.SpreadingTypeXOL} readOnly onChange={() => undefined} />
            )}
          </div>
        )}
        {bernama ? (
          <>
            <div className="table-wrap trin__limit-grid">
              <table className="trin__tabel">
                <thead>
                  <tr>
                    {KOLOM_SPREADING.map((k) => (
                      <th key={k} scope="col">
                        {k}
                      </th>
                    ))}
                  </tr>
                </thead>
                <tbody>
                  {b.SpreadingListXOL.length === 0 && (
                    <tr>
                      <td colSpan={2}>{SHARE_NP.tanpaBaris}</td>
                    </tr>
                  )}
                  {b.SpreadingListXOL.map((s, i) => (
                    <tr key={i}>
                      {/* `.ReinsTypeName` · `.Pct` pxTextInput `Auto`, nol aksi. */}
                      <td>
                        <Medan
                          label=""
                          nilai={s.ReinsTypeName}
                          bisaUbah={modeUbah}
                          onUbah={(v) => onUbah({ ...b, SpreadingListXOL: ganti(b.SpreadingListXOL, i, { ...s, ReinsTypeName: v }) })}
                        />
                      </td>
                      <td className="tl-angka">
                        <Medan
                          label=""
                          nilai={modeUbah ? s.Pct : persen(s.Pct)}
                          bisaUbah={modeUbah}
                          onUbah={(v) => onUbah({ ...b, SpreadingListXOL: ganti(b.SpreadingListXOL, i, { ...s, Pct: v }) })}
                        />
                      </td>
                    </tr>
                  ))}
                </tbody>
                <tfoot>
                  <tr>
                    <th scope="row">{SHARE_NP.totalPct}</th>
                    <td className="tl-angka">{persen(b.SpreadingTotalPctXOL)}</td>
                  </tr>
                </tfoot>
              </table>
            </div>
            <p className="tl-share-catatan">
              {SHARE_NP.totalSpreadingPct} <strong>{persen(b.SpreadingTotalPctXOL)}</strong>
            </p>
            {/* Blok `hidden, reference` — TAMPIL bersama Spreading bernama
                (`NOHEADER`, nol syarat sendiri): dasar · OR · R/I per baris. */}
            <div className="tl-share-total-grid tl-share-rujukan">
              {GRID_RINCIAN_SPREADING.flat().map((g) => (
                <GridTotal key={g.kunci} judul={g.judul} baris={b[g.kunci as keyof BarisShareNP] as readonly NilaiShare[]} satuKepala />
              ))}
            </div>
          </>
        ) : (
          <>
            <div className="table-wrap trin__limit-grid">
              <table className="trin__tabel">
                <thead>
                  <tr>
                    {KOLOM_SPREADING_MANUAL.map((k) => (
                      <th key={k} scope="col">
                        {k}
                      </th>
                    ))}
                    {modeUbah && (
                      <th scope="col" className="tl-kolom-aksi">
                        <TombolTambah label={SHARE_NP.tambah} onClick={() => hitung('spreading-tambah')} />
                      </th>
                    )}
                  </tr>
                </thead>
                <tbody>
                  {b.SpreadingListXOL.length === 0 && (
                    <tr>
                      <td colSpan={2 + (modeUbah ? 1 : 0)}>{SHARE_NP.tanpaBaris}</td>
                    </tr>
                  )}
                  {b.SpreadingListXOL.map((s, i) => (
                    <tr key={i}>
                      <td>
                        {modeUbah ? (
                          <select
                            className="field__input"
                            aria-label={KOLOM_SPREADING_MANUAL[0]}
                            value={s.ReinsTypeID}
                            onChange={(e) => {
                              // postValue saja — namanya diisi `SetSpreadingXOL`.
                              onUbah({ ...b, SpreadingListXOL: ganti(b.SpreadingListXOL, i, { ...s, ReinsTypeID: e.target.value }) })
                            }}
                          >
                            <option value="">{SHARE_NP.pilihKosong}</option>
                            {indukManual.map((p) => (
                              <option key={p.reinsTypeId} value={p.reinsTypeId}>
                                {p.reinsTypeName}
                              </option>
                            ))}
                          </select>
                        ) : (
                          s.ReinsTypeName
                        )}
                      </td>
                      <td>
                        <Medan
                          label=""
                          nilai={modeUbah ? s.Pct : persen(s.Pct)}
                          bisaUbah={modeUbah}
                          onUbah={(v) => onUbah({ ...b, SpreadingListXOL: ganti(b.SpreadingListXOL, i, { ...s, Pct: v }) })}
                          onLepas={() => hitung('spreading-pct')}
                        />
                      </td>
                      {modeUbah && (
                        <td>
                          <TombolHapus
                            label={SHARE_NP.hapus}
                            labelAkses={`${SHARE_NP.hapus} ${SHARE_NP.spreading} ${i + 1}`}
                            onClick={() => hitung('spreading-hapus', { baris: i })}
                          />
                        </td>
                      )}
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
            <p className="tl-share-catatan">
              {SHARE_NP.totalSharePct} <strong>{persen(b.SpreadingTotalPctXOL)}</strong>
            </p>
          </>
        )}
      </Bagian>
    </div>
  )
}

export default function TabShareNonProp({
  share: awal,
  layers,
  mode,
  idKontrak = '',
  commencement = '',
  edmState = '',
  edmJenisMaterial = '',
  reasuradur: reasuradurAwal,
  induk: indukAwal,
  opsi: opsiAwal,
  onUbah,
}: {
  share?: ShareNP
  /** Layer Limits Non-Prop — sumber Update Summary. */
  layers: readonly SimpulLimit[]
  mode: ModeForm
  idKontrak?: string
  /** `TreatyIn.Commencement`, YYYYMMDD. */
  commencement?: string
  edmState?: string
  /** `EDMMaterialType` — `2` mematikan isian dan Update Summary/Total. */
  edmJenisMaterial?: string
  /** Isi autocomplete; bila tidak diberikan, diminta sendiri. */
  reasuradur?: readonly PilihanWarisan[]
  /** Isi dropdown Spreading Type per Treaty Group; bila tidak diberikan, diminta. */
  induk?: Readonly<Record<string, readonly SusunanSpreading[]>>
  /** Opsi dropdown panel rincian; bila tidak diberikan, diminta sendiri. */
  opsi?: OpsiLimits
  /** Isi tab terkini, dilaporkan tiap kali berubah — form menyimpannya lintas pindah tab. */
  onUbah?: (s: ShareNP) => void
}) {
  const [opsi, setOpsi] = useState<OpsiLimits>(opsiAwal ?? OPSI_KOSONG)
  const [share, setShare] = useState<ShareNP>(awal ?? SHARE_NP_KOSONG)
  const onUbahTerkini = useRef(onUbah)
  onUbahTerkini.current = onUbah
  useEffect(() => {
    onUbahTerkini.current?.(share)
  }, [share])
  const [pesan, setPesan] = useState<string[]>([])
  // Pesan per baris Share, menurut indeks baris.
  const [pesanBaris, setPesanBaris] = useState<Readonly<Record<number, readonly string[]>>>({})
  const [gagal, setGagal] = useState('')
  const [buka, setBuka] = useState<number | null>(null)
  const [sub, setSub] = useState<string>(SUB_TAB_SHARE[0])
  const [reasuradur, setReasuradur] = useState<readonly PilihanWarisan[]>(reasuradurAwal ?? [])
  const [induk, setInduk] = useState<Readonly<Record<string, readonly SusunanSpreading[]>>>(indukAwal ?? {})
  // `TreatyIn.ViewState != 1` — tombol dan isian hidup di mode Edit saja …
  const modeUbah = mode === 'ubah'
  // … dan `pyDisabledWhen … || TreatyIn.EDMMaterialType = 2`.
  const bisaUbah = modeUbah && edmJenisMaterial.trim() !== '2'
  const adaFac = lebihDariNol(share.FacultativeShare)

  useEffect(() => {
    if (reasuradurAwal !== undefined || !modeUbah) return
    let dibuang = false
    ambilReasuradurShare()
      .then((d) => {
        if (!dibuang) setReasuradur(d)
      })
      .catch(() => undefined)
    return () => {
      dibuang = true
    }
  }, [reasuradurAwal, modeUbah])

  useEffect(() => {
    if (opsiAwal !== undefined || !modeUbah) return
    let dibuang = false
    ambilOpsiLimits()
      .then((o) => {
        if (!dibuang) setOpsi(o)
      })
      .catch(() => undefined)
    return () => {
      dibuang = true
    }
  }, [opsiAwal, modeUbah])

  // Kedua dropdown spreading — diminta sekali per Treaty Group saat panel
  // dibuka: grup baris (Spreading Type) dan grup KOSONG (Reins Type manual).
  const grupBuka = buka === null ? null : (share.Share[buka]?.TreatyGroupList[0]?.TreatyGroupID ?? '')
  useEffect(() => {
    if (indukAwal !== undefined || !modeUbah || grupBuka === null) return
    let dibuang = false
    for (const g of new Set([grupBuka, ''])) {
      if (induk[g] !== undefined) continue
      ambilIndukSpreading(g, commencement)
        .then((d) => {
          if (!dibuang) setInduk((k) => ({ ...k, [g]: d }))
        })
        .catch(() => undefined)
    }
    return () => {
      dibuang = true
    }
  }, [indukAwal, modeUbah, grupBuka, commencement, induk])

  /** Menjalankan rantai Activity satu aksi atas isi `dasar`. */
  const hitung = (aksi: AksiShareNP, opsi: { indeks?: number; baris?: number; sts?: string; dasar?: ShareNP } = {}) => {
    if (!modeUbah) return
    setGagal('')
    hitungShareNP({
      aksi,
      share: opsi.dasar ?? share,
      layers,
      indeks: opsi.indeks ?? 0,
      baris: opsi.baris ?? 0,
      sts: opsi.sts ?? '',
      idKontrak,
      commencement,
    })
      .then((h) => {
        setShare(h.share)
        setPesan(h.pesan)
        const baru: Record<number, string[]> = {}
        for (const p of h.pesanBaris) {
          ;(baru[p.indeks] ??= []).push(p.pesan)
        }
        // Aksi satu baris hanya mengganti pesan baris itu; aksi tab mengganti seluruhnya.
        setPesanBaris((kini) =>
          AKSI_BARIS.includes(aksi) && opsi.indeks !== undefined ? { ...kini, [opsi.indeks]: baru[opsi.indeks] ?? [] } : baru,
        )
      })
      .catch((e: unknown) => {
        setGagal(e instanceof Error ? e.message : String(e))
      })
  }
  const ubahAkar = (medan: keyof ShareNP, v: string) => {
    const baru = { ...share, [medan]: v }
    setShare(baru)
    return baru
  }

  const tombolRingkasan: ReactNode = modeUbah && (
    <button
      type="button"
      className="btn btn--primary btn--sm"
      disabled={!bisaUbah}
      onClick={() => {
        hitung('summary')
      }}
    >
      {SHARE_NP.perbaruiRingkasan}
    </button>
  )

  return (
    <Panel judul={SHARE_NP.judul}>
      <div className="tl-rincian">
        <KepalaBagian judul={SHARE_NP.judul} aksi={tombolRingkasan} />
        {share.IsProRate === 'true' && <p className="tl-share-catatan">{SHARE_NP.nonProRate}</p>}
        <div className="tl-share-akar">
          <div className="tl-share-kolom">
            <Medan
              label={SHARE_NP.persenRnm}
              nilai={bisaUbah ? share.RNMShare : persen(share.RNMShare)}
              bisaUbah={bisaUbah}
              placeholder={SHARE_NP.placeholderAngka}
              onUbah={(v) => ubahAkar('RNMShare', v)}
              onLepas={() => hitung('rnm')}
            />
            <Medan
              label={SHARE_NP.persenBrokerage}
              nilai={bisaUbah ? share.BrokeragePercent : persen(share.BrokeragePercent)}
              bisaUbah={bisaUbah}
              placeholder={SHARE_NP.placeholderAngka}
              onUbah={(v) => ubahAkar('BrokeragePercent', v)}
              onLepas={() => hitung('brokerage')}
            />
            <label className="trin__limit-medan tl-share-centang">
              <input
                type="checkbox"
                checked={share.RNMShareAcrossTheBoard === 'true'}
                disabled={!bisaUbah}
                onChange={(e) => {
                  const baru = ubahAkar('RNMShareAcrossTheBoard', e.target.checked ? 'true' : 'false')
                  hitung('centang', { dasar: baru })
                }}
              />{' '}
              {SHARE_NP.acrossTheBoard}
            </label>
          </div>
          <div className="tl-share-kolom">
            <Medan
              label={SHARE_NP.shareKeRetro}
              nilai={bisaUbah ? share.FacultativeShare : persen(share.FacultativeShare)}
              bisaUbah={bisaUbah}
              placeholder={SHARE_NP.placeholderAngka}
              onUbah={(v) => ubahAkar('FacultativeShare', v)}
              onLepas={() => hitung('fac')}
            />
            {adaFac && (
              <Medan
                label={SHARE_NP.brokerageRetro}
                nilai={bisaUbah ? share.FacultativeShareBrokerage : persen(share.FacultativeShareBrokerage)}
                bisaUbah={bisaUbah}
                placeholder={SHARE_NP.placeholderAngka}
                onUbah={(v) => ubahAkar('FacultativeShareBrokerage', v)}
                onLepas={() => hitung('fac')}
              />
            )}
          </div>
        </div>

        <div className="tl-share-dua">
          <GridReins
            kolom={KOLOM_REINS}
            baris={share.ShareReins}
            pilihan={reasuradur}
            bisaUbah={modeUbah}
            onUbah={(b) => setShare({ ...share, ShareReins: b })}
          />
          {/* `pyContainerVisibleWhen TreatyIn.FacultativeShare > 0`. */}
          {adaFac && (
            <GridReins
              kolom={KOLOM_FAC_REINS}
              baris={share.ShareFacultativeReinsurers}
              pilihan={reasuradur}
              bisaUbah={modeUbah}
              onUbah={(b) => setShare({ ...share, ShareFacultativeReinsurers: b })}
            />
          )}
        </div>

        {pesan.length > 0 && (
          <ul className="tl-pesan" role="alert">
            {pesan.map((p, i) => (
              <li key={i}>{p}</li>
            ))}
          </ul>
        )}
        {gagal !== '' && (
          <p className="tl-pesan" role="alert">
            {gagal}
          </p>
        )}

        <StripTabNavigasi tab={SUB_TAB_SHARE} aktif={SUB_TAB_SHARE[0]} onPilih={setSub} />

        {sub === SUB_TAB_SHARE[0] && (
          <>
            {adaFac && (
              <p className="tl-share-rnm">
                {SHARE_NP.shareKeRnm} <strong>{share.RnmShareDeducted}</strong> {SHARE_NP.persen}
              </p>
            )}
            {share.Share.length === 0 ? (
              <Kosong pesan={SHARE_NP.tanpaBaris} petunjuk={SHARE_NP.petunjukKosong} />
            ) : (
              <div className="table-wrap trin__limit-grid tl-share-grid">
                <table className="trin__tabel">
                  <thead>
                    <tr>
                      <th scope="col" className="tl-kolom-buka" aria-label={SHARE_NP.bukaRincian} />
                      {KEPALA_SHARE.slice(0, 5).map((k, i) => (
                        <th key={i} scope="col" aria-hidden={k === '' ? true : undefined}>
                          {k}
                        </th>
                      ))}
                      {/* Kepala ekspor hanya di atas sel nilai; mata uangnya
                          berkepala kosong, jadi satu kepala merangkul keduanya. */}
                      <th scope="col" colSpan={2}>
                        {KEPALA_SHARE[6]}
                      </th>
                      <th scope="col" colSpan={2}>
                        {KEPALA_SHARE[8]}
                      </th>
                      <th scope="col" colSpan={2}>
                        {KEPALA_SHARE[10]}
                      </th>
                      <th scope="col" colSpan={2}>
                        {KEPALA_SHARE[12]}
                      </th>
                      <th scope="col">{KEPALA_SHARE[13]}</th>
                    </tr>
                  </thead>
                  <tbody>
                    {share.Share.map((b, i) => {
                      const terbuka = buka === i
                      const [c1, v1] = selPasangan(b.RnmLimitList, 0)
                      const [c2, v2] = selPasangan(b.RnmLimitList, 1)
                      const [m1, n1] = selPasangan(b.GrossPremiumList, 0)
                      const [m2, n2] = selPasangan(b.GrossPremiumList, 1)
                      return [
                        <tr key={`b${i}`} className={terbuka ? 'tl-share-baris--buka' : undefined}>
                          <td>
                            {/* Navigasi, bukan `<button>` — tetap hidup di mode lihat. */}
                            <TombolNavigasi
                              className="btn btn--ghost btn--sm tl-share-buka"
                              terbuka={terbuka}
                              label={`${SHARE_NP.bukaRincian} ${judulBarisShare(b)}`}
                              onKlik={() => {
                                setBuka(terbuka ? null : i)
                              }}
                            >
                              <IkonChevronKanan />
                            </TombolNavigasi>
                            {(pesanBaris[i]?.length ?? 0) > 0 && (
                              <span className="tl-share-tanda" role="img" aria-label={SHARE_NP.adaPesan} title={pesanBaris[i]?.join(' · ')}>
                                !
                              </span>
                            )}
                          </td>
                          <td>{b.LayerType}</td>
                          <td>{b.Layer}</td>
                          <td className="tl-share-of">{SHARE_NP.bagianOf}</td>
                          <td>{b.LayerPartType}</td>
                          <td>{b.LayerPart}</td>
                          <td className="tl-share-mu">{c1}</td>
                          <td className="tl-angka">{v1}</td>
                          <td className="tl-share-mu">{c2}</td>
                          <td className="tl-angka">{v2}</td>
                          <td className="tl-share-mu">{m1}</td>
                          <td className="tl-angka">{n1}</td>
                          <td className="tl-share-mu">{m2}</td>
                          <td className="tl-angka">{n2}</td>
                          <td className="tl-angka">{persen(b.RNMShare)}</td>
                        </tr>,
                        terbuka && (
                          <tr key={`r${i}`} className="tl-share-panel">
                            <td colSpan={15}>
                              <RincianShare
                                b={b}
                                induk={induk[b.TreatyGroupList[0]?.TreatyGroupID ?? ''] ?? []}
                                indukManual={induk[''] ?? []}
                                opsi={opsi}
                                pesan={pesanBaris[i] ?? []}
                                bisaUbah={bisaUbah}
                                modeUbah={modeUbah}
                                onUbah={(baru) => setShare({ ...share, Share: ganti(share.Share, i, baru) })}
                                hitung={(aksi, o = {}) => {
                                  const dasar = o.dasar === undefined ? share : { ...share, Share: ganti(share.Share, i, o.dasar) }
                                  hitung(aksi, { indeks: i, baris: o.baris, sts: o.sts, dasar })
                                }}
                              />
                            </td>
                          </tr>
                        ),
                      ]
                    })}
                  </tbody>
                </table>
              </div>
            )}

            <Bagian judul={SHARE_NP.ringkasan}>
              <div className="table-wrap trin__limit-grid">
                <table className="trin__tabel">
                  <thead>
                    <tr>
                      {KOLOM_RINGKASAN_SHARE.map((k) => (
                        <th key={k.kunci} scope="col">
                          {k.label}
                        </th>
                      ))}
                    </tr>
                  </thead>
                  <tbody>
                    {share.LimitShareSummaryList.length === 0 && (
                      <tr>
                        <td colSpan={KOLOM_RINGKASAN_SHARE.length}>{SHARE_NP.tanpaBaris}</td>
                      </tr>
                    )}
                    {share.LimitShareSummaryList.map((r, i) => (
                      <tr key={i}>
                        {KOLOM_RINGKASAN_SHARE.map((k) => (
                          <td key={k.kunci} className={k.kunci === 'Note' ? undefined : 'tl-angka'}>
                            {k.kunci === 'Note' ? r.Note : uang(r[k.kunci])}
                          </td>
                        ))}
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </Bagian>

            <Bagian
              judul={SHARE_NP.totalSemua}
              aksi={
                modeUbah && (
                  <div className="trin__aksi">
                    <button type="button" className="btn btn--primary btn--sm" disabled={!bisaUbah} onClick={() => hitung('total')}>
                      {SHARE_NP.perbaruiTotal}
                    </button>
                    {/* `TreatyIn.ViewState != '1' && TreatyMasterInEDM`. */}
                    {kontrakRevisi(edmState) && (
                      <button type="button" className="btn btn--sm" disabled={!bisaUbah} onClick={() => hitung('nilai-share')}>
                        {SHARE_NP.perbaruiNilai}
                      </button>
                    )}
                  </div>
                )
              }
            >
              <div className="tl-share-total-grid">
                {GRID_TOTAL_SHARE.flatMap((baris, r) =>
                  baris.map((g, c) =>
                    g === null ? (
                      <div key={`${r}-${c}`} aria-hidden="true" />
                    ) : (
                      <GridTotal key={g.kunci} judul={g.judul} baris={share.Total[g.kunci] ?? []} />
                    ),
                  ),
                )}
              </div>
            </Bagian>
          </>
        )}
      </div>
    </Panel>
  )
}
