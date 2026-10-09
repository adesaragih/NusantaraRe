// Tab **Co-Ins Scale** — `Section/TreatyInTabsProportional.xml`.
//
// Dua saksi bentuknya, dan keduanya tangkapan layar Pega yang berjalan:
//   gambar 18 — mode lihat (tanpa tombol);
//   tangkapan layar pemakai 6 Oktober 2026 — mode ubah (tombol `Tambah`).
//
//   Co-Insurance Share │ % Treaty Limit │ [+ Tambah]     (TreatyIn.CoInScale)
//   No items
//   Max Co-Insurance Panel (Non Group)  [     ]          (sel 277)
//   Max Co-Insurance Panel (Group)      [     ]          (sel 278)
//
// ---------------------------------------------------------------------
// Yang dibaca dari XML, sel demi sel
// ---------------------------------------------------------------------
//   wadah   `pyContainerFormat`/`pyContainerType` = `NOHEADER` — TANPA judul
//           di atas grid. Bentuk sebelumnya memakai `Panel judul=…`.
//   grid    `pyGridNumbering` = false — tanpa kolom nomor baris (berbeda
//           dari Rate of Exchange); `pyHideGridHeaderWhenNoRows` = false —
//           kepala tetap tampil saat kosong; `pyFieldValueForNoRows` =
//           `GridNoResultsOnLoad` — teks bawaan Pega "No items".
//   sel 260/261  kepala kolom, `pyWidth` 206 · 206.
//   sel 262 tombol tambah DI SEL KEPALA KETIGA (`pyCellHeader` true,
//           `pyWidth` 101), `addRow`, `pyCondition` `TreatyIn.IsEditData!='1'`.
//   sel 264 `.CoInShare` `pxTextInput`; sel 265 `.PctLimit` `pxNumber`
//           ber-`pySymbol` `%` di KANAN. Keduanya baca-saja bila
//           `TreatyIn.IsEditData='1'` — jadi dapat disunting DI BARISNYA
//           pada mode ubah.
//   sel 266 tombol hapus per baris, `deleteRow`, syarat sama dengan 262.
//
// ⛔ SUNTING DI BARIS, BUKAN MODAL. Ekspor memuat `pyEditingDetailsType =
// modal` dan `pyNextGenRowEditing = masterDetail`, tetapi keduanya milik
// grid "next gen" — dan grid ini `pyNextGenGrid = false` dengan
// `pyRowEditing = row`. Grid Portfolio membawa pasangan nilai yang persis
// sama, dan Reporting Period (pola baca-saja-bersyarat yang sama) disunting
// di barisnya. Section modalnya (`TreatyInCoInScaleDetails!pyGridModalTemplate`)
// juga TIDAK ada di ekspor.
//
// ---------------------------------------------------------------------
// ⚠️ NILAI KEDUA MEDAN BELUM DAPAT DIBACA
// ---------------------------------------------------------------------
// `MaxCoNonGroup` terisi di 173 dan `MaxCoGroup` di 87 dari 1.855 dokumen
// `M_TREATY_IN` — tetapi tidak ada satu pun kolom pendaratan yang memuatnya
// (diperiksa di `ALL_TAB_COLUMNS` POOLDATA; `PROPORTIONALARRG.COINS_MAX`
// milik tabel induk aturan proporsional, BUKAN nilai kontrak).
// ⭐ 7 Oktober 2026: rumahnya DITETAPKAN — `T_TREATY_HAZARD_LIMIT`
// (`MAXCOGROUP`/`MAXCONONGROUP`, diagram v2 `TreatyIn [BATAS_BAHAYA]`,
// migrasi `446`). Belum terpasang, dan nilai kontrak lama hanya ada di
// dokumen JSON yang dilarang dibaca — jadi medannya tetap kosong.
// Catatan "belum terjangkau" yang dulu tampil di bawah medan DICABUT dari
// layar 6 Oktober 2026 atas permintaan pemakai ("samakan seperti yang ada
// di pega"); celahnya tercatat di sini dan di laporan pencocokan.
//
// ⛔ Masukan TIDAK diformat saat diketik — memformat di tiap ketukan pernah
// membuat koma desimal mustahil diketik (grid Rate of Exchange). Nilainya
// di keadaan layar saja sampai Save/Submit — `KEPUTUSAN-SASARAN-TULIS.md`.

import type { BarisSkalaKoasuransiWarisan } from '../api'
import { CO_INS_SCALE, JENIS_COIN_SCALE, KOLOM_COIN_SCALE } from '../labels'
import { useProperti } from '../halaman'
import type { ModeForm } from '../mode'
import { selAngka } from './angka'
import { TombolHapus, TombolTambah } from './limitsUI'

/**
 * Satu baris `TreatyIn.CoInScale` — `.CoInShare` dan `.PctLimit`, ejaan
 * Pega (isi penampung halaman, `halaman.tsx`).
 */
export interface BarisCoInScale {
  CoInShare: string
  PctLimit: string
}
type Baris = BarisCoInScale

