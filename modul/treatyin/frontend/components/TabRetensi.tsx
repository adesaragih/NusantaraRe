// Tab **Maximum Retention** cabang NON-PROPORSIONAL — bentuk dan rumus
// dari ekspor.
//
// ---------------------------------------------------------------------
// ⛔ YANG DIGANTI, DAN MENGAPA
// ---------------------------------------------------------------------
// Sebelumnya: `TabGridWarisan` lima kolom hanya-baca + `PanelTotalRetensi`
// yang tombol `Update Total`-nya MATI. Ekspor memperlihatkan tab yang dapat
// disunting:
//
//   Section/TreatyInTabsNonProportional.xml   wadah TABBED ke-1
//     panel  `Maximum Retention`     grid 3 kolom + Add + Delete
//     panel  `Total Retention Amount` | `Value`   per mata uang
//     tombol `Update Total`  → TreatyInNPSetTotal(type=retention)
//   Section/MaxRetention.xml                  rincian: Treaty Group ·
//     Amount (mata uang + nilai) · Note
//
// ---------------------------------------------------------------------
// ⚠️ SATU TOMBOL, BUKAN DUA — DAN ITU BUKAN KELALAIAN
// ---------------------------------------------------------------------
// Tab EGNPI punya `Update Total` DAN `Update EGNPI Value`; tab ini hanya
// yang pertama. Sebabnya: retensi nol punya kolom `Amount in IDR`, jadi nol
// konversi kurs yang perlu dijalankan. Menambahkan tombol kedua di sini
// berarti membuat tombol yang di Pega tidak ada dan tidak punya pekerjaan.
//
// ---------------------------------------------------------------------
// ⭐ RUMUSNYA DI SERVICES, BUKAN DI SINI
// ---------------------------------------------------------------------
// `backend/services/hitung_retensi.go` — disalin dari Activity dan diuji.
// Layar ini mengirim isian dan menampilkan jawabannya; nol aritmetika di
// berkas ini, sama seperti tab Limits dan EGNPI.
//
// ⛔ Aturan B pemilik proses tetap berlaku: mengetik TIDAK menyentuh basis
// data. Rute yang dipanggil ber-awalan `/hitung/`, yang menyatakan dirinya
// nol tulis.

import { useEffect, useId, useState } from 'react'

import { FieldAngka } from '../../../../inti/frontend/components/ui/dasar'
import {
  ambilOpsiLimits,
  hitungRetensi,
  type AksiRetensi,
  type BarisRetensi,
  type OpsiLimits,
  type PilihanWarisan,
} from '../api'
import { TOTAL_RETENSI } from '../labels'
import { DESIMAL_RETENSI, RETENSI, totalTerkunci } from '../labelsRetensi'
import { useProperti } from '../halaman'
import type { ModeForm } from '../mode'
import { DropdownDaftar } from './IsianAuto'
import { BlokPega, DeretTombolPega, GridNilaiPega, GridPega } from './gridPega'
import { TataPegaBlok } from './tataPega'
import { formatLimit } from './TabLimitsProp'

/** Angka tampil — pemformat modul, bukan pemformat kedua. */
const tampil = (nilai: string, desimal: number) => formatLimit('uang', desimal, nilai)

interface BarisTotal {
  Currency: string
  CurrencyID: string
  Value: string
}

/**
 * Grid nilai baca-saja — `Total Retention Amount` | `Value`. ⭐ Judul kolom
 * pertama ADALAH nama panelnya — @152591 — dan isinya mata uang
 * (`.Currency` @161105). Bentuk Pega datar, selebar `pyWidth` (193 · 349).
 */
function GridTotal({ baris }: { baris: readonly BarisTotal[] }) {
  return (
    <GridNilaiPega
      judul={TOTAL_RETENSI.judul}
      nilai={TOTAL_RETENSI.kolomNilai}
      baris={baris}
      lebar={[193, 349]}
      tampil={(v) => tampil(v, DESIMAL_RETENSI.nilaiTotal)}
    />
  )
}

