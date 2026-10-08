// ⛔ DIPINDAHKAN dari `pages/FormKontrakTreatyIn.tsx` 5 Oktober 2026 —
// pemindahan MURNI, nol perubahan perilaku.

import { useEffect, useRef, useState } from 'react'

import { Gagal, Kosong, Modal, Panel } from '../../../../inti/frontend/components/ui/dasar'
import { pesanGalat } from '../../../../inti/frontend/klien'
import {
  ambilPanelLampiran,
  ambilTautanLampiran,
  hapusLampiran,
  ubahKategoriLampiran,
  unduhLampiran,
  unggahLampiran,
  type BarisKategoriLampiran,
  type BarisLampiranWarisan,
  type PanelLampiranAPI,
} from '../api'
import {
  KOLOM_LAMPIRAN,
  KOLOM_LIHAT_BERKAS,
  LAMPIRAN,
} from '../labels'
import { gabungBerkas, unggahBerurutan } from '../unggahBerkas'

/**
 * Panel Attachment — Category · Count · Upload file · View File.
 *
 * ⛔ Disalin dari `Section/WorkAttachments.xml`: keempat kolomnya, kedua
 * tombolnya, dan spanduk birunya. Datanya dari `M_ATTACHMENTTREATY_2`,
 * tabel WARISAN — nol tabel baru dibuat untuk panel ini.
 *
 * ⛔ SELURUH kategori tampil, termasuk yang nol berkas: kolom `Count` tidak
 * akan pernah berbunyi `0` kalau barisnya disembunyikan saat kosong.
 *
 * ⭐ `Upload file` HIDUP sejak 8 Oktober 2026 — keputusan pemakai: unggahan
 * masuk `M_ATTACHMENTTREATY_2`. Rantainya `Section/WorkAttachments.xml`:
 * tombol per kategori (syarat `TreatyIn.ViewState !='1' ||
 * TreatyIn.RevisionState='1'`) → `SetkategoriDoc` → modal `ASM Attach
 * Content` (`TreatyAttachContent`) → Attach (`TreatySaveAttachment`) →
 * Refresh (`GetMasterTreatyCategory_Act`). Kode kategorinya dari
 * `M_KATEGORIMASTERTREATY` — kesebelasnya kini dipastikan.
 *
 * ⭐ BENTUK modal unggah disamakan dengan menu Master Product Name Life
 * (permintaan pemakai 8 Oktober 2026, tim memakai bentuk itu): kotak
 * seret-lepas, berkas terpilih bisa dibuang satu-satu, diunggah SATU PER
 * SATU dengan progres `2/5 nama`, dan berkas yang gagal tinggal di pilihan
 * beserta alasannya untuk diulang. Modal tertutup sendiri bila semua
 * berhasil.
 */
