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

import { FieldAngka, Kosong, Panel } from '../../../../inti/frontend/components/ui/dasar'
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
import { Bagian, Chip, KartuLipat, KepalaBagian, TombolHapus, TombolTambah } from './limitsUI'
import { formatLimit } from './TabLimitsProp'

/** Angka tampil — pemformat modul, bukan pemformat kedua. */
const tampil = (nilai: string, desimal: number) => formatLimit('uang', desimal, nilai)

interface BarisTotal {
  Currency: string
  CurrencyID: string
  Value: string
}

/** Grid dua kolom hanya-baca — panel `Total Retention Amount`. */
function GridTotal({ baris }: { baris: readonly BarisTotal[] }) {
  return (
    <div className="table-wrap trin__share-total">
      <table className="trin__tabel">
        <thead>
          <tr>
            {/* ⭐ Judul kolom pertama ADALAH nama panelnya — @152591 — dan
                isinya mata uang (`.Currency` @161105). */}
            <th scope="col">{TOTAL_RETENSI.judul}</th>
            <th scope="col">{TOTAL_RETENSI.kolomNilai}</th>
          </tr>
        </thead>
        <tbody>
          {baris.length === 0 && (
            <tr>
              <td colSpan={2}>{TOTAL_RETENSI.tanpaBaris}</td>
            </tr>
          )}
          {baris.map((b, i) => (
            <tr key={String(i)}>
              <td>{b.Currency}</td>
              <td>{tampil(b.Value, DESIMAL_RETENSI.nilaiTotal)}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
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
    <Bagian judul={RETENSI.panel}>
      <div className="trin__retensi-medan">
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
        <div className="trin__retensi-jumlah">
          <DropdownDaftar
            label={RETENSI.jumlah}
            nilai={b.Currency}
            pilihan={mataUang}
            bisaUbah={bisaUbah}
            onPilih={(nama, id) => {
              onUbah({ ...b, Currency: nama, CurrencyID: id })
            }}
          />
          <FieldAngka
            label=""
            value={b.Amount}
            desimal={DESIMAL_RETENSI.jumlah}
            readOnly={!bisaUbah}
            onChange={(v) => {
              onUbah({ ...b, Amount: v })
            }}
          />
        </div>

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
      </div>
    </Bagian>
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

  return (
    <Panel judul={RETENSI.judul}>
      <div className="tl-rincian">
        <KepalaBagian
          judul={RETENSI.panel}
          jumlah={rows.length}
          aksi={
            bisaUbah && (
              <TombolTambah
                label={RETENSI.tambah}
                onClick={() => {
                  // `TreatyInNonAddItem(retention)` — baris KOSONG, `ID=""`.
                  jalankan('tambah')
                }}
              />
            )
          }
        />

        {rows.length === 0 ? (
          <Kosong pesan={RETENSI.petunjukKosong} />
        ) : (
          <div className="tl-daftar">
            {rows.map((b, i) => (
              <KartuLipat
                key={i}
                nomor={i + 1}
                judul={b.TreatyGroup}
                judulKosong={RETENSI.barisBaru}
                bukaAwal={b.TreatyGroup === ''}
                meta={
                  <>
                    <Chip label={RETENSI.mataUang} nilai={b.Currency} />
                    <Chip label={RETENSI.jumlah} nilai={tampil(b.Amount, DESIMAL_RETENSI.jumlah)} />
                  </>
                }
                aksi={
                  bisaUbah && (
                    <TombolHapus
                      label={RETENSI.hapus}
                      labelAkses={`${RETENSI.hapus} ${RETENSI.treatyGroup} ${i + 1}`}
                      onClick={() => {
                        jalankan('hapus', i)
                      }}
                    />
                  )
                }
              >
                <RincianBaris
                  b={b}
                  bisaUbah={bisaUbah}
                  mataUang={opsi.mataUang}
                  kelompokTreaty={opsi.kelompokTreaty}
                  onUbah={(baru) => {
                    setRows(rows.map((x, j) => (j === i ? baru : x)))
                  }}
                />
              </KartuLipat>
            ))}
          </div>
        )}

        {gagal !== '' && (
          <p className="tl-pesan" role="alert">
            {gagal}
          </p>
        )}

        <Bagian
          judul={TOTAL_RETENSI.judul}
          aksi={
            bisaUbah && (
              <div className="trin__aksi">
                {/* ⛔ MATI bila `TreatyIn.EDMMaterialType = 2` —
                    `pyDisabledWhen` @202558. Tombol kedua tepat di atasnya
                    di ekspor ber-`pyCondition 1=2`: ia MATI dan tidak
                    dibangun. Yang dibangun yang hidup. */}
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
              </div>
            )
          }
        >
          <GridTotal baris={total} />
        </Bagian>
      </div>
    </Panel>
  )
}