/** Rincian satu baris — `Section/MaxRetention.xml`, urut dan hak ubahnya. */
function RincianBaris({
  b,
  bisaUbah,
  mataUang,
  kelompokTreaty,
  onUbah,
}: {
  b: BarisRetensi
  bisaUbah: boolean
  mataUang: readonly PilihanWarisan[]
  kelompokTreaty: readonly PilihanWarisan[]
  onUbah: (baru: BarisRetensi) => void
}) {
  const idNote = useId()
  return (
    <>
      {/* `MaxRetention.xml` @459 `Stacked with labels left`: Treaty Group ·
          baris Amount (`Inline grid 30 70` @1105 = baris [30% `Stacked
          with labels left` @1403 mata uang berlabel Amount | 70% nilai]) ·
          Note. Tangkapan layar 27 Non-Prop. */}
      <TataPegaBlok tata="kiri">
        {/* ⚠️ EKSPOR BERSELISIH DENGAN DIRINYA SENDIRI, sama seperti di tab
            EGNPI: grid membuat `.TreatyGroup` dapat diubah, rinciannya
            `pyReadOnly=true`. Yang dipakai bentuk GRID — `Add` melahirkan
            baris KOSONG, dan baris yang kelompok treaty-nya tidak dapat
            diisi tidak berguna bagi siapa pun. */}
        <DropdownDaftar
          label={RETENSI.treatyGroup}
          nilai={b.TreatyGroup}
          pilihan={kelompokTreaty}
          bisaUbah={bisaUbah}
          onPilih={(nama, id) => {
            onUbah({ ...b, TreatyGroup: nama, TreatyGroupID: id })
          }}
        />

        {/* ⭐ SATU BARIS BERLABEL `Amount` memuat DUA medan — pilihan mata
            uang lalu kotak nilai. Itu bunyi ekspornya
            (`pyLabelFieldValue = Amount` menempel pada `.Currency`, dan
            `.Amount` di sebelahnya tanpa label sendiri) dan itu pula yang
            terlihat di tangkapan layar pemilik proses. */}
        <TataPegaBlok tata="t3070">
          <TataPegaBlok tata="kiri">
            <DropdownDaftar
              label={RETENSI.jumlah}
              nilai={b.Currency}
              pilihan={mataUang}
              bisaUbah={bisaUbah}
              onPilih={(nama, id) => {
                onUbah({ ...b, Currency: nama, CurrencyID: id })
              }}
            />
          </TataPegaBlok>
          <FieldAngka
            label=""
            value={b.Amount}
            desimal={DESIMAL_RETENSI.jumlah}
            readOnly={!bisaUbah}
            onChange={(v) => {
              onUbah({ ...b, Amount: v })
            }}
          />
        </TataPegaBlok>

        {/* ⛔ RALAT 7 Oktober 2026 — medan ini SEMPAT DIKUNCI MATI.
            `Section/MaxRetention.xml` berbunyi `pyReadOnly = true`, tetapi
            sel yang SAMA juga memuat
            `pyReadOnlyCondition = TreatyIn.IsEditData = 1` — dan SYARAT
            itulah yang berlaku.

            ⚠️ SYARATNYA BUKAN `ViewState`, beda dengan tab EGNPI yang medan
            serupanya memakai `TreatyIn.ViewState = 1`. Nol tambalan
            menyeluruh untuk keduanya.

            ⛔ `IsEditData` NOL dibawa aplikasi — ia tidak ada di
            `KontrakWarisan` maupun di mana pun. Jadi ia DIDEKATI dengan mode
            Edit, mengikuti preseden `TabPortofolio.tsx` (sel 112·113·114,
            `pyReadOnlyCondition = TreatyIn.IsEditData = 1`) yang sudah
            melakukan hal yang sama dan diuji begitu. Pendekatan ini salah
            HANYA ketika `IsEditData = 1`, keadaan yang belum pernah terukur
            di aplikasi.

            ⭐ Tetap `textarea` ber-`readOnly` sungguhan di mode lihat —
            `Area` inti nol punya `readOnly`, dan ia tidak diubah dari sini. */}
        <div className="field">
          <label className="field__label" htmlFor={idNote}>
            {RETENSI.keterangan}
          </label>
          <textarea
            id={idNote}
            className="field__input"
            value={b.Note}
            rows={4}
            readOnly={!bisaUbah}
            onChange={(e) => {
              onUbah({ ...b, Note: e.target.value })
            }}
          />
        </div>
      </TataPegaBlok>
    </>
  )
}

