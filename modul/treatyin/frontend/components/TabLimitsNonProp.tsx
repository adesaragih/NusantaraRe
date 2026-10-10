// Tab **Limits** cabang NON-PROPORSIONAL — dibaca dari ekspor
// (`labelsLimitsNP.ts`): grid layer → rincian `Layers` (Treaty Group →
// `CoBList`, limit dua mata uang, Reinstatement, Premium Earned, MDP, ROL)
// → Summary of Limit → Total All Layers.
//
// ⭐ Rumusnya hidup di services (`hitung_limit_np.go`), disalin dari
// Activity dan diukur terhadap data Pega. Layar ini hanya mengirim isian
// dan menggabungkan medan yang Activity TULIS — nol rumus di sini.
//
// ⭐ Tampilan (permintaan pemakai 6 Oktober 2026, "UI/UX hebat"): tiap
// layer satu KARTU bernomor dengan ringkasan limit/deductible/ROL di
// kepalanya; rinciannya dikelompokkan per bagian, dan limit dua mata uang
// disusun sebagai matriks (baris = 100 % Limit · Agregate · Deductible,
// kolom = mata uang 1 · mata uang 2) — tiga baris Pega yang mengikat SATU
// properti mata uang kini memilihnya SEKALI di kepala kolom.
//
// ⛔ Hasil suntingan dan hitungan hidup di salinan pohon layar ini; jalur
// Save menunggu keputusan pemilik proses.

import { Fragment, useEffect, useRef, useState } from 'react'

import { PemicuUbah, usePemicuUbah } from './pemicuUbah'

import { Field, FieldAngka } from '../../../../inti/frontend/components/ui/dasar'
import { PilihCari as Pilih } from '../../../../inti/frontend/components/ui/pilihSaring'
import {
  ambilKelasBisnis,
  ambilOpsiLimits,
  hitungLimitNP,
  type AksiLimitNP,
  type EgnpiLimitNP,
  type KursLimitNP,
  type LimitsAkar,
  type OpsiLimits,
  type OpsiPilihan,
  type PilihanWarisan,
  type SimpulLimit,
} from '../api'
import { GRID_TOTAL, KOLOM_LAYER, KOLOM_REINSTATEMENT, KOLOM_RINGKASAN, LIMITS_NP, kontrakRevisi, type KolomNP } from '../labelsLimitsNP'
import type { ModeForm } from '../mode'
import { DropdownDaftar } from './IsianAuto'
import { BlokPega, GridPega, TeksPega, type KolomPega } from './gridPega'
import { saringAngka } from './saringAngka'
import { SelKosongPega, TataPegaBlok } from './tataPega'
import { formatLimit, teksDari } from './TabLimitsProp'

function larikDari(s: SimpulLimit, kunci: string): SimpulLimit[] {
  const v = s[kunci]
  return Array.isArray(v) ? v : []
}

function ganti<T>(larik: readonly T[], i: number, baru: T): T[] {
  return larik.map((x, j) => (j === i ? baru : x))
}

/** Angka tampil (mode lihat) — pemformat modul, bukan pemformat kedua. */
const tampil = (nilai: string, desimal: number | null) => formatLimit('uang', desimal, nilai)

/**
 * Medan yang tiap aksi TULIS — hanya ini yang digabung kembali ke layer.
 * Menggabung seluruh jawaban akan menimpa suntingan yang tidak dikirim.
 */
export const DITULIS_LIMIT_NP: Readonly<Record<AksiLimitNP, readonly string[]>> = {
  egnpi: ['EgnpiTotalList'],
  reinstatement: ['Reinstatement_List'],
  'reinst-jumlah': ['Reinstatement_List'],
  'reinst-persen': ['Reinstatement_List'],
  'reinst-tambahan': ['Reinstatement_List'],
  adj: ['PremiumEarnedList', 'ROLPct'],
  // ⭐ `Reinstatement_List` — sesudah MDP berubah, Reinstatement Premium
  // Amount disegarkan dengan rumus `ReCalculateReinstatement` yang SAMA
  // (keputusan pemakai 8 Oktober 2026; `segarkanReinstatement`).
  mdp: ['MDPList', 'MDPMinList', 'ROLPct', 'Reinstatement_List'],
  total: [],
  'nilai-list': ['EgnpiTotalList', 'PremiumEarnedList', 'MDPList', 'MDPMinList', 'ROLPct', 'Reinstatement_List'],
}

/** Aksi yang menulis SEMUA layer; selainnya hanya layer `indeks`. */
const SEMUA_LAYER: readonly AksiLimitNP[] = ['egnpi', 'nilai-list']

/** Dropdown `associated` — nilai tersimpan apa adanya bila di luar domain. */
function PilihNP({
  label,
  nilai,
  opsi,
  kosong,
  bisaUbah,
  onUbah,
}: {
  label: string
  nilai: string
  opsi: readonly OpsiPilihan[]
  kosong?: string
  bisaUbah: boolean
  onUbah: (v: string) => void
}) {
  if (!bisaUbah) {
    const t = opsi.find((o) => o.value === nilai)?.label ?? nilai
    return <Field label={label} value={t} readOnly onChange={() => undefined} />
  }
  return <Pilih label={label} value={nilai} kosong={kosong ?? LIMITS_NP.pilihKosong} opsi={[...opsi]} onChange={onUbah} />
}

