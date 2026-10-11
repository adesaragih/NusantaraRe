// Blok "Attachment File" case NB FacIn (tiket 67) - port `NB FacIn\Section\AttachmentGridReas.xml`
// (`ASM-FW-GISFW-DATA-FACOFFER!ATTACHMENTGRIDREAS`; Pega menampilkannya di popup Preview R/I Slip, di sini di layar case di bawah
// riwayat akseptasi - work owner 08-10-2026 "dipega ada add attachment … tambahin dong"):
// - grid Category · Count · Upload File · View File; kategori = `CATEGORY_ATTACH_REAS` (order by NOTE), Count = jumlah dokumen
//   case per kategori (`DOCUMENT_POLIS` KATEGORI_1 "Reas", KATEGORI_2 = kategori);
// - Upload File (`AttachContentGIS` / `AttachFileGIS`, banyak berkas sekaligus): berkas dikirim satu per satu (multipart, pola
//   Treaty Contract Out), disimpan ke Google Storage (ADR-0010) + `T_STORAGE_IMAGE` + `DOCUMENT_POLIS` di backend;
// - View File (harness `ViewPictureList`): File · Note · Upload Date · Delete; klik nama = unduh; Delete hanya bila
//   `bolehHapus` (persis Pega: FlagOnGoingPolicy 0 atau IT) - hapus FISIK seperti Pega (pengecualian sadar ADR-0031, work
//   owner 08-10-2026 "ikutin pega aja"); status unggah asinkron (ADR-0008) menunggu / gagal + Ulangi.
// Keputusan work owner 08-10-2026: tabel Pega DOCUMENT_POLIS + T_STORAGE_IMAGE; upload di semua tahap oleh pemegang workbasket.

import { useCallback, useEffect, useRef, useState, type ChangeEvent } from 'react'

import { Gagal, Modal } from '../../../../inti/frontend/components/ui/dasar'
import { unduhBerkasBeridentitas } from '../../../../inti/frontend/klien'
import {
  ambilLampiran,
  daftarLampiran,
  hapusLampiran,
  jalurIsiLampiran,
  ulangiLampiran,
  unggahLampiran,
  type BarisLampiran,
  type TampilanLampiran,
} from '../api'
import { LAMPIRAN as L, TEKS_LAMPIRAN as T } from '../labels'
import { tanggalRiwayat } from './BlokEmailType'
import HalamanPopup from './HalamanPopup'
import PenampilLampiran, { jenisPenampil } from './PenampilLampiran'

/** Kategori per halaman (work owner 08-10-2026: "buatin page, 5 per page dan rapikan biar user friendly"). */
export const UKURAN_KATEGORI_LAMPIRAN = 5

/** Ikon panah unggah / unduh (Pega: tombol ikon tanpa teks); teks tombol ada di aria-label / title. */
function Ikon({ arah }: { arah: 'naik' | 'turun' }) {
  return (
    <svg viewBox="0 0 24 24" width="16" height="16" aria-hidden="true" focusable="false">
      <path
        d={arah === 'naik' ? 'M12 16V4m0 0-5 5m5-5 5 5M5 20h14' : 'M12 4v12m0 0-5-5m5 5 5-5M5 20h14'}
        fill="none"
        stroke="currentColor"
        strokeWidth="2"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
    </svg>
  )
}

