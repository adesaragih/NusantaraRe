// Popup Copy Old Data (keputusan work owner 04-10-2026: "buat copy old data", tombol hanya untuk superadmin): berkas
// Pega yang header, detail, atau riwayatnya masih tertinggal di JSON lama. Setiap baris dapat dicentang; Process Copy
// menyalin yang dicentang - satu transaksi per berkas, isi yang sudah ada di tabel tidak ditimpa. Sesudah Process Copy
// daftar dibaca ulang: yang tersalin hilang dari popup. Bentuknya sama dengan Copy Old Company Detail.

import { useCallback, useEffect, useMemo, useRef, useState } from 'react'

import { Gagal, Kosong, Memuat, Modal } from '../../../../inti/frontend/components/ui/dasar'
import { ambilLama, salinLama, type BerkasLama, type HasilSalinLama, type StatusSalinLama } from '../api'
import { catatanLama, potong, ringkasSalin, saringLama, UKURAN_SALIN_LAMA } from '../aturan'
import { BDX, teksStatus } from '../labels'

const TEKS_STATUS: Record<StatusSalinLama, string> = {
  disalin: BDX.statusDisalin,
  sudahAda: BDX.statusSudahAda,
  ditolak: BDX.statusDitolak,
  gagal: BDX.statusGagal,
}

