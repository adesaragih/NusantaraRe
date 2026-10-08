// Tab `Share` cabang PROPORSIONAL — bentuknya dari ekspor Pega.
//
// ---------------------------------------------------------------------
// ⛔ MENGAPA KOMPONEN TERPISAH DARI `TabShare.tsx`
// ---------------------------------------------------------------------
// Kedua cabang memakai tab bernama sama dan isi yang BERBEDA. Satu komponen
// untuk keduanya akan menampilkan kolom yang di cabang seberang tidak pernah
// ada — kesalahan yang sama yang membuat strip tab dipisah sejak awal
// (sebelas tab prop lawan dua belas non-prop, hanya enam namanya sama).
//
// ---------------------------------------------------------------------
// ⭐ BENTUKNYA DISALIN, BUKAN DIKARANG
// ---------------------------------------------------------------------
// `Section/TreatyInTabsProportional.xml` tab Share + `Section/TreatyInShareProp.xml`:
//
//   panel `Total Share`  Refresh · `% RNM Share` (RNMShareP) · `% Brokerage`
//                        (BrokeragePercentP) · `Option` (OptionLimit)
//   sub-tab `RNM Share`
//     `Share to RNM :`   RnmShareDeducted, bila `TreatyIn.FacultativeShare >0`
//     grid `Kind of Treaty` = page list `TreatyIn.Limits`, rincian baris
//       `TotalLimits`  → grid `.Detail` Treaty Group · % RNM Share (+ShareNote)
//         rincian `DetailShare` → % RNM Share · grid RNMShareList ·
//         Spreading (Spreading Type / SpreadingList) · Value Spreading OR / R/I
//     grid `Total Share RNM Limit` · `Total Value Spreading OR` · `… R/I`
//
// ---------------------------------------------------------------------
// ⭐ 7 Oktober 2026 — KETERGANTUNGAN DENGAN TAB LIMITS
// ---------------------------------------------------------------------
// Grid `Kind of Treaty` ADALAH `TreatyIn.Limits` tab Limits — dibaca dari
// penampung halaman (`../halaman.tsx`), jadi layer yang baru ditambah di tab
// Limits langsung tampil di sini, dan RNM Share yang Refresh hitung tampil
// pula di tab Limits/Achievement. Sebelum ini grid itu membaca proyeksi lain
// (`M_TREATY_IN2` per layer) dan Refresh mati.
//
// Rumus: rute `/hitung/share-prop` (`hitung_share_prop.go`) —
//   Refresh, `% RNM Share`, `Option` (change)   → `TreatyInPropshare`
//   `% RNM Share` rincian Detail (change)       → `TreatyInPropshareDetail`
//   `Spreading Type` (change)                   → `FetchQSfromMaster`
//
// ⭐ KEPUTUSAN PEMILIK PROSES 8 Oktober 2026 — Spreading DIPILIH, bukan
// diketik: *"pctnya di ambil dari table PROPORTIONALARRG ... select nya
// menggunakan rd"*. Rincian Detail SELALU memakai tata letak cabang
// `.SpreadingTypeID != ''` (`DetailShare`): dropdown `Spreading Type` dari
// RD `BrowseTreatyArrangement_ParentReinsMasterTrt`, lalu grid Reins Type ·
// Pct BACA-SAJA dari `PROPORTIONALARRG` (`BrowseTreatyArrangement_Limit_
// MstTrt_RD`) dan Value Spreading OR/R/I dari `FetchQSfromMaster`.
// Grid spreading MANUAL (`.SpreadingTypeID == ''`, `SetSpreadName`,
// `AddDelSpreadingTreatyin`) TIDAK dirender lagi.
//
// ⛔ BELUM: `Share to Other Retro`, `Brokerage From Other Retro`, dan grid
// `Facultative Reinsurers` — semuanya bersyarat `TreatyIn.IsMultipleRetro`
// dan milik Retro, keputusan §17 (`KEPUTUSAN-PENYELARASAN-REPO.md`).

