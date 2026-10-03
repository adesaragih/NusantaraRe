// Popup `Copy Old` - permintaan work owner 03-10-2026: tombol di samping `Add`; popup memuat SEMUA produk tabel JSON
// lama yang belum ada di tabel flat, setiap baris dapat dicentang, dan `Process Copy` (kaki popup) menyalin yang
// dicentang ke tabel flat. Server menyalin per produk (satu transaksi per produk) dengan aturan alat pindah; baris yang
// tidak dapat disalin tampil dengan alasannya dan tidak dapat dicentang. Sesudah `Process Copy` daftar dibaca ulang:
// yang tersalin hilang dari popup dan tampil di grid produk.

import { useCallback, useEffect, useMemo, useRef, useState } from 'react'

import { Gagal, Kosong, Memuat, Modal } from '../../../../inti/frontend/components/ui/dasar'
import { ambilProdukLama, salinProdukLama, type HasilSalinLama, type ProdukLama, type StatusSalinLama } from '../api'
import { idBolehDisalin, ringkasSalin, saringLama } from '../bentuk'
import { GRID_MPNL, LAIN_MPNL, PEMILIH_MPNL, TOMBOL_MPNL, UMUM_MPNL } from '../labels'

const TEKS_STATUS: Record<StatusSalinLama, string> = {
  disalin: LAIN_MPNL.statusDisalin,
  sudahAda: LAIN_MPNL.statusSudahAda,
  ditolak: LAIN_MPNL.statusDitolak,
  gagal: LAIN_MPNL.statusGagal,
}