/** Satu medan angka — `pxNumber` di ekspor, label di KIRI kotaknya. */
function MedanAngka({
  label,
  nilai,
  bisaUbah,
  onUbah,
}: {
  label: string
  nilai: string
  bisaUbah: boolean
  onUbah: (v: string) => void
}) {
  return (
    <div className="field">
      <label className="field__label">{label}</label>
      <input
        className="field__input"
        type="text"
        inputMode="decimal"
        aria-label={label}
        value={nilai}
        // ⛔ Baca-saja bila `TreatyIn.ViewState ='1'` — pyReadOnlyCondition
        // sel 277/278. `lihat` di layar ini.
        readOnly={!bisaUbah}
        onChange={(e) => {
          onUbah(e.target.value)
        }}
      />
    </div>
  )
}

export default function TabCoInsScale({
  baris,
  mode = 'lihat',
}: {
  baris: readonly BarisSkalaKoasuransiWarisan[]
  mode?: ModeForm
}) {
  // ⭐ Penampung halaman — `TreatyIn.CoInScale`, `MaxCoNonGroup`, `MaxCoGroup`
  // bertahan saat pindah tab.
  const [isi, setIsi] = useProperti<Baris[]>('CoInScale', () =>
    baris.map((b) => ({ CoInShare: b.bagianKoasuransi, PctLimit: b.persenLimit })),
  )
  const [nonGroup, setNonGroup] = useProperti('MaxCoNonGroup', '')
  const [group, setGroup] = useProperti('MaxCoGroup', '')
  const bisaUbah = mode === 'ubah'

  const ubah = (i: number, ruas: keyof Baris, v: string): void => {
    setIsi(isi.map((r, x) => (x === i ? { ...r, [ruas]: v } : r)))
  }

  return (
    // ⭐ KARTU TANPA JUDUL — wadah ekspor `NOHEADER` (bukan `Panel`, yang
    // selalu berjudul). Tanpa kartu, grid dan kedua medan melayang di latar
    // halaman tanpa jarak bawah dan menempel ke kartu Attachment (tangkapan
    // layar pemakai 7 Oktober 2026).
    <section className="panel trin__coin" aria-label="Co-Ins Scale">
      <div className="trin__coin-bungkus">
        <table className="trin__tabel">
          <thead>
            <tr>
              <th scope="col">{KOLOM_COIN_SCALE[0]}</th>
              <th scope="col">{KOLOM_COIN_SCALE[1]}</th>
              {/* Sel 262 — tombol tambah DUDUK di sel kepala ketiga. */}
              <th scope="col">
                {bisaUbah && (
                  <TombolTambah
                    label={CO_INS_SCALE.tambah}
                    onClick={() => {
                      // `addRow`, `pyPosition` After — baris baru di ujung.
                      setIsi([...isi, { CoInShare: '', PctLimit: '' }])
                    }}
                  />
                )}
              </th>
            </tr>
          </thead>
          <tbody>
            {isi.length === 0 && (
              <tr>
                <td colSpan={3} className="trin__coin-kosong">
                  {CO_INS_SCALE.kosong}
                </td>
              </tr>
            )}
            {isi.map((b, i) => (
              <tr key={i}>
                {/* Sel 264 — `.CoInShare`, `pxTextInput`. */}
                <td>
                  {bisaUbah ? (
                    <input
                      className="field__input"
                      type="text"
                      aria-label={KOLOM_COIN_SCALE[0]}
                      value={b.CoInShare}
                      onChange={(e) => {
                        ubah(i, 'CoInShare', e.target.value)
                      }}
                    />
                  ) : (
                    selAngka(JENIS_COIN_SCALE[0] ?? 'teks', b.CoInShare)
                  )}
                </td>
                {/* Sel 265 — `.PctLimit`, `pxNumber`, `%` di kanan. */}
                <td>
                  {bisaUbah ? (
                    <span className="trin__coin-persen">
                      <input
                        className="field__input"
                        type="text"
                        inputMode="decimal"
                        aria-label={KOLOM_COIN_SCALE[1]}
                        value={b.PctLimit}
                        onChange={(e) => {
                          ubah(i, 'PctLimit', e.target.value)
                        }}
                      />
                      <span aria-hidden="true">{CO_INS_SCALE.simbolPersen}</span>
                    </span>
                  ) : (
                    selAngka(JENIS_COIN_SCALE[1] ?? 'teks', b.PctLimit)
                  )}
                </td>
                {/* Sel 266 — `deleteRow`, syarat tampil sama dengan sel 262. */}
                <td>
                  {bisaUbah && (
                    <TombolHapus
                      label={CO_INS_SCALE.hapus}
                      labelAkses={`${CO_INS_SCALE.hapus} ${String(i + 1)}`}
                      onClick={() => {
                        setIsi(isi.filter((_, x) => x !== i))
                      }}
                    />
                  )}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
      {/* ⛔ TANPA judul panel kedua: di kedua tangkapan layar kedua medan
          duduk langsung di bawah grid, label di kiri kotaknya. */}
      <div className="trin__kolom trin__coin-medan">
        <MedanAngka label={CO_INS_SCALE.nonGroup} nilai={nonGroup} bisaUbah={bisaUbah} onUbah={setNonGroup} />
        <MedanAngka label={CO_INS_SCALE.group} nilai={group} bisaUbah={bisaUbah} onUbah={setGroup} />
      </div>
    </section>
  )
}
