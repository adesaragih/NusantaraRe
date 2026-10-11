// Layar kasus Endorsement Life - FlowAction `InputEDMLife` (b166) → `Section/InputEDMLife.xml`.
//
// Kepala seluruhnya baca-saja (sel `ro`); grid peserta b11899 (EdmType 1) / b17500 (EdmType 3).
// `Save` b37202 → `SetPremi_EDM` (tiket 05/06): centang `.EdmBatal` b15753 dan `DELETE ALL` b13607
// hanya di grid Perubahan Data; rekap mata uang per `.Type` (b23064 …) sesudah simpan.
// Unggah CSV (tiket 07) - wadah b8698 `.EdmType=1 && .EditInput=1`: `Upload CSV` b8973 dan
// `Add CSV Data` b10405 mati bila `.EditInput1=1` (b8965/b10403, ada baris `New`); `View Upload` b9340.
// Keputusan (tiket 08) - `ConfirmSection` tampil sesudah `Save` (b35518 `.IsJsonPolis=1`); Confirm
// menutup dengan `ConfirmSubmitEDM`, Decline kembali ke kotak masuk (`finishAssignment` b38385).

import { useCallback, useEffect, useRef, useState, type ReactNode } from 'react'

import { Gagal, Halaman, Kosong, Memuat } from '../../../../inti/frontend/components/ui/dasar'
import {
  ambilKasus,
  ambilPeserta,
  simpanKasus,
  tambahCSV,
  type HalamanEDM,
  type HasilPutusanEDM,
  type KasusEDM,
  type PeriksaCSVEDM,
  type PesertaEDM,
  unduhTemplatCSV,
} from '../api'
import HasilCSV from '../components/HasilCSV'
import KeputusanEDM from '../components/KeputusanEDM'
import KonfirmasiProses from '../components/KonfirmasiProses'
import PerbandinganRekap from '../components/PerbandinganRekap'
import PolisLama, { tipePopup } from '../components/PolisLama'
import RincianPeserta from '../components/RincianPeserta'
import { TabelKorpus } from '../components/TabelKorpus'
import TerimaKasih from '../components/TerimaKasih'
import UnggahCSV, { TabelPesanCSV } from '../components/UnggahCSV'
import { barisUnggahan } from '../csv'
import { MATA_UANG } from '../kolomKorpus'
import { KASUS_EDM, PERBANDINGAN_EDM, POLIS_LAMA_EDM, PROSES_EDM, SIMPAN_EDM, UMUM_EDM, UNGGAH_EDM } from '../labels'
import {
  adaPerubahan,
  alihBaris,
  alihSemua,
  bagianKepalaKasus,
  bolehCentang,
  PILIHAN_KOSONG,
  sel,
  selAngka,
  selTanggal,
  tercentang,
  type PilihanHapusEDM,
  UKURAN_HALAMAN_EDM,
  kalimatKonfirmasiProses,
  kolomGrid,
} from '../tampilan'
import '../endorsementlife.css'

// Kolom grid peserta dipindah ke `tampilan.ts` (10-10-2026) supaya popup View Upload (`HasilCSV`) memakai sumber yang
// sama tanpa impor komponen → halaman; diekspor ulang di sini agar impor lama tetap berlaku.
export { KOLOM_GRID_DI_RINCIAN, kolomGrid } from '../tampilan'

function isiSel(p: PesertaEDM, kolom: string, jenis: 't' | 'd' | 'n'): string {
  const v = p.nilai[kolom] ?? ''
  if (jenis === 'd') return selTanggal(v)
  if (jenis === 'n') return selAngka(v)
  return sel(v)
}

/**
 * Satu baris grid peserta; `Details` membuka rincian di bawahnya (`pyEditingMode` `expandPane`,
 * `pyEditAction` `PL_DetailAction`).
 */
