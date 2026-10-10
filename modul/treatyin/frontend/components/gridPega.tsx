// BENTUK PEGA untuk tab Treaty In — grid, blok berjudul, dan sel label.
//
// ⭐ Permintaan pemakai 8 Oktober 2026: *"ubah semua design yang ada di
// treaty in dan treaty in adjustment, samakan designnya jangan ada bentuk
// aneh aneh, samakan posisi tiap property dan bentuk tablenya … jangan ada
// design tambahan"* — dijawab "Persis Pega", warna/huruf tema Treaty
// Exchange Yearly TETAP.
//
// Bentuknya meniru `GridEkspor` / `RenderKerangka` layar Treaty In Adjustment
// (`modul/treatyinadjustment/frontend/komponen/KerangkaTab.tsx`), yang
// dibangkitkan dari Section ekspor YANG SAMA (`TreatyInTabsNonProportional`,
// `TreatyInTabsProportional`). Dengan begitu kedua layar sama:
//   - kolom dan urutannya = kolom ekspor; lebar ekspor dipakai sebagai
//     perbandingan (persen), seperti Adjustment;
//   - grid `expandPane`/`masterDetail`: kolom ▸/▾ di kiri, rincian sebaris
//     penuh di bawah barisnya;
//   - tombol `Add` di sel KEPALA kolom tombol, `Delete` di sel baris;
//   - grid tanpa baris: satu sel `No items`.
//
// ⛔ Ditiru, BUKAN diimpor — modul tidak boleh saling impor.
// ⛔ Menggantikan kartu lipat bernomor, chip ringkasan, dan lencana jumlah
// (`limitsUI.tsx`, 6 Oktober 2026) — bentuk yang tidak ada di Pega.

import { Fragment, useState, type ReactNode } from 'react'

import { klikBaris, TombolNavigasi } from './navigasi'

/** Satu kolom data grid. */
export interface KolomPega<T> {
  /** Kepala kolom — teks ekspor apa adanya (`''` = kepala kosong). */
  judul: string
  /** Lebar ekspor (px) — perbandingan, bukan piksel. */
  lebar: number
  /** Angka rata kanan. */
  angka?: boolean
  isi: (b: T, i: number) => ReactNode
}

/** Kolom tombol — `Add` di kepala, `Delete` di baris. */
export interface TombolGridPega<T> {
  lebar: number
  /**
   * Posisi kolom tombol di antara kolom data (indeks sisip). Bawaan: paling
   * kanan. Grid Treaty Group Limits menaruhnya DI TENGAH (`Treaty Group` ·
   * tombol · `ROL Profile`), seperti ekspor.
   */
  sisip?: number
  tambah?: { label: string; onKlik: () => void; mati?: boolean }
  hapus?: { label: string; onKlik: (i: number) => void; akses?: (b: T, i: number) => string; mati?: boolean }
}