/** Medan angka/teks satu layer; `onLepas` = peristiwa `change` Pega. */
function MedanNP({
  label,
  nilai,
  desimal,
  bisaUbah,
  onUbah,
  onLepas,
}: {
  label: string
  nilai: string
  desimal: number | null
  bisaUbah: boolean
  onUbah: (v: string) => void
  onLepas?: () => void
}) {
  // ⭐ Peristiwa `change` Pega — Activity berjalan HANYA bila nilainya
  // BERUBAH sejak medan dimasuki (8 Oktober 2026). Sebelumnya setiap
  // meninggalkan medan menjalankannya: melewati `Reinstatement` saja
  // membangun ulang grid Reinstatement dan menghapus persen yang sudah
  // disunting; melewati `Adjustment Rate`/`MDP %` menimpa Premium Earned/MDP
  // yang diketik tangan.
  // Peristiwa `change` Pega — `pemicuUbah.tsx`.
  const pemicu = usePemicuUbah(nilai, bisaUbah ? onLepas : undefined)
  return (
    <div className="trin__limit-medan" onFocus={pemicu.masuk} onBlur={pemicu.keluar}>
      {/* ⭐ PEMISAH RIBUAN SAAT MENGETIK — permintaan pemilik proses
          7 Oktober 2026. `FieldAngka` menerima dan mengembalikan bentuk
          KABEL (titik desimal), jadi nol pemanggil berubah selain nama
          komponennya.

          ⛔ RALAT 7 Oktober 2026, pukul berikutnya. Bentuk pertama menulis
          `desimal ?? 2` — memaksa SETIAP medan ber-`desimal === null`
          menjadi uang berdua desimal. Akibatnya `Layer` dan `Part` di kartu
          Layers tampil `1,00`, dan pemilik proses melaporkannya: *"pada tab
          limits di layer nya tidak perlu format nya krn itu memang angka
          biasa"*.

          ⚠️ `desimal === null` memang BUKAN "pakai bawaan". Ia berarti
          presisi kolomnya TIDAK TERBACA di gambar mana pun — dan medan
          seperti `Layer`/`Part` bukan uang sama sekali, melainkan nomor
          urut. Yang presisinya tidak diketahui LEWAT APA ADANYA, persis
          seperti sebelum pemisah ribuan ada. */}
      {desimal === null ? (
        // Seluruh MedanNP kontrol Number (Layer, Part, AdjRate, MDP%, …) —
        // huruf/simbol ditolak (permintaan pemakai 8 Oktober 2026).
        <Field
          label={label}
          value={bisaUbah ? nilai : tampil(nilai, desimal)}
          readOnly={!bisaUbah}
          onChange={(v) => {
            onUbah(saringAngka(v))
          }}
        />
      ) : (
        <FieldAngka label={label} value={nilai} desimal={desimal} readOnly={!bisaUbah} onChange={onUbah} />
      )}
    </div>
  )
}

/** Kotak centang bernilai teks `"true"`/`"false"`, seperti dokumen. */
function Centang({ label, nilai, bisaUbah, onUbah }: { label: string; nilai: string; bisaUbah: boolean; onUbah: (v: string) => void }) {
  return (
    <label className="trin__limit-medan">
      <input
        type="checkbox"
        checked={nilai === 'true'}
        disabled={!bisaUbah}
        onChange={(e) => {
          onUbah(e.target.checked ? 'true' : 'false')
        }}
      />{' '}
      {label}
    </label>
  )
}

/**
 * Grid `Currency · Value` — Premium Earned, Min Premium Amt, MDP, Egnpi this
 * layer, total. Bentuk Pega (8 Oktober 2026): kepala `judul · Value`
 * (`judulNilai`), `Add MDP` di sel kepala kolom tombol, `Delete` per baris,
 * selebar `pyWidth` ekspor.
 */
function GridNilai({
  judul,
  judulNilai = LIMITS_NP.nilai,
  lebar = [198, 363, 100],
  baris,
  bisaUbahMataUang,
  bisaUbahNilai,
  tambah,
  mataUang,
  onUbah,
}: {
  judul: string
  judulNilai?: string
  /** `pyWidth` ekspor: mata uang · nilai · kolom tombol. */
  lebar?: readonly [number, number, number?]
  baris: readonly SimpulLimit[]
  bisaUbahMataUang: boolean
  bisaUbahNilai: boolean
  tambah: boolean
  mataUang: readonly PilihanWarisan[]
  onUbah: (b: SimpulLimit[]) => void
}) {
  return (
    <GridPega
      label={judul}
      kelas="trin__tabel--nilai"
      lebarTetap
      kolom={[
        {
          judul,
          lebar: lebar[0],
          isi: (b, r) =>
            bisaUbahMataUang ? (
              <DropdownDaftar
                label=""
                nilai={teksDari(b, 'Currency')}
                pilihan={mataUang}
                bisaUbah
                onPilih={(nama, id) => {
                  onUbah(ganti(baris, r, { ...b, Currency: nama, CurrencyID: id }))
                }}
              />
            ) : (
              teksDari(b, 'Currency')
            ),
        },
        {
          judul: judulNilai,
          lebar: lebar[1],
          angka: true,
          isi: (b, r) =>
            bisaUbahNilai ? (
              /* ⭐ Sel angka berpemisah ribuan — sama dengan grid
                 `100% Limit` / `Retention` / `Cession to R/I` cabang Prop. */
              <FieldAngka
                label=""
                value={teksDari(b, 'Value')}
                desimal={2}
                onChange={(v) => {
                  onUbah(ganti(baris, r, { ...b, Value: v }))
                }}
              />
            ) : (
              tampil(teksDari(b, 'Value'), 2)
            ),
        },
      ]}
      baris={baris}
      tombol={
        tambah
          ? {
              lebar: lebar[2] ?? 100,
              tambah: {
                label: LIMITS_NP.tambahMDP,
                onKlik: () => {
                  onUbah([...baris, { Currency: '', CurrencyID: '', Value: '' }])
                },
              },
              hapus: {
                label: LIMITS_NP.hapusMDP,
                akses: (_, r) => `${LIMITS_NP.hapusMDP} ${judul} ${String(r + 1)}`,
                onKlik: (r) => {
                  onUbah(baris.filter((_, j) => j !== r))
                },
              },
            }
          : undefined
      }
    />
  )
}

