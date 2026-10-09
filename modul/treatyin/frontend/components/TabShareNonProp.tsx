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
// MENIMPA `pyReadOnly`. Di panel rincian, Cover · Deduction Details ·
// Spreading Type terkunci HANYA bila `TreatyIn.ViewState = 1`.
// ⛔ KECUALI sel ber-`pyEditOptions = Read-only` yang modenya
// `pyDisabledNew = always` (8 Oktober 2026; kerangka Adjustment `baca:
// "selalu"`): Layer Type · Layer · Layer Part Type · Layer Part (@63296 dst.),
// grid Treaty Group (`CoBListReadOnly`), dan grid spreading bernama — SELALU
// baca-saja, seperti Pega.
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

import { formatNumber } from '../../../../inti/frontend/lib/format'
import { Field } from '../../../../inti/frontend/components/ui/dasar'
import { PilihCari as Pilih } from '../../../../inti/frontend/components/ui/pilihSaring'
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
  type RingkasanShareNP,
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
// ⛔ `DropdownDaftar`, BUKAN `IsianAuto`: yang kedua `<datalist>` —
// ketikan bebas lolos jadi nilai, dan daftarnya tidak menyaring.
// Permintaan pemilik proses 8 Oktober 2026 untuk seluruh dropdown
// Treaty In dan Adjustment: mengetik harus terasa seperti mencari.
import { DropdownDaftar } from './IsianAuto'
import { klikBaris, StripTabNavigasi, TombolNavigasi } from './navigasi'
import { PemicuUbah, usePemicuUbah } from './pemicuUbah'
import { BlokPega, DeretTombolPega, GridNilaiPega, GridPega, TeksPega } from './gridPega'
import { TataPegaBlok } from './tataPega'
import { Bagian, TombolHapus, TombolTambah } from './limitsUI'
import { saringAngka } from './saringAngka'
import { formatLimit } from './TabLimitsProp'

const uang = (v: string) => formatLimit('uang', 2, v)
/** `pxNumber` tanpa desimal ekspor — angka apa adanya, berpemisah ribuan. */
/**
 * Sel angka `Summarry of RNM Share` — DUA desimal, dan nol di ekor dibuang.
 *
 * ⛔ Dahulu `formatLimit('uang', -1, …)` = APA ADANYA, dan hasil pembagian
 * membawa ekornya utuh ke layar: `11.652.832,7749999503125`. Angka seperti
 * itu bukan hanya sukar dibaca, ia MEMBUNGKUS ke baris berikutnya dan
 * merusak barisnya (laporan pemilik proses 8 Oktober 2026).
 *
 * ⭐ `formatNumber(v, 2)` membulatkan ke dua desimal LALU membuang nol di
 * ekor, jadi catatan ekspor tetap terpenuhi: `175.000.000` dan `0` tampil
 * tanpa `,00`, persis seperti Pega menampilkannya. Yang berubah hanya nilai
 * yang memang berpecahan.
 */