function Baris({
  p,
  kolom,
  centang,
  terbuka,
  onBuka,
  children,
}: {
  p: PesertaEDM
  kolom: ReturnType<typeof kolomGrid>
  /** `null` = grid Batal b17500, tanpa kolom kotak centang. */
  centang: { boleh: boolean; aktif: boolean; onAlih: () => void } | null
  terbuka: boolean
  onBuka: () => void
  children: ReactNode
}) {
  return (
    <>
      <tr>
        {centang !== null && (
          <td>
            <input
              type="checkbox"
              aria-label={UMUM_EDM.tandaiHapus}
              checked={centang.boleh ? centang.aktif : p.edmStatus === 'Delete'}
              disabled={!centang.boleh}
              onChange={centang.onAlih}
            />
          </td>
        )}
        {kolom.map(([label, kol, jenis]) => (
          <td key={label}>{isiSel(p, kol, jenis)}</td>
        ))}
        <td>
          <button type="button" className="btn btn--ghost btn--sm" aria-expanded={terbuka} onClick={onBuka}>
            {UMUM_EDM.rinci}
          </button>
        </td>
      </tr>
      {terbuka && (
        <tr className="edm-baris-rinci">
          <td colSpan={kolom.length + (centang === null ? 1 : 2)}>{children}</td>
        </tr>
      )}
    </>
  )
}

/**
 * Judul layar + baris chip di bawahnya (permintaan work owner 08-10-2026, pola visual kepala Premium List Detail -
 * kelas sendiri `edm-kepala__*`, tanpa impor premiumlistlife): chip utama Case ID kasus (EDML-…). Chip
 * Endorsement No BELUM ada - `NO_ENDORS` tidak termasuk respons BacaKasus (`kolomKasus`), backend tidak diubah.
 */
export function JudulKasus({ id }: { id: string }) {
  return (
    <div className="edm-kepala__judul">
      <h2 className="inbox__judul">{KASUS_EDM.judul}</h2>
      {id !== '' && (
        <div className="edm-kepala__meta">
          <span className="edm-kepala__chip edm-kepala__chip--utama">{id}</span>
        </div>
      )}
    </div>
  )
}