export default function PanelLampiran({
  kategori: kategoriAwal,
  berkas: berkasAwal,
  idKontrak = '',
  bisaUnggah = false,
  statusAkseptasi = '',
}: {
  kategori: readonly BarisKategoriLampiran[]
  berkas: readonly BarisLampiranWarisan[]
  /** `TreatyIn.ID` — kosong untuk kontrak yang belum tersimpan. */
  idKontrak?: string
  /** `TreatyIn.ViewState !='1' || TreatyIn.RevisionState='1'`. */
  bisaUnggah?: boolean
  /** `TreatyIn.StatusAkseptasi` — syarat wadah Change Category. */
  statusAkseptasi?: string
}) {
  // Isi panel — dari kontrak yang dimuat, lalu dari Refresh/Attach sendiri
  // (panel saja; isian form yang belum di-Save tidak ikut dibaca ulang).
  const [kategori, setKategori] = useState<readonly BarisKategoriLampiran[]>(kategoriAwal)
  const [berkas, setBerkas] = useState<readonly BarisLampiranWarisan[]>(berkasAwal)
  useEffect(() => {
    setKategori(kategoriAwal)
    setBerkas(berkasAwal)
  }, [kategoriAwal, berkasAwal])
  // Modal `ASM Attach Content` — kategori yang dipilih (`SetkategoriDoc`).
  const [unggahKe, setUnggahKe] = useState<BarisKategoriLampiran | null>(null)
  // Bentuk Product Name Life: berkas terpilih, berkas yang sedang diunggah
  // ("2/5 nama"), kegagalan per berkas, seret di atas kotak.
  const [dipilih, setDipilih] = useState<File[]>([])
  const [proses, setProses] = useState<string | null>(null)
  const [gagalUnggah, setGagalUnggah] = useState<{ berkas: File; galat: unknown }[]>([])
  const [seret, setSeret] = useState(false)
  const [galatUnggah, setGalatUnggah] = useState<unknown>(null)
  const [sibuk, setSibuk] = useState(false)
  // Galat Refresh / unduh / Delete / Change Category.
  const [galatPanel, setGalatPanel] = useState('')
  const galat = (e: unknown) => {
    setGalatPanel(e instanceof Error ? e.message : String(e))
  }
  const terapkanPanel = (p: PanelLampiranAPI) => {
    setKategori(p.kategoriLampiran)
    setBerkas(p.lampiran)
  }
  const segarkan = () => {
    if (idKontrak === '') return
    setSibuk(true)
    setGalatPanel('')
    ambilPanelLampiran(idKontrak)
      .then(terapkanPanel)
      .catch(galat)
      .finally(() => {
        setSibuk(false)
      })
  }
  // ⭐ Modal `View File` — `ShowAttachmentTreaty`.
  // Change Category: `StatusDoc.CARI30` 1 = sel kategori menjadi dropdown.
  const [gantiKategori, setGantiKategori] = useState(false)
  const [kategoriBaru, setKategoriBaru] = useState<Record<string, string>>({})
  const bolehGantiKategori = statusAkseptasi !== 'Resolve Complete' && statusAkseptasi !== 'Decline'
  /**
   * Tautan nama berkas (`DownloadAttachmentTreaty`) — isi berkas dialirkan
   * backend dan diunduh lewat `fetch` beridentitas.
   * ⛔ BUKAN membuka URL di tab baru: penjaga lintas modul
   * `claimlife/frontend/unduhdokumen.test.ts` melarang `window.open` dan
   * navigasi lewat skrip (ralat 8 Oktober 2026).
   */
  const unduh = (b: BarisLampiranWarisan) => {
    setSibuk(true)
    setGalatPanel('')
    unduhLampiran(idKontrak, b.id, b.namaBerkas)
      .catch(galat)
      .finally(() => {
        setSibuk(false)
      })
  }
  /**
   * `View Office Online` — penampil kantor di POPUP berbingkai, pola Master
   * Product Name Life: popup dibuka lebih dulu (tanpa jendela baru, jadi
   * tidak diblokir peramban), lalu bingkainya diisi URL penampil dari
   * backend. Jawaban untuk popup yang sudah ditutup dibuang (`giliran`).
   */
  const [penampil, setPenampil] = useState<{ nama: string; siap: boolean } | null>(null)
  const bingkai = useRef<HTMLIFrameElement>(null)
  const giliran = useRef(0)
  const tutupPenampil = () => {
    giliran.current++
    setPenampil(null)
  }
  const bukaOffice = (b: BarisLampiranWarisan) => {
    const ke = ++giliran.current
    setPenampil({ nama: b.namaBerkas, siap: false })
    ambilTautanLampiran(idKontrak, b.id, true)
      .then((j) => {
        if (ke !== giliran.current) return
        if (bingkai.current !== null) bingkai.current.src = j.url
        setPenampil((p) => (p === null ? p : { ...p, siap: true }))
      })
      .catch((e: unknown) => {
        if (ke === giliran.current) tutupPenampil()
        galat(e)
      })
  }
  const hapus = (idLampiran: string) => {
    setSibuk(true)
    setGalatPanel('')
    hapusLampiran(idKontrak, idLampiran)
      .then(terapkanPanel)
      .catch(galat)
      .finally(() => {
        setSibuk(false)
      })
  }
  const simpanKategori = () => {
    const perubahan = Object.entries(kategoriBaru).map(([id, kode]) => ({ id, kategori: kode }))
    setSibuk(true)
    setGalatPanel('')
    ubahKategoriLampiran(idKontrak, perubahan)
      .then((p) => {
        terapkanPanel(p)
        setGantiKategori(false)
        setKategoriBaru({})
      })
      .catch(galat)
      .finally(() => {
        setSibuk(false)
      })
  }
  const bukaUnggah = (k: BarisKategoriLampiran) => {
    setUnggahKe(k)
    setDipilih([])
    setGagalUnggah([])
    setGalatUnggah(null)
    setProses(null)
  }
  /**
   * Attach (`TreatySaveAttachment`) — bentuk Product Name Life: SATU berkas
   * per permintaan, berurutan. Tiap jawaban membawa panel yang dibaca ulang,
   * jadi `Count` naik selagi berkas berikutnya dikirim.
   */
  const lampirkan = () => {
    if (unggahKe === null || sibuk) return
    const kode = unggahKe.kode
    const antre = dipilih
    setSibuk(true)
    setGalatUnggah(null)
    setGagalUnggah([])
    void (async () => {
      try {
        if (antre.length === 0) {
          // Kalimat backend apa adanya: "Tidak ada berkas yang dipilih."
          await unggahLampiran(idKontrak, kode, [])
          return
        }
        const gagal = await unggahBerurutan(
          antre,
          async (f) => {
            const h = await unggahLampiran(idKontrak, kode, [f])
            terapkanPanel(h)
            // Berkas DITOLAK (jenis tak dikenal, terlalu besar …) dijawab 200
            // dengan `berhasil: false` — ia tetap kegagalan berkas ini.
            const b = h.berkas.at(0)
            if (b !== undefined && !b.berhasil) throw new Error(b.pesan)
          },
          (f, i) => {
            setProses(`${String(i + 1)}/${String(antre.length)} ${f.name}`)
          },
        )
        // Yang berhasil keluar dari pilihan (tampil di Count); yang gagal
        // tinggal untuk diperbaiki / diulang.
        setGagalUnggah(gagal)
        setDipilih(gagal.map((g) => g.berkas))
        if (gagal.length === 0) setUnggahKe(null)
      } catch (e) {
        setGalatUnggah(e)
      } finally {
        setProses(null)
        setSibuk(false)
      }
    })()
  }
  /**
   * ⭐ PERMINTAAN PERUBAHAN, satu-satunya tempat ronde ini menyimpang dari
   * Pega — dan pemilik proses yang memintanya, pada keterangan gambar `25`:
   *
   *   "dan view upload saran dibuatkan pop up dan ini diubah menjadi
   *    design nya bagus"
   *
   * ⛔ Yang gambar 25 perlihatkan adalah JENDELA CHROME TERPISAH
   * (`ShowAttachmentTreaty`), lengkap dengan bilah alamat `appdev...`. Di
   * sini ia menjadi modal di dalam halaman.
   *
   * ⚠️ `Modal` DIPANGGIL dari `inti`, tidak ditulis ulang — modal kedua di
   * repositori ini berarti dua perilaku tutup, dua perangkap fokus, dan dua
   * tempat untuk salah.
   */
  const [berkasDilihat, setBerkasDilihat] = useState<string | null>(null)
  const kategoriDilihat = kategori.find((k) => k.kode === berkasDilihat) ?? null
  return (
    <Panel judul={LAMPIRAN.judul}>
      {/* Spanduk biru — ATURAN nama berkas, bukan hiasan. Ditegakkan di
          services (`NamaBerkasAman`) dan dinyatakan di sini. */}
      <span className="trin__spanduk" role="note">
        {LAMPIRAN.spanduk}
      </span>
      {/* `Refresh` — `GetMasterTreatyCategory_Act` (`!pyIsMobile`).
          `Download All` ber-`pyCondition never` di ekspor: tidak dirender. */}
      <div className="trin__aksi">
        <button type="button" className="btn btn--sm" disabled={sibuk || idKontrak === ''} onClick={segarkan}>
          {LAMPIRAN.segarkan}
        </button>
      </div>
      {galatPanel !== '' && berkasDilihat === null && (
        <div className="alert alert--error" role="alert">
          {galatPanel}
        </div>
      )}
      {bisaUnggah && idKontrak === '' && <p className="trin__redup">{LAMPIRAN.simpanDulu}</p>}

      <div className="table-wrap">
        <table className="trin__tabel">
          <thead>
            <tr>
              {KOLOM_LAMPIRAN.map((k) => (
                <th key={k} scope="col">
                  {k}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {kategori.length === 0 && (
              <tr>
                <td colSpan={KOLOM_LAMPIRAN.length}>
                  <Kosong pesan={LAMPIRAN.tanpaLampiran} petunjuk={LAMPIRAN.petunjukLampiran} />
                </td>
              </tr>
            )}
            {kategori.map((k) => (
              <tr key={k.kode}>
                {/* ⭐ NAMA apa adanya — kesebelasnya, persis tangkapan layar
                    pemilik proses 6 Oktober 2026. Sebelumnya empat baris
                    berbunyi KODE-nya (`00003 nama kategori belum
                    dipastikan`), sebab panel dirender dari katalog basis
                    data alih-alih dari daftar layar.

                    ⭐ 8 Oktober 2026 pasangan kode↔nama TERJAWAB oleh
                    `M_KATEGORIMASTERTREATY` (katalog RD Pega) —
                    `docs/PERTANYAAN-TERBUKA-KODE-KATEGORI-LAMPIRAN.md`. */}
                <td>{k.nama}</td>
                {/* ⛔ Cacah TIDAK diformat — ia butir, bukan uang. */}
                <td>{String(k.cacah)}</td>
                <td>
                  {/* ⭐ `Upload file` — tampil bila `TreatyIn.ViewState !='1'
                      || TreatyIn.RevisionState='1'`; butuh kode kategori dan
                      kontrak yang sudah ber-ID. */}
                  {bisaUnggah && (
                    <button
                      type="button"
                      className="btn btn--sm"
                      disabled={sibuk || idKontrak === '' || k.kode === ''}
                      onClick={() => {
                        bukaUnggah(k)
                      }}
                    >
                      {LAMPIRAN.unggah}
                    </button>
                  )}
                </td>
                <td>
                  {/* ⭐ `View File` HIDUP — ia hanya menampilkan apa yang
                      sudah dibaca, dan menampilkan itu yang panel baca-saja
                      memang boleh lakukan. */}
                  <button
                    type="button"
                    className="btn btn--sm"
                    onClick={() => {
                      setGalatPanel('')
                      setGantiKategori(false)
                      setKategoriBaru({})
                      setBerkasDilihat(k.kode)
                    }}
                  >
                    {LAMPIRAN.lihatBerkas}
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {unggahKe !== null && (
        <Modal
          judul={LAMPIRAN.judulUnggah}
          onTutup={() => {
            setUnggahKe(null)
          }}
          labelBatal={LAMPIRAN.batal}
          aksi={
            <button type="button" className="btn btn--primary" disabled={sibuk} onClick={lampirkan}>
              {LAMPIRAN.lampirkan}
            </button>
          }
        >
          {/* ⭐ Bentuk modal `Add attachment` menu Master Product Name Life. */}
          <p>
            <strong>{unggahKe.nama}</strong>
          </p>
          {galatUnggah !== null && <Gagal galat={galatUnggah} />}
          <label
            className={'trin__unggah' + (seret ? ' trin__unggah--seret' : '')}
            onDragOver={(e) => {
              e.preventDefault()
              setSeret(true)
            }}
            onDragLeave={() => {
              setSeret(false)
            }}
            onDrop={(e) => {
              e.preventDefault()
              setSeret(false)
              const jatuh = Array.from(e.dataTransfer.files)
              setDipilih((b) => gabungBerkas(b, jatuh))
            }}
          >
            <span>{LAMPIRAN.seretBerkas}</span>
            <input
              type="file"
              multiple
              aria-label={LAMPIRAN.pilihBerkas}
              disabled={sibuk}
              onChange={(e) => {
                const baru = Array.from(e.target.files ?? [])
                setDipilih((b) => gabungBerkas(b, baru))
                // Dikosongkan: memilih berkas yang sama lagi tetap memicu onChange.
                e.target.value = ''
              }}
            />
          </label>
          {dipilih.length > 0 && (
            <ul className="trin__unggah-daftar">
              {dipilih.map((f) => (
                <li key={f.name}>
                  <span>{f.name}</span>
                  <button
                    type="button"
                    className="btn btn--ghost btn--sm"
                    disabled={sibuk}
                    onClick={() => {
                      setDipilih((b) => b.filter((x) => x !== f))
                    }}
                  >
                    {LAMPIRAN.buangPilihan}
                  </button>
                </li>
              ))}
            </ul>
          )}
          {proses !== null && (
            <p className="muted">
              {LAMPIRAN.mengunggah} {proses}
            </p>
          )}
          {gagalUnggah.length > 0 && (
            <ul className="trin__unggah-gagal">
              {gagalUnggah.map((g) => (
                <li key={g.berkas.name}>
                  {g.berkas.name}: {pesanGalat(g.galat) ?? (g.galat instanceof Error ? g.galat.message : LAMPIRAN.gagal)}
                </li>
              ))}
            </ul>
          )}
        </Modal>
      )}

      {berkasDilihat !== null && (
        <Modal
          judul={LAMPIRAN.judulLihatBerkas}
          onTutup={() => {
            setBerkasDilihat(null)
          }}
          labelBatal={LAMPIRAN.tutup}
          lebar
        >
          {/* `Change Category` / `Save` — wadah bersyarat `StatusAkseptasi`
              bukan Resolve Complete / Decline; tombolnya bergantian menurut
              `StatusDoc.CARI30`. */}
          {bolehGantiKategori && idKontrak !== '' && (
            <div className="trin__aksi">
              {gantiKategori ? (
                <button type="button" className="btn btn--primary btn--sm" disabled={sibuk} onClick={simpanKategori}>
                  {LAMPIRAN.simpanKategori}
                </button>
              ) : (
                <button
                  type="button"
                  className="btn btn--sm"
                  disabled={sibuk}
                  onClick={() => {
                    setGantiKategori(true)
                  }}
                >
                  {LAMPIRAN.gantiKategori}
                </button>
              )}
            </div>
          )}
          {galatPanel !== '' && (
            <div className="alert alert--error" role="alert">
              {galatPanel}
            </div>
          )}
          <div className="table-wrap">
            <table className="trin__tabel trin__lihat-berkas">
              {/* ⭐ Dirapikan 8 Oktober 2026 (permintaan pemakai): nama
                  berkas lebar, sel Type memuat jenis + tombolnya sebaris. */}
              <colgroup>
                <col className="trin__lihat-nama" />
                <col className="trin__lihat-tipe" />
              </colgroup>
              <thead>
                <tr>
                  {/* ⛔ DUA kolom, persis gambar 25 — `File Name` dan
                      `Type`. Tidak tiga, tidak empat. */}
                  {KOLOM_LIHAT_BERKAS.map((k) => (
                    <th key={k} scope="col">
                      {k}
                    </th>
                  ))}
                </tr>
              </thead>
              <tbody>
                {berkasKategori(berkas, kategoriDilihat).length === 0 && (
                  <tr>
                    <td colSpan={KOLOM_LIHAT_BERKAS.length}>{LAMPIRAN.tanpaBaris}</td>
                  </tr>
                )}
                {berkasKategori(berkas, kategoriDilihat).map((b) => (
                  <tr key={b.id}>
                    {/* Nama berkas = tautan `DownloadAttachmentTreaty`. */}
                    <td>
                      <button
                        type="button"
                        className="trin__tautan"
                        disabled={idKontrak === '' || sibuk}
                        onClick={() => {
                          unduh(b)
                        }}
                      >
                        {b.namaBerkas}
                      </button>
                    </td>
                    <td>
                      <div className="trin__lihat-jenis">
                      <span>{b.jenisMime}</span>
                      <div className="trin__lihat-aksi">
                      {/* `View Office Online` — xls/xlsx/doc/docx/ppt/pptx. */}
                      {JENIS_OFFICE.includes(b.jenisMime.toLowerCase()) && (
                        <button
                          type="button"
                          className="btn btn--ghost btn--sm"
                          disabled={idKontrak === ''}
                          onClick={() => {
                            bukaOffice(b)
                          }}
                        >
                          {LAMPIRAN.viewOffice}
                        </button>
                      )}
                      {/* Mode Change Category: kategori menjadi dropdown. */}
                      {gantiKategori && (
                        <select
                          className="field__input"
                          aria-label={LAMPIRAN.gantiKategori}
                          value={kategoriBaru[b.id] ?? b.kodeKategori}
                          onChange={(e) => {
                            setKategoriBaru((x) => ({ ...x, [b.id]: e.target.value }))
                          }}
                        >
                          {kategori
                            .filter((k) => k.kode !== '')
                            .map((k) => (
                              <option key={k.kode} value={k.kode}>
                                {k.nama}
                              </option>
                            ))}
                        </select>
                      )}
                      {/* `Delete` — `ViewState !='1' || RevisionState='1'`. */}
                      {bisaUnggah && idKontrak !== '' && (
                        <button
                          type="button"
                          className="btn btn--sm tl-hapus"
                          disabled={sibuk}
                          onClick={() => {
                            hapus(b.id)
                          }}
                        >
                          {LAMPIRAN.hapus}
                        </button>
                      )}
                      </div>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </Modal>
      )}

      {/* `View Office Online` — penampil kantor di bingkai popup. */}
      {penampil !== null && (
        <Modal judul={penampil.nama} onTutup={tutupPenampil} labelBatal={LAMPIRAN.tutup} penuh>
          {!penampil.siap && <p className="trin__redup">{LAMPIRAN.memuatPenampil}</p>}
          <iframe ref={bingkai} title={penampil.nama} className="trin__penampil" />
        </Modal>
      )}
    </Panel>
  )
}

/** Syarat `View Office Online` — `ShowAttachmentTreaty`. */
const JENIS_OFFICE: readonly string[] = ['xls', 'xlsx', 'doc', 'docx', 'ppt', 'pptx']

/**
 * Berkas milik SATU kategori.
 *
 * ⚠️ Disaring menurut `kode`, bukan menurut nama: nama kategori yang
 * pasangannya belum dipastikan ditampilkan sebagai kodenya, dan menyaring
 * menurut yang tampil akan menyaring menurut dua hal yang berbeda.
 */
function berkasKategori(
  berkas: readonly BarisLampiranWarisan[],
  kategori: BarisKategoriLampiran | null,
): readonly BarisLampiranWarisan[] {
  if (kategori === null) return []
  // ⛔ Kode bila ia diketahui, NAMA bila tidak. Keempat kategori yang
  // kodenya belum dipastikan akan selalu berbunyi `No items` kalau
  // penyaringnya hanya kode — dan kosong yang salah terbaca persis seperti
  // kosong yang benar.
  if (kategori.kode !== '') return berkas.filter((b) => b.kodeKategori === kategori.kode)
  return berkas.filter((b) => b.namaKategori === kategori.nama)
}