const angkaRingkasan = (v: string) => formatNumber(v, 2)
/** Lebar kolom `Summarry of RNM Share` dari ekspor (`pyWidth`). */
const LEBAR_RINGKASAN_SHARE = [293, 199, 204, 123, 119, 205, 210, 150, 144] as const
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
  angka = false,
  onUbah,
  onLepas,
}: {
  label: string
  nilai: string
  bisaUbah: boolean
  placeholder?: string
  /** Kontrol Number di ekspor — huruf/simbol ditolak (`saringAngka`). */
  angka?: boolean
  onUbah: (v: string) => void
  onLepas?: () => void
}) {
  // Peristiwa `change` Pega — hanya bila nilainya berubah (`pemicuUbah.tsx`).
  const pemicu = usePemicuUbah(nilai, bisaUbah ? onLepas : undefined)
  return (
    <div className="trin__limit-medan" onFocus={pemicu.masuk} onBlur={pemicu.keluar}>
      <Field
        label={label}
        value={nilai}
        readOnly={!bisaUbah}
        placeholder={bisaUbah ? placeholder : undefined}
        onChange={(v) => {
          onUbah(angka ? saringAngka(v) : v)
        }}
      />
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

/**
 * Dropdown reasuradur grid Reinsurer / Facultative Reinsurers.
 *
 * Nilai pilihan = `ReinsID`. Nama kembar (satu nama, beberapa pengenal —
 * `PilihanWarisan.kembar`) diberi pengenalnya supaya yang memilih tidak
 * menebak. Baris lama yang `ReinsID`-nya kosong dicocokkan lewat nama; yang
 * tetap tidak ditemukan tampil apa adanya sebagai pilihan tersendiri, tidak
 * dikosongkan.
 */
export function nilaiReins(b: Pick<BarisReinsShare, 'ReinsID' | 'ReinsName'>, pilihan: readonly PilihanWarisan[]): string {
  if (b.ReinsID !== '') return b.ReinsID
  return pilihan.find((o) => o.nama === b.ReinsName)?.id ?? b.ReinsName
}
export function opsiReins(b: Pick<BarisReinsShare, 'ReinsID' | 'ReinsName'>, pilihan: readonly PilihanWarisan[]): OpsiPilihan[] {
  const opsi = pilihan.map((o) => ({ value: o.id, label: o.kembar ? `${o.nama} (${o.id})` : o.nama }))
  const kini = nilaiReins(b, pilihan)
  if (kini !== '' && !opsi.some((o) => o.value === kini)) opsi.unshift({ value: kini, label: b.ReinsName || kini })
  return opsi
}
export function pilihReins<T extends Pick<BarisReinsShare, 'ReinsID' | 'ReinsName'>>(b: T, pilihan: readonly PilihanWarisan[], v: string): T {
  const p = pilihan.find((o) => o.id === v)
  if (p !== undefined) return { ...b, ReinsID: p.id, ReinsName: p.nama }
  if (v === '') return { ...b, ReinsID: '', ReinsName: '' }
  return b
}

/**
 * Grid reasuradur — `ShareReins` atau `ShareFacultativeReinsurers`.
 *
 * ⭐ Bentuk Pega (8 Oktober 2026): grid `row` berkolom ekspor (lebar 424 ·
 * 180 · 188, kolom tombol 159), `Add` di sel kepala kolom tombol, `Delete`
 * per baris — bukan kepala berlencana dengan tombol Add merah di luar grid.
 */
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
    <GridPega
      label={kolom[0]}
      kolom={[
        {
          judul: kolom[0] ?? '',
          lebar: 424,
          // `.ReinsName` — daftar `BrowseAgentNusaRe_RD`: `.ClientName`,
          // `.ID → .ReinsID`.
          // ⛔ DROPDOWN, bukan autocomplete — permintaan pemakai 9 Oktober
          // 2026: *"jangan auto complete melainkan dropdown seperti yang
          // lainnya"* (ekspor: pxAutoComplete). Lihat `pilihReins`.
          isi: (b, i) =>
            bisaUbah ? (
              <PilihMedan
                label=""
                nilai={nilaiReins(b, pilihan)}
                opsi={opsiReins(b, pilihan)}
                bisaUbah
                onUbah={(v) => {
                  onUbah(ganti(baris, i, pilihReins(b, pilihan, v)))
                }}
              />
            ) : (
              <Field label="" value={b.ReinsName} readOnly onChange={() => undefined} />
            ),
        },
        {
          judul: kolom[1] ?? '',
          lebar: 180,
          isi: (b, i) => (
            <Medan
              label=""
              nilai={b.Layer}
              bisaUbah={bisaUbah}
              placeholder={SHARE_NP.placeholderTeks}
              onUbah={(v) => onUbah(ganti(baris, i, { ...b, Layer: v }))}
            />
          ),
        },
        {
          judul: kolom[2] ?? '',
          lebar: 188,
          isi: (b, i) => (
            <Medan
              label=""
              nilai={bisaUbah ? b.SharePct : persen(b.SharePct)}
              bisaUbah={bisaUbah}
              placeholder={SHARE_NP.placeholderAngka}
              angka
              onUbah={(v) => onUbah(ganti(baris, i, { ...b, SharePct: v }))}
            />
          ),
        },
      ]}
      baris={baris}
      tombol={
        bisaUbah
          ? {
              lebar: 159,
              // `TreatyInNonAddItem(sharereins | sharefacname)`: baris ber-ID kosong.
              tambah: {
                label: SHARE_NP.tambah,
                onKlik: () => {
                  onUbah([...baris, { ID: '', ReinsID: '', ReinsName: '', Layer: '', SharePct: '' }])
                },
              },
              hapus: {
                label: SHARE_NP.hapus,
                akses: (_, i) => `${SHARE_NP.hapus} ${kolom[0] ?? ''} ${String(i + 1)}`,
                onKlik: (i) => {
                  onUbah(baris.filter((_, j) => j !== i))
                },
              },
            }
          : undefined
      }
    />
  )
}