export default function DialogCopyOld({ onTutup }: { onTutup: (adaYangDisalin: boolean) => void }) {
  const [daftar, setDaftar] = useState<BerkasLama[] | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [kata, setKata] = useState('')
  const [terpilih, setTerpilih] = useState<ReadonlySet<string>>(new Set())
  const [proses, setProses] = useState(false)
  const [hasil, setHasil] = useState<HasilSalinLama[] | null>(null)
  /** Kemajuan Process Copy bertahap: ID yang sudah dikirim / seluruhnya. */
  const [kemajuan, setKemajuan] = useState<{ selesai: number; total: number } | null>(null)
  // Cacah tersalin sepanjang popup terbuka: daftar berkas dimuat ulang saat popup ditutup bila ada.
  const disalin = useRef(0)

  const muat = useCallback(async () => {
    try {
      const d = await ambilLama()
      setDaftar(d.daftar)
      setGalat(null)
      setTerpilih((t) => new Set(d.daftar.filter((x) => t.has(x.id)).map((x) => x.id)))
    } catch (e) {
      setGalat(e)
    }
  }, [])

  useEffect(() => {
    void muat()
  }, [muat])

  const tampil = useMemo(() => saringLama(daftar ?? [], kata), [daftar, kata])
  const semuaTercentang = tampil.length > 0 && tampil.every((d) => terpilih.has(d.id))
  const kirim = (daftar ?? []).filter((d) => terpilih.has(d.id)).map((d) => d.id)

  function centang(id: string, ya: boolean): void {
    setTerpilih((t) => {
      const b = new Set(t)
      if (ya) b.add(id)
      else b.delete(id)
      return b
    })
  }

  function centangSemua(ya: boolean): void {
    setTerpilih((t) => {
      const b = new Set(t)
      for (const d of tampil) {
        if (ya) b.add(d.id)
        else b.delete(d.id)
      }
      return b
    })
  }

  // Dikirim per UKURAN_SALIN_LAMA ID supaya setiap permintaan jauh di bawah batas waktunya. Hasil terkumpul per
  // kelompok; galat menghentikan sisanya.
  async function prosesCopy(): Promise<void> {
    if (proses) return
    setProses(true)
    setGalat(null)
    setHasil(null)
    const semua: HasilSalinLama[] = []
    let selesai = 0
    try {
      for (const kelompok of potong(kirim, UKURAN_SALIN_LAMA)) {
        setKemajuan({ selesai, total: kirim.length })
        const j = await salinLama(kelompok)
        disalin.current += j.disalin
        semua.push(...j.hasil)
        selesai += kelompok.length
        setHasil([...semua])
      }
    } catch (e) {
      setGalat(e)
    } finally {
      setKemajuan(null)
      await muat()
      setProses(false)
    }
  }

  const ringkas = hasil === null ? null : ringkasSalin(hasil)
  const bukanDisalin = (hasil ?? []).filter((h) => h.status !== 'disalin')

  return (
    <Modal
      judul={BDX.copyOld}
      lebar
      onTutup={() => {
        onTutup(disalin.current > 0)
      }}
      labelBatal={BDX.close}
      aksi={
        <button type="button" className="btn btn--primary" disabled={kirim.length === 0 || proses} onClick={() => void prosesCopy()}>
          {proses ? (kemajuan === null ? BDX.menyalin : BDX.menyalinKemajuan(kemajuan.selesai, kemajuan.total)) : BDX.prosesCopy}
        </button>
      }
    >
      <div className="bordereaux__lama">
        <p className="muted">{BDX.copyOldKeterangan}</p>
        {ringkas !== null && (
          <div className="bordereaux__lama-ringkas" role="status">
            {(Object.keys(TEKS_STATUS) as StatusSalinLama[])
              .filter((s) => ringkas[s] > 0)
              .map((s) => (
                <span key={s} className={s === 'disalin' ? 'bordereaux__lama-jumlah' : 'bordereaux__lama-jumlah bordereaux__lama-jumlah--masalah'}>
                  {ringkas[s]} {TEKS_STATUS[s]}
                </span>
              ))}
          </div>
        )}
        {bukanDisalin.length > 0 && (
          <ul className="bordereaux__lama-masalah">
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
        {daftar !== null && daftar.length === 0 && <Kosong pesan={BDX.copyOldKosong} />}
        {daftar !== null && daftar.length > 0 && (
          <>
            <div className="bordereaux__lama-alat">
              <input
                className="field__input bordereaux__lama-cari"
                type="search"
                placeholder={BDX.cariLama}
                aria-label={BDX.cariLama}
                value={kata}
                onChange={(e) => {
                  setKata(e.target.value)
                }}
              />
              <span aria-live="polite">
                {kirim.length} {BDX.dipilih}
              </span>
            </div>
            <div className="bordereaux__gulir bordereaux__lama-tabel">
              <table className="inbox__tabel bordereaux__tabel">
                <thead>
                  <tr>
                    <th className="bordereaux__lama-centang">
                      <input
                        type="checkbox"
                        aria-label={BDX.pilihSemua}
                        title={BDX.pilihSemua}
                        checked={semuaTercentang}
                        disabled={tampil.length === 0 || proses}
                        onChange={(e) => {
                          centangSemua(e.target.checked)
                        }}
                      />
                    </th>
                    <th>{BDX.id}</th>
                    <th>{BDX.type}</th>
                    <th>{BDX.business}</th>
                    <th>{BDX.cedingCo}</th>
                    <th>{BDX.treatyName}</th>
                    <th>{BDX.status}</th>
                    <th className="bordereaux__lama-catatan">{BDX.akanDisalin}</th>
                  </tr>
                </thead>
                <tbody>
                  {tampil.map((d) => (
                    <tr key={d.id} className="inbox__baris">
                      <td className="bordereaux__lama-centang">
                        <input
                          type="checkbox"
                          aria-label={`${BDX.pilihBaris} ${d.id}`}
                          checked={terpilih.has(d.id)}
                          disabled={proses}
                          onChange={(e) => {
                            centang(d.id, e.target.checked)
                          }}
                        />
                      </td>
                      <td>{d.id}</td>
                      <td>{d.type}</td>
                      <td>{d.business}</td>
                      <td>{d.ceding}</td>
                      <td>{d.treaty}</td>
                      <td>{teksStatus(d.status)}</td>
                      <td className="bordereaux__lama-catatan">{catatanLama(d)}</td>
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
