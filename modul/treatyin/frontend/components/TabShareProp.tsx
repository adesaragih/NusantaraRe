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

import { useEffect, useState } from 'react'

import { FieldAngka } from '../../../../inti/frontend/components/ui/dasar'
import { PilihCari as Pilih } from '../../../../inti/frontend/components/ui/pilihSaring'
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
import { BlokPega, DeretTombolPega, GridNilaiPega, GridPega } from './gridPega'
import { StripTabNavigasi } from './navigasi'
import { PemicuUbah } from './pemicuUbah'
import { saringAngka } from './saringAngka'
import { SelKosongPega, TataPegaBlok } from './tataPega'
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

/**
 * `CountTotalPctSpead` — Σ `.Pct` baris `.SpreadingList` (langkah terakhir
 * `AddDelSpreadingTreatyin`). Ketikan berkoma desimal (`25,5`) ikut terbaca.
 */
export function jumlahPct(ls: readonly Simpul[]): string {
  let total = 0
  for (const r of ls) {
    const t = teksDari(r, 'Pct').trim()
    const n = Number(t.includes(',') ? t.replace(/\./g, '').replace(',', '.') : t)
    if (Number.isFinite(n)) total += n
  }
  return String(Math.round(total * 1e6) / 1e6)
}

/** Ganti satu Detail di pohon — pohon baru, yang lain utuh. */
function gantiDetail(ls: readonly Simpul[], i: number, j: number, f: (d: Simpul) => Simpul): Simpul[] {
  return ls.map((l, x) => (x !== i ? l : { ...l, Detail: larik(l, 'Detail').map((d, y) => (y !== j ? d : f(d))) }))
}

/** Uang dua desimal — gambar 17: `3.000.000.000,00`. */
const uang = (v: string) => selAngka(['uang', 2], v)

/**
 * Satu grid nilai Pega: kolom mata uang berjudul, kolom nilai (`192 · 347`).
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
  lebarTetap = true,
}: {
  judul: string
  baris: readonly NilaiTotalProp[] | readonly Simpul[]
  judulNilai?: string
  lebarTetap?: boolean
}) {
  return (
    <GridNilaiPega
      judul={judul}
      nilai={judulNilai}
      baris={baris.map((b) => ({
        Currency: typeof b.Currency === 'string' ? b.Currency : '',
        Value: typeof b.Value === 'string' ? b.Value : '',
      }))}
      lebar={[192, 347]}
      tampil={uang}
      lebarTetap={lebarTetap}
    />
  )
}

/**
 * Dropdown induk spreading — `BrowseTreatyArrangement_ParentReinsMasterTrt`.
 *
 * ⛔ PARAMETERNYA PERSIS `Section/DetailShare.xml`, dan daftar itu dibaca ULANG
 * 8 Oktober 2026 sesudah pembacaan pertama keliru:
 *
 *	TreatyYear · TreatyGroupID · TreatyDescID "10001" · StartDate ·
 *	ReinsTypeID "10246"
 *
 * ⚠️ Pembacaan pertama menyimpulkan `TreatyGroupID` dan `TreatyDescID` TIDAK
 * dikirim, dan keduanya sempat dibuang dari sini. Sebabnya jendela baca yang
 * terlalu sempit — ia memotong daftar parameter dan hanya menangkap dua yang
 * terakhir. Keduanya dikembalikan.
 *
 * ⭐ Dropdown yang menawarkan induk DI LUAR Treaty Group barisnya bukan
 * perbaikan: `FetchQSfromMaster` mencari induk itu dengan filter grup yang
 * sama, jadi pilihan di luar grup menghasilkan spreading KOSONG — menyesatkan,
 * bukan menolong.
 */
const DESC_SPREADING_PROP = '10001'