export function GridPega<T>({
  kolom,
  baris,
  rincian,
  tombol,
  labelBuka = 'Detail',
  kosong = 'No items',
  label,
  kelas,
  lebarTetap = false,
}: {
  kolom: readonly KolomPega<T>[]
  baris: readonly T[]
  /** Rincian baris (`expandPane`) — tertutup semula, seperti Pega. */
  rincian?: (b: T, i: number) => ReactNode
  tombol?: TombolGridPega<T>
  labelBuka?: string
  kosong?: string
  /** `aria-label` tabel. */
  label?: string
  /** Kelas tambahan tabel (mis. `trin__tabel--nilai`). */
  kelas?: string
  /**
   * Grid selebar jumlah `pyWidth`-nya, bukan selebar wadah — grid nilai
   * sempit Pega (`Total … | Value`).
   */
  lebarTetap?: boolean
}) {
  const [buka, setBuka] = useState<ReadonlySet<number>>(() => new Set())
  const t = tombol !== undefined && (tombol.tambah !== undefined || tombol.hapus !== undefined) ? tombol : undefined
  const total = kolom.reduce((a, k) => a + k.lebar, 0) + (t?.lebar ?? 0) || 1
  const persen = (w: number) => `${((w / total) * 100).toFixed(2)}%`
  const rentang = kolom.length + (t !== undefined ? 1 : 0) + (rincian !== undefined ? 1 : 0)
  // Urutan sel: indeks kolom data, dan `'tombol'` di posisi sisipnya.
  const urutan: (number | 'tombol')[] = kolom.map((_, i) => i)
  if (t !== undefined) urutan.splice(Math.min(Math.max(t.sisip ?? kolom.length, 0), kolom.length), 0, 'tombol')
  return (
    <div className="trin__grid" style={lebarTetap ? { maxWidth: `${String(total)}px` } : undefined}>
      <div className="table-wrap">
        <table
          className={
            'trin__tabel trin__tabel--pega' + (rincian !== undefined ? ' trin__tabel--rinci' : '') + (kelas !== undefined ? ` ${kelas}` : '')
          }
          aria-label={label}
        >
          <colgroup>
            {rincian !== undefined && <col className="trin__buka-kolom" />}
            {urutan.map((u) => (
              <col key={u} style={{ width: persen(u === 'tombol' ? (t?.lebar ?? 0) : (kolom[u]?.lebar ?? 0)) }} />
            ))}
          </colgroup>
          <thead>
            <tr>
              {rincian !== undefined && <th scope="col" aria-label={labelBuka} />}
              {urutan.map((u) =>
                u === 'tombol' ? (
                  <th key={u} scope="col">
                    {t?.tambah !== undefined && (
                      <button type="button" className="btn btn--sm" onClick={t.tambah.onKlik} disabled={t.tambah.mati}>
                        {t.tambah.label}
                      </button>
                    )}
                  </th>
                ) : (
                  <th key={u} scope="col" className={kolom[u]?.angka === true ? 'trin__angka' : undefined}>
                    {kolom[u]?.judul}
                  </th>
                ),
              )}
            </tr>
          </thead>
          <tbody>
            {baris.length === 0 && (
              <tr>
                <td colSpan={rentang} className="trin__kosong-pega">
                  {kosong}
                </td>
              </tr>
            )}
            {baris.map((b, r) => {
              const terbuka = buka.has(r)
              return (
                <Fragment key={r}>
                  <tr
                    className={rincian !== undefined ? 'trin__baris-buka' : undefined}
                    onClick={
                      rincian !== undefined
                        ? (e) => {
                            // ⭐ Klik baris = klik panah (`klikBaris`).
                            klikBaris(e, () => {
                              setBuka((x) => {
                                const y = new Set(x)
                                if (y.has(r)) y.delete(r)
                                else y.add(r)
                                return y
                              })
                            })
                          }
                        : undefined
                    }
                  >
                    {rincian !== undefined && (
                      <td className="trin__buka">
                        {/* Navigasi, bukan `<button>` — tetap dapat dibuka di mode lihat. */}
                        <TombolNavigasi
                          className="trin__buka-tombol"
                          label={`${labelBuka} ${String(r + 1)}`}
                          terbuka={terbuka}
                          onKlik={() => {
                            setBuka((x) => {
                              const y = new Set(x)
                              if (y.has(r)) y.delete(r)
                              else y.add(r)
                              return y
                            })
                          }}
                        >
                          {terbuka ? '▾' : '▸'}
                        </TombolNavigasi>
                      </td>
                    )}
                    {urutan.map((u) =>
                      u === 'tombol' ? (
                        <td key={u}>
                          {t?.hapus !== undefined && (
                            <button
                              type="button"
                              className="btn btn--sm"
                              aria-label={t.hapus.akses?.(b, r)}
                              disabled={t.hapus.mati}
                              onClick={() => {
                                // Indeks bergeser — rincian yang terbuka ditutup.
                                setBuka(new Set())
                                t.hapus?.onKlik(r)
                              }}
                            >
                              {t.hapus.label}
                            </button>
                          )}
                        </td>
                      ) : (
                        <td key={u} className={kolom[u]?.angka === true ? 'trin__angka' : undefined}>
                          {kolom[u]?.isi(b, r)}
                        </td>
                      ),
                    )}
                  </tr>
                  {rincian !== undefined && (
                    // ⛔ Rincian yang tertutup TETAP di DOM (`hidden`): sub-tab
                    // aktif dan isiannya tidak hilang saat dilipat.
                    <tr className="trin__rincian" hidden={!terbuka}>
                      <td colSpan={rentang}>
                        <div className="trin__blok">{rincian(b, r)}</div>
                      </td>
                    </tr>
                  )}
                </Fragment>
              )
            })}
          </tbody>
        </table>
      </div>
    </div>
  )
}

/**
 * Grid nilai baca-saja `{Currency, Value}` (`readOnly` ekspor) — kepala
 * kolom pertama = judul grid (mis. `Total EGNPI Amount`), kedua `Value`.
 */
export function GridNilaiPega({
  judul,
  nilai = 'Value',
  baris,
  tampil,
  lebar = [193, 349],
  lebarTetap = true,
}: {
  judul: string
  nilai?: string
  baris: readonly { Currency: string; Value: string }[]
  tampil: (v: string) => string
  lebar?: readonly [number, number]
  /** `false` — selebar wadahnya (wadah yang membatasi, mis. total Share Prop). */
  lebarTetap?: boolean
}) {
  return (
    <GridPega
      kolom={[
        { judul, lebar: lebar[0], isi: (b) => b.Currency },
        { judul: nilai, lebar: lebar[1], angka: true, isi: (b) => tampil(b.Value) },
      ]}
      baris={baris}
      kelas="trin__tabel--nilai"
      lebarTetap={lebarTetap}
    />
  )
}

/** Blok layout ekspor — judul (bila ada) lalu isi yang mengalir. */
export function BlokPega({ judul, children }: { judul?: string; children: ReactNode }) {
  return (
    <div className="trin__blok">
      {judul !== undefined && judul !== '' && <h5 className="trin__subjudul">{judul}</h5>}
      {children}
    </div>
  )
}

/** Sel label ekspor (`Part of`, `Total ROL`, `%`, `IDR`) — teks polos. */
export function TeksPega({ children }: { children: ReactNode }) {
  return <span className="trin__teks-sel">{children}</span>
}

/** Deret tombol ekspor di bawah blok (`Update Total`, `Update Value in List`). */
export function DeretTombolPega({ children }: { children: ReactNode }) {
  return <div className="trin__deret-tombol">{children}</div>
}
