// Popup `Copy Old` (perintah work owner 04-10-2026: "buatkan fungsinya copy old seperti pada productname; copy old hanya
// muncul untuk superadmin"): organisasi dokumen Pega yang datanya belum pindah ke tabel datar; setiap baris dapat
// dicentang, dan `Process Copy` menyalin yang dicentang - satu transaksi per organisasi, aturan alat pindah. Sesudah
// Process Copy daftar dibaca ulang: yang tersalin hilang dari popup. Bentuknya sama dengan Copy Old Product Name.

import { useCallback, useEffect, useMemo, useRef, useState } from 'react'

import { Gagal, Kosong, Memuat, Modal } from '../../../../inti/frontend/components/ui/dasar'
import { ambilLama, salinLama, type HasilSalinLama, type OrgLama, type StatusSalinLama } from '../api'
import { catatanLama, potong, ringkasSalin, saringLama, UKURAN_SALIN_LAMA } from '../aturan'
import { CD } from '../labels'

const TEKS_STATUS: Record<StatusSalinLama, string> = {
  disalin: CD.statusDisalin,
  sudahAda: CD.statusSudahAda,
  ditolak: CD.statusDitolak,
  gagal: CD.statusGagal,
}

export default function DialogCopyOld({ onTutup }: { onTutup: (adaYangDisalin: boolean) => void }) {
  const [daftar, setDaftar] = useState<OrgLama[] | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [kata, setKata] = useState('')
  const [terpilih, setTerpilih] = useState<ReadonlySet<string>>(new Set())
  const [proses, setProses] = useState(false)
  const [hasil, setHasil] = useState<HasilSalinLama[] | null>(null)
  /** Kemajuan Process Copy bertahap: ID yang sudah dikirim / seluruhnya. */
  const [kemajuan, setKemajuan] = useState<{ selesai: number; total: number } | null>(null)
  // Cacah tersalin sepanjang popup terbuka: daftar organisasi dimuat ulang saat popup ditutup bila ada.
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

  // Dikirim per UKURAN_SALIN_LAMA ID: satu permintaan untuk ratusan organisasi melewati batas waktu layar
  // ("signal is aborted without reason", 04-10-2026). Hasil terkumpul per kelompok; galat menghentikan sisanya.
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
      judul={CD.copyOld}
      lebar
      onTutup={() => {
        onTutup(disalin.current > 0)
      }}
      labelBatal={CD.tutup}
      aksi={
        <button type="button" className="btn btn--primary" disabled={kirim.length === 0 || proses} onClick={() => void prosesCopy()}>
          {proses ? (kemajuan === null ? CD.menyalin : CD.menyalinKemajuan(kemajuan.selesai, kemajuan.total)) : CD.prosesCopy}
        </button>
      }
    >
      <div className="companydetail__lama">
        <p className="muted">{CD.copyOldKeterangan}</p>
        {ringkas !== null && (
          <div className="companydetail__lama-ringkas" role="status">
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
          <ul className="companydetail__lama-masalah">
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
        {daftar !== null && daftar.length === 0 && <Kosong pesan={CD.copyOldKosong} />}
        {daftar !== null && daftar.length > 0 && (
          <>
            <div className="companydetail__lama-alat">
              <input
                className="field__input companydetail__cari"
                type="search"
                placeholder={CD.cari}
                aria-label={CD.cari}
                value={kata}
                onChange={(e) => {
                  setKata(e.target.value)
                }}
              />
              <label className="companydetail__centang">
                <input
                  type="checkbox"
                  checked={semuaTercentang}
                  disabled={tampil.length === 0 || proses}
                  onChange={(e) => {
                    centangSemua(e.target.checked)
                  }}
                />
                {CD.pilihSemua}
              </label>
              <span aria-live="polite">
                {kirim.length} {CD.dipilih}
              </span>
            </div>
            <div className="companydetail__gulir companydetail__lama-tabel">
              <table className="inbox__tabel companydetail__tabel">
                <thead>
                  <tr>
                    <th className="companydetail__lama-centang">
                      <input
                        type="checkbox"
                        aria-label={CD.pilihSemua}
                        checked={semuaTercentang}
                        disabled={tampil.length === 0 || proses}
                        onChange={(e) => {
                          centangSemua(e.target.checked)
                        }}
                      />
                    </th>
                    <th>{CD.kolomId}</th>
                    <th>{CD.organizationName}</th>
                    <th>{CD.kolomCatatan}</th>
                  </tr>
                </thead>
                <tbody>
                  {tampil.map((d) => (
                    <tr key={d.id} className="inbox__baris">
                      <td className="companydetail__lama-centang">
                        <input
                          type="checkbox"
                          aria-label={`${CD.pilihBaris} ${d.idView}`}
                          checked={terpilih.has(d.id)}
                          disabled={proses}
                          onChange={(e) => {
                            centang(d.id, e.target.checked)
                          }}
                        />
                      </td>
                      <td>{d.idView}</td>
                      <td>{d.nama}</td>
                      <td>{catatanLama(d)}</td>
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