import { Fragment, useEffect, useState } from 'react'

import { FieldAngka, Kosong, Panel, Pilih, StripTab } from '../../../../inti/frontend/components/ui/dasar'
import { ambilIndukSpreading, hitungShareProp, type MasukanShareProp, type NilaiTotalProp, type SimpulLimit, type SusunanSpreading } from '../api'
import { useProperti } from '../halaman'
import {
  GRID_TOTAL_RNM_SHARE,
  KOLOM_KIND_OF_TREATY_SHARE,
  SUB_TAB_SHARE,
  TOTAL_SHARE,
} from '../labels'
import { DETAIL_SHARE, KOLOM_TOTAL_LIMITS } from '../labelsShareProp'
import type { ModeForm } from '../mode'
import { selAngka } from './angka'
import { teksDari } from './TabLimitsProp'

/**
 * Pilihan `Option` — rule Property `OptionLimit`.
 *
 * ⭐ EKSPORNYA DATANG 7 Oktober 2026, dan menutup pertanyaan terbuka yang
 * berdiri sejak 6 Oktober. Pemilik proses mengirim tangkapan layar rule-nya:
 *
 *   Property  `OptionLimit`   kelas `ASM-FW-GISFW-Int-TREATY_IN`
 *   UI Control `pxTextInput`  Table type `Prompt List`
 *   Prompt values   1 → `Of Cession to R/I`
 *                   2 → `Of 100% Limit`
 *
 * ⛔ Sebelum ini label `2` ditampilkan sebagai kodenya apa adanya, BUKAN
 * karangan. Catatan lamanya berbunyi: *"pilihan berlabel karangan terbaca
 * benar sampai seseorang memilihnya"*. Menunggunya ternyata benar: tebakan
 * yang masuk akal (`Of Cession to R/I` lawan sesuatu tentang cession) akan
 * meleset — pasangannya ternyata `100% Limit`, bukan cession sama sekali.
 */
const OPSI_SHARE = [
  { value: '1', label: 'Of Cession to R/I' },
  { value: '2', label: 'Of 100% Limit' },
]

type Simpul = SimpulLimit

/** Larik anak satu simpul — yang tidak ada → kosong. */
function larik(s: Simpul, kunci: string): Simpul[] {
  const v = s[kunci]
  return Array.isArray(v) ? v : []
}

/** Ganti satu Detail di pohon — pohon baru, yang lain utuh. */
function gantiDetail(ls: readonly Simpul[], i: number, j: number, f: (d: Simpul) => Simpul): Simpul[] {
  return ls.map((l, x) => (x !== i ? l : { ...l, Detail: larik(l, 'Detail').map((d, y) => (y !== j ? d : f(d))) }))
}

/** Uang dua desimal — gambar 17: `3.000.000.000,00`. */
const uang = (v: string) => selAngka(['uang', 2], v)

/**
 * Satu grid nilai: kolom mata uang berjudul, kolom nilai.
 *
 * ⭐ `judulNilai` — grid TOTAL akar (`TreatyInShareProp`) berjudul `Value`;
 * grid rincian Detail (`DetailShare`: RNM Share, Value Spreading OR/R/I)
 * berkepala KOSONG di kolom keduanya — sel LABEL tanpa teks di ekspor, dan
 * begitu pula di layar Pega.
 */
