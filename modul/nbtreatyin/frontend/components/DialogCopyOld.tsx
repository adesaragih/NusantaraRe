// Popup Copy Old NB Treaty In - perintah work owner 07-10-2026 ("nb ttreatyin tobol copy untuk data lama mana?" ->
// "langsung anda kerjakan!"); bukan layar Pega. Pola sama dengan Copy Old EDM Treaty In
// (`modul/edmtreatyin/frontend/components/DialogCopyOld.tsx`) dan Product Name Life. Popup memuat dokumen polis NB
// lama (JSON_POLIS generasi 0) yang belum ada di tabel flat; setiap baris dapat dicentang dan `Process Copy` menyalin
// yang dicentang (satu transaksi per dokumen, aturan pemuat dokumen lama). Baris yang tidak dapat disalin tampil
// dengan alasannya dan tidak dapat dicentang. Sesudah `Process Copy` daftar dibaca ulang: yang tersalin hilang dari
// popup dan tampil di tab Resolved portal.

import { useCallback, useEffect, useMemo, useRef, useState } from 'react'

import { Gagal, Kosong, Memuat, Modal } from '../../../../inti/frontend/components/ui/dasar'
import { ambilDokumenLama, salinDokumenLama, type DokumenLama, type HasilSalinLama, type StatusSalinLama } from '../api'
import { COPY_OLD } from '../labels'
import { idBolehDisalin, ringkasSalin, saringLama } from '../lama'
import { sajikan } from '../sajian'

const URUT_STATUS: StatusSalinLama[] = ['disalin', 'sudahAda', 'ditolak', 'gagal']