/**
 * Grid baca-saja `Currency · Value` — grid nilai Pega yang sempit.
 * Grid total berkepala `judul · Value`; grid rincian panel (`satuKepala`)
 * berkepala `judul · ''`, seperti di ekspor (`RNM Limit`, `''`).
 */
function GridTotal({ judul, baris, satuKepala = false }: { judul: string; baris: readonly NilaiShare[]; satuKepala?: boolean }) {
  return (
    <GridNilaiPega
      judul={judul}
      nilai={satuKepala ? '' : SHARE_NP.nilai}
      baris={baris}
      lebar={[194, 352]}
      tampil={uang}
    />
  )
}

/** Sel pasangan mata uang · nilai ke-`i` sebuah larik. */
function selPasangan(daftar: readonly NilaiShare[], i: number): [string, string] {
  const v = daftar[i]
  return v === undefined ? ['', ''] : [v.Currency, uang(v.Value)]
}

/** Teks tampil sebuah nilai dropdown `associated` (label opsi, atau nilainya). */
function teksPilihan(opsi: readonly OpsiPilihan[], nilai: string): string {
  return opsi.find((o) => o.value === nilai)?.label ?? nilai
}

/** Satu baris `TreatyGroupList` panel Share, berikut CoB bila kontraknya membawa. */
type GrupShareBaca = BarisShareNP['TreatyGroupList'][number] & {
  ClassOfBusinessList?: readonly { ClassOfBusiness?: string }[]
}

/**
 * «CoBListReadOnly» (`Section/CoBListReadOnly.xml`) — rincian satu baris grid
 * Treaty Group panel Share: `Treaty Group` baca-saja BERLABEL (@15187,
 * `pyEditOptions = Read-only`, `pyLabelFieldValue = Treaty Group`) lalu grid
 * `Class of Business` baca-saja (@45602, `readOnly`, tanpa tombol).
 *
 * ⚠️ `ClassOfBusinessList` belum ikut kontrak baris Share (`GrupShareNP`
 * services hanya `TreatyGroup`/`TreatyGroupID`) — dibaca bila ada; selain itu
 * grid menampilkan "No items", seperti Pega tanpa baris.
 */