function GridTotal({
  judul,
  baris,
  judulNilai = DETAIL_SHARE.value,
}: {
  judul: string
  baris: readonly NilaiTotalProp[] | readonly Simpul[]
  judulNilai?: string
}) {
  return (
    <div className="table-wrap trin__share-total">
      <table className="trin__tabel">
        <thead>
          <tr>
            <th scope="col">{judul}</th>
            <th scope="col">{judulNilai}</th>
          </tr>
        </thead>
        <tbody>
          {baris.length === 0 && (
            <tr>
              <td colSpan={2}>{TOTAL_SHARE.tanpaBaris}</td>
            </tr>
          )}
          {baris.map((b, i) => (
            <tr key={String(i)}>
              <td>{typeof b.Currency === 'string' ? b.Currency : ''}</td>
              <td className="trin__angka">{uang(typeof b.Value === 'string' ? b.Value : '')}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}

/** Tombol ▸/▾ rincian baris (`expandPane`), tertutup semula seperti Pega. */
function TombolRincian({ buka, label, onKlik }: { buka: boolean; label: string; onKlik: () => void }) {
  return (
    <button type="button" className="trin__buka" aria-expanded={buka} aria-label={label} onClick={onKlik}>
      {buka ? '▾' : '▸'}
    </button>
  )
}

/** Dropdown induk spreading — `BrowseTreatyArrangement_ParentReinsMasterTrt`. */
function useIndukSpreading(treatyGroupId: string, commencement: string, aktif: boolean): SusunanSpreading[] {
  const [daftar, setDaftar] = useState<SusunanSpreading[]>([])
  useEffect(() => {
    if (!aktif) return
    let dibuang = false
    ambilIndukSpreading(treatyGroupId, commencement)
      .then((d) => {
        if (!dibuang) setDaftar(d)
      })
      .catch(() => undefined)
    return () => {
      dibuang = true
    }
  }, [treatyGroupId, commencement, aktif])
  return daftar
}

/** Rincian `DetailShare` satu Treaty Group. */
function RincianDetailShare({
  d,
  terkunci,
  commencement,
  onUbah,
  onRumus,
}: {
  d: Simpul
  terkunci: boolean
  commencement: string
  /** Ganti Detail ini tanpa rumus. */
  onUbah: (baru: Simpul) => void
  /** Ganti Detail ini lalu jalankan rumus `aksi` atas pohon barunya. */
  onRumus: (aksi: MasukanShareProp['aksi'], baru: Simpul) => void
}) {
  const tipe = teksDari(d, 'SpreadingTypeID')
  // RD `BrowseTreatyArrangement_ParentReinsMasterTrt` — Treaty Group Detail
  // ini, `TreatyIn.Commencement`, kecuali `10246` (parameter dropdown).
  const induk = useIndukSpreading(teksDari(d, 'TreatyGroupID'), commencement, !terkunci)
  const sebar = larik(d, 'SpreadingList')
  return (
    <div className="trin__share-rincian">
      {/* `.RNMShare` — change → `TreatyInPropshareDetail` (saat ditinggalkan). */}
      <div
        onBlur={() => {
          if (!terkunci) onRumus('detail', d)
        }}
      >
        <FieldAngka
          label={DETAIL_SHARE.persenRnmShare}
          value={teksDari(d, 'RNMShare')}
          desimal={2}
          readOnly={terkunci}
          onChange={(v) => {
            onUbah({ ...d, RNMShare: v })
          }}
        />
      </div>
      <GridTotal judul={DETAIL_SHARE.rnmShare} baris={larik(d, 'RNMShareList')} judulNilai="" />

      <h5 className="trin__subjudul">{DETAIL_SHARE.spreading}</h5>
      {terkunci ? (
        <p className="trin__teks-sel">
          {DETAIL_SHARE.jenisSpreading}: {teksDari(d, 'SpreadingType') || tipe}
        </p>
      ) : (
        <Pilih
          label={DETAIL_SHARE.jenisSpreading}
          value={tipe}
          onChange={(v) => {
            // change → `FetchQSfromMaster(ParentReinsTypeID = .SpreadingTypeID)`:
            // Reins Type · Pct dari `PROPORTIONALARRG`, Value OR/R/I dari RNM Share.
            onRumus('spreading', { ...d, SpreadingTypeID: v, ...(v === '' ? { SpreadingType: '' } : {}) })
          }}
          opsi={[
            ...induk.map((x) => ({ value: x.reinsTypeId, label: x.reinsTypeName })),
            ...(tipe === '' || induk.some((x) => x.reinsTypeId === tipe) ? [] : [{ value: tipe, label: teksDari(d, 'SpreadingType') || tipe }]),
          ]}
        />
      )}
      {/* Reins Type · Pct — BACA-SAJA, hasil RD anak (`PROPORTIONALARRG`). */}
      <div className="table-wrap">
        <table className="trin__tabel">
          <thead>
            <tr>
              {DETAIL_SHARE.kolomSpreading.map((k) => (
                <th key={k} scope="col">
                  {k}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {sebar.length === 0 && (
              <tr>
                <td colSpan={2}>{TOTAL_SHARE.tanpaBaris}</td>
              </tr>
            )}
            {sebar.map((r, k) => (
              <tr key={k}>
                <td>{teksDari(r, 'ReinsTypeName')}</td>
                <td>{teksDari(r, 'Pct')}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
      <p className="trin__teks-sel">
        {DETAIL_SHARE.totalSpreadingPct} {selAngka(['persen', 2], teksDari(d, 'SpreadingTotalPct'))}
      </p>
      <GridTotal judul={DETAIL_SHARE.sebaranOR} baris={larik(d, 'RNMSpreadedList')} judulNilai="" />
      <GridTotal judul={DETAIL_SHARE.sebaranRI} baris={larik(d, 'RNMSpreadedListRI')} judulNilai="" />
    </div>
  )
}

export default function TabShareProp({
  pohon = [],
  petunjukKosong,
  mode = 'lihat',
  commencement = '',
  edmJenisMaterial = '',
}: {
  /** `TreatyIn.Limits` kontrak yang dimuat — nilai awal penampung bila tab Limits belum dibuka. */
  pohon?: readonly SimpulLimit[]
  petunjukKosong: string
  mode?: ModeForm
  /** `TreatyIn.Commencement` kepala (`YYYYMMDD`) — saringan RD spreading. */
  commencement?: string
  edmJenisMaterial?: string
}) {
  // ⭐ Penampung halaman — `Limits` milik bersama dengan tab Limits.
  const [limits, setLimits] = useProperti<SimpulLimit[]>('Limits', () => [...pohon])
  const [rnmShare, setRnmShare] = useProperti('RNMShareP', '')
  const [brokerage, setBrokerage] = useProperti('BrokeragePercentP', '')
  const [opsi, setOpsi] = useProperti('OptionLimit', '1')
  const [lintasBoard] = useProperti('RNMShareAcrossTheBoard', '')
  const [fakultatif] = useProperti('FacultativeShare', '')
  const [dipotong] = useProperti('RnmShareDeducted', '')
  const [totShare, setTotShare] = useProperti<NilaiTotalProp[]>('TotalShareRnmProp', [])
  const [totOR, setTotOR] = useProperti<NilaiTotalProp[]>('TotalSpreadedRnmProp', [])
  const [totRI, setTotRI] = useProperti<NilaiTotalProp[]>('TotalSpreadedRnmRIProp', [])
  const [sub, setSub] = useState<string>(SUB_TAB_SHARE[0])
  const [pesan, setPesan] = useState<string[]>([])
  const [sibuk, setSibuk] = useState(false)
  const [bukaKind, setBukaKind] = useState<ReadonlySet<number>>(() => new Set())
  const [bukaDetail, setBukaDetail] = useState<ReadonlySet<string>>(() => new Set())
  const daftar: readonly string[] = SUB_TAB_SHARE
  const tampil = daftar.includes(sub) ? sub : SUB_TAB_SHARE[0]
  const bisaUbah = mode === 'ubah'
  const terkunci = !bisaUbah || edmJenisMaterial.trim() === '2'

  const jalankan = (aksi: MasukanShareProp['aksi'], ubahan: { Limits?: SimpulLimit[]; RNMShareP?: string; OptionLimit?: string } = {}, i = 0, j = 0) => {
    setSibuk(true)
    setPesan([])
    hitungShareProp({
      aksi,
      Limits: ubahan.Limits ?? limits,
      RNMShareP: ubahan.RNMShareP ?? rnmShare,
      BrokeragePercentP: brokerage,
      OptionLimit: ubahan.OptionLimit ?? opsi,
      RNMShareAcrossTheBoard: lintasBoard,
      Commencement: commencement,
      TotalShareRnmProp: totShare,
      TotalSpreadedRnmProp: totOR,
      TotalSpreadedRnmRIProp: totRI,
      indeksLimit: i,
      indeksDetail: j,
    })
      .then((h) => {
        setLimits(h.Limits)
        setTotShare(h.TotalShareRnmProp)
        setTotOR(h.TotalSpreadedRnmProp)
        setTotRI(h.TotalSpreadedRnmRIProp)
        setPesan(h.pesan)
      })
      .catch((e: unknown) => {
        setPesan([e instanceof Error ? e.message : String(e)])
      })
      .finally(() => {
        setSibuk(false)
      })
  }
  const balik = <T,>(s: ReadonlySet<T>, k: T): ReadonlySet<T> => {
    const y = new Set(s)
    if (y.has(k)) y.delete(k)
    else y.add(k)
    return y
  }
  const totalKosong = totShare.length + totOR.length + totRI.length === 0

  return (
    <>
      <Panel judul={TOTAL_SHARE.judul}>
        <div className="trin__panel-kepala">
          <span className="trin__redup">{petunjukKosong}</span>
          {/* `Refresh` — `TreatyInPropshare`, `TreatyIn.ViewState !='1'`. */}
          <button
            type="button"
            className="btn btn--ghost btn--sm"
            disabled={terkunci || sibuk}
            onClick={() => {
              jalankan('share')
            }}
          >
            {TOTAL_SHARE.segarkan}
          </button>
        </div>
        <div className="trin__kolom trin__share-medan">
          {/* ⭐ `pxNumber` 2 desimal TANPA simbol, `pyPlaceholder` `%` —
              `TreatyInTabsProportional`. Nilai tampil `1,28` (bukan `1,28%`);
              `%` hanya placeholder saat kosong, persis layar Pega.
              `% RNM Share` change → `TreatyInPropshare` (saat ditinggalkan). */}
          <div
            onBlur={() => {
              if (!terkunci) jalankan('share')
            }}
          >
            <FieldAngka label={TOTAL_SHARE.persenRnmShare} value={rnmShare} desimal={2} placeholder="%" readOnly={terkunci} onChange={setRnmShare} />
          </div>
          <FieldAngka label={TOTAL_SHARE.persenBrokerage} value={brokerage} desimal={2} placeholder="%" readOnly={terkunci} onChange={setBrokerage} />
          <Pilih
            label={TOTAL_SHARE.opsi}
            value={opsi}
            onChange={(v) => {
              setOpsi(v)
              // `Option` change → `TreatyInPropshare`.
              if (!terkunci) jalankan('share', { OptionLimit: v })
            }}
            opsi={OPSI_SHARE}
          />
        </div>
        {pesan.length > 0 && (
          <p className="trin__galat" role="alert">
            {pesan.join(' · ')}
          </p>
        )}
      </Panel>

      <StripTab tab={daftar} aktif={tampil} onPilih={setSub} />

      {tampil === 'RNM Share' && (
        <Panel judul={SUB_TAB_SHARE[0]}>
          {/* `TreatyIn.FacultativeShare >0` — `Share to RNM :` + `%`. */}
          {Number(fakultatif) > 0 && (
            <p className="trin__teks-sel">
              {DETAIL_SHARE.rnmShareDeducted} {selAngka(['persen', 2], dipotong)}
            </p>
          )}
          {/* ⭐ Grid `Kind of Treaty` = `TreatyIn.Limits` tab Limits. */}
          <div className="table-wrap trin__share-total">
            <table className="trin__tabel">
              <thead>
                <tr>
                  <th scope="col" aria-label={DETAIL_SHARE.rincian} />
                  <th scope="col">{KOLOM_KIND_OF_TREATY_SHARE[0]}</th>
                </tr>
              </thead>
              <tbody>
                {limits.length === 0 && (
                  <tr>
                    <td colSpan={2}>{TOTAL_SHARE.tanpaBaris}</td>
                  </tr>
                )}
                {limits.map((l, i) => (
                  <Fragment key={i}>
                    <tr>
                      <td className="trin__buka-sel">
                        <TombolRincian
                          buka={bukaKind.has(i)}
                          label={`${DETAIL_SHARE.rincian} ${String(i + 1)}`}
                          onKlik={() => {
                            setBukaKind((x) => balik(x, i))
                          }}
                        />
                      </td>
                      <td>{teksDari(l, 'TreatyType') || DETAIL_SHARE.kindBelumDipilih}</td>
                    </tr>
                    {bukaKind.has(i) && (
                      <tr className="trin__rincian">
                        <td colSpan={2}>
                          {/* Rincian `TotalLimits` — grid `.Detail`. */}
                          <table className="trin__tabel">
                            <thead>
                              <tr>
                                <th scope="col" aria-label={DETAIL_SHARE.rincian} />
                                <th scope="col">{KOLOM_TOTAL_LIMITS[0]}</th>
                                <th scope="col">{KOLOM_TOTAL_LIMITS[1]}</th>
                              </tr>
                            </thead>
                            <tbody>
                              {larik(l, 'Detail').length === 0 && (
                                <tr>
                                  <td colSpan={3}>{TOTAL_SHARE.tanpaBaris}</td>
                                </tr>
                              )}
                              {larik(l, 'Detail').map((d, j) => {
                                const kunci = `${String(i)}:${String(j)}`
                                return (
                                  <Fragment key={kunci}>
                                    <tr>
                                      <td className="trin__buka-sel">
                                        <TombolRincian
                                          buka={bukaDetail.has(kunci)}
                                          label={`${DETAIL_SHARE.rincian} ${String(i + 1)}.${String(j + 1)}`}
                                          onKlik={() => {
                                            setBukaDetail((x) => balik(x, kunci))
                                          }}
                                        />
                                      </td>
                                      <td>{teksDari(d, 'TreatyGroup')}</td>
                                      {/* `.RNMShare` (`pxNumber` 2 desimal, TANPA simbol — Section
                                          `TotalLimits`) + `.ShareNote`: `1,28 of 100% of 100%`. */}
                                      <td>
                                        {selAngka(['uang', 2], teksDari(d, 'RNMShare'))}
                                        {teksDari(d, 'ShareNote')}
                                      </td>
                                    </tr>
                                    {bukaDetail.has(kunci) && (
                                      <tr className="trin__rincian">
                                        <td colSpan={3}>
                                          <RincianDetailShare
                                            d={d}
                                            terkunci={terkunci}
                                            commencement={commencement}
                                            onUbah={(baru) => {
                                              setLimits((ls) => gantiDetail(ls, i, j, () => baru))
                                            }}
                                            onRumus={(aksi, baru) => {
                                              const ls = gantiDetail(limits, i, j, () => baru)
                                              setLimits(ls)
                                              jalankan(aksi, { Limits: ls }, i, j)
                                            }}
                                          />
                                        </td>
                                      </tr>
                                    )}
                                  </Fragment>
                                )
                              })}
                            </tbody>
                          </table>
                        </td>
                      </tr>
                    )}
                  </Fragment>
                ))}
              </tbody>
            </table>
          </div>

          <GridTotal judul={GRID_TOTAL_RNM_SHARE[0]} baris={totShare} />
          <GridTotal judul={GRID_TOTAL_RNM_SHARE[1]} baris={totOR} />
          <GridTotal judul={GRID_TOTAL_RNM_SHARE[2]} baris={totRI} />
          {/* ⚠️ Kosong sebelum Refresh: nilai tersimpannya belum punya tabel. */}
          {totalKosong && <Kosong pesan={TOTAL_SHARE.tanpaBaris} petunjuk={TOTAL_SHARE.petunjukTotal} />}
        </Panel>
      )}
    </>
  )
}