export default function TabRetensi({
  baris,
  totalAwal = [],
  opsiAwal,
  edmJenisMaterial = '',
  mode = 'lihat',
}: {
  baris: readonly BarisRetensi[]
  /**
   * Total TERSIMPAN dari dokumen — yang tampil SEBELUM `Update Total`
   * ditekan.
   *
   * ⭐ Pega memperlihatkan total yang sudah ada di clipboard kontrak, bukan
   * panel kosong: membiarkannya kosong sampai tombol ditekan akan terbaca
   * sebagai "kontrak ini nol retensi", dan itu berbeda dari "belum
   * dihitung ulang".
   */
  totalAwal?: readonly BarisTotal[]
  /** Isi dropdown; bila nol diberikan, tab ini mengambilnya sendiri. */
  opsiAwal?: OpsiLimits
  /** `TreatyIn.EDMMaterialType` — `2` mematikan `Update Total`. */
  edmJenisMaterial?: string
  mode?: ModeForm
}) {
  // ⭐ PENAMPUNG HALAMAN (`../halaman.tsx`) — 7 Oktober 2026. Isian bertahan
  // saat pindah tab, dan tab lain membaca properti yang SAMA.
  //
  // ⛔ PERBAIKAN "input tiba-tiba hilang": efek `setRows([...baris])` atas
  // perubahan `baris` DICABUT. `baris` dibuat ulang (`.map`) tiap form
  // dirender, jadi efek itu mengembalikan baris ke data kontrak setiap kali
  // apa pun di form berubah. Kontrak lain = penampung dikosongkan + tab
  // dilahirkan ulang (`key` fieldset form).
  const [rows, setRows] = useProperti<BarisRetensi[]>('Retention', () => [...baris])
  const [total, setTotal] = useProperti<BarisTotal[]>('TotalRetentionAmountNP', () => [...totalAwal])
  const [gagal, setGagal] = useState('')
  const [opsi, setOpsi] = useState<OpsiLimits>(opsiAwal ?? { jenisTreaty: [], kelompokTreaty: [], mataUang: [] })

  const bisaUbah = mode === 'ubah'

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

  function jalankan(aksi: AksiRetensi, indeks = 0, dasar: readonly BarisRetensi[] = rows) {
    setGagal('')
    hitungRetensi({ aksi, retensi: dasar, indeks })
      .then((h) => {
        setRows(h.retensi)
        setTotal(h.TotalRetentionAmountNP ?? [])
      })
      .catch((e: unknown) => {
        setGagal(e instanceof Error ? e.message : String(e))
      })
  }

  // ⭐ Urutan = Section `TreatyInTabsNonProportional` tab `Maximum Retention`
  // (bentuk Pega, 8 Oktober 2026): blok `Maximum Retention` — grid
  // `masterDetail` Treaty Group · Currency · Amount (lebar 245 · 127 · 379,
  // kolom tombol 153), rincian `MaxRetention` — lalu grid `Total Retention
  // Amount` dan tombol `Update Total` DI BAWAHNYA. Kartu berchip DIGANTI.
  return (
    <div className="trin__blok trin__tab">
      <BlokPega judul={RETENSI.panel}>
        <GridPega
          label={RETENSI.panel}
          kolom={[
            { judul: RETENSI.treatyGroup, lebar: 245, isi: (b) => (b.TreatyGroup === '' ? RETENSI.barisBaru : b.TreatyGroup) },
            { judul: RETENSI.mataUang, lebar: 127, isi: (b) => b.Currency },
            { judul: RETENSI.jumlah, lebar: 379, angka: true, isi: (b) => tampil(b.Amount, DESIMAL_RETENSI.jumlah) },
          ]}
          baris={rows}
          tombol={
            bisaUbah
              ? {
                  lebar: 153,
                  // `TreatyInNonAddItem(retention)` — baris KOSONG, `ID=""`.
                  tambah: {
                    label: RETENSI.tambah,
                    onKlik: () => {
                      jalankan('tambah')
                    },
                  },
                  hapus: {
                    label: RETENSI.hapus,
                    akses: (_, i) => `${RETENSI.hapus} ${RETENSI.treatyGroup} ${String(i + 1)}`,
                    onKlik: (i) => {
                      jalankan('hapus', i)
                    },
                  },
                }
              : undefined
          }
          rincian={(b, i) => (
            <RincianBaris
              b={b}
              bisaUbah={bisaUbah}
              mataUang={opsi.mataUang}
              kelompokTreaty={opsi.kelompokTreaty}
              onUbah={(baru) => {
                setRows(rows.map((x, j) => (j === i ? baru : x)))
              }}
            />
          )}
        />
      </BlokPega>

      {gagal !== '' && (
        <p className="tl-pesan" role="alert">
          {gagal}
        </p>
      )}

      {/* Wadah tanpa kepala @3976 (`Default` @4169): grid Total Retention
          Amount lalu tombol Update Total di bawahnya. */}
      <TataPegaBlok tata="tumpuk">
        <GridTotal baris={total} />
        {bisaUbah && (
          <DeretTombolPega>
            {/* ⛔ MATI bila `TreatyIn.EDMMaterialType = 2` — `pyDisabledWhen`
                @202558. Tombol kedua tepat di atasnya di ekspor ber-`pyCondition
                1=2`: ia MATI dan tidak dibangun. Yang dibangun yang hidup. */}
            <button
              type="button"
              className="btn btn--primary btn--sm"
              disabled={totalTerkunci(edmJenisMaterial)}
              onClick={() => {
                jalankan('total')
              }}
            >
              {TOTAL_RETENSI.perbarui}
            </button>
          </DeretTombolPega>
        )}
      </TataPegaBlok>
    </div>
  )
}