export default function DialogCopyOld({ onTutup }: { onTutup: (adaYangDisalin: boolean) => void }) {
  const [daftar, setDaftar] = useState<ProdukLama[] | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [kata, setKata] = useState('')
  const [terpilih, setTerpilih] = useState<ReadonlySet<string>>(new Set())
  const [proses, setProses] = useState(false)
  const [hasil, setHasil] = useState<HasilSalinLama[] | null>(null)
  // Cacah tersalin sepanjang popup terbuka: grid produk dimuat ulang saat popup ditutup bila ada.
  const disalin = useRef(0)

  const muat = useCallback(async () => {
    try {
      const d = await ambilProdukLama()
      setDaftar(d.daftar)
      setGalat(null)
      // Pilihan yang tidak lagi ada di daftar (sudah tersalin) dibuang.
      setTerpilih((t) => new Set(d.daftar.filter((x) => t.has(x.id)).map((x) => x.id)))
    } catch (e) {
      setGalat(e)
    }
  }, [])

  useEffect(() => {
    void muat()
  }, [muat])

  const tampil = useMemo(() => saringLama(daftar ?? [], kata), [daftar, kata])
  const bolehTampil = tampil.filter((d) => d.bolehDisalin)
  const semuaTercentang = bolehTampil.length > 0 && bolehTampil.every((d) => terpilih.has(d.id))
  const cacahKirim = idBolehDisalin(daftar ?? [], terpilih).length

  function centang(id: string, ya: boolean): void {
    setTerpilih((t) => {
      const b = new Set(t)
      if (ya) b.add(id)
      else b.delete(id)
      return b
    })
  }

  /** Kepala kolom: centang / lepas SEMUA baris yang tampil (hasil `Search`) dan dapat disalin. */
  function centangSemua(ya: boolean): void {
    setTerpilih((t) => {
      const b = new Set(t)
      for (const d of bolehTampil) {
        if (ya) b.add(d.id)
        else b.delete(d.id)
      }
      return b
    })
  }

  async function prosesCopy(): Promise<void> {
    if (proses) return
    setProses(true)
    setGalat(null)
    try {
      const j = await salinProdukLama(idBolehDisalin(daftar ?? [], terpilih))
      disalin.current += j.disalin
      setHasil(j.hasil)
      await muat()
    } catch (e) {
      setGalat(e)
    } finally {
      setProses(false)
    }
  }

  const ringkas = hasil === null ? null : ringkasSalin(hasil)
  const bukanDisalin = (hasil ?? []).filter((h) => h.status !== 'disalin')

  return (
    <Modal
      judul={LAIN_MPNL.copyOld}
      lebar
      onTutup={() => {
        onTutup(disalin.current > 0)
      }}
      labelBatal={TOMBOL_MPNL.close}
      aksi={
        <button type="button" className="btn btn--primary" disabled={cacahKirim === 0 || proses} onClick={() => void prosesCopy()}>
          {LAIN_MPNL.prosesCopy}
        </button>
      }
    >
      <div className="mpnl-lama">
        <p className="muted mpnl-lama__keterangan">{LAIN_MPNL.copyOldKeterangan}</p>
        {ringkas !== null && (
          <div className="mpnl-lama__ringkas" role="status">
            {(Object.keys(TEKS_STATUS) as StatusSalinLama[])
              .filter((s) => ringkas[s] > 0)
              .map((s) => (
                <span key={s}>
                  {ringkas[s]} {TEKS_STATUS[s]}
                </span>
              ))}
          </div>
        )}
        {bukanDisalin.length > 0 && (
          <ul className="mpnl-lama__masalah">
            {bukanDisalin.map((h) => (
              <li key={h.id}>
                <strong>{h.id}</strong> {TEKS_STATUS[h.status]}
                {h.pesan.length > 0 && ` - ${h.pesan.join('; ')}`}
              </li>
            ))}
          </ul>
        )}
        {galat !== null && <Gagal galat={galat} />}
        {daftar === null && galat === null && <Memuat />}
        {daftar !== null && daftar.length === 0 && <Kosong pesan={LAIN_MPNL.copyOldKosong} />}
        {daftar !== null && daftar.length > 0 && (
          <>
            <div className="mpnl-lama__alat">
              <input
                className="field__input mpnl-lama__cari"
                type="search"
                placeholder={PEMILIH_MPNL.search}
                aria-label={PEMILIH_MPNL.search}
                value={kata}
                onChange={(e) => {
                  setKata(e.target.value)
                }}
              />
              <span className="mpnl-lama__cacah" aria-live="polite">
                {cacahKirim} {LAIN_MPNL.dipilih}
              </span>
            </div>
            <div className="mpnl-tabel">
              <table className="inbox__tabel">
                <thead>
                  <tr>
                    <th className="mpnl-lama__centang">
                      <input
                        type="checkbox"
                        aria-label={LAIN_MPNL.pilihSemua}
                        checked={semuaTercentang}
                        disabled={bolehTampil.length === 0 || proses}
                        onChange={(e) => {
                          centangSemua(e.target.checked)
                        }}
                      />
                    </th>
                    <th>{GRID_MPNL.kolomId}</th>
                    <th>{UMUM_MPNL.productName}</th>
                    <th>{GRID_MPNL.kolomCeding}</th>
                    <th>{GRID_MPNL.kolomTreatyNumber}</th>
                    <th>{GRID_MPNL.kolomTreatyName}</th>
                    <th>{GRID_MPNL.kolomCreateOp}</th>
                    <th>{LAIN_MPNL.kolomCatatan}</th>
                  </tr>
                </thead>
                <tbody>
                  {tampil.map((d) => (
                    <tr key={d.id} className={d.bolehDisalin ? 'inbox__baris' : 'inbox__baris mpnl-lama__ditolak'}>
                      <td className="mpnl-lama__centang">
                        <input
                          type="checkbox"
                          aria-label={`${LAIN_MPNL.pilihBaris} ${d.id}`}
                          checked={terpilih.has(d.id)}
                          disabled={!d.bolehDisalin}
                          onChange={(e) => {
                            centang(d.id, e.target.checked)
                          }}
                        />
                      </td>
                      <td>{d.id}</td>
                      <td>{d.productName}</td>
                      <td>{d.ceding}</td>
                      <td>{d.treatyNumber}</td>
                      <td>{d.inwardName}</td>
                      <td>{d.createOp}</td>
                      <td className="mpnl-lama__catatan">
                        {d.bolehDisalin ? (
                          d.catatan.join('; ')
                        ) : (
                          <span className="mpnl-status mpnl-status--gagal">
                            {TEKS_STATUS.ditolak}: {d.alasan.join('; ')}
                          </span>
                        )}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </>
        )}
      </div>
    </Modal>
  )
}