/** Grid baca-saja berkolom ekspor (`Summary of Limit`). */
function GridBaca({ kolom, lebar, baris }: { kolom: readonly KolomNP[]; lebar: readonly number[]; baris: readonly SimpulLimit[] }) {
  return (
    <GridPega
      kolom={kolom.map((k, c) => ({
        judul: k.label,
        lebar: lebar[c] ?? 150,
        angka: k.desimal !== null,
        isi: (b: SimpulLimit) => (k.desimal === null ? teksDari(b, k.kunci) : tampil(teksDari(b, k.kunci), k.desimal)),
      }))}
      baris={baris}
    />
  )
}

/**
 * Rincian satu Treaty Group layer — `Section/CoBList.xml`: Treaty Group lalu
 * grid `Class of Business`. Baris gridnya sendiri di `GridGrup`.
 */
function RincianGrup({
  g,
  opsi,
  bisaUbah,
  onUbah,
  onGantiGrup,
}: {
  g: SimpulLimit
  opsi: OpsiLimits
  bisaUbah: boolean
  onUbah: (g: SimpulLimit) => void
  onGantiGrup: (g: SimpulLimit) => void
}) {
  const [kelas, setKelas] = useState<PilihanWarisan[]>([])
  const idGrup = teksDari(g, 'TreatyGroupID')
  useEffect(() => {
    if (!bisaUbah || idGrup === '') return
    let dibuang = false
    ambilKelasBisnis(idGrup)
      .then((d) => {
        if (!dibuang) setKelas(d)
      })
      .catch(() => undefined)
    return () => {
      dibuang = true
    }
  }, [bisaUbah, idGrup])
  const pilihGrup = (nama: string, id: string) => {
    onGantiGrup({ ...g, TreatyGroup: nama, TreatyGroupID: id })
  }
  const daftarCoB = larikDari(g, 'ClassOfBusinessList')
  return (
    <>
      <DropdownDaftar label={LIMITS_NP.treatyGroup} nilai={teksDari(g, 'TreatyGroup')} pilihan={opsi.kelompokTreaty} bisaUbah={bisaUbah} onPilih={pilihGrup} />
      {/* Add CoB / Delete MATI di ekspor (`ViewState !='1' && 1=2`). */}
      <GridPega
        label={LIMITS_NP.kelasBisnis}
        kolom={[
          {
            judul: LIMITS_NP.kelasBisnis,
            lebar: 722,
            isi: (c, r) => (
              <DropdownDaftar
                label=""
                nilai={teksDari(c, 'ClassOfBusiness')}
                pilihan={kelas}
                bisaUbah={bisaUbah}
                onPilih={(nama, id) => {
                  // DT `SetCoB` — nama dan kode Class of Business.
                  onUbah({ ...g, ClassOfBusinessList: ganti(daftarCoB, r, { ...c, ClassOfBusiness: nama, ClassOfBusinessID: id }) })
                }}
              />
            ),
          },
        ]}
        baris={daftarCoB}
      />
    </>
  )
}

/**
 * Grid Treaty Group layer — `masterDetail`, rincian `CoBList`. Bentuk Pega:
 * `Treaty Group` · kolom tombol (`Add Treaty Group` di kepala, `Delete` per
 * baris) · `ROL Profile` — kolom tombol DI TENGAH, seperti ekspor (lebar 202 ·
 * 147 · 104).
 */
function GridGrup({
  grup,
  opsi,
  bisaUbah,
  onTambah,
  onUbahBaris,
  onHapus,
  onGantiBaris,
}: {
  grup: readonly SimpulLimit[]
  opsi: OpsiLimits
  bisaUbah: boolean
  onTambah: () => void
  onUbahBaris: (r: number, g: SimpulLimit) => void
  onHapus: (r: number) => void
  onGantiBaris: (r: number, g: SimpulLimit) => void
}) {
  return (
    <GridPega
      label={LIMITS_NP.treatyGroup}
      kolom={[
        { judul: LIMITS_NP.treatyGroup, lebar: 202, isi: (g) => (teksDari(g, 'TreatyGroup') === '' ? LIMITS_NP.belumDipilih : teksDari(g, 'TreatyGroup')) },
        {
          judul: LIMITS_NP.rolProfile,
          lebar: 104,
          isi: (g, r) => (
            <input
              type="checkbox"
              aria-label={`${LIMITS_NP.rolProfile} ${String(r + 1)}`}
              checked={teksDari(g, 'IsROLProfile') === 'true'}
              disabled={!bisaUbah}
              onChange={(e) => {
                onUbahBaris(r, { ...g, IsROLProfile: e.target.checked ? 'true' : 'false' })
              }}
            />
          ),
        },
      ]}
      baris={grup}
      tombol={
        bisaUbah
          ? {
              lebar: 147,
              sisip: 1,
              tambah: { label: LIMITS_NP.tambahGrup, onKlik: onTambah },
              hapus: {
                label: LIMITS_NP.hapus,
                akses: (_, r) => `${LIMITS_NP.hapus} ${LIMITS_NP.treatyGroup} ${String(r + 1)}`,
                onKlik: onHapus,
              },
            }
          : undefined
      }
      rincian={(g, r) => (
        <RincianGrup
          g={g}
          opsi={opsi}
          bisaUbah={bisaUbah}
          onUbah={(baru) => {
            onUbahBaris(r, baru)
          }}
          onGantiGrup={(baru) => {
            onGantiBaris(r, baru)
          }}
        />
      )}
    />
  )
}