export default function InputEDMLife({ kasusId, onTutup }: { kasusId: string; onTutup: () => void }) {
  const [kasus, setKasus] = useState<KasusEDM | null>(null)
  const [peserta, setPeserta] = useState<HalamanEDM<PesertaEDM> | null>(null)
  const [halaman, setHalaman] = useState(1)
  const [galat, setGalat] = useState<unknown>(null)
  const [buka, setBuka] = useState('')
  const [polisLama, setPolisLama] = useState(false)
  const [pilihan, setPilihan] = useState<PilihanHapusEDM>(PILIHAN_KOSONG)
  const [menyimpan, setMenyimpan] = useState(false)
  const [konfirmasi, setKonfirmasi] = useState(false)
  const [galatTemplat, setGalatTemplat] = useState<unknown>(null)
  // Summary Comparison dimuat ulang sesudah Process Endorsement / Add CSV Data (08-10-2026).
  const [versiBanding, setVersiBanding] = useState(0)
  const [galatSimpan, setGalatSimpan] = useState<unknown>(null)
  const [jendelaCSV, setJendelaCSV] = useState<'unggah' | 'hasil' | null>(null)
  const [berkasCSV, setBerkasCSV] = useState<File | null>(null)
  const [barisCSV, setBarisCSV] = useState<Array<Record<string, string>>>([])
  const [periksa, setPeriksa] = useState<PeriksaCSVEDM | null>(null)
  const [menambah, setMenambah] = useState(false)
  const [disimpanCSV, setDisimpanCSV] = useState<number | null>(null)
  const [putusan, setPutusan] = useState<HasilPutusanEDM | null>(null)
  // Berkas CSV yang berlaku sekarang - hasil baca berkas lama tidak boleh menimpa berkas yang lebih baru.
  const berkasKini = useRef<File | null>(null)

  const muatKasus = useCallback(async () => {
    setGalat(null)
    try {
      setKasus(await ambilKasus(kasusId))
    } catch (e) {
      setGalat(e)
    }
  }, [kasusId])

  const muatPeserta = useCallback(
    async (h: number) => {
      try {
        setPeserta(await ambilPeserta(kasusId, h))
      } catch (e) {
        setGalat(e)
      }
    },
    [kasusId],
  )

  useEffect(() => {
    void muatKasus()
  }, [muatKasus])
  useEffect(() => {
    void muatPeserta(halaman)
  }, [muatPeserta, halaman])

  const simpan = async () => {
    setMenyimpan(true)
    setGalatSimpan(null)
    try {
      await simpanKasus(kasusId, pilihan)
      setPilihan(PILIHAN_KOSONG)
      await muatKasus()
      await muatPeserta(halaman)
      setVersiBanding((v) => v + 1)
    } catch (e) {
      setGalatSimpan(e)
    } finally {
      setMenyimpan(false)
    }
  }

  const pilihCSV = (berkas: File, hasil: PeriksaCSVEDM) => {
    berkasKini.current = berkas
    setBerkasCSV(berkas)
    setPeriksa(hasil)
    setDisimpanCSV(null)
    setBarisCSV([])
  }

  // `View Upload`: berkas diurai di peramban baru saat popup dibuka (unggahan tanpa batas baris).
  const bukaHasilCSV = () => {
    setJendelaCSV('hasil')
    const b = berkasKini.current
    if (b === null) return
    void b.text().then((t) => {
      if (berkasKini.current === b) setBarisCSV(barisUnggahan(t))
    })
  }

  const tambah = async () => {
    if (berkasCSV === null) return
    setMenambah(true)
    setGalatSimpan(null)
    try {
      const j = await tambahCSV(kasusId, berkasCSV)
      if (j.jenis === 'ditolak') {
        setPeriksa(j.periksa)
        return
      }
      setDisimpanCSV(j.hasil.disimpan)
      setPeriksa(null)
      berkasKini.current = null
      setBerkasCSV(null)
      setBarisCSV([])
      await muatKasus()
      await muatPeserta(halaman)
      setVersiBanding((v) => v + 1)
    } catch (e) {
      setGalatSimpan(e)
    } finally {
      setMenambah(false)
    }
  }

  if (galat !== null) return <Gagal galat={galat} />
  if (kasus === null) return <Memuat pesan={UMUM_EDM.memuat} />

  const k = kasus.kepala
  // Kepala dikelompokkan tiga bagian berjudul (keputusan work owner 07-10-2026) - `bagianKepalaKasus`.
  const bagian = bagianKepalaKasus(kasus)
  const kolom = kolomGrid()
  const denganCentang = kasus.edmType === '1'
  const terbukaKasus = kasus.status === ''
  // `DELETE ALL` mengikuti kotak centang: hidup selama kasus terbuka dan belum disimpan.
  const bolehHapus = bolehCentang({ edmStatus: 'Old' }, kasus)
  const kolomUang = MATA_UANG[tipePopup(k.TYPE ?? '')]
  // Proteksi tanpa perubahan (08-10-2026): Perubahan Data tanpa centang dan tanpa peserta New tidak dapat diproses.
  const tanpaPerubahan = !adaPerubahan(kasus, pilihan)

  return (
    <section className="panel edm-kasus">
      <header className="inbox__kepala edm-kasus__kepala">
        <JudulKasus id={kasus.id} />
        <div className="edm-aksi">
          {/* `View Old Policy` - dipindah dari `ShowLifePremiumSummary_EDM` (R02); popup sesuai `.Type`. */}
          <button type="button" className="btn btn--ghost btn--sm" onClick={() => setPolisLama(true)}>
            {POLIS_LAMA_EDM.viewOldPolicy}
          </button>
          {/* Save dipindah ke bawah (Process Endorsement); Back dihapus - kembali ke kotak masuk lewat menu
              Endorsement Life (keputusan work owner 08-10-2026). */}
        </div>
      </header>
      {/* Kotak info tepat di bawah judul (07-10-2026). */}
      <p className="alert alert--info edm-kasus__info" role="note">
        {UMUM_EDM.terkunci}
      </p>
      {galatSimpan !== null && <Gagal galat={galatSimpan} />}
      <div className="edm-kasus__tigabagian">
        {bagian.map((b) => (
          <section key={b.judul} className="edm-kasus__bagian">
            <h3 className="edm-subjudul">{b.judul}</h3>
            <dl className="edm-kasus__grid">
              {b.medan.map(([label, nilai, lebar]) => (
                <div key={label} className={lebar ? 'edm-kepala__medan edm-kasus__lebar' : 'edm-kepala__medan'}>
                  <dt>{label}</dt>
                  <dd>{nilai}</dd>
                </div>
              ))}
            </dl>
          </section>
        ))}
      </div>
      {/* Bagian peserta berjudul; tombol CSV menjadi toolbar di baris judul. Sesudah Process Endorsement
          (`kasus.sudahSimpan`) ketiganya NONAKTIF + teks bantuan - langkah final (keputusan work owner 08-10-2026). */}
      <div className="edm-kasus__kepalapeserta">
        <h3 className="edm-subjudul">{KASUS_EDM.peserta}</h3>
        {denganCentang && terbukaKasus && (
        <div className="edm-aksi">
          {kasus.sudahSimpan && (
            <span className="edm-catatan" id="edm-csv-bantuan" aria-live="polite">
              {PROSES_EDM.sudahDiproses}
            </span>
          )}
          {/* Download Template (keputusan work owner 10-10-2026): kiri Upload CSV, mati pada kondisi yang sama. */}
          <button
            type="button"
            className="btn btn--ghost btn--sm"
            disabled={kasus.csvTerkunci || kasus.sudahSimpan}
            aria-describedby={kasus.sudahSimpan ? 'edm-csv-bantuan' : undefined}
            onClick={() => {
              setGalatTemplat(null)
              unduhTemplatCSV(kasusId).catch((e: unknown) => setGalatTemplat(e))
            }}
          >
            {UNGGAH_EDM.downloadTemplate}
          </button>
          <button
            type="button"
            className="btn btn--ghost btn--sm"
            disabled={kasus.csvTerkunci || kasus.sudahSimpan}
            aria-describedby={kasus.sudahSimpan ? 'edm-csv-bantuan' : undefined}
            onClick={() => setJendelaCSV('unggah')}
          >
            {UNGGAH_EDM.uploadCsv}
          </button>
          <button
            type="button"
            className="btn btn--ghost btn--sm"
            disabled={kasus.sudahSimpan}
            aria-describedby={kasus.sudahSimpan ? 'edm-csv-bantuan' : undefined}
            onClick={bukaHasilCSV}
          >
            {UNGGAH_EDM.viewUpload}
          </button>
          <button
            type="button"
            className="btn btn--sm"
            disabled={kasus.csvTerkunci || kasus.sudahSimpan || berkasCSV === null || menambah}
            aria-describedby={kasus.sudahSimpan ? 'edm-csv-bantuan' : undefined}
            onClick={() => void tambah()}
          >
            {menambah ? UMUM_EDM.mengunggah : UNGGAH_EDM.addCsvData}
          </button>
          {disimpanCSV !== null && (
            <span className="edm-catatan">
              {UMUM_EDM.barisDisimpan}: {disimpanCSV}
            </span>
          )}
        </div>
        )}
      </div>
      {galatTemplat !== null && <Gagal galat={galatTemplat} />}
      {periksa !== null && jendelaCSV === null && <TabelPesanCSV hasil={periksa} />}
      {peserta !== null && peserta.baris.length === 0 && <Kosong pesan={UMUM_EDM.kosong} />}
      {peserta !== null && peserta.baris.length > 0 && (
        <div className="edm-gulir">
          <table className="inbox__tabel">
            <thead>
              <tr>
                {denganCentang && (
                  <th>
                    <button type="button" className="btn btn--ghost btn--sm" disabled={!bolehHapus} onClick={() => setPilihan(alihSemua(pilihan))}>
                      {SIMPAN_EDM.deleteAll}
                    </button>
                  </th>
                )}
                {kolom.map(([label]) => (
                  <th key={label}>{label}</th>
                ))}
                <th aria-label={UMUM_EDM.rinci} />
              </tr>
            </thead>
            <tbody>
              {peserta.baris.map((p) => (
                <Baris
                  key={p.id}
                  p={p}
                  kolom={kolom}
                  centang={
                    denganCentang
                      ? { boleh: bolehCentang(p, kasus), aktif: tercentang(pilihan, p.id), onAlih: () => setPilihan(alihBaris(pilihan, p.id)) }
                      : null
                  }
                  terbuka={buka === p.id}
                  onBuka={() => setBuka(buka === p.id ? '' : p.id)}
                >
                  <RincianPeserta kasusId={kasusId} pesertaId={p.id} tipe={k.TYPE ?? ''} />
                </Baris>
              ))}
            </tbody>
          </table>
        </div>
      )}
      {peserta !== null && peserta.total > 0 && (
        <Halaman halaman={halaman} ukuran={UKURAN_HALAMAN_EDM} total={peserta.total} onPindah={setHalaman} />
      )}
      {/* `Save` b37202 → "Process Endorsement" (keputusan work owner 08-10-2026): di bawah tabel peserta + paging,
          di atas rekap mata uang dan Keputusan. Kondisi TETAP: `VIS .EditInput=1` (kasus terbuka), mati bila
          `.IsJsonPolis=1` b37200; klik membuka konfirmasi, Continue menjalankan `simpan` yang sama. */}
      {terbukaKasus && (
        <div className="edm-kasus__proses">
          {/* Teks bantuan di samping tombol nonaktif - pola Submit Create Addendum (`edm-buat__bantuan`). */}
          {!kasus.sudahSimpan && tanpaPerubahan && (
            <p className="edm-catatan edm-kasus__bantuan" id="edm-proses-bantuan" aria-live="polite">
              {PROSES_EDM.tanpaPerubahan}
            </p>
          )}
          <button
            type="button"
            className="btn btn--primary"
            disabled={kasus.sudahSimpan || menyimpan || tanpaPerubahan}
            aria-describedby={!kasus.sudahSimpan && tanpaPerubahan ? 'edm-proses-bantuan' : undefined}
            onClick={() => setKonfirmasi(true)}
          >
            {menyimpan ? PROSES_EDM.memproses : PROSES_EDM.processEndorsement}
          </button>
        </div>
      )}
      {/* Summary Comparison (08-10-2026, bukan korpus): tepat di atas rekap mata uang, sebelum Keputusan. Hanya
          SESUDAH Process Endorsement (keputusan work owner 10-10-2026) - sebelumnya tidak dirender dan API-nya
          tidak dipanggil; sudahSimpan naik lewat muatKasus sesudah simpan, jadi tabel langsung muncul. */}
      {kasus.sudahSimpan && <PerbandinganRekap kasusId={kasusId} tipe={k.TYPE ?? ''} versi={versiBanding} />}
      {/* Rekap T_PREMIUM_LIST_SUMMARY berjudul (tampilan 08-10-2026; keterangan dihapus 08-10-2026). */}
      {kasus.rekap.length > 0 && (
        <section className="edm-rekap">
          <h3 className="edm-subjudul">{PERBANDINGAN_EDM.judulRekap}</h3>
          <TabelKorpus kolom={kolomUang} baris={kasus.rekap} kunci={(b, i) => `${b.CURRENCY ?? ''}-${i}`} tandaiNegatif />
        </section>
      )}
      {terbukaKasus && kasus.sudahSimpan && (
        <section className="edm-keputusan">
          <h3 className="edm-subjudul">{PERBANDINGAN_EDM.judulKeputusan}</h3>
          <KeputusanEDM
            kasusId={kasusId}
            riwayat={kasus.riwayat}
            // Accept: popup hasil (nomor kosong pun tetap tampil "—"); Decline: langsung ke kotak masuk (09-10-2026).
            onSelesai={(h, status) => (status === '1' ? setPutusan(h) : onTutup())}
          />
        </section>
      )}
      {putusan !== null && <TerimaKasih hasil={putusan} kasus={kasus} onTutup={onTutup} />}
      {konfirmasi && (
        <KonfirmasiProses
          kalimat={kalimatKonfirmasiProses(kasus, pilihan)}
          onTutup={() => setKonfirmasi(false)}
          onLanjut={() => {
            setKonfirmasi(false)
            void simpan()
          }}
        />
      )}
      {jendelaCSV === 'unggah' && <UnggahCSV kasusId={kasusId} onBerkas={pilihCSV} onTutup={() => setJendelaCSV(null)} />}
      {jendelaCSV === 'hasil' && <HasilCSV baris={barisCSV} onTutup={() => setJendelaCSV(null)} />}
      {polisLama && <PolisLama kasusId={kasusId} tipe={k.TYPE ?? ''} onTutup={() => setPolisLama(false)} />}
    </section>
  )
}