function useIndukSpreading(commencement: string, aktif: boolean): SusunanSpreading[] {
  const [daftar, setDaftar] = useState<SusunanSpreading[]>([])
  useEffect(() => {
    if (!aktif) return
    let dibuang = false
    ambilIndukSpreading('', DESC_SPREADING_PROP, commencement)
      .then((d) => {
        if (!dibuang) setDaftar(d)
      })
      .catch(() => undefined)
    return () => {
      dibuang = true
    }
  }, [commencement, aktif])
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
  // ⛔ RALAT 8 Oktober 2026 (sore) — DUA cabang, persis `DetailShare.xml`:
  //   `.SpreadingTypeID != ''` → dropdown Spreading Type + grid BACA (L12);
  //   `.SpreadingTypeID == ''` → grid MANUAL (L16) yang dapat ditambah.
  // Pemilik proses: *"spread nya bisa ditambah, cek lagi xml nya"*.
  const manual = tipe === ''
  // RD `BrowseTreatyArrangement_ParentReinsMasterTrt` — Treaty Group Detail
  // ini, `TreatyIn.Commencement`, kecuali `10246` (parameter dropdown).
  const induk = useIndukSpreading(commencement, !terkunci && !manual)
  // Dropdown `Reins Type` grid manual — RD yang sama dengan `TreatyGroupID =
  // TempSprd.TreatyGroupID`, halaman yang tak pernah diisi: filter grup
  // DILEWATI (sama dengan `setSpreadName` backend).
  const indukManual = useIndukSpreading(commencement, !terkunci && manual)
  const sebar = larik(d, 'SpreadingList')
  /** `CountTotalPctSpead` — Σ `.Pct` (langkah terakhir `AddDelSpreadingTreatyin`). */
  const denganTotal = (ls: Simpul[]): Simpul => ({ ...d, SpreadingList: ls, SpreadingTotalPct: jumlahPct(ls) })
  const gantiSebar = (i: number, isi: Simpul): Simpul => denganTotal(sebar.map((r, x) => (x === i ? { ...r, ...isi } : r)))
  // ⭐ Bentuk Pega (tangkapan layar pemakai 8 Oktober 2026): `% RNM Share`
  // berlabel di kiri, grid `RNM Share` sempit, judul `Spreading`, `Spreading
  // Type`, grid `Reins Type · Pct` (200 · 192), `Total Spreading Pct :`,
  // lalu `Value Spreading OR` dan `Value Spreading R/I`.
  // ⭐ DUA SPREADING, PERSIS PEGA — keputusan pemakai 9 Oktober 2026 (*"di
  // pega ada terdapat 2 spreading … yang tidak bisa ditambah itu khusus untuk
  // spreading lama"*), menggantikan "satu grid untuk kedua cabang" 8 Oktober.
  //
  //   `.SpreadingTypeID != ''` → spreading LAMA: dropdown Spreading Type +
  //       grid BACA (L12, nol tombol) — isinya `FetchQSfromMaster`
  //       (`PROPORTIONALARRG`), tidak dapat ditambah/diubah.
  //   `.SpreadingTypeID == ''` → spreading BARU: grid MANUAL (L16)
  //       ber-`AddDelSpreadingTreatyin`, dropdown Spreading Type TIDAK ada.
  //
  // ⚠️ Activity pembentuk Share (`TreatyInPropshare`/`…Detail`) TIDAK PERNAH
  // mengisi Spreading Type — ia hanya menjalankan `FetchQSfromMaster` bila
  // sudah terisi, `SetSpreadName` bila kosong. Jadi spreading lama = kontrak
  // warisan yang menyimpan Spreading Type; kontrak baru selalu manual.
  const gridSebar = (lama: boolean) => (
    <GridPega
      label={DETAIL_SHARE.spreading}
      lebarTetap
      kosong="No items"
      kolom={[
        {
          judul: (lama ? DETAIL_SHARE.kolomSpreading : DETAIL_SHARE.kolomSpreadingManual)[0],
          lebar: 260,
          isi: (r: Simpul, i: number) =>
            terkunci || lama ? (
              teksDari(r, 'ReinsTypeName')
            ) : (
              <select
                className="field__input"
                aria-label={`${DETAIL_SHARE.kolomSpreadingManual[0]} ${String(i + 1)}`}
                value={teksDari(r, 'ReinsTypeID')}
                onChange={(e) => {
                  const id = e.target.value
                  const nama = indukManual.find((x) => x.reinsTypeId === id)?.reinsTypeName ?? ''
                  onRumus('sebar-nama', gantiSebar(i, { ReinsTypeID: id, ReinsTypeName: nama }))
                }}
              >
                <option value="">{DETAIL_SHARE.pilihReins}</option>
                {indukManual.map((x) => (
                  <option key={x.reinsTypeId} value={x.reinsTypeId}>
                    {x.reinsTypeName}
                  </option>
                ))}
                {teksDari(r, 'ReinsTypeID') !== '' && !indukManual.some((x) => x.reinsTypeId === teksDari(r, 'ReinsTypeID')) && (
                  <option value={teksDari(r, 'ReinsTypeID')}>{teksDari(r, 'ReinsTypeName') || teksDari(r, 'ReinsTypeID')}</option>
                )}
              </select>
            ),
        },
        {
          judul: (lama ? DETAIL_SHARE.kolomSpreading : DETAIL_SHARE.kolomSpreadingManual)[1],
          lebar: 140,
          angka: true,
          isi: (r: Simpul, i: number) =>
            terkunci || lama ? (
              selAngka(['uang', 2], teksDari(r, 'Pct'))
            ) : (
              <PemicuUbah
                nilai={teksDari(r, 'Pct')}
                aksi={() => {
                  onRumus('sebar-nama', d)
                }}
              >
                <input
                  className="field__input"
                  type="text"
                  inputMode="decimal"
                  aria-label={`${DETAIL_SHARE.kolomSpreadingManual[1]} ${String(i + 1)}`}
                  value={teksDari(r, 'Pct')}
                  onChange={(e) => {
                    onUbah(gantiSebar(i, { Pct: saringAngka(e.target.value) }))
                  }}
                />
              </PemicuUbah>
            ),
        },
      ]}
      baris={sebar}
      rincian={lama ? undefined : (r: Simpul) => (
        <GridPega
          label={DETAIL_SHARE.kolomPecahan[0]}
          lebarTetap
          kolom={[
            { judul: DETAIL_SHARE.kolomPecahan[0], lebar: 140, isi: (p: Simpul) => teksDari(p, 'ReinsName') },
            { judul: DETAIL_SHARE.kolomPecahan[1], lebar: 100, isi: (p: Simpul) => teksDari(p, 'Currency') },
            { judul: DETAIL_SHARE.kolomPecahan[2], lebar: 100, angka: true, isi: (p: Simpul) => selAngka(['uang', 2], teksDari(p, 'SharePct')) },
            { judul: DETAIL_SHARE.kolomPecahan[3], lebar: 180, angka: true, isi: (p: Simpul) => uang(teksDari(p, 'Amount')) },
          ]}
          baris={larik(r, 'BreakDownSprdList')}
        />
      )}
      tombol={
        terkunci || lama
          ? undefined
          : {
              lebar: 90,
              tambah: {
                label: DETAIL_SHARE.tambahSpread,
                onKlik: () => {
                  onUbah(denganTotal([...sebar, { Pct: '0' }]))
                },
              },
              hapus: {
                label: DETAIL_SHARE.hapusSpread,
                onKlik: (i: number) => {
                  onUbah(denganTotal(sebar.filter((_, x) => x !== i)))
                },
              },
            }
      }
    />
  )

  return (
    <>
      <div className="trin__kolom">
        {/* `.RNMShare` — change → `TreatyInPropshareDetail` (saat ditinggalkan). */}
        <PemicuUbah
          nilai={teksDari(d, 'RNMShare')}
          aktif={!terkunci}
          aksi={() => {
            onRumus('detail', d)
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
        </PemicuUbah>
      </div>
      <GridTotal judul={DETAIL_SHARE.rnmShare} baris={larik(d, 'RNMShareList')} judulNilai="" />

      <h5 className="trin__subjudul">{DETAIL_SHARE.spreading}</h5>
      {manual ? (
        <>
          {/* L16 — grid MANUAL. Reins Type / Pct Share: change →
              `SetSpreadName` (aksi `sebar-nama`); Add / Delete →
              `AddDelSpreadingTreatyin` (tambah `Pct = 0` / hapus baris, lalu
              `CountTotalPctSpead`). Tiap baris dapat dibuka ▸ ke panel
              `SpreadingTPDtl` — pecahan `.BreakDownSprdList`. */}
          {gridSebar(false)}
          <p className="trin__teks-sel">
            {DETAIL_SHARE.totalSharePct} {selAngka(['persen', 2], teksDari(d, 'SpreadingTotalPct'))}
          </p>
          <GridTotal judul={DETAIL_SHARE.sebaranOR} baris={larik(d, 'RNMSpreadedList')} judulNilai="" />
          <GridTotal judul={DETAIL_SHARE.sebaranRI} baris={larik(d, 'RNMSpreadedListRI')} judulNilai="" />
        </>
      ) : (
      <>
      <div className="trin__kolom">
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
      </div>
      {/* L12 — spreading LAMA: grid BACA, nol tombol (`gridSebar`). */}
      {gridSebar(true)}
      <p className="trin__teks-sel">
        {DETAIL_SHARE.totalSpreadingPct} {selAngka(['persen', 2], teksDari(d, 'SpreadingTotalPct'))}
      </p>
      <GridTotal judul={DETAIL_SHARE.sebaranOR} baris={larik(d, 'RNMSpreadedList')} judulNilai="" />
      <GridTotal judul={DETAIL_SHARE.sebaranRI} baris={larik(d, 'RNMSpreadedListRI')} judulNilai="" />
      </>
      )}
    </>
  )
}

export default function TabShareProp({
  pohon = [],
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
  const [dipotong] = useProperti('RnmShareDeducted', '')
  const [totShare, setTotShare] = useProperti<NilaiTotalProp[]>('TotalShareRnmProp', [])
  const [totOR, setTotOR] = useProperti<NilaiTotalProp[]>('TotalSpreadedRnmProp', [])
  const [totRI, setTotRI] = useProperti<NilaiTotalProp[]>('TotalSpreadedRnmRIProp', [])
  const [sub, setSub] = useState<string>(SUB_TAB_SHARE[0])
  const [pesan, setPesan] = useState<string[]>([])
  const [sibuk, setSibuk] = useState(false)
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
  // ⭐ `Share to RNM :` — tangkapan layar Pega pemakai 8 Oktober 2026
  // memperlihatkannya pada kontrak TANPA retro (`Share to Other Retro`
  // tidak tampil di blok Total Share) dengan nilai = `% RNM Share` (20,00).
  // Jadi baris ini SELALU tampil; nilainya `RnmShareDeducted` tersimpan,
  // atau `% RNM Share` bila tidak ada potongan retro. Nol aritmetika.
  const shareKeRnm = dipotong.trim() !== '' ? dipotong : rnmShare

  // ⭐ Bentuk Pega (8 Oktober 2026), urut ekspor `TreatyInTabsProportional`
  // tab `Share` + `TreatyInShareProp`: blok `Total Share` (Refresh di atas,
  // medan berlabel kiri), strip `RNM Share`, `Share to RNM :`, grid
  // `masterDetail` Kind of Treaty → grid Treaty Group · % RNM Share →
  // `DetailShare`, lalu tiga grid total BERTUMPUK.
  return (
    <div className="trin__blok trin__tab">
      <BlokPega judul={TOTAL_SHARE.judul}>
        {/* ⭐ TATA LETAK PEGA (8 Oktober 2026) — `Section/TreatyInTabsProportional.xml`
            panel `Total Share` L22230: `Inline grid double` L22729 —
            separuh KIRI `Stacked with labels left` L23028 [Refresh, lalu
            `Inline grid double` L23622: (`Stacked with labels left` L23920:
            % RNM Share · % Brokerage · Option) | kotak centang `1=2`],
            separuh KANAN `Stacked with labels left` L25916 (medan retro
            `IsMultipleRetro` — belum dibangun, §17). Jadi ketiga medan
            selebar ±seperempat panel, persis gambar Pega 16. */}
        <TataPegaBlok tata="g2">
        <TataPegaBlok tata="kiri">
        <DeretTombolPega>
          {/* `Refresh` — `TreatyInPropshare`, `TreatyIn.ViewState !='1'`. */}
          <button
            type="button"
            className="btn btn--sm"
            disabled={terkunci || sibuk}
            onClick={() => {
              jalankan('share')
            }}
          >
            {TOTAL_SHARE.segarkan}
          </button>
        </DeretTombolPega>
        <TataPegaBlok tata="g2">
        <div className="trin__kolom trin__share-medan">
          {/* ⭐ `pxNumber` 2 desimal TANPA simbol, `pyPlaceholder` `%` —
              `TreatyInTabsProportional`. Nilai tampil `1,28` (bukan `1,28%`);
              `%` hanya placeholder saat kosong, persis layar Pega.
              `% RNM Share` change → `TreatyInPropshare` (saat ditinggalkan). */}
          <PemicuUbah
            nilai={rnmShare}
            aktif={!terkunci}
            aksi={() => {
              jalankan('share')
            }}
          >
            <FieldAngka label={TOTAL_SHARE.persenRnmShare} value={rnmShare} desimal={2} placeholder="%" readOnly={terkunci} onChange={setRnmShare} />
          </PemicuUbah>
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
        <SelKosongPega />
        </TataPegaBlok>
        </TataPegaBlok>
        <SelKosongPega />
        </TataPegaBlok>
        {pesan.length > 0 && (
          <p className="trin__galat" role="alert">
            {pesan.join(' · ')}
          </p>
        )}
      </BlokPega>

      <StripTabNavigasi tab={daftar} aktif={tampil} onPilih={setSub} />

      {tampil === 'RNM Share' && (
        <BlokPega>
          <p className="trin__share-rnm-prop">
            <strong>{DETAIL_SHARE.rnmShareDeducted}</strong> {selAngka(['uang', 2], shareKeRnm)} <span>%</span>
          </p>
          {/* ⭐ Grid `Kind of Treaty` = `TreatyIn.Limits` tab Limits. */}
          <GridPega
            label={KOLOM_KIND_OF_TREATY_SHARE[0]}
            labelBuka={DETAIL_SHARE.rincian}
            kolom={[
              {
                judul: KOLOM_KIND_OF_TREATY_SHARE[0],
                lebar: 800,
                isi: (l: Simpul) => teksDari(l, 'TreatyType') || DETAIL_SHARE.kindBelumDipilih,
              },
            ]}
            baris={limits}
            rincian={(l, i) => (
              // Rincian `TotalLimits` — grid `.Detail` Treaty Group · % RNM Share.
              <GridPega
                label={KOLOM_TOTAL_LIMITS[0]}
                labelBuka={DETAIL_SHARE.rincian}
                lebarTetap
                kolom={[
                  { judul: KOLOM_TOTAL_LIMITS[0], lebar: 399, isi: (d: Simpul) => teksDari(d, 'TreatyGroup') },
                  {
                    judul: KOLOM_TOTAL_LIMITS[1],
                    lebar: 216,
                    // `.RNMShare` (`pxNumber` 2 desimal, TANPA simbol — Section
                    // `TotalLimits`) + `.ShareNote`: `20,00 of 100% of 100%`.
                    isi: (d: Simpul) => `${selAngka(['uang', 2], teksDari(d, 'RNMShare'))} ${teksDari(d, 'ShareNote')}`.trim(),
                  },
                ]}
                baris={larik(l, 'Detail')}
                rincian={(d, j) => (
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
                )}
              />
            )}
          />

          {/* Tiga grid total BERTUMPUK, selebar ±separuh layar — tangkapan Pega.
              ⭐ Sebabnya kini terbaca: `Section/TreatyInShareProp.xml`
              `Inline grid double` L4158 berisi grid · tombol `3=4` · grid ·
              tombol `3=4` · grid — sel tersembunyi tetap memesan slot, jadi
              tiap grid duduk di separuh KIRI barisnya sendiri. */}
          <div className="trin__share-prop-total">
            <GridTotal judul={GRID_TOTAL_RNM_SHARE[0]} baris={totShare} lebarTetap={false} />
            <GridTotal judul={GRID_TOTAL_RNM_SHARE[1]} baris={totOR} lebarTetap={false} />
            <GridTotal judul={GRID_TOTAL_RNM_SHARE[2]} baris={totRI} lebarTetap={false} />
          </div>
        </BlokPega>
      )}
    </div>
  )
}