/** Lebar kolom grid `.Reinstatement_List` (`Layers.xml` @770156, `pyWidth`). */
const LEBAR_REINST: Readonly<Record<string, number>> = {
  ReinstatementValue: 156,
  ReinstatementPct: 146,
  ReinstatementNote: 319,
  AdditionalAmount1: 236,
  AdditionalAmount2: 231,
  AdditionalPct: 132,
  ReinstatementAmount1: 302,
  ReinstatementAmount2: 304,
}

/**
 * Lebar grid `Min Premium Amt` / `MDP` (`Layers.xml` @992599 · @1089868):
 * mata uang 198 · nilai 363 seperti ekspor; kolom tombol 120, bukan 100/101.
 *
 * ⚠️ Simpangan kecil yang disengaja (8 Oktober 2026): huruf tombol tema ini
 * (12px tebal) lebih lebar dari skin Pega (±11px), dan tabel `table-layout:
 * fixed` tidak melebarkan kolom mengikuti isinya — pada 100/661 tombol `Add
 * MDP` di kolom TERAKHIR meluber dan terpotong tepi grid menjadi "Add MD".
 */
const LEBAR_GRID_MDP: readonly [number, number, number] = [198, 363, 120]

/** Rincian satu layer — `Section/Layers.xml`. Diekspor untuk uji. */
export function RincianLayer({
  l,
  opsi,
  bisaUbah,
  modeUbah,
  onUbah,
  hitung,
  onGantiGrup,
}: {
  l: SimpulLimit
  opsi: OpsiLimits
  bisaUbah: boolean
  /** Mode Edit saja — No RIP dan ROL % tidak ikut kunci materialitas. */
  modeUbah: boolean
  onUbah: (l: SimpulLimit) => void
  hitung: (aksi: AksiLimitNP, baris?: number) => void
  onGantiGrup: (l: SimpulLimit) => void
}) {
  const set = (kunci: string) => (v: string) => {
    onUbah({ ...l, [kunci]: v })
  }
  const grup = larikDari(l, 'TreatyGroupList')
  const reinst = larikDari(l, 'Reinstatement_List')
  // Kolom Reinstatement Amount USD SELALU tampil, walau Limit2 kosong —
  // syarat `.Limit2 != 0` @801378 tidak diikuti: layar Pega produksi tetap
  // menampilkannya (keputusan pemakai 9 Oktober 2026).
  const kolomReinst = KOLOM_REINSTATEMENT

  /** Sel mata uang di kepala kolom matriks — satu properti per kolom. */
  const mataUang = (kunciMU: 'Currency' | 'Currency2', label: string) => (
    <DropdownDaftar
      label={label}
      nilai={teksDari(l, kunciMU)}
      pilihan={opsi.mataUang}
      bisaUbah={bisaUbah}
      onPilih={(nama, id) => {
        // Pega: ketiga baris satu mata uang mengikat SATU properti, dan
        // keduanya menimpa `.CurrencyID` (`col:.ID->.CurrencyID`).
        onUbah({ ...l, [kunciMU]: nama, CurrencyID: id })
      }}
    />
  )
  const selMatriks = (kunci: string, desimal: number | null, label: string) =>
    bisaUbah ? (
      <input
        className="field__input"
        value={teksDari(l, kunci)}
        aria-label={label}
        onChange={(e) => {
          onUbah({ ...l, [kunci]: saringAngka(e.target.value) })
        }}
      />
    ) : (
      tampil(teksDari(l, kunci), desimal)
    )
  const barisMatriks: readonly { label: string; k1: string; k2: string; desimal: number | null }[] = [
    { label: LIMITS_NP.limit100, k1: 'Limit', k2: 'Limit2', desimal: null },
    { label: LIMITS_NP.agregat, k1: 'AgregateLimit', k2: 'AgregateLimit2', desimal: null },
    { label: LIMITS_NP.deductible, k1: 'Deductible', k2: 'Deductible2', desimal: 2 },
  ]

  // ⭐ Urutan = Section `Layers` (bentuk Pega, 8 Oktober 2026), TANPA judul
  // bagian tambahan: jenis/nomor layer · `Part of` · grid Treaty Group ·
  // `Egnpi this layer` · Cover · Currency Relation · matriks limit · No RIP ·
  // `Reinstatement` · grid Reinstatement · Adj Rate/Premium Earned · Min
  // Premium · MDP · Combine MDP · ROL %.
  return (
    <>
      {/* `Layers.xml` `Inline` @11943 — jenis/nomor layer · teks `Part of`
          (sel LABEL @34562) · bagiannya, sebaris dan TANPA label: keempat
          sel berlabel kosong di ekspor, dan tangkapan layar 31 memperlihatkan
          `[Layer ▾] [1] Part of [Layer ▾] [1]`. */}
      <TataPegaBlok tata="alir">
        <PilihNP label="" nilai={teksDari(l, 'LayerType')} opsi={opsi.jenisLayer ?? []} bisaUbah={bisaUbah} onUbah={set('LayerType')} />
        <MedanNP label="" nilai={teksDari(l, 'Layer')} desimal={null} bisaUbah={bisaUbah} onUbah={set('Layer')} />
        <TeksPega>{LIMITS_NP.partOf}</TeksPega>
        <PilihNP label="" nilai={teksDari(l, 'LayerPartType')} opsi={opsi.jenisLayer ?? []} bisaUbah={bisaUbah} onUbah={set('LayerPartType')} />
        <MedanNP label="" nilai={teksDari(l, 'LayerPart')} desimal={null} bisaUbah={bisaUbah} onUbah={set('LayerPart')} />
      </TataPegaBlok>

      {/* Grid Treaty Group — `masterDetail`, rincian `CoBList`. */}
      <GridGrup
        grup={grup}
        opsi={opsi}
        bisaUbah={bisaUbah}
        onTambah={() => {
          onUbah({ ...l, TreatyGroupList: [...grup, { TreatyGroup: '', TreatyGroupID: '', IsROLProfile: 'false', ClassOfBusinessList: [] }] })
        }}
        onUbahBaris={(r, baru) => {
          onUbah({ ...l, TreatyGroupList: ganti(grup, r, baru) })
        }}
        onHapus={(r) => {
          // deleteRow → `TotalEgnpi`
          onGantiGrup({ ...l, TreatyGroupList: grup.filter((_, j) => j !== r) })
        }}
        onGantiBaris={(r, baru) => {
          // DT `SetIndexLayer_DT` + `TotalEgnpi`
          onGantiGrup({ ...l, TreatyGroupList: ganti(grup, r, baru) })
        }}
      />
      <GridNilai
        judul={LIMITS_NP.egnpiLayer}
        judulNilai=""
        lebar={[196, 184]}
        baris={larikDari(l, 'EgnpiTotalList')}
        bisaUbahMataUang={false}
        bisaUbahNilai={false}
        tambah={false}
        mataUang={opsi.mataUang}
        onUbah={() => undefined}
      />

      {/* `Default` @7873 bertumpuk: Cover · Currency Relation (`Stacked
          with labels left` @8165) · matriks · No RIP · Reinstatement · grid ·
          Adjustment Rate · Premium Earned · Min Premium · MDP · Combine MDP ·
          ROL % — SEMUANYA satu per baris (tangkapan layar 32). */}
      <TataPegaBlok tata="kiri">
          <PilihNP label={LIMITS_NP.cover} nilai={teksDari(l, 'Cover')} opsi={opsi.cover ?? []} bisaUbah={bisaUbah} onUbah={set('Cover')} />
          <PilihNP
            label={LIMITS_NP.relasi}
            nilai={teksDari(l, 'CurrencyRelation')}
            opsi={opsi.relasiMataUang ?? []}
            kosong={LIMITS_NP.pilihRelasi}
            bisaUbah={bisaUbah}
            onUbah={set('CurrencyRelation')}
          />
        </TataPegaBlok>
        {/* `Inline grid double` @9155: [ `Inline grid double` @9444 mata uang 1
            | `Inline grid double` @15826 mata uang 2 ]. Tiap sisi berbaris
            [`Stacked with labels left` mata uang berlabel | nilai], lalu dua
            "Spacer" `1=2` yang TETAP memakan slot — baris kosong di antara
            100 % Limit · Agregate Year Limit · Deductible (tangkapan layar
            32). Ketiga mata uang satu sisi mengikat SATU properti, seperti
            Pega (`.Currency` / `.Currency2`). */}
        <TataPegaBlok tata="g2">
          {(['Currency', 'Currency2'] as const).map((mu, s) => (
            <TataPegaBlok key={mu} tata="g2">
              {barisMatriks.map((b) => (
                <Fragment key={b.k1}>
                  <TataPegaBlok tata="kiri">{mataUang(mu, b.label)}</TataPegaBlok>
                  <div>{selMatriks(s === 0 ? b.k1 : b.k2, b.desimal, `${b.label} ${String(s + 1)}`)}</div>
                  <SelKosongPega />
                  <SelKosongPega />
                </Fragment>
              ))}
            </TataPegaBlok>
          ))}
        </TataPegaBlok>

        <TataPegaBlok tata="alir">
          <Centang label={LIMITS_NP.noRIP} nilai={teksDari(l, 'NoRIPCalculation')} bisaUbah={modeUbah} onUbah={set('NoRIPCalculation')} />
        </TataPegaBlok>
        {/* `ReinstatementValue` → `SetReinstatementPct` saat berubah.
            Blok `Inline` @698608: sel LABEL `Reinstatement` (@704972) lalu
            isian TANPA label (@714520) — sebaris, SESUDAH No RIP (gambar 32). */}
        <TataPegaBlok tata="alir">
          <TeksPega>{LIMITS_NP.reinstatement}</TeksPega>
          <MedanNP
            label=""
            nilai={teksDari(l, 'ReinstatementValue')}
            desimal={0}
            bisaUbah={bisaUbah}
            onUbah={set('ReinstatementValue')}
            onLepas={() => {
              hitung('reinstatement')
            }}
          />
        </TataPegaBlok>
        {/* Grid `.Reinstatement_List` @770156 — grid `row` Pega berkolom
            ekspor (lebar `LEBAR_REINST`), kepala tipis dan "No items" — bukan
            tabel tema yang menggulir sendiri (gambar 32). */}
        <GridPega
          label={LIMITS_NP.reinstatement}
          kolom={kolomReinst.map((k) => ({
            judul: k.label,
            lebar: LEBAR_REINST[k.kunci] ?? 150,
            angka: k.desimal !== null,
            isi: (b: SimpulLimit, r: number) => {
              const v = teksDari(b, k.kunci)
              const ubahSel = (x: string) => {
                onUbah({ ...l, Reinstatement_List: ganti(reinst, r, { ...b, [k.kunci]: x }) })
              }
              // DT per sel: % Additional Premium → ReCalculateReinstatement;
              // Reinstatement % → CalculateReinstatement; Amount IDR →
              // CalculateReinstatementPct.
              const pemicu: Partial<Record<string, AksiLimitNP>> = {
                ReinstatementPct: 'reinst-tambahan',
                AdditionalPct: 'reinst-jumlah',
                ReinstatementAmount1: 'reinst-persen',
              }
              const aksi = pemicu[k.kunci]
              // ⛔ `Note` DIDAHULUKAN — keluhan pemilik proses 8 Oktober 2026:
              // *"di nonprop bagian ReinstatementNote itu seharusnya mengambil
              // Prompt value bukan standard value"*.
              //
              // Dahulu cabang ini duduk SESUDAH `!bisaUbah`, jadi mode lihat
              // mengembalikan `v` — kode tersimpannya (`asamount`) — dan tidak
              // pernah sampai ke penerjemah. `PilihNP` mode baca yang
              // menerjemahkannya (`opsi.find(...)?.label ?? nilai`), jadi
              // selnya memakai komponen yang sama di KEDUA mode; yang berbeda
              // hanya boleh atau tidaknya diubah.
              //
              // ⚠️ Nilai di luar daftar tetap tampil apa adanya — tidak
              // ditebak, tidak dikosongkan.
              if (k.kunci === 'ReinstatementNote') {
                return <PilihNP label="" nilai={v} opsi={opsi.catatanReinstatement ?? []} bisaUbah={bisaUbah} onUbah={ubahSel} />
              }
              if (k.kunci === 'ReinstatementValue' || !bisaUbah) {
                return k.desimal === null ? v : tampil(v, k.desimal)
              }
              return (
                /* ⭐ `change` Pega: DT baris berjalan HANYA bila sel ini
                   berubah sejak dimasuki — sel yang hanya dilewati tidak
                   menimpa nilai baris yang sudah disunting. */
                <PemicuUbah
                  nilai={v}
                  aktif={aksi !== undefined}
                  aksi={() => {
                    if (aksi !== undefined) hitung(aksi, r)
                  }}
                >
                  <input
                    className="field__input"
                    value={v}
                    aria-label={k.label}
                    onChange={(e) => {
                      ubahSel(saringAngka(e.target.value))
                    }}
                  />
                </PemicuUbah>
              )
            },
          }))}
          baris={reinst}
        />
        <TataPegaBlok tata="tumpuk">
          <div>
            <MedanNP
              label={LIMITS_NP.adjRate}
              nilai={teksDari(l, 'AdjRate')}
              desimal={null}
              bisaUbah={bisaUbah}
              onUbah={set('AdjRate')}
              onLepas={() => {
                hitung('adj')
              }}
            />
            <GridNilai
              judul={LIMITS_NP.premiEarned}
              baris={larikDari(l, 'PremiumEarnedList')}
              bisaUbahMataUang={false}
              bisaUbahNilai={bisaUbah}
              tambah={false}
              mataUang={opsi.mataUang}
              onUbah={(b) => {
                onUbah({ ...l, PremiumEarnedList: b })
              }}
            />
          </div>
          <div>
            <MedanNP
              label={LIMITS_NP.mdpMinPct}
              nilai={teksDari(l, 'MDPMinPct')}
              desimal={null}
              bisaUbah={bisaUbah}
              onUbah={set('MDPMinPct')}
              onLepas={() => {
                hitung('mdp')
              }}
            />
            <GridNilai
              judul={LIMITS_NP.mdpMinAmt}
              lebar={LEBAR_GRID_MDP}
              baris={larikDari(l, 'MDPMinList')}
              bisaUbahMataUang={bisaUbah}
              bisaUbahNilai={bisaUbah}
              tambah={bisaUbah}
              mataUang={opsi.mataUang}
              onUbah={(b) => {
                onUbah({ ...l, MDPMinList: b })
              }}
            />
          </div>
          <div>
            <MedanNP
              label={LIMITS_NP.mdpPct}
              nilai={teksDari(l, 'MDPPct')}
              desimal={null}
              bisaUbah={bisaUbah}
              onUbah={set('MDPPct')}
              onLepas={() => {
                hitung('mdp')
              }}
            />
            <GridNilai
              judul={LIMITS_NP.mdp}
              lebar={LEBAR_GRID_MDP}
              baris={larikDari(l, 'MDPList')}
              bisaUbahMataUang={bisaUbah}
              bisaUbahNilai={bisaUbah}
              tambah={bisaUbah}
              mataUang={opsi.mataUang}
              onUbah={(b) => {
                onUbah({ ...l, MDPList: b })
              }}
            />
          </div>
        </TataPegaBlok>
        <TataPegaBlok tata="tumpuk">
          <Centang label={LIMITS_NP.combineMDP} nilai={teksDari(l, 'IsCombineMDP')} bisaUbah={bisaUbah} onUbah={set('IsCombineMDP')} />
          {/* ⛔ ROL % TIDAK PERNAH diketik — Section `Layers` / `LayersEDM`
              `.ROLPct`: `pyDisabled=true`, `pyDisabledNew=always`. Nilainya
              hanya dari `DetailCalculationROL` (Adjustment Rate / MDP %).
              Dulu medan ini terbuka di mode Edit; "1,659" yang diketik
              berkoma terbaca NOL oleh services, dan `Total ROL` jadi 0,00
              (laporan pemakai 8 Oktober 2026). */}
          <MedanNP label={LIMITS_NP.rol} nilai={teksDari(l, 'ROLPct')} desimal={null} bisaUbah={false} onUbah={set('ROLPct')} />
        </TataPegaBlok>
    </>
  )
}

