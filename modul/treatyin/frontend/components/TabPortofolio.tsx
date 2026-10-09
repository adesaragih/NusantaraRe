// Tab **Portfolio** — grid `TreatyIn.Portfolio`,
// `Section/TreatyInTabsProportional.xml`.
//
// ⛔ DIPISAHKAN dari `TabGridWarisan` 6 Oktober 2026. Grid umum itu
// merender SETIAP sel sebagai kotak teks, dan tab ini tidak begitu di
// Pega: dua kolom pertamanya daftar pilihan, yang ketiga area teks. Grid
// umum tetap dipakai tab lain; yang berubah hanya tab ini.
//
// ---------------------------------------------------------------------
// BENTUKNYA DARI EKSPOR, SEL DEMI SEL
// ---------------------------------------------------------------------
// Layout `TreatyIn.Portfolio`, `pyPageListPropertyClass` =
// `ASM-FW-GISFW-Data-TreatyInPortfolio`, dua baris sel yang berpasangan
// lurus:
//
//   kolom  judul            sel   properti          sel   kendali
//     1    Portfolio Type   107   .TypePortfolio    112   warisan properti
//     2    Premium/LossType 108   .Type             113   warisan properti
//     3    Description      109   .Description      114   pxTextArea
//     4    Add              110   Delete            115   pxButton
//
// ⚠️ Kolom 1 memakai `TypePortfolio` dan kolom 2 memakai `Type` — bukan
// sebaliknya. Layar ini pernah menukarnya; alasan lengkap dan ketiga
// saksinya ada di `labels.ts` pada `KOLOM_PORTOFOLIO`.
//
// ---------------------------------------------------------------------
// SATU SAKELAR MENGATUR SELURUHNYA: `TreatyIn.IsEditData`
// ---------------------------------------------------------------------
//   sel 112·113·114  `pyReadOnlyCondition` = `TreatyIn.IsEditData= 1`
//   sel 110·115      `pyCondition`         = `TreatyIn.IsEditData!='1'`
//
// Satu tanda, dua arah: ketika baris menjadi baca-saja, tombolnya hilang;
// ketika baris dapat diubah, tombolnya muncul. Itu persis `mode` di layar
// ini — `ubah` ⇔ `IsEditData != '1'`, `lihat` ⇔ `IsEditData = '1'`.
//
// ⚠️ Tab ini memakai `IsEditData`, sementara grid lain di seksi yang sama
// memakai `TreatyIn.ViewState` (terukur di ekspor: 9 sel `ViewState`, 3
// sel `IsEditData`). Keduanya dipetakan ke `mode` yang sama di sini, dan
// perbedaan itu dicatat supaya tidak terbaca sebagai kelalaian.
//
// ⛔ YANG TIDAK DIBAWA, dan sebabnya: `pyDisabledWhen` =
// `TreatyIn.EDMMaterialType = 2` ada pada keempat sel. Medan
// `EDMMaterialType` belum terbaca di layar ini — tidak ada di
// `KontrakWarisan` — jadi menyalakannya berarti menebak nilainya. Ia
// dicatat di sini, tidak dikarang di kode.

import { Kosong, Panel } from '../../../../inti/frontend/components/ui/dasar'
import type { BarisPortofolioWarisan } from '../api'
import { FORM_KONTRAK, GRID_TAMBAH, KOLOM_PORTOFOLIO, PORTOFOLIO } from '../labels'
import { useProperti } from '../halaman'
import type { ModeForm } from '../mode'

/**
 * Satu baris `TreatyIn.Portfolio` — EJAAN PEGA, sebab ia isi penampung
 * halaman (`halaman.tsx`) yang kelak dipetakan Save.
 */
export type BarisPortfolio = {
  /** Kolom 1 `Portfolio Type`. */
  TypePortfolio: string
  /** Kolom 2 `Premium / Loss Type`. */
  Type: string
  /** Kolom 3 `Description`. */
  Description: string
}
type Baris = BarisPortfolio

/** Daftar pilihan satu sel — `<select>` telanjang, bukan `Pilih`. */
function SelPilih({
  label,
  nilai,
  opsi,
  onUbah,
}: {
  label: string
  nilai: string
  opsi: readonly string[]
  onUbah: (v: string) => void
}) {
  // ⚠️ Nilai yang TIDAK ada di daftar tetap ditawarkan, di urutan paling
  // bawah. Himpunan pilihan di sini DIUKUR dari data, bukan dibaca dari
  // aturan properti (yang tidak ikut diekspor) — jadi nilai yang sah
  // tetapi belum pernah terpakai mungkin ada. Menjatuhkannya diam-diam
  // akan MENGUBAH kontrak yang dibuka tanpa seorang pun menyentuhnya.
  const asing = nilai !== '' && !opsi.includes(nilai)
  return (
    <select
      className="field__input"
      aria-label={label}
      value={nilai}
      onChange={(e) => {
        onUbah(e.target.value)
      }}
    >
      <option value="">{PORTOFOLIO.belumDipilih}</option>
      {opsi.map((o) => (
        <option key={o} value={o}>
          {o}
        </option>
      ))}
      {asing && (
        <option key={nilai} value={nilai}>
          {nilai}
        </option>
      )}
    </select>
  )
}