/** `awal` = tampilan awal (uji / pratinjau); tanpa itu dimuat dari server. */
export default function BlokLampiran({ caseId, awal }: { caseId: string; awal?: TampilanLampiran }) {
  const [v, setV] = useState<TampilanLampiran | null>(awal ?? null)
  const [galat, setGalat] = useState<unknown>(null)
  const [pesan, setPesan] = useState('')
  const [mengunggah, setMengunggah] = useState('')
  const [lihat, setLihat] = useState<string | null>(null)
  const [berkasKategori, setBerkasKategori] = useState<BarisLampiran[] | null>(null)
  const [hapus, setHapus] = useState<BarisLampiran | null>(null)
  // Penampil View / View Office Online (work owner 11-10-2026) dan galatnya - tampil di popup View File.
  const [penampil, setPenampil] = useState<BarisLampiran | null>(null)
  const [galatLihat, setGalatLihat] = useState<unknown>(null)
  const [hal, setHal] = useState(1)
  const pilihBerkas = useRef<HTMLInputElement>(null)
  const kategoriUnggah = useRef('')

  const muat = useCallback(() => {
    ambilLampiran(caseId).then(setV, (err: unknown) => setGalat(err))
  }, [caseId])

  useEffect(() => {
    muat()
  }, [muat])

  async function muatBerkas(kategori: string) {
    setBerkasKategori(null)
    try {
      setBerkasKategori((await daftarLampiran(caseId, kategori)).baris)
    } catch (err) {
      setGalat(err)
    }
  }

  function bukaUnggah(kategori: string) {
    kategoriUnggah.current = kategori
    pilihBerkas.current?.click()
  }

  /** Banyak berkas sekaligus (Pega drag-drop multi-berkas): dikirim berurutan, berhenti pada galat pertama. */
  async function unggah(e: ChangeEvent<HTMLInputElement>) {
    const daftar = Array.from(e.target.files ?? [])
    e.target.value = ''
    const kategori = kategoriUnggah.current
    if (daftar.length === 0 || kategori === '') return
    setGalat(null)
    setPesan('')
    let n = 0
    try {
      for (const f of daftar) {
        setMengunggah(kategori)
        await unggahLampiran(caseId, kategori, f)
        n++
      }
      setPesan(T.terunggah(n, kategori))
    } catch (err) {
      setGalat(err)
      if (n > 0) setPesan(T.terunggah(n, kategori))
    } finally {
      setMengunggah('')
      muat()
      if (lihat === kategori) void muatBerkas(kategori)
    }
  }

  async function ulangi(b: BarisLampiran) {
    if (lihat === null) return
    try {
      await ulangiLampiran(caseId, b.id)
      await muatBerkas(lihat)
    } catch (err) {
      setGalat(err)
    }
  }

  async function lakukanHapus() {
    if (hapus === null || lihat === null) return
    const kategori = lihat
    setHapus(null)
    try {
      await hapusLampiran(caseId, hapus.id)
      muat()
      await muatBerkas(kategori)
    } catch (err) {
      setGalat(err)
    }
  }

  if (v === null) return galat ? <Gagal galat={galat} /> : null
  const totalHal = Math.max(1, Math.ceil(v.kategori.length / UKURAN_KATEGORI_LAMPIRAN))
  const kini = Math.min(Math.max(1, hal), totalHal)
  const halamanIni = v.kategori.slice((kini - 1) * UKURAN_KATEGORI_LAMPIRAN, kini * UKURAN_KATEGORI_LAMPIRAN)

  return (
    <section className="panel nbf-lampiran">
      <h3 className="nbf-lampiran__judul">{L.judul}</h3>
      <Gagal galat={galat} />
      {pesan !== '' && <div className="alert alert--ok">{pesan}</div>}
      <input ref={pilihBerkas} type="file" multiple hidden aria-label={L.unggah} onChange={(e) => void unggah(e)} />
      {v.kategori.length > UKURAN_KATEGORI_LAMPIRAN && (
        <HalamanPopup halaman={kini} ukuran={UKURAN_KATEGORI_LAMPIRAN} total={v.kategori.length} onPindah={setHal} />
      )}
      <div className="nbf-cov-wrap">
        <table className="nbf-tabel nbf-lampiran__tabel">
          <colgroup>
            <col />
            <col className="nbf-lampiran__kolom-jumlah" />
            <col className="nbf-lampiran__kolom-aksi" />
            <col className="nbf-lampiran__kolom-aksi" />
          </colgroup>
          <thead>
            <tr>
              <th scope="col">{L.kategori}</th>
              <th scope="col" className="nbf-lampiran__tengah">
                {L.jumlah}
              </th>
              <th scope="col" className="nbf-lampiran__tengah">
                {L.unggah}
              </th>
              <th scope="col" className="nbf-lampiran__tengah">
                {L.lihat}
              </th>
            </tr>
          </thead>
          <tbody>
            {v.kategori.length === 0 && (
              <tr>
                <td colSpan={4}>{T.kosong}</td>
              </tr>
            )}
            {halamanIni.map((k) => (
              <tr key={k.kategori}>
                <td className="nbf-lampiran__nama">{k.kategori}</td>
                <td className="nbf-lampiran__tengah">
                  <span className="nbf-lampiran__jumlah" data-ada={k.jumlah > 0 ? 'ya' : undefined}>
                    {k.jumlah}
                  </span>
                </td>
                <td className="nbf-lampiran__tengah">
                  <button
                    type="button"
                    className="btn btn--sm nbf-lampiran__ikon nbf-lampiran__ikon--unggah"
                    aria-label={`${L.unggah} - ${k.kategori}`}
                    title={L.unggah}
                    disabled={!v.bolehUnggah || mengunggah !== ''}
                    onClick={() => bukaUnggah(k.kategori)}
                  >
                    {mengunggah === k.kategori ? T.mengunggah : <Ikon arah="naik" />}
                  </button>
                </td>
                <td className="nbf-lampiran__tengah">
                  <button
                    type="button"
                    className="btn btn--sm nbf-lampiran__ikon nbf-lampiran__ikon--lihat"
                    aria-label={`${L.lihat} - ${k.kategori}`}
                    title={L.lihat}
                    data-navigasi=""
                    onClick={() => {
                      setLihat(k.kategori)
                      void muatBerkas(k.kategori)
                    }}
                  >
                    <Ikon arah="turun" />
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {/* Popup penampil MENGGANTI popup View File selama terbuka (key berbeda - pola PanelLampiranReas, bug Bordereaux 08-10-2026). */}
      {lihat !== null && penampil !== null && (
        <PenampilLampiran
          key="penampil"
          caseId={caseId}
          berkas={penampil}
          onTutup={(err) => {
            setPenampil(null)
            if (err !== undefined) setGalatLihat(err)
          }}
        />
      )}

      {lihat !== null && penampil === null && (
        <Modal
          key="daftar"
          judul={`${L.lihat} - ${lihat}`}
          onTutup={() => {
            setLihat(null)
            setGalatLihat(null)
          }}
          labelBatal={T.tutup}
          lebar
        >
          <Gagal galat={galatLihat} />
          {berkasKategori === null ? (
            <p className="muted">{T.memuat}</p>
          ) : (
            <div className="table-wrap">
              <table className="nbf-tabel">
                <thead>
                  <tr>
                    <th scope="col">{L.berkas}</th>
                    <th scope="col">{L.catatan}</th>
                    <th scope="col">{L.tanggal}</th>
                    <th scope="col" aria-label={L.view} />
                    {v.bolehHapus && <th scope="col" />}
                  </tr>
                </thead>
                <tbody>
                  {berkasKategori.length === 0 && (
                    <tr>
                      <td colSpan={v.bolehHapus ? 5 : 4}>{T.kosong}</td>
                    </tr>
                  )}
                  {berkasKategori.map((b) => (
                    <tr key={b.id}>
                      <td>
                        <button
                          type="button"
                          className="btn nbf-tema__tombol btn--ghost btn--sm"
                          data-navigasi=""
                          onClick={() => void unduhBerkasBeridentitas(jalurIsiLampiran(caseId, b.id), b.namaFile).catch(setGalat)}
                        >
                          {b.namaFile}
                        </button>
                      </td>
                      <td>
                        {b.catatan}
                        {b.status === 'menunggu' && <span className="muted"> · {T.menunggu}</span>}
                        {b.status === 'gagal' && (
                          <>
                            <span className="field__error"> · {T.gagal} </span>
                            <button type="button" className="btn nbf-tema__tombol btn--ghost btn--sm" onClick={() => void ulangi(b)}>
                              {T.ulangi}
                            </button>
                          </>
                        )}
                      </td>
                      <td>{tanggalRiwayat(b.tanggal)}</td>
                      <td>
                        {/* Hanya berkas yang sudah tersimpan di penyimpanan; jenis lain tetap lewat unduh (klik nama). */}
                        {b.status === '' && jenisPenampil(b.namaFile) !== null && (
                          <button
                            type="button"
                            className="btn nbf-tema__tombol btn--ghost btn--sm"
                            onClick={() => {
                              setGalatLihat(null)
                              setPenampil(b)
                            }}
                          >
                            {jenisPenampil(b.namaFile) === 'office' ? L.viewOffice : L.view}
                          </button>
                        )}
                      </td>
                      {v.bolehHapus && (
                        <td className="table__actions">
                          <button type="button" className="btn nbf-tema__tombol btn--sm" onClick={() => setHapus(b)}>
                            {L.hapus}
                          </button>
                        </td>
                      )}
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </Modal>
      )}

      {hapus !== null && (
        <Modal
          judul={L.hapus}
          onTutup={() => setHapus(null)}
          labelBatal={T.batal}
          aksi={
            <button type="button" className="btn nbf-tema__tombol nbf-tema__tombol--utama btn--primary" onClick={() => void lakukanHapus()}>
              {T.ya}
            </button>
          }
        >
          <p>{T.yakinHapus(hapus.namaFile)}</p>
        </Modal>
      )}
    </section>
  )
}