/**
 * Kolom grid layer — ekspor `TreatyIn.Limits` (`masterDetail`): `Layers` ·
 * ` ` · `of` · ` ` · ` ` (jenis/nomor layer dan bagiannya, lebar 102 · 102 ·
 * 102 · 102 · 105) lalu keempat nilai `KOLOM_LAYER` (185 · 187 · 185 · 187).
 */
const KOLOM_GRID_LAYER: readonly KolomPega<SimpulLimit>[] = [
  { judul: LIMITS_NP.layers, lebar: 102, isi: (l) => teksDari(l, 'LayerType') },
  { judul: '', lebar: 102, isi: (l) => teksDari(l, 'Layer') },
  { judul: '', lebar: 102, isi: () => LIMITS_NP.of },
  { judul: '', lebar: 102, isi: (l) => teksDari(l, 'LayerPartType') },
  { judul: '', lebar: 105, isi: (l) => teksDari(l, 'LayerPart') },
  ...KOLOM_LAYER.map((k, c) => ({
    judul: k.label,
    lebar: [185, 187, 185, 187][c] ?? 185,
    angka: true,
    isi: (l: SimpulLimit) => (teksDari(l, k.kunci) === '' ? '' : tampil(teksDari(l, k.kunci), k.desimal)),
  })),
]