/** Tabel popup: kepala kolom mencentang / melepas SEMUA baris tampil yang dapat disalin. */
export function TabelLama({
  tampil,
  terpilih,
  proses,
  onCentang,
  onCentangSemua,
}: {
  tampil: readonly DokumenLama[]
  terpilih: ReadonlySet<string>
  proses: boolean
  onCentang: (id: string, ya: boolean) => void
  onCentangSemua: (ya: boolean) => void
}) {
  const boleh = tampil.filter((d) => d.bolehDisalin)
  const semua = boleh.length > 0 && boleh.every((d) => terpilih.has(d.id))
  const k = COPY_OLD.kolom
  return (
    <div className="table-wrap nbti__lama-tabel">
      <table>
        <thead>
          <tr>
            <th scope="col" className="nbti__lama-centang">
              <input
                type="checkbox"
                aria-label={COPY_OLD.pilihSemua}
                checked={semua}
                disabled={boleh.length === 0 || proses}
                onChange={(e) => onCentangSemua(e.target.checked)}
              />
            </th>
            <th scope="col">{k.id}</th>
            <th scope="col">{k.noOffer}</th>
            <th scope="col">{k.noPolis}</th>
            <th scope="col">{k.insured}</th>
            <th scope="col">{k.bisnis}</th>
            <th scope="col">{k.sob}</th>
            <th scope="col">{k.ceding}</th>
            <th scope="col">{k.tglProd}</th>
            <th scope="col">{k.catatan}</th>
          </tr>
        </thead>
        <tbody>
          {tampil.map((d, i) => (
            <tr key={`${d.id}-${i}`} className={d.bolehDisalin ? undefined : 'nbti__lama-ditolak'}>
              <td className="nbti__lama-centang">
                <input
                  type="checkbox"
                  aria-label={`${COPY_OLD.pilihBaris} ${d.id}`}
                  checked={terpilih.has(d.id)}
                  disabled={!d.bolehDisalin || proses}
                  onChange={(e) => onCentang(d.id, e.target.checked)}
                />
              </td>
              <td>{d.id}</td>
              <td>{d.noOffer}</td>
              <td>{d.noPolis}</td>
              <td>{d.insuredName}</td>
              <td>{d.businessName}</td>
              <td>{d.sobName}</td>
              <td>{d.cedingCoName}</td>
              <td>{sajikan(d.tglProd, 'tanggal')}</td>
              <td className="nbti__lama-catatan">
                {!d.bolehDisalin && (
                  <span className="nbti__lama-alasan">
                    {COPY_OLD.status.ditolak}: {d.alasan.join('; ')}
                  </span>
                )}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}

export default function DialogCopyOld({ onTutup }: { onTutup: (adaYangDisalin: boolean) => void }) {
  const [daftar, setDaftar] = useState<DokumenLama[] | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [kata, setKata] = useState('')
  const [terpilih, setTerpilih] = useState<ReadonlySet<string>>(new Set())
  const [proses, setProses] = useState(false)
  const [hasil, setHasil] = useState<HasilSalinLama[] | null>(null)
  // Cacah tersalin sepanjang popup terbuka: portal dimuat ulang saat popup ditutup bila ada.
  const disalin = useRef(0)

  const muat = useCallback(async () => {
    try {
      const d = await ambilDokumenLama()
      setDaftar(d)
      setGalat(null)
      // Pilihan yang tidak lagi ada di daftar (sudah tersalin) dibuang.
      setTerpilih((t) => new Set(d.filter((x) => t.has(x.id)).map((x) => x.id)))
    } catch (e) {
      setGalat(e)
    }
  }, [])

  useEffect(() => {
    void muat()
  }, [muat])

  const tampil = useMemo(() => saringLama(daftar ?? [], kata), [daftar, kata])
  const cacahKirim = idBolehDisalin(daftar ?? [], terpilih).length

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
      for (const d of tampil.filter((x) => x.bolehDisalin)) {
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
      const j = await salinDokumenLama(idBolehDisalin(daftar ?? [], terpilih))
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
      judul={COPY_OLD.tombol}
      lebar
      onTutup={() => onTutup(disalin.current > 0)}
      labelBatal={COPY_OLD.tutup}
      aksi={
        <button
          type="button"
          className="btn btn--primary"
          disabled={cacahKirim === 0 || proses}
          onClick={() => void prosesCopy()}
        >
          {COPY_OLD.prosesCopy}
        </button>
      }
    >
      <div className="nbti__lama">
        <p className="muted">{COPY_OLD.keterangan}</p>
        {ringkas !== null && (
          <div className="nbti__lama-ringkas" role="status">
            {URUT_STATUS.filter((s) => ringkas[s] > 0).map((s) => (
              <span key={s}>
                {ringkas[s]} {COPY_OLD.status[s]}
              </span>
            ))}
          </div>
        )}
        {bukanDisalin.length > 0 && (
          <ul className="nbti__lama-masalah">
            {bukanDisalin.map((h) => (
              <li key={h.id}>
                <strong>{h.id}</strong> {COPY_OLD.status[h.status]}
                {h.pesan.length > 0 && ` - ${h.pesan.join('; ')}`}
              </li>
            ))}
          </ul>
        )}
        {galat !== null && <Gagal galat={galat} />}
        {daftar === null && galat === null && <Memuat />}
        {daftar !== null && daftar.length === 0 && <Kosong pesan={COPY_OLD.kosong} />}
        {daftar !== null && daftar.length > 0 && (
          <>
            <div className="nbti__lama-alat">
              <input
                className="field__input nbti__lama-cari"
                type="search"
                placeholder={COPY_OLD.cari}
                aria-label={COPY_OLD.cari}
                value={kata}
                onChange={(e) => setKata(e.target.value)}
              />
              <span className="nbti__lama-cacah" aria-live="polite">
                {cacahKirim} {COPY_OLD.dipilih}
              </span>
            </div>
            <TabelLama
              tampil={tampil}
              terpilih={terpilih}
              proses={proses}
              onCentang={centang}
              onCentangSemua={centangSemua}
            />
          </>
        )}
      </div>
    </Modal>
  )
}