/**
 * Tab Portfolio.
 *
 * ⚠️ Baris disalin ke keadaan SEKALI saat tab lahir — sama dengan
 * `TabGridWarisan`. Pemanggilnya memberi `key` per kontrak supaya kontrak
 * lain melahirkan tab baru.
 *
 * ⛔ Perubahan di sini BELUM punya jalur simpan, persis seperti grid
 * lain, dan itu dikatakan di layar lewat `PORTOFOLIO`/`GRID_TAMBAH`
 * petunjuknya — bukan disembunyikan.
 */
export default function TabPortofolio({
  baris: barisWarisan,
  mode = 'lihat',
  petunjukKosong,
}: {
  baris: readonly BarisPortofolioWarisan[]
  mode?: ModeForm
  petunjukKosong: string
}) {
  // ⭐ Page list `TreatyIn.Portfolio` di penampung halaman — bertahan saat
  // pindah tab. Nilai awalnya baris kontrak yang dimuat.
  const [isi, setIsi] = useProperti<Baris[]>('Portfolio', () =>
    barisWarisan.map((b) => ({
      TypePortfolio: b.jenisPortfolio,
      Type: b.jenis,
      Description: b.keterangan,
    })),
  )
  const bisaUbah = mode === 'ubah'

  const ubah = (i: number, ruas: keyof Baris, v: string): void => {
    setIsi(isi.map((r, x) => (x === i ? { ...r, [ruas]: v } : r)))
  }

  return (
    <Panel judul={PORTOFOLIO.judul}>
      {/* Sel 110 — `pyCondition` = `TreatyIn.IsEditData!='1'`. */}
      {bisaUbah && (
        <div className="trin__panel-kepala">
          <span className="trin__redup">{GRID_TAMBAH.petunjuk}</span>
          <button
            type="button"
            className="btn btn--ghost btn--sm"
            onClick={() => {
              // ⛔ Baris baru KOSONG bertiga. `Activity/TreatyInPropAdd.xml`
              // hanya menetapkan `.Description = ""`; kedua pilihannya
              // memang lahir belum terisi.
              setIsi([...isi, { TypePortfolio: '', Type: '', Description: '' }])
            }}
          >
            {PORTOFOLIO.tambah}
          </button>
        </div>
      )}
      <div className="table-wrap">
        <table className="trin__tabel">
          <thead>
            <tr>
              {KOLOM_PORTOFOLIO.map((k) => (
                <th key={k} scope="col">
                  {k}
                </th>
              ))}
              {bisaUbah && <th scope="col" aria-label={PORTOFOLIO.hapus} />}
            </tr>
          </thead>
          <tbody>
            {isi.length === 0 && (
              <tr>
                <td colSpan={KOLOM_PORTOFOLIO.length + (bisaUbah ? 1 : 0)}>
                  <Kosong pesan={FORM_KONTRAK.tanpaBaris} petunjuk={petunjukKosong} />
                </td>
              </tr>
            )}
            {isi.map((b, i) => (
              <tr key={i}>
                {/* Sel 112 — kolom 1, `.TypePortfolio`. */}
                <td>
                  {bisaUbah ? (
                    <SelPilih
                      label={KOLOM_PORTOFOLIO[0]}
                      nilai={b.TypePortfolio}
                      opsi={PORTOFOLIO.opsiArah}
                      onUbah={(v) => {
                        ubah(i, 'TypePortfolio', v)
                      }}
                    />
                  ) : (
                    b.TypePortfolio
                  )}
                </td>
                {/* Sel 113 — kolom 2, `.Type`. */}
                <td>
                  {bisaUbah ? (
                    <SelPilih
                      label={KOLOM_PORTOFOLIO[1]}
                      nilai={b.Type}
                      opsi={PORTOFOLIO.opsiJenis}
                      onUbah={(v) => {
                        ubah(i, 'Type', v)
                      }}
                    />
                  ) : (
                    b.Type
                  )}
                </td>
                {/* Sel 114 — `pyFormat` = `pxTextArea`, bukan kotak sebaris. */}
                <td>
                  {bisaUbah ? (
                    <textarea
                      className="field__input"
                      rows={2}
                      aria-label={KOLOM_PORTOFOLIO[2]}
                      value={b.Description}
                      onChange={(e) => {
                        ubah(i, 'Description', e.target.value)
                      }}
                    />
                  ) : (
                    b.Description
                  )}
                </td>
                {/* Sel 115 — `Embed-SelectedContextAPI-DeleteRow`. */}
                {bisaUbah && (
                  <td>
                    <button
                      type="button"
                      className="btn btn--ghost btn--sm"
                      onClick={() => {
                        setIsi(isi.filter((_, x) => x !== i))
                      }}
                    >
                      {PORTOFOLIO.hapus}
                    </button>
                  </td>
                )}
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </Panel>
  )
}