/** Lebar kolom `Summary of Limit` dari ekspor (`pyWidth`). */
const LEBAR_RINGKASAN = [252, 170, 172, 100, 100, 170, 170, 170, 170] as const

const OPSI_KOSONG: OpsiLimits = { jenisTreaty: [], kelompokTreaty: [], mataUang: [] }

/** Larik akar kosong — kontrak tanpa Summary/Total. */
const AKAR_KOSONG: LimitsAkar = { LimitSummaryList: [], Total: {}, TotalLimitsROL: '' }

export default function TabLimitsNonProp({
  petunjukKosong,
  pohon,
  akar: akarAwal,
  egnpi,
  kurs,
  edmState = '',
  edmJenisMaterial = '',
  mode,
  opsi: opsiAwal,
  onUbah,
}: {
  /** Petunjuk bila kontrak tak punya layer — layar membedakan sebabnya. */
  petunjukKosong?: string
  pohon: readonly SimpulLimit[]
  akar?: LimitsAkar
  egnpi: readonly EgnpiLimitNP[]
  kurs: readonly KursLimitNP[]
  edmState?: string
  /**
   * `EDMMaterialType` — bila `2`, hampir seluruh kontrol tab ini dimatikan
   * (`pyDisabledWhen TreatyIn.EDMMaterialType = 2` di setiap sel kecuali
   * No RIP dan ROL %).
   */
  edmJenisMaterial?: string
  mode: ModeForm
  /** Isi dropdown; bila tidak diberikan, diminta sendiri. */
  opsi?: OpsiLimits
  /**
   * Keadaan terkini layer dan larik akar, dilaporkan tiap kali berubah —
   * form menyimpannya supaya ia hidup lintas pindah tab dan tab Share
   * membaca layer yang SAMA (`TreatyIn.Limits` di clipboard Pega).
   */
  onUbah?: (layers: SimpulLimit[], akar: LimitsAkar) => void
}) {
  const [layers, setLayers] = useState<SimpulLimit[]>(() => [...pohon])
  const [akar, setAkar] = useState<LimitsAkar>(akarAwal ?? AKAR_KOSONG)
  const onUbahTerkini = useRef(onUbah)
  onUbahTerkini.current = onUbah
  useEffect(() => {
    onUbahTerkini.current?.(layers, akar)
  }, [layers, akar])
  const [opsi, setOpsi] = useState<OpsiLimits>(opsiAwal ?? OPSI_KOSONG)
  const [pesan, setPesan] = useState<string[]>([])
  const [gagal, setGagal] = useState('')
  // `TreatyIn.ViewState != 1` — tombol dan isian hidup di mode Edit saja.
  const modeUbah = mode === 'ubah'
  // … dan dimatikan pula oleh kunci materialitas.
  const bisaUbah = modeUbah && edmJenisMaterial.trim() !== '2'

  useEffect(() => {
    if (opsiAwal !== undefined) return
    let dibuang = false
    ambilOpsiLimits()
      .then((o) => {
        if (!dibuang) setOpsi(o)
      })
      .catch(() => undefined)
    return () => {
      dibuang = true
    }
  }, [opsiAwal])

  /**
   * Menjalankan satu Activity di services lalu menggabung HANYA medan yang
   * Activity itu tulis (`DITULIS_LIMIT_NP`).
   */
  const hitung = (aksi: AksiLimitNP, indeks = 0, baris = 0, dasar: readonly SimpulLimit[] = layers) => {
    if (!modeUbah) return
    setGagal('')
    hitungLimitNP({ aksi, layers: dasar, egnpi, kurs, indeks, baris })
      .then((h) => {
        setPesan(h.pesan)
        const ditulis = DITULIS_LIMIT_NP[aksi]
        if (ditulis.length > 0) {
          setLayers((kini) =>
            kini.map((l, i) => {
              if (!SEMUA_LAYER.includes(aksi) && i !== indeks) return l
              const hasil = h.layers[i]
              if (hasil === undefined) return l
              const baru: SimpulLimit = { ...l }
              for (const k of ditulis) {
                const v = hasil[k]
                if (v !== undefined) baru[k] = v
              }
              return baru
            }),
          )
        }
        if (h.Total !== undefined) {
          setAkar({
            LimitSummaryList: h.LimitSummaryList ?? [],
            Total: h.Total,
            TotalLimitsROL: h.TotalLimitsROL ?? '',
          })
        }
      })
      .catch((e: unknown) => {
        setGagal(e instanceof Error ? e.message : String(e))
      })
  }

  const layerKosong: SimpulLimit = {
    TreatyGroupList: [],
    EgnpiTotalList: [],
    PremiumEarnedList: [],
    MDPList: [],
    MDPMinList: [],
    Reinstatement_List: [],
  }

  // ⭐ Urutan = Section `TreatyInTabsNonProportional` tab `Limits` (bentuk
  // Pega, 8 Oktober 2026): grid layer `masterDetail` (rincian `Layers`), blok
  // `Summary of Limit`, blok `Total All Layers` dengan `Total ROL` dan
  // tombolnya DI BAWAH grid. Kartu bernomor berchip sebelumnya DIGANTI grid.
  return (
    <div className="trin__blok trin__tab">
      <GridPega
        label={LIMITS_NP.layers}
        kolom={KOLOM_GRID_LAYER}
        baris={layers}
        kosong={petunjukKosong ?? LIMITS_NP.tanpaBaris}
        tombol={
          bisaUbah
            ? {
                lebar: 108,
                tambah: {
                  label: LIMITS_NP.tambahLayer,
                  onKlik: () => {
                    // `TreatyInNonAddItem(limits)`: ID = jumlah layer.
                    // `CopyLastLimitNP` tidak tercapai (Activity keluar
                    // sesudah langkah 5).
                    setLayers([...layers, { ...layerKosong, ID: String(layers.length + 1) }])
                  },
                },
                hapus: {
                  label: LIMITS_NP.hapus,
                  akses: (_, i) => `${LIMITS_NP.hapus} ${LIMITS_NP.layers} ${String(i + 1)}`,
                  onKlik: (i) => {
                    setLayers(layers.filter((_, j) => j !== i))
                  },
                },
              }
            : undefined
        }
        rincian={(l, i) => (
          <RincianLayer
            l={l}
            opsi={opsi}
            bisaUbah={bisaUbah}
            modeUbah={modeUbah}
            onUbah={(baru) => {
              setLayers(ganti(layers, i, baru))
            }}
            hitung={(aksi, baris = 0) => {
              hitung(aksi, i, baris)
            }}
            onGantiGrup={(baru) => {
              // Treaty Group berubah/dihapus → `TotalEgnpi` SEMUA layer;
              // ID layer = indeksnya (`TreatyTypeSetIndex`).
              const semua = ganti(layers, i, { ...baru, ID: String(i + 1) })
              setLayers(semua)
              hitung('egnpi', i, 0, semua)
            }}
          />
        )}
      />
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

      <BlokPega judul={LIMITS_NP.ringkasan}>
        <GridBaca kolom={KOLOM_RINGKASAN} lebar={LEBAR_RINGKASAN} baris={akar.LimitSummaryList} />
      </BlokPega>

      <BlokPega judul={LIMITS_NP.totalSemua}>
        {/* `Total All Layers` @44149 (`Default` @44346): keempat grid
            BERTUMPUK, satu per baris (bukan berdampingan). */}
        <TataPegaBlok tata="tumpuk">
          {GRID_TOTAL.map((g) => (
            <GridNilai
              key={g.kunci}
              judul={g.judul}
              lebar={[193, 349]}
              baris={akar.Total[g.kunci] ?? []}
              bisaUbahMataUang={false}
              bisaUbahNilai={false}
              tambah={false}
              mataUang={[]}
              onUbah={() => undefined}
            />
          ))}
        </TataPegaBlok>
        {/* `Inline grid double` @53680: [ `Inline grid 30 70` @53979 Total
            ROL | tombol `1=2` ×3 | `Inline grid double` @55990 tombol ]. Sel
            tersembunyi TETAP memakan slotnya (tangkapan layar 30/33): tombol
            turun ke baris KETIGA, separuh kiri. */}
        <TataPegaBlok tata="g2">
          {/* Baris [label 30% | nilai 70%]. */}
          <TataPegaBlok tata="t3070">
            <TeksPega>{LIMITS_NP.totalROL}</TeksPega>
            <Field label="" value={tampil(akar.TotalLimitsROL, 2)} readOnly onChange={() => undefined} />
          </TataPegaBlok>
          <SelKosongPega />
          <SelKosongPega />
          <SelKosongPega />
          {modeUbah && (
            <TataPegaBlok tata="g2">
              <div>
                <button
                  type="button"
                  className="btn btn--primary btn--sm"
                  disabled={!bisaUbah}
                  onClick={() => {
                    hitung('total')
                  }}
                >
                  {LIMITS_NP.perbaruiTotal}
                </button>
              </div>
              {kontrakRevisi(edmState) && (
                <div>
                  <button
                    type="button"
                    className="btn btn--sm"
                    disabled={!bisaUbah}
                    onClick={() => {
                      hitung('nilai-list')
                    }}
                  >
                    {LIMITS_NP.perbaruiList}
                  </button>
                </div>
              )}
            </TataPegaBlok>
          )}
        </TataPegaBlok>
      </BlokPega>
    </div>
  )
}