function RincianCoBBaca({ g }: { g: GrupShareBaca }) {
  return (
    <TataPegaBlok tata="tumpuk">
      <Field label={SHARE_NP.treatyGroup} value={g.TreatyGroup} readOnly onChange={() => undefined} />
      <GridPega
        label={SHARE_NP.kelasBisnis}
        kolom={[{ judul: SHARE_NP.kelasBisnis, lebar: 196, isi: (c: { ClassOfBusiness?: string }) => c.ClassOfBusiness ?? '' }]}
        baris={g.ClassOfBusinessList ?? []}
      />
    </TataPegaBlok>
  )
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
  // ⭐ DUA SPREADING, PERSIS PEGA — keputusan pemakai 9 Oktober 2026 (*"di
  // pega ada terdapat 2 spreading … yang tidak bisa ditambah itu khusus untuk
  // spreading lama"*), menggantikan "satu grid untuk kedua cabang" 8 Oktober.
  //
  //   `.SpreadingTypeXOL != ''` → spreading LAMA: dropdown Spreading Type +
  //       grid `readOnly` (@408855, nol tombol) — isi `FetchQSfromMasterXOL`.
  //   `.SpreadingTypeXOL == ''` → spreading BARU: grid manual ber-Add/Delete
  //       (@54364, `AddSpreadingXOL`), dropdown Spreading Type TIDAK ada.
  //
  // ⚠️ `TreatyInXOLAddSpreading` TIDAK PERNAH mengisi Spreading Type — ia
  // menjalankan `FetchQSfromMasterXOL` bila sudah terisi, `SetSpreadingXOL`
  // bila kosong. Spreading lama = kontrak warisan yang menyimpannya.
  const gridSebar = (lama: boolean) => {
    const kolom = lama ? KOLOM_SPREADING : KOLOM_SPREADING_MANUAL
    const sunting = modeUbah && !lama
    return (
            <div className="table-wrap trin__limit-grid">
              <table className="trin__tabel trin__tabel--pega">
                <thead>
                  <tr>
                    {kolom.map((k) => (
                      <th key={k} scope="col">
                        {k}
                      </th>
                    ))}
                    {sunting && (
                      <th scope="col" className="tl-kolom-aksi">
                        <TombolTambah label={SHARE_NP.tambah} onClick={() => hitung('spreading-tambah')} />
                      </th>
                    )}
                  </tr>
                </thead>
                <tbody>
                  {b.SpreadingListXOL.length === 0 && (
                    <tr>
                      <td colSpan={2 + (sunting ? 1 : 0)}>{SHARE_NP.tanpaBaris}</td>
                    </tr>
                  )}
                  {b.SpreadingListXOL.map((s, i) => (
                    <tr key={i}>
                      <td>
                        {sunting ? (
                          <select
                            className="field__input"
                            aria-label={kolom[0]}
                            value={s.ReinsTypeID}
                            onChange={(e) => {
                              // postValue saja — namanya diisi `SetSpreadingXOL`.
                              onUbah({ ...b, SpreadingListXOL: ganti(b.SpreadingListXOL, i, { ...s, ReinsTypeID: e.target.value }) })
                            }}
                          >
                            <option value="">{SHARE_NP.pilihKosong}</option>
                            {/* ⛔ Nilai TERSIMPAN yang tidak ada di daftar tetap
                                ditawarkan, bukan dijatuhkan: baris hasil
                                `FetchQSfromMasterXOL` ber-`ReinsTypeID` ANAK
                                (`QS (OR)`, `QS (R/I)`) sementara daftar ini
                                berisi INDUK. Tanpa ini namanya hilang dari
                                layar begitu gridnya dapat disunting. */}
                            {s.ReinsTypeID !== '' && !indukManual.some((p) => p.reinsTypeId === s.ReinsTypeID) && (
                              <option value={s.ReinsTypeID}>{s.ReinsTypeName || s.ReinsTypeID}</option>
                            )}
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
                          nilai={sunting ? s.Pct : persen(s.Pct)}
                          bisaUbah={sunting}
                          angka
                          onUbah={(v) => onUbah({ ...b, SpreadingListXOL: ganti(b.SpreadingListXOL, i, { ...s, Pct: v }) })}
                          onLepas={() => hitung('spreading-pct')}
                        />
                      </td>
                      {sunting && (
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
    )
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
      {/* Urutan & tata = `Share.xml` `Default` @546: % RNM Share (`Inline
          labels left` @848) · baris layer (`Inline` @1965) · grid Treaty
          Group · Cover · Deduction Details · Spreading Type · Spreading. */}
      <TataPegaBlok tata="alir">
        <Medan
          label={SHARE_NP.persenRnm}
          nilai={b.RNMShare}
          bisaUbah={bisaUbah}
          placeholder={SHARE_NP.placeholderAngka}
          angka
          onUbah={(v) => onUbah({ ...b, RNMShare: v })}
          // `TreatyInXOLAddSpreadingDetail(idx)`.
          onLepas={() => hitung('rnm-baris')}
        />
      </TataPegaBlok>
      {/* Layer Type · Layer · "Part of" · Layer Part Type · Layer Part —
          `Inline` @62575. SELALU baca-saja (`pyEditOptions = Read-only`,
          `pyDisabledNew = always` @65740 · @71745 · @81918 · @87930;
          kerangka Adjustment `baca: "selalu"`) dan TANPA label
          (`pyLabelReserveSpace = false`, `pyIncludeLabel = false`) — teks
          sebaris, seperti baris grid Pega di atasnya. */}
      <TataPegaBlok tata="alir">
        <TeksPega>{teksPilihan(opsi.jenisLayer ?? [], b.LayerType)}</TeksPega>
        <TeksPega>{b.Layer}</TeksPega>
        <TeksPega>{SHARE_NP.bagianOf}</TeksPega>
        <TeksPega>{teksPilihan(opsi.jenisLayer ?? [], b.LayerPartType)}</TeksPega>
        <TeksPega>{b.LayerPart}</TeksPega>
      </TataPegaBlok>

      {/* `.TreatyGroupList` @117429 — grid `masterDetail` baca-saja, SATU
          kolom `Treaty Group` (lebar 196), TANPA judul blok (wadah @101614
          tanpa kepala) dan tanpa tombol; ▸ membuka «CoBListReadOnly»
          (`templatBaris … TreatyInLimitLayer!pyGridRowDetails`). */}
      <GridPega
        label={SHARE_NP.treatyGroup}
        labelBuka={SHARE_NP.bukaRincian}
        kolom={[{ judul: SHARE_NP.treatyGroup, lebar: 196, isi: (g: GrupShareBaca) => g.TreatyGroup }]}
        baris={b.TreatyGroupList}
        rincian={(g: GrupShareBaca) => <RincianCoBBaca g={g} />}
      />

      {/* `.Cover` pxDropdown @162739 — sel `Default` sesudah grid Treaty
          Group: label `Cover` di ATAS (`pyLabelReserveSpace`, label properti),
          kotak SELEBAR ISINYA (dropdown Pega tidak melar selebar panel).
          Terkunci hanya bila `ViewState = 1`. */}
      <div style={{ width: 'fit-content', maxWidth: '100%' }}>
        <PilihMedan
          label={SHARE_NP.cover}
          nilai={b.Cover}
          opsi={opsi.cover ?? []}
          bisaUbah={modeUbah}
          onUbah={(v) => onUbah({ ...b, Cover: v })}
        />
      </div>

      {/* `Deduction Details` @171027 (wadah berjudul) — grid `row` berkolom
          ekspor (203 · 195 · 199 · 30 · 210 · 117, kolom tombol 118): `Add`
          di sel KEPALA kolom tombol (@217024), `Remove` per baris (@269980);
          keduanya `ViewState != '1'` dan mati bila `EDMMaterialType = 2`. */}
      <BlokPega judul={SHARE_NP.deduksi}>
        <GridPega
          label={SHARE_NP.deduksi}
          kolom={[
            {
              judul: KOLOM_DEDUKSI_SHARE[0],
              lebar: 203,
              isi: (d: BarisDeduksiShare, i: number) => (
                <Medan label="" nilai={d.Comment} bisaUbah={bisaUbah} onUbah={(v) => ubahDeduksi(i, { ...d, Comment: v })} />
              ),
            },
            {
              judul: KOLOM_DEDUKSI_SHARE[1],
              lebar: 195,
              isi: (d: BarisDeduksiShare, i: number) => (
                /* `.Currency` pxAutoComplete `BrowseCurrency_RD` — dipanggil
                   tanpa parameter, jadi filter `Currency`/`ID` dilewati:
                   isinya = daftar mata uang tab Limits (`!= ITL`).
                   `change` → `CalculateDeduction(index)`. */
                <PemicuUbah
                  className="trin__limit-medan"
                  nilai={d.Currency}
                  aktif={bisaUbah}
                  aksi={() => hitung('deduksi', { baris: i, sts: '' })}
                >
                  <DropdownDaftar
                    label=""
                    nilai={d.Currency}
                    pilihan={opsi.mataUang}
                    bisaUbah={bisaUbah}
                    onPilih={(nama, id) => ubahDeduksi(i, { ...d, Currency: nama, CurrencyID: id })}
                  />
                </PemicuUbah>
              ),
            },
            {
              judul: KOLOM_DEDUKSI_SHARE[2],
              lebar: 199,
              isi: (d: BarisDeduksiShare, i: number) => (
                <Medan
                  label=""
                  nilai={bisaUbah ? d.Deduction : uang(d.Deduction)}
                  bisaUbah={bisaUbah}
                  angka
                  onUbah={(v) => ubahDeduksi(i, { ...d, Deduction: v })}
                  onLepas={() => hitung('deduksi', { baris: i, sts: 'val' })}
                />
              ),
            },
            { judul: KOLOM_DEDUKSI_SHARE[3], lebar: 30, isi: () => <span className="tl-share-atau">or</span> },
            {
              judul: KOLOM_DEDUKSI_SHARE[4],
              lebar: 210,
              isi: (d: BarisDeduksiShare, i: number) => (
                <Medan
                  label=""
                  nilai={bisaUbah ? d.DeductionPct : persen(d.DeductionPct)}
                  bisaUbah={bisaUbah}
                  angka
                  onUbah={(v) => ubahDeduksi(i, { ...d, DeductionPct: v })}
                  onLepas={() => hitung('deduksi', { baris: i, sts: 'pct' })}
                />
              ),
            },
            {
              judul: KOLOM_DEDUKSI_SHARE[5],
              lebar: 117,
              isi: (d: BarisDeduksiShare, i: number) => (
                <input
                  type="checkbox"
                  aria-label={KOLOM_DEDUKSI_SHARE[5]}
                  checked={d.DeductionPctCalculate === 'true'}
                  disabled={!bisaUbah}
                  onChange={(e) => ubahDeduksi(i, { ...d, DeductionPctCalculate: e.target.checked ? 'true' : 'false' })}
                />
              ),
            },
          ]}
          baris={b.DeductionList}
          tombol={
            modeUbah
              ? {
                  lebar: 118,
                  // `AddDeduction` — baris kosong.
                  tambah: {
                    label: SHARE_NP.tambah,
                    mati: !bisaUbah,
                    onKlik: () => {
                      onUbah({
                        ...b,
                        DeductionList: [...b.DeductionList, { Comment: '', Currency: '', CurrencyID: '', Deduction: '', DeductionPct: '' }],
                      })
                    },
                  },
                  hapus: {
                    label: SHARE_NP.hapusBaris,
                    mati: !bisaUbah,
                    akses: (_, i) => `${SHARE_NP.hapusBaris} ${SHARE_NP.deduksi} ${String(i + 1)}`,
                    onKlik: (i) => {
                      // deleteRow, lalu `CalculateDeduction(index)`.
                      const baru = { ...b, DeductionList: b.DeductionList.filter((_, j) => j !== i) }
                      onUbah(baru)
                      hitung('deduksi', { baris: i, sts: '', dasar: baru })
                    },
                  },
                }
              : undefined
          }
        />
      </BlokPega>

      {/* Spreading Type — layout sendiri (`Stacked with labels left`
          @326811) DI ATAS blok Spreading, bukan di dalamnya: label KIRI. */}
      {/* ⛔ HANYA di cabang spreading lama — Pega: blok `.SpreadingTypeXOL!=''`. */}
      {bernama && (
        <TataPegaBlok tata="kiri">
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
        </TataPegaBlok>
      )}
      {/* Spreading bernama: `Inline grid 30 70` @11860 = baris [30% blok
          `Spreading` @11969 | 70% blok `hidden, reference` @17096] — blok
          rujukan SAUDARA kotak Spreading, di sebelah kanannya. Spreading
          manual: blok `Spreading` tersendiri @54364 (`.SpreadingTypeXOL = ''`). */}
      {bernama ? (
        <TataPegaBlok tata="t3070">
          <Bagian judul={SHARE_NP.spreading}>
            {/* @408855 — spreading LAMA: grid `readOnly`, nol tombol (`gridSebar`). */}
            {gridSebar(true)}
            {/* Kaki grid @442259: label `Total Pct` (@444655) + sel
                `.pyTemplateInputBox` pxPercentage baca-saja TANPA properti —
                Pega selalu menampilkannya kosong, jadi hanya labelnya. */}
            <TataPegaBlok tata="alir">
              <TeksPega>{SHARE_NP.totalPct}</TeksPega>
            </TataPegaBlok>
            {/* `Inline grid double` @464235 → `Inline labels left` @473224:
                `Spreading Total Pct` (`pyLabelFieldValue` @475xxx, Decimal
                baca-saja `.SpreadingTotalPctXOL`) · teks `%` @486235. Tombol
                @479141 ber-`NEVER`. Slot kanan grid ganda kosong. */}
            <TataPegaBlok tata="g2">
              <TataPegaBlok tata="alir">
                <Field label={SHARE_NP.spreadingTotalPct} value={uang(b.SpreadingTotalPctXOL)} readOnly onChange={() => undefined} />
                <TeksPega>{SHARE_NP.persen}</TeksPega>
              </TataPegaBlok>
            </TataPegaBlok>
          </Bagian>
          <TataPegaBlok tata="tumpuk">
            {/* Blok `hidden, reference` @17096 — TAMPIL bersama Spreading
                bernama (`NOHEADER`, nol syarat sendiri): lima layout `Inline
                grid triple` bertumpuk (@17588 · @24795 · @31991 · @39187 ·
                @46426), dasar · OR · R/I per baris. */}
            {GRID_RINCIAN_SPREADING.map((barisGrid, r) => (
              <TataPegaBlok key={r} tata="g3">
                {barisGrid.map((g) => (
                  <GridTotal key={g.kunci} judul={g.judul} baris={b[g.kunci as keyof BarisShareNP] as readonly NilaiShare[]} satuKepala />
                ))}
              </TataPegaBlok>
            ))}
          </TataPegaBlok>
        </TataPegaBlok>
      ) : (
        <Bagian judul={SHARE_NP.spreading}>
          {gridSebar(false)}
          <p className="tl-share-catatan">
            {SHARE_NP.totalSharePct} <strong>{persen(b.SpreadingTotalPctXOL)}</strong>
          </p>
        </Bagian>
      )}
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
    // ⛔ FILTER TREATY GROUP SENGAJA DIBUANG — PENYIMPANGAN DARI PEGA, dan ini
    // pernyataannya.
    //
    // `Section/Share.xml` dan `Section/DetailShare.xml` MENGIRIM
    // `TreatyGroupID`, begitu pula `FetchQSfromMaster(XOL)`. Mengikutinya berarti
    // hanya induk yang terdaftar di Treaty Group baris itu yang dapat dipilih —
    // dan untuk kontrak yang dilaporkan 8 Oktober 2026 daftarnya kosong, padahal
    // susunannya ADA di `PROPORTIONALARRG` (dinyatakan pemilik proses).
    //
    // Keputusan pemilik proses, 8 Oktober 2026: *"gimana pun caranya asal itu ada
    // isinya"*.
    //
    // ⚠️ DIBUANG DI KEDUA TEMPAT SEKALIGUS, dan itu syaratnya. Membuangnya hanya
    // di dropdown — yang sempat terjadi — membuat layar menawarkan induk yang
    // pencariannya sendiri tidak dapat menemukan: Spreading Type terpilih,
    // grid spreading tetap kosong. Setengah penyimpangan lebih buruk daripada
    // keduanya, sebab ia terbaca seperti berhasil.
    //
    // ⭐ `TreatyDescID = "10001"` TETAP dikirim: nol pemanggil yang menghilangkannya,
    // dan ia tidak pernah menjadi sebab daftar kosong.
    //
    // Pasangannya di backend: `fetchQS` (`hitung_share_np.go`).
    for (const g of ['']) {
      if (induk[g] !== undefined) continue
      ambilIndukSpreading('', '10001', commencement)
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

  // ⭐ Urutan = Section `TreatyInTabsNonProportional` tab `Share` (bentuk
  // Pega, 8 Oktober 2026): blok `Share` — medan, `Update Summary`, grid
  // reasuradur — lalu `RNM Share`, `Summarry of RNM Share`, dan blok
  // `Total All Layers RNM Share` dengan tombolnya DI BAWAH grid.
  return (
    <div className="trin__blok trin__tab">
      <BlokPega judul={SHARE_NP.judul}>
        {share.IsProRate === 'true' && <p className="tl-share-catatan">{SHARE_NP.nonProRate}</p>}
        {/* Tata ekspor (`TreatyInTabsNonProportional.xml`): `Inline grid
            double` @58958 = [ `Inline grid double` @59257 ( `Stacked with
            labels left` @59555 RNM · Brokerage | centang Across The Board )
            | `Stacked with labels left` @61418 retro ]. */}
        <TataPegaBlok tata="g2">
          <TataPegaBlok tata="g2">
            <TataPegaBlok tata="kiri">
              <Medan
                label={SHARE_NP.persenRnm}
                nilai={bisaUbah ? share.RNMShare : persen(share.RNMShare)}
                bisaUbah={bisaUbah}
                placeholder={SHARE_NP.placeholderAngka}
                angka
                onUbah={(v) => ubahAkar('RNMShare', v)}
                onLepas={() => hitung('rnm')}
              />
              <Medan
                label={SHARE_NP.persenBrokerage}
                nilai={bisaUbah ? share.BrokeragePercent : persen(share.BrokeragePercent)}
                bisaUbah={bisaUbah}
                placeholder={SHARE_NP.placeholderAngka}
                angka
                onUbah={(v) => ubahAkar('BrokeragePercent', v)}
                onLepas={() => hitung('brokerage')}
              />
            </TataPegaBlok>
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
          </TataPegaBlok>
          <TataPegaBlok tata="kiri">
            <Medan
              label={SHARE_NP.shareKeRetro}
              nilai={bisaUbah ? share.FacultativeShare : persen(share.FacultativeShare)}
              bisaUbah={bisaUbah}
              placeholder={SHARE_NP.placeholderAngka}
              angka
              onUbah={(v) => ubahAkar('FacultativeShare', v)}
              onLepas={() => hitung('fac')}
            />
            {adaFac && (
              <Medan
                label={SHARE_NP.brokerageRetro}
                nilai={bisaUbah ? share.FacultativeShareBrokerage : persen(share.FacultativeShareBrokerage)}
                bisaUbah={bisaUbah}
                placeholder={SHARE_NP.placeholderAngka}
                angka
                onUbah={(v) => ubahAkar('FacultativeShareBrokerage', v)}
                onLepas={() => hitung('fac')}
              />
            )}
          </TataPegaBlok>
        </TataPegaBlok>

        {tombolRingkasan !== false && <DeretTombolPega>{tombolRingkasan}</DeretTombolPega>}

        {/* `Inline grid double` @63466: grid Reinsurer Name | grid
            Facultative Reinsurers — tiap grid separuh lebar; bila grid
            kedua tersembunyi, separuh kanan kosong seperti Pega. */}
        <TataPegaBlok tata="g2">
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
        </TataPegaBlok>
      </BlokPega>

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
            {/* ⭐ Grid Pega tetap tampil tanpa baris — kepala + "No items".
                ⛔ TANPA `.table-wrap` (8 Oktober 2026): pembungkus itu
                menggulir sendiri (`max-height: 62vh`, `overflow: auto`) —
                panel rincian yang dibuka terpotong dengan penggulir di dalam
                grid. Pega tidak punya penggulir dalam: tabel `table-layout:
                fixed` selebar wadah, rincian sebaris penuh di bawah barisnya. */}
            <div className="trin__limit-grid tl-share-grid">
                <table className="trin__tabel trin__tabel--pega">
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
                    {share.Share.length === 0 && (
                      <tr>
                        <td colSpan={15} className="trin__kosong-pega">
                          {SHARE_NP.tanpaBaris}
                        </td>
                      </tr>
                    )}
                    {share.Share.map((b, i) => {
                      const terbuka = buka === i
                      const [c1, v1] = selPasangan(b.RnmLimitList, 0)
                      const [c2, v2] = selPasangan(b.RnmLimitList, 1)
                      const [m1, n1] = selPasangan(b.GrossPremiumList, 0)
                      const [m2, n2] = selPasangan(b.GrossPremiumList, 1)
                      return [
                        <tr
                          key={`b${i}`}
                          className={'trin__baris-buka' + (terbuka ? ' tl-share-baris--buka' : '')}
                          onClick={(e) => {
                            // ⭐ Klik baris = klik panah (`klikBaris`).
                            klikBaris(e, () => {
                              setBuka(terbuka ? null : i)
                            })
                          }}
                        >
                          <td>
                            {/* Navigasi, bukan `<button>` — tetap hidup di mode lihat. */}
                            <TombolNavigasi
                              className="trin__buka-tombol tl-share-buka"
                              terbuka={terbuka}
                              label={`${SHARE_NP.bukaRincian} ${judulBarisShare(b)}`}
                              onKlik={() => {
                                setBuka(terbuka ? null : i)
                              }}
                            >
                              {terbuka ? '▾' : '▸'}
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
                          <tr key={`r${i}`} className="tl-share-panel trin__rincian">
                            <td colSpan={15}>
                              <RincianShare
                                b={b}
                                induk={induk[''] ?? []}
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

            {/* `Summarry of RNM Share` — `pxNumber` TANPA desimal di ekspor:
                Pega menampilkan `4.000.000.000` dan `0`, bukan `,00`. */}
            <BlokPega judul={SHARE_NP.ringkasan}>
              <GridPega
                label={SHARE_NP.ringkasan}
                kolom={KOLOM_RINGKASAN_SHARE.map((k, c) => ({
                  judul: k.label,
                  lebar: LEBAR_RINGKASAN_SHARE[c] ?? 150,
                  angka: k.kunci !== 'Note',
                  isi: (r: RingkasanShareNP) => (k.kunci === 'Note' ? r.Note : angkaRingkasan(r[k.kunci])),
                }))}
                baris={share.LimitShareSummaryList}
              />
            </BlokPega>

            {/* `Total All Layers RNM Share` @84885 (`Default` @85082) — LIMA
                layout `Inline grid triple` bertumpuk, satu per baris
                `GRID_TOTAL_SHARE`: @85381 · @92542 · @95289 · @98036 ·
                @100784. Baris berbutir satu mengisi kolom pertama; slot
                sisanya kosong (layout terpisah, jadi tanpa sel pengisi).
                Tombol DI BAWAH. */}
            <BlokPega judul={SHARE_NP.totalSemua}>
              {GRID_TOTAL_SHARE.map((baris, r) => (
                <TataPegaBlok key={r} tata="g3">
                  {baris.map((g) =>
                    g === null ? null : <GridTotal key={g.kunci} judul={g.judul} baris={share.Total[g.kunci] ?? []} />,
                  )}
                </TataPegaBlok>
              ))}
              {modeUbah && (
                <DeretTombolPega>
                  <button type="button" className="btn btn--primary btn--sm" disabled={!bisaUbah} onClick={() => hitung('total')}>
                    {SHARE_NP.perbaruiTotal}
                  </button>
                  {/* `TreatyIn.ViewState != '1' && TreatyMasterInEDM`. */}
                  {kontrakRevisi(edmState) && (
                    <button type="button" className="btn btn--sm" disabled={!bisaUbah} onClick={() => hitung('nilai-share')}>
                      {SHARE_NP.perbaruiNilai}
                    </button>
                  )}
                </DeretTombolPega>
              )}
            </BlokPega>
          </>
        )}
    </div>
  )
}
